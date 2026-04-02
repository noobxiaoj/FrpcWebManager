package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const proxyNameRandomCharset = "abcdefghijklmnopqrstuvwxyz0123456789"

// proxyNameSettings 只解析代理命名所需的最小设置字段。
// 这里不依赖完整的 settings 模型，避免把工具层与更高层结构强耦合。
type proxyNameSettings struct {
	ConnectionIdentifier string `json:"connectionIdentifier"` // 连接标识，用于拼接代理名前缀
}

// LoadConnectionIdentifier 从 settings.json 中读取连接标识。
// 如果文件不存在、字段缺失或解析失败，则返回空字符串，由上层决定回退策略。
//
// @param settingsPath - settings.json 的绝对或相对路径
// @returns string 返回去除首尾空格后的连接标识；读取失败时返回空字符串
func LoadConnectionIdentifier(settingsPath string) string {
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return ""
	}

	var settings proxyNameSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return ""
	}

	return strings.TrimSpace(settings.ConnectionIdentifier)
}

// GenerateProxyName 生成 frpc 代理名称。
// 命名规则：
// 1. 优先使用“连接标识-任务名-端口名”
// 2. 当连接标识为空时，自动回退为“4位随机字符-任务名-端口名”
// 3. 各片段会做轻量规范化，避免空格、斜杠等字符影响 frpc 日志可读性
//
// @param connectionIdentifier - 设置中的连接标识，可为空
// @param taskName - 任务名称
// @param proxyName - 端口名/代理名
// @returns string 返回最终写入 toml 的代理名称
func GenerateProxyName(connectionIdentifier string, taskName string, proxyName string) string {
	prefix := normalizeProxyNameSegment(connectionIdentifier)
	if prefix == "" {
		prefix = generateRandomProxyPrefix(4)
	}

	taskSegment := normalizeProxyNameSegment(taskName)
	if taskSegment == "" {
		taskSegment = "task"
	}

	proxySegment := normalizeProxyNameSegment(proxyName)
	if proxySegment == "" {
		proxySegment = "proxy"
	}

	return fmt.Sprintf("%s-%s-%s", prefix, taskSegment, proxySegment)
}

// normalizeProxyNameSegment 对代理名片段做轻量收敛。
// 这里只替换容易破坏可读性的分隔字符，尽量保留用户原始命名语义。
//
// @param value - 原始片段
// @returns string 规范化后的片段
func normalizeProxyNameSegment(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}

	replacer := strings.NewReplacer(
		" ", "-",
		"/", "-",
		"\\", "-",
		"_", "-",
	)

	normalized := replacer.Replace(trimmed)

	// 连续的连接符会降低日志可读性，这里统一压缩。
	for strings.Contains(normalized, "--") {
		normalized = strings.ReplaceAll(normalized, "--", "-")
	}

	return strings.Trim(normalized, "-")
}

// generateRandomProxyPrefix 生成固定长度的随机前缀。
// 当连接标识为空时，用这个值兜底，既满足唯一性，也便于在日志中区分一次配置生成。
//
// @param length - 随机字符串长度
// @returns string 返回指定长度的随机小写字母数字串
func generateRandomProxyPrefix(length int) string {
	if length <= 0 {
		return "rand"
	}

	buffer := make([]byte, length)
	randomBytes := make([]byte, length)
	if _, err := rand.Read(randomBytes); err != nil {
		// 如果随机源异常，退回到时间相关的稳定字符串并截断。
		fallback := GenerateUUID()
		if len(fallback) >= length {
			return fallback[:length]
		}
		return fallback
	}

	for index, randomByte := range randomBytes {
		buffer[index] = proxyNameRandomCharset[int(randomByte)%len(proxyNameRandomCharset)]
	}

	return string(buffer)
}
