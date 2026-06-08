import { ref, onMounted } from 'vue';
import api from '../../../api.js';
import Swal from 'sweetalert2';
import imageCompression from 'browser-image-compression';

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

    // HELPER FUNGSI: Biar bisa ngebaca link Cloud Supabase (https://...) ATAU sisa data Lokal
    const getCleanUrl = (url) => {
        if (!url) return null;
        if (url.startsWith('http://') || url.startsWith('https://')) {
            return url; // Sudah dari Cloud, balikin mentah!
        }
        return API_BASE_URL + url;
    };

    const fetchProfile = async () => {
        isLoading.value = true;
        try {
            const response = await api.get('/me');
            const data = response.data;

            role.value = data.role;
            profileForm.value.name = data.name || '';
            profileForm.value.no_hp = data.no_hp || '';
            profileForm.value.tempat_lahir = data.tempat_lahir || '';
            profileForm.value.tanggal_lahir = data.tanggal_lahir ? data.tanggal_lahir.substring(0, 10) : '';

            fotoPreview.value = getCleanUrl(data.foto_url);
            bioPreview.value = getCleanUrl(data.biometric_url);
            
        } catch (error) {
            Swal.fire('Gagal', 'Tidak dapat mengambil data profil', 'error');
        } finally {
            isLoading.value = false;
        }
    };

    // 🚀 FIX: handleFileChange sekarang mendukung kompresi background sebelum disimpan di state form
    const handleFileChange = async (type, file) => {
        if (!file) return;

        // Jika tipe file bukan gambar (misal file biner lain), lewati kompresi
        if (!file.type.startsWith('image/')) {
            if (type === 'foto') {
                profileForm.value.foto = file;
                fotoPreview.value = URL.createObjectURL(file);
            } else if (type === 'bio') {
                profileForm.value.biometric_file = file;
                bioPreview.value = URL.createObjectURL(file);
            }
            return;
        }

        // Tentukan aturan kompresi dinamis berdasarkan jenis dokumen
        const options = {
            maxSizeMB: type === 'foto' ? 0.3 : 1.0,       // Avatar cukup 300KB, Biometrik/Nota max 1MB
            maxWidthOrHeight: type === 'foto' ? 1024 : 1920, // Batasi resolusi piksel maksimal
            useWebWorker: true,
            fileType: file.type
        };

        try {
            // Tampilkan loading screen mini agar kasir tahu gambar sedang diproses browser
            const compressedFile = await imageCompression(file, options);

            // Buat preview URL dari hasil file yang sudah dikompresi
            const previewUrl = URL.createObjectURL(compressedFile);

            if (type === 'foto') {
                profileForm.value.foto = compressedFile;
                fotoPreview.value = previewUrl;
            } else if (type === 'bio') {
                profileForm.value.biometric_file = compressedFile;
                bioPreview.value = previewUrl;
            }
        } catch (error) {
            console.error("Gagal melakukan optimasi gambar:", error);
            // Fallback: Jika kompresi gagal total, gunakan file asli agar aplikasi tidak macet
            const fallbackUrl = URL.createObjectURL(file);
            if (type === 'foto') {
                profileForm.value.foto = file;
                fotoPreview.value = fallbackUrl;
            } else if (type === 'bio') {
                profileForm.value.biometric_file = file;
                bioPreview.value = fallbackUrl;
            }
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
            
            // 🚀 SINKRONISASI BACKEND: Tambahkan pengiriman file biometrik yang kemarin tertinggal di FormData!
            if (profileForm.value.biometric_file) {
                formData.append('biometric_file', profileForm.value.biometric_file);
            }

            const res = await api.put('/profile', formData, {
                headers: { 'Content-Type': 'multipart/form-data' }
            });

            const responseData = res.data.data || res.data;

            if (responseData.name) {
                localStorage.setItem('name', responseData.name);
            }
            if (responseData.foto_url) {
                localStorage.setItem('foto_url', responseData.foto_url);
            }

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