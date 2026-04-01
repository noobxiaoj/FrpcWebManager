<template>
  <div class="task-create">
    <div class="page-header">
      <AppButton class="btn-back" preserve-style @click="goBack">
        <template #icon>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="19" y1="12" x2="5" y2="12"></line>
            <polyline points="12 19 5 12 12 5"></polyline>
          </svg>
        </template>
        返回
      </AppButton>
      <h1 class="page-title">创建 FRPC 任务</h1>
      <div></div>
    </div>

    <form class="task-form" @submit.prevent="handleSubmit">
      <!-- 基本信息 -->
      <section class="form-section">
        <h2 class="section-title">基本信息</h2>

        <div class="form-group">
          <label for="name">任务名称 *</label>
          <input
            id="name"
            v-model="form.name"
            type="text"
            required
            placeholder="例如: 我的NAS穿透"
          />
        </div>

        <div class="form-group">
          <label for="description">任务描述</label>
          <textarea
            id="description"
            v-model="form.description"
            rows="3"
            placeholder="简要描述这个任务的用途..."
          ></textarea>
        </div>
      </section>

      <!-- 服务器配置 -->
      <section class="form-section">
        <h2 class="section-title">FRPS 服务器配置</h2>

        <div class="form-group">
          <label for="serverId">选择服务器 *</label>
          <select
            id="serverId"
            v-model="form.serverId"
            required
            :disabled="loadingServers"
            @change="handleServerChange"
          >
            <option value="" disabled>请选择FRPS服务器</option>
            <option
              v-for="server in availableServers"
              :key="server.id"
              :value="server.id"
            >
              {{ server.name }} ({{ server.address }})
            </option>
          </select>
          <div v-if="loadingServers" class="loading-servers">加载服务器列表中...</div>
          <div v-if="!loadingServers && availableServers.length === 0" class="no-servers">
            暂无可用服务器，请先在
            <router-link to="/">服务器页面</router-link>
            添加服务器
          </div>
        </div>

        <div v-if="selectedServer" class="server-info">
          <div class="info-item">
            <span class="info-label">服务器名称:</span>
            <span class="info-value">{{ selectedServer.name }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">服务器地址:</span>
            <span class="info-value">{{ selectedServer.address }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">状态:</span>
            <span
              class="status-badge"
              :class="{
                'status-online': selectedServer.status === 'online',
                'status-offline': selectedServer.status === 'offline',
                'status-no-task': selectedServer.status === 'no_task'
              }"
            >
              {{ getStatusText(selectedServer.status) }}
            </span>
          </div>
        </div>
      </section>

      <!-- Frpc配置 -->
      <section class="form-section">
        <div class="section-header">
          <h2 class="section-title">FRPC 配置</h2>
          <AppButton type="button" class="btn-add" preserve-style @click="addProxy">
            <template #icon>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="12" y1="5" x2="12" y2="19"></line>
                <line x1="5" y1="12" x2="19" y2="12"></line>
              </svg>
            </template>
            添加端口
          </AppButton>
        </div>

        <div v-if="form.proxies.length === 0" class="empty-proxies">
          <p>暂无Frpc配置,点击上方按钮添加</p>
        </div>

        <ProxyConfigCard
          v-for="(proxy, index) in form.proxies"
          :key="index"
          :proxy="proxy"
          :index="index"
          @check="checkProxy"
          @toggle-edit="toggleEdit"
          @remove="removeProxy"
        />
      </section>

      <!-- 提交按钮 -->
      <div class="form-actions">
        <AppButton type="button" class="btn-secondary" preserve-style @click="goBack">取消</AppButton>
        <AppButton type="submit" class="btn-primary" preserve-style :disabled="submitting">
          {{ submitting ? '创建中...' : '创建任务' }}
        </AppButton>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useTaskStore } from '@/stores/task'
import { serverAPI } from '@/api/server'
import AppButton from '@/components/AppButton.vue'
import ProxyConfigCard from '@/components/ProxyConfigCard.vue'

const router = useRouter()
const taskStore = useTaskStore()

const submitting = ref(false)
const loadingServers = ref(false)
const availableServers = ref([])
const selectedServer = ref(null)

const form = reactive({
  name: '',
  description: '',
  serverId: '',
  serverAddr: '',
  serverPort: 7000,
  authToken: '',
  proxies: []
})

// 加载服务器列表
const loadServers = async () => {
  try {
    loadingServers.value = true
    const servers = await serverAPI.listServers()
    availableServers.value = servers
  } catch (error) {
    console.error('加载服务器列表失败:', error)
    alert('加载服务器列表失败: ' + error.message)
  } finally {
    loadingServers.value = false
  }
}

// 处理服务器选择变化
const handleServerChange = () => {
  const server = availableServers.value.find(s => s.id === form.serverId)
  if (server) {
    selectedServer.value = server
    // 解析服务器地址和端口
    const [addr, port] = server.address.split(':')
    form.serverAddr = addr
    form.serverPort = parseInt(port) || 7000
    // 如果服务器有token，使用服务器的token（这里假设服务器对象中有token字段）
    // 注意：根据当前API响应，服务器对象可能不包含token，所以这里暂时不处理
  } else {
    selectedServer.value = null
    form.serverAddr = ''
    form.serverPort = 7000
  }
}

const goBack = () => {
  router.back()
}

const addProxy = () => {
  form.proxies.push(
    reactive({
      name: '',
      type: 'tcp',
      localIP: '127.0.0.1',
      localPort: '',
      remotePort: '',
      customDomains: '',
      subdomain: '',
      isEditing: true,
      isValid: false,
      errors: reactive({})
    })
  )
}

const removeProxy = (index) => {
  form.proxies.splice(index, 1)
}

const checkProxy = (index) => {
  // 验证逻辑已在 ProxyConfigCard 组件中处理
  console.log('Proxy checked:', index)
}

const toggleEdit = (index) => {
  // 编辑逻辑已在 ProxyConfigCard 组件中处理
  console.log('Toggle edit:', index)
}

const handleSubmit = async () => {
  if (form.proxies.length === 0) {
    alert('请至少添加一个映射配置')
    return
  }

  if (!form.serverId) {
    alert('请选择FRPS服务器')
    return
  }

  submitting.value = true

  try {
    // 处理代理数据
    const proxies = form.proxies.map(p => ({
      ...p,
      localPort: parseInt(p.localPort),
      remotePort: p.remotePort ? parseInt(p.remotePort) : 0,
      customDomains: p.customDomains ? p.customDomains.split(',').map(d => d.trim()).filter(d => d) : []
    }))

    const newTask = await taskStore.createTask({
      name: form.name,
      description: form.description,
      serverAddr: form.serverAddr,
      serverPort: parseInt(form.serverPort),
      authToken: form.authToken,
      proxies
    })

    // 跳转到新创建任务的详情页
    router.push(`/tasks/${newTask.id}`)
  } catch (err) {
    alert('创建任务失败: ' + err.message)
  } finally {
    submitting.value = false
  }
}

const getStatusText = (status) => {
  switch (status) {
    case 'online':
      return '在线'
    case 'offline':
      return '离线'
    case 'no_task':
      return '无任务'
    default:
      return '未知'
  }
}

onMounted(() => {
  loadServers()
})
</script>

<style scoped>
.task-create {
  max-width: 900px;
  margin: 0 auto;
  padding: 2rem;
}

.page-header {
  display: grid;
  grid-template-columns: auto 1fr auto;
  gap: 1rem;
  align-items: center;
  margin-bottom: 2rem;
}

.page-title {
  font-size: 2rem;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
  text-align: center;
}

.btn-back {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.3s;
}

.btn-back:hover {
  background: var(--accent-color-bg);
  color: var(--accent-color);
}

.btn-back svg {
  width: 1rem;
  height: 1rem;
}

.task-form {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.form-section {
  background: var(--card-bg);
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 2px 8px var(--shadow-color);
}

.section-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
  line-height: 2.25rem;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.btn-add {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  background: var(--success-color-bg);
  border: 1px solid var(--success-color);
  border-radius: 6px;
  color: var(--success-color);
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-add:hover {
  background: var(--success-color-bg);
  opacity: 0.8;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px var(--shadow-hover);
}

.btn-add svg {
  width: 1rem;
  height: 1rem;
}

.form-group {
  margin-bottom: 1.25rem;
}

.form-group label {
  display: block;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 0.5rem;
}

.form-group input,
.form-group textarea,
.form-group select {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 0.9rem;
  transition: all 0.3s;
}

.form-group select:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.loading-servers {
  margin-top: 0.5rem;
  font-size: 0.875rem;
  color: var(--text-secondary);
  font-style: italic;
}

.no-servers {
  margin-top: 0.5rem;
  font-size: 0.875rem;
  color: var(--warning-color);
  padding: 0.75rem;
  background: var(--warning-color-bg);
  border-radius: 6px;
  border-left: 3px solid var(--warning-color);
}

.no-servers a {
  color: var(--accent-color);
  text-decoration: none;
  font-weight: 600;
}

.no-servers a:hover {
  text-decoration: underline;
}

.server-info {
  margin-top: 1.25rem;
  padding: 1rem;
  background: var(--bg-secondary);
  border-radius: 8px;
  border: 1px solid var(--border-color);
}

.info-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.75rem;
}

.info-item:last-child {
  margin-bottom: 0;
}

.info-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-secondary);
  min-width: 80px;
}

.info-value {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-primary);
}

.status-badge {
  padding: 0.2rem 0.6rem;
  border-radius: 10px;
  font-size: 0.7rem;
  font-weight: 600;
  text-transform: uppercase;
}

.status-online {
  background: var(--success-color-bg);
  color: var(--success-color);
}

.status-offline {
  background: var(--danger-color-bg);
  color: var(--danger-color);
}

.status-no-task {
  background: var(--text-secondary-bg, rgba(148, 163, 184, 0.1));
  color: var(--text-secondary);
}

.no-servers-warning {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  padding: 1.5rem;
  background: var(--warning-color-bg);
  border: 2px solid var(--warning-color);
  border-radius: 12px;
  margin-top: 1rem;
  box-shadow: 0 2px 8px var(--shadow-color);
}

.no-servers-warning svg {
  width: 2.5rem;
  height: 2.5rem;
  color: var(--warning-color);
  flex-shrink: 0;
}

.warning-content h3 {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--warning-color);
  margin: 0 0 0.5rem 0;
}

.warning-content p {
  font-size: 0.9rem;
  color: var(--text-primary);
  margin: 0;
  line-height: 1.5;
}

.warning-content a {
  color: var(--accent-color);
  text-decoration: none;
  font-weight: 600;
}

.warning-content a:hover {
  text-decoration: underline;
}

/* 移除数字输入框的上下箭头 */
.form-group input[type="number"] {
  -moz-appearance: textfield;
}

.form-group input[type="number"]::-webkit-outer-spin-button,
.form-group input[type="number"]::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

.form-group input:focus,
.form-group textarea:focus,
.form-group select:focus {
  outline: none;
  border-color: var(--accent-color);
  box-shadow: 0 0 0 3px var(--accent-color-bg);
}

.form-group textarea {
  resize: vertical;
  min-height: 80px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.empty-proxies {
  text-align: center;
  padding: 2rem;
  color: var(--text-secondary);
  background: var(--bg-secondary);
  border-radius: 8px;
  border: 2px dashed var(--border-color);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  padding: 1.5rem;
  background: var(--card-bg);
  border-radius: 12px;
  box-shadow: 0 2px 8px var(--shadow-color);
}

.btn-primary,
.btn-secondary {
  padding: 0.75rem 2rem;
  border: none;
  border-radius: 6px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-primary {
  background: var(--accent-color);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.35);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-secondary {
  background: var(--bg-secondary);
  color: var(--text-primary);
  border: 1px solid var(--border-color);
}

.btn-secondary:hover {
  background: var(--accent-color-bg);
}

@media (max-width: 768px) {
  .task-create {
    padding: 1rem;
  }

  .page-header {
    grid-template-columns: 1fr;
    gap: 1rem;
  }

  .page-title {
    font-size: 1.5rem;
  }

  .form-row {
    grid-template-columns: 1fr;
  }

  .form-actions {
    flex-direction: column;
  }

  .form-actions button {
    width: 100%;
  }
}
</style>
