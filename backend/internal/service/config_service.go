package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/xiaoj/frpc_webmanager/internal/model"
)

// ConfigService 配置服务
type ConfigService struct {
	configDir   string
	configMutex sync.Mutex
}

// NewConfigService 创建配置服务
func NewConfigService(dataDir string) (*ConfigService, error) {
	// 创建配置目录
	configDir := filepath.Join(dataDir, "configs")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}

	return &ConfigService{
		configDir: configDir,
	}, nil
}

// GenerateTaskConfig 为任务生成配置文件
func (s *ConfigService) GenerateTaskConfig(task *model.Task) (string, error) {
	return "", fmt.Errorf("功能待重构")
}

// GenerateMergedConfig 为服务器组生成合并配置文件
func (s *ConfigService) GenerateMergedConfig(serverAddr string, serverPort int, authToken string, tasks []*model.Task) (string, error) {
	// 构建 TOML 配置内容
	var configBuilder strings.Builder

	// 写入服务器通用配置
	configBuilder.WriteString(fmt.Sprintf("serverAddr = \"%s\"\n", serverAddr))
	configBuilder.WriteString(fmt.Sprintf("serverPort = %d\n", serverPort))
	if authToken != "" {
		configBuilder.WriteString(fmt.Sprintf("auth.token = \"%s\"\n", authToken))
	}
	configBuilder.WriteString("\n")

	// 合并所有任务的代理配置
	for _, task := range tasks {
		if task.Status != model.TaskStatusRunning {
			continue // 只包含运行中的任务
		}

		for _, proxy := range task.Proxies {
			// 使用任务ID_代理名称确保唯一性
			proxyName := fmt.Sprintf("%s_%s", task.ID, proxy.Name)
			if proxy.Name == "" {
				proxyName = fmt.Sprintf("%s_proxy_%s", task.ID, proxy.ID)
			}

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

	return configBuilder.String(), nil
}

// GetTaskConfigPath 获取任务配置文件路径
func (s *ConfigService) GetTaskConfigPath(taskID string) string {
	return s.getTaskConfigPath(taskID)
}

// GetServerGroupConfigPath 获取服务器组配置文件路径
func (s *ConfigService) GetServerGroupConfigPath(serverAddr string, serverPort int) string {
	return s.getServerGroupConfigPath(serverAddr, serverPort)
}

// getTaskConfigPath 获取任务配置文件路径(内部方法)
func (s *ConfigService) getTaskConfigPath(taskID string) string {
	return filepath.Join(s.configDir, fmt.Sprintf("%s.toml", taskID))
}

// getServerGroupConfigPath 获取服务器组配置文件路径(内部方法)
func (s *ConfigService) getServerGroupConfigPath(serverAddr string, serverPort int) string {
	// 将服务器地址中的点替换为下划线,避免文件名问题
	safeAddr := strings.ReplaceAll(serverAddr, ".", "_")
	return filepath.Join(s.configDir, fmt.Sprintf("merged_%s_%d.toml", safeAddr, serverPort))
}

// ValidateConfig 验证配置文件
func (s *ConfigService) ValidateConfig(configPath string) error {
	// 检查文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("配置文件不存在: %s", configPath)
	}

	return nil
}

// SaveMergedConfig 保存服务器组的合并配置文件
func (s *ConfigService) SaveMergedConfig(serverAddr string, serverPort int, authToken string, tasks []*model.Task) error {
	s.configMutex.Lock()
	defer s.configMutex.Unlock()

	// 生成配置内容
	configContent, err := s.GenerateMergedConfig(serverAddr, serverPort, authToken, tasks)
	if err != nil {
		return fmt.Errorf("生成配置失败: %w", err)
	}

	// 获取配置文件路径
	configPath := s.getServerGroupConfigPath(serverAddr, serverPort)

	// 写入文件
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

// SaveMergedConfigAll 保存服务器组的合并配置文件(包含所有状态的任务)
func (s *ConfigService) SaveMergedConfigAll(serverAddr string, serverPort int, authToken string, tasks []*model.Task) error {
	s.configMutex.Lock()
	defer s.configMutex.Unlock()

	// 构建 TOML 配置内容
	var configBuilder strings.Builder

	// 写入服务器通用配置
	configBuilder.WriteString(fmt.Sprintf("serverAddr = \"%s\"\n", serverAddr))
	configBuilder.WriteString(fmt.Sprintf("serverPort = %d\n", serverPort))
	if authToken != "" {
		configBuilder.WriteString(fmt.Sprintf("auth.token = \"%s\"\n", authToken))
	}
	configBuilder.WriteString("\n")

	// 合并所有任务的代理配置(包含所有状态的任务)
	for _, task := range tasks {
		for _, proxy := range task.Proxies {
			// 使用任务ID_代理名称确保唯一性
			proxyName := fmt.Sprintf("%s_%s", task.ID, proxy.Name)
			if proxy.Name == "" {
				proxyName = fmt.Sprintf("%s_proxy_%s", task.ID, proxy.ID)
			}

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

	// 获取配置文件路径
	configPath := s.getServerGroupConfigPath(serverAddr, serverPort)

	// 写入文件
	if err := os.WriteFile(configPath, []byte(configBuilder.String()), 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

// DeleteServerConfig 删除服务器配置文件
func (s *ConfigService) DeleteServerConfig(serverAddr string, serverPort int) error {
	s.configMutex.Lock()
	defer s.configMutex.Unlock()

	configPath := s.getServerGroupConfigPath(serverAddr, serverPort)

	// 检查文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil // 文件不存在,不算错误
	}

	// 删除文件
	if err := os.Remove(configPath); err != nil {
		return fmt.Errorf("删除配置文件失败: %w", err)
	}

	return nil
}
