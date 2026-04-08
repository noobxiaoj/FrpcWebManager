<template>
  <div class="server-header">
    <div class="server-title">
      <h3 class="server-name">{{ server.name }}</h3>
      <div
        class="status-badge"
        :class="{
          'status-online': server.status === 'online',
          'status-offline': server.status === 'offline',
          'status-no-task': server.status === 'no_task',
          'status-fault': server.status === 'fault',
          'status-suspected-abnormal': server.status === 'suspected_abnormal'
        }"
      >
        <span class="status-indicator"></span>
        <span class="status-text">{{ statusText }}</span>
      </div>
    </div>

    <div class="server-meta">
      <div class="server-meta-details">
        <div class="meta-item">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="2" y="2" width="20" height="8" rx="2" ry="2"></rect>
            <rect x="2" y="14" width="20" height="8" rx="2" ry="2"></rect>
            <line x1="6" y1="6" x2="6" y2="6"></line>
            <line x1="6" y1="18" x2="6" y2="18"></line>
          </svg>
          <span>{{ server.address }}</span>
        </div>
      </div>

      <div
        v-if="runningTaskCount > 0"
        class="meta-item meta-item--running-task"
        :title="'正在运行的任务: ' + runningTaskCount"
      >
        <span>{{ runningTaskCount }} 个任务</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  server: {
    type: Object,
    required: true
  },
  /**
   * 首页卡片会把“正在运行的任务数”作为单独字段传入。
   * 这里默认回退为 0，避免父组件异步渲染时出现空值闪烁。
   */
  runningTaskCount: {
    type: Number,
    default: 0
  }
})

// 状态文本
const statusText = computed(() => {
  switch (props.server.status) {
    case 'online':
      return '在线'
    case 'offline':
      return '离线'
    case 'no_task':
      return '无任务'
    case 'fault':
      return '故障'
    case 'suspected_abnormal':
      return '疑似异常'
    default:
      return '未知'
  }
})

</script>

<style scoped>
/* 服务器头部 */
.server-header {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.server-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
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

.status-fault {
  background: var(--warning-color-bg);
  color: var(--warning-color);
}

.status-fault .status-indicator {
  background: var(--warning-color);
  box-shadow: 0 0 8px var(--warning-color);
}

.status-fault .status-indicator::after {
  background: var(--warning-color);
}

.status-suspected-abnormal {
  background: var(--info-color-bg);
  color: var(--info-color);
}

.status-suspected-abnormal .status-indicator {
  background: var(--info-color);
  box-shadow: 0 0 8px var(--info-color);
}

.status-suspected-abnormal .status-indicator::after {
  background: var(--info-color);
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
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}

.server-meta-details {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  color: var(--text-secondary);
  font-size: 0.875rem;
}

.meta-item span {
  line-height: 1;
  padding-top: 1px;
}

.meta-item svg {
  width: 1rem;
  height: 1rem;
  flex-shrink: 0;
}

.meta-item--running-task {
  display: inline-flex;
  align-items: center;
  font-weight: 600;
  color: var(--text-primary);
}

@media (max-width: 768px) {
  .server-meta {
    flex-direction: column;
    gap: 0.5rem;
  }

  .meta-item--running-task {
    margin-left: 0;
  }
}
</style>
