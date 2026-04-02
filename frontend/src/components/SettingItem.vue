<template>
  <div class="setting-item">
    <div class="setting-info">
      <label class="setting-label" :for="id">
        <slot name="label">{{ label }}</slot>
      </label>
      <p v-if="hasDescription" class="setting-description">
        <slot name="description">{{ description }}</slot>
      </p>
      <p v-if="warning" class="setting-warning">
        <span class="setting-warning-icon" aria-hidden="true">!</span>
        <span>{{ warning }}</span>
      </p>
    </div>
    <div class="setting-control">
      <slot>
        <!-- 默认插槽：控制组件 -->
      </slot>
    </div>
  </div>
</template>

<script setup>
import { computed, useSlots } from 'vue'

const props = defineProps({
  id: {
    type: String,
    default: ''
  },
  label: {
    type: String,
    default: ''
  },
  description: {
    type: String,
    default: ''
  },
  warning: {
    type: String,
    default: ''
  }
})

const slots = useSlots()

/**
 * 仅在存在描述文本或描述插槽时渲染说明区域，
 * 避免设置项移除简介后仍然保留空白占位。
 *
 * @returns {boolean} 是否显示简介区域
 */
const hasDescription = computed(() => {
  return Boolean(props.description || slots.description)
})
</script>

<style scoped>
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

.setting-warning {
  display: inline-flex;
  align-items: center;
  gap: 0.45rem;
  font-size: 0.85rem;
  color: var(--muted-warning-text);
  margin: 0.5rem 0 0 0;
  font-weight: 400;
  opacity: 0.8;
}

.setting-warning-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 1rem;
  height: 1rem;
  border-radius: 999px;
  border: 1px solid currentColor;
  font-size: 0.72rem;
  font-weight: 700;
  line-height: 1;
  flex-shrink: 0;
}

.setting-control {
  flex-shrink: 0;
  position: relative;
  display: flex;
  align-items: center;
}

@media (max-width: 768px) {
  .setting-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }

  .setting-control {
    align-self: flex-end;
  }
}
</style>
