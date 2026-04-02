package middleware

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xiaoj/frpc_webmanager/internal/service"
)

const passwordAuthCookieName = "frpc_webmanager_auth"

// PasswordAuthMiddleware 提供基于会话 Cookie 的页面访问认证。
// 这里使用内存中的随机密钥对 Cookie 内容签名，避免客户端伪造登录态。
// Cookie 不设置持久化过期时间，因此浏览器关闭后会自动失效，符合“重新打开网页再次输入密码”的需求。
type PasswordAuthMiddleware struct {
	settingsService *service.SettingsService
	secretKey       []byte
}

// NewPasswordAuthMiddleware 创建密码认证中间件。
//
// @param settingsService 系统设置服务，用于读取密码启用状态
// @returns *PasswordAuthMiddleware 返回初始化后的中间件实例
// @returns error 随机密钥生成失败时返回错误
func NewPasswordAuthMiddleware(settingsService *service.SettingsService) (*PasswordAuthMiddleware, error) {
	secretKey := make([]byte, 32)
	if _, err := rand.Read(secretKey); err != nil {
		return nil, err
	}

	return &PasswordAuthMiddleware{
		settingsService: settingsService,
		secretKey:       secretKey,
	}, nil
}

// Middleware 拦截受保护的 API 请求。
// 当系统开启密码认证且当前请求未携带有效会话 Cookie 时，直接返回未授权响应。
//
// @returns gin.HandlerFunc Gin 中间件函数
func (m *PasswordAuthMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.Next()
			return
		}

		if m.isPublicPath(c.Request.URL.Path, c.Request.Method) {
			c.Next()
			return
		}

		passwordAuthInfo, err := m.settingsService.GetPasswordAuthInfo()
		if err != nil {
			c.JSON(200, gin.H{
				"code":    1003,
				"message": "读取认证配置失败",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		if !passwordAuthInfo.Enabled {
			c.Next()
			return
		}

		if _, err := m.ValidateAuthCookie(c); err != nil {
			c.JSON(200, gin.H{
				"code":    1004,
				"message": "未登录或登录已失效",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SetAuthCookie 写入认证成功后的会话 Cookie。
// Cookie 中仅保存用户名和签名，不保存明文密码。
//
// @param c Gin 上下文
// @param username 当前登录用户名
func (m *PasswordAuthMiddleware) SetAuthCookie(c *gin.Context, username string) {
	trimmedUsername := strings.TrimSpace(username)
	signature := m.sign(trimmedUsername)
	cookieValue := base64.StdEncoding.EncodeToString([]byte(trimmedUsername + ":" + signature))

	// MaxAge 设置为 0，表示浏览器会话级 Cookie。
	c.SetCookie(passwordAuthCookieName, cookieValue, 0, "/", "", false, true)
}

// ClearAuthCookie 清理当前登录态 Cookie。
//
// @param c Gin 上下文
func (m *PasswordAuthMiddleware) ClearAuthCookie(c *gin.Context) {
	c.SetCookie(passwordAuthCookieName, "", -1, "/", "", false, true)
}

// ValidateAuthCookie 校验请求中的认证 Cookie 是否有效。
//
// @param c Gin 上下文
// @returns string 返回已认证的用户名
// @returns error Cookie 缺失、格式错误或签名不匹配时返回错误
func (m *PasswordAuthMiddleware) ValidateAuthCookie(c *gin.Context) (string, error) {
	cookieValue, err := c.Cookie(passwordAuthCookieName)
	if err != nil {
		return "", errors.New("认证信息不存在")
	}

	decoded, err := base64.StdEncoding.DecodeString(cookieValue)
	if err != nil {
		return "", errors.New("认证信息格式无效")
	}

	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 {
		return "", errors.New("认证信息格式错误")
	}

	username := strings.TrimSpace(parts[0])
	signature := parts[1]
	expectedSignature := m.sign(username)
	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return "", errors.New("认证签名校验失败")
	}

	return username, nil
}

// isPublicPath 判断是否为无需登录即可访问的公开接口。
// 这里只放登录流程和健康检查必须用到的端点，避免未登录状态下访问其他业务数据。
//
// @param path 请求路径
// @param method 请求方法
// @returns bool true 表示放行
func (m *PasswordAuthMiddleware) isPublicPath(path string, method string) bool {
	if path == "/api/health" && method == "GET" {
		return true
	}

	if path == "/api/settings/password" && method == "POST" {
		return true
	}

	if path == "/api/auth/status" && method == "GET" {
		return true
	}

	if path == "/api/auth/login" && method == "POST" {
		return true
	}

	if path == "/api/auth/logout" && method == "POST" {
		return true
	}

	return false
}

// sign 对用户名生成 HMAC-SHA256 签名。
// 这样客户端即使能看到 Cookie 内容，也无法伪造有效登录态。
//
// @param username 需要签名的用户名
// @returns string 返回 Base64 URL 安全编码后的签名结果
func (m *PasswordAuthMiddleware) sign(username string) string {
	mac := hmac.New(sha256.New, m.secretKey)
	mac.Write([]byte(username))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
