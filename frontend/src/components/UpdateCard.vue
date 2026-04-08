<template>
  <div class="update-card markdown-content" :class="{ 'current-update': highlighted }">
    <div v-if="status === 'success' && content" v-html="renderedContent"></div>
    <div v-else-if="status === 'loading'" class="state-message">{{ t('common.loading') }}</div>
    <div v-else-if="status === 'error'" class="state-message state-message--error">{{ errorMessage || t('common.loadFailed') }}</div>
    <div v-else class="state-message">{{ emptyMessage || t('common.noContent') }}</div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import MarkdownIt from 'markdown-it'
import { useI18n } from '@/utils/i18n'

const props = defineProps({
  /**
   * Markdown 原始内容。
   * 仅在状态为 success 时参与渲染，避免用空字符串混淆“加载中”和“暂无内容”。
   */
  content: {
    type: String,
    default: ''
  },
  /**
   * 卡片当前展示状态。
   * - loading: 正在加载
   * - success: 加载成功，渲染 Markdown
   * - error: 加载失败，显示错误提示
   * - empty: 没有内容可展示
   */
  status: {
    type: String,
    default: 'success'
  },
  /**
   * 错误态文案。
   * 由上层传入，避免组件内部拼接业务相关错误描述。
   */
  errorMessage: {
    type: String,
    default: ''
  },
  /**
   * 空态文案。
   * 当内容为空且不是加载中时，用于向用户明确表达“暂无内容”。
   */
  emptyMessage: {
    type: String,
    default: ''
  },
  highlighted: {
    type: Boolean,
    default: false
  }
})

const { t } = useI18n()

/**
 * Markdown 渲染器实例。
 *
 * 安全策略：
 * 1. 显式关闭原生 HTML 渲染，避免 changelog 中混入的 HTML 片段直接进入 DOM
 * 2. 保留 linkify / typographer，继续支持常规 Markdown 展示体验
 *
 * 说明：
 * 当前 about 文档来自仓库内置的静态 Markdown，关闭 html 不会影响现有文档展示，
 * 但可以为未来文档来源变化提供更稳妥的默认安全边界。
 */
const md = new MarkdownIt({
  html: false,
  linkify: true,
  typographer: true
})

/**
 * 将 Markdown 文本转换为安全性更高的 HTML 片段。
 *
 * @returns {string} 返回供 `v-html` 渲染的 HTML 字符串
 */
const renderedContent = computed(() => {
  return md.render(props.content)
})
</script>

<style scoped>
.update-card {
  background: var(--card-bg);
  padding: 1rem;
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-md);
  margin-bottom: 0.5rem;
  transition: transform 0.2s, box-shadow 0.2s;
}

.update-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
}

.current-update {
  border-left: 4px solid var(--accent-color);
}

/* 统一状态文案样式，避免 loading / error / empty 各自分散维护。 */
.state-message {
  color: var(--text-secondary);
  text-align: center;
  padding: 2rem;
  font-style: italic;
}

.state-message--error {
  color: var(--danger-color);
}

/* Markdown 内容样式 */
.markdown-content {
  line-height: 1.8;
  color: var(--text-secondary);
}

.markdown-content :deep(h4) {
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--accent-color);
  background: var(--accent-color-bg);
  padding: 0.25rem 0.75rem;
  border-radius: var(--radius-xs);
  margin-bottom: 1rem;
  display: inline-block;
}

.markdown-content :deep(ul) {
  list-style: none;
  padding: 0;
  margin: 0;
}

.markdown-content :deep(li) {
  padding: 0.5rem 0;
  color: var(--text-secondary);
  border-bottom: 1px solid var(--border-color);
}

.markdown-content :deep(li:last-child) {
  border-bottom: none;
}

.markdown-content :deep(strong) {
  color: var(--text-primary);
  font-weight: 600;
}

.markdown-content :deep(code) {
  background: var(--accent-color-bg);
  color: var(--accent-color);
  padding: 0.2rem 0.4rem;
  border-radius: var(--radius-xs);
  font-size: 0.9em;
}

.markdown-content :deep(p) {
  margin: 0.5rem 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .markdown-content :deep(h4) {
    font-size: 1rem;
  }
}
</style>
