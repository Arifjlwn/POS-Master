<script setup>
import { ref, computed } from 'vue';
import SidebarFnB from './SidebarFnB.vue';
import api from '../../api.js';
import Swal from 'sweetalert2';

// --- STATE UTAMA F&B ---
const activeCategory = ref('semua');
const searchQuery = ref('');
const tipeOrder = ref('DINE_IN'); // DINE_IN atau TAKE_AWAY
const nomorMeja = ref('');
const isMobileCartOpen = ref(false);

// --- STATE DATA (Biasanya di-fetch dari Golang)
const categories = ref([
    { id: 'semua', nama: 'Semua Menu' },
    { id: 'makanan', nama: 'Makanan Utama' },
    { id: 'minuman', nama: 'Minuman Segar' },
    { id: 'cemilan', nama: 'Cemilan / Dessert' }
]);

const products = ref([
    { id: 1, nama: 'Nasi Goreng Arzu Spesial', harga: 25000, kategori: 'makanan', stok: 20, gambar: 'https://images.unsplash.com/photo-1512058564366-18510be2db19?w=400' },
    { id: 2, nama: 'Mie Ayam Pangsit Badas', harga: 22000, kategori: 'makanan', stok: 15, gambar: 'https://images.unsplash.com/photo-1569718212165-3a8278d5f624?w=400' },
    { id: 3, nama: 'Es Teh Manis Jumbo (Gula Asli)', harga: 6000, kategori: 'minuman', stok: 100, gambar: 'https://images.unsplash.com/photo-1556679343-c7306c1976bc?w=400' },
    { id: 4, nama: 'Kopi Susu Gula Aren Dewa', harga: 18000, kategori: 'minuman', stok: 40, gambar: 'https://images.unsplash.com/photo-1541167760496-1628856ab772?w=400' },
    { id: 5, nama: 'Kentang Goreng Krispi Crispy', harga: 15000, kategori: 'cemilan', stok: 25, gambar: 'https://images.unsplash.com/photo-1573080496219-bb080dd4f877?w=400' }
]);

const cart = ref([]);

// --- UTILITIES ---
const formatRupiah = (angka) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(angka || 0);

// --- FILTER & PENCARIAN MENU ---
const filteredProducts = computed(() => {
    return products.value.filter(p => {
        const matchCategory = activeCategory.value === 'semua' || p.kategori === activeCategory.value;
        const matchSearch = p.nama.toLowerCase().includes(searchQuery.value.toLowerCase());
        return matchCategory && matchSearch;
    });
});

const totalBelanja = computed(() => {
    return cart.value.reduce((sum, item) => sum + (item.harga * item.qty), 0);
});

// --- CORE LOGIC CART F&B ---
const addToCart = (product) => {
    const existingItem = cart.value.find(item => item.id === product.id);
    if (existingItem) {
        existingItem.qty++;
    } else {
        cart.value.push({
            ...product,
            qty: 1,
            notes: '' // 🚀 Tambahan field notes khusus F&B untuk koki!
        });
    }
};

// 🚀 LOGIKA DECREASE QTY YANG KAMU KIRIM (DIPERTAHANKAN 100% AMAN BEB)
const decreaseQty = (product) => {
    const existingItem = cart.value.find(item => item.id === product.id);
    if (existingItem) {
        if (existingItem.qty > 1) {
            existingItem.qty--;
        } else {
            cart.value = cart.value.filter(item => item.id !== product.id);
            if (cart.value.length === 0) isMobileCartOpen.value = false;
        }
    }
};

// --- SIMPAN ORDER KE BACKEND GOLANG ---
const checkoutOrder = async () => {
    if (cart.value.length === 0) return Swal.fire('Keranjang Kosong', 'Pilih menu dulu dong beb!', 'warning');
    if (tipeOrder.value === 'DINE_IN' && !nomorMeja.value) return Swal.fire('Nomor Meja?', 'Dine In wajib isi nomor meja ya sayang!', 'warning');

    try {
        const payload = {
            session_id: Number(activeSessionId), // 🚀 Bawa session_id biar core_transaction gak menolak
            tipe_order: tipeOrder.value,
            nomor_meja: tipeOrder.value === 'DINE_IN' ? nomorMeja.value : '-',
            sub_total: totalBelanja.value, // Petakan nominal subtotal dasar
            pajak: 0, // Bisa kamu kalkulasi jika resto menggunakan PPN
            total_harga: totalBelanja.value,
            items: cart.value.map(item => ({
                product_id: item.id,
                qty: item.qty,
                sub_total: item.harga * item.qty,
                notes: item.notes // Bawa catatan koki ke database modul
            }))
        };

        await api.post('/fnb/order', payload); // Jalur endpoint F&B kamu

        Swal.fire({ 
            icon: 'success', 
            title: response.data.message, // "Pesanan berhasil diteruskan ke Monitor Dapur!"
            showConfirmButton: false, 
            timer: 1500 
        });

        // Reset Cart
        cart.value = [];
        nomorMeja.value = '';
        isMobileCartOpen.value = false;
    } catch (error) {
        Swal.fire('Gagal!', 'Terjadi kendala saat memproses pesanan koki.', 'error');
    }
};
</script>

<template>
    <SidebarFnB>
        <div class="flex-1 flex h-full bg-slate-50/50 overflow-hidden relative">
            
            <div class="flex-1 flex flex-col h-full overflow-hidden">
                <div class="p-4 md:p-6 bg-white border-b border-slate-200 shrink-0 space-y-4 shadow-sm">
                    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                        <div>
                            <h1 class="text-xl md:text-2xl font-black tracking-tighter uppercase text-slate-800">Kasir Resto & Cafe</h1>
                            <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest mt-0.5">Sistem POS FnB Dinamis</p>
                        </div>
                        <div class="relative w-full sm:w-72">
                            <input v-model="searchQuery" type="text" placeholder="Cari menu makanan/minuman..." class="w-full pl-10 pr-4 py-2.5 bg-slate-50 border-2 border-slate-200 rounded-xl focus:bg-white focus:border-indigo-500 font-bold text-xs text-slate-800 outline-none transition-all placeholder:text-slate-300">
                            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none"><svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg></div>
                        </div>
                    </div>

                    <div class="flex gap-2 overflow-x-auto hide-scrollbar pb-1">
                        <button v-for="cat in categories" :key="cat.id" @click="activeCategory = cat.id" :class="activeCategory === cat.id ? 'bg-indigo-600 text-white shadow-md shadow-indigo-100' : 'bg-slate-100 text-slate-500 hover:bg-slate-200'" class="whitespace-nowrap px-5 py-2 rounded-xl text-[10px] font-black uppercase tracking-widest transition-all">
                            {{ cat.nama }}
                        </button>
                    </div>
                </div>

                <div class="flex-1 overflow-y-auto custom-scrollbar p-4 md:p-6 pb-24">
                    <div v-if="filteredProducts.length === 0" class="text-center py-20 text-slate-400 font-bold text-xs uppercase tracking-widest">Menu tidak ditemukan beb.</div>
                    <div class="grid grid-cols-2 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
                        <div v-for="prod in filteredProducts" :key="prod.id" @click="addToCart(prod)" class="bg-white border border-slate-200 rounded-2xl overflow-hidden shadow-sm hover:shadow-xl hover:border-indigo-300 transition-all cursor-pointer flex flex-col group relative">
                            <div class="w-full h-32 md:h-40 bg-slate-100 relative overflow-hidden shrink-0">
                                <img :src="prod.gambar" :alt="prod.nama" class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500">
                                <div class="absolute inset-0 bg-slate-900/10 group-hover:bg-transparent transition-colors"></div>
                            </div>
                            <div class="p-3 md:p-4 flex flex-col flex-1 justify-between gap-2">
                                <h3 class="font-black text-xs md:text-sm text-slate-800 uppercase leading-tight line-clamp-2">{{ prod.nama }}</h3>
                                <p class="font-black text-xs md:text-base text-indigo-600 tracking-tight">{{ formatRupiah(prod.harga) }}</p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div class="hidden xl:flex w-96 bg-white border-l border-slate-200 flex-col h-full shadow-2xl shrink-0 z-20">
                <div class="p-5 border-b border-slate-100 space-y-4">
                    <h2 class="font-black text-xs uppercase tracking-[0.2em] text-slate-800 flex items-center gap-2"><svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-indigo-500" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M6 2 3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z"/><line x1="3" y1="6" x2="21" y2="6"/><path d="M16 10a4 4 0 0 1-8 0"/></svg> Daftar Pesanan</h2>
                    
                    <div class="grid grid-cols-2 gap-2 bg-slate-100 p-1 rounded-xl">
                        <button @click="tipeOrder = 'DINE_IN'" :class="tipeOrder === 'DINE_IN' ? 'bg-white text-slate-800 shadow-sm' : 'text-slate-400'" class="py-2 text-[10px] font-black uppercase tracking-widest rounded-lg transition-all">Dine In</button>
                        <button @click="tipeOrder = 'TAKE_AWAY'" :class="tipeOrder === 'TAKE_AWAY' ? 'bg-white text-slate-800 shadow-sm' : 'text-slate-400'" class="py-2 text-[10px] font-black uppercase tracking-widest rounded-lg transition-all">Take Away</button>
                    </div>

                    <div v-if="tipeOrder === 'DINE_IN'" class="animate-[fadeIn_0.2s_ease-out]">
                        <label class="block text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1.5">Nomor Meja Pelanggan</label>
                        <input v-model="nomorMeja" type="text" placeholder="Contoh: Meja 05, Meja VIP" class="w-full px-4 py-2.5 bg-slate-50 border-2 border-slate-200 rounded-xl outline-none font-bold text-xs text-slate-800 focus:bg-white focus:border-indigo-500 transition-all">
                    </div>
                </div>

                <div class="flex-1 overflow-y-auto custom-scrollbar p-5 space-y-4 bg-slate-50/50">
                    <div v-if="cart.length === 0" class="h-full flex flex-col items-center justify-center text-slate-300 py-20 text-center">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-12 h-12 mb-2 opacity-40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><path d="M8 14s1.5 2 4 2 4-2 4-2"/><line x1="9" y1="9" x2="9.01" y2="9"/><line x1="15" y1="9" x2="15.01" y2="9"/></svg>
                        <p class="text-[10px] font-black uppercase tracking-widest">Keranjang Kosong</p>
                    </div>
                    
                    <div v-for="item in cart" :key="item.id" class="bg-white p-4 rounded-xl border border-slate-200 shadow-sm space-y-3">
                        <div class="flex justify-between items-start gap-2">
                            <div>
                                <h4 class="font-black text-xs text-slate-800 uppercase leading-tight">{{ item.nama }}</h4>
                                <p class="text-[11px] font-bold text-indigo-600 mt-0.5">{{ formatRupiah(item.harga) }}</p>
                            </div>
                            <div class="flex items-center bg-slate-100 rounded-lg p-1 shrink-0">
                                <button @click="decreaseQty(item)" class="w-6 h-6 rounded bg-white font-black text-xs flex items-center justify-center shadow-sm text-slate-600 hover:bg-rose-50 hover:text-rose-600 transition-colors">-</button>
                                <span class="px-2.5 font-black text-xs text-slate-800 min-w-[20px] text-center">{{ item.qty }}</span>
                                <button @click="addToCart(item)" class="w-6 h-6 rounded bg-white font-black text-xs flex items-center justify-center shadow-sm text-slate-600 hover:bg-indigo-50 hover:text-indigo-600 transition-colors">+</button>
                            </div>
                        </div>

                        <div class="relative">
                            <input v-model="item.notes" type="text" placeholder="Tambahkan catatan koki (misal: pedas/es dikit)..." class="w-full pl-7 pr-3 py-1.5 bg-slate-50 hover:bg-slate-100 rounded-lg text-[10px] font-bold text-slate-600 outline-none focus:bg-white border border-transparent focus:border-slate-200 transition-all">
                            <div class="absolute inset-y-0 left-0 pl-2.5 flex items-center pointer-events-none"><svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3 text-slate-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 1 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg></div>
                        </div>
                    </div>
                </div>

                <div class="p-5 border-t border-slate-100 space-y-4 shrink-0 bg-white shadow-[0_-8px_24px_rgba(0,0,0,0.02)]">
                    <div class="flex justify-between items-baseline">
                        <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Total Bayar</span>
                        <span class="text-xl font-black text-slate-800 tracking-tight">{{ formatRupiah(totalBelanja) }}</span>
                    </div>
                    <button @click="checkoutOrder" class="w-full bg-slate-900 hover:bg-indigo-600 text-white py-4 rounded-xl font-black text-xs uppercase tracking-[0.2em] shadow-xl shadow-slate-200 transition-all active:scale-95 flex items-center justify-center gap-2">
                        Kirim Order Ke Dapur ➔
                    </button>
                </div>
            </div>

            <div v-if="cart.length > 0" class="xl:hidden fixed bottom-6 left-1/2 -translate-x-1/2 z-40 w-[90%] max-w-sm">
                <button @click="isMobileCartOpen = true" class="w-full bg-slate-900 text-white p-4 rounded-xl font-black text-xs uppercase tracking-widest shadow-2xl flex justify-between items-center animate-[fadeIn_0.2s_ease-out]">
                    <div class="flex items-center gap-2">
                        <span class="bg-indigo-600 text-white w-5 h-5 rounded-md flex items-center justify-center text-[10px] font-black">{{ cart.length }}</span>
                        <span>Lihat Keranjang</span>
                    </div>
                    <span>{{ formatRupiah(totalBelanja) }}</span>
                </button>
            </div>
            
        </div>
    </SidebarFnB>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { height: 6px; width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
.hide-scrollbar::-webkit-scrollbar { display: none; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
</style>