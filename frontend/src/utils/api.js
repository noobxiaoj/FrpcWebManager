/**
 * API 工具函数
 */

/**
 * 获取 API 基础 URL
 * 根据环境自动选择 API 地址
 * @returns {string} API 基础 URL
 */
export function getApiBaseUrl() {
  // 优先使用环境变量配置的 API 地址（开发环境）
  if (import.meta.env.VITE_API_URL) {
    return import.meta.env.VITE_API_URL
  }

  // 生产环境：使用相对路径（前后端同服务）
  return ''
}
