<template>
  <div class="about">
    <h1 class="about-title">关于项目</h1>

    <!-- 错误提示 -->
    <div v-if="loadError" class="error-message">
      <h3>⚠️ 加载失败</h3>
      <pre>{{ loadError }}</pre>
    </div>

    <!-- 内容简介 -->
    <section class="section">
      <h2 class="section-title">内容简介</h2>
      <UpdateCard :content="introductionContent" highlighted />
    </section>

    <!-- 当前更新 -->
    <section class="section">
      <h2 class="section-title">当前更新</h2>
      <UpdateCard :content="latestUpdate" highlighted />
    </section>

    <!-- 以前更新 -->
    <section class="section">
      <h2 class="section-title" id="history-title">历史更新</h2>
      <div class="update-history">
        <UpdateCard v-for="(update, index) in paginatedUpdates" :key="index" :content="update" />
      </div>

      <!-- 分页控件 -->
      <div v-if="totalPages > 1" class="pagination">
        <AppButton
          @click="goToPage(1)"
          :disabled="currentPage === 1"
          class="pagination-btn"
          title="第一页"
          preserve-style
        >
          <template #icon>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <polyline points="11 17 6 12 11 7"></polyline>
              <polyline points="18 17 13 12 18 7"></polyline>
            </svg>
          </template>
        </AppButton>

        <AppButton
          @click="goToPage(currentPage - 1)"
          :disabled="currentPage === 1"
          class="pagination-btn"
          title="上一页"
          preserve-style
        >
          <template #icon>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <polyline points="15 18 9 12 15 6"></polyline>
            </svg>
          </template>
        </AppButton>

        <div class="pagination-info">
          {{ currentPage }} / {{ totalPages }}
        </div>

        <AppButton
          @click="goToPage(currentPage + 1)"
          :disabled="currentPage === totalPages"
          class="pagination-btn"
          title="下一页"
          preserve-style
        >
          <template #icon>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <polyline points="9 18 15 12 9 6"></polyline>
            </svg>
          </template>
        </AppButton>

        <AppButton
          @click="goToPage(totalPages)"
          :disabled="currentPage === totalPages"
          class="pagination-btn"
          title="最后一页"
          preserve-style
        >
          <template #icon>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <polyline points="13 17 18 12 13 7"></polyline>
              <polyline points="6 17 11 12 6 7"></polyline>
            </svg>
          </template>
        </AppButton>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import AppButton from '@/components/AppButton.vue'
import UpdateCard from '@/components/UpdateCard.vue'

const introductionContent = ref('')
const latestUpdate = ref('')
const historyUpdates = ref([])
const loadError = ref('')

// 分页状态
const currentPage = ref(1)
const pageSize = 10

// 计算总页数
const totalPages = computed(() => {
  return Math.ceil(historyUpdates.value.length / pageSize)
})

// 计算当前页显示的更新
const paginatedUpdates = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  const end = start + pageSize
  return historyUpdates.value.slice(start, end)
})

// 滚动到历史更新区域
const scrollToHistory = () => {
  const historyTitle = document.getElementById('history-title')
  if (historyTitle) {
    // 获取元素相对于文档顶部的位置
    const elementPosition = historyTitle.getBoundingClientRect().top + window.pageYOffset
    // 减去偏移量(预留顶部空间)
    const offsetPosition = elementPosition - 100

    // 一次性滚动到目标位置
    window.scrollTo({
      top: offsetPosition,
      behavior: 'smooth'
    })
  }
}

// 分页处理函数
const goToPage = (page) => {
  currentPage.value = page
  scrollToHistory()
}

// 版本号比较函数
const compareVersions = (a, b) => {
  const parseVersion = (version) => {
    // 移除 'v' 前缀,处理 beta 等标识
    let cleanVersion = version.replace(/^v/i, '').toLowerCase()

    // 分离版本号和预发布标识
    const prereleaseMatch = cleanVersion.match(/(alpha|beta|rc)(\d*)$/)
    let prerelease = null
    let prereleaseNum = 0

    if (prereleaseMatch) {
      prerelease = prereleaseMatch[1]
      // 提取预发布版本号后的数字
      prereleaseNum = prereleaseMatch[2] ? parseInt(prereleaseMatch[2]) : 0
      // 移除预发布部分,保留主版本号
      cleanVersion = cleanVersion.replace(/(alpha|beta|rc)\d*$/i, '')
    }

    const parts = cleanVersion.split('.').map(num => {
      const parsed = parseInt(num)
      return isNaN(parsed) ? 0 : parsed
    })

    // 补齐版本号到3位
    while (parts.length < 3) {
      parts.push(0)
    }

    return { parts, prerelease, prereleaseNum, version: cleanVersion }
  }

  const verA = parseVersion(a)
  const verB = parseVersion(b)

  // 比较主版本号
  for (let i = 0; i < 3; i++) {
    if (verA.parts[i] !== verB.parts[i]) {
      return verB.parts[i] - verA.parts[i] // 降序
    }
  }

  // 如果主版本号相同,比较预发布标识
  const prereleaseOrder = { 'rc': 1, 'beta': 2, 'alpha': 3, 'undefined': 0 }
  const orderA = verA.prerelease ? prereleaseOrder[verA.prerelease] : 0
  const orderB = verB.prerelease ? prereleaseOrder[verB.prerelease] : 0

  if (orderA !== orderB) {
    return orderA - orderB
  }

  // 如果预发布标识也相同,比较预发布版本号
  return verB.prereleaseNum - verA.prereleaseNum // 降序
}

// 获取 API 基础 URL
const getApiUrl = () => {
  // 如果配置了环境变量,使用环境变量（开发环境）
  if (import.meta.env.VITE_API_URL) {
    return import.meta.env.VITE_API_URL
  }

  // 生产环境：使用相对路径（前后端同服务）
  return ''
}

// 加载 markdown 文件
const loadMarkdownFile = async (filename) => {
  try {
    const apiUrl = getApiUrl()
    const url = `${apiUrl}/api/changelog/file/${filename}`
    console.log('加载文件使用的 API:', url)
    console.log('环境变量 VITE_API_URL:', import.meta.env.VITE_API_URL)

    const response = await fetch(url)
    console.log('响应状态:', response.status, response.statusText)

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`)
    }

    const data = await response.json()
    console.log('响应数据:', data)
    return data.data.content
  } catch (error) {
    console.error(`Error loading ${filename}:`, error)
    console.error('错误堆栈:', error.stack)
    return null
  }
}

// 从文件名提取版本号
const extractVersionFromFilename = (filename) => {
  // 匹配文件名中的版本号,如 v1.0.0.md, v1.0.0beta1.md 等
  const match = filename.match(/^v?([\d.]+(?:alpha|beta|rc)?\d*)/i)
  return match ? match[1] : null
}

// 加载所有更新日志
const loadChangelog = async () => {
  try {
    console.log('开始加载更新日志...')
    console.log('当前页面地址:', window.location.href)
    loadError.value = ''

    // 加载简介
    const intro = await loadMarkdownFile('简介.md')
    console.log('简介内容:', intro ? '加载成功' : '加载失败')
    if (intro) {
      introductionContent.value = intro
    }

    // 从后端 API 获取版本文件列表
    const apiUrl = getApiUrl()
    console.log('使用 API 地址:', apiUrl)
    console.log('完整请求 URL:', `${apiUrl}/api/changelog/files`)

    const versionsResponse = await fetch(`${apiUrl}/api/changelog/files`)
    console.log('响应状态:', versionsResponse.status, versionsResponse.statusText)

    if (!versionsResponse.ok) {
      throw new Error(`无法连接到后端服务器 (${apiUrl}) - 状态码: ${versionsResponse.status}`)
    }
    const data = await versionsResponse.json()
    const updateFiles = data.data.files || []
    console.log('从API获取的版本文件列表:', updateFiles)

    // 过滤出有效的更新文件(排除简介等特殊文件)
    const updates = updateFiles
      .map(filename => {
        const version = extractVersionFromFilename(filename)
        return version ? { version, filename } : null
      })
      .filter(item => item !== null)

    console.log('找到的更新文件:', updates)

    // 按版本号排序
    updates.sort((a, b) => compareVersions(a.version, b.version))
    console.log('排序后的更新列表:', updates)

    // 加载所有更新的内容
    const updateContents = await Promise.all(
      updates.map(async (update) => {
        const content = await loadMarkdownFile(update.filename)
        console.log(`加载 ${update.filename}:`, content ? '成功' : '失败')
        return { content, filename: update.filename }
      })
    )

    // 过滤掉加载失败的内容
    const validUpdates = updateContents.filter(item => item.content !== null)
    console.log('有效的更新数量:', validUpdates.length)

    // 第一个作为最新更新
    if (validUpdates.length > 0) {
      latestUpdate.value = validUpdates[0].content
      historyUpdates.value = validUpdates.slice(1).map(item => item.content)
      console.log('最新更新已设置,历史更新数量:', historyUpdates.value.length)
    }
  } catch (error) {
    console.error('Error loading changelog:', error)
    const apiUrl = getApiUrl()
    loadError.value = `加载失败: ${error.message}\n\n调试信息:\n- API 地址: ${apiUrl}\n- 当前页面: ${window.location.href}\n- 主机名: ${window.location.hostname}\n\n请检查:\n1. 后端服务器是否在 ${apiUrl} 运行\n2. 浏览器控制台查看详细错误(F12)\n3. 确保后端和前端在同一台机器上`
  }
}

onMounted(() => {
  console.log('===== AboutView 组件已挂载 =====')
  console.log('环境变量检查:')
  console.log('  import.meta.env.VITE_API_URL:', import.meta.env.VITE_API_URL)
  console.log('  window.location.hostname:', window.location.hostname)
  console.log('  window.location.protocol:', window.location.protocol)
  console.log('  getApiUrl():', getApiUrl())
  console.log('==============================')
  loadChangelog()
})
</script>

<style scoped>
.about {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.about-title {
  font-size: 2rem;
  font-weight: 700;
  margin-bottom: 2rem;
  color: var(--text-primary);
}

.section {
  margin-bottom: 2.5rem;
}

.section-title {
  font-size: 1.5rem;
  font-weight: 600;
  margin-bottom: 1rem;
  color: var(--text-primary);
  border-left: 4px solid var(--accent-color);
  padding-left: 1rem;
}

.update-history {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

/* 错误提示样式 */
.error-message {
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 8px;
  padding: 1.5rem;
  margin-bottom: 2rem;
}

.error-message h3 {
  color: #ef4444;
  margin: 0 0 1rem 0;
  font-size: 1.2rem;
}

.error-message pre {
  color: var(--text-secondary);
  white-space: pre-wrap;
  word-wrap: break-word;
  margin: 0;
  font-family: monospace;
  font-size: 0.9rem;
  line-height: 1.6;
}

/* 分页控件样式 */
.pagination {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 0.75rem;
  margin-top: 1rem;
}

.pagination-btn {
  position: relative;
  background: linear-gradient(135deg, var(--accent-color) 0%, var(--accent-color) 100%);
  color: white;
  border: none;
  padding: 0.4rem 0.6rem;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  line-height: 1;
  box-shadow: 0 2px 6px rgba(99, 102, 241, 0.25);
  overflow: hidden;
  min-width: 2rem;
  height: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.pagination-btn svg {
  width: 1.125rem;
  height: 1.125rem;
}

.pagination-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s;
}

.pagination-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 3px 8px rgba(99, 102, 241, 0.35);
}

.pagination-btn:hover:not(:disabled)::before {
  left: 100%;
}

.pagination-btn:active:not(:disabled) {
  transform: translateY(0);
  box-shadow: 0 1px 3px rgba(99, 102, 241, 0.2);
}

.pagination-btn:disabled {
  background: linear-gradient(135deg, var(--border-color) 0%, var(--border-color) 100%);
  color: var(--text-secondary);
  cursor: not-allowed;
  box-shadow: none;
  opacity: 0.6;
  transform: none;
}

.pagination-info {
  position: relative;
  color: var(--text-primary);
  font-size: 0.8125rem;
  font-weight: 500;
  min-width: 3.5rem;
  text-align: center;
  padding: 0.35rem 0.75rem;
  background: var(--accent-color-bg);
  border-radius: 6px;
}

.pagination-info::before {
  content: none;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .about {
    padding: 1rem;
  }

  .about-title {
    font-size: 1.5rem;
    margin-bottom: 1.5rem;
  }

  .section-title {
    font-size: 1.25rem;
  }

  .pagination {
    flex-wrap: wrap;
    gap: 0.5rem;
  }

  .pagination-btn {
    padding: 0.35rem 0.5rem;
    min-width: 1.75rem;
    height: 1.75rem;
  }

  .pagination-btn svg {
    width: 1rem;
    height: 1rem;
  }

  .pagination-info {
    font-size: 0.75rem;
    min-width: 3rem;
    padding: 0.3rem 0.625rem;
  }
}
</style>
