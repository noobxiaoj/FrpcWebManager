package repository

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/xiaoj/frpc_webmanager/internal/model"
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

	return r.settings, nil
}

// UpdateSettings 更新系统设置
func (r *SettingsRepository) UpdateSettings(settings *model.SystemSettings) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.settings = settings
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

	// 如果 JSON 文件中没有 FrontendPort 字段（零值），使用默认值
	if settings.FrontendPort == 0 {
		settings.FrontendPort = 4500
	}

	r.settings = &settings
	return nil
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
