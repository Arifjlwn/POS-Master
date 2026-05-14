<script setup>
import { ref, onMounted } from 'vue';
import api from '../api.js';
import Sidebar from '../components/Sidebar.vue';
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
    foto: null
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
        <div class="p-6 md:p-10 max-w-7xl mx-auto font-sans">
            <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
                <div>
                    <h1 class="text-3xl font-black text-slate-900 tracking-tighter">Human Resources</h1>
                    <p class="text-sm text-slate-500 font-bold uppercase tracking-widest text-[10px]">Database & Akses Tim Toko</p>
                </div>
                <button @click="openAddModal" class="bg-blue-600 hover:bg-blue-700 text-white px-8 py-4 rounded-[24px] font-black text-sm flex items-center gap-3 shadow-xl shadow-blue-200 transition-all active:scale-95">
                    <span>👤</span> Tambah Karyawan Baru
                </button>
            </div>

            <div class="bg-white rounded-[35px] shadow-xl shadow-slate-200/50 border border-white overflow-hidden">
                <table class="w-full text-left whitespace-nowrap">
                    <thead class="bg-slate-50 border-b border-slate-100">
                        <tr>
                            <th class="px-8 py-5 text-[10px] font-black text-slate-400 uppercase tracking-[0.2em]">Karyawan</th>
                            <th class="px-8 py-5 text-[10px] font-black text-slate-400 uppercase tracking-[0.2em]">ID Login (NIK)</th>
                            <th class="px-8 py-5 text-[10px] font-black text-slate-400 uppercase tracking-[0.2em]">Jabatan</th>
                            <th class="px-8 py-5 text-[10px] font-black text-slate-400 uppercase tracking-[0.2em]">No. HP</th>
                            <th class="px-8 py-5 text-[10px] font-black text-slate-400 uppercase tracking-[0.2em] text-right">Aksi</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-slate-50 font-bold">
                        <tr v-if="isLoading">
                            <td colspan="5" class="px-8 py-10 text-center text-slate-400 italic">Memuat data tim...</td>
                        </tr>
                        <tr v-else-if="karyawan.length === 0">
                            <td colspan="5" class="px-8 py-10 text-center text-slate-400 italic">Belum ada data karyawan terdaftar.</td>
                        </tr>
                        <tr v-for="user in karyawan" :key="user.id" class="hover:bg-slate-50/50 transition-colors">
                            <td class="px-8 py-5">
                                <div class="flex items-center gap-4">
                                    <img 
                                        v-if="user.foto_url" 
                                        :src="API_BASE_URL + user.foto_url"
                                        class="w-12 h-12 rounded-2xl object-cover border-2 border-white shadow-md"
                                    >
                                    <div v-else class="w-12 h-12 rounded-2xl bg-blue-100 flex items-center justify-center text-blue-600 font-black text-xs">
                                        {{ user.name.substring(0, 2).toUpperCase() }}
                                    </div>
                                    <div>
                                        <div class="text-slate-900 uppercase text-sm">{{ user.name }}</div>
                                        <div class="text-[10px] text-slate-400">{{ user.tempat_lahir || 'N/A' }}, {{ user.tanggal_lahir || 'N/A' }}</div>
                                    </div>
                                </div>
                            </td>
                            <td class="px-8 py-5 text-blue-600 font-black tracking-widest">{{ user.nik }}</td>
                            
                            <td class="px-8 py-5">
                                <span v-if="user.role === 'owner'" class="bg-slate-900 text-white font-black px-3 py-1.5 rounded-xl text-[9px] uppercase tracking-wider shadow-sm">
                                    👑 OWNER
                                </span>
                                <span v-else-if="user.role === 'manager'" class="bg-purple-100 text-purple-700 font-black px-3 py-1.5 rounded-xl text-[9px] uppercase tracking-wider shadow-sm">
                                    💼 MANAGER
                                </span>
                                <span v-else-if="user.role === 'supervisor'" class="bg-blue-100 text-blue-700 font-black px-3 py-1.5 rounded-xl text-[9px] uppercase tracking-wider shadow-sm">
                                    ⚡ SUPERVISOR
                                </span>
                                <span v-else class="bg-green-100 text-green-700 font-black px-3 py-1.5 rounded-xl text-[9px] uppercase tracking-wider shadow-sm">
                                    📦 STAFF
                                </span>
                            </td>

                            <td class="px-8 py-5 text-slate-500 text-sm">{{ user.no_hp || '-' }}</td>
                            <td class="px-8 py-5 text-right space-x-2">
                                <button v-if="user.role !== 'owner'" @click="openEditModal(user)" class="text-amber-600 bg-amber-50 px-4 py-2 rounded-xl text-xs font-black hover:bg-amber-100 transition-all">EDIT</button>
                                <button v-if="user.role !== 'owner'" @click="deleteKaryawan(user.id)" class="text-red-500 bg-red-50 px-4 py-2 rounded-xl text-xs hover:bg-red-500 hover:text-white transition-all">PECAT</button>
                                <span v-else class="text-slate-300 text-xs italic font-medium pr-4">Sistem Terkunci</span>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>

        <div v-if="showModal" class="fixed inset-0 bg-slate-900/60 backdrop-blur-md z-[100] flex items-center justify-center p-4 overflow-y-auto">
            <div class="bg-white rounded-[40px] shadow-2xl w-full max-w-2xl my-8 overflow-hidden">
                <div class="p-8 border-b border-slate-100 flex justify-between items-center bg-slate-50">
                    <h3 class="font-black text-2xl text-slate-800 tracking-tighter uppercase">{{ isEditMode ? 'Update Data' : 'Registrasi' }} Karyawan</h3>
                    <button @click="closeModal" class="text-slate-400 hover:text-red-500 font-black text-2xl">✕</button>
                </div>
                
                <form @submit.prevent="submit" class="p-8 grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div class="space-y-4">
                        <div>
                            <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Nama Lengkap</label>
                            <input v-model="form.name" type="text" required class="w-full p-4 bg-slate-50 rounded-2xl outline-none border-2 border-transparent focus:border-blue-600 font-bold">
                        </div>
                        
                        <div>
                            <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Jabatan / Role Tim</label>
                            <select v-model="form.role" required class="w-full p-4 bg-slate-50 rounded-2xl outline-none border-2 border-transparent focus:border-blue-600 font-black text-sm text-gray-700 bg-white cursor-pointer">
                                <option value="manager">MANAGER</option>
                                <option value="supervisor">SUPERVISOR</option>
                                <option value="staff">STAFF /KASIR / PRAMUNIAGA</option>
                            </select>
                        </div>

                        <div class="grid grid-cols-2 gap-4">
                            <div>
                                <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Tempat Lahir</label>
                                <input v-model="form.tempat_lahir" type="text" required class="w-full p-4 bg-slate-50 rounded-2xl outline-none font-bold">
                            </div>
                            <div>
                                <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Tgl Lahir</label>
                                <input v-model="form.tanggal_lahir" type="date" required class="w-full p-4 bg-slate-50 rounded-2xl outline-none font-bold text-xs uppercase">
                            </div>
                        </div>
                        <div>
                            <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">No. Handphone (WhatsApp)</label>
                            <input v-model="form.no_hp" type="text" placeholder="08..." required class="w-full p-4 bg-slate-50 rounded-2xl outline-none font-bold">
                        </div>
                        <div>
                            <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Password {{ isEditMode ? '(Isi jika ingin ganti)' : 'Login' }}</label>
                            <input v-model="form.password" type="password" :required="!isEditMode" class="w-full p-4 bg-slate-50 rounded-2xl outline-none font-bold border-2 border-dashed border-slate-200">
                        </div>
                    </div>

                    <div class="flex flex-col items-center justify-center border-2 border-dashed border-slate-200 rounded-[35px] p-6 bg-slate-50">
                        <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-4">Foto Wajah (Close-Up)</label>
                        <div @click="$refs.fileInput.click()" class="w-48 h-48 bg-white rounded-[30px] shadow-inner flex items-center justify-center cursor-pointer overflow-hidden border-4 border-white group relative">
                            <img v-if="fotoPreview" :src="fotoPreview" class="w-full h-full object-cover">
                            <span v-else class="text-4xl text-slate-200">📸</span>
                            <div class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 flex items-center justify-center transition-opacity">
                                <span class="text-white text-[10px] font-black uppercase">Ganti Foto</span>
                            </div>
                        </div>
                        <input type="file" ref="fileInput" @change="onFileChange" class="hidden" accept="image/*">
                        <p class="text-[9px] text-slate-400 mt-4 text-center px-4 font-bold uppercase leading-relaxed">Pastikan wajah terlihat jelas untuk sistem absensi wajah nantinya.</p>
                    </div>

                    <div class="md:col-span-2 pt-4">
                        <button type="submit" :disabled="isProcessing" class="w-full bg-slate-900 hover:bg-blue-600 text-white py-5 rounded-[25px] font-black text-lg shadow-xl transition-all active:scale-95 disabled:opacity-50 uppercase tracking-widest">
                            {{ isProcessing ? 'MENGUPLOAD DATA...' : (isEditMode ? '💾 SIMPAN PERUBAHAN' : '🚀 DAFTARKAN KARYAWAN') }}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    </Sidebar>
</template>