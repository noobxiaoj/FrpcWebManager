<template>
  <div class="about">
    <PageHeader :title="t('about.title')" />

    <!-- 错误提示 -->
    <div v-if="loadError" class="error-message">
      <h3>{{ t('about.loadFailedTitle') }}</h3>
      <pre>{{ loadError }}</pre>
    </div>

    <!-- 内容简介 -->
    <section class="section">
      <h2 class="section-title">{{ t('about.overview') }}</h2>
      <UpdateCard
        :content="introductionContent"
        :status="introductionStatus"
        :error-message="t('about.introLoadFailed')"
        :empty-message="t('about.introEmpty')"
        highlighted
      />
    </section>

    <!-- 当前更新 -->
    <section class="section">
      <h2 class="section-title">{{ t('about.latest') }}</h2>
      <UpdateCard
        :content="latestUpdate"
        :status="latestUpdateStatus"
        :error-message="t('about.latestLoadFailed')"
        :empty-message="t('about.latestEmpty')"
        highlighted
      />
    </section>

    <!-- 以前更新 -->
    <section class="section">
      <h2 class="section-title" id="history-title">{{ t('about.history') }}</h2>
      <div class="update-history">
        <UpdateCard
          v-for="(update, index) in paginatedUpdates"
          :key="index"
          :content="update"
          status="success"
        />
        <UpdateCard
          v-if="historyUpdatesStatus !== 'success' || historyUpdates.length === 0"
          :status="historyUpdatesStatus"
          :error-message="t('about.historyLoadFailed')"
          :empty-message="t('about.historyEmpty')"
        />
      </div>

      <!-- 分页控件 -->
      <div v-if="totalPages > 1" class="pagination">
        <AppButton
          @click="goToPage(1)"
          :disabled="currentPage === 1"
          class="pagination-btn"
          :title="t('common.firstPage')"
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
          :title="t('common.previousPage')"
          preserve-style
        >
          <template #icon>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <polyline points="15 18 9 12 15 6"></polyline>
            </svg>
          </template>
        </AppButton>

        <div class="pagination-info">
          {{ t('common.page', { current: currentPage, total: totalPages }) }}
        </div>

        <AppButton
          @click="goToPage(currentPage + 1)"
          :disabled="currentPage === totalPages"
          class="pagination-btn"
          :title="t('common.nextPage')"
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
          :title="t('common.lastPage')"
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
import { ref, computed, onMounted, watch } from 'vue'
import AppButton from '@/components/AppButton.vue'
import PageHeader from '@/components/PageHeader.vue'
import UpdateCard from '@/components/UpdateCard.vue'
import { useI18n } from '@/utils/i18n'

/**
 * 关于页内容状态常量。
 * 统一状态字面量，避免页面与子组件之间出现拼写不一致。
 */
const CONTENT_STATUS = {
  LOADING: 'loading',
  SUCCESS: 'success',
  ERROR: 'error',
  EMPTY: 'empty'
}

const { locale, t } = useI18n()

const introductionContent = ref('')
const latestUpdate = ref('')
const historyUpdates = ref([])
const loadError = ref('')
const introductionStatus = ref(CONTENT_STATUS.LOADING)
const latestUpdateStatus = ref(CONTENT_STATUS.LOADING)
const historyUpdatesStatus = ref(CONTENT_STATUS.LOADING)
const latestLoadRequestId = ref(0)

/**
 * 仅在开发环境输出调试日志。
 * 这样保留排障能力，同时避免生产环境噪声和内部信息暴露过多。
 *
 * @param {...unknown} args - 需要输出的调试内容
 * @returns {void}
 */
const debugLog = (...args) => {
  if (import.meta.env.DEV) {
    console.log(...args)
  }
}

/**
 * 仅在开发环境输出错误细节。
 * 用户可见区域只展示简洁提示，详细错误保留在开发态控制台。
 *
 * @param {...unknown} args - 需要输出的错误详情
 * @returns {void}
 */
const debugError = (...args) => {
  if (import.meta.env.DEV) {
    console.error(...args)
  }
}

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

/**
 * 读取当前文档语言参数。
 * 关于页文档不再把语言写死在文件名里，而是通过接口参数交给后端做语言版本选择，
 * 这样前端始终只认“基础文件名”，切换语言时也更容易统一回退策略。
 *
 * @returns {string} 当前接口请求使用的语言代码
 */
const getDocsLanguage = () => {
  return locale.value || 'zh-CN'
}

// 加载 markdown 文件
const loadMarkdownFile = async (filename) => {
  try {
    const apiUrl = getApiUrl()
    const params = new URLSearchParams({
      lang: getDocsLanguage()
    })
    const url = `${apiUrl}/api/changelog/file/${filename}?${params.toString()}`
    debugLog('加载文件使用的 API:', url)
    debugLog('环境变量 VITE_API_URL:', import.meta.env.VITE_API_URL)

    const response = await fetch(url, {
      credentials: 'include'
    })
    debugLog('响应状态:', response.status, response.statusText)

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`)
    }

    const data = await response.json()
    if (!data?.data || typeof data.data.content !== 'string') {
      throw new Error(t('about.invalidData'))
    }

    debugLog('响应数据:', data)
    return data.data.content
  } catch (error) {
    debugError(`Error loading ${filename}:`, error)
    debugError('错误堆栈:', error?.stack)
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
  const requestId = Date.now()
  latestLoadRequestId.value = requestId

  try {
    debugLog('开始加载更新日志...')
    debugLog('当前文档语言:', getDocsLanguage())
    debugLog('当前页面地址:', window.location.href)
    loadError.value = ''
    introductionStatus.value = CONTENT_STATUS.LOADING
    latestUpdateStatus.value = CONTENT_STATUS.LOADING
    historyUpdatesStatus.value = CONTENT_STATUS.LOADING
    introductionContent.value = ''
    latestUpdate.value = ''
    historyUpdates.value = []

    // 加载简介
    const intro = await loadMarkdownFile('简介.md')
    debugLog('简介内容:', intro ? '加载成功' : '加载失败')
    if (typeof intro === 'string' && intro.trim()) {
      introductionContent.value = intro
      introductionStatus.value = CONTENT_STATUS.SUCCESS
    } else if (intro === '') {
      introductionStatus.value = CONTENT_STATUS.EMPTY
    } else {
      introductionStatus.value = CONTENT_STATUS.ERROR
    }

    // 从后端 API 获取版本文件列表
    const apiUrl = getApiUrl()
    const params = new URLSearchParams({
      lang: getDocsLanguage()
    })
    debugLog('使用 API 地址:', apiUrl)
    debugLog('完整请求 URL:', `${apiUrl}/api/changelog/files?${params.toString()}`)

    const versionsResponse = await fetch(`${apiUrl}/api/changelog/files?${params.toString()}`, {
      credentials: 'include'
    })
    debugLog('响应状态:', versionsResponse.status, versionsResponse.statusText)

    if (!versionsResponse.ok) {
      throw new Error(`更新日志列表加载失败（状态码: ${versionsResponse.status}）`)
    }

    const data = await versionsResponse.json()
    const updateFiles = data.data.files || []
    debugLog('从API获取的版本文件列表:', updateFiles)

    // 过滤出有效的更新文件(排除简介等特殊文件)
    const updates = updateFiles
      .map(filename => {
        const version = extractVersionFromFilename(filename)
        return version ? { version, filename } : null
      })
      .filter(item => item !== null)

    debugLog('找到的更新文件:', updates)

    // 按版本号排序
    updates.sort((a, b) => compareVersions(a.version, b.version))
    debugLog('排序后的更新列表:', updates)

    // 加载所有更新的内容
    const updateContents = await Promise.all(
      updates.map(async (update) => {
        const content = await loadMarkdownFile(update.filename)
        debugLog(`加载 ${update.filename}:`, content ? '成功' : '失败')
        return { content, filename: update.filename }
      })
    )

    // 过滤掉加载失败的内容
    const validUpdates = updateContents.filter(item => typeof item.content === 'string' && item.content.trim())
    const failedUpdates = updateContents.filter(item => item.content === null)
    debugLog('有效的更新数量:', validUpdates.length)

    // 如果用户在加载过程中切换了语言，则丢弃旧请求结果，
    // 避免英文与中文内容交叉覆盖。
    if (latestLoadRequestId.value !== requestId) {
      return
    }

    // 第一个作为最新更新
    if (validUpdates.length > 0) {
      latestUpdate.value = validUpdates[0].content
      historyUpdates.value = validUpdates.slice(1).map(item => item.content)
      latestUpdateStatus.value = CONTENT_STATUS.SUCCESS
      historyUpdatesStatus.value = historyUpdates.value.length > 0 ? CONTENT_STATUS.SUCCESS : CONTENT_STATUS.EMPTY
      debugLog('最新更新已设置,历史更新数量:', historyUpdates.value.length)
    } else {
      latestUpdateStatus.value = updates.length > 0 ? CONTENT_STATUS.ERROR : CONTENT_STATUS.EMPTY
      historyUpdatesStatus.value = updates.length > 1 ? CONTENT_STATUS.ERROR : CONTENT_STATUS.EMPTY
    }

    if (failedUpdates.length > 0 && !loadError.value) {
      loadError.value = t('about.partialLoadFailed')
    }

    if (updates.length === 0) {
      latestUpdateStatus.value = CONTENT_STATUS.EMPTY
      historyUpdatesStatus.value = CONTENT_STATUS.EMPTY
    }
  } catch (error) {
    debugError('Error loading changelog:', error)
    loadError.value = t('about.historyLoadFailedWithMessage', {
      message: error.message || t('about.tryLater')
    })
    if (latestUpdateStatus.value === CONTENT_STATUS.LOADING) {
      latestUpdateStatus.value = CONTENT_STATUS.ERROR
    }
    if (historyUpdatesStatus.value === CONTENT_STATUS.LOADING) {
      historyUpdatesStatus.value = CONTENT_STATUS.ERROR
    }
  }
}

onMounted(() => {
  debugLog('===== AboutView 组件已挂载 =====')
  debugLog('环境变量检查:')
  debugLog('  import.meta.env.VITE_API_URL:', import.meta.env.VITE_API_URL)
  debugLog('  window.location.hostname:', window.location.hostname)
  debugLog('  window.location.protocol:', window.location.protocol)
  debugLog('  getApiUrl():', getApiUrl())
  debugLog('==============================')
  loadChangelog()
})

watch(
  locale,
  () => {
    currentPage.value = 1
    loadChangelog()
  }
)
</script>

<style scoped>
.about {
  /* 关于页改为引用全局页面宽度变量，便于后续统一维护。 */
  max-width: var(--page-content-width);
  margin: 0 auto;
  padding: 2rem;
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
  background: color-mix(in srgb, var(--danger-color) 10%, transparent);
  border: 1px solid color-mix(in srgb, var(--danger-color) 30%, transparent);
  border-radius: var(--radius-md);
  padding: 1.5rem;
  margin-bottom: 2rem;
}

.error-message h3 {
  color: var(--danger-color);
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
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  line-height: 1;
  box-shadow: var(--shadow-md);
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
  box-shadow: var(--shadow-lg);
}

.pagination-btn:hover:not(:disabled)::before {
  left: 100%;
}

.pagination-btn:active:not(:disabled) {
  transform: translateY(0);
  box-shadow: var(--shadow-sm);
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
  border-radius: var(--radius-sm);
}

.pagination-info::before {
  content: none;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .about {
    padding: 1rem;
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
