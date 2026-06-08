import { autoAnimatePlugin } from '@formkit/auto-animate/vue';
import { createApp } from 'vue';
import App from './App.vue';
import router from './router'; // Panggil Router
import './style.css'; // Panggil Tailwind

// --- DAFTARKAN SERVICE WORKER UNTUK PWA ---
if ('serviceWorker' in navigator) {
	window.addEventListener('load', () => {
		navigator.serviceWorker
			.register('/sw.js')
			.then(() => console.log('Arzura POS PWA: Service Worker Aktif'))
			.catch((err) => console.log('Arzura POS PWA: Service Worker Gagal', err));
	});
}
// ------------------------------------------

const app = createApp(App);

app.use(router);
app.use(autoAnimatePlugin);
app.mount('#app');
