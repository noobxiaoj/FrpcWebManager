<template>
  <div class="task-detail">
    <!-- 加载状态 -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>加载中...</p>
    </div>

    <div v-else-if="task" class="detail-container">
      <!-- 主卡片 -->
      <section class="card main-card">
        <!-- 头部导航栏 -->
        <div class="card-header">
          <div class="header-left">
            <AppButton class="btn-back" preserve-style @click="goBack" title="返回列表">
              <template #icon>
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <line x1="19" y1="12" x2="5" y2="12"></line>
                  <polyline points="12 19 5 12 12 5"></polyline>
                </svg>
              </template>
            </AppButton>
          </div>

          <div class="header-center">
            <h1 class="page-title">{{ task.name }}</h1>
            <TaskStatusIndicator v-if="task" :status="task.status" />
          </div>

          <div class="header-actions">
            <AppButton
              v-if="task.status === 'stopped'"
              class="btn-action start"
              preserve-style
              @click="startTask"
              title="启动任务"
            >
              <template #icon>
                <svg viewBox="0 0 24 24" fill="currentColor">
                  <polygon points="5 3 19 12 5 21 5 3"></polygon>
                </svg>
              </template>
            </AppButton>

            <AppButton
              v-if="task.status === 'running'"
              class="btn-action stop"
              preserve-style
              @click="stopTask"
              title="停止任务"
            >
              <template #icon>
                <svg viewBox="0 0 24 24" fill="currentColor">
                  <rect x="6" y="4" width="4" height="16"></rect>
                  <rect x="14" y="4" width="4" height="16"></rect>
                </svg>
              </template>
            </AppButton>

            <AppButton
              v-if="task.status === 'running'"
              class="btn-action reload"
              preserve-style
              @click="reloadTask"
              title="重载配置"
            >
              <template #icon>
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="23 4 23 10 17 10"></polyline>
                  <polyline points="1 20 1 14 7 14"></polyline>
                  <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"></path>
                </svg>
              </template>
            </AppButton>

            <AppButton
              class="btn-action edit"
              preserve-style
              @click="editTask"
              title="编辑任务"
            >
              <template #icon>
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
                  <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
                </svg>
              </template>
            </AppButton>

            <AppButton
              class="btn-action delete"
              preserve-style
              @click="deleteTask"
              title="删除任务"
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

        <!-- 内容区域 -->
        <div class="card-content">
          <!-- 信息网格 -->
          <div class="info-grid">
            <div class="info-item">
              <div class="info-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="2" y="2" width="20" height="8" rx="2" ry="2"></rect>
                  <rect x="2" y="14" width="20" height="8" rx="2" ry="2"></rect>
                </svg>
              </div>
              <div class="info-content">
                <span class="info-label">服务器</span>
                <span class="info-value">{{ getServerDisplayText(task) }}</span>
              </div>
            </div>

            <div class="info-item">
              <div class="info-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
                  <polyline points="22 4 12 14.01 9 11.01"></polyline>
                </svg>
              </div>
              <div class="info-content">
                <span class="info-label">端口数量</span>
                <span class="info-value">{{ task.proxies?.length || 0 }} 个</span>
              </div>
            </div>

            <div class="info-item">
              <div class="info-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
                  <line x1="16" y1="2" x2="16" y2="6"></line>
                  <line x1="8" y1="2" x2="8" y2="6"></line>
                  <line x1="3" y1="10" x2="21" y2="10"></line>
                </svg>
              </div>
              <div class="info-content">
                <span class="info-label">创建时间</span>
                <span class="info-value">{{ formatDate(task.createdAt) }}</span>
              </div>
            </div>
          </div>

          <!-- 任务简介 -->
          <div class="task-description">
            <p>{{ task.description || '暂无描述' }}</p>
          </div>

          <!-- 端口配置 -->
          <div class="proxies-section">
            <h3 class="section-title">端口配置</h3>

            <div v-if="!task.proxies || task.proxies.length === 0" class="empty-state">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="12" y1="8" x2="12" y2="12"></line>
                <line x1="12" y1="16" x2="12.01" y2="16"></line>
              </svg>
              <p>暂无端口配置</p>
            </div>

            <div v-else>
              <div v-for="(proxy, index) in task.proxies" :key="index" class="proxy-card proxy-card-compact">
                <div class="proxy-header">
                  <h3 class="proxy-name">
                    <span class="proxy-type-badge" :class="`badge-${proxy.type}`">{{ proxy.type.toUpperCase() }}</span>
                    {{ proxy.name || '未命名配置' }}
                  </h3>

                  <!-- TCP/UDP 显示本地IP和远程端口 -->
                  <div v-if="proxy.type === 'tcp' || proxy.type === 'udp'" class="proxy-compact-info">
                    <div class="info-item">本地IP: <span class="info-value">{{ proxy.localIP }}:{{ proxy.localPort }}</span></div>
                    <div class="info-item">远程端口: <span class="info-value">{{ proxy.remotePort }}</span></div>
                  </div>

                  <!-- HTTP/HTTPS 显示本地IP和域名 -->
                  <div v-else class="proxy-compact-info proxy-compact-info-http">
                    <div class="info-item">本地IP: <span class="info-value">{{ proxy.localIP }}:{{ proxy.localPort }}</span></div>
                    <div class="info-item">域名: <span class="info-value">{{ proxy.customDomains && proxy.customDomains.length ? proxy.customDomains.join(', ') : (proxy.subdomain || '-') }}</span></div>
                  </div>
                </div>
              </div>
            </div>
          </div>

        </div>
      </section>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useTaskStore } from '@/stores/task'
import AppButton from '@/components/AppButton.vue'
import TaskStatusIndicator from '@/components/TaskStatusIndicator.vue'
import { getApiBaseUrl } from '@/utils/api'

const router = useRouter()
const route = useRoute()
const taskStore = useTaskStore()

const task = ref(null)
const loading = ref(true)

const taskId = computed(() => route.params.id)

watch(taskId, async (newId) => {
  if (newId) {
    await loadTask()
  }
})

onMounted(async () => {
  await loadSettings()
  await taskStore.fetchTasks()
  await taskStore.fetchServers()
  await loadTask()
  loading.value = false
})

// 加载设置
const loadSettings = async () => {
  try {
    const response = await fetch(`${getApiBaseUrl()}/api/settings`)
    if (response.ok) {
      const data = await response.json()
      if (data.data && data.data.showServerName !== undefined) {
        taskStore.showServerName = data.data.showServerName
      }
    }
  } catch (error) {
    console.error('加载设置失败:', error)
  }
}

const loadTask = async () => {
  try {
    task.value = await taskStore.fetchTask(taskId.value)
  } catch (err) {
    alert('加载任务失败: ' + err.message)
    router.push('/tasks')
  }
}

const goBack = () => {
  router.push('/tasks')
}

const startTask = async () => {
  try {
    await taskStore.startTask(taskId.value)
    await loadTask()
    await taskStore.fetchTasks()
  } catch (err) {
    alert('启动任务失败: ' + err.message)
  }
}

const stopTask = async () => {
  if (!confirm('确定要停止这个任务吗?')) return

  try {
    await taskStore.stopTask(taskId.value)
    await loadTask()
    await taskStore.fetchTasks()
  } catch (err) {
    alert('停止任务失败: ' + err.message)
  }
}

const reloadTask = async () => {
  try {
    await taskStore.reloadTask(taskId.value)
    alert('重载成功')
    await loadTask()
    await taskStore.fetchTasks()
  } catch (err) {
    alert('重载任务失败: ' + err.message)
  }
}

const deleteTask = async () => {
  if (!confirm('确定要删除这个任务吗?此操作不可恢复。')) return

  try {
    await taskStore.deleteTask(taskId.value)
    router.push('/tasks')
  } catch (err) {
    alert('删除任务失败: ' + err.message)
  }
}

const editTask = () => {
  router.push(`/tasks/${taskId.value}/edit`)
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })
}

const getServerDisplayText = (task) => {
  if (!task) return '-'
  if (taskStore.showServerName) {
    const serverName = taskStore.getServerNameByAddress(task.serverAddr, task.serverPort)
    return serverName || `${task.serverAddr}:${task.serverPort}`
  }
  return `${task.serverAddr}:${task.serverPort}`
}
</script>

<style scoped>
.task-detail {
  min-height: 100vh;
  background: var(--bg-primary);
}

/* 头部导航 */
.card-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.25rem 1.5rem;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
  position: sticky;
  top: 0;
  z-index: 10;
}

.btn-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2.5rem;
  height: 2.5rem;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  background: var(--bg-secondary);
  cursor: pointer;
  transition: all 0.2s;
  padding: 0;
}

.btn-icon:hover {
  background: var(--accent-color-bg);
  border-color: var(--accent-color);
}

.btn-icon svg {
  width: 1.25rem;
  height: 1.25rem;
  stroke: var(--text-primary);
}

.btn-icon:hover svg {
  stroke: var(--accent-color);
}

.header-center {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 1rem;
}

.page-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.btn-action {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2.25rem;
  height: 2.25rem;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  background: var(--bg-secondary);
  cursor: pointer;
  transition: all 0.2s;
  padding: 0;
}

.btn-action svg {
  width: 1rem;
  height: 1rem;
}

.btn-action:hover {
  transform: translateY(-1px);
}

.btn-action.start {
  border-color: var(--success-color);
  background: var(--success-color-bg);
}

.btn-action.start:hover {
  opacity: 0.8;
}

.btn-action.start svg {
  stroke: var(--success-color);
}

.btn-action.stop {
  border-color: var(--danger-color);
  background: var(--danger-color-bg);
}

.btn-action.stop:hover {
  opacity: 0.8;
}

.btn-action.stop svg {
  stroke: var(--danger-color);
}

.btn-action.reload {
  border-color: var(--accent-color);
  background: var(--accent-color-bg);
}

.btn-action.reload:hover {
  opacity: 0.8;
}

.btn-action.reload svg {
  stroke: var(--accent-color);
}

.btn-action.edit {
  border-color: var(--edit-color);
  background: var(--edit-color-bg);
}

.btn-action.edit:hover {
  opacity: 0.8;
}

.btn-action.edit svg {
  stroke: var(--edit-color);
}

.btn-action.edit:hover svg {
  stroke: var(--edit-color);
}

.btn-action.delete {
  border-color: var(--danger-color);
  background: var(--danger-color-bg);
}

.btn-action.delete:hover {
  opacity: 0.8;
  border-color: var(--danger-color);
}

.btn-action.delete svg {
  stroke: var(--danger-color);
}

.btn-action.delete:hover svg {
  stroke: var(--danger-color);
}

.btn-primary {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.625rem 1.25rem;
  background: var(--accent-color);
  color: white;
  border: none;
  border-radius: var(--radius-md);
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary:hover {
  opacity: 0.9;
  transform: translateY(-1px);
}

.btn-primary svg {
  width: 1rem;
  height: 1rem;
}

/* 加载状态 */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  gap: 1rem;
  color: var(--text-secondary);
  min-height: 50vh;
}

.detail-container {
  padding: 1.5rem 2rem;
  max-width: 1400px;
  margin: 0 auto;
}

.card-content {
  padding: 1.5rem;
}

/* 头部区域 */
.header-left {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.btn-back {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  color: var(--text-primary);
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-back:hover {
  background: var(--accent-color-bg);
  border-color: var(--accent-color);
  color: var(--accent-color);
}

.btn-back svg {
  width: 1rem;
  height: 1rem;
}

.header-center {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 1rem;
  min-width: 0;
}

.page-title {
  font-size: 1.375rem;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.section-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 1.25rem 0;
}

.spinner {
  width: 3rem;
  height: 3rem;
  border: 3px solid var(--border-color);
  border-top-color: var(--accent-color);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* 主布局 */
.detail-layout {
  display: grid;
  grid-template-columns: 1fr 320px;
  gap: 1.5rem;
  padding: 1.5rem 2rem;
  max-width: 1600px;
  margin: 0 auto;
}

.main-content {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

/* 卡片样式 */
.card {
  background: var(--card-bg);
  border-radius: var(--radius-lg);
  border: 1px solid var(--border-color);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
}

.main-card {
  box-shadow: var(--shadow-lg);
}

.btn-icon-small {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  background: var(--bg-primary);
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
}

.btn-icon-small:hover {
  background: var(--accent-color-bg);
  border-color: var(--accent-color);
  color: var(--accent-color);
}

.btn-icon-small svg {
  width: 0.875rem;
  height: 0.875rem;
}

/* 任务信息卡片 */
.task-info-card {
  padding: 2rem;
}

.task-description {
  margin-bottom: 1.5rem;
  text-align: left;
}

.task-description p {
  font-size: 1rem;
  color: var(--text-primary);
  line-height: 1.6;
  margin: 0;
  text-align: left;
}

/* 信息网格 */
.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 0.75rem;
  margin-bottom: 1.5rem;
  align-items: stretch;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.875rem 1rem;
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
  transition: all 0.2s;
}

.info-item:hover {
  background: var(--accent-color-bg);
}

.info-icon {
  flex-shrink: 0;
  width: 2rem;
  height: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--accent-color-bg);
  border-radius: var(--radius-sm);
  margin-top: 0;
}

.info-icon svg {
  width: 1rem;
  height: 1rem;
  stroke: var(--accent-color);
}

.info-content {
  display: flex;
  flex-direction: column;
  gap: 0.375rem;
  flex: 1;
  min-width: 0;
}

.info-label {
  font-size: 0.6875rem;
  font-weight: 600;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.info-value {
  font-size: 0.875rem;
  color: var(--text-primary);
  font-weight: 500;
}

/* 端口配置 */
.proxies-section {
  margin-top: 1.5rem;
}

/* 卡片样式 */
.proxy-card {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  padding: 1.5rem;
  margin-bottom: 1rem;
  transition: all 0.3s;
  box-shadow: var(--shadow-md);
}

.proxy-card:hover {
  border-color: var(--accent-color);
  box-shadow: var(--shadow-lg);
}

/* 缩略模式 */
.proxy-card-compact {
  padding: 0.75rem 1rem;
}

.proxy-card-compact:hover {
  border-color: var(--accent-color);
  box-shadow: var(--shadow-lg);
}

.proxy-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
  flex-wrap: wrap;
  gap: 1rem;
}

.proxy-card-compact .proxy-header {
  margin-bottom: 0;
  gap: 1rem;
}

.proxy-name {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  flex-shrink: 0;
  min-width: 200px;
}

.proxy-name span:last-child {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 150px;
}

/* 协议徽章样式 */
.proxy-type-badge {
  font-size: 0.75rem;
  padding: 0.25rem 0.75rem;
  border-radius: var(--radius-xs);
  font-weight: 600;
}

/* TCP - 蓝色 */
.proxy-type-badge.badge-tcp {
  background: var(--protocol-tcp-bg);
  color: var(--protocol-tcp-color);
}

/* UDP - 红色 */
.proxy-type-badge.badge-udp {
  background: var(--protocol-udp-bg);
  color: var(--protocol-udp-color);
}

/* HTTP - 紫色 */
.proxy-type-badge.badge-http {
  background: var(--protocol-http-bg);
  color: var(--protocol-http-color);
}

/* HTTPS - 绿色 */
.proxy-type-badge.badge-https {
  background: var(--protocol-https-bg);
  color: var(--protocol-https-color);
}

/* TCPMUX - 粉色 */
.proxy-type-badge.badge-tcpmux {
  background: var(--protocol-tcpmux-bg);
  color: var(--protocol-tcpmux-color);
}

/* 缩略信息区域 */
.proxy-compact-info {
  display: flex;
  align-items: center;
  gap: 2.5rem;
  flex: 1;
  font-size: 0.875rem;
}

.proxy-compact-info .info-item {
  min-width: 200px;
  color: var(--text-secondary);
  font-weight: 500;
}

.proxy-compact-info .info-value {
  color: var(--text-primary);
  font-weight: 600;
}

/* 提示卡片 */
.tips-card {
  padding: 1.5rem;
}

.tips-content {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.tips-content p {
  font-size: 0.875rem;
  color: var(--text-secondary);
  line-height: 1.6;
  margin: 0;
}

.tips-content p::before {
  content: "💡 ";
  margin-right: 0.25rem;
}

/* 任务列表卡片 */
.tasks-card {
  padding: 1.5rem;
}

.task-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.task-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem 1rem;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.2s;
}

.task-item:hover {
  background: var(--accent-color-bg);
  border-color: var(--accent-color);
}

.task-item.active {
  background: var(--accent-color-bg);
  border-color: var(--accent-color);
}

.task-item.running {
  border-left: 3px solid var(--success-color);
}

.task-info {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  flex: 1;
  min-width: 0;
}

.task-name {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.task-server {
  font-size: 0.75rem;
  color: var(--text-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.task-status {
  padding: 0.25rem 0.5rem;
  border-radius: var(--radius-pill);
  font-size: 0.625rem;
  font-weight: 600;
  text-transform: uppercase;
  flex-shrink: 0;
  margin-left: 0.5rem;
}

/* 空状态 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem 2rem;
  gap: 1rem;
  color: var(--text-secondary);
}

.empty-state svg {
  width: 3rem;
  height: 3rem;
  opacity: 0.5;
}

.empty-state p {
  font-size: 0.875rem;
  margin: 0;
}

.empty-tasks {
  text-align: center;
  padding: 2rem;
  color: var(--text-secondary);
  font-size: 0.875rem;
}

.btn-small {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.375rem 0.75rem;
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  color: var(--text-primary);
  font-size: 0.75rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-small:hover {
  background: var(--accent-color-bg);
  border-color: var(--accent-color);
}

.btn-small svg {
  width: 0.875rem;
  height: 0.875rem;
}

/* 响应式设计 - 平板 */
@media (max-width: 1024px) {
  .detail-container {
    padding: 1rem;
  }

  .card-header {
    padding: 1rem;
  }

  .info-grid {
    grid-template-columns: 1fr;
  }
}

/* 响应式设计 - 手机 */
@media (max-width: 768px) {
  /* 头部导航 */
  .card-header {
    padding: 1rem;
    gap: 0.75rem;
    flex-wrap: nowrap;
    align-items: center;
  }

  .header-left {
    order: 1;
    flex-shrink: 0;
  }

  .header-center {
    order: 2;
    flex: 1;
    min-width: 0;
    justify-content: flex-start;
    gap: 0.5rem;
    align-items: center;
  }

  .header-actions {
    order: 3;
    flex-shrink: 0;
    flex-wrap: wrap;
    gap: 0.25rem;
    justify-content: flex-end;
    max-width: 140px;
    align-items: center;
  }

  .btn-back {
    padding: 0.4rem;
  }

  .btn-back span {
    display: none;
  }

  .btn-back svg {
    width: 1rem;
    height: 1rem;
  }

  .page-title {
    font-size: 1rem;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  /* 更紧凑的操作按钮 */
  .btn-action {
    width: 1.75rem;
    height: 1.75rem;
    flex-shrink: 0;
    border-radius: 4px;
  }

  .btn-action svg {
    width: 0.75rem;
    height: 0.75rem;
  }

  /* 内容区域 */
  .detail-container {
    padding: 0.5rem;
  }

  .card-content {
    padding: 0.75rem;
  }

  /* 信息网格 */
  .info-grid {
    grid-template-columns: 1fr;
    gap: 0.625rem;
  }

  .info-item {
    padding: 0.5rem 0.625rem;
  }

  .info-icon {
    width: 1.75rem;
    height: 1.75rem;
  }

  .info-icon svg {
    width: 0.875rem;
    height: 0.875rem;
  }

  .info-label {
    font-size: 0.625rem;
  }

  .info-value {
    font-size: 0.8125rem;
  }

  /* 任务描述 */
  .task-description {
    margin-bottom: 1rem;
  }

  .task-description p {
    font-size: 0.875rem;
  }

  /* 端口配置 */
  .proxies-section {
    margin-top: 1rem;
  }

  .section-title {
    font-size: 1rem;
    margin-bottom: 1rem;
  }

  /* 代理卡片 */
  .proxy-card-compact {
    padding: 0.875rem;
  }

  .proxy-header {
    gap: 0.75rem;
    flex-direction: column;
    align-items: flex-start;
    padding-left: 0;
    margin-left: 0;
  }

  .proxy-name {
    font-size: 0.875rem;
    min-width: auto;
    width: auto;
    margin: 0 0 0.5rem 0;
    padding: 0;
    display: flex;
    align-items: center;
  }

  .proxy-name span:last-child {
    max-width: calc(100vw - 180px);
  }

  .proxy-type-badge {
    font-size: 0.6875rem;
    padding: 0.1875rem 0.5rem;
    flex-shrink: 0;
  }

  .proxy-compact-info {
    gap: 0.5rem;
    font-size: 0.8125rem;
    flex-direction: column;
    flex-wrap: nowrap;
    align-items: flex-start;
    width: 100%;
    padding-left: 0;
    margin-left: 0;
  }

  .proxy-compact-info .info-item {
    min-width: auto;
    flex: none;
    width: 100%;
    display: flex;
    justify-content: flex-start;
    align-items: center;
    gap: 0.5rem;
    padding: 0.25rem 0;
    box-sizing: border-box;
    margin-left: 0;
  }

  /* HTTP 类型卡片特殊处理 */
  .proxy-compact-info-http {
    flex-direction: column;
    flex-wrap: nowrap;
    align-items: stretch;
  }

  /* 空状态 */
  .empty-state {
    padding: 2rem 1rem;
  }

  .empty-state svg {
    width: 2.5rem;
    height: 2.5rem;
  }

  .empty-state p {
    font-size: 0.8125rem;
  }
}

/* 超小屏幕优化 */
@media (max-width: 375px) {
  .card-header {
    padding: 0.75rem;
    gap: 0.5rem;
  }

  .page-title {
    font-size: 0.875rem;
  }

  .btn-action {
    width: 1.875rem;
    height: 1.875rem;
  }

  .btn-action svg {
    width: 0.75rem;
    height: 0.75rem;
  }

  .info-item {
    padding: 0.625rem;
  }

  .proxy-compact-info {
    font-size: 0.75rem;
    gap: 0.5rem;
  }

  .proxy-name span:last-child {
    max-width: 100px;
  }
}
</style>
