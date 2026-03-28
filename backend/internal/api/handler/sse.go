package handler

import (
	"encoding/json"
	"io"
	"sync"

	"github.com/gin-gonic/gin"
)

// SSEManager SSE 连接管理器
type SSEManager struct {
	clients map[chan Event]bool
	mu      sync.RWMutex
}

// Event SSE 事件
type Event struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

// NewSSEManager 创建 SSE 管理器
func NewSSEManager() *SSEManager {
	return &SSEManager{
		clients: make(map[chan Event]bool),
	}
}

// AddClient 添加客户端
func (m *SSEManager) AddClient(client chan Event) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clients[client] = true
}

// RemoveClient 移除客户端
func (m *SSEManager) RemoveClient(client chan Event) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.clients[client]; ok {
		delete(m.clients, client)
		close(client)
	}
}

// Broadcast 广播事件到所有客户端
func (m *SSEManager) Broadcast(event Event) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for client := range m.clients {
		select {
		case client <- event:
		default:
			// 客户端通道已满,跳过
		}
	}
}

// BroadcastSettingsUpdated 广播设置更新事件
func (m *SSEManager) BroadcastSettingsUpdated(settings interface{}) {
	m.Broadcast(Event{
		Event: "settings-updated",
		Data:  settings,
	})
}

// 全局 SSE 管理器实例
var SSEManagerInstance = NewSSEManager()

// HandleSettingsEvents 处理设置更新的 SSE 连接
// @Summary 设置更新事件流
// @Description 通过 SSE 接收设置更新通知
// @Tags 系统设置
// @Produce text/event-stream
// @Router /api/settings/events [get]
func HandleSettingsEvents(c *gin.Context) {
	// 设置 SSE 响应头
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	// 创建客户端通道
	clientChan := make(chan Event, 10)

	// 添加客户端
	SSEManagerInstance.AddClient(clientChan)
	defer SSEManagerInstance.RemoveClient(clientChan)

	// 保持连接并监听事件
	c.Stream(func(w io.Writer) bool {
		select {
		case event := <-clientChan:
			// 序列化事件
			eventData, err := json.Marshal(event.Data)
			if err != nil {
				return false
			}

			// 格式化为 SSE 格式
			c.SSEvent(event.Event, string(eventData))
			return true
		case <-c.Request.Context().Done():
			// 客户端断开连接
			return false
		}
	})
}
