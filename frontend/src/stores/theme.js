import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useThemeStore = defineStore('theme', () => {
  // 从 localStorage 读取保存的主题，默认为 'light'
  const theme = ref(localStorage.getItem('theme') || 'light')

  // 切换主题
  const toggleTheme = () => {
    theme.value = theme.value === 'light' ? 'dark' : 'light'
    localStorage.setItem('theme', theme.value)
    applyTheme()
  }

  // 应用主题到 document
  const applyTheme = () => {
    document.documentElement.setAttribute('data-theme', theme.value)
  }

  // 设置指定主题
  const setTheme = (newTheme) => {
    theme.value = newTheme
    localStorage.setItem('theme', theme.value)
    applyTheme()
  }

  // 初始化时应用主题
  applyTheme()

  return {
    theme,
    toggleTheme,
    setTheme,
    applyTheme
  }
})
