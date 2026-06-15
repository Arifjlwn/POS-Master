import api from '../../../../api.js';

export const laundryCatalogService = {
    // Paket Jasa Layanan
    async getServices() {
        const res = await api.get('/laundry/services');
        return res.data || [];
    },
    async createService(payload) {
        return await api.post('/laundry/services', payload);
    },
    async updateService(id, payload) {
        return await api.put(`/laundry/services/${id}`, payload);
    },
    async deleteService(id) {
        return await api.delete(`/laundry/services/${id}`);
    },

    // Add-on Parfum Premium
    async getPerfumes() {
        const res = await api.get('/laundry/perfumes');
        return res.data || [];
    },
    async createPerfume(payload) {
        return await api.post('/laundry/perfumes', payload);
    },
    async deletePerfume(id) {
        return await api.delete(`/laundry/perfumes/${id}`);
    }
};