<template>
  <button
    class="app-button"
    :class="[
      `app-button--${variant}`,
      { 'app-button--preserve-style': preserveStyle },
      { 'app-button--icon-only': !hasContent },
      { 'app-button--loading': loading }
    ]"
    :disabled="disabled || loading"
    :title="title"
  >
    <span v-if="loading" class="app-button__spinner"></span>
    <span v-else-if="$slots.icon" class="app-button__icon">
      <slot name="icon"></slot>
    </span>
    <span v-if="hasContent" class="app-button__content">
      <slot></slot>
    </span>
  </button>
</template>

<script setup>
import { computed, useSlots } from 'vue'

const props = defineProps({
  variant: {
    type: String,
    default: 'primary',
    validator: (value) => ['primary', 'secondary', 'danger', 'ghost', 'icon'].includes(value)
  },
  disabled: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: ''
  },
  /**
   * 是否尽量保留调用方原有按钮样式。
   * 开启后会重置 AppButton 的默认视觉样式，让页面原本的 class 继续主导外观，
   * 这样可以在统一替换为 AppButton 的同时，尽量避免影响现有设计。
   */
  preserveStyle: {
    type: Boolean,
    default: false
  }
})

const slots = useSlots()

const hasContent = computed(() => {
  return slots.default && slots.default().length > 0
})
</script>

<style scoped>
.app-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: var(--radius-md);
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  box-shadow: var(--shadow-md);
  position: relative;
  overflow: hidden;
}

.app-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.app-button--preserve-style {
  all: unset;
  box-sizing: border-box;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.app-button--preserve-style:disabled {
  cursor: not-allowed;
}

.app-button__icon {
  display: flex;
  align-items: center;
  justify-content: center;
}

.app-button__icon :deep(svg) {
  width: 1.25rem;
  height: 1.25rem;
}

.app-button__content {
  display: flex;
  align-items: center;
}

.app-button__spinner {
  width: 1.25rem;
  height: 1.25rem;
  border: 2px solid currentColor;
  border-top-color: transparent;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

/* Primary variant - 主要按钮（蓝色/强调色） */
.app-button--primary {
  background: var(--accent-color);
  color: white;
}

.app-button--primary:hover:not(:disabled) {
  opacity: 0.9;
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
}

/* Secondary variant - 次要按钮（边框样式） */
.app-button--secondary {
  background: var(--bg-primary);
  color: var(--text-primary);
  border: 1.5px solid var(--border-color);
  box-shadow: none;
}

.app-button--secondary:hover:not(:disabled) {
  background: var(--bg-secondary);
  border-color: var(--text-secondary);
}

/* Danger variant - 危险按钮（红色） */
.app-button--danger {
  background: var(--danger-color);
  color: white;
  box-shadow: var(--shadow-md);
}

.app-button--danger:hover:not(:disabled) {
  opacity: 0.9;
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
}

/* Ghost variant - 幽灵按钮（透明背景） */
.app-button--ghost {
  background: transparent;
  color: var(--text-secondary);
  box-shadow: none;
  padding: 0.5rem;
}

.app-button--ghost:hover:not(:disabled) {
  background: var(--accent-color-bg);
  color: var(--accent-color);
}

/* Icon variant - 纯图标按钮 */
.app-button--icon {
  width: 2.5rem;
  height: 2.5rem;
  padding: 0;
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  box-shadow: none;
}

.app-button--icon:hover:not(:disabled) {
  background: var(--accent-color-bg);
  border-color: var(--accent-color);
  transform: rotate(90deg);
}

/* Icon-only button (no content) */
.app-button--icon-only {
  padding: 0.5rem;
}

.app-button--icon-only.app-button--primary,
.app-button--icon-only.app-button--danger {
  width: 2.5rem;
  height: 2.5rem;
}

/* Loading state */
.app-button--loading {
  pointer-events: none;
}

/* 响应式 */
@media (max-width: 768px) {
  .app-button {
    padding: 0.625rem 1.25rem;
    font-size: 0.9rem;
  }
}
</style>
