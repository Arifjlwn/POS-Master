import { ref, onMounted } from 'vue';
import api from '../../../api.js';
import Swal from 'sweetalert2';

export function useStoreSetting() {
    const isLoading = ref(true);
    const isSaving = ref(false);
    const activeTab = ref('basic'); // basic, payment, tax
    const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

    const form = ref({
        nama_toko: '', telepon: '', business_type: '',
        alamat: '', provinsi: '', kota: '', kecamatan: '', kelurahan: '', kode_pos: '',
        logo_url: null, qris_image: null, qris_name: '',
        is_tax_active: false, pajak_persen: 0, receipt_footer: ''
    });

    const logoPreview = ref(null);
    const qrisPreview = ref(null);

    const fetchStoreSettings = async () => {
        isLoading.value = true;
        try {
            // Asumsi rute GET /retail/store/settings sudah lu buat di Golang
            const response = await api.get('/retail/store/settings');
            const data = response.data.data;
            
            form.value = { ...form.value, ...data };
            
            if (data.logo_url) logoPreview.value = API_BASE_URL + data.logo_url;
            if (data.qris_image) qrisPreview.value = API_BASE_URL + data.qris_image;
        } catch (error) {
            Swal.fire('Gagal', 'Tidak dapat mengambil data toko', 'error');
        } finally {
            isLoading.value = false;
        }
    };

    const handleFileChange = (type, file, previewUrl) => {
        if (type === 'logo') {
            form.value.logo_url = file;
            logoPreview.value = previewUrl;
        } else if (type === 'qris') {
            form.value.qris_image = file;
            qrisPreview.value = previewUrl;
        }
    };

    const saveSettings = async () => {
        isSaving.value = true;
        const formData = new FormData();
        
        Object.keys(form.value).forEach(key => {
            if (key !== 'logo_url' && key !== 'qris_image') {
                formData.append(key, form.value[key]);
            }
        });

        if (form.value.logo_url instanceof File) formData.append('logo', form.value.logo_url);
        if (form.value.qris_image instanceof File) formData.append('qris', form.value.qris_image);

        try {
            // 🚀 TANGKAP RESPONSE-NYA
            const response = await api.put('/retail/store/settings', formData);
            
            // 🚀 UPDATE LOCAL STORAGE LANGSUNG BIAR SIDEBAR & HEADER BERUBAH
            const updatedData = response.data.data;
            localStorage.setItem('storeName', updatedData.nama_toko || '');
            localStorage.setItem('storeLogo', updatedData.logo_url || '');
            
            // 🚀 PICU EVENT SAKTI BIAR SIDEBAR LANGSUNG RE-RENDER TANPA REFRESH
            window.dispatchEvent(new Event('storage'));

            Swal.fire({ icon: 'success', title: 'Tersimpan!', text: 'Pengaturan toko berhasil diperbarui.', timer: 2000, showConfirmButton: false });
        } catch (error) {
            Swal.fire('Gagal Menyimpan', error.response?.data?.error || 'Terjadi kesalahan', 'error');
        } finally {
            isSaving.value = false;
        }
    };

    onMounted(fetchStoreSettings);

    return {
        isLoading, isSaving, activeTab, form, logoPreview, qrisPreview,
        handleFileChange, saveSettings
    };
}