import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { taskApi } from '@/api'
import { serverAPI } from '@/api/server'

export const useTaskStore = defineStore('task', () => {
  // 状态
  const tasks = ref([])
  const currentTask = ref(null)
  const loading = ref(false)
  const error = ref(null)
  const servers = ref([])
  const showServerName = ref(false) // 是否显示服务器名称而非IP
  const taskOrder = ref([]) // 任务排序ID列表

  // 从本地存储加载排序
  const loadTaskOrder = () => {
    try {
      const saved = localStorage.getItem('taskOrder')
      if (saved) {
        taskOrder.value = JSON.parse(saved)
      }
    } catch (err) {
      console.error('加载任务排序失败:', err)
    }
  }

  // 保存排序到本地存储
  const saveTaskOrder = () => {
    try {
      localStorage.setItem('taskOrder', JSON.stringify(taskOrder.value))
    } catch (err) {
      console.error('保存任务排序失败:', err)
    }
  }

  // 获取排序后的任务列表
  const sortedTasks = computed(() => {
    if (taskOrder.value.length === 0) {
      return tasks.value
    }
    // 根据排序ID创建有序的任务列表
    const ordered = []
    const remaining = new Set(tasks.value.map(t => t.id))

    // 先添加已排序的任务
    for (const id of taskOrder.value) {
      const task = tasks.value.find(t => t.id === id)
      if (task) {
        ordered.push(task)
        remaining.delete(id)
      }
    }

    // 添加新任务(不在排序中的)
    for (const task of tasks.value) {
      if (remaining.has(task.id)) {
        ordered.push(task)
      }
    }

    return ordered
  })

  // 获取任务列表
  const fetchTasks = async () => {
    loading.value = true
    error.value = null
    try {
      const response = await taskApi.list()
      tasks.value = response.data.tasks || []
      loadTaskOrder()
    } catch (err) {
      error.value = err.message
      console.error('获取任务列表失败:', err)
    } finally {
      loading.value = false
    }
  }

  // 获取任务详情
  const fetchTask = async (id) => {
    loading.value = true
    error.value = null
    try {
      const response = await taskApi.get(id)
      currentTask.value = response.data.task
      return response.data.task
    } catch (err) {
      error.value = err.message
      console.error('获取任务详情失败:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  // 创建任务
  const createTask = async (taskData) => {
    loading.value = true
    error.value = null
    try {
      const response = await taskApi.create(taskData)
      const newTask = response.data.task
      tasks.value.push(newTask)
      // 新任务添加到排序末尾
      taskOrder.value.push(newTask.id)
      saveTaskOrder()
      return newTask
    } catch (err) {
      error.value = err.message
      console.error('创建任务失败:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  // 更新任务
  const updateTask = async (id, taskData) => {
    loading.value = true
    error.value = null
    try {
      const response = await taskApi.update(id, taskData)
      const updatedTask = response.data.task
      const index = tasks.value.findIndex(t => t.id === id)
      if (index !== -1) {
        tasks.value[index] = updatedTask
      }
      return updatedTask
    } catch (err) {
      error.value = err.message
      console.error('更新任务失败:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  // 删除任务
  const deleteTask = async (id) => {
    loading.value = true
    error.value = null
    try {
      await taskApi.delete(id)
      tasks.value = tasks.value.filter(t => t.id !== id)
      // 从排序中移除
      taskOrder.value = taskOrder.value.filter(taskId => taskId !== id)
      saveTaskOrder()
      if (currentTask.value?.id === id) {
        currentTask.value = null
      }
    } catch (err) {
      error.value = err.message
      console.error('删除任务失败:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  // 启动任务
  const startTask = async (id) => {
    loading.value = true
    error.value = null
    try {
      await taskApi.start(id)
      // 刷新任务列表以获取最新状态
      await fetchTasks()
    } catch (err) {
      error.value = err.message
      console.error('启动任务失败:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  // 停止任务
  const stopTask = async (id) => {
    loading.value = true
    error.value = null
    try {
      await taskApi.stop(id)
      // 刷新任务列表以获取最新状态
      await fetchTasks()
    } catch (err) {
      error.value = err.message
      console.error('停止任务失败:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  // 重载任务
  const reloadTask = async (id) => {
    loading.value = true
    error.value = null
    try {
      await taskApi.reload(id)
      // 刷新任务列表以获取最新状态
      await fetchTasks()
    } catch (err) {
      error.value = err.message
      console.error('重载任务失败:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  // 获取任务状态
  const fetchTaskStatus = async (id) => {
    try {
      const response = await taskApi.status(id)
      return response.data
    } catch (err) {
      console.error('获取任务状态失败:', err)
      throw err
    }
  }

  // 获取任务配置
  const fetchTaskConfig = async (id) => {
    try {
      const response = await taskApi.config(id)
      return response.data
    } catch (err) {
      console.error('获取任务配置失败:', err)
      throw err
    }
  }

  // 根据 ID 获取任务
  const getTaskById = (id) => {
    return tasks.value.find(t => t.id === id)
  }

  // 获取运行中的任务
  const getRunningTasks = () => {
    return tasks.value.filter(t => t.status === 'running')
  }

  // 获取服务器列表
  const fetchServers = async () => {
    try {
      servers.value = await serverAPI.listServers()
    } catch (err) {
      console.error('获取服务器列表失败:', err)
    }
  }

  // 根据IP和端口查找服务器名称
  const getServerNameByAddress = (serverAddr, serverPort) => {
    const server = servers.value.find(s => {
      const [addr, port] = s.address.split(':')
      return addr === serverAddr && parseInt(port) === serverPort
    })
    return server?.name || null
  }

  // 切换显示模式
  const toggleDisplayMode = () => {
    showServerName.value = !showServerName.value
  }

  // 更新任务排序
  const updateTaskOrder = (newOrder) => {
    taskOrder.value = newOrder
    saveTaskOrder()
  }

  return {
    // 状态
    tasks,
    currentTask,
    loading,
    error,
    servers,
    showServerName,
    taskOrder,
    sortedTasks,
    // 方法
    fetchTasks,
    fetchTask,
    createTask,
    updateTask,
    deleteTask,
    startTask,
    stopTask,
    reloadTask,
    fetchTaskStatus,
    fetchTaskConfig,
    getTaskById,
    getRunningTasks,
    fetchServers,
    getServerNameByAddress,
    toggleDisplayMode,
    updateTaskOrder,
    loadTaskOrder,
    saveTaskOrder
  }
})
