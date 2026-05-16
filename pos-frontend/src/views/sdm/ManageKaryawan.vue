<script setup>
import { ref, onMounted, computed } from 'vue'; // 🚀 FIX: Tambah computed untuk Live Search
import api from '../../api.js';
import Sidebar from '../../components/Sidebar.vue';
import Swal from 'sweetalert2';

// --- STATE DATA ---
const karyawan = ref([]);
const isLoading = ref(true);
const showModal = ref(false);
const isProcessing = ref(false);
const isEditMode = ref(false);
const selectedId = ref(null);
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

// 🚀 STATE PENCARIAN KARYAWAN
const searchQuery = ref('');

const form = ref({
    name: '',
    password: '',
    tempat_lahir: '',
    tanggal_lahir: '',
    no_hp: '',
    role: 'staff',
    foto: null,           // Untuk Foto Profil
    biometric_file: null  // Untuk Foto Absensi/Face
});

const fotoPreview = ref(null);
const fileInput = ref(null);

// --- KALKULATOR MASA KERJA (TENURE) ---
const hitungMasaKerja = (tanggalDibuat) => {
    if (!tanggalDibuat) return 'Baru Bergabung';
    
    const start = new Date(tanggalDibuat);
    const end = new Date();
    
    let years = end.getFullYear() - start.getFullYear();
    let months = end.getMonth() - start.getMonth();
    
    if (months < 0) {
        years--;
        months += 12;
    }
    
    if (years === 0 && months === 0) return 'Baru Bergabung';
    
    let result = '';
    if (years > 0) result += `${years} Tahun `;
    if (months > 0) result += `${months} Bulan`;
    
    return result.trim();
};

const formatDate = (dateString) => {
    if (!dateString) return '-';
    const d = new Date(dateString);
    return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' });
};

// --- FUNGSI PREVIEW FOTO ---
const onFileChange = (e) => {
    const file = e.target.files[0];
    if (file) {
        form.value.foto = file;
        fotoPreview.value = URL.createObjectURL(file);
    }
};

// --- TARIK DATA DARI GOLANG ---
const fetchKaryawan = async () => {
    isLoading.value = true;
    try {
        const response = await api.get('/employees');
        karyawan.value = response.data.data || [];
    } catch (error) {
        console.error("Gagal menarik data:", error);
    } finally {
        isLoading.value = false;
    }
};

onMounted(() => fetchKaryawan());

// 🚀 FITUR PENCARIAN REALTIME (Computed)
const filteredKaryawan = computed(() => {
    if (!searchQuery.value) return karyawan.value;
    const query = searchQuery.value.toLowerCase();
    return karyawan.value.filter(user => 
        (user.name && user.name.toLowerCase().includes(query)) || 
        (user.nik && String(user.nik).toLowerCase().includes(query)) ||
        (user.role && user.role.toLowerCase().includes(query))
    );
});

// --- MODAL CONTROL ---
const openAddModal = () => {
    isEditMode.value = false;
    selectedId.value = null;
    form.value = { name: '', password: '', tempat_lahir: '', tanggal_lahir: '', no_hp: '', role: 'staff', foto: null };
    fotoPreview.value = null;
    showModal.value = true;
};

const openEditModal = (user) => {
    isEditMode.value = true;
    selectedId.value = user.id;
    form.value = {
        name: user.name,
        password: '',
        tempat_lahir: user.tempat_lahir || '',
        tanggal_lahir: user.tanggal_lahir ? user.tanggal_lahir.substring(0,10) : '',
        no_hp: user.no_hp || '',
        role: user.role || 'staff',
        foto: null
    };
    fotoProfilPreview.value = user.foto_url ? import.meta.env.VITE_API_BASE_URL + user.foto_url : null;
    showModal.value = true;
};

const closeModal = () => {
    showModal.value = false;
    stopCamera();
};

// --- STATE TAMBAHAN UNTUK FOTO & BIOMETRIK ---
const attendanceMethod = ref('face');
const video = ref(null);
const canvas = ref(null);
const isCameraOpen = ref(false);
const fotoProfilPreview = ref(null);
const fotoBiometricPreview = ref(null);

// --- LOGIKA KAMERA ---
const startCamera = async () => {
    isCameraOpen.value = true;
    try {
        const stream = await navigator.mediaDevices.getUserMedia({ video: true });
        video.value.srcObject = stream;
        video.value.play();
    } catch (err) {
        Swal.fire('Gagal!', 'Akses kamera ditolak atau tidak ditemukan.', 'error');
        isCameraOpen.value = false;
    }
};

const capturePhoto = () => {
    const context = canvas.value.getContext('2d');
    canvas.value.width = video.value.videoWidth;
    canvas.value.height = video.value.videoHeight;
    context.drawImage(video.value, 0, 0, canvas.value.width, canvas.value.height);
    
    canvas.value.toBlob((blob) => {
        const file = new File([blob], "face_registration.jpg", { type: "image/jpeg" });
        form.value.biometric_file = file;
        fotoBiometricPreview.value = URL.createObjectURL(file);
        stopCamera();
    }, 'image/jpeg');
};

const stopCamera = () => {
    if (video.value && video.value.srcObject) {
        const stream = video.value.srcObject;
        const tracks = stream.getTracks();
        tracks.forEach(track => track.stop());
    }
    isCameraOpen.value = false;
};

// --- FUNGSI PILIH FOTO PROFIL (UPLOAD BIASA) ---
const onProfileFileChange = (e) => {
    const file = e.target.files[0];
    if (file) {
        form.value.foto = file;
        fotoProfilPreview.value = URL.createObjectURL(file);
    }
};

// --- SUBMIT (POST & PUT) ---
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
            await api.put(`/employees/${selectedId.value}`, formData, {
                headers: { 'Content-Type': 'multipart/form-data' }
            });
            Swal.fire('Berhasil!', 'Data karyawan telah diperbarui.', 'success');
        } else {
            const response = await api.post('/employees', formData, {
                headers: { 'Content-Type': 'multipart/form-data' }
            });
            const newEmp = response.data.data;
            Swal.fire('Berhasil!', `Karyawan dengan NIK: ${newEmp.nik} berhasil dibuat.`, 'success');
        }
        
        closeModal();
        fetchKaryawan();
    } catch (error) {
        Swal.fire('Gagal!', error.response?.data?.error || 'Terjadi kesalahan sistem', 'error');
    } finally {
        isProcessing.value = false;
    }
};

// --- HAPUS KARYAWAN ---
const deleteKaryawan = async (id) => {
    const result = await Swal.fire({
        title: 'Yakin mau pecat?',
        text: "Karyawan ini tidak akan bisa login lagi!",
        icon: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#d33',
        confirmButtonText: 'Ya, Pecat!'
    });

    if (result.isConfirmed) {
        try {
            await api.delete(`/employees/${id}`);
            Swal.fire('Dihapus!', 'Karyawan telah diberhentikan.', 'success');
            fetchKaryawan();
        } catch (error) {
            Swal.fire('Gagal!', 'Gagal menghapus data.', 'error');
        }
    }
};

</script>

<template>
    <Sidebar>
        <div class="p-4 md:p-8 lg:p-10 max-w-7xl mx-auto font-sans bg-[#f8fafc] min-h-screen">
            
            <div class="bg-gradient-to-br from-slate-900 via-slate-800 to-indigo-900 rounded-[32px] p-6 md:p-10 mb-6 text-white shadow-2xl flex flex-col md:flex-row items-center justify-between relative overflow-hidden border border-slate-800 gap-6">
                <svg xmlns="http://www.w3.org/2000/svg" class="absolute -right-10 -bottom-10 w-64 h-64 text-indigo-500 opacity-10 pointer-events-none" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><line x1="19" y1="8" x2="19" y2="14"/><line x1="22" y1="11" x2="16" y2="11"/></svg>
                
                <div class="z-10 text-center md:text-left">
                    <h1 class="text-3xl md:text-4xl font-black tracking-tighter mb-1 uppercase italic">Human <span class="text-blue-400">Resources</span></h1>
                    <p class="text-slate-300 font-bold text-[10px] md:text-xs uppercase tracking-[0.2em] flex items-center justify-center md:justify-start gap-2">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-blue-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M17 21v2a4 4 0 0 1-4 4H5a4 4 0 0 1-4-4v-2a4 4 0 0 1 4-4h1v-1a4 4 0 0 1 4-4h2a4 4 0 0 1 4 4v1h1a4 4 0 0 1 4 4v2z"/><path d="M9 11v-1a2 2 0 0 1 2-2h2a2 2 0 0 1 2 2v1"/></svg>
                        Employee Directory & Access
                    </p>
                </div>
                
                <div class="z-10 w-full md:w-auto">
                    <button @click="openAddModal" class="w-full md:w-auto bg-blue-600 hover:bg-blue-500 text-white px-8 py-4 rounded-[20px] md:rounded-[24px] font-black text-xs uppercase tracking-[0.2em] shadow-xl shadow-blue-900/50 flex items-center justify-center gap-3 transition-all active:scale-95 border border-blue-400/30">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><line x1="19" y1="8" x2="19" y2="14"/><line x1="22" y1="11" x2="16" y2="11"/></svg>
                        Add Employee
                    </button>
                </div>
            </div>

            <div class="mb-8 relative group">
                <div class="absolute inset-y-0 left-0 pl-5 flex items-center pointer-events-none">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-400 group-focus-within:text-blue-600 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
                </div>
                <input v-model="searchQuery" type="text" placeholder="Cari Nama, NIK, atau Jabatan..." class="block w-full pl-14 pr-4 py-4 bg-white rounded-2xl border-2 border-slate-100 shadow-sm focus:border-blue-600 outline-none font-bold text-sm transition-all text-slate-700">
            </div>

            <div v-if="isLoading" class="py-20 flex flex-col items-center justify-center bg-white rounded-[32px] border border-slate-100 shadow-sm">
                <div class="w-12 h-12 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin mb-4"></div>
                <p class="text-slate-400 font-black text-xs uppercase tracking-widest animate-pulse">Menarik Data Karyawan...</p>
            </div>

            <div v-else-if="filteredKaryawan.length === 0" class="flex flex-col items-center justify-center py-20 bg-white/50 rounded-[32px] border-2 border-dashed border-slate-300">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-20 h-20 text-slate-300 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" /><line x1="3" y1="3" x2="21" y2="21" stroke-width="2"/></svg>
                <p class="text-slate-400 font-black text-sm uppercase tracking-widest">Karyawan Tidak Ditemukan</p>
            </div>

            <div v-else>
                <div class="lg:hidden flex flex-col gap-4">
                    <div v-for="user in filteredKaryawan" :key="user.id" class="bg-white p-5 rounded-[24px] shadow-sm border border-slate-100 flex flex-col gap-4 relative overflow-hidden">
                        
                        <div class="flex items-center gap-4">
                            <img v-if="user.foto_url" :src="API_BASE_URL + user.foto_url" class="w-14 h-14 rounded-2xl object-cover border-2 border-slate-100 shadow-sm shrink-0">
                            <div v-else class="w-14 h-14 rounded-2xl bg-slate-100 flex items-center justify-center text-slate-400 border-2 border-slate-200 shrink-0">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
                            </div>
                            
                            <div class="flex-1 min-w-0">
                                <h3 class="font-black text-base text-slate-800 uppercase truncate">{{ user.name }}</h3>
                                <div class="flex items-center gap-1.5 mt-1">
                                    <span v-if="user.role === 'owner'" class="inline-flex items-center gap-1 bg-slate-900 text-white font-black px-2 py-0.5 rounded text-[8px] uppercase tracking-widest">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-2.5 h-2.5 text-amber-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/></svg>
                                        Owner
                                    </span>
                                    <span v-else-if="user.role === 'manager'" class="inline-flex items-center gap-1 bg-purple-50 text-purple-600 border border-purple-100 font-black px-2 py-0.5 rounded text-[8px] uppercase tracking-widest">
                                        Manager
                                    </span>
                                    <span v-else-if="user.role === 'supervisor'" class="inline-flex items-center gap-1 bg-blue-50 text-blue-600 border border-blue-100 font-black px-2 py-0.5 rounded text-[8px] uppercase tracking-widest">
                                        Supervisor
                                    </span>
                                    <span v-else class="inline-flex items-center gap-1 bg-emerald-50 text-emerald-600 border border-emerald-100 font-black px-2 py-0.5 rounded text-[8px] uppercase tracking-widest">
                                        Staff
                                    </span>
                                    <span class="text-[9px] font-black text-blue-600 tracking-widest">{{ user.nik }}</span>
                                </div>
                            </div>
                        </div>

                        <div class="grid grid-cols-2 gap-3 pt-3 border-t border-dashed border-slate-100">
                            <div>
                                <p class="text-[8px] font-black text-slate-400 uppercase tracking-widest mb-1">Tgl Bergabung</p>
                                <p class="font-bold text-xs text-slate-700">{{ formatDate(user.created_at) }}</p>
                            </div>
                            <div>
                                <p class="text-[8px] font-black text-slate-400 uppercase tracking-widest mb-1">Masa Kerja</p>
                                <p class="font-black text-xs text-indigo-600">{{ hitungMasaKerja(user.created_at) }}</p>
                            </div>
                            <div class="col-span-2">
                                <p class="text-[8px] font-black text-slate-400 uppercase tracking-widest mb-1">Kontak / Ttl</p>
                                <p class="font-bold text-xs text-slate-600">{{ user.no_hp || '-' }} • {{ user.tempat_lahir || '-' }}, {{ user.tanggal_lahir || '-' }}</p>
                            </div>
                        </div>

                        <div class="flex gap-2 mt-2" v-if="user.role !== 'owner'">
                            <button @click="openEditModal(user)" class="flex-1 bg-slate-100 text-slate-500 hover:bg-blue-600 hover:text-white py-3 rounded-xl transition-colors font-black text-[10px] uppercase tracking-widest flex items-center justify-center gap-2">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M17 3a2.828 2.828 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5L17 3z"/></svg>
                                Edit
                            </button>
                            <button @click="deleteKaryawan(user.id)" class="flex-1 bg-red-50 text-red-500 hover:bg-red-500 hover:text-white py-3 rounded-xl transition-colors font-black text-[10px] uppercase tracking-widest flex items-center justify-center gap-2">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/></svg>
                                Pecat
                            </button>
                        </div>
                    </div>
                </div>

                <div class="hidden lg:block bg-white rounded-[32px] shadow-sm border border-slate-100 overflow-hidden">
                    <div class="overflow-x-auto custom-scrollbar">
                        <table class="w-full text-left whitespace-nowrap">
                            <thead class="bg-slate-50/80 border-b border-slate-100">
                                <tr>
                                    <th class="px-6 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">Employee Profile</th>
                                    <th class="px-6 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">Login ID (NIK)</th>
                                    <th class="px-6 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">Tenure & Date Joined</th>
                                    <th class="px-6 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">Contact Info</th>
                                    <th class="px-6 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest text-right">Actions</th>
                                </tr>
                            </thead>
                            <tbody class="divide-y divide-slate-50">
                                <tr v-for="user in filteredKaryawan" :key="user.id" class="hover:bg-slate-50/50 transition-colors group">
                                    <td class="px-6 py-5">
                                        <div class="flex items-center gap-4">
                                            <img 
                                                v-if="user.foto_url" 
                                                :src="API_BASE_URL + user.foto_url"
                                                class="w-14 h-14 rounded-[14px] object-cover border-2 border-slate-100 shadow-sm"
                                            >
                                            <div v-else class="w-14 h-14 rounded-[14px] bg-slate-100 flex items-center justify-center text-slate-400 border-2 border-slate-200">
                                                <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
                                            </div>
                                            <div>
                                                <div class="text-slate-800 font-black text-sm uppercase tracking-tight">{{ user.name }}</div>
                                                <div class="mt-1.5">
                                                    <span v-if="user.role === 'owner'" class="inline-flex items-center gap-1 bg-slate-900 text-white font-black px-2.5 py-0.5 rounded text-[9px] uppercase tracking-widest shadow-sm">
                                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-2.5 h-2.5 text-amber-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/></svg>
                                                        Owner
                                                    </span>
                                                    <span v-else-if="user.role === 'manager'" class="inline-flex items-center gap-1 bg-purple-50 text-purple-600 border border-purple-100 font-black px-2.5 py-0.5 rounded text-[9px] uppercase tracking-widest">
                                                        Manager
                                                    </span>
                                                    <span v-else-if="user.role === 'supervisor'" class="inline-flex items-center gap-1 bg-blue-50 text-blue-600 border border-blue-100 font-black px-2.5 py-0.5 rounded text-[9px] uppercase tracking-widest">
                                                        Supervisor
                                                    </span>
                                                    <span v-else class="inline-flex items-center gap-1 bg-emerald-50 text-emerald-600 border border-emerald-100 font-black px-2.5 py-0.5 rounded text-[9px] uppercase tracking-widest">
                                                        Staff
                                                    </span>
                                                </div>
                                            </div>
                                        </div>
                                    </td>
                                    
                                    <td class="px-6 py-5">
                                        <div class="inline-flex items-center gap-2 bg-slate-50 px-3 py-2 rounded-xl border border-slate-200">
                                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><rect width="18" height="12" x="3" y="4" rx="2" ry="2"/><path d="M2 20h20"/><path d="M12 16v4"/></svg>
                                            <span class="text-blue-600 font-black tracking-widest text-xs">{{ user.nik }}</span>
                                        </div>
                                    </td>

                                    <td class="px-6 py-5">
                                        <div class="font-black text-indigo-600 text-sm">{{ hitungMasaKerja(user.created_at) }}</div>
                                        <div class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-1">Sejak: {{ formatDate(user.created_at) }}</div>
                                    </td>

                                    <td class="px-6 py-5">
                                        <div class="text-slate-600 font-bold text-xs">{{ user.no_hp || '-' }}</div>
                                        <div class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-1">{{ user.tempat_lahir || '-' }}, {{ user.tanggal_lahir || '-' }}</div>
                                    </td>
                                    
                                    <td class="px-6 py-5 text-right">
                                        <div class="flex items-center justify-end gap-2">
                                            <template v-if="user.role !== 'owner'">
                                                <button @click="openEditModal(user)" class="bg-slate-50 border border-slate-200 text-slate-500 hover:bg-blue-600 hover:border-blue-600 hover:text-white p-2.5 rounded-xl transition-colors shadow-sm" title="Edit Data">
                                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M17 3a2.828 2.828 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5L17 3z"/></svg>
                                                </button>
                                                <button @click="deleteKaryawan(user.id)" class="bg-red-50 text-red-500 hover:bg-red-500 hover:text-white p-2.5 rounded-xl transition-colors shadow-sm" title="Pecat Karyawan">
                                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/></svg>
                                                </button>
                                            </template>
                                            <span v-else class="text-[9px] font-black text-slate-300 uppercase tracking-widest bg-slate-50 border border-slate-100 px-3 py-2 rounded-xl shadow-sm">Protected</span>
                                        </div>
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>

        <div v-if="showModal" class="fixed inset-0 bg-slate-900/80 backdrop-blur-sm z-[100] flex items-center justify-center p-4">
            <div class="bg-white rounded-[32px] md:rounded-[40px] shadow-2xl w-full max-w-4xl flex flex-col max-h-[95vh] overflow-hidden border-[8px] md:border-[10px] border-slate-900/5">
                
                <div class="p-6 md:p-8 border-b border-slate-50 flex justify-between items-center bg-slate-50/80 shrink-0">
                    <h3 class="font-black text-xl md:text-2xl text-slate-800 tracking-tighter uppercase italic">{{ isEditMode ? 'Update Profile' : 'New Registration' }}</h3>
                    <button @click="closeModal" class="text-slate-400 hover:text-red-500 transition-colors bg-white p-2 md:p-2.5 rounded-xl shadow-sm border border-slate-200">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 md:w-6 md:h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                    </button>
                </div>
                
                <div class="overflow-y-auto custom-scrollbar p-6 md:p-8">
                    <form @submit.prevent="submit" class="flex flex-col gap-6 md:gap-8">
                        
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-5 md:gap-6">
                            <div>
                                <label class="block text-[9px] md:text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1.5 md:mb-2">Nama Lengkap</label>
                                <input v-model="form.name" type="text" required class="w-full p-3.5 md:p-4 bg-slate-50 rounded-2xl outline-none border-2 border-slate-100 focus:border-blue-600 focus:bg-white font-bold text-sm transition-all text-slate-800">
                            </div>
                            
                            <div>
                                <label class="block text-[9px] md:text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1.5 md:mb-2">Jabatan / Role</label>
                                <select v-model="form.role" required class="w-full p-3.5 md:p-4 bg-slate-50 rounded-2xl outline-none border-2 border-slate-100 focus:border-blue-600 focus:bg-white font-black text-xs text-slate-700 cursor-pointer transition-all uppercase tracking-widest appearance-none">
                                    <option value="manager">Manager</option>
                                    <option value="supervisor">Supervisor</option>
                                    <option value="staff">Staff / Cashier</option>
                                </select>
                            </div>

                            <div class="grid grid-cols-2 gap-3 md:gap-4">
                                <div>
                                    <label class="block text-[9px] md:text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1.5 md:mb-2">Tempat Lahir</label>
                                    <input v-model="form.tempat_lahir" type="text" required class="w-full p-3.5 md:p-4 bg-slate-50 rounded-2xl outline-none border-2 border-slate-100 focus:border-blue-600 focus:bg-white font-bold text-sm transition-all text-slate-800">
                                </div>
                                <div>
                                    <label class="block text-[9px] md:text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1.5 md:mb-2">Tgl Lahir</label>
                                    <input v-model="form.tanggal_lahir" type="date" required class="w-full p-3.5 md:p-4 bg-slate-50 rounded-2xl outline-none border-2 border-slate-100 focus:border-blue-600 focus:bg-white font-bold text-xs uppercase tracking-widest transition-all text-slate-800">
                                </div>
                            </div>
                            
                            <div>
                                <label class="block text-[9px] md:text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1.5 md:mb-2">No. Handphone (WA)</label>
                                <input v-model="form.no_hp" type="text" placeholder="08..." required class="w-full p-3.5 md:p-4 bg-slate-50 rounded-2xl outline-none border-2 border-slate-100 focus:border-blue-600 focus:bg-white font-bold text-sm transition-all text-slate-800">
                            </div>
                            
                            <div class="md:col-span-2">
                                <label class="block text-[9px] md:text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1.5 md:mb-2">Password <span v-if="isEditMode" class="text-amber-500 italic lowercase tracking-normal">(Kosongi jika tidak diubah)</span></label>
                                <input v-model="form.password" type="password" :required="!isEditMode" class="w-full p-3.5 md:p-4 bg-slate-50 rounded-2xl outline-none border-2 border-slate-100 focus:border-blue-600 focus:bg-white font-bold text-sm transition-all placeholder:text-slate-300 text-slate-800" placeholder="••••••••">
                            </div>
                        </div>

                        <div class="grid grid-cols-1 md:grid-cols-2 gap-5 md:gap-6 pt-6 md:pt-8 border-t border-slate-100">
                            
                            <div class="flex flex-col items-center p-6 bg-slate-50 rounded-[32px] border-2 border-slate-100">
                                <label class="text-[10px] font-black text-slate-500 uppercase tracking-widest mb-4 flex items-center gap-2">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-blue-500" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
                                    1. Foto Profil (Aesthetic)
                                </label>
                                <div @click="$refs.profileInput.click()" class="w-36 h-36 md:w-40 md:h-40 bg-white rounded-[24px] shadow-sm flex items-center justify-center cursor-pointer overflow-hidden border-2 border-slate-200 group relative transition-all hover:border-blue-500 hover:shadow-blue-100 hover:shadow-xl">
                                    <img v-if="fotoProfilPreview" :src="fotoProfilPreview" class="w-full h-full object-cover">
                                    <div v-else class="text-center p-4">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-10 h-10 mx-auto text-slate-300 group-hover:text-blue-500 transition-colors" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect width="18" height="18" x="3" y="3" rx="2" ry="2"/><circle cx="9" cy="9" r="2"/><path d="m21 15-3.086-3.086a2 2 0 0 0-2.828 0L6 21"/></svg>
                                    </div>
                                    <div class="absolute inset-0 bg-slate-900/60 opacity-0 group-hover:opacity-100 flex items-center justify-center transition-opacity backdrop-blur-sm">
                                        <span class="text-white text-[10px] font-black uppercase tracking-widest">Pilih File</span>
                                    </div>
                                </div>
                                <input type="file" ref="profileInput" @change="onProfileFileChange" class="hidden" accept="image/*">
                            </div>

                            <div class="flex flex-col items-center p-6 bg-slate-900 rounded-[32px] relative overflow-hidden shadow-2xl border-4 border-slate-800">
                                <label class="text-[10px] font-black text-indigo-400 uppercase tracking-widest mb-4 flex items-center gap-2 relative z-10">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10"/></svg>
                                    2. Foto Biometrik (Face AI)
                                </label>

                                <div class="w-full flex flex-col items-center gap-4 relative z-10">
                                    <div class="w-36 h-36 md:w-40 md:h-40 bg-slate-950 rounded-[24px] overflow-hidden relative border-2 border-slate-700 shadow-inner">
                                        <video ref="video" v-show="isCameraOpen" autoplay muted playsinline class="w-full h-full object-cover scale-x-[-1]"></video>
                                        <img v-if="fotoBiometricPreview && !isCameraOpen" :src="fotoBiometricPreview" class="w-full h-full object-cover">
                                        <canvas ref="canvas" class="hidden"></canvas>
                                        
                                        <div v-if="isCameraOpen" class="absolute inset-0 border-2 border-indigo-500/50 rounded-[24px] pointer-events-none">
                                            <div class="w-full h-0.5 bg-indigo-500 absolute top-0 animate-[scan_2s_infinite] shadow-[0_0_8px_#6366f1]"></div>
                                        </div>
                                    </div>

                                    <div class="flex gap-2 w-full justify-center">
                                        <button v-if="!isCameraOpen" @click.prevent="startCamera" class="bg-indigo-600 hover:bg-indigo-500 text-white px-6 py-3 md:py-3.5 rounded-xl font-black text-[9px] md:text-[10px] uppercase tracking-widest flex items-center justify-center gap-2 transition-colors w-3/4">
                                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M23 19a2 2 0 0 1-2-2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z"/><circle cx="12" cy="13" r="4"/></svg>
                                            Nyalakan Kamera
                                        </button>
                                        <button v-else @click.prevent="capturePhoto" class="bg-emerald-500 hover:bg-emerald-400 text-white px-6 py-3 md:py-3.5 rounded-xl font-black text-[9px] md:text-[10px] uppercase tracking-widest animate-pulse transition-colors w-3/4 shadow-lg shadow-emerald-500/30">
                                            Jepret Wajah!
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <div class="pt-6 border-t border-slate-100 mt-2">
                            <button type="submit" :disabled="isProcessing" class="w-full bg-blue-600 hover:bg-slate-900 text-white py-4 md:py-5 rounded-[20px] md:rounded-[24px] font-black text-xs md:text-sm shadow-xl shadow-blue-200 transition-all active:scale-95 disabled:opacity-50 uppercase tracking-[0.2em] flex items-center justify-center gap-3">
                                <template v-if="isProcessing">
                                    <div class="w-5 h-5 border-4 border-white/30 border-t-white rounded-full animate-spin"></div>
                                    MENYIMPAN DATA...
                                </template>
                                <template v-else-if="isEditMode">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/><polyline points="17 21 17 13 7 13 7 21"/><polyline points="7 3 7 8 15 8"/></svg>
                                    SIMPAN PERUBAHAN
                                </template>
                                <template v-else>
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><line x1="19" y1="8" x2="19" y2="14"/><line x1="22" y1="11" x2="16" y2="11"/></svg>
                                    DAFTARKAN KARYAWAN
                                </template>
                            </button>
                        </div>
                        
                    </form>
                </div>
            </div>
        </div>

    </Sidebar>
</template>

<style scoped>
/* Custom Scrollbar */
.custom-scrollbar::-webkit-scrollbar { height: 6px; width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #94a3b8; }

/* Menyembunyikan panah default select bawaan browser */
select {
    -webkit-appearance: none;
    -moz-appearance: none;
    appearance: none;
}

/* Animasi untuk garis scanner biometrik */
@keyframes scan {
  0% { top: 0%; opacity: 0; }
  10% { opacity: 1; }
  90% { opacity: 1; }
  100% { top: 100%; opacity: 0; }
}
</style>