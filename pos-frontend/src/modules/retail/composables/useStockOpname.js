import { ref, onBeforeUnmount, onMounted } from 'vue';
import api from '../../../api.js';
import Swal from 'sweetalert2';
import { Html5Qrcode } from "html5-qrcode";

export function useStockOpname() {
    const activeTab = ref('SO'); 
    const soStep = ref('COUNTING'); 
    
    // 🚀 DETEKSI ROLE LOGIN
    const role = localStorage.getItem('role') || 'staff';
    const isOwner = role.toLowerCase() === 'owner';

    // 🚀 STATE TIME-LOCK KLAIM BARANG
    const isKlaimEligible = ref(false);
    const daysLeftKlaim = ref(0);

    const notes = ref(`Stock Opname ${new Date().toLocaleDateString('id-ID')}`);
    const searchQuery = ref('');
    const products = ref([]);
    const cartSO = ref([]);
    const cartKlaim = ref([]);
    const isSubmitting = ref(false);

    const showScanner = ref(false);
    let html5QrCode = null;

    // --- 🚀 LOGIKA CEK ELIGIBILITY KLAIM (7 HARI) ---
    const checkKlaimEligibility = async () => {
        try {
            const res = await api.get('/retail/stock-opname/last-status');

            // 🐛 RADAR DEBUGGING: Cek apa yang dikirim Backend
            console.log("👉 DATA DARI BACKEND:", res.data);

            if (res.data && res.data.last_so_date) {
                const lastSO = new Date(res.data.last_so_date);
                const now = new Date();
                
                // Pake Math.abs biar aman dari selisih milidetik timezone
                const diffTime = Math.abs(now - lastSO); 
                const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

                // 🐛 RADAR DEBUGGING: Cek hasil hitungan
                console.log("🗓️ Tanggal SO:", lastSO.toLocaleDateString());
                console.log("🗓️ Tanggal Skrg:", now.toLocaleDateString());
                console.log("⏳ Selisih Hari:", diffDays);

                // Syarat: Selisih harus 7 hari atau kurang
                if (diffDays <= 7) {
                    isKlaimEligible.value = true;
                    // Kalau diff 0 (baru aja SO hari ini), berarti sisa 7 hari
                    daysLeftKlaim.value = diffDays === 0 ? 7 : (8 - diffDays); 
                } else {
                    console.log("❌ Backend ngirim null / gak ada tanggal");
                    isKlaimEligible.value = false;
                }
            } else {
                isKlaimEligible.value = false;
            }
        } catch (e) {
            console.error("Gagal cek status Klaim:", e);
            isKlaimEligible.value = false; 
        }
    };

    // Jalankan pengecekan saat komponen dimuat
    onMounted(() => {
        checkKlaimEligibility();
    });

    const startScanner = async () => {
        showScanner.value = true;
        setTimeout(async () => {
            try {
                html5QrCode = new Html5Qrcode("reader-so");
                await html5QrCode.start(
                    { facingMode: "environment" }, 
                    { fps: 10, qrbox: { width: 250, height: 100 } }, 
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

    let searchTimer = null;
    const searchProduct = (isFromScanner = false) => {
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
                    Swal.fire({ toast: true, position: 'top-end', icon: 'success', title: `${foundData[0].nama_produk} ditambahkan!`, showConfirmButton: false, timer: 1500 });
                } else {
                    Swal.fire('Tidak Ditemukan!', 'Barcode ini tidak ada di Master.', 'error');
                }
                searchQuery.value = '';
            } else { products.value = foundData; }
        } catch (err) { console.error(err); }
    };

    const addToCart = (product) => {
        const targetCart = activeTab.value === 'SO' ? cartSO : cartKlaim;
        const existing = targetCart.value.find(item => item.product_id === product.id);
        
        if (existing) {
            existing.actual_qty++;
        } else {
            targetCart.value.unshift({
                product_id: product.id,
                nama_produk: product.nama_produk,
                sku: product.sku,
                system_qty: product.stok,
                actual_qty: activeTab.value === 'SO' ? 0 : 1, 
                alasan: ''
            });
        }
        searchQuery.value = '';
        products.value = [];
    };

    const removeItem = (index) => {
        activeTab.value === 'SO' ? cartSO.value.splice(index, 1) : cartKlaim.value.splice(index, 1);
    };

    // 🚀 ALUR BARU: TARIK BARANG YANG GAK KE-SCAN
    const proceedToReview = async () => {
        if (cartSO.value.length === 0) return Swal.fire('Kosong', 'Belum ada barang yang dihitung!', 'warning');
        
        const belumDihitung = cartSO.value.some(i => i.actual_qty === 0 && i.system_qty > 0);
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
            const res = await api.get('/retail/products?limit=10000');
            const allProducts = res.data.data || [];

            allProducts.forEach(p => {
                const isScanned = cartSO.value.find(item => item.product_id === p.id);
                if (!isScanned && p.stok > 0) {
                    cartSO.value.push({
                        product_id: p.id,
                        nama_produk: p.nama_produk,
                        sku: p.sku,
                        system_qty: p.stok,
                        actual_qty: 0, // PAKSA JADI 0
                        alasan: '' 
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
            await api.post('/retail/stock-opname', {
                notes: notes.value,
                items: cartSO.value.map(i => ({
                    product_id: i.product_id,
                    system_qty: i.system_qty,
                    actual_qty: i.actual_qty,
                    selisih: i.actual_qty - i.system_qty,
                    alasan: i.alasan
                })),
                status: isOwner ? 'APPROVED' : 'PENDING_APPROVAL'
            });

            Swal.fire('Selesai!', 'Data SO berhasil diproses.', 'success');
            cartSO.value = [];
            soStep.value = 'COUNTING';
            checkKlaimEligibility(); // 🚀 Refresh status klaim setelah SO berhasil
        } catch (err) { Swal.fire('Gagal', 'Terjadi kesalahan sistem', 'error'); } 
        finally { isSubmitting.value = false; }
    };

    const submitKlaimTemuan = async () => {
        if (cartKlaim.value.length === 0) return Swal.fire('Kosong', 'Keranjang klaim kosong', 'warning');
        
        // Alasan wajib buat klaim barang nyempil
        const kosong = cartKlaim.value.some(i => !i.alasan);
        if (kosong) return Swal.fire('Wajib Diisi', 'Keterangan barang temuan wajib diisi!', 'warning');

        isSubmitting.value = true;
        try {
            await api.post('/retail/stock-adjustment/request', {
                notes: 'Klaim Barang Nyempil',
                items: cartKlaim.value.map(i => ({ product_id: i.product_id, qty: i.actual_qty, alasan: i.alasan }))
            });
            Swal.fire('Terkirim!', 'Klaim barang telah dikirim ke Owner untuk di-Approve.', 'success');
            cartKlaim.value = [];
        } catch (err) { Swal.fire('Gagal', 'Sistem gagal memproses', 'error'); } 
        finally { isSubmitting.value = false; }
    };

    return {
        activeTab, soStep, notes, searchQuery, products, cartSO, cartKlaim, isSubmitting, showScanner, isOwner,
        isKlaimEligible, daysLeftKlaim, // 🚀 Wajib di-export ke komponen
        startScanner, stopScanner, searchProduct, addToCart, removeItem, 
        proceedToReview, submitSOFinal, submitKlaimTemuan
    };
}