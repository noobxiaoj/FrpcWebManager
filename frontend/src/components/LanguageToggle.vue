<template>
  <div ref="menuRef" class="language-toggle-wrapper">
    <button
      class="language-toggle"
      type="button"
      :aria-label="t('language.switchAria', { language: currentLanguageLabel })"
      :title="t('language.current', { language: currentLanguageLabel })"
      @click="toggleMenu"
    >
      <svg class="icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path
          d="M12 3C16.9706 3 21 7.02944 21 12C21 16.9706 16.9706 21 12 21C7.02944 21 3 16.9706 3 12C3 7.02944 7.02944 3 12 3Z"
          stroke="currentColor"
          stroke-width="2"
        />
        <path
          d="M3.6 9H20.4M3.6 15H20.4M12 3C14.2 5.2 15.5 8.4 15.5 12C15.5 15.6 14.2 18.8 12 21M12 3C9.8 5.2 8.5 8.4 8.5 12C8.5 15.6 9.8 18.8 12 21"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
        />
      </svg>
    </button>

    <div v-if="menuOpen" class="language-menu">
      <button
        v-for="option in languageOptions"
        :key="option.value"
        type="button"
        class="language-menu-item"
        :class="{ active: currentLanguage === option.value }"
        :disabled="saving"
        @click="selectLanguage(option.value)"
      >
        <span>{{ t(option.labelKey) }}</span>
        <span v-if="currentLanguage === option.value" class="language-menu-check">✓</span>
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { getApiBaseUrl } from '@/utils/api'
import { languageOptions, setLanguage, syncLanguageFromSettings, useI18n } from '@/utils/i18n'

/**
 * 顶部导航栏语言切换按钮。
 * 当前阶段仅负责读取并保存语言选择到 settings.json，
 * 不负责真正切换页面文案，后续接入 i18n 时可直接复用该组件。
 */
const menuRef = ref(null)
const menuOpen = ref(false)
const saving = ref(false)
const { locale, t } = useI18n()
const currentLanguage = locale

/**
 * 根据当前语言代码返回展示文案。
 *
 * @returns {string} 当前语言对应的按钮提示文本
 */
const currentLanguageLabel = computed(() => {
  const currentOption = languageOptions.find((option) => option.value === currentLanguage.value)
  return currentOption ? t(currentOption.labelKey) : t('language.option.zhCN')
})

/**
 * 读取系统设置中的语言值。
 * 如果后端尚未包含 language 字段，后端归一化逻辑会自动补齐默认值，
 * 这里仍保留前端兜底，避免接口异常时按钮状态为空。
 *
 * @returns {Promise<void>}
 */
const loadLanguage = async () => {
  try {
    const response = await fetch(`${getApiBaseUrl()}/api/settings`, {
      credentials: 'include'
    })

    if (!response.ok) {
      throw new Error(t('language.loadFailed'))
    }

    const result = await response.json()
    syncLanguageFromSettings(result.data?.language)
  } catch (error) {
    console.error('加载语言设置失败:', error)
    setLanguage('zh-CN')
  }
}

/**
 * 切换菜单开关状态。
 *
 * @returns {void}
 */
const toggleMenu = () => {
  menuOpen.value = !menuOpen.value
}

/**
 * 处理点击空白区域关闭下拉菜单。
 *
 * @param {MouseEvent} event - 浏览器点击事件对象
 * @returns {void}
 */
const handleClickOutside = (event) => {
  if (!menuRef.value?.contains(event.target)) {
    menuOpen.value = false
  }
}

/**
 * 将用户选中的语言保存回后端 settings.json。
 * 保存前会先读取完整设置对象，再仅覆盖 language 字段，
 * 这样不会破坏其他设置项，也兼容旧配置自动补齐后的结构。
 *
 * @param {string} language - 目标语言代码，仅允许 zh-CN 或 en-US
 * @returns {Promise<void>}
 */
const selectLanguage = async (language) => {
  if (saving.value || currentLanguage.value === language) {
    menuOpen.value = false
    return
  }

  saving.value = true

  try {
    const getResponse = await fetch(`${getApiBaseUrl()}/api/settings`, {
      credentials: 'include'
    })

    if (!getResponse.ok) {
      throw new Error('获取当前设置失败')
    }

    const getResult = await getResponse.json()
    const nextSettings = {
      ...getResult.data,
      language,
      passwordAuth: {
        enabled: getResult.data?.passwordAuth?.enabled || false,
        username: getResult.data?.passwordAuth?.username || ''
      }
    }

    const saveResponse = await fetch(`${getApiBaseUrl()}/api/settings`, {
      method: 'PUT',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(nextSettings)
    })

    if (!saveResponse.ok) {
      throw new Error('保存语言设置失败')
    }

    setLanguage(language)
    menuOpen.value = false
  } catch (error) {
    console.error('保存语言设置失败:', error)
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadLanguage()
  document.addEventListener('click', handleClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.language-toggle-wrapper {
  position: relative;
}

.language-toggle {
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

.language-toggle:hover {
  background: var(--bg-tertiary);
  border-color: var(--accent-color);
  transform: scale(1.05);
  box-shadow: var(--shadow-sm);
}

.language-toggle:active {
  transform: scale(0.95);
}

.language-menu {
  position: absolute;
  top: calc(100% + 0.5rem);
  right: 0;
  min-width: 9rem;
  padding: 0.35rem;
  border-radius: var(--radius-md);
  border: 1px solid var(--border-color);
  background: var(--card-bg);
  box-shadow: var(--shadow-lg);
  z-index: 120;
}

.language-menu-item {
  width: 100%;
  border: none;
  background: transparent;
  color: var(--text-primary);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.75rem;
  padding: 0.65rem 0.8rem;
  border-radius: var(--radius-sm);
  text-align: left;
}

.language-menu-item:hover:not(:disabled) {
  background: var(--accent-color-bg);
  color: var(--accent-color);
}

.language-menu-item.active {
  background: var(--accent-color-bg);
  color: var(--accent-color);
}

.language-menu-item:disabled {
  opacity: 0.6;
  cursor: wait;
}

.language-menu-check {
  font-size: 0.9rem;
  font-weight: 700;
}

.icon {
  width: 20px;
  height: 20px;
}
</style>
