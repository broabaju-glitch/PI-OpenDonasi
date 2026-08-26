import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

export default defineConfig({
  plugins: [vue()],

  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },

  server: {
    port: 5174,
    strictPort: true,

    // ── Proxy: forward /api/* → backend on :8080 ──────────────────
    // This eliminates ALL CORS issues during development.
    // The browser only ever talks to localhost:5174; Vite proxies to :8080.
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
        // Uncomment if backend doesn't mount routes under /api prefix:
        // rewrite: (p) => p.replace(/^\/api/, '')
      },
    },
  },
})
