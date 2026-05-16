import { createApp } from 'vue'
import App from './App.vue'
import router from './router' // Panggil Router
import './style.css' // Panggil Tailwind
import { autoAnimatePlugin } from '@formkit/auto-animate/vue'

const app = createApp(App)

app.use(router)
app.use(autoAnimatePlugin)
app.mount('#app')