<script setup>
import { ref, computed, watch, onMounted } from 'vue';
import SidebarFnB from '../components/SidebarFnB.vue';
import api from '../../../api.js';
import Swal from 'sweetalert2';

// --- STATE UTAMA ---
const activeCategory = ref('semua');
const searchQuery = ref('');
const tipeOrder = ref('DINE_IN'); 
const nomorMeja = ref('');
const namaPemesan = ref('');
const isMobileCartOpen = ref(false);
const storeName = ref('Kasir POS'); 

// --- STATE PEMBAYARAN & STRUK ---
const metodeBayar = ref('CASH'); 
const uangBayarRaw = ref(0);      // Angka murni untuk kalkulasi
const uangBayarDisplay = ref(''); // String tampilan format ribuan
const isReceiptModalOpen = ref(false);
const lastOrderData = ref(null); 

// --- DRAG SCROLL CATEGORY (DESKTOP ONLY) ---
const scrollContainer = ref(null);

let isDragging = false;
let startX = 0;
let scrollStart = 0;
let movedDistance = 0;

const onMouseDown = (e) => {
    // hanya desktop
    if (window.innerWidth < 1024) return;

    isDragging = true;
    movedDistance = 0;

    startX = e.pageX;
    scrollStart = scrollContainer.value.scrollLeft;

    scrollContainer.value.classList.add('cursor-grabbing');
};

const onMouseMove = (e) => {
    if (!isDragging) return;

    e.preventDefault();

    const walk = e.pageX - startX;

    movedDistance = Math.abs(walk);

    scrollContainer.value.scrollLeft = scrollStart - walk;
};

const stopDragging = () => {
    isDragging = false;

    if (scrollContainer.value) {
        scrollContainer.value.classList.remove('cursor-grabbing');
    }
};

const onMouseUp = () => {
    stopDragging();
};

const onMouseLeave = () => {
    stopDragging();
};
// --- TOUCH SUPPORT ---
const onTouchStart = (e) => {
    const touch = e.touches[0];

    isDragging = true;
    movedDistance = 0;

    startX = touch.pageX;
    scrollStart = scrollContainer.value.scrollLeft;
};

const onTouchMove = (e) => {
    if (!isDragging) return;

    const touch = e.touches[0];

    const walk = touch.pageX - startX;

    movedDistance = Math.abs(walk);

    scrollContainer.value.scrollLeft = scrollStart - walk;
};

const onTouchEnd = () => {
    stopDragging();
};

watch(isMobileCartOpen, (newVal) => {
    if (typeof window !== 'undefined') document.body.style.overflow = newVal ? 'hidden' : 'auto';
});

// 🚀 FORMAT UANG REALTIME
const handleUangInput = (e) => {
    let val = e.target.value.replace(/\D/g, '');
    uangBayarRaw.value = Number(val);
    uangBayarDisplay.value = val ? new Intl.NumberFormat('id-ID').format(val) : '';
};

// --- DATA MASTER ---
const categories = ref([
    { id: 'semua', nama: 'Semua', icon: 'M4 6h16M4 12h16M4 18h16' },
    { id: 'paket', nama: 'Paket Combo', icon: 'M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z' },
    { id: 'makanan', nama: 'Makanan', icon: 'M12 2v20M17 5H9.5a4.5 4.5 0 0 0 0 9H15a4.5 4.5 0 0 1 0 9H6.5' },
    { id: 'minuman', nama: 'Minuman', icon: 'M17 10v12H7V10l5-5 5 5z M12 5v5' },
    { id: 'snack', nama: 'Snack', icon: 'M12 21a9 9 0 1 0 0-18 9 9 0 0 0 0 18Z M12 11a2 2 0 1 0 0-4 2 2 0 0 0 0 4Z M8 16a2 2 0 1 0 0-4 2 2 0 0 0 0 4Z M16 16a2 2 0 1 0 0-4 2 2 0 0 0 0 4Z' },
    { id: 'dessert', nama: 'Dessert', icon: 'M6 8a6 6 0 0 1 12 0c0 7-3 9-3 9H9s-3-2-3-9' }
]);

const products = ref([]);
const cart = ref([]);

const fetchMenuProducts = async () => {
    try {
        const response = await api.get('/fnb/products');
        products.value = response.data.map(p => ({
            ...p,
            nama: p.nama_produk, 
            harga: p.harga_jual,
            kategori: p.kategori,
            gambar: p.gambar || 'https://images.unsplash.com/photo-1546069901-ba9599a7e63c?auto=format&fit=crop&w=400&q=80'
        }));
    } catch (error) { console.error("Gagal menarik data menu dari server", error); }
};

onMounted(() => {
    fetchMenuProducts();
    const savedStore = localStorage.getItem('storeName');
    if (savedStore) storeName.value = savedStore;
});

const toggleMenuStatus = async (product, event) => {
    event.stopPropagation(); 
    try {
        const response = await api.put(`/fnb/products/${product.id}/toggle`);
        Swal.fire({ toast: true, position: 'top-end', icon: response.data.is_available ? 'success' : 'warning', title: response.data.message, showConfirmButton: false, timer: 2000, customClass: { popup: 'rounded-xl font-sans' } });
        fetchMenuProducts();
    } catch (error) { Swal.fire('Gagal!', 'Tidak dapat mengubah status menu.', 'error'); }
};

const formatRupiah = (angka) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(angka || 0);

const filteredProducts = computed(() => {
    return products.value.filter(p => {
        const matchCategory = activeCategory.value === 'semua' || p.kategori === activeCategory.value;
        const matchSearch = p.nama.toLowerCase().includes(searchQuery.value.toLowerCase());
        return matchCategory && matchSearch;
    });
});

// --- LOGIKA KERANJANG & KEMBALIAN ---
const totalBelanja = computed(() => cart.value.reduce((sum, item) => sum + (item.harga * item.qty), 0));
const kembalian = computed(() => Math.max(0, uangBayarRaw.value - totalBelanja.value));
const getQtyInCart = (productId) => { const item = cart.value.find(i => i.id === productId); return item ? item.qty : 0; };

const addToCart = (product) => {
    const existingItem = cart.value.find(item => item.id === product.id);
    if (existingItem) existingItem.qty++;
    else cart.value.push({ ...product, qty: 1, notes: '' });
};

const decreaseQty = (product) => {
    const existingItem = cart.value.find(item => item.id === product.id);
    if (existingItem) {
        if (existingItem.qty > 1) existingItem.qty--;
        else {
            cart.value = cart.value.filter(item => item.id !== product.id);
            if (cart.value.length === 0) isMobileCartOpen.value = false;
        }
    }
};

// 🚀 FUNGSI CHECKOUT
const checkoutOrder = async () => {
    if (cart.value.length === 0) return Swal.fire('Keranjang Kosong', 'Pilih menu dulu dong beb!', 'warning');
    if (tipeOrder.value === 'DINE_IN' && !nomorMeja.value) return Swal.fire('Nomor Meja?', 'Dine In wajib isi nomor meja!', 'warning');
    if (tipeOrder.value === 'TAKE_AWAY' && !namaPemesan.value) return Swal.fire('Nama Pemesan?', 'Wajib isi nama pemesan Take Away!', 'warning');
    if (metodeBayar.value === 'CASH' && uangBayarRaw.value < totalBelanja.value) return Swal.fire('Uang Kurang!', 'Nominal uang pembayaran tidak cukup!', 'error');

    try {
        const activeSessionId = localStorage.getItem('session_id') || 1; 
        const payload = {
            session_id: Number(activeSessionId), 
            tipe_order: tipeOrder.value,
            nomor_meja: tipeOrder.value === 'DINE_IN' ? nomorMeja.value : '-',
            nama_pemesan: tipeOrder.value === 'TAKE_AWAY' ? namaPemesan.value : '-', 
            metode_bayar: metodeBayar.value, 
            uang_diterima: metodeBayar.value === 'CASH' ? uangBayarRaw.value : totalBelanja.value,
            kembalian: metodeBayar.value === 'CASH' ? kembalian.value : 0,
            sub_total: totalBelanja.value, 
            pajak: 0, 
            total_harga: totalBelanja.value,
            items: cart.value.map(item => ({
                product_id: item.id, qty: item.qty, sub_total: item.harga * item.qty, notes: item.notes
            }))
        };

        const response = await api.post('/fnb/order', payload); 

        lastOrderData.value = {
            ...payload,
            invoice: response.data?.invoice || `INV-${Math.floor(Date.now() / 1000)}`,
            waktu: new Date().toLocaleString('id-ID', { dateStyle: 'short', timeStyle: 'short' }),
            kasir: localStorage.getItem('name') || 'Kasir',
            items: cart.value.map(item => ({...item, total_harga: item.harga * item.qty })) 
        };

        isReceiptModalOpen.value = true;
        cart.value = []; nomorMeja.value = ''; namaPemesan.value = ''; uangBayarRaw.value = 0; uangBayarDisplay.value = ''; isMobileCartOpen.value = false;

    } catch (error) { Swal.fire('Gagal!', 'Terjadi kendala saat memproses pesanan.', 'error'); }
};

const printReceipt = () => { window.print(); };
</script>

<template>
    <SidebarFnB>
        <div class="flex-1 flex h-full lg:h-screen bg-slate-50 overflow-hidden font-sans relative">
            
            <div class="flex-1 flex flex-col h-full overflow-hidden relative z-10 print:hidden">
                <div class="bg-white border-b border-slate-200 shrink-0 z-20 shadow-sm">
                    <div class="p-4 lg:p-6 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                        <div class="flex items-center gap-3 w-full sm:w-auto">
                            <div class="w-10 h-10 lg:w-12 lg:h-12 bg-gradient-to-br from-indigo-500 to-indigo-600 rounded-xl lg:rounded-2xl flex items-center justify-center shrink-0 shadow-lg shadow-indigo-200">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 lg:w-6 lg:h-6 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M3 2v7c0 1.1.9 2 2 2h4a2 2 0 0 0 2-2V2"/><path d="M7 2v20"/><path d="M21 15V2v0a5 5 0 0 0-5 5v6c0 1.1.9 2 2 2h3Zm0 0v7"/></svg>
                            </div>
                            <div class="overflow-hidden flex-1">
                                <h1 class="text-lg lg:text-2xl font-black tracking-tighter uppercase text-slate-800 leading-none truncate">{{ storeName }}</h1>
                                <p class="text-[9px] lg:text-[10px] font-black text-slate-400 uppercase tracking-widest mt-1">Sistem Kasir Utama</p>
                            </div>
                        </div>
                        
                        <div class="relative w-full sm:w-[280px] lg:w-[320px] group">
                            <input v-model="searchQuery" type="text" placeholder="Cari Nasi Goreng, Es Teh..." class="w-full pl-11 pr-4 py-3 bg-slate-100/50 border-2 border-slate-100 rounded-xl focus:bg-white focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 font-bold text-xs lg:text-sm text-slate-800 outline-none transition-all placeholder:text-slate-400 shadow-sm">
                            <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-slate-400 group-focus-within:text-indigo-500 transition-colors">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 lg:w-5 lg:h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
                            </div>
                        </div>
                    </div>

                    <div ref="scrollContainer" 
                         @mousedown="onMouseDown" @mouseleave="onMouseLeave" @mouseup="onMouseUp" @mousemove="onMouseMove"
                         class="px-4 lg:px-6 pb-4 flex gap-2 lg:gap-3 overflow-x-auto hide-scrollbar cursor-grab select-none">
                        <button 
                            v-for="cat in categories" 
                            :key="cat.id" 
                            @click="movedDistance < 8 ? activeCategory = cat.id : null" 
                            class="whitespace-nowrap px-4 py-2 lg:py-2.5 rounded-xl text-[10px] lg:text-xs font-black uppercase tracking-widest transition-all flex items-center gap-1.5 lg:gap-2 border-2 active:scale-95 shrink-0 pointer-events-auto"
                            :class="activeCategory === cat.id ? 'bg-indigo-600 text-white border-indigo-600 shadow-lg shadow-indigo-200' : 'bg-white text-slate-500 border-slate-100 hover:border-slate-300 hover:bg-slate-50'">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 lg:w-4 lg:h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path :d="cat.icon"/></svg>
                            {{ cat.nama }}
                        </button>
                    </div>
                </div>

                <div class="flex-1 min-h-0 overflow-y-auto custom-scrollbar p-4 lg:p-6 pb-32 xl:pb-6 relative">
                    <div v-if="filteredProducts.length === 0" class="h-full flex flex-col items-center justify-center text-slate-400">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-16 h-16 mb-4 opacity-20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><circle cx="12" cy="12" r="10"/><path d="m21.5 21.5-2.1-2.1"/><path d="M12 8v4l3 3"/></svg>
                        <span class="font-black text-xs uppercase tracking-widest">Menu tidak ditemukan</span>
                    </div>
                    
                    <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 2xl:grid-cols-5 gap-3 lg:gap-4">
                        <div v-for="prod in filteredProducts" :key="prod.id" class="bg-white border-2 border-slate-100/80 rounded-2xl lg:rounded-3xl overflow-hidden shadow-sm hover:shadow-xl transition-all flex flex-col group relative" :class="prod.is_available ? 'hover:border-indigo-300 hover:-translate-y-1' : 'opacity-75 grayscale-[30%]' ">
                            
                            <button @click="(e) => toggleMenuStatus(prod, e)" class="absolute top-2 left-2 z-20 px-2 py-1.5 rounded-lg text-[9px] font-black uppercase tracking-widest shadow-md transition-all active:scale-95 border border-white/20 backdrop-blur-md" :class="prod.is_available ? 'bg-emerald-500/90 text-white hover:bg-emerald-600' : 'bg-rose-500/90 text-white hover:bg-rose-600'">
                                {{ prod.is_available ? 'Tersedia' : 'Habis' }}
                            </button>

                            <div v-if="getQtyInCart(prod.id) > 0" class="absolute top-2 right-2 z-10 w-8 h-8 bg-indigo-600 text-white rounded-full flex items-center justify-center font-black text-sm border-2 border-white shadow-lg animate-bounce-short">
                                {{ getQtyInCart(prod.id) }}
                            </div>

                            <div class="w-full aspect-[4/3] bg-slate-100 relative overflow-hidden shrink-0 cursor-pointer" @click="prod.is_available ? addToCart(prod) : null">
                                <img :src="prod.gambar" :alt="prod.nama" class="w-full h-full object-cover transition-transform duration-700" :class="prod.is_available ? 'group-hover:scale-110' : ''">
                                <div class="absolute inset-0 bg-gradient-to-t from-slate-900/80 via-slate-900/10 to-transparent"></div>
                                <div v-if="!prod.is_available" class="absolute inset-0 flex items-center justify-center backdrop-blur-[2px]">
                                    <span class="bg-rose-600 text-white px-4 py-2 rounded-xl font-black text-sm uppercase tracking-widest -rotate-12 border-2 border-white shadow-xl">Habis Terjual</span>
                                </div>
                                <div class="absolute bottom-3 left-3 right-3">
                                    <span class="font-black text-sm lg:text-lg text-white tracking-tight drop-shadow-md">{{ formatRupiah(prod.harga) }}</span>
                                </div>
                            </div>

                            <div class="p-3 lg:p-4 flex-1 flex flex-col justify-between gap-3 bg-white">
                                <h3 class="font-black text-xs lg:text-sm text-slate-800 uppercase leading-snug line-clamp-2" :class="prod.is_available ? 'group-hover:text-indigo-600' : ''">{{ prod.nama }}</h3>
                                <button v-if="prod.is_available" @click="addToCart(prod)" class="w-full bg-indigo-50 hover:bg-indigo-600 text-indigo-600 hover:text-white border border-indigo-100 py-2.5 rounded-xl font-black text-[10px] lg:text-xs uppercase tracking-widest transition-colors active:scale-95 flex justify-center items-center gap-1.5">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg> Tambah
                                </button>
                                <button v-else disabled class="w-full bg-slate-100 text-slate-400 border border-slate-200 py-2.5 rounded-xl font-black text-[10px] lg:text-xs uppercase tracking-widest cursor-not-allowed">
                                    Kosong
                                </button>
                            </div>
                        </div>
                    </div>
                </div>

                <div v-if="cart.length > 0" class="xl:hidden absolute bottom-4 left-4 right-4 z-30 pointer-events-none print:hidden">
                    <button @click="isMobileCartOpen = true" class="w-full bg-slate-900 text-white p-4 lg:p-5 rounded-2xl font-black text-xs lg:text-sm uppercase tracking-widest shadow-[0_10px_30px_rgba(15,23,42,0.6)] flex justify-between items-center pointer-events-auto active:scale-95 transition-transform border border-slate-700 animate-slide-up">
                        <div class="flex items-center gap-3">
                            <span class="bg-indigo-500 text-white w-8 h-8 rounded-lg flex items-center justify-center text-sm font-black shadow-inner border border-indigo-400">{{ cart.length }}</span>
                            <span class="mt-0.5">Lihat Tagihan</span>
                        </div>
                        <span class="text-base mt-0.5">{{ formatRupiah(totalBelanja) }}</span>
                    </button>
                </div>
            </div>

            <aside class="hidden xl:flex w-[360px] 2xl:w-[420px] bg-white border-l border-slate-200 flex-col h-full z-20 shadow-2xl relative print:hidden">
                
                <div class="p-5 border-b border-slate-100 bg-white shrink-0 flex items-center gap-3">
                    <div class="w-8 h-8 rounded-lg bg-indigo-50 text-indigo-600 flex items-center justify-center shadow-inner">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M6 2 3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z"/><line x1="3" y1="6" x2="21" y2="6"/><path d="M16 10a4 4 0 0 1-8 0"/></svg>
                    </div>
                    <h2 class="font-black text-xs uppercase tracking-[0.2em] text-slate-800 mt-0.5">Rincian Pesanan</h2>
                </div>

                <div class="flex-1 min-h-0 overflow-y-auto custom-scrollbar p-5 space-y-5 bg-slate-50/50">
                    
                    <div class="bg-white p-3 rounded-2xl border border-slate-100 shadow-sm space-y-3">
                        <div class="flex p-1 bg-slate-100 rounded-xl relative">
                            <div class="absolute inset-y-1 w-[calc(50%-4px)] bg-white rounded-lg shadow-sm transition-all duration-300" :class="tipeOrder === 'TAKE_AWAY' ? 'left-[calc(50%+2px)]' : 'left-1'"></div>
                            <button @click="tipeOrder = 'DINE_IN'" class="flex-1 py-2.5 text-[10px] font-black uppercase tracking-widest z-10 transition-colors rounded-lg" :class="tipeOrder === 'DINE_IN' ? 'text-indigo-700' : 'text-slate-500'">Makan Sini</button>
                            <button @click="tipeOrder = 'TAKE_AWAY'" class="flex-1 py-2.5 text-[10px] font-black uppercase tracking-widest z-10 transition-colors rounded-lg" :class="tipeOrder === 'TAKE_AWAY' ? 'text-indigo-700' : 'text-slate-500'">Bungkus</button>
                        </div>
                        <div v-if="tipeOrder === 'DINE_IN'" class="animate-fade-in relative">
                            <input v-model="nomorMeja" type="text" placeholder="Nomor Meja (Cth: 05)" class="w-full pl-11 pr-4 py-3 bg-slate-50 border border-slate-100 rounded-xl outline-none font-black text-xs text-slate-800 focus:border-indigo-500 focus:bg-white transition-all">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-400 absolute left-4 top-1/2 -translate-y-1/2" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><rect width="18" height="18" x="3" y="3" rx="2"/><path d="M3 9h18"/><path d="M9 21v-9"/><path d="M15 21v-9"/></svg>
                        </div>
                        <div v-if="tipeOrder === 'TAKE_AWAY'" class="animate-fade-in relative">
                            <input v-model="namaPemesan" type="text" placeholder="Nama Pemesan (Cth: Kak Budi)" class="w-full pl-11 pr-4 py-3 bg-indigo-50/50 border border-indigo-100 rounded-xl outline-none font-black text-xs text-indigo-900 focus:bg-white focus:border-indigo-400 transition-all placeholder:text-indigo-300">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-indigo-400 absolute left-4 top-1/2 -translate-y-1/2" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
                        </div>
                    </div>

                    <div v-if="cart.length === 0" class="flex flex-col items-center justify-center py-10 text-slate-300">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-16 h-16 mb-4 opacity-30" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><circle cx="12" cy="12" r="10"/><path d="M8 14s1.5 2 4 2 4-2 4-2"/><line x1="9" y1="9" x2="9.01" y2="9"/><line x1="15" y1="9" x2="15.01" y2="9"/></svg>
                        <p class="text-[10px] font-black uppercase tracking-widest">Keranjang Kosong</p>
                    </div>

                    <div class="space-y-3">
                        <div v-for="item in cart" :key="'desk-'+item.id" class="bg-white p-4 rounded-2xl border border-slate-100 shadow-sm relative group overflow-hidden">
                            <div class="absolute left-0 top-0 bottom-0 w-1 bg-indigo-500 opacity-0 group-hover:opacity-100 transition-opacity"></div>
                            
                            <div class="flex justify-between items-start gap-2 pl-1">
                                <div class="flex-1">
                                    <h4 class="font-black text-xs text-slate-800 uppercase leading-snug">{{ item.nama }}</h4>
                                    <p class="text-[11px] font-bold text-indigo-600 mt-1">{{ formatRupiah(item.harga) }}</p>
                                </div>
                                <div class="flex items-center bg-slate-50 rounded-lg p-1 shrink-0 border border-slate-100">
                                    <button @click="decreaseQty(item)" class="w-7 h-7 rounded bg-white font-black text-sm flex items-center justify-center shadow-sm text-slate-600 active:bg-rose-100 transition-colors">-</button>
                                    <span class="px-3 font-black text-xs text-slate-800 min-w-[28px] text-center">{{ item.qty }}</span>
                                    <button @click="addToCart(item)" class="w-7 h-7 rounded bg-white font-black text-sm flex items-center justify-center shadow-sm text-slate-600 active:bg-indigo-100 transition-colors">+</button>
                                </div>
                            </div>
                            
                            <div class="mt-3 relative pl-1">
                                <div class="absolute inset-y-0 left-3 flex items-center pointer-events-none">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 text-slate-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M17 3a2.828 2.828 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5L17 3z"/></svg>
                                </div>
                                <input v-model="item.notes" type="text" placeholder="Catatan opsional (pedas, dll)..." class="w-full pl-8 pr-3 py-2 bg-slate-50 hover:bg-slate-100 focus:bg-white rounded-lg text-[10px] font-bold text-slate-700 outline-none border border-transparent focus:border-indigo-300 transition-all placeholder:text-slate-400">
                            </div>
                        </div>
                    </div>
                </div>

                <div class="p-5 border-t border-slate-200 bg-white shrink-0 shadow-[0_-10px_20px_rgba(0,0,0,0.02)]">
                    
                    <div class="flex gap-2 mb-4">
                        <button @click="metodeBayar = 'CASH'" class="flex-1 py-3 rounded-xl border-2 font-black text-[10px] uppercase tracking-widest transition-all flex items-center justify-center gap-2" :class="metodeBayar === 'CASH' ? 'border-indigo-600 bg-indigo-50 text-indigo-700' : 'border-slate-100 text-slate-400 hover:border-slate-300'">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><rect width="20" height="12" x="2" y="6" rx="2"/><circle cx="12" cy="12" r="2"/><path d="M6 12h.01M18 12h.01"/></svg>
                            TUNAI
                        </button>
                        <button @click="metodeBayar = 'QRIS'" class="flex-1 py-3 rounded-xl border-2 font-black text-[10px] uppercase tracking-widest transition-all flex items-center justify-center gap-2" :class="metodeBayar === 'QRIS' ? 'border-sky-500 bg-sky-50 text-sky-600' : 'border-slate-100 text-slate-400 hover:border-slate-300'">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><rect width="5" height="5" x="3" y="3" rx="1"/><rect width="5" height="5" x="16" y="3" rx="1"/><rect width="5" height="5" x="3" y="16" rx="1"/><path d="M21 16h-3a2 2 0 0 0-2 2v3"/><path d="M21 21v.01"/><path d="M12 7v3a2 2 0 0 1-2 2H7"/><path d="M3 12h.01"/><path d="M12 3h.01"/><path d="M12 16v.01"/><path d="M16 12h1"/><path d="M21 12v.01"/><path d="M12 21v-1"/></svg>
                            QRIS
                        </button>
                    </div>

                    <div v-if="metodeBayar === 'CASH'" class="mb-4 flex flex-col gap-2 animate-fade-in bg-slate-50 p-3.5 rounded-2xl border border-slate-100">
                        <label class="block text-[9px] font-black text-slate-400 uppercase tracking-widest ml-1">Uang Diterima</label>
                        <div class="relative">
                            <input :value="uangBayarDisplay" @input="handleUangInput" type="text" placeholder="Ketik Nominal..." class="w-full pl-10 pr-4 py-2.5 bg-white border border-slate-200 rounded-xl outline-none font-black text-sm text-slate-800 focus:border-indigo-500 transition-all">
                            <span class="absolute left-3.5 top-1/2 -translate-y-1/2 font-black text-slate-400 text-sm">Rp</span>
                        </div>
                        <div class="flex justify-between items-center pt-2 border-t border-slate-200 border-dashed mt-1">
                            <span class="text-[9px] font-black text-slate-400 uppercase tracking-widest">Kembalian</span>
                            <span class="text-xs font-black" :class="kembalian > 0 ? 'text-emerald-500' : 'text-slate-400'">{{ formatRupiah(kembalian) }}</span>
                        </div>
                    </div>

                    <div class="flex justify-between items-end mb-4 px-1">
                        <div>
                            <span class="block text-[9px] font-black text-slate-400 uppercase tracking-widest mb-0.5">Total Tagihan</span>
                            <span class="text-[10px] font-bold text-slate-400">{{ cart.length }} Item Terpilih</span>
                        </div>
                        <span class="text-2xl font-black text-indigo-600 tracking-tighter">{{ formatRupiah(totalBelanja) }}</span>
                    </div>
                    <button @click="checkoutOrder" class="w-full bg-slate-900 hover:bg-indigo-600 text-white py-4 rounded-xl font-black text-xs uppercase tracking-[0.2em] shadow-xl shadow-slate-200 transition-all active:scale-95 flex items-center justify-center gap-2">
                        Bayar & Cetak ➔
                    </button>
                </div>
            </aside>
            
        </div>

        <div v-if="isMobileCartOpen" @click="isMobileCartOpen = false" class="xl:hidden fixed inset-0 bg-slate-900/60 backdrop-blur-sm z-[60] transition-opacity animate-fade-in print:hidden"></div>
        
        <div class="xl:hidden fixed inset-x-0 bottom-0 z-[70] flex flex-col bg-slate-50 rounded-t-[32px] shadow-[0_-20px_40px_rgba(0,0,0,0.3)] transition-transform duration-300 ease-out print:hidden h-[85vh]"
             :class="isMobileCartOpen ? 'translate-y-0' : 'translate-y-full'">
            
            <div class="p-4 shrink-0 bg-white rounded-t-[32px] border-b border-slate-100 flex flex-col items-center">
                <div @click="isMobileCartOpen = false" class="w-12 h-1.5 bg-slate-300 rounded-full mb-3 cursor-pointer"></div>
                <div class="w-full flex justify-between items-center px-1">
                    <h2 class="font-black text-sm uppercase tracking-widest text-slate-800">Detail Keranjang</h2>
                    <button @click="isMobileCartOpen = false" class="w-8 h-8 bg-slate-100 text-slate-500 rounded-full flex items-center justify-center font-black active:scale-90 transition-transform">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                    </button>
                </div>
            </div>

            <div class="flex-1 min-h-0 overflow-y-auto custom-scrollbar p-4 space-y-4">
                
                <div class="bg-white p-3 rounded-2xl border border-slate-100 shadow-sm space-y-3">
                    <div class="flex p-1 bg-slate-100 rounded-xl relative">
                        <div class="absolute inset-y-1 w-[calc(50%-4px)] bg-white rounded-lg shadow-sm transition-all duration-300" :class="tipeOrder === 'TAKE_AWAY' ? 'left-[calc(50%+2px)]' : 'left-1'"></div>
                        <button @click="tipeOrder = 'DINE_IN'" class="flex-1 py-2.5 text-[10px] font-black uppercase tracking-widest z-10 transition-colors rounded-lg" :class="tipeOrder === 'DINE_IN' ? 'text-indigo-700' : 'text-slate-500'">Makan Sini</button>
                        <button @click="tipeOrder = 'TAKE_AWAY'" class="flex-1 py-2.5 text-[10px] font-black uppercase tracking-widest z-10 transition-colors rounded-lg" :class="tipeOrder === 'TAKE_AWAY' ? 'text-indigo-700' : 'text-slate-500'">Bungkus</button>
                    </div>

                    <div v-if="tipeOrder === 'DINE_IN'" class="animate-fade-in relative">
                        <input v-model="nomorMeja" type="text" placeholder="Nomor Meja (Cth: 05)" class="w-full pl-11 pr-4 py-3 bg-slate-50 border border-slate-100 rounded-xl outline-none font-black text-xs text-slate-800 focus:border-indigo-500 focus:bg-white transition-all">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-400 absolute left-4 top-1/2 -translate-y-1/2" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><rect width="18" height="18" x="3" y="3" rx="2"/><path d="M3 9h18"/><path d="M9 21v-9"/><path d="M15 21v-9"/></svg>
                    </div>
                    <div v-if="tipeOrder === 'TAKE_AWAY'" class="animate-fade-in relative">
                        <input v-model="namaPemesan" type="text" placeholder="Nama Pemesan (Cth: Kak Budi)" class="w-full pl-11 pr-4 py-3 bg-indigo-50/50 border border-indigo-100 rounded-xl outline-none font-black text-xs text-indigo-900 focus:bg-white focus:border-indigo-400 transition-all placeholder:text-indigo-300">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-indigo-400 absolute left-4 top-1/2 -translate-y-1/2" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
                    </div>
                </div>

                <div class="space-y-3">
                    <div v-for="item in cart" :key="'mob-'+item.id" class="bg-white p-4 rounded-2xl border border-slate-100 shadow-sm overflow-hidden relative">
                        <div class="flex justify-between items-start gap-2">
                            <div class="flex-1">
                                <h4 class="font-black text-[11px] sm:text-xs text-slate-800 uppercase leading-snug">{{ item.nama }}</h4>
                                <p class="text-[10px] sm:text-xs font-bold text-indigo-600 mt-1">{{ formatRupiah(item.harga) }}</p>
                            </div>
                            <div class="flex items-center bg-slate-50 rounded-lg p-1 shrink-0 border border-slate-100">
                                <button @click="decreaseQty(item)" class="w-6 h-6 sm:w-7 sm:h-7 rounded-md bg-white font-black text-sm flex items-center justify-center shadow-sm text-slate-600 active:bg-rose-100 transition-colors">-</button>
                                <span class="px-2 sm:px-3 font-black text-xs text-slate-800 min-w-[24px] text-center">{{ item.qty }}</span>
                                <button @click="addToCart(item)" class="w-6 h-6 sm:w-7 sm:h-7 rounded-md bg-white font-black text-sm flex items-center justify-center shadow-sm text-slate-600 active:bg-indigo-100 transition-colors">+</button>
                            </div>
                        </div>
                        
                        <div class="mt-3 relative">
                            <div class="absolute inset-y-0 left-2.5 flex items-center pointer-events-none">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 text-slate-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M17 3a2.828 2.828 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5L17 3z"/></svg>
                            </div>
                            <input v-model="item.notes" type="text" placeholder="Catatan opsional..." class="w-full pl-8 pr-3 py-2 bg-slate-50 hover:bg-slate-100 focus:bg-white rounded-lg text-[10px] font-bold text-slate-700 outline-none border border-transparent focus:border-indigo-300 transition-all placeholder:text-slate-400">
                        </div>
                    </div>
                </div>
            </div>

            <div class="p-4 border-t border-slate-200 bg-white shrink-0 space-y-3 pb-safe shadow-[0_-10px_20px_rgba(0,0,0,0.03)]">
                <div class="flex gap-2">
                    <button @click="metodeBayar = 'CASH'" class="flex-1 py-3 rounded-xl border-2 font-black text-[10px] uppercase tracking-widest transition-all flex items-center justify-center gap-1.5" :class="metodeBayar === 'CASH' ? 'border-indigo-600 bg-indigo-50 text-indigo-700' : 'border-slate-100 text-slate-400'">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><rect width="20" height="12" x="2" y="6" rx="2"/><circle cx="12" cy="12" r="2"/><path d="M6 12h.01M18 12h.01"/></svg>
                        TUNAI
                    </button>
                    <button @click="metodeBayar = 'QRIS'" class="flex-1 py-3 rounded-xl border-2 font-black text-[10px] uppercase tracking-widest transition-all flex items-center justify-center gap-1.5" :class="metodeBayar === 'QRIS' ? 'border-sky-500 bg-sky-50 text-sky-600' : 'border-slate-100 text-slate-400'">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><rect width="5" height="5" x="3" y="3" rx="1"/><rect width="5" height="5" x="16" y="3" rx="1"/><rect width="5" height="5" x="3" y="16" rx="1"/><path d="M21 16h-3a2 2 0 0 0-2 2v3"/><path d="M21 21v.01"/><path d="M12 7v3a2 2 0 0 1-2 2H7"/><path d="M3 12h.01"/><path d="M12 3h.01"/><path d="M12 16v.01"/><path d="M16 12h1"/><path d="M21 12v.01"/><path d="M12 21v-1"/></svg>
                        QRIS
                    </button>
                </div>
                
                <div v-if="metodeBayar === 'CASH'" class="flex items-center gap-3 bg-slate-50 p-2.5 rounded-xl border border-slate-100 animate-fade-in">
                    <div class="w-[55%] relative">
                        <input :value="uangBayarDisplay" @input="handleUangInput" type="text" placeholder="Nominal..." class="w-full pl-8 pr-2 py-2 bg-white border border-slate-200 rounded-lg outline-none font-bold text-xs">
                        <span class="absolute left-2.5 top-1/2 -translate-y-1/2 font-black text-slate-400 text-[10px]">Rp</span>
                    </div>
                    <div class="w-[45%] text-right pr-2">
                        <p class="text-[8px] font-black text-slate-400 uppercase tracking-widest">Kembali</p>
                        <p class="text-xs font-black text-emerald-500">{{ formatRupiah(kembalian) }}</p>
                    </div>
                </div>

                <div class="flex justify-between items-center px-1 pt-1">
                    <div>
                        <span class="block text-[9px] font-black text-slate-400 uppercase tracking-widest mb-0.5">Total Tagihan</span>
                        <span class="text-xl font-black text-indigo-600 tracking-tighter leading-none">{{ formatRupiah(totalBelanja) }}</span>
                    </div>
                    <button @click="checkoutOrder" class="w-32 bg-slate-900 text-white py-3.5 rounded-xl font-black text-xs uppercase tracking-widest shadow-xl shadow-slate-900/20 active:scale-95 transition-transform flex items-center justify-center gap-1.5">
                        Bayar
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="5" y1="12" x2="19" y2="12"/><polyline points="12 5 19 12 12 19"/></svg>
                    </button>
                </div>
            </div>
        </div>

        <div v-if="isReceiptModalOpen" class="fixed inset-0 z-[100] flex items-center justify-center p-4">
            <div @click="isReceiptModalOpen = false" class="absolute inset-0 bg-slate-900/80 backdrop-blur-sm print:hidden animate-fade-in"></div>
            
            <div class="bg-white w-full max-w-sm rounded-[24px] shadow-2xl relative z-10 animate-slide-up flex flex-col overflow-hidden max-h-[90vh]">
                
                <div id="receipt-paper" class="p-6 md:p-8 bg-white text-slate-900 font-mono text-xs leading-relaxed flex-1 overflow-y-auto custom-scrollbar">
                    
                    <div class="text-center mb-6">
                        <h2 class="font-black text-lg md:text-xl uppercase tracking-widest mb-1">{{ storeName }}</h2>
                        <p class="text-[10px] text-slate-500">Sistem POS FnB Terintegrasi</p>
                        <p class="text-[10px] text-slate-500 mt-1">Terima Kasih Atas Kunjungan Anda!</p>
                        <div class="border-b-2 border-dashed border-slate-300 mt-4 mb-4"></div>
                    </div>

                    <div class="space-y-1.5 mb-4 text-[10px] sm:text-xs">
                        <div class="flex justify-between"><span>No. Order</span><span class="font-bold">{{ lastOrderData?.invoice }}</span></div>
                        <div class="flex justify-between"><span>Waktu</span><span>{{ lastOrderData?.waktu }}</span></div>
                        <div class="flex justify-between"><span>Kasir</span><span>{{ lastOrderData?.kasir }}</span></div>
                        <div class="flex justify-between items-start mt-2">
                            <span>Pesanan</span>
                            <div class="text-right">
                                <span class="font-bold border border-slate-800 px-1.5 rounded">{{ lastOrderData?.tipe_order === 'DINE_IN' ? 'DINE IN' : 'TAKE AWAY' }}</span>
                                <p class="font-black mt-1 uppercase">{{ lastOrderData?.tipe_order === 'DINE_IN' ? lastOrderData?.nomor_meja : lastOrderData?.nama_pemesan }}</p>
                            </div>
                        </div>
                        <div class="border-b-2 border-dashed border-slate-300 pt-3 mb-3"></div>
                    </div>

                    <div class="space-y-3 mb-4">
                        <div v-for="(item, idx) in lastOrderData?.items" :key="idx">
                            <div class="flex justify-between font-bold text-[10px] sm:text-xs">
                                <span class="uppercase w-2/3">{{ item.nama }}</span>
                                <span>{{ formatRupiah(item.total_harga) }}</span>
                            </div>
                            <div class="flex justify-between text-[10px] text-slate-500 mt-0.5">
                                <span>{{ item.qty }} x {{ formatRupiah(item.harga) }}</span>
                            </div>
                            <div v-if="item.notes" class="text-[9px] text-slate-400 mt-0.5 italic">* {{ item.notes }}</div>
                        </div>
                    </div>

                    <div class="border-b-2 border-dashed border-slate-300 mb-4"></div>

                    <div class="space-y-1.5 text-[10px] sm:text-xs">
                        <div class="flex justify-between font-black text-sm mb-2">
                            <span>TOTAL BAYAR</span>
                            <span>{{ formatRupiah(lastOrderData?.total_harga) }}</span>
                        </div>
                        <div class="flex justify-between text-slate-600">
                            <span>Pembayaran</span>
                            <span class="font-bold">{{ lastOrderData?.metode_bayar }}</span>
                        </div>
                        <div class="flex justify-between text-slate-600" v-if="lastOrderData?.metode_bayar === 'CASH'">
                            <span>Tunai Terima</span>
                            <span>{{ formatRupiah(lastOrderData?.uang_diterima) }}</span>
                        </div>
                        <div class="flex justify-between font-bold mt-1" v-if="lastOrderData?.metode_bayar === 'CASH'">
                            <span>KEMBALI</span>
                            <span>{{ formatRupiah(lastOrderData?.kembalian) }}</span>
                        </div>
                    </div>

                    <div class="mt-8 text-center text-[9px] text-slate-400 space-y-1">
                        <p>Barang yang sudah dibeli tidak dapat ditukar.</p>
                        <p>Powered by RestoPOS System</p>
                    </div>
                </div>

                <div class="p-4 border-t border-slate-100 bg-slate-50 flex gap-3 shrink-0 print:hidden">
                    <button @click="isReceiptModalOpen = false" class="flex-1 py-3.5 bg-white border border-slate-200 text-slate-600 rounded-xl font-black text-[10px] uppercase tracking-widest hover:bg-slate-100 active:scale-95 transition-all">Tutup</button>
                    <button @click="printReceipt" class="flex-1 py-3.5 bg-indigo-600 text-white rounded-xl font-black text-[10px] uppercase tracking-widest shadow-lg shadow-indigo-200 hover:bg-slate-900 active:scale-95 transition-all flex items-center justify-center gap-2">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="6 9 6 2 18 2 18 9"/><path d="M6 18H4a2 2 0 0 1-2-2v-5a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v5a2 2 0 0 1-2 2h-2"/><rect width="12" height="8" x="6" y="14"/></svg>
                        Print Struk
                    </button>
                </div>
            </div>
        </div>

    </SidebarFnB>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 5px; height: 5px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
.hide-scrollbar::-webkit-scrollbar { display: none; }

.animate-fade-in { animation: fadeIn 0.3s cubic-bezier(0.4, 0, 0.2, 1); }
.animate-slide-up { animation: slideUp 0.4s cubic-bezier(0.34, 1.56, 0.64, 1); }
.animate-bounce-short { animation: bounceShort 0.4s ease-out; }

@keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
@keyframes slideUp { from { opacity: 0; transform: translateY(30px); } to { opacity: 1; transform: translateY(0); } }
@keyframes bounceShort { 0%, 100% { transform: scale(1); } 50% { transform: scale(1.25); } }

.cursor-grab { cursor: grab; }
.cursor-grabbing { cursor: grabbing !important; }

/* CSS KHUSUS PRINT STRUK THERMAL */
@media print {
    body * { visibility: hidden; }
    #receipt-paper, #receipt-paper * { visibility: visible; }
    #receipt-paper {
        position: absolute;
        left: 0;
        top: 0;
        width: 100%;
        max-width: 58mm; /* Ukuran kertas struk thermal */
        margin: 0 auto;
        padding: 10px;
        box-shadow: none;
    }
}

/* Biar padding bawah aman di iPhone (Notch) */
.pb-safe { padding-bottom: env(safe-area-inset-bottom, 1rem); }
</style>