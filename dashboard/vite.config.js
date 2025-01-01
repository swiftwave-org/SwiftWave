import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  base: '/dashboard/',
  server: {
    proxy: {
     // Proxy requests not starting with /dashboard to localhost:3333
     '^/(?!dashboard).*': {
        target: 'http://localhost:3333',
        changeOrigin: true,
        rewrite: (path) => path,
      },
    },
  },
  plugins: [
    vue(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  }
})
