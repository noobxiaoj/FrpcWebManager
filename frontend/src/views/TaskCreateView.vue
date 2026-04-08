<template>
  <TaskForm
    :title="t('taskCreate.title')"
    :submit-text="t('taskCreate.submit')"
    :submitting-text="t('taskCreate.submitting')"
    :submitting="submitting"
    :initial-form="initialForm"
    @submit="handleSubmit"
    @cancel="goBack"
  />
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useTaskStore } from '@/stores/task'
import TaskForm from '@/components/TaskForm.vue'
import { useI18n } from '@/utils/i18n'

const router = useRouter()
const taskStore = useTaskStore()
const { t } = useI18n()

const submitting = ref(false)
const initialForm = {
  name: '',
  description: '',
  serverId: '',
  serverAddr: '',
  serverPort: 7000,
  authToken: '',
  proxies: []
}

const goBack = () => {
  router.back()
}

/**
 * 处理创建任务提交。
 * 通用表单组件已经完成了字段校验与数据清洗，这里只负责调用创建接口和跳转。
 *
 * @param {object} payload - 通用表单组件整理后的任务数据
 * @returns {Promise<void>} 返回创建流程对应的 Promise
 */
const handleSubmit = async (payload) => {
  submitting.value = true

  try {
    const newTask = await taskStore.createTask({
      name: payload.name,
      description: payload.description,
      serverAddr: payload.serverAddr,
      serverPort: payload.serverPort,
      authToken: payload.authToken,
      proxies: payload.proxies
    })

    router.push(`/tasks/${newTask.id}`)
  } catch (err) {
    alert(`${t('taskCreate.messages.failed')}: ${err.message}`)
  } finally {
    submitting.value = false
  }
}
</script>
