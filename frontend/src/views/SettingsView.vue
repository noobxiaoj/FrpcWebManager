<template>
  <div class="settings">
    <PageHeader :title="t('settings.pageTitle')">
      <template #actions>
        <AppButton
          @click="saveSettings"
          :disabled="saving || !hasChanges"
          class="btn-save page-header-action-button"
          preserve-style
        >
          <svg v-if="saving" class="spinner" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10" stroke-dasharray="31.4" stroke-dashoffset="10"></circle>
          </svg>
          <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"></path>
            <polyline points="17 21 17 13 7 13 7 21"></polyline>
            <polyline points="7 3 7 8 15 8"></polyline>
          </svg>
          {{ saving ? t('common.saving') : t('common.saveSettings') }}
        </AppButton>
        <AppButton
          v-if="hasChanges"
          @click="resetSettings"
          :disabled="saving"
          class="btn-reset page-header-action-button"
          preserve-style
        >
          {{ t('common.reset') }}
        </AppButton>
      </template>
    </PageHeader>

    <div class="settings-container">
      <SettingSection
        :title="t('settings.sections.serverDisplay')"
        section-id="server-display-settings-content"
      >
        <SettingItem
          id="showServerPort"
          :label="t('settings.fields.showServerPort')"
        >
          <Switch id="showServerPort" v-model="settings.showServerPort" :disabled="saving" />
        </SettingItem>

        <SettingItem
          id="showServerName"
          :label="t('settings.fields.showServerName')"
        >
          <Switch id="showServerName" v-model="settings.showServerName" :disabled="saving" />
        </SettingItem>

        <SettingItem
          id="showRefreshTime"
          :label="t('settings.fields.showRefreshTime')"
        >
          <Switch id="showRefreshTime" v-model="settings.showRefreshTime" :disabled="saving" />
        </SettingItem>
      </SettingSection>

      <SettingSection
        :title="t('settings.sections.dataRefresh')"
        section-id="data-refresh-settings-content"
      >
        <SettingItem
          id="refreshInterval"
          :label="t('settings.fields.refreshInterval')"
        >
          <select
            id="refreshInterval"
            v-model.number="settings.refreshInterval"
            :disabled="saving"
            class="interval-select"
          >
            <option :value="0">{{ t('settings.refreshOptions.noRefresh') }}</option>
            <option :value="1">{{ t('settings.refreshOptions.seconds', { count: 1 }) }}</option>
            <option :value="3">{{ t('settings.refreshOptions.seconds', { count: 3 }) }}</option>
            <option :value="10">{{ t('settings.refreshOptions.seconds', { count: 10 }) }}</option>
            <option :value="30">{{ t('settings.refreshOptions.seconds', { count: 30 }) }}</option>
            <option :value="60">{{ t('settings.refreshOptions.oneMinute') }}</option>
            <option :value="120">{{ t('settings.refreshOptions.minutes', { count: 2 }) }}</option>
          </select>
        </SettingItem>
      </SettingSection>

      <SettingSection
        :title="t('settings.sections.frontendService')"
        section-id="frontend-service-settings-content"
      >
        <SettingItem
          id="frontendPort"
          :label="t('settings.fields.frontendPort')"
          :warning="t('settings.warnings.frontendPort')"
        >
          <input
            id="frontendPort"
            type="number"
            v-model.number="settings.frontendPort"
            :disabled="saving"
            class="port-input"
            min="1024"
            max="65535"
            step="1"
          />
        </SettingItem>
      </SettingSection>

      <SettingSection
        :title="t('settings.sections.accessControl')"
        section-id="access-control-settings-content"
      >
        <SettingItem
          id="enableIPWhitelist"
          :label="t('settings.fields.enableIPWhitelist')"
          :warning="t('settings.warnings.ipWhitelist')"
        >
          <Switch id="enableIPWhitelist" v-model="settings.enableIPWhitelist" :disabled="saving" />
        </SettingItem>

        <!-- IP白名单管理 -->
        <div v-if="settings.enableIPWhitelist" class="whitelist-manager">
          <div class="setting-item">
            <div class="setting-info full-width">
              <label class="setting-label">{{ t('settings.fields.ipWhitelist') }}</label>

              <div class="whitelist-input-group">
                <input
                  ref="ipInput"
                  type="text"
                  v-model="newIP"
                  @keyup.enter="addIP"
                  :placeholder="t('settings.ip.placeholder')"
                  class="ip-input"
                  :disabled="saving"
                />
                <AppButton
                  @click="addIP"
                  :disabled="saving || !newIP.trim()"
                  class="btn-add-ip"
                  preserve-style
                >
                  <template #icon>
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <line x1="12" y1="5" x2="12" y2="19"></line>
                      <line x1="5" y1="12" x2="19" y2="12"></line>
                    </svg>
                  </template>
                  {{ t('settings.ip.add') }}
                </AppButton>
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
                  <AppButton
                    @click="removeIP(index)"
                    class="btn-remove-ip"
                    :title="t('settings.ip.remove')"
                    preserve-style
                  >
                    <template #icon>
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <line x1="18" y1="6" x2="6" y2="18"></line>
                        <line x1="6" y1="6" x2="18" y2="18"></line>
                      </svg>
                    </template>
                  </AppButton>
                </div>
              </div>

              <div v-else class="whitelist-empty">
                {{ t('settings.ip.empty') }}
              </div>
            </div>
          </div>
        </div>

        <SettingItem id="passwordSetting" :label="t('settings.fields.passwordSetting')">
          <div class="password-actions">
            <AppButton
              v-if="!settings.passwordAuth.enabled"
              class="password-action-button"
              @click="openPasswordModal('add')"
              :disabled="saving || passwordModalSaving"
            >
              {{ t('passwordModal.addPassword') }}
            </AppButton>

            <template v-else>
              <AppButton
                variant="secondary"
                class="password-action-button"
                @click="openPasswordModal('change')"
                :disabled="saving || passwordModalSaving"
              >
                {{ t('passwordModal.changePassword') }}
              </AppButton>
              <AppButton
                variant="danger"
                class="password-action-button password-action-button--danger"
                @click="openPasswordModal('delete')"
                :disabled="saving || passwordModalSaving"
              >
                {{ t('passwordModal.deletePassword') }}
              </AppButton>
            </template>
          </div>
        </SettingItem>
      </SettingSection>
    </div>

    <div v-if="saveMessage" class="save-message" :class="{ 'success': saveSuccess, 'error': !saveSuccess }">
      {{ saveMessage }}
    </div>

    <PasswordSettingsModal
      :visible="passwordModalVisible"
      :mode="passwordModalMode"
      :saving="passwordModalSaving"
      :current-username="settings.passwordAuth.username"
      @close="closePasswordModal"
      @submit="handlePasswordSubmit"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useTaskStore } from '@/stores/task'
import { getApiBaseUrl } from '@/utils/api'
import { setLanguage, syncLanguageFromSettings, useI18n } from '@/utils/i18n'
import AppButton from '@/components/AppButton.vue'
import PageHeader from '@/components/PageHeader.vue'
import SettingItem from '@/components/SettingItem.vue'
import Switch from '@/components/AppSwitch.vue'
import SettingSection from '@/components/SettingSection.vue'
import PasswordSettingsModal from '@/components/PasswordSettingsModal.vue'

const taskStore = useTaskStore()
const { t } = useI18n()

const settings = ref({
  showServerPort: true,
  refreshInterval: 10,
  showRefreshTime: true,
  showServerName: false,
  language: 'zh-CN',
  frontendPort: 4500,
  enableIPWhitelist: false,
  ipWhitelist: [],
  passwordAuth: {
    enabled: false,
    username: ''
  }
})

const saving = ref(false)
const saveMessage = ref('')
const saveSuccess = ref(false)
const newIP = ref('')
const ipError = ref('')
const originalSettings = ref(null)
const passwordModalVisible = ref(false)
const passwordModalMode = ref('add')
const passwordModalSaving = ref(false)

// 加载设置
const loadSettings = async () => {
  try {
    const response = await fetch(`${getApiBaseUrl()}/api/settings`, {
      credentials: 'include'
    })
    if (!response.ok) {
      throw new Error('获取设置失败')
    }
    const data = await response.json()
    if (data.data) {
      const normalizedSettings = {
        ...JSON.parse(JSON.stringify(data.data)),
        language: data.data.language || 'zh-CN',
        passwordAuth: {
          enabled: data.data.passwordAuth?.enabled || false,
          username: data.data.passwordAuth?.username || ''
        }
      }
      settings.value = normalizedSettings
      originalSettings.value = JSON.parse(JSON.stringify(normalizedSettings))
      syncLanguageFromSettings(normalizedSettings.language)
      // 同步到 taskStore
      taskStore.showServerName = data.data.showServerName || false
    }
  } catch (error) {
    console.error('加载设置失败:', error)
    showSaveMessage(t('settings.messages.loadFailed'), false)
  }
}

// 检测是否有修改
const hasChanges = computed(() => {
  if (!originalSettings.value) return false
  return JSON.stringify(settings.value) !== JSON.stringify(originalSettings.value)
})

// 重置设置到原始状态
const resetSettings = () => {
  if (originalSettings.value) {
    settings.value = JSON.parse(JSON.stringify(originalSettings.value))
  }
}

// 打开密码设置弹窗
const openPasswordModal = (mode) => {
  passwordModalMode.value = mode
  passwordModalVisible.value = true
}

// 关闭密码设置弹窗
const closePasswordModal = () => {
  if (passwordModalSaving.value) return
  passwordModalVisible.value = false
}

// 端口验证函数
const validatePort = (port) => {
  const portNum = Number(port)
  if (!Number.isInteger(portNum) || portNum < 1024 || portNum > 65535) {
    return false
  }
  return true
}

// 保存设置
const saveSettings = async () => {
  // 端口验证
  if (!validatePort(settings.value.frontendPort)) {
    showSaveMessage(t('settings.messages.invalidPort'), false)
    return
  }

  saving.value = true
  saveMessage.value = ''

  try {
    const response = await fetch(`${getApiBaseUrl()}/api/settings`, {
      method: 'PUT',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(settings.value)
    })

    if (!response.ok) {
      throw new Error('保存设置失败')
    }

    // 同步 showServerName 到 taskStore
    taskStore.showServerName = settings.value.showServerName || false
    setLanguage(settings.value.language)

    // 更新原始设置
    originalSettings.value = JSON.parse(JSON.stringify(settings.value))

    showSaveMessage(t('settings.messages.saved'), true)
  } catch (error) {
    console.error('保存设置失败:', error)
    showSaveMessage(t('settings.messages.saveFailed'), false)
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
const addIP = () => {
  const trimmed = newIP.value.trim()

  if (!trimmed) {
    ipError.value = t('settings.ip.errors.empty')
    return
  }

  if (!validateIP(trimmed)) {
    ipError.value = t('settings.ip.errors.invalid')
    return
  }

  // 检查是否已存在
  if (settings.value.ipWhitelist && settings.value.ipWhitelist.includes(trimmed)) {
    ipError.value = t('settings.ip.errors.duplicate')
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
}

// 从白名单删除IP
const removeIP = (index) => {
  if (settings.value.ipWhitelist && settings.value.ipWhitelist.length > 0) {
    settings.value.ipWhitelist.splice(index, 1)

    // 如果删除了所有IP，自动关闭白名单
    if (settings.value.ipWhitelist.length === 0) {
      settings.value.enableIPWhitelist = false
    }
  }
}

/**
 * 统一处理密码设置弹窗提交。
 * 根据当前模式调用不同接口，并在成功后同步本地设置状态。
 *
 * @param {Object} payload - 弹窗提交的数据
 * @returns {Promise<void>}
 */
const handlePasswordSubmit = async (payload) => {
  passwordModalSaving.value = true

  try {
    let response

    if (payload.mode === 'add') {
      response = await fetch(`${getApiBaseUrl()}/api/settings/password`, {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          username: payload.username,
          password: payload.password
        })
      })
    } else if (payload.mode === 'change') {
      response = await fetch(`${getApiBaseUrl()}/api/settings/password`, {
        method: 'PUT',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          oldPassword: payload.oldPassword,
          newPassword: payload.newPassword
        })
      })
    } else {
      response = await fetch(`${getApiBaseUrl()}/api/settings/password`, {
        method: 'DELETE',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          username: payload.username,
          password: payload.password
        })
      })
    }

    const data = await response.json()
    if (!response.ok || (data.code !== undefined && data.code !== 0)) {
      throw new Error(data.message || '密码设置失败')
    }

    settings.value.passwordAuth = {
      enabled: data.data?.passwordAuth?.enabled || false,
      username: data.data?.passwordAuth?.username || ''
    }

    if (originalSettings.value) {
      originalSettings.value.passwordAuth = JSON.parse(JSON.stringify(settings.value.passwordAuth))
    }

    closePasswordModal()
    showSaveMessage(data.data?.message || '密码设置已更新', true)
  } catch (error) {
    console.error('密码设置失败:', error)
    showSaveMessage(error.message || '密码设置失败', false)
  } finally {
    passwordModalSaving.value = false
  }
}

onMounted(() => {
  loadSettings()
})
</script>

<style scoped>
.settings {
  /* 设置页与常规信息页保持一致的外层宽度，减少页面切换时的跳变感。 */
  max-width: var(--page-content-width);
  margin: 0 auto;
  padding: 2rem;
}

.settings-container {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

/* 保存按钮样式 */
.btn-save {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.6rem 1.25rem;
  background: var(--accent-color);
  color: white;
  border: none;
  border-radius: var(--radius-md);
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  white-space: nowrap;
}

.btn-save:hover:not(:disabled) {
  opacity: 0.9;
  transform: translateY(-1px);
  box-shadow: var(--shadow-lg);
}

.btn-save:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

.btn-save svg {
  width: 1rem;
  height: 1rem;
}

.btn-save .spinner {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.btn-reset {
  padding: 0.6rem 1.25rem;
  background: transparent;
  color: var(--text-secondary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
  white-space: nowrap;
}

.btn-reset:hover:not(:disabled) {
  border-color: var(--accent-color);
  color: var(--accent-color);
}

.btn-reset:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.setting-section {
  background: var(--card-bg);
  border-radius: var(--radius-xl);
  padding: 1.5rem;
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-lg);
}

.section-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 1.5rem 0;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid var(--border-color);
}

/* 下拉选择框样式 */
.interval-select {
  min-width: 150px;
  padding: 0.5rem 2rem 0.5rem 0.75rem;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
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
  box-shadow: var(--focus-ring);
}

.interval-select:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.interval-select:hover {
  border-color: var(--accent-color);
}

/* 保存消息 */
.save-message {
  position: fixed;
  bottom: 2rem;
  right: 2rem;
  padding: 1rem 1.5rem;
  border-radius: var(--radius-md);
  font-size: 0.875rem;
  font-weight: 500;
  box-shadow: var(--shadow-lg);
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
  border-radius: var(--radius-md);
  font-size: 1rem;
  background: var(--input-bg);
  color: var(--text-primary);
  text-align: center;
  transition: border-color 0.3s, box-shadow 0.3s;
}

.port-input:focus {
  outline: none;
  border-color: var(--accent-color);
  box-shadow: var(--focus-ring);
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
  color: var(--muted-warning-text);
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
  border-radius: var(--radius-md);
  font-size: 0.9rem;
  background: var(--input-bg);
  color: var(--text-primary);
  transition: border-color 0.3s, box-shadow 0.3s;
  font-family: 'SF Mono', 'Monaco', 'Consolas', monospace;
}

.ip-input:focus {
  outline: none;
  border-color: var(--accent-color);
  box-shadow: var(--focus-ring);
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
  border-radius: var(--radius-md);
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  white-space: nowrap;
}

.btn-add-ip:hover:not(:disabled) {
  opacity: 0.9;
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
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
  border-radius: var(--radius-sm);
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
  border-radius: var(--radius-md);
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
  border-radius: var(--radius-sm);
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
  border-radius: var(--radius-md);
}

.password-actions {
  display: flex;
  gap: 0.75rem;
  flex-shrink: 0;
}

.password-action-button {
  min-height: 2.1rem;
  padding: 0.45rem 0.9rem;
  font-size: 0.82rem;
  font-weight: 600;
  box-shadow: none;
}

.password-action-button--danger {
  background: linear-gradient(
    135deg,
    var(--accent-color) 0%,
    color-mix(in srgb, var(--accent-color) 82%, black) 100%
  );
  color: #fff;
  box-shadow: 0 8px 18px color-mix(in srgb, var(--accent-color) 28%, transparent);
}

.password-action-button--danger:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 10px 22px color-mix(in srgb, var(--accent-color) 36%, transparent);
}

@media (max-width: 768px) {
  .settings {
    padding: 1rem;
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

  .password-actions {
    width: 100%;
    flex-direction: column;
  }
}
</style>
