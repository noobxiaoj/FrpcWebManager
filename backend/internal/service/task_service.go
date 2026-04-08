package service

import (
	"fmt"
	"time"

	"github.com/xiaoj/frpc_webmanager/internal/frpc"
	"github.com/xiaoj/frpc_webmanager/internal/model"
	"github.com/xiaoj/frpc_webmanager/internal/repository"
	"github.com/xiaoj/frpc_webmanager/internal/utils"
)

// TaskService 任务服务
type TaskService struct {
	taskRepo      *repository.TaskRepository
	processRepo   *repository.ProcessRepository
	frpcManager   *frpc.FrpcManager
	serverService *ServerService
}

// NewTaskService 创建任务服务
func NewTaskService(
	taskRepo *repository.TaskRepository,
	processRepo *repository.ProcessRepository,
	frpcManager *frpc.FrpcManager,
	serverService *ServerService,
) *TaskService {
	return &TaskService{
		taskRepo:      taskRepo,
		processRepo:   processRepo,
		frpcManager:   frpcManager,
		serverService: serverService,
	}
}

// CreateTask 创建任务
func (s *TaskService) CreateTask(req *model.CreateTaskRequest) (*model.Task, error) {
	// 为每个代理配置生成ID(如果前端没有提供)
	proxies := make([]model.Proxy, len(req.Proxies))
	for i, proxy := range req.Proxies {
		if proxy.ID == "" {
			proxy.ID = utils.GenerateUUID()
		}
		proxies[i] = proxy
	}

	// 创建任务模型
	task := &model.Task{
		ID:          utils.GenerateUUID(),
		Name:        req.Name,
		Description: req.Description,
		ServerAddr:  req.ServerAddr,
		ServerPort:  req.ServerPort,
		AuthToken:   req.AuthToken,
		Proxies:     proxies,
		Status:      model.TaskStatusStopped,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 保存到数据库
	if err := s.taskRepo.Create(task); err != nil {
		return nil, err
	}

	// 更新服务器的任务数量
	if err := s.serverService.UpdateServerTaskCountByAddr(task.ServerAddr, task.ServerPort); err != nil {
		fmt.Printf("警告: 更新服务器任务数量失败: %v\n", err)
	}

	return task, nil
}

// GetTask 获取任务
func (s *TaskService) GetTask(id string) (*model.Task, error) {
	return s.taskRepo.Get(id)
}

// ListTasks 获取任务列表
func (s *TaskService) ListTasks() ([]*model.Task, error) {
	return s.taskRepo.List()
}

// UpdateTask 更新任务
func (s *TaskService) UpdateTask(id string, req *model.UpdateTaskRequest) (*model.Task, error) {
	// 获取现有任务
	task, err := s.taskRepo.Get(id)
	if err != nil {
		return nil, err
	}

	// 记录原始服务器地址和端口
	oldServerAddr := task.ServerAddr
	oldServerPort := task.ServerPort

	// 更新字段
	if req.Name != nil {
		task.Name = *req.Name
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.ServerAddr != nil {
		task.ServerAddr = *req.ServerAddr
	}
	if req.ServerPort != nil {
		task.ServerPort = *req.ServerPort
	}
	if req.AuthToken != nil {
		task.AuthToken = *req.AuthToken
	}
	if req.Proxies != nil {
		// 为每个代理配置生成ID(如果前端没有提供)
		proxies := make([]model.Proxy, len(*req.Proxies))
		for i, proxy := range *req.Proxies {
			if proxy.ID == "" {
				proxy.ID = utils.GenerateUUID()
			}
			proxies[i] = proxy
		}
		task.Proxies = proxies
	}
	task.UpdatedAt = time.Now()

	// 保存更新
	if err := s.taskRepo.Update(task); err != nil {
		return nil, err
	}

	// 如果服务器地址或端口发生变化,需要更新两个服务器的任务数量
	if task.ServerAddr != oldServerAddr || task.ServerPort != oldServerPort {
		// 更新旧服务器的任务数量
		if err := s.serverService.UpdateServerTaskCountByAddr(oldServerAddr, oldServerPort); err != nil {
			fmt.Printf("警告: 更新旧服务器任务数量失败: %v\n", err)
		}
	}

	// 更新新服务器的任务数量
	if err := s.serverService.UpdateServerTaskCountByAddr(task.ServerAddr, task.ServerPort); err != nil {
		fmt.Printf("警告: 更新服务器任务数量失败: %v\n", err)
	}

	return task, nil
}

// DeleteTask 删除任务
func (s *TaskService) DeleteTask(id string) error {
	// 获取任务信息
	task, err := s.taskRepo.Get(id)
	if err != nil {
		return err
	}

	// 记录服务器地址和端口
	serverAddr := task.ServerAddr
	serverPort := task.ServerPort

	// 如果任务正在运行,先尝试停止
	if task.Status == model.TaskStatusRunning {
		// 尝试停止任务,但不因为停止失败而阻止删除
		if err := s.StopTask(id); err != nil {
			// 记录错误但继续删除
			fmt.Printf("警告: 停止任务失败,继续删除: %v\n", err)
		}
	}

	// 删除进程记录
	s.processRepo.DeleteProcess(id)

	// 删除任务
	if err := s.taskRepo.Delete(id); err != nil {
		return err
	}

	// 获取该服务器的剩余任务数量
	tasks, err := s.taskRepo.GetByServer(serverAddr, serverPort)
	if err != nil {
		fmt.Printf("警告: 获取服务器任务失败: %v\n", err)
		// 继续更新服务器状态
	}

	// 根据剩余任务数量决定如何处理配置文件
	if len(tasks) == 0 {
		// 没有任务了,删除配置文件
		if err := s.serverService.DeleteServerConfigByAddr(serverAddr, serverPort); err != nil {
			fmt.Printf("警告: 删除服务器配置文件失败: %v\n", err)
		}
	}

	// 更新服务器的任务数量
	if err := s.serverService.UpdateServerTaskCountByAddr(serverAddr, serverPort); err != nil {
		fmt.Printf("警告: 更新服务器任务数量失败: %v\n", err)
	}

	return nil
}

// StartTask 启动任务
func (s *TaskService) StartTask(id string) error {
	// 获取任务
	task, err := s.taskRepo.Get(id)
	if err != nil {
		return err
	}

	// 检查任务是否已在运行
	if task.Status == model.TaskStatusRunning {
		return fmt.Errorf("任务已在运行中")
	}

	// 获取连接到同一服务器的所有任务
	serverTasks, err := s.taskRepo.GetByServer(task.ServerAddr, task.ServerPort)
	if err != nil {
		return err
	}

	// 计算运行中的任务数量（不包括当前任务）
	runningCount := 0
	for _, t := range serverTasks {
		if t.ID != id && t.Status == model.TaskStatusRunning {
			runningCount++
		}
	}

	// 先更新任务状态为running
	if err := s.taskRepo.UpdateStatus(id, model.TaskStatusRunning); err != nil {
		return fmt.Errorf("更新任务状态失败: %w", err)
	}

	isPaused, err := s.serverService.IsServerPausedByAddr(task.ServerAddr, task.ServerPort)
	if err == nil && isPaused {
		return nil
	}

	// 情况1：服务器已有其他任务运行 -> 使用热重载
	if runningCount > 0 {
		if err := s.reloadServerWithTasks(task.ServerAddr, task.ServerPort, serverTasks); err != nil {
			// 回滚任务状态
			s.taskRepo.UpdateStatus(id, model.TaskStatusStopped)
			return err
		}
		return nil
	}

	// 情况2：服务器没有其他任务 -> 启动新进程
	if err := s.startServerProcess(task.ServerAddr, task.ServerPort, serverTasks); err != nil {
		// 回滚任务状态
		s.taskRepo.UpdateStatus(id, model.TaskStatusStopped)
		return err
	}

	return nil
}

// reloadServerWithTasks 热重载服务器
func (s *TaskService) reloadServerWithTasks(serverAddr string, serverPort int, tasks []*model.Task) error {
	// 获取第一个任务的authToken（假设同一服务器的任务使用相同token）
	authToken := ""
	if tasks != nil && len(tasks) > 0 {
		authToken = tasks[0].AuthToken
	}

	// 使用frpcManager热重载
	if err := s.frpcManager.ReloadServer(serverAddr, serverPort, authToken, tasks); err != nil {
		return fmt.Errorf("热重载失败: %w", err)
	}

	return nil
}

// startServerProcess 启动服务器进程
func (s *TaskService) startServerProcess(serverAddr string, serverPort int, tasks []*model.Task) error {
	// 获取authToken
	authToken := ""
	if tasks != nil && len(tasks) > 0 {
		authToken = tasks[0].AuthToken
	}

	// 使用frpcManager启动进程
	if err := s.frpcManager.StartServer(serverAddr, serverPort, authToken, tasks); err != nil {
		return fmt.Errorf("启动frpc进程失败: %w", err)
	}

	return nil
}

// StopTask 停止任务
func (s *TaskService) StopTask(id string) error {
	// 获取任务
	task, err := s.taskRepo.Get(id)
	if err != nil {
		return err
	}

	// 检查任务是否在运行
	if task.Status != model.TaskStatusRunning {
		return fmt.Errorf("任务未运行")
	}

	// 获取连接到同一服务器的所有任务
	serverTasks, err := s.taskRepo.GetByServer(task.ServerAddr, task.ServerPort)
	if err != nil {
		return err
	}

	// 计算停止该任务后，还在运行的任务数量
	remainingRunningCount := 0
	for _, t := range serverTasks {
		if t.ID != id && t.Status == model.TaskStatusRunning {
			remainingRunningCount++
		}
	}

	// 先更新任务状态为stopped
	if err := s.taskRepo.UpdateStatus(id, model.TaskStatusStopped); err != nil {
		return fmt.Errorf("更新任务状态失败: %w", err)
	}

	isPaused, err := s.serverService.IsServerPausedByAddr(task.ServerAddr, task.ServerPort)
	if err == nil && isPaused {
		return nil
	}

	// 情况1：停止后还有其他任务运行 -> 热重载移除该任务
	if remainingRunningCount > 0 {
		if err := s.reloadServerWithTasks(task.ServerAddr, task.ServerPort, serverTasks); err != nil {
			// 回滚任务状态
			s.taskRepo.UpdateStatus(id, model.TaskStatusRunning)
			return err
		}
		return nil
	}

	// 情况2：没有其他任务运行 -> 停止整个服务器进程
	if err := s.stopServerProcess(task.ServerAddr, task.ServerPort); err != nil {
		// 回滚任务状态
		s.taskRepo.UpdateStatus(id, model.TaskStatusRunning)
		return err
	}

	return nil
}

// stopServerProcess 停止服务器进程
func (s *TaskService) stopServerProcess(serverAddr string, serverPort int) error {
	if err := s.frpcManager.StopServer(serverAddr, serverPort); err != nil {
		return fmt.Errorf("停止frpc进程失败: %w", err)
	}
	return nil
}

// fromTask 从任务列表中移除指定任务
func fromTask(taskID string, tasks []*model.Task) []*model.Task {
	var result []*model.Task
	for _, t := range tasks {
		if t.ID != taskID {
			result = append(result, t)
		}
	}
	return result
}

// ReloadTask 重载任务
func (s *TaskService) ReloadTask(id string) error {
	// 获取任务
	task, err := s.taskRepo.Get(id)
	if err != nil {
		return err
	}

	// 检查任务是否在运行
	if task.Status != model.TaskStatusRunning {
		return fmt.Errorf("任务未运行")
	}

	isPaused, err := s.serverService.IsServerPausedByAddr(task.ServerAddr, task.ServerPort)
	if err == nil && isPaused {
		return fmt.Errorf("服务器当前处于暂停状态，无法重载")
	}

	// 简单实现：先停止再启动
	if err := s.StopTask(id); err != nil {
		return err
	}

	return s.StartTask(id)
}

// reloadServerGroup 重载服务器组
func (s *TaskService) reloadServerGroup(serverAddr string, serverPort int) error {
	// 获取连接到该服务器的所有任务
	tasks, err := s.taskRepo.GetByServer(serverAddr, serverPort)
	if err != nil {
		return err
	}

	return s.reloadServerWithTasks(serverAddr, serverPort, tasks)
}

// GetTaskStatus 获取任务状态
func (s *TaskService) GetTaskStatus(id string) (*model.TaskStatusResponse, error) {
	// 获取任务
	task, err := s.taskRepo.Get(id)
	if err != nil {
		return nil, err
	}

	resp := &model.TaskStatusResponse{
		TaskID:    task.ID,
		Status:    task.Status,
		ServerKey: fmt.Sprintf("%s:%d", task.ServerAddr, task.ServerPort),
	}

	return resp, nil
}

// RestoreRunningTasks 恢复所有状态为running的任务
// 用于后端启动时自动恢复之前运行中的任务
func (s *TaskService) RestoreRunningTasks() error {
	// 获取所有任务
	tasks, err := s.taskRepo.List()
	if err != nil {
		return fmt.Errorf("获取任务列表失败: %w", err)
	}

	// 筛选出状态为running的任务
	var runningTasks []*model.Task
	for _, task := range tasks {
		if task.Status == model.TaskStatusRunning {
			runningTasks = append(runningTasks, task)
		}
	}

	if len(runningTasks) == 0 {
		fmt.Println("没有需要恢复的运行中任务")
		return nil
	}

	fmt.Printf("发现 %d 个运行中的任务,开始恢复...\n", len(runningTasks))

	// 按服务器分组任务
	serverGroups := make(map[string][]*model.Task)
	for _, task := range runningTasks {
		serverKey := fmt.Sprintf("%s:%d", task.ServerAddr, task.ServerPort)
		serverGroups[serverKey] = append(serverGroups[serverKey], task)
	}

	// 为每个服务器组启动进程
	for serverKey, tasks := range serverGroups {
		if len(tasks) == 0 {
			continue
		}

		serverAddr := tasks[0].ServerAddr
		serverPort := tasks[0].ServerPort
		authToken := tasks[0].AuthToken

		isPaused, pauseErr := s.serverService.IsServerPausedByAddr(serverAddr, serverPort)
		if pauseErr == nil && isPaused {
			fmt.Printf("服务器 %s 处于暂停状态，跳过自动恢复\n", serverKey)
			continue
		}

		fmt.Printf("恢复服务器 %s 的 %d 个任务...\n", serverKey, len(tasks))

		// 启动服务器进程
		if err := s.frpcManager.StartServer(serverAddr, serverPort, authToken, tasks); err != nil {
			fmt.Printf("警告: 恢复服务器 %s 的任务失败: %v\n", serverKey, err)
			// 继续处理其他服务器
			continue
		}

		fmt.Printf("服务器 %s 的任务恢复成功\n", serverKey)
	}

	fmt.Printf("成功恢复 %d 个服务器的运行中任务\n", len(serverGroups))
	return nil
}
