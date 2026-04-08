<template>
  <div class="task-form-page">
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>{{ loadingText || t('common.loading') }}</p>
    </div>

    <div v-else>
      <div class="page-header">
        <AppButton class="btn-back" preserve-style @click="handleCancel">
          <template #icon>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="19" y1="12" x2="5" y2="12"></line>
              <polyline points="12 19 5 12 12 5"></polyline>
            </svg>
          </template>
          {{ t('common.back') }}
        </AppButton>
        <h1 class="page-title">{{ title }}</h1>
        <AppButton
          type="submit"
          class="btn-primary page-header-submit"
          preserve-style
          :disabled="submitting"
          @click="handleSubmit"
        >
          <template #icon>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"></path>
              <polyline points="17 21 17 13 7 13 7 21"></polyline>
              <polyline points="7 3 7 8 15 8"></polyline>
            </svg>
          </template>
          {{ submitting ? submittingText : submitText }}
        </AppButton>
      </div>

      <form class="task-form" @submit.prevent="handleSubmit">
        <!-- 基本信息 -->
        <section class="form-section">
          <h2 class="section-title">{{ t('taskForm.basicInfo') }}</h2>

          <div class="form-group">
            <label for="name">{{ t('taskForm.name') }}</label>
            <input
              id="name"
              v-model="form.name"
              type="text"
              required
              :placeholder="t('taskForm.placeholders.name')"
            />
          </div>

          <div class="form-group">
            <label for="description">{{ t('taskForm.description') }}</label>
            <textarea
              id="description"
              v-model="form.description"
              rows="3"
              :placeholder="t('taskForm.placeholders.description')"
            ></textarea>
          </div>
        </section>

        <!-- 服务器配置 -->
        <section class="form-section">
          <h2 class="section-title">{{ t('taskForm.serverConfig') }}</h2>

          <div class="form-group">
            <label for="serverId">{{ t('taskForm.selectServer') }}</label>
            <select
              id="serverId"
              v-model="form.serverId"
              required
              :disabled="loadingServers"
              @change="handleServerChange"
            >
              <option value="" disabled>{{ t('taskForm.selectServerPlaceholder') }}</option>
              <option
                v-for="server in availableServers"
                :key="server.id"
                :value="server.id"
              >
                {{ server.name }} ({{ server.address }})
              </option>
            </select>
            <div v-if="loadingServers" class="loading-servers">{{ t('taskForm.loadingServers') }}</div>
            <div v-if="!loadingServers && availableServers.length === 0" class="no-servers">
              {{ t('taskForm.noServersPrefix') }}
              <router-link to="/">{{ t('taskForm.noServersLink') }}</router-link>
              {{ t('taskForm.noServersSuffix') }}
            </div>
          </div>

          <div v-if="selectedServer" class="server-info">
            <div class="info-item">
              <span class="info-label">{{ t('taskForm.serverName') }}</span>
              <span class="info-value">{{ selectedServer.name }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">{{ t('taskForm.serverAddress') }}</span>
              <span class="info-value">{{ selectedServer.address }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">{{ t('taskForm.serverStatus') }}</span>
              <span
                class="status-badge"
                :class="{
                  'status-online': selectedServer.status === 'online',
                  'status-offline': selectedServer.status === 'offline',
                  'status-no-task': selectedServer.status === 'no_task',
                  'status-fault': selectedServer.status === 'fault',
                  'status-suspected-abnormal': selectedServer.status === 'suspected_abnormal'
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
            <h2 class="section-title">{{ t('taskForm.frpcConfig') }}</h2>
            <AppButton type="button" class="btn-add" preserve-style @click="addProxy">
              <template #icon>
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="12" y1="5" x2="12" y2="19"></line>
                  <line x1="5" y1="12" x2="19" y2="12"></line>
                </svg>
              </template>
              {{ t('taskForm.addProxy') }}
            </AppButton>
          </div>

          <div v-if="form.proxies.length === 0" class="empty-proxies">
            <p>{{ t('taskForm.emptyProxies') }}</p>
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

      </form>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, watch, onMounted } from 'vue'
import { serverAPI } from '@/api/server'
import AppButton from '@/components/AppButton.vue'
import ProxyConfigCard from '@/components/ProxyConfigCard.vue'
import { useI18n } from '@/utils/i18n'

/**
 * 创建空的代理配置对象。
 * 这里单独抽成工厂函数，避免创建页与编辑页对代理默认值的定义再次分散。
 *
 * @returns {object} 返回用于表单编辑的默认代理对象
 */
const createEmptyProxy = () => ({
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

/**
 * 创建空的任务表单对象。
 * 该函数统一了通用表单的初始结构，方便后续在创建态和编辑态之间复用。
 *
 * @returns {object} 返回任务表单的默认值
 */
const createEmptyForm = () => ({
  name: '',
  description: '',
  serverId: '',
  serverAddr: '',
  serverPort: 7000,
  authToken: '',
  proxies: []
})

const props = defineProps({
  title: {
    type: String,
    required: true
  },
  submitText: {
    type: String,
    required: true
  },
  submittingText: {
    type: String,
    required: true
  },
  loadingText: {
    type: String,
    default: ''
  },
  initialForm: {
    type: Object,
    default: () => ({})
  },
  loading: {
    type: Boolean,
    default: false
  },
  submitting: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['submit', 'cancel'])
const { t } = useI18n()

const loadingServers = ref(false)
const availableServers = ref([])
const selectedServer = ref(null)
const form = reactive(createEmptyForm())

/**
 * 将外部传入的初始表单数据归一化为组件内部结构。
 * 这里会补齐默认字段，并确保 proxies 永远是数组，降低模板与提交逻辑的分支复杂度。
 *
 * @param {object} source - 外部传入的初始表单数据
 * @returns {object} 返回可直接写入响应式表单的标准结构
 */
const normalizeFormData = (source = {}) => ({
  ...createEmptyForm(),
  ...source,
  proxies: Array.isArray(source.proxies) ? source.proxies.map(proxy => ({
    ...proxy,
    errors: reactive(proxy.errors || {})
  })) : []
})

/**
 * 根据当前表单中的 serverId 更新选中的服务器详情。
 * 当用户从下拉框切换服务器时，会同步写入 serverAddr 与 serverPort，
 * 确保提交时使用的是和界面一致的服务器地址。
 *
 * @returns {void}
 */
const handleServerChange = () => {
  const server = availableServers.value.find(item => item.id === form.serverId)

  if (server) {
    selectedServer.value = server
    const [addr, port] = server.address.split(':')
    form.serverAddr = addr
    form.serverPort = parseInt(port) || 7000
    return
  }

  selectedServer.value = null
  form.serverAddr = ''
  form.serverPort = 7000
}

/**
 * 根据已有的 serverAddr 与 serverPort 反查服务器下拉项。
 * 这个逻辑主要服务编辑场景，因为编辑接口返回的是地址与端口，
 * 页面需要把它映射回 serverId 才能正确展示下拉框选中状态。
 *
 * @returns {void}
 */
const syncSelectedServerByAddress = () => {
  if (!form.serverAddr) {
    selectedServer.value = null
    return
  }

  const matchingServer = availableServers.value.find(server => {
    const [addr, port] = server.address.split(':')
    return addr === form.serverAddr && parseInt(port) === form.serverPort
  })

  if (matchingServer) {
    form.serverId = matchingServer.id
    selectedServer.value = matchingServer
    return
  }

  selectedServer.value = null
}

/**
 * 将外部初始值同步进当前响应式表单。
 * 这里采用字段级赋值而不是直接替换整个 reactive 对象，
 * 这样可以保留模板中对 form 的响应式引用。
 *
 * @param {object} source - 外部初始表单数据
 * @returns {void}
 */
const applyInitialForm = (source = {}) => {
  const normalized = normalizeFormData(source)

  form.name = normalized.name
  form.description = normalized.description
  form.serverId = normalized.serverId
  form.serverAddr = normalized.serverAddr
  form.serverPort = normalized.serverPort
  form.authToken = normalized.authToken
  form.proxies = normalized.proxies

  if (form.serverId) {
    handleServerChange()
  } else {
    syncSelectedServerByAddress()
  }
}

/**
 * 加载服务器列表，并在加载完成后尝试恢复当前表单的选中状态。
 * 创建页依赖它展示服务器下拉框，编辑页依赖它将旧任务地址映射回具体服务器。
 *
 * @returns {Promise<void>} 返回异步加载流程
 */
const loadServers = async () => {
  try {
    loadingServers.value = true
    const servers = await serverAPI.listServers()
    availableServers.value = servers
    applyInitialForm(props.initialForm)
  } catch (error) {
    console.error('加载服务器列表失败:', error)
    alert(`${t('taskForm.alerts.loadServersFailed')}: ${error.message}`)
  } finally {
    loadingServers.value = false
  }
}

/**
 * 添加一条新的代理配置。
 * 新增代理默认进入编辑态，方便用户直接填写配置内容。
 *
 * @returns {void}
 */
const addProxy = () => {
  form.proxies.push(createEmptyProxy())
}

/**
 * 删除指定索引的代理配置。
 *
 * @param {number} index - 要删除的代理项索引
 * @returns {void}
 */
const removeProxy = (index) => {
  form.proxies.splice(index, 1)
}

/**
 * 代理卡片内部已经处理了校验，此处保留事件钩子用于调试和未来扩展。
 *
 * @param {number} index - 被校验的代理项索引
 * @returns {void}
 */
const checkProxy = (index) => {
  console.log('Proxy checked:', index)
}

/**
 * 代理卡片内部已经处理了编辑态切换，此处保留事件钩子用于调试和未来扩展。
 *
 * @param {number} index - 被切换编辑态的代理项索引
 * @returns {void}
 */
const toggleEdit = (index) => {
  console.log('Toggle edit:', index)
}

/**
 * 将页面上的编辑态表单转换为接口所需的提交数据。
 * 主要负责端口数字化、域名字符串拆分等数据清洗工作，
 * 避免创建页与编辑页各自维护一套相同的转换逻辑。
 *
 * @returns {object} 返回可直接提交给父组件的任务数据
 */
const buildSubmitPayload = () => ({
  name: form.name,
  description: form.description,
  serverId: form.serverId,
  serverAddr: form.serverAddr,
  serverPort: parseInt(form.serverPort),
  authToken: form.authToken,
  proxies: form.proxies.map(proxy => ({
    ...proxy,
    localPort: parseInt(proxy.localPort),
    remotePort: proxy.remotePort ? parseInt(proxy.remotePort) : 0,
    customDomains: proxy.customDomains
      ? proxy.customDomains.split(',').map(domain => domain.trim()).filter(domain => domain)
      : []
  }))
})

/**
 * 统一处理表单提交前校验，并将清洗后的数据抛给父组件。
 * 父组件只负责“提交到哪里”和“提交成功后跳转到哪里”，
 * 通用组件只负责“表单是否合法以及如何组装请求数据”。
 *
 * @returns {void}
 */
const handleSubmit = () => {
  if (form.proxies.length === 0) {
    alert(t('taskForm.alerts.noProxy'))
    return
  }

  if (!form.serverId) {
    alert(t('taskForm.alerts.noServer'))
    return
  }

  emit('submit', buildSubmitPayload())
}

/**
 * 统一触发取消事件，由父组件决定返回行为。
 *
 * @returns {void}
 */
const handleCancel = () => {
  emit('cancel')
}

/**
 * 将服务器状态值转换为中文文案。
 *
 * @param {string} status - 服务器状态值
 * @returns {string} 返回状态对应的中文显示文本
 */
const getStatusText = (status) => {
  switch (status) {
    case 'online':
      return t('status.server.online')
    case 'offline':
      return t('status.server.offline')
    case 'no_task':
      return t('status.server.noTask')
    case 'fault':
      return t('status.server.fault')
    case 'suspected_abnormal':
      return t('status.server.suspectedAbnormal')
    default:
      return t('status.server.unknown')
  }
}

watch(
  () => props.initialForm,
  (value) => {
    applyInitialForm(value)
  },
  { deep: true, immediate: true }
)

onMounted(() => {
  loadServers()
})
</script>

<style scoped>
.task-form-page {
  max-width: 900px;
  margin: 0 auto;
  padding: 2rem;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  gap: 1rem;
  color: var(--text-secondary);
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
  border-radius: var(--radius-sm);
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
  border-radius: var(--radius-lg);
  padding: 1.5rem;
  box-shadow: var(--shadow-md);
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
  border-radius: var(--radius-sm);
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
  box-shadow: var(--shadow-md);
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
  border-radius: var(--radius-sm);
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
  border-radius: var(--radius-sm);
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
  border-radius: var(--radius-md);
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
  border-radius: var(--radius-pill);
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
  background: var(--neutral-soft-bg);
  color: var(--text-secondary);
}

.status-fault {
  background: var(--warning-color-bg);
  color: var(--warning-color);
}

.status-suspected-abnormal {
  background: var(--info-color-bg);
  color: var(--info-color);
}

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
  box-shadow: var(--focus-ring);
}

.form-group textarea {
  resize: vertical;
  min-height: 80px;
}

.empty-proxies {
  text-align: center;
  padding: 2rem;
  color: var(--text-secondary);
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
  border: 2px dashed var(--border-color);
}

.btn-primary,
.btn-secondary {
  padding: 0.75rem 2rem;
  border: none;
  border-radius: var(--radius-sm);
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
  box-shadow: var(--shadow-lg);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.page-header-submit {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  justify-self: end;
  padding: 0.5rem 1rem;
  font-size: inherit;
  line-height: inherit;
}

.page-header-submit svg {
  width: 1rem;
  height: 1rem;
}

.btn-secondary {
  background: var(--bg-secondary);
  color: var(--text-primary);
  border: 1px solid var(--border-color);
}

.btn-secondary:hover {
  background: var(--accent-color-bg);
}

@media (max-width: 1024px) {
  .task-form-page {
    padding: 1.5rem;
  }

  .page-title {
    font-size: 1.75rem;
  }
}

@media (max-width: 768px) {
  .task-form-page {
    padding: 0.75rem;
  }

  .page-header {
    padding: 0.75rem 0;
    margin-bottom: 1rem;
    gap: 0.75rem;
    grid-template-columns: 1fr;
  }

  .btn-back {
    padding: 0.5rem 0.75rem;
    font-size: 0.8125rem;
  }

  .btn-back svg {
    width: 0.875rem;
    height: 0.875rem;
  }

  .page-header-submit {
    width: 100%;
  }

  .page-title {
    font-size: 1.25rem;
    text-align: center;
  }

  .task-form {
    gap: 1rem;
  }

  .form-section {
    padding: 1rem;
    border-radius: var(--radius-md);
  }

  .section-title {
    font-size: 1rem;
    line-height: 2rem;
  }

  .section-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.75rem;
    margin-bottom: 1rem;
  }

  .btn-add {
    width: 100%;
    justify-content: center;
    padding: 0.625rem 1rem;
  }

  .form-group {
    margin-bottom: 1rem;
  }

  .form-group label {
    font-size: 0.8125rem;
    margin-bottom: 0.375rem;
  }

  .form-group input,
  .form-group textarea,
  .form-group select {
    padding: 0.625rem;
    font-size: 0.875rem;
  }

  .form-group textarea {
    min-height: 70px;
  }

  .empty-proxies {
    padding: 1.5rem 1rem;
    font-size: 0.875rem;
  }

  .btn-primary,
  .btn-secondary {
    width: 100%;
    padding: 0.75rem 1.5rem;
    font-size: 0.9375rem;
  }
}

@media (max-width: 375px) {
  .task-form-page {
    padding: 0.5rem;
  }

  .page-header {
    margin-bottom: 0.75rem;
  }

  .btn-back {
    padding: 0.375rem 0.625rem;
    font-size: 0.75rem;
  }

  .btn-back svg {
    width: 0.75rem;
    height: 0.75rem;
  }

  .page-header-submit {
    min-width: 0;
  }

  .page-title {
    font-size: 1.125rem;
  }

  .form-section {
    padding: 0.875rem;
  }

  .section-title {
    font-size: 0.9375rem;
  }

  .form-group label {
    font-size: 0.75rem;
  }

  .form-group input,
  .form-group textarea,
  .form-group select {
    padding: 0.5rem;
    font-size: 0.8125rem;
  }

  .btn-add {
    padding: 0.5rem 0.75rem;
    font-size: 0.8125rem;
  }

  .btn-add svg {
    width: 0.875rem;
    height: 0.875rem;
  }
}
</style>
