<script setup>

import { ref, computed, onMounted, nextTick, onUnmounted } from 'vue';
<<<<<<< HEAD
import { Head } from '@inertiajs/vue3';
=======
import { Head, router, usePage } from '@inertiajs/vue3';
>>>>>>> fe5108ac44404475b1c66a1cb99fadaf9190e970
import AppLayout from '@/Layouts/AppLayout.vue';
import api from '../../api.js';

// Mengambil data page dari Inertia (untuk QRIS dan Struk)
const page = usePage();

// Fungsi Jam Realtime
const currentTime = ref('');
let timer;

// State utama untuk keranjang, pembayaran, dan tampilan
const products = ref([]);
const isLoadingProducts = ref(true);
const cart = ref([]);
const heldOrders = ref([]);
const payAmount = ref(0);
const paymentMethod = ref('Cash')
const showReceipt = ref(false);
const showQrisModal = ref(false);
const lastTransaction = ref(null);

//Pencarian Produk & Barcode Scan
const searchQuery = ref('');
const searchInput = ref(null);

// --- FUNGSI BARU: Tarik Data dari Golang ---
const fetchProducts = async () => {
    try {
        const response = await api.get('/products');
        // Petakan nama kolom Golang (nama_produk) ke nama variabel Vue Mas (name)
        products.value = response.data.data.map(p => ({
            id: p.id,
            sku: p.sku || `SKU-${p.id}`,
            name: p.nama_produk,
            price: p.harga_jual,
            stock: p.stok,
            image: null
        }));
    } catch (error) {
        console.error("Gagal narik produk:", error);
        if (error.response && error.response.status === 401) {
            window.location.href = '/login'; // Lempar ke login kalau tiket mati
        }
    } finally {
        isLoadingProducts.value = false;
    }
};

// Fokus ke input pencarian saat halaman dimuat
onMounted(() => {
    fetchProducts();

    if (searchInput.value) searchInput.value.focus();
    // Jalankan Jam Setiap Detik
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


// Filter produk secara otomatis berdasarkan ketikan kasir
const filteredProducts = computed(() => {
    if (!searchQuery.value) return products.value;
    const query = searchQuery.value.toLocaleLowerCase();
    return products.value.filter(product =>
        product.name.toLowerCase().includes(query) ||
        (product.sku && product.sku.toLowerCase().includes(query))
    );
});


// Fungsi Otomatis oleh scanner barcode
const handleBarcodeScan = () => {
    if (!searchQuery.value) return;
    const query = String(searchQuery.value).trim().toLowerCase();

    // cek apakah input cocok dengan SKU produk
    const exactMatch = props.products.find(p =>
        p.sku && String(p.sku).toLowerCase() === query);

    if (exactMatch) {
        addToCart(exactMatch); //masuk keranjang
        searchQuery.value = ''; // reset input
    } else if (filteredProducts.value.length === 1) {
        addToCart(filteredProducts.value[0]);
        searchQuery.value = '';
    }

    // Kursor fokus ke kotak pencarian
    nextTick(() => {
        if (searchInput.value) searchInput.value.focus();
    });
};

// Fungsi Keranjang
const addToCart = (product) => {
    if (product.stock <= 0) {
        alert("Stok barang habis Bos!");
        return;
    }
    const existingItem = cart.value.find(item => item.id === product.id);
    if (existingItem) {
        if (existingItem.qty < product.stock) {
            existingItem.qty++;
        } else {
            alert("Kuantitas melebihi stok yang ada!");
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

// Fungsi Void
const clearCart = () => {
    if (cart.value.length > 0 && confirm("Apakah Anda yakin ingin membatalkan transaksi ini?")) {
        cart.value = [];
        payAmount.value = 0;
        setPaymentMethod('Cash');
    }
};

// fungsi Hold Order
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
<<<<<<< HEAD
=======
    setPaymentMethod('Cash');
>>>>>>> fe5108ac44404475b1c66a1cb99fadaf9190e970
};



const resumeOrder = (order) => {
    if (cart.value.length > 0 && !confirm('Keranjang saat ini akan diganti. Lanjutkan?')) return;

    cart.value = [...order.items];

    heldOrders.value = heldOrders.value.filter(h => h.id !== order.id);
};

// Perhitungan Total, Kembalian, dan Pajak
const totalBelanja = computed(() => {
    return cart.value.reduce((total, item) => total + (item.price * item.qty), 0);
});

const kembalian = computed(() => {
    return payAmount.value - totalBelanja.value;
});

// Proses Bayar
const setPaymentMethod = (method) => {
    paymentMethod.value = method;

    // Kalau non-tunai, uang bayar otomatis nge-pas dengan total belanja
    if (method !== 'Cash') {
        payAmount.value = totalBelanja.value;
    } else {
        payAmount.value = 0; // Reset kalau balik ke Tunai
    }
};


<<<<<<< HEAD

// Fungsi Eksekusi ke database (dipisah agar bisa dipanggil dari tombol tunai atau qris)
const executeCheckout = async() => {
    const payloadItems = cart.value.map(item => ({
        product_id: item.id,
        kuantitas: item.qty
    }));

    try {
        // Tembak mesin golang
        const response = await api.post('/checkout', {
            items: payloadItems,
            nominal_bayar: payAmount.value
        });

        // Ambil data respon dari Golang
        // Hitung PPN Mundur (Include TAX 11%)
        const dppAmount = Math.round(totalBelanja.value / 1.11);
        const ppnAmount = totalBelanja.value - dppAmount;

        lastTransaction.value = {
            invoice: response.data.invoice, // No Invoice Asli dari DB!
            cart: [...cart.value],
            total: response.data.tagihan,   // Total fix setelah diolah Golang
            pay: payAmount.value,
            return: response.data.kembali,  // Kembalian fix dari Golang
            method: paymentMethod.value,
            dpp: dppAmount,
            ppn: ppnAmount,
            date: new Date().toLocaleString('id-ID', { year: '2-digit', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' }).replace(/\//g, '.')
        };
        showQrisModal.value = false;
        showReceipt.value = true;
        cart.value = [];
        payAmount.value = 0;
        paymentMethod.value = 'Cash';

        // Refresh katalog barang biar stok terbaru langsung ngesync
        fetchProducts();
    } catch (error) {
        if (error.response && error.response.data.error) {
            alert("Transaksi Gagal: " + error.response.data.error);
        } else {
            alert("Koneksi ke server terputus");
=======
    router.post('/pos/checkout', {
        cart: cart.value,
        total_price: totalBelanja.value,
        pay_amount: payAmount.value,
        payment_method: paymentMethod.value
    }, {
        onSuccess: () => {
            lastTransaction.value = {
                cart: [...cart.value],
                total: totalBelanja.value,
                pay: payAmount.value,
                return: kembalian.value,
                method: paymentMethod.value,
                dpp: dppAmount,
                ppn: ppnAmount,
                date: new Date().toLocaleString('id-ID', { year: '2-digit', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' }).replace(/\//g, '.')
            };
            showQrisModal.value = false;
            showReceipt.value = true;
            cart.value = [];
            payAmount.value = 0;
            paymentMethod.value = 'Cash';
            // Fokusin lagi ke search bar setelah bayar
            nextTick(() => { if (searchInput.value) searchInput.value.focus(); });
>>>>>>> fe5108ac44404475b1c66a1cb99fadaf9190e970
        }
    }
};

// Proses bayar Utama
<<<<<<< HEAD
const formattedPayAmount = computed({
    get() {
        // Saat nampil di layar: ubah angka jadi format lokal (ada titiknya)
        // Kalau nilainya 0, tampilkan string kosong agar kasir gampang ngetik
        return payAmount.value === 0 ? '' : payAmount.value.toLocaleString('id-ID');
    },
    set(newValue) {
        // Saat kasir ngetik: bersihkan semua huruf/titik, sisakan angka murni saja
        const cleanValue = String(newValue).replace(/\D/g, '');

        // Simpan angka murninya ke variabel asli payAmount
        payAmount.value = cleanValue ? parseInt(cleanValue, 10) : 0;
    }
});
=======
const formatInputRupiah = (event) => {
    let rawValue = event.target.value.replace(/\D/g, '');
    payAmount.value = rawValue ? parseInt(rawValue, 10) : 0;
    event.target.value = payAmount.value === 0 ? '' : payAmount.value.toLocaleString('id-ID');
}
>>>>>>> fe5108ac44404475b1c66a1cb99fadaf9190e970



const processCheckout = () => {
    if (payAmount.value < totalBelanja.value) {
        alert("Maaf, uang pembayaran kurang!");
        return;
    }

    if (paymentMethod.value === 'QRIS') {
        // Tampilkan modal QRIS
        showQrisModal.value = true;
    } else {
        // Proses checkout langsung untuk Cash atau Debit
        executeCheckout();
    }

};

const printReceipt = () => {
    window.print();
};

</script>

<template>
    <Head title="Kasir POS" />

    <AppLayout>
        <div class="p-4 md:p-6 bg-gray-50 min-h-screen">

            <div
                class="bg-blue-800 text-white rounded-xl shadow-md mb-6 flex flex-col md:flex-row overflow-hidden border border-blue-900">
                    <div
                        class="p-3 px-5 bg-blue-900 md:w-1/4 flex flex-col justify-center">
                        <div
                            class="font-black text-lg tracking-wider text-blue-50 text-center">POS KASIR
                        </div>
                        <div
                            class="text-center text-[10px] text-blue-300 uppercase mt-0.5 tracking-widest font-semibold">{{ $page.props.auth.store?.nama_toko || 'Toko Belum di Setting'}}
                        </div>
                    </div>

                <div
                    class="flex-1 flex flex-col justify-between border-r border-blue-700">
                        <div
                            class="bg-red-600 text-white text-xs font-bold py-1 px-4 text-center tracking-widest overflow-hidden whitespace-nowrap">
                            <marquee scrollamount="6">INFO KASIR: PASTIKAN SCAN BARCODE DENGAN BENAR | JANGAN LUPA TAWARKAN PRODUK BEST SELLER KEPADA KONSUMEN</marquee>
                        </div>

                    <div class="flex justify-between items-center px-4 py-2 bg-blue-800">
                        <div class="text-xs md:text-sm font-semibold text-blue-100 flex gap-3 md:gap-6">
                            <span class="flex items-center gap-1">
                                <span class="opacity-70">Kasir:</span>
                                {{ $page.props.auth.user.name }}
                            </span>
                            <span class="flex items-center gap-1">
                                <span class="opacity-70">
                                    Shift:
                                </span>
                                1
                            </span>
                            <span class="flex items-center gap-1">
                                <span class="opacity-70">
                                    Station:
                                </span>
                                01
                            </span>
                        </div>
                        <div
                            class="font-mono text-sm md:text-lg font-bold text-yellow-300 tracking-wider">
                            {{ currentTime }}
                        </div>
                    </div>
                </div>

                <div class="p-2 w-full md:w-40 bg-white flex items-center justify-center shrink-0 border-b-4 border-red-600 md:border-b-0 md:border-l-4">
                    <div class="text-center flex gap-1 items-center">
                        <span class="text-red-600 font-black text-xl italic leading-none">Indo</span>
                        <span class="text-blue-600 font-black text-xl italic leading-none">UMKM</span>
                    </div>
                </div>
            </div>

            <div class="flex flex-col lg:flex-row gap-6">

                <div class="w-full lg:w-8/12 flex flex-col">
                    <div class="mb-5 relative group">
                        <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                            <span class="text-gray-400 text-xl group-focus-within:text-blue-500 transition-colors">🔍</span>
                        </div>
                        <input
                            ref="searchInput"
                            type="text"
                            v-model="searchQuery"
                            @keydown.enter.prevent="handleBarcodeScan"
                            placeholder="Ketik nama produk atau Scan Barcode (Tekan Enter)..."
                            class="w-full pl-12 pr-4 py-3.5 rounded-2xl border-gray-200 shadow-sm focus:ring-4 focus:ring-blue-500/20 focus:border-blue-500 text-gray-800 font-bold bg-white text-lg transition-all"
                        >
                    </div>

                    <div v-if="filteredProducts.length === 0" class="flex-1 flex items-center justify-center bg-white rounded-2xl border border-dashed border-gray-300 min-h-[400px]">
                        <div class="text-center">
                            <span class="text-4xl block mb-3">📦</span>
                            <p class="text-gray-500 font-bold text-lg">Barang tidak ditemukan</p>
                        </div>
                    </div>

                    <div v-else class="grid grid-cols-2 sm:grid-cols-3 xl:grid-cols-4 gap-4 pb-6 overflow-y-auto max-h-[65vh] pr-2 custom-scrollbar">
                        <div v-for="product in filteredProducts" :key="product.id" @click="addToCart(product)"
                            class="bg-white rounded-2xl shadow-sm hover:shadow-xl hover:border-blue-500 transition-all duration-200 overflow-hidden cursor-pointer border border-gray-100 group flex flex-col h-full transform hover:-translate-y-1">
                            <div class="relative pt-3 px-3">
                                <div class="bg-gray-50 rounded-xl overflow-hidden">
                                    <img :src="product.image || 'https://placehold.co/150x150?text=No+Image'" :alt="product.name" class="w-full h-28 object-contain mix-blend-multiply p-2 group-hover:scale-110 transition-transform duration-300">
                                </div>
                            </div>
                            <div class="p-4 flex flex-col flex-1">
                                <h2 class="font-bold text-gray-800 text-sm line-clamp-2 leading-snug flex-1" :title="product.name">{{ product.name }}</h2>
                                <div class="mt-3 flex justify-between items-end">
                                    <p class="text-blue-700 font-black text-base">Rp {{ product.price.toLocaleString('id-ID') }}</p>
                                </div>
                                <div class="mt-2 text-[11px] font-bold px-2 py-1 rounded-md inline-block w-max"
                                    :class="product.stock > 10 ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'">
                                    Stok: {{ product.stock }}
                                </div>
                            </div>
                        </div>
                    </div>

                    <div v-if="heldOrders.length > 0" class="mt-auto pt-4 border-t border-gray-200">
                        <h3 class="text-xs font-black text-gray-500 uppercase tracking-widest mb-3 flex items-center gap-2">
                            <span class="w-2 h-2 rounded-full bg-orange-500 animate-pulse"></span> Pesanan Tertunda
                        </h3>
                        <div class="flex flex-wrap gap-3">
                            <div v-for="order in heldOrders" :key="order.id" @click="resumeOrder(order)"
                                class="bg-gradient-to-br from-orange-50 to-orange-100 border border-orange-200 p-3 rounded-xl cursor-pointer hover:shadow-md transition-all w-40 relative overflow-hidden group">
                                <div class="absolute -right-4 -top-4 text-4xl opacity-10 group-hover:rotate-12 transition-transform">⏳</div>
                                <div class="text-[10px] font-black text-orange-600 tracking-wider">{{ order.time }}</div>
                                <div class="text-sm font-black text-gray-900 mt-1">Rp {{ order.total.toLocaleString('id-ID') }}</div>
                                <div class="text-[11px] text-gray-600 font-medium mt-1">{{ order.items.length }} Item di keranjang</div>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="w-full lg:w-4/12">
                    <div class="bg-white rounded-3xl shadow-xl border border-gray-100 flex flex-col h-[78vh] overflow-hidden">

                        <div class="p-5 border-b border-gray-100 bg-white flex justify-between items-center z-10 shadow-sm">
                            <h2 class="text-lg font-black text-gray-800 flex items-center gap-2">
                                🛒 Rincian Pesanan
                            </h2>
                            <button @click="clearCart" class="text-red-500 hover:text-red-700 hover:bg-red-50 px-3 py-1.5 rounded-lg text-xs font-black uppercase tracking-wider transition-colors">
                                Void (Batal)
                            </button>
                        </div>

                        <div class="p-5 flex-1 overflow-y-auto bg-gray-50/50 custom-scrollbar">
                            <div v-if="cart.length === 0" class="h-full flex flex-col items-center justify-center opacity-50">
                                <span class="text-5xl mb-3">🧾</span>
                                <p class="text-gray-500 font-bold">Keranjang masih kosong</p>
                            </div>

                            <div v-for="item in cart" :key="item.id" class="flex justify-between items-center mb-3 p-3 bg-white rounded-xl border border-gray-100 shadow-sm hover:border-blue-200 transition-colors">
                                <div class="flex-1 pr-3">
                                    <h3 class="font-bold text-sm text-gray-800 leading-tight mb-1 line-clamp-1">{{ item.name }}</h3>
                                    <p class="text-xs font-bold text-blue-600">Rp {{ item.price.toLocaleString('id-ID') }}</p>
                                </div>
                                <div class="flex flex-col items-end gap-2">
                                    <div class="font-black text-sm text-gray-900">Rp {{ (item.price * item.qty).toLocaleString('id-ID') }}</div>
                                    <div class="flex items-center bg-gray-100 rounded-lg p-1 border border-gray-200">
                                        <button @click="decreaseQty(item)" class="w-7 h-7 flex items-center justify-center bg-white rounded shadow-sm text-gray-600 hover:text-red-600 font-black transition-colors">-</button>
                                        <input type="number" v-model.number="item.qty" @change="validateQty(item)" class="w-10 text-center text-xs font-black text-gray-800 bg-transparent border-none focus:ring-0 p-0 [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none">
                                        <button @click="addToCart(item)" class="w-7 h-7 flex items-center justify-center bg-white rounded shadow-sm text-gray-600 hover:text-blue-600 font-black transition-colors">+</button>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <div class="p-5 bg-white border-t border-gray-200 shadow-[0_-10px_20px_-10px_rgba(0,0,0,0.05)]">
                            <button @click="holdTransaction" :disabled="cart.length === 0" class="w-full mb-4 bg-orange-50 text-orange-600 py-2.5 rounded-xl font-bold text-xs uppercase tracking-widest border border-orange-200 hover:bg-orange-100 disabled:opacity-50 transition-colors flex justify-center items-center gap-2">
                                ⏸️ Simpan Transaksi (Hold)
                            </button>

                            <div class="mb-5">
                                <span class="font-black text-[11px] text-gray-400 block mb-2 uppercase tracking-widest">Metode Pembayaran</span>
                                <div class="grid grid-cols-3 gap-2">
                                    <button @click="setPaymentMethod('Cash')" :class="paymentMethod === 'Cash' ? 'bg-blue-600 text-white shadow-md ring-2 ring-blue-600 ring-offset-1' : 'bg-gray-50 text-gray-600 border border-gray-200 hover:bg-gray-100'" class="py-2.5 rounded-xl font-black text-xs transition-all flex flex-col items-center gap-1">
                                        <span class="text-base">💵</span> Tunai
                                    </button>
                                    <button @click="setPaymentMethod('QRIS')" :class="paymentMethod === 'QRIS' ? 'bg-blue-600 text-white shadow-md ring-2 ring-blue-600 ring-offset-1' : 'bg-gray-50 text-gray-600 border border-gray-200 hover:bg-gray-100'" class="py-2.5 rounded-xl font-black text-xs transition-all flex flex-col items-center gap-1">
                                        <span class="text-base">📱</span> QRIS
                                    </button>
                                    <button @click="setPaymentMethod('Debit')" :class="paymentMethod === 'Debit' ? 'bg-blue-600 text-white shadow-md ring-2 ring-blue-600 ring-offset-1' : 'bg-gray-50 text-gray-600 border border-gray-200 hover:bg-gray-100'" class="py-2.5 rounded-xl font-black text-xs transition-all flex flex-col items-center gap-1">
                                        <span class="text-base">💳</span> Debit
                                    </button>
                                </div>
                            </div>

                            <div class="space-y-3 mb-5">
                                <div class="flex justify-between items-end border-b border-dashed border-gray-200 pb-3">
                                    <span class="font-bold text-sm text-gray-500">Total Tagihan</span>
                                    <span class="text-2xl font-black text-blue-900 leading-none">Rp {{ totalBelanja.toLocaleString('id-ID') }}</span>
                                </div>

                                <div class="flex justify-between items-center bg-gray-50 p-3 rounded-xl border border-gray-200">
                                    <span class="font-bold text-sm text-gray-700">Nominal Bayar</span>
                                    <div class="relative">
                                        <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500 font-bold text-sm">Rp</span>
                                        <input
                                            type="text"
                                            :value="payAmount === 0 ? '' : payAmount.toLocaleString('id-ID')"
                                            @input="formatInputRupiah"
                                            :disabled="paymentMethod !== 'Cash'"
                                            :class="paymentMethod !== 'Cash' ? 'bg-gray-200/50 text-gray-500 cursor-not-allowed' : 'bg-white text-gray-900 focus:ring-2 focus:ring-blue-500 border-gray-300'"
                                            class="w-36 text-right text-base font-black rounded-lg py-1.5 pl-8 pr-3 transition-all shadow-sm"
                                            placeholder="0">
                                    </div>
                                </div>

                                <div class="flex justify-between items-center pt-2">
                                    <span class="font-bold text-sm text-gray-500">Kembalian</span>
                                    <span class="text-xl font-black" :class="kembalian >= 0 ? 'text-green-600' : 'text-red-500'">
                                        Rp {{ kembalian.toLocaleString('id-ID') }}
                                    </span>
                                </div>
                            </div>

<<<<<<< HEAD
                            <div class="grid grid-cols-1 mb-3">
                                <button @click="holdTransaction" :disabled="cart.length === 0" class="bg-orange-100 text-orange-700 py-2 rounded-lg font-bold text-xs uppercase border border-orange-200 hover:bg-orange-200 disabled:opacity-50 transition-colors">
                                    ⏸️ Hold / Simpan Pesanan
                                </button>
                            </div>

                            <div class="flex justify-between items-center mb-2">
                                <span class="font-semibold text-sm text-gray-600">Nominal Bayar</span>
                                <input type="text" v-model="formattedPayAmount" :disabled="paymentMethod !== 'Cash'" :class="paymentMethod !== 'Cash' ? 'bg-gray-100 text-gray-500' : 'bg-white text-gray-800 focus:ring-blue-500 focus:border-blue-500'" class="w-32 text-right text-sm font-bold border-gray-300 rounded-md p-1 transition-colors" placeholder="0">
                            </div>

                            <div class="flex justify-between items-center mb-4 pt-2 border-t border-dashed border-gray-300">
                                <span class="font-bold text-gray-600">Kembalian</span>
                                <span class="text-xl font-black" :class="kembalian >= 0 ? 'text-green-600' : 'text-red-500'">
                                    Rp {{ kembalian.toLocaleString('id-ID') }}
                                </span>
                            </div>

                            <button @click="processCheckout" :disabled="cart.length === 0 || payAmount < totalBelanja" class="w-full bg-blue-600 hover:bg-blue-700 text-white font-bold py-3 px-4 rounded-lg transition-colors flex justify-center items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed shadow-md">
                                Proses Pembayaran
=======
                            <button @click="processCheckout" :disabled="cart.length === 0 || payAmount < totalBelanja"
                                class="w-full bg-blue-700 hover:bg-blue-800 text-white font-black py-4 px-4 rounded-xl transition-all flex justify-center items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed shadow-lg hover:shadow-xl transform hover:-translate-y-0.5 text-lg">
                                Bayar Sekarang 🚀
>>>>>>> fe5108ac44404475b1c66a1cb99fadaf9190e970
                            </button>
                        </div>
                    </div>
                </div>

            </div>
        </div>

        <div v-if="showQrisModal" class="fixed inset-0 bg-gray-900/80 flex items-center justify-center z-50 p-4 backdrop-blur-sm">
            <div class="bg-white p-8 rounded-3xl shadow-2xl w-full max-w-md text-center transform transition-all flex flex-col">
                <div class="w-16 h-16 bg-blue-100 text-blue-600 rounded-full flex items-center justify-center mx-auto mb-4 text-3xl">📱</div>
                <h2 class="text-2xl font-black text-gray-900 mb-1">Scan QRIS</h2>
                <p class="text-gray-500 text-sm mb-6 font-medium">Arahkan kamera ke kode QR di bawah ini</p>

                <div class="bg-white p-3 rounded-2xl border-2 border-dashed border-gray-300 w-full mb-6 shadow-inner flex justify-center items-center min-h-[250px]">
                    <img :src="$page.props.auth.store?.qris_image ? '/storage/' + $page.props.auth.store.qris_image : 'https://placehold.co/300x300?text=QRIS+Belum+Diatur'" alt="QRIS Toko" class="w-full h-full max-h-64 object-contain mx-auto rounded-xl">
                </div>

                <div class="bg-blue-50 text-blue-900 p-4 rounded-2xl mb-8 border border-blue-100">
                    <span class="block text-xs font-bold uppercase tracking-widest opacity-70 mb-1">Total Tagihan</span>
                    <span class="text-4xl font-black">Rp {{ totalBelanja.toLocaleString('id-ID') }}</span>
                </div>

                <div class="flex gap-4">
                    <button @click="showQrisModal = false" class="flex-1 bg-gray-100 py-3.5 rounded-xl font-black text-gray-600 hover:bg-gray-200 transition-colors">Batal</button>
                    <button @click="executeCheckout" class="flex-1 bg-blue-600 py-3.5 rounded-xl font-black text-white hover:bg-blue-700 transition-colors shadow-lg hover:shadow-xl">Selesai Dibayar</button>
                </div>
            </div>
        </div>

        <div v-if="showReceipt" class="fixed inset-0 bg-gray-900/80 flex items-center justify-center z-50 p-4 backdrop-blur-sm">
            <div class="bg-gray-200 p-4 rounded-2xl shadow-2xl w-full max-w-sm overflow-hidden border-t-8 border-gray-800">
                <div id="print-area" class="text-left font-mono text-[11px] leading-tight uppercase text-black bg-white p-4 mx-auto" style="width: 58mm;">

                    <div class="text-center mb-3">
                        <h2 class="font-black text-sm mb-1">{{ $page.props.auth.store?.nama_toko || 'NAMA TOKO' }}</h2>
                        <p class="font-medium">{{ $page.props.auth.store?.alamat_toko || 'ALAMAT BELUM DIATUR' }}</p>
                        <p v-if="$page.props.auth.store?.telepon" class="font-medium mt-0.5">TELP: {{ $page.props.auth.store?.telepon }}</p>
                    </div>

                    <div class="text-center my-2 font-bold tracking-widest border-y border-black py-1">
                        <p>S T R U K   B E L A N J A</p>
                    </div>

                    <div class="mb-2 text-[10px] font-bold">
                        <p>{{ lastTransaction?.date }} / {{ $page.props.auth.user.name.split(' ')[0] }} / 01</p>
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
                        <p>PPN INC: DPP={{ lastTransaction?.dpp.toLocaleString('id-ID') }} PPN={{ lastTransaction?.ppn.toLocaleString('id-ID') }}</p>
                        <p class="mt-1">TRX-ID: {{ Date.now().toString().slice(-8) }}</p>
                    </div>

                    <div class="text-center mt-4 font-bold">
                        <p>=== TERIMA KASIH ===</p>
                        <p>BARANG YANG SUDAH DIBELI</p>
                        <p>TIDAK DAPAT DITUKAR/DIKEMBALIKAN</p>
                    </div>
                </div>

                <div class="mt-4 flex gap-3 no-print">
                    <button @click="printReceipt" class="flex-1 bg-gray-800 text-white py-3 rounded-xl font-black hover:bg-gray-900 transition-colors shadow-md">🖨️ CETAK</button>
                    <button @click="showReceipt = false" class="flex-1 bg-white border border-gray-300 py-3 rounded-xl font-black text-gray-800 hover:bg-gray-50 transition-colors">TUTUP</button>
                </div>
            </div>
        </div>

    </AppLayout>
</template>

<style>
/* Custom Scrollbar biar estetik */
.custom-scrollbar::-webkit-scrollbar { width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #94a3b8; }

/* Sembunyikan tombol saat mencetak */
@media print {
    body * { visibility: hidden; }
    #print-area, #print-area * { visibility: visible; }
    #print-area { position: absolute; left: 0; top: 0; width: 58mm; padding: 0; margin: 0; }
    @page { margin: 0; }
    .no-print { display: none !important; }
}
</style>
