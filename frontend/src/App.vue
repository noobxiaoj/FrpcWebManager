<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { RouterLink, RouterView } from 'vue-router'
import { getApiBaseUrl } from '@/utils/api'
import { syncLanguageFromSettings, useI18n } from '@/utils/i18n'
import ThemeToggle from '@/components/ThemeToggle.vue'
import LanguageToggle from '@/components/LanguageToggle.vue'
import PasswordLoginGate from '@/components/PasswordLoginGate.vue'

const mobileMenuOpen = ref(false)
const authLoading = ref(true)
const loginLoading = ref(false)
const loginErrorMessage = ref('')
const authState = ref({
  passwordAuth: {
    enabled: false,
    username: ''
  },
  authenticated: true
})
const appSettings = ref({
  language: 'zh-CN',
  navigationBar: {
    showAboutButton: true,
    showLockButton: true,
    showLanguageButton: true,
    showThemeButton: true
  }
})

const toggleMobileMenu = () => {
  mobileMenuOpen.value = !mobileMenuOpen.value
}

const isPasswordProtected = computed(() => authState.value.passwordAuth?.enabled)
const isAuthenticated = computed(() => authState.value.authenticated)
const showPasswordGate = computed(() => !authLoading.value && isPasswordProtected.value && !isAuthenticated.value)
const { t } = useI18n()
const navigationBarSettings = computed(() => appSettings.value.navigationBar || {
  showAboutButton: true,
  showLockButton: true,
  showLanguageButton: true,
  showThemeButton: true
})

/**
 * 统一处理“登录态已失效”事件。
 * 当任意业务请求返回未授权时，根组件会重新拉取认证状态并切回密码页。
 */
const handleAuthExpired = async () => {
  await loadAuthStatus()
}

/**
 * 获取当前浏览器的认证状态。
 * 页面首次加载、登录成功、退出登录后都会复用这个方法统一刷新界面。
 */
const loadAuthStatus = async () => {
  authLoading.value = true

  try {
    const response = await fetch(`${getApiBaseUrl()}/api/auth/status`, {
      credentials: 'include'
    })
    const result = await response.json()

    if (result.code !== 0) {
      throw new Error(result.message || t('app.auth.loadStatusFailed'))
    }

    authState.value = {
      passwordAuth: result.data.passwordAuth || {
        enabled: false,
        username: ''
      },
      authenticated: Boolean(result.data.authenticated)
    }
  } catch (error) {
    console.error('获取认证状态失败:', error)
    loginErrorMessage.value = t('app.auth.initFailed')
  } finally {
    authLoading.value = false
  }
}

/**
 * 提交登录表单并写入浏览器会话 Cookie。
 * 成功后重新读取认证状态，让应用恢复原本页面内容。
 */
const handleLogin = async ({ username, password }) => {
  loginLoading.value = true
  loginErrorMessage.value = ''

  try {
    const response = await fetch(`${getApiBaseUrl()}/api/auth/login`, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        username,
        password
      })
    })
    const result = await response.json()

    if (result.code !== 0) {
      throw new Error(result.message || t('app.auth.loginFailed'))
    }

    await loadAuthStatus()
    await loadAppSettings()
  } catch (error) {
    console.error('登录失败:', error)
    loginErrorMessage.value = error.message || t('app.auth.loginFallbackError')
  } finally {
    loginLoading.value = false
  }
}

/**
 * 读取后端保存的界面设置，并同步到根组件状态。
 * 这样页面刷新、重新登录后都能恢复语言与顶部导航按钮显示状态。
 *
 * @returns {Promise<void>}
 */
const loadAppSettings = async () => {
  try {
    const response = await fetch(`${getApiBaseUrl()}/api/settings`, {
      credentials: 'include'
    })

    if (!response.ok) {
      return
    }

    const result = await response.json()
    applyAppSettings(result.data)
  } catch (error) {
    console.error('获取应用设置失败:', error)
  }
}

/**
 * 将后端返回的设置对象同步到根组件中。
 * 目前这里主要消费语言与导航栏按钮显隐配置。
 *
 * @param {Object} settings - 后端返回的系统设置
 */
const applyAppSettings = (settings = {}) => {
  appSettings.value = {
    language: settings.language || 'zh-CN',
    navigationBar: {
      showAboutButton: settings.navigationBar?.showAboutButton ?? true,
      showLockButton: settings.navigationBar?.showLockButton ?? true,
      showLanguageButton: settings.navigationBar?.showLanguageButton ?? true,
      showThemeButton: settings.navigationBar?.showThemeButton ?? true
    }
  }

  syncLanguageFromSettings(appSettings.value.language)
}

/**
 * 响应设置页保存后的本地广播事件。
 * 这样顶部导航按钮可以在保存后立即切换，不需要手动刷新页面。
 *
 * @param {CustomEvent} event - 携带最新设置数据的浏览器事件
 */
const handleAppSettingsUpdated = (event) => {
  applyAppSettings(event.detail || {})
}

/**
 * 主动退出当前浏览器的登录态。
 * 退出后如果系统仍启用了密码，将重新回到密码验证页面。
 *
 * @returns {Promise<void>}
 */
const handleLogout = async () => {
  try {
    await fetch(`${getApiBaseUrl()}/api/auth/logout`, {
      method: 'POST',
      credentials: 'include'
    })
  } catch (error) {
    console.error('退出登录失败:', error)
  } finally {
    await loadAuthStatus()
  }
}

onMounted(async () => {
  await loadAuthStatus()

  if (!showPasswordGate.value) {
    await loadAppSettings()
  }

  /**
   * 当业务接口返回“未登录”时，统一切回密码验证页面。
   * 这样即使会话失效，也不会让界面停留在不可操作状态。
   */
  window.addEventListener('auth-expired', handleAuthExpired)
  window.addEventListener('app-settings-updated', handleAppSettingsUpdated)
})

onUnmounted(() => {
  window.removeEventListener('auth-expired', handleAuthExpired)
  window.removeEventListener('app-settings-updated', handleAppSettingsUpdated)
})
</script>

<template>
  <div v-if="authLoading" class="auth-loading-screen">
    <div class="auth-loading-card">
      <p class="auth-loading-label">{{ t('app.auth.checkingTitle') }}</p>
      <p class="auth-loading-text">{{ t('app.auth.checkingText') }}</p>
    </div>
  </div>

  <PasswordLoginGate
    v-else-if="showPasswordGate"
    :loading="loginLoading"
    :error-message="loginErrorMessage"
    :default-username="authState.passwordAuth.username"
    @submit="handleLogin"
  />

  <div v-else class="app">
    <header class="top-header">
      <div class="header-container">
        <!-- 移动端菜单按钮 -->
        <button class="mobile-menu-btn" @click="toggleMobileMenu" :aria-label="t('app.auth.mobileMenu')">
          <span></span>
          <span></span>
          <span></span>
        </button>

        <!-- 左侧组：Logo + 导航 -->
        <div class="left-group">
          <!-- Logo / 标题 -->
          <div class="logo">
            <h1>{{ t('app.name') }}</h1>
          </div>

          <!-- 桌面端导航 -->
          <nav class="desktop-nav">
            <RouterLink to="/">{{ t('app.nav.servers') }}</RouterLink>
            <RouterLink to="/tasks">{{ t('app.nav.tasks') }}</RouterLink>
            <RouterLink to="/settings">{{ t('app.nav.settings') }}</RouterLink>
            <RouterLink v-if="navigationBarSettings.showAboutButton" to="/about">{{ t('app.nav.about') }}</RouterLink>
          </nav>
        </div>

        <!-- 顶部功能按钮：语言切换 + 主题切换 -->
        <div class="theme-wrapper">
          <button
            v-if="isPasswordProtected && navigationBarSettings.showLockButton"
            class="logout-button"
            @click="handleLogout"
            :aria-label="t('app.auth.logout')"
            :title="t('app.auth.logout')"
          >
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
              <path
                d="M8 10V7.5C8 5.01472 10.0147 3 12.5 3C14.9853 3 17 5.01472 17 7.5V10"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
              />
              <rect
                x="5"
                y="10"
                width="15"
                height="11"
                rx="2"
                stroke="currentColor"
                stroke-width="2"
              />
              <circle cx="12.5" cy="15.5" r="1.2" fill="currentColor" />
            </svg>
          </button>
          <LanguageToggle v-if="navigationBarSettings.showLanguageButton" />
          <ThemeToggle v-if="navigationBarSettings.showThemeButton" />
        </div>
      </div>

      <!-- 移动端导航菜单 -->
      <nav class="mobile-nav" :class="{ open: mobileMenuOpen }">
        <RouterLink to="/" @click="mobileMenuOpen = false">{{ t('app.nav.servers') }}</RouterLink>
        <RouterLink to="/tasks" @click="mobileMenuOpen = false">{{ t('app.nav.tasks') }}</RouterLink>
        <RouterLink to="/settings" @click="mobileMenuOpen = false">{{ t('app.nav.settings') }}</RouterLink>
        <RouterLink v-if="navigationBarSettings.showAboutButton" to="/about" @click="mobileMenuOpen = false">{{ t('app.nav.about') }}</RouterLink>
      </nav>
    </header>

    <main class="main-content">
      <RouterView />
    </main>
  </div>
</template>

<style scoped>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

.app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: var(--bg-primary);
  color: var(--text-primary);
  transition: background-color 0.3s, color 0.3s;
}

.auth-loading-screen {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  background:
    radial-gradient(circle at top left, var(--accent-color-bg), transparent 28%),
    linear-gradient(180deg, var(--bg-secondary), var(--bg-primary));
}

.auth-loading-card {
  min-width: min(100%, 360px);
  padding: 1.5rem;
  border-radius: var(--radius-xl);
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-lg);
}

.auth-loading-label {
  color: var(--accent-color);
  font-size: 0.875rem;
  font-weight: 700;
}

.auth-loading-text {
  margin-top: 0.75rem;
  color: var(--text-secondary);
}

/* 顶部导航栏 */
.top-header {
  background: var(--header-bg);
  box-shadow: var(--shadow-md);
  position: sticky;
  top: 0;
  z-index: 100;
  width: 100%;
  transition: background-color 0.3s;
}

.header-container {
  max-width: 100%;
  margin: 0 auto;
  padding: 1rem 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

@media (max-width: 768px) {
  .header-container {
    justify-content: space-between;
  }
}

/* 左侧组：Logo + 导航 */
.left-group {
  display: flex;
  align-items: center;
  gap: 20px;
}

/* Logo */
.logo h1 {
  color: var(--header-text);
  font-size: 1.2rem;
  font-weight: 600;
}

/* 主题切换按钮容器 */
.theme-wrapper {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.logout-button {
  width: 40px;
  min-width: 40px;
  height: 40px;
  padding: 0;
  border-radius: var(--radius-md);
  background: transparent;
  color: var(--header-text);
  display: flex;
  align-items: center;
  justify-content: center;
}

.logout-button:hover {
  background: var(--accent-color-bg);
}

.logout-button svg {
  width: 18px;
  height: 18px;
}

/* 桌面端导航 */
.desktop-nav {
  display: flex;
  /* 收紧顶部导航项之间的间距，让“服务器 / 任务 / 设置 / 关于”更紧凑。 */
  gap: 12px;
  align-items: center;
}

.desktop-nav a {
  color: var(--header-text);
  text-decoration: none;
  font-weight: 500;
  /* 略微减小横向内边距，避免导航项之间显得过于松散。 */
  padding: 0.5rem 0.75rem;
  border-radius: var(--radius-xs);
  transition: all 0.3s;
  position: relative;
}

.desktop-nav a:hover,
.desktop-nav a.router-link-exact-active {
  color: var(--accent-color);
}

.desktop-nav a.router-link-exact-active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 80%;
  height: 2px;
  background: var(--accent-color);
  border-radius: 1px;
}

/* 移动端菜单按钮 */
.mobile-menu-btn {
  display: none;
  flex-direction: column;
  justify-content: space-around;
  width: 30px;
  height: 30px;
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 0;
}

.mobile-menu-btn span {
  width: 30px;
  height: 3px;
  background: var(--header-text);
  border-radius: 3px;
  transition: all 0.3s;
}

/* 移动端导航 */
.mobile-nav {
  display: none;
  flex-direction: column;
  padding: 0 2rem;
  max-height: 0;
  overflow: hidden;
  transition: max-height 0.3s ease-out;
  background: var(--header-bg);
}

.mobile-nav.open {
  display: flex;
  max-height: 300px;
  padding: 1rem 2rem;
}

.mobile-nav a {
  color: var(--header-text);
  text-decoration: none;
  font-weight: 500;
  padding: 0.75rem 0;
  border-bottom: 1px solid var(--border-color);
  transition: all 0.3s;
}

.mobile-nav a:last-child {
  border-bottom: none;
}

.mobile-nav a:hover,
.mobile-nav a.router-link-exact-active {
  color: var(--accent-color);
}

/* 主内容区域 */
.main-content {
  flex: 1;
  padding: 2rem;
  max-width: 100%;
}

/* 响应式设计 - 平板 */
@media (max-width: 768px) {
  .desktop-nav {
    display: none;
  }

  .mobile-menu-btn {
    display: flex;
  }

  .left-group {
    flex: 1;
    margin-left: 0.5em;
  }

  .header-container {
    padding: 1rem;
  }

  .main-content {
    padding: 1rem;
  }
}

/* 响应式设计 - 手机 */
@media (max-width: 480px) {
  .logo h1 {
    font-size: 1.2rem;
  }

  .header-container {
    padding: 0.75rem 1rem;
  }
}
</style>
