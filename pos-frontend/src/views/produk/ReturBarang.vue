<script setup>
import { ref, onMounted, onUnmounted, computed, watch, nextTick } from 'vue';
import { useRouter } from 'vue-router';
import Sidebar from '../../components/Sidebar.vue';
import api from '../../api.js';
import Swal from 'sweetalert2';
import { Html5Qrcode } from "html5-qrcode";

const router = useRouter();

// --- STATE DATA ---
const products = ref([]);
const cart = ref([]); 
const isLoading = ref(true);
const isSubmitting = ref(false);

// 🚀 STATE CETAK LANGSUNG
const lastSubmittedReturn = ref(null); // Nyimpen data sementara buat dicetak

// 🚀 STATE PENCARIAN PRODUK
const searchProductQuery = ref('');
const isDropdownOpen = ref(false);
const selectedProduct = ref(null);

// 🚀 STATE KAMERA SCANNER
const isScannerOpen = ref(false);
const cameras = ref([]);
const selectedCamera = ref('');
let html5QrCode = null;

const form = ref({
    product_id: '',
    qty: 1,
    alasan: '',
    catatan: ''
});

const alasanOptions = [
    { value: 'Expired / Basi', label: 'Expired / Basi' },
    { value: 'Rusak Fisik / Pecah', label: 'Rusak Fisik / Pecah' },
    { value: 'Retur ke Supplier', label: 'Retur ke Supplier' },
];

// --- FETCH DATA ---
const fetchProducts = async () => {
    isLoading.value = true;
    try {
        const resProd = await api.get('/products');
        const allProducts = resProd.data.data || [];
        products.value = allProducts.filter(p => p.stok > 0);
    } catch (error) {
        Swal.fire('Error', 'Gagal memuat data produk.', 'error');
    } finally {
        isLoading.value = false;
    }
};

onMounted(() => fetchProducts());

onUnmounted(() => {
    if (html5QrCode && html5QrCode.isScanning) {
        html5QrCode.stop().then(() => html5QrCode.clear());
    }
});

// 🚀 FITUR KAMERA SCANNER
const getCameras = async () => {
    try {
        const devices = await Html5Qrcode.getCameras();
        if (devices && devices.length) {
            cameras.value = devices;
            selectedCamera.value = devices[1].id;
        }
    } catch (err) {
        console.error("Gagal mendapatkan kamera:", err);
    }
};

const startScanner = async () => {
    isScannerOpen.value = true;
    await nextTick();
    await getCameras();

    if (cameras.value.length === 0) {
        Swal.fire('Error', 'Tidak ada kamera terdeteksi!', 'error');
        isScannerOpen.value = false;
        return;
    }

    html5QrCode = new Html5Qrcode("reader");
    startScanning();
};

const startScanning = () => {
    if (!selectedCamera.value || !html5QrCode) return;
    const config = { fps: 10, qrbox: { width: 250, height: 150 } };
    
    html5QrCode.start(
        selectedCamera.value,
        config,
        (decodedText) => {
            searchProductQuery.value = decodedText;
            const audio = new Audio('https://www.soundjay.com/buttons/beep-07.wav');
            audio.play().catch(() => {});
            stopScanner();
        },
        (errorMessage) => {}
    ).catch(err => console.error("Scanner error:", err));
};

const switchCamera = async () => {
    if (html5QrCode && html5QrCode.isScanning) {
        await html5QrCode.stop();
        startScanning();
    }
};

const stopScanner = async () => {
    if (html5QrCode && html5QrCode.isScanning) {
        await html5QrCode.stop();
        html5QrCode.clear();
    }
    isScannerOpen.value = false;
};

// 🚀 PENCARIAN PRODUK REALTIME
const filteredProducts = computed(() => {
    if (!searchProductQuery.value) return [];
    const query = searchProductQuery.value.toLowerCase();
    return products.value.filter(p => 
        (p.nama_produk && p.nama_produk.toLowerCase().includes(query)) || 
        (p.sku && p.sku.toLowerCase().includes(query))
    ).slice(0, 10);
});

const selectProduct = (prod) => {
    selectedProduct.value = prod;
    form.value.product_id = prod.id;
    searchProductQuery.value = ''; 
    isDropdownOpen.value = false;
};

watch(searchProductQuery, (newVal) => {
    if (newVal) {
        isDropdownOpen.value = true;
        const exactMatch = products.value.find(p => p.sku === newVal);
        if (exactMatch) { selectProduct(exactMatch); }
    } else {
        isDropdownOpen.value = false;
    }
});

const clearSelectedProduct = () => {
    selectedProduct.value = null;
    form.value.product_id = '';
};

// 🚀 KERANJANG LOGIC
const addToCart = () => {
    if (!form.value.product_id || !form.value.alasan || form.value.qty < 1) {
        return Swal.fire('Data Kurang', 'Pilih produk, alasan, dan qty wajib diisi!', 'warning');
    }

    if (form.value.qty > selectedProduct.value.stok) {
        return Swal.fire('Stok Tidak Cukup', `Sisa stok ${selectedProduct.value.nama_produk} hanya ${selectedProduct.value.stok}!`, 'error');
    }

    const existingIndex = cart.value.findIndex(item => item.product_id === form.value.product_id && item.alasan === form.value.alasan);
    
    if (existingIndex !== -1) {
        if ((cart.value[existingIndex].qty + form.value.qty) > selectedProduct.value.stok) {
             return Swal.fire('Melebihi Stok', `Total di keranjang + input baru melebihi stok yang ada!`, 'error');
        }
        cart.value[existingIndex].qty += form.value.qty;
    } else {
        cart.value.push({
            product_id: form.value.product_id,
            nama_produk: selectedProduct.value.nama_produk,
            sku: selectedProduct.value.sku,
            qty: form.value.qty,
            alasan: form.value.alasan,
            catatan: form.value.catatan
        });
    }

    clearSelectedProduct();
    form.value.qty = 1;
    form.value.alasan = '';
    form.value.catatan = '';
};

const removeFromCart = (index) => {
    cart.value.splice(index, 1);
};

// 🚀 SUBMIT BATCH & CETAK LANGSUNG
const submitBatchReturn = async () => {
    if (cart.value.length === 0) return;

    const confirm = await Swal.fire({
        title: 'Proses Berita Acara?',
        html: `Ada <b>${cart.value.length} item</b> di keranjang. Stok akan dipotong permanen!`,
        icon: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#e11d48',
        cancelButtonColor: '#94a3b8',
        confirmButtonText: 'Ya, Proses Sekarang!'
    });

    if (!confirm.isConfirmed) return;

    isSubmitting.value = true;
    try {
        const payload = {
            items: cart.value.map(item => ({
                product_id: item.product_id,
                qty: item.qty,
                alasan: item.alasan,
                catatan: item.catatan
            }))
        };

        const res = await api.post('/returns', payload);
        
        // 🚀 SIMPAN DATA SEMENTARA BUAT DICETAK SEKARANG JUGA
        lastSubmittedReturn.value = {
            return_no: res.data.return_no,
            created_at: new Date(),
            items: [...cart.value],
            total_qty: cart.value.reduce((acc, curr) => acc + curr.qty, 0),
            user: { name: localStorage.getItem('name') || 'Kasir' },
            storeName: localStorage.getItem('storeName') || 'POS UMKM'
        };

        // Reset keranjang
        cart.value = [];
        fetchProducts(); 

        // 🚀 TAMPILKAN SWEETALERT DENGAN TOMBOL CETAK
        const resultPrint = await Swal.fire({
            icon: 'success',
            title: 'Berhasil Diproses!',
            html: `Dokumen: <b>${res.data.return_no}</b><br/>Stok telah berhasil dipotong dari sistem.`,
            showCancelButton: true,
            confirmButtonText: '<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 inline-block mr-1" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="6 9 6 2 18 2 18 9"/><path d="M6 18H4a2 2 0 0 1-2-2v-5a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v5a2 2 0 0 1-2 2h-2"/><rect width="12" height="8" x="6" y="14"/></svg> Cetak Bukti',
            cancelButtonText: 'Tutup',
            confirmButtonColor: '#4f46e5', // Warna Indigo/Biru
            cancelButtonColor: '#94a3b8',
            customClass: {
                confirmButton: 'rounded-xl font-black px-6 py-3 flex items-center gap-2',
                cancelButton: 'rounded-xl font-black px-6 py-3'
            }
        });

        // Kalau kasir ngeklik Cetak Bukti
        if (resultPrint.isConfirmed) {
            setTimeout(() => {
                window.print();
            }, 300);
        }

    } catch (error) {
        Swal.fire('Gagal!', error.response?.data?.error || 'Sistem gagal memproses retur.', 'error');
    } finally {
        isSubmitting.value = false;
    }
};

const getBadgeClass = (alasan) => {
    if (alasan.includes('Expired')) return 'bg-rose-50 text-rose-600 border-rose-200';
    if (alasan.includes('Rusak')) return 'bg-amber-50 text-amber-600 border-amber-200';
    if (alasan.includes('Retur')) return 'bg-blue-50 text-blue-600 border-blue-200';
    return 'bg-slate-50 text-slate-600 border-slate-200';
};
</script>

<template>
    <Sidebar>
        <div class="p-4 md:p-8 lg:p-10 max-w-7xl mx-auto font-sans bg-[#f8fafc] min-h-screen print:bg-white print:p-0">
            
            <div class="bg-gradient-to-br from-rose-900 via-rose-800 to-slate-900 rounded-[32px] p-6 md:p-10 mb-8 text-white shadow-2xl flex flex-col md:flex-row items-center justify-between relative overflow-hidden border border-rose-800 gap-6 print:hidden">
                <svg xmlns="http://www.w3.org/2000/svg" class="absolute -right-10 -bottom-10 w-64 h-64 text-rose-500 opacity-10 pointer-events-none" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/></svg>
                
                <div class="z-10 text-center md:text-left">
                    <h1 class="text-3xl md:text-4xl font-black tracking-tighter mb-1 uppercase italic">Waste & <span class="text-rose-400">Return</span></h1>
                    <p class="text-rose-200 font-bold text-[10px] md:text-xs uppercase tracking-[0.2em] flex items-center justify-center md:justify-start gap-2">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
                        Keranjang Pemusnahan & Retur Barang
                    </p>
                </div>

                <div class="z-10">
                    <router-link to="/returns/report" class="bg-white/10 hover:bg-white/20 border border-white/20 backdrop-blur-sm text-white px-6 py-3.5 rounded-2xl font-black text-[10px] uppercase tracking-[0.2em] flex items-center justify-center gap-2 transition-all active:scale-95 shadow-lg">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/><polyline points="10 9 9 9 8 9"/></svg>
                        Riwayat & Cetak B.A
                    </router-link>
                </div>
            </div>

            <div class="grid grid-cols-1 xl:grid-cols-2 gap-8 print:hidden">
                
                <div class="xl:col-span-1">
                    <div class="bg-white rounded-[32px] p-6 md:p-8 shadow-sm border border-slate-100 xl:sticky xl:top-24">
                        <h2 class="text-lg font-black text-slate-800 uppercase tracking-tight mb-6 flex items-center gap-2 border-b border-slate-100 pb-4">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-rose-500" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M11 2a2 2 0 0 0-2 2v5H4a2 2 0 0 0-2 2v2c0 1.1.9 2 2 2h5v5c0 1.1.9 2 2 2h2a2 2 0 0 0 2-2v-5h5a2 2 0 0 0 2-2v-2a2 2 0 0 0-2-2h-5V4a2 2 0 0 0-2-2h-2z"/></svg>
                            Pilih Barang
                        </h2>

                        <form @submit.prevent="addToCart" class="flex flex-col gap-5">
                            <div class="relative">
                                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block">Pilih Barang Bermasalah</label>
                                
                                <div v-if="selectedProduct" class="flex items-center justify-between p-4 bg-rose-50 border-2 border-rose-200 rounded-2xl">
                                    <div class="flex flex-col">
                                        <span class="text-rose-900 font-black text-sm uppercase leading-tight truncate">{{ selectedProduct.nama_produk }}</span>
                                        <span class="text-[10px] font-bold text-rose-500 mt-0.5">Stok Tersedia: {{ selectedProduct.stok }}</span>
                                    </div>
                                    <button type="button" @click="clearSelectedProduct" class="w-8 h-8 rounded-xl bg-white text-rose-400 hover:text-rose-600 hover:bg-rose-100 flex items-center justify-center transition-all shadow-sm">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                                    </button>
                                </div>

                                <div v-else class="relative">
                                    <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M4 5v14M8 5v14M12 5v14M16 5v14M20 5v14"/><line x1="2" y1="12" x2="22" y2="12" stroke="red" stroke-width="1.5"/></svg>
                                    </div>
                                    <input 
                                        v-model="searchProductQuery" 
                                        type="text" 
                                        class="w-full pl-12 pr-14 py-4 bg-slate-50 border-2 border-slate-100 rounded-2xl focus:border-rose-500 focus:bg-white outline-none font-bold text-slate-800 text-sm transition-all placeholder:text-slate-300" 
                                        placeholder="Ketik nama atau SKU..."
                                    >
                                    <button type="button" @click="startScanner" class="absolute inset-y-0 right-2 my-auto w-10 h-10 bg-slate-200 hover:bg-rose-100 text-slate-500 hover:text-rose-600 rounded-xl flex items-center justify-center transition-all">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M23 19a2 2 0 0 1-2-2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z"/><circle cx="12" cy="13" r="4"/></svg>
                                    </button>
                                    
                                    <div v-if="isDropdownOpen && filteredProducts.length > 0" class="absolute z-50 w-full mt-2 bg-white border border-slate-100 rounded-2xl shadow-xl overflow-hidden max-h-60 overflow-y-auto">
                                        <div 
                                            v-for="prod in filteredProducts" 
                                            :key="prod.id" 
                                            @click="selectProduct(prod)"
                                            class="p-4 border-b border-slate-50 hover:bg-rose-50 cursor-pointer transition-colors flex justify-between items-center group"
                                        >
                                            <div>
                                                <div class="font-black text-slate-800 text-xs uppercase group-hover:text-rose-700">{{ prod.nama_produk }}</div>
                                                <div class="text-[9px] font-bold text-slate-400 mt-1">{{ prod.sku }}</div>
                                            </div>
                                            <span class="text-[10px] font-black bg-slate-100 text-slate-500 px-2 py-1 rounded-lg group-hover:bg-rose-100 group-hover:text-rose-600">Sisa: {{ prod.stok }}</span>
                                        </div>
                                    </div>
                                </div>
                            </div>

                            <div class="grid grid-cols-2 gap-4">
                                <div>
                                    <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block">Qty Buang</label>
                                    <input v-model="form.qty" type="number" min="1" required class="w-full px-4 py-3.5 bg-slate-50 border-2 border-slate-100 rounded-2xl focus:border-rose-500 focus:bg-white outline-none font-black text-rose-600 text-lg transition-all disabled:opacity-50" placeholder="1" :disabled="!selectedProduct">
                                </div>

                                <div>
                                    <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block">Alasan</label>
                                    <div class="relative">
                                        <select v-model="form.alasan" required class="w-full px-4 py-4 bg-slate-50 border-2 border-slate-100 rounded-2xl focus:border-rose-500 focus:bg-white outline-none font-bold text-slate-800 text-xs appearance-none cursor-pointer transition-all disabled:opacity-50" :disabled="!selectedProduct">
                                            <option value="" disabled selected hidden>Pilih...</option>
                                            <option v-for="opt in alasanOptions" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
                                        </select>
                                        <div class="absolute inset-y-0 right-3 flex items-center pointer-events-none">
                                            <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" /></svg>
                                        </div>
                                    </div>
                                </div>
                            </div>

                            <div>
                                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block">Keterangan Opsional</label>
                                <textarea v-model="form.catatan" rows="1" class="w-full px-4 py-3.5 bg-slate-50 border-2 border-slate-100 rounded-2xl focus:border-rose-500 focus:bg-white outline-none font-medium text-slate-700 text-sm transition-all resize-none placeholder:text-slate-300 disabled:opacity-50" placeholder="Contoh: Jatuh pecah..." :disabled="!selectedProduct"></textarea>
                            </div>

                            <button type="submit" :disabled="!selectedProduct" class="w-full bg-slate-100 hover:bg-rose-50 border-2 border-slate-200 hover:border-rose-200 text-slate-600 hover:text-rose-600 py-4 rounded-[20px] font-black text-xs uppercase tracking-[0.2em] transition-all active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed mt-2 flex items-center justify-center gap-2">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
                                TAMBAH KE KERANJANG
                            </button>
                        </form>
                    </div>
                </div>

                <div class="xl:col-span-1 flex flex-col h-full">
                    <div class="bg-white rounded-[32px] shadow-sm border border-slate-100 overflow-hidden flex flex-col flex-1 h-[600px] xl:h-auto xl:sticky xl:top-24">
                        <div class="p-6 border-b border-slate-50 bg-slate-50/50 flex items-center justify-between shrink-0">
                            <h3 class="font-black text-slate-800 text-sm uppercase tracking-widest flex items-center gap-2">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-rose-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z"/><line x1="3" y1="6" x2="21" y2="6"/><path d="M16 10a4 4 0 0 1-8 0"/></svg>
                                Keranjang Retur
                                <span v-if="cart.length > 0" class="ml-2 bg-rose-600 text-white w-6 h-6 rounded-full flex items-center justify-center text-[10px]">{{ cart.length }}</span>
                            </h3>
                        </div>

                        <div v-if="cart.length === 0" class="flex-1 flex flex-col items-center justify-center text-center p-8">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-20 h-20 text-slate-200 mb-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z"/><line x1="3" y1="6" x2="21" y2="6"/><path d="M16 10a4 4 0 0 1-8 0"/></svg>
                            <p class="text-slate-400 font-black text-sm uppercase tracking-widest">Keranjang Kosong</p>
                            <p class="text-slate-400 font-medium text-xs mt-2 max-w-[200px]">Pilih barang di form samping untuk memproses retur.</p>
                        </div>

                        <div v-auto-animate class="flex-1 overflow-y-auto custom-scrollbar p-4 bg-slate-50/30 flex flex-col gap-3">
                            <div v-for="(item, index) in cart" :key="index" class="bg-white p-4 rounded-[20px] shadow-sm border border-slate-100 flex gap-3 relative group">
                                <div class="w-12 h-12 bg-rose-50 text-rose-600 rounded-xl flex items-center justify-center font-black text-lg border border-rose-100 shrink-0">
                                    {{ item.qty }}
                                </div>
                                <div class="flex-1 pr-8">
                                    <div class="font-black text-slate-800 text-sm uppercase leading-tight">{{ item.nama_produk }}</div>
                                    <div class="mt-1 flex flex-wrap gap-1">
                                        <span class="inline-block px-2 py-0.5 rounded text-[8px] font-black uppercase tracking-widest border" :class="getBadgeClass(item.alasan)">{{ item.alasan }}</span>
                                    </div>
                                    <p v-if="item.catatan" class="text-[10px] text-slate-500 mt-1.5 italic truncate">"{{ item.catatan }}"</p>
                                </div>
                                <button @click="removeFromCart(index)" class="absolute top-4 right-4 w-8 h-8 flex items-center justify-center rounded-lg bg-red-50 text-red-400 hover:text-white hover:bg-red-500 transition-colors opacity-100 xl:opacity-0 xl:group-hover:opacity-100">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/></svg>
                                </button>
                            </div>
                        </div>

                        <div v-if="cart.length > 0" class="p-6 bg-white border-t border-slate-100 shrink-0">
                            <button @click="submitBatchReturn" :disabled="isSubmitting" class="w-full bg-rose-600 hover:bg-slate-900 text-white py-5 rounded-[24px] font-black text-sm uppercase tracking-[0.2em] shadow-xl shadow-rose-200 transition-all active:scale-95 disabled:opacity-50 flex items-center justify-center gap-2">
                                <template v-if="isSubmitting">
                                    <div class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                                    MENYIMPAN...
                                </template>
                                <template v-else>
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"/></svg>
                                    PROSES PEMUSNAHAN
                                </template>
                            </button>
                        </div>
                    </div>
                </div>
            </div>

            <div v-if="lastSubmittedReturn" id="printable-area" class="hidden print:block print:p-0 print:text-black print:absolute print:inset-0 print:w-full">
                <div class="text-center border-b-[3px] border-slate-800 pb-6 mb-6">
                    <h1 class="text-2xl font-black uppercase tracking-widest text-slate-900">{{ lastSubmittedReturn.storeName }}</h1>
                    <p class="text-sm font-medium text-slate-600 mt-1">BERITA ACARA PEMUSNAHAN / RETUR BARANG</p>
                </div>

                <div class="flex justify-between items-end mb-8">
                    <div>
                        <table class="text-xs font-bold text-slate-700">
                            <tbody><tr><td class="pb-2 pr-4 text-slate-500 uppercase tracking-widest text-[9px]">No. Dokumen</td><td class="pb-2">: {{ lastSubmittedReturn.return_no }}</td></tr>
                            <tr><td class="pb-2 pr-4 text-slate-500 uppercase tracking-widest text-[9px]">Tanggal</td><td class="pb-2">: {{ new Date(lastSubmittedReturn.created_at).toLocaleDateString('id-ID', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' }) }}</td></tr>
                            <tr><td class="pr-4 text-slate-500 uppercase tracking-widest text-[9px]">Operator</td><td>: {{ lastSubmittedReturn.user.name }}</td></tr>
                            </tbody>
                        </table>
                    </div>
                    <div class="text-right">
                        <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1">Total Kuantitas</p>
                        <p class="text-3xl font-black text-slate-900 tracking-tighter">{{ lastSubmittedReturn.total_qty }}</p>
                    </div>
                </div>

                <table class="w-full text-left border-collapse mb-12">
                    <thead>
                        <tr class="border-y-2 border-slate-800">
                            <th class="py-3 px-2 text-[10px] font-black text-slate-800 uppercase tracking-widest w-12 text-center">No</th>
                            <th class="py-3 px-2 text-[10px] font-black text-slate-800 uppercase tracking-widest">SKU</th>
                            <th class="py-3 px-2 text-[10px] font-black text-slate-800 uppercase tracking-widest">Nama Barang</th>
                            <th class="py-3 px-2 text-[10px] font-black text-slate-800 uppercase tracking-widest text-center">Qty</th>
                            <th class="py-3 px-2 text-[10px] font-black text-slate-800 uppercase tracking-widest">Alasan / Klasifikasi</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-slate-200">
                        <tr v-for="(item, index) in lastSubmittedReturn.items" :key="index">
                            <td class="py-3 px-2 text-xs font-bold text-slate-600 text-center">{{ index + 1 }}</td>
                            <td class="py-3 px-2 text-[10px] font-bold text-slate-500 uppercase tracking-wider">{{ item.sku || '-' }}</td>
                            <td class="py-3 px-2 text-xs font-black text-slate-800 uppercase">{{ item.nama_produk }}</td>
                            <td class="py-3 px-2 text-sm font-black text-slate-900 text-center">{{ item.qty }}</td>
                            <td class="py-3 px-2">
                                <div class="text-xs font-bold text-slate-700">{{ item.alasan }}</div>
                                <div v-if="item.catatan" class="text-[10px] font-medium text-slate-500 italic mt-0.5">Catatan: {{ item.catatan }}</div>
                            </td>
                        </tr>
                    </tbody>
                </table>

                <div class="grid grid-cols-2 gap-8 mt-16 pt-8 break-inside-avoid">
                    <div class="text-center">
                        <p class="text-[10px] font-bold text-slate-500 uppercase tracking-widest mb-16">Dibuat Oleh,</p>
                        <p class="text-sm font-black text-slate-800 uppercase underline">{{ lastSubmittedReturn.user.name }}</p>
                        <p class="text-[9px] font-bold text-slate-400 mt-1">Staff / Kasir</p>
                    </div>
                    <div class="text-center">
                        <p class="text-[10px] font-bold text-slate-500 uppercase tracking-widest mb-16">Mengetahui,</p>
                        <p class="text-sm font-black text-slate-800 uppercase underline">...................................</p>
                        <p class="text-[9px] font-bold text-slate-400 mt-1">Manager / Owner</p>
                    </div>
                </div>
            </div>

        </div>

        <div v-if="isScannerOpen" class="fixed inset-0 bg-slate-900/90 flex items-center justify-center z-[100] p-4 backdrop-blur-sm print:hidden">
            <div class="bg-white rounded-[32px] overflow-hidden w-full max-w-sm shadow-2xl border-[8px] border-slate-800 flex flex-col">
                <div class="p-5 border-b border-slate-100 flex justify-between items-center bg-white shrink-0">
                    <div>
                        <h3 class="font-black text-slate-800 uppercase tracking-tighter text-lg italic">Scan Barcode</h3>
                        <p class="text-[9px] text-rose-500 font-black uppercase tracking-widest mt-0.5">Retur / Buang Barang</p>
                    </div>
                    <button @click="stopScanner" class="w-10 h-10 rounded-xl bg-slate-100 text-slate-400 hover:text-red-500 hover:bg-red-50 transition-all flex items-center justify-center">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                    </button>
                </div>
                
                <div v-if="cameras.length > 0" class="p-3 bg-slate-50 border-b border-slate-100 flex gap-2 items-center">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-400 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M14.5 4h-5L7 7H4a2 2 0 0 0-2 2v9a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2h-3l-2.5-3z"/><circle cx="12" cy="13" r="3"/></svg>
                    <select v-model="selectedCamera" @change="switchCamera" class="w-full bg-white border border-slate-200 text-[10px] font-black uppercase tracking-widest text-slate-600 rounded-lg p-2 outline-none cursor-pointer">
                        <option v-for="cam in cameras" :key="cam.id" :value="cam.id">
                            {{ cam.label || `Kamera ${cam.id.substring(0, 5)}...` }}
                        </option>
                    </select>
                </div>
                
                <div class="relative bg-black w-full aspect-square flex items-center justify-center">
                    <div id="reader" class="w-full h-full object-cover"></div>
                    <div class="absolute inset-0 border-[12px] border-black/40 pointer-events-none z-10"></div>
                    <div class="absolute inset-x-8 inset-y-16 border-2 border-white/30 rounded-[20px] pointer-events-none z-20">
                        <div class="absolute top-0 left-0 w-8 h-8 border-t-4 border-l-4 border-rose-500 rounded-tl-[18px]"></div>
                        <div class="absolute top-0 right-0 w-8 h-8 border-t-4 border-r-4 border-rose-500 rounded-tr-[18px]"></div>
                        <div class="absolute bottom-0 left-0 w-8 h-8 border-b-4 border-l-4 border-rose-500 rounded-bl-[18px]"></div>
                        <div class="absolute bottom-0 right-0 w-8 h-8 border-b-4 border-r-4 border-rose-500 rounded-br-[18px]"></div>
                        <div class="w-full h-0.5 bg-rose-500 absolute top-0 animate-[scan_2s_infinite] shadow-[0_0_8px_#f43f5e]"></div>
                    </div>
                </div>
            </div>
        </div>
    </Sidebar>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { height: 6px; width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #94a3b8; }

input[type=number]::-webkit-inner-spin-button, 
input[type=number]::-webkit-outer-spin-button { 
    -webkit-appearance: none; 
    margin: 0; 
}
input[type=number] {
    -moz-appearance: textfield;
}

select {
    -webkit-appearance: none;
    -moz-appearance: none;
    appearance: none;
}

@keyframes scan {
  0% { top: 0%; opacity: 0; }
  10% { opacity: 1; }
  90% { opacity: 1; }
  100% { top: 100%; opacity: 0; }
}

:deep(#reader) { border: none !important; }
:deep(#reader video) { object-fit: cover !important; }

/* 🚀 CSS SAKTI BUAT CETAK PDF (Cetak Langsung dari Halaman Cart) */
@media print {
    body * {
        visibility: hidden;
    }
    
    #printable-area, #printable-area * {
        visibility: visible;
    }
    
    #printable-area {
        position: absolute;
        left: 0;
        top: 0;
        width: 100%;
        margin: 0;
        padding: 0;
    }

    body {
        background-color: white !important;
    }

    table { page-break-inside: auto; }
    tr    { page-break-inside: avoid; page-break-after: auto; }
    
    @page {
        margin: 20mm;
        size: A4 portrait;
    }
}
</style>