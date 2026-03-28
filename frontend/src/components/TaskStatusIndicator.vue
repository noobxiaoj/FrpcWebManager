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
  border-radius: 50%;
  display: inline-block;
  flex-shrink: 0;
}

/* 运行状态 - 绿色 */
.status-running {
  background-color: #22c55e;
  box-shadow: 0 0 0 3px rgba(34, 197, 94, 0.15);
}

/* 停止状态 - 红色 */
.status-stopped {
  background-color: #ef4444;
  box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.15);
}

/* 错误状态 - 红色 */
.status-error {
  background-color: #ef4444;
  box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.15);
}
</style>
