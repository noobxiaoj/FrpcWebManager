<template>
  <div class="home">
    <div class="page-header">
      <h1 class="page-title">服务器列表</h1>
      <div class="header-actions">
        <AppButton
          variant="icon"
          @click="handleRefresh"
          :disabled="loading"
          :loading="loading"
          title="刷新日志"
        >
          <template #icon>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="23 4 23 10 17 10"></polyline>
              <polyline points="1 20 1 14 7 14"></polyline>
              <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"></path>
            </svg>
          </template>
        </AppButton>
        <AppButton variant="primary" @click="addServer" title="新建服务器">
          <template #icon>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="12" y1="5" x2="12" y2="19"></line>
              <line x1="5" y1="12" x2="19" y2="12"></line>
            </svg>
          </template>
          新建服务器
        </AppButton>
      </div>
    </div>

    <!-- 服务器日志列表 -->
    <div v-if="servers.length > 0 || loading" class="server-list">
      <!-- 全局加载指示器 -->
      <div v-if="loading && servers.length === 0" class="global-loading">
        <div class="loading-spinner"></div>
        <p>加载服务器列表中...</p>
      </div>

      <div
        v-for="(server, index) in servers"
        :key="server.id"
        :data-server-id="server.id"
        class="server-log-card"
        :class="{
          'dragging': isDragging && draggedIndex === index,
          'drag-over': dragOverIndex === index,
          'locked': server.locked
        }"
        :draggable="!server.locked"
        @dragstart="handleDragStart(index, $event)"
        @dragend="handleDragEnd"
        @dragover="handleDragOver(index, $event)"
        @dragleave="handleDragLeave"
        @drop="handleDrop(index, $event)"
      >
        <!-- 服务器信息头部 -->
        <ServerInfo
          :server="server"
          :settings="settings"
          @toggle-lock="toggleServerLock(server)"
          @delete="confirmDeleteServer(server)"
        />

        <!-- 日志内容区域 -->
        <div class="log-container" v-show="server.logs && server.logs.length > 0 || server.loadingLogs">
          <!-- 日志内容包装器 -->
          <div class="log-content-wrapper">
            <!-- 加载指示器 -->
            <div v-if="server.loadingLogs" class="log-loading">
              <div class="loading-spinner"></div>
              <p>加载日志中...</p>
            </div>

            <div v-else class="log-content">
              <div v-if="server.logs && server.logs.length > 0" class="log-lines">
                <div
                  v-for="(log, index) in server.logs"
                  :key="index"
                  class="log-line"
                  :class="getLogLineClass(log)"
                >
                  <span class="log-timestamp">{{ log.timestamp }}</span>
                  <span class="log-level" :class="`level-${log.level}`">{{ log.level }}</span>
                  <span class="log-message">{{ stripTimestampFromMessage(log.message) }}</span>
                </div>
              </div>
            </div>

          </div>

          <!-- 日志统计 -->
          <div class="log-footer">
            <div class="log-stats">
              <span class="stat-item">
                <span class="stat-label">总行数:</span>
                <span class="stat-value">{{ server.logs?.length || 0 }}</span>
              </span>
              <span class="stat-item stat-error">
                <span class="stat-label">错误:</span>
                <span class="stat-value">{{ getLogCount(server, 'error') }}</span>
              </span>
              <span class="stat-item stat-warn">
                <span class="stat-label">警告:</span>
                <span class="stat-value">{{ getLogCount(server, 'warn') }}</span>
              </span>
            </div>

            <div class="log-actions-right">
              <AppButton
                class="btn-clear"
                preserve-style
                @click="clearLog(server)"
                title="清空日志"
              >
                <template #icon>
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="3 6 5 6 21 6"></polyline>
                    <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                  </svg>
                </template>
              </AppButton>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else-if="isInitialized && servers.length === 0" class="empty-state">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <rect x="2" y="2" width="20" height="8" rx="2" ry="2"></rect>
        <rect x="2" y="14" width="20" height="8" rx="2" ry="2"></rect>
        <line x1="6" y1="6" x2="6" y2="6"></line>
        <line x1="6" y1="18" x2="6" y2="18"></line>
      </svg>
      <h2>暂无服务器</h2>
      <p>请先添加 FRPC 服务器</p>
    </div>

    <!-- 新建服务器模态框 -->
    <div v-if="showAddServerModal" class="modal-overlay" @click="closeAddServerModal">
      <div class="modal-card" @click.stop>
        <div class="modal-header">
          <h3 class="modal-title">新建服务器</h3>
          <AppButton class="modal-close" preserve-style @click="closeAddServerModal" title="关闭">
            <template #icon>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
              </svg>
            </template>
          </AppButton>
        </div>

        <div class="modal-body">
          <form @submit.prevent="submitAddServer" class="server-form">
            <div class="form-group">
              <label for="serverName" class="form-label">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                  <circle cx="12" cy="7" r="4"></circle>
                </svg>
                服务器名称
              </label>
              <input
                id="serverName"
                v-model="newServer.name"
                type="text"
                class="form-input"
                :class="{ 'input-error': addServerFormErrors.name }"
                placeholder="例如：生产环境服务器"
                @input="addServerFormErrors.name = ''"
              />
              <span v-if="addServerFormErrors.name" class="form-error">
                {{ addServerFormErrors.name }}
              </span>
            </div>

            <div class="form-group">
              <label for="serverAddress" class="form-label">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="2" y="2" width="20" height="8" rx="2" ry="2"></rect>
                  <rect x="2" y="14" width="20" height="8" rx="2" ry="2"></rect>
                  <line x1="6" y1="6" x2="6.01" y2="6"></line>
                  <line x1="6" y1="18" x2="6.01" y2="18"></line>
                </svg>
                服务器地址
              </label>
              <input
                id="serverAddress"
                v-model="newServer.address"
                type="text"
                class="form-input"
                :class="{ 'input-error': addServerFormErrors.address }"
                placeholder="例如：192.168.1.100 或 example.com"
                @input="addServerFormErrors.address = ''"
              />
              <span v-if="addServerFormErrors.address" class="form-error">
                {{ addServerFormErrors.address }}
              </span>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label for="serverPort" class="form-label">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="3"></circle>
                    <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
                  </svg>
                  端口
                </label>
                <input
                  id="serverPort"
                  v-model="newServer.port"
                  type="text"
                  class="form-input"
                  :class="{ 'input-error': addServerFormErrors.port }"
                  placeholder="7000"
                  @input="addServerFormErrors.port = ''"
                />
                <span v-if="addServerFormErrors.port" class="form-error">
                  {{ addServerFormErrors.port }}
                </span>
              </div>
            </div>

            <div class="form-group">
              <label for="serverToken" class="form-label">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
                  <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
                </svg>
                密钥（可选）
              </label>
              <input
                id="serverToken"
                v-model="newServer.token"
                type="password"
                class="form-input"
                :class="{ 'input-error': addServerFormErrors.token }"
                placeholder="请输入服务器密钥（可选）"
                @input="addServerFormErrors.token = ''"
              />
              <span v-if="addServerFormErrors.token" class="form-error">
                {{ addServerFormErrors.token }}
              </span>
            </div>

            <div class="form-actions">
              <AppButton variant="secondary" type="button" @click="closeAddServerModal">
                取消
              </AppButton>
              <AppButton variant="primary" type="submit">
                <template #icon>
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <line x1="12" y1="5" x2="12" y2="19"></line>
                    <line x1="5" y1="12" x2="19" y2="12"></line>
                  </svg>
                </template>
                添加服务器
              </AppButton>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- 删除确认模态框 -->
    <div v-if="showDeleteConfirmModal" class="modal-overlay" @click="closeDeleteConfirmModal">
      <div class="modal-card modal-delete-confirm" @click.stop>
        <div class="modal-header">
          <h3 class="modal-title">确认删除</h3>
          <AppButton class="modal-close" preserve-style @click="closeDeleteConfirmModal" title="关闭">
            <template #icon>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
              </svg>
            </template>
          </AppButton>
        </div>

        <div class="modal-body">
          <div class="delete-confirm-content">
            <div class="delete-warning-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path>
                <line x1="12" y1="9" x2="12" y2="13"></line>
                <line x1="12" y1="17" x2="12.01" y2="17"></line>
              </svg>
            </div>
            <h4 class="delete-confirm-title">删除服务器</h4>

            <!-- 如果有任务，显示任务列表 -->
            <div v-if="serverTasks.length > 0">
              <p class="delete-confirm-message">
                服务器 <strong>{{ serverToDelete?.name }}</strong> 还有 <strong>{{ serverTasks.length }}</strong> 个任务：
              </p>
              <ul class="task-list">
                <li v-for="(task, index) in serverTasks" :key="index" class="task-item">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="9 11 12 14 22 4"></polyline>
                    <path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"></path>
                  </svg>
                  {{ task }}
                </li>
              </ul>
              <p class="delete-confirm-warning">确认删除将同时删除以上任务，此操作无法撤销！</p>
            </div>

            <!-- 如果没有任务，显示普通删除确认 -->
            <div v-else>
              <p class="delete-confirm-message">
                确定要删除服务器 <strong>{{ serverToDelete?.name }}</strong> 吗？
              </p>
              <p class="delete-confirm-warning">此操作无法撤销，删除后所有相关日志将被永久删除。</p>
            </div>
          </div>

          <div class="form-actions">
            <AppButton variant="secondary" type="button" @click="closeDeleteConfirmModal">
              取消
            </AppButton>
            <AppButton variant="danger" type="button" @click="deleteServer">
              <template #icon>
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="3 6 5 6 21 6"></polyline>
                  <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                </svg>
              </template>
              {{ serverTasks.length > 0 ? '确认删除服务器和任务' : '确认删除' }}
            </AppButton>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { serverAPI } from '@/api/server'
import { getApiBaseUrl } from '@/utils/api'
import AppButton from '@/components/AppButton.vue'
import ServerInfo from '@/components/ServerInfo.vue'

const loading = ref(false)
const refreshIntervals = ref({})
const currentTime = ref(Date.now())
const isInitialized = ref(false) // 标记是否已经初始化加载过

// 系统设置
const settings = ref({
  showServerPort: true,
  refreshInterval: 10,
  showRefreshTime: true
})

// 服务器添加模态框状态
const showAddServerModal = ref(false)
const newServer = ref({
  name: '',
  address: '',
  port: '',
  token: ''
})
const addServerFormErrors = ref({})

// 删除确认模态框状态
const showDeleteConfirmModal = ref(false)
const serverToDelete = ref(null)
const serverTasks = ref([]) // 服务器关联的任务列表

// 服务器日志数据（从后端API获取）
const servers = ref([])

// 拖放相关
const draggedIndex = ref(null)
const dragOverIndex = ref(null)
const isDragging = ref(false)

const handleRefresh = async () => {
  await loadServers()
}

const addServer = () => {
  // 打开新建服务器模态框
  showAddServerModal.value = true
  // 重置表单
  newServer.value = {
    name: '',
    address: '',
    port: '',
    token: ''
  }
  addServerFormErrors.value = {}
}

const closeAddServerModal = () => {
  showAddServerModal.value = false
  newServer.value = {
    name: '',
    address: '',
    port: '',
    token: ''
  }
  addServerFormErrors.value = {}
}

const validateServerForm = () => {
  const errors = {}

  if (!newServer.value.name.trim()) {
    errors.name = '请输入服务器名称'
  }

  if (!newServer.value.address.trim()) {
    errors.address = '请输入服务器地址'
  } else if (!/^[\w\.-]+$/.test(newServer.value.address)) {
    errors.address = '请输入有效的服务器地址'
  }

  if (!newServer.value.port.trim()) {
    errors.port = '请输入端口号'
  } else if (!/^\d{1,5}$/.test(newServer.value.port) || parseInt(newServer.value.port) > 65535) {
    errors.port = '请输入有效的端口号 (1-65535)'
  }

  // 密钥是可选项，不需要验证

  addServerFormErrors.value = errors
  return Object.keys(errors).length === 0
}

const submitAddServer = async () => {
  if (!validateServerForm()) {
    return
  }

  try {
    // 调用API创建服务器
    const server = await serverAPI.createServer({
      name: newServer.value.name,
      address: newServer.value.address,
      port: newServer.value.port,
      token: newServer.value.token
    })

    // 添加到本地列表
    servers.value.push(server)

    // 关闭模态框
    closeAddServerModal()

    console.log('服务器已添加:', server)
  } catch (error) {
    console.error('创建服务器失败:', error)
    alert('创建服务器失败: ' + error.message)
  }
}

const confirmDeleteServer = async (server) => {
  serverToDelete.value = server
  serverTasks.value = []

  // 先尝试删除，检查是否有任务
  try {
    const result = await serverAPI.deleteServer(server.id, false)

    // 如果有任务，显示任务列表
    if (result.hasTasks) {
      serverTasks.value = result.tasks
      showDeleteConfirmModal.value = true
    } else {
      // 如果没有任务，直接删除成功
      console.log('服务器已删除:', server.name)

      // 清理该服务器的自动刷新定时器
      if (refreshIntervals.value[server.id]) {
        clearInterval(refreshIntervals.value[server.id])
        delete refreshIntervals.value[server.id]
      }

      // 从列表中删除服务器
      const index = servers.value.findIndex(s => s.id === server.id)
      if (index > -1) {
        servers.value.splice(index, 1)
      }

      serverToDelete.value = null
    }
  } catch (error) {
    console.error('检查服务器任务失败:', error)
    alert('检查服务器任务失败: ' + error.message)
    serverToDelete.value = null
  }
}

const closeDeleteConfirmModal = () => {
  showDeleteConfirmModal.value = false
  serverToDelete.value = null
  serverTasks.value = []
}

const deleteServer = async () => {
  if (!serverToDelete.value) return

  // 保存服务器名称，用于后续日志输出
  const serverName = serverToDelete.value.name
  const serverId = serverToDelete.value.id

  try {
    // 调用API删除服务器（force=true 表示同时删除关联任务）
    const result = await serverAPI.deleteServer(serverId, true)

    // 清理该服务器的自动刷新定时器
    if (refreshIntervals.value[serverId]) {
      clearInterval(refreshIntervals.value[serverId])
      delete refreshIntervals.value[serverId]
    }

    // 从列表中删除服务器
    const index = servers.value.findIndex(s => s.id === serverId)
    if (index > -1) {
      servers.value.splice(index, 1)
    }

    // 关闭模态框
    closeDeleteConfirmModal()

    // 显示删除成功的消息
    const deletedTasks = result.tasks || []
    if (deletedTasks.length > 0) {
      alert(`服务器 "${serverName}" 已删除，同时删除了 ${deletedTasks.length} 个任务:\n${deletedTasks.join('\n')}`)
    } else {
      console.log('服务器已删除:', serverName)
    }
  } catch (error) {
    console.error('删除服务器失败:', error)
    alert('删除服务器失败: ' + error.message)
  }
}

const updateServerLog = async (server) => {
  try {
    // 检查用户是否在查看历史日志（滚动位置不在底部）
    const wasNearBottom = isLogNearBottom(server)

    // 从后端获取最新的日志
    const logs = await serverAPI.getLogs(server.id, 100)

    // 更新服务器日志
    server.logs = logs

    // 只有当用户之前就在底部时，才自动滚动到底部
    // 这样用户查看历史日志时不会被自动滚动打断
    if (wasNearBottom) {
      scrollLogToBottom(server)
    }
  } catch (error) {
    console.error('更新日志失败:', error)
    // 如果获取失败，添加错误日志
    const errorLog = {
      timestamp: new Date().toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        hour12: false
      }).replace(/\//g, '-'),
      level: 'error',
      message: `获取日志失败: ${error.message}`
    }
    server.logs = server.logs || []
    server.logs.push(errorLog)
    if (server.logs.length > 100) {
      server.logs = server.logs.slice(-100)
    }
  }

  // 同时更新服务器的运行时间和刷新时间
  try {
    const servers = await serverAPI.listServers()
    const updatedServer = servers.find(s => s.id === server.id)
    if (updatedServer && updatedServer.uptime) {
      server.uptime = updatedServer.uptime
    }
    // 更新刷新时间戳
    server.lastRefreshTime = new Date().toISOString()
  } catch (error) {
    console.error('更新运行时间失败:', error)
  }
}

const clearLog = async (server) => {
  try {
    // 调用API清空日志
    await serverAPI.clearLogs(server.id)

    // 更新本地状态
    server.logs = []
  } catch (error) {
    console.error('清空日志失败:', error)
    alert('清空日志失败: ' + error.message)
  }
}

// 滚动日志到底部
const scrollLogToBottom = (server) => {
  nextTick(() => {
    const logContainer = document.querySelector(`[data-server-id="${server.id}"] .log-content`)
    if (logContainer) {
      logContainer.scrollTop = logContainer.scrollHeight
    }
  })
}

// 检查日志是否已经在底部（接近底部50px以内）
const isLogNearBottom = (server) => {
  const logContainer = document.querySelector(`[data-server-id="${server.id}"] .log-content`)
  if (!logContainer) return true

  const threshold = 50 // 距离底部50px以内认为是在底部
  return logContainer.scrollHeight - logContainer.scrollTop - logContainer.clientHeight < threshold
}

const getLogLineClass = (log) => {
  return `log-line-${log.level}`
}

// 去除日志消息中的时间戳和ANSI颜色代码
const stripTimestampFromMessage = (message) => {
  if (!message) return ''
  let result = message

  // 去除ANSI颜色代码（如 [0m[1;34m 等）
  // 匹配 \x1b[ 或 \033[ 或 [ 开头的转义序列
  result = result.replace(/\x1b\[[0-9;]*m/g, '')
  result = result.replace(/\033\[[0-9;]*m/g, '')
  result = result.replace(/\[[0-9;]*m/g, '')

  // 去除时间戳
  // 匹配常见的时间戳格式：
  // 2026-02-06 04:57:19]
  // [2026-02-06 04:57:19]
  // [2026-02-06 04:57:19
  // 2026-02-06 04:57:19
  result = result.replace(/^\[?\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2}(\.\d+)?\]?\s*/, '')

  return result
}

const getLogCount = (server, level) => {
  if (!server.logs) return 0
  return server.logs.filter(log => log.level === level).length
}

// 切换服务器锁定状态
const toggleServerLock = async (server) => {
  try {
    const newLockedState = !server.locked
    server.locked = newLockedState

    // 调用API保存锁定状态
    await serverAPI.updateServerLock(server.id, newLockedState)

    console.log(`服务器 ${server.name} ${newLockedState ? '已锁定' : '已解锁'}`)
  } catch (error) {
    console.error('更新锁定状态失败:', error)
    // 如果失败,恢复原来的状态
    server.locked = !server.locked
    alert('更新锁定状态失败: ' + error.message)
  }
}

// 拖放事件处理
const handleDragStart = (index, event) => {
  draggedIndex.value = index
  isDragging.value = true
  event.dataTransfer.effectAllowed = 'move'
  event.dataTransfer.setData('text/html', event.target.innerHTML)

  // 创建自定义拖动预览
  const dragElement = event.target.cloneNode(true)
  dragElement.style.opacity = '0.8'
  dragElement.style.transform = 'rotate(3deg)'
  dragElement.style.boxShadow = '0 10px 30px rgba(0,0,0,0.3)'
  document.body.appendChild(dragElement)
  event.dataTransfer.setDragImage(dragElement, 200, 50)

  // 延迟移除克隆元素
  setTimeout(() => {
    document.body.removeChild(dragElement)
  }, 0)
}

const handleDragEnd = () => {
  isDragging.value = false
  draggedIndex.value = null
  dragOverIndex.value = null
}

const handleDragOver = (index, event) => {
  event.preventDefault()
  event.dataTransfer.dropEffect = 'move'
  dragOverIndex.value = index
}

const handleDragLeave = () => {
  // 延迟清除,避免闪烁
  setTimeout(() => {
    dragOverIndex.value = null
  }, 100)
}

const handleDrop = async (dropIndex, event) => {
  event.preventDefault()
  const dragIndex = draggedIndex.value

  if (dragIndex === null || dragIndex === dropIndex) {
    dragOverIndex.value = null
    return
  }

  try {
    // 保存被拖动的服务器
    const [draggedServer] = servers.value.splice(dragIndex, 1)
    servers.value.splice(dropIndex, 0, draggedServer)

    // 调用API保存新的顺序
    const serverOrder = servers.value.map(s => s.id)
    await serverAPI.updateServerOrder(serverOrder)

    console.log('服务器顺序已更新')
  } catch (error) {
    console.error('更新服务器顺序失败:', error)
    // 如果失败,恢复原来的顺序
    const serverOrder = await serverAPI.listServers()
    servers.value = serverOrder
  }

  isDragging.value = false
  draggedIndex.value = null
  dragOverIndex.value = null
}

// 监听刷新间隔设置的变化
watch(() => settings.value.refreshInterval, (newInterval, oldInterval) => {
  console.log(`watch 触发: oldInterval=${oldInterval}, newInterval=${newInterval}`)

  // 如果刷新间隔发生变化，重新设置所有服务器的定时器
  if (oldInterval !== undefined && newInterval !== oldInterval) {
    console.log(`刷新间隔从 ${oldInterval}s 改为 ${newInterval}s，重新设置定时器`)
    console.log('当前服务器列表:', servers.value)

    // 重新设置所有服务器的定时器
    servers.value.forEach(server => {
      console.log(`处理服务器 ${server.id}: ${server.name}`)

      // 清除旧的定时器
      if (refreshIntervals.value[server.id]) {
        console.log(`清除服务器 ${server.id} 的旧定时器`)
        clearInterval(refreshIntervals.value[server.id])
        delete refreshIntervals.value[server.id]
      }

      // 如果新的间隔大于0，启动新的定时器
      if (newInterval > 0) {
        console.log(`为服务器 ${server.id} 启动新定时器，间隔 ${newInterval}s`)
        const intervalMs = newInterval * 1000
        refreshIntervals.value[server.id] = setInterval(async () => {
          await updateServerLog(server)
        }, intervalMs)
      } else {
        console.log(`服务器 ${server.id} 不启动定时器（间隔为 0）`)
      }
    })

    console.log('所有定时器:', refreshIntervals.value)
  }
})

onMounted(() => {
  // 页面加载时先获取系统设置，再加载服务器列表
  loadSettings().then(() => {
    loadServers()
  })

  // 每3秒更新当前时间，用于计算刷新时长
  const timeUpdateInterval = setInterval(() => {
    currentTime.value = Date.now()
  }, 3000)

  // 保存定时器引用以便清理
  refreshIntervals.value['_timeUpdate'] = timeUpdateInterval

  // 监听设置更新事件
  const eventSource = new EventSource(`${getApiBaseUrl()}/api/settings/events`)
  eventSource.addEventListener('settings-updated', () => {
    console.log('收到设置更新通知')
    loadSettings()
  })
  eventSource.onerror = (error) => {
    console.error('SSE 连接错误:', error)
    eventSource.close()
  }

  // 保存 EventSource 引用以便清理
  refreshIntervals.value['_settingsEventSource'] = eventSource
})

// 加载系统设置
const loadSettings = async () => {
  try {
    const API_BASE_URL = getApiBaseUrl()
    const response = await fetch(`${API_BASE_URL}/api/settings`)
    if (!response.ok) {
      throw new Error('获取设置失败')
    }
    const data = await response.json()
    if (data.data) {
      console.log('加载到设置:', data.data)
      settings.value = data.data
      console.log('当前刷新间隔:', settings.value.refreshInterval)
    }
  } catch (error) {
    console.error('加载设置失败:', error)
    // 使用默认值
    settings.value = {
      showServerPort: true,
      refreshInterval: 10,
      showRefreshTime: true
    }
    console.log('使用默认设置，刷新间隔:', settings.value.refreshInterval)
  }
}

// 加载单个服务器的日志
const loadServerLogs = async (server) => {
  try {
    // 设置加载状态
    server.loadingLogs = true
    const logs = await serverAPI.getLogs(server.id, 100)
    server.logs = logs
    server.lastRefreshTime = new Date().toISOString()

    // 滚动到底部
    scrollLogToBottom(server)
  } catch (error) {
    console.error(`加载服务器 ${server.name} 日志失败:`, error)
    server.logs = server.logs || []
  } finally {
    // 清除加载状态
    server.loadingLogs = false
  }
}

// 加载服务器列表
const loadServers = async () => {
  try {
    loading.value = true
    const data = await serverAPI.listServers()

    // 保留现有的 lastRefreshTime 和日志
    const existingServers = servers.value.reduce((map, server) => {
      map[server.id] = {
        lastRefreshTime: server.lastRefreshTime,
        logs: server.logs
      }
      return map
    }, {})

    servers.value = data

    // 先设置每个服务器的基本信息和定时器（不等待日志加载）
    servers.value.forEach(server => {
      // 初始化加载状态
      server.loadingLogs = true

      // 恢复之前的日志数据（如果有）
      if (existingServers[server.id]) {
        server.lastRefreshTime = existingServers[server.id].lastRefreshTime
        server.logs = existingServers[server.id].logs
      }

      // 清除该服务器的旧定时器（如果存在）
      if (refreshIntervals.value[server.id]) {
        clearInterval(refreshIntervals.value[server.id])
        delete refreshIntervals.value[server.id]
      }

      // 为每个服务器启动自动刷新，使用配置的刷新间隔
      // 如果 refreshInterval 为 0，则不启动自动刷新
      const interval = settings.value.refreshInterval ?? 10
      if (interval > 0) {
        const intervalMs = interval * 1000
        refreshIntervals.value[server.id] = setInterval(async () => {
          await updateServerLog(server)
        }, intervalMs)
      }
    })

    // 标记已初始化，立即结束 loading 状态，让用户看到服务器列表
    isInitialized.value = true
    loading.value = false

    // 在后台异步加载所有服务器的日志（不阻塞UI渲染）
    // 使用 setTimeout 让出主线程，确保UI先渲染
    setTimeout(async () => {
      const loadPromises = servers.value.map(server => loadServerLogs(server))
      await Promise.all(loadPromises)
    }, 0)
  } catch (error) {
    console.error('加载服务器列表失败:', error)
    isInitialized.value = true
    loading.value = false
  }
}

onUnmounted(() => {
  // 清理所有定时器
  Object.values(refreshIntervals.value).forEach(interval => {
    clearInterval(interval)
  })

  // 关闭 EventSource 连接
  if (refreshIntervals.value['_settingsEventSource']) {
    refreshIntervals.value['_settingsEventSource'].close()
  }
})
</script>

<style scoped>
.home {
  max-width: 1600px;
  margin: 0 auto;
  padding: 2rem;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
  margin-bottom: 2rem;
}

.page-title {
  font-size: 2rem;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.server-list {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.server-log-card {
  background: var(--card-bg);
  border-radius: 16px;
  box-shadow: 0 4px 12px var(--shadow-color);
  overflow: hidden;
  border: 1px solid var(--border-color);
  transition: all 0.3s;
  position: relative;
}

.server-log-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg,
    var(--accent-color) 0%,
    var(--success-color) 50%,
    var(--accent-color) 100%);
  opacity: 0;
  transition: opacity 0.3s;
}

.server-log-card:hover {
  box-shadow: 0 12px 24px var(--shadow-color);
  transform: translateY(-2px);
  border-color: var(--accent-color);
}

.server-log-card:hover::before {
  opacity: 1;
}

.server-log-card.dragging {
  opacity: 0.4;
  transform: scale(0.98) rotate(1deg);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
  cursor: grabbing;
}

.server-log-card.drag-over {
  border-color: var(--accent-color);
  background: var(--accent-color-bg);
  transform: scale(1.01);
  box-shadow: 0 0 0 3px var(--accent-color-bg), 0 12px 24px var(--shadow-color);
}

.server-log-card[draggable="true"] {
  cursor: grab;
}

.server-log-card[draggable="true"]:active {
  cursor: grabbing;
}

/* 锁定状态的卡片样式 */
.server-log-card.locked {
  opacity: 0.95;
  cursor: not-allowed;
}

/* 日志容器 */
.log-container {
  display: flex;
  flex-direction: column;
}

/* 日志内容包装器 */
.log-content-wrapper {
  position: relative;
  display: flex;
  flex-direction: column;
}

/* 日志内容 - 重新设计 */
.log-content {
  overflow-y: auto;
  background: linear-gradient(180deg, var(--bg-primary) 0%, var(--bg-secondary) 100%);
  padding: 0;
  font-family: 'JetBrains Mono', 'Fira Code', 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 0.8rem;
  line-height: 1.4;
  position: relative;
  /* 设置最大高度为10行的高度 (每行约1.4rem + padding) */
  max-height: calc(1.4rem * 10 + 1.5rem);
}

.log-content::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(90deg,
    transparent 0%,
    var(--border-color) 50%,
    transparent 100%);
}

.log-content::-webkit-scrollbar {
  width: 8px;
}

.log-content::-webkit-scrollbar-track {
  background: var(--bg-primary);
  border-left: 1px solid var(--border-color);
}

.log-content::-webkit-scrollbar-thumb {
  background: linear-gradient(180deg, var(--border-color), var(--text-secondary));
  border-radius: 4px;
  border: 2px solid var(--bg-primary);
}

.log-content::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(180deg, var(--text-secondary), var(--accent-color));
}

.log-lines {
  display: flex;
  flex-direction: column;
  padding: 0.75rem 1rem;
}

.log-line {
  display: grid;
  grid-template-columns: auto auto 1fr;
  gap: 0.75rem;
  padding: 0.25rem 0.5rem;
  margin: 0;
  border-radius: 4px;
  transition: all 0.2s ease;
  position: relative;
  border-left: 2px solid transparent;
}

.log-line::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 2px;
  background: var(--border-color);
  border-radius: 2px 0 0 2px;
  transition: all 0.2s;
}

.log-line:hover {
  background: var(--bg-primary);
  transform: translateX(2px);
  box-shadow: 0 1px 4px var(--shadow-color);
}

.log-line:hover::before {
  background: var(--accent-color);
  width: 3px;
}

.log-line-log-line-error {
  background: rgba(var(--danger-color-rgb, 239, 68, 68), 0.05);
}

.log-line-log-line-error:hover {
  background: rgba(var(--danger-color-rgb, 239, 68, 68), 0.1);
}

.log-line-log-line-error::before {
  background: var(--danger-color);
}

.log-line-log-line-warn {
  background: rgba(var(--warning-color-rgb, 245, 158, 11), 0.05);
}

.log-line-log-line-warn:hover {
  background: rgba(var(--warning-color-rgb, 245, 158, 11), 0.1);
}

.log-line-log-line-warn::before {
  background: var(--warning-color);
}

.log-timestamp {
  color: var(--text-secondary);
  flex-shrink: 0;
  font-size: 0.75rem;
  font-family: 'SF Mono', 'Monaco', 'Consolas', monospace;
  opacity: 0.7;
  letter-spacing: 0.2px;
  /* 在 grid 布局中顶部对齐，不随内容行数拉伸 */
  align-self: start;
}

.log-level {
  padding: 0.2rem 0.55rem;
  border-radius: 6px;
  font-size: 0.65rem;
  font-weight: 600;
  flex-shrink: 0;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
  min-width: 48px;
  text-align: center;
  transition: all 0.25s ease;
  /* 在 grid 布局中顶部对齐，不随内容行数拉伸 */
  align-self: start;
  /* 固定高度，不随内容变化 */
  height: fit-content;
  position: relative;
  overflow: hidden;
  border: 2px solid transparent;
}

/* 所有日志级别的光效效果 */
.log-level::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    90deg,
    transparent,
    rgba(255, 255, 255, 0.2),
    transparent
  );
  transition: left 0.5s ease;
}

.log-level:hover::before {
  left: 100%;
}

/* INFO 级别 - 蓝色渐变 */
.level-info {
  background: linear-gradient(135deg, var(--info-color) 0%, color-mix(in srgb, var(--info-color) 85%, black) 100%);
  color: white;
  box-shadow:
    0 2px 4px rgba(var(--info-color-rgb), 0.3),
    inset 0 1px 0 rgba(255, 255, 255, 0.2);
  border-color: color-mix(in srgb, var(--info-color) 60%, white);
}

.level-info:hover {
  transform: translateY(-1px);
  box-shadow:
    0 4px 8px rgba(var(--info-color-rgb), 0.4),
    inset 0 1px 0 rgba(255, 255, 255, 0.3);
}

/* WARN 级别 - 橙色渐变 */
.level-warn {
  background: linear-gradient(135deg, var(--warning-color) 0%, color-mix(in srgb, var(--warning-color) 85%, black) 100%);
  color: white;
  box-shadow:
    0 2px 4px rgba(var(--warning-color-rgb), 0.3),
    inset 0 1px 0 rgba(255, 255, 255, 0.2);
  border-color: color-mix(in srgb, var(--warning-color) 60%, white);
}

.level-warn:hover {
  transform: translateY(-1px);
  box-shadow:
    0 4px 8px rgba(var(--warning-color-rgb), 0.4),
    inset 0 1px 0 rgba(255, 255, 255, 0.3);
}

/* ERROR 级别 - 红色/紫色渐变（根据主题） */
.level-error {
  background: linear-gradient(135deg, var(--danger-color) 0%, color-mix(in srgb, var(--danger-color) 85%, black) 100%);
  color: white;
  box-shadow:
    0 2px 4px rgba(var(--danger-color-rgb), 0.35),
    inset 0 1px 0 rgba(255, 255, 255, 0.25);
  border-color: color-mix(in srgb, var(--danger-color) 60%, white);
  animation: errorPulse 2.5s ease-in-out infinite;
}

.level-error:hover {
  transform: translateY(-1px);
  box-shadow:
    0 5px 10px rgba(var(--danger-color-rgb), 0.5),
    inset 0 1px 0 rgba(255, 255, 255, 0.35);
}

@keyframes errorPulse {
  0%, 100% {
    box-shadow:
      0 2px 4px rgba(var(--danger-color-rgb), 0.35),
      inset 0 1px 0 rgba(255, 255, 255, 0.25);
  }
  50% {
    box-shadow:
      0 3px 8px rgba(var(--danger-color-rgb), 0.5),
      0 0 15px rgba(var(--danger-color-rgb), 0.3),
      inset 0 1px 0 rgba(255, 255, 255, 0.3);
  }
}

/* DEBUG 级别 - 灰色渐变 */
.level-debug {
  background: linear-gradient(135deg, var(--debug-color) 0%, color-mix(in srgb, var(--debug-color) 85%, black) 100%);
  color: white;
  box-shadow:
    0 2px 4px rgba(var(--debug-color-rgb), 0.25),
    inset 0 1px 0 rgba(255, 255, 255, 0.15);
  border-color: color-mix(in srgb, var(--debug-color) 60%, white);
}

.level-debug:hover {
  transform: translateY(-1px);
  box-shadow:
    0 4px 7px rgba(var(--debug-color-rgb), 0.35),
    inset 0 1px 0 rgba(255, 255, 255, 0.25);
}

.log-message {
  color: var(--text-primary);
  word-break: break-all;
  line-height: 1.5;
  font-weight: 400;
  letter-spacing: 0.1px;
}

/* 日志行根据级别有不同的消息颜色 */
.log-line-log-line-error .log-message {
  color: var(--danger-color);
  font-weight: 500;
}

.log-line-log-line-warn .log-message {
  color: var(--warning-color);
  font-weight: 450;
}

.log-line-log-line-info .log-message {
  color: var(--info-color);
}

.log-line-log-line-debug .log-message {
  color: var(--debug-color);
  opacity: 0.85;
}

/* 日志响应式布局 - 移动设备优化 */
@media (max-width: 768px) {
  .log-line {
    grid-template-columns: auto 1fr;
    gap: 0.5rem;
  }

  .log-timestamp {
    font-size: 0.65rem;
  }

  .log-level {
    min-width: 40px;
    font-size: 0.6rem;
    padding: 0.15rem 0.45rem;
  }

  .log-message {
    grid-column: 1 / -1;
    word-break: break-word;
    overflow-wrap: break-word;
    white-space: pre-wrap;
  }
}

/* 日志加载状态 */
.log-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem 2rem;
  min-height: 150px;
}

.loading-spinner {
  width: 2.5rem;
  height: 2.5rem;
  border: 3px solid var(--border-color);
  border-top-color: var(--accent-color);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

.log-loading p {
  margin-top: 1rem;
  font-size: 0.85rem;
  color: var(--text-secondary);
  opacity: 0.7;
}

/* 日志底部 */
.log-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem 1rem;
  background: linear-gradient(to top, var(--bg-secondary), var(--bg-primary));
  border-top: 1px solid var(--border-color);
  font-size: 0.8rem;
  position: relative;
  gap: 1rem;
}

.log-footer::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(90deg,
    transparent 0%,
    var(--border-color) 50%,
    transparent 100%);
}

.log-stats {
  display: flex;
  gap: 1.5rem;
  align-items: center;
}

.stat-item {
  display: flex;
  gap: 0.4rem;
  align-items: center;
  padding: 0.25rem 0.6rem;
  background: var(--bg-primary);
  border-radius: 5px;
  border: 1px solid var(--border-color);
  transition: all 0.3s;
}

.stat-item:hover {
  border-color: var(--accent-color);
  transform: translateY(-1px);
  box-shadow: 0 1px 4px var(--shadow-color);
}

.stat-label {
  color: var(--text-secondary);
  font-size: 0.75rem;
  font-weight: 500;
}

.stat-value {
  color: var(--text-primary);
  font-weight: 700;
  font-size: 0.8rem;
}

.stat-error {
  border-color: var(--danger-color);
  background: rgba(var(--danger-color-rgb, 239, 68, 68), 0.05);
}

.stat-error .stat-value {
  color: var(--danger-color);
}

.stat-error:hover {
  background: rgba(var(--danger-color-rgb, 239, 68, 68), 0.1);
  box-shadow: 0 1px 4px rgba(var(--danger-color-rgb, 239, 68, 68), 0.3);
}

.stat-warn {
  border-color: var(--warning-color);
  background: rgba(var(--warning-color-rgb, 245, 158, 11), 0.05);
}

.stat-warn .stat-value {
  color: var(--warning-color);
}

.stat-warn:hover {
  background: rgba(var(--warning-color-rgb, 245, 158, 11), 0.1);
  box-shadow: 0 1px 4px rgba(var(--warning-color-rgb, 245, 158, 11), 0.3);
}

.log-actions-right {
  display: flex;
  align-items: center;
  gap: 0;
}

.btn-refresh,
.btn-clear {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  padding: 0;
  margin: 0;
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 5px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.3s;
  position: relative;
  overflow: hidden;
}

.btn-refresh + .btn-clear {
  margin-left: 0.5rem;
}

.btn-refresh::before,
.btn-clear::before {
  content: '';
  position: absolute;
  inset: 0;
  background: var(--accent-color);
  opacity: 0;
  transition: opacity 0.3s;
}

.btn-refresh:hover::before,
.btn-clear:hover::before {
  opacity: 0.1;
}

.btn-refresh:hover,
.btn-clear:hover {
  border-color: var(--accent-color);
  color: var(--accent-color);
  transform: translateY(-1px);
  box-shadow: 0 2px 6px var(--shadow-color);
}

.btn-refresh svg,
.btn-clear svg {
  width: 0.9rem;
  height: 0.9rem;
  position: relative;
  z-index: 1;
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 4rem 2rem;
  color: var(--text-secondary);
}

/* 全局加载状态 */
.global-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  min-height: 300px;
}

.global-loading .loading-spinner {
  width: 3rem;
  height: 3rem;
  border: 3px solid var(--border-color);
  border-top-color: var(--accent-color);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

.global-loading p {
  margin-top: 1rem;
  font-size: 0.9rem;
  color: var(--text-secondary);
  opacity: 0.7;
}

.empty-state svg {
  width: 6rem;
  height: 6rem;
  margin: 0 auto 1.5rem;
  opacity: 0.5;
}

.empty-state h2 {
  font-size: 1.5rem;
  margin: 0 0 0.5rem 0;
  color: var(--text-primary);
}

.empty-state p {
  font-size: 1rem;
  margin: 0;
}

/* 模态框样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.modal-card {
  background: var(--card-bg);
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  width: 90%;
  max-width: 520px;
  max-height: 90vh;
  overflow: hidden;
  border: 1px solid var(--border-color);
  animation: slideUp 0.3s ease;
}

@keyframes slideUp {
  from {
    transform: translateY(20px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.5rem 1.75rem;
  border-bottom: 1px solid var(--border-color);
  background: linear-gradient(135deg, var(--header-bg), var(--bg-primary));
}

.modal-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.modal-close {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  padding: 0;
  background: transparent;
  border: none;
  border-radius: 8px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.3s;
}

.modal-close:hover {
  background: var(--danger-color-bg);
  color: var(--danger-color);
  transform: rotate(90deg);
}

.modal-close svg {
  width: 1.25rem;
  height: 1.25rem;
}

.modal-body {
  padding: 1.75rem;
  overflow-y: auto;
  max-height: calc(90vh - 120px);
}

.modal-body::-webkit-scrollbar {
  width: 6px;
}

.modal-body::-webkit-scrollbar-track {
  background: var(--bg-primary);
}

.modal-body::-webkit-scrollbar-thumb {
  background: var(--border-color);
  border-radius: 3px;
}

.modal-body::-webkit-scrollbar-thumb:hover {
  background: var(--text-secondary);
}

.server-form {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr;
  gap: 1rem;
}

.form-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-primary);
}

.form-label svg {
  width: 1rem;
  height: 1rem;
  color: var(--accent-color);
}

.form-input {
  padding: 0.75rem 1rem;
  background: var(--bg-primary);
  border: 1.5px solid var(--border-color);
  border-radius: 8px;
  font-size: 0.9rem;
  color: var(--text-primary);
  transition: all 0.3s;
  outline: none;
}

.form-input::placeholder {
  color: var(--text-secondary);
  opacity: 0.6;
}

.form-input:hover {
  border-color: var(--text-secondary);
}

.form-input:focus {
  border-color: var(--accent-color);
  box-shadow: 0 0 0 3px var(--accent-color-bg);
}

.form-input.input-error {
  border-color: var(--danger-color);
}

.form-input.input-error:focus {
  box-shadow: 0 0 0 3px var(--danger-color-bg);
}

.form-error {
  font-size: 0.75rem;
  color: var(--danger-color);
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.form-error::before {
  content: '!';
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 14px;
  height: 14px;
  background: var(--danger-color);
  color: white;
  border-radius: 50%;
  font-size: 0.65rem;
  font-weight: 700;
}

.form-actions {
  display: flex;
  gap: 0.75rem;
  margin-top: 0.5rem;
  padding-top: 1.5rem;
  border-top: 1px solid var(--border-color);
}

/* 删除确认模态框样式 */
.modal-delete-confirm {
  max-width: 440px;
}

.delete-confirm-content {
  text-align: center;
  padding: 1rem 0;
}

.delete-warning-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 4rem;
  height: 4rem;
  background: var(--danger-color-bg);
  border-radius: 50%;
  color: var(--danger-color);
  margin-bottom: 1rem;
}

.delete-warning-icon svg {
  width: 2.5rem;
  height: 2.5rem;
}

.delete-confirm-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 0.75rem 0;
}

.delete-confirm-message {
  font-size: 0.95rem;
  color: var(--text-primary);
  margin: 0 0 0.5rem 0;
  line-height: 1.5;
}

.delete-confirm-message strong {
  font-weight: 600;
  color: var(--danger-color);
}

.delete-confirm-warning {
  font-size: 0.85rem;
  color: var(--text-secondary);
  margin: 0;
  padding: 0.75rem;
  background: var(--warning-color-bg);
  border-radius: 8px;
  border-left: 3px solid var(--warning-color);
}

/* 任务列表样式 */
.task-list {
  list-style: none;
  padding: 0;
  margin: 1rem 0;
  max-height: 200px;
  overflow-y: auto;
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
}

.task-list::-webkit-scrollbar {
  width: 6px;
}

.task-list::-webkit-scrollbar-track {
  background: var(--bg-primary);
}

.task-list::-webkit-scrollbar-thumb {
  background: var(--border-color);
  border-radius: 3px;
}

.task-list::-webkit-scrollbar-thumb:hover {
  background: var(--text-secondary);
}

.task-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.6rem 0.75rem;
  border-bottom: 1px solid var(--border-color);
  font-size: 0.9rem;
  color: var(--text-primary);
  transition: all 0.2s;
}

.task-item:last-child {
  border-bottom: none;
}

.task-item:hover {
  background: var(--bg-secondary);
}

.task-item svg {
  width: 1rem;
  height: 1rem;
  color: var(--accent-color);
  flex-shrink: 0;
}

@media (max-width: 768px) {
  .home {
    padding: 1rem;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .page-title {
    font-size: 1.5rem;
  }

  .server-meta {
    flex-direction: column;
    gap: 0.5rem;
  }

  .log-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.75rem;
  }

  .log-footer {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.75rem;
  }

  .log-content {
    max-height: 300px;
  }

  .log-line {
    flex-wrap: wrap;
  }

  .modal-card {
    width: 95%;
    max-width: none;
  }

  .modal-header {
    padding: 1.25rem 1.5rem;
  }

  .modal-body {
    padding: 1.25rem 1.5rem;
  }

  .form-actions {
    flex-direction: column;
  }
}
</style>
