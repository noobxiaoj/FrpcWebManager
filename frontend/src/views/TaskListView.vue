<template>
  <div class="task-list">
    <div class="page-header">
      <h1 class="page-title">任务列表</h1>
      <div class="header-actions">
        <button
          class="btn-refresh"
          @click="handleRefresh"
          :disabled="taskStore.loading"
          title="刷新列表"
        >
          <svg
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            :class="{ spinning: taskStore.loading }"
          >
            <polyline points="23 4 23 10 17 10"></polyline>
            <polyline points="1 20 1 14 7 14"></polyline>
            <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"></path>
          </svg>
        </button>
        <button class="btn-primary" @click="goToCreate">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="5" x2="12" y2="19"></line>
            <line x1="5" y1="12" x2="19" y2="12"></line>
          </svg>
          创建任务
        </button>
      </div>
    </div>

    <!-- 错误提示 -->
    <div v-if="taskStore.error" class="error-message">
      {{ taskStore.error }}
      <button @click="taskStore.error = null">×</button>
    </div>

    <!-- 加载状态 -->
    <div v-if="taskStore.loading && tasks.length === 0" class="loading">
      加载中...
    </div>

    <!-- 任务列表 -->
    <div v-else-if="tasks.length > 0" class="task-grid">
      <div
        v-for="(task, index) in tasks"
        :key="task.id"
        class="task-card"
        :class="{
          'task-running': task.status === 'running',
          'dragging': isDragging && draggedIndex === index,
          'drag-over': dragOverIndex === index
        }"
        :draggable="true"
        @click="viewTask(task.id)"
        @dragstart="handleDragStart(index, $event)"
        @dragend="handleDragEnd"
        @dragover="handleDragOver(index, $event)"
        @dragleave="handleDragLeave"
        @drop="handleDrop(index, $event)"
      >
        <div class="task-card-header">
          <h3 class="task-name">{{ task.name }}</h3>
          <TaskStatusIndicator :status="task.status" />
        </div>

        <div class="task-info">
          <div class="info-row">
            <div class="info-item">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="2" y="2" width="20" height="8" rx="2" ry="2"></rect>
                <rect x="2" y="14" width="20" height="8" rx="2" ry="2"></rect>
                <line x1="6" y1="6" x2="6" y2="6"></line>
                <line x1="6" y1="18" x2="6" y2="18"></line>
              </svg>
              <span>{{ getDisplayText(task) }}</span>
            </div>

            <div class="info-item">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="2" y1="12" x2="22" y2="12"></line>
                <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path>
              </svg>
              <span>{{ task.proxies?.length || 0 }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else class="empty-state">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="12" r="10"></circle>
        <line x1="12" y1="8" x2="12" y2="12"></line>
        <line x1="12" y1="16" x2="12.01" y2="16"></line>
      </svg>
      <h2>暂无任务</h2>
      <p>点击上方按钮创建第一个 FRPC 任务</p>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, watch, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useTaskStore } from '@/stores/task'
import TaskStatusIndicator from '@/components/TaskStatusIndicator.vue'
import { getApiBaseUrl } from '@/utils/api'

const router = useRouter()
const route = useRoute()
const taskStore = useTaskStore()

const tasks = computed(() => taskStore.sortedTasks)

const handleRefresh = () => {
  taskStore.fetchTasks()
}

onMounted(async () => {
  await loadSettings()
  taskStore.loadTaskOrder()
  taskStore.fetchTasks()
  taskStore.fetchServers()
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

// 监听路由变化,每次进入页面时刷新数据
watch(
  () => route.path,
  (newPath) => {
    if (newPath === '/tasks') {
      taskStore.fetchTasks()
    }
  }
)

const goToCreate = () => {
  router.push('/tasks/create')
}

const viewTask = (id) => {
  router.push(`/tasks/${id}`)
}

const getDisplayText = (task) => {
  if (taskStore.showServerName) {
    const serverName = taskStore.getServerNameByAddress(task.serverAddr, task.serverPort)
    return serverName || `${task.serverAddr}:${task.serverPort}`
  }
  return `${task.serverAddr}:${task.serverPort}`
}

// 拖放相关
const draggedIndex = ref(null)
const dragOverIndex = ref(null)
const isDragging = ref(false)

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
  event.dataTransfer.setDragImage(dragElement, 150, 100)

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

const handleDrop = (dropIndex, event) => {
  event.preventDefault()
  const dragIndex = draggedIndex.value

  if (dragIndex === null || dragIndex === dropIndex) {
    dragOverIndex.value = null
    return
  }

  // 创建新的排序数组
  const newOrder = [...taskStore.taskOrder]
  const taskIds = tasks.value.map(t => t.id)

  // 如果还没有排序,先初始化
  if (newOrder.length === 0) {
    newOrder.push(...taskIds)
  }

  // 移动任务
  const [movedId] = newOrder.splice(dragIndex, 1)
  newOrder.splice(dropIndex, 0, movedId)

  // 更新排序
  taskStore.updateTaskOrder(newOrder)

  isDragging.value = false
  draggedIndex.value = null
  dragOverIndex.value = null
}
</script>

<style scoped>
.task-list {
  max-width: 1400px;
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

.btn-refresh {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2.5rem;
  height: 2.5rem;
  padding: 0;
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  color: var(--text-primary);
  cursor: pointer;
  transition: all 0.3s;
}

.btn-refresh:hover:not(:disabled) {
  background: var(--accent-color-bg);
  border-color: var(--accent-color);
  transform: rotate(90deg);
}

.btn-refresh:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-refresh svg {
  width: 1.25rem;
  height: 1.25rem;
  transition: transform 0.6s ease;
}

.btn-refresh svg.spinning {
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.btn-primary {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: var(--accent-color);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  box-shadow: 0 2px 8px var(--shadow-color);
}

.btn-primary:hover {
  opacity: 0.9;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px var(--shadow-hover);
}

.btn-primary svg {
  width: 1.25rem;
  height: 1.25rem;
}

.error-message {
  background: var(--danger-color-bg);
  border: 1px solid var(--danger-color);
  border-radius: 8px;
  padding: 1rem;
  margin-bottom: 1.5rem;
  color: var(--danger-color);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.error-message button {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: var(--danger-color);
  padding: 0;
  width: 1.5rem;
  height: 1.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.loading {
  text-align: center;
  padding: 3rem;
  color: var(--text-secondary);
  font-style: italic;
}

.task-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 1.5rem;
}

.task-card {
  position: relative;
  background: var(--card-bg);
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 2px 8px var(--shadow-color);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: pointer;
  border: 2px solid transparent;
}

.task-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 16px var(--shadow-color);
}

.task-card.task-running {
  border-color: var(--success-color);
}

.task-card.dragging {
  opacity: 0.4;
  transform: scale(0.95) rotate(2deg);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
  cursor: grabbing;
}

.task-card.drag-over {
  border-color: var(--accent-color);
  background: var(--accent-color-bg);
  transform: scale(1.02);
  box-shadow: 0 0 0 3px var(--accent-color-bg), 0 12px 24px var(--shadow-color);
}

.task-card[draggable="true"] {
  cursor: grab;
}

.task-card[draggable="true"]:active {
  cursor: grabbing;
}

.task-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.75rem;
}

.task-name {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.task-info {
  margin-top: 0.75rem;
}

.info-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  color: var(--text-secondary);
  font-size: 0.875rem;
}

.info-item span {
  line-height: 1;
  padding-top: 1px;
}

.info-item svg {
  width: 1rem;
  height: 1rem;
  flex-shrink: 0;
}

.empty-state {
  text-align: center;
  padding: 4rem 2rem;
  color: var(--text-secondary);
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

@media (max-width: 768px) {
  .task-list {
    padding: 1rem;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }

  .page-title {
    font-size: 1.5rem;
  }

  .task-grid {
    grid-template-columns: 1fr;
  }
}
</style>
