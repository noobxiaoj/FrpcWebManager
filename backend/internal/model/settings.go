package model

import "encoding/json"

const (
	// LanguageZhCN 表示简体中文。
	// 后续如果接入 i18n，可直接复用这个稳定值作为语言资源键。
	LanguageZhCN = "zh-CN"

	// LanguageEnUS 表示英语。
	// 这里使用 en-US 作为统一英文标识，避免后续多英文区域扩展时出现歧义。
	LanguageEnUS = "en-US"
)

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

// NavigationBarSettings 顶部导航栏按钮显示设置。
// 这里集中控制主界面头部中可见的快捷入口，避免用户只能通过改代码隐藏按钮。
type NavigationBarSettings struct {
	ShowAboutButton    bool `json:"showAboutButton"`    // 是否显示“关于”入口
	ShowLockButton     bool `json:"showLockButton"`     // 是否显示“锁定/退出登录”入口
	ShowLanguageButton bool `json:"showLanguageButton"` // 是否显示语言切换入口
	ShowThemeButton    bool `json:"showThemeButton"`    // 是否显示主题切换入口
}

// SystemSettings 系统设置模型
type SystemSettings struct {
	ShowServerPort    bool                 `json:"showServerPort"`    // 是否在服务器页面显示进程端口
	RefreshInterval   int                  `json:"refreshInterval"`   // 日志和运行时间的刷新间隔（秒）
	ShowRefreshTime   bool                 `json:"showRefreshTime"`   // 是否在服务器页面显示刷新时间
	ShowServerName    bool                 `json:"showServerName"`    // 是否在任务中显示服务器名称而非IP
	Language          string               `json:"language"`          // 当前界面语言，仅保存选择结果，暂不触发翻译逻辑
	FrontendPort      int                  `json:"frontendPort"`      // 前端服务端口（默认 4500）
	EnableIPWhitelist bool                 `json:"enableIPWhitelist"` // 是否启用IP白名单
	IPWhitelist       []string             `json:"ipWhitelist"`       // IP白名单列表
	NavigationBar     NavigationBarSettings `json:"navigationBar"`     // 顶部导航栏按钮显示设置
	PasswordAuth      PasswordAuthSettings `json:"passwordAuth"`      // 密码认证设置
}

// DefaultSystemSettings 返回默认系统设置
func DefaultSystemSettings() *SystemSettings {
	return &SystemSettings{
		ShowServerPort:    true,  // 默认显示端口
		RefreshInterval:   10,    // 默认刷新间隔 10 秒
		ShowRefreshTime:   true,  // 默认显示刷新时间
		ShowServerName:    false, // 默认不显示服务器名称
		Language:          LanguageZhCN, // 默认语言为简体中文
		FrontendPort:      4500,  // 默认前端端口 4500
		EnableIPWhitelist: false, // 默认不启用IP白名单
		IPWhitelist:       []string{}, // 默认白名单为空
		NavigationBar: NavigationBarSettings{
			ShowAboutButton:    true,
			ShowLockButton:     true,
			ShowLanguageButton: true,
			ShowThemeButton:    true,
		},
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
		Language:          s.Language,
		FrontendPort:      s.FrontendPort,
		EnableIPWhitelist: s.EnableIPWhitelist,
		IPWhitelist:       ipWhitelist,
		NavigationBar: NavigationBarSettings{
			ShowAboutButton:    s.NavigationBar.ShowAboutButton,
			ShowLockButton:     s.NavigationBar.ShowLockButton,
			ShowLanguageButton: s.NavigationBar.ShowLanguageButton,
			ShowThemeButton:    s.NavigationBar.ShowThemeButton,
		},
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

// UnmarshalJSON 为系统设置提供默认值兜底。
// 这样无论是加载旧版 settings.json，还是接收缺少部分字段的接口请求，
// 都能自动补齐新增字段，尤其是默认应显示的导航栏按钮配置。
func (s *SystemSettings) UnmarshalJSON(data []byte) error {
	type Alias SystemSettings

	defaultSettings := DefaultSystemSettings()
	alias := Alias(*defaultSettings)
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	*s = SystemSettings(alias)
	return nil
}
