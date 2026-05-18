<script setup>
import { ref, onMounted, computed } from 'vue';
import SidebarLaundry from './SidebarLaundry.vue';
import api from '../../../api.js';
import Swal from 'sweetalert2';

const services = ref([]);
const isLoading = ref(false);

// 🚀 STATE LIVE SEARCH & EDIT
const searchQuery = ref('');
const isEditing = ref(false);
const editId = ref(null);

// State Form Tambah/Edit
const showForm = ref(false);
const form = ref({
    nama_produk: '',
    harga_jual: '',
    satuan_dasar: 'KG',
    estimasi: '1 Hari'
});

const formatRupiah = (angka) => {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(angka || 0);
};

// 🚀 LIVE SEARCH COMPUTED
const filteredServices = computed(() => {
    if (!searchQuery.value) return services.value;
    return services.value.filter(s => 
        s.nama_produk.toLowerCase().includes(searchQuery.value.toLowerCase())
    );
});

// Ambil Data
const fetchServices = async () => {
    isLoading.value = true;
    try {
        const response = await api.get('/laundry/services');
        services.value = response.data;
    } catch (error) {
        Swal.fire('Gagal!', 'Gagal mengambil data layanan.', 'error');
    } finally {
        isLoading.value = false;
    }
};

onMounted(() => fetchServices());

// 🚀 MASUK MODE EDIT
const triggerEdit = (item) => {
    isEditing.value = true;
    editId.value = item.id;
    form.value = {
        nama_produk: item.nama_produk,
        harga_jual: item.harga_jual,
        satuan_dasar: item.satuan_dasar,
        estimasi: item.estimasi || '1 Hari'
    };
    showForm.value = true;
    // Auto scroll ke atas biar enak
    window.scrollTo({ top: 0, behavior: 'smooth' });
};

// 🚀 BATALKAN FORM
const cancelForm = () => {
    showForm.value = false;
    setTimeout(() => {
        isEditing.value = false;
        editId.value = null;
        form.value = { nama_produk: '', harga_jual: '', satuan_dasar: 'KG', estimasi: '1 Hari' };
    }, 300); // Tunggu animasi kelar
};

// Simpan Data (POST atau PUT)
const saveService = async () => {
    if(!form.value.nama_produk || !form.value.harga_jual) {
        return Swal.fire('Oops!', 'Nama dan Harga wajib diisi!', 'warning');
    }

    try {
        const payload = {
            nama_produk: form.value.nama_produk,
            harga_jual: parseFloat(form.value.harga_jual),
            satuan_dasar: form.value.satuan_dasar,
            estimasi: form.value.estimasi
        };

        if (isEditing.value) {
            await api.put(`/laundry/services/${editId.value}`, payload);
            Swal.fire({ toast: true, position: 'top-end', icon: 'success', title: 'Perubahan Disimpan!', showConfirmButton: false, timer: 1500 });
        } else {
            await api.post('/laundry/services', payload);
            Swal.fire({ toast: true, position: 'top-end', icon: 'success', title: 'Layanan Baru Disimpan!', showConfirmButton: false, timer: 1500 });
        }
        
        cancelForm();
        fetchServices();
    } catch (error) {
        Swal.fire('Gagal!', 'Terjadi kesalahan sistem.', 'error');
    }
};

// Hapus Data
const deleteService = (id, nama) => {
    Swal.fire({
        title: `Hapus ${nama}?`, 
        text: 'Tindakan ini tidak bisa dibatalkan!', 
        icon: 'warning',
        showCancelButton: true, 
        confirmButtonColor: '#e11d48', 
        cancelButtonColor: '#94a3b8',
        confirmButtonText: 'Ya, Hapus'
    }).then(async (result) => {
        if (result.isConfirmed) {
            try {
                await api.delete(`/laundry/services/${id}`);
                fetchServices();
                Swal.fire({ toast: true, position: 'top-end', icon: 'success', title: 'Terhapus!', showConfirmButton: false, timer: 1500 });
            } catch (error) {
                Swal.fire('Gagal!', 'Gagal menghapus layanan.', 'error');
            }
        }
    });
};
</script>

<template>
    <SidebarLaundry>
        <div class="flex-1 flex flex-col h-full bg-[#F8FAFC] overflow-hidden relative">
            
            <div class="p-5 md:p-8 shrink-0 bg-white border-b border-slate-100 flex flex-col z-10 shadow-sm relative gap-5">
                
                <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                    <div class="flex items-center gap-4">
                        <div class="w-12 h-12 bg-indigo-50 border border-indigo-100 rounded-2xl flex items-center justify-center shrink-0">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-indigo-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/><polyline points="3.27 6.96 12 12.01 20.73 6.96"/><line x1="12" y1="22.08" x2="12" y2="12"/></svg>
                        </div>
                        <div>
                            <h1 class="text-xl md:text-2xl font-black tracking-tighter uppercase text-slate-800 leading-tight">Master Layanan</h1>
                            <p class="text-[9px] md:text-[10px] font-black text-slate-400 mt-1 uppercase tracking-widest">Katalog Harga Cuci Kiloan & Satuan</p>
                        </div>
                    </div>
                    
                    <button @click="showForm ? cancelForm() : (showForm = true)" :class="showForm ? 'bg-rose-500 hover:bg-rose-600 shadow-rose-200' : 'bg-indigo-600 hover:bg-indigo-700 shadow-indigo-200'" class="w-full sm:w-auto text-white px-5 py-3.5 rounded-xl font-black text-xs uppercase tracking-widest transition-all active:scale-95 shadow-lg flex items-center justify-center gap-2">
                        <svg v-if="!showForm" xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
                        <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                        {{ showForm ? 'Batal Form' : 'Tambah Jasa' }}
                    </button>
                </div>

                <div class="relative w-full group max-w-2xl">
                    <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-400 group-focus-within:text-indigo-500 transition-colors" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
                    </div>
                    <input v-model="searchQuery" type="text" placeholder="Ketik nama layanan yang dicari..." class="w-full pl-12 pr-4 py-3 bg-slate-50 border-2 border-slate-200 rounded-xl focus:bg-white focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 outline-none font-bold text-sm text-slate-800 transition-all placeholder:text-slate-300 shadow-inner">
                </div>
            </div>

            <div class="flex-1 overflow-y-auto custom-scrollbar p-5 md:p-8">
                
                <div v-if="showForm" class="bg-white p-6 md:p-8 rounded-[24px] border border-slate-200 shadow-xl shadow-slate-200/50 mb-8 animate-[fadeInDown_0.3s_ease-out]">
                    <div class="flex items-center gap-3 mb-6 border-b border-slate-100 pb-4">
                        <div class="p-2 bg-indigo-50 text-indigo-600 rounded-lg">
                            <svg v-if="!isEditing" xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/></svg>
                            <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
                        </div>
                        <h3 class="text-sm font-black text-slate-800 uppercase tracking-widest">
                            {{ isEditing ? 'Edit Layanan' : 'Input Layanan Baru' }}
                        </h3>
                    </div>
                    
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-5 md:gap-6">
                        <div>
                            <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Nama Layanan / Paket</label>
                            <div class="relative group">
                                <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none"><svg xmlns="http://www.w3.org/2000/svg" class="w-4.5 h-4.5 text-slate-400 group-focus-within:text-indigo-500 transition-colors" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z"/><line x1="7" y1="7" x2="7.01" y2="7"/></svg></div>
                                <input v-model="form.nama_produk" type="text" placeholder="Contoh: Cuci Karpet Tebal" class="w-full pl-11 pr-4 py-3.5 bg-slate-50 border-2 border-slate-200 rounded-xl outline-none focus:bg-white focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 font-bold text-sm text-slate-800 transition-all placeholder:text-slate-300">
                            </div>
                        </div>
                        
                        <div>
                            <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Harga Jual</label>
                            <div class="relative group">
                                <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none"><span class="text-sm font-bold text-slate-400 group-focus-within:text-indigo-500 transition-colors">Rp</span></div>
                                <input v-model="form.harga_jual" type="number" placeholder="0" class="w-full pl-11 pr-4 py-3.5 bg-slate-50 border-2 border-slate-200 rounded-xl outline-none focus:bg-white focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 font-black text-sm text-slate-800 transition-all placeholder:text-slate-300">
                            </div>
                        </div>
                        
                        <div>
                            <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Tipe Satuan</label>
                            <div class="relative group">
                                <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none"><svg xmlns="http://www.w3.org/2000/svg" class="w-4.5 h-4.5 text-slate-400 group-focus-within:text-indigo-500 transition-colors" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><rect x="2" y="7" width="20" height="14" rx="2" ry="2"/><path d="M16 21V5a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16"/></svg></div>
                                <select v-model="form.satuan_dasar" class="w-full pl-11 pr-10 py-3.5 bg-slate-50 border-2 border-slate-200 rounded-xl outline-none focus:bg-white focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 font-bold text-sm text-slate-800 appearance-none cursor-pointer transition-all">
                                    <option value="KG">Kiloan (KG)</option>
                                    <option value="PCS">Satuan (PCS / Helai)</option>
                                    <option value="M2">Meter Persegi (M2)</option>
                                </select>
                                <div class="absolute inset-y-0 right-4 flex items-center pointer-events-none"><svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="m6 9 6 6 6-6"/></svg></div>
                            </div>
                        </div>
                        
                        <div>
                            <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Estimasi Waktu Pengerjaan</label>
                            <div class="relative group">
                                <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none"><svg xmlns="http://www.w3.org/2000/svg" class="w-4.5 h-4.5 text-slate-400 group-focus-within:text-indigo-500 transition-colors" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg></div>
                                <select v-model="form.estimasi" class="w-full pl-11 pr-10 py-3.5 bg-slate-50 border-2 border-slate-200 rounded-xl outline-none focus:bg-white focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 font-bold text-sm text-slate-800 appearance-none cursor-pointer transition-all">
                                    <option value="3 Jam">Express (3 Jam)</option>
                                    <option value="6 Jam">Express (6 Jam)</option>
                                    <option value="1 Hari">Kilat (1 Hari)</option>
                                    <option value="2 Hari">Reguler (2 Hari)</option>
                                    <option value="3 Hari">Standar (3 Hari)</option>
                                    <option value="7 Hari">Santai (7 Hari)</option>
                                </select>
                                <div class="absolute inset-y-0 right-4 flex items-center pointer-events-none"><svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="m6 9 6 6 6-6"/></svg></div>
                            </div>
                        </div>
                    </div>
                    
                    <div class="mt-8 flex justify-end gap-3">
                        <button v-if="isEditing" @click="cancelForm" class="bg-slate-100 hover:bg-slate-200 text-slate-600 px-6 py-4 rounded-xl font-black text-xs uppercase tracking-widest transition-all active:scale-95">Batal</button>
                        <button @click="saveService" class="w-full md:w-auto bg-slate-900 hover:bg-slate-800 text-white px-8 py-4 rounded-xl font-black text-xs uppercase tracking-[0.15em] transition-all active:scale-95 shadow-xl shadow-slate-900/20 flex items-center justify-center gap-2">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/><polyline points="17 21 17 13 7 13 7 21"/><polyline points="7 3 7 8 15 8"/></svg>
                            {{ isEditing ? 'Simpan Perubahan' : 'Simpan Katalog Jasa' }}
                        </button>
                    </div>
                </div>

                <div v-if="isLoading" class="flex flex-col items-center justify-center py-20 text-slate-400">
                    <div class="w-12 h-12 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin mb-4"></div>
                    <p class="font-black text-xs uppercase tracking-[0.2em] animate-pulse">Memuat Katalog...</p>
                </div>
                
                <div v-else-if="filteredServices.length === 0" class="flex flex-col items-center justify-center py-20 text-slate-400 bg-white rounded-3xl border border-slate-100 border-dashed border-2">
                    <div class="w-20 h-20 bg-slate-50 rounded-full flex items-center justify-center mb-4">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-10 h-10 text-slate-300" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/><polyline points="3.27 6.96 12 12.01 20.73 6.96"/><line x1="12" y1="22.08" x2="12" y2="12"/></svg>
                    </div>
                    <p class="font-black text-sm uppercase tracking-widest text-slate-500 mb-1">Pencarian Tidak Ditemukan</p>
                    <p class="text-xs font-bold">Coba ketik kata kunci yang lain</p>
                </div>
                
                <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4 md:gap-5">
                    <div v-for="item in filteredServices" :key="item.id" class="bg-white rounded-[24px] border-2 border-slate-100 p-5 md:p-6 relative group hover:border-indigo-400 transition-all shadow-sm hover:shadow-xl hover:shadow-indigo-100 overflow-hidden flex flex-col justify-between h-40 md:h-48">
                        
                        <div class="absolute -right-6 -bottom-6 opacity-[0.03] group-hover:opacity-10 group-hover:scale-110 transition-all duration-500 pointer-events-none">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-32 h-32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/></svg>
                        </div>

                        <div class="absolute top-4 right-4 flex gap-2 sm:opacity-0 group-hover:opacity-100 transition-all z-50">
                            <button @click.stop="triggerEdit(item)" class="w-8 h-8 bg-white border border-slate-100 text-slate-400 rounded-full flex items-center justify-center hover:bg-sky-500 hover:text-white hover:border-sky-500 shadow-sm transition-colors">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
                            </button>
                            <button @click.stop="deleteService(item.id, item.nama_produk)" class="w-8 h-8 bg-white border border-slate-100 text-slate-400 rounded-full flex items-center justify-center hover:bg-rose-500 hover:text-white hover:border-rose-500 shadow-sm transition-colors">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/><line x1="10" y1="11" x2="10" y2="17"/><line x1="14" y1="11" x2="14" y2="17"/></svg>
                            </button>
                        </div>

                        <div class="relative z-10 pointer-events-none">
                            <div class="mb-3">
                                <span class="text-[9px] font-black bg-indigo-50 text-indigo-600 px-2.5 py-1 rounded-lg uppercase tracking-widest border border-indigo-100 inline-block">Est: {{ item.estimasi || 'Standar' }}</span>
                            </div>
                            <h3 class="text-sm font-black text-slate-800 uppercase leading-tight line-clamp-2 pr-16">{{ item.nama_produk }}</h3>
                        </div>
                        
                        <div class="relative z-10 pt-4 border-t border-slate-100/50 mt-auto flex items-end justify-between pointer-events-none">
                            <p class="text-lg md:text-xl font-black text-emerald-500 leading-none">{{ formatRupiah(item.harga_jual) }}</p>
                            <span class="text-[9px] md:text-[10px] text-slate-400 font-bold uppercase tracking-widest">/ {{ item.satuan_dasar }}</span>
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
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #94a3b8; }

input[type=number]::-webkit-inner-spin-button, 
input[type=number]::-webkit-outer-spin-button { -webkit-appearance: none; margin: 0; }
input[type=number] { -moz-appearance: textfield; }

@keyframes fadeInDown {
    from { opacity: 0; transform: translateY(-15px); }
    to { opacity: 1; transform: translateY(0); }
}
</style>