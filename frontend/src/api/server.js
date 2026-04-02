// 服务器API客户端
// 生产环境使用相对路径，开发环境可通过环境变量配置
const getApiBaseUrl = () => {
  // 如果环境变量中指定了 API 地址,直接使用（开发环境）
  if (import.meta.env.VITE_API_URL) {
    return import.meta.env.VITE_API_URL
  }

  // 生产环境：使用相对路径（前后端同服务）
  return ''
}

const API_BASE_URL = getApiBaseUrl()

/**
 * 服务器API
 */
export const serverAPI = {
  /**
   * 获取服务器列表
   */
  async listServers() {
    try {
      const response = await fetch(`${API_BASE_URL}/api/servers`, {
        credentials: 'include'
      })
      const result = await response.json()
      if (result.code === 1004) {
        window.dispatchEvent(new CustomEvent('auth-expired'))
      }
      if (result.code === 0) {
        return result.data.servers || []
      }
      throw new Error(result.message || '获取服务器列表失败')
    } catch (error) {
      console.error('获取服务器列表失败:', error)
      throw error
    }
  },

  /**
   * 获取服务器详情
   */
  async getServer(id) {
    try {
      const response = await fetch(`${API_BASE_URL}/api/servers/${id}`, {
        credentials: 'include'
      })
      const result = await response.json()
      if (result.code === 1004) {
        window.dispatchEvent(new CustomEvent('auth-expired'))
      }
      if (result.code === 0) {
        return result.data
      }
      throw new Error(result.message || '获取服务器详情失败')
    } catch (error) {
      console.error('获取服务器详情失败:', error)
      throw error
    }
  },

  /**
   * 创建服务器
   */
  async createServer(serverData) {
    try {
      const response = await fetch(`${API_BASE_URL}/api/servers`, {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(serverData)
      })
      const result = await response.json()
      if (result.code === 1004) {
        window.dispatchEvent(new CustomEvent('auth-expired'))
      }
      if (result.code === 0) {
        return result.data
      }
      throw new Error(result.message || '创建服务器失败')
    } catch (error) {
      console.error('创建服务器失败:', error)
      throw error
    }
  },

  /**
   * 更新服务器
   */
  async updateServer(id, serverData) {
    try {
      const response = await fetch(`${API_BASE_URL}/api/servers/${id}`, {
        method: 'PUT',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(serverData)
      })
      const result = await response.json()
      if (result.code === 1004) {
        window.dispatchEvent(new CustomEvent('auth-expired'))
      }
      if (result.code === 0) {
        return result.data
      }
      throw new Error(result.message || '更新服务器失败')
    } catch (error) {
      console.error('更新服务器失败:', error)
      throw error
    }
  },

  /**
   * 删除服务器
   * @param {string} id - 服务器ID
   * @param {boolean} force - 是否强制删除（同时删除关联任务）
   */
  async deleteServer(id, force = false) {
    try {
      const url = `${API_BASE_URL}/api/servers/${id}${force ? '?force=true' : ''}`
      const response = await fetch(url, {
        method: 'DELETE',
        credentials: 'include'
      })
      const result = await response.json()
      if (result.code === 1004) {
        window.dispatchEvent(new CustomEvent('auth-expired'))
      }

      // 如果服务器有任务，返回任务信息
      if (result.code === 1 && result.data?.hasTasks) {
        return {
          hasTasks: true,
          tasks: result.data.tasks,
          taskCount: result.data.taskCount
        }
      }

      if (result.code === 0) {
        return {
          success: true,
          tasks: result.data?.tasks || []
        }
      }

      throw new Error(result.message || '删除服务器失败')
    } catch (error) {
      console.error('删除服务器失败:', error)
      throw error
    }
  },

  /**
   * 更新服务器排序
   */
  async updateServerOrder(serverOrder) {
    try {
      const response = await fetch(`${API_BASE_URL}/api/servers/order`, {
        method: 'PUT',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ order: serverOrder })
      })
      const result = await response.json()
      if (result.code === 1004) {
        window.dispatchEvent(new CustomEvent('auth-expired'))
      }
      if (result.code === 0) {
        return true
      }
      throw new Error(result.message || '更新服务器排序失败')
    } catch (error) {
      console.error('更新服务器排序失败:', error)
      throw error
    }
  },

  /**
   * 更新服务器锁定状态
   */
  async updateServerLock(id, locked) {
    try {
      const response = await fetch(`${API_BASE_URL}/api/servers/${id}/lock`, {
        method: 'PUT',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ locked })
      })
      const result = await response.json()
      if (result.code === 1004) {
        window.dispatchEvent(new CustomEvent('auth-expired'))
      }
      if (result.code === 0) {
        return true
      }
      throw new Error(result.message || '更新服务器锁定状态失败')
    } catch (error) {
      console.error('更新服务器锁定状态失败:', error)
      throw error
    }
  },

  /**
   * 添加日志
   */
  async addLog(id, log) {
    try {
      const response = await fetch(`${API_BASE_URL}/api/servers/${id}/logs`, {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(log)
      })
      const result = await response.json()
      if (result.code === 1004) {
        window.dispatchEvent(new CustomEvent('auth-expired'))
      }
      if (result.code === 0) {
        return true
      }
      throw new Error(result.message || '添加日志失败')
    } catch (error) {
      console.error('添加日志失败:', error)
      throw error
    }
  },

  /**
   * 清空日志
   */
  async clearLogs(id) {
    try {
      const response = await fetch(`${API_BASE_URL}/api/servers/${id}/logs`, {
        method: 'DELETE',
        credentials: 'include'
      })
      const result = await response.json()
      if (result.code === 1004) {
        window.dispatchEvent(new CustomEvent('auth-expired'))
      }
      if (result.code === 0) {
        return true
      }
      throw new Error(result.message || '清空日志失败')
    } catch (error) {
      console.error('清空日志失败:', error)
      throw error
    }
  },

  /**
   * 获取服务器日志
   */
  async getLogs(id, limit = 100) {
    try {
      const response = await fetch(`${API_BASE_URL}/api/servers/${id}/logs?limit=${limit}`, {
        credentials: 'include'
      })
      const result = await response.json()
      if (result.code === 1004) {
        window.dispatchEvent(new CustomEvent('auth-expired'))
      }
      if (result.code === 0) {
        return result.data.logs || []
      }
      throw new Error(result.message || '获取日志失败')
    } catch (error) {
      console.error('获取日志失败:', error)
      throw error
    }
  }
}
