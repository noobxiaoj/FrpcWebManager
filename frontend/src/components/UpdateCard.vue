<template>
  <div class="update-card markdown-content" :class="{ 'current-update': highlighted }">
    <div v-if="content" v-html="renderedContent"></div>
    <div v-else class="loading">加载中...</div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import MarkdownIt from 'markdown-it'

const props = defineProps({
  content: {
    type: String,
    default: ''
  },
  highlighted: {
    type: Boolean,
    default: false
  }
})

const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true
})

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

/* 加载状态样式 */
.loading {
  color: var(--text-secondary);
  text-align: center;
  padding: 2rem;
  font-style: italic;
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
