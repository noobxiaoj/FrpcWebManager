package repository

import (
	"fmt"
	"sync"

	"github.com/xiaoj/frpc_webmanager/internal/model"
)

// ProcessRepository 进程数据仓储
type ProcessRepository struct {
	store             *JSONStore
	processes         map[string]*model.ProcessInfo // taskID -> ProcessInfo
	serverGroups      map[string]*model.ServerGroup // serverKey (addr:port) -> ServerGroup
	processesMutex    sync.RWMutex
	serverGroupsMutex sync.RWMutex
}

// NewProcessRepository 创建进程仓储
func NewProcessRepository(store *JSONStore) *ProcessRepository {
	repo := &ProcessRepository{
		store:        store,
		processes:    make(map[string]*model.ProcessInfo),
		serverGroups: make(map[string]*model.ServerGroup),
	}

	// 加载已有数据
	repo.loadProcesses()
	repo.loadServerGroups()

	return repo
}

// loadProcesses 加载进程信息
func (r *ProcessRepository) loadProcesses() error {
	var processList []*model.ProcessInfo
	if err := r.store.Load("processes.json", &processList); err != nil {
		return err
	}

	r.processesMutex.Lock()
	defer r.processesMutex.Unlock()

	r.processes = make(map[string]*model.ProcessInfo)
	for _, process := range processList {
		r.processes[process.TaskID] = process
	}

	return nil
}

// saveProcesses 保存进程信息
func (r *ProcessRepository) saveProcesses() error {
	r.processesMutex.Lock()
	defer r.processesMutex.Unlock()

	return r.saveProcessesLocked()
}

// loadServerGroups 加载服务器组信息
func (r *ProcessRepository) loadServerGroups() error {
	var groupList []*model.ServerGroup
	if err := r.store.Load("server_groups.json", &groupList); err != nil {
		return err
	}

	r.serverGroupsMutex.Lock()
	defer r.serverGroupsMutex.Unlock()

	r.serverGroups = make(map[string]*model.ServerGroup)
	for _, group := range groupList {
		key := r.getServerKey(group.ServerAddr, group.ServerPort)
		r.serverGroups[key] = group
	}

	return nil
}

// saveServerGroups 保存服务器组信息
func (r *ProcessRepository) saveServerGroups() error {
	r.serverGroupsMutex.RLock()
	defer r.serverGroupsMutex.RUnlock()

	groupList := make([]*model.ServerGroup, 0, len(r.serverGroups))
	for _, group := range r.serverGroups {
		groupList = append(groupList, group)
	}

	return r.store.Save("server_groups.json", groupList)
}

// getServerKey 获取服务器标识
func (r *ProcessRepository) getServerKey(serverAddr string, serverPort int) string {
	return fmt.Sprintf("%s:%d", serverAddr, serverPort)
}

// SaveProcess 保存进程信息
func (r *ProcessRepository) SaveProcess(process *model.ProcessInfo) error {
	r.processesMutex.Lock()
	defer r.processesMutex.Unlock()

	r.processes[process.TaskID] = process

	return r.saveProcesses()
}

// GetProcess 获取进程信息
func (r *ProcessRepository) GetProcess(taskID string) (*model.ProcessInfo, error) {
	r.processesMutex.RLock()
	defer r.processesMutex.RUnlock()

	process, exists := r.processes[taskID]
	if !exists {
		return nil, fmt.Errorf("进程不存在: %s", taskID)
	}

	return process, nil
}

// DeleteProcess 删除进程信息
func (r *ProcessRepository) DeleteProcess(taskID string) error {
	r.processesMutex.Lock()
	defer r.processesMutex.Unlock()

	delete(r.processes, taskID)

	return r.saveProcessesLocked()
}

// saveProcessesLocked 保存进程信息(调用者必须持有锁)
func (r *ProcessRepository) saveProcessesLocked() error {
	processList := make([]*model.ProcessInfo, 0, len(r.processes))
	for _, process := range r.processes {
		processList = append(processList, process)
	}

	return r.store.Save("processes.json", processList)
}

// StopProcess 停止进程
func (r *ProcessRepository) StopProcess(pid int) error {
	// TODO: 实现进程停止逻辑
	// 这里需要使用 os/process 包来终止进程
	return fmt.Errorf("停止进程功能待实现")
}

// ListProcesses 获取所有进程信息
func (r *ProcessRepository) ListProcesses() ([]*model.ProcessInfo, error) {
	r.processesMutex.RLock()
	defer r.processesMutex.RUnlock()

	processes := make([]*model.ProcessInfo, 0, len(r.processes))
	for _, process := range r.processes {
		processes = append(processes, process)
	}

	return processes, nil
}

// SaveServerGroup 保存服务器组
func (r *ProcessRepository) SaveServerGroup(group *model.ServerGroup) error {
	r.serverGroupsMutex.Lock()
	defer r.serverGroupsMutex.Unlock()

	key := r.getServerKey(group.ServerAddr, group.ServerPort)
	r.serverGroups[key] = group

	return r.saveServerGroups()
}

// GetServerGroup 获取服务器组
func (r *ProcessRepository) GetServerGroup(serverAddr string, serverPort int) (*model.ServerGroup, error) {
	r.serverGroupsMutex.RLock()
	defer r.serverGroupsMutex.RUnlock()

	key := r.getServerKey(serverAddr, serverPort)
	group, exists := r.serverGroups[key]
	if !exists {
		return nil, fmt.Errorf("服务器组不存在: %s", key)
	}

	return group, nil
}

// DeleteServerGroup 删除服务器组
func (r *ProcessRepository) DeleteServerGroup(serverAddr string, serverPort int) error {
	r.serverGroupsMutex.Lock()
	defer r.serverGroupsMutex.Unlock()

	key := r.getServerKey(serverAddr, serverPort)
	delete(r.serverGroups, key)

	return r.saveServerGroups()
}

// ListServerGroups 获取所有服务器组
func (r *ProcessRepository) ListServerGroups() ([]*model.ServerGroup, error) {
	r.serverGroupsMutex.RLock()
	defer r.serverGroupsMutex.RUnlock()

	groups := make([]*model.ServerGroup, 0, len(r.serverGroups))
	for _, group := range r.serverGroups {
		groups = append(groups, group)
	}

	return groups, nil
}

// IsServerGroupRunning 检查服务器组是否正在运行
func (r *ProcessRepository) IsServerGroupRunning(serverAddr string, serverPort int) bool {
	r.serverGroupsMutex.RLock()
	defer r.serverGroupsMutex.RUnlock()

	key := r.getServerKey(serverAddr, serverPort)
	_, exists := r.serverGroups[key]
	return exists
}
