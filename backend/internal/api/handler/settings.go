package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaoj/frpc_webmanager/internal/model"
	"github.com/xiaoj/frpc_webmanager/internal/service"
)

// SettingsHandler 系统设置处理器
type SettingsHandler struct {
	service        *service.SettingsService
	authMiddleware interface {
		SetAuthCookie(c *gin.Context, username string)
		ClearAuthCookie(c *gin.Context)
		ValidateAuthCookie(c *gin.Context) (string, error)
	}
}

// PasswordCreateRequest 添加密码请求
type PasswordCreateRequest struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

// PasswordChangeRequest 修改密码请求
type PasswordChangeRequest struct {
	OldPassword string `json:"oldPassword"` // 旧密码
	NewPassword string `json:"newPassword"` // 新密码
}

// PasswordDeleteRequest 删除密码请求
type PasswordDeleteRequest struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

// NewSettingsHandler 创建设置处理器
func NewSettingsHandler(service *service.SettingsService, authMiddleware interface {
	SetAuthCookie(c *gin.Context, username string)
	ClearAuthCookie(c *gin.Context)
	ValidateAuthCookie(c *gin.Context) (string, error)
}) *SettingsHandler {
	return &SettingsHandler{
		service:        service,
		authMiddleware: authMiddleware,
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

	SuccessResponse(c, gin.H{
		"showServerPort":    settings.ShowServerPort,
		"refreshInterval":   settings.RefreshInterval,
		"showRefreshTime":   settings.ShowRefreshTime,
		"showServerName":    settings.ShowServerName,
		"language":          settings.Language,
		"frontendPort":      settings.FrontendPort,
		"connectionIdentifier": settings.ConnectionIdentifier,
		"enableIPWhitelist": settings.EnableIPWhitelist,
		"ipWhitelist":       settings.IPWhitelist,
		"navigationBar":     settings.NavigationBar,
		"passwordAuth":      settings.ToPasswordAuthInfo(),
	})
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

	// 保存完成后重新读取一次最终设置。
	// 这样可以把仓库层自动补齐的默认值一并返回给前端，
	// 避免旧客户端未传 language 等新字段时，响应内容与落盘结果不一致。
	savedSettings, err := h.service.GetSettings()
	if err != nil {
		ErrorResponse(c, CodeInternalError, "读取最新设置失败", err)
		return
	}

	// 广播设置更新事件到所有连接的客户端
	SSEManagerInstance.BroadcastSettingsUpdated(gin.H{
		"showServerPort":    savedSettings.ShowServerPort,
		"refreshInterval":   savedSettings.RefreshInterval,
		"showRefreshTime":   savedSettings.ShowRefreshTime,
		"showServerName":    savedSettings.ShowServerName,
		"language":          savedSettings.Language,
		"frontendPort":      savedSettings.FrontendPort,
		"connectionIdentifier": savedSettings.ConnectionIdentifier,
		"enableIPWhitelist": savedSettings.EnableIPWhitelist,
		"ipWhitelist":       savedSettings.IPWhitelist,
		"navigationBar":     savedSettings.NavigationBar,
	})

	SuccessResponse(c, gin.H{
		"message": "设置已更新",
		"settings": gin.H{
			"showServerPort":    savedSettings.ShowServerPort,
			"refreshInterval":   savedSettings.RefreshInterval,
			"showRefreshTime":   savedSettings.ShowRefreshTime,
			"showServerName":    savedSettings.ShowServerName,
			"language":          savedSettings.Language,
			"frontendPort":      savedSettings.FrontendPort,
			"connectionIdentifier": savedSettings.ConnectionIdentifier,
			"enableIPWhitelist": savedSettings.EnableIPWhitelist,
			"ipWhitelist":       savedSettings.IPWhitelist,
			"navigationBar":     savedSettings.NavigationBar,
		},
	})
}

// CreatePassword 添加密码设置
func (h *SettingsHandler) CreatePassword(c *gin.Context) {
	var request PasswordCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, CodeBadRequest, "请求参数错误", err)
		return
	}

	if err := h.service.SetPassword(request.Username, request.Password); err != nil {
		ErrorResponse(c, CodeBadRequest, err.Error(), err)
		return
	}

	// 首次设置密码时，默认将当前浏览器会话标记为已登录，
	// 避免管理员刚设置完密码后立刻被访问控制拦截。
	if h.authMiddleware != nil {
		h.authMiddleware.SetAuthCookie(c, request.Username)
	}

	settings, err := h.service.GetSettings()
	if err != nil {
		ErrorResponse(c, CodeInternalError, "获取密码状态失败", err)
		return
	}

	SSEManagerInstance.BroadcastSettingsUpdated(gin.H{
		"passwordAuth": settings.ToPasswordAuthInfo(),
	})

	SuccessResponse(c, gin.H{
		"message":      "密码已添加",
		"passwordAuth": settings.ToPasswordAuthInfo(),
	})
}

// UpdatePassword 修改密码设置
func (h *SettingsHandler) UpdatePassword(c *gin.Context) {
	var request PasswordChangeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, CodeBadRequest, "请求参数错误", err)
		return
	}

	if err := h.service.ChangePassword(request.OldPassword, request.NewPassword); err != nil {
		ErrorResponse(c, CodeBadRequest, err.Error(), err)
		return
	}

	settings, err := h.service.GetSettings()
	if err != nil {
		ErrorResponse(c, CodeInternalError, "获取密码状态失败", err)
		return
	}

	SSEManagerInstance.BroadcastSettingsUpdated(gin.H{
		"passwordAuth": settings.ToPasswordAuthInfo(),
	})

	SuccessResponse(c, gin.H{
		"message":      "密码已修改",
		"passwordAuth": settings.ToPasswordAuthInfo(),
	})
}

// DeletePassword 删除密码设置
func (h *SettingsHandler) DeletePassword(c *gin.Context) {
	var request PasswordDeleteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, CodeBadRequest, "请求参数错误", err)
		return
	}

	if err := h.service.DeletePassword(request.Username, request.Password); err != nil {
		ErrorResponse(c, CodeBadRequest, err.Error(), err)
		return
	}

	// 密码已删除时，主动清理当前浏览器的登录态 Cookie，避免残留无效状态。
	if h.authMiddleware != nil {
		h.authMiddleware.ClearAuthCookie(c)
	}

	settings, err := h.service.GetSettings()
	if err != nil {
		ErrorResponse(c, CodeInternalError, "获取密码状态失败", err)
		return
	}

	SSEManagerInstance.BroadcastSettingsUpdated(gin.H{
		"passwordAuth": settings.ToPasswordAuthInfo(),
	})

	SuccessResponse(c, gin.H{
		"message":      "密码已删除",
		"passwordAuth": settings.ToPasswordAuthInfo(),
	})
}

// GetAuthStatus 获取当前密码认证状态。
// 该接口用于前端应用启动时判断是否需要先显示登录页。
func (h *SettingsHandler) GetAuthStatus(c *gin.Context) {
	passwordAuthInfo, err := h.service.GetPasswordAuthInfo()
	if err != nil {
		ErrorResponse(c, CodeInternalError, "获取认证状态失败", err)
		return
	}

	authenticated := false
	if passwordAuthInfo.Enabled && h.authMiddleware != nil {
		if _, err := h.authMiddleware.ValidateAuthCookie(c); err == nil {
			authenticated = true
		}
	}

	SuccessResponse(c, gin.H{
		"passwordAuth":  passwordAuthInfo,
		"authenticated": authenticated || !passwordAuthInfo.Enabled,
	})
}

// Login 使用系统设置中的账号密码创建当前浏览器会话。
// 登录成功后会下发会话级 Cookie，浏览器关闭后自动失效。
func (h *SettingsHandler) Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorResponse(c, CodeBadRequest, "请求参数错误", err)
		return
	}

	if err := h.service.VerifyLogin(request.Username, request.Password); err != nil {
		ErrorResponse(c, CodeBadRequest, err.Error(), err)
		return
	}

	if h.authMiddleware != nil {
		h.authMiddleware.SetAuthCookie(c, request.Username)
	}

	passwordAuthInfo, err := h.service.GetPasswordAuthInfo()
	if err != nil {
		ErrorResponse(c, CodeInternalError, "获取认证状态失败", err)
		return
	}

	SuccessResponse(c, gin.H{
		"message":       "登录成功",
		"passwordAuth":  passwordAuthInfo,
		"authenticated": true,
	})
}

// Logout 清理当前浏览器的登录会话。
func (h *SettingsHandler) Logout(c *gin.Context) {
	if h.authMiddleware != nil {
		h.authMiddleware.ClearAuthCookie(c)
	}

	SuccessResponse(c, gin.H{
		"message": "已退出登录",
	})
}
