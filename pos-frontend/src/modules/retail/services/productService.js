import api from '../../../api.js';

export const productService = {
	// 🚀 Ambil katalog produk retail (Server-side Pagination & Filter)
	getProducts: (params) => api.get('/retail/products', { params }),

	// 🚀 Ambil daftar string kategori unik untuk dropdown filter
	getCategories: () => api.get('/retail/categories'),

	// 🚀 Simpan produk baru (Multipart karena ada upload file gambar)
	createProduct: (formData) => api.post('/retail/products', formData, { headers: { 'Content-Type': 'multipart/form-data' } }),

	// 🚀 Update data produk lama via ULID publik (Multipart Form-Data)
	updateProduct: (publicId, formData) => api.put(`/retail/products/${publicId}`, formData, { headers: { 'Content-Type': 'multipart/form-data' } }),

	// 🚀 Hapus produk permanen dari sistem gudang via ULID publik
	deleteProduct: (publicId) => api.delete(`/retail/products/${publicId}`),

	// 🚀 Unggah file CSV untuk mass import data produk
	importCSV: (formData) => api.post('/retail/products/import', formData, { headers: { 'Content-Type': 'multipart/form-data' } }),

	// 🚀 Download file CSV massal untuk master data produk (🔐 SAFE ACCESSED VIA AXIOS BLOB)
	exportProducts: () => api.get('/retail/products/export', { responseType: 'blob' }),
};
