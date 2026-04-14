import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src'),
      '@wails': path.resolve(__dirname, 'wailsjs')
    }
  },
  optimizeDeps: {
    include: ['naive-ui']
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          'naive-ui': ['naive-ui'],
          'monaco-editor': ['monaco-editor'],
          'vue-vendor': ['vue', 'pinia']
        }
      }
    }
  }
})
