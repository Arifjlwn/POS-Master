import api from '../../../api.js';

export const settingService = {
    // Ambil data konfigurasi toko laundry dari backend Go bray
    async getSettings() {
        const response = await api.get('/laundry/setting');
        return response.data;
    },

    // Kirim perubahan profile ruko ke database
    async updateSettings(payload) {
        const response = await api.put('/laundry/setting', payload);
        return response.data;
    }
};