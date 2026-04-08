<template>
  <TaskForm
    :title="t('taskEdit.title')"
    :submit-text="t('taskEdit.submit')"
    :submitting-text="t('taskEdit.submitting')"
    :loading-text="t('taskEdit.loading')"
    :loading="loading"
    :submitting="submitting"
    :initial-form="initialForm"
    @submit="handleSubmit"
    @cancel="goBack"
  />
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useTaskStore } from '@/stores/task'
import TaskForm from '@/components/TaskForm.vue'
import { useI18n } from '@/utils/i18n'

const router = useRouter()
const route = useRoute()
const taskStore = useTaskStore()
const { t } = useI18n()

const loading = ref(true)
const submitting = ref(false)
const taskId = route.params.id
const taskStatus = ref('') // 保存任务状态,用于判断更新后是否需要重载
const initialForm = ref({
  name: '',
  description: '',
  serverId: '',
  serverAddr: '',
  serverPort: 7000,
  authToken: '',
  proxies: []
})

const goBack = () => {
  router.push(`/tasks/${taskId}`)
}

/**
 * 加载任务详情，并将其转换为通用表单组件需要的初始值。
 * 编辑页只负责取数和场景特有逻辑，表单展示与基础数据整理交给通用组件处理。
 *
 * @returns {Promise<void>} 返回任务加载流程
 */
const loadTask = async () => {
  try {
    const task = await taskStore.fetchTask(taskId)
    taskStatus.value = task.status || ''

    initialForm.value = {
      name: task.name || '',
      description: task.description || '',
      serverAddr: task.serverAddr || '',
      serverPort: task.serverPort || 7000,
      authToken: task.authToken || '',
      proxies: (task.proxies || []).map(proxy => ({
        name: proxy.name || '',
        type: proxy.type || 'tcp',
        localIP: proxy.localIP || '127.0.0.1',
        localPort: proxy.localPort || '',
        remotePort: proxy.remotePort || '',
        customDomains: Array.isArray(proxy.customDomains) ? proxy.customDomains.join(', ') : '',
        subdomain: proxy.subdomain || '',
        isEditing: false,
        isValid: true,
        errors: {}
      }))
    }
  } catch (err) {
    alert(`${t('taskEdit.messages.loadFailed')}: ${err.message}`)
    router.push('/tasks')
  } finally {
    loading.value = false
  }
}

/**
 * 处理编辑任务提交。
 * 除了更新任务外，这里还保留了“运行中任务修改后自动重载”的页面特有行为。
 *
 * @param {object} payload - 通用表单组件整理后的任务数据
 * @returns {Promise<void>} 返回更新流程对应的 Promise
 */
const handleSubmit = async (payload) => {
  submitting.value = true

  try {
    await taskStore.updateTask(taskId, {
      name: payload.name,
      description: payload.description,
      serverAddr: payload.serverAddr,
      serverPort: payload.serverPort,
      authToken: payload.authToken,
      proxies: payload.proxies
    })

    if (taskStatus.value === 'running') {
      try {
        await taskStore.reloadTask(taskId)
        console.log(t('taskEdit.messages.autoReloaded'))
      } catch (reloadError) {
        console.error(`${t('taskEdit.messages.autoReloadFailed')}:`, reloadError)
      }
    }

    router.push(`/tasks/${taskId}`)
  } catch (err) {
    alert(`${t('taskEdit.messages.updateFailed')}: ${err.message}`)
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadTask()
})
</script>
