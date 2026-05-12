<script setup>
import { ref, computed, onMounted, nextTick, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import api from '../api.js';
import Swal from 'sweetalert2';

const router = useRouter();

// --- SETUP USER & ROLE ---
const getUserInfo = () => {
    const token = localStorage.getItem('token');
    const role = localStorage.getItem('role') || 'kasir';
    let name = localStorage.getItem('name'); 

    if (token) {
        try {
            const payload = JSON.parse(atob(token.split('.')[1]));
            // 🚀 Kalau di localStorage kosong, coba ambil dari payload token
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

// Fungsi Jam Realtime
const currentTime = ref('');
let timer;

// State utama untuk keranjang, pembayaran, dan tampilan
const products = ref([]);
const isLoadingProducts = ref(true);
const cart = ref([]);
const heldOrders = ref([]);
const payAmount = ref(0);
const paymentMethod = ref('Cash');
const showReceipt = ref(false);
const showQrisModal = ref(false);
const lastTransaction = ref(null);

// 🚀 Update BaseURL kalau pakai IP Address HP (Ganti sesuai IP Mas)
const getImageUrl = (path) => {
    if (!path) return null;
    return `http://localhost:8080${path}`;
}

// Pencarian Produk & Barcode Scan
const searchQuery = ref('');
const searchInput = ref(null);

// Tarik Data dari Golang
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
        if (error.response && error.response.status === 401) {
            router.push('/login');
        }
    } finally {
        isLoadingProducts.value = false;
    }
};

onMounted(() => {
    fetchProducts();
    if (searchInput.value) searchInput.value.focus();
    
    timer = setInterval(() => {
        const now = new Date();
        currentTime.value = now.toLocaleString('id-ID', {
            day: '2-digit', month: '2-digit', year: 'numeric',
            hour: '2-digit', minute: '2-digit', second: '2-digit'
        }).replace(/\//g, '.');
    }, 1000);
});

onUnmounted(() => {
    clearInterval(timer);
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
        Swal.fire({
            icon: 'error',
            title: 'Stok Habis !',
            text: `Stok ${product.name} sudah kosong`,
            confirmButtonColor: '#2563eb'
        });
        return;
    }
    const existingItem = cart.value.find(item => item.id === product.id);
    if (existingItem) {
        if (existingItem.qty < product.stock) {
            existingItem.qty++;
        } else {
            Swal.fire({
                icon: 'warning',
                title: 'Stok Terbatas',
                text: 'Kuantitas tidak bisa melebihi stok tersedia!',
                confirmButtonColor: '#2563eb'
            });
        }
    } else {
        cart.value.push({ id: product.id, name: product.name, price: product.price, qty: 1 });
    }
};

const decreaseQty = (product) => {
    const existingItem = cart.value.find(item => item.id === product.id);
    if (existingItem) {
        if (existingItem.qty > 1) {
            existingItem.qty--;
        } else {
            cart.value = cart.value.filter(item => item.id !== product.id);
        }
    }
};

const validateQty = (item) => {
    if (!item.qty || item.qty < 1) item.qty = 1;
};

const clearCart = () => {
    if (cart.value.length === 0) return;

    Swal.fire({
        title: 'Batalkan Transaksi?',
        text: "Semua barang di keranjang akan dihapus!",
        icon: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#ef4444', // Merah
        cancelButtonColor: '#64748b',
        confirmButtonText: 'Ya, Hapus Semua!',
        cancelButtonText: 'Jangan Jadi'
    }).then((result) => {
        if (result.isConfirmed) {
            cart.value = [];
            payAmount.value = 0;
            setPaymentMethod('Cash');
            Swal.fire('Terhapus!', 'Keranjang sudah bersih kembali.', 'success');
        }
    });
};
const holdTransaction = () => {
    if (cart.value.length === 0) return;
    heldOrders.value.push({
        id: Date.now(),
        items: [...cart.value],
        time: new Date().toLocaleTimeString('id-ID'),
        total: totalBelanja.value
    });
    cart.value = [];
    payAmount.value = 0;
    setPaymentMethod('Cash');
};

const resumeOrder = (order) => {
    if (cart.value.length > 0 && !confirm('Keranjang saat ini akan diganti. Lanjutkan?')) return;
    cart.value = [...order.items];
    heldOrders.value = heldOrders.value.filter(h => h.id !== order.id);
};

const pajakPersen = ref(0); 

const subTotalBelanja = computed(() => {
    return cart.value.reduce((total, item) => total + (item.price * item.qty), 0);
});

const nilaiPajak = computed(() => {
    return (pajakPersen.value / 100) * subTotalBelanja.value;
});

const totalBelanja = computed(() => {
    const rawTotal = subTotalBelanja.value + nilaiPajak.value;
    return Math.round(rawTotal / 100) * 100; 
});

const kembalian = computed(() => {
    return payAmount.value - totalBelanja.value;
});

const setPaymentMethod = (method) => {
    paymentMethod.value = method;
    if (method !== 'Cash') {
        payAmount.value = totalBelanja.value;
    } else {
        payAmount.value = 0;
    }
};

const executeCheckout = async() => {
    const payloadItems = cart.value.map(item => ({
        product_id: item.id,
        kuantitas: item.qty
    }));

    try {
        const response = await api.post('/checkout', {
            items: payloadItems,
            nominal_bayar: payAmount.value
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
        showReceipt.value = true;
        cart.value = [];
        payAmount.value = 0;
        paymentMethod.value = 'Cash';

        fetchProducts();
        nextTick(() => { if (searchInput.value) searchInput.value.focus(); });
    } catch (error) {
        if (error.response && error.response.data.error) {
            alert("Transaksi Gagal: " + error.response.data.error);
        } else {
            alert("Koneksi ke server terputus");
        }
    }
};

const formattedPayAmount = computed({
    get() {
        return payAmount.value === 0 ? '' : payAmount.value.toLocaleString('id-ID');
    },
    set(newValue) {
        const cleanValue = String(newValue).replace(/\D/g, '');
        payAmount.value = cleanValue ? parseInt(cleanValue, 10) : 0;
    }
});

const formatInputRupiah = (event) => {
    let rawValue = event.target.value.replace(/\D/g, '');
    payAmount.value = rawValue ? parseInt(rawValue, 10) : 0;
    event.target.value = payAmount.value === 0 ? '' : payAmount.value.toLocaleString('id-ID');
};

const processCheckout = () => {
    if (payAmount.value < totalBelanja.value) {
        Swal.fire({
            icon: 'error',
            title: 'Uang Kurang!',
            text: `Nominal bayar masih kurang Rp ${(totalBelanja.value - payAmount.value).toLocaleString('id-ID')}`,
            confirmButtonColor: '#2563eb'
        });
        return;
    }
    if (paymentMethod.value === 'QRIS') {
        showQrisModal.value = true;
    } else {
        executeCheckout();
    }
};

const printReceipt = () => {
    window.print();
};

// 🚀 LOGIKA KEMBALI KE DASHBOARD (Diperbaiki)
const goToDashboard = () => {
    if (currentUser.value.role === 'owner') {
        router.push('/dashboard');
    } else {
        router.push('/absensi');
    }
};

// 🚀 LOGOUT CEPAT UNTUK KASIR
const logout = () => {
    Swal.fire({
        title: 'Akhiri Shift?',
        text: "Pastikan semua transaksi sudah selesai dicatat.",
        icon: 'question',
        showCancelButton: true,
        confirmButtonColor: '#2563eb',
        cancelButtonColor: '#64748b',
        confirmButtonText: 'Ya, Logout Sekarang',
        cancelButtonText: 'Batal'
    }).then((result) => {
        if (result.isConfirmed) {
            localStorage.clear();
            router.push('/login');
        }
    });
};
</script>

<template>
    <div class="p-3 md:p-5 bg-slate-50 min-h-screen flex flex-col font-sans">
        
        <div class="bg-blue-800 rounded-2xl shadow-lg mb-5 flex flex-col md:flex-row overflow-hidden border border-blue-900 shrink-0">
            <div @click="goToDashboard" class="p-3 bg-blue-900 md:w-56 flex flex-col justify-center items-center cursor-pointer hover:bg-blue-950 transition-colors border-b md:border-b-0 md:border-r border-blue-800">
                <div class="font-black text-xl tracking-wider text-white flex items-center gap-2">
                    <span class="text-2xl">🛒</span> POS KASIR
                </div>
                <div v-if="currentUser.role === 'owner'" class="text-center text-[10px] text-yellow-400 uppercase mt-1 tracking-widest font-black bg-blue-800/50 px-3 py-0.5 rounded-full">
                    ⬅️ Ke Dashboard
                </div>
                <div v-else class="text-center text-[10px] text-green-400 uppercase mt-1 tracking-widest font-black bg-blue-800/50 px-3 py-0.5 rounded-full">
                    🟢 Shift Aktif
                </div>
            </div>

            <div class="flex-1 flex flex-col">
                <div class="bg-red-600 text-white text-[11px] font-bold py-1.5 px-4 text-center tracking-widest overflow-hidden whitespace-nowrap shadow-inner">
                    <marquee scrollamount="5">INFO KASIR: PASTIKAN SCAN BARCODE DENGAN BENAR | SELALU TAWARKAN PRODUK BEST SELLER KEPADA KONSUMEN SEBELUM CHECKOUT</marquee>
                </div>

                <div class="flex justify-between items-center px-4 md:px-6 py-2.5 flex-1 bg-blue-800">
                    <div class="flex items-center gap-3 md:gap-5 text-xs md:text-sm font-medium text-blue-100">
                        <div class="flex items-center gap-2 bg-blue-900/60 px-3 py-1.5 rounded-lg border border-blue-700/50">
                            <span>👤</span>
                            <span class="font-black text-white uppercase tracking-wide">{{ currentUser.name }}</span>
                        </div>
                        
                        <div class="hidden sm:flex items-center gap-4 border-l border-blue-700 pl-4">
                            <span class="flex items-center gap-1 opacity-90">🕒 Shift: <strong class="text-white">1</strong></span>
                            <span class="flex items-center gap-1 opacity-90">💻 Station: <strong class="text-white">01</strong></span>
                        </div>

                        <button @click="logout" title="Akhiri Shift (Logout)" class="bg-red-500 hover:bg-red-600 text-white p-1.5 md:px-3 md:py-1.5 rounded-lg transition-colors shadow-sm ml-1 md:ml-2 active:scale-95 flex items-center gap-1">
                            <span class="text-sm">🚪</span>
                            <span class="text-[10px] font-black uppercase tracking-wider hidden md:block">Keluar</span>
                        </button>
                    </div>

                    <div class="font-mono text-sm md:text-lg font-black text-yellow-400 tracking-wider bg-blue-900/40 px-3 py-1 rounded-lg border border-blue-700/30">
                        {{ currentTime }}
                    </div>
                </div>
            </div>
        </div>

        <div class="flex flex-col lg:flex-row gap-5 flex-1 items-start">
            
            <div class="w-full lg:w-8/12 xl:w-9/12 flex flex-col gap-4">
                <div class="relative group">
                    <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                        <span class="text-gray-400 text-xl group-focus-within:text-blue-500 transition-colors">🔍</span>
                    </div>
                    <input
                        ref="searchInput"
                        type="text"
                        v-model="searchQuery"
                        @keydown.enter.prevent="handleBarcodeScan"
                        placeholder="Ketik nama produk atau Scan Barcode (Lalu tekan Enter)..."
                        class="w-full pl-12 pr-4 py-4 rounded-2xl border-0 shadow-md ring-1 ring-gray-200 focus:ring-4 focus:ring-blue-500/20 text-gray-800 font-bold bg-white text-lg transition-all outline-none"
                    >
                </div>

                <div v-if="filteredProducts.length === 0" class="flex-1 flex flex-col items-center justify-center bg-white rounded-2xl border-2 border-dashed border-gray-200 min-h-[50vh]">
                    <span class="text-6xl mb-4 grayscale opacity-50">📦</span>
                    <p class="text-gray-500 font-bold text-xl">Barang tidak ditemukan</p>
                    <p class="text-gray-400 text-sm mt-1">Coba kata kunci lain atau cek Master Produk.</p>
                </div>

                <div v-else class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-3 xl:grid-cols-4 gap-4 pb-4 overflow-y-auto custom-scrollbar max-h-[75vh] pr-2">
                    <div v-for="product in filteredProducts" :key="product.id" @click="addToCart(product)"
                        class="bg-white rounded-2xl shadow-sm hover:shadow-xl hover:ring-2 hover:ring-blue-500 transition-all duration-200 overflow-hidden cursor-pointer border border-gray-100 group flex flex-col transform hover:-translate-y-1">
                        
                        <div class="relative pt-3 px-3">
                            <div class="bg-slate-50 rounded-xl overflow-hidden aspect-square flex items-center justify-center">
                                <img 
                                    :src="getImageUrl(product.image) || 'https://placehold.co/150x150?text=No+Image'"
                                    :alt="product.name"
                                    class="w-full h-32 object-contain mix-blend-multiply p-3 group-hover:scale-110 transition-transform duration-300"
                                >
                            </div>
                            <div class="absolute top-5 right-5 text-[10px] font-black px-2 py-1 rounded-md shadow-sm"
                                :class="product.stock > 10 ? 'bg-green-100 text-green-700' : 'bg-red-500 text-white animate-pulse'">
                                Sisa: {{ product.stock }}
                            </div>
                        </div>

                        <div class="p-4 flex flex-col flex-1 text-center justify-between">
                            <h2 class="font-bold text-gray-700 text-xs md:text-sm line-clamp-2 leading-snug mb-2" :title="product.name">{{ product.name }}</h2>
                            <p class="text-blue-700 font-black text-base md:text-lg">Rp {{ product.price.toLocaleString('id-ID') }}</p>
                        </div>
                    </div>
                </div>
            </div>

            <div class="w-full lg:w-4/12 xl:w-3/12 bg-white rounded-3xl shadow-xl border border-gray-200 flex flex-col h-auto lg:h-[82vh] lg:sticky lg:top-4 overflow-hidden shrink-0">
                
                <div class="p-5 border-b border-gray-100 bg-slate-50 flex justify-between items-center shrink-0">
                    <h2 class="text-lg font-black text-gray-800 flex items-center gap-2">🛒 Rincian Pesanan</h2>
                    <button @click="clearCart" :disabled="cart.length===0" class="text-red-500 hover:text-white hover:bg-red-500 px-3 py-1.5 rounded-lg text-[10px] font-black uppercase tracking-widest transition-colors disabled:opacity-50 border border-red-200">
                        Void (Batal)
                    </button>
                </div>

                <div class="p-4 flex-1 overflow-y-auto bg-white custom-scrollbar">
                    <div v-if="cart.length === 0" class="h-full flex flex-col items-center justify-center opacity-40 py-10">
                        <span class="text-6xl mb-4">🧾</span>
                        <p class="text-gray-600 font-bold text-lg">Keranjang Kosong</p>
                    </div>

                    <div v-for="item in cart" :key="item.id" class="flex flex-col mb-3 p-3 bg-slate-50 rounded-xl border border-gray-100 shadow-sm hover:border-blue-200 transition-colors">
                        <div class="flex justify-between items-start mb-2">
                            <h3 class="font-bold text-sm text-gray-800 leading-tight pr-2 line-clamp-2">{{ item.name }}</h3>
                            <div class="font-black text-sm text-gray-900 whitespace-nowrap">Rp {{ (item.price * item.qty).toLocaleString('id-ID') }}</div>
                        </div>
                        <div class="flex justify-between items-center">
                            <p class="text-xs font-bold text-gray-500">@ Rp {{ item.price.toLocaleString('id-ID') }}</p>
                            <div class="flex items-center bg-white rounded-lg p-0.5 border border-gray-200 shadow-sm">
                                <button @click="decreaseQty(item)" class="w-7 h-7 flex items-center justify-center rounded text-gray-600 hover:bg-red-50 hover:text-red-600 font-black transition-colors">-</button>
                                <input type="number" v-model.number="item.qty" @change="validateQty(item)" class="w-10 text-center text-xs font-black text-gray-800 bg-transparent border-none focus:ring-0 p-0 [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none">
                                <button @click="addToCart(item)" class="w-7 h-7 flex items-center justify-center rounded text-gray-600 hover:bg-blue-50 hover:text-blue-600 font-black transition-colors">+</button>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="p-5 bg-white border-t border-gray-200 shadow-[0_-10px_20px_-10px_rgba(0,0,0,0.03)] shrink-0 z-10">
                    
                    <div class="mb-4">
                        <span class="font-black text-[10px] text-gray-400 block mb-2 uppercase tracking-widest text-center">Metode Pembayaran</span>
                        <div class="grid grid-cols-3 gap-2">
                            <button @click="setPaymentMethod('Cash')" :class="paymentMethod === 'Cash' ? 'bg-blue-600 text-white shadow-md ring-2 ring-blue-200' : 'bg-slate-50 text-gray-600 border border-gray-200 hover:bg-gray-100'" class="py-2 rounded-xl font-black text-xs transition-all flex flex-col items-center gap-1">
                                <span class="text-base">💵</span> Tunai
                            </button>
                            <button @click="setPaymentMethod('QRIS')" :class="paymentMethod === 'QRIS' ? 'bg-blue-600 text-white shadow-md ring-2 ring-blue-200' : 'bg-slate-50 text-gray-600 border border-gray-200 hover:bg-gray-100'" class="py-2 rounded-xl font-black text-xs transition-all flex flex-col items-center gap-1">
                                <span class="text-base">📱</span> QRIS
                            </button>
                            <button @click="setPaymentMethod('Debit')" :class="paymentMethod === 'Debit' ? 'bg-blue-600 text-white shadow-md ring-2 ring-blue-200' : 'bg-slate-50 text-gray-600 border border-gray-200 hover:bg-gray-100'" class="py-2 rounded-xl font-black text-xs transition-all flex flex-col items-center gap-1">
                                <span class="text-base">💳</span> Debit
                            </button>
                        </div>
                    </div>

                    <div class="space-y-3 mb-4">
                        <div class="flex justify-between items-end border-b border-dashed border-gray-200 pb-2">
                            <span class="font-bold text-xs text-gray-500 uppercase tracking-wider">Total Tagihan</span>
                            <span class="text-3xl font-black text-blue-800 leading-none tracking-tight">Rp {{ totalBelanja.toLocaleString('id-ID') }}</span>
                        </div>

                        <div class="flex justify-between items-center bg-slate-50 p-3 rounded-xl border border-gray-200 focus-within:ring-2 focus-within:ring-blue-500 transition-all">
                            <span class="font-bold text-xs text-gray-700 uppercase tracking-wider">Bayar</span>
                            <div class="relative flex-1 ml-4">
                                <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 font-bold text-sm">Rp</span>
                                <input
                                    type="text"
                                    :value="payAmount === 0 ? '' : payAmount.toLocaleString('id-ID')"
                                    @input="formatInputRupiah"
                                    :disabled="paymentMethod !== 'Cash'"
                                    :class="paymentMethod !== 'Cash' ? 'bg-gray-200/50 text-gray-500 cursor-not-allowed border-transparent' : 'bg-white text-gray-900 border-gray-300 shadow-sm'"
                                    class="w-full text-right text-lg font-black rounded-lg py-1.5 pl-8 pr-3 transition-all outline-none border"
                                    placeholder="0">
                            </div>
                        </div>

                        <div class="flex justify-between items-center pt-1 px-1">
                            <span class="font-bold text-xs text-gray-500 uppercase tracking-wider">Kembali</span>
                            <span class="text-xl font-black" :class="kembalian >= 0 ? 'text-green-500' : 'text-red-500'">
                                Rp {{ kembalian.toLocaleString('id-ID') }}
                            </span>
                        </div>
                    </div>

                    <button @click="processCheckout" :disabled="cart.length === 0 || payAmount < totalBelanja"
                        class="w-full bg-green-500 hover:bg-green-600 text-white font-black py-4 px-4 rounded-xl transition-all flex justify-center items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed shadow-lg hover:shadow-xl transform hover:-translate-y-0.5 text-lg uppercase tracking-wider">
                        Bayar Sekarang 💸
                    </button>
                </div>
            </div>
        </div>

        <div v-if="showQrisModal" class="fixed inset-0 bg-gray-900/80 flex items-center justify-center z-50 p-4 backdrop-blur-sm">
            <div class="bg-white p-8 rounded-3xl shadow-2xl w-full max-w-md text-center transform transition-all flex flex-col">
                <div class="w-16 h-16 bg-blue-100 text-blue-600 rounded-full flex items-center justify-center mx-auto mb-4 text-3xl">📱</div>
                <h2 class="text-2xl font-black text-gray-900 mb-1">Scan QRIS</h2>
                <p class="text-gray-500 text-sm mb-6 font-medium">Arahkan kamera ke kode QR di bawah ini</p>
                <div class="bg-white p-3 rounded-2xl border-2 border-dashed border-gray-300 w-full mb-6 shadow-inner flex justify-center items-center min-h-[250px]">
                    <img src="https://placehold.co/300x300?text=QRIS+Toko" alt="QRIS Toko" class="w-full h-full max-h-64 object-contain mx-auto rounded-xl">
                </div>
                <div class="bg-blue-50 text-blue-900 p-4 rounded-2xl mb-8 border border-blue-100">
                    <span class="block text-xs font-bold uppercase tracking-widest opacity-70 mb-1">Total Tagihan</span>
                    <span class="text-4xl font-black">Rp {{ totalBelanja.toLocaleString('id-ID') }}</span>
                </div>
                <div class="flex gap-4">
                    <button @click="showQrisModal = false" class="flex-1 bg-gray-100 py-3.5 rounded-xl font-black text-gray-600 hover:bg-gray-200 transition-colors">Batal</button>
                    <button @click="executeCheckout" class="flex-1 bg-blue-600 py-3.5 rounded-xl font-black text-white hover:bg-blue-700 transition-colors shadow-lg">Selesai Dibayar</button>
                </div>
            </div>
        </div>

        <div v-if="showReceipt" class="fixed inset-0 bg-gray-900/80 flex items-center justify-center z-50 p-4 backdrop-blur-sm">
            <div class="bg-gray-200 p-4 rounded-2xl shadow-2xl w-full max-w-sm overflow-hidden border-t-8 border-gray-800">
                <div id="print-area" class="text-left font-mono text-[11px] leading-tight uppercase text-black bg-white p-4 mx-auto" style="width: 58mm;">
                    <div class="text-center mb-3">
                        <h2 class="font-black text-sm mb-1">INDO UMKM</h2>
                        <p class="font-medium">JL. KEBON KOSONG NO 56 F</p>
                    </div>
                    <div class="text-center my-2 font-bold tracking-widest border-y border-black py-1">
                        <p>S T R U K   B E L A N J A</p>
                    </div>
                    <div class="mb-2 text-[10px] font-bold">
                        <p>{{ lastTransaction?.date }} / {{ currentUser.name.split(' ')[0] }} / 01</p>
                    </div>
                    <p class="border-b border-dashed border-black mb-2"></p>
                    <div v-for="item in lastTransaction?.cart" :key="item.id" class="mb-1.5 font-bold">
                        <div class="truncate w-full pr-2">{{ item.name }}</div>
                        <div class="flex justify-between pl-4 text-[10px]">
                            <span>{{ item.qty }} x {{ item.price.toLocaleString('id-ID') }}</span>
                            <span>{{ (item.price * item.qty).toLocaleString('id-ID') }}</span>
                        </div>
                    </div>
                    <p class="border-t border-dashed border-black mt-2 pt-2"></p>
                    <div class="flex justify-between font-black text-xs mb-2">
                        <span>TOTAL BELANJA :</span>
                        <span>{{ lastTransaction?.total.toLocaleString('id-ID') }}</span>
                    </div>
                    <p class="border-b border-dashed border-black mb-2"></p>
                    <div class="flex justify-between mb-1 font-bold">
                        <span>{{ lastTransaction?.method === 'Cash' ? 'TUNAI' : 'QRIS/DEBIT' }} :</span>
                        <span>{{ lastTransaction?.pay.toLocaleString('id-ID') }}</span>
                    </div>
                    <div v-if="lastTransaction?.method === 'Cash'" class="flex justify-between mb-2 font-bold">
                        <span>KEMBALIAN :</span>
                        <span>{{ lastTransaction?.return.toLocaleString('id-ID') }}</span>
                    </div>
                    <div class="mt-4 text-[9px] font-medium text-center border-t border-dashed border-black pt-2">
                        <p>SUBTOTAL: {{ lastTransaction?.subtotal.toLocaleString('id-ID') }} | PAJAK: {{ lastTransaction?.pajak.toLocaleString('id-ID') }}</p>
                        <p class="mt-1">TRX-ID: {{ lastTransaction?.invoice }}</p>
                    </div>
                    <div class="text-center mt-4 font-bold">
                        <p>=== TERIMA KASIH ===</p>
                        <p>BARANG SUDAH DIBELI TIDAK DAPAT DITUKAR</p>
                    </div>
                </div>
                <div class="mt-4 flex gap-3 no-print">
                    <button @click="printReceipt" class="flex-1 bg-gray-800 text-white py-3 rounded-xl font-black hover:bg-gray-900 transition-colors shadow-md">🖨️ CETAK</button>
                    <button @click="showReceipt = false" class="flex-1 bg-white border border-gray-300 py-3 rounded-xl font-black text-gray-800 hover:bg-gray-50 transition-colors">TUTUP</button>
                </div>
            </div>
        </div>

    </div>
</template>

<style>
/* CSS Scrollbar dan Print tetap sama */
.custom-scrollbar::-webkit-scrollbar { width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #94a3b8; }

@media print {
    body * { visibility: hidden; }
    #print-area, #print-area * { visibility: visible; }
    #print-area { position: absolute; left: 0; top: 0; width: 58mm; padding: 0; margin: 0; }
    @page { margin: 0; }
    .no-print { display: none !important; }
}

.swal2-popup {
    font-family: 'Inter', sans-serif !important;
    border-radius: 24px !important;
}
.swal2-confirm {
    border-radius: 12px !important;
    font-weight: 800 !important;
    text-transform: uppercase !important;
    letter-spacing: 0.5px !important;
}
</style>