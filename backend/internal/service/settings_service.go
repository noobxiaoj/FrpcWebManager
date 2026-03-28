package service

import (
	"github.com/xiaoj/frpc_webmanager/internal/model"
	"github.com/xiaoj/frpc_webmanager/internal/repository"
)

// SettingsService 系统设置服务
type SettingsService struct {
	repo *repository.SettingsRepository
}

// NewSettingsService 创建设置服务
func NewSettingsService(repo *repository.SettingsRepository) *SettingsService {
	return &SettingsService{
		repo: repo,
	}
}

// GetSettings 获取系统设置
func (s *SettingsService) GetSettings() (*model.SystemSettings, error) {
	return s.repo.GetSettings()
}

// UpdateSettings 更新系统设置
func (s *SettingsService) UpdateSettings(settings *model.SystemSettings) error {
	return s.repo.UpdateSettings(settings)
}
