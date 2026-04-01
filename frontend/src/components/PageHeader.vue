<script setup>
/**
 * 通用页面头部组件
 *
 * 功能：
 * 1. 统一页面标题区域的布局、间距与分隔线样式
 * 2. 通过 `actions` 具名插槽承载右侧按钮区域
 * 3. 支持可选的 `description` 描述文案，未传入时仅渲染标题
 *
 * @param {string} title - 页面标题
 * @param {string} [description=''] - 页面描述，可选
 *
 * 使用示例：
 * <PageHeader title="任务列表">
 *   <template #actions>
 *     <AppButton>创建任务</AppButton>
 *   </template>
 * </PageHeader>
 */
defineProps({
  title: {
    type: String,
    required: true
  },
  description: {
    type: String,
    default: ''
  }
})
</script>

<template>
  <header class="page-header">
    <div class="page-header-left">
      <h1 class="page-title">{{ title }}</h1>
      <p v-if="description" class="page-description">{{ description }}</p>
    </div>

    <!--
      操作区使用具名插槽承载，便于不同页面注入刷新、保存、新建等动作按钮。
      如果页面没有传入 actions 插槽，则不渲染右侧容器，避免产生空白占位。
    -->
    <div v-if="$slots.actions" class="page-header-actions">
      <slot name="actions" />
    </div>
  </header>
</template>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 1rem;
  margin-bottom: 2rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid var(--border-color);
}

.page-header-left {
  min-width: 0;
}

.page-title {
  font-size: 2rem;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
  line-height: 1.2;
}

.page-description {
  margin: 0.35rem 0 0;
  color: var(--text-secondary);
  font-size: 0.95rem;
}

.page-header-actions {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  flex-wrap: wrap;
  justify-content: flex-end;
}

/*
 * 统一页头操作按钮的尺寸基线。
 * 这里使用 :deep 作用到插槽中的按钮，并通过更高的选择器优先级统一高度、字号和图标大小，
 * 避免各页面的局部样式继续把按钮拉成不同尺寸。
 */
.page-header-actions :deep(.page-header-action-button) {
  min-height: var(--page-header-button-height) !important;
  padding-top: 0.625rem !important;
  padding-bottom: 0.625rem !important;
  padding-left: var(--page-header-button-padding-x) !important;
  padding-right: var(--page-header-button-padding-x) !important;
  font-size: 0.95rem !important;
  line-height: 1 !important;
  white-space: nowrap;
}

.page-header-actions :deep(.page-header-action-button--icon) {
  width: var(--page-header-button-height) !important;
  min-width: var(--page-header-button-height) !important;
  height: var(--page-header-button-height) !important;
  padding: 0 !important;
}

.page-header-actions :deep(.page-header-action-button svg),
.page-header-actions :deep(.page-header-action-button .spinner),
.page-header-actions :deep(.page-header-action-button .app-button__spinner) {
  width: var(--page-header-button-icon-size) !important;
  height: var(--page-header-button-icon-size) !important;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }

  .page-title {
    font-size: 1.5rem;
  }

  .page-header-actions {
    width: 100%;
    justify-content: flex-start;
  }

  .page-header-actions :deep(.page-header-action-button) {
    min-height: 2.5rem !important;
    padding-top: 0.55rem !important;
    padding-bottom: 0.55rem !important;
    padding-left: 1rem !important;
    padding-right: 1rem !important;
    font-size: 0.9rem !important;
  }

  .page-header-actions :deep(.page-header-action-button--icon) {
    width: 2.5rem !important;
    min-width: 2.5rem !important;
    height: 2.5rem !important;
  }
}
</style>
