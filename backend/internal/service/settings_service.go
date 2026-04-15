package service

import (
	"errors"
	"strings"

	"github.com/xiaoj/frpc_webmanager/internal/model"
	"github.com/xiaoj/frpc_webmanager/internal/repository"
	"github.com/xiaoj/frpc_webmanager/internal/security"
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
	currentSettings, err := s.repo.GetSettings()
	if err != nil {
		return err
	}

	updatedSettings := settings.Clone()
	updatedSettings.PasswordAuth = currentSettings.PasswordAuth

	return s.repo.UpdateSettings(updatedSettings)
}

// SetPassword 首次设置密码
func (s *SettingsService) SetPassword(username string, password string) error {
	currentSettings, err := s.repo.GetSettings()
	if err != nil {
		return err
	}

	if currentSettings.PasswordAuth.Enabled {
		return errors.New("密码已存在，请使用修改密码功能")
	}

	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	if username == "" || password == "" {
		return errors.New("用户名和密码不能为空")
	}

	passwordHash, err := security.HashPassword(password)
	if err != nil {
		return err
	}

	sessionVersion, err := security.GenerateSessionVersion()
	if err != nil {
		return err
	}

	currentSettings.PasswordAuth = model.PasswordAuthSettings{
		Enabled:        true,
		Username:       username,
		PasswordHash:   passwordHash,
		Password:       "",
		SessionVersion: sessionVersion,
	}

	return s.repo.UpdateSettings(currentSettings)
}

// ChangePassword 修改已存在的密码
func (s *SettingsService) ChangePassword(oldPassword string, newPassword string) error {
	currentSettings, err := s.repo.GetSettings()
	if err != nil {
		return err
	}

	if !currentSettings.PasswordAuth.Enabled {
		return errors.New("当前未设置密码")
	}

	if strings.TrimSpace(oldPassword) == "" || strings.TrimSpace(newPassword) == "" {
		return errors.New("旧密码和新密码不能为空")
	}

	if !security.VerifyPassword(currentSettings.PasswordAuth.PasswordHash, oldPassword) {
		return errors.New("旧密码不正确")
	}

	passwordHash, err := security.HashPassword(newPassword)
	if err != nil {
		return err
	}

	sessionVersion, err := security.GenerateSessionVersion()
	if err != nil {
		return err
	}

	currentSettings.PasswordAuth.PasswordHash = passwordHash
	currentSettings.PasswordAuth.Password = ""
	currentSettings.PasswordAuth.SessionVersion = sessionVersion
	return s.repo.UpdateSettings(currentSettings)
}

// DeletePassword 删除密码设置
func (s *SettingsService) DeletePassword(username string, password string) error {
	currentSettings, err := s.repo.GetSettings()
	if err != nil {
		return err
	}

	if !currentSettings.PasswordAuth.Enabled {
		return errors.New("当前未设置密码")
	}

	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		return errors.New("用户名和密码不能为空")
	}

	if currentSettings.PasswordAuth.Username != username || !security.VerifyPassword(currentSettings.PasswordAuth.PasswordHash, password) {
		return errors.New("用户名或密码不正确")
	}

	currentSettings.PasswordAuth = model.PasswordAuthSettings{
		Enabled:        false,
		Username:       "",
		PasswordHash:   "",
		Password:       "",
		SessionVersion: "",
	}

	return s.repo.UpdateSettings(currentSettings)
}

// VerifyLogin 校验登录用户名和密码是否正确。
// 该方法仅用于页面访问认证，不会修改任何持久化数据。
//
// @param username 前端提交的登录用户名
// @param password 前端提交的登录密码
// @returns error 校验失败时返回对应错误，成功返回 nil
func (s *SettingsService) VerifyLogin(username string, password string) error {
	currentSettings, err := s.repo.GetSettings()
	if err != nil {
		return err
	}

	if !currentSettings.PasswordAuth.Enabled {
		return nil
	}

	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	if username == "" || password == "" {
		return errors.New("用户名和密码不能为空")
	}

	if currentSettings.PasswordAuth.Username != username {
		return errors.New("用户名或密码不正确")
	}

	if !security.VerifyPassword(currentSettings.PasswordAuth.PasswordHash, password) {
		return errors.New("用户名或密码不正确")
	}

	return nil
}

// GetPasswordAuthInfo 获取当前密码认证的展示信息。
// 该方法用于登录状态查询接口，只返回前端可安全使用的数据。
//
// @returns model.PasswordAuthInfo 返回是否启用密码以及当前用户名
// @returns error 获取设置失败时返回错误
func (s *SettingsService) GetPasswordAuthInfo() (model.PasswordAuthInfo, error) {
	currentSettings, err := s.repo.GetSettings()
	if err != nil {
		return model.PasswordAuthInfo{}, err
	}

	return currentSettings.ToPasswordAuthInfo(), nil
}
