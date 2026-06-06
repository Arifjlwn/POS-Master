import { ref, onBeforeUnmount, onMounted, computed } from 'vue'; 
import api from '../../../api.js';
import Swal from 'sweetalert2';
import { Html5Qrcode, Html5QrcodeSupportedFormats } from "html5-qrcode";

export function useStockOpname() {
    const activeTab = ref('SO'); 
    const soStep = ref('COUNTING'); 
    
    // DETEKSI ROLE LOGIN KASIR
    const role = localStorage.getItem('role') || 'staff';
    const isOwner = role.toLowerCase() === 'owner';

    // STATE TIMELOCK CONTROL
    const isKlaimEligible = ref(false);
    const daysLeftKlaim = ref(0);
    const isSOLockedThisMonth = ref(false); 

    const notes = ref(`Stock Opname ${new Date().toLocaleDateString('id-ID')}`);
    const searchQuery = ref('');
    const products = ref([]);
    const cartSO = ref([]);
    const cartKlaim = ref([]);
    const isSubmitting = ref(false);

    const showScanner = ref(false);
	let html5QrCode = null;

    // --- 🚀 AUDIT ELIGIBILITY TIMELOCK CHECKER ---
    const checkKlaimEligibility = async () => {
        try {
            const res = await api.get('/retail/stock-opname/last-status');

            if (res.data && res.data.last_so_date) {
                const lastSO = new Date(res.data.last_so_date);
                const now = new Date();
                const hasClaimed = res.data.has_claimed; 
                
                // Hitung sisa kalender hari secara presisi 
                const diffTime = now.getTime() - lastSO.getTime(); 
                const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));

                if (diffDays <= 7 && !hasClaimed) {
                    isKlaimEligible.value = true;
                    isKlaimEligible.value = true;
                    daysLeftKlaim.value = (7 - diffDays) <= 0 ? 1 : (7 - diffDays); 
                    fetchMinusItems(); 
                } else {
                    isKlaimEligible.value = false;
                }

                // Gembok otomatis jika bulan dan tahun SO terakhir sama dengan bulan berjalan 
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
            console.error("Gagal verifikasi status timelock Klaim/SO:", e);
        }
    };

    const fetchMinusItems = async () => {
        try {
            const res = await api.get('/retail/stock-opname/last-minus');
            const minusData = res.data.data || [];
            
            cartKlaim.value = minusData.map(d => {
                const p = d.product || d.Product || d;
                return {
                    product_id: p.id,
                    nama_produk: p.nama_produk || p.name,
                    sku: p.sku || p.SKU,
                    system_qty: p.stok || p.stock || 0,
                    max_klaim: Math.abs(d.selisih || d.Selisih || 0), 
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
                };
            });
        } catch (e) { 
            console.error("Gagal menarik daftar produk minus historis:", e); 
        }
    };

    onMounted(() => {
        checkKlaimEligibility();
    });

    // --- 🛡️ FIXED CORE ENGINE SCANNER DEVICE MANAGEMENT ---
    const startScanner = async () => {
        if (activeTab.value === 'SO' && isSOLockedThisMonth.value) {
            return Swal.fire('Akses Ditolak', 'Sistem mengunci Stock Opname baru untuk bulan ini .', 'warning');
        }
        showScanner.value = true;
        setTimeout(async () => {
            try {
                // FIX ID LINKING: Ditarget murni ke ID "reader" agar klop dengan div StockOpname.vue  !
                html5QrCode = new Html5Qrcode("reader");
                await html5QrCode.start(
                    { facingMode: "environment" }, 
                    { 
                        fps: 15, 
                        qrbox: { width: 280, height: 160 },
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
                console.error(err);
                Swal.fire('Camera Error', 'Gagal memicu hardware stream kamera perangkat.', 'error');
                stopScanner();
            }
        }, 300);
    };

    const stopScanner = () => {
        if (html5QrCode) {
            html5QrCode.stop().then(() => { 
                html5QrCode.clear(); 
                showScanner.value = false; 
            }).catch(() => { 
                showScanner.value = false; 
            });
        } else { 
            showScanner.value = false; 
        }
    };

    onBeforeUnmount(() => { if (showScanner.value) stopScanner(); });

    // --- AUTOCOMPLETE SEARCH AUTOFILTER ---
    let searchTimer = null;
    const searchProduct = (isFromScanner = false) => {
        if (activeTab.value === 'SO' && isSOLockedThisMonth.value) {
            Swal.fire('Terkunci!', 'Daftar audit stock opname terkunci bulan ini.', 'warning');
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
            const foundData = res.data.data || res.data || [];

            if (isFromScanner) {
                if (foundData.length > 0) {
                    addToCart(foundData[0]);
                    Swal.fire({ toast: true, position: 'top-end', icon: 'success', title: `${foundData[0].nama_produk} masuk antrean!`, showConfirmButton: false, timer: 1500 });
                } else {
                    Swal.fire('Barcode Gaib', 'Nomor SKU Barcode tidak terdaftar di sistem master data .', 'error');
                }
                searchQuery.value = '';
            } else { 
                products.value = foundData; 
            }
        } catch (err) { 
            console.error(err); 
        }
    };

    const addToCart = (product) => {
        if (activeTab.value === 'SO' && isSOLockedThisMonth.value) {
            return Swal.fire('Terkunci!', 'Sistem audit opname bulan berjalan telah dikunci.', 'warning');
        }

        const targetCart = activeTab.value === 'SO' ? cartSO : cartKlaim;
        const existing = targetCart.value.find(item => item.product_id === product.id);
        
        if (existing) {
            existing.qty_dasar++;
        } else {
			const p = product;
            targetCart.value.unshift({
                product_id: p.id,
                nama_produk: p.nama_produk || p.name,
                sku: p.sku || p.SKU,
                system_qty: p.stok ?? p.stock ?? 0,
                alasan: '',
                qty_besar: 0,
                qty_tengah: 0,
                qty_dasar: activeTab.value === 'SO' ? 0 : 1, 
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
        searchQuery.value = '';
        products.value = [];
    };

    const removeItem = (index) => {
        activeTab.value === 'SO' ? cartSO.value.splice(index, 1) : cartKlaim.value.splice(index, 1);
    };

    // --- MATHEMATICS UOM RECONCILIATION LAYER ---
    const hitungTotalFisik = (item) => {
        const qBesar = parseInt(item.qty_besar, 10) || 0;
        const qTengah = parseInt(item.qty_tengah, 10) || 0;
        const qDasar = parseInt(item.qty_dasar, 10) || 0;

        if (item.is_nested) {
            const stokDariBesar = qBesar * (parseInt(item.isi_per_besar, 10) || 1); 
            const stokDariTengah = qTengah * (parseInt(item.isi_tengah_ke_dasar, 10) || 1); 
            return stokDariBesar + stokDariTengah + qDasar;
        } else {
            const stokDariBesar = qBesar * (parseInt(item.isi_per_besar, 10) || 1);
            return stokDariBesar + qDasar;
        }
    };

    // --- REVIEW FLOW SYSTEM CONTROL ---
    const proceedToReview = async () => {
        if (cartSO.value.length === 0) return Swal.fire('Kosong', 'Daftar laci hitung stock opname masih kosong melompong !', 'warning');
        
        const belumDihitung = cartSO.value.some(i => hitungTotalFisik(i) === 0 && i.system_qty > 0);
        if (belumDihitung) return Swal.fire('Audit Gantung', 'Wajib isi kuantitas fisik riil rak. Jika barang emang hilang, isi angka 0 !', 'info');

        const confirm = await Swal.fire({
            title: 'Kunci Mode Scan Audit?',
            html: "Daftar akan dikunci .<br><br><b class='text-rose-600'>Barang master retail lainnya yang tidak masuk daftar scan otomatis akan digilas ditarik menjadi 0 (Dianggap Hilang Toko)!</b>",
            icon: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#4f46e5',
            cancelButtonText: 'Batal',
            confirmButtonText: 'Ya, Hitung Selisih!'
        });

        if (!confirm.isConfirmed) return;

        Swal.fire({ title: 'Menyisir Selisih Data Master DB...', allowOutsideClick: false, didOpen: () => Swal.showLoading() });
        
        try {
            // SAFE PAGINATION MITIGATION: Gunakan limit rasional  biar browser kasir ga mendadak Force Close RAM leak!
            const res = await api.get(`/retail/products?limit=1500`);
            const allProducts = res.data.data || res.data || [];

            allProducts.forEach(p => {
                const isScanned = cartSO.value.find(item => item.product_id === p.id);
                if (!isScanned && (p.stok || p.stock) > 0) {
                    cartSO.value.push({
                        product_id: p.id,
                        nama_produk: p.nama_produk || p.name,
                        sku: p.sku || p.SKU,
                        system_qty: p.stok || p.stock || 0,
                        alasan: 'Tidak Ditemukan Saat Audit Fisik Rak',
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
            Swal.fire('Gateway Timeout', 'Gagal melakukan sinkronisasi pencocokan silisih barang dengan cloud .', 'error');
        }
    };

    const submitSOFinal = async () => {
        if (isSOLockedThisMonth.value) {
            return Swal.fire('Akses Ditolak!', 'Stock Opname bulanan cabang retail Anda sudah terkunci .', 'error');
        }

        const confirm = await Swal.fire({
			title: isOwner ? 'Finalisasi & Eksekusi Update?' : 'Kirim Berkas Pengajuan SO?',
            text: isOwner ? 'Angka selisih variance akan langsung membakar merubah stok master DB permanen !' : 'Berkas audit akan dikirim ke laptop Owner untuk divalidasi.',
            icon: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#e11d48',
			cancelButtonText: 'Cek Lagi'
        });

        if (!confirm.isConfirmed) return;

        isSubmitting.value = true;
        try {
            const payloadItems = cartSO.value.map(i => {
                const totalFisik = hitungTotalFisik(i);
                return {
                    product_id: parseInt(i.product_id, 10),
                    system_qty: parseInt(i.system_qty, 10),
                    actual_qty: totalFisik,
                    selisih: totalFisik - parseInt(i.system_qty, 10),
                    alasan: i.alasan || 'Stock Opname Rutin Bulanan'
                };
            });

            await api.post('/retail/stock-opname', {
                notes: notes.value,
                items: payloadItems,
                status: isOwner ? 'APPROVED' : 'PENDING_APPROVAL'
            });

            Swal.fire('Sukses!', 'Dokumen berkas audit SO sukses diproses server!', 'success');
            cartSO.value = [];
            soStep.value = 'COUNTING';
            await checkKlaimEligibility(); 
        } catch (err) { 
			Swal.fire('Gagal Simpan', err.response?.data?.error || 'Koneksi database backend terputus.', 'error'); 
		} finally { 
			isSubmitting.value = false; 
		}
    };

    const submitKlaimTemuan = async () => {
        if (cartKlaim.value.length === 0) return Swal.fire('Kosong', 'Tidak ada daftar barang koreksi negatif yang bisa diklaim .', 'warning');
        
        const overclaim = cartKlaim.value.find(i => hitungTotalFisik(i) > i.max_klaim);
        if (overclaim) {
            return Swal.fire('Klaim Ditolak', `Barang ${overclaim.nama_produk} limit kehilangan maks adalah ${overclaim.max_klaim} PCS! Angka fisik yang Anda masukkan over-limit !`, 'error');
        }

        const payloadItems = cartKlaim.value
            .filter(i => hitungTotalFisik(i) > 0)
            .map(i => ({ 
                product_id: parseInt(i.product_id, 10), 
                qty: hitungTotalFisik(i), 
                alasan: i.alasan 
            }));

        if (payloadItems.length === 0) return Swal.fire('Info', 'Isi kuantitas fisik barang temuan nyempil minimal 1 item !', 'info');

        const kosong = payloadItems.some(i => !i.alasan || i.alasan.trim() === "");
        if (kosong) return Swal.fire('Keterangan Wajib', 'Wajib isi kolom keterangan kronologi/lokasi penemuan barang nyempil !', 'warning');

        isSubmitting.value = true;
        try {
            await api.post('/retail/stock-adjustment/request', {
                notes: 'Koreksi Penyesuaian Klaim Barang Hasil Stock Opname',
                items: payloadItems
            });
            Swal.fire('Pengajuan Dikirim!', 'Berkas klaim penyesuaian stok temuan berhasil diteruskan ke komputer Owner !', 'success');
            cartKlaim.value = [];
            await checkKlaimEligibility(); 
        } catch (err) { 
			Swal.fire('Gagal Simpan', 'Server gagal memproses kuitansi penyesuaian penemuan .', 'error'); 
		} finally { 
			isSubmitting.value = false; 
		}
    };

    // 🚀 BALIKIN KE NEGARA ASAL , BERSIH DARI HURUF SILUMAN!
    return {
        activeTab, soStep, notes, searchQuery, products, cartSO, cartKlaim, isSubmitting, showScanner, isOwner,
        isKlaimEligible, daysLeftKlaim, isSOLockedThisMonth, 
        startScanner, stopScanner, searchProduct, addToCart, removeItem, hitungTotalFisik, 
        proceedToReview, submitSOFinal, submitKlaimTemuan
    };
}