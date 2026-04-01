<template>
  <div class="server-header">
    <div class="server-title">
      <h3 class="server-name">{{ server.name }}</h3>
      <div
        class="status-badge"
        :class="{
          'status-online': server.status === 'online',
          'status-offline': server.status === 'offline',
          'status-no-task': server.status === 'no_task'
        }"
      >
        <span class="status-indicator"></span>
        <span class="status-text">{{ statusText }}</span>
      </div>
    </div>

    <div class="server-meta">
      <div class="meta-item">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="2" y="2" width="20" height="8" rx="2" ry="2"></rect>
          <rect x="2" y="14" width="20" height="8" rx="2" ry="2"></rect>
          <line x1="6" y1="6" x2="6" y2="6"></line>
          <line x1="6" y1="18" x2="6" y2="18"></line>
        </svg>
        <span>{{ server.address }}</span>
      </div>

      <div class="meta-item" v-if="settings.showServerPort && server.webServerPort > 0" :title="'frpc webServer 端口: ' + server.webServerPort">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="3"></circle>
          <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
        </svg>
        <span>端口: {{ server.webServerPort }}</span>
      </div>

      <div class="meta-item">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"></circle>
          <polyline points="12 6 12 12 16 14"></polyline>
        </svg>
        <span>{{ server.uptime || '未知' }}</span>
        <span v-if="settings.showRefreshTime && server.lastRefreshTime" class="refresh-time" :title="'距离上次刷新: ' + refreshTimeDiff">
          ({{ refreshTimeDiff }})
        </span>
      </div>

      <AppButton
        class="btn-lock-server"
        @click="$emit('toggle-lock')"
        :title="server.locked ? '解锁服务器' : '锁定服务器'"
        preserve-style
      >
        <template #icon>
          <svg v-if="!server.locked" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
            <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
          </svg>
          <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
            <path d="M7 11V7a5 5 0 0 1 9.9-1"></path>
          </svg>
        </template>
      </AppButton>

      <AppButton
        class="btn-delete-server"
        @click="$emit('delete')"
        title="删除服务器"
        preserve-style
      >
        <template #icon>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="3 6 5 6 21 6"></polyline>
            <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
            <line x1="10" y1="11" x2="10" y2="17"></line>
            <line x1="14" y1="11" x2="14" y2="17"></line>
          </svg>
        </template>
      </AppButton>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import AppButton from '@/components/AppButton.vue'

const props = defineProps({
  server: {
    type: Object,
    required: true
  },
  settings: {
    type: Object,
    default: () => ({
      showServerPort: true,
      showRefreshTime: true
    })
  }
})

defineEmits(['toggle-lock', 'delete'])

// 当前时间，用于计算刷新时间差
const currentTime = ref(Date.now())

// 每3秒更新当前时间
let timeInterval = null
watch(() => props.server, () => {
  if (!timeInterval) {
    timeInterval = setInterval(() => {
      currentTime.value = Date.now()
    }, 3000)
  }
}, { immediate: true })

// 状态文本
const statusText = computed(() => {
  switch (props.server.status) {
    case 'online':
      return '在线'
    case 'offline':
      return '离线'
    case 'no_task':
      return '无任务'
    default:
      return '未知'
  }
})

// 计算刷新时间差
const refreshTimeDiff = computed(() => {
  if (!props.server.lastRefreshTime) return ''

  const now = currentTime.value
  const then = new Date(props.server.lastRefreshTime).getTime()
  const diffMs = now - then
  const diffSecs = Math.floor(diffMs / 1000)
  const diffMins = Math.floor(diffSecs / 60)
  const diffHours = Math.floor(diffMins / 60)
  const diffDays = Math.floor(diffHours / 24)

  if (diffDays > 0) {
    return `${diffDays}天`
  } else if (diffHours > 0) {
    return `${diffHours}小时`
  } else if (diffMins > 0) {
    return `${diffMins}分钟`
  } else if (diffSecs > 0) {
    return `${diffSecs}秒`
  } else {
    return '刚刚'
  }
})
</script>

<style scoped>
/* 服务器头部 */
.server-header {
  padding: 1rem 1.25rem;
  padding-top: 1.25rem;
  border-bottom: 1px solid var(--border-color);
  background: linear-gradient(135deg, var(--header-bg), var(--bg-primary));
  position: relative;
}

.server-header::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(90deg,
    transparent 0%,
    var(--border-color) 20%,
    var(--border-color) 80%,
    transparent 100%);
}

.server-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.server-name {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.status-badge {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  padding: 0.35rem 0.75rem;
  border-radius: var(--radius-pill);
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  transition: all 0.3s ease;
}

.status-indicator {
  position: relative;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.status-indicator::after {
  content: '';
  position: absolute;
  inset: -4px;
  border-radius: 50%;
  opacity: 0.3;
  animation: pulse 2s ease-in-out infinite;
}

.status-online {
  background: var(--success-color-bg);
  color: var(--success-color);
}

.status-online .status-indicator {
  background: var(--success-color);
  box-shadow: 0 0 8px var(--success-color);
}

.status-online .status-indicator::after {
  background: var(--success-color);
}

.status-offline {
  background: var(--danger-color-bg);
  color: var(--danger-color);
}

.status-offline .status-indicator {
  background: var(--danger-color);
  box-shadow: 0 0 8px var(--danger-color);
}

.status-offline .status-indicator::after {
  background: var(--danger-color);
  animation: none;
}

.status-no-task {
  background: var(--neutral-soft-bg);
  color: var(--text-secondary);
}

.status-no-task .status-indicator {
  background: var(--text-secondary);
  box-shadow: 0 0 8px var(--text-secondary);
}

.status-no-task .status-indicator::after {
  background: var(--text-secondary);
  animation: none;
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
    opacity: 0.3;
  }
  50% {
    transform: scale(1.8);
    opacity: 0.1;
  }
}

.server-meta {
  display: flex;
  gap: 1.25rem;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  color: var(--text-secondary);
  font-size: 0.8rem;
}

.refresh-time {
  font-size: 0.75rem;
  color: var(--text-secondary);
  opacity: 0.8;
  margin-left: 0.2rem;
}

.meta-item svg {
  width: 0.9rem;
  height: 0.9rem;
}

/* 锁定服务器按钮 */
.btn-lock-server {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  padding: 0;
  background: transparent;
  border: none;
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.3s;
}

.btn-lock-server:hover {
  background: var(--accent-color-bg);
  color: var(--accent-color);
  transform: scale(1.1);
}

.btn-lock-server svg {
  width: 0.9rem;
  height: 0.9rem;
}

/* 删除服务器按钮 */
.btn-delete-server {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  padding: 0;
  background: transparent;
  border: none;
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.3s;
  margin-left: auto;
}

.btn-delete-server:hover {
  background: var(--danger-color-bg);
  color: var(--danger-color);
  transform: scale(1.1);
}

.btn-delete-server svg {
  width: 0.9rem;
  height: 0.9rem;
}

@media (max-width: 768px) {
  .server-meta {
    flex-direction: column;
    gap: 0.5rem;
  }
}
</style>
