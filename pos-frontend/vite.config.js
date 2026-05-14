import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    host: true, // Supaya bisa diakses via HP (IP lokal)
    port: 5173,
    proxy: {
      // 🚀 OPER JALUR STATIC UPLOADS KE BACKEND GO
      '/uploads': {
        target: 'http://localhost:8080', // Sesuaikan dengan port/IP backend Go Mas Arif
        changeOrigin: true,
        secure: false,
      },
      // Jalur API biasa (jika ada)
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      }
    }
  }
})