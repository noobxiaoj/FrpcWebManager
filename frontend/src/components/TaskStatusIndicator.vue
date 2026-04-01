<template>
  <span class="task-status-indicator" :class="`status-${status}`"></span>
</template>

<script setup>
defineProps({
  status: {
    type: String,
    required: true,
    validator: (value) => ['running', 'stopped', 'error'].includes(value)
  }
})
</script>

<style scoped>
.task-status-indicator {
  width: 8px;
  height: 8px;
  border-radius: var(--radius-pill);
  display: inline-block;
  flex-shrink: 0;
}

/* 运行状态 - 绿色 */
.status-running {
  background-color: var(--success-color);
  box-shadow: 0 0 0 3px var(--success-color-bg);
}

/* 停止状态 - 红色 */
.status-stopped {
  background-color: var(--danger-color);
  box-shadow: 0 0 0 3px var(--danger-color-bg);
}

/* 错误状态 - 红色 */
.status-error {
  background-color: var(--danger-color);
  box-shadow: 0 0 0 3px var(--danger-color-bg);
}
</style>
