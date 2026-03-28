package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/xiaoj/frpc_webmanager/internal/model"
	"github.com/xiaoj/frpc_webmanager/internal/service"
)

// TaskHandler 任务API处理器
type TaskHandler struct {
	taskService *service.TaskService
}

// NewTaskHandler 创建任务处理器
func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

// CreateTask 创建任务
// POST /api/tasks
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req model.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, CodeBadRequest, "请求参数错误", err)
		return
	}

	// 验证代理配置数量
	if len(req.Proxies) == 0 {
		ErrorResponse(c, CodeBadRequest, "至少需要一个代理配置", nil)
		return
	}

	// 验证服务器配置
	if err := validateServerConfig(req.ServerAddr, req.ServerPort); err != nil {
		ErrorResponse(c, CodeBadRequest, fmt.Sprintf("服务器配置无效: %s", err.Error()), nil)
		return
	}

	// 验证每个代理配置
	for i, proxy := range req.Proxies {
		if err := validateProxyConfig(proxy); err != nil {
			ErrorResponse(c, CodeBadRequest, fmt.Sprintf("第 %d 个代理配置无效: %s", i+1, err.Error()), nil)
			return
		}
	}

	// 创建任务
	task, err := h.taskService.CreateTask(&req)
	if err != nil {
		ErrorResponse(c, CodeTaskError, "创建任务失败", err)
		return
	}
	SuccessResponse(c, gin.H{
		"task": task,
	})
}

// GetTask 获取任务详情
// GET /api/tasks/:id
func (h *TaskHandler) GetTask(c *gin.Context) {
	id := c.Param("id")
	task, err := h.taskService.GetTask(id)
	if err != nil {
		ErrorResponse(c, CodeNotFound, "任务不存在", err)
		return
	}
	SuccessResponse(c, gin.H{
		"task": task,
	})
}

// ListTasks 获取任务列表
// GET /api/tasks
func (h *TaskHandler) ListTasks(c *gin.Context) {
	tasks, err := h.taskService.ListTasks()
	if err != nil {
		ErrorResponse(c, CodeInternalError, "获取任务列表失败", err)
		return
	}
	SuccessResponse(c, gin.H{
		"tasks": tasks,
		"total": len(tasks),
	})
}

// UpdateTask 更新任务
// PUT /api/tasks/:id
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var req model.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, CodeBadRequest, "请求参数错误", err)
		return
	}

	// 验证代理配置数量
	if req.Proxies != nil && len(*req.Proxies) == 0 {
		ErrorResponse(c, CodeBadRequest, "至少需要一个代理配置", nil)
		return
	}

	// 验证服务器配置（如果提供了）
	if req.ServerAddr != nil && req.ServerPort != nil {
		if err := validateServerConfig(*req.ServerAddr, *req.ServerPort); err != nil {
			ErrorResponse(c, CodeBadRequest, fmt.Sprintf("服务器配置无效: %s", err.Error()), nil)
			return
		}
	}

	// 验证每个代理配置
	if req.Proxies != nil {
		for i, proxy := range *req.Proxies {
			if err := validateProxyConfig(proxy); err != nil {
				ErrorResponse(c, CodeBadRequest, fmt.Sprintf("第 %d 个代理配置无效: %s", i+1, err.Error()), nil)
				return
			}
		}
	}

	task, err := h.taskService.UpdateTask(id, &req)
	if err != nil {
		ErrorResponse(c, CodeTaskError, "更新任务失败", err)
		return
	}
	SuccessResponse(c, gin.H{
		"task": task,
	})
}

// DeleteTask 删除任务
// DELETE /api/tasks/:id
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")

	if err := h.taskService.DeleteTask(id); err != nil {
		ErrorResponse(c, CodeTaskError, "删除任务失败", err)
		return
	}

	SuccessResponse(c, gin.H{
		"message": "任务已删除",
	})
}

// StartTask 启动任务
// POST /api/tasks/:id/start
func (h *TaskHandler) StartTask(c *gin.Context) {
	id := c.Param("id")
	if err := h.taskService.StartTask(id); err != nil {
		ErrorResponse(c, CodeTaskError, "启动任务失败", err)
		return
	}
	SuccessResponse(c, gin.H{
		"message": "任务已启动",
	})
}

// StopTask 停止任务
// POST /api/tasks/:id/stop
func (h *TaskHandler) StopTask(c *gin.Context) {
	id := c.Param("id")
	if err := h.taskService.StopTask(id); err != nil {
		ErrorResponse(c, CodeTaskError, "停止任务失败", err)
		return
	}
	SuccessResponse(c, gin.H{
		"message": "任务已停止",
	})
}

// ReloadTask 重载任务
// POST /api/tasks/:id/reload
func (h *TaskHandler) ReloadTask(c *gin.Context) {
	id := c.Param("id")
	if err := h.taskService.ReloadTask(id); err != nil {
		ErrorResponse(c, CodeTaskError, "重载任务失败", err)
		return
	}
	SuccessResponse(c, gin.H{
		"message": "任务已重载",
	})
}

// GetTaskStatus 获取任务状态
// GET /api/tasks/:id/status
func (h *TaskHandler) GetTaskStatus(c *gin.Context) {
	id := c.Param("id")
	status, err := h.taskService.GetTaskStatus(id)
	if err != nil {
		ErrorResponse(c, CodeNotFound, "获取任务状态失败", err)
		return
	}
	SuccessResponse(c, status)
}

// GetTaskConfig 获取任务配置内容
// GET /api/tasks/:id/config
func (h *TaskHandler) GetTaskConfig(c *gin.Context) {
	_ = c.Param("id")
	// TODO: 实现获取配置文件内容的逻辑
	SuccessResponse(c, gin.H{
		"message": "功能待实现",
	})
}
