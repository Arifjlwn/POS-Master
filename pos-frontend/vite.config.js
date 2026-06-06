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
        target: 'http://127.0.0.1:8080', // Gunakan IP loopback murni atau IP komputer lu 
        changeOrigin: true,
        secure: false,
      },
      // 🚀 FIX MUTLAK: Arahkan target ke 127.0.0.1 agar Proxy Vite bisa meneruskan request HP dengan mulus
      '/api': {
        target: 'http://127.0.0.1:8080', 
        changeOrigin: true,
      }
    }
  }
})