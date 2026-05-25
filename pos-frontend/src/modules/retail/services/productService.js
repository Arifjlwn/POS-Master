import api from '../../../api.js';

export const productService = {
    // 🚀 Ambil katalog produk retail (Server-side Pagination & Filter)
    async getProducts(params) {
        return await api.get('/retail/products', { params });
    },

    // 🚀 Ambil daftar string kategori unik untuk dropdown filter
    async getCategories() {
        return await api.get('/retail/categories');
    },

    // 🚀 Simpan produk baru (Multipart karena ada upload file gambar)
    async createProduct(formData) {
        return await api.post('/retail/products', formData, {
            headers: { 'Content-Type': 'multipart/form-data' }
        });
    },

    // 🚀 Update data produk lama via form id (Multipart Form-Data)
    async updateProduct(id, formData) {
        return await api.put(`/retail/products/${id}`, formData, {
            headers: { 'Content-Type': 'multipart/form-data' }
        });
    },

    // 🚀 Hapus produk permanen dari sistem gudang
    async deleteProduct(id) {
        return await api.delete(`/retail/products/${id}`);
    },

    // 🚀 Unggah file CSV untuk mass import data produk
    async importCSV(formData) {
        return await api.post('/retail/products/import', formData, {
            headers: { 'Content-Type': 'multipart/form-data' }
        });
    }
};