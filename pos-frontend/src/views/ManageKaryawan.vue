<script setup>
import { ref, onMounted } from 'vue';
import api from '../api.js';
import Sidebar from '../components/Sidebar.vue';

// State Data Karyawan
const karyawan = ref([]);
const isLoading = ref(true);

// State Modal dan Form
const showModal = ref(false);
const isProcessing = ref(false);

const form = ref({
    name: '',
    password: '',
});

// --- FUNGSI TARIK DATA (API GOLANG) ---
const fetchKaryawan = async () => {
    isLoading.value = true;
    try {
        const response = await api.get('/employees');
        karyawan.value = response.data.data || [];
    } catch (error) {
        console.error("Gagal menarik data karyawan:", error);
    } finally {
        isLoading.value = false;
    }
};

onMounted(() => fetchKaryawan());

// --- FUNGSI TAMBAH KARYAWAN BARU ---
const submit = async () => {
    isProcessing.value = true;
    try {
        const response = await api.post('/employees', {
            name: form.value.name,
            password: form.value.password,
        });

        const newEmployee = response.data.data;
        
        // 🚀 Tampilkan NIK yang digenerate Golang biar Bos bisa nyatet!
        alert(`BERHASIL! 🤝\nNama: ${newEmployee.nama}\nNIK: ${newEmployee.nik}\n\nSilakan berikan NIK dan Password tersebut kepada kasir untuk login.`);
        
        form.value = { name: '', password: '' };
        showModal.value = false;
        fetchKaryawan();
    } catch (error) {
        alert(error.response?.data?.error || 'Gagal mendaftarkan karyawan');
        console.error("Error submit:", error);
    } finally {
        isProcessing.value = false;
    }
};

// --- FUNGSI HAPUS KARYAWAN ---
const deleteKaryawan = async (id) => {
    if (confirm('Yakin ingin memecat dan menghapus akun karyawan ini? 😭')) {
        try {
            await api.delete(`/employees/${id}`);
            alert('Karyawan berhasil dihapus.');
            fetchKaryawan();
        } catch (error) {
            alert('Gagal menghapus karyawan.');
        }
    }
};
</script>

<template>
    <Sidebar>
        <div class="p-6 md:p-8 w-full max-w-full mx-auto font-sans">
            <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
                <div>
                    <h1 class="text-3xl font-black text-gray-900 leading-tight">Manajemen Karyawan</h1>
                    <p class="text-sm text-gray-500 font-medium mt-1">Kelola akses kasir dan tim tokomu di sini.</p>
                </div>
                <button @click="showModal = true" class="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-2xl font-black text-sm flex items-center gap-2 shadow-lg shadow-blue-200 transition-all active:scale-95">
                    <span>➕</span> Tambah Kasir Baru
                </button>
            </div>

            <div class="bg-white rounded-3xl shadow-sm border border-gray-100 overflow-hidden">
                <table class="w-full text-left">
                    <thead class="bg-gray-50 border-b border-gray-100">
                        <tr>
                            <th class="px-6 py-4 text-xs font-black text-gray-400 uppercase tracking-widest">Nama</th>
                            <th class="px-6 py-4 text-xs font-black text-gray-400 uppercase tracking-widest">NIK (ID Login)</th>
                            <th class="px-6 py-4 text-xs font-black text-gray-400 uppercase tracking-widest">Jabatan</th>
                            <th class="px-6 py-4 text-xs font-black text-gray-400 uppercase tracking-widest text-right">Aksi</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-gray-50">
                        <tr v-if="isLoading">
                            <td colspan="4" class="px-6 py-12 text-center text-gray-400 font-medium">Sedang memuat data...</td>
                        </tr>
                        <tr v-else-if="karyawan.length === 0">
                            <td colspan="4" class="px-6 py-12 text-center text-gray-400 font-medium italic">Belum ada karyawan yang didaftarkan.</td>
                        </tr>
                        <tr v-for="user in karyawan" :key="user.id" class="hover:bg-gray-50 transition-colors">
                            <td class="px-6 py-4">
                                <div class="flex items-center gap-3">
                                    <div class="w-9 h-9 rounded-full bg-blue-100 text-blue-600 flex items-center justify-center font-black text-xs uppercase">
                                        {{ user.name.substring(0,2) }}
                                    </div>
                                    <span class="font-bold text-gray-800">{{ user.name }}</span>
                                </div>
                            </td>
                            <td class="px-6 py-4 text-sm text-gray-900 font-black tracking-widest">{{ user.nik }}</td>
                            <td class="px-6 py-4">
                                <span class="bg-blue-100 text-blue-700 text-[10px] font-black px-2.5 py-1 rounded-lg uppercase tracking-wider">
                                    {{ user.role === 'kasir' ? 'Kasir' : user.role }}
                                </span>
                            </td>
                            <td class="px-6 py-4 text-right">
                                <button @click="deleteKaryawan(user.id)" class="text-red-400 hover:text-red-600 transition-colors bg-red-50 p-2 rounded-lg text-sm font-bold">
                                    🗑️ Pecat
                                </button>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>

        <div v-if="showModal" class="fixed inset-0 bg-gray-900/60 backdrop-blur-sm z-[100] flex items-center justify-center p-4">
            <div class="bg-white rounded-3xl shadow-2xl w-full max-w-md overflow-hidden transform transition-all border-t-8 border-blue-600">
                <div class="p-6 border-b border-gray-100 flex justify-between items-center">
                    <h3 class="font-black text-xl text-gray-800">Daftarkan Kasir</h3>
                    <button @click="showModal = false" class="text-gray-400 hover:text-red-500 font-black">✕</button>
                </div>
                <form @submit.prevent="submit" class="p-6 space-y-4">
                    <div>
                        <label class="block text-xs font-black text-gray-400 uppercase tracking-widest mb-1.5">Nama Lengkap</label>
                        <input v-model="form.name" type="text" required placeholder="Contoh: Budi Santoso" class="w-full px-4 py-3 bg-white border border-gray-200 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 font-bold text-sm outline-none">
                    </div>
                    <div>
                        <label class="block text-xs font-black text-gray-400 uppercase tracking-widest mb-1.5">Password Sementara</label>
                        <input v-model="form.password" type="password" required placeholder="Minimal 6 karakter" class="w-full px-4 py-3 bg-white border border-gray-200 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 font-bold text-sm outline-none">
                    </div>
                    <div class="pt-2">
                        <button type="submit" :disabled="isProcessing" class="w-full bg-blue-600 hover:bg-blue-700 text-white py-4 rounded-xl font-black text-sm shadow-lg shadow-blue-100 transition-all active:scale-95 disabled:opacity-50">
                            {{ isProcessing ? 'Menyimpan...' : 'Generate NIK & Daftarkan 🚀' }}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    </Sidebar>
</template>