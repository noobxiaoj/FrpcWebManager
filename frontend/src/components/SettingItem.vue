<template>
  <div class="setting-item">
    <div class="setting-info">
      <label class="setting-label" :for="id">
        <slot name="label">{{ label }}</slot>
      </label>
      <p class="setting-description">
        <slot name="description">{{ description }}</slot>
      </p>
      <p v-if="warning" class="setting-warning">
        {{ warning }}
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
defineProps({
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
  font-size: 0.85rem;
  color: #78716c;
  margin: 0.5rem 0 0 0;
  font-weight: 400;
  opacity: 0.8;
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
