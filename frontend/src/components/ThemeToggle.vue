<template>
  <button class="theme-toggle" @click="toggleTheme" :aria-label="themeAriaLabel" :title="themeAriaLabel">
    <svg v-if="theme === 'light'" class="icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="M12 3V4M12 20V21M4 12H3M21 12H20M5.636 5.636L4.929 4.929M19.071 19.071L18.364 18.364M5.636 18.364L4.929 19.071M19.071 4.929L18.364 5.636" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
      <circle cx="12" cy="12" r="4" stroke="currentColor" stroke-width="2"/>
    </svg>
    <svg v-else class="icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
    </svg>
  </button>
</template>

<script setup>
import { computed } from 'vue'
import { useThemeStore } from '@/stores/theme'
import { storeToRefs } from 'pinia'
import { useI18n } from '@/utils/i18n'

const themeStore = useThemeStore()
const { theme } = storeToRefs(themeStore)
const { t } = useI18n()

/**
 * 根据当前主题返回切换按钮的可访问性文案。
 * 当前是亮色时提示“切到暗色”，反之提示“切到亮色”。
 */
const themeAriaLabel = computed(() => {
  return theme.value === 'light'
    ? t('theme.switchToDark')
    : t('theme.switchToLight')
})

const toggleTheme = () => {
  themeStore.toggleTheme()
}
</script>

<style scoped>
.theme-toggle {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  padding: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s;
  color: var(--header-text);
}

.theme-toggle:hover {
  background: var(--bg-tertiary);
  border-color: var(--accent-color);
  transform: scale(1.05);
  box-shadow: var(--shadow-sm);
}

.theme-toggle:active {
  transform: scale(0.95);
}

.icon {
  width: 20px;
  height: 20px;
}
</style>
