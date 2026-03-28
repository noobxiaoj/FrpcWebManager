package frpc

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/xiaoj/frpc_webmanager/internal/model"
)

// LogCollector 日志收集器
type LogCollector struct {
	logDir       string
	logBuffers   map[string][]string // serverKey -> log lines
	logMutex     sync.RWMutex
	logCallbacks map[string]func([]model.LogEntry)
}

// NewLogCollector 创建日志收集器
func NewLogCollector(logDir string) *LogCollector {
	return &LogCollector{
		logDir:       logDir,
		logBuffers:   make(map[string][]string),
		logCallbacks: make(map[string]func([]model.LogEntry)),
	}
}

// CollectLogs 收集进程日志
func (lc *LogCollector) CollectLogs(serverKey string, reader io.Reader, stream string) {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		lc.addLog(serverKey, line, stream)
	}

	if err := scanner.Err(); err != nil {
		lc.addLog(serverKey, fmt.Sprintf("日志读取错误: %v", err), "error")
	}
}

// addLog 添加日志
func (lc *LogCollector) addLog(serverKey string, line string, stream string) {
	lc.logMutex.Lock()
	// 解析日志级别
	level := lc.parseLogLevel(line)

	// 存储到缓冲区（最多保留1000条）
	if len(lc.logBuffers[serverKey]) > 1000 {
		lc.logBuffers[serverKey] = lc.logBuffers[serverKey][1:]
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] [%s] %s", timestamp, level, line)
	lc.logBuffers[serverKey] = append(lc.logBuffers[serverKey], logEntry)

	// 复制回调函数和日志（在锁内）
	var callback func([]model.LogEntry)
	var logs []model.LogEntry
	if cb, exists := lc.logCallbacks[serverKey]; exists {
		logs = lc.GetLogs(serverKey, 100) // 获取最新100条
		callback = cb
	}

	lc.logMutex.Unlock()

	// 在锁外调用回调，避免死锁
	if callback != nil {
		go callback(logs)
	}
}

// parseLogLevel 解析日志级别
func (lc *LogCollector) parseLogLevel(line string) string {
	lineLower := strings.ToLower(line)

	// 按优先级检查关键词
	if strings.Contains(lineLower, "error") || strings.Contains(lineLower, "fatal") || strings.Contains(lineLower, "panic") {
		return "error"
	}
	if strings.Contains(lineLower, "warn") {
		return "warn"
	}
	if strings.Contains(lineLower, "debug") || strings.Contains(lineLower, "trace") {
		return "debug"
	}

	// 如果包含info关键词，或者没有明确级别，默认为info
	if strings.Contains(lineLower, "info") {
		return "info"
	}

	return "info"
}

// WriteLog 写入日志（用于系统消息）
func (lc *LogCollector) WriteLog(serverKey string, level string, message string) {
	lc.addLog(serverKey, message, level)
}

// GetLogs 获取日志
func (lc *LogCollector) GetLogs(serverKey string, limit int) []model.LogEntry {
	lc.logMutex.RLock()

	// 先从内存缓冲区获取
	lines, exists := lc.logBuffers[serverKey]
	lc.logMutex.RUnlock()

	// 如果内存中没有日志，尝试从日志文件读取
	if !exists || len(lines) == 0 {
		lines = lc.readLogsFromFile(serverKey)
	}

	// 转换为LogEntry格式
	var logs []model.LogEntry
	start := 0
	if len(lines) > limit {
		start = len(lines) - limit
	}

	for i := start; i < len(lines); i++ {
		// 解析日志行 [timestamp] [level] message
		line := lines[i]
		parts := strings.SplitN(line, "] ", 3)
		if len(parts) >= 3 {
			timestamp := strings.Trim(parts[0], "[]")
			level := strings.Trim(parts[1], "[]")
			message := parts[2]

			logs = append(logs, model.LogEntry{
				Timestamp: timestamp,
				Level:     level,
				Message:   message,
			})
		}
	}

	return logs
}

// RegisterCallback 注册日志回调
func (lc *LogCollector) RegisterCallback(serverKey string, callback func([]model.LogEntry)) {
	lc.logMutex.Lock()
	defer lc.logMutex.Unlock()

	lc.logCallbacks[serverKey] = callback
}

// UnregisterCallback 取消注册日志回调
func (lc *LogCollector) UnregisterCallback(serverKey string) {
	lc.logMutex.Lock()
	defer lc.logMutex.Unlock()

	delete(lc.logCallbacks, serverKey)
}

// StartLogCleanup 启动日志清理协程
func (lc *LogCollector) StartLogCleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		lc.cleanupOldLogs()
	}
}

// cleanupOldLogs 清理旧日志文件
func (lc *LogCollector) cleanupOldLogs() {
	// 清理7天前的日志文件
	cutoff := time.Now().AddDate(0, 0, -7)

	err := filepath.Walk(lc.logDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && info.ModTime().Before(cutoff) {
			os.Remove(path)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("清理日志文件失败: %v\n", err)
	}
}

// ClearLogs 清空指定服务器的日志缓冲区
func (lc *LogCollector) ClearLogs(serverKey string) {
	lc.logMutex.Lock()
	defer lc.logMutex.Unlock()

	// 清空内存中的日志缓冲
	delete(lc.logBuffers, serverKey)

	// 删除日志文件
	logFile := filepath.Join(lc.logDir, fmt.Sprintf("frpc_%s.log", serverKey))
	if err := os.Remove(logFile); err != nil && !os.IsNotExist(err) {
		fmt.Printf("删除日志文件失败: %v\n", err)
	}
}

// readLogsFromFile 从日志文件读取日志
func (lc *LogCollector) readLogsFromFile(serverKey string) []string {
	// 构造日志文件路径
	logFile := filepath.Join(lc.logDir, fmt.Sprintf("frpc_%s.log", serverKey))

	// 检查文件是否存在
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		return []string{}
	}

	// 打开文件
	file, err := os.Open(logFile)
	if err != nil {
		fmt.Printf("打开日志文件失败: %v\n", err)
		return []string{}
	}
	defer file.Close()

	// 读取文件内容
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("读取日志文件失败: %v\n", err)
		return []string{}
	}

	return lines
}
