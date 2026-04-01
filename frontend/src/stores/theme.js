import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useThemeStore = defineStore('theme', () => {
  /**
   * 当前主题状态。
   * 优先读取 localStorage 中用户上次手动选择的主题；
   * 如果用户是首次访问或本地没有缓存，则默认使用 dark 黑色主题。
   * @type {import('vue').Ref<'light' | 'dark'>}
   */
  const theme = ref(localStorage.getItem('theme') || 'dark')

  /**
   * 在亮色和暗色主题之间切换。
   * 每次切换后都会立即同步到 localStorage，并重新应用到页面根节点。
   * @returns {void}
   */
  const toggleTheme = () => {
    theme.value = theme.value === 'light' ? 'dark' : 'light'
    localStorage.setItem('theme', theme.value)
    applyTheme()
  }

  /**
   * 将当前主题写入到 document 根节点，供全局 CSS 变量选择器使用。
   * 这里统一通过 data-theme 属性驱动样式切换，避免组件各自维护主题状态。
   * @returns {void}
   */
  const applyTheme = () => {
    document.documentElement.setAttribute('data-theme', theme.value)
  }

  /**
   * 显式设置指定主题。
   * @param {'light' | 'dark'} newTheme - 目标主题名称
   * @returns {void}
   */
  const setTheme = (newTheme) => {
    theme.value = newTheme
    localStorage.setItem('theme', theme.value)
    applyTheme()
  }

  // 初始化时立即应用主题，确保页面首屏就使用正确的配色。
  applyTheme()

  return {
    theme,
    toggleTheme,
    setTheme,
    applyTheme
  }
})
