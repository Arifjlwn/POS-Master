import Swal from 'sweetalert2';
import { ref } from 'vue';
import api from '../../../api.js';

export function useStoreSetting() {
	const isLoading = ref(true);
	const isSaving = ref(false);
	const activeTab = ref('basic');
	const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

	const form = ref({
		nama_toko: '',
		telepon: '',
		business_type: '',
		alamat: '',
		provinsi: '',
		kota: '',
		kecamatan: '',
		kelurahan: '',
		kode_pos: '',
		latitude: 0,
		longitude: 0,
		logo_url: null,
		qris_image: null,
		qris_name: '',
		is_tax_active: false,
		pajak_persen: 0,
		receipt_footer: '',
		payment_type: 'qris_static',
		midtrans_server_key: '',
		midtrans_client_key: '',
		printer_width: '58mm',
		printer_type: 'bluetooth',
		delete_logo: false,
		delete_qris: false,
	});

	const logoPreview = ref(null);
	const qrisPreview = ref(null);

	// HELPER: Mencegah bug string "null" dan mengamankan link Cloud Supabase (https://)
	const getCleanUrl = (url) => {
		if (!url || url === 'null' || url === 'undefined' || url === '') return null;
		if (url.startsWith('http://') || url.startsWith('https://')) return url;
		return API_BASE_URL + url;
	};

	const fetchSettings = async () => {
		isLoading.value = true;
		try {
			const response = await api.get('/retail/store/settings');
			const data = response.data.data;

			Object.assign(form.value, data);

			if (!form.value.payment_type) form.value.payment_type = 'qris_static';
			if (!form.value.printer_width) form.value.printer_width = '58mm';
			if (!form.value.printer_type) form.value.printer_type = 'bluetooth';

			logoPreview.value = getCleanUrl(data.logo_url);
			qrisPreview.value = getCleanUrl(data.qris_image);
		} catch (error) {
			if (error.response && (error.response.status === 402 || error.response.status === 403)) {
				console.warn('Akses ditolak karena lisensi toko tidak aktif.');
				return;
			}
			Swal.fire('Gagal', 'Tidak dapat mengambil data konfigurasi toko', 'error');
		} finally {
			isLoading.value = false;
		}
	};

	const handleFileChange = (type, file, previewUrl) => {
		if (type === 'logo') {
			form.value.logo_url = file;
			logoPreview.value = previewUrl;
			form.value.delete_logo = false;
		} else if (type === 'qris') {
			form.value.qris_image = file;
			qrisPreview.value = previewUrl;
			form.value.delete_qris = false;
		}
	};

	const removeLogo = () => {
		form.value.logo_url = null;
		logoPreview.value = null;
		form.value.delete_logo = true;
	};

	const removeQris = () => {
		form.value.qris_image = null;
		qrisPreview.value = null;
		form.value.delete_qris = true;
	};

	const saveSettings = async () => {
		isSaving.value = true;
		const formData = new FormData();
		const uppercaseFields = ['nama_toko', 'alamat', 'provinsi', 'kota', 'kecamatan', 'kelurahan', 'qris_name', 'receipt_footer'];

		Object.keys(form.value).forEach((key) => {
			if (key !== 'logo_url' && key !== 'qris_image') {
				let val = form.value[key];

				// Pengaman untuk koordinat peta
				if (key === 'latitude' || key === 'longitude') {
					val = parseFloat(val) || 0;
				} else if (typeof val === 'string') {
					val = uppercaseFields.includes(key) ? val.toUpperCase() : val.trim();
				} else if (typeof val === 'boolean') {
					val = String(val);
				}

				formData.append(key, val);
				console.log('DATA YANG AKAN DIKIRIM KE BACKEND:');
				for (let pair of formData.entries()) {
					console.log(pair[0] + ': ' + pair[1]);
				}
			}
		});

		// PENGAMAN BARU: Hanya tempelkan ke form jika user benar-benar mengunggah file baru
		if (form.value.logo_url instanceof File) {
			formData.append('logo', form.value.logo_url);
		}

		if (form.value.qris_image instanceof File) {
			formData.append('qris', form.value.qris_image);
		}

		try {
			const response = await api.put('/retail/store/settings', formData);
			const updatedData = response.data.data;

			if (updatedData.nama_toko) {
				localStorage.setItem('storeName', updatedData.nama_toko);
			}

			if (form.value.delete_logo) {
				localStorage.setItem('storeLogo', '');
			} else if (updatedData.logo_url) {
				localStorage.setItem('storeLogo', updatedData.logo_url);
			}

			window.dispatchEvent(new Event('store-updated'));
			window.dispatchEvent(new Event('storage'));

			Swal.fire({ icon: 'success', title: 'Tersimpan!', text: 'Pengaturan toko berhasil diperbarui.', timer: 2000, showConfirmButton: false });
		} catch (error) {
			Swal.fire('Gagal Menyimpan', error.response?.data?.error || 'Terjadi kesalahan sistem', 'error');
		} finally {
			isSaving.value = false;
		}
	};

	// onMounted dicabut dari sini karena sudah dipanggil langsung di komponen StoreSetting.vue

	// Pastikan fetchSettings terekspor agar bisa di-destructure di komponen
	return { isLoading, isSaving, activeTab, form, logoPreview, qrisPreview, handleFileChange, removeLogo, removeQris, fetchSettings, saveSettings };
}
