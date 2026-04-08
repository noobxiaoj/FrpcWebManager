package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xiaoj/frpc_webmanager/internal/model"
	"github.com/xiaoj/frpc_webmanager/internal/service"
)

// ServerHandler 服务器处理器
type ServerHandler struct {
	service *service.ServerService
}

// NewServerHandler 创建服务器处理器
func NewServerHandler(service *service.ServerService) *ServerHandler {
	return &ServerHandler{
		service: service,
	}
}

// ListServers 获取服务器列表
// @Summary 获取服务器列表
// @Description 获取所有FRPC服务器列表
// @Tags 服务器管理
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Router /api/servers [get]
func (h *ServerHandler) ListServers(c *gin.Context) {
	servers, err := h.service.ListServers()
	if err != nil {
		ErrorResponse(c, CodeInternalError, "获取服务器列表失败", err)
		return
	}

	SuccessResponse(c, gin.H{
		"servers": servers,
	})
}

// GetServer 获取服务器详情
// @Summary 获取服务器详情
// @Description 根据ID获取服务器详细信息
// @Tags 服务器管理
// @Accept json
// @Produce json
// @Param id path string true "服务器ID"
// @Success 200 {object} Response
// @Router /api/servers/{id} [get]
func (h *ServerHandler) GetServer(c *gin.Context) {
	id := c.Param("id")

	server, err := h.service.GetServer(id)
	if err != nil {
		ErrorResponse(c, CodeNotFound, "服务器不存在", err)
		return
	}

	SuccessResponse(c, server)
}

// CreateServer 创建服务器
// @Summary 创建服务器
// @Description 创建新的FRPC服务器
// @Tags 服务器管理
// @Accept json
// @Produce json
// @Param server body model.CreateServerRequest true "服务器信息"
// @Success 200 {object} Response
// @Router /api/servers [post]
func (h *ServerHandler) CreateServer(c *gin.Context) {
	var req model.CreateServerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, CodeBadRequest, "请求参数错误", err)
		return
	}

	server, err := h.service.CreateServer(&req)
	if err != nil {
		ErrorResponse(c, CodeInternalError, "创建服务器失败", err)
		return
	}

	SuccessResponse(c, server)
}

// UpdateServer 更新服务器
// @Summary 更新服务器
// @Description 更新服务器信息
// @Tags 服务器管理
// @Accept json
// @Produce json
// @Param id path string true "服务器ID"
// @Param server body model.UpdateServerRequest true "服务器信息"
// @Success 200 {object} Response
// @Router /api/servers/{id} [put]
func (h *ServerHandler) UpdateServer(c *gin.Context) {
	id := c.Param("id")

	var req model.UpdateServerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, CodeBadRequest, "请求参数错误", err)
		return
	}

	server, err := h.service.UpdateServer(id, &req)
	if err != nil {
		ErrorResponse(c, CodeInternalError, "更新服务器失败", err)
		return
	}

	SuccessResponse(c, server)
}

// RestartServer 重启服务器进程
// @Summary 重启服务器进程
// @Description 按服务器当前关联任务重新启动 frpc 进程组
// @Tags 服务器管理
// @Accept json
// @Produce json
// @Param id path string true "服务器ID"
// @Success 200 {object} Response
// @Router /api/servers/{id}/restart [post]
func (h *ServerHandler) RestartServer(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.RestartServer(id); err != nil {
		ErrorResponse(c, CodeInternalError, "重启服务器失败", err)
		return
	}

	SuccessResponse(c, gin.H{
		"message": "服务器已重启",
	})
}

// DeleteServer 删除服务器
// @Summary 删除服务器
// @Description 删除指定的服务器，如果有任务需要强制删除
// @Tags 服务器管理
// @Accept json
// @Produce json
// @Param id path string true "服务器ID"
// @Param force query bool false "是否强制删除（同时删除关联任务）" default(false)
// @Success 200 {object} Response
// @Router /api/servers/{id} [delete]
func (h *ServerHandler) DeleteServer(c *gin.Context) {
	id := c.Param("id")

	// 获取 forceDelete 参数，默认为 false
	forceDeleteStr := c.Query("force")
	forceDelete := forceDeleteStr == "true" || forceDeleteStr == "1"

	taskNames, err := h.service.DeleteServer(id, forceDelete)
	if err != nil {
		// 如果服务器有任务且不是强制删除，返回任务列表
		if len(taskNames) > 0 {
			c.JSON(200, gin.H{
				"code":    1,
				"message": err.Error(),
				"data": gin.H{
					"hasTasks":  true,
					"tasks":     taskNames,
					"taskCount": len(taskNames),
				},
			})
			return
		}
		ErrorResponse(c, CodeInternalError, "删除服务器失败", err)
		return
	}

	SuccessResponse(c, gin.H{
		"message": "服务器已删除",
		"tasks":   taskNames,
	})
}

// UpdateServerOrder 更新服务器排序
// @Summary 更新服务器排序
// @Description 更新服务器的显示顺序
// @Tags 服务器管理
// @Accept json
// @Produce json
// @Param order body object{order []string} true "服务器ID顺序列表"
// @Success 200 {object} Response
// @Router /api/servers/order [put]
func (h *ServerHandler) UpdateServerOrder(c *gin.Context) {
	var req struct {
		Order []string `json:"order" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, CodeBadRequest, "请求参数错误", err)
		return
	}

	if err := h.service.UpdateServerOrder(req.Order); err != nil {
		ErrorResponse(c, CodeInternalError, "更新服务器排序失败", err)
		return
	}

	SuccessResponse(c, gin.H{
		"message": "服务器排序已更新",
	})
}

// UpdateServerLock 更新服务器锁定状态
// @Summary 更新服务器锁定状态
// @Description 更新服务器的锁定状态，锁定后不可拖动
// @Tags 服务器管理
// @Accept json
// @Produce json
// @Param id path string true "服务器ID"
// @Param lock body model.UpdateServerLockRequest true "锁定状态"
// @Success 200 {object} Response
// @Router /api/servers/{id}/lock [put]
func (h *ServerHandler) UpdateServerLock(c *gin.Context) {
	id := c.Param("id")

	var req model.UpdateServerLockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, CodeBadRequest, "请求参数错误", err)
		return
	}

	if err := h.service.UpdateServerLock(id, req.Locked); err != nil {
		ErrorResponse(c, CodeInternalError, "更新锁定状态失败", err)
		return
	}

	SuccessResponse(c, gin.H{
		"message": "锁定状态已更新",
	})
}

// AddLog 添加日志
// @Summary 添加日志
// @Description 向服务器添加日志条目
// @Tags 服务器管理
// @Accept json
// @Produce json
// @Param id path string true "服务器ID"
// @Param log body model.LogEntry true "日志条目"
// @Success 200 {object} Response
// @Router /api/servers/{id}/logs [post]
func (h *ServerHandler) AddLog(c *gin.Context) {
	id := c.Param("id")

	var log model.LogEntry
	if err := c.ShouldBindJSON(&log); err != nil {
		ErrorResponse(c, CodeBadRequest, "请求参数错误", err)
		return
	}

	if err := h.service.AddLog(id, log); err != nil {
		ErrorResponse(c, CodeInternalError, "添加日志失败", err)
		return
	}

	SuccessResponse(c, gin.H{
		"message": "日志已添加",
	})
}

// ClearLogs 清空日志
// @Summary 清空日志
// @Description 清空服务器的所有日志
// @Tags 服务器管理
// @Accept json
// @Produce json
// @Param id path string true "服务器ID"
// @Success 200 {object} Response
// @Router /api/servers/{id}/logs [delete]
func (h *ServerHandler) ClearLogs(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.ClearLogs(id); err != nil {
		ErrorResponse(c, CodeInternalError, "清空日志失败", err)
		return
	}

	SuccessResponse(c, gin.H{
		"message": "日志已清空",
	})
}

// GenerateConfig 生成服务器配置文件
// @Summary 生成服务器配置文件
// @Description 为指定服务器生成 frpc 配置文件,整合该服务器的所有任务
// @Tags 服务器管理
// @Accept json
// @Produce json
// @Param id path string true "服务器ID"
// @Success 200 {object} Response
// @Router /api/servers/{id}/config [post]
func (h *ServerHandler) GenerateConfig(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.GenerateAndSaveConfig(id); err != nil {
		ErrorResponse(c, CodeInternalError, "生成配置失败", err)
		return
	}

	// 获取配置文件路径
	configPath, err := h.service.GetConfigPath(id)
	if err != nil {
		ErrorResponse(c, CodeInternalError, "获取配置路径失败", err)
		return
	}

	SuccessResponse(c, gin.H{
		"message": "配置文件已生成",
		"configPath": configPath,
	})
}

// GetConfig 获取服务器配置文件路径
// @Summary 获取服务器配置文件路径
// @Description 获取指定服务器的 frpc 配置文件路径
// @Tags 服务器管理
// @Accept json
// @Produce json
// @Param id path string true "服务器ID"
// @Success 200 {object} Response
// @Router /api/servers/{id}/config [get]
func (h *ServerHandler) GetConfig(c *gin.Context) {
	id := c.Param("id")

	configPath, err := h.service.GetConfigPath(id)
	if err != nil {
		ErrorResponse(c, CodeInternalError, "获取配置路径失败", err)
		return
	}

	SuccessResponse(c, gin.H{
		"configPath": configPath,
	})
}

// GetLogs 获取服务器日志
// @Summary 获取服务器日志
// @Description 获取服务器的实时日志
// @Tags 服务器管理
// @Accept json
// @Produce json
// @Param id path string true "服务器ID"
// @Param limit query int false "日志条数限制" default(100)
// @Success 200 {object} Response
// @Router /api/servers/{id}/logs [get]
func (h *ServerHandler) GetLogs(c *gin.Context) {
	id := c.Param("id")

	limit := 100 // 默认值
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	logs, err := h.service.GetServerLogs(id, limit)
	if err != nil {
		ErrorResponse(c, CodeInternalError, "获取日志失败", err)
		return
	}

	SuccessResponse(c, gin.H{
		"logs": logs,
	})
}
