<script setup>
import { ref, onBeforeUnmount } from 'vue';
import Sidebar from '../../components/Sidebar.vue';
import api from '../../api.js';
import Swal from 'sweetalert2';
import { Html5Qrcode } from "html5-qrcode"; // 🚀 IMPORT SCANNER

const notes = ref('Stock Opname Reguler');
const searchQuery = ref('');
const products = ref([]);
const cartSO = ref([]);
const isSubmitting = ref(false);

// --- 🚀 LOGIKA KAMERA SCANNER BARCODE ---
const showScanner = ref(false);
let html5QrCode = null;

const startScanner = async () => {
    showScanner.value = true;
    setTimeout(async () => {
        try {
            html5QrCode = new Html5Qrcode("reader-so");
            await html5QrCode.start(
                { facingMode: "environment" }, 
                { fps: 10, qrbox: { width: 250, height: 100 } }, 
                (decodedText) => {
                    searchQuery.value = decodedText; 
                    stopScanner();
                    
                    // Langsung cari dan auto-add kalau dari scanner
                    searchProduct(true); 
                    
                    const audio = new Audio('https://www.soundjay.com/buttons/beep-07a.mp3');
                    audio.play().catch(()=>{}); 
                },
                (errorMessage) => {} 
            );
        } catch (err) {
            console.error(err);
            Swal.fire('Oops!', 'Gagal menyalakan kamera. Pastikan browser diizinkan mengakses kamera.', 'error');
            stopScanner();
        }
    }, 300);
};

const stopScanner = () => {
    if (html5QrCode) {
        html5QrCode.stop().then(() => {
            html5QrCode.clear();
            showScanner.value = false;
        }).catch(err => {
            showScanner.value = false;
        });
    } else {
        showScanner.value = false;
    }
};

onBeforeUnmount(() => {
    if (showScanner.value) stopScanner();
});
// ------------------------------------------

// --- SEARCH LOGIC DENGAN AUTO-ADD ---
const searchProduct = async (isFromScanner = false) => {
    if (searchQuery.value.length < 2) return;
    try {
        // Blind Counting: Hanya ambil nama dan ID, JANGAN tampilkan stok sistem
        const res = await api.get(`/products?search=${searchQuery.value}`);
        const foundData = res.data.data;

        if (isFromScanner) {
            if (foundData.length > 0) {
                addToCart(foundData[0]);
                Swal.fire({
                    toast: true,
                    position: 'top-end',
                    icon: 'success',
                    title: `${foundData[0].nama_produk} masuk hitungan!`,
                    showConfirmButton: false,
                    timer: 1500
                });
            } else {
                Swal.fire('Tidak Ditemukan!', 'Barcode ini tidak ada di Master Produk.', 'error');
                searchQuery.value = '';
            }
        } else {
            products.value = foundData;
        }
    } catch (err) { console.error(err); }
};

const addToCart = (product) => {
    const existing = cartSO.value.find(item => item.product_id === product.id);
    if (existing) {
        existing.actual_qty++; // Khas scanner: tiap scan tambah 1
    } else {
        // unshift biar barang yang baru di-scan selalu muncul paling atas!
        cartSO.value.unshift({
            product_id: product.id,
            nama_produk: product.nama_produk,
            actual_qty: 1
        });
    }
    searchQuery.value = '';
    products.value = [];
};

const removeItem = (index) => cartSO.value.splice(index, 1);

// --- SUBMIT STOCK OPNAME ---
const submitSO = async () => {
    if (cartSO.value.length === 0) {
        return Swal.fire('Kosong!', 'Input barang yang dihitung dulu Mas.', 'warning');
    }

    const result = await Swal.fire({
        title: 'Finalisasi Hitung Fisik?',
        text: "Pastikan jumlah fisik sudah benar. Stok sistem akan di-override secara permanen!",
        icon: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#2563eb',
        confirmButtonText: 'Ya, Finalisasi & Kunci!'
    });

    if (!result.isConfirmed) return;

    isSubmitting.value = true;
    try {
        await api.post('/stock-opname', {
            notes: notes.value,
            items: cartSO.value
        });

        Swal.fire({
            icon: 'success',
            title: 'Audit Selesai!',
            text: 'Stok master telah disesuaikan dengan fisik gudang.',
            timer: 2500,
            showConfirmButton: false
        });
        cartSO.value = [];
        notes.value = 'Stock Opname Reguler';
    } catch (err) {
        Swal.fire('Gagal', err.response?.data?.error || 'Terjadi kesalahan sistem', 'error');
    } finally {
        isSubmitting.value = false;
    }
};
</script>

<template>
    <Sidebar>
        <div class="p-4 md:p-8 lg:p-10 max-w-7xl mx-auto font-sans bg-[#f8fafc] min-h-screen">
            
            <div class="bg-gradient-to-br from-slate-900 via-slate-800 to-indigo-900 rounded-[40px] p-8 md:p-10 mb-8 text-white shadow-2xl flex flex-col md:flex-row items-center justify-between relative overflow-hidden border border-white/10">
                <div class="absolute -left-10 -bottom-10 opacity-10">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-64 h-64" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1"><path stroke-linecap="round" stroke-linejoin="round" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" /></svg>
                </div>
                
                <div class="z-10 text-center md:text-left">
                    <h1 class="text-3xl md:text-4xl font-black tracking-tighter mb-2 uppercase italic leading-none">Stock <span class="text-indigo-400">Opname</span></h1>
                    <p class="text-slate-300 font-bold text-[10px] uppercase tracking-[0.2em] flex items-center justify-center md:justify-start gap-2">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
                        Penyesuaian Stok Fisik vs Sistem
                    </p>
                </div>

                <div class="z-10 mt-6 md:mt-0 flex items-center gap-3 bg-amber-500/20 px-5 py-3 rounded-2xl border border-amber-500/30 backdrop-blur-md">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" /></svg>
                    <div class="flex flex-col">
                        <span class="text-[9px] font-black text-amber-200 uppercase tracking-widest">Metode Audit</span>
                        <span class="text-sm font-black text-amber-400 uppercase tracking-widest">BLIND COUNTING</span>
                    </div>
                </div>
            </div>

            <div class="grid grid-cols-1 xl:grid-cols-3 gap-8">
                <div class="xl:col-span-2 space-y-6">
                    
                    <div class="bg-white p-6 rounded-[32px] shadow-sm border border-slate-100 relative">
                        <div class="flex items-center gap-3">
                            <div class="relative flex-1 group">
                                <div class="absolute inset-y-0 left-0 pl-5 flex items-center pointer-events-none">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-400 group-focus-within:text-indigo-600 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
                                </div>
                                <input v-model="searchQuery" @input="searchProduct(false)" type="text" placeholder="Scan Barcode / Ketik Nama Barang..." class="w-full pl-14 pr-4 py-5 bg-slate-50 border-2 border-transparent rounded-[24px] outline-none font-bold text-sm focus:border-indigo-500 focus:bg-white transition-all text-slate-800 placeholder:text-slate-400">
                            </div>
                            
                            <button @click="startScanner" class="shrink-0 p-5 bg-indigo-50 hover:bg-indigo-600 text-indigo-600 hover:text-white rounded-[24px] transition-all border border-indigo-100 shadow-sm" title="Scan Barcode">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" /><path stroke-linecap="round" stroke-linejoin="round" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
                            </button>
                        </div>
                        
                        <div v-if="products.length && searchQuery.length > 0" class="absolute left-6 right-6 mt-3 bg-white border border-slate-100 rounded-[28px] shadow-2xl z-[100] overflow-hidden">
                            <div v-for="p in products" :key="p.id" @click="addToCart(p)" class="p-5 hover:bg-indigo-50 cursor-pointer border-b last:border-0 flex justify-between items-center group transition-all">
                                <div class="flex items-center gap-4">
                                    <div class="w-10 h-10 bg-slate-100 rounded-xl flex items-center justify-center text-xl shadow-inner group-hover:bg-indigo-100 transition-colors">📦</div>
                                    <div class="font-black text-slate-800 uppercase text-sm group-hover:text-indigo-600 transition-colors">{{ p.nama_produk }}</div>
                                </div>
                                <span class="text-[9px] bg-slate-900 text-white px-4 py-2 rounded-xl font-black uppercase tracking-widest shadow-sm">Pilih</span>
                            </div>
                        </div>
                    </div>

                    <div class="bg-white rounded-[40px] shadow-xl shadow-slate-200/50 border border-slate-100 overflow-hidden">
                        <div class="overflow-x-auto custom-scrollbar">
                            <table class="w-full text-left whitespace-nowrap border-collapse min-w-[500px]">
                                <thead>
                                    <tr class="bg-slate-50/50 text-[10px] font-black text-slate-400 uppercase tracking-[0.2em] border-b border-slate-100">
                                        <th class="p-6">Informasi Barang</th>
                                        <th class="p-6 text-center w-48">Jumlah Fisik (Real)</th>
                                        <th class="p-6 text-center w-24">Aksi</th>
                                    </tr>
                                </thead>
                                <tbody class="divide-y divide-slate-50">
                                    <tr v-if="cartSO.length === 0">
                                        <td colspan="3" class="p-20 text-center">
                                            <div class="flex flex-col items-center opacity-30">
                                                <svg xmlns="http://www.w3.org/2000/svg" class="w-20 h-20 mb-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" /></svg>
                                                <p class="text-xs font-black text-slate-600 uppercase tracking-[0.3em]">Belum ada barang dihitung</p>
                                            </div>
                                        </td>
                                    </tr>
                                    <tr v-for="(item, index) in cartSO" :key="index" class="hover:bg-slate-50/50 transition-colors group">
                                        <td class="p-6">
                                            <div class="font-black text-slate-800 text-sm uppercase tracking-tight">{{ item.nama_produk }}</div>
                                            <div class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-1">Item #{{ item.product_id }}</div>
                                        </td>
                                        <td class="p-6 text-center">
                                            <div class="relative w-28 mx-auto">
                                                <input v-model.number="item.actual_qty" type="number" min="0" class="w-full p-3 bg-white border-2 border-slate-200 focus:border-indigo-600 rounded-xl text-center font-black text-indigo-700 text-lg outline-none transition-all shadow-inner">
                                            </div>
                                        </td>
                                        <td class="p-6 text-center">
                                            <button @click="removeItem(index)" class="p-3 bg-slate-50 hover:bg-red-50 rounded-2xl transition-all group-hover:border-red-100 border border-transparent">
                                                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-300 hover:text-red-500 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h14" /></svg>
                                            </button>
                                        </td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>

                <div class="xl:col-span-1 space-y-6">
                    <div class="bg-white p-8 rounded-[32px] shadow-xl shadow-slate-200/50 border border-slate-100">
                        <h3 class="font-black text-slate-800 text-[10px] uppercase tracking-[0.2em] mb-6 flex items-center gap-3">
                            <span class="w-1.5 h-6 bg-indigo-600 rounded-full"></span>
                            Catatan Audit
                        </h3>
                        <textarea v-model="notes" rows="4" placeholder="Misal: Audit Lorong A, Rak 2..." class="w-full p-4 bg-slate-50 border-2 border-slate-100 focus:border-indigo-500 rounded-2xl outline-none font-bold text-sm transition-all text-slate-700 resize-none shadow-inner"></textarea>
                        
                        <div class="mt-6 p-4 bg-indigo-50 rounded-2xl border border-indigo-100">
                            <h4 class="text-[9px] font-black text-indigo-600 uppercase tracking-widest mb-1">Total Item Dihitung</h4>
                            <p class="text-2xl font-black text-indigo-900">{{ cartSO.length }} <span class="text-sm">SKU</span></p>
                        </div>
                    </div>

                    <button @click="submitSO" :disabled="isSubmitting" class="w-full bg-slate-900 hover:bg-indigo-600 text-white py-6 rounded-[30px] font-black text-xs uppercase tracking-[0.3em] shadow-2xl transition-all active:scale-[0.98] disabled:opacity-50 flex items-center justify-center gap-3 border-b-4 border-slate-800 hover:border-indigo-800">
                        <template v-if="isSubmitting">
                            <div class="w-5 h-5 border-4 border-white/30 border-t-white rounded-full animate-spin"></div>
                            Sinkronisasi...
                        </template>
                        <template v-else>
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                            Finalisasi Audit SO
                        </template>
                    </button>
                </div>
            </div>

            <div v-if="showScanner" class="fixed inset-0 bg-slate-900/90 flex items-center justify-center z-[150] p-4 backdrop-blur-sm">
                <div class="bg-white rounded-[32px] shadow-2xl w-full max-w-md overflow-hidden flex flex-col">
                    <div class="p-6 border-b border-slate-50 flex justify-between items-center bg-slate-50/50">
                        <h2 class="text-lg font-black text-slate-800 uppercase tracking-widest flex items-center gap-2">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm14 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z" /></svg>
                            Scan Barcode
                        </h2>
                        <button @click="stopScanner" class="p-2 rounded-xl bg-slate-100 text-slate-400 hover:text-red-500 hover:bg-red-50 transition-all">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                        </button>
                    </div>
                    <div class="p-6 bg-black relative">
                        <div id="reader-so" class="w-full rounded-2xl overflow-hidden border-2 border-slate-700"></div>
                        <p class="text-[10px] font-bold text-white/50 text-center mt-4 uppercase tracking-widest animate-pulse">Arahkan kamera ke barcode produk</p>
                    </div>
                </div>
            </div>

        </div>
    </Sidebar>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { height: 8px; width: 8px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #e2e8f0; border-radius: 20px; border: 2px solid transparent; background-clip: content-box; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #cbd5e1; background-clip: content-box; }
input[type=number]::-webkit-inner-spin-button, input[type=number]::-webkit-outer-spin-button { -webkit-appearance: none; margin: 0; }
</style>