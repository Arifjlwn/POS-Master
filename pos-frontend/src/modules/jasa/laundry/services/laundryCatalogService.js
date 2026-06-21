import api from '../../../../api.js';

export const laundryCatalogService = {
	// Paket Jasa Layanan
	async getServices() {
		const res = await api.get('/laundry/services');
		return res.data.data || res.data || [];
	},
	async createService(payload) {
		const res = await api.post('/laundry/services', payload);
		return res.data;
	},
	async updateService(id, payload) {
		const res = await api.put(`/laundry/services/${id}`, payload);
		return res.data;
	},
	async deleteService(id) {
		const res = await api.delete(`/laundry/services/${id}`);
		return res.data;
	},

	// Add-on Parfum Premium
	async getPerfumes() {
		const res = await api.get('/laundry/perfumes');
		return res.data.data || res.data || [];
	},
	async createPerfume(payload) {
		const res = await api.post('/laundry/perfumes', payload);
		return res.data;
	},
	async deletePerfume(id) {
		const res = await api.delete(`/laundry/perfumes/${id}`);
		return res.data;
	},
};
