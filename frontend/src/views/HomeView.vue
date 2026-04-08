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

    <div v-if="loading && servers.length === 0" class="global-loading">
      <div class="loading-spinner"></div>
      <p>{{ t('home.loading') }}</p>
    </div>

    <div v-else-if="servers.length > 0" class="server-grid">
      <article
        v-for="server in servers"
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
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
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

const showAddServerModal = ref(false)
const submittingServer = ref(false)

/**
 * 拉取首页卡片所需的服务器与任务数据。
 * 当前首页仅展示服务器列表，因此只需要刷新服务器数据即可。
 */
const loadServers = async () => {
  try {
    loading.value = true
    servers.value = await serverAPI.listServers()
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
