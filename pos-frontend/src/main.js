import { createApp } from 'vue'
import './style.css' // Panggil Tailwind
import App from './App.vue'
import router from './router' // Panggil Router

createApp(App).use(router).mount('#app')