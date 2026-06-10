import axios from 'axios';
import Swal from 'sweetalert2';

const api = axios.create({
	baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
	withCredentials: true,
});

// Otomatis bawa tiket JWT setiap kali request dikirim
api.interceptors.request.use(
	async (config) => {
		const token = localStorage.getItem('token');

		if (token) {
			config.headers.Authorization = `Bearer ${token}`;
		}

		return config;
	},
	(error) => {
		return Promise.reject(error);
	}
);

// 🚀 SECURITY & UX PATCH: Response Interceptor untuk menangani Token Expired (401)
api.interceptors.response.use(
	(response) => response,
	(error) => {
		const originalRequest = error.config;

		// Cek jika response ada dan statusnya 401 Unauthorized
		if (error.response && error.response.status === 401) {
			if (originalRequest.url && originalRequest.url.includes('/login')) {
				return Promise.reject(error);
			}
			// =======================================================================

			// 🛒 Jalur Normal Kasir (Kalau lagi di dashboard ruko, baru jalankan proteksi logout otomatis ini)
			localStorage.removeItem('token');

			Swal.fire({
				icon: 'error',
				title: 'Sesi Berakhir',
				text: 'Silakan login kembali untuk melanjutkan.',
				confirmButtonColor: '#4f46e5',
				customClass: { popup: 'rounded-[32px]' },
			}).then(() => {
				window.location.href = '/login';
			});
		}
		return Promise.reject(error);
	}
);

export default api;
