package handler

import "github.com/gin-gonic/gin"

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`              // 业务状态码
	Message string      `json:"message"`           // 响应消息
	Data    interface{} `json:"data,omitempty"`    // 响应数据
	Error   string      `json:"error,omitempty"`   // 错误详情
}

// 业务状态码
const (
	CodeSuccess      = 0     // 成功
	CodeBadRequest   = 1001  // 请求参数错误
	CodeNotFound     = 1002  // 资源不存在
	CodeInternalError = 1003 // 内部错误
	CodeUnauthorized = 1004  // 未授权或未登录
	CodeTaskError    = 2001  // 任务相关错误
	CodeProcessError = 2002  // 进程相关错误
	CodeConfigError  = 2003  // 配置相关错误
)

// SuccessResponse 成功响应
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// ErrorResponse 错误响应
func ErrorResponse(c *gin.Context, code int, message string, err error) {
	resp := Response{
		Code:    code,
		Message: message,
	}

	if err != nil {
		resp.Error = err.Error()
	}

	c.JSON(200, resp)
}
