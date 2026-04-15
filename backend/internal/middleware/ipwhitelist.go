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
	remoteIP := parseRemoteIP(c.Request.RemoteAddr)

	// 只有请求确实来自本机反向代理时，才信任代理转发的客户端 IP。
	// 如果服务直接暴露公网，X-Real-IP / X-Forwarded-For 都可以被客户端伪造，
	// 因此不能无条件拿它们参与白名单判断。
	if isTrustedForwardingProxy(remoteIP) {
		if xRealIP := strings.TrimSpace(c.GetHeader("X-Real-IP")); xRealIP != "" {
			return xRealIP
		}

		if xForwardedFor := c.GetHeader("X-Forwarded-For"); xForwardedFor != "" {
			// X-Forwarded-For 格式: "clientIP, proxy1IP, proxy2IP"，取第一个IP作为真实客户端IP。
			if idx := strings.Index(xForwardedFor, ","); idx != -1 {
				return strings.TrimSpace(xForwardedFor[:idx])
			}
			return strings.TrimSpace(xForwardedFor)
		}
	}

	if remoteIP != "" {
		return remoteIP
	}

	return c.Request.RemoteAddr
}

// parseRemoteIP 从 RemoteAddr 中提取远端 IP。
// RemoteAddr 通常是 "IP:port" 格式；如果解析失败，则退回原始字符串，便于后续 ParseIP 再判断。
//
// @param remoteAddr 请求的 RemoteAddr
// @returns string 提取出的 IP 或原始地址
func parseRemoteIP(remoteAddr string) string {
	if ip, _, err := net.SplitHostPort(remoteAddr); err == nil {
		return ip
	}

	return strings.TrimSpace(remoteAddr)
}

// isTrustedForwardingProxy 判断是否信任该来源提供的转发 Header。
// 当前采用安全默认值：只信任本机代理。若后续需要信任 Docker 网桥或内网反代，
// 建议改为从显式配置读取可信代理网段，而不是默认信任所有内网地址。
//
// @param remoteIP 直接连接到服务的远端 IP
// @returns bool 是否允许读取 X-Real-IP / X-Forwarded-For
func isTrustedForwardingProxy(remoteIP string) bool {
	parsedIP := net.ParseIP(remoteIP)
	if parsedIP == nil {
		return false
	}

	return parsedIP.IsLoopback()
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
