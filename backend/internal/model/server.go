package model

import "time"

// ServerStatus 服务器状态类型
type ServerStatus string

const (
	ServerStatusOnline  ServerStatus = "online"  // 在线
	ServerStatusOffline ServerStatus = "offline" // 离线
	ServerStatusNoTask  ServerStatus = "no_task" // 无任务
)

// LogEntry 日志条目模型
type LogEntry struct {
	Timestamp string `json:"timestamp"` // 日志时间戳
	Level     string `json:"level"`     // 日志级别 (info, warn, error, debug)
	Message   string `json:"message"`   // 日志消息
}

// Server 服务器模型
type Server struct {
	ID               string       `json:"id"`               // 服务器唯一标识
	Name             string       `json:"name"`             // 服务器名称
	Address          string       `json:"address"`          // 服务器地址（包含端口）
	Status           ServerStatus `json:"status"`           // 服务器状态
	TaskCount        int          `json:"taskCount"`        // 任务数量
	Uptime           string       `json:"uptime"`           // 运行时长
	LogMaxHeight     string       `json:"logMaxHeight"`     // 日志最大高度
	WebServerPort    int          `json:"webServerPort"`    // frpc webServer 端口（用于热重载）
	Logs             []LogEntry   `json:"logs"`             // 日志列表
	Locked           bool         `json:"locked"`           // 是否锁定（锁定后不可拖动）
	CreatedAt        time.Time    `json:"createdAt"`        // 创建时间
	UpdatedAt        time.Time    `json:"updatedAt"`        // 更新时间
}

// CreateServerRequest 创建服务器请求
type CreateServerRequest struct {
	Name    string `json:"name" binding:"required"`    // 服务器名称
	Address string `json:"address" binding:"required"` // 服务器地址
	Port    string `json:"port" binding:"required"`    // 端口
	Token   string `json:"token"`                      // 密钥（可选）
}

// UpdateServerRequest 更新服务器请求
type UpdateServerRequest struct {
	Name    *string `json:"name"`    // 服务器名称
	Address *string `json:"address"` // 服务器地址
	Port    *string `json:"port"`    // 端口
	Token   *string `json:"token"`   // 密钥
}

// UpdateServerLockRequest 更新服务器锁定状态请求
type UpdateServerLockRequest struct {
	Locked bool `json:"locked"` // 是否锁定
}
