import api from '../../../api.js'; // Sesuaikan path-nya ke file api.js lu

export const fnbService = {
    // Master Menu
    async getProducts() {
        const res = await api.get('/fnb/products');
        return res.data;
    },

    async createProduct(payload) {
        const res = await api.post('/fnb/products', payload);
        return res.data;
    },

    async updateProduct(id, payload) {
        const res = await api.put(`/fnb/products/${id}`, payload);
        return res.data;
    },

    async deleteProduct(id) {
        const res = await api.delete(`/fnb/products/${id}`);
        return res.data;
    },

    async toggleStatus(id) {
        const res = await api.put(`/fnb/products/${id}/toggle`);
        return res.data;
    },

    // Order & Kasir
    async submitOrder(payload) {
        const res = await api.post('/fnb/order', payload);
        return res.data;
    },
    
    // Kitchen (Kalau nanti lu mau rapiin dapur juga)
    async getKitchenQueue() {
        const res = await api.get('/fnb/kitchen');
        return res.data;
    },

    async completeOrder(id) {
        const res = await api.put(`/fnb/kitchen/${id}/selesai`);
        return res.data;
    }
};