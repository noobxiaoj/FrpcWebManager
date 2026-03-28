<template>
  <div class="settings">
    <div class="page-header">
      <h1 class="page-title">系统设置</h1>
      <p class="page-description">配置 frpc 管理器的各项参数</p>
    </div>

    <div class="settings-container">
      <div class="setting-section">
        <h2 class="section-title">服务器显示设置</h2>

        <div class="setting-item">
          <div class="setting-info">
            <label class="setting-label" for="showServerPort">
              显示进程端口
            </label>
            <p class="setting-description">
              控制是否在服务器页面显示 frpc webServer 端口信息
            </p>
          </div>
          <div class="setting-control">
            <label class="switch">
              <input
                id="showServerPort"
                type="checkbox"
                v-model="settings.showServerPort"
                @change="handleToggleChange"
                :disabled="saving"
              />
              <span class="slider"></span>
            </label>
          </div>
        </div>

        <div class="setting-item">
          <div class="setting-info">
            <label class="setting-label" for="showServerName">
              任务显示服务器名称
            </label>
            <p class="setting-description">
              控制在任务列表和任务详情中显示服务器名称而非IP地址
            </p>
          </div>
          <div class="setting-control">
            <label class="switch">
              <input
                id="showServerName"
                type="checkbox"
                v-model="settings.showServerName"
                @change="handleToggleChange"
                :disabled="saving"
              />
              <span class="slider"></span>
            </label>
          </div>
        </div>
      </div>

      <div class="setting-section">
        <h2 class="section-title">数据刷新设置</h2>

        <div class="setting-item">
          <div class="setting-info">
            <label class="setting-label" for="refreshInterval">
              刷新间隔
            </label>
            <p class="setting-description">
              日志和运行时间的自动刷新间隔，"不刷新"表示手动刷新
            </p>
          </div>
          <div class="setting-control">
            <select
              id="refreshInterval"
              v-model.number="settings.refreshInterval"
              @change="handleIntervalChange"
              :disabled="saving"
              class="interval-select"
            >
              <option :value="0">不刷新</option>
              <option :value="1">1 秒</option>
              <option :value="3">3 秒</option>
              <option :value="10">10 秒</option>
              <option :value="30">30 秒</option>
              <option :value="60">1 分钟</option>
              <option :value="120">2 分钟</option>
            </select>
          </div>
        </div>

        <div class="setting-item">
          <div class="setting-info">
            <label class="setting-label" for="showRefreshTime">
              显示刷新时间
            </label>
            <p class="setting-description">
              控制是否在服务器页面显示距离上次刷新的时间
            </p>
          </div>
          <div class="setting-control">
            <label class="switch">
              <input
                id="showRefreshTime"
                type="checkbox"
                v-model="settings.showRefreshTime"
                @change="handleToggleChange"
                :disabled="saving"
              />
              <span class="slider"></span>
            </label>
          </div>
        </div>
      </div>

      <div class="setting-section">
        <h2 class="section-title">前端服务设置</h2>

        <div class="setting-item">
          <div class="setting-info">
            <label class="setting-label" for="frontendPort">
              前端服务端口
            </label>
            <p class="setting-description">
              前端开发服务器监听的端口号（范围：1024-65535）
            </p>
            <p class="setting-warning">
              修改端口后需要重启前端服务才能生效
            </p>
          </div>
          <div class="setting-control">
            <input
              id="frontendPort"
              type="number"
              v-model.number="settings.frontendPort"
              @change="handlePortChange"
              :disabled="saving"
              class="port-input"
              min="1024"
              max="65535"
              step="1"
            />
          </div>
        </div>
      </div>

      <div class="setting-section">
        <h2 class="section-title">访问控制设置</h2>

        <div class="setting-item">
          <div class="setting-info">
            <label class="setting-label" for="enableIPWhitelist">
              启用IP白名单
            </label>
            <p class="setting-description">
              开启后，只有白名单中的IP地址可以访问系统
            </p>
            <p class="setting-warning">
              白名单设置修改后需要重启容器才能生效
            </p>
          </div>
          <div class="setting-control">
            <label class="switch">
              <input
                id="enableIPWhitelist"
                type="checkbox"
                v-model="settings.enableIPWhitelist"
                @change="handleToggleChange"
                :disabled="saving"
              />
              <span class="slider"></span>
            </label>
          </div>
        </div>

        <!-- IP白名单管理 -->
        <div v-if="settings.enableIPWhitelist" class="whitelist-manager">
          <div class="setting-item">
            <div class="setting-info full-width">
              <label class="setting-label">IP白名单</label>
              <p class="setting-description">
                添加允许访问系统的IP地址或CIDR网段（例如：192.168.1.100 或 192.168.1.0/24）
              </p>
              <p class="setting-description">
                修改白名单后请重启容器使设置生效
              </p>

              <div class="whitelist-input-group">
                <input
                  ref="ipInput"
                  type="text"
                  v-model="newIP"
                  @keyup.enter="addIP"
                  placeholder="输入IP地址或CIDR（如：192.168.1.100）"
                  class="ip-input"
                  :disabled="saving"
                />
                <button
                  @click="addIP"
                  :disabled="saving || !newIP.trim()"
                  class="btn-add-ip"
                >
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <line x1="12" y1="5" x2="12" y2="19"></line>
                    <line x1="5" y1="12" x2="19" y2="12"></line>
                  </svg>
                  添加
                </button>
              </div>

              <div v-if="ipError" class="ip-error">
                {{ ipError }}
              </div>

              <div v-if="settings.ipWhitelist && settings.ipWhitelist.length > 0" class="whitelist-items">
                <div
                  v-for="(ip, index) in settings.ipWhitelist"
                  :key="index"
                  class="whitelist-item"
                >
                  <span class="ip-address">{{ ip }}</span>
                  <button
                    @click="removeIP(index)"
                    class="btn-remove-ip"
                    title="删除"
                  >
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <line x1="18" y1="6" x2="6" y2="18"></line>
                      <line x1="6" y1="6" x2="18" y2="18"></line>
                    </svg>
                  </button>
                </div>
              </div>

              <div v-else class="whitelist-empty">
                白名单为空，请添加IP地址
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-if="saveMessage" class="save-message" :class="{ 'success': saveSuccess, 'error': !saveSuccess }">
      {{ saveMessage }}
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useTaskStore } from '@/stores/task'
import { getApiBaseUrl } from '@/utils/api'

const taskStore = useTaskStore()

const settings = ref({
  showServerPort: true,
  refreshInterval: 10,
  showRefreshTime: true,
  showServerName: false,
  frontendPort: 4500,
  enableIPWhitelist: false,
  ipWhitelist: []
})

const saving = ref(false)
const saveMessage = ref('')
const saveSuccess = ref(false)
const newIP = ref('')
const ipError = ref('')

// 加载设置
const loadSettings = async () => {
  try {
    const response = await fetch(`${getApiBaseUrl()}/api/settings`)
    if (!response.ok) {
      throw new Error('获取设置失败')
    }
    const data = await response.json()
    if (data.data) {
      settings.value = data.data
      // 同步到 taskStore
      taskStore.showServerName = data.data.showServerName || false
    }
  } catch (error) {
    console.error('加载设置失败:', error)
    showSaveMessage('加载设置失败', false)
  }
}

// 处理开关变化
const handleToggleChange = async () => {
  // 同步 showServerName 到 taskStore
  taskStore.showServerName = settings.value.showServerName
  await saveSettings()
}

// 处理刷新间隔变化
const handleIntervalChange = async () => {
  await saveSettings()
}

// 端口验证函数
const validatePort = (port) => {
  const portNum = Number(port)
  if (!Number.isInteger(portNum) || portNum < 1024 || portNum > 65535) {
    return false
  }
  return true
}

// 处理端口变化
const handlePortChange = async () => {
  if (!validatePort(settings.value.frontendPort)) {
    showSaveMessage('端口号必须在 1024-65535 之间', false)
    // 重新加载设置恢复原值
    await loadSettings()
    return
  }

  await saveSettings()
  if (saveSuccess.value) {
    showSaveMessage('端口已保存，请重启前端服务以生效', true)
  }
}

// 保存设置
const saveSettings = async () => {
  saving.value = true
  saveMessage.value = ''

  try {
    const response = await fetch(`${getApiBaseUrl()}/api/settings`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(settings.value)
    })

    if (!response.ok) {
      throw new Error('保存设置失败')
    }

    return true
  } catch (error) {
    console.error('保存设置失败:', error)
    showSaveMessage('保存设置失败', false)
    return false
  } finally {
    saving.value = false
  }
}

// 显示保存消息
const showSaveMessage = (message, success) => {
  saveMessage.value = message
  saveSuccess.value = success
  setTimeout(() => {
    saveMessage.value = ''
  }, 3000)
}

// 验证IP地址或CIDR格式
const validateIP = (ip) => {
  const trimmed = ip.trim()

  // 检查CIDR格式
  if (trimmed.includes('/')) {
    const parts = trimmed.split('/')
    if (parts.length !== 2) return false

    const ipPart = parts[0]
    const maskPart = parts[1]

    // 验证IP部分
    const ipRegex = /^(\d{1,3}\.){3}\d{1,3}$/
    if (!ipRegex.test(ipPart)) return false

    // 验证IP段范围
    const segments = ipPart.split('.')
    for (const seg of segments) {
      const num = parseInt(seg, 10)
      if (num < 0 || num > 255) return false
    }

    // 验证子网掩码
    const mask = parseInt(maskPart, 10)
    if (isNaN(mask) || mask < 0 || mask > 32) return false

    return true
  }

  // 检查单个IP格式
  const ipRegex = /^(\d{1,3}\.){3}\d{1,3}$/
  if (!ipRegex.test(trimmed)) return false

  // 验证IP段范围
  const segments = trimmed.split('.')
  for (const seg of segments) {
    const num = parseInt(seg, 10)
    if (num < 0 || num > 255) return false
  }

  return true
}

// 添加IP到白名单
const addIP = async () => {
  const trimmed = newIP.value.trim()

  if (!trimmed) {
    ipError.value = '请输入IP地址'
    return
  }

  if (!validateIP(trimmed)) {
    ipError.value = 'IP地址格式无效，请输入有效的IP地址或CIDR（例如：192.168.1.100 或 192.168.1.0/24）'
    return
  }

  // 检查是否已存在
  if (settings.value.ipWhitelist && settings.value.ipWhitelist.includes(trimmed)) {
    ipError.value = '该IP地址已在白名单中'
    return
  }

  // 添加到白名单
  if (!settings.value.ipWhitelist) {
    settings.value.ipWhitelist = []
  }
  settings.value.ipWhitelist.push(trimmed)

  // 清空输入框和错误信息
  newIP.value = ''
  ipError.value = ''

  // 保存设置
  const success = await saveSettings()
  if (success) {
    showSaveMessage('IP已添加到白名单', true)
  }
}

// 从白名单删除IP
const removeIP = async (index) => {
  if (settings.value.ipWhitelist && settings.value.ipWhitelist.length > 0) {
    const removed = settings.value.ipWhitelist.splice(index, 1)

    // 如果删除了所有IP，自动关闭白名单
    if (settings.value.ipWhitelist.length === 0) {
      settings.value.enableIPWhitelist = false
    }

    // 保存设置
    const success = await saveSettings()
    if (success) {
      showSaveMessage(`已从白名单删除: ${removed}`, true)
    }
  }
}

onMounted(() => {
  loadSettings()
})
</script>

<style scoped>
.settings {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;
}

.page-header {
  margin-bottom: 2rem;
}

.page-title {
  font-size: 2rem;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 0.5rem 0;
}

.page-description {
  font-size: 1rem;
  color: var(--text-secondary);
  margin: 0;
}

.settings-container {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.setting-section {
  background: var(--card-bg);
  border-radius: 16px;
  padding: 1.5rem;
  border: 1px solid var(--border-color);
  box-shadow: 0 4px 12px var(--shadow-color);
}

.section-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 1.5rem 0;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid var(--border-color);
}

.setting-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 2rem;
  padding: 1rem 0;
}

.setting-info {
  flex: 1;
}

.setting-label {
  display: block;
  font-size: 1rem;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 0.5rem;
}

.setting-description {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin: 0;
  line-height: 1.5;
}

.setting-control {
  flex-shrink: 0;
}

/* 下拉选择框样式 */
.interval-select {
  min-width: 150px;
  padding: 0.5rem 2rem 0.5rem 0.75rem;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 1rem;
  background: var(--input-bg);
  color: var(--text-primary);
  cursor: pointer;
  transition: border-color 0.3s, box-shadow 0.3s;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 12 12'%3E%3Cpath fill='%23666' d='M6 9L1 4h10z'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 0.75rem center;
  padding-right: 2rem;
}

.interval-select:focus {
  outline: none;
  border-color: var(--accent-color);
  box-shadow: 0 0 0 3px var(--accent-color-bg);
}

.interval-select:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.interval-select:hover {
  border-color: var(--accent-color);
}

/* 开关样式 */
.switch {
  position: relative;
  display: inline-block;
  width: 50px;
  height: 26px;
}

.switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: var(--border-color);
  transition: 0.3s;
  border-radius: 26px;
}

.slider:before {
  position: absolute;
  content: "";
  height: 20px;
  width: 20px;
  left: 3px;
  bottom: 3px;
  background-color: white;
  transition: 0.3s;
  border-radius: 50%;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

input:checked + .slider {
  background-color: var(--accent-color);
}

input:checked + .slider:before {
  transform: translateX(24px);
}

input:disabled + .slider {
  opacity: 0.5;
  cursor: not-allowed;
}

.slider:hover {
  box-shadow: 0 0 0 3px var(--accent-color-bg);
}

/* 保存消息 */
.save-message {
  position: fixed;
  bottom: 2rem;
  right: 2rem;
  padding: 1rem 1.5rem;
  border-radius: 8px;
  font-size: 0.875rem;
  font-weight: 500;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  animation: slideIn 0.3s ease;
  z-index: 1000;
}

@keyframes slideIn {
  from {
    transform: translateY(100px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.save-message.success {
  background-color: var(--success-color);
  color: white;
}

.save-message.error {
  background-color: var(--danger-color);
  color: white;
}

/* 端口输入框样式 */
.port-input {
  width: 120px;
  padding: 0.5rem 0.75rem;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 1rem;
  background: var(--input-bg);
  color: var(--text-primary);
  text-align: center;
  transition: border-color 0.3s, box-shadow 0.3s;
}

.port-input:focus {
  outline: none;
  border-color: var(--accent-color);
  box-shadow: 0 0 0 3px var(--accent-color-bg);
}

.port-input:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.port-input:hover {
  border-color: var(--accent-color);
}

/* 设置警告文字 */
.setting-warning {
  font-size: 0.85rem;
  color: #78716c;  /* 更柔和的灰色提示 */
  margin: 0.5rem 0 0 0;
  font-weight: 400;
  opacity: 0.8;
}

/* 白名单管理器 */
.whitelist-manager {
  margin-top: 1rem;
  padding-top: 1rem;
  border-top: 1px dashed var(--border-color);
}

.setting-info.full-width {
  flex: 1;
  width: 100%;
}

/* IP输入框组 */
.whitelist-input-group {
  display: flex;
  gap: 0.5rem;
  margin-top: 1rem;
}

.ip-input {
  flex: 1;
  padding: 0.6rem 1rem;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 0.9rem;
  background: var(--input-bg);
  color: var(--text-primary);
  transition: border-color 0.3s, box-shadow 0.3s;
  font-family: 'SF Mono', 'Monaco', 'Consolas', monospace;
}

.ip-input:focus {
  outline: none;
  border-color: var(--accent-color);
  box-shadow: 0 0 0 3px var(--accent-color-bg);
}

.ip-input:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-add-ip {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.6rem 1.25rem;
  background: var(--accent-color);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  white-space: nowrap;
}

.btn-add-ip:hover:not(:disabled) {
  opacity: 0.9;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px var(--shadow-color);
}

.btn-add-ip:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-add-ip svg {
  width: 1rem;
  height: 1rem;
}

/* IP错误提示 */
.ip-error {
  margin-top: 0.5rem;
  padding: 0.5rem 0.75rem;
  background: var(--danger-color-bg);
  color: var(--danger-color);
  border-radius: 6px;
  font-size: 0.85rem;
  border-left: 3px solid var(--danger-color);
}

/* 白名单项容器 */
.whitelist-items {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  margin-top: 1rem;
  max-height: 300px;
  overflow-y: auto;
}

.whitelist-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem 1rem;
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  transition: all 0.2s;
}

.whitelist-item:hover {
  border-color: var(--accent-color);
  background: var(--accent-color-bg);
}

.ip-address {
  font-family: 'SF Mono', 'Monaco', 'Consolas', monospace;
  font-size: 0.9rem;
  color: var(--text-primary);
  font-weight: 500;
}

.btn-remove-ip {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  padding: 0;
  background: transparent;
  border: none;
  border-radius: 6px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
}

.btn-remove-ip:hover {
  background: var(--danger-color-bg);
  color: var(--danger-color);
  transform: scale(1.1);
}

.btn-remove-ip svg {
  width: 0.9rem;
  height: 0.9rem;
}

/* 空状态 */
.whitelist-empty {
  margin-top: 1rem;
  padding: 2rem;
  text-align: center;
  color: var(--text-secondary);
  font-size: 0.9rem;
  background: var(--bg-primary);
  border: 1px dashed var(--border-color);
  border-radius: 8px;
}

@media (max-width: 768px) {
  .settings {
    padding: 1rem;
  }

  .setting-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }

  .setting-control {
    align-self: flex-end;
  }

  .save-message {
    right: 1rem;
    left: 1rem;
    bottom: 1rem;
  }

  .whitelist-input-group {
    flex-direction: column;
  }

  .btn-add-ip {
    width: 100%;
    justify-content: center;
  }

  .ip-input {
    width: 100%;
  }
}
</style>
