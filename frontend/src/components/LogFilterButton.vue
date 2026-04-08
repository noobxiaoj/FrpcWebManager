<template>
  <button
    type="button"
    class="log-filter-button"
    :class="[
      `log-filter-button--${tone}`,
      { 'log-filter-button--active': active }
    ]"
    :aria-pressed="String(active)"
    @click="emit('toggle')"
  >
    {{ label }} {{ count }}
  </button>
</template>

<script setup>
/**
 * 日志筛选按钮组件。
 * 负责统一渲染日志筛选项的文案、颜色和选中态，
 * 让父组件只关心“当前按钮是否选中”和“点击后如何切换”。
 */
const props = defineProps({
  /**
   * 按钮显示文案。
   * 例如：信息、错误、警告。
   */
  label: {
    type: String,
    required: true
  },
  /**
   * 当前类型对应的日志数量。
   * 会直接展示在按钮文案后方。
   */
  count: {
    type: Number,
    required: true
  },
  /**
   * 按钮是否处于选中状态。
   * 选中后会呈现更明显的按下视觉效果。
   */
  active: {
    type: Boolean,
    default: false
  },
  /**
   * 按钮的视觉语义类型。
   * info / error / warn 分别映射到不同的颜色方案。
   */
  tone: {
    type: String,
    required: true,
    validator: (value) => ['info', 'error', 'warn'].includes(value)
  }
})

/**
 * 对外抛出切换事件。
 * 父组件收到后负责更新筛选状态。
 */
const emit = defineEmits(['toggle'])

void props
</script>

<style scoped>
.log-filter-button {
  appearance: none;
  outline: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 2rem;
  padding: 0.3rem 0.7rem;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-pill);
  background: var(--bg-secondary);
  color: var(--text-secondary);
  font-weight: 500;
  font-size: inherit;
  font-family: inherit;
  line-height: 1;
  cursor: pointer;
  opacity: 0.72;
  box-shadow: inset 0 0 0 1px transparent;
  transition:
    border-color 0.2s ease,
    background 0.2s ease,
    color 0.2s ease,
    transform 0.2s ease,
    box-shadow 0.2s ease,
    opacity 0.2s ease;
}

.log-filter-button:focus,
.log-filter-button:focus-visible {
  outline: none;
}

.log-filter-button:hover {
  opacity: 0.92;
  transform: translateY(-1px);
}

.log-filter-button:active {
  transform: translateY(1px) scale(0.99);
}

.log-filter-button--active {
  border-color: var(--accent-color);
  background: var(--accent-color-bg);
  color: var(--accent-color);
  opacity: 1;
  font-weight: 700;
  transform: translateY(1px);
  box-shadow:
    inset 0 2px 6px rgba(15, 23, 42, 0.12),
    0 0 0 2px color-mix(in srgb, var(--accent-color) 22%, transparent 78%);
}

.log-filter-button--error {
  border-color: color-mix(in srgb, var(--danger-color) 22%, var(--border-color) 78%);
  background: color-mix(in srgb, var(--danger-color-bg) 45%, var(--bg-secondary) 55%);
  color: var(--danger-color);
}

.log-filter-button--info {
  border-color: var(--log-filter-info-border);
  background: var(--log-filter-info-bg);
  color: var(--log-filter-info-color);
}

.log-filter-button--warn {
  border-color: var(--log-filter-warn-border);
  background: var(--log-filter-warn-bg);
  color: var(--log-filter-warn-color);
}

.log-filter-button--active.log-filter-button--info {
  border-color: var(--log-filter-info-color);
  background: var(--log-filter-info-active-bg);
  color: var(--log-filter-info-active-color);
  box-shadow:
    inset 0 2px 6px var(--log-filter-info-active-shadow),
    0 0 0 2px var(--log-filter-info-active-ring);
}

.log-filter-button--active.log-filter-button--error {
  border-color: var(--danger-color);
  background: color-mix(in srgb, var(--danger-color-bg) 88%, var(--bg-secondary) 12%);
  color: var(--danger-color);
  box-shadow:
    inset 0 2px 6px color-mix(in srgb, var(--danger-color) 14%, transparent 86%),
    0 0 0 2px color-mix(in srgb, var(--danger-color) 22%, transparent 78%);
}

.log-filter-button--active.log-filter-button--warn {
  border-color: var(--warning-color);
  background: var(--log-filter-warn-active-bg);
  color: var(--log-filter-warn-active-color);
  box-shadow:
    inset 0 2px 6px var(--log-filter-warn-active-shadow),
    0 0 0 2px var(--log-filter-warn-active-ring);
}
</style>
