<script setup>
import { ref, watch, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import Swal from 'sweetalert2';
import api from '../../api.js';

const router = useRouter();

// --- DATA MASTER KATEGORI & DETAIL BISNIS ---
const kategoriOptions = [
    { id: 'Retail', label: 'Retail & Barang', icon: 'M16 11V7a4 4 0 0 0-8 0v4M5 9h14l1 12H4L5 9z' },
    { id: 'F&B', label: 'Food & Beverage', icon: 'M18 8h1a4 4 0 0 1 0 8h-1M2 8h16v9a4 4 0 0 1-4 4H6a4 4 0 0 1-4-4V8z M6 1v3 M10 1v3 M14 1v3' },
    { id: 'Jasa', label: 'Layanan & Jasa', icon: 'M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z M3.27 6.96L12 12.01l8.73-5.05 M12 22.08V12' },
    { id: 'Lainnya', label: 'Bisnis Lainnya', icon: 'M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z M22 6l-10 7L2 6' }
];

const detailOptions = {
    'Retail': ['Minimarket / Toko Kelontong', 'Toko Pakaian / Butik', 'Elektronik / Gadget', 'Apotek / Farmasi', 'Toko Bangunan', 'Pet Shop', 'Lainnya'],
    'F&B': ['Restoran / Rumah Makan', 'Cafe / Coffee Shop', 'Bakery / Toko Roti', 'Food Court / Kaki Lima', 'Katering', 'Lainnya'],
    'Jasa': ['Laundry', 'Barbershop / Salon', 'Cuci Mobil / Motor', 'Bengkel Otomotif', 'Klinik / Praktek', 'Lainnya'],
    'Lainnya': ['Bisnis Umum / Lainnya']
};

// --- API WILAYAH INDONESIA ---
const BASE_WILAYAH_API = 'https://www.emsifa.com/api-wilayah-indonesia/api';

const listProvinsi = ref([]);
const listKota = ref([]);
const listKecamatan = ref([]);
const listKelurahan = ref([]);

// State sementara buat nangkep ID wilayah
const regIds = ref({
    provinsi: '',
    kota: '',
    kecamatan: '',
    kelurahan: ''
});

// Load daftar Provinsi pas halaman pertama dibuka
const loadProvinsi = async () => {
    try {
        const res = await fetch(`${BASE_WILAYAH_API}/provinces.json`);
        listProvinsi.value = await res.json();
    } catch (e) { console.error("Gagal load Provinsi"); }
};

onMounted(() => loadProvinsi());

// 🚀 Rantai Watcher Wilayah (Auto-Cascading)
watch(() => regIds.value.provinsi, async (newId) => {
    // Reset bawahan
    regIds.value.kota = ''; listKota.value = [];
    regIds.value.kecamatan = ''; listKecamatan.value = [];
    regIds.value.kelurahan = ''; listKelurahan.value = [];
    form.value.provinsi = listProvinsi.value.find(p => p.id === newId)?.name || '';
    
    if (newId) {
        const res = await fetch(`${BASE_WILAYAH_API}/regencies/${newId}.json`);
        listKota.value = await res.json();
    }
});

watch(() => regIds.value.kota, async (newId) => {
    regIds.value.kecamatan = ''; listKecamatan.value = [];
    regIds.value.kelurahan = ''; listKelurahan.value = [];
    form.value.kota = listKota.value.find(p => p.id === newId)?.name || '';
    
    if (newId) {
        const res = await fetch(`${BASE_WILAYAH_API}/districts/${newId}.json`);
        listKecamatan.value = await res.json();
    }
});

watch(() => regIds.value.kecamatan, async (newId) => {
    regIds.value.kelurahan = ''; listKelurahan.value = [];
    form.value.kecamatan = listKecamatan.value.find(p => p.id === newId)?.name || '';
    
    if (newId) {
        const res = await fetch(`${BASE_WILAYAH_API}/villages/${newId}.json`);
        listKelurahan.value = await res.json();
    }
});

watch(() => regIds.value.kelurahan, (newId) => {
    form.value.kelurahan = listKelurahan.value.find(p => p.id === newId)?.name || '';
});

// --- STATE FORM UTAMA ---
const form = ref({
    nama_toko: '',
    kategori_bisnis: 'Retail', 
    detail_bisnis: '',
    telepon: '',
    fitur_opsional: ['kasir', 'inventory'], 
    // Pecahan Alamat
    alamat_jalan: '',
    provinsi: '',
    kota: '',
    kecamatan: '',
    kelurahan: '',
    kode_pos: ''
});

const isLoading = ref(false);

watch(() => form.value.kategori_bisnis, (newVal) => {
    form.value.detail_bisnis = detailOptions[newVal][0]; 
}, { immediate: true });

const formatNoHp = () => {
    let val = String(form.value.telepon);
    if (val.startsWith('0')) val = val.substring(1);
    if (val.startsWith('62')) val = val.substring(2);
    form.value.telepon = val;
};

const submit = async () => {
    if (!regIds.value.kelurahan) {
        return Swal.fire('Data Kurang!', 'Harap lengkapi pilihan Kelurahan / Desa!', 'warning');
    }

    isLoading.value = true;
    try {
        const finalTipeBisnis = `${form.value.kategori_bisnis} - ${form.value.detail_bisnis}`;
        
        // 🚀 JAHIT SEMUA ALAMAT JADI 1 KALIMAT RAPI
        const alamatLengkap = `${form.value.alamat_jalan}, Kel. ${form.value.kelurahan}, Kec. ${form.value.kecamatan}, ${form.value.kota}, Prov. ${form.value.provinsi}, ${form.value.kode_pos}`;

        const payload = {
            nama_toko: form.value.nama_toko,
            tipe_bisnis: finalTipeBisnis,
            alamat_toko: alamatLengkap, // Dikirim full
            telepon: `62${form.value.telepon}`,
            fitur_aktif: form.value.fitur_opsional 
        };

        const response = await api.post('/setup', payload);
        
        if (response.data && response.data.token) {
            localStorage.setItem('token', response.data.token);
            localStorage.setItem('storeName', form.value.nama_toko);
            localStorage.setItem('storeType', form.value.kategori_bisnis); 
        }
        
        await Swal.fire({
            icon: 'success',
            title: 'Infrastruktur Siap!',
            text: `Selamat datang di era digital, Bos ${localStorage.getItem('name')}!`,
            confirmButtonColor: '#4f46e5',
            customClass: { popup: 'rounded-[32px]' }
        });
        
        router.push('/dashboard'); 
    } catch (error) {
        Swal.fire({
            icon: 'error',
            title: 'Gagal Setup Toko',
            text: error.response?.data?.error || 'Terjadi kesalahan sistem.',
            confirmButtonColor: '#ef4444'
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

        <div class="sm:mx-auto sm:w-full sm:max-w-2xl text-center relative z-10 px-4">
            <div class="w-20 h-20 bg-gradient-to-br from-indigo-600 to-blue-600 rounded-[24px] flex items-center justify-center mx-auto shadow-2xl shadow-indigo-200 mb-6 transform -rotate-6 transition-transform hover:rotate-0 duration-500 border-4 border-white">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-10 h-10 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M3 9h18v10a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V9Z"/><path d="m3 9 2.45-4.9A2 2 0 0 1 7.24 3h9.52a2 2 0 0 1 1.8 1.1L21 9"/><path d="M12 3v6"/></svg>
            </div>
            <h2 class="text-3xl md:text-4xl font-black text-slate-900 tracking-tighter">Setup Infrastruktur <span class="text-indigo-600">Bisnis</span></h2>
            <p class="mt-3 text-slate-400 font-bold text-[10px] md:text-xs uppercase tracking-[0.2em]">Konfigurasi Sistem Menyesuaikan Alur Kerja Anda</p>
        </div>

        <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-3xl px-4 relative z-10">
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
                                <input v-model="form.nama_toko" type="text" required class="input-modern text-lg" placeholder="Contoh: Indomaret, Laundry Bersih, dsb...">
                            </div>

                            <div class="md:col-span-2">
                                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-3 block">Kategori Industri</label>
                                <div class="grid grid-cols-2 lg:grid-cols-4 gap-3">
                                    <div v-for="kat in kategoriOptions" :key="kat.id" 
                                         @click="form.kategori_bisnis = kat.id"
                                         class="p-4 rounded-[20px] border-2 cursor-pointer transition-all flex flex-col items-center text-center gap-3"
                                         :class="form.kategori_bisnis === kat.id ? 'border-indigo-600 bg-indigo-50 shadow-md shadow-indigo-100' : 'border-slate-100 bg-white hover:border-slate-300 hover:bg-slate-50'">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8" :class="form.kategori_bisnis === kat.id ? 'text-indigo-600' : 'text-slate-400'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                                            <path :d="kat.icon" />
                                        </svg>
                                        <span class="font-black text-[10px] uppercase tracking-widest" :class="form.kategori_bisnis === kat.id ? 'text-indigo-800' : 'text-slate-500'">{{ kat.label }}</span>
                                    </div>
                                </div>
                            </div>

                            <div class="md:col-span-2 animate-[fadeIn_0.3s_ease-out]">
                                <label class="text-[10px] font-black text-indigo-500 uppercase tracking-widest ml-1 mb-2 flex items-center gap-2">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="9 18 15 12 9 6"/></svg>
                                    Spesifikasi {{ form.kategori_bisnis }}
                                </label>
                                <div class="relative">
                                    <select v-model="form.detail_bisnis" required class="input-modern bg-slate-50/50 cursor-pointer appearance-none text-indigo-900 border-indigo-100">
                                        <option v-for="opt in detailOptions[form.kategori_bisnis]" :key="opt" :value="opt">{{ opt }}</option>
                                    </select>
                                    <div class="absolute inset-y-0 right-5 flex items-center pointer-events-none">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-indigo-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="m6 9 6 6 6-6"/></svg>
                                    </div>
                                </div>
                            </div>

                            <div class="md:col-span-2">
                                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block">No. WhatsApp Bisnis</label>
                                <div class="flex items-center bg-white border-2 border-slate-200 rounded-2xl focus-within:border-indigo-500 focus-within:ring-4 focus-within:ring-indigo-500/10 transition-all shadow-sm overflow-hidden">
                                    <div class="pl-5 pr-4 py-4 bg-slate-50 border-r border-slate-200 flex items-center justify-center">
                                        <span class="text-slate-500 font-black text-sm">+62</span>
                                    </div>
                                    <input 
                                        v-model="form.telepon" 
                                        @input="formatNoHp" 
                                        type="number" 
                                        required 
                                        class="w-full px-4 py-4 bg-transparent outline-none font-black text-slate-800 placeholder:text-slate-300 placeholder:font-bold" 
                                        placeholder="81234567890"
                                    >
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="flex flex-col gap-6 pt-4 border-t border-slate-100">
                        <div class="flex items-center gap-3 border-b border-slate-100 pb-3">
                            <div class="w-8 h-8 rounded-full bg-emerald-50 flex items-center justify-center text-emerald-600 font-black text-xs border border-emerald-100 shadow-sm">2</div>
                            <h3 class="text-lg font-black text-slate-800 uppercase tracking-tight">Lokasi Operasional</h3>
                        </div>

                        <div class="grid grid-cols-1 md:grid-cols-2 gap-5 md:gap-6">
                            
                            <div class="md:col-span-2">
                                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block">Detail Alamat (Jalan, No, RT/RW)</label>
                                <textarea v-model="form.alamat_jalan" rows="2" required class="input-modern resize-none" placeholder="Contoh: Jl. Jendral Sudirman Kav 21, Gedung X Lantai 3, RT 01 / RW 02..."></textarea>
                            </div>

                            <div class="relative">
                                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block">Provinsi</label>
                                <select v-model="regIds.provinsi" required class="input-modern cursor-pointer appearance-none">
                                    <option value="" disabled selected hidden>Pilih Provinsi...</option>
                                    <option v-for="prov in listProvinsi" :key="prov.id" :value="prov.id">{{ prov.name }}</option>
                                </select>
                                <div class="absolute inset-y-0 right-5 mt-7 flex items-center pointer-events-none">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" /></svg>
                                </div>
                            </div>

                            <div class="relative">
                                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block">Kota / Kabupaten</label>
                                <select v-model="regIds.kota" :disabled="!regIds.provinsi" required class="input-modern cursor-pointer appearance-none disabled:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-60">
                                    <option value="" disabled selected hidden>{{ regIds.provinsi ? 'Pilih Kota...' : 'Pilih Provinsi Dulu' }}</option>
                                    <option v-for="kota in listKota" :key="kota.id" :value="kota.id">{{ kota.name }}</option>
                                </select>
                                <div class="absolute inset-y-0 right-5 mt-7 flex items-center pointer-events-none">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" /></svg>
                                </div>
                            </div>

                            <div class="relative">
                                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block">Kecamatan</label>
                                <select v-model="regIds.kecamatan" :disabled="!regIds.kota" required class="input-modern cursor-pointer appearance-none disabled:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-60">
                                    <option value="" disabled selected hidden>{{ regIds.kota ? 'Pilih Kecamatan...' : 'Pilih Kota Dulu' }}</option>
                                    <option v-for="kec in listKecamatan" :key="kec.id" :value="kec.id">{{ kec.name }}</option>
                                </select>
                                <div class="absolute inset-y-0 right-5 mt-7 flex items-center pointer-events-none">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" /></svg>
                                </div>
                            </div>

                            <div class="grid grid-cols-2 gap-3">
                                <div class="relative">
                                    <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block flex items-center gap-1"><span class="truncate">Desa / Kel</span></label>
                                    <select v-model="regIds.kelurahan" :disabled="!regIds.kecamatan" required class="w-full px-4 py-4 bg-white border-2 border-slate-200 rounded-2xl focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 outline-none font-black text-slate-800 transition-all text-xs cursor-pointer appearance-none disabled:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-60">
                                        <option value="" disabled selected hidden>{{ regIds.kecamatan ? 'Pilih Kelurahan...' : 'Pilih Kelurahan' }}</option>
                                        <option v-for="kel in listKelurahan" :key="kel.id" :value="kel.id" class="text-xs">{{ kel.name }}</option>
                                    </select>
                                    <div class="absolute inset-y-0 right-4 mt-7 flex items-center pointer-events-none">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" /></svg>
                                    </div>
                                </div>
                                
                                <div>
                                    <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block flex items-center gap-1">
                                        Kode Pos
                                    </label>
                                    <input 
                                        v-model="form.kode_pos" 
                                        :disabled="!regIds.kelurahan"
                                        type="number" 
                                        required 
                                        class="w-full px-4 py-4 bg-white border-2 border-slate-200 rounded-2xl focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 outline-none font-black text-slate-800 transition-all text-sm placeholder:text-slate-300 disabled:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-60 disabled:placeholder:text-slate-400" 
                                        :placeholder="regIds.kelurahan ? 'Ketik Kode Pos...' : 'Ketik Kode Pos'"
                                    >
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="flex flex-col gap-6 pt-4 border-t border-slate-100">
                        <div class="flex items-center gap-3 border-b border-slate-100 pb-3">
                            <div class="w-8 h-8 rounded-full bg-blue-50 flex items-center justify-center text-blue-600 font-black text-xs border border-blue-100 shadow-sm">3</div>
                            <h3 class="text-lg font-black text-slate-800 uppercase tracking-tight">Modul SaaS & Ekosistem</h3>
                        </div>

                        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <label class="feature-card border-indigo-200 bg-indigo-50/50 cursor-not-allowed">
                                <input type="checkbox" value="kasir" v-model="form.fitur_opsional" disabled class="w-5 h-5 text-indigo-600 rounded-lg border-gray-300">
                                <div class="ml-4 flex-1">
                                    <span class="font-black text-indigo-900 block text-sm uppercase">Cloud POS Kasir</span>
                                    <span class="text-indigo-600 text-[9px] font-bold uppercase tracking-widest italic block mt-0.5">Modul Inti (Wajib)</span>
                                </div>
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-indigo-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect width="16" height="20" x="4" y="2" rx="2" ry="2"/><line x1="12" y1="18" x2="12.01" y2="18"/></svg>
                            </label>

                            <label class="feature-card border-slate-200 hover:bg-slate-50 cursor-pointer transition-all" :class="form.fitur_opsional.includes('inventory') ? 'border-blue-400 bg-blue-50/30' : ''">
                                <input type="checkbox" value="inventory" v-model="form.fitur_opsional" class="w-5 h-5 text-blue-600 rounded flex-shrink-0 cursor-pointer">
                                <div class="ml-4 flex-1">
                                    <span class="font-black text-slate-800 block text-sm uppercase">Master Inventory</span>
                                    <span class="text-slate-400 text-[9px] font-bold uppercase tracking-widest block mt-0.5">Katalog & Stok Gudang</span>
                                </div>
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-slate-300" :class="form.fitur_opsional.includes('inventory') ? 'text-blue-500' : ''" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/><polyline points="3.27 6.96 12 12.01 20.73 6.96"/><line x1="12" y1="22.08" x2="12" y2="12"/></svg>
                            </label>

                            <label class="feature-card border-slate-200 hover:bg-slate-50 cursor-pointer transition-all" :class="form.fitur_opsional.includes('absensi') ? 'border-emerald-400 bg-emerald-50/30' : ''">
                                <input type="checkbox" value="absensi" v-model="form.fitur_opsional" class="w-5 h-5 text-emerald-600 rounded flex-shrink-0 cursor-pointer">
                                <div class="ml-4 flex-1">
                                    <span class="font-black text-slate-800 block text-sm uppercase">Biometric Attendance</span>
                                    <span class="text-slate-400 text-[9px] font-bold uppercase tracking-widest block mt-0.5">Presensi Face AI</span>
                                </div>
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-slate-300" :class="form.fitur_opsional.includes('absensi') ? 'text-emerald-500' : ''" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10"/><path d="m9 12 2 2 4-4"/></svg>
                            </label>

                            <label class="feature-card border-slate-200 hover:bg-slate-50 cursor-pointer transition-all" :class="form.fitur_opsional.includes('jadwal') ? 'border-amber-400 bg-amber-50/30' : ''">
                                <input type="checkbox" value="jadwal" v-model="form.fitur_opsional" class="w-5 h-5 text-amber-600 rounded flex-shrink-0 cursor-pointer">
                                <div class="ml-4 flex-1">
                                    <span class="font-black text-slate-800 block text-sm uppercase">Shift & Schedule</span>
                                    <span class="text-slate-400 text-[9px] font-bold uppercase tracking-widest block mt-0.5">Manajemen Jadwal Tim</span>
                                </div>
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-slate-300" :class="form.fitur_opsional.includes('jadwal') ? 'text-amber-500' : ''" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
                            </label>
                        </div>
                    </div>

                    <div class="pt-6 mt-2 border-t border-slate-100">
                        <button type="submit" :disabled="isLoading" class="btn-submit">
                            <template v-if="!isLoading">
                                Luncurkan Bisnis Sekarang
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 ml-2" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M5 12h14"/><path d="m12 5 7 7-7 7"/></svg>
                            </template>
                            <template v-else>
                                <div class="w-5 h-5 border-4 border-white/30 border-t-white rounded-full animate-spin mr-3"></div>
                                MENGKONFIGURASI SISTEM...
                            </template>
                        </button>
                    </div>

                </form>
            </div>
            <p class="mt-10 mb-6 text-center text-[9px] font-black text-slate-400 uppercase tracking-[0.3em]">Integrated Business Intelligence &copy; 2026</p>
        </div>
    </div>
</template>

<style scoped>
.input-modern {
    @apply block w-full px-5 py-4 bg-white border-2 border-slate-200 rounded-2xl focus:bg-white focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 outline-none font-black text-slate-800 transition-all placeholder:text-slate-300 placeholder:font-bold shadow-sm;
}

.feature-card {
    @apply flex items-center p-4 md:p-5 border-2 rounded-[20px] transition-all;
}

.btn-submit {
    @apply w-full flex items-center justify-center py-5 md:py-6 px-6 rounded-[24px] shadow-2xl shadow-indigo-200/50 text-xs md:text-sm font-black text-white bg-indigo-600 hover:bg-slate-900 transition-all active:scale-95 disabled:opacity-50 uppercase tracking-[0.2em];
}

/* Custom Scrollbar */
textarea::-webkit-scrollbar { width: 6px; }
textarea::-webkit-scrollbar-track { background: transparent; }
textarea::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }

@keyframes fadeIn {
    from { opacity: 0; transform: translateY(-5px); }
    to { opacity: 1; transform: translateY(0); }
}

/* Hilangkan arrow up/down di input type number */
input[type=number]::-webkit-inner-spin-button, 
input[type=number]::-webkit-outer-spin-button { 
    -webkit-appearance: none; 
    margin: 0; 
}
input[type=number] {
    -moz-appearance: textfield;
}
</style>