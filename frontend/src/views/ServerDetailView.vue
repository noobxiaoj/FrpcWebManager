<template>
  <div class="server-detail">
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>{{ t('common.loading') }}</p>
    </div>

    <div v-else-if="server" class="detail-container">
      <section class="card main-card">
        <div class="card-header">
          <div class="header-left">
            <AppButton class="btn-back" preserve-style @click="goBack" :title="t('serverDetail.backTitle')">
              <template #icon>
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <line x1="19" y1="12" x2="5" y2="12"></line>
                  <polyline points="12 19 5 12 12 5"></polyline>
                </svg>
              </template>
            </AppButton>
          </div>

          <div class="header-center">
            <h1 class="page-title">{{ server.name }}</h1>
            <div
              class="status-badge"
              :class="{
                'status-online': server.status === 'online',
                'status-offline': server.status === 'offline',
                'status-no-task': server.status === 'no_task',
                'status-fault': server.status === 'fault',
                'status-suspected-abnormal': server.status === 'suspected_abnormal'
              }"
            >
              <span class="status-indicator"></span>
              <span>{{ getStatusText(server.status) }}</span>
            </div>
          </div>

          <div class="header-actions">
            <AppButton
              class="btn-action refresh"
              preserve-style
              @click="refreshServerData"
              :disabled="actionLoading"
              :title="t('serverDetail.refreshTitle')"
            >
              <template #icon>
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="23 4 23 10 17 10"></polyline>
                  <polyline points="1 20 1 14 7 14"></polyline>
                  <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"></path>
                </svg>
              </template>
            </AppButton>

            <AppButton
              class="btn-action delete"
              preserve-style
              @click="confirmDeleteServer"
              :disabled="actionLoading"
              :title="t('serverDetail.deleteTitle')"
            >
              <template #icon>
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="3 6 5 6 21 6"></polyline>
                  <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                </svg>
              </template>
            </AppButton>
          </div>
        </div>

        <div class="card-content">
          <div class="info-grid">
            <div class="info-item">
              <div class="info-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="2" y="2" width="20" height="8" rx="2" ry="2"></rect>
                  <rect x="2" y="14" width="20" height="8" rx="2" ry="2"></rect>
                </svg>
              </div>
              <div class="info-content">
                <span class="info-label">{{ t('serverDetail.address') }}</span>
                <span class="info-value">{{ server.address }}</span>
              </div>
            </div>

            <div class="info-item">
              <div class="info-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="3"></circle>
                  <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
                </svg>
              </div>
              <div class="info-content">
                <span class="info-label">{{ t('serverDetail.processPort') }}</span>
                <span class="info-value">{{ processPortText }}</span>
              </div>
            </div>

            <div class="info-item">
              <div class="info-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
                  <line x1="16" y1="2" x2="16" y2="6"></line>
                  <line x1="8" y1="2" x2="8" y2="6"></line>
                  <line x1="3" y1="10" x2="21" y2="10"></line>
                </svg>
              </div>
              <div class="info-content">
                <span class="info-label">{{ t('serverDetail.createdAt') }}</span>
                <span class="info-value">{{ formatDate(server.createdAt) }}</span>
              </div>
            </div>

            <div class="info-item">
              <div class="info-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
                  <polyline points="22 4 12 14.01 9 11.01"></polyline>
                </svg>
              </div>
              <div class="info-content">
                <span class="info-label">{{ t('serverDetail.portCount') }}</span>
                <span class="info-value">{{ t('common.countUnit', { count: relatedPortCount }) }}</span>
              </div>
            </div>

            <div class="info-item">
              <div class="info-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10"></circle>
                  <polyline points="12 6 12 12 16 14"></polyline>
                </svg>
              </div>
              <div class="info-content">
                <span class="info-label">{{ t('serverDetail.uptime') }}</span>
                <span class="info-value">{{ server.uptime || t('common.unknown') }}</span>
              </div>
            </div>

            <div class="info-item">
              <div class="info-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M12 20h9"></path>
                  <path d="M16.5 3.5a2.121 2.121 0 1 1 3 3L7 19l-4 1 1-4Z"></path>
                </svg>
              </div>
              <div class="info-content">
                <span class="info-label">{{ t('serverDetail.lastRefresh') }}</span>
                <span class="info-value">{{ lastRefreshText }}</span>
              </div>
            </div>
          </div>

          <div class="section-block">
            <div class="section-heading">
              <h2 class="section-title">{{ t('serverDetail.relatedTasks', { count: relatedTasks.length }) }}</h2>
              <AppButton
                class="section-toggle"
                preserve-style
                type="button"
                :aria-expanded="String(taskSectionExpanded)"
                :title="taskSectionExpanded ? t('serverDetail.collapseTasks') : t('serverDetail.expandTasks')"
                @click="toggleTaskSection"
              >
                <span>{{ taskSectionExpanded ? t('common.collapse') : t('common.expand') }}</span>
                <svg class="section-toggle-icon" :class="{ 'is-expanded': taskSectionExpanded }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="6 9 12 15 18 9"></polyline>
                </svg>
              </AppButton>
            </div>

            <div v-if="taskSectionExpanded && relatedTasks.length === 0" class="empty-state compact-empty">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="12" y1="8" x2="12" y2="12"></line>
                <line x1="12" y1="16" x2="12.01" y2="16"></line>
              </svg>
              <p>{{ t('serverDetail.noRelatedTasks') }}</p>
            </div>

            <div
              v-else-if="taskSectionExpanded"
              ref="taskChipListRef"
              class="task-chip-list"
              :style="{ '--task-chip-columns': taskChipColumns }"
            >
              <button
                v-for="task in relatedTasks"
                :key="task.id"
                type="button"
                class="task-chip"
                @click="viewTask(task.id)"
              >
                <span class="task-chip-name">{{ task.name }}</span>
                <span
                  class="task-chip-status"
                  :class="`task-chip-status--${task.status}`"
                  :title="getTaskStatusText(task.status)"
                  :aria-label="t('serverDetail.taskStatusAria', { status: getTaskStatusText(task.status) })"
                ></span>
              </button>
            </div>
          </div>

          <div class="section-block log-section-card" :class="{ 'log-section-card-collapsed': !logSectionExpanded }">
            <div class="section-heading" :class="{ 'section-heading-collapsed': !logSectionExpanded }">
              <h2 class="section-title">{{ t('serverDetail.runtimeLogs') }}</h2>
              <AppButton
                class="section-toggle"
                preserve-style
                type="button"
                :aria-expanded="String(logSectionExpanded)"
                :title="logSectionExpanded ? t('serverDetail.collapseLogs') : t('serverDetail.expandLogs')"
                @click="toggleLogSection"
              >
                <span>{{ logSectionExpanded ? t('common.collapse') : t('common.expand') }}</span>
                <svg class="section-toggle-icon" :class="{ 'is-expanded': logSectionExpanded }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="6 9 12 15 18 9"></polyline>
                </svg>
              </AppButton>
            </div>

            <div v-if="logSectionExpanded" class="log-panel">
              <div class="section-subheading">
                <div class="section-meta">
                  <span class="section-meta-total">{{ t('serverDetail.totalLines', { filtered: filteredLogCount, total: logs.length }) }}</span>
                  <LogFilterButton
                    v-for="filter in logFilterOptions"
                    :key="filter.value"
                    :label="filter.label"
                    :count="getLogCount(filter.value)"
                    :tone="filter.value"
                    :active="isLogFilterSelected(filter.value)"
                    @toggle="toggleLogFilter(filter.value)"
                  />
                </div>

                <div class="section-actions">
                  <AppButton
                    class="btn-action delete log-clear-action"
                    preserve-style
                    type="button"
                    :disabled="actionLoading"
                    :title="t('serverDetail.clearLogs')"
                    :aria-label="t('serverDetail.clearLogs')"
                    @click="clearLogs"
                  >
                    <template #icon>
                      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 48 48" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="3" aria-hidden="true">
                        <path d="M11.732 38.242c4.055 3.835 7.865 4.258 7.865 4.258c3.178-2.839 4.91-5.375 6.842-10.807c0 0-7.548-7.321-9.091-9.625c0 0-3.509 1.8-6.223 2.048c-2.045.187-5.536-.86-5.536-.86c.028 2.918 2.023 11.088 6.143 14.986"></path>
                        <path d="M17.348 22.068c1.903-1.076 5.383-3.554 6.994-2.15a26 26 0 0 1 3.477 3.421c1.099 1.434-.418 5.077-1.38 8.354"></path>
                        <path d="M26.055 21.448C29.422 18.428 42.411 5.5 42.411 5.5M7.03 30.106c4.237.317 7.03-.612 7.03-.612m-3.142 7.893c6.19-.867 8.333-3.97 8.333-3.97"></path>
                      </svg>
                    </template>
                  </AppButton>
                </div>
              </div>

              <div v-if="logsLoading" class="log-loading">
                <div class="spinner spinner-small"></div>
                <p>{{ t('common.loadingLogs') }}</p>
              </div>

              <div v-else-if="logs.length > 0" class="log-scroll-shell">
                <div
                  ref="logContainerRef"
                  class="log-content"
                >
                  <div
                    v-for="(log, index) in displayedLogs"
                    :key="index"
                    class="log-line"
                    :class="getLogLineClass(log)"
                  >
                    <span class="log-timestamp">{{ log.timestamp }}</span>
                    <span class="log-level" :class="`level-${log.level}`">{{ log.level }}</span>
                    <span class="log-message">{{ stripTimestampFromMessage(log.message) }}</span>
                  </div>
                </div>
              </div>

              <div v-else class="empty-state compact-empty">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M3 3v18h18"></path>
                  <path d="M19 9l-5 5-4-4-3 3"></path>
                </svg>
                <p>{{ logEmptyText }}</p>
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>

    <div v-if="showDeleteConfirmModal" class="modal-overlay" @click="closeDeleteConfirmModal">
      <div class="modal-card" @click.stop>
        <div class="modal-header modal-header--dialog">
          <h3 class="modal-title">{{ t('serverDetail.deleteConfirmTitle') }}</h3>
          <AppButton class="modal-close" preserve-style @click="closeDeleteConfirmModal" :title="t('common.close')">
            <template #icon>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
              </svg>
            </template>
          </AppButton>
        </div>

        <div class="modal-body">
          <div class="delete-confirm-content">
            <div class="delete-warning-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path>
                <line x1="12" y1="9" x2="12" y2="13"></line>
                <line x1="12" y1="17" x2="12.01" y2="17"></line>
              </svg>
            </div>
            <h4 class="delete-confirm-title">{{ t('serverDetail.deleteServerTitle') }}</h4>

            <div v-if="serverTasks.length > 0">
              <p class="delete-confirm-message">
                {{ t('serverDetail.deleteWithTasks', { name: server?.name, count: serverTasks.length }) }}
              </p>
              <ul class="task-list">
                <li v-for="(task, index) in serverTasks" :key="index" class="task-item">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="9 11 12 14 22 4"></polyline>
                    <path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"></path>
                  </svg>
                  {{ task }}
                </li>
              </ul>
              <p class="delete-confirm-warning">{{ t('serverDetail.deleteWithTasksWarning') }}</p>
            </div>

            <div v-else>
              <p class="delete-confirm-message">
                {{ t('serverDetail.deleteWithoutTasks', { name: server?.name }) }}
              </p>
              <p class="delete-confirm-warning">{{ t('serverDetail.deleteWithoutTasksWarning') }}</p>
            </div>
          </div>

          <div class="dialog-actions">
            <AppButton variant="secondary" type="button" @click="closeDeleteConfirmModal">{{ t('common.cancel') }}</AppButton>
            <AppButton variant="danger" type="button" @click="deleteServer">{{ t('common.confirmDelete') }}</AppButton>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { serverAPI } from '@/api/server'
import { useTaskStore } from '@/stores/task'
import { getApiBaseUrl } from '@/utils/api'
import AppButton from '@/components/AppButton.vue'
import LogFilterButton from '@/components/LogFilterButton.vue'
import { useI18n } from '@/utils/i18n'

const router = useRouter()
const route = useRoute()
const taskStore = useTaskStore()
const { locale, t } = useI18n()
const AVAILABLE_LOG_FILTERS = ['info', 'error', 'warn']

const loading = ref(true)
const logsLoading = ref(false)
const actionLoading = ref(false)
const server = ref(null)
const logs = ref([])
const lastRefreshTime = ref('')
const logContainerRef = ref(null)
const taskChipListRef = ref(null)
const refreshTimer = ref(null)
const settingsEventSource = ref(null)
const currentTime = ref(Date.now())
const currentTimeTimer = ref(null)
const taskChipColumns = ref(1)
const taskChipResizeObserver = ref(null)
const taskSectionExpanded = ref(true)
const logSectionExpanded = ref(true)
const selectedLogFilters = ref([...AVAILABLE_LOG_FILTERS])
const showDeleteConfirmModal = ref(false)
const serverTasks = ref([])
const settings = ref({
  refreshInterval: 10
})

const serverId = computed(() => route.params.id)
const logFilterOptions = computed(() => AVAILABLE_LOG_FILTERS.map((value) => ({
  value,
  label: t(`serverDetail.logFilters.${value}`)
})))

/**
 * 将形如 `host:port` 的地址拆成主机和端口。
 * 这里使用最后一个冒号分割，兼容当前项目里服务端返回的存储格式。
 * @param {string} address - 后端返回的服务器地址
 * @returns {{ host: string, port: number|null }} 主机和端口
 */
const parseServerAddress = (address = '') => {
  const separatorIndex = address.lastIndexOf(':')
  if (separatorIndex === -1) {
    return {
      host: address,
      port: null
    }
  }

  return {
    host: address.slice(0, separatorIndex),
    port: Number(address.slice(separatorIndex + 1))
  }
}

/**
 * 计算字符串的“视觉长度”。
 * 中文、全角字符在卡片中通常占用更宽，因此按 2 个单位估算；
 * 英文、数字、半角符号按 1 个单位估算。
 * 这里不追求像素级精确，而是为 3/2/1 列布局提供稳定且可预期的判断依据。
 * @param {string} text - 需要估算宽度的文本。
 * @returns {number} 返回估算后的视觉长度单位。
 */
const getVisualTextLength = (text = '') => {
  return Array.from(text).reduce((total, char) => {
    return total + (/[\u0000-\u00ff]/.test(char) ? 1 : 2)
  }, 0)
}

/**
 * 获取任务用于排序的“数量”指标。
 * 当前详情页里最稳定且符合业务含义的数量字段是 proxies 数量，
 * 代理配置越多，任务复杂度通常越高，因此按它进行“数量多 > 数量少”的排序。
 * @param {object} task - 任务对象。
 * @returns {number} 返回代理配置数量。
 */
const getTaskProxyCount = (task) => {
  return Array.isArray(task?.proxies) ? task.proxies.length : 0
}

/**
 * 获取任务创建时间的时间戳。
 * 创建时间越早，排序优先级越高；如果时间缺失或非法，则视为最新，避免意外排到前面。
 * @param {object} task - 任务对象。
 * @returns {number} 返回创建时间戳，非法时返回 Infinity。
 */
const getTaskCreatedTimestamp = (task) => {
  const timestamp = new Date(task?.createdAt || '').getTime()
  return Number.isNaN(timestamp) ? Number.POSITIVE_INFINITY : timestamp
}

/**
 * 根据当前关联任务估算单个任务卡片的最小宽度。
 * 估算值会同时考虑：
 * 1. 任务名称长度，避免名称在多列时明显挤压；
 * 2. 卡片内边距和状态点所占空间；
 * 3. 一个合理的上下限，避免极端名称导致布局失控。
 * @param {Array<object>} tasks - 当前服务器的关联任务列表。
 * @returns {number} 返回单个卡片建议的最小宽度（像素）。
 */
const getTaskChipMinWidth = (tasks) => {
  const longestNameLength = tasks.reduce((maxLength, task) => {
    return Math.max(maxLength, getVisualTextLength(task?.name || t('common.unnamed')))
  }, 0)

  const estimatedWidth = 120 + longestNameLength * 10
  return Math.max(220, Math.min(estimatedWidth, 420))
}

/**
 * 根据容器宽度和任务名称长度，动态决定任务卡片显示为 3 / 2 / 1 列。
 * 规则按需求优先尝试 3 列，再回退到 2 列，最后保证至少 1 列。
 * 成功选中列数后，CSS Grid 会让每列自动均分剩余空间，实现“自动填充长度”。
 */
const updateTaskChipColumns = () => {
  const container = taskChipListRef.value
  if (!container) {
    taskChipColumns.value = 1
    return
  }

  const tasks = relatedTasks.value
  if (tasks.length === 0) {
    taskChipColumns.value = 1
    return
  }

  const containerWidth = container.clientWidth
  const columnGap = 12
  const taskChipMinWidth = getTaskChipMinWidth(tasks)

  for (let columns = Math.min(3, tasks.length); columns >= 1; columns -= 1) {
    const requiredWidth = columns * taskChipMinWidth + (columns - 1) * columnGap
    if (requiredWidth <= containerWidth) {
      taskChipColumns.value = columns
      return
    }
  }

  taskChipColumns.value = 1
}

const relatedTasks = computed(() => {
  if (!server.value) return []

  const { host, port } = parseServerAddress(server.value.address)
  return taskStore.tasks
    .filter((task) => task.serverAddr === host && task.serverPort === port)
    .slice()
    .sort((taskA, taskB) => {
      const runningWeightDiff = Number(taskB.status === 'running') - Number(taskA.status === 'running')
      if (runningWeightDiff !== 0) return runningWeightDiff

      const proxyCountDiff = getTaskProxyCount(taskB) - getTaskProxyCount(taskA)
      if (proxyCountDiff !== 0) return proxyCountDiff

      return getTaskCreatedTimestamp(taskA) - getTaskCreatedTimestamp(taskB)
    })
})

/**
 * 汇总当前服务器下所有关联任务的端口配置数量。
 * 这里将每个任务的 proxies 数量累加，得到更符合详情页语义的“端口数”。
 * @returns {number} 返回当前服务器关联的总端口数。
 */
const relatedPortCount = computed(() => {
    return relatedTasks.value.reduce((total, task) => total + getTaskProxyCount(task), 0)
})

/**
 * 提供“进程端口”的展示文案。
 * 当前详情页里的进程端口来自服务器地址中的监听端口，
 * 如果地址不完整或端口解析失败，则统一显示“未分配”。
 * @returns {string|number} 返回端口号或占位文案。
 */
const processPortText = computed(() => {
  if (!server.value) {
    return t('common.notAssigned')
  }

  const { port } = parseServerAddress(server.value.address)
  return Number.isInteger(port) && port > 0 ? port : t('common.notAssigned')
})

/**
 * 供界面渲染使用的日志列表。
 * 这里会先按当前筛选条件决定显示全部 / 错误 / 警告，
 * 再在展示层反转顺序，实现“最新在上、最早在下”。
 * 原始日志数组保持不变，避免影响统计、刷新与后续数据处理逻辑。
 * @returns {Array<object>} 返回筛选后且按最新优先排序的日志数组。
 */
const displayedLogs = computed(() => {
  const filteredLogs = logs.value.filter(log => selectedLogFilters.value.includes(log.level))

  return filteredLogs.slice().reverse()
})

/**
 * 生成日志区空状态文案。
 * 当用户启用了错误或警告筛选，但当前结果为空时，
 * 需要给出更明确的提示，避免和“完全没有日志”混淆。
 * @returns {string} 返回当前日志列表对应的空状态文案。
 */
const logEmptyText = computed(() => {
  if (selectedLogFilters.value.length === 0) {
    return t('serverDetail.noLogFilters')
  }

  if (selectedLogFilters.value.length === AVAILABLE_LOG_FILTERS.length) {
    return t('serverDetail.noLogs')
  }

  const selectedFilterLabels = AVAILABLE_LOG_FILTERS
    .filter(filter => selectedLogFilters.value.includes(filter))
    .map(filter => t(`serverDetail.logFilters.${filter}`))

  return t('serverDetail.noFilteredLogs', {
    filters: selectedFilterLabels.join(locale.value === 'en-US' ? ' or ' : '或')
  })
})

/**
 * 当前筛选结果数量。
 * 这里直接复用展示层日志数组长度，保证“总行数 x/y”中的 x
 * 始终与界面当前实际渲染的数据条数保持一致。
 * @returns {number} 返回当前筛选后的日志数量。
 */
const filteredLogCount = computed(() => displayedLogs.value.length)

const lastRefreshText = computed(() => {
  if (!lastRefreshTime.value) return t('serverDetail.notRefreshed')

  const diffMs = currentTime.value - new Date(lastRefreshTime.value).getTime()
  const diffSeconds = Math.max(0, Math.floor(diffMs / 1000))

  if (diffSeconds < 5) return t('serverDetail.justNow')
  if (diffSeconds < 60) return t('serverDetail.secondsAgo', { count: diffSeconds })

  const diffMinutes = Math.floor(diffSeconds / 60)
  if (diffMinutes < 60) return t('serverDetail.minutesAgo', { count: diffMinutes })

  const diffHours = Math.floor(diffMinutes / 60)
  if (diffHours < 24) return t('serverDetail.hoursAgo', { count: diffHours })

  const diffDays = Math.floor(diffHours / 24)
  return t('serverDetail.daysAgo', { count: diffDays })
})

const getStatusText = (status) => {
  switch (status) {
    case 'online':
      return t('status.server.online')
    case 'offline':
      return t('status.server.offline')
    case 'no_task':
      return t('status.server.noTask')
    case 'fault':
      return t('status.server.fault')
    case 'suspected_abnormal':
      return t('status.server.suspectedAbnormal')
    default:
      return t('status.server.unknown')
  }
}

/**
 * 获取任务状态对应的中文文案。
 * @param {string} status - 任务状态值，例如 running / stopped / error。
 * @returns {string} 返回用于提示信息和无障碍描述的中文状态文本。
 */
const getTaskStatusText = (status) => {
  switch (status) {
    case 'running':
      return t('status.task.running')
    case 'stopped':
      return t('status.task.stopped')
    case 'error':
      return t('status.task.error')
    default:
      return status || t('status.task.unknown')
  }
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'

  const date = new Date(dateStr)
  if (Number.isNaN(date.getTime())) return '-'

  return date.toLocaleString(locale.value, {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getLogLineClass = (log) => `log-line-${log.level}`

const getLogCount = (level) => logs.value.filter(log => log.level === level).length

/**
 * 判断某个日志级别筛选项是否处于选中状态。
 * @param {'info'|'error'|'warn'} filter - 需要判断的筛选项。
 * @returns {boolean} 返回当前筛选项是否已选中。
 */
const isLogFilterSelected = (filter) => selectedLogFilters.value.includes(filter)

/**
 * 切换日志筛选条件。
 * 信息、错误和警告都支持独立开关，因此可以同时选中；
 * 默认状态为全部选中，但用户也可以手动全部取消，此时界面显示空结果。
 * @param {'info'|'error'|'warn'} filter - 需要切换的筛选项。
 * @returns {void} 无返回值。
 */
const toggleLogFilter = (filter) => {
  if (!AVAILABLE_LOG_FILTERS.includes(filter)) {
    return
  }

  if (isLogFilterSelected(filter)) {
    selectedLogFilters.value = selectedLogFilters.value.filter(item => item !== filter)
  } else {
    selectedLogFilters.value = [...selectedLogFilters.value, filter]
  }

  selectedLogFilters.value = AVAILABLE_LOG_FILTERS.filter(item => selectedLogFilters.value.includes(item))

  scrollLogToTop()
}

const stripTimestampFromMessage = (message) => {
  if (!message) return ''

  let result = message
  result = result.replace(/\x1b\[[0-9;]*m/g, '')
  result = result.replace(/\033\[[0-9;]*m/g, '')
  result = result.replace(/\[[0-9;]*m/g, '')
  result = result.replace(/^\[?\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2}(\.\d+)?\]?\s*/, '')
  return result
}

/**
 * 判断用户是否仍停留在日志顶部附近。
 * 当界面改为“最新在上”后，顶部就是最新日志区域；
 * 因此只有用户仍在顶部附近时，自动刷新才继续贴顶，避免打断查看旧日志。
 * @returns {boolean} 是否接近日志顶部
 */
const isLogNearTop = () => {
  const container = logContainerRef.value
  if (!container) return true

  const threshold = 50
  return container.scrollTop < threshold
}

/**
 * 将日志滚动区域定位到顶部。
 * 在“最新在上”的布局下，顶部就是最新日志所在位置，
 * 因此刷新完成后需要贴顶，而不是再滚到底部。
 * @returns {void} 无返回值。
 */
const scrollLogToTop = () => {
  nextTick(() => {
    const container = logContainerRef.value
    if (container) {
      container.scrollTop = 0
    }
  })
}

/**
 * 切换“关联任务”区域的折叠状态。
 * 展开后需要等待 DOM 渲染完成，再重新计算任务卡片列数，
 * 这样可以避免容器从隐藏状态恢复时宽度计算不准确。
 * @returns {Promise<void>} 无返回值。
 */
const toggleTaskSection = async () => {
  taskSectionExpanded.value = !taskSectionExpanded.value

  if (taskSectionExpanded.value) {
    await nextTick()
    updateTaskChipColumns()
  }
}

/**
 * 切换“运行日志”区域的折叠状态。
 * 当日志重新展开时，如果用户原本停留在底部，则继续滚动到底部，
 * 保证实时日志场景下的阅读体验更自然。
 * @returns {Promise<void>} 无返回值。
 */
const toggleLogSection = async () => {
  logSectionExpanded.value = !logSectionExpanded.value

  if (logSectionExpanded.value) {
    await nextTick()
    scrollLogToTop()
  }
}

const loadSettings = async () => {
  try {
    const response = await fetch(`${getApiBaseUrl()}/api/settings`, {
      credentials: 'include'
    })

    if (!response.ok) {
      throw new Error(t('serverDetail.messages.loadSettingsFailed'))
    }

    const data = await response.json()
    if (data.data) {
      settings.value = {
        refreshInterval: data.data.refreshInterval ?? 10
      }
    }
  } catch (error) {
    console.error('加载设置失败:', error)
    settings.value = {
      refreshInterval: 10
    }
  }
}

const loadServer = async () => {
  server.value = await serverAPI.getServer(serverId.value)
}

const loadLogs = async () => {
  const shouldStickTop = isLogNearTop()

  try {
    logsLoading.value = true
    logs.value = await serverAPI.getLogs(serverId.value, 100)
    lastRefreshTime.value = new Date().toISOString()

    if (shouldStickTop) {
      scrollLogToTop()
    }
  } catch (error) {
    console.error('加载日志失败:', error)
    alert(`${t('serverDetail.messages.loadLogsFailed')}: ${error.message}`)
  } finally {
    logsLoading.value = false
  }
}

/**
 * 同步刷新详情页的核心数据：
 * 1. 服务器详情
 * 2. 任务列表（用于关联任务展示）
 * 3. 服务器日志
 */
const refreshServerData = async () => {
  try {
    actionLoading.value = true
    await Promise.all([
      loadServer(),
      taskStore.fetchTasks(),
      loadLogs()
    ])
  } catch (error) {
    console.error('刷新服务器详情失败:', error)
  } finally {
    actionLoading.value = false
  }
}

const resetRefreshTimer = () => {
  if (refreshTimer.value) {
    clearInterval(refreshTimer.value)
    refreshTimer.value = null
  }

  if (settings.value.refreshInterval > 0) {
    refreshTimer.value = setInterval(() => {
      refreshServerData()
    }, settings.value.refreshInterval * 1000)
  }
}

const goBack = () => {
  router.push('/')
}

const viewTask = (id) => {
  router.push(`/tasks/${id}`)
}

const clearLogs = async () => {
  if (!server.value) return

  try {
    actionLoading.value = true
    await serverAPI.clearLogs(server.value.id)
    logs.value = []
    lastRefreshTime.value = new Date().toISOString()
  } catch (error) {
    console.error('清空日志失败:', error)
    alert(`${t('serverDetail.messages.clearLogsFailed')}: ${error.message}`)
  } finally {
    actionLoading.value = false
  }
}

const confirmDeleteServer = async () => {
  if (!server.value) return

  try {
    actionLoading.value = true
    const result = await serverAPI.deleteServer(server.value.id, false)

    if (result.hasTasks) {
      serverTasks.value = result.tasks
      showDeleteConfirmModal.value = true
      return
    }

    router.push('/')
  } catch (error) {
    console.error('检查服务器任务失败:', error)
    alert(`${t('serverDetail.messages.checkTasksFailed')}: ${error.message}`)
  } finally {
    actionLoading.value = false
  }
}

const closeDeleteConfirmModal = () => {
  showDeleteConfirmModal.value = false
  serverTasks.value = []
}

const deleteServer = async () => {
  if (!server.value) return

  try {
    actionLoading.value = true
    await serverAPI.deleteServer(server.value.id, true)
    closeDeleteConfirmModal()
    router.push('/')
  } catch (error) {
    console.error('删除服务器失败:', error)
    alert(`${t('serverDetail.messages.deleteFailed')}: ${error.message}`)
  } finally {
    actionLoading.value = false
  }
}

const initPage = async () => {
  try {
    loading.value = true
    await loadSettings()
    await Promise.all([
      taskStore.fetchTasks(),
      loadServer(),
      loadLogs()
    ])
    resetRefreshTimer()
  } catch (error) {
    console.error('初始化服务器详情失败:', error)
    alert(`${t('serverDetail.messages.initFailed')}: ${error.message}`)
    router.push('/')
  } finally {
    loading.value = false
  }
}

watch(serverId, async (newId, oldId) => {
  if (newId && newId !== oldId) {
    await initPage()
  }
})

watch(() => settings.value.refreshInterval, () => {
  resetRefreshTimer()
})

watch(logSectionExpanded, async (expanded) => {
  if (!expanded) return

  await nextTick()
})

watch(relatedTasks, async () => {
  if (!taskSectionExpanded.value) return

  await nextTick()
  taskChipResizeObserver.value?.disconnect()
  if (taskChipResizeObserver.value && taskChipListRef.value) {
    taskChipResizeObserver.value.observe(taskChipListRef.value)
  }
  updateTaskChipColumns()
}, { deep: true })

onMounted(async () => {
  await initPage()
  await nextTick()
  updateTaskChipColumns()

  if (typeof ResizeObserver !== 'undefined') {
    taskChipResizeObserver.value = new ResizeObserver(() => {
      updateTaskChipColumns()
    })

    if (taskChipListRef.value) {
      taskChipResizeObserver.value.observe(taskChipListRef.value)
    }
  }

  currentTimeTimer.value = setInterval(() => {
    currentTime.value = Date.now()
  }, 3000)

  settingsEventSource.value = new EventSource(`${getApiBaseUrl()}/api/settings/events`, {
    withCredentials: true
  })
  settingsEventSource.value.addEventListener('settings-updated', async () => {
    await loadSettings()
  })
  settingsEventSource.value.onerror = () => {
    settingsEventSource.value?.close()
  }
})

onUnmounted(() => {
  if (refreshTimer.value) {
    clearInterval(refreshTimer.value)
  }

  if (currentTimeTimer.value) {
    clearInterval(currentTimeTimer.value)
  }

  taskChipResizeObserver.value?.disconnect()
  settingsEventSource.value?.close()
})
</script>

<style scoped>
.server-detail {
  min-height: 100vh;
  background: var(--bg-primary);
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 50vh;
  gap: 1rem;
  color: var(--text-secondary);
}

.spinner {
  width: 3rem;
  height: 3rem;
  border: 3px solid var(--border-color);
  border-top-color: var(--accent-color);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

.spinner-small {
  width: 1.5rem;
  height: 1.5rem;
  border-width: 2px;
}

.detail-container {
  padding: 1.5rem 2rem;
  max-width: 1400px;
  margin: 0 auto;
}

.card {
  background: var(--card-bg);
  border-radius: var(--radius-lg);
  border: 1px solid var(--border-color);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
}

.main-card {
  box-shadow: var(--shadow-lg);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.25rem 1.5rem;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
  position: sticky;
  top: 0;
  z-index: 10;
}

.header-left {
  display: flex;
  align-items: center;
}

.btn-back {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2.5rem;
  height: 2.5rem;
  border-radius: var(--radius-md);
  border: 1px solid var(--border-color);
  background: var(--bg-primary);
  color: var(--text-primary);
}

.btn-back:hover {
  background: var(--accent-color-bg);
  border-color: var(--accent-color);
  color: var(--accent-color);
}

.btn-back svg {
  width: 1rem;
  height: 1rem;
}

.header-center {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 1rem;
  min-width: 0;
}

.page-title {
  margin: 0;
  font-size: 1.4rem;
  font-weight: 700;
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  padding: 0.35rem 0.75rem;
  border-radius: var(--radius-pill);
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
}

.status-indicator {
  width: 8px;
  height: 8px;
  border-radius: 999px;
}

.status-online {
  background: var(--success-color-bg);
  color: var(--success-color);
}

.status-online .status-indicator {
  background: var(--success-color);
}

.status-offline {
  background: var(--danger-color-bg);
  color: var(--danger-color);
}

.status-offline .status-indicator {
  background: var(--danger-color);
}

.status-no-task {
  background: var(--bg-primary);
  color: var(--text-secondary);
}

.status-no-task .status-indicator {
  background: var(--text-secondary);
}

.status-fault {
  background: rgba(245, 158, 11, 0.12);
  color: #d97706;
}

.status-fault .status-indicator {
  background: #d97706;
}

.status-suspected-abnormal {
  background: rgba(99, 102, 241, 0.12);
  color: #4f46e5;
}

.status-suspected-abnormal .status-indicator {
  background: #4f46e5;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.btn-action {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2.25rem;
  height: 2.25rem;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  background: var(--bg-primary);
  color: var(--text-primary);
}

.btn-action svg {
  width: 1rem;
  height: 1rem;
}

.btn-action:hover {
  transform: translateY(-1px);
}

.btn-action.refresh {
  border-color: var(--accent-color);
  background: var(--accent-color-bg);
  color: var(--accent-color);
}

.btn-action.clear,
.btn-action.delete {
  border-color: var(--danger-color);
  background: var(--danger-color-bg);
  color: var(--danger-color);
}

.card-content {
  padding: 1.5rem;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 0.75rem;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.95rem 1rem;
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
}

.info-icon {
  width: 2rem;
  height: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-sm);
  background: var(--accent-color-bg);
  flex-shrink: 0;
}

.info-icon svg {
  width: 1rem;
  height: 1rem;
  stroke: var(--accent-color);
}

.info-content {
  display: flex;
  flex-direction: column;
  gap: 0.3rem;
  min-width: 0;
}

.info-label {
  font-size: 0.7rem;
  font-weight: 600;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.info-value {
  font-size: 0.92rem;
  font-weight: 500;
  color: var(--text-primary);
  word-break: break-all;
}

.section-block {
  margin-top: 1.5rem;
}

.log-section-card {
  padding: 1rem;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  background:
    linear-gradient(180deg, rgba(15, 23, 42, 0.02) 0%, rgba(15, 23, 42, 0.01) 100%),
    var(--card-bg);
  box-shadow: var(--shadow-sm);
}

.log-section-card-collapsed {
  padding-bottom: 0.75rem;
}

.section-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 1rem;
  margin-bottom: 1rem;
}

.section-heading-collapsed {
  margin-bottom: 0;
}

.section-title {
  margin: 0;
  font-size: 1.08rem;
  color: var(--text-primary);
}

.section-meta {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 1rem;
  font-size: 0.88rem;
  color: var(--text-secondary);
}

.section-meta-total {
  color: var(--text-primary);
  font-weight: 700;
  line-height: 1;
}

.section-subheading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.75rem;
  margin-bottom: 1rem;
}

/* 将筛选条与日志列表收拢到同一块面板内，形成更清晰的整体边界。 */
.log-panel {
  padding: 1rem;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  background: color-mix(in srgb, var(--bg-secondary) 72%, var(--card-bg) 28%);
}

.section-actions {
  display: inline-flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.section-toggle {
  display: inline-flex;
  align-items: center;
  gap: 0.45rem;
  padding: 0.45rem 0.75rem;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-pill);
  background: var(--bg-secondary);
  color: var(--text-secondary);
  font-size: 0.85rem;
  font-weight: 600;
  line-height: 1;
  white-space: nowrap;
}

.section-toggle:hover {
  border-color: var(--accent-color);
  color: var(--accent-color);
  background: var(--accent-color-bg);
}

.log-clear-action svg {
  width: 1.575rem;
  height: 1.575rem;
}

:global([data-theme="light"] .log-clear-action) {
  border-color: #9c27b0 !important;
  background: rgba(156, 39, 176, 0.1) !important;
  color: #9c27b0 !important;
}

:global([data-theme="dark"] .log-clear-action) {
  border-color: #ef4444 !important;
  background: rgba(239, 68, 68, 0.15) !important;
  color: #ef4444 !important;
}

.section-toggle-icon {
  width: 0.95rem;
  height: 0.95rem;
  transition: transform 0.2s ease;
}

.section-toggle-icon.is-expanded {
  transform: rotate(180deg);
}

.task-chip-list {
  display: grid;
  grid-template-columns: repeat(var(--task-chip-columns, 1), minmax(0, 1fr));
  gap: 0.75rem;
}

.task-chip {
  width: 100%;
  display: inline-flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 0.95rem 1rem;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  background: linear-gradient(180deg, var(--bg-secondary) 0%, var(--card-bg) 100%);
  color: var(--text-primary);
  box-shadow: var(--shadow-sm);
  cursor: pointer;
  transition: transform 0.2s ease, border-color 0.2s ease, box-shadow 0.2s ease;
}

.task-chip:hover {
  transform: translateY(-2px);
  border-color: var(--accent-color);
  box-shadow: var(--shadow-md);
}

.task-chip-name {
  font-weight: 600;
  text-align: left;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.task-chip-status {
  width: 0.6rem;
  height: 0.6rem;
  border-radius: var(--radius-pill);
  display: inline-block;
  flex-shrink: 0;
  background: var(--neutral-soft-bg);
}

.task-chip-status--running {
  background: var(--success-color);
  box-shadow: 0 0 0 3px var(--success-color-bg);
}

.task-chip-status--stopped {
  background: var(--danger-color);
  box-shadow: 0 0 0 3px var(--danger-color-bg);
}

.task-chip-status--error {
  background: var(--warning-color);
  box-shadow: 0 0 0 3px var(--warning-color-bg);
}

.log-loading,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
  min-height: 180px;
  color: var(--text-secondary);
  border: 1px dashed var(--border-color);
  border-radius: var(--radius-lg);
  background: var(--bg-secondary);
}

.compact-empty {
  min-height: 150px;
}

.empty-state svg {
  width: 2rem;
  height: 2rem;
  color: var(--accent-color);
}

.empty-state p {
  margin: 0;
}

.log-scroll-shell {
  position: relative;
  min-width: 0;
}

.log-content {
  max-height: 480px;
  overflow-y: auto;
  padding: 0.5rem 1rem 0.5rem 0.5rem;
  border: 1px solid color-mix(in srgb, var(--border-color) 82%, var(--bg-primary) 18%);
  border-radius: var(--radius-lg);
  background: color-mix(in srgb, var(--card-bg) 78%, var(--bg-secondary) 22%);
  scrollbar-width: thin;
  scrollbar-color: color-mix(in srgb, var(--accent-color) 58%, var(--bg-primary) 42%) color-mix(in srgb, var(--bg-primary) 82%, transparent 18%);
}

.log-content::-webkit-scrollbar {
  width: 10px;
  height: 10px;
}

.log-content::-webkit-scrollbar-thumb {
  background: color-mix(in srgb, var(--accent-color) 58%, var(--bg-primary) 42%);
  border-radius: 999px;
}

.log-content::-webkit-scrollbar-track {
  background: color-mix(in srgb, var(--bg-primary) 82%, transparent 18%);
  border-radius: 999px;
}

.log-line {
  display: grid;
  grid-template-columns: 165px 68px 1fr;
  gap: 0.75rem;
  align-items: start;
  padding: 0.75rem 0.875rem;
  border-radius: var(--radius-md);
  transition: background 0.2s ease, transform 0.2s ease;
}

.log-line + .log-line {
  margin-top: 0.35rem;
}

.log-line:hover {
  background: var(--bg-primary);
  transform: translateX(2px);
}

.log-timestamp {
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
  font-size: 0.8rem;
  color: var(--text-secondary);
}

.log-level {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 1.5rem;
  padding: 0.15rem 0.45rem;
  border-radius: var(--radius-pill);
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: uppercase;
}

.level-info {
  background: var(--log-level-info-bg);
  color: var(--log-level-info-color);
}

.level-warn {
  background: var(--log-level-warn-bg);
  color: var(--log-level-warn-color);
}

.level-error {
  background: rgba(239, 68, 68, 0.12);
  color: var(--danger-color);
}

.log-message {
  color: var(--text-primary);
  word-break: break-word;
  line-height: 1.6;
}

.log-line-info {
  background: var(--log-line-info-bg);
}

.log-line-error {
  background: rgba(239, 68, 68, 0.04);
}

.log-line-warn {
  background: var(--log-line-warn-bg);
}

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1.25rem;
  z-index: 1000;
}

.modal-card {
  width: min(100%, 560px);
  background: var(--card-bg);
  border-radius: var(--radius-lg);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-overlay);
}

.modal-header--dialog {
  position: static;
}

.modal-title {
  margin: 0;
  color: var(--text-primary);
}

.modal-close {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  border-radius: var(--radius-sm);
  color: var(--text-secondary);
}

.modal-close:hover {
  background: var(--bg-primary);
  color: var(--text-primary);
}

.modal-body {
  padding: 1.5rem;
}

.delete-confirm-content {
  text-align: center;
}

.delete-warning-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 3.5rem;
  height: 3.5rem;
  margin-bottom: 1rem;
  border-radius: 50%;
  background: var(--danger-color-bg);
  color: var(--danger-color);
}

.delete-warning-icon svg {
  width: 1.5rem;
  height: 1.5rem;
}

.delete-confirm-title {
  margin: 0 0 0.75rem;
  color: var(--text-primary);
}

.delete-confirm-message {
  margin: 0 0 1rem;
  color: var(--text-primary);
  line-height: 1.6;
}

.delete-confirm-warning {
  margin: 0;
  color: var(--danger-color);
  font-size: 0.9rem;
}

.task-list {
  margin: 1rem 0;
  padding: 0;
  list-style: none;
  text-align: left;
}

.task-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.65rem 0.75rem;
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
  color: var(--text-primary);
}

.task-item + .task-item {
  margin-top: 0.5rem;
}

.task-item svg {
  width: 0.95rem;
  height: 0.95rem;
  color: var(--success-color);
}

.dialog-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 1.25rem;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

@media (max-width: 768px) {
  .detail-container {
    padding: 1rem;
  }

  .card-header {
    flex-wrap: wrap;
    align-items: flex-start;
  }

  .header-center {
    width: 100%;
    flex-direction: column;
    align-items: flex-start;
  }

  .header-actions {
    width: 100%;
    justify-content: flex-start;
  }

  .info-grid {
    grid-template-columns: 1fr;
  }

  .section-heading {
    flex-direction: column;
    align-items: flex-start;
  }

  .section-subheading {
    flex-direction: column;
    align-items: flex-start;
  }

  .log-line {
    grid-template-columns: 1fr;
  }

  .dialog-actions {
    flex-direction: column;
  }
}
</style>
