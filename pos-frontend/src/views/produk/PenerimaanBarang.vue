<script setup>
import { ref, computed, onBeforeUnmount } from 'vue';
import Sidebar from '../../components/Sidebar.vue';
import api from '../../api.js';
import Swal from 'sweetalert2';
import { Html5Qrcode } from "html5-qrcode";

// --- STATE USER & ROLE ---
const getPayloadFromToken = () => {
    const token = localStorage.getItem('token');
    const role = localStorage.getItem('role') || 'staff';
    if (!token) return { role: 'staff' };
    try {
        return { role: role.toLowerCase() };
    } catch (e) { return { role: 'staff' }; }
};
const currentUser = ref(getPayloadFromToken());

// 🚀 CEK ROLE OWNER
const isOwner = computed(() => {
    return currentUser.value.role === 'owner';
});

// --- STATE ---
const supplierName = ref('');
const noFaktur = ref('');
const searchQuery = ref('');
const products = ref([]); 
const cartLPB = ref([]); 
const isSubmitting = ref(false);

// --- 🚀 LOGIKA KAMERA SCANNER BARCODE ---
const showScanner = ref(false);
let html5QrCode = null;

const startScanner = async () => {
    showScanner.value = true;
    setTimeout(async () => {
        try {
            html5QrCode = new Html5Qrcode("reader");
            await html5QrCode.start(
                { facingMode: "environment" }, 
                { fps: 10, qrbox: { width: 250, height: 100 } }, 
                (decodedText) => {
                    searchQuery.value = decodedText; 
                    stopScanner();
                    
                    // 🚀 PANGGIL FUNGSI PENCARIAN DENGAN MODE "AUTO-ADD"
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

// --- SEARCH LOGIC (DIPERBARUI UNTUK AUTO-ADD) ---
const searchProduct = async (isFromScanner = false) => {
    if (searchQuery.value.length < 2) return;
    try {
        const res = await api.get(`/products?search=${searchQuery.value}`);
        const foundData = res.data.data;

        // 🚀 JIKA DARI SCANNER, LANGSUNG TAMBAHKAN OTOMATIS
        if (isFromScanner) {
            if (foundData.length > 0) {
                addToCart(foundData[0]); // Ambil hasil pertama
                
                // Notifikasi keren di pojok kanan atas
                Swal.fire({
                    toast: true,
                    position: 'top-end',
                    icon: 'success',
                    title: `${foundData[0].nama_produk} ditambahkan!`,
                    showConfirmButton: false,
                    timer: 1500
                });
            } else {
                Swal.fire('Tidak Ditemukan!', 'Barcode ini tidak ada di Master Produk.', 'error');
                searchQuery.value = '';
            }
        } else {
            // Kalau cuma ngetik biasa, tampilin dropdown
            products.value = foundData;
        }
    } catch (err) { console.error(err); }
};

const addToCart = (product) => {
    const existing = cartLPB.value.find(item => item.product_id === product.id);
    if (existing) {
        existing.qty_masuk++;
    } else {
        cartLPB.value.push({
            product_id: product.id,
            nama_produk: product.nama_produk,
            qty_masuk: 1,
            harga_jual_saat_ini: product.harga_jual || 0,
            harga_modal_database: product.harga_modal || 0,
            harga_beli_input: product.satuan_besar ? (product.harga_modal * product.isi_per_besar) : product.harga_modal, 
            satuan_dasar: product.satuan_dasar || 'PCS',
            satuan_besar: product.satuan_besar || null,
            isi_per_besar: product.isi_per_besar || 1,
            has_satuan_besar: !!product.satuan_besar && product.isi_per_besar > 1,
            jenis_satuan: product.satuan_besar ? 'BESAR' : 'DASAR' 
        });
    }
    searchQuery.value = '';
    products.value = [];
};

const hitungTotalStok = (item) => {
    return item.jenis_satuan === 'BESAR' ? (item.qty_masuk * item.isi_per_besar) : item.qty_masuk;
};

const hitungModalPerPcs = (item) => {
    // Harga total dibagi total keseluruhan pcs
    const totalPcsMasuk = hitungTotalStok(item);

    if (totalPcsMasuk === 0) return 0

    return Math.round(item.harga_beli_input / totalPcsMasuk);
};

const removeItem = (index) => cartLPB.value.splice(index, 1);

const submitLPB = async () => {
    if (!supplierName.value || cartLPB.value.length === 0) {
        return Swal.fire('Oops!', 'Lengkapi data!', 'warning');
    }

    const adaYangRugi = isOwner.value && cartLPB.value.some(item => hitungModalPerPcs(item) >= item.harga_jual_saat_ini);
    
    const result = await Swal.fire({
        title: 'Konfirmasi LPB',
        text: adaYangRugi ? "⚠️ Modal baru rugi! Lanjut?" : "Simpan penerimaan barang?",
        icon: 'question',
        showCancelButton: true,
        confirmButtonText: 'Ya, Simpan!'
    });

    if (!result.isConfirmed) return;

    isSubmitting.value = true;
    try {
        const payloadItems = cartLPB.value.map(item => ({
            product_id: item.product_id,
            qty_masuk: hitungTotalStok(item),
            harga_modal: hitungModalPerPcs(item) 
        }));

        await api.post('/purchases', {
            supplier_name: supplierName.value,
            no_faktur: noFaktur.value,
            items: payloadItems
        });

        Swal.fire({
            icon: 'success',
            title: 'Berhasil!',
            text: isOwner.value ? 'Stok dan Modal di Master Produk telah diperbarui.' : 'Stok dan Data Faktur berhasil dicatat.',
            timer: 2500,
            showConfirmButton: false
        });

        supplierName.value = ''; noFaktur.value = ''; cartLPB.value = [];
    } catch (err) {
        console.error("DETAIL ERROR:", err.response?.data || err);
        Swal.fire('Gagal!', err.response?.data?.error || 'Error dari server backend!', 'error');
    } finally {
        isSubmitting.value = false;
    }
};
</script>

<template>
    <Sidebar>
        <div class="p-4 md:p-8 lg:p-10 max-w-7xl mx-auto font-sans bg-[#f8fafc] min-h-screen">
            
            <div class="bg-gradient-to-br from-slate-900 to-indigo-900 rounded-[40px] p-8 md:p-10 mb-10 text-white shadow-2xl flex flex-col md:flex-row items-center justify-between relative overflow-hidden border border-white/10">
                <div class="z-10">
                    <h1 class="text-3xl font-black tracking-tighter mb-2 uppercase italic leading-none">Inbound <span class="text-blue-400">Logistics</span></h1>
                    <p class="text-blue-200 font-bold text-[10px] uppercase tracking-[0.2em] flex items-center gap-2">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" /></svg>
                        Penerimaan Stok Barang
                    </p>
                </div>
                <div class="z-10 mt-4 md:mt-0 px-4 py-2 bg-white/10 rounded-2xl border border-white/20 backdrop-blur-md">
                    <span class="text-[9px] font-black uppercase tracking-widest text-blue-300">Mode Akses: {{ currentUser.role }}</span>
                </div>
            </div>

            <div class="grid grid-cols-1 xl:grid-cols-4 gap-8">
                <div class="xl:col-span-1 space-y-6">
                    <div class="bg-white p-8 rounded-[32px] shadow-xl shadow-slate-200/50 border border-slate-100">
                        <h3 class="font-black text-slate-800 text-[10px] uppercase tracking-[0.2em] mb-6 flex items-center gap-2">Data Pengiriman</h3>
                        <div class="space-y-4">
                            <input v-model="supplierName" type="text" placeholder="Nama Supplier..." class="w-full p-4 bg-slate-50 border-2 border-slate-100 rounded-2xl outline-none font-bold text-sm focus:border-blue-500">
                            <input v-model="noFaktur" type="text" placeholder="No Faktur..." class="w-full p-4 bg-slate-50 border-2 border-slate-100 rounded-2xl outline-none font-bold text-sm focus:border-blue-500 uppercase">
                        </div>
                    </div>
                </div>

                <div class="xl:col-span-3 space-y-6">
                    <div class="bg-white p-6 rounded-[32px] shadow-sm border border-slate-100 relative">
                        <div class="flex items-center gap-3">
                            <div class="relative flex-1">
                                <div class="absolute inset-y-0 left-0 pl-5 flex items-center pointer-events-none">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-400 focus-within:text-blue-500 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
                                </div>
                                <input v-model="searchQuery" @input="searchProduct(false)" type="text" placeholder="Scan Barcode / Cari Produk..." class="w-full pl-14 pr-4 py-5 bg-slate-50 border-2 border-transparent rounded-[24px] outline-none font-bold text-sm focus:border-blue-500 focus:bg-white transition-all">
                            </div>
                            
                            <button @click="startScanner" class="shrink-0 p-5 bg-blue-50 hover:bg-blue-600 text-blue-600 hover:text-white rounded-[24px] transition-all border border-blue-100 shadow-sm" title="Scan pakai Kamera">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" /><path stroke-linecap="round" stroke-linejoin="round" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
                            </button>
                        </div>
                        
                        <div v-if="products.length && searchQuery.length > 0" class="absolute left-6 right-6 mt-3 bg-white border border-slate-100 rounded-[28px] shadow-2xl z-[100] overflow-hidden">
                            <div v-for="p in products" :key="p.id" @click="addToCart(p)" class="p-5 hover:bg-blue-50 cursor-pointer border-b last:border-0 flex justify-between items-center group transition-all">
                                <div class="flex items-center gap-4">
                                    <div class="w-10 h-10 bg-slate-100 rounded-xl flex items-center justify-center text-xl shadow-inner group-hover:bg-blue-100 transition-colors">📦</div>
                                    <div class="font-black text-slate-800 uppercase text-sm group-hover:text-blue-600 transition-colors">{{ p.nama_produk }}</div>
                                </div>
                                <span class="text-[9px] bg-slate-100 text-slate-500 px-3 py-1.5 rounded-lg font-black border border-slate-200">STOK: {{ p.stok }}</span>
                            </div>
                        </div>
                    </div>

                    <div class="bg-white rounded-[40px] shadow-xl shadow-slate-200/50 border border-slate-100 overflow-hidden">
                        <div class="overflow-x-auto custom-scrollbar">
                            <table class="w-full text-left whitespace-nowrap border-collapse">
                                <thead>
                                    <tr class="bg-slate-50/50 text-[10px] font-black text-slate-400 uppercase tracking-[0.2em] border-b border-slate-100">
                                        <th class="p-6">Informasi Barang</th>
                                        <th class="p-6 text-center">Qty Masuk</th>
                                        <th class="p-6 text-center bg-blue-50/30 text-blue-600">Total Pcs</th>
                                        <th class="p-6 text-right">Harga Nota</th>
                                        <th v-if="isOwner" class="p-6 text-right">Modal Baru (PCS)</th>
                                        <th class="p-6 text-center">Aksi</th>
                                    </tr>
                                </thead>
                                <tbody class="divide-y divide-slate-50">
                                    <tr v-if="cartLPB.length === 0">
                                        <td :colspan="isOwner ? 6 : 5" class="p-16 text-center text-slate-400 font-black text-[10px] uppercase tracking-[0.3em] opacity-50">Belum ada barang di scan</td>
                                    </tr>
                                    <tr v-for="(item, index) in cartLPB" :key="item.product_id" 
                                        :class="(isOwner && hitungModalPerPcs(item) > item.harga_jual_saat_ini) ? 'bg-red-50/50' : 'hover:bg-slate-50/50'">
                                        
                                        <td class="p-6">
                                            <div class="font-black text-slate-800 text-sm uppercase tracking-tight">{{ item.nama_produk }}</div>
                                            <div v-if="isOwner && hitungModalPerPcs(item) > item.harga_jual_saat_ini" class="mt-1 text-[8px] font-black text-red-600 animate-pulse uppercase">⚠️ Potensi Rugi!</div>
                                        </td>

                                        <td class="p-6 text-center">
                                            <div class="flex items-center justify-center gap-2">
                                                <input v-model.number="item.qty_masuk" type="number" min="1" class="w-16 p-3 bg-white border-2 border-slate-100 rounded-xl text-center font-black text-blue-600 outline-none focus:border-blue-500">
                                                <select v-model="item.jenis_satuan" class="p-3 bg-white border-2 border-slate-100 rounded-xl font-black text-[10px] uppercase cursor-pointer">
                                                    <option value="DASAR">{{ item.satuan_dasar }}</option>
                                                    <option v-if="item.has_satuan_besar" value="BESAR">{{ item.satuan_besar }}</option>
                                                </select>
                                            </div>
                                        </td>

                                        <td class="p-6 text-center bg-blue-50/20">
                                            <div class="flex flex-col items-center">
                                                <span class="text-lg font-black text-blue-700 leading-none">{{ hitungTotalStok(item) }}</span>
                                                <span class="text-[8px] font-black text-blue-400 uppercase mt-1 tracking-widest">{{ item.satuan_dasar }}</span>
                                            </div>
                                        </td>

                                        <td class="p-6 min-w-[180px]"> <div class="relative w-full">
        <span class="absolute inset-y-0 left-0 pl-4 flex items-center text-xs font-black text-slate-400 italic">Rp</span>
        <input v-model.number="item.harga_beli_input" type="number" min="0" class="w-full pl-11 pr-4 py-3 bg-white border-2 border-slate-100 rounded-xl text-right font-black text-slate-800 outline-none focus:border-blue-500 transition-all text-sm shadow-inner">
    </div>
    <div class="text-[8px] text-slate-400 text-right mt-1.5 uppercase font-bold tracking-widest">
    SUBTOTAL UNTUK {{ item.qty_masuk }} {{ item.jenis_satuan === 'BESAR' ? item.satuan_besar : item.satuan_dasar }}
</div>
</td>

                                        <td v-if="isOwner" class="p-6 text-right">
                                            <div class="text-base font-black tracking-tight text-emerald-600" :class="{ 'text-red-600': hitungModalPerPcs(item) > item.harga_jual_saat_ini }">
                                                Rp {{ hitungModalPerPcs(item).toLocaleString('id-ID') }}
                                            </div>
                                        </td>

                                        <td class="p-6 text-center">
                                            <button @click="removeItem(index)" class="p-3 bg-slate-50 hover:bg-red-50 rounded-2xl transition-all group border border-transparent hover:border-red-100">
                                                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-300 group-hover:text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h14" /></svg>
                                            </button>
                                        </td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                    </div>

                    <button @click="submitLPB" :disabled="isSubmitting" class="w-full bg-blue-600 hover:bg-slate-900 text-white py-6 rounded-[30px] font-black text-xs uppercase tracking-[0.3em] shadow-2xl transition-all active:scale-[0.98] disabled:opacity-50 flex items-center justify-center gap-4">
                        <template v-if="isSubmitting"><div class="w-5 h-5 border-4 border-white/30 border-t-white rounded-full animate-spin"></div>Sinkronisasi Stok...</template>
                        <template v-else>Posting Penerimaan Barang</template>
                    </button>
                </div>
            </div>

            <div v-if="showScanner" class="fixed inset-0 bg-slate-900/90 flex items-center justify-center z-[150] p-4 backdrop-blur-sm">
                <div class="bg-white rounded-[32px] shadow-2xl w-full max-w-md overflow-hidden flex flex-col">
                    <div class="p-6 border-b border-slate-50 flex justify-between items-center bg-slate-50/50">
                        <h2 class="text-lg font-black text-slate-800 uppercase tracking-widest flex items-center gap-2">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm14 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z" /></svg>
                            Scan Barcode Produk
                        </h2>
                        <button @click="stopScanner" class="p-2 rounded-xl bg-slate-100 text-slate-400 hover:text-red-500 hover:bg-red-50 transition-all">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                        </button>
                    </div>
                    <div class="p-6 bg-black relative">
                        <div id="reader" class="w-full rounded-2xl overflow-hidden border-2 border-slate-700"></div>
                        <p class="text-[10px] font-bold text-white/50 text-center mt-4 uppercase tracking-widest animate-pulse">Arahkan kamera ke barcode</p>
                    </div>
                </div>
            </div>

        </div>
    </Sidebar>
</template>

<style scoped>
/* ... (Style tetap sama) ... */
.custom-scrollbar::-webkit-scrollbar { height: 8px; width: 8px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #e2e8f0; border-radius: 20px; border: 2px solid transparent; background-clip: content-box; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #cbd5e1; background-clip: content-box; }
input[type=number]::-webkit-inner-spin-button, input[type=number]::-webkit-outer-spin-button { -webkit-appearance: none; margin: 0; }
</style>