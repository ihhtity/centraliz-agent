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
    host: '0.0.0.0', // 监听所有网络接口，允许外部访问
    port: 4100, // 开发服务器端口
    proxy: {
      // 代理 /api/v1 开头的请求到后端服务
      '/api/v1': {
        target: 'http://localhost:3300', // 后端本地开发地址
        changeOrigin: true,
        secure: false, // 关闭SSL验证
        rewrite: (path) => path.replace(/^\/api\/v1/, '/api/v1')
      },
    }
  }
})