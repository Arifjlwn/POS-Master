// src/modules/retail/services/posService.js
import api from '../../../api.js';

export const posService = {
	// --- MASTER DATA PRODUK ---
	async getProducts() {
		const response = await api.get('/retail/products');
		return response.data;
	},

	// --- SESI AKTIF (Buka, Cek, Tutup Sesi Laci Kasir / Cash Drawer) ---
	// FIX AMAN: Buang parameter token manual. Biarkan api.js (Axios Interceptor) global lu
	// yang otomatis menyuntikkan tiket JWT Bearer Token secara terpusat dan aman!
	async checkSession() {
		const response = await api.get('/retail/pos/check-session');
		return response.data;
	},

	async openSession(payload) {
		const response = await api.post('/retail/pos/open-session', payload);
		return response.data;
	},

	async closeSession(sessionId, payload) {
		const response = await api.post(`/retail/pos/close-session/${sessionId}`, payload);
		return response.data;
	},

	// --- TRANSAKSI CHECKOUT MESIN KASIR ---
	async checkout(payload) {
		// FIX ARSITEKTUR SaaS: Arahkan ke rute terproteksi /retail/pos/checkout
		// supaya kasir ga bisa nembak transaksi belanja kalau laci kasir belum di-inisialisasi (Bypass Modal Awal)!
		const response = await api.post('/retail/pos/checkout', payload);
		return response.data;
	},
};
