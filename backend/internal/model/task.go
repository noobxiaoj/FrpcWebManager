package model

import "time"

// TaskStatus 任务状态类型
type TaskStatus string

const (
	TaskStatusStopped TaskStatus = "stopped" // 已停止
	TaskStatusRunning TaskStatus = "running" // 运行中
	TaskStatusError   TaskStatus = "error"   // 错误状态
)

// ProxyType 代理类型
type ProxyType string

const (
	ProxyTypeTCP    ProxyType = "tcp"
	ProxyTypeUDP    ProxyType = "udp"
	ProxyTypeHTTP   ProxyType = "http"
	ProxyTypeHTTPS  ProxyType = "https"
	ProxyTypeTCPMUX ProxyType = "tcpmux"
	ProxyTypeSTCP   ProxyType = "stcp"
	ProxyTypeSUDP   ProxyType = "sudp"
)

// Task 任务模型
type Task struct {
	ID          string       `json:"id"`           // 任务唯一标识
	Name        string       `json:"name"`         // 任务名称
	Description string       `json:"description"`  // 任务简介
	ServerAddr  string       `json:"serverAddr"`   // frps 服务器地址
	ServerPort  int          `json:"serverPort"`   // frps 服务器端口
	AuthToken   string       `json:"authToken"`    // 认证令牌(可选)
	Proxies     []Proxy      `json:"proxies"`      // 代理配置列表
	Status      TaskStatus   `json:"status"`       // 任务状态
	CreatedAt   time.Time    `json:"createdAt"`    // 创建时间
	UpdatedAt   time.Time    `json:"updatedAt"`    // 更新时间
}

// Proxy 代理配置模型
type Proxy struct {
	ID                string            `json:"id"`                // 代理配置唯一标识
	Name              string            `json:"name"`              // 代理名称(在toml中使用)
	Type              ProxyType         `json:"type"`              // 代理类型
	LocalIP           string            `json:"localIP"`           // 本地IP
	LocalPort         int               `json:"localPort"`         // 本地端口
	RemotePort        int               `json:"remotePort"`        // 远程端口(tcp/udp)
	CustomDomains     []string          `json:"customDomains"`     // 自定义域名(http/https)
	Subdomain         string            `json:"subdomain"`         // 子域名
	Locations         []string          `json:"locations"`         // URL路径前缀(http)
	HostHeaderRewrite string            `json:"hostHeaderRewrite"` // Host头重写
	Extra             map[string]string `json:"extra"`             // 额外配置参数
}

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	Name        string  `json:"name" binding:"required"`        // 任务名称
	Description string  `json:"description"`                    // 任务简介
	ServerAddr  string  `json:"serverAddr" binding:"required"`  // frps 服务器地址
	ServerPort  int     `json:"serverPort" binding:"required"`  // frps 服务器端口
	AuthToken   string  `json:"authToken"`                      // 认证令牌
	Proxies     []Proxy `json:"proxies" binding:"required"`     // 代理配置列表
}

// UpdateTaskRequest 更新任务请求
type UpdateTaskRequest struct {
	Name        *string  `json:"name"`        // 任务名称
	Description *string  `json:"description"` // 任务简介
	ServerAddr  *string  `json:"serverAddr"`  // frps 服务器地址
	ServerPort  *int     `json:"serverPort"`  // frps 服务器端口
	AuthToken   *string  `json:"authToken"`   // 认证令牌
	Proxies     *[]Proxy `json:"proxies"`     // 代理配置列表
}
