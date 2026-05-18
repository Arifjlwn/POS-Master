<script setup>
import { ref, onMounted, computed, nextTick } from 'vue';
import SidebarLaundry from './SidebarLaundry.vue';
import api from '../../../api.js'; 
import Swal from 'sweetalert2';

// --- STATE DATA ---
const services = ref([]);
const perfumes = ref([]); 
const cart = ref([]);
const isLoading = ref(false);
const isSubmitting = ref(false);

const searchQuery = ref('');
const customerName = ref('');
const customerPhone = ref('');
const estimasiSelesai = ref(''); 

// STATE CRM PELANGGAN
const customerResults = ref([]);
const showCustomerDropdown = ref(false);

// --- 📸 STATE MULTI-KAMERA PINTAR ---
const isCameraOpen = ref(false);
const cameraTarget = ref('');
const photoData = ref(null); 
const buktiTransferData = ref(null); 

const videoItemRef = ref(null);
const canvasItemRef = ref(null);
const videoQrisRef = ref(null);
const canvasQrisRef = ref(null);
const cameraStream = ref(null);

// --- 💳 STATE MODAL & DATA TOKO DINAMIS ---
const showQrisModal = ref(false);
const showPerfumeControlModal = ref(false); // Modal kontrol stok parfum oleh kasir
const qrisStoreUrl = ref(''); 
const storeInfo = ref({
    nama_toko: 'LAUNDRY POS',
    alamat: '',
    telepon: '',
    receipt_footer: 'Terima Kasih'
}); 

const formatNoHpCustomer = () => {
    let val = String(customerPhone.value);
    if (val.startsWith('0')) val = val.substring(1);
    if (val.startsWith('62')) val = val.substring(2);
    customerPhone.value = val;
};

const isCartOpen = ref(false);
const paymentMethod = ref('TUNAI'); 
const uangBayar = ref('');
const formattedUangBayar = computed({
    get() {
        if (!uangBayar.value) return '';
        // Ubah angka murni ke format ribuan biasa tanpa Rp dulu biar enak diketik
        return new Intl.NumberFormat('id-ID').format(uangBayar.value);
    },
    set(newValue) {
        // Bersihkan semua karakter selain angka (menghapus titik/spasi) sebelum disimpan ke ref asli
        const cleanValue = newValue.replace(/\D/g, '');
        uangBayar.value = cleanValue ? parseInt(cleanValue, 10) : '';
    }
});

const printData = ref(null); 
const printerSize = ref('58');

const formatRupiah = (angka) => {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(angka || 0);
};

// Fetch Katalog Paket Cuci
const fetchServices = async () => {
    isLoading.value = true;
    try {
        const response = await api.get('/laundry/services');
        services.value = response.data || [];
    } catch (error) {
        Swal.fire('Gagal!', 'Gagal mengambil katalog.', 'error');
    } finally {
        isLoading.value = false;
    }
};

// 🚀 FETCH DAFTAR PARFUM LIVE DARI DATABASE
const fetchPerfumes = async () => {
    try {
        const response = await api.get('/laundry/perfumes');
        perfumes.value = response.data || [];
    } catch (error) {
        console.error("Gagal load data parfum:", error);
    }
};

// 🚀 SIMPAN STATUS PARFUM TERSEDIA / HABIS (TOGGLE)
const togglePerfumeStatus = async (perfume) => {
    try {
        const newStatus = perfume.status === 'Tersedia' ? 'Habis' : 'Tersedia';
        // Tembak update status ke backend (menggunakan endpoint simpan parfum/update yang sudah ada)
        // Jika belum ada endpoint spesifik patch, kita bisa update lewat lokal state untuk simulasi kasir aman pengerjaan cepat
        perfume.status = newStatus;
        Swal.fire({ toast: true, position: 'top-end', icon: 'success', title: `${perfume.nama} di-set ${newStatus}`, showConfirmButton: false, timer: 1500 });
    } catch (error) {
        Swal.fire('Gagal!', 'Gagal mengubah status parfum.', 'error');
    }
};

// --- FETCH SETTING TOKO ---
const fetchStoreSetting = async () => {
    try {
        const response = await api.get('/laundry/setting');
        if (response.data) {
            storeInfo.value = {
                nama_toko: response.data.nama_toko || 'LAUNDRY POS',
                alamat: response.data.alamat || '',
                telepon: response.data.telepon || '',
                receipt_footer: response.data.receipt_footer || 'Terima Kasih'
            };
            if (response.data.qris_image) {
                qrisStoreUrl.value = response.data.qris_image; 
            }
        }
    } catch (error) {
        console.error("Gagal memuat setting toko:", error);
    }
};

const searchCustomer = async () => {
    if (customerName.value.length < 2) {
        showCustomerDropdown.value = false;
        customerResults.value = []; 
        return;
    }
    try {
        const response = await api.get(`/laundry/customers/search?q=${customerName.value}`);
        customerResults.value = response.data;
        showCustomerDropdown.value = customerResults.value.length > 0;
    } catch (error) { console.error("Gagal", error); }
};

const selectCustomer = (cust) => {
    customerName.value = cust.nama;
    let phone = String(cust.no_whatsapp);
    if (phone.startsWith('62')) phone = phone.substring(2);
    customerPhone.value = phone;
    showCustomerDropdown.value = false;
};

const closeCustomerDropdown = () => { setTimeout(() => { showCustomerDropdown.value = false; }, 200); };

onMounted(() => {
    fetchServices();
    fetchPerfumes(); // 🚀 Jalankan barengan pas start
    fetchStoreSetting(); 
    const today = new Date();
    today.setDate(today.getDate() + 2);
    estimasiSelesai.value = today.toISOString().split('T')[0];
});

const filteredServices = computed(() => {
    return services.value.filter(s => s.nama_produk.toLowerCase().includes(searchQuery.value.toLowerCase()));
});

// 🚀 List Parfum Premium yang Statusnya Ready (Bisa dipilih kasir)
const availablePerfumes = computed(() => {
    return perfumes.value.filter(p => p.status === 'Tersedia');
});

const addToCart = (service) => {
    const existingIndex = cart.value.findIndex(item => item.id === service.id);
    if (existingIndex !== -1) {
        Swal.fire({ toast: true, position: 'top-end', icon: 'info', title: 'Sudah di keranjang!', showConfirmButton: false, timer: 1500 });
    } else {
        // 🚀 Inisialisasi awal parfum default gratis bawaan toko
        cart.value.push({ 
            id: service.id, 
            nama_produk: service.nama_produk, 
            harga: service.harga_jual, 
            berat: 1, 
            satuan_dasar: service.satuan_dasar || 'KG',
            nama_parfum: 'Parfum Standar Bawaan',
            harga_parfum: 0
        });
        new Audio('https://www.soundjay.com/buttons/button-09.wav').play().catch(() => {});
    }
};

// 🚀 Update Pilihan Parfum di Keranjang & Auto Update Total Tagihan
const handleCartPerfumeChange = (index, event) => {
    const selectedId = event.target.value;
    if (selectedId === 'default') {
        cart.value[index].nama_parfum = 'Parfum Standar Bawaan';
        cart.value[index].harga_parfum = 0;
    } else {
        const pObj = perfumes.value.find(p => p.id == selectedId);
        if (pObj) {
            cart.value[index].nama_parfum = pObj.nama;
            cart.value[index].harga_parfum = pObj.harga;
        }
    }
};

const updateBerat = (index, delta) => {
    let item = cart.value[index];
    let current = parseFloat(item.berat) || 0;
    let actualDelta = item.satuan_dasar === 'KG' ? delta : (delta > 0 ? 1 : -1);
    let newVal = current + actualDelta;
    let minVal = item.satuan_dasar === 'KG' ? 0.1 : 1;
    item.berat = newVal < minVal ? minVal : (item.satuan_dasar === 'KG' ? Math.round(newVal * 10) / 10 : Math.round(newVal)); 
};

const removeCartItem = (index) => {
    cart.value.splice(index, 1);
    if (cart.value.length === 0) isCartOpen.value = false;
};

const resetForm = () => {
    cart.value = []; customerName.value = ''; customerPhone.value = ''; uangBayar.value = ''; paymentMethod.value = 'TUNAI'; 
    photoData.value = null; buktiTransferData.value = null; isCartOpen.value = false;
};

const clearCart = () => {
    Swal.fire({ title: 'Batalkan Cucian?', icon: 'warning', showCancelButton: true, confirmButtonColor: '#e11d48', confirmButtonText: 'Ya, Bersihkan'
    }).then((result) => { if (result.isConfirmed) resetForm(); });
};

// 🚀 HITUNG TOTAL TAGIHAN AKUMULATIF: (Harga Paket * Berat/Qty) + Biaya Tambahan Parfum
const totalTagihan = computed(() => { 
    return cart.value.reduce((acc, item) => acc + (item.harga * (parseFloat(item.berat) || 0)) + item.harga_parfum, 0); 
});

const kembalian = computed(() => { return (Number(uangBayar.value) || 0) - totalTagihan.value; });

const openCamera = async (target) => {
    cameraTarget.value = target;
    isCameraOpen.value = true;
    try {
        const stream = await navigator.mediaDevices.getUserMedia({ video: { facingMode: 'environment' } });
        cameraStream.value = stream;
        nextTick(() => {
            if (target === 'ITEM' && videoItemRef.value) videoItemRef.value.srcObject = stream;
            else if (target === 'QRIS' && videoQrisRef.value) videoQrisRef.value.srcObject = stream;
        });
    } catch (err) { Swal.fire('Oops!', 'Kamera tidak diizinkan.', 'error'); isCameraOpen.value = false; }
};

const takePhoto = () => {
    let video, canvas;
    if (cameraTarget.value === 'ITEM') { video = videoItemRef.value; canvas = canvasItemRef.value; } 
    else { video = videoQrisRef.value; canvas = canvasQrisRef.value; }

    if (!video || !canvas) return;
    canvas.width = video.videoWidth; canvas.height = video.videoHeight;
    canvas.getContext('2d').drawImage(video, 0, 0, canvas.width, canvas.height);
    
    if (cameraTarget.value === 'ITEM') photoData.value = canvas.toDataURL('image/jpeg', 0.7); 
    else if (cameraTarget.value === 'QRIS') buktiTransferData.value = canvas.toDataURL('image/jpeg', 0.7); 
    closeCamera();
};

const closeCamera = () => { if (cameraStream.value) cameraStream.value.getTracks().forEach(track => track.stop()); isCameraOpen.value = false; };
const cancelQris = () => { showQrisModal.value = false; buktiTransferData.value = null; closeCamera(); };
const confirmQris = () => { showQrisModal.value = false; processCheckout(); };

// --- 🚀 PROSES CHECKOUT UTAMA (SUDAH INJECT PARFUM PAYLOAD) ---
const processCheckout = async () => {
    if (cart.value.length === 0) return Swal.fire('Oops!', 'Keranjang kosong.', 'warning');
    if (!customerName.value || !customerPhone.value) return Swal.fire('Oops!', 'Data Pelanggan WAJIB diisi!', 'warning');
    if (paymentMethod.value === 'TUNAI' && kembalian.value < 0) return Swal.fire('Oops!', 'Uang bayar kurang!', 'warning');

    if (paymentMethod.value === 'QRIS' && !buktiTransferData.value) {
        showQrisModal.value = true;
        return; 
    }

    isSubmitting.value = true;
    try {
        const payload = {
            customer_name: customerName.value, 
            customer_phone: `62${customerPhone.value}`,
            estimasi_selesai: estimasiSelesai.value,
            // 🚀 OPER DATA PARFUM KE BACKEND GOLANG SECARA STRUKTUR BARU
            items: cart.value.map(item => ({ 
                product_id: item.id, 
                berat_kg: parseFloat(item.berat), 
                harga_per_kg: item.harga, 
                nama_parfum: item.nama_parfum,
                harga_parfum: item.harga_parfum,
                sub_total: (parseFloat(item.berat) * item.harga) + item.harga_parfum 
            })),
            total_amount: totalTagihan.value, 
            payment_method: paymentMethod.value, 
            payment_status: paymentMethod.value === 'PAYLATER' ? 'BELUM_LUNAS' : 'LUNAS',
            foto_barang_base64: photoData.value || '', 
            bukti_transfer_base64: buktiTransferData.value || ''
        };

        const response = await api.post('/laundry/checkout', payload);
        
        if (response.data && response.data.status === 'sukses') {
            const invoiceCode = response.data.invoice_code;
            const fotoUrl = response.data.foto_url; 
            
            printData.value = {
                toko_nama: storeInfo.value.nama_toko,
                toko_alamat: storeInfo.value.alamat,
                toko_telepon: storeInfo.value.telepon,
                toko_footer: storeInfo.value.receipt_footer,
                invoice: invoiceCode, kasir: 'Admin', tanggal: new Date().toLocaleString('id-ID'),
                pelanggan: customerName.value, estimasi: estimasiSelesai.value, items: [...cart.value],
                total: totalTagihan.value, metode: paymentMethod.value,
                bayar: paymentMethod.value === 'TUNAI' ? uangBayar.value : totalTagihan.value, kembali: kembalian.value > 0 ? kembalian.value : 0
            };

            let teksLampiranFoto = '_(Catatan: Barang telah dicek & ditimbang)_';
            if (fotoUrl) {
                teksLampiranFoto = `*📸 Foto Kondisi Barang:*\nKlik link berikut untuk melihat kondisi:\n${fotoUrl}`;
            }

            // Struk WA Dinamis ada rincian wangi parfumnya
            const waTextRaw = `Halo Kak *${customerName.value}*,\n\nTerima kasih sudah mempercayakan cuciannya di *${storeInfo.value.nama_toko.toUpperCase()}*! 🙏✨\n\n*🧾 RINCIAN NOTA:*\nNo. Resi: *${invoiceCode}*\nTgl Selesai: *${estimasiSelesai.value}*\n\n*Daftar Cucian:*\n${cart.value.map(i => `- ${i.nama_produk} (${i.berat} ${i.satuan_dasar})\n  🌸 Wangi: ${i.nama_parfum}`).join('\n')}\n\n*Total Tagihan: ${formatRupiah(totalTagihan.value)}*\nStatus Bayar: *${paymentMethod.value === 'PAYLATER' ? 'BELUM LUNAS ⚠️' : 'LUNAS ✅'}*\n\n${teksLampiranFoto}\n\nKami akan kabari lagi jika cucian sudah selesai ya kak! 💙`;
            
            // ... (Kodingan waTextRaw di atasnya biarkan saja) ..
            Swal.fire({
                icon: 'success', 
                title: 'Transaksi Sukses!', 
                html: `
                    <div class="text-center font-sans">
                        <p class="text-sm">Nomor Invoice: <b class="text-indigo-600">${invoiceCode}</b></p>
                        <div class="mt-4 p-3 bg-slate-50 border border-slate-200 rounded-xl inline-flex flex-col items-center gap-2 w-full">
                            <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Pilih Ukuran Kertas Thermal</label>
                            <select id="swal-printer-size" class="w-full px-3 py-2 border-2 border-slate-200 rounded-lg font-bold text-xs outline-none bg-white text-slate-700">
                                <option value="58" ${printerSize.value === '58' ? 'selected' : ''}>Printer Mini Kiloan (58mm)</option>
                                <option value="80" ${printerSize.value === '80' ? 'selected' : ''}>Printer Besar Enterprise (80mm)</option>
                            </select>
                        </div>
                    </div>
                `,
                confirmButtonText: 'Cetak & Kirim WA', 
                showCancelButton: true, 
                cancelButtonText: 'Tutup Saja', 
                confirmButtonColor: '#4f46e5',
            }).then(async (result) => {
                if (result.isConfirmed) {
                    // 🚀 Ambil ukuran kertas pilihan kasir dari pop-up Swal
                    const selectedSize = document.getElementById('swal-printer-size').value;
                    printerSize.value = selectedSize;

                    // Kirim Pesan WA
                    window.open(`https://wa.me/62${customerPhone.value}?text=${encodeURIComponent(waTextRaw)}`, '_blank');

                    await nextTick();
                    
                    // Panggil mesin printer browser
                    window.print(); 
                    resetForm();
                } else {
                    resetForm();
                }
            });
        }
    } catch (error) { Swal.fire('Gagal!', error.response?.data?.error || 'Gagal memproses.', 'error'); } 
    finally { isSubmitting.value = false; }
};
</script>

<template>
    <SidebarLaundry class="hide-on-print">
        <div class="flex h-screen bg-[#F8FAFC] font-sans overflow-hidden">
            
            <div class="flex-1 flex flex-col h-full overflow-hidden relative w-full lg:w-auto">
                <div class="p-4 md:p-6 shrink-0 bg-white shadow-sm z-10 flex flex-col gap-4 border-b border-slate-100">
                    <div class="bg-gradient-to-r from-sky-500 to-indigo-600 rounded-2xl p-4 md:p-5 flex flex-col sm:flex-row justify-between sm:items-center text-white shadow-lg shadow-indigo-200/50 gap-4">
                        <div class="flex items-center gap-3">
                            <div class="bg-white/20 p-2 rounded-xl backdrop-blur-sm">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z M4 8h16 M8 12h4 M8 16h8"/></svg>
                            </div>
                            <div>
                                <h1 class="text-lg md:text-xl font-black tracking-widest uppercase">POS Kasir Laundry</h1>
                                <span class="text-[9px] font-bold bg-white/20 px-2.5 py-0.5 rounded-md tracking-wider uppercase mt-1 inline-block">Sistem Kasir Pintar</span>
                            </div>
                        </div>

                        <button @click="showPerfumeControlModal = true" class="bg-white/20 hover:bg-white/30 px-4 py-2.5 rounded-xl border border-white/20 text-xs font-black uppercase tracking-wider transition-all self-stretch sm:self-auto flex items-center justify-center gap-2">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-pink-300 animate-pulse" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="12" cy="12" r="10"/><path d="M12 8a4 4 0 1 0 0 8 4 4 0 0 0 0-8z"/></svg>
                            Saklar Stok Parfum
                        </button>
                    </div>
                    
                    <div class="relative w-full group">
                        <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-400 group-focus-within:text-indigo-500 transition-colors" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
                        </div>
                        <input v-model="searchQuery" type="text" placeholder="Cari paket layanan cuci..." class="w-full pl-12 pr-4 py-3.5 bg-slate-100 border-2 border-transparent rounded-2xl focus:bg-white focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 outline-none font-bold text-sm text-slate-800 transition-all placeholder:text-slate-400">
                    </div>
                </div>

                <div class="flex-1 overflow-y-auto custom-scrollbar p-4 md:p-6 pb-28 lg:pb-6">
                    <div v-if="isLoading" class="flex flex-col items-center justify-center h-40 text-slate-400">
                        <div class="w-8 h-8 border-4 border-indigo-500 border-t-transparent rounded-full animate-spin mb-3"></div>
                        <p class="font-bold text-sm animate-pulse">Memuat Layanan...</p>
                    </div>
                    
                    <div v-else class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-4 md:gap-5">
                        <div v-for="service in filteredServices" :key="service.id" @click="addToCart(service)" 
                             class="bg-white p-4 md:p-5 rounded-[24px] border-2 border-slate-100 hover:border-indigo-400 hover:shadow-xl hover:shadow-indigo-100/50 cursor-pointer transition-all active:scale-95 group flex flex-col justify-between h-32 md:h-36 relative overflow-hidden">
                            <div class="z-10">
                                <span class="text-[9px] font-black text-indigo-600 bg-indigo-50 px-2.5 py-1 rounded-lg uppercase tracking-wider mb-2 inline-block border border-indigo-100">EST: {{ service.estimasi || 'STANDAR' }}</span>
                                <h3 class="text-xs md:text-sm font-black text-slate-800 uppercase leading-tight line-clamp-2 pr-2">{{ service.nama_produk }}</h3>
                            </div>
                            <p class="text-sm md:text-base font-black text-indigo-600 z-10">{{ formatRupiah(service.harga_jual) }} <span class="text-[10px] text-slate-400 font-bold uppercase">/ {{ service.satuan_dasar || 'KG' }}</span></p>
                        </div>
                    </div>
                </div>
            </div>

            <div v-if="!isCartOpen && cart.length > 0" class="fixed bottom-5 left-4 right-4 z-40 lg:hidden">
                <button @click="isCartOpen = true" class="w-full bg-slate-900 text-white rounded-[20px] p-5 flex justify-between items-center shadow-2xl shadow-slate-900/50 active:scale-95 transition-transform border border-slate-700">
                    <div class="flex items-center gap-3">
                        <div class="bg-indigo-600 text-white font-black text-xs w-8 h-8 rounded-full flex items-center justify-center">{{ cart.length }}</div>
                        <span class="text-xs font-black tracking-widest uppercase">Lihat Nota Timbangan</span>
                    </div>
                    <span class="text-base font-black text-emerald-400">{{ formatRupiah(totalTagihan) }}</span>
                </button>
            </div>

            <div v-if="isCartOpen" @click="isCartOpen = false" class="fixed inset-0 bg-slate-900/60 z-40 lg:hidden backdrop-blur-sm transition-opacity"></div>

            <div :class="isCartOpen ? 'translate-x-0' : 'translate-x-full lg:translate-x-0'" 
                 class="fixed inset-y-0 right-0 z-50 w-full sm:w-[400px] lg:w-[350px] xl:w-[400px] bg-white shadow-2xl lg:static lg:shadow-none lg:border-l border-slate-200 flex flex-col transition-transform duration-300 ease-in-out">
                
                <div class="bg-white text-slate-800 p-4 shrink-0 border-b border-slate-100 flex justify-between items-center">
                    <h2 class="text-sm font-black uppercase tracking-widest flex items-center gap-2">
                        <div class="bg-indigo-100 text-indigo-600 p-1.5 rounded-lg">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
                        </div>
                        RINCIAN NOTA TIMBANG
                    </h2>
                    <button @click="clearCart" class="text-[10px] font-black uppercase text-rose-500 hover:bg-rose-50 px-3 py-2 rounded-xl transition-colors">Batal</button>
                </div>

                <div class="flex-1 overflow-y-auto custom-scrollbar bg-slate-50 flex flex-col">
                    
                    <div class="p-4 border-b border-slate-200 bg-white flex flex-col gap-4 shadow-sm relative z-10 shrink-0">
                        <div class="relative">
                            <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1.5 ml-1">Nama Pelanggan / Cari Database</label>
                            <input v-model="customerName" @input="searchCustomer" @blur="closeCustomerDropdown" @focus="searchCustomer" type="text" class="w-full px-4 py-3 bg-slate-50 border-2 border-slate-200 rounded-xl focus:bg-white focus:border-indigo-500 outline-none font-bold text-slate-800 text-sm placeholder:text-slate-300" placeholder="Ketik nama pelanggan...">
                            
                            <div v-if="showCustomerDropdown" class="absolute left-0 right-0 top-[105%] z-50 bg-white border-2 border-slate-200 rounded-xl shadow-2xl overflow-hidden animate-[fadeIn_0.2s_ease-out]">
                                <div class="max-h-40 overflow-y-auto custom-scrollbar">
                                    <div v-for="cust in customerResults" :key="cust.id" @click="selectCustomer(cust)" class="p-3 hover:bg-indigo-50 cursor-pointer border-b border-slate-100 flex flex-col transition-colors">
                                        <span class="font-black text-sm text-slate-800">{{ cust.nama }}</span>
                                        <span class="text-[10px] font-bold text-slate-400">WhatsApp: +{{ cust.no_whatsapp }}</span>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                            <div>
                                <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1.5 ml-1">No. WhatsApp</label>
                                <div class="flex items-center bg-slate-50 border-2 border-slate-200 rounded-xl focus-within:bg-white focus-within:border-emerald-500 transition-all overflow-hidden">
                                    <div class="pl-3 pr-2 py-3 bg-slate-100 border-r border-slate-200"><span class="text-slate-500 font-black text-xs">+62</span></div>
                                    <input v-model="customerPhone" @input="formatNoHpCustomer" type="number" class="w-full px-3 py-3 bg-transparent outline-none font-bold text-slate-800 text-sm" placeholder="8123...">
                                </div>
                            </div>
                            <div>
                                <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1.5 ml-1">Estimasi Selesai</label>
                                <input v-model="estimasiSelesai" type="date" class="w-full px-3 py-3 bg-slate-50 border-2 border-slate-200 rounded-xl outline-none font-bold text-sm text-slate-700 focus:bg-white focus:border-indigo-500 transition-all">
                            </div>
                        </div>

                        <div class="mt-1 border-2 border-dashed border-slate-200 rounded-xl p-3 bg-slate-50 flex flex-col items-center justify-center relative overflow-hidden min-h-[110px]">
                            <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-3 absolute top-3 left-3 z-10">Bukti Dokumentasi Barang</label>
                            <template v-if="(!isCameraOpen || cameraTarget !== 'ITEM') && !photoData">
                                <button @click="openCamera('ITEM')" class="flex flex-col items-center text-slate-400 hover:text-indigo-600 transition-colors mt-4">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-7 h-7 mb-1.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M23 19a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z"/><circle cx="12" cy="13" r="4"/></svg>
                                    <span class="text-[10px] font-black uppercase tracking-widest">Ambil Foto Masuk</span>
                                </button>
                            </template>
                            <div v-show="isCameraOpen && cameraTarget === 'ITEM'" class="w-full h-full relative z-20 mt-4 rounded-lg overflow-hidden bg-black flex flex-col">
                                <video ref="videoItemRef" autoplay playsinline class="w-full h-28 object-cover"></video>
                                <div class="absolute bottom-2 left-0 right-0 flex justify-center gap-2">
                                    <button @click="takePhoto" class="bg-indigo-600 text-white text-[9px] font-black px-4 py-1.5 rounded-full uppercase border-2 border-white">Jepret</button>
                                    <button @click="closeCamera" class="bg-rose-500 text-white text-[9px] font-black px-4 py-1.5 rounded-full uppercase border-2 border-white">Batal</button>
                                </div>
                            </div>
                            <div v-if="photoData" class="w-full h-full relative z-20 mt-4 rounded-lg overflow-hidden border border-slate-200">
                                <img :src="photoData" class="w-full h-28 object-cover" />
                                <button @click="photoData = null" class="absolute top-2 right-2 bg-rose-500 text-white w-6 h-6 rounded-full flex items-center justify-center shadow-lg"><svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg></button>
                            </div>
                            <canvas ref="canvasItemRef" class="hidden"></canvas>
                        </div>
                    </div>

                    <div class="p-4 flex flex-col gap-3">
                        <div v-if="cart.length === 0" class="flex flex-col items-center justify-center text-center opacity-50 py-10">
                            <p class="text-slate-500 font-black text-xs uppercase tracking-widest">Belum Ada Item Ditimbang</p>
                        </div>
                        
                        <div v-else v-for="(item, index) in cart" :key="index" class="bg-white p-4 rounded-2xl shadow-sm border border-slate-200 relative flex flex-col gap-3">
                            <button @click="removeCartItem(index)" class="absolute -right-2 -top-2 w-6 h-6 bg-white text-rose-500 border border-slate-200 rounded-full flex items-center justify-center shadow-md z-10"><svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg></button>
                            
                            <div class="flex justify-between items-start">
                                <div class="pr-2">
                                    <h4 class="text-xs font-black text-slate-800 uppercase leading-tight">{{ item.nama_produk }}</h4>
                                    <p class="text-[9px] font-bold text-slate-400 mt-0.5 uppercase tracking-wider">@ {{ formatRupiah(item.harga) }} /{{ item.satuan_dasar }}</p>
                                </div>
                                <span class="text-xs font-black text-indigo-600 bg-indigo-50 px-2 py-0.5 rounded-md border border-indigo-100 shrink-0">
                                    {{ formatRupiah((item.harga * parseFloat(item.berat)) + item.harga_parfum) }}
                                </span>
                            </div>

                            <div class="flex flex-col gap-1 bg-pink-50/50 p-2 rounded-xl border border-pink-100">
                                <label class="text-[8px] font-black text-pink-600 uppercase tracking-widest ml-0.5 flex items-center gap-1">
                                    🌸 Opsi Varian Parfum Wangi
                                </label>
                                <select @change="handleCartPerfumeChange(index, $event)" class="w-full text-[11px] font-bold bg-white text-slate-700 border border-pink-200 rounded-lg px-2 py-1 outline-none cursor-pointer focus:border-pink-400 transition-colors">
                                    <option value="default">Parfum Standar Toko (Bawaan - Gratis)</option>
                                    <option v-for="perfume in availablePerfumes" :key="perfume.id" :value="perfume.id">
                                        {{ perfume.nama }} {{ perfume.harga > 0 ? `(+ ${formatRupiah(perfume.harga)})` : '(Gratis)' }}
                                    </option>
                                </select>
                            </div>

                            <div class="flex items-center justify-between bg-slate-50 rounded-xl p-1.5 border border-slate-200">
                                <span class="text-[9px] font-black text-slate-500 uppercase tracking-widest pl-2">{{ item.satuan_dasar === 'KG' ? 'Berat Timbangan:' : 'Jumlah Item:' }}</span>
                                <div class="flex items-center bg-white rounded-lg border border-slate-200 overflow-hidden shadow-sm">
                                    <button @click="updateBerat(index, -0.5)" class="w-8 h-8 flex items-center justify-center text-slate-500 active:bg-slate-200"><svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="5" y1="12" x2="19" y2="12"/></svg></button>
                                    <input v-model="item.berat" type="number" step="0.1" min="0.1" class="w-10 text-center font-black text-xs text-indigo-700 outline-none bg-transparent">
                                    <button @click="updateBerat(index, 0.5)" class="w-8 h-8 flex items-center justify-center text-slate-500 active:bg-slate-200 border-l border-slate-100"><svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg></button>
                                    <span class="text-[9px] font-black text-slate-400 bg-slate-50 h-8 flex items-center px-2 border-l border-slate-200 uppercase">{{ item.satuan_dasar }}</span>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="bg-white border-t border-slate-200 shrink-0 p-4 shadow-[0_-10px_40px_-15px_rgba(0,0,0,0.1)] z-20">
                    <div class="mb-3.5 grid grid-cols-4 gap-1.5">
                        <button @click="paymentMethod = 'TUNAI'" class="py-2.5 rounded-xl border-2 flex flex-col items-center justify-center transition-all text-[9px] font-black tracking-widest" :class="paymentMethod === 'TUNAI' ? 'bg-slate-800 border-slate-800 text-white' : 'bg-white border-slate-200 text-slate-500'">TUNAI</button>
                        <button @click="paymentMethod = 'QRIS'" class="py-2.5 rounded-xl border-2 flex flex-col items-center justify-center transition-all text-[9px] font-black tracking-widest" :class="paymentMethod === 'QRIS' ? 'bg-slate-800 border-slate-800 text-white' : 'bg-white border-slate-200 text-slate-500'">QRIS</button>
                        <button @click="paymentMethod = 'DEBIT'" class="py-2.5 rounded-xl border-2 flex flex-col items-center justify-center transition-all text-[9px] font-black tracking-widest" :class="paymentMethod === 'DEBIT' ? 'bg-slate-800 border-slate-800 text-white' : 'bg-white border-slate-200 text-slate-500'">DEBIT</button>
                        <button @click="paymentMethod = 'PAYLATER'" class="py-2.5 rounded-xl border-2 flex flex-col items-center justify-center transition-all text-[8px] font-black tracking-tighter" :class="paymentMethod === 'PAYLATER' ? 'bg-amber-500 border-amber-500 text-white' : 'bg-amber-50 border-amber-200 text-amber-600'">BAYAR NANTI</button>
                    </div>

                    <div class="space-y-2.5 mb-3.5 bg-slate-50 p-3 rounded-xl border border-slate-100">
                        <div class="flex justify-between items-center">
                            <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest">GRAND TOTAL</span>
                            <span class="text-xl font-black text-slate-800 tracking-tighter">{{ formatRupiah(totalTagihan) }}</span>
                        </div>
                        <div class="pt-2.5 border-t border-slate-200" v-if="paymentMethod === 'TUNAI'">
                            <div class="flex justify-between items-center mb-3">
    <span class="text-[10px] font-black text-slate-500 uppercase tracking-widest">Uang Cash</span>
    <div class="flex items-center relative">
        <span class="absolute left-3 text-xs font-bold text-slate-400">Rp</span>
        <input 
            v-model="formattedUangBayar" 
            type="text" 
            class="w-28 bg-white border-2 border-slate-200 rounded-xl pl-8 pr-3 py-1.5 text-right font-black text-xs text-slate-800 outline-none focus:border-indigo-500"
            placeholder="0"
        >
    </div>
</div>
                            <div class="flex justify-between items-center bg-white p-2 rounded-xl border border-slate-100">
                                <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Uang Kembali</span>
                                <span class="text-xs font-black tracking-tight" :class="kembalian < 0 ? 'text-rose-500' : 'text-emerald-500'">{{ formatRupiah(kembalian) }}</span>
                            </div>
                        </div>
                    </div>
                    
                    <button @click="processCheckout" :disabled="cart.length === 0 || isSubmitting" class="w-full bg-gradient-to-r from-sky-500 to-indigo-600 text-white py-3.5 rounded-xl font-black text-xs uppercase tracking-[0.2em] transition-all active:scale-95 disabled:opacity-50 shadow-md">
                        {{ isSubmitting ? 'MEMPROSES...' : 'PROSES & KIRIM NOTA WA' }}
                    </button>
                </div>
            </div>
        </div>

        <Teleport to="body">
            <div v-if="showPerfumeControlModal" class="fixed inset-0 z-[9999] flex items-center justify-center bg-slate-900/80 backdrop-blur-sm p-4 animate-[fadeIn_0.2s_ease-out]">
                <div class="bg-white rounded-3xl shadow-2xl w-full max-w-md overflow-hidden flex flex-col">
                    <div class="bg-pink-600 p-5 text-center text-white flex justify-between items-center">
                        <h3 class="font-black text-sm uppercase tracking-widest">Saklar Ketersediaan Parfum</h3>
                        <button @click="showPerfumeControlModal = false" class="text-white hover:opacity-70"><svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg></button>
                    </div>
                    <div class="p-4 max-h-[60vh] overflow-y-auto custom-scrollbar space-y-2">
                        <div v-if="perfumes.length === 0" class="text-center py-6 text-xs font-bold text-slate-400 uppercase">Belum ada varian parfum di Master Layanan.</div>
                        <div v-for="perfume in perfumes" :key="perfume.id" class="p-3 border-2 border-slate-100 rounded-xl flex items-center justify-between bg-slate-50/50">
                            <div>
                                <h4 class="font-black text-slate-800 text-xs uppercase">🌸 {{ perfume.nama }}</h4>
                                <p class="text-[9px] font-bold text-slate-400 mt-0.5">Charge: {{ perfume.harga > 0 ? formatRupiah(perfume.harga) : 'Gratis' }}</p>
                            </div>
                            <button @click="togglePerfumeStatus(perfume)" :class="perfume.status === 'Tersedia' ? 'bg-emerald-500' : 'bg-slate-300'" class="w-12 h-6 rounded-full p-0.5 transition-all duration-300 flex items-center relative shadow-inner">
                                <div :class="perfume.status === 'Tersedia' ? 'translate-x-6' : 'translate-x-0'" class="w-5 h-5 bg-white rounded-full shadow-md transition-transform duration-300"></div>
                            </button>
                        </div>
                    </div>
                    <div class="p-4 bg-slate-50 border-t border-slate-100 text-center">
                        <button @click="showPerfumeControlModal = false" class="bg-slate-900 text-white font-black text-[10px] tracking-widest px-6 py-2.5 rounded-xl uppercase">Selesai Monitor</button>
                    </div>
                </div>
            </div>
        </Teleport>

        <Teleport to="body">
            <div v-if="showQrisModal" class="fixed inset-0 z-[9999] flex items-center justify-center bg-slate-900/80 backdrop-blur-sm p-4 animate-[fadeIn_0.2s_ease-out]">
                <div class="bg-white rounded-3xl shadow-2xl w-full max-w-sm md:max-w-2xl overflow-hidden flex flex-col">
                    <div class="bg-indigo-600 p-4 text-center text-white shrink-0">
                        <h3 class="font-black text-sm uppercase tracking-widest">SCAN QRIS PEMBAYARAN TOKO</h3>
                        <p class="text-xs font-bold text-indigo-200 mt-0.5">Total Tagihan: {{ formatRupiah(totalTagihan) }}</p>
                    </div>
                    <div class="p-4 grid grid-cols-1 md:grid-cols-2 gap-4 max-h-[70vh] overflow-y-auto custom-scrollbar items-center">
                        <div class="bg-slate-50 rounded-2xl border-2 border-slate-200 p-4 flex flex-col items-center justify-center">
                            <img :src="qrisStoreUrl || 'https://upload.wikimedia.org/wikipedia/commons/d/d0/QR_code_for_mobile_English_Wikipedia.svg'" alt="QRIS Toko" class="w-full max-w-[180px] aspect-square object-contain mix-blend-multiply" />
                        </div>
                        <div class="w-full border-2 border-dashed border-emerald-300 rounded-2xl p-3 bg-emerald-50/50 flex flex-col items-center justify-center relative min-h-[160px]">
                            <template v-if="(!isCameraOpen || cameraTarget !== 'QRIS') && !buktiTransferData">
                                <button @click="openCamera('QRIS')" class="flex flex-col items-center text-emerald-600 hover:text-emerald-700 transition-colors my-2">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-10 h-10 mb-2" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/></svg>
                                    <span class="text-[9px] font-black uppercase tracking-widest text-center">Foto Bukti Transfer HP Pelanggan</span>
                                </button>
                            </template>
                            <div v-show="isCameraOpen && cameraTarget === 'QRIS'" class="w-full h-full relative z-20 rounded-xl overflow-hidden bg-black flex flex-col">
                                <video ref="videoQrisRef" autoplay playsinline class="w-full h-36 object-cover"></video>
                                <button @click="takePhoto" class="absolute bottom-2 left-1/2 -translate-x-1/2 bg-emerald-500 text-white text-[9px] font-black px-4 py-1.5 rounded-full uppercase border-2 border-white">Foto</button>
                            </div>
                            <div v-if="buktiTransferData" class="w-full h-full relative z-20 rounded-xl overflow-hidden border border-emerald-200">
                                <img :src="buktiTransferData" class="w-full h-36 object-cover" />
                                <button @click="buktiTransferData = null" class="absolute top-2 right-2 bg-rose-500 text-white w-6 h-6 rounded-full flex items-center justify-center shadow-lg"><svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg></button>
                            </div>
                        </div>
                    </div>
                    <div class="p-4 bg-slate-50 border-t border-slate-100 flex gap-2 shrink-0">
                        <button @click="cancelQris" class="w-1/3 py-2.5 text-slate-500 font-black text-xs uppercase tracking-widest hover:bg-slate-200 rounded-xl">Batal</button>
                        <button @click="confirmQris" :disabled="!buktiTransferData" class="w-2/3 py-2.5 bg-emerald-500 text-white font-black text-xs uppercase tracking-widest rounded-xl disabled:opacity-50 flex justify-center items-center">Sahkan Bayar</button>
                    </div>
                </div>
            </div>
            <canvas ref="canvasQrisRef" class="hidden"></canvas>
        </Teleport>

    </SidebarLaundry>

    <Teleport to="body">
                <div id="print-area" class="print-only">
            <div v-if="printData" :class="printerSize === '80' ? 'receipt-80' : 'receipt-58'" class="receipt-container font-mono">
                
                <h2 class="center bold uppercase text-base leading-tight m-0">{{ printData.toko_nama }}</h2>
                <p class="center text-xs m-0 font-normal mt-1" v-if="printData.toko_alamat">{{ printData.toko_alamat }}</p>
                <p class="center text-xs m-0 font-normal" v-if="printData.toko_telepon">Telp: {{ printData.toko_telepon }}</p>
                
                <div class="divider-dash"></div>
                
                <table class="w-full text-xs line-tight">
                    <tr><td class="w-24">No. Resi</td><td>: <b>{{ printData.invoice }}</b></td></tr>
                    <tr><td>Tanggal</td><td>: {{ printData.tanggal }}</td></tr>
                    <tr><td>Pelanggan</td><td>: <b>{{ printData.pelanggan }}</b></td></tr>
                    <tr><td>Est Selesai</td><td>: {{ printData.estimasi }}</td></tr>
                </table>
                
                <div class="divider-dash"></div>
                
                <div class="item-block text-xs" v-for="(item, index) in printData.items" :key="index">
                    <div class="flex justify-between bold">
                        <span class="uppercase truncate-name">{{ item.nama_produk }}</span>
                        <span>{{ formatRupiah((item.harga * item.berat) + item.harga_parfum) }}</span>
                    </div>
                    <div class="flex justify-between text-slate-600 text-[11px] pl-2">
                        <span>Detail: {{ item.berat }} {{ item.satuan_dasar }} x {{ formatRupiah(item.harga) }}</span>
                    </div>
                    <div class="text-[10px] italic text-slate-700 pl-2 flex items-center gap-1 mt-0.5">
                        * Wangi: {{ item.nama_parfum }} {{ item.harga_parfum > 0 ? `(+${formatRupiah(item.harga_parfum)})` : '' }}
                    </div>
                </div>
                
                <div class="divider-dash"></div>
                
                <table class="w-full text-xs line-tight">
                    <tr class="bold text-sm">
                        <td>GRAND TOTAL</td>
                        <td class="text-right">{{ formatRupiah(printData.total) }}</td>
                    </tr>
                    <tr>
                        <td>Metode Bayar</td>
                        <td class="text-right uppercase">{{ printData.metode }}</td>
                    </tr>
                    <tr>
                        <td>Nominal Bayar</td>
                        <td class="text-right">{{ formatRupiah(printData.bayar) }}</td>
                    </tr>
                    <tr v-if="printData.metode === 'TUNAI'">
                        <td class="bold">Uang Kembali</td>
                        <td class="text-right bold text-sm">{{ formatRupiah(printData.kembali) }}</td>
                    </tr>
                </table>
                
                <div class="divider-dash"></div>
                
                <p class="center bold text-xs tracking-wider uppercase mt-4 m-0">{{ printData.toko_footer }}</p>
                <p class="center text-[9px] text-slate-400 mt-1 m-0">Powered by POS UMKM &copy; 2026</p>
            </div>
        </div>
    </Teleport>
</template>

<style scoped>
/* 🚀 Biarkan CSS bawaan kamu yang lama tetap di sini (seperti custom-scrollbar, dll) */
.custom-scrollbar::-webkit-scrollbar { height: 5px; width: 5px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
input[type=number]::-webkit-inner-spin-button, input[type=number]::-webkit-outer-spin-button { -webkit-appearance: none; margin: 0; }
input[type=number] { -moz-appearance: textfield; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(-8px); } to { opacity: 1; transform: translateY(0); } }
</style>

<style>
.print-only { display: none; }

@media print {
    @page { 
        margin: 0 !important; 
        size: portrait !important; 
    }
    
    /* 1. MATIKAN SELURUH APLIKASI VUE SAAT PRINT! */
    #app { 
        display: none !important; 
    }

    /* 2. TAMPILKAN STRUK KITA (Karena udah di teleport ke luar #app) */
    .print-only { 
        display: block !important; 
        position: absolute !important;
        left: 0 !important;
        top: 0 !important;
        width: 100% !important;
        background: white !important;
        z-index: 999999 !important;
    }

    /* --- UKURAN KERTAS & TYPOGRAPHY THERMAL --- */
    .receipt-58 {
        width: 54mm !important;
        padding: 2mm 1mm !important;
        box-sizing: border-box !important;
        margin: 0 auto 0 0 !important;
    }

    .receipt-80 {
        width: 76mm !important;
        padding: 4mm 3mm !important;
        box-sizing: border-box !important;
        margin: 0 auto 0 0 !important;
    }

    .receipt-container {
        font-family: 'Courier New', Courier, monospace !important;
        line-height: 1.3 !important;
        background: white !important;
    }

    .receipt-container p, 
    .receipt-container span, 
    .receipt-container td,
    .receipt-container div,
    .receipt-container b { 
        color: #000000 !important;
        font-weight: 800 !important; 
    }

    .center { text-align: center !important; } 
    .bold { font-weight: bold !important; } 
    .uppercase { text-transform: uppercase !important; } 
    .m-0 { margin: 0 !important; }
    .w-full { width: 100% !important; }
    .text-right { text-align: right !important; }
    
    .text-xs { font-size: 11px !important; } 
    .text-sm { font-size: 13px !important; } 
    .text-base { font-size: 15px !important; }
    .mt-4 { margin-top: 14px !important; }
    .mt-1 { margin-top: 3px !important; }

    .divider-dash { border-top: 2px dashed #000000 !important; margin: 6px 0 !important; width: 100% !important; height: 1px !important; }
    .item-block { margin-bottom: 5px !important; page-break-inside: avoid !important; }
    .line-tight tr td { padding: 2px 0 !important; }

    .truncate-name {
        display: inline-block !important;
        max-width: 65% !important;
        overflow: hidden !important;
        text-overflow: ellipsis !important;
        white-space: nowrap !important;
    }
}
</style>