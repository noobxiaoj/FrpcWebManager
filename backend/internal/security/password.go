package security

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

/**
 * HashPassword 使用 bcrypt 对密码进行哈希。
 * bcrypt 会自动生成随机盐，因此不需要额外维护 salt 字段。
 *
 * @param password 原始明文密码
 * @returns string 返回可持久化保存的哈希字符串
 * @returns error 哈希失败时返回错误
 */
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

/**
 * VerifyPassword 校验明文密码是否与哈希匹配。
 *
 * @param hashedPassword 已保存的 bcrypt 哈希
 * @param plainPassword 待校验的明文密码
 * @returns bool 匹配返回 true，否则返回 false
 */
func VerifyPassword(hashedPassword string, plainPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword)) == nil
}

/**
 * GenerateSessionVersion 生成密码认证会话版本。
 * 该值会参与登录 Cookie 签名，修改密码时轮换后旧 Cookie 会自然失效。
 *
 * @returns string 返回 URL 安全的随机会话版本
 * @returns error 随机数生成失败时返回错误
 */
func GenerateSessionVersion() (string, error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(randomBytes), nil
}
