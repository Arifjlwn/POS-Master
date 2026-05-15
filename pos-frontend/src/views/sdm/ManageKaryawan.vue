<script setup>
import { ref, onMounted } from 'vue';
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
        tanggal_lahir: user.tanggal_lahir || '',
        no_hp: user.no_hp || '',
        role: user.role || 'staff',
        foto: null
    };
    fotoPreview.value = user.foto_url ? import.meta.env.VITE_API_BASE_URL + user.foto_url : null;
    showModal.value = true;
};

const closeModal = () => {
    showModal.value = false;
};

// --- STATE TAMBAHAN UNTUK FOTO & BIOMETRIK ---
const attendanceMethod = ref('face'); // Opsi: 'face', 'finger', 'other'
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
    const stream = video.value.srcObject;
    const tracks = stream.getTracks();
    tracks.forEach(track => track.stop());
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
    formData.append('role', form.value.role); // 🚀 KIRIM JABATAN KE GOLANG
    
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
            Swal.fire('Berhasil!', `Karyawan dengan NIK: ${newEmp.nik} & Jabatan: ${newEmp.jabatan.toUpperCase()} berhasil dibuat.`, 'success');
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
        <div class="p-4 md:p-8 lg:p-10 max-w-7xl mx-auto font-sans">
            <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
                <div>
                    <h1 class="text-3xl font-black text-slate-900 tracking-tighter uppercase italic">Human <span class="text-blue-600">Resources</span></h1>
                    <p class="text-xs text-slate-400 font-black uppercase tracking-[0.2em] mt-1">Employee Directory & Access</p>
                </div>
                <button @click="openAddModal" class="bg-slate-900 hover:bg-blue-600 text-white px-6 md:px-8 py-3.5 md:py-4 rounded-[20px] font-black text-xs uppercase tracking-widest flex items-center justify-center gap-3 shadow-xl transition-all active:scale-95">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><line x1="19" y1="8" x2="19" y2="14"/><line x1="22" y1="11" x2="16" y2="11"/></svg>
                    Add Employee
                </button>
            </div>

            <div class="bg-white rounded-[32px] shadow-sm border border-slate-100 overflow-hidden">
                <div class="overflow-x-auto custom-scrollbar">
                    <table class="w-full text-left whitespace-nowrap">
                        <thead class="bg-slate-50/50 border-b border-slate-100">
                            <tr>
                                <th class="px-6 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">Employee Profile</th>
                                <th class="px-6 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">Login ID (NIK)</th>
                                <th class="px-6 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">Role</th>
                                <th class="px-6 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">Contact</th>
                                <th class="px-6 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest text-right">Actions</th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-slate-50">
                            <tr v-if="isLoading">
                                <td colspan="5" class="px-6 py-12 text-center text-slate-400 font-black text-xs uppercase tracking-widest">Loading Directory...</td>
                            </tr>
                            <tr v-else-if="karyawan.length === 0">
                                <td colspan="5" class="px-6 py-12 text-center text-slate-400 font-black text-xs uppercase tracking-widest opacity-50">No Employees Found</td>
                            </tr>
                            <tr v-for="user in karyawan" :key="user.id" class="hover:bg-slate-50/80 transition-colors group">
                                <td class="px-6 py-4">
                                    <div class="flex items-center gap-4">
                                        <img 
                                            v-if="user.foto_url" 
                                            :src="API_BASE_URL + user.foto_url"
                                            class="w-12 h-12 rounded-xl object-cover border border-slate-200 shadow-sm"
                                        >
                                        <div v-else class="w-12 h-12 rounded-xl bg-slate-100 flex items-center justify-center text-slate-400 border border-slate-200">
                                            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
                                        </div>
                                        <div>
                                            <div class="text-slate-800 font-black text-sm uppercase">{{ user.name }}</div>
                                            <div class="text-[10px] text-slate-400 font-bold uppercase tracking-widest mt-0.5">{{ user.tempat_lahir || 'N/A' }}, {{ user.tanggal_lahir || 'N/A' }}</div>
                                        </div>
                                    </div>
                                </td>
                                <td class="px-6 py-4">
                                    <div class="inline-flex items-center gap-2 bg-slate-100 px-3 py-1.5 rounded-lg border border-slate-200">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 text-slate-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><rect width="18" height="12" x="3" y="4" rx="2" ry="2"/><path d="M2 20h20"/><path d="M12 16v4"/></svg>
                                        <span class="text-blue-600 font-black tracking-widest text-xs">{{ user.nik }}</span>
                                    </div>
                                </td>
                                <td class="px-6 py-4">
                                    <span v-if="user.role === 'owner'" class="inline-flex items-center gap-1.5 bg-slate-900 text-white font-black px-3 py-1.5 rounded-lg text-[9px] uppercase tracking-widest shadow-sm">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3 text-amber-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/></svg>
                                        Owner
                                    </span>
                                    <span v-else-if="user.role === 'manager'" class="inline-flex items-center gap-1.5 bg-purple-50 text-purple-600 border border-purple-100 font-black px-3 py-1.5 rounded-lg text-[9px] uppercase tracking-widest">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><rect width="20" height="14" x="2" y="7" rx="2" ry="2"/><path d="M16 21V5a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16"/></svg>
                                        Manager
                                    </span>
                                    <span v-else-if="user.role === 'supervisor'" class="inline-flex items-center gap-1.5 bg-blue-50 text-blue-600 border border-blue-100 font-black px-3 py-1.5 rounded-lg text-[9px] uppercase tracking-widest">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/></svg>
                                        Supervisor
                                    </span>
                                    <span v-else class="inline-flex items-center gap-1.5 bg-emerald-50 text-emerald-600 border border-emerald-100 font-black px-3 py-1.5 rounded-lg text-[9px] uppercase tracking-widest">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M22 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
                                        Staff
                                    </span>
                                </td>
                                <td class="px-6 py-4 text-slate-500 font-bold text-xs">{{ user.no_hp || '-' }}</td>
                                <td class="px-6 py-4 text-right">
                                    <div class="flex items-center justify-end gap-2">
                                        <template v-if="user.role !== 'owner'">
                                            <button @click="openEditModal(user)" class="bg-slate-100 text-slate-500 hover:bg-blue-600 hover:text-white p-2 rounded-lg transition-colors" title="Edit Data">
                                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M17 3a2.828 2.828 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5L17 3z"/></svg>
                                            </button>
                                            <button @click="deleteKaryawan(user.id)" class="bg-red-50 text-red-500 hover:bg-red-500 hover:text-white p-2 rounded-lg transition-colors" title="Pecat Karyawan">
                                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/></svg>
                                            </button>
                                        </template>
                                        <span v-else class="text-[9px] font-black text-slate-300 uppercase tracking-widest bg-slate-50 px-3 py-1.5 rounded-lg">Protected</span>
                                    </div>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>

        <div v-if="showModal" class="fixed inset-0 bg-slate-900/80 backdrop-blur-sm z-[100] flex items-center justify-center p-4">
            <div class="bg-white rounded-[32px] shadow-2xl w-full max-w-4xl flex flex-col max-h-[95vh] overflow-hidden border border-slate-100">
                
                <div class="p-6 border-b border-slate-50 flex justify-between items-center bg-slate-50/80 shrink-0">
                    <h3 class="font-black text-xl text-slate-800 tracking-tighter uppercase italic">{{ isEditMode ? 'Update Profile' : 'New Registration' }}</h3>
                    <button @click="closeModal" class="text-slate-400 hover:text-red-500 transition-colors bg-white p-2 rounded-xl shadow-sm border border-slate-100">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                    </button>
                </div>
                
                <div class="overflow-y-auto custom-scrollbar p-6 md:p-8">
                    <form @submit.prevent="submit" class="flex flex-col gap-6">
                        
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
                            <div>
                                <label class="block text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1.5">Nama Lengkap</label>
                                <input v-model="form.name" type="text" required class="w-full p-3.5 bg-slate-50 rounded-2xl outline-none border border-slate-200 focus:border-blue-600 focus:ring-4 focus:ring-blue-600/10 font-bold text-sm transition-all">
                            </div>
                            
                            <div>
                                <label class="block text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1.5">Jabatan / Role</label>
                                <select v-model="form.role" required class="w-full p-3.5 bg-slate-50 rounded-2xl outline-none border border-slate-200 focus:border-blue-600 focus:ring-4 focus:ring-blue-600/10 font-bold text-sm text-slate-700 cursor-pointer transition-all uppercase">
                                    <option value="manager">Manager</option>
                                    <option value="supervisor">Supervisor</option>
                                    <option value="staff">Staff / Cashier</option>
                                </select>
                            </div>

                            <div class="grid grid-cols-2 gap-3">
                                <div>
                                    <label class="block text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1.5">Tempat Lahir</label>
                                    <input v-model="form.tempat_lahir" type="text" required class="w-full p-3.5 bg-slate-50 rounded-2xl outline-none border border-slate-200 focus:border-blue-600 font-bold text-sm transition-all">
                                </div>
                                <div>
                                    <label class="block text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1.5">Tgl Lahir</label>
                                    <input v-model="form.tanggal_lahir" type="date" required class="w-full p-3.5 bg-slate-50 rounded-2xl outline-none border border-slate-200 focus:border-blue-600 font-bold text-xs uppercase tracking-widest transition-all">
                                </div>
                            </div>
                            
                            <div>
                                <label class="block text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1.5">No. Handphone (WA)</label>
                                <input v-model="form.no_hp" type="text" placeholder="08..." required class="w-full p-3.5 bg-slate-50 rounded-2xl outline-none border border-slate-200 focus:border-blue-600 font-bold text-sm transition-all">
                            </div>
                            
                            <div class="md:col-span-2">
                                <label class="block text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1.5">Password <span v-if="isEditMode" class="text-amber-500 italic lowercase tracking-normal">(Kosongi jika tidak diubah)</span></label>
                                <input v-model="form.password" type="password" :required="!isEditMode" class="w-full p-3.5 bg-slate-50 rounded-2xl outline-none border border-slate-200 focus:border-blue-600 font-bold text-sm transition-all placeholder:text-slate-300" placeholder="••••••••">
                            </div>
                        </div>

                        <div class="grid grid-cols-1 md:grid-cols-2 gap-5 pt-6 border-t border-slate-100">
                            <div class="flex flex-col items-center p-5 bg-slate-50 rounded-[28px] border border-slate-200">
                                <label class="text-[9px] font-black text-slate-500 uppercase tracking-widest mb-4 flex items-center gap-2">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-blue-500" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
                                    1. Foto Profil (ID Card)
                                </label>
                                <div @click="$refs.profileInput.click()" class="w-36 h-36 bg-white rounded-[20px] shadow-sm flex items-center justify-center cursor-pointer overflow-hidden border border-slate-200 group relative transition-all hover:border-blue-500">
                                    <img v-if="fotoProfilPreview" :src="fotoProfilPreview" class="w-full h-full object-cover">
                                    <div v-else class="text-center p-4">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8 mx-auto text-slate-300 group-hover:text-blue-500 transition-colors" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect width="18" height="18" x="3" y="3" rx="2" ry="2"/><circle cx="9" cy="9" r="2"/><path d="m21 15-3.086-3.086a2 2 0 0 0-2.828 0L6 21"/></svg>
                                    </div>
                                    <div class="absolute inset-0 bg-slate-900/60 opacity-0 group-hover:opacity-100 flex items-center justify-center transition-opacity backdrop-blur-sm">
                                        <span class="text-white text-[9px] font-black uppercase tracking-widest">Pilih File</span>
                                    </div>
                                </div>
                                <input type="file" ref="profileInput" @change="onProfileFileChange" class="hidden" accept="image/*">
                            </div>

                            <div v-if="attendanceMethod !== 'other'" class="flex flex-col items-center p-5 bg-slate-900 rounded-[28px] relative overflow-hidden shadow-inner">
                                <label class="text-[9px] font-black text-blue-400 uppercase tracking-widest mb-4 flex items-center gap-2 relative z-10">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10"/></svg>
                                    2. Biometrik (Face)
                                </label>

                                <div v-if="attendanceMethod === 'face'" class="w-full flex flex-col items-center gap-4 relative z-10">
                                    <div class="w-36 h-36 bg-slate-950 rounded-[20px] overflow-hidden relative border-2 border-slate-700 shadow-2xl">
                                        <video ref="video" v-show="isCameraOpen" autoplay muted playsinline class="w-full h-full object-cover scale-x-[-1]"></video>
                                        <img v-if="fotoBiometricPreview && !isCameraOpen" :src="fotoBiometricPreview" class="w-full h-full object-cover">
                                        <canvas ref="canvas" class="hidden"></canvas>
                                        
                                        <div v-if="isCameraOpen" class="absolute inset-0 border-2 border-blue-500/50 rounded-[20px] pointer-events-none">
                                            <div class="w-full h-0.5 bg-blue-500 absolute top-0 animate-[scan_2s_infinite]"></div>
                                        </div>
                                    </div>

                                    <div class="flex gap-2">
                                        <button v-if="!isCameraOpen" @click="startCamera" type="button" class="bg-blue-600 hover:bg-blue-500 text-white px-5 py-2.5 rounded-xl font-black text-[9px] uppercase tracking-widest flex items-center gap-2 transition-colors">
                                            <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M23 19a2 2 0 0 1-2-2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z"/><circle cx="12" cy="13" r="4"/></svg>
                                            Kamera
                                        </button>
                                        <button v-else @click="capturePhoto" type="button" class="bg-emerald-500 hover:bg-emerald-400 text-white px-6 py-2.5 rounded-xl font-black text-[9px] uppercase tracking-widest animate-pulse transition-colors">
                                            Jepret
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <div class="pt-4 border-t border-slate-100">
                            <button type="submit" :disabled="isProcessing" class="w-full bg-blue-600 hover:bg-slate-900 text-white py-4 rounded-[20px] font-black text-xs shadow-xl shadow-blue-200 transition-all active:scale-95 disabled:opacity-50 uppercase tracking-[0.2em] flex items-center justify-center gap-3">
                                <template v-if="isProcessing">
                                    <div class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                                    MEMPROSES...
                                </template>
                                <template v-else-if="isEditMode">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/><polyline points="17 21 17 13 7 13 7 21"/><polyline points="7 3 7 8 15 8"/></svg>
                                    SIMPAN PERUBAHAN
                                </template>
                                <template v-else>
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><line x1="19" y1="8" x2="19" y2="14"/><line x1="22" y1="11" x2="16" y2="11"/></svg>
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
/* Custom Scrollbar for tables & modals */
.custom-scrollbar::-webkit-scrollbar { height: 6px; width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #94a3b8; }

/* Animasi untuk garis scanner biometrik */
@keyframes scan {
  0% { top: 0%; opacity: 0; }
  10% { opacity: 1; }
  90% { opacity: 1; }
  100% { top: 100%; opacity: 0; }
}
</style>