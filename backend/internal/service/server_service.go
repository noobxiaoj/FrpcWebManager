package service

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/xiaoj/frpc_webmanager/internal/frpc"
	"github.com/xiaoj/frpc_webmanager/internal/model"
	"github.com/xiaoj/frpc_webmanager/internal/repository"
	"github.com/xiaoj/frpc_webmanager/internal/utils"
)

// ServerService 服务器服务
type ServerService struct {
	repo          *repository.ServerRepository
	taskRepo      *repository.TaskRepository
	configService *ConfigService
	frpcManager   *frpc.FrpcManager
}

// NewServerService 创建服务器服务
func NewServerService(
	repo *repository.ServerRepository,
	taskRepo *repository.TaskRepository,
	configService *ConfigService,
	frpcManager *frpc.FrpcManager,
) *ServerService {
	return &ServerService{
		repo:          repo,
		taskRepo:      taskRepo,
		configService: configService,
		frpcManager:   frpcManager,
	}
}

// ListServers 获取服务器列表
func (s *ServerService) ListServers() ([]model.Server, error) {
	servers, err := s.repo.ListServers()
	if err != nil {
		return nil, err
	}

	// 更新每个服务器的任务数量和状态（不保存到数据库，只返回最新状态）
	for i := range servers {
		if err := s.updateServerTaskCountWithoutSave(&servers[i]); err != nil {
			// 记录错误但继续处理其他服务器
			fmt.Printf("警告: 更新服务器任务数量失败: %v\n", err)
		}
	}

	return servers, nil
}

// GetServer 获取服务器详情
func (s *ServerService) GetServer(id string) (*model.Server, error) {
	server, err := s.repo.GetServer(id)
	if err != nil {
		return nil, err
	}

	// 更新任务数量
	if err := s.updateServerTaskCount(server); err != nil {
		// 记录错误但继续返回服务器信息
		fmt.Printf("警告: 更新服务器任务数量失败: %v\n", err)
	}

	return server, nil
}

// CreateServer 创建服务器
func (s *ServerService) CreateServer(req *model.CreateServerRequest) (*model.Server, error) {
	// 验证输入
	if req.Name == "" {
		return nil, errors.New("服务器名称不能为空")
	}
	if req.Address == "" {
		return nil, errors.New("服务器地址不能为空")
	}
	if req.Port == "" {
		return nil, errors.New("端口不能为空")
	}

	// 创建服务器对象
	now := time.Now()
	server := &model.Server{
		ID:           utils.GenerateUUID(),
		Name:         req.Name,
		Address:      fmt.Sprintf("%s:%s", req.Address, req.Port),
		Status:       model.ServerStatusNoTask, // 新服务器默认为无任务状态
		Uptime:       "",
		LogMaxHeight: "none",
		Logs:         []model.LogEntry{}, // 空日志,后续只记录连接成功的日志
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// 保存到数据库
	if err := s.repo.CreateServer(server); err != nil {
		return nil, err
	}

	return server, nil
}

// UpdateServer 更新服务器
func (s *ServerService) UpdateServer(id string, req *model.UpdateServerRequest) (*model.Server, error) {
	// 获取现有服务器
	server, err := s.repo.GetServer(id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Name != nil {
		server.Name = *req.Name
	}
	if req.Address != nil && req.Port != nil {
		server.Address = fmt.Sprintf("%s:%s", *req.Address, *req.Port)
	}
	server.UpdatedAt = time.Now()

	// 保存更新
	if err := s.repo.UpdateServer(server); err != nil {
		return nil, err
	}

	return server, nil
}

// DeleteServer 删除服务器
func (s *ServerService) DeleteServer(id string, forceDelete bool) ([]string, error) {
	// 获取服务器信息
	server, err := s.repo.GetServer(id)
	if err != nil {
		return nil, err
	}

	// 从服务器地址中解析地址和端口
	addr, port, err := parseServerAddress(server.Address)
	if err != nil {
		return nil, err
	}

	// 获取服务器的任务数量
	tasks, err := s.taskRepo.GetByServer(addr, port)
	if err != nil {
		return nil, fmt.Errorf("获取服务器任务失败: %w", err)
	}

	// 定义任务名称列表
	taskNames := make([]string, 0)

	// 如果服务器有任务,返回任务列表
	if len(tasks) > 0 {
		for _, task := range tasks {
			taskNames = append(taskNames, task.Name)
		}

		// 如果不是强制删除，返回任务列表让前端决定
		if !forceDelete {
			return taskNames, fmt.Errorf("服务器还有 %d 个任务", len(tasks))
		}

		// 如果是强制删除，先删除所有任务
		for _, task := range tasks {
			if err := s.taskRepo.Delete(task.ID); err != nil {
				return nil, fmt.Errorf("删除任务 %s 失败: %w", task.Name, err)
			}
		}
	}

	// 删除服务器
	if err := s.repo.DeleteServer(id); err != nil {
		return nil, err
	}

	return taskNames, nil
}

// AddLog 添加日志到服务器
func (s *ServerService) AddLog(serverID string, log model.LogEntry) error {
	server, err := s.repo.GetServer(serverID)
	if err != nil {
		return err
	}

	server.Logs = append(server.Logs, log)

	// 限制日志条数,最多保留100条
	if len(server.Logs) > 100 {
		server.Logs = server.Logs[len(server.Logs)-100:]
	}

	server.UpdatedAt = time.Now()

	return s.repo.UpdateServer(server)
}

// ClearLogs 清空服务器日志
func (s *ServerService) ClearLogs(serverID string) error {
	server, err := s.repo.GetServer(serverID)
	if err != nil {
		return err
	}

	// 从服务器地址中解析地址和端口
	addr, port, err := parseServerAddress(server.Address)
	if err != nil {
		return err
	}

	// 获取服务器key
	serverKey := fmt.Sprintf("%s:%d", addr, port)

	// 清空日志收集器中的内存缓冲和日志文件
	s.frpcManager.GetLogCollector().ClearLogs(serverKey)

	// 清空数据库中的日志
	server.Logs = []model.LogEntry{}
	server.UpdatedAt = time.Now()

	return s.repo.UpdateServer(server)
}

// GenerateAndSaveConfig 为指定服务器生成并保存 frpc 配置文件
// 根据服务器地址和端口找到所有关联的任务,整合生成一个 toml 配置文件
func (s *ServerService) GenerateAndSaveConfig(serverID string) error {
	// 获取服务器信息
	server, err := s.repo.GetServer(serverID)
	if err != nil {
		return fmt.Errorf("获取服务器失败: %w", err)
	}

	// 从服务器地址中解析地址和端口
	addr, port, err := parseServerAddress(server.Address)
	if err != nil {
		return err
	}

	// 获取所有连接到该服务器的任务
	tasks, err := s.taskRepo.GetByServer(addr, port)
	if err != nil {
		return fmt.Errorf("获取服务器任务失败: %w", err)
	}

	// 如果没有任务,则不生成配置
	if len(tasks) == 0 {
		return fmt.Errorf("服务器没有关联的任务")
	}

	// 使用第一个任务的 authToken(假设同一服务器的所有任务使用相同的 token)
	authToken := ""
	if len(tasks) > 0 && tasks[0].AuthToken != "" {
		authToken = tasks[0].AuthToken
	}

	// 保存配置文件(包含所有状态的任务)
	if err := s.configService.SaveMergedConfigAll(addr, port, authToken, tasks); err != nil {
		return fmt.Errorf("保存配置文件失败: %w", err)
	}

	return nil
}

// GenerateAndSaveConfigByAddr 根据服务器地址和端口生成并保存 frpc 配置文件
func (s *ServerService) GenerateAndSaveConfigByAddr(serverAddr string, serverPort int) error {
	// 获取所有连接到该服务器的任务
	tasks, err := s.taskRepo.GetByServer(serverAddr, serverPort)
	if err != nil {
		return fmt.Errorf("获取服务器任务失败: %w", err)
	}

	// 如果没有任务,则不生成配置
	if len(tasks) == 0 {
		return nil // 没有任务不算错误
	}

	// 使用第一个任务的 authToken(假设同一服务器的所有任务使用相同的 token)
	authToken := ""
	if len(tasks) > 0 && tasks[0].AuthToken != "" {
		authToken = tasks[0].AuthToken
	}

	// 保存配置文件(包含所有状态的任务)
	if err := s.configService.SaveMergedConfigAll(serverAddr, serverPort, authToken, tasks); err != nil {
		return fmt.Errorf("保存配置文件失败: %w", err)
	}

	return nil
}

// GetConfigPath 获取服务器配置文件的路径
func (s *ServerService) GetConfigPath(serverID string) (string, error) {
	// 获取服务器信息
	server, err := s.repo.GetServer(serverID)
	if err != nil {
		return "", fmt.Errorf("获取服务器失败: %w", err)
	}

	// 从服务器地址中解析地址和端口
	addr, port, err := parseServerAddress(server.Address)
	if err != nil {
		return "", err
	}

	// 返回配置文件路径
	return s.configService.GetServerGroupConfigPath(addr, port), nil
}

// updateServerTaskCount 更新服务器的任务数量和状态
func (s *ServerService) updateServerTaskCount(server *model.Server) error {
	if err := s.updateServerTaskCountWithoutSave(server); err != nil {
		return err
	}

	// 保存更新
	return s.repo.UpdateServer(server)
}

// updateServerTaskCountWithoutSave 更新服务器的任务数量和状态（不保存到数据库）
func (s *ServerService) updateServerTaskCountWithoutSave(server *model.Server) error {
	// 从服务器地址中解析地址和端口
	addr, port, err := parseServerAddress(server.Address)
	if err != nil {
		return err
	}

	// 获取连接到该服务器的所有任务
	tasks, err := s.taskRepo.GetByServer(addr, port)
	if err != nil {
		return fmt.Errorf("获取服务器任务失败: %w", err)
	}

	// 更新任务数量
	taskCount := len(tasks)
	server.TaskCount = taskCount

	// 统计运行中的任务数量。
	// 这里使用任务自身的状态字段，而不是日志或端口推断，
	// 这样可以保证“服务器状态”完全由“任务状态 + frpc 进程状态”共同决定。
	runningTaskCount := 0
	for _, task := range tasks {
		if task.Status == model.TaskStatusRunning {
			runningTaskCount++
		}
	}

	// 预先缓存 frpc 进程运行状态，避免同一次状态计算中重复访问进程管理器。
	// 后续规则严格按照产品定义执行：
	// 1. taskCount == 0 -> 无任务
	// 2. taskCount > 0 且至少一个任务在运行，且 frpc 在运行 -> 在线
	// 3. taskCount > 0 且没有任务在运行，且 frpc 未运行 -> 离线
	// 4. taskCount > 0 且“任务运行状态”与“frpc 运行状态”不一致 -> 故障
	// 5. 其他未覆盖情况 -> 疑似异常
	isFrpcRunning := s.frpcManager.IsServerRunning(addr, port)
	hasRunningTask := runningTaskCount > 0

	// 根据任务数量和进程状态更新服务器状态
	if taskCount == 0 {
		// 场景 1：该服务器完全没有任务。
		server.Status = model.ServerStatusNoTask
	} else if hasRunningTask && isFrpcRunning {
		// 场景 2：至少一个任务在运行，同时 frpc 进程也在运行。
		server.Status = model.ServerStatusOnline
	} else if !hasRunningTask && !isFrpcRunning {
		// 场景 3：有任务，但没有任何任务在运行，且 frpc 进程也未运行。
		server.Status = model.ServerStatusOffline
	} else if (hasRunningTask && !isFrpcRunning) || (!hasRunningTask && isFrpcRunning) {
		// 场景 4：任务运行状态和 frpc 运行状态“打架”了。
		// 例如：
		// - 有任务显示运行中，但 frpc 实际没起来
		// - 没有任务在运行，但 frpc 进程却还活着
		server.Status = model.ServerStatusFault
	} else {
		// 场景 5：理论上不应进入这里。
		// 为了避免把未知场景误标成正常状态，统一标记为疑似异常，便于排查。
		server.Status = model.ServerStatusSuspectedAbnormal
	}

	// 尝试获取 webServer 端口
	if webServerPort, err := s.frpcManager.GetAllocatedWebServerPort(addr, port); err == nil {
		server.WebServerPort = webServerPort
	} else {
		// 如果没有分配端口，设置为0
		server.WebServerPort = 0
	}

	// 更新运行时长
	uptime := s.frpcManager.GetServerUptime(addr, port)
	if uptime != "" {
		server.Uptime = uptime
	} else if server.Status == model.ServerStatusOffline || server.Status == model.ServerStatusFault {
		// 离线或故障且拿不到运行时长时，统一显示“未连接”，
		// 这样可以明确告诉用户当前没有拿到可用的 frpc 运行时间信息。
		server.Uptime = "未连接"
	}

	return nil
}

// UpdateServerTaskCountByAddr 根据服务器地址和端口更新任务数量
func (s *ServerService) UpdateServerTaskCountByAddr(serverAddr string, serverPort int) error {
	// 获取所有服务器
	servers, err := s.repo.ListServers()
	if err != nil {
		return err
	}

	// 找到匹配的服务器并更新
	for i := range servers {
		addr, port, err := parseServerAddress(servers[i].Address)
		if err != nil {
			continue
		}

		if addr == serverAddr && port == serverPort {
			if err := s.updateServerTaskCount(&servers[i]); err != nil {
				return err
			}
			break
		}
	}

	return nil
}

// DeleteServerConfigByAddr 根据服务器地址和端口删除配置文件
func (s *ServerService) DeleteServerConfigByAddr(serverAddr string, serverPort int) error {
	return s.configService.DeleteServerConfig(serverAddr, serverPort)
}

// parseServerAddress 解析服务器地址,返回地址和端口
func parseServerAddress(address string) (string, int, error) {
	// 尝试使用 net.SplitHostPort 解析
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return "", 0, fmt.Errorf("解析服务器地址失败: %w", err)
	}

	// 转换端口
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return "", 0, fmt.Errorf("解析端口失败: %w", err)
	}

	return host, portNum, nil
}

// GetServerLogs 获取服务器日志
func (s *ServerService) GetServerLogs(serverID string, limit int) ([]model.LogEntry, error) {
	// 获取服务器信息
	server, err := s.repo.GetServer(serverID)
	if err != nil {
		return nil, fmt.Errorf("获取服务器失败: %w", err)
	}

	// 从服务器地址中解析地址和端口
	addr, port, err := parseServerAddress(server.Address)
	if err != nil {
		return nil, err
	}

	// 获取服务器key
	serverKey := fmt.Sprintf("%s:%d", addr, port)

	// 从frpcManager获取日志
	return s.frpcManager.GetLogCollector().GetLogs(serverKey, limit), nil
}

// AddServerLog 添加系统日志到服务器
func (s *ServerService) AddServerLog(serverID string, level string, message string) error {
	server, err := s.repo.GetServer(serverID)
	if err != nil {
		return err
	}

	// 从服务器地址中解析地址和端口
	addr, port, err := parseServerAddress(server.Address)
	if err != nil {
		return err
	}

	// 获取服务器key
	serverKey := fmt.Sprintf("%s:%d", addr, port)

	// 写入日志到收集器
	s.frpcManager.GetLogCollector().WriteLog(serverKey, level, message)

	return nil
}

// UpdateServerUptime 更新服务器运行时长
func (s *ServerService) UpdateServerUptime(server *model.Server) error {
	// 从服务器地址中解析地址和端口
	addr, port, err := parseServerAddress(server.Address)
	if err != nil {
		return err
	}

	// 获取运行时长
	uptime := s.frpcManager.GetServerUptime(addr, port)
	server.Uptime = uptime

	// 保存更新
	return s.repo.UpdateServer(server)
}

// UpdateServerOrder 更新服务器排序
func (s *ServerService) UpdateServerOrder(order []string) error {
	return s.repo.UpdateServerOrder(order)
}

// UpdateServerLock 更新服务器锁定状态
func (s *ServerService) UpdateServerLock(id string, locked bool) error {
	// 获取服务器信息
	server, err := s.repo.GetServer(id)
	if err != nil {
		return err
	}

	// 更新锁定状态
	server.Locked = locked
	server.UpdatedAt = time.Now()

	// 保存更新
	return s.repo.UpdateServer(server)
}
