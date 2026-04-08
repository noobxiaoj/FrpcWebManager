<template>
  <div class="home">
    <PageHeader title="服务器列表">
      <template #actions>
        <AppButton
          class="page-header-action-button page-header-action-button--icon"
          variant="icon"
          @click="handleRefresh"
          :disabled="loading"
          :loading="loading"
          title="刷新列表"
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
          title="新建服务器"
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
      <p>加载服务器列表中...</p>
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
      <h2>暂无服务器</h2>
      <p>请先添加 FRPC 服务器</p>
    </div>

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

  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { serverAPI } from '@/api/server'
import AppButton from '@/components/AppButton.vue'
import PageHeader from '@/components/PageHeader.vue'
import ServerInfo from '@/components/ServerInfo.vue'

const router = useRouter()

const loading = ref(false)
const isInitialized = ref(false)
const servers = ref([])

const showAddServerModal = ref(false)
const newServer = ref({
  name: '',
  address: '',
  port: '',
  token: ''
})
const addServerFormErrors = ref({})

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

/**
 * 校验服务器创建表单，避免无效参数直接发往后端。
 * @returns {boolean} 是否通过校验
 */
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

  addServerFormErrors.value = errors
  return Object.keys(errors).length === 0
}

const submitAddServer = async () => {
  if (!validateServerForm()) return

  try {
    await serverAPI.createServer({
      name: newServer.value.name,
      address: newServer.value.address,
      port: newServer.value.port,
      token: newServer.value.token
    })

    closeAddServerModal()
    await loadServers()
  } catch (error) {
    console.error('创建服务器失败:', error)
    alert('创建服务器失败: ' + error.message)
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

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1.25rem;
  z-index: 1000;
}

.modal-card {
  width: min(100%, 560px);
  max-height: calc(100vh - 2.5rem);
  overflow: auto;
  background: var(--card-bg);
  border-radius: var(--radius-lg);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-overlay);
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 1.25rem 1.5rem;
  border-bottom: 1px solid var(--border-color);
}

.modal-title {
  margin: 0;
  font-size: 1.15rem;
  color: var(--text-primary);
}

.modal-close {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  border-radius: var(--radius-sm);
  color: var(--text-secondary);
}

.modal-close:hover {
  background: var(--bg-secondary);
  color: var(--text-primary);
}

.modal-body {
  padding: 1.5rem;
}

.server-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
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
  font-size: 0.95rem;
  font-weight: 600;
  color: var(--text-primary);
}

.form-label svg {
  width: 1rem;
  height: 1rem;
  color: var(--accent-color);
}

.form-input {
  width: 100%;
  padding: 0.8rem 0.9rem;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 0.95rem;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.form-input:focus {
  outline: none;
  border-color: var(--accent-color);
  box-shadow: var(--focus-ring);
}

.input-error {
  border-color: var(--danger-color);
}

.form-error {
  font-size: 0.85rem;
  color: var(--danger-color);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 0.5rem;
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

  .form-actions {
    flex-direction: column;
  }
}
</style>
