<script setup>
import Swal from 'sweetalert2';
import { computed, onMounted, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import api from '../../../api.js';

// 🚀 IMPORT CHILD COMPONENT
import RegionSelector from './RegionSelector.vue'; // Sesuaikan path foldernya !

const route = useRoute();
const router = useRouter();
const isExpansion = route.query.is_expansion === 'true';

const isResumingPayment = ref(false);

const pendingIndustry = localStorage.getItem('pendingIndustry') || 'retail';
const pendingPlan = localStorage.getItem('pendingPlan') || 'trial';

const kategoriOptions = [
    { id: 'Retail', label: 'Retail & Barang', icon: 'M16 11V7a4 4 0 0 0-8 0v4M5 9h14l1 12H4L5 9z' },
    { id: 'F&B', label: 'Food & Beverage', icon: 'M18 8h1a4 4 0 0 1 0 8h-1M2 8h16v9a4 4 0 0 1-4 4H6a4 4 0 0 1-4-4V8z M6 1v3 M10 1v3 M14 1v3' },
    { id: 'Jasa', label: 'Layanan & Jasa', icon: 'M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z M3.27 6.96L12 12.01l8.73-5.05 M12 22.08V12' },
    { id: 'Lainnya', label: 'Bisnis Lainnya', icon: 'M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z M22 6l-10 7L2 6' },
];

// 🚀 SINKRONISASI TOTAL: Peta Industri Berbasis Kode UPPERCASE Anti-Fraud
const detailOptions = {
    Retail: [
        { value: 'MINIMARKET', label: 'Minimarket / Toko Kelontong', isReady: true },
        { value: 'BOUTIQUE', label: 'Toko Pakaian / Butik Fashion', isReady: true },
        { value: 'ELECTRONIC', label: 'Elektronik / Toko Gadget', isReady: true },
        { value: 'PHARMACY', label: 'Apotek / Toko Farmasi', isReady: true },
        { value: 'BUILDING', label: 'Toko Bahan Bangunan / Hardware', isReady: true },
        { value: 'PETSHOP', label: 'Pet Shop / Perlengkapan Hewan', isReady: true },
        { value: 'RETAIL_LAINNYA', label: 'Retail Bisnis Lainnya', isReady: true }
    ],
    'F&B': [
        { value: 'CAFE', label: 'Cafe / Coffee Shop (Soon)', isReady: false },
        { value: 'RESTAURANT', label: 'Restoran / Rumah Makan (Soon)', isReady: false },
        { value: 'BAKERY', label: 'Bakery / Toko Kue (Soon)', isReady: false },
        { value: 'STREET_FOOD', label: 'Food Court / Kedai Kaki Lima (Soon)', isReady: false },
        { value: 'CATERING', label: 'Jasa Katering Masakan (Soon)', isReady: false }
    ],
    Jasa: [
        { value: 'LAUNDRY', label: 'Jasa Laundry Ruko (LIVE)', isReady: true },
        { value: 'BARBER', label: 'Barbershop / Salon Kecantikan (Soon)', isReady: false },
        { value: 'BENGKEL', label: 'Bengkel Otomotif / Service (Soon)', isReady: false },
        { value: 'CAR_WASH', label: 'Cuci Mobil & Motor (Soon)', isReady: false },
        { value: 'CLINIC', label: 'Klinik / Praktek Dokter Mandiri (Soon)', isReady: false },
        { value: 'SPA_REFLEXOLOGY', label: 'Spa & Tempat Refleksi (Soon)', isReady: false },
        { value: 'JASA_LAINNYA', label: 'Layanan Jasa Umum Lainnya (Soon)', isReady: false }
    ],
    Lainnya: [
        { value: 'UMUM', label: 'Bisnis Umum / Tipe Lainnya (Soon)', isReady: false }
    ],
};

onMounted(async () => {
    // 1. Injeksi SDK Midtrans
    if (!document.getElementById('midtrans-script-owner')) {
        const midtransEnv = import.meta.env.VITE_MIDTRANS_ENV || 'sandbox';
        const clientKey = import.meta.env.VITE_MIDTRANS_CLIENT_KEY;
        const snapUrl = midtransEnv === 'production' ? 'https://app.midtrans.com/snap/snap.js' : 'https://app.sandbox.midtrans.com/snap/snap.js';

        const script = document.createElement('script');
        script.id = 'midtrans-script-owner';
        script.src = snapUrl;
        script.setAttribute('data-client-key', clientKey);
        document.head.appendChild(script);
    }

    // 2. Set Kategori Bisnis awal dari local storage
    const pendingIndustry = localStorage.getItem('pendingIndustry');
    if (pendingIndustry === 'fnb') form.value.kategori_bisnis = 'F&B';
    else if (pendingIndustry === 'jasa') form.value.kategori_bisnis = 'Jasa';
    else form.value.kategori_bisnis = 'Retail';

    // 3. RESUME CONTROLLER TRIGGER
    const resumeStoreId = route.query.resume_store_id;
    if (resumeStoreId) {
        isResumingPayment.value = true;
        await handleResumePendingStore(resumeStoreId);
    }
});

const handleResumePendingStore = async (storeId) => {
    try {
        Swal.fire({
            title: 'Menghubungkan Server...',
            text: 'Membuka gerbang billing pembayaran, mohon tunggu sebentar .',
            allowOutsideClick: false,
            didOpen: () => {
                Swal.showLoading();
            },
        });

        const res = await api.post('/re-trigger-payment', { store_id: Number(storeId) });
        Swal.close();

        if (res.data.store_name) form.value.nama_toko = res.data.store_name;

        window.snap.pay(res.data.snap_token, {
            onSuccess: function (result) {
                Swal.fire({
                    icon: 'success',
                    title: 'Pembayaran Berhasil',
                    text: 'Gerai cabang Anda resmi diaktifkan global!',
                    timer: 2000,
                    showConfirmButton: false,
                    customClass: { popup: 'rounded-[32px]' },
                }).then(() => {
                    localStorage.removeItem('pendingIndustry');
                    localStorage.removeItem('pendingPlan');
                    localStorage.removeItem('temp_stores');
                    window.location.href = '/select-store';
                });
            },
            onPending: function (result) {
                Swal.fire('Menunggu Pembayaran', 'Segera selesaikan transaksi invoice Anda .', 'info').then(() => {
                    window.location.href = '/select-store';
                });
            },
            onError: function (result) {
                Swal.fire('Gagal', 'Sistem pembayaran mendeteksi anomali perbankan.', 'error').then(() => {
                    window.location.href = '/select-store';
                });
            },
            onClose: function () {
                Swal.fire({
                    title: 'Aktivasi Tertunda',
                    text: 'Tenang , konfigurasi toko tersimpan aman. Lu bisa selesaikan pembayaran kapan pun lewat menu Pilih Toko.',
                    icon: 'warning',
                    confirmButtonColor: '#4f46e5',
                    customClass: { popup: 'rounded-[32px]' },
                }).then(() => {
                    window.location.href = '/select-store';
                });
            },
        });
    } catch (error) {
        Swal.close();
        console.error('Gagal memuat billing resume:', error);
        Swal.fire('Error', 'Infrastruktur billing delay. Coba beberapa saat lagi.', 'error');
        isResumingPayment.value = false;
    }
};

const form = ref({
    nama_toko: '',
    kategori_bisnis: 'Retail',
    detail_bisnis: '',
    telepon: '',
    alamat: '',
    provinsi: '',
    kota: '',
    kecamatan: '',
    kelurahan: '',
    kode_pos: '',
});

const isLoading = ref(false);

const activeKategori = computed(() => {
    return kategoriOptions.find((k) => k.id === form.value.kategori_bisnis) || kategoriOptions[0];
});

// Watcher pintar: Auto-select opsi pertama yang "isReady" murni value UPPERCASE bray
watch(
    () => form.value.kategori_bisnis,
    (newVal) => {
        const options = detailOptions[newVal] || [];
        const readyOption = options.find(o => o.isReady) || options[0];
        form.value.detail_bisnis = readyOption ? readyOption.value : '';
    },
    { immediate: true }
);

const formatNoHp = () => {
    let val = String(form.value.telepon);
    if (val.startsWith('0')) val = val.substring(1);
    if (val.startsWith('62')) val = val.substring(2);
    form.value.telepon = val;
};

const bukaSnapMidtrans = (snapToken) => {
    window.snap.pay(snapToken, {
        onSuccess: function (result) {
            Swal.fire({
                icon: 'success',
                title: 'Pembayaran Berhasil',
                text: 'Infrastruktur premium Anda telah aktif sepenuhnya.',
                timer: 2000,
                showConfirmButton: false,
                allowOutsideClick: false,
                customClass: { popup: 'rounded-[32px]' },
            }).then(() => {
                localStorage.removeItem('temp_stores');
                localStorage.removeItem('pendingIndustry');
                localStorage.removeItem('pendingPlan');
                window.location.href = '/select-store';
            });
        },
        onPending: function (result) {
            Swal.fire('Menunggu Pembayaran', 'Segera selesaikan transaksi Anda sebelum invoice kedaluwarsa.', 'info').then(() => {
                window.location.href = '/select-store';
            });
        },
        onError: function (result) {
            Swal.fire('Pembayaran Gagal', 'Terjadi kesalahan sistem perbankan.', 'error').then(() => {
                window.location.href = '/select-store';
            });
        },
        onClose: function () {
            Swal.fire({
                title: 'Aktivasi Ditunda',
                text: 'Konfigurasi Toko Anda tersimpan aman. Selesaikan aktivasi kapan pun melalui halaman Pilih Toko.',
                icon: 'warning',
                confirmButtonColor: '#4f46e5',
                customClass: { popup: 'rounded-[32px]' },
            }).then(() => {
                window.location.href = '/select-store';
            });
        },
    });
};

const submit = async () => {
    if (!form.value.kelurahan) {
        return Swal.fire('Data Kurang', 'Harap lengkapi pilihan Kelurahan atau Desa terlebih dahulu.', 'warning');
    }
    isLoading.value = true;

    try {
        // 🚀 FIX CRITICAL: Ambil nilai asli UPPERCASE (Contoh: 'LAUNDRY'), jangan di-toLowerCase() bray!
        const finalTipeBisnis = String(form.value.detail_bisnis || form.value.kategori_bisnis.toUpperCase());
        const currentPendingIndustry = localStorage.getItem('pendingIndustry') || 'retail';
        const currentPendingPlan = (localStorage.getItem('pendingPlan') || 'trial').toLowerCase();
        const existingOwnerToken = localStorage.getItem('token');

        const payload = {
            nama_toko: form.value.nama_toko,
            telepon: `62${form.value.telepon}`,
            business_type: finalTipeBisnis, // Bersih murni UPPERCASE bray!
            industry: currentPendingIndustry,
            plan: currentPendingPlan,
            alamat_toko: form.value.alamat,
            provinsi: form.value.provinsi,
            kota: form.value.kota,
            kecamatan: form.value.kecamatan,
            kelurahan: form.value.kelurahan,
            kode_pos: String(form.value.kode_pos),
        };

        const response = await api.post('/setup', payload);
        const tokenTerupdate = response.data?.token;

        if (tokenTerupdate) {
            api.defaults.headers.common['Authorization'] = `Bearer ${tokenTerupdate}`;
            localStorage.setItem('token', tokenTerupdate);
            localStorage.setItem('store_id', response.data.store_id);
            localStorage.setItem('storeName', response.data.store_name || 'POS UMKM');
            localStorage.setItem('subscriptionPlan', response.data.subscription_plan || 'basic');
            localStorage.setItem('role', 'owner');

            const oldStoresRaw = localStorage.getItem('temp_stores');
            let currentStores = oldStoresRaw ? JSON.parse(oldStoresRaw) : [];

            const newStoreObj = {
                id: response.data.store_id,
                nama_toko: form.value.nama_toko,
                industry: currentPendingIndustry,
                subscription_plan: currentPendingPlan,
                subscription_status: response.data.subscription_status || 'pending',
                kota: form.value.kota || 'Lokasi Belum Diatur',
            };

            currentStores.push(newStoreObj);
            localStorage.setItem('temp_stores', JSON.stringify(currentStores));
        }

        if (currentPendingPlan === 'basic' || currentPendingPlan === 'pro' || currentPendingPlan === 'premium') {
            isLoading.value = false;

            Swal.fire({
                title: 'Menyiapkan Pembayaran',
                text: 'Menghubungkan ke gerbang aman Midtrans untuk verifikasi lisensi...',
                allowOutsideClick: false,
                didOpen: () => {
                    Swal.showLoading();
                },
            });

            try {
                const activeToken = tokenTerupdate || existingOwnerToken;
                if (!activeToken) throw new Error('Sesi token owner tidak ditemukan di sistem. Harap login ulang.');

                const payRes = await api.post('/retail/subscription/upgrade', { plan_name: currentPendingPlan }, { headers: { Authorization: `Bearer ${activeToken}` } });
                const snapToken = payRes.data.token;
                Swal.close();
                bukaSnapMidtrans(snapToken);
            } catch (err) {
                Swal.close();
                if (err.response && err.response.status === 402) {
                    const backupSnapToken = err.response.data?.token || err.response.data?.snap_token;
                    if (backupSnapToken) {
                        bukaSnapMidtrans(backupSnapToken);
                        return;
                    }
                }
                Swal.fire('Error Gateway', err.response?.data?.error || err.message || 'Gagal memanggil gerbang pembayaran.', 'error').then(() => {
                    window.location.href = '/retail/account';
                });
            }
        } else {
            localStorage.removeItem('pendingIndustry');
            localStorage.removeItem('pendingPlan');

            await Swal.fire({
                icon: 'success',
                title: 'Infrastruktur Ready !',
                text: 'Selamat menikmati fasilitas Free Trial selama 14 hari.',
                confirmButtonColor: '#4f46e5',
                customClass: { popup: 'rounded-[32px]' },
            });

            const kat = form.value.kategori_bisnis;
            const det = (form.value.detail_bisnis || '').toUpperCase(); // Menyesuaikan pengecekan dengan UPPERCASE bray!

            if (kat === 'Retail' || kat === 'Lainnya') window.location.href = '/retail/dashboard';
            else if (kat === 'F&B') window.location.href = '/fnb/dashboard';
            else if (kat === 'Jasa') {
                if (det.includes('LAUNDRY')) window.location.href = '/laundry/laporan';
                else if (det.includes('BENGKEL')) window.location.href = '/bengkel/dashboard';
                else if (det.includes('BARBER')) window.location.href = '/salon/dashboard';
                else if (det.includes('CAR_WASH')) window.location.href = '/cuci-kendaraan/dashboard';
                else window.location.href = '/retail/dashboard';
            }
        }
    } catch (error) {
        Swal.fire({
            icon: 'error',
            title: 'Gagal Setup Toko',
            text: error.response?.data?.error || error.message || 'Terjadi kesalahan sistem database.',
            confirmButtonColor: '#ef4444',
        });
    } finally {
        isLoading.value = false;
    }
};
</script>

<template>
    <div class="min-h-screen bg-[#F8FAFC] flex flex-col justify-center py-10 md:py-16 sm:px-6 lg:px-8 font-sans relative overflow-hidden">
        <div class="absolute -top-24 -left-24 w-[30rem] h-[30rem] bg-indigo-200/40 rounded-full blur-3xl pointer-events-none"></div>
        <div class="absolute -bottom-24 -right-24 w-[30rem] h-[30rem] bg-blue-200/40 rounded-full blur-3xl pointer-events-none"></div>

        <div v-if="isResumingPayment" class="sm:mx-auto sm:w-full sm:max-w-xl px-4 relative z-10 animate-fade-in-up">
            <div class="bg-white p-8 md:p-12 text-center shadow-2xl rounded-[40px] border border-slate-100/50 flex flex-col items-center">
                <div class="w-14 h-14 border-4 border-indigo-100 border-t-indigo-600 rounded-full animate-spin mb-6"></div>
                <h3 class="font-black text-slate-900 uppercase text-base tracking-wider">Menyambungkan Ke Midtrans</h3>
                <p class="text-[11px] font-bold text-slate-400 mt-2 uppercase leading-relaxed max-w-xs mx-auto">Sistem mendeteksi transaksi tertunda. Mengembalikan data invoice gerai Anda secara realtime...</p>
            </div>
        </div>

        <div v-else class="w-full flex flex-col items-center">
            <div class="sm:mx-auto sm:w-full sm:max-w-2xl text-center relative z-10 px-4">
                <div class="w-20 h-20 bg-gradient-to-br from-indigo-600 to-blue-600 rounded-[24px] flex items-center justify-center mx-auto shadow-2xl shadow-indigo-200 mb-6 transform -rotate-6 transition-transform hover:rotate-0 duration-500 border-4 border-white">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-10 h-10 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M3 9h18v10a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V9Z" />
                        <path d="m3 9 2.45-4.9A2 2 0 0 1 7.24 3h9.52a2 2 0 0 1 1.8 1.1L21 9" />
                        <path d="M12 3v6" />
                    </svg>
                </div>
                <h2 class="text-3xl md:text-4xl font-black text-slate-900 tracking-tighter">
                    Setup Infrastruktur
                    <span class="text-indigo-600">Bisnis</span>
                </h2>
                <p class="mt-3 text-slate-400 font-bold text-[10px] md:text-xs uppercase tracking-[0.2em]">Konfigurasi Sistem Menyesuaikan Alur Kerja Anda</p>
                <div class="mt-4 inline-flex items-center gap-2 bg-indigo-100 text-indigo-700 px-4 py-2 rounded-full font-black text-[10px] uppercase tracking-widest shadow-sm">Paket Aktif: {{ pendingIndustry }} - {{ pendingPlan }}</div>
            </div>

            <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-3xl px-4 relative z-10 w-full">
                <div class="bg-white/90 backdrop-blur-xl p-6 md:p-10 shadow-2xl shadow-slate-200/50 rounded-[32px] md:rounded-[40px] border border-white">
                    <form @submit.prevent="submit" class="flex flex-col gap-10">
                        <div class="flex flex-col gap-6">
                            <div class="flex items-center gap-3 border-b border-slate-100 pb-3">
                                <div class="w-8 h-8 rounded-full bg-indigo-50 flex items-center justify-center text-indigo-600 font-black text-xs border border-indigo-100 shadow-sm">1</div>
                                <h3 class="text-lg font-black text-slate-800 uppercase tracking-tight">Identitas Bisnis</h3>
                            </div>

                            <div class="grid grid-cols-1 md:grid-cols-2 gap-5 md:gap-6">
                                <div class="md:col-span-2">
                                    <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block">Nama Brand / Toko</label>
                                    <input v-model="form.nama_toko" type="text" required class="input-modern text-lg" placeholder="Contoh: Ruko JKT2 Laundry Barokah..." />
                                </div>

                                <div class="md:col-span-2">
                                    <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block">Kategori Industri</label>
                                    <div class="p-4 rounded-[20px] border-2 border-indigo-200 bg-indigo-50/50 flex items-center gap-4 select-none opacity-90">
                                        <div class="w-12 h-12 bg-white rounded-xl shadow-sm border border-indigo-100 flex items-center justify-center">
                                            <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-indigo-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                                                <path :d="activeKategori.icon" />
                                            </svg>
                                        </div>
                                        <div class="flex flex-col">
                                            <span class="font-black text-slate-800 uppercase tracking-widest">{{ activeKategori.label }}</span>
                                            <span class="text-[10px] font-bold text-indigo-500 uppercase tracking-widest">Pilihan Dari Landing Page</span>
                                        </div>
                                        <div class="ml-auto text-slate-400" title="Kategori telah dikunci">
                                            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                                                <rect width="18" height="11" x="3" y="11" rx="2" ry="2" />
                                                <path d="M7 11V7a5 5 0 0 1 10 0v4" />
                                            </svg>
                                        </div>
                                    </div>
                                </div>

                                <div class="md:col-span-2 animate-[fadeIn_0.3s_ease-out]">
                                    <label class="text-[10px] font-black text-indigo-500 uppercase tracking-widest ml-1 mb-2 flex items-center gap-2">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                                            <polyline points="9 18 15 12 9 6" />
                                        </svg>
                                        Spesifikasi Model Bisnis {{ form.kategori_bisnis }}
                                    </label>
                                    <div class="relative">
                                        <select v-model="form.detail_bisnis" required class="input-modern bg-white cursor-pointer appearance-none text-indigo-900 border-slate-200 focus:border-indigo-500">
                                            <option 
                                                v-for="opt in detailOptions[form.kategori_bisnis]" 
                                                :key="opt.value" 
                                                :value="opt.value"
                                                :disabled="!opt.isReady"
                                                :class="!opt.isReady ? 'text-slate-300 bg-slate-50/50' : 'text-slate-800 font-bold'"
                                            >
                                                {{ opt.label }}
                                            </option>
                                        </select>
                                        <div class="absolute inset-y-0 right-5 flex items-center pointer-events-none">
                                            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-indigo-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                                                <path d="m6 9 6 6 6-6" />
                                            </svg>
                                        </div>
                                    </div>
                                </div>

                                <div class="md:col-span-2">
                                    <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block">No. WhatsApp Bisnis</label>
                                    <div class="flex items-center bg-white border-2 border-slate-200 rounded-2xl focus-within:border-indigo-500 focus-within:ring-4 focus-within:ring-indigo-500/10 transition-all shadow-sm overflow-hidden">
                                        <div class="pl-5 pr-4 py-4 bg-slate-50 border-r border-slate-200 flex items-center justify-center select-none">
                                            <span class="text-slate-500 font-black text-sm">+62</span>
                                        </div>
                                        <input v-model="form.telepon" @input="formatNoHp" type="number" required class="w-full px-4 py-4 bg-transparent outline-none font-black text-slate-800 placeholder:text-slate-300 placeholder:font-bold" placeholder="81234567890" />
                                    </div>
                                </div>
                            </div>
                        </div>

                        <RegionSelector :formData="form" />

                        <div class="pt-6 mt-2 border-t border-slate-100">
                            <button type="submit" :disabled="isLoading" class="btn-submit">
                                <template v-if="!isLoading">
                                    Luncurkan Bisnis Sekarang
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 ml-2" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                                        <path d="M5 12h14" />
                                        <path d="m12 5 7 7-7 7" />
                                    </svg>
                                </template>
                                <template v-else>
                                    <div class="w-5 h-5 border-4 border-white/30 border-t-white rounded-full animate-spin mr-3"></div>
                                    MENGKONFIGURASI SISTEM...
                                </template>
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
        <p class="mt-10 mb-6 text-center text-[9px] font-black text-slate-400 uppercase tracking-[0.3em]">ARZURA POS Operations &copy; 2026</p>
    </div>
</template>

<style scoped>
.input-modern {
    @apply block w-full px-5 py-4 bg-white border-2 border-slate-200 rounded-2xl focus:bg-white focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 outline-none font-black text-slate-800 transition-all placeholder:text-slate-300 placeholder:font-bold shadow-sm;
}
.btn-submit {
    @apply w-full flex items-center justify-center py-5 md:py-6 px-6 rounded-[24px] shadow-2xl shadow-indigo-200/50 text-xs md:text-sm font-black text-white bg-indigo-600 hover:bg-slate-900 transition-all active:scale-95 disabled:opacity-50 uppercase tracking-[0.2em];
}
input[type='number']::-webkit-inner-spin-button,
input[type='number']::-webkit-outer-spin-button {
    -webkit-appearance: none;
    margin: 0;
}
input[type='number'] {
    -moz-appearance: textfield;
}
</style>