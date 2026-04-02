package repository

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/xiaoj/frpc_webmanager/internal/model"
	"github.com/xiaoj/frpc_webmanager/internal/security"
)

// SettingsRepository 系统设置存储仓库
type SettingsRepository struct {
	settingsPath string
	settings     *model.SystemSettings
	mu           sync.RWMutex
}

// NewSettingsRepository 创建设置仓库
func NewSettingsRepository(dataDir string) (*SettingsRepository, error) {
	settingsPath := filepath.Join(dataDir, "settings.json")

	repo := &SettingsRepository{
		settingsPath: settingsPath,
		settings:     model.DefaultSystemSettings(),
	}

	// 尝试加载已有设置
	if err := repo.load(); err != nil {
		// 如果文件不存在，创建默认设置文件
		if os.IsNotExist(err) {
			if err := repo.save(); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return repo, nil
}

// GetSettings 获取系统设置
func (r *SettingsRepository) GetSettings() (*model.SystemSettings, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.settings.Clone(), nil
}

// UpdateSettings 更新系统设置
func (r *SettingsRepository) UpdateSettings(settings *model.SystemSettings) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	normalized, err := normalizeSettings(settings)
	if err != nil {
		return err
	}

	r.settings = normalized
	return r.save()
}

// load 从文件加载设置
func (r *SettingsRepository) load() error {
	data, err := os.ReadFile(r.settingsPath)
	if err != nil {
		return err
	}

	var settings model.SystemSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return err
	}

	normalized, err := normalizeSettings(&settings)
	if err != nil {
		return err
	}
	r.settings = normalized

	// 启动迁移逻辑：
	// 如果旧版本 settings.json 中缺少新字段，或仍使用明文密码，
	// normalizeSettings 会补齐默认值并转换为 bcrypt 哈希。
	// 这里通过与原始文件内容比较，决定是否回写文件，确保迁移能自动落盘。
	normalizedData, err := json.Marshal(normalized)
	if err == nil && string(normalizedData) != string(compactJSON(data)) {
		return r.save()
	}

	return nil
}

// compactJSON 将原始 JSON 压缩为紧凑格式，便于与标准 Marshal 结果比较
func compactJSON(data []byte) []byte {
	var compacted bytes.Buffer
	if err := json.Compact(&compacted, data); err != nil {
		return data
	}

	return compacted.Bytes()
}

// normalizeSettings 统一补齐系统设置中的默认值。
// 这里既用于旧配置迁移，也用于接口写入时兜底，保证 settings.json 始终结构完整。
func normalizeSettings(settings *model.SystemSettings) (*model.SystemSettings, error) {
	defaultSettings := model.DefaultSystemSettings()
	if settings == nil {
		return defaultSettings, nil
	}

	normalized := settings.Clone()

	if normalized.FrontendPort == 0 {
		normalized.FrontendPort = defaultSettings.FrontendPort
	}

	// 新增的连接标识字段允许为空，但仍统一去掉首尾空格，
	// 避免用户误填空格后保存出看不见的差异。
	normalized.ConnectionIdentifier = strings.TrimSpace(normalized.ConnectionIdentifier)

	// 统一收敛语言字段，确保旧版 settings.json 缺失该字段时自动补齐默认值，
	// 也避免未来前端传入非法值导致配置文件中出现不可识别内容。
	switch strings.TrimSpace(normalized.Language) {
	case model.LanguageZhCN, "":
		normalized.Language = defaultSettings.Language
	case model.LanguageEnUS:
		normalized.Language = model.LanguageEnUS
	default:
		normalized.Language = defaultSettings.Language
	}

	if normalized.IPWhitelist == nil {
		normalized.IPWhitelist = []string{}
	}

	if !normalized.PasswordAuth.Enabled {
		normalized.PasswordAuth.Username = ""
		normalized.PasswordAuth.PasswordHash = ""
		normalized.PasswordAuth.Password = ""
		return normalized, nil
	}

	// 兼容旧版明文密码迁移：
	// 如果旧 settings.json 中存在 password 且尚未生成 passwordHash，
	// 则在加载阶段自动转换为 bcrypt 哈希，并清空明文字段。
	if normalized.PasswordAuth.PasswordHash == "" && normalized.PasswordAuth.Password != "" {
		passwordHash, err := security.HashPassword(normalized.PasswordAuth.Password)
		if err != nil {
			return nil, err
		}

		normalized.PasswordAuth.PasswordHash = passwordHash
		normalized.PasswordAuth.Password = ""
	}

	return normalized, nil
}

// save 保存设置到文件
func (r *SettingsRepository) save() error {
	data, err := json.MarshalIndent(r.settings, "", "  ")
	if err != nil {
		return err
	}

	// 确保目录存在
	dir := filepath.Dir(r.settingsPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(r.settingsPath, data, 0644)
}
