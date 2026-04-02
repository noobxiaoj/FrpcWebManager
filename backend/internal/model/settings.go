package model

// PasswordAuthSettings 密码认证设置
// 当前使用 bcrypt 哈希存储密码。
// Password 字段仅用于兼容旧版 settings.json 中的明文密码迁移，迁移完成后会被清空。
type PasswordAuthSettings struct {
	Enabled      bool   `json:"enabled"`               // 是否已启用密码
	Username     string `json:"username"`              // 登录用户名
	PasswordHash string `json:"passwordHash"`          // bcrypt 哈希后的密码
	Password     string `json:"password,omitempty"`    // 旧版明文密码，仅用于自动迁移
}

// PasswordAuthInfo 返回给前端的密码设置信息
// 出于安全考虑，只暴露是否启用与用户名，不返回密码内容。
type PasswordAuthInfo struct {
	Enabled  bool   `json:"enabled"`  // 是否已启用密码
	Username string `json:"username"` // 当前用户名
}

// SystemSettings 系统设置模型
type SystemSettings struct {
	ShowServerPort    bool                 `json:"showServerPort"`    // 是否在服务器页面显示进程端口
	RefreshInterval   int                  `json:"refreshInterval"`   // 日志和运行时间的刷新间隔（秒）
	ShowRefreshTime   bool                 `json:"showRefreshTime"`   // 是否在服务器页面显示刷新时间
	ShowServerName    bool                 `json:"showServerName"`    // 是否在任务中显示服务器名称而非IP
	FrontendPort      int                  `json:"frontendPort"`      // 前端服务端口（默认 4500）
	EnableIPWhitelist bool                 `json:"enableIPWhitelist"` // 是否启用IP白名单
	IPWhitelist       []string             `json:"ipWhitelist"`       // IP白名单列表
	PasswordAuth      PasswordAuthSettings `json:"passwordAuth"`      // 密码认证设置
}

// DefaultSystemSettings 返回默认系统设置
func DefaultSystemSettings() *SystemSettings {
	return &SystemSettings{
		ShowServerPort:    true,  // 默认显示端口
		RefreshInterval:   10,    // 默认刷新间隔 10 秒
		ShowRefreshTime:   true,  // 默认显示刷新时间
		ShowServerName:    false, // 默认不显示服务器名称
		FrontendPort:      4500,  // 默认前端端口 4500
		EnableIPWhitelist: false, // 默认不启用IP白名单
		IPWhitelist:       []string{}, // 默认白名单为空
		PasswordAuth: PasswordAuthSettings{
			Enabled:  false,
			Username: "",
			PasswordHash: "",
			Password:     "",
		},
	}
}

// Clone 深拷贝系统设置，避免直接暴露内部引用导致状态被意外修改
func (s *SystemSettings) Clone() *SystemSettings {
	if s == nil {
		return DefaultSystemSettings()
	}

	ipWhitelist := make([]string, len(s.IPWhitelist))
	copy(ipWhitelist, s.IPWhitelist)

	return &SystemSettings{
		ShowServerPort:    s.ShowServerPort,
		RefreshInterval:   s.RefreshInterval,
		ShowRefreshTime:   s.ShowRefreshTime,
		ShowServerName:    s.ShowServerName,
		FrontendPort:      s.FrontendPort,
		EnableIPWhitelist: s.EnableIPWhitelist,
		IPWhitelist:       ipWhitelist,
		PasswordAuth: PasswordAuthSettings{
			Enabled:      s.PasswordAuth.Enabled,
			Username:     s.PasswordAuth.Username,
			PasswordHash: s.PasswordAuth.PasswordHash,
			Password:     s.PasswordAuth.Password,
		},
	}
}

// ToPasswordAuthInfo 将内部密码配置转换为前端可安全使用的展示信息
func (s *SystemSettings) ToPasswordAuthInfo() PasswordAuthInfo {
	if s == nil {
		return PasswordAuthInfo{
			Enabled:  false,
			Username: "",
		}
	}

	return PasswordAuthInfo{
		Enabled:  s.PasswordAuth.Enabled,
		Username: s.PasswordAuth.Username,
	}
}
