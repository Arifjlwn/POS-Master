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

// 🛠️ TAMBAHKAN KODE INI UNTUK MEMPERBAIKI EROR BANNER PWA
window.addEventListener('beforeinstallprompt', (e) => {
    // Biarkan browser tahu kita akan menghandle install prompt-nya
    e.preventDefault();
    
    // Simpan event-nya ke window supaya tidak hilang dan erornya lenyap
    window.deferredPrompt = e;
    
    // (Opsional) Jika ingin langsung memunculkan banner bawaan browser tanpa eror:
    // e.prompt();
});
// ------------------------------------------

const app = createApp(App);

app.use(router);
app.use(autoAnimatePlugin);
app.mount('#app');
