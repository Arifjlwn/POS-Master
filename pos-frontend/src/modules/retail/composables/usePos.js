import { ref, computed, onMounted, nextTick, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import { posService } from '../services/posService.js';
import Swal from 'sweetalert2';
import { Html5Qrcode } from "html5-qrcode";

export function usePos() {
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

    // STATE KHUSUS UNTUK HP (MOBILE CART DRAWER)
    const isMobileCartOpen = ref(false);

    const getImageUrl = (path) => {
        if (!path) return null;
        return `${import.meta.env.VITE_API_BASE_URL}${path}`;
    };

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

    // --- LOGIKA DATA PRODUK & PENCARIAN ---
    const fetchProducts = async () => {
        try {
            const response = await posService.getProducts();
            products.value = response.data.map(p => ({
                id: p.id,
                sku: p.sku || `SKU-${p.id}`,
                name: p.nama_produk,
                price: p.harga_jual,
                stock: p.stok,
                image: p.gambar,
                satuan_dasar: p.satuan_dasar || 'PCS',
                satuan_besar: p.satuan_besar || null,
                isi_per_besar: p.isi_per_besar || 0,
                harga_jual_besar: p.harga_jual_besar || 0
            }));
        } catch (error) {
            console.error("Gagal narik produk:", error);
        } finally {
            isLoadingProducts.value = false;
        }
    };

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
        nextTick(() => { if (searchInput.value) searchInput.value.focus(); });
    };

    // --- LOGIKA KERANJANG BELANJA ---
    const addToCart = (product) => {
        if (product.stock <= 0) {
            Swal.fire({ icon: 'error', title: 'Stok Habis!', text: `Stok ${product.name} kosong` });
            return;
        }

        const defaultUom = product.satuan_dasar;
        const existingItem = cart.value.find(item => item.id === product.id && item.selected_uom === defaultUom);
        
        if (existingItem) {
            // Cek stok berdasarkan multiplier (Kalau satuan dasar, multiplier = 1)
            if ((existingItem.qty + 1) * existingItem.uom_multiplier <= product.stock) {
                existingItem.qty++;
            } else {
                Swal.fire({ icon: 'warning', title: 'Stok Terbatas', text: 'Kuantitas melebihi stok!' });
                return;
            }
        } else {
            // 🚀 Bikin properti canggih buat menampung status konversi di keranjang
            cart.value.unshift({ 
                id: product.id, 
                name: product.name, 
                price: product.price, // Harga aktif saat ini
                qty: 1, 
                
                // Variabel Sistem Konversi
                selected_uom: defaultUom, // Saat masuk pertama, selalu pakai eceran
                uom_multiplier: 1, // Pengali potong stok
                has_grosir: !!product.satuan_besar, 
                satuan_dasar: product.satuan_dasar,
                satuan_besar: product.satuan_besar,
                isi_per_besar: product.isi_per_besar,
                harga_jual_besar: product.harga_jual_besar,
                harga_dasar: product.price
            });
        }
        
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

    // 🚀 FUNGSI BARU: TOGGLE UBAH SATUAN ECERAN -> GROSIR
    const toggleUom = (item) => {
        if (!item.has_grosir) return; // Kalau ga punya satuan besar, cuekin aja
        
        const prodMaster = products.value.find(p => p.id === item.id);
        
        if (item.selected_uom === item.satuan_dasar) {
            // Mau diubah ke GROSIR (Misal PCS -> KARTON)
            // Cek dulu stoknya cukup ga buat 1 karton?
            if ((item.qty * item.isi_per_besar) > prodMaster.stock) {
                Swal.fire({ icon: 'error', title: 'Stok Tidak Cukup!', text: `Stok fisik tidak cukup untuk dijual dalam bentuk ${item.satuan_besar}` });
                return;
            }
            item.selected_uom = item.satuan_besar;
            item.uom_multiplier = item.isi_per_besar;
            item.price = item.harga_jual_besar; // Pake harga spesial grosir
        } else {
            // Kembalikan ke ECERAN (Misal KARTON -> PCS)
            item.selected_uom = item.satuan_dasar;
            item.uom_multiplier = 1;
            item.price = item.harga_dasar; // Pake harga normal
        }
    };

    const decreaseQty = (item) => {
        // Harus bener ngecek item spesifik beserta satuannya
        const existingItem = cart.value.find(i => i.id === item.id && i.selected_uom === item.selected_uom);
        if (existingItem) {
            if (existingItem.qty > 1) {
                existingItem.qty--;
            } else {
                cart.value = cart.value.filter(i => !(i.id === item.id && i.selected_uom === item.selected_uom));
                if (cart.value.length === 0) isMobileCartOpen.value = false;
            }
        }
    };

    const increaseQty = (item) => {
        const existingItem = cart.value.find(i => i.id === item.id && i.selected_uom === item.selected_uom);
        const prodMaster = products.value.find(p => p.id === item.id);
        if (existingItem && prodMaster) {
            // Cek stok berdasarkan hitungan konversi (Misal 2 Slop * 192 Batang)
            if ((existingItem.qty + 1) * existingItem.uom_multiplier <= prodMaster.stock) {
                existingItem.qty++;
            } else {
                Swal.fire({ icon: 'warning', title: 'Stok Terbatas', text: 'Kuantitas melebihi stok!' });
            }
        }
    };

    const validateQty = (item) => {
        const existingItem = cart.value.find(i => i.id === item.id && i.selected_uom === item.selected_uom);
        if (existingItem && existingItem.qty < 1) existingItem.qty = 1;
        
        // Proteksi kalau input manualnya bablas
        const prodMaster = products.value.find(p => p.id === item.id);
        if (existingItem && prodMaster && (existingItem.qty * existingItem.uom_multiplier > prodMaster.stock)) {
            Swal.fire({ icon: 'warning', title: 'Overstock!', text: 'Kuantitas dikembalikan ke batas stok' });
            existingItem.qty = Math.floor(prodMaster.stock / existingItem.uom_multiplier);
        }
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

    // --- LOGIKA HOLD & RESUME TRANSAKSI ---
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

    // --- PROSES CHECKOUT ---
    const isProcessingCheckout = ref(false);

    const executeCheckout = async() => {
        if (isProcessingCheckout.value) return;

        isProcessingCheckout.value = true;
        
        // 🚀 BAGIAN PALING KRUSIAL: Mapping stok yang dikirim ke Backend!
        const payloadItems = cart.value.map(item => ({ 
            product_id: item.id, 
            // Kuantitas yang dikirim ke DB adalah TOTAL SATUAN DASAR! (Biar stok dipotong bener)
            // Misal beli 1 Slop (isi 12) -> Kirim 12.
            kuantitas: item.qty * item.uom_multiplier 
        }));

        try {
            const response = await posService.checkout({
                session_id: currentSession.value.id,
                items: payloadItems,
                nominal_bayar: payAmount.value,
                metode_bayar: paymentMethod.value
            });

            lastTransaction.value = {
                invoice: response.invoice, 
                cart: [...cart.value],
                total: response.tagihan, 
                pay: payAmount.value,
                return: response.kembali, 
                method: paymentMethod.value,
                subtotal: subTotalBelanja.value, 
                pajak: nilaiPajak.value,          
                date: new Date().toLocaleString('id-ID', { year: '2-digit', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' }).replace(/\//g, '.')
            };
            
            showQrisModal.value = false;
            isMobileCartOpen.value = false; 
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

    // --- LOGIKA INITIAL OPEN SESSION (UNTUK KASIRAN BARU) ---
    const openSession = async (stationNumber, modalAwalValue) => {
        try {
            const res = await posService.openSession({
                station_number: stationNumber,
                modal_awal: parseFloat(modalAwalValue)
            });
            return res;
        } catch (error) {
            throw error;
        }
    };

    // --- LOGIKA CLOSING SHIFT ---
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
            const res = await posService.closeSession(currentSession.value.id, {
                total_aktual: totalUangFisik.value,
                pecahan: pecahan.value
            });
            Swal.fire('Closing Berhasil!', 'Struk closing akan dicetak.', 'success');
            lastClosingData.value = res; 
            showClosingModal.value = false;
            showReceiptClosing.value = true;
        } catch (error) {
            Swal.fire('Gagal Closing', error.response?.data?.error, 'error');
        }
    };

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

    // --- LIFECYCLE HOOKS ---
    onMounted(async () => {
        const token = localStorage.getItem('token');
        if (!token) { router.push('/login'); return; }

        try {
            const res = await posService.checkSession(token);
            if (!res.has_session) {
                Swal.fire('Akses Ditolak', 'Isi modal awal atau absen dulu ya!', 'warning');
                router.push('/retail/pos/buka-kasir');
                return;
            }

            currentSession.value = res.session;
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

    onUnmounted(() => {
        if (showScanner.value) stopScanner();
        clearInterval(timer);
    });

    // PENTING: EXPORT SEMUA VARIABEL & FUNGSI UNTUK VIEW
    return {
        currentUser, 
        currentSession, 
        currentTime, 
        products, 
        isLoadingProducts, 
        cart, 
        heldOrders,
        showHeldModal, 
        payAmount, 
        paymentMethod, 
        showReceipt, 
        showQrisModal, 
        lastTransaction,
        showReceiptClosing, 
        lastClosingData, 
        isMobileCartOpen, 
        searchQuery, 
        searchInput,
        showScanner, 
        pecahan, 
        totalUangFisik, 
        filteredProducts, 
        subTotalBelanja, 
        nilaiPajak,
        totalBelanja, 
        kembalian, 
        isProcessingCheckout,
        showClosingModal,  // 🟢 WAJIB ADA DI SINI BIAR KEBACA DI VIEW!
        getImageUrl, 
        startScanner, 
        stopScanner, 
        handleBarcodeScan, 
        addToCart, 
        decreaseQty, 
        increaseQty, 
        validateQty, 
        clearCart, 
        holdTransaction, 
        resumeOrder,
        setPaymentMethod, 
        executeCheckout, 
        formatInputRupiah,
        processCheckout, 
        handleClosing, 
        logout,
        toggleUom
    };
}