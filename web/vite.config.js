import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue()
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:11223', // 后端地址
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '/api/v1'),
        timeout: 300000, // 5分钟超时
        proxyTimeout: 300000,
        headers: {
          'Connection': 'keep-alive'
        }
      }
    }
  }
})
