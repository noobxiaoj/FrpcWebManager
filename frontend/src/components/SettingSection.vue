<template>
  <div class="setting-section">
    <button
      type="button"
      class="section-toggle"
      :aria-expanded="expanded"
      :aria-controls="contentId"
      @click="toggleExpanded"
    >
      <div class="section-toggle-main">
        <h2 class="section-title">{{ title }}</h2>
        <p v-if="summary" class="section-summary">{{ summary }}</p>
      </div>

      <div class="section-toggle-side">
        <span class="toggle-text">{{ expanded ? t('common.collapse') : t('common.expand') }}</span>
        <svg
          class="toggle-icon"
          :class="{ expanded }"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          aria-hidden="true"
        >
          <polyline points="6 9 12 15 18 9"></polyline>
        </svg>
      </div>
    </button>

    <Transition name="section-collapse">
      <div v-show="expanded" :id="contentId" class="section-content">
        <slot></slot>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useI18n } from '@/utils/i18n'

const props = defineProps({
  title: {
    type: String,
    required: true
  },
  summary: {
    type: String,
    default: ''
  },
  sectionId: {
    type: String,
    default: ''
  },
  defaultExpanded: {
    type: Boolean,
    default: false
  }
})

const { t } = useI18n()

/**
 * 控制分组的折叠状态。
 * 默认值来自 defaultExpanded，便于父组件根据页面场景决定初始展开行为。
 */
const expanded = ref(props.defaultExpanded)

/**
 * 生成可访问性所需的内容区域 id。
 * 如果父组件显式传入 sectionId，则优先使用，方便页面内调试与定位；
 * 否则根据标题生成一个稳定的兜底 id。
 *
 * @returns {string} 返回当前折叠内容区域的唯一 id。
 */
const contentId = computed(() => {
  if (props.sectionId) {
    return props.sectionId
  }

  return `settings-section-${props.title.replace(/\s+/g, '-').toLowerCase()}`
})

/**
 * 切换折叠面板的展开状态。
 * 该方法只负责本组件的 UI 展示，不会影响父组件中的设置数据。
 *
 * @returns {void}
 */
const toggleExpanded = () => {
  expanded.value = !expanded.value
}
</script>

<style scoped>
.setting-section {
  background: var(--card-bg);
  border-radius: var(--radius-xl);
  padding: 1.5rem;
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-lg);
}

.section-toggle {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 0;
  border: none;
  background: transparent;
  color: inherit;
  text-align: left;
  cursor: pointer;
  outline: none;
  box-shadow: none;
  -webkit-tap-highlight-color: transparent;
}

.section-toggle:focus,
.section-toggle:focus-visible,
.section-toggle:active {
  outline: none;
  box-shadow: none;
}

.section-toggle-main {
  min-width: 0;
  flex: 1;
}

.section-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.section-summary {
  margin: 0.5rem 0 0 0;
  color: var(--text-secondary);
  font-size: 0.9rem;
  line-height: 1.5;
}

.section-toggle-side {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: var(--text-secondary);
  flex-shrink: 0;
}

.toggle-text {
  font-size: 0.9rem;
  font-weight: 500;
}

.toggle-icon {
  width: 1rem;
  height: 1rem;
  transition: transform 0.25s ease;
}

.toggle-icon.expanded {
  transform: rotate(180deg);
}

.section-content {
  margin-top: 1.25rem;
  padding-top: 0.75rem;
  border-top: 1px solid var(--border-color);
}

.section-collapse-enter-active,
.section-collapse-leave-active {
  transition: opacity 0.25s ease, transform 0.25s ease;
  transform-origin: top;
}

.section-collapse-enter-from,
.section-collapse-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}

@media (max-width: 768px) {
  .section-toggle {
    align-items: flex-start;
  }

  .section-toggle-side {
    padding-top: 0.2rem;
  }
}
</style>
