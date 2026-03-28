// API 基础配置
// 容器生产环境：前后端同源，使用相对路径
// 开发环境：可通过环境变量配置
const getApiBaseUrl = () => {
  // 如果环境变量中指定了 API 地址,直接使用（开发环境）
  if (import.meta.env.VITE_API_URL) {
    return import.meta.env.VITE_API_URL
  }

  // 生产环境：使用相对路径（前后端同服务）
  return ''
}

const API_BASE_URL = getApiBaseUrl()

// API 请求封装
async function request(url, options = {}) {
  const config = {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers
    }
  }

  try {
    const response = await fetch(`${API_BASE_URL}${url}`, config)
    const data = await response.json()

    // 检查业务状态码
    if (data.code !== 0 && data.code !== undefined) {
      throw new Error(data.message || data.error || '请求失败')
    }

    return data
  } catch (error) {
    console.error('API 请求错误:', error)
    throw error
  }
}

// 任务管理 API
export const taskApi = {
  // 获取任务列表
  list: () => request('/api/tasks'),

  // 获取任务详情
  get: (id) => request(`/api/tasks/${id}`),

  // 创建任务
  create: (data) => request('/api/tasks', {
    method: 'POST',
    body: JSON.stringify(data)
  }),

  // 更新任务
  update: (id, data) => request(`/api/tasks/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data)
  }),

  // 删除任务
  delete: (id) => request(`/api/tasks/${id}`, {
    method: 'DELETE'
  }),

  // 启动任务
  start: (id) => request(`/api/tasks/${id}/start`, {
    method: 'POST'
  }),

  // 停止任务
  stop: (id) => request(`/api/tasks/${id}/stop`, {
    method: 'POST'
  }),

  // 重载任务
  reload: (id) => request(`/api/tasks/${id}/reload`, {
    method: 'POST'
  }),

  // 获取任务状态
  status: (id) => request(`/api/tasks/${id}/status`),

  // 获取任务配置
  config: (id) => request(`/api/tasks/${id}/config`)
}

export default request
