import { ref, onMounted, onUnmounted, computed, watch, nextTick } from 'vue';
import api from '../../../api.js';
import Swal from 'sweetalert2';
import { Html5Qrcode } from "html5-qrcode";

export function useReturBarang() {
    // --- STATE DATA ---
    const products = ref([]);
    const cart = ref([]); 
    const isLoading = ref(true);
    const isSubmitting = ref(false);

    // 🚀 STATE CETAK LANGSUNG
    const lastSubmittedReturn = ref(null); 

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
            const resProd = await api.get('/retail/products');
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
                selectedCamera.value = devices[1]?.id || devices[0]?.id; // Fallback kalau kamera cuma 1
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

            const res = await api.post('/retail/returns', payload);
            
            lastSubmittedReturn.value = {
                return_no: res.data.return_no,
                created_at: new Date(),
                items: [...cart.value],
                total_qty: cart.value.reduce((acc, curr) => acc + curr.qty, 0),
                user: { name: localStorage.getItem('name') || 'Kasir' },
                storeName: localStorage.getItem('storeName') || 'POS UMKM'
            };

            cart.value = [];
            fetchProducts(); 

            const resultPrint = await Swal.fire({
                icon: 'success',
                title: 'Berhasil Diproses!',
                html: `Dokumen: <b>${res.data.return_no}</b><br/>Stok telah berhasil dipotong dari sistem.`,
                showCancelButton: true,
                confirmButtonText: '<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 inline-block mr-1" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="6 9 6 2 18 2 18 9"/><path d="M6 18H4a2 2 0 0 1-2-2v-5a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v5a2 2 0 0 1-2 2h-2"/><rect width="12" height="8" x="6" y="14"/></svg> Cetak Bukti',
                cancelButtonText: 'Tutup',
                confirmButtonColor: '#4f46e5',
                cancelButtonColor: '#94a3b8',
                customClass: {
                    confirmButton: 'rounded-xl font-black px-6 py-3 flex items-center gap-2',
                    cancelButton: 'rounded-xl font-black px-6 py-3'
                }
            });

            if (resultPrint.isConfirmed) {
                setTimeout(() => { window.print(); }, 300);
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

    return {
        products, cart, isLoading, isSubmitting, lastSubmittedReturn,
        searchProductQuery, isDropdownOpen, selectedProduct, isScannerOpen,
        cameras, selectedCamera, form, alasanOptions, filteredProducts,
        startScanner, stopScanner, switchCamera, selectProduct, clearSelectedProduct,
        addToCart, removeFromCart, submitBatchReturn, getBadgeClass
    };
}