<script setup>
import { ref, onMounted } from 'vue';
import SidebarLaundry from './SidebarLaundry.vue';
import api from '../../../api.js'; 
import Swal from 'sweetalert2';

const kasirList = ref([]);
const isLoading = ref(false);
const isSaving = ref(false);

const form = ref({
    name: '',
    no_hp: '',
    password: ''
});

// Ambil daftar kasir dari backend
const fetchKasir = async () => {
    isLoading.value = true;
    try {
        const response = await api.get('/laundry/kasir');
        kasirList.value = response.data || [];
    } catch (error) {
        Swal.fire('Gagal!', 'Tidak dapat memuat data karyawan.', 'error');
    } finally {
        isLoading.value = false;
    }
};

onMounted(() => {
    fetchKasir();
});

// Tambah kasir baru
const simpanKasir = async () => {
    isSaving.value = true;
    try {
        const payload = {
            name: form.value.name,
            no_hp: String(form.value.no_hp || ''),
            password: form.value.password
        };

        const res = await api.post('/laundry/kasir', payload);
        
        Swal.fire({
            icon: 'success',
            title: 'Berhasil!',
            text: `Kasir ditambahkan. Email login: ${res.data.email}`,
            confirmButtonColor: '#4f46e5'
        });

        // Reset form & Refresh data
        form.value = { name: '', no_hp: '', password: '' };
        fetchKasir();

    } catch (error) {
        Swal.fire('Gagal!', error.response?.data?.error || 'Gagal menambahkan kasir.', 'error');
    } finally {
        isSaving.value = false;
    }
};

// Hapus/Pecat kasir
const hapusKasir = async (id, name) => {
    const confirm = await Swal.fire({
        title: 'Cabut Akses?',
        text: `Apakah Anda yakin ingin menghapus akses login untuk kasir ${name}?`,
        icon: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#ef4444',
        cancelButtonColor: '#94a3b8',
        confirmButtonText: 'Ya, Hapus!',
        cancelButtonText: 'Batal'
    });

    if (confirm.isConfirmed) {
        try {
            await api.delete(`/laundry/kasir/${id}`);
            Swal.fire({ toast: true, position: 'top-end', icon: 'success', title: 'Akses kasir dicabut!', showConfirmButton: false, timer: 1500 });
            fetchKasir();
        } catch (error) {
            Swal.fire('Gagal!', 'Terjadi kesalahan saat menghapus data.', 'error');
        }
    }
};
</script>

<template>
    <SidebarLaundry>
        <div class="flex-1 flex flex-col h-full bg-[#F8FAFC] overflow-hidden">
            
            <div class="p-5 md:p-8 shrink-0 bg-white border-b border-slate-100 flex items-center gap-4 z-10 shadow-sm">
                <div class="w-12 h-12 bg-indigo-50 border border-indigo-100 rounded-2xl flex items-center justify-center shrink-0">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-indigo-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M22 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
                </div>
                <div>
                    <h1 class="text-xl md:text-2xl font-black tracking-tighter uppercase text-slate-800 leading-tight">Manajemen Karyawan</h1>
                    <p class="text-[9px] md:text-[10px] font-black text-slate-400 mt-1 uppercase tracking-widest">Kelola Akses Kasir & Staff Operasional</p>
                </div>
            </div>

            <div class="flex-1 overflow-y-auto custom-scrollbar p-5 md:p-8">
                <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 md:gap-8 items-start">
                    
                    <div class="bg-white border border-slate-200 p-6 md:p-8 rounded-[24px] shadow-sm">
                        <h3 class="text-xs font-black text-slate-400 uppercase tracking-widest border-b border-slate-100 pb-3 mb-6">Tambah Kasir Baru</h3>
                        
                        <form @submit.prevent="simpanKasir" class="space-y-5">
                            <div>
                                <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Nama Panggilan</label>
                                <input v-model="form.name" type="text" required class="w-full px-4 py-3 bg-slate-50 border-2 border-slate-200 rounded-xl outline-none focus:bg-white focus:border-indigo-500 font-bold text-sm text-slate-800 transition-all" placeholder="Contoh: Siti">
                            </div>

                            <div>
                                <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Nomor HP (Opsional)</label>
                                <input v-model="form.no_hp" type="number" class="w-full px-4 py-3 bg-slate-50 border-2 border-slate-200 rounded-xl outline-none focus:bg-white focus:border-indigo-500 font-bold text-sm text-slate-800 transition-all" placeholder="08123...">
                            </div>

                            <div>
                                <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Password Login</label>
                                <input v-model="form.password" type="text" class="w-full px-4 py-3 bg-slate-50 border-2 border-slate-200 rounded-xl outline-none focus:bg-white focus:border-indigo-500 font-bold text-sm text-slate-800 transition-all" placeholder="Default: kasir123">
                                <span class="text-[9px] font-bold text-slate-400 block mt-1">Kosongkan jika ingin menggunakan password default.</span>
                            </div>

                            <button type="submit" :disabled="isSaving" class="w-full mt-4 bg-indigo-600 hover:bg-indigo-700 text-white py-4 rounded-xl font-black text-xs uppercase tracking-[0.2em] transition-all active:scale-95 shadow-xl shadow-indigo-600/20 flex items-center justify-center gap-2">
                                <template v-if="!isSaving">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
                                    DAFTARKAN KASIR
                                </template>
                                <template v-else>
                                    <div class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                                    MENYIMPAN...
                                </template>
                            </button>
                        </form>
                    </div>

                    <div class="lg:col-span-2 bg-white border border-slate-200 rounded-[24px] shadow-sm overflow-hidden flex flex-col">
                        <div class="p-6 md:p-8 border-b border-slate-100 flex justify-between items-center bg-slate-50/50">
                            <h3 class="text-xs font-black text-slate-400 uppercase tracking-widest">Daftar Akun Kasir Aktif</h3>
                            <span class="bg-indigo-100 text-indigo-700 font-black text-[10px] px-3 py-1 rounded-full uppercase tracking-widest">{{ kasirList.length }} Karyawan</span>
                        </div>
                        
                        <div v-if="isLoading" class="flex flex-col items-center justify-center py-20">
                            <div class="w-8 h-8 border-4 border-slate-200 border-t-indigo-600 rounded-full animate-spin mb-3"></div>
                            <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Memuat Data...</span>
                        </div>

                        <div v-else-if="kasirList.length === 0" class="flex flex-col items-center justify-center py-20 text-center px-4">
                            <div class="w-16 h-16 bg-slate-50 rounded-full flex items-center justify-center mb-4">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8 text-slate-300" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><line x1="19" y1="8" x2="19" y2="14"/><line x1="22" y1="11" x2="16" y2="11"/></svg>
                            </div>
                            <h4 class="text-sm font-black text-slate-600 uppercase">Belum Ada Kasir</h4>
                            <p class="text-[10px] font-bold text-slate-400 mt-1 max-w-xs">Tambahkan karyawan kasir pertama Anda melalui form di sebelah kiri.</p>
                        </div>

                        <div v-else class="flex-1 overflow-y-auto custom-scrollbar p-6">
                            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                                <div v-for="kasir in kasirList" :key="kasir.id" class="p-4 border-2 border-slate-100 rounded-2xl flex flex-col gap-4 hover:border-indigo-200 hover:shadow-md transition-all group bg-white">
                                    <div class="flex items-start justify-between">
                                        <div class="flex items-center gap-3">
                                            <div class="w-10 h-10 rounded-xl bg-gradient-to-br from-indigo-500 to-blue-600 flex items-center justify-center text-white font-black text-lg shadow-inner">
                                                {{ kasir.name.charAt(0).toUpperCase() }}
                                            </div>
                                            <div>
                                                <h4 class="font-black text-slate-800 text-sm capitalize">{{ kasir.name }}</h4>
                                                <span class="text-[9px] font-black text-emerald-600 bg-emerald-50 border border-emerald-100 px-2 py-0.5 rounded uppercase mt-1 inline-block">Role: Kasir</span>
                                            </div>
                                        </div>
                                        <button @click="hapusKasir(kasir.id, kasir.name)" class="text-slate-300 hover:text-red-500 transition-colors p-2 bg-slate-50 hover:bg-red-50 rounded-lg" title="Pecat Kasir">
                                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/></svg>
                                        </button>
                                    </div>
                                    
                                    <div class="bg-slate-50 p-3 rounded-xl border border-slate-100 space-y-2">
                                        <div class="flex items-center gap-2">
                                            <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 text-slate-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><rect width="20" height="16" x="2" y="4" rx="2"/><path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7"/></svg>
                                            <span class="text-[10px] font-bold text-slate-600 tracking-wide">{{ kasir.email }}</span>
                                        </div>
                                        <div class="flex items-center gap-2" v-if="kasir.no_hp">
                                            <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 text-slate-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><rect width="14" height="20" x="5" y="2" rx="2" ry="2"/><path d="M12 18h.01"/></svg>
                                            <span class="text-[10px] font-bold text-slate-600 tracking-wide">{{ kasir.no_hp }}</span>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>

                </div>
            </div>
        </div>
    </SidebarLaundry>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { height: 6px; width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
input[type=number]::-webkit-inner-spin-button, input[type=number]::-webkit-outer-spin-button { -webkit-appearance: none; margin: 0; }
input[type=number] { -moz-appearance: textfield; }
</style>