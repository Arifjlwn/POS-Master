<script setup>
import { ref, computed, onMounted, nextTick, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import api from '../../api.js';
import Swal from 'sweetalert2';
import { Html5Qrcode } from "html5-qrcode";

const router = useRouter();

// --- SETUP USER & ROLE ---
const getUserInfo = () => {
    const token = localStorage.getItem('token');
    const role = localStorage.getItem('role') || 'kasir';
    let name = localStorage.getItem('name'); 

    if (token) {
        try {
            const payload = JSON.parse(atob(token.split('.')[1]));
            if (!name || name === 'undefined' || name === '') {
                name = payload.name || payload.username || 'Kasir Toko';
            }
            return { userId: payload.user_id, role, name };
        } catch (e) {
            return { userId: 0, role, name: 'Kasir Toko' };
        }
    }
    return { userId: 0, role, name: 'Kasir Toko' };
};

const currentUser = ref(getUserInfo());
const currentSession = ref(null);

// Fungsi Jam Realtime
const currentTime = ref('');
let timer;

// State utama
const products = ref([]);
const isLoadingProducts = ref(true);
const cart = ref([]);
const heldOrders = ref([]);
const showHeldModal = ref(false);
const payAmount = ref(0);
const paymentMethod = ref('Cash');
const showReceipt = ref(false);
const showQrisModal = ref(false);
const lastTransaction = ref(null);
const showReceiptClosing = ref(false);
const lastClosingData = ref(null);

// 🚀 STATE KHUSUS UNTUK HP (MOBILE CART DRAWER)
const isMobileCartOpen = ref(false);

const getImageUrl = (path) => {
    if (!path) return null;
    return `${import.meta.env.VITE_API_BASE_URL}${path}`;
}

const searchQuery = ref('');
const searchInput = ref(null);

// --- LOGIKA KAMERA SCANNER KASIR ---
const showScanner = ref(false);
let html5QrCode = null;

const startScanner = async () => {
    showScanner.value = true;
    setTimeout(async () => {
        try {
            html5QrCode = new Html5Qrcode("reader-kasir");
            await html5QrCode.start(
                { facingMode: "environment" }, 
                { fps: 15, qrbox: { width: 250, height: 100 } }, 
                (decodedText) => {
                    searchQuery.value = decodedText; 
                    stopScanner();
                    handleBarcodeScan(); 
                    const audio = new Audio('https://www.soundjay.com/buttons/beep-07a.mp3');
                    audio.play().catch(()=>{}); 
                },
                (errorMessage) => {} 
            );
        } catch (err) {
            console.error(err);
            Swal.fire('Oops!', 'Gagal menyalakan kamera. Pastikan izin kamera aktif.', 'error');
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

onUnmounted(() => {
    if (showScanner.value) stopScanner();
    clearInterval(timer);
});

const fetchProducts = async () => {
    try {
        const response = await api.get('/products');
        products.value = response.data.data.map(p => ({
            id: p.id,
            sku: p.sku || `SKU-${p.id}`,
            name: p.nama_produk,
            price: p.harga_jual,
            stock: p.stok,
            image: p.gambar
        }));
    } catch (error) {
        console.error("Gagal narik produk:", error);
    } finally {
        isLoadingProducts.value = false;
    }
};

onMounted(async () => {
    const token = localStorage.getItem('token');
    if (!token) { router.push('/login'); return; }

    try {
        const res = await api.get('/pos/check-session', {
            headers: { Authorization: `Bearer ${token}` }
        });

        if (!res.data.has_session) {
            Swal.fire('Akses Ditolak', 'Isi modal awal atau absen dulu ya!', 'warning');
            router.push('/pos/buka-kasir');
            return;
        }

        currentSession.value = res.data.session;
        await fetchProducts();
        if (searchInput.value) searchInput.value.focus();

    } catch (error) {
        if (error.response?.status === 401) {
            router.push('/login');
        }
    }

    timer = setInterval(() => {
        const now = new Date();
        currentTime.value = now.toLocaleString('id-ID', {
            day: '2-digit', month: '2-digit', year: 'numeric',
            hour: '2-digit', minute: '2-digit', second: '2-digit'
        }).replace(/\//g, '.');
    }, 1000);
});

const filteredProducts = computed(() => {
    if (!searchQuery.value) return products.value;
    const query = searchQuery.value.toLocaleLowerCase();
    return products.value.filter(product =>
        product.name.toLowerCase().includes(query) ||
        (product.sku && product.sku.toLowerCase().includes(query))
    );
});

const handleBarcodeScan = () => {
    if (!searchQuery.value) return;
    const query = String(searchQuery.value).trim().toLowerCase();
    const exactMatch = products.value.find(p => p.sku && String(p.sku).toLowerCase() === query);

    if (exactMatch) {
        addToCart(exactMatch);
        searchQuery.value = '';
    } else if (filteredProducts.value.length === 1) {
        addToCart(filteredProducts.value[0]);
        searchQuery.value = '';
    }

    nextTick(() => {
        if (searchInput.value) searchInput.value.focus();
    });
};

const addToCart = (product) => {
    if (product.stock <= 0) {
        Swal.fire({ icon: 'error', title: 'Stok Habis!', text: `Stok ${product.name} kosong` });
        return;
    }
    const existingItem = cart.value.find(item => item.id === product.id);
    if (existingItem) {
        if (existingItem.qty < product.stock) {
            existingItem.qty++;
        } else {
            Swal.fire({ icon: 'warning', title: 'Stok Terbatas', text: 'Kuantitas melebihi stok!' });
            return;
        }
    } else {
        cart.value.unshift({ id: product.id, name: product.name, price: product.price, qty: 1 });
    }
    
    // 🚀 Feedback visual halus buat HP pas masukin barang tapi keranjang lagi ditutup
    if (window.innerWidth < 1024 && !isMobileCartOpen.value) {
        Swal.fire({
            toast: true,
            position: 'top',
            icon: 'success',
            title: `${product.name} Masuk Keranjang`,
            showConfirmButton: false,
            timer: 800,
            timerProgressBar: true
        });
    }
};

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

const validateQty = (product) => {
    const existingItem = cart.value.find(item => item.id === product.id);
    if (existingItem && existingItem.qty < 1) existingItem.qty = 1;
};

const clearCart = () => {
    if (cart.value.length === 0) return;
    Swal.fire({
        title: 'Batalkan Transaksi?',
        text: "Semua barang di keranjang akan dihapus!",
        icon: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#ef4444',
        confirmButtonText: 'Ya, Hapus Semua!',
        cancelButtonText: 'Batal'
    }).then((result) => {
        if (result.isConfirmed) {
            cart.value = [];
            payAmount.value = 0;
            setPaymentMethod('Cash');
            isMobileCartOpen.value = false;
        }
    });
};

const holdTransaction = () => {
    if (cart.value.length === 0) return;
    heldOrders.value.push({
        id: Date.now(),
        customer: `Pelanggan ${heldOrders.value.length + 1}`,
        items: [...cart.value],
        time: new Date().toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' }),
        total: totalBelanja.value
    });
    cart.value = [];
    payAmount.value = 0;
    setPaymentMethod('Cash');
    isMobileCartOpen.value = false;
    Swal.fire({ toast: true, position: 'top-end', icon: 'info', title: 'Pesanan ditunda!', showConfirmButton: false, timer: 1500 });
};

const resumeOrder = (order) => {
    if (cart.value.length > 0) {
        Swal.fire({
            title: 'Timpa Keranjang?',
            text: "Ada barang di keranjang saat ini. Lanjutkan memuat pesanan tertunda?",
            icon: 'warning',
            showCancelButton: true,
            confirmButtonText: 'Ya, Timpa!'
        }).then((res) => {
            if (res.isConfirmed) processResume(order);
        });
    } else {
        processResume(order);
    }
};

const processResume = (order) => {
    cart.value = [...order.items];
    heldOrders.value = heldOrders.value.filter(h => h.id !== order.id);
    showHeldModal.value = false;
    if (window.innerWidth < 1024) isMobileCartOpen.value = true;
};

const pajakPersen = ref(0); 
const subTotalBelanja = computed(() => cart.value.reduce((total, item) => total + (item.price * item.qty), 0));
const nilaiPajak = computed(() => (pajakPersen.value / 100) * subTotalBelanja.value);
const totalBelanja = computed(() => Math.round((subTotalBelanja.value + nilaiPajak.value) / 100) * 100);
const kembalian = computed(() => payAmount.value - totalBelanja.value);

const setPaymentMethod = (method) => {
    paymentMethod.value = method;
    payAmount.value = method !== 'Cash' ? totalBelanja.value : 0;
};

// PROSES CHECKOUT
// State
const isProcessingCheckout = ref(false);

const executeCheckout = async() => {
    if (isProcessingCheckout.value) return;

    isProcessingCheckout.value = true;
    const payloadItems = cart.value.map(item => ({ product_id: item.id, kuantitas: item.qty }));
    try {
        const response = await api.post('/checkout', {
            session_id: currentSession.value.id,
            items: payloadItems,
            nominal_bayar: payAmount.value,
            metode_bayar: paymentMethod.value
        });

        lastTransaction.value = {
            invoice: response.data.invoice, 
            cart: [...cart.value],
            total: response.data.tagihan, 
            pay: payAmount.value,
            return: response.data.kembali, 
            method: paymentMethod.value,
            subtotal: subTotalBelanja.value, 
            pajak: nilaiPajak.value,         
            date: new Date().toLocaleString('id-ID', { year: '2-digit', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' }).replace(/\//g, '.')
        };
        
        showQrisModal.value = false;
        isMobileCartOpen.value = false; // Tutup drawer HP
        showReceipt.value = true;
        cart.value = [];
        payAmount.value = 0;
        paymentMethod.value = 'Cash';

        fetchProducts();
        nextTick(() => { if (searchInput.value) searchInput.value.focus(); });
    } catch (error) {
        Swal.fire('Gagal!', error.response?.data?.error || 'Koneksi terputus', 'error');
    } finally {
        isProcessingCheckout.value = false;
    }
};

const formatInputRupiah = (event) => {
    let rawValue = event.target.value.replace(/\D/g, '');
    payAmount.value = rawValue ? parseInt(rawValue, 10) : 0;
    event.target.value = payAmount.value === 0 ? '' : payAmount.value.toLocaleString('id-ID');
};

const processCheckout = () => {
    if (payAmount.value < totalBelanja.value) {
        Swal.fire({ icon: 'error', title: 'Uang Kurang!', text: `Kurang Rp ${(totalBelanja.value - payAmount.value).toLocaleString('id-ID')}` });
        return;
    }
    paymentMethod.value === 'QRIS' ? showQrisModal.value = true : executeCheckout();
};

const showClosingModal = ref(false);
const pecahan = ref({
    p100k: 0, p50k: 0, p20k: 0, p10k: 0, p5k: 0, p2k: 0, p1k: 0,
    p500: 0, p200: 0, p100: 0, p50: 0, p25: 0
});

const totalUangFisik = computed(() => {
    return (pecahan.value.p100k * 100000) + (pecahan.value.p50k * 50000) + 
           (pecahan.value.p20k * 20000) + (pecahan.value.p10k * 10000) + 
           (pecahan.value.p5k * 5000) + (pecahan.value.p2k * 2000) + 
           (pecahan.value.p1k * 1000) + (pecahan.value.p500 * 500) + 
           (pecahan.value.p200 * 200) + (pecahan.value.p100 * 100) + 
           (pecahan.value.p50 * 50) + (pecahan.value.p25 * 25);
});

const handleClosing = async () => {
    try {
        const res = await api.post(`/pos/close-session/${currentSession.value.id}`, {
            total_aktual: totalUangFisik.value,
            pecahan: pecahan.value
        });
        Swal.fire('Closing Berhasil!', 'Struk closing akan dicetak.', 'success');
        lastClosingData.value = res.data; 
        showClosingModal.value = false;
        showReceiptClosing.value = true;
    } catch (error) {
        Swal.fire('Gagal Closing', error.response?.data?.error, 'error');
    }
};

const printReceipt = () => window.print();

const printClosing = () => {
    const printContent = document.getElementById('print-closing').innerHTML;
    const originalContent = document.body.innerHTML;
    document.body.innerHTML = printContent;
    window.print();
    document.body.innerHTML = originalContent;
    window.location.reload(); 
};

const finishClosing = () => router.push('/absensi');
const goToDashboard = () => currentUser.value.role === 'owner' ? router.push('/dashboard') : router.push('/absensi');

const logout = () => {
    Swal.fire({
        title: 'Akhiri Shift?',
        text: "Hitung uang laci (Cash Count) sebelum tutup shift.",
        icon: 'question',
        showCancelButton: true,
        confirmButtonColor: '#2563eb',
        confirmButtonText: 'Ya, Tutup Shift'
    }).then((result) => {
        if (result.isConfirmed) {
            Object.keys(pecahan.value).forEach(k => pecahan.value[k] = 0);
            showClosingModal.value = true;
        }
    });
};
</script>

<template>
    <div class="h-[100dvh] flex flex-col bg-slate-100 font-sans overflow-hidden">
        
        <div class="p-2 md:p-4 shrink-0 pb-0">
            <div class="bg-indigo-900 rounded-[20px] md:rounded-[24px] shadow-lg flex flex-col md:flex-row overflow-hidden border border-indigo-800">
                <div @click="goToDashboard" class="p-3 md:p-4 bg-indigo-950 md:w-64 flex flex-col justify-center items-center cursor-pointer hover:bg-indigo-900 transition-colors border-b md:border-b-0 md:border-r border-indigo-800 shrink-0">
                    <div class="font-black text-xl md:text-2xl tracking-tight text-white flex items-center gap-2">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 md:w-8 md:h-8 text-indigo-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" /></svg>
                        NEXA POS
                    </div>
                    <div class="text-center text-[9px] md:text-[10px] uppercase mt-1 md:mt-2 tracking-widest font-black px-3 py-1 rounded-lg" :class="currentUser.role === 'owner' ? 'text-amber-400 bg-amber-400/10' : 'text-emerald-400 bg-emerald-400/10'">
                        {{ currentUser.role === 'owner' ? 'Ke Dashboard' : '🟢 SHIFT AKTIF' }}
                    </div>
                </div>

                <div class="flex-1 flex flex-col min-w-0">
                    <div class="bg-rose-600 text-white text-[9px] md:text-[10px] font-black py-1.5 text-center tracking-[0.2em] overflow-hidden whitespace-nowrap shadow-inner flex w-full">
                        <div class="animate-marquee min-w-full">
                            INFO KASIR: PASTIKAN SCAN BARCODE DENGAN BENAR | SELALU TAWARKAN PRODUK PROMO KEPADA KONSUMEN SEBELUM CHECKOUT
                        </div>
                    </div>

                    <div class="flex flex-wrap md:flex-nowrap justify-between items-center px-3 md:px-5 py-2 md:py-3 flex-1 bg-indigo-900 gap-2 md:gap-4">
                        
                        <div class="flex items-center gap-2 md:gap-4 text-xs md:text-sm font-medium text-indigo-100 min-w-0">
                            <div class="flex items-center gap-1.5 md:gap-2 bg-indigo-950/50 px-2.5 py-1.5 rounded-lg md:rounded-xl border border-indigo-700/50 truncate max-w-[120px] md:max-w-[200px]">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 md:w-4 md:h-4 text-indigo-400 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" /></svg>
                                <span class="font-black text-white uppercase tracking-wider text-[10px] md:text-xs truncate" :title="currentUser.name">{{ currentUser.name.split(' ')[0] }}</span>
                            </div>
                            
                            <div class="flex items-center gap-3 sm:gap-5 border-l border-indigo-800 pl-3 sm:pl-5">
                                <span class="flex items-center gap-1 sm:gap-2 opacity-90 text-[10px] uppercase tracking-widest font-black whitespace-nowrap">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 md:w-4 md:h-4 text-indigo-400 hidden sm:block" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                                    <span>Shift:</span><strong class="text-white text-xs">1</strong>
                                </span>
                                <span class="flex items-center gap-1 sm:gap-2 opacity-90 text-[10px] uppercase tracking-widest font-black whitespace-nowrap">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 md:w-4 md:h-4 text-indigo-400 hidden sm:block" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" /></svg>
                                    <span>POS:</span><strong class="text-white text-xs">{{ currentSession?.station_number || '00' }}</strong>
                                </span>
                            </div>
                        </div>

                        <div class="flex items-center gap-2 md:gap-3 shrink-0">
                            <div class="font-mono text-xs lg:text-sm font-black text-amber-400 tracking-wider bg-indigo-950/50 px-2 md:px-3 py-1.5 rounded-lg md:rounded-xl border border-indigo-800/50 hidden lg:block">
                                {{ currentTime }}
                            </div>
                            <button @click="logout" class="bg-rose-500 hover:bg-rose-600 text-white p-2 px-3 md:px-4 md:py-2 rounded-lg md:rounded-xl transition-all shadow-lg shadow-rose-900/50 active:scale-95 flex items-center gap-2 border border-rose-400 whitespace-nowrap">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 md:w-5 md:h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" /></svg>
                                <span class="text-[9px] md:text-[10px] font-black uppercase tracking-widest">Tutup Shift</span>
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div class="flex-1 flex overflow-hidden p-2 md:p-4 pt-2 md:pt-4 gap-4 relative">
            
            <div class="flex-1 flex flex-col min-h-0 w-full lg:w-8/12 xl:w-9/12 transition-all duration-300" 
                 :class="cart.length > 0 ? 'pb-20 lg:pb-0' : 'pb-0'">
                <div class="flex gap-2 md:gap-3 shrink-0 mb-3 md:mb-4 items-stretch h-12 md:h-14">
                    <div class="relative flex-1 group h-full">
                        <div class="absolute inset-y-0 left-0 pl-4 md:pl-5 flex items-center pointer-events-none">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 md:w-6 md:h-6 text-slate-400 group-focus-within:text-indigo-600 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
                        </div>
                        <input
                            ref="searchInput"
                            type="text"
                            v-model="searchQuery"
                            @keydown.enter.prevent="handleBarcodeScan"
                            placeholder="Cari atau Scan..."
                            class="w-full h-full pl-12 md:pl-14 pr-4 rounded-[16px] md:rounded-[20px] border-2 border-slate-200 focus:border-indigo-600 shadow-sm text-slate-800 font-bold bg-white text-xs md:text-sm transition-all outline-none"
                        >
                    </div>
                    
                    <button v-if="heldOrders.length > 0" @click="showHeldModal = true" class="lg:hidden shrink-0 bg-amber-100 hover:bg-amber-500 text-amber-600 hover:text-white px-4 md:px-5 rounded-[16px] md:rounded-[20px] transition-all border-2 border-amber-200 hover:border-amber-500 flex items-center justify-center shadow-sm relative animate-pulse" title="Pesanan Tertunda">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                        <span class="absolute -top-1.5 -right-1.5 flex h-5 w-5 items-center justify-center rounded-full bg-rose-500 text-[10px] font-black text-white shadow-md">{{ heldOrders.length }}</span>
                    </button>
                    
                    <button @click="startScanner" class="shrink-0 bg-indigo-100 hover:bg-indigo-600 text-indigo-600 hover:text-white px-4 md:px-5 rounded-[16px] md:rounded-[20px] transition-all border-2 border-indigo-200 hover:border-indigo-600 flex items-center justify-center shadow-sm h-full" title="Scan Kamera">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 md:w-7 md:h-7" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" /><path stroke-linecap="round" stroke-linejoin="round" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
                    </button>
                </div>

                <div v-if="filteredProducts.length === 0" class="flex-1 flex flex-col items-center justify-center bg-white/50 rounded-[24px] md:rounded-[32px] border-2 border-dashed border-slate-300">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-16 h-16 md:w-24 md:h-24 text-slate-300 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" /></svg>
                    <p class="text-slate-400 font-black text-sm md:text-lg uppercase tracking-widest text-center">Produk Tidak Ditemukan</p>
                </div>

                <div v-else class="flex-1 overflow-y-auto custom-scrollbar pr-2 pb-4">
                    <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 xl:grid-cols-5 gap-3 md:gap-4">
                        <div v-for="product in filteredProducts" :key="product.id" @click="addToCart(product)"
                            class="bg-white rounded-[16px] md:rounded-[24px] shadow-sm hover:shadow-xl hover:ring-2 hover:ring-indigo-500 transition-all duration-200 overflow-hidden cursor-pointer border border-slate-100 group flex flex-col transform hover:-translate-y-1">
                            
                            <div class="relative pt-2 px-2 md:pt-3 md:px-3">
                                <div class="bg-slate-50 rounded-xl md:rounded-2xl overflow-hidden aspect-square flex items-center justify-center border border-slate-100">
                                    <img 
                                        :src="getImageUrl(product.image) || 'https://placehold.co/150x150?text=No+Foto'"
                                        :alt="product.name"
                                        class="w-full h-full object-contain mix-blend-multiply p-3 md:p-4 group-hover:scale-110 transition-transform duration-300"
                                    >
                                </div>
                                <div v-if="product.stock <= 0" 
                                     class="absolute top-3 right-3 md:top-5 md:right-5 text-[8px] md:text-[9px] font-black px-1.5 md:px-2 py-0.5 md:py-1 rounded-md shadow-sm tracking-widest bg-rose-500 text-white animate-pulse uppercase">
                                    HABIS
                                </div>
                            </div>

                            <div class="p-2 md:p-4 flex flex-col flex-1 text-center justify-between gap-1 md:gap-2">
                                <h2 class="font-bold text-slate-700 text-[10px] md:text-[11px] line-clamp-2 leading-tight uppercase" :title="product.name">{{ product.name }}</h2>
                                <p class="text-indigo-700 font-black text-xs md:text-sm">Rp {{ product.price.toLocaleString('id-ID') }}</p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div :class="isMobileCartOpen ? 'fixed inset-0 z-[100] bg-slate-900/40 backdrop-blur-sm flex justify-end transition-all' : 'hidden lg:flex lg:relative w-4/12 xl:w-3/12'">
                <div :class="isMobileCartOpen ? 'w-[85%] sm:w-[360px] h-full animate-slide-in-right' : 'w-full h-full'" 
                     class="bg-white lg:rounded-[32px] shadow-2xl border-l lg:border border-slate-100 flex flex-col shrink-0 overflow-hidden">
                    
                    <div class="lg:hidden p-4 bg-indigo-900 text-white flex justify-between items-center shrink-0">
                        <h2 class="font-black tracking-widest uppercase text-sm flex items-center gap-2">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z" /></svg>
                            Keranjang
                        </h2>
                        <button @click="isMobileCartOpen = false" class="bg-white/20 p-2 rounded-xl text-white hover:bg-rose-500 transition-colors">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
                        </button>
                    </div>

                    <div v-auto-animate class="p-4 md:p-5 border-b border-slate-100 bg-slate-50/80 hidden lg:flex justify-between items-center shrink-0">
                        <h2 class="text-sm md:text-base font-black text-slate-800 flex items-center gap-2 uppercase tracking-widest">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z" /></svg>
                            Keranjang
                        </h2>
                        <div class="flex gap-1.5 md:gap-2">
                            <button @click="showHeldModal = true" class="p-2 bg-amber-50 hover:bg-amber-100 text-amber-600 rounded-xl transition-colors relative" title="Lihat Pesanan Tertunda">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 md:w-5 md:h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                                <span v-if="heldOrders.length" class="absolute -top-1 -right-1 flex h-4 w-4 items-center justify-center rounded-full bg-rose-500 text-[8px] font-black text-white">{{ heldOrders.length }}</span>
                            </button>
                            <button @click="holdTransaction" :disabled="cart.length===0" class="p-2 bg-indigo-50 hover:bg-indigo-100 text-indigo-600 rounded-xl transition-colors disabled:opacity-50" title="Hold Pesanan Ini">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 md:w-5 md:h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M10 9v6m4-6v6m7-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                            </button>
                            <button @click="clearCart" :disabled="cart.length===0" class="p-2 bg-rose-50 hover:bg-rose-100 text-rose-600 rounded-xl transition-colors disabled:opacity-50" title="Kosongkan Keranjang">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 md:w-5 md:h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h14" /></svg>
                            </button>
                        </div>
                    </div>

                    <div class="p-3 md:p-4 flex-1 overflow-y-auto bg-white custom-scrollbar min-h-0">
                        <div v-if="cart.length > 0" class="flex gap-2 lg:hidden mb-4">
                            <button @click="holdTransaction" class="flex-1 py-2 bg-amber-50 text-amber-600 rounded-lg text-[9px] font-black uppercase tracking-widest border border-amber-100">Hold Order</button>
                            <button @click="clearCart" class="flex-1 py-2 bg-rose-50 text-rose-600 rounded-lg text-[9px] font-black uppercase tracking-widest border border-rose-100">Bersihkan</button>
                        </div>

                        <div v-if="cart.length === 0" class="h-full flex flex-col items-center justify-center opacity-30 py-10">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-12 h-12 md:w-16 md:h-16 mb-4 text-slate-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" /></svg>
                            <p class="text-slate-600 font-black text-xs md:text-sm uppercase tracking-widest">Keranjang Kosong</p>
                        </div>

                        <div v-for="item in cart" :key="item.id" class="flex flex-col mb-3 p-3 bg-slate-50/50 rounded-xl md:rounded-2xl border border-slate-100 shadow-sm hover:border-indigo-200 transition-colors">
                            <div class="flex justify-between items-start mb-2">
                                <h3 class="font-bold text-[10px] md:text-[11px] text-slate-800 leading-tight pr-2 line-clamp-2 uppercase">{{ item.name }}</h3>
                                <div class="font-black text-[11px] md:text-xs text-indigo-700 whitespace-nowrap">Rp {{ (item.price * item.qty).toLocaleString('id-ID') }}</div>
                            </div>
                            <div class="flex justify-between items-center">
                                <p class="text-[9px] md:text-[10px] font-bold text-slate-400">@ Rp {{ item.price.toLocaleString('id-ID') }}</p>
                                <div class="flex items-center bg-white rounded-lg md:rounded-xl p-1 border border-slate-200 shadow-sm">
                                    <button @click="decreaseQty(item)" class="w-6 h-6 md:w-7 md:h-7 flex items-center justify-center rounded-md md:rounded-lg text-slate-400 hover:bg-rose-50 hover:text-rose-600 font-black transition-colors">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M20 12H4" /></svg>
                                    </button>
                                    <input type="number" v-model.number="item.qty" @change="validateQty(item)" class="w-8 md:w-10 text-center text-xs font-black text-slate-800 bg-transparent border-none focus:ring-0 p-0 [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none">
                                    <button @click="addToCart(item)" class="w-6 h-6 md:w-7 md:h-7 flex items-center justify-center rounded-md md:rounded-lg text-slate-400 hover:bg-indigo-50 hover:text-indigo-600 font-black transition-colors">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" /></svg>
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="p-3 md:p-5 bg-white border-t border-slate-100 shadow-[0_-10px_20px_-10px_rgba(0,0,0,0.05)] shrink-0 z-10 lg:rounded-b-[32px]">
                        
                        <div class="mb-3 md:mb-4">
                            <span class="font-black text-[8px] md:text-[9px] text-slate-400 block mb-1.5 md:mb-2 uppercase tracking-widest text-center">Metode Pembayaran</span>
                            <div class="grid grid-cols-3 gap-1.5 md:gap-2">
                                <button @click="setPaymentMethod('Cash')" :class="paymentMethod === 'Cash' ? 'bg-indigo-600 text-white shadow-md shadow-indigo-200 border-indigo-600' : 'bg-slate-50 text-slate-500 border-slate-200 hover:bg-slate-100'" class="py-2.5 rounded-lg md:rounded-xl font-black text-[9px] md:text-[10px] uppercase transition-all flex flex-col items-center gap-1.5 border">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 md:w-5 md:h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z" /></svg>
                                    Tunai
                                </button>
                                <button @click="setPaymentMethod('QRIS')" :class="paymentMethod === 'QRIS' ? 'bg-indigo-600 text-white shadow-md shadow-indigo-200 border-indigo-600' : 'bg-slate-50 text-slate-500 border-slate-200 hover:bg-slate-100'" class="py-2.5 rounded-lg md:rounded-xl font-black text-[9px] md:text-[10px] uppercase transition-all flex flex-col items-center gap-1.5 border">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 md:w-5 md:h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm14 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z" /></svg>
                                    QRIS
                                </button>
                                <button @click="setPaymentMethod('Debit')" :class="paymentMethod === 'Debit' ? 'bg-indigo-600 text-white shadow-md shadow-indigo-200 border-indigo-600' : 'bg-slate-50 text-slate-500 border-slate-200 hover:bg-slate-100'" class="py-2.5 rounded-lg md:rounded-xl font-black text-[9px] md:text-[10px] uppercase transition-all flex flex-col items-center gap-1.5 border">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 md:w-5 md:h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z" /></svg>
                                    Debit
                                </button>
                            </div>
                        </div>

                        <div class="space-y-3 md:space-y-4 mb-4 md:mb-5">
                            <div class="flex justify-between items-end border-b border-dashed border-slate-200 pb-1.5 md:pb-2">
                                <span class="font-black text-[9px] md:text-[10px] text-slate-400 uppercase tracking-widest">Total Tagihan</span>
                                <span class="text-2xl md:text-3xl font-black text-indigo-800 leading-none tracking-tighter">Rp {{ totalBelanja.toLocaleString('id-ID') }}</span>
                            </div>

                            <div class="flex justify-between items-center bg-slate-50 p-2 md:p-2.5 rounded-xl md:rounded-2xl border-2 border-slate-100 focus-within:border-indigo-500 transition-all">
                                <span class="font-black text-[9px] md:text-[10px] text-slate-600 uppercase tracking-widest pl-2">Bayar</span>
                                <div class="relative flex-1 ml-2 md:ml-4">
                                    <span class="absolute left-2.5 md:left-3 top-1/2 -translate-y-1/2 text-slate-400 font-bold text-xs md:text-sm italic">Rp</span>
                                    <input
                                        type="text"
                                        :value="payAmount === 0 ? '' : payAmount.toLocaleString('id-ID')"
                                        @input="formatInputRupiah"
                                        :disabled="paymentMethod !== 'Cash'"
                                        :class="paymentMethod !== 'Cash' ? 'bg-slate-200/50 text-slate-400 cursor-not-allowed border-transparent' : 'bg-white text-slate-900 border-slate-200 shadow-sm'"
                                        class="w-full text-right text-base md:text-lg font-black rounded-lg md:rounded-xl py-1.5 md:py-2 pl-8 pr-3 transition-all outline-none border"
                                        placeholder="0">
                                </div>
                            </div>

                            <div class="flex justify-between items-center px-1">
                                <span class="font-black text-[9px] md:text-[10px] text-slate-400 uppercase tracking-widest">Kembali</span>
                                <span class="text-lg md:text-xl font-black" :class="kembalian >= 0 ? 'text-emerald-500' : 'text-rose-500'">
                                    Rp {{ kembalian.toLocaleString('id-ID') }}
                                </span>
                            </div>
                        </div>

                        <button @click="processCheckout" :disabled="cart.length === 0 || payAmount < totalBelanja || isProcessingCheckout"
                            class="w-full bg-emerald-500 hover:bg-emerald-600 text-white font-black py-3 md:py-4 px-4 rounded-xl md:rounded-2xl transition-all flex justify-center items-center gap-2 md:gap-3 disabled:opacity-50 disabled:cursor-not-allowed shadow-xl shadow-emerald-200 hover:shadow-emerald-300 active:scale-95 text-xs md:text-sm uppercase tracking-widest">
                            <template v-if = "isProcessingCheckout">
                                <div class="w-4 h-4 md:w-5 md:h-5 border-2 border-white/30 border-t-white rounded-full animate-spin">
                                </div>
                                Memproses...
                            </template>
                            <template v-else>
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 md:w-5 md:h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" /></svg>
                                Proses Bayar
                            </template>
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <div v-if="cart.length > 0" class="lg:hidden fixed bottom-0 left-0 right-0 p-3 bg-white/90 backdrop-blur-md border-t border-slate-200 z-40 shadow-[0_-10px_20px_rgba(0,0,0,0.05)]">
            <button @click="isMobileCartOpen = true" class="w-full bg-indigo-600 text-white p-3 md:p-4 rounded-xl shadow-xl shadow-indigo-200 flex justify-between items-center active:scale-95 transition-all">
                <div class="flex items-center gap-2 md:gap-3">
                    <span class="bg-white text-indigo-600 font-black px-2 md:px-3 py-1 rounded-lg text-xs md:text-sm">{{ cart.length }} Item</span>
                    <span class="font-bold text-[10px] md:text-xs uppercase tracking-widest hidden sm:block">Lihat Keranjang</span>
                </div>
                <div class="font-black text-sm md:text-lg tracking-tight">
                    Rp {{ totalBelanja.toLocaleString('id-ID') }} ➔
                </div>
            </button>
        </div>

        <div v-if="showHeldModal" class="fixed inset-0 bg-slate-900/80 flex items-center justify-center z-[150] p-4 backdrop-blur-sm">
            <div class="bg-white p-6 md:p-8 rounded-[32px] shadow-2xl w-full max-w-xl max-h-[80vh] flex flex-col">
                <div class="flex justify-between items-center mb-6">
                    <h2 class="text-lg md:text-xl font-black text-slate-800 uppercase tracking-widest flex items-center gap-2">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 md:w-6 md:h-6 text-amber-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                        Pesanan Tertunda
                    </h2>
                    <button @click="showHeldModal = false" class="p-2 bg-slate-100 hover:bg-rose-100 text-slate-400 hover:text-rose-600 rounded-xl transition-colors">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
                    </button>
                </div>
                
                <div class="flex-1 overflow-y-auto custom-scrollbar pr-2">
                    <div v-if="heldOrders.length === 0" class="text-center py-10 text-slate-400 font-bold text-[10px] md:text-xs uppercase tracking-widest italic">
                        Tidak ada pesanan yang ditunda.
                    </div>
                    <div v-for="order in heldOrders" :key="order.id" class="p-3 md:p-4 bg-slate-50 border border-slate-200 rounded-xl md:rounded-2xl mb-3 flex justify-between items-center group hover:border-indigo-300 transition-colors">
                        <div>
                            <p class="font-black text-xs md:text-sm text-slate-800 uppercase">{{ order.customer }}</p>
                            <p class="text-[9px] md:text-[10px] font-bold text-slate-400 mt-0.5 md:mt-1">Jam: {{ order.time }} | {{ order.items.length }} Item</p>
                            <p class="text-indigo-600 font-black mt-1 text-xs md:text-sm">Rp {{ order.total.toLocaleString('id-ID') }}</p>
                        </div>
                        <button @click="resumeOrder(order)" class="bg-indigo-100 text-indigo-600 hover:bg-indigo-600 hover:text-white px-3 md:px-4 py-2 rounded-lg md:rounded-xl font-black text-[9px] md:text-[10px] uppercase tracking-widest transition-all">
                            Lanjutkan
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <div v-if="showScanner" class="fixed inset-0 bg-slate-900/90 flex items-center justify-center z-[150] p-4 backdrop-blur-sm">
            <div class="bg-white rounded-[24px] md:rounded-[32px] shadow-2xl w-full max-w-md overflow-hidden flex flex-col">
                <div class="p-4 md:p-6 border-b border-slate-50 flex justify-between items-center bg-slate-50/50">
                    <h2 class="text-base md:text-lg font-black text-slate-800 uppercase tracking-widest flex items-center gap-2">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm14 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z" /></svg>
                        Scan Barcode
                    </h2>
                    <button @click="stopScanner" class="p-2 rounded-xl bg-slate-100 text-slate-400 hover:text-rose-500 hover:bg-rose-50 transition-all">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                    </button>
                </div>
                <div class="p-4 md:p-6 bg-black relative">
                    <div id="reader-kasir" class="w-full rounded-2xl overflow-hidden border-2 border-slate-700"></div>
                    <p class="text-[9px] md:text-[10px] font-bold text-white/50 text-center mt-3 md:mt-4 uppercase tracking-widest animate-pulse">Arahkan kamera ke barcode produk</p>
                </div>
            </div>
        </div>

        <div v-if="showQrisModal" class="fixed inset-0 bg-slate-900/80 flex items-center justify-center z-[150] p-4 backdrop-blur-sm">
            <div class="bg-white p-6 md:p-8 rounded-[32px] shadow-2xl w-full max-w-sm text-center flex flex-col border-t-8 border-indigo-600">
                <div class="w-12 h-12 md:w-16 md:h-16 bg-indigo-100 text-indigo-600 rounded-full flex items-center justify-center mx-auto mb-3 md:mb-4">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 md:w-8 md:h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm14 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z" /></svg>
                </div>
                <h2 class="text-lg md:text-xl font-black text-slate-900 uppercase tracking-widest mb-1">Bayar via QRIS</h2>
                <p class="text-slate-500 text-[9px] md:text-[10px] font-bold uppercase tracking-widest mb-4 md:mb-6">Minta pelanggan scan kode ini</p>
                
                <div class="bg-white p-3 rounded-2xl border-2 border-dashed border-slate-300 w-full mb-4 md:mb-6 shadow-inner flex justify-center items-center min-h-[180px] md:min-h-[200px]">
                    <img src="https://placehold.co/300x300?text=QRIS+TOKO" alt="QRIS Toko" class="w-full h-full max-h-40 md:max-h-48 object-contain mx-auto rounded-xl">
                </div>
                
                <div class="bg-indigo-50 text-indigo-900 p-3 md:p-4 rounded-2xl mb-6 md:mb-8 border border-indigo-100">
                    <span class="block text-[8px] md:text-[9px] font-black uppercase tracking-[0.2em] opacity-70 mb-1">Total Tagihan</span>
                    <span class="text-2xl md:text-3xl font-black tracking-tighter">Rp {{ totalBelanja.toLocaleString('id-ID') }}</span>
                </div>
                
                <div class="flex gap-2 md:gap-3">
                    <button @click="showQrisModal = false" :disabled="isProcessingCheckout" class="flex-1 bg-slate-100 py-3 md:py-4 rounded-xl font-black text-[9px] md:text-[10px] uppercase tracking-widest text-slate-500 hover:bg-slate-200 transition-colors disabled:opacity-50">Batal</button>
                    
                    <button @click="executeCheckout" :disabled="isProcessingCheckout" class="flex-1 bg-indigo-600 py-3 md:py-4 rounded-xl font-black text-[9px] md:text-[10px] uppercase tracking-widest text-white hover:bg-indigo-700 transition-colors shadow-lg shadow-indigo-200 disabled:opacity-50 flex items-center justify-center gap-2">
                        <template v-if="isProcessingCheckout">
                            <div class="w-3 h-3 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                            Proses...
                        </template>
                        <template v-else>
                            Lunas
                        </template>
                    </button>
                </div>
            </div>
        </div>

        <div v-if="showReceipt" class="fixed inset-0 bg-slate-950/90 flex items-center justify-center z-[200] p-4 backdrop-blur-sm no-print">
            <div class="bg-white p-6 md:p-8 rounded-[32px] md:rounded-[40px] shadow-2xl w-full max-w-sm flex flex-col max-h-[90vh] border-[6px] md:border-[8px] border-slate-800">
                <div class="overflow-y-auto custom-scrollbar bg-white p-2 mx-auto" id="print-area" style="width: 58mm;">
                    <div class="text-center mb-4 font-mono leading-none">
                        <h2 class="font-black text-sm uppercase tracking-tighter mb-1 italic">{{ currentSession?.store?.nama_toko || 'ARZU STORE' }}</h2>
                        <p class="text-[8px] font-bold uppercase tracking-widest opacity-80">{{ currentSession?.store?.alamat || 'JAKARTA, INDONESIA' }}</p>
                    </div>
                    
                    <div class="border-y border-black py-1.5 text-center font-black mb-3 font-mono text-[9px] tracking-[0.2em] uppercase bg-slate-100">Struk Belanja</div>
                    
                    <div class="mb-3 text-[8px] font-bold font-mono uppercase space-y-0.5">
                        <div class="flex justify-between"><span>WAKTU:</span><span>{{ lastTransaction?.date }}</span></div>
                        <div class="flex justify-between"><span>KASIR:</span><span>{{ currentUser.name.split(' ')[0] }} / POS-{{ currentSession?.station_number }}</span></div>
                    </div>
                    
                    <div class="border-b border-black border-dashed mb-2"></div>
                    
                    <div v-for="item in lastTransaction?.cart" :key="item.id" class="mb-2 font-bold font-mono text-[9px] leading-tight uppercase">
                        <div class="truncate w-full pr-2">{{ item.name }}</div>
                        <div class="flex justify-between pl-2 text-[8px] mt-0.5">
                            <span>{{ item.qty }} x {{ item.price.toLocaleString('id-ID') }}</span>
                            <span class="font-black text-[9px]">{{ (item.price * item.qty).toLocaleString('id-ID') }}</span>
                        </div>
                    </div>
                    
                    <div class="border-t border-black border-dashed mt-2 pt-2"></div>
                    
                    <div class="flex justify-between font-black text-[11px] mb-2 font-mono uppercase italic">
                        <span>TOTAL:</span><span>Rp{{ lastTransaction?.total.toLocaleString('id-ID') }}</span>
                    </div>
                    
                    <div class="border-b border-black border-dashed mb-2"></div>
                    
                    <div class="flex justify-between mb-1 font-bold font-mono text-[8px] uppercase">
                        <span>BAYAR ({{ lastTransaction?.method }}):</span><span>Rp{{ lastTransaction?.pay.toLocaleString('id-ID') }}</span>
                    </div>
                    <div v-if="lastTransaction?.method === 'Cash'" class="flex justify-between font-black font-mono text-[9px] uppercase italic text-black">
                        <span>KEMBALI:</span><span>Rp{{ lastTransaction?.return?.toLocaleString('id-ID') }}</span>
                    </div>
                    
                    <div class="mt-5 text-[7px] font-bold text-center border-t border-black border-dashed pt-2 font-mono uppercase space-y-1">
                        <p>INV: {{ lastTransaction?.invoice }}</p>
                    </div>
                    <div class="text-center mt-4 font-black font-mono text-[8px] border-2 border-black p-1.5 uppercase">
                        Terima Kasih!<br>Barang tidak dapat ditukar.
                    </div>
                </div>
                
                <div class="mt-4 md:mt-6 flex flex-col gap-2 md:gap-3 shrink-0 no-print">
                    <button @click="printReceipt" class="w-full bg-indigo-600 text-white py-3 md:py-4 rounded-xl md:rounded-2xl font-black text-[9px] md:text-[10px] uppercase tracking-[0.2em] shadow-xl shadow-indigo-200 flex items-center justify-center gap-2 active:scale-95 transition-all">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 md:w-5 md:h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z" /></svg>
                        Cetak Struk
                    </button>
                    <button @click="showReceipt = false" class="w-full bg-slate-100 text-slate-500 py-3 md:py-3.5 rounded-xl md:rounded-2xl font-black text-[9px] md:text-[10px] uppercase tracking-widest hover:bg-slate-200 transition-all">Tutup</button>
                </div>
            </div>
        </div>

        <div v-if="showClosingModal" class="fixed inset-0 bg-slate-900/90 flex items-center justify-center z-[160] p-4 backdrop-blur-md">
            <div class="bg-white rounded-[24px] md:rounded-[32px] w-full max-w-2xl overflow-hidden shadow-2xl border-[6px] md:border-[8px] border-slate-800 flex flex-col max-h-[95vh] md:max-h-[90vh]">
                <div class="bg-slate-800 p-4 md:p-6 text-center shrink-0 relative">
                    <button @click="showClosingModal = false" class="absolute top-4 md:top-6 right-4 md:right-6 text-slate-400 hover:text-white transition-colors">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 md:w-6 md:h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
                    </button>
                    <h2 class="text-white font-black text-xl md:text-2xl tracking-tighter uppercase italic">Cash Count</h2>
                    <p class="text-slate-400 text-[9px] md:text-[10px] font-bold uppercase tracking-widest mt-1">Hitung Uang Fisik Laci Kasir</p>
                </div>

                <div class="p-4 md:p-6 grid grid-cols-2 md:grid-cols-3 gap-2 md:gap-3 overflow-y-auto custom-scrollbar flex-1 bg-slate-50">
                    <div v-for="denon in [
                        { label: '100.000', key: 'p100k' }, { label: '50.000', key: 'p50k' }, { label: '20.000', key: 'p20k' },
                        { label: '10.000', key: 'p10k' }, { label: '5.000', key: 'p5k' }, { label: '2.000', key: 'p2k' },
                        { label: '1.000', key: 'p1k' }, { label: '500', key: 'p500' }, { label: '200', key: 'p200' },
                        { label: '100', key: 'p100' }, { label: '50', key: 'p50' }
                    ]" :key="denon.key" class="bg-white p-2 md:p-3 rounded-xl md:rounded-2xl border-2 border-slate-200 focus-within:border-indigo-500 shadow-sm transition-colors">
                        <label class="text-[8px] md:text-[9px] font-black text-slate-400 uppercase tracking-widest block mb-0.5 md:mb-1">Pecahan {{ denon.label }}</label>
                        <input type="number" v-model.number="pecahan[denon.key]" class="w-full bg-transparent text-lg md:text-xl font-black text-slate-800 outline-none" placeholder="0" min="0">
                    </div>
                </div>

                <div class="p-4 md:p-6 bg-white border-t border-slate-200 shrink-0">
                    <div class="flex justify-between items-center mb-4 md:mb-5">
                        <span class="text-slate-500 text-[9px] md:text-[10px] font-black uppercase tracking-widest">Total Fisik:</span>
                        <span class="text-2xl md:text-3xl font-black text-indigo-700 tracking-tighter">Rp {{ totalUangFisik.toLocaleString('id-ID') }}</span>
                    </div>
                    <button @click="handleClosing" class="w-full bg-indigo-600 text-white py-3.5 md:py-4 rounded-xl md:rounded-2xl font-black text-[10px] md:text-xs tracking-widest shadow-xl shadow-indigo-200 hover:bg-indigo-700 active:scale-95 transition-all flex items-center justify-center gap-2 uppercase">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 md:w-5 md:h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                        Proses Tutup Shift
                    </button>
                </div>
            </div>
        </div>

        <div v-if="showReceiptClosing" class="fixed inset-0 bg-slate-900/90 flex items-center justify-center z-[200] p-4 backdrop-blur-sm">
            <div class="bg-white p-6 md:p-8 rounded-[32px] md:rounded-[40px] shadow-2xl w-full max-w-sm overflow-hidden border-t-8 border-indigo-600">
                <div id="print-closing" class="text-left font-mono text-[9px] md:text-[10px] leading-tight uppercase text-black bg-white p-2 mx-auto" style="width: 58mm;">
                    <div class="text-center mb-3">
                        <h2 class="font-black text-xs md:text-sm mb-1 italic">{{ currentSession?.store?.nama_toko || 'ARZU STORE' }}</h2>
                        <p class="font-bold text-[7px] md:text-[8px] tracking-widest border-y border-black py-1">CLOSING REPORT - POS {{ currentSession?.station_number }}</p>
                    </div>
            
                    <div class="flex justify-between mb-1 font-bold"><span>MULAI :</span><span>{{ lastClosingData?.start_time }}</span></div>
                    <div class="flex justify-between mb-1 font-bold"><span>SELESAI :</span><span>{{ lastClosingData?.end_time }}</span></div>
                    <div class="flex justify-between mb-2 font-bold"><span>KASIR :</span><span>{{ currentUser.name.split(' ')[0] }}</span></div>
            
                    <p class="border-b border-dashed border-black mb-2"></p>
                    <div class="flex justify-between font-bold"><span>SALES KOTOR :</span><span>{{ lastClosingData?.sales_gross?.toLocaleString('id-ID') || 0 }}</span></div>
                    <p class="border-b border-black my-1"></p>
                    <div class="flex justify-between font-black text-[10px] md:text-[11px]"><span>NET SALES :</span><span>{{ lastClosingData?.net_sales?.toLocaleString('id-ID') || 0 }}</span></div>
            
                    <p class="border-b border-dashed border-black my-2"></p>
                    <div class="flex justify-between font-bold"><span>MODAL AWAL :</span><span>{{ currentSession?.modal_awal?.toLocaleString('id-ID') || 0 }}</span></div>
                    <div class="flex justify-between font-bold"><span>SALES TUNAI :</span><span>{{ lastClosingData?.sales_cash?.toLocaleString('id-ID') || 0 }}</span></div>

                    <div class="border-t border-dashed border-black my-1"></div>
                    <div class="flex justify-between font-bold"><span>SALES NON-TUNAI :</span><span>{{ lastClosingData?.sales_non_tunai?.toLocaleString('id-ID') || 0 }}</span></div>
                    <div class="flex justify-between font-black text-[10px] md:text-[11px] mt-1"><span>TOTAL MASUK :</span><span>{{ lastClosingData?.total_expected?.toLocaleString('id-ID') || 0 }}</span></div>
            
                    <p class="border-b border-dashed border-black my-2"></p>
                    <div class="flex justify-between font-bold"><span>UANG FISIK :</span><span>{{ lastClosingData?.total_actual?.toLocaleString('id-ID') || 0 }}</span></div>
                    <div class="flex justify-between font-black text-[11px] md:text-[12px] mt-1" :class="lastClosingData?.selisih < 0 ? 'text-black' : 'text-black'">
                        <span>SELISIH :</span><span>{{ lastClosingData?.selisih?.toLocaleString('id-ID') || 0 }}</span>
                    </div>
            
                    <p class="border-b border-dashed border-black mt-4"></p>
                    <div class="text-center mt-2 font-bold text-[7px] md:text-[8px] tracking-widest">
                        <p>=== SHIFT SELESAI ===</p>
                        <p class="mt-2">DIPERIKSA OLEH : ___________</p>
                    </div>
                </div>
                <div class="mt-5 md:mt-6 flex gap-2 md:gap-3 no-print">
                    <button @click="printClosing" class="flex-1 bg-indigo-600 hover:bg-indigo-700 text-white py-3 md:py-3.5 rounded-xl md:rounded-2xl font-black text-[9px] md:text-[10px] uppercase tracking-widest shadow-lg transition-colors flex items-center justify-center gap-2">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z" /></svg>
                        Struk
                    </button>
                    <button @click="finishClosing" class="flex-1 bg-slate-100 hover:bg-slate-200 py-3 md:py-3.5 rounded-xl md:rounded-2xl font-black text-slate-600 text-[9px] md:text-[10px] uppercase tracking-widest transition-colors flex items-center justify-center gap-2">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" /></svg>
                        Pulang
                    </button>
                </div>
            </div>
        </div>

    </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
@media (min-width: 768px) { .custom-scrollbar::-webkit-scrollbar { width: 6px; } }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #94a3b8; }

/* CSS Animasi untuk Teks Berjalan */
.animate-marquee {
    display: inline-block;
    padding-left: 100%;
    animation: marquee 20s linear infinite;
}
@keyframes marquee {
    0% { transform: translateX(0); }
    100% { transform: translateX(-100%); }
}

/* Animasi Drawer Keranjang Mobile */
@keyframes slide-in-right {
    from { transform: translateX(100%); }
    to { transform: translateX(0); }
}
.animate-slide-in-right {
    animation: slide-in-right 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
}

@media print {
    body * { visibility: hidden; }
    #print-area, #print-area *, #print-closing, #print-closing * { visibility: visible; }
    #print-area, #print-closing { position: absolute; left: 0; top: 0; width: 58mm; padding: 0; margin: 0; }
    @page { margin: 0; }
    .no-print { display: none !important; }
}

:deep(.swal2-popup) {
    font-family: 'Inter', sans-serif !important;
    border-radius: 28px !important;
}
:deep(.swal2-confirm), :deep(.swal2-cancel) {
    border-radius: 14px !important;
    font-weight: 800 !important;
    text-transform: uppercase !important;
    letter-spacing: 0.5px !important;
    padding: 12px 24px !important;
}
</style>