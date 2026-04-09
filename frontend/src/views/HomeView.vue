<template>
  <div class="home">
    <PageHeader :title="t('home.title')">
      <template #actions>
        <AppButton
          class="page-header-action-button page-header-action-button--icon"
          variant="icon"
          @click="handleRefresh"
          :disabled="loading"
          :loading="loading"
          :title="t('home.refreshTitle')"
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
          class="btn-primary page-header-action-button page-header-action-button--icon"
          preserve-style
          @click="openAddServerModal"
          :title="t('home.createTitle')"
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

    <div v-if="loading && sortedServers.length === 0" class="global-loading">
      <div class="loading-spinner"></div>
      <p>{{ t('home.loading') }}</p>
    </div>

    <div v-else-if="sortedServers.length > 0" class="server-grid">
      <article
        v-for="server in sortedServers"
        :key="server.id"
        class="server-card"
        :class="{ 'server-card-running': server.status === 'online' }"
        @click="viewServer(server.id)"
      >
        <ServerInfo :server="server" />
      </article>
    </div>

    <div v-else-if="isInitialized" class="empty-state">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <rect x="2" y="2" width="20" height="8" rx="2" ry="2"></rect>
        <rect x="2" y="14" width="20" height="8" rx="2" ry="2"></rect>
        <line x1="6" y1="6" x2="6" y2="6"></line>
        <line x1="6" y1="18" x2="6" y2="18"></line>
      </svg>
      <h2>{{ t('home.emptyTitle') }}</h2>
      <p>{{ t('home.emptyDescription') }}</p>
    </div>

    <ServerFormModal
      :visible="showAddServerModal"
      mode="create"
      :submitting="submittingServer"
      @close="closeAddServerModal"
      @submit="submitAddServer"
    />
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { taskApi } from '@/api'
import { serverAPI } from '@/api/server'
import AppButton from '@/components/AppButton.vue'
import PageHeader from '@/components/PageHeader.vue'
import ServerInfo from '@/components/ServerInfo.vue'
import ServerFormModal from '@/components/ServerFormModal.vue'
import { useI18n } from '@/utils/i18n'

const router = useRouter()
const { t } = useI18n()

const loading = ref(false)
const isInitialized = ref(false)
const servers = ref([])
const tasks = ref([])

const showAddServerModal = ref(false)
const submittingServer = ref(false)

/**
 * 服务器状态排序权重。
 * 需求顺序为：异常 > 在线 > 离线 > 无任务。
 * 当前系统里“异常”包含 fault 和 suspected_abnormal，“暂停”更接近离线语义，因此归入离线档位。
 *
 * @param {string} status - 服务器状态值
 * @returns {number} 排序权重，值越大越靠前
 */
const getServerStatusWeight = (status) => {
  switch (status) {
    case 'fault':
    case 'suspected_abnormal':
      return 4
    case 'online':
      return 3
    case 'offline':
    case 'paused':
      return 2
    case 'no_task':
      return 1
    default:
      return 0
  }
}

/**
 * 统计指定服务器关联的端口总数。
 * 首页服务器接口没有直接返回端口数量，因此这里通过任务列表中 proxies 的长度累加得到。
 *
 * @param {{address?: string}} server - 服务器对象
 * @returns {number} 该服务器关联的端口总数
 */
const getServerPortCount = (server) => {
  const serverAddress = server?.address || ''
  if (!serverAddress) {
    return 0
  }

  return tasks.value.reduce((total, task) => {
    const taskAddress = `${task.serverAddr || ''}:${task.serverPort || ''}`
    if (taskAddress !== serverAddress) {
      return total
    }

    return total + (Array.isArray(task.proxies) ? task.proxies.length : 0)
  }, 0)
}

/**
 * 首页服务器统一排序：
 * 1. 状态：异常 > 在线 > 离线 > 无任务
 * 2. 端口数量：从大到小
 * 3. 创建时间：从晚到早
 *
 * @returns {Array<object>} 排序后的服务器列表副本
 */
const sortedServers = computed(() => {
  return [...servers.value].sort((serverA, serverB) => {
    const statusWeightDiff = getServerStatusWeight(serverB.status) - getServerStatusWeight(serverA.status)
    if (statusWeightDiff !== 0) {
      return statusWeightDiff
    }

    const portCountDiff = getServerPortCount(serverB) - getServerPortCount(serverA)
    if (portCountDiff !== 0) {
      return portCountDiff
    }

    const createdAtA = new Date(serverA.createdAt || 0).getTime()
    const createdAtB = new Date(serverB.createdAt || 0).getTime()
    return createdAtB - createdAtA
  })
})

/**
 * 拉取首页卡片所需的服务器与任务数据。
 * 虽然首页只展示服务器卡片，但排序依赖“端口数量”，
 * 所以这里会额外拉取任务列表，用于统计每台服务器关联的端口总数。
 */
const loadServers = async () => {
  try {
    loading.value = true
    const [serverResult, taskResult] = await Promise.allSettled([
      serverAPI.listServers(),
      taskApi.list()
    ])

    if (serverResult.status === 'fulfilled') {
      servers.value = serverResult.value
    } else {
      throw serverResult.reason
    }

    // 任务列表仅用于辅助首页排序。
    // 如果这里失败，首页仍然展示服务器列表，只是端口数量排序临时回退为 0。
    tasks.value = taskResult.status === 'fulfilled'
      ? (taskResult.value?.data?.tasks || [])
      : []
  } catch (error) {
    console.error('加载服务器列表失败:', error)
  } finally {
    loading.value = false
    isInitialized.value = true
  }
}

const handleRefresh = async () => {
  await loadServers()
}

const viewServer = (id) => {
  router.push(`/servers/${id}`)
}

const openAddServerModal = () => {
  showAddServerModal.value = true
}

const closeAddServerModal = () => {
  showAddServerModal.value = false
}

/**
 * 提交新建服务器表单。
 * 弹窗内部已经完成本地校验，这里只负责调用接口并刷新首页列表。
 * @param {{name: string, address: string, port: string, token: string}} payload - 标准化后的服务器表单数据
 * @returns {Promise<void>} 无返回值
 */
const submitAddServer = async (payload) => {
  try {
    submittingServer.value = true
    await serverAPI.createServer(payload)
    closeAddServerModal()
    await loadServers()
  } catch (error) {
    console.error('创建服务器失败:', error)
    alert(`${t('home.messages.createFailed')}: ${error.message}`)
  } finally {
    submittingServer.value = false
  }
}

onMounted(async () => {
  await loadServers()
})
</script>

<style scoped>
.home {
  max-width: var(--page-content-width);
  margin: 0 auto;
  padding: 2rem;
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

.global-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 320px;
  gap: 1rem;
  color: var(--text-secondary);
}

.loading-spinner {
  width: 2rem;
  height: 2rem;
  border: 3px solid var(--border-color);
  border-top-color: var(--accent-color);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

.server-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 1.5rem;
}

.server-card {
  position: relative;
  background: var(--card-bg);
  border-radius: var(--radius-lg);
  padding: 1.5rem;
  box-shadow: var(--shadow-md);
  border: 2px solid transparent;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.server-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-xl);
}

.server-card-running {
  border-color: var(--success-color);
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 320px;
  text-align: center;
  color: var(--text-secondary);
  gap: 0.75rem;
}

.empty-state svg {
  width: 3rem;
  height: 3rem;
  color: var(--accent-color);
}

.empty-state h2 {
  margin: 0;
  color: var(--text-primary);
  font-size: 1.5rem;
}

.empty-state p {
  margin: 0;
}

@media (max-width: 768px) {
  .home {
    padding: 1.25rem;
  }
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

@media (max-width: 768px) {
  .home {
    padding: 1rem;
  }

  .server-grid {
    grid-template-columns: 1fr;
  }
}
</style>
