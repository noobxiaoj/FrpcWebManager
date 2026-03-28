package repository

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/xiaoj/frpc_webmanager/internal/model"
)

// TaskRepository 任务数据仓储
type TaskRepository struct {
	store      *JSONStore
	tasks      map[string]*model.Task
	tasksMutex sync.RWMutex
}

// NewTaskRepository 创建任务仓储
func NewTaskRepository(store *JSONStore) *TaskRepository {
	repo := &TaskRepository{
		store: store,
		tasks: make(map[string]*model.Task),
	}

	// 加载已有任务
	repo.loadTasks()

	return repo
}

// loadTasks 从文件加载任务列表
func (r *TaskRepository) loadTasks() error {
	var taskList []*model.Task
	if err := r.store.Load("tasks.json", &taskList); err != nil {
		return err
	}

	r.tasksMutex.Lock()
	defer r.tasksMutex.Unlock()

	r.tasks = make(map[string]*model.Task)
	for _, task := range taskList {
		r.tasks[task.ID] = task
	}

	return nil
}

// saveTasks 保存任务列表到文件
// 注意:调用此函数时必须已经持有 tasksMutex 锁
func (r *TaskRepository) saveTasks() error {
	// 不再获取锁,因为调用者已经持有锁
	taskList := make([]*model.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		taskList = append(taskList, task)
	}

	return r.store.Save("tasks.json", taskList)
}

// List 获取所有任务
func (r *TaskRepository) List() ([]*model.Task, error) {
	r.tasksMutex.RLock()
	defer r.tasksMutex.RUnlock()

	tasks := make([]*model.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Get 根据ID获取任务
func (r *TaskRepository) Get(id string) (*model.Task, error) {
	r.tasksMutex.RLock()
	defer r.tasksMutex.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return nil, fmt.Errorf("任务不存在: %s", id)
	}

	return task, nil
}

// Create 创建任务
func (r *TaskRepository) Create(task *model.Task) error {
	r.tasksMutex.Lock()
	defer r.tasksMutex.Unlock()

	// 检查ID是否已存在
	if _, exists := r.tasks[task.ID]; exists {
		return fmt.Errorf("任务ID已存在: %s", task.ID)
	}

	r.tasks[task.ID] = task

	// 保存到文件
	if err := r.saveTasks(); err != nil {
		delete(r.tasks, task.ID)
		return err
	}

	// 保存详细配置到单独文件
	return r.store.Save(filepath.Join("tasks", task.ID+".json"), task)
}

// Update 更新任务
func (r *TaskRepository) Update(task *model.Task) error {
	r.tasksMutex.Lock()
	defer r.tasksMutex.Unlock()

	// 检查任务是否存在
	if _, exists := r.tasks[task.ID]; !exists {
		return fmt.Errorf("任务不存在: %s", task.ID)
	}

	r.tasks[task.ID] = task

	// 保存到文件
	if err := r.saveTasks(); err != nil {
		return err
	}

	// 更新详细配置文件
	return r.store.Save(filepath.Join("tasks", task.ID+".json"), task)
}

// Delete 删除任务
func (r *TaskRepository) Delete(id string) error {
	r.tasksMutex.Lock()
	defer r.tasksMutex.Unlock()

	// 检查任务是否存在
	if _, exists := r.tasks[id]; !exists {
		return fmt.Errorf("任务不存在: %s", id)
	}

	delete(r.tasks, id)

	// 保存到文件
	if err := r.saveTasks(); err != nil {
		return err
	}

	// 删除详细配置文件
	taskFile := filepath.Join("tasks", id+".json")
	return r.store.Delete(taskFile)
}

// GetByServer 获取连接到指定服务器的所有任务
func (r *TaskRepository) GetByServer(serverAddr string, serverPort int) ([]*model.Task, error) {
	r.tasksMutex.RLock()
	defer r.tasksMutex.RUnlock()

	var tasks []*model.Task
	for _, task := range r.tasks {
		if task.ServerAddr == serverAddr && task.ServerPort == serverPort {
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}

// UpdateStatus 更新任务状态
func (r *TaskRepository) UpdateStatus(id string, status model.TaskStatus) error {
	r.tasksMutex.Lock()
	defer r.tasksMutex.Unlock()

	task, exists := r.tasks[id]
	if !exists {
		return fmt.Errorf("任务不存在: %s", id)
	}

	task.Status = status

	return r.saveTasks()
}
