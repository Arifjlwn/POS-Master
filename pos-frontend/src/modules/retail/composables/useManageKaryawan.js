import { ref, computed, onMounted } from 'vue';
import api from '../../../api.js';
import Swal from 'sweetalert2';

export function useManageKaryawan() {
    const karyawan = ref([]);
    const isLoading = ref(true);
    const showModal = ref(false);
    const isProcessing = ref(false);
    const isEditMode = ref(false);
    const selectedId = ref(null);
    const searchQuery = ref('');
    const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

    const form = ref({
        name: '', password: '', tempat_lahir: '', tanggal_lahir: '',
        no_hp: '', role: 'staff', foto: null, biometric_file: null
    });

    const isCameraOpen = ref(false);
    const fotoProfilPreview = ref(null);
    const fotoBiometricPreview = ref(null);

    const hitungMasaKerja = (tanggalDibuat) => {
        if (!tanggalDibuat) return 'Baru Bergabung';
        const start = new Date(tanggalDibuat);
        const end = new Date();
        let years = end.getFullYear() - start.getFullYear();
        let months = end.getMonth() - start.getMonth();
        if (months < 0) { years--; months += 12; }
        if (years === 0 && months === 0) return 'Baru Bergabung';
        
        let result = '';
        if (years > 0) result += `${years} Tahun `;
        if (months > 0) result += `${months} Bulan`;
        return result.trim();
    };

    const formatDate = (dateString) => {
        if (!dateString) return '-';
        return new Date(dateString).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' });
    };

    const fetchKaryawan = async () => {
        isLoading.value = true;
        try {
            const response = await api.get('/retail/employees');
            karyawan.value = response.data.data || [];
        } catch (error) {
            console.error("Gagal menarik data:", error);
        } finally {
            isLoading.value = false;
        }
    };

    const filteredKaryawan = computed(() => {
        if (!searchQuery.value) return karyawan.value;
        const query = searchQuery.value.toLowerCase();
        return karyawan.value.filter(user => 
            (user.name && user.name.toLowerCase().includes(query)) || 
            (user.nik && String(user.nik).toLowerCase().includes(query)) ||
            (user.role && user.role.toLowerCase().includes(query))
        );
    });

    const openAddModal = () => {
        isEditMode.value = false;
        selectedId.value = null;
        form.value = { name: '', password: '', tempat_lahir: '', tanggal_lahir: '', no_hp: '', role: 'staff', foto: null, biometric_file: null };
        fotoProfilPreview.value = null;
        fotoBiometricPreview.value = null;
        showModal.value = true;
    };

    const openEditModal = (user) => {
        isEditMode.value = true;
        selectedId.value = user.id;
        form.value = {
            name: user.name, password: '', tempat_lahir: user.tempat_lahir || '',
            tanggal_lahir: user.tanggal_lahir ? user.tanggal_lahir.substring(0,10) : '',
            no_hp: user.no_hp || '', role: user.role || 'staff', foto: null, biometric_file: null
        };
        fotoProfilPreview.value = user.foto_url ? API_BASE_URL + user.foto_url : null;
        fotoBiometricPreview.value = user.biometric_url ? API_BASE_URL + user.biometric_url : null;
        showModal.value = true;
    };

    const closeModal = () => {
        showModal.value = false;
    };

    const handleUpdateProfile = (file, previewUrl) => {
        form.value.foto = file;
        fotoProfilPreview.value = previewUrl;
    };

    const handleUpdateBiometric = (file, previewUrl) => {
        form.value.biometric_file = file;
        fotoBiometricPreview.value = previewUrl;
    };

    const submit = async () => {
        isProcessing.value = true;
        const formData = new FormData();
        formData.append('name', form.value.name);
        formData.append('tempat_lahir', form.value.tempat_lahir);
        formData.append('tanggal_lahir', form.value.tanggal_lahir);
        formData.append('no_hp', form.value.no_hp);
        formData.append('role', form.value.role); 
        
        if (form.value.password) formData.append('password', form.value.password);
        if (form.value.foto) formData.append('foto', form.value.foto);
        if (form.value.biometric_file) formData.append('biometric_file', form.value.biometric_file);

        try {
            if (isEditMode.value) {
                // 🚀 HAPUS HEADERS MANUAL! Biar Axios otomatis ngasih boundary
                await api.put(`/retail/employees/${selectedId.value}`, formData);
                Swal.fire('Berhasil!', 'Data karyawan telah diperbarui.', 'success');
            } else {
                // 🚀 HAPUS JUGA DI SINI!
                const response = await api.post('/retail/employees', formData);
                Swal.fire('Berhasil!', `Karyawan dengan NIK: ${response.data.data.nik} berhasil dibuat.`, 'success');
            }
            closeModal();
            fetchKaryawan();
        } catch (error) {
            Swal.fire('Gagal!', error.response?.data?.error || 'Terjadi kesalahan sistem', 'error');
        } finally {
            isProcessing.value = false;
        }
    };

    const deleteKaryawan = async (id) => {
        const result = await Swal.fire({
            title: 'Yakin mau pecat?', text: "Karyawan ini tidak akan bisa login lagi!",
            icon: 'warning', showCancelButton: true, confirmButtonColor: '#d33', confirmButtonText: 'Ya, Pecat!'
        });
        if (result.isConfirmed) {
            try {
                await api.delete(`/retail/employees/${id}`);
                Swal.fire('Dihapus!', 'Karyawan telah diberhentikan.', 'success');
                fetchKaryawan();
            } catch (error) {
                Swal.fire('Gagal!', 'Gagal menghapus data.', 'error');
            }
        }
    };

    onMounted(() => fetchKaryawan());

    return {
        karyawan, isLoading, showModal, isProcessing, isEditMode, searchQuery,
        form, isCameraOpen, fotoProfilPreview, fotoBiometricPreview, API_BASE_URL,
        filteredKaryawan, hitungMasaKerja, formatDate, openAddModal, openEditModal,
        closeModal, handleUpdateProfile, handleUpdateBiometric, submit, deleteKaryawan
    };
}