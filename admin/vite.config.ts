import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': '/src'
    }
  },
  css: {
    preprocessorOptions: {
      scss: {
        silenceDeprecations: ['legacy-js-api']
      }
    }
  },
  server: {
    host: '0.0.0.0', // 监听所有网络接口，允许外部访问
    port: 4000, // 开发服务器端口
    proxy: { // 代理 /api 开头的请求到后端服务
      '/api': { // 代理 /api 开头的请求到后端服务
        target: 'http://localhost:3300', // 后端本地开发地址
        changeOrigin: true, // 改变源地址
      }
    }
  },
  build: {
    sourcemap: false,
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true,
        drop_debugger: true
      }
    },
    chunkSizeWarningLimit: 1000,
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['react', 'react-dom', 'react-router-dom'],
          antd: ['antd'],
          axios: ['axios'],
          xlsx: ['xlsx'],
          fileSaver: ['file-saver']
        }
      }
    }
  }
})