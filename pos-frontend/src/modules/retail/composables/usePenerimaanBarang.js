import { ref, computed, onBeforeUnmount } from 'vue';
import api from '../../../api.js';
import Swal from 'sweetalert2';
import { Html5Qrcode, Html5QrcodeSupportedFormats } from "html5-qrcode";

export function usePenerimaanBarang() {
    const getPayloadFromToken = () => {
        const token = localStorage.getItem('token');
        const role = localStorage.getItem('role') || 'staff';
        if (!token) return { role: 'staff' };
        try { return { role: role.toLowerCase() }; } 
        catch (e) { return { role: 'staff' }; }
    };
    
    const currentUser = ref(getPayloadFromToken());
    const isOwner = computed(() => currentUser.value.role === 'owner');

    const supplierName = ref('');
    const noFaktur = ref('');
    const searchQuery = ref('');
    const products = ref([]); 
    const cartLPB = ref([]); 
    const isSubmitting = ref(false);

    // --- LOGIKA SCANNER KOKOH ---
    const showScanner = ref(false);
    let html5QrCode = null;

    const startScanner = async () => {
        showScanner.value = true;
        setTimeout(async () => {
            try {
                html5QrCode = new Html5Qrcode("reader-lpb");
                await html5QrCode.start(
                    { facingMode: "environment" }, 
                    { 
                        fps: 15, 
                        qrbox: { width: 320, height: 150 },
                        formatsToSupport: [
                            Html5QrcodeSupportedFormats.EAN_13,
                            Html5QrcodeSupportedFormats.EAN_8,
                            Html5QrcodeSupportedFormats.UPC_A,
                            Html5QrcodeSupportedFormats.UPC_E,
                            Html5QrcodeSupportedFormats.CODE_128
                        ]
                    }, 
                    (decodedText) => {
                        searchQuery.value = decodedText; 
                        stopScanner();
                        searchProduct(true); // Auto-add mode
                        const audio = new Audio('https://www.soundjay.com/buttons/beep-07a.mp3');
                        audio.play().catch(()=>{}); 
                    },
                    () => {} 
                );
            } catch (err) {
                console.error(err);
                Swal.fire('Oops!', 'Gagal menyalakan kamera scanner.', 'error');
                stopScanner();
            }
        }, 400);
    };

    const stopScanner = () => {
        if (html5QrCode) {
            html5QrCode.stop().then(() => {
                html5QrCode.clear();
                showScanner.value = false;
            }).catch(() => { showScanner.value = false; });
        } else {
            showScanner.value = false;
        }
    };

    onBeforeUnmount(() => { if (showScanner.value) stopScanner(); });

    // --- PENCARIAN & KERANJANG ---
    let searchTimer = null;
    const searchProduct = async (isFromScanner = false) => {
        clearTimeout(searchTimer);
        if (isFromScanner) {
            return executeSearch(true);
        }
        searchTimer = setTimeout(() => {
            executeSearch(false);
        }, 300);
    };

    // 🚀 Fungsi Inti Pencarian (Biar kode lebih DRY - Don't Repeat Yourself)
    const executeSearch = async (isFromScanner) => {
        if (!isFromScanner && searchQuery.value.length < 2) {
            products.value = [];
            return;
        }

        try {
            const res = await api.get(`/retail/products?search=${searchQuery.value}`);
            const foundData = res.data.data || [];

            if (isFromScanner) {
                if (foundData.length > 0) {
                    addToCart(foundData[0]); // Ambil hasil pertama
                    Swal.fire({ 
                        toast: true, 
                        position: 'top-end', 
                        icon: 'success', 
                        title: `${foundData[0].nama_produk} ditambahkan!`, 
                        showConfirmButton: false, 
                        timer: 1500 
                    });
                } else {
                    Swal.fire({
                        icon: 'error',
                        title: 'Oops!',
                        text: 'Barcode tidak terdaftar di Master Produk.',
                        timer: 2000
                    });
                }
                searchQuery.value = ''; // Reset input setelah scan
                products.value = [];    // Bersihkan hasil dropdown
            } else {
                products.value = foundData;
            }
        } catch (err) { 
            console.error("Gagal cari produk:", err);
            products.value = [];
        }
    };

    // 🚀 MASUKIN KE KERANJANG PENERIMAAN (LPB) DENGAN STRUKTUR 3 LAPIS
    const addToCart = (product) => {
        const existing = cartLPB.value.find(item => item.product_id === product.id);
        
        if (existing) {
            existing.qty_besar++; // Default nambah kemasan paling besar
        } else {
            // Kita bikin variabel penampung harga modal yang pinter
            let defaultHargaModal = Number(product.harga_modal) || 0;
            let hargaBeliInputDefault = defaultHargaModal;

            // Kalau dia punya kemasan besar, default form LPB minta input harga 1 Dus/Slop
            if (product.satuan_besar) {
                hargaBeliInputDefault = defaultHargaModal * Number(product.isi_per_besar);
            }

            cartLPB.value.push({
                product_id: product.id,
                nama_produk: product.nama_produk,
                
                // Form Input Qty (Sekarang ada 3 kemungkinan)
                qty_besar: 1,
                qty_tengah: 0,
                qty_dasar: 0, // Dulu qty_eceran
                
                // Variabel Keuangan
                harga_jual_saat_ini: Number(product.harga_jual) || 0,
                harga_modal_database: defaultHargaModal,
                harga_beli_input: hargaBeliInputDefault, // Modal per Faktur (1 Dus / 1 Slop)
                
                // 🚀 Data 3 Lapis dari Master
                satuan_dasar: product.satuan_dasar || 'PCS',
                has_satuan_besar: !!product.satuan_besar && Number(product.isi_per_besar) > 1,
                satuan_besar: product.satuan_besar || null,
                isi_per_besar: Number(product.isi_per_besar) || 1,
                
                is_nested: product.is_nested_uom || false,
                satuan_tengah: product.satuan_tengah || null,
                isi_besar_ke_tengah: Number(product.isi_besar_ke_tengah) || 0,
                isi_tengah_ke_dasar: Number(product.isi_tengah_ke_dasar) || 0
            });
        }
        searchQuery.value = '';
        products.value = [];
    };

    // 🚀 RUMUS SILUMAN UNTUK NGITUNG TOTAL STOK DASAR DARI LPB
    const hitungTotalStok = (item) => {
        const qBesar = Number(item.qty_besar) || 0;
        const qTengah = Number(item.qty_tengah) || 0;
        const qDasar = Number(item.qty_dasar) || 0;

        if (item.is_nested) {
            // Mode 3 Lapis (Rokok: Slop -> Bungkus -> Batang)
            const stokDariBesar = qBesar * item.isi_per_besar; // 1 Slop x 160 Batang
            const stokDariTengah = qTengah * item.isi_tengah_ke_dasar; // 1 Bungkus x 16 Batang
            return stokDariBesar + stokDariTengah + qDasar;
        } else {
            // Mode 2 Lapis (Indomie: Dus -> Pcs) atau Normal
            const stokDariBesar = qBesar * item.isi_per_besar;
            return stokDariBesar + qDasar;
        }
    };

    const hitungModalPerPcs = (item) => {
        const totalDasarMasuk = hitungTotalStok(item);
        const modalInputFaktur = Number(item.harga_beli_input) || 0;
        if (totalDasarMasuk === 0) return 0;
        
        // Modal Input di form LPB itu adalah modal kemasan paling besar (kalau dia punya kemasan besar)
        if (item.has_satuan_besar) {
            // Modal Faktur / Total Isi dari kemasan paling besarnya
            // (Jadi kita dapetin HPP per satuan dasarnya)
            return Math.round(modalInputFaktur / item.isi_per_besar);
        }

        // Kalau ga punya kemasan besar, berarti modal input faktur emang untuk satuan dasar
        return Math.round(modalInputFaktur / totalDasarMasuk); 
    };

    const removeItem = (index) => cartLPB.value.splice(index, 1);

    const submitLPB = async () => {
        if (!supplierName.value || cartLPB.value.length === 0) {
            return Swal.fire('Oops!', 'Lengkapi data supplier & barang!', 'warning');
        }

        // 🚀 Cek ada yang rugi ga modal sama harga jual ecerannya?
        const adaYangRugi = isOwner.value && cartLPB.value.some(item => hitungModalPerPcs(item) >= item.harga_jual_saat_ini);
        
        const result = await Swal.fire({
            title: 'Konfirmasi LPB',
            text: adaYangRugi ? "⚠️ Peringatan: Ada modal barang yang lebih mahal dari harga jual eceran! Tetap lanjutkan?" : "Simpan data penerimaan barang ke gudang?",
            icon: adaYangRugi ? 'warning' : 'question',
            showCancelButton: true,
            confirmButtonColor: adaYangRugi ? '#ef4444' : '#2563eb',
            confirmButtonText: 'Ya, Simpan!'
        });

        if (!result.isConfirmed) return;

        isSubmitting.value = true;
        try {
            const payloadItems = cartLPB.value.map(item => ({
                product_id: item.product_id,
                qty_masuk: hitungTotalStok(item), // Stok yang disuntik ke Master Produk = Total Satuan Dasar
                harga_modal: hitungModalPerPcs(item) // Ini yang dikirim ke Backend untuk dirata-rata!
            }));

            await api.post('/retail/purchases', {
                supplier_name: supplierName.value,
                no_faktur: noFaktur.value,
                items: payloadItems
            });

            Swal.fire({
                icon: 'success', title: 'Berhasil!',
                text: isOwner.value ? 'Stok dan Modal Rata-rata di Master Produk telah diperbarui.' : 'Stok dan Data Faktur berhasil dicatat.',
                timer: 2500, showConfirmButton: false
            });

            supplierName.value = ''; noFaktur.value = ''; cartLPB.value = [];
        } catch (err) {
            Swal.fire('Gagal!', err.response?.data?.error || 'Error dari server backend!', 'error');
        } finally {
            isSubmitting.value = false;
        }
    };

    return {
        currentUser, isOwner, supplierName, noFaktur, searchQuery, products, cartLPB, isSubmitting, showScanner,
        startScanner, stopScanner, searchProduct, addToCart, hitungTotalStok, hitungModalPerPcs, removeItem, submitLPB
    };
}