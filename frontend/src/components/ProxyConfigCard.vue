<template>
  <div class="proxy-card" :class="{ 'proxy-card-compact': !proxy.isEditing }">
    <div class="proxy-header">
      <h3 class="proxy-name">
        <span v-if="!proxy.isEditing" class="proxy-type-badge" :class="`badge-${proxy.type}`">{{ proxy.type.toUpperCase() }}</span>
        {{ proxy.name || t('proxy.unnamed') }}
      </h3>

      <!-- 缩略模式 - 在标题行显示配置信息 -->
      <template v-if="!proxy.isEditing">
        <div class="proxy-compact-info" v-if="proxy.type === 'tcp' || proxy.type === 'udp'">
          <div class="info-item">{{ t('proxy.localIp') }}: <span class="info-value">{{ proxy.localIP }}:{{ proxy.localPort }}</span></div>
          <div class="info-item">{{ t('proxy.remotePort') }}: <span class="info-value">{{ proxy.remotePort }}</span></div>
        </div>
        <div class="proxy-compact-info proxy-compact-info-http" v-else>
          <div class="info-item">{{ t('proxy.localIp') }}: <span class="info-value">{{ proxy.localIP }}:{{ proxy.localPort }}</span></div>
          <div class="info-item">{{ t('proxy.domain') }}: <span class="info-value">{{ proxy.customDomains || proxy.subdomain || '-' }}</span></div>
        </div>
      </template>

      <div class="proxy-actions">
        <!-- 编辑模式显示验证按钮和删除按钮 -->
        <AppButton
          v-if="proxy.isEditing"
          type="button"
          class="btn-check"
          preserve-style
          @click="handleCheck"
          :title="t('proxy.check')"
        >
          <template #icon>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="20 6 9 17 4 12"></polyline>
            </svg>
          </template>
        </AppButton>
        <AppButton
          v-else
          type="button"
          class="btn-edit"
          preserve-style
          @click="handleToggleEdit"
          :title="t('proxy.edit')"
        >
          <template #icon>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
              <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
            </svg>
          </template>
        </AppButton>
        <AppButton v-if="proxy.isEditing" type="button" class="btn-remove" preserve-style @click="handleRemove" :title="t('proxy.remove')">
          <template #icon>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </template>
        </AppButton>
      </div>
    </div>

    <!-- 编辑模式 - 完整表单 -->
    <div v-if="proxy.isEditing" class="proxy-content">
      <div class="form-row-4">
        <div class="form-group">
          <label>{{ t('proxy.protocolType') }}</label>
          <select v-model="proxy.type" required>
            <option value="tcp">TCP</option>
            <option value="udp">UDP</option>
            <option value="http">HTTP</option>
            <option value="https">HTTPS</option>
            <option value="tcpmux">TCP Multiplexer</option>
          </select>
        </div>

        <div class="form-group" :class="{ 'has-error': proxy.errors.name }">
          <label>{{ t('proxy.name') }}</label>
          <input v-model="proxy.name" type="text" required :placeholder="t('proxy.placeholders.name')" />
          <span v-if="proxy.errors.name" class="error-message">{{ proxy.errors.name }}</span>
        </div>

        <div class="form-group" :class="{ 'has-error': proxy.errors.localIP }">
          <label>{{ t('proxy.localIp') }} *</label>
          <input v-model="proxy.localIP" type="text" required :placeholder="t('proxy.placeholders.localIp')" />
          <span v-if="proxy.errors.localIP" class="error-message">{{ proxy.errors.localIP }}</span>
        </div>

        <div class="form-group" :class="{ 'has-error': proxy.errors.localPort }">
          <label>{{ t('proxy.localPort') }} *</label>
          <input v-model="proxy.localPort" type="number" required min="1" max="65535" :placeholder="t('proxy.placeholders.localPort')" />
          <span v-if="proxy.errors.localPort" class="error-message">{{ proxy.errors.localPort }}</span>
        </div>
      </div>

      <div class="form-row" v-if="proxy.type === 'tcp' || proxy.type === 'udp'">
        <div class="form-group" :class="{ 'has-error': proxy.errors.remotePort }">
          <label>{{ t('proxy.remotePort') }} *</label>
          <input v-model="proxy.remotePort" type="number" required min="1" max="65535" :placeholder="t('proxy.placeholders.remotePort')" />
          <span v-if="proxy.errors.remotePort" class="error-message">{{ proxy.errors.remotePort }}</span>
        </div>
      </div>

      <div class="form-row" v-if="proxy.type === 'http' || proxy.type === 'https'">
        <div class="form-group" :class="{ 'has-error': proxy.errors.customDomains }">
          <label>{{ t('proxy.customDomains') }}</label>
          <input v-model="proxy.customDomains" type="text" :placeholder="t('proxy.placeholders.customDomains')" />
          <span v-if="proxy.errors.customDomains" class="error-message">{{ proxy.errors.customDomains }}</span>
        </div>

        <div class="form-group" :class="{ 'has-error': proxy.errors.subdomain }">
          <label>{{ t('proxy.subdomain') }}</label>
          <input v-model="proxy.subdomain" type="text" :placeholder="t('proxy.placeholders.subdomain')" />
          <span v-if="proxy.errors.subdomain" class="error-message">{{ proxy.errors.subdomain }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import AppButton from '@/components/AppButton.vue'
import { useI18n } from '@/utils/i18n'

const props = defineProps({
  proxy: {
    type: Object,
    required: true
  },
  index: {
    type: Number,
    required: true
  }
})

const emit = defineEmits(['check', 'toggle-edit', 'remove'])
const { t } = useI18n()

const validateProxy = (proxy) => {
  const errors = {}

  // 基本字段验证
  if (!proxy.name || proxy.name.trim() === '') {
    errors.name = t('proxy.validation.nameRequired')
  }

  if (!proxy.localIP || proxy.localIP.trim() === '') {
    errors.localIP = t('proxy.validation.localIpRequired')
  }

  if (!proxy.localPort || proxy.localPort === '') {
    errors.localPort = t('proxy.validation.localPortRequired')
  } else {
    const localPort = parseInt(proxy.localPort)
    if (isNaN(localPort) || localPort < 1 || localPort > 65535) {
      errors.localPort = t('proxy.validation.localPortInvalid')
    }
  }

  // 根据类型验证特定字段
  if (proxy.type === 'tcp' || proxy.type === 'udp') {
    if (!proxy.remotePort || proxy.remotePort === '') {
      errors.remotePort = t('proxy.validation.remotePortRequired')
    } else {
      const remotePort = parseInt(proxy.remotePort)
      if (isNaN(remotePort) || remotePort < 1 || remotePort > 65535) {
        errors.remotePort = t('proxy.validation.remotePortInvalid')
      }
    }
  }

  if (proxy.type === 'http' || proxy.type === 'https') {
    if (!proxy.customDomains && !proxy.subdomain) {
      errors.customDomains = t('proxy.validation.domainOrSubdomainRequired')
      errors.subdomain = t('proxy.validation.domainOrSubdomainRequired')
    }
  }

  return errors
}

const handleCheck = () => {
  const errors = validateProxy(props.proxy)

  if (Object.keys(errors).length > 0) {
    // 设置错误信息 - 使用 Vue 3 的方式确保响应式
    Object.assign(props.proxy, { errors })
    return
  }

  // 验证通过,切换到缩略卡片模式
  Object.assign(props.proxy, {
    isEditing: false,
    isValid: true,
    errors: {}
  })

  emit('check', props.index)
}

const handleToggleEdit = () => {
  const newEditingState = !props.proxy.isEditing

  // 切换编辑状态并清除错误
  Object.assign(props.proxy, {
    isEditing: newEditingState,
    errors: newEditingState ? {} : props.proxy.errors
  })

  emit('toggle-edit', props.index)
}

const handleRemove = () => {
  emit('remove', props.index)
}
</script>

<style scoped>
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

.proxy-card-compact {
  padding: 0.75rem 1rem;
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

.proxy-actions {
  display: flex;
  gap: 0.5rem;
  flex-shrink: 0;
}

.btn-check,
.btn-edit,
.btn-remove {
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  border-radius: var(--radius-sm);
  width: 2rem;
  height: 2rem;
  padding: 0;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-check {
  background: var(--success-color-bg);
  color: var(--success-color);
  border: 1px solid var(--success-color);
}

.btn-check:hover {
  background: var(--success-color-bg);
  opacity: 0.8;
  transform: scale(1.1);
  box-shadow: var(--shadow-md);
}

.btn-check:active {
  transform: scale(1);
}

.btn-check svg {
  width: 1.25rem;
  height: 1.25rem;
  display: block;
}

.btn-edit {
  background: var(--edit-color-bg);
  color: var(--edit-color);
  border: 1px solid var(--edit-color);
}

.btn-edit:hover {
  background: var(--edit-color-bg);
  opacity: 0.8;
  transform: scale(1.1);
  box-shadow: var(--shadow-md);
}

.btn-edit:active {
  transform: scale(1);
}

.btn-edit svg {
  width: 1.25rem;
  height: 1.25rem;
  display: block;
}

.btn-remove {
  background: var(--danger-color-bg);
  color: var(--danger-color);
  border: 1px solid var(--danger-color);
}

.btn-remove:hover {
  background: var(--danger-color-bg);
  opacity: 0.8;
  transform: scale(1.1);
  box-shadow: var(--shadow-md);
}

.btn-remove:active {
  transform: scale(1);
}

.btn-remove svg {
  width: 1.25rem;
  height: 1.25rem;
  display: block;
}

.proxy-content {
  margin-top: 1rem;
}

.form-group {
  margin-bottom: 1.25rem;
}

.form-group.has-error input,
.form-group.has-error select {
  border-color: var(--danger-color);
  background: color-mix(in srgb, var(--danger-color) 8%, transparent);
  animation: shake 0.3s ease-in-out;
}

.form-group.has-error label {
  color: var(--danger-color);
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-5px); }
  75% { transform: translateX(5px); }
}

.error-message {
  display: block;
  margin-top: 0.375rem;
  font-size: 0.75rem;
  color: var(--danger-color);
  font-weight: 500;
}

.form-group label {
  display: block;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 0.5rem;
}

.form-group input,
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
.form-group select:focus {
  outline: none;
  border-color: var(--accent-color);
  box-shadow: var(--focus-ring);
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.form-row-4 {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr 1fr;
  gap: 1rem;
}

@media (max-width: 768px) {
  /* 移动端减小卡片内边距 */
  .proxy-card {
    padding: 0.75rem;
  }

  .proxy-card-compact {
    padding: 0.5rem 0.75rem;
  }

  .form-row {
    grid-template-columns: 1fr;
  }

  .form-row-4 {
    grid-template-columns: 1fr;
  }

  /* 小卡片响应式布局 */
  .proxy-card-compact .proxy-header {
    flex-direction: column;
    align-items: stretch;
    gap: 0.75rem;
  }

  .proxy-card-compact .proxy-name {
    width: 100%;
    min-width: auto;
    margin-bottom: 0.5rem;
  }

  .proxy-compact-info {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.5rem;
    width: 100%;
  }

  .proxy-compact-info .info-item {
    min-width: auto;
  }
}
</style>
