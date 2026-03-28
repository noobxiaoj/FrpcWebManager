package utils

import (
	"crypto/rand"
	"fmt"
	"time"
)

// GenerateUUID 生成一个简单的 UUID
func GenerateUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		// 如果随机数生成失败,使用时间戳作为fallback
		return fmt.Sprintf("task-%d", time.Now().UnixNano())
	}

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
