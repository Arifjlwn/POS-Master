<script setup>
import { ref, onMounted, computed, nextTick } from 'vue';
import SidebarLaundry from './SidebarLaundry.vue';
import api from '../../../api.js'; 
import Swal from 'sweetalert2';

// --- STATE DATA ---
const services = ref([]);
const cart = ref([]);
const isLoading = ref(false);
const isSubmitting = ref(false);

const searchQuery = ref('');
const customerName = ref('');
const customerPhone = ref('')
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
const printData = ref(null); 

const formatRupiah = (angka) => {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(angka || 0);
};

const fetchServices = async () => {
    isLoading.value = true;
    try {
        const response = await api.get('/laundry/services');
        services.value = response.data;
    } catch (error) {
        Swal.fire('Gagal!', 'Gagal mengambil katalog.', 'error');
    } finally {
        isLoading.value = false;
    }
};

// --- 🚀 FETCH SETTING TOKO ---
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
    fetchStoreSetting(); 
    const today = new Date();
    today.setDate(today.getDate() + 2);
    estimasiSelesai.value = today.toISOString().split('T')[0];
});

const filteredServices = computed(() => {
    return services.value.filter(s => s.nama_produk.toLowerCase().includes(searchQuery.value.toLowerCase()));
});

const addToCart = (service) => {
    const existingIndex = cart.value.findIndex(item => item.id === service.id);
    if (existingIndex !== -1) {
        Swal.fire({ toast: true, position: 'top-end', icon: 'info', title: 'Sudah di keranjang!', showConfirmButton: false, timer: 1500 });
    } else {
        cart.value.push({ id: service.id, nama_produk: service.nama_produk, harga: service.harga_jual, berat: 1, satuan_dasar: service.satuan_dasar || 'KG' });
        new Audio('https://www.soundjay.com/buttons/button-09.wav').play().catch(() => {});
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

const totalTagihan = computed(() => { return cart.value.reduce((acc, item) => acc + (item.harga * (parseFloat(item.berat) || 0)), 0); });
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

// --- 🚀 PROSES CHECKOUT UTAMA ---
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
            customer_name: customerName.value, customer_phone: `62${customerPhone.value}`,
            estimasi_selesai: estimasiSelesai.value,
            items: cart.value.map(item => ({ product_id: item.id, berat_kg: parseFloat(item.berat), harga_per_kg: item.harga, sub_total: parseFloat(item.berat) * item.harga })),
            total_amount: totalTagihan.value, payment_method: paymentMethod.value, payment_status: paymentMethod.value === 'NANTI_AJA' ? 'BELUM_LUNAS' : 'LUNAS',
            foto_barang_base64: photoData.value || '', bukti_transfer_base64: buktiTransferData.value || ''
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

            const waTextRaw = `Halo Kak *${customerName.value}*,\n\nTerima kasih sudah mempercayakan cuciannya di *${storeInfo.value.nama_toko.toUpperCase()}*! 🙏✨\n\n*🧾 RINCIAN NOTA:*\nNo. Resi: *${invoiceCode}*\nTgl Selesai: *${estimasiSelesai.value}*\n\n*Daftar Cucian:*\n${cart.value.map(i => `- ${i.nama_produk} (${i.berat} ${i.satuan_dasar})`).join('\n')}\n\n*Total Tagihan: ${formatRupiah(totalTagihan.value)}*\nStatus Bayar: *${paymentMethod.value === 'NANTI_AJA' ? 'BELUM LUNAS ⚠️' : 'LUNAS ✅'}*\n\n${teksLampiranFoto}\n\nKami akan kabari lagi jika cucian sudah selesai ya kak! 💙`;
            
            Swal.fire({
                icon: 'success', title: 'Transaksi Sukses!', 
                html: `Nomor Invoice: <b>${invoiceCode}</b><br><br>Klik tombol di bawah ini untuk mencetak struk & kirim WA.`,
                confirmButtonText: 'Cetak & Kirim WA', showCancelButton: true, cancelButtonText: 'Tutup Saja', confirmButtonColor: '#4f46e5',
            }).then((result) => {
                if (result.isConfirmed) {
                    window.open(`https://wa.me/62${customerPhone.value}?text=${encodeURIComponent(waTextRaw)}`, '_blank');
                    setTimeout(() => { window.print(); resetForm(); }, 500);
                } else resetForm();
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
                    <div class="bg-gradient-to-r from-sky-500 to-indigo-600 rounded-2xl p-5 md:p-6 flex justify-between items-center text-white shadow-lg shadow-indigo-200/50">
                        <div>
                            <div class="flex items-center gap-3">
                                <div class="bg-white/20 p-2 rounded-xl backdrop-blur-sm"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 md:w-7 md:h-7 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M12 22v-5"/><path d="M9 7V2"/><path d="M15 7V2"/><rect width="16" height="15" x="4" y="7" rx="2" ry="2"/></svg></div>
                                <h1 class="text-xl md:text-2xl font-black tracking-widest uppercase">POS Kasir</h1>
                            </div>
                            <span class="text-[10px] md:text-xs font-bold bg-white/20 px-3 py-1 rounded-full mt-3 inline-block tracking-widest">LAUNDRY KILOAN & SATUAN</span>
                        </div>
                    </div>
                    <div class="relative w-full group">
                        <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none"><svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-400 group-focus-within:text-indigo-500 transition-colors" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg></div>
                        <input v-model="searchQuery" type="text" placeholder="Cari paket layanan cuci..." class="w-full pl-12 pr-4 py-4 bg-slate-100 border-2 border-transparent rounded-2xl focus:bg-white focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 outline-none font-bold text-sm text-slate-800 transition-all placeholder:text-slate-400">
                    </div>
                </div>

                <div class="flex-1 overflow-y-auto custom-scrollbar p-4 md:p-6 pb-28 lg:pb-6">
                    <div v-if="isLoading" class="flex flex-col items-center justify-center h-40 text-slate-400">
                        <div class="w-8 h-8 border-4 border-indigo-500 border-t-transparent rounded-full animate-spin mb-3"></div><p class="font-bold text-sm animate-pulse">Memuat Layanan...</p>
                    </div>
                    <div v-else class="grid grid-cols-2 sm:grid-cols-3 xl:grid-cols-4 gap-4 md:gap-5">
                        <div v-for="service in filteredServices" :key="service.id" @click="addToCart(service)" 
                             class="bg-white p-4 md:p-5 rounded-[24px] border-2 border-slate-100 hover:border-indigo-400 hover:shadow-xl hover:shadow-indigo-100 cursor-pointer transition-all active:scale-95 group flex flex-col justify-between h-36 md:h-40 relative overflow-hidden">
                            <div class="absolute -right-4 -bottom-4 opacity-5 group-hover:opacity-10 group-hover:scale-110 transition-all duration-500"><svg xmlns="http://www.w3.org/2000/svg" class="w-28 h-28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22v-5"/><path d="M9 7V2"/><path d="M15 7V2"/><rect width="16" height="15" x="4" y="7" rx="2" ry="2"/></svg></div>
                            <div class="z-10">
                                <span class="text-[9px] md:text-[10px] font-black text-indigo-600 bg-indigo-50 px-2.5 py-1 rounded-lg uppercase tracking-wider mb-2 md:mb-3 inline-block border border-indigo-100">EST: {{ service.estimasi || 'STANDAR' }}</span>
                                <h3 class="text-xs md:text-sm font-black text-slate-800 uppercase leading-tight line-clamp-2 pr-2">{{ service.nama_produk }}</h3>
                            </div>
                            <p class="text-sm md:text-base font-black text-indigo-600 z-10">{{ formatRupiah(service.harga_jual) }} <span class="text-[10px] text-slate-400 font-bold uppercase">/{{ service.satuan_dasar || 'KG' }}</span></p>
                        </div>
                    </div>
                </div>
            </div>

            <div v-if="!isCartOpen && cart.length > 0" class="fixed bottom-5 left-4 right-4 z-40 lg:hidden">
                <button @click="isCartOpen = true" class="w-full bg-slate-900 text-white rounded-[20px] p-5 flex justify-between items-center shadow-2xl shadow-slate-900/50 active:scale-95 transition-transform border border-slate-700">
                    <div class="flex items-center gap-3">
                        <div class="bg-white/20 text-white font-black text-xs w-8 h-8 rounded-full flex items-center justify-center">{{ cart.length }}</div>
                        <span class="text-xs font-black tracking-widest uppercase">Lihat Timbangan</span>
                    </div>
                    <span class="text-base font-black text-emerald-400">{{ formatRupiah(totalTagihan) }}</span>
                </button>
            </div>

            <div v-if="isCartOpen" @click="isCartOpen = false" class="fixed inset-0 bg-slate-900/60 z-40 lg:hidden backdrop-blur-sm transition-opacity"></div>

            <div :class="isCartOpen ? 'translate-x-0' : 'translate-x-full lg:translate-x-0'" 
                 class="fixed inset-y-0 right-0 z-50 w-full sm:w-[400px] lg:w-[340px] xl:w-[400px] bg-white shadow-2xl lg:static lg:shadow-none lg:border-l border-slate-200 flex flex-col transition-transform duration-300 ease-in-out">
                
                <div class="bg-white text-slate-800 p-4 md:p-5 flex justify-between items-center shrink-0 border-b border-slate-100">
                    <h2 class="text-sm md:text-base font-black uppercase tracking-widest flex items-center gap-2">
                        <div class="bg-indigo-100 text-indigo-600 p-1.5 rounded-lg"><svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/><polyline points="10 9 9 9 8 9"/></svg></div>
                        NOTA CUCIAN
                    </h2>
                    <div class="flex items-center gap-3">
                        <button @click="clearCart" class="text-[10px] font-black uppercase text-rose-500 hover:bg-rose-50 px-3 py-2 rounded-xl transition-colors">Batal</button>
                    </div>
                </div>

                <div class="flex-1 overflow-y-auto custom-scrollbar bg-slate-50 flex flex-col">
                    <div class="p-4 md:p-5 border-b border-slate-200 bg-white flex flex-col gap-4 shadow-sm relative z-10 shrink-0">
                        <div class="relative">
                            <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1.5 ml-1">Nama Pelanggan / Cari Lama</label>
                            <input v-model="customerName" @input="searchCustomer" @blur="closeCustomerDropdown" @focus="searchCustomer" type="text" class="w-full px-4 py-3 bg-slate-50 border-2 border-slate-200 rounded-xl focus:bg-white focus:border-indigo-500 outline-none font-bold text-slate-800 text-sm placeholder:text-slate-300" placeholder="Ketik nama atau nomor...">
                            <div v-if="showCustomerDropdown" class="absolute left-0 right-0 top-[105%] z-50 bg-white border-2 border-slate-200 rounded-xl shadow-2xl overflow-hidden animate-[fadeIn_0.2s_ease-out]">
                                <div class="px-3 py-2 bg-slate-50 border-b border-slate-100 flex items-center gap-2">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 text-indigo-500" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
                                    <span class="text-[9px] font-black text-slate-400 uppercase tracking-widest">Daftar Pelanggan Langganan</span>
                                </div>
                                <div class="max-h-40 overflow-y-auto custom-scrollbar">
                                    <div v-for="cust in customerResults" :key="cust.id" @click="selectCustomer(cust)" class="p-3 hover:bg-indigo-50 cursor-pointer border-b border-slate-100 flex flex-col transition-colors group">
                                        <span class="font-black text-sm text-slate-800 group-hover:text-indigo-700">{{ cust.nama }}</span>
                                        <span class="text-[10px] font-bold text-slate-400 mt-0.5">Wa: +{{ cust.no_whatsapp }}</span>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div>
                                <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1.5 ml-1">No. WhatsApp</label>
                                <div class="flex items-center bg-slate-50 border-2 border-slate-200 rounded-xl focus-within:bg-white focus-within:border-emerald-500 transition-all overflow-hidden">
                                    <div class="pl-3 pr-2 py-3 bg-slate-100 border-r border-slate-200"><span class="text-slate-500 font-black text-xs">+62</span></div>
                                    <input v-model="customerPhone" @input="formatNoHpCustomer" type="number" class="w-full px-3 py-3 bg-transparent outline-none font-bold text-slate-800 text-sm" placeholder="8123...">
                                </div>
                            </div>
                            <div>
                                <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1.5 ml-1">Tgl Selesai</label>
                                <input v-model="estimasiSelesai" type="date" class="w-full px-3 py-3 bg-slate-50 border-2 border-slate-200 rounded-xl outline-none font-bold text-sm text-slate-700 focus:bg-white focus:border-indigo-500 transition-all">
                            </div>
                        </div>

                        <div class="mt-2 border-2 border-dashed border-slate-200 rounded-xl p-3 bg-slate-50 flex flex-col items-center justify-center relative overflow-hidden min-h-[120px]">
                            <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-3 absolute top-3 left-3 z-10">Bukti Barang / Baju</label>
                            
                            <template v-if="(!isCameraOpen || cameraTarget !== 'ITEM') && !photoData">
                                <button @click="openCamera('ITEM')" class="flex flex-col items-center text-slate-400 hover:text-indigo-600 transition-colors mt-4">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8 mb-2" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M23 19a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z"/><circle cx="12" cy="13" r="4"/></svg>
                                    <span class="text-[10px] font-black uppercase tracking-widest">Foto Baju</span>
                                </button>
                            </template>

                            <div v-show="isCameraOpen && cameraTarget === 'ITEM'" class="w-full h-full relative z-20 mt-4 rounded-lg overflow-hidden bg-black flex flex-col">
                                <video ref="videoItemRef" autoplay playsinline class="w-full h-32 object-cover"></video>
                                <div class="absolute bottom-2 left-0 right-0 flex justify-center gap-2">
                                    <button @click="takePhoto" class="bg-indigo-600 text-white text-[10px] font-black px-4 py-1.5 rounded-full uppercase shadow-lg border-2 border-white">Jepret Baju</button>
                                    <button @click="closeCamera" class="bg-rose-500 text-white text-[10px] font-black px-4 py-1.5 rounded-full uppercase shadow-lg border-2 border-white">Batal</button>
                                </div>
                            </div>

                            <div v-if="photoData" class="w-full h-full relative z-20 mt-4 rounded-lg overflow-hidden border border-slate-200">
                                <img :src="photoData" class="w-full h-32 object-cover" />
                                <button @click="photoData = null" class="absolute top-2 right-2 bg-rose-500 text-white w-6 h-6 rounded-full flex items-center justify-center shadow-lg"><svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg></button>
                            </div>
                            <canvas ref="canvasItemRef" class="hidden"></canvas>
                        </div>
                    </div>

                    <div class="p-4 md:p-5 flex flex-col gap-3">
                        <div v-if="cart.length === 0" class="flex flex-col items-center justify-center text-center opacity-60 py-5">
                            <p class="text-slate-500 font-black text-xs uppercase tracking-widest">Keranjang Kosong</p>
                        </div>
                        <div v-else v-for="(item, index) in cart" :key="index" class="bg-white p-4 rounded-2xl shadow-sm border border-slate-200 relative group">
                            <button @click="removeCartItem(index)" class="absolute -right-2 -top-2 w-7 h-7 bg-white text-rose-500 border border-slate-200 rounded-full flex items-center justify-center shadow-md z-10"><svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg></button>
                            <div class="flex justify-between items-start mb-4">
                                <div class="pr-4">
                                    <h4 class="text-sm font-black text-slate-800 uppercase leading-tight">{{ item.nama_produk }}</h4>
                                    <p class="text-[10px] font-bold text-slate-400 mt-1 uppercase tracking-wider">@ {{ formatRupiah(item.harga) }} /{{ item.satuan_dasar }}</p>
                                </div>
                                <span class="text-sm font-black text-indigo-600 bg-indigo-50 px-2.5 py-1 rounded-lg border border-indigo-100">{{ formatRupiah(item.harga * (parseFloat(item.berat) || 0)) }}</span>
                            </div>
                            <div class="flex items-center justify-between bg-slate-50 rounded-xl p-1.5 border border-slate-200">
                                <span class="text-[10px] font-black text-slate-500 uppercase tracking-widest pl-3">{{ item.satuan_dasar === 'KG' ? 'Berat:' : 'Jumlah:' }}</span>
                                <div class="flex items-center bg-white rounded-lg border border-slate-200 overflow-hidden shadow-sm">
                                    <button @click="updateBerat(index, -0.5)" class="w-9 h-9 flex items-center justify-center text-slate-500 active:bg-slate-200"><svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="5" y1="12" x2="19" y2="12"/></svg></button>
                                    <input v-model="item.berat" type="number" step="0.1" min="0.1" class="w-12 md:w-14 text-center font-black text-sm text-indigo-700 outline-none bg-transparent">
                                    <button @click="updateBerat(index, 0.5)" class="w-9 h-9 flex items-center justify-center text-slate-500 active:bg-slate-200 border-l border-slate-100"><svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg></button>
                                    <span class="text-[10px] font-black text-slate-400 bg-slate-50 h-9 flex items-center px-2 border-l border-slate-200">{{ item.satuan_dasar }}</span>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="bg-white border-t border-slate-200 shrink-0 p-4 md:p-5 shadow-[0_-10px_40px_-15px_rgba(0,0,0,0.1)] z-20">
                    <div class="mb-4 grid grid-cols-4 gap-2">
                        <button @click="paymentMethod = 'TUNAI'" class="py-2 md:py-3 rounded-xl border-2 flex flex-col items-center gap-1 transition-all" :class="paymentMethod === 'TUNAI' ? 'bg-slate-800 border-slate-800 text-white shadow-md' : 'bg-white border-slate-200 text-slate-500'"><span class="text-[9px] font-black tracking-widest mt-1">TUNAI</span></button>
                        <button @click="paymentMethod = 'QRIS'" class="py-2 md:py-3 rounded-xl border-2 flex flex-col items-center gap-1 transition-all" :class="paymentMethod === 'QRIS' ? 'bg-slate-800 border-slate-800 text-white shadow-md' : 'bg-white border-slate-200 text-slate-500'"><span class="text-[9px] font-black tracking-widest mt-1">QRIS</span></button>
                        <button @click="paymentMethod = 'DEBIT'" class="py-2 md:py-3 rounded-xl border-2 flex flex-col items-center gap-1 transition-all" :class="paymentMethod === 'DEBIT' ? 'bg-slate-800 border-slate-800 text-white shadow-md' : 'bg-white border-slate-200 text-slate-500'"><span class="text-[9px] font-black tracking-widest mt-1">DEBIT</span></button>
                        <button @click="paymentMethod = 'NANTI_AJA'" class="py-2 md:py-3 rounded-xl border-2 flex flex-col items-center gap-1 transition-all" :class="paymentMethod === 'NANTI_AJA' ? 'bg-amber-500 border-amber-500 text-white shadow-md' : 'bg-amber-50 border-amber-200 text-amber-600'"><span class="text-[8px] font-black tracking-tighter mt-1">BELAKANGAN</span></button>
                    </div>

                    <div class="space-y-3 mb-4 bg-slate-50 p-3 md:p-4 rounded-xl border border-slate-100">
                        <div class="flex justify-between items-center">
                            <span class="text-[10px] md:text-xs font-black text-slate-400 uppercase tracking-widest">TOTAL</span>
                            <span class="text-xl md:text-2xl font-black text-slate-800 tracking-tighter">{{ formatRupiah(totalTagihan) }}</span>
                        </div>
                        <div class="pt-3 border-t border-slate-200" v-if="paymentMethod === 'TUNAI'">
                            <div class="flex justify-between items-center mb-3">
                                <span class="text-[10px] font-black text-slate-500 uppercase tracking-widest">Dibayar</span>
                                <div class="flex items-center relative">
                                    <span class="absolute left-3 text-sm font-bold text-slate-400">Rp</span>
                                    <input v-model="uangBayar" type="number" class="w-32 md:w-36 bg-white border-2 border-slate-200 rounded-xl pl-9 pr-3 py-2 text-right font-black text-sm md:text-base text-slate-800 outline-none focus:border-indigo-500">
                                </div>
                            </div>
                            <div class="flex justify-between items-center bg-white p-2 md:p-3 rounded-xl border border-slate-100">
                                <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Kembali</span>
                                <span class="text-sm md:text-base font-black tracking-tight" :class="kembalian < 0 ? 'text-rose-500' : 'text-emerald-500'">{{ formatRupiah(kembalian) }}</span>
                            </div>
                        </div>
                    </div>
                    
                    <button @click="processCheckout" :disabled="cart.length === 0 || isSubmitting" class="w-full bg-gradient-to-r from-sky-500 to-indigo-600 hover:from-sky-400 hover:to-indigo-500 text-white py-3.5 md:py-4 rounded-xl font-black text-xs md:text-sm uppercase tracking-[0.2em] transition-all active:scale-95 disabled:opacity-50 shadow-lg shadow-indigo-200/50">
                        {{ isSubmitting ? 'MEMPROSES...' : 'SIMPAN & CETAK NOTA' }}
                    </button>
                </div>
            </div>
        </div>

        <Teleport to="body">
            <div v-if="showQrisModal" class="fixed inset-0 z-[9999] flex items-center justify-center bg-slate-900/80 backdrop-blur-sm p-4 animate-[fadeIn_0.2s_ease-out]">
                <div class="bg-white rounded-3xl shadow-2xl w-full max-w-sm md:max-w-3xl overflow-hidden flex flex-col">
                    
                    <div class="bg-indigo-600 p-5 text-center text-white shrink-0">
                        <h3 class="font-black text-lg md:text-xl uppercase tracking-widest">SCAN QRIS PEMBAYARAN</h3>
                        <p class="text-sm font-bold text-indigo-200 mt-1">Total Tagihan: {{ formatRupiah(totalTagihan) }}</p>
                    </div>

                    <div class="p-5 md:p-6 grid grid-cols-1 md:grid-cols-2 gap-4 md:gap-6 overflow-y-auto max-h-[80vh] custom-scrollbar items-stretch">
                        
                        <div class="bg-slate-50 rounded-2xl border-2 border-slate-200 p-4 flex flex-col items-center justify-center shadow-inner">
                            <img :src="qrisStoreUrl || 'https://upload.wikimedia.org/wikipedia/commons/d/d0/QR_code_for_mobile_English_Wikipedia.svg'" alt="QRIS Toko" class="w-full max-w-[220px] md:max-w-[260px] aspect-square object-contain mix-blend-multiply opacity-90" />
                            <p class="text-[10px] md:text-[11px] text-center font-bold text-slate-500 uppercase tracking-widest leading-relaxed mt-4">
                                Minta pelanggan scan kode ini
                            </p>
                        </div>

                        <div class="w-full border-2 border-dashed border-emerald-300 rounded-2xl p-3 bg-emerald-50/50 flex flex-col items-center justify-center relative overflow-hidden min-h-[200px] transition-all" :class="buktiTransferData ? 'border-solid border-emerald-500 bg-emerald-50' : ''">
                            <template v-if="(!isCameraOpen || cameraTarget !== 'QRIS') && !buktiTransferData">
                                <button @click="openCamera('QRIS')" class="flex flex-col items-center text-emerald-600 hover:text-emerald-700 transition-colors my-2">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-12 h-12 mb-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/></svg>
                                    <span class="text-[10px] font-black uppercase tracking-widest text-center">Jepret Layar HP<br>Pelanggan (Bukti)</span>
                                </button>
                            </template>

                            <div v-show="isCameraOpen && cameraTarget === 'QRIS'" class="w-full h-full relative z-20 rounded-xl overflow-hidden bg-black flex flex-col">
                                <video ref="videoQrisRef" autoplay playsinline class="w-full h-full object-cover min-h-[200px]"></video>
                                <div class="absolute bottom-3 left-0 right-0 flex justify-center gap-2 z-30">
                                    <button @click="takePhoto" class="bg-emerald-500 text-white text-[10px] font-black px-4 py-2 rounded-full uppercase shadow-lg border-2 border-white active:scale-95">Jepret Bukti</button>
                                </div>
                            </div>

                            <div v-if="buktiTransferData" class="w-full h-full relative z-20 rounded-xl overflow-hidden border border-emerald-200 shadow-sm flex items-center justify-center bg-black/5">
                                <img :src="buktiTransferData" class="w-full h-full object-cover min-h-[200px]" />
                                <button @click="buktiTransferData = null" class="absolute top-3 right-3 bg-rose-500 text-white w-8 h-8 rounded-full flex items-center justify-center shadow-lg hover:bg-rose-600 transition-colors"><svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg></button>
                            </div>
                        </div>
                    </div>

                    <div class="p-5 bg-slate-50 border-t border-slate-100 flex gap-3 shrink-0">
                        <button @click="cancelQris" class="w-1/3 py-3.5 text-slate-500 font-black text-xs uppercase tracking-widest hover:bg-slate-200 rounded-xl transition-colors">Batal</button>
                        <button @click="confirmQris" :disabled="!buktiTransferData" class="w-2/3 py-3.5 bg-emerald-500 text-white font-black text-xs uppercase tracking-widest rounded-xl hover:bg-emerald-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors shadow-lg shadow-emerald-200 flex justify-center items-center gap-2">
                            Sahkan Transaksi
                        </button>
                    </div>
                </div>
            </div>
            <canvas ref="canvasQrisRef" class="hidden"></canvas>
        </Teleport>

    <div id="print-area" class="print-only">
        <div class="receipt" v-if="printData">
            <h2 class="center bold uppercase">{{ printData.toko_nama }}</h2>
            <p class="center text-sm" v-if="printData.toko_alamat">{{ printData.toko_alamat }}</p>
            <p class="center text-sm" v-if="printData.toko_telepon">Telp: {{ printData.toko_telepon }}</p>
            <hr class="dash">
            <p>Resi    : {{ printData.invoice }}</p>
            <p>Tanggal : {{ printData.tanggal }}</p>
            <hr class="dash">
            <p>Pelanggan : <b>{{ printData.pelanggan }}</b></p>
            <p>Selesai   : {{ printData.estimasi }}</p>
            <hr class="dash">
            <div class="item-list" v-for="(item, index) in printData.items" :key="index">
                <p>{{ item.nama_produk }}</p>
                <div class="item-calc">
                    <span>{{ item.berat }} {{ item.satuan_dasar }} x {{ formatRupiah(item.harga) }}</span>
                    <span>{{ formatRupiah(item.berat * item.harga) }}</span>
                </div>
            </div>
            <hr class="dash">
            <div class="total-row bold text-lg"><span>TOTAL</span><span>{{ formatRupiah(printData.total) }}</span></div>
            <div class="total-row"><span>Metode</span><span>{{ printData.metode }}</span></div>
            <div class="total-row"><span>Bayar</span><span>{{ formatRupiah(printData.bayar) }}</span></div>
            <div class="total-row" v-if="printData.metode === 'TUNAI'"><span>Kembali</span><span>{{ formatRupiah(printData.kembali) }}</span></div>
            <hr class="dash">
            <p class="center bold text-sm mt-3">{{ printData.toko_footer }}</p>
        </div>
    </div>
    </SidebarLaundry>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { height: 6px; width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
input[type=number]::-webkit-inner-spin-button, input[type=number]::-webkit-outer-spin-button { -webkit-appearance: none; margin: 0; }
.print-only { display: none; }
@media print {
    @page { margin: 0; size: 80mm auto; }
    body * { visibility: hidden; }
    .hide-on-print { display: none !important; }
    .print-only, .print-only * { visibility: visible; }
    .print-only { display: block; position: absolute; left: 0; top: 0; width: 78mm; padding: 2mm; background: white; }
    .receipt { font-family: 'Courier New', Courier, monospace; color: black; line-height: 1.3; }
    .receipt p, .receipt span { font-size: 12px; margin: 2px 0; }
    .center { text-align: center; } .bold { font-weight: bold; } .text-sm { font-size: 10px !important; } .text-lg { font-size: 14px !important; } .mt-3 { margin-top: 12px; }
    hr.dash { border: none; border-top: 1px dashed black; margin: 8px 0; }
    .item-list { margin-bottom: 6px; } .item-calc { display: flex; justify-content: space-between; padding-left: 10px; } .total-row { display: flex; justify-content: space-between; margin-bottom: 2px; }
}
@keyframes fadeIn { from { opacity: 0; transform: translateY(-10px); } to { opacity: 1; transform: translateY(0); } }
</style>