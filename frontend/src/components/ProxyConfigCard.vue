<template>
  <div class="proxy-card" :class="{ 'proxy-card-compact': !proxy.isEditing }">
    <div class="proxy-header">
      <h3 class="proxy-name">
        <span v-if="!proxy.isEditing" class="proxy-type-badge" :class="`badge-${proxy.type}`">{{ proxy.type.toUpperCase() }}</span>
        {{ proxy.name || '未命名配置' }}
      </h3>

      <!-- 缩略模式 - 在标题行显示配置信息 -->
      <template v-if="!proxy.isEditing">
        <div class="proxy-compact-info" v-if="proxy.type === 'tcp' || proxy.type === 'udp'">
          <div class="info-item">本地IP: <span class="info-value">{{ proxy.localIP }}:{{ proxy.localPort }}</span></div>
          <div class="info-item">远程端口: <span class="info-value">{{ proxy.remotePort }}</span></div>
        </div>
        <div class="proxy-compact-info proxy-compact-info-http" v-else>
          <div class="info-item">本地IP: <span class="info-value">{{ proxy.localIP }}:{{ proxy.localPort }}</span></div>
          <div class="info-item">域名: <span class="info-value">{{ proxy.customDomains || proxy.subdomain || '-' }}</span></div>
        </div>
      </template>

      <div class="proxy-actions">
        <!-- 编辑模式显示验证按钮和删除按钮 -->
        <button
          v-if="proxy.isEditing"
          type="button"
          class="btn-check"
          @click="handleCheck"
          title="验证并完成"
        >
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="20 6 9 17 4 12"></polyline>
          </svg>
        </button>
        <button
          v-else
          type="button"
          class="btn-edit"
          @click="handleToggleEdit"
          title="编辑"
        >
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
            <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
          </svg>
        </button>
        <button v-if="proxy.isEditing" type="button" class="btn-remove" @click="handleRemove" title="删除此Frpc">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>
    </div>

    <!-- 编辑模式 - 完整表单 -->
    <div v-if="proxy.isEditing" class="proxy-content">
      <div class="form-row-4">
        <div class="form-group">
          <label>协议类型 *</label>
          <select v-model="proxy.type" required>
            <option value="tcp">TCP</option>
            <option value="udp">UDP</option>
            <option value="http">HTTP</option>
            <option value="https">HTTPS</option>
            <option value="tcpmux">TCP Multiplexer</option>
          </select>
        </div>

        <div class="form-group" :class="{ 'has-error': proxy.errors.name }">
          <label>名称 *</label>
          <input v-model="proxy.name" type="text" required placeholder="例如: ssh" />
          <span v-if="proxy.errors.name" class="error-message">{{ proxy.errors.name }}</span>
        </div>

        <div class="form-group" :class="{ 'has-error': proxy.errors.localIP }">
          <label>本地IP *</label>
          <input v-model="proxy.localIP" type="text" required placeholder="例如: 127.0.0.1" />
          <span v-if="proxy.errors.localIP" class="error-message">{{ proxy.errors.localIP }}</span>
        </div>

        <div class="form-group" :class="{ 'has-error': proxy.errors.localPort }">
          <label>本地端口 *</label>
          <input v-model="proxy.localPort" type="number" required min="1" max="65535" placeholder="例如: 22" />
          <span v-if="proxy.errors.localPort" class="error-message">{{ proxy.errors.localPort }}</span>
        </div>
      </div>

      <div class="form-row" v-if="proxy.type === 'tcp' || proxy.type === 'udp'">
        <div class="form-group" :class="{ 'has-error': proxy.errors.remotePort }">
          <label>远程端口 *</label>
          <input v-model="proxy.remotePort" type="number" required min="1" max="65535" placeholder="例如: 6000" />
          <span v-if="proxy.errors.remotePort" class="error-message">{{ proxy.errors.remotePort }}</span>
        </div>
      </div>

      <div class="form-row" v-if="proxy.type === 'http' || proxy.type === 'https'">
        <div class="form-group" :class="{ 'has-error': proxy.errors.customDomains }">
          <label>自定义域名(多个用逗号分隔)</label>
          <input v-model="proxy.customDomains" type="text" placeholder="例如: www.example.com,test.example.com" />
          <span v-if="proxy.errors.customDomains" class="error-message">{{ proxy.errors.customDomains }}</span>
        </div>

        <div class="form-group" :class="{ 'has-error': proxy.errors.subdomain }">
          <label>子域名</label>
          <input v-model="proxy.subdomain" type="text" placeholder="例如: myapp" />
          <span v-if="proxy.errors.subdomain" class="error-message">{{ proxy.errors.subdomain }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>

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

const validateProxy = (proxy) => {
  const errors = {}

  // 基本字段验证
  if (!proxy.name || proxy.name.trim() === '') {
    errors.name = '请填写名称'
  }

  if (!proxy.localIP || proxy.localIP.trim() === '') {
    errors.localIP = '请填写本地IP'
  }

  if (!proxy.localPort || proxy.localPort === '') {
    errors.localPort = '请填写本地端口'
  } else {
    const localPort = parseInt(proxy.localPort)
    if (isNaN(localPort) || localPort < 1 || localPort > 65535) {
      errors.localPort = '本地端口必须是1-65535之间的数字'
    }
  }

  // 根据类型验证特定字段
  if (proxy.type === 'tcp' || proxy.type === 'udp') {
    if (!proxy.remotePort || proxy.remotePort === '') {
      errors.remotePort = '请填写远程端口'
    } else {
      const remotePort = parseInt(proxy.remotePort)
      if (isNaN(remotePort) || remotePort < 1 || remotePort > 65535) {
        errors.remotePort = '远程端口必须是1-65535之间的数字'
      }
    }
  }

  if (proxy.type === 'http' || proxy.type === 'https') {
    if (!proxy.customDomains && !proxy.subdomain) {
      errors.customDomains = '必须填写自定义域名或子域名'
      errors.subdomain = '必须填写自定义域名或子域名'
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
  border-radius: 12px;
  padding: 1.5rem;
  margin-bottom: 1rem;
  transition: all 0.3s;
  box-shadow: 0 2px 8px var(--shadow-color);
}

.proxy-card:hover {
  border-color: var(--accent-color);
  box-shadow: 0 2px 12px var(--shadow-hover);
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
  border-radius: 4px;
  font-weight: 600;
}

/* TCP - 蓝色 */
.proxy-type-badge.badge-tcp {
  background: rgba(59, 130, 246, 0.15);
  color: #3b82f6;
}

/* UDP - 红色 */
.proxy-type-badge.badge-udp {
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
}

/* HTTP - 紫色 */
.proxy-type-badge.badge-http {
  background: rgba(168, 85, 247, 0.15);
  color: #a855f7;
}

/* HTTPS - 绿色 */
.proxy-type-badge.badge-https {
  background: rgba(34, 197, 94, 0.15);
  color: #22c55e;
}

/* TCPMUX - 粉色 */
.proxy-type-badge.badge-tcpmux {
  background: rgba(236, 72, 153, 0.15);
  color: #ec4899;
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
  border-radius: 6px;
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
  box-shadow: 0 2px 8px var(--shadow-hover);
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
  box-shadow: 0 2px 8px var(--shadow-hover);
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
  box-shadow: 0 2px 8px var(--shadow-hover);
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
  border-color: #ef4444;
  background: rgba(239, 68, 68, 0.05);
  animation: shake 0.3s ease-in-out;
}

.form-group.has-error label {
  color: #ef4444;
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
  color: #ef4444;
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
  border-radius: 6px;
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
  box-shadow: 0 0 0 3px var(--accent-color-bg);
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
