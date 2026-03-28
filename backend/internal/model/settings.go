package model

// SystemSettings 系统设置模型
type SystemSettings struct {
	ShowServerPort bool `json:"showServerPort"` // 是否在服务器页面显示进程端口
	RefreshInterval int `json:"refreshInterval"` // 日志和运行时间的刷新间隔（秒）
	ShowRefreshTime bool `json:"showRefreshTime"` // 是否在服务器页面显示刷新时间
	ShowServerName bool `json:"showServerName"` // 是否在任务中显示服务器名称而非IP
	FrontendPort int `json:"frontendPort"` // 前端服务端口（默认 4500）
	EnableIPWhitelist bool `json:"enableIPWhitelist"` // 是否启用IP白名单
	IPWhitelist []string `json:"ipWhitelist"` // IP白名单列表
}

// DefaultSystemSettings 返回默认系统设置
func DefaultSystemSettings() *SystemSettings {
	return &SystemSettings{
		ShowServerPort: true, // 默认显示端口
		RefreshInterval: 10, // 默认刷新间隔 10 秒
		ShowRefreshTime: true, // 默认显示刷新时间
		ShowServerName: false, // 默认不显示服务器名称
		FrontendPort: 4500, // 默认前端端口 4500
		EnableIPWhitelist: false, // 默认不启用IP白名单
		IPWhitelist: []string{}, // 默认白名单为空
	}
}
