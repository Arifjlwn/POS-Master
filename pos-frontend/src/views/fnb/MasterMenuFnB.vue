<script setup>
import { ref, onMounted } from 'vue';
import SidebarFnB from './SidebarFnB.vue';
import api from '../../api.js'; 
import Swal from 'sweetalert2';

const products = ref([]);
const isLoading = ref(false);
const isModalOpen = ref(false);
const isSubmitting = ref(false);
const userRole = ref('');

// State Form Tambah Menu (Stok dihilangkan)
const form = ref({
    nama: '',
    harga: '',
    kategori: 'makanan', // Default
    gambar: ''
});

const formatRupiah = (angka) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(angka || 0);

/// 🚀 1. TARIK DATA DARI SERVER (Di-upgrade biar pintar baca Inggris/Indo)
const fetchProducts = async () => {
    isLoading.value = true;
    try {
        const response = await api.get('/fnb/products');
        products.value = response.data.map(p => ({
            ...p,
            nama: p.nama_produk,
            harga: p.harga_jual,
            kategori: p.kategori,
            gambar: p.gambar || 'https://images.unsplash.com/photo-1546069901-ba9599a7e63c?auto=format&fit=crop&w=400&q=80'
        }));
    } catch (error) {
        Swal.fire('Gagal!', 'Tidak dapat memuat katalog.', 'error');
    } finally {
        isLoading.value = false;
    }
};

onMounted(() => {
    userRole.value = localStorage.getItem('role') || ''; // Tarik role dari login
    fetchProducts();
});

// 🚀 2. SIMPAN MENU BARU
const submitMenu = async () => {
    // PROTEKSI OWNER DI SISI CLIENT
    if (userRole.value !== 'owner') {
        return Swal.fire('Akses Ditolak!', 'Hanya Owner yang bisa menambah menu!', 'error');
    }

    isSubmitting.value = true;
    try {
        await api.post('/fnb/products', {
            nama: form.value.nama,
            harga: Number(form.value.harga),
            kategori: form.value.kategori,
            gambar: form.value.gambar || 'https://images.unsplash.com/photo-1546069901-ba9599a7e63c?auto=format&fit=crop&w=400&q=80'
        });
        
        Swal.fire('Sukses!', 'Menu berhasil ditambahkan!', 'success');
        isModalOpen.value = false;
        form.value = { nama: '', harga: '', kategori: 'makanan', gambar: '' };
        fetchProducts();
    } catch (error) {
        Swal.fire('Error', 'Gagal menyimpan menu.', 'error');
    } finally {
        isSubmitting.value = false;
    }
};

// 🚀 3. TOGGLE KETERSEDIAAN MENU (ON/OFF)
const toggleStatus = async (product) => {
    try {
        const res = await api.put(`/fnb/products/${product.id}/toggle`);
        Swal.fire({ toast: true, position: 'top-end', icon: res.data.is_available ? 'success' : 'warning', title: res.data.message, showConfirmButton: false, timer: 1500 });
        fetchProducts();
    } catch (error) {
        Swal.fire('Gagal!', 'Tidak dapat mengubah status.', 'error');
    }
};
</script>

<template>
    <SidebarFnB>
        <div class="flex-1 flex flex-col h-screen bg-[#F8FAFC] overflow-hidden relative font-sans">
            
            <div class="p-5 md:p-8 bg-white border-b border-slate-200 shrink-0 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 z-10 shadow-sm">
                <div>
                    <h1 class="text-xl md:text-3xl font-black tracking-tighter uppercase text-slate-800 leading-none">Master Menu</h1>
                    <p class="text-[10px] md:text-xs font-black text-slate-400 mt-1.5 uppercase tracking-widest">Katalog Produk & Status Ketersediaan</p>
                </div>
                <button v-if="userRole === 'owner'" @click="isModalOpen = true" class="w-full sm:w-auto bg-indigo-600 hover:bg-slate-900 text-white px-6 py-3.5 rounded-xl font-black text-xs uppercase tracking-widest shadow-lg shadow-indigo-200 transition-all active:scale-95 flex items-center justify-center gap-2">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
                    Tambah Menu
                </button>
            </div>

            <div class="flex-1 overflow-y-auto custom-scrollbar p-5 md:p-8">
                
                <div v-if="isLoading" class="flex flex-col items-center justify-center h-64">
                    <div class="w-8 h-8 rounded-full border-4 border-slate-200 border-t-indigo-600 animate-spin mb-4"></div>
                    <span class="font-black text-xs uppercase tracking-widest text-slate-400">Memuat Katalog...</span>
                </div>

                <div v-else-if="products.length === 0" class="flex flex-col items-center justify-center h-64 bg-white rounded-[32px] border-2 border-dashed border-slate-200">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-16 h-16 text-slate-300 mb-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M12 2v20M17 5H9.5a4.5 4.5 0 0 0 0 9H15a4.5 4.5 0 0 1 0 9H6.5"/></svg>
                    <h2 class="font-black text-lg text-slate-800 uppercase tracking-widest mb-1">Katalog Masih Kosong</h2>
                    <p class="text-xs font-bold text-slate-400 uppercase tracking-widest text-center">Yuk, tambahkan menu andalan resto kamu!</p>
                </div>

                <div v-else class="bg-white rounded-[24px] border border-slate-200 shadow-sm overflow-hidden">
                    <div class="overflow-x-auto">
                        <table class="w-full text-left border-collapse">
                            <thead>
                                <tr class="bg-slate-50 border-b border-slate-200">
                                    <th class="p-4 text-[10px] font-black text-slate-400 uppercase tracking-widest">Menu</th>
                                    <th class="p-4 text-[10px] font-black text-slate-400 uppercase tracking-widest">Kategori</th>
                                    <th class="p-4 text-[10px] font-black text-slate-400 uppercase tracking-widest">Harga</th>
                                    <th class="p-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-center">Status Jualan</th>
                                </tr>
                            </thead>
                            <tbody class="divide-y divide-slate-100">
                                <tr v-for="prod in products" :key="prod.id" class="hover:bg-slate-50/50 transition-colors group" :class="!prod.is_available ? 'bg-slate-50/50 grayscale-[20%]' : ''">
                                    <td class="p-4">
                                        <div class="flex items-center gap-4">
                                            <div class="w-12 h-12 rounded-xl bg-slate-100 overflow-hidden shrink-0 border border-slate-200">
                                                <img :src="prod.gambar" alt="Foto Menu" class="w-full h-full object-cover" :class="!prod.is_available ? 'opacity-70' : ''">
                                            </div>
                                            <div>
                                                <p class="font-black text-xs md:text-sm text-slate-800 uppercase">{{ prod.nama }}</p>
                                                <p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-0.5">ID: #{{ prod.id }}</p>
                                            </div>
                                        </div>
                                    </td>
                                    <td class="p-4">
                                        <span class="bg-slate-100 text-slate-600 px-3 py-1 rounded-lg text-[9px] font-black uppercase tracking-widest border border-slate-200">
                                            {{ prod.kategori }}
                                        </span>
                                    </td>
                                    <td class="p-4 font-black text-sm text-indigo-600">{{ formatRupiah(prod.harga) }}</td>
                                    <td class="p-4 text-center">
                                        <button @click="toggleStatus(prod)" class="relative inline-flex h-7 w-12 items-center rounded-full transition-colors focus:outline-none" :class="prod.is_available ? 'bg-emerald-500' : 'bg-slate-300'">
                                            <span class="inline-block h-5 w-5 transform rounded-full bg-white transition-transform shadow-sm" :class="prod.is_available ? 'translate-x-6' : 'translate-x-1'"></span>
                                        </button>
                                        <p class="text-[8px] font-black mt-1 uppercase tracking-widest" :class="prod.is_available ? 'text-emerald-600' : 'text-rose-500'">{{ prod.is_available ? 'Tersedia' : 'Habis' }}</p>
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>

        </div>

        <div v-if="isModalOpen" class="fixed inset-0 z-[100] flex items-center justify-center px-4">
            <div @click="isModalOpen = false" class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm animate-[fadeIn_0.2s_ease-out]"></div>
            
            <div class="bg-white rounded-[32px] shadow-2xl w-full max-w-sm md:max-w-md relative z-10 animate-[slideUp_0.3s_ease-out] overflow-hidden flex flex-col max-h-[90vh]">
                
                <div class="p-6 border-b border-slate-100 flex justify-between items-center bg-slate-50/50 shrink-0">
                    <h2 class="font-black text-base uppercase tracking-widest text-slate-800">Rilis Menu Baru</h2>
                    <button @click="isModalOpen = false" class="w-8 h-8 bg-white border border-slate-200 text-slate-400 rounded-full flex items-center justify-center active:scale-90 hover:text-slate-800 transition-all shadow-sm">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                    </button>
                </div>

                <div class="p-6 overflow-y-auto custom-scrollbar flex-1">
                    <form id="formAddMenu" @submit.prevent="submitMenu" class="space-y-5">
                        
                        <div>
                            <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2 ml-1">Nama Menu</label>
                            <input v-model="form.nama" type="text" required placeholder="Cth: Ayam Bakar Madu" class="w-full px-4 py-3.5 bg-slate-50 border-2 border-slate-100 rounded-2xl outline-none font-bold text-sm text-slate-800 focus:bg-white focus:border-indigo-500 transition-all placeholder:text-slate-300">
                        </div>

                        <div>
                            <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2 ml-1">Kategori</label>
                            <div class="relative">
                                <select v-model="form.kategori" class="w-full pl-4 pr-10 py-3.5 bg-slate-50 border-2 border-slate-100 rounded-2xl outline-none font-bold text-sm text-slate-800 focus:bg-white focus:border-indigo-500 transition-all appearance-none cursor-pointer">
                                    <option value="makanan">Makanan Utama</option>
                                    <option value="minuman">Minuman Segar</option>
                                    <option value="paket">Paket Combo</option>
                                    <option value="snack">Snack & Ala Carte</option>
                                    <option value="dessert">Dessert & Manisan</option>
                                </select>
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-400 absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="m6 9 6 6 6-6"/></svg>
                            </div>
                        </div>

                        <div>
                            <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2 ml-1">Harga Jual (Rp)</label>
                            <input v-model="form.harga" type="number" required placeholder="Cth: 25000" class="w-full px-4 py-3.5 bg-slate-50 border-2 border-slate-100 rounded-2xl outline-none font-black text-base text-indigo-600 focus:bg-white focus:border-indigo-500 transition-all placeholder:text-slate-300">
                        </div>

                        <div>
                            <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2 ml-1">Link Gambar (Opsional)</label>
                            <input v-model="form.gambar" type="url" placeholder="https://..." class="w-full px-4 py-3.5 bg-slate-50 border-2 border-slate-100 rounded-2xl outline-none font-bold text-xs text-slate-600 focus:bg-white focus:border-indigo-500 transition-all placeholder:text-slate-300">
                            <p class="text-[9px] font-bold text-slate-400 mt-2 ml-1">Kosongkan jika ingin menggunakan gambar bawaan sistem.</p>
                        </div>
                    </form>
                </div>

                <div class="p-6 border-t border-slate-100 bg-white shrink-0">
                    <button type="submit" form="formAddMenu" :disabled="isSubmitting" class="w-full bg-indigo-600 hover:bg-slate-900 text-white py-4.5 rounded-2xl font-black text-xs uppercase tracking-widest shadow-xl shadow-indigo-200 transition-all active:scale-95 disabled:opacity-50 flex items-center justify-center gap-2">
                        <template v-if="!isSubmitting">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M5 12h14"/><path d="m12 5 7 7-7 7"/></svg>
                            Simpan Menu Baru
                        </template>
                        <template v-else>
                            <div class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                            MENYIMPAN...
                        </template>
                    </button>
                </div>

            </div>
        </div>
    </SidebarFnB>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 6px; height: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }

@keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }
@keyframes slideUp { from { opacity: 0; transform: translateY(20px) scale(0.95); } to { opacity: 1; transform: translateY(0) scale(1); } }
</style>