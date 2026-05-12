import { createApp } from 'vue'
import App from './App.vue'
import router from './router' // Panggil Router
import './style.css' // Panggil Tailwind

const app = createApp(App)

app.use(router)
app.mount('#app')