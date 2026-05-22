// src/modules/retail/services/posService.js
import api from '../../../api.js';

export const posService = {
    // --- PRODUK ---
    async getProducts() {
        const response = await api.get('/retail/products');
        return response.data;
    },

    // --- SESI (Buka, Cek, Tutup) ---
    async checkSession(token) {
        // Pake header kalau token ada
        const config = token ? { headers: { Authorization: `Bearer ${token}` } } : {};
        const response = await api.get('/retail/pos/check-session', config);
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

    // --- TRANSAKSI ---
    async checkout(payload) {
        const response = await api.post('/retail/checkout', payload);
        return response.data;
    }
};