import { defineConfig } from 'vite'
import uni from '@dcloudio/vite-plugin-uni'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    uni(),
  ],
  css: {
    preprocessorOptions: {
      scss: {
        silenceDeprecations: ["legacy-js-api", "import", "global-builtin"],
      },
    },
  },
  server: {
    proxy: {
      // 代理 /api/v1 开头的请求到后端服务
      '/api/v1': {
        target: 'http://localhost:8080', // 后端本地开发地址
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/v1/, '/api/v1')
      },
    }
  }
})