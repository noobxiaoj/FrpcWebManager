package middleware

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xiaoj/frpc_webmanager/internal/model"
	"github.com/xiaoj/frpc_webmanager/internal/service"
)

// IPWhitelistMiddleware IP白名单中间件
type IPWhitelistMiddleware struct {
	settingsService *service.SettingsService
	// 启动时的白名单配置（重启后生效）
	enableWhitelist bool
	whitelist       []string
}

// NewIPWhitelistMiddleware 创建IP白名单中间件
func NewIPWhitelistMiddleware(settingsService *service.SettingsService) *IPWhitelistMiddleware {
	// 在创建时读取一次配置
	settings, err := settingsService.GetSettings()
	if err != nil {
		// 如果读取失败，默认关闭白名单
		settings = model.DefaultSystemSettings()
	}

	return &IPWhitelistMiddleware{
		settingsService: settingsService,
		enableWhitelist: settings.EnableIPWhitelist,
		whitelist:       settings.IPWhitelist,
	}
}

// GetClientIP 获取客户端真实IP
func (m *IPWhitelistMiddleware) GetClientIP(c *gin.Context) string {
	// 1. 检查 X-Real-IP 头（Nginx 代理常用）
	if xRealIP := c.GetHeader("X-Real-IP"); xRealIP != "" {
		return xRealIP
	}

	// 2. 检查 X-Forwarded-For 头（可能包含多个IP，取第一个）
	if xForwardedFor := c.GetHeader("X-Forwarded-For"); xForwardedFor != "" {
		// X-Forwarded-For 格式: "clientIP, proxy1IP, proxy2IP"
		// 取第一个IP作为真实客户端IP
		if idx := strings.Index(xForwardedFor, ","); idx != -1 {
			return strings.TrimSpace(xForwardedFor[:idx])
		}
		return strings.TrimSpace(xForwardedFor)
	}

	// 3. 使用 RemoteAddr（直接连接时的IP）
	// RemoteAddr 格式: "IP:port"，需要提取IP部分
	if ip, _, err := net.SplitHostPort(c.Request.RemoteAddr); err == nil {
		return ip
	}

	return c.Request.RemoteAddr
}

// IsIPInWhitelist 检查IP是否在白名单中
func (m *IPWhitelistMiddleware) IsIPInWhitelist(clientIP string, whitelist []string) bool {
	// 如果白名单为空，拒绝所有访问
	if len(whitelist) == 0 {
		return false
	}

	// 解析客户端IP
	clientNetIP := net.ParseIP(clientIP)
	if clientNetIP == nil {
		return false
	}

	for _, whitelistItem := range whitelist {
		whitelistItem = strings.TrimSpace(whitelistItem)

		// 支持单个IP（如：192.168.1.100）
		if strings.Contains(whitelistItem, "/") {
			// CIDR 格式（如：192.168.1.0/24）
			_, ipNet, err := net.ParseCIDR(whitelistItem)
			if err != nil {
				continue
			}
			if ipNet.Contains(clientNetIP) {
				return true
			}
		} else {
			// 单个IP
			whitelistIP := net.ParseIP(whitelistItem)
			if whitelistIP != nil && whitelistIP.Equal(clientNetIP) {
				return true
			}
		}
	}

	return false
}

// Middleware 返回 Gin 中间件函数
func (m *IPWhitelistMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果未启用白名单，直接放行
		// 注意：这里的配置是容器启动时读取的，修改后需要重启容器才能生效
		if !m.enableWhitelist {
			c.Next()
			return
		}

		// 如果启用了白名单但白名单为空，为了防止锁死，允许所有访问
		if m.whitelist == nil || len(m.whitelist) == 0 {
			c.Next()
			return
		}

		// 获取客户端IP
		clientIP := m.GetClientIP(c)

		// 检查是否在白名单中
		if !m.IsIPInWhitelist(clientIP, m.whitelist) {
			// 直接关闭连接，不返回任何内容
			c.Abort()
			return
		}

		// 在白名单中，放行
		c.Next()
	}
}
