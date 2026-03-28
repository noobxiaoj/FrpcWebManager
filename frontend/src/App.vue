<script setup>
import { ref } from 'vue'
import { RouterLink, RouterView } from 'vue-router'
import ThemeToggle from '@/components/ThemeToggle.vue'

const mobileMenuOpen = ref(false)

const toggleMobileMenu = () => {
  mobileMenuOpen.value = !mobileMenuOpen.value
}
</script>

<template>
  <div class="app">
    <header class="top-header">
      <div class="header-container">
        <!-- 移动端菜单按钮 -->
        <button class="mobile-menu-btn" @click="toggleMobileMenu" aria-label="菜单">
          <span></span>
          <span></span>
          <span></span>
        </button>

        <!-- 左侧组：Logo + 导航 -->
        <div class="left-group">
          <!-- Logo / 标题 -->
          <div class="logo">
            <h1>FRPC 管理器</h1>
          </div>

          <!-- 桌面端导航 -->
          <nav class="desktop-nav">
            <RouterLink to="/">服务器</RouterLink>
            <RouterLink to="/tasks">任务</RouterLink>
            <RouterLink to="/settings">设置</RouterLink>
            <RouterLink to="/about">关于</RouterLink>
          </nav>
        </div>

        <!-- 主题切换按钮 -->
        <div class="theme-wrapper">
          <ThemeToggle />
        </div>
      </div>

      <!-- 移动端导航菜单 -->
      <nav class="mobile-nav" :class="{ open: mobileMenuOpen }">
        <RouterLink to="/" @click="mobileMenuOpen = false">服务器</RouterLink>
        <RouterLink to="/tasks" @click="mobileMenuOpen = false">任务</RouterLink>
        <RouterLink to="/settings" @click="mobileMenuOpen = false">设置</RouterLink>
        <RouterLink to="/about" @click="mobileMenuOpen = false">关于</RouterLink>
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

/* 顶部导航栏 */
.top-header {
  background: var(--header-bg);
  box-shadow: 0 2px 8px var(--shadow-color);
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
}

/* 桌面端导航 */
.desktop-nav {
  display: flex;
  gap: 20px;
  align-items: center;
}

.desktop-nav a {
  color: var(--header-text);
  text-decoration: none;
  font-weight: 500;
  padding: 0.5rem 1rem;
  border-radius: 4px;
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
