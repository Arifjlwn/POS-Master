import { ref, onMounted } from 'vue';
import api from '../../../api.js';
import Swal from 'sweetalert2';

export function useAccount() {
    const isLoading = ref(true);
    const isSaving = ref(false);
    const activeTab = ref('profile'); // profile, security, billing
    const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';
    const role = ref('');

    const profileForm = ref({
        name: '', no_hp: '', tempat_lahir: '', tanggal_lahir: '',
        foto: null, biometric_file: null
    });

    const passwordForm = ref({
        old_password: '', new_password: '', confirm_password: ''
    });

    const fotoPreview = ref(null);
    const bioPreview = ref(null);

    const fetchProfile = async () => {
        isLoading.value = true;
        try {
            // Narik data dari fungsi GetMe di Golang lu
            const response = await api.get('/me');
            const data = response.data;

            role.value = data.role;
            profileForm.value.name = data.name || '';
            profileForm.value.no_hp = data.no_hp || '';
            profileForm.value.tempat_lahir = data.tempat_lahir || '';
            profileForm.value.tanggal_lahir = data.tanggal_lahir ? data.tanggal_lahir.substring(0, 10) : '';

            if (data.foto_url) fotoPreview.value = API_BASE_URL + data.foto_url;
            if (data.biometric_url) bioPreview.value = API_BASE_URL + data.biometric_url;
        } catch (error) {
            Swal.fire('Gagal', 'Tidak dapat mengambil data profil', 'error');
        } finally {
            isLoading.value = false;
        }
    };

    const handleFileChange = (type, file, previewUrl) => {
        if (type === 'foto') {
            profileForm.value.foto = file;
            fotoPreview.value = previewUrl;
        } else if (type === 'bio') {
            profileForm.value.biometric_file = file;
            bioPreview.value = previewUrl;
        }
    };

    const saveProfile = async () => {
        isSaving.value = true;
        try {
            const formData = new FormData();
            formData.append('name', profileForm.value.name);
            formData.append('no_hp', profileForm.value.no_hp);
            formData.append('tempat_lahir', profileForm.value.tempat_lahir);
            formData.append('tanggal_lahir', profileForm.value.tanggal_lahir);

            if (profileForm.value.foto) {
                formData.append('foto', profileForm.value.foto);
            }

            const res = await api.put('/profile', formData, {
                headers: { 'Content-Type': 'multipart/form-data' }
            });

            // 🚀 ANTISIPASI FORMAT RESPONSE GOLANG (Bisa res.data.data atau res.data langsung)
            const responseData = res.data.data || res.data;

            // 1. Simpan nama baru
            if (responseData.name) {
                localStorage.setItem('name', responseData.name);
            }

            // 2. Simpan URL foto baru (Ini yang bikin refresh balik lagi kalau gagal)
            if (responseData.foto_url) {
                localStorage.setItem('foto_url', responseData.foto_url);
            }

            // 3. TEMBAK SINYAL CUSTOM KE SIDEBAR!
            window.dispatchEvent(new Event('profile-updated'));

            Swal.fire({
                icon: 'success',
                title: 'Profil Diperbarui!',
                text: 'Perubahan telah disimpan ke dalam sistem.',
                confirmButtonColor: '#4f46e5',
                customClass: { popup: 'rounded-[32px]' }
            });

        } catch (error) {
            Swal.fire({
                icon: 'error',
                title: 'Gagal Menyimpan',
                text: error.response?.data?.error || 'Terjadi kesalahan sistem',
                confirmButtonColor: '#ef4444',
                customClass: { popup: 'rounded-[32px]' }
            });
        } finally {
            isSaving.value = false;
        }
    };

    const updatePassword = async () => {
        if (passwordForm.value.new_password !== passwordForm.value.confirm_password) {
            return Swal.fire('Error', 'Konfirmasi password baru tidak cocok!', 'error');
        }

        isSaving.value = true;
        try {
            // Asumsi rute ini bakal kita buat di Golang nanti
            await api.put('/password', {
                old_password: passwordForm.value.old_password,
                new_password: passwordForm.value.new_password
            });
            Swal.fire('Berhasil!', 'Password berhasil diubah. Harap ingat password baru Anda.', 'success');
            passwordForm.value = { old_password: '', new_password: '', confirm_password: '' };
        } catch (error) {
            Swal.fire('Gagal', error.response?.data?.error || 'Password lama salah.', 'error');
        } finally {
            isSaving.value = false;
        }
    };

    onMounted(fetchProfile);

    return {
        isLoading, isSaving, activeTab, role, profileForm, passwordForm,
        fotoPreview, bioPreview, handleFileChange, saveProfile, updatePassword
    };
}
