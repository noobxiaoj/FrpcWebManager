import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    host: '0.0.0.0', // 允许外部访问
    port: 5173,
    strictPort: true, // 如果端口被占用则失败而不是尝试下一个端口
    // 如果需要 HTTPS，可以取消下面的注释
    // https: true
  },
  build: {
    // 使用 esbuild 进行压缩（无需额外安装依赖）
    minify: 'esbuild',
    // 生产环境移除 console
    target: 'es2015'
  }
})
