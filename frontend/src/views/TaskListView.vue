<template>
  <div class="task-list">
    <PageHeader :title="t('taskList.title')">
      <template #actions>
        <AppButton
          class="btn-refresh page-header-action-button page-header-action-button--icon"
          preserve-style
          @click="handleRefresh"
          :disabled="taskStore.loading"
          :title="t('taskList.refreshTitle')"
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
        </AppButton>
        <AppButton
          class="btn-primary page-header-action-button page-header-action-button--icon"
          preserve-style
          @click="goToCreate"
          :title="t('taskList.createTitle')"
        >
          <template #icon>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="12" y1="5" x2="12" y2="19"></line>
              <line x1="5" y1="12" x2="19" y2="12"></line>
            </svg>
          </template>
        </AppButton>
      </template>
    </PageHeader>

    <!-- 错误提示 -->
    <div v-if="taskStore.error" class="error-message">
      {{ taskStore.error }}
      <AppButton preserve-style @click="taskStore.error = null">×</AppButton>
    </div>

    <!-- 加载状态 -->
    <div v-if="taskStore.loading && tasks.length === 0" class="loading">
      {{ t('common.loading') }}
    </div>

    <!-- 任务列表 -->
    <div v-else-if="tasks.length > 0" class="task-grid">
      <div
        v-for="task in tasks"
        :key="task.id"
        class="task-card"
        :class="{ 'task-running': task.status === 'running' }"
        @click="viewTask(task.id)"
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
      <h2>{{ t('taskList.emptyTitle') }}</h2>
      <p>{{ t('taskList.emptyDescription') }}</p>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useTaskStore } from '@/stores/task'
import AppButton from '@/components/AppButton.vue'
import PageHeader from '@/components/PageHeader.vue'
import TaskStatusIndicator from '@/components/TaskStatusIndicator.vue'
import { getApiBaseUrl } from '@/utils/api'
import { useI18n } from '@/utils/i18n'

const router = useRouter()
const route = useRoute()
const taskStore = useTaskStore()
const { t } = useI18n()

const tasks = computed(() => taskStore.tasks)

const handleRefresh = () => {
  taskStore.fetchTasks()
}

onMounted(async () => {
  await loadSettings()
  taskStore.fetchTasks()
  taskStore.fetchServers()
})

// 加载设置
const loadSettings = async () => {
  try {
    const response = await fetch(`${getApiBaseUrl()}/api/settings`, {
      credentials: 'include'
    })
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

</script>

<style scoped>
.task-list {
  /* 任务列表页与首页、关于页、设置页共用统一内容宽度。 */
  max-width: var(--page-content-width);
  margin: 0 auto;
  padding: 2rem;
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
  border-radius: var(--radius-md);
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
  border-radius: var(--radius-md);
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  box-shadow: var(--shadow-md);
}

.btn-primary:hover {
  opacity: 0.9;
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
}

.btn-primary svg {
  width: 1.25rem;
  height: 1.25rem;
}

.error-message {
  background: var(--danger-color-bg);
  border: 1px solid var(--danger-color);
  border-radius: var(--radius-md);
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
  border-radius: var(--radius-lg);
  padding: 1.5rem;
  box-shadow: var(--shadow-md);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: pointer;
  border: 2px solid transparent;
}

.task-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-xl);
}

.task-card.task-running {
  border-color: var(--success-color);
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

  .task-grid {
    grid-template-columns: 1fr;
  }
}
</style>
