import { ref, onMounted } from 'vue';
import api from '../../../api.js';
import Swal from 'sweetalert2';

export function useStoreSetting() {
    const isLoading = ref(true);
    const isSaving = ref(false);
    const activeTab = ref('basic'); // basic, payment, tax
    const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

    // 🚀 1. STATE FORM (UDAH DITAMBAHIN VARIABEL MIDTRANS, PRINTER & DELETE FLAG)
    const form = ref({
        nama_toko: '', telepon: '', business_type: '',
        alamat: '', provinsi: '', kota: '', kecamatan: '', kelurahan: '', kode_pos: '',
        logo_url: null, qris_image: null, qris_name: '',
        is_tax_active: false, pajak_persen: 0, receipt_footer: '',
        
        // --- DATA BARU ---
        payment_type: 'qris_static', // qris_static atau midtrans
        midtrans_server_key: '',
        midtrans_client_key: '',
        printer_width: '58mm',
        printer_type: 'bluetooth',
        
        // Flag untuk kasih tau Golang kalau user sengaja ngehapus gambar
        delete_logo: false, 
        delete_qris: false
    });

    const logoPreview = ref(null);
    const qrisPreview = ref(null);

    const fetchStoreSettings = async () => {
        isLoading.value = true;
        try {
            const response = await api.get('/retail/store/settings');
            const data = response.data.data;
            
            // Gabungin data dari backend ke state form
            form.value = { ...form.value, ...data };
            
            // Set Default kalau dari database masih kosong (Toko Baru)
            if (!form.value.payment_type) form.value.payment_type = 'qris_static';
            if (!form.value.printer_width) form.value.printer_width = '58mm';
            if (!form.value.printer_type) form.value.printer_type = 'bluetooth';

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
            form.value.delete_logo = false; // Batal hapus karena user milih file baru
        } else if (type === 'qris') {
            form.value.qris_image = file;
            qrisPreview.value = previewUrl;
            form.value.delete_qris = false; // Batal hapus karena user milih file baru
        }
    };

    // 🚀 2. FUNGSI BARU BUAT HAPUS GAMBAR
    const removeLogo = () => {
        form.value.logo_url = null;
        logoPreview.value = null;
        form.value.delete_logo = true; // Nyalain flag biar backend tau ini minta dihapus
    };

    const removeQris = () => {
        form.value.qris_image = null;
        qrisPreview.value = null;
        form.value.delete_qris = true; // Nyalain flag biar backend tau ini minta dihapus
    };

    const saveSettings = async () => {
        isSaving.value = true;
        const formData = new FormData();
        
        Object.keys(form.value).forEach(key => {
            if (key !== 'logo_url' && key !== 'qris_image') {
                formData.append(key, form.value[key]);
            }
        });

        // Cek apakah itu beneran File baru, kalau iya baru di-append ke form data
        if (form.value.logo_url instanceof File) formData.append('logo', form.value.logo_url);
        if (form.value.qris_image instanceof File) formData.append('qris', form.value.qris_image);

        try {
            const response = await api.put('/retail/store/settings', formData);
            
            const updatedData = response.data.data;
            localStorage.setItem('storeName', updatedData.nama_toko || '');
            localStorage.setItem('storeLogo', updatedData.logo_url || '');
            
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
        handleFileChange, removeLogo, removeQris, saveSettings // 🚀 EKSPOR FUNGSI REMOVE-NYA
    };
}