package frpc

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/xiaoj/frpc_webmanager/internal/model"
	"github.com/xiaoj/frpc_webmanager/internal/utils"
)

const (
	// BaseWebServerPort webServer端口的基础值
	BaseWebServerPort = 7400
)

// FrpcBinaryManager 管理frpc二进制文件
type FrpcBinaryManager struct {
	binDir string
}

// GetFrpcBinary 获取当前平台的frpc可执行文件路径
func (m *FrpcBinaryManager) GetFrpcBinary() (string, error) {
	runtimeOS := runtime.GOOS
	runtimeArch := runtime.GOARCH

	// 定义可能的文件名列表（按优先级排序）
	var possibleNames []string

	switch {
	case runtimeOS == "darwin" && runtimeArch == "arm64":
		possibleNames = []string{"frpc", "frpc_darwin_arm64"}
	case runtimeOS == "darwin" && runtimeArch == "amd64":
		possibleNames = []string{"frpc", "frpc_darwin_amd64"}
	case runtimeOS == "linux":
		possibleNames = []string{"frpc"}
	case runtimeOS == "windows":
		possibleNames = []string{"frpc.exe", "frpc_windows_amd64.exe"}
	default:
		return "", fmt.Errorf("unsupported platform: %s_%s", runtimeOS, runtimeArch)
	}

	// 将 binDir 转换为绝对路径
	absBinDir, err := filepath.Abs(m.binDir)
	if err != nil {
		return "", fmt.Errorf("获取bin目录绝对路径失败: %w", err)
	}

	// 按优先级查找文件
	for _, binaryName := range possibleNames {
		binaryPath := filepath.Join(absBinDir, binaryName)
		if info, err := os.Stat(binaryPath); err == nil {
			// 文件存在，检查是否可执行
			if !info.Mode().IsRegular() {
				continue
			}
			// 在 Unix 系统上检查可执行权限
			if runtimeOS != "windows" {
				if info.Mode().Perm()&0111 == 0 {
					continue
				}
			}
			return binaryPath, nil
		}
	}

	return "", fmt.Errorf("frpc binary not found in %s (tried: %v)", absBinDir, possibleNames)
}

// ManagedProcess 被管理的frpc进程
type ManagedProcess struct {
	ServerAddr string
	ServerPort int
	Cmd        *exec.Cmd
	PID        int
	ConfigPath string
	LogPath    string
	StartTime  time.Time
	CancelFunc context.CancelFunc
}

// FrpcManager frpc进程管理器
type FrpcManager struct {
	binaryManager     *FrpcBinaryManager
	processes         map[string]*ManagedProcess // serverKey -> ManagedProcess
	processMutex      sync.RWMutex
	logService        *LogCollector
	configDir         string
	settingsPath      string
	logDir            string
	allocatedWebPorts map[int]string // webServer端口 -> serverKey, 用于跟踪已分配的端口
	allocatedPorts    sync.RWMutex   // 保护 allocatedWebPorts 的读写锁
}

// NewFrpcManager 创建frpc管理器
func NewFrpcManager(configDir string) (*FrpcManager, error) {
	// 将 configDir 转换为绝对路径
	absConfigDir, err := filepath.Abs(configDir)
	if err != nil {
		return nil, fmt.Errorf("获取配置目录绝对路径失败: %w", err)
	}

	// bin 目录在项目根目录下的 bin
	// absConfigDir = "xxx/data/configs"
	// 我们需要 "xxx/bin"
	projectRoot := filepath.Dir(filepath.Dir(absConfigDir))
	binDir := filepath.Join(projectRoot, "bin")
	logDir := filepath.Join(filepath.Dir(absConfigDir), "logs")

	// 创建日志目录
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %w", err)
	}

	manager := &FrpcManager{
		binaryManager:     &FrpcBinaryManager{binDir: binDir},
		processes:         make(map[string]*ManagedProcess),
		configDir:         absConfigDir,
		settingsPath:      filepath.Join(filepath.Dir(absConfigDir), "settings.json"),
		logDir:            logDir,
		allocatedWebPorts: make(map[int]string),
	}

	// 初始化日志收集器
	manager.logService = NewLogCollector(logDir)

	// 启动日志清理协程
	go manager.logService.StartLogCleanup()

	return manager, nil
}

// StartServer 启动frpc服务器组进程
func (m *FrpcManager) StartServer(serverAddr string, serverPort int, authToken string, tasks []*model.Task) error {
	serverKey := m.getServerKey(serverAddr, serverPort)

	m.processMutex.Lock()
	defer m.processMutex.Unlock()

	// 检查是否已存在运行中的进程
	// 需要检查进程是否真的在运行,避免使用已退出的进程记录
	if process, exists := m.processes[serverKey]; exists {
		// 检查进程是否还在运行
		if process.Cmd.Process != nil {
			err := process.Cmd.Process.Signal(syscall.Signal(0))
			if err == nil {
				// 进程确实在运行
				return fmt.Errorf("服务器 %s 已有frpc进程运行中 (PID: %d)", serverKey, process.PID)
			}
			// 进程已退出,清理记录
			m.logService.WriteLog(serverKey, "info", fmt.Sprintf("发现已退出的进程记录 (PID: %d), 清理中...", process.PID))
		}
		// 进程记录存在但进程已退出,清理记录
		delete(m.processes, serverKey)
	}

	// 生成配置文件
	configPath, err := m.generateConfig(serverAddr, serverPort, authToken, tasks)
	if err != nil {
		return fmt.Errorf("生成配置文件失败: %w", err)
	}

	// 获取frpc可执行文件路径
	binaryPath, err := m.binaryManager.GetFrpcBinary()
	if err != nil {
		return fmt.Errorf("获取frpc可执行文件失败: %w", err)
	}

	// 创建日志文件
	logPath := filepath.Join(m.logDir, fmt.Sprintf("frpc_%s.log", serverKey))
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("创建日志文件失败: %w", err)
	}

	// 创建context用于取消
	ctx, cancel := context.WithCancel(context.Background())

	// 构建命令
	cmd := exec.CommandContext(ctx, binaryPath, "-c", configPath)

	// 设置工作目录，确保相对路径能正确解析
	// configDir 是 "xxx/data/configs"，项目根目录是 "xxx"
	projectRoot := filepath.Dir(filepath.Dir(m.configDir))
	cmd.Dir = projectRoot

	// 创建管道用于实时日志收集
	stdoutReader, stdoutWriter := io.Pipe()
	stderrReader, stderrWriter := io.Pipe()

	// 同时输出到文件和管道
	cmd.Stdout = io.MultiWriter(logFile, stdoutWriter)
	cmd.Stderr = io.MultiWriter(logFile, stderrWriter)

	// 启动进程
	if err := cmd.Start(); err != nil {
		cancel()
		logFile.Close()
		stdoutWriter.Close()
		stderrWriter.Close()
		return fmt.Errorf("启动frpc进程失败: %w", err)
	}

	// 记录启动信息
	m.logService.WriteLog(serverKey, "info", fmt.Sprintf("frpc进程已启动 (PID: %d, 工作目录: %s)", cmd.Process.Pid, cmd.Dir))

	// 创建进程对象
	process := &ManagedProcess{
		ServerAddr: serverAddr,
		ServerPort: serverPort,
		Cmd:        cmd,
		PID:        cmd.Process.Pid,
		ConfigPath: configPath,
		LogPath:    logPath,
		StartTime:  time.Now(),
		CancelFunc: cancel,
	}

	// 存储进程
	m.processes[serverKey] = process

	// 启动日志收集协程
	go m.logService.CollectLogs(serverKey, stdoutReader, "stdout")
	go m.logService.CollectLogs(serverKey, stderrReader, "stderr")

	// 启动进程监控协程
	go m.monitorProcess(serverKey, process)

	// 记录启动日志
	m.logService.WriteLog(serverKey, "info", fmt.Sprintf("frpc进程已启动 (PID: %d)", process.PID))

	return nil
}

// ReloadServer 热重载frpc配置
func (m *FrpcManager) ReloadServer(serverAddr string, serverPort int, authToken string, tasks []*model.Task) error {
	serverKey := m.getServerKey(serverAddr, serverPort)

	m.processMutex.RLock()
	process, exists := m.processes[serverKey]
	m.processMutex.RUnlock()

	if !exists {
		return fmt.Errorf("服务器 %s 没有运行中的frpc进程", serverKey)
	}

	// 生成新的配置文件
	configPath, err := m.generateConfig(serverAddr, serverPort, authToken, tasks)
	if err != nil {
		return fmt.Errorf("生成配置文件失败: %w", err)
	}

	// 获取frpc可执行文件路径
	binaryPath, err := m.binaryManager.GetFrpcBinary()
	if err != nil {
		return fmt.Errorf("获取frpc可执行文件失败: %w", err)
	}

	// 执行reload命令
	reloadCmd := exec.Command(binaryPath, "reload", "-c", configPath)
	output, err := reloadCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("热重载失败: %w, output: %s", err, string(output))
	}

	// 更新进程的配置路径
	m.processMutex.Lock()
	process.ConfigPath = configPath
	m.processMutex.Unlock()

	// 记录日志
	m.logService.WriteLog(serverKey, "info", fmt.Sprintf("配置已热重载: %s", string(output)))

	return nil
}

// StopServer 停止frpc服务器组进程
func (m *FrpcManager) StopServer(serverAddr string, serverPort int) error {
	serverKey := m.getServerKey(serverAddr, serverPort)

	m.processMutex.Lock()
	process, exists := m.processes[serverKey]
	if !exists {
		m.processMutex.Unlock()
		// 进程不存在，可能已经自然退出，不视为错误
		m.logService.WriteLog(serverKey, "info", "frpc进程不存在，可能已经退出")
		return nil
	}
	// 不在这里删除，让 monitorProcess 来清理
	m.processMutex.Unlock()

	// 记录停止日志
	m.logService.WriteLog(serverKey, "info", fmt.Sprintf("正在停止frpc进程 (PID: %d)...", process.PID))

	// 发送SIGTERM信号
	if process.Cmd.Process != nil {
		if err := process.Cmd.Process.Signal(syscall.SIGTERM); err != nil {
			// 如果SIGTERM失败，尝试Kill
			m.logService.WriteLog(serverKey, "warn", fmt.Sprintf("SIGTERM失败，尝试强制终止: %v", err))
			if err := process.Cmd.Process.Kill(); err != nil {
				// 进程可能已经退出了
				m.logService.WriteLog(serverKey, "warn", fmt.Sprintf("强制终止失败，进程可能已退出: %v", err))
			}
		}
	}

	// 取消context
	process.CancelFunc()

	// 等待进程结束（最多5秒）
	done := make(chan error, 1)
	go func() {
		done <- process.Cmd.Wait()
	}()

	select {
	case <-time.After(5 * time.Second):
		// 超时，强制Kill
		m.logService.WriteLog(serverKey, "warn", "停止进程超时，强制终止")
		if process.Cmd.Process != nil {
			process.Cmd.Process.Kill()
		}
	case err := <-done:
		if err != nil {
			// 进程已退出，可能是正常退出
			m.logService.WriteLog(serverKey, "info", fmt.Sprintf("frpc进程已停止: %v", err))
		} else {
			m.logService.WriteLog(serverKey, "info", "frpc进程已正常停止")
		}
	}

	return nil
}

// IsServerRunning 检查服务器是否运行
func (m *FrpcManager) IsServerRunning(serverAddr string, serverPort int) bool {
	serverKey := m.getServerKey(serverAddr, serverPort)

	m.processMutex.RLock()
	defer m.processMutex.RUnlock()

	process, exists := m.processes[serverKey]
	if !exists {
		return false
	}

	// 检查进程是否还在运行
	if process.Cmd.Process == nil {
		return false
	}

	// 发送信号0检查进程是否存在
	err := process.Cmd.Process.Signal(syscall.Signal(0))
	return err == nil
}

// GetServerProcess 获取服务器进程信息
func (m *FrpcManager) GetServerProcess(serverAddr string, serverPort int) (*ManagedProcess, error) {
	serverKey := m.getServerKey(serverAddr, serverPort)

	m.processMutex.RLock()
	defer m.processMutex.RUnlock()

	process, exists := m.processes[serverKey]
	if !exists {
		return nil, fmt.Errorf("服务器 %s 没有运行中的frpc进程", serverKey)
	}

	return process, nil
}

// ListProcesses 列出所有运行的进程
func (m *FrpcManager) ListProcesses() []*ManagedProcess {
	m.processMutex.RLock()
	defer m.processMutex.RUnlock()

	processes := make([]*ManagedProcess, 0, len(m.processes))
	for _, process := range m.processes {
		processes = append(processes, process)
	}

	return processes
}

// monitorProcess 监控进程状态
func (m *FrpcManager) monitorProcess(serverKey string, process *ManagedProcess) {
	err := process.Cmd.Wait()

	// 进程已退出，清理资源
	m.processMutex.Lock()
	delete(m.processes, serverKey)
	m.processMutex.Unlock()

	// 释放已分配的webServer端口
	m.releaseWebServerPort(serverKey)

	// 记录退出原因
	if err != nil {
		m.logService.WriteLog(serverKey, "error", fmt.Sprintf("frpc进程异常退出: %v", err))
	} else {
		m.logService.WriteLog(serverKey, "info", "frpc进程正常退出")
	}
}

// generateConfig 生成配置文件（只包含运行中的任务）
func (m *FrpcManager) generateConfig(serverAddr string, serverPort int, authToken string, tasks []*model.Task) (string, error) {
	// 替换服务器地址中的点，避免文件名问题
	safeAddr := strings.ReplaceAll(serverAddr, ".", "_")

	// 使用绝对路径生成配置文件
	configPath := filepath.Join(m.configDir, fmt.Sprintf("merged_%s_%d.toml", safeAddr, serverPort))

	// 构建TOML配置内容
	var configBuilder strings.Builder
	connectionIdentifier := utils.LoadConnectionIdentifier(m.settingsPath)

	// 写入服务器通用配置
	configBuilder.WriteString(fmt.Sprintf("serverAddr = \"%s\"\n", serverAddr))
	configBuilder.WriteString(fmt.Sprintf("serverPort = %d\n", serverPort))
	if authToken != "" {
		configBuilder.WriteString(fmt.Sprintf("auth.token = \"%s\"\n", authToken))
	}

	// 配置webServer端口（用于热重载）
	// 动态分配可用端口，避免冲突
	serverKey := m.getServerKey(serverAddr, serverPort)
	webServerPort, err := m.getOrAllocateWebServerPort(BaseWebServerPort, serverKey)
	if err != nil {
		return "", fmt.Errorf("分配 webServer 端口失败: %w", err)
	}
	configBuilder.WriteString(fmt.Sprintf("webServer.addr = \"127.0.0.1\"\n"))
	configBuilder.WriteString(fmt.Sprintf("webServer.port = %d\n", webServerPort))
	configBuilder.WriteString("\n")

	// 只合并运行中的任务
	for _, task := range tasks {
		if task.Status != model.TaskStatusRunning {
			continue // 关键：跳过非运行状态的任务
		}

		for _, proxy := range task.Proxies {
			// 运行态生成的 frpc 配置必须与导出配置保持同一命名规则，
			// 这样日志里的代理名才能稳定对应到用户在设置中看到的连接标识。
			proxyName := utils.GenerateProxyName(connectionIdentifier, task.Name, proxy.Name)

			configBuilder.WriteString("[[proxies]]\n")
			configBuilder.WriteString(fmt.Sprintf("name = \"%s\"\n", proxyName))
			configBuilder.WriteString(fmt.Sprintf("type = \"%s\"\n", proxy.Type))
			configBuilder.WriteString(fmt.Sprintf("localIP = \"%s\"\n", proxy.LocalIP))
			configBuilder.WriteString(fmt.Sprintf("localPort = %d\n", proxy.LocalPort))

			// 根据代理类型写入特定配置
			switch proxy.Type {
			case model.ProxyTypeTCP, model.ProxyTypeUDP:
				configBuilder.WriteString(fmt.Sprintf("remotePort = %d\n", proxy.RemotePort))
			case model.ProxyTypeHTTP, model.ProxyTypeHTTPS:
				if len(proxy.CustomDomains) > 0 {
					configBuilder.WriteString(fmt.Sprintf("customDomains = [\"%s\"]\n",
						strings.Join(proxy.CustomDomains, "\", \"")))
				}
				if proxy.Subdomain != "" {
					configBuilder.WriteString(fmt.Sprintf("subdomain = \"%s\"\n", proxy.Subdomain))
				}
			}

			// 写入额外配置参数
			for key, value := range proxy.Extra {
				configBuilder.WriteString(fmt.Sprintf("%s = \"%s\"\n", key, value))
			}

			configBuilder.WriteString("\n")
		}
	}

	// 确保配置文件目录存在
	if err := os.MkdirAll(m.configDir, 0755); err != nil {
		return "", fmt.Errorf("创建配置目录失败: %w", err)
	}

	// 写入配置文件
	if err := os.WriteFile(configPath, []byte(configBuilder.String()), 0644); err != nil {
		return "", fmt.Errorf("写入配置文件失败: %w", err)
	}

	return configPath, nil
}

// GetLogCollector 获取日志收集器
func (m *FrpcManager) GetLogCollector() *LogCollector {
	return m.logService
}

// GetAllocatedWebServerPort 获取服务器分配的 webServer 端口
func (m *FrpcManager) GetAllocatedWebServerPort(serverAddr string, serverPort int) (int, error) {
	serverKey := m.getServerKey(serverAddr, serverPort)

	m.allocatedPorts.RLock()
	defer m.allocatedPorts.RUnlock()

	for port, key := range m.allocatedWebPorts {
		if key == serverKey {
			return port, nil
		}
	}

	return 0, fmt.Errorf("服务器 %s 未分配 webServer 端口", serverKey)
}

// getServerKey 获取服务器标识
func (m *FrpcManager) getServerKey(serverAddr string, serverPort int) string {
	return fmt.Sprintf("%s:%d", serverAddr, serverPort)
}

// isPortAvailable 检查端口是否可用
func isPortAvailable(port int) bool {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}
	listener.Close()
	return true
}

// getOrAllocateWebServerPort 获取或分配 webServer 端口
// 如果该服务器已经分配了端口，则复用；否则分配新端口
// 从基础端口开始尝试，如果被占用则随机增加1-20，最多尝试100次
// 检查系统端口占用和已分配端口记录，避免端口冲突
func (m *FrpcManager) getOrAllocateWebServerPort(basePort int, serverKey string) (int, error) {
	// 首先检查是否已经为该服务器分配了端口
	m.allocatedPorts.RLock()
	for port, key := range m.allocatedWebPorts {
		if key == serverKey {
			// 已分配端口，直接复用
			m.allocatedPorts.RUnlock()
			return port, nil
		}
	}
	m.allocatedPorts.RUnlock()

	// 未分配端口，进行分配
	maxAttempts := 100
	currentPort := basePort
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < maxAttempts; i++ {
		// 检查端口是否已被系统占用
		if !isPortAvailable(currentPort) {
			// 端口被系统占用，尝试下一个
			currentPort += rand.Intn(20) + 1
			continue
		}

		// 检查端口是否已分配给其他服务器
		m.allocatedPorts.Lock()
		allocatedKey, exists := m.allocatedWebPorts[currentPort]
		if exists && allocatedKey != serverKey {
			// 端口已分配给其他服务器，尝试下一个
			m.allocatedPorts.Unlock()
			currentPort += rand.Intn(20) + 1
			continue
		}

		// 分配端口
		m.allocatedWebPorts[currentPort] = serverKey
		m.allocatedPorts.Unlock()

		return currentPort, nil
	}

	return 0, fmt.Errorf("无法分配可用的 webServer 端口，已尝试 %d 次", maxAttempts)
}

// releaseWebServerPort 释放已分配的 webServer 端口
func (m *FrpcManager) releaseWebServerPort(serverKey string) {
	m.allocatedPorts.Lock()
	defer m.allocatedPorts.Unlock()

	// 查找并释放该服务器占用的所有端口
	for port, key := range m.allocatedWebPorts {
		if key == serverKey {
			delete(m.allocatedWebPorts, port)
		}
	}
}

// GetServerUptime 获取服务器进程运行时长
// 如果服务器未运行,返回空字符串
func (m *FrpcManager) GetServerUptime(serverAddr string, serverPort int) string {
	serverKey := m.getServerKey(serverAddr, serverPort)

	m.processMutex.RLock()
	defer m.processMutex.RUnlock()

	process, exists := m.processes[serverKey]
	if !exists {
		return ""
	}

	// 检查进程是否还在运行
	if process.Cmd.Process == nil {
		return ""
	}

	// 发送信号0检查进程是否存在
	if err := process.Cmd.Process.Signal(syscall.Signal(0)); err != nil {
		return ""
	}

	// 计算运行时长
	uptime := time.Since(process.StartTime)

	// 格式化时长
	return formatDuration(uptime)
}

// formatDuration 格式化时间长度为可读字符串
func formatDuration(d time.Duration) string {
	// 小于1分钟
	if d < time.Minute {
		return fmt.Sprintf("%d秒", int(d.Seconds()))
	}

	// 小于1小时
	if d < time.Hour {
		minutes := int(d.Minutes())
		seconds := int(d.Seconds()) % 60
		if seconds > 0 {
			return fmt.Sprintf("%d分%d秒", minutes, seconds)
		}
		return fmt.Sprintf("%d分钟", minutes)
	}

	// 小于1天
	if d < 24*time.Hour {
		hours := int(d.Hours())
		minutes := int(d.Minutes()) % 60
		if minutes > 0 {
			return fmt.Sprintf("%d小时%d分", hours, minutes)
		}
		return fmt.Sprintf("%d小时", hours)
	}

	// 大于等于1天
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	if hours > 0 {
		return fmt.Sprintf("%d天%d小时", days, hours)
	}
	return fmt.Sprintf("%d天", days)
}
