package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaoj/frpc_webmanager/internal/model"
	"github.com/xiaoj/frpc_webmanager/internal/service"
)

// SettingsHandler 系统设置处理器
type SettingsHandler struct {
	service *service.SettingsService
}

// NewSettingsHandler 创建设置处理器
func NewSettingsHandler(service *service.SettingsService) *SettingsHandler {
	return &SettingsHandler{
		service: service,
	}
}

// GetSettings 获取系统设置
// @Summary 获取系统设置
// @Description 获取系统配置
// @Tags 系统设置
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Router /api/settings [get]
func (h *SettingsHandler) GetSettings(c *gin.Context) {
	settings, err := h.service.GetSettings()
	if err != nil {
		ErrorResponse(c, CodeInternalError, "获取设置失败", err)
		return
	}

	SuccessResponse(c, settings)
}

// UpdateSettings 更新系统设置
// @Summary 更新系统设置
// @Description 更新系统配置
// @Tags 系统设置
// @Accept json
// @Produce json
// @Param settings body model.SystemSettings true "设置信息"
// @Success 200 {object} Response
// @Router /api/settings [put]
func (h *SettingsHandler) UpdateSettings(c *gin.Context) {
	var settings model.SystemSettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		ErrorResponse(c, CodeBadRequest, "请求参数错误", err)
		return
	}

	if err := h.service.UpdateSettings(&settings); err != nil {
		ErrorResponse(c, CodeInternalError, "更新设置失败", err)
		return
	}

	// 广播设置更新事件到所有连接的客户端
	SSEManagerInstance.BroadcastSettingsUpdated(settings)

	SuccessResponse(c, gin.H{
		"message": "设置已更新",
		"settings": settings,
	})
}
