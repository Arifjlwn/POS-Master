import { ref, onBeforeUnmount, onMounted, computed } from 'vue'; // 🚀 TAMBAHIN COMPUTED
import api from '../../../api.js';
import Swal from 'sweetalert2';
import { Html5Qrcode, Html5QrcodeSupportedFormats } from "html5-qrcode";

export function useStockOpname() {
    const activeTab = ref('SO'); 
    const soStep = ref('COUNTING'); 
    
    // 🚀 DETEKSI ROLE LOGIN
    const role = localStorage.getItem('role') || 'staff';
    const isOwner = role.toLowerCase() === 'owner';

    // 🚀 STATE TIME-LOCK KLAIM BARANG & GEMBOK SO BULANAN
    const isKlaimEligible = ref(false);
    const daysLeftKlaim = ref(0);
    const isSOLockedThisMonth = ref(false); // 🚀 GEMBOK BARU

    const notes = ref(`Stock Opname ${new Date().toLocaleDateString('id-ID')}`);
    const searchQuery = ref('');
    const products = ref([]);
    const cartSO = ref([]);
    const cartKlaim = ref([]);
    const isSubmitting = ref(false);

    const showScanner = ref(false);
    let html5QrCode = null;

    // --- 🚀 CEK ELIGIBILITY KLAIM & GEMBOK SO BULAN INI ---
    const checkKlaimEligibility = async () => {
        try {
            const res = await api.get('/retail/stock-opname/last-status');

            if (res.data && res.data.last_so_date) {
                const lastSO = new Date(res.data.last_so_date);
                const now = new Date();
                const hasClaimed = res.data.has_claimed; // 🚀 NANGKEP STATUS KLAIM DARI GOLANG
                
                // 1. Cek Eligibility Klaim (7 Hari & Belum Pernah Klaim)
                const diffTime = Math.abs(now - lastSO); 
                const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

                // 🚀 Kalau dalam 7 hari DAN belum pernah klaim, baru dibuka!
                if (diffDays <= 7 && !hasClaimed) {
                    isKlaimEligible.value = true;
                    daysLeftKlaim.value = diffDays === 0 ? 7 : (8 - diffDays); 
                } else {
                    isKlaimEligible.value = false;
                }

                // 2. Cek Gembok SO Bulanan
                if (lastSO.getMonth() === now.getMonth() && lastSO.getFullYear() === now.getFullYear()) {
                    isSOLockedThisMonth.value = true;
                } else {
                    isSOLockedThisMonth.value = false;
                }

            } else {
                isKlaimEligible.value = false;
                isSOLockedThisMonth.value = false;
            }
        } catch (e) {
            console.error("Gagal cek status Klaim/SO:", e);
            isKlaimEligible.value = false; 
            isSOLockedThisMonth.value = false;
        }
    };

    onMounted(() => {
        checkKlaimEligibility();
    });

    // --- 🚀 SCANNER ---
    const startScanner = async () => {
        if (activeTab.value === 'SO' && isSOLockedThisMonth.value) {
            return Swal.fire('Terkunci!', 'Stock Opname sudah dilakukan bulan ini.', 'warning');
        }
        showScanner.value = true;
        setTimeout(async () => {
            try {
                html5QrCode = new Html5Qrcode("reader-so");
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
                        executeSearch(true); 
                        const audio = new Audio('https://www.soundjay.com/buttons/beep-07a.mp3');
                        audio.play().catch(()=>{}); 
                    }, () => {} 
                );
            } catch (err) {
                Swal.fire('Oops!', 'Gagal menyalakan kamera scanner.', 'error');
                stopScanner();
            }
        }, 300);
    };

    const stopScanner = () => {
        if (html5QrCode) {
            html5QrCode.stop().then(() => { html5QrCode.clear(); showScanner.value = false; }).catch(() => { showScanner.value = false; });
        } else { showScanner.value = false; }
    };

    onBeforeUnmount(() => { if (showScanner.value) stopScanner(); });

    // --- 🚀 PENCARIAN & KERANJANG (DI-UPGRADE BUAT 3 LAPIS) ---
    let searchTimer = null;
    const searchProduct = (isFromScanner = false) => {
        if (activeTab.value === 'SO' && isSOLockedThisMonth.value) {
            Swal.fire('Terkunci!', 'Stock Opname sudah dilakukan bulan ini.', 'warning');
            searchQuery.value = '';
            return;
        }
        clearTimeout(searchTimer);
        if (isFromScanner) return executeSearch(true);
        searchTimer = setTimeout(() => executeSearch(false), 300);
    };

    const executeSearch = async (isFromScanner) => {
        if (!isFromScanner && searchQuery.value.length < 2) { products.value = []; return; }
        try {
            const res = await api.get(`/retail/products?search=${searchQuery.value}`);
            const foundData = res.data.data || [];

            if (isFromScanner) {
                if (foundData.length > 0) {
                    addToCart(foundData[0]);
                    Swal.fire({ toast: true, position: 'top-end', icon: 'success', title: `${foundData[0].nama_produk} masuk daftar!`, showConfirmButton: false, timer: 1500 });
                } else {
                    Swal.fire('Tidak Ditemukan!', 'Barcode ini tidak ada di Master.', 'error');
                }
                searchQuery.value = '';
            } else { products.value = foundData; }
        } catch (err) { console.error(err); }
    };

    const addToCart = (product) => {
        if (activeTab.value === 'SO' && isSOLockedThisMonth.value) {
            Swal.fire('Terkunci!', 'Stock Opname sudah dilakukan bulan ini.', 'warning');
            return;
        }

        const targetCart = activeTab.value === 'SO' ? cartSO : cartKlaim;
        const existing = targetCart.value.find(item => item.product_id === product.id);
        
        if (existing) {
            existing.qty_dasar++;
        } else {
            targetCart.value.unshift({
                product_id: product.id,
                nama_produk: product.nama_produk,
                sku: product.sku,
                system_qty: product.stok,
                alasan: '',
                qty_besar: 0,
                qty_tengah: 0,
                qty_dasar: activeTab.value === 'SO' ? 0 : 1, 
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

    const removeItem = (index) => {
        activeTab.value === 'SO' ? cartSO.value.splice(index, 1) : cartKlaim.value.splice(index, 1);
    };

    const hitungTotalFisik = (item) => {
        const qBesar = Number(item.qty_besar) || 0;
        const qTengah = Number(item.qty_tengah) || 0;
        const qDasar = Number(item.qty_dasar) || 0;

        if (item.is_nested) {
            const stokDariBesar = qBesar * item.isi_per_besar; 
            const stokDariTengah = qTengah * item.isi_tengah_ke_dasar; 
            return stokDariBesar + stokDariTengah + qDasar;
        } else {
            const stokDariBesar = qBesar * item.isi_per_besar;
            return stokDariBesar + qDasar;
        }
    };

    // --- 🚀 ALUR REVIEW & FINALISASI SO ---
    const proceedToReview = async () => {
        if (cartSO.value.length === 0) return Swal.fire('Kosong', 'Belum ada barang yang dihitung!', 'warning');
        
        const belumDihitung = cartSO.value.some(i => hitungTotalFisik(i) === 0 && i.system_qty > 0);
        if (belumDihitung) return Swal.fire('Perhatian', 'Pastikan semua angka fisik terisi, ada barang yang masih 0.', 'info');

        const confirm = await Swal.fire({
            title: 'Kunci Mode Scan?',
            html: "Setelah ini Anda tidak bisa scan barang baru.<br><br><b class='text-rose-600'>Barang yang tidak di-scan akan otomatis ditarik, dianggap 0 (Hilang), dan wajib dipertanggungjawabkan!</b>",
            icon: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#4f46e5',
            confirmButtonText: 'Ya, Kalkulasi Sekarang!'
        });

        if (!confirm.isConfirmed) return;

        Swal.fire({ title: 'Mengkalkulasi Selisih...', allowOutsideClick: false, didOpen: () => Swal.showLoading() });
        
        try {
            const res = await api.get(`/retail/products?limit=10000`);
            const allProducts = res.data.data || [];

            allProducts.forEach(p => {
                const isScanned = cartSO.value.find(item => item.product_id === p.id);
                if (!isScanned && p.stok > 0) {
                    cartSO.value.push({
                        product_id: p.id,
                        nama_produk: p.nama_produk,
                        sku: p.sku,
                        system_qty: p.stok,
                        alasan: '',
                        qty_besar: 0, qty_tengah: 0, qty_dasar: 0, 
                        satuan_dasar: p.satuan_dasar || 'PCS',
                        has_satuan_besar: !!p.satuan_besar && Number(p.isi_per_besar) > 1,
                        satuan_besar: p.satuan_besar || null,
                        isi_per_besar: Number(p.isi_per_besar) || 1,
                        is_nested: p.is_nested_uom || false,
                        satuan_tengah: p.satuan_tengah || null,
                        isi_besar_ke_tengah: Number(p.isi_besar_ke_tengah) || 0,
                        isi_tengah_ke_dasar: Number(p.isi_tengah_ke_dasar) || 0
                    });
                }
            });

            Swal.close();
            soStep.value = 'REVIEW'; 
        } catch (err) {
            Swal.close();
            Swal.fire('Error', 'Gagal memuat data master produk dari server.', 'error');
        }
    };

    const submitSOFinal = async () => {
        if (isSOLockedThisMonth.value) {
            return Swal.fire('Sistem Terkunci!', 'Stock Opname hanya bisa dilakukan 1x dalam sebulan.', 'error');
        }

        const confirm = await Swal.fire({
            title: isOwner ? 'Finalisasi & Kunci Stok?' : 'Kirim Pengajuan SO?',
            text: isOwner ? 'Stok akan langsung diperbarui permanen!' : 'Menunggu persetujuan Owner sebelum stok berubah.',
            icon: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#e11d48'
        });

        if (!confirm.isConfirmed) return;

        isSubmitting.value = true;
        try {
            const payloadItems = cartSO.value.map(i => {
                const totalFisik = hitungTotalFisik(i);
                return {
                    product_id: i.product_id,
                    system_qty: i.system_qty,
                    actual_qty: totalFisik,
                    selisih: totalFisik - i.system_qty,
                    alasan: i.alasan
                };
            });

            await api.post('/retail/stock-opname', {
                notes: notes.value,
                items: payloadItems,
                status: isOwner ? 'APPROVED' : 'PENDING_APPROVAL'
            });

            Swal.fire('Selesai!', 'Data SO berhasil diproses.', 'success');
            cartSO.value = [];
            soStep.value = 'COUNTING';
            checkKlaimEligibility(); // Ini bakal otomatis nge-set isSOLockedThisMonth jadi true
        } catch (err) { Swal.fire('Gagal', err.response?.data?.error || 'Terjadi kesalahan sistem', 'error'); } 
        finally { isSubmitting.value = false; }
    };

    const submitKlaimTemuan = async () => {
        if (cartKlaim.value.length === 0) return Swal.fire('Kosong', 'Keranjang klaim kosong', 'warning');
        
        const kosong = cartKlaim.value.some(i => !i.alasan);
        if (kosong) return Swal.fire('Wajib Diisi', 'Keterangan barang temuan wajib diisi!', 'warning');

        isSubmitting.value = true;
        try {
            const payloadItems = cartKlaim.value.map(i => ({ 
                product_id: i.product_id, 
                qty: hitungTotalFisik(i), 
                alasan: i.alasan 
            }));

            await api.post('/retail/stock-adjustment/request', {
                notes: 'Klaim Barang Nyempil',
                items: payloadItems
            });
            Swal.fire('Terkirim!', 'Klaim barang telah dikirim ke Owner untuk di-Approve.', 'success');
            cartKlaim.value = [];
        } catch (err) { Swal.fire('Gagal', 'Sistem gagal memproses', 'error'); } 
        finally { isSubmitting.value = false; }
    };

    return {
        activeTab, soStep, notes, searchQuery, products, cartSO, cartKlaim, isSubmitting, showScanner, isOwner,
        isKlaimEligible, daysLeftKlaim, isSOLockedThisMonth, // 🚀 EXPORT GEMBOKNYA
        startScanner, stopScanner, searchProduct, addToCart, removeItem, hitungTotalFisik, 
        proceedToReview, submitSOFinal, submitKlaimTemuan
    };
}