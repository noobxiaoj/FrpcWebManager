package model

import "time"

// ProcessInfo 进程信息
type ProcessInfo struct {
	TaskID     string    `json:"taskId"`     // 关联的任务ID
	PID        int       `json:"pid"`        // 进程ID
	ConfigPath string    `json:"configPath"` // 配置文件路径
	LogPath    string    `json:"logPath"`    // 日志文件路径
	StartTime  time.Time `json:"startTime"`  // 启动时间
}

// ServerGroup 服务器分组(用于配置合并)
type ServerGroup struct {
	ServerAddr string   `json:"serverAddr"` // 服务器地址
	ServerPort int      `json:"serverPort"` // 服务器端口
	AuthToken  string   `json:"authToken"`  // 认证令牌
	TaskIDs    []string `json:"taskIds"`    // 关联的任务ID列表
	ConfigPath string   `json:"configPath"` // 合并后的配置文件路径
	LogPath    string   `json:"logPath"`    // 日志文件路径
	PID        int      `json:"pid"`        // 共享进程ID
	StartTime  time.Time `json:"startTime"` // 启动时间
}

// TaskStatusResponse 任务状态响应
type TaskStatusResponse struct {
	TaskID    string       `json:"taskId"`    // 任务ID
	Status    TaskStatus   `json:"status"`    // 任务状态
	Process   *ProcessInfo `json:"process"`   // 进程信息(如果运行中)
	ServerKey string       `json:"serverKey"` // 服务器组标识 (addr:port格式)
}
