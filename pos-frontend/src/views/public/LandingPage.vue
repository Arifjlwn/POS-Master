<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();
const isMenuOpen = ref(false);
const isScrolled = ref(false);
const activePricingTab = ref('retail');

// Menangani efek blur navigasi saat di-scroll
onMounted(() => {
    window.addEventListener('scroll', () => {
        isScrolled.value = window.scrollY > 20;
    });
});

const features = [
    {
        title: "Omnichannel POS",
        desc: "Sistem kasir cerdas yang beradaptasi dengan alur kerja Anda. Dari scan barcode, manajemen meja, hingga tracking status layanan.",
        icon: `<svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="16" height="20" x="4" y="2" rx="2" ry="2"/><line x1="12" y1="18" x2="12.01" y2="18"/></svg>`
    },
    {
        title: "Smart Inventory & BOM",
        desc: "Kontrol stok presisi tinggi. Mendukung manajemen multi-gudang dan sistem Resep (BOM) otomatis untuk industri F&B.",
        icon: `<svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m7.5 4.27 9 5.15"/><path d="M21 8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16Z"/><path d="m3.3 7 8.7 5 8.7-5"/><path d="M12 22V12"/></svg>`
    },
    {
        title: "Workforce Management",
        desc: "Kelola rotasi shift, sistem komisi (bagi hasil), dan pantau kehadiran karyawan dengan teknologi Face AI & Geofencing.",
        icon: `<svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M22 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>`
    },
    {
        title: "Automated CRM & WA",
        desc: "Bangun loyalitas pelanggan. Kirim e-receipt, status pesanan, hingga promo eksklusif otomatis via integrasi WhatsApp.",
        icon: `<svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"/></svg>`
    }
];

const industries = [
    { id: 'retail', title: "Retail & Distribusi", desc: "Supermarket, Butik, Elektronik", icon: "M16 11V7a4 4 0 0 0-8 0v4M5 9h14l1 12H4L5 9z" },
    { id: 'fnb', title: "Food & Beverage", desc: "Cafe, Restoran, Franchise", icon: "M18 8h1a4 4 0 0 1 0 8h-1M2 8h16v9a4 4 0 0 1-4 4H6a4 4 0 0 1-4-4V8z M6 1v3 M10 1v3 M14 1v3" },
    { id: 'jasa', title: "Layanan & Jasa", desc: "Laundry, Barbershop, Bengkel", icon: "M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z M3.27 6.96L12 12.01l8.73-5.05" }
];

// 🚀 MASTER DATA PRICING & FITUR (SIAP UNTUK SCALING)
const pricingPlans = {
    retail: [
        { id: 'trial', name: 'Starter Trial', price: 'Rp 0', duration: '14 Hari', desc: 'Validasi kesesuaian sistem dengan ekosistem bisnis Anda.', features: ['POS Kasir Retail', 'Master Data Produk', 'Scan Barcode Reader', 'Tanpa Kartu Kredit'] },
        { id: 'basic', name: 'Retail Basic', price: '49k', duration: '/Bulan', desc: 'Solusi solid untuk toko kelontong dengan 1 titik kasir.', features: ['Kasir Tanpa Batas', 'Manajemen Stok Dasar', 'Struk Thermal Bluetooth', 'Laporan Penjualan Harian'] },
        { id: 'pro', name: 'Retail Pro', price: '149k', duration: '/Bulan', desc: 'Cocok untuk minimarket yang mulai mengelola karyawan.', features: ['Semua Fitur Basic', 'Manajemen Hak Akses Kasir', 'Smart Attendance & Shift', 'Laporan Ekspor (Excel/PDF)'] },
        { id: 'premium', name: 'Retail Enterprise', price: '299k', duration: '/Bulan', desc: 'Kendali penuh untuk bisnis multi-cabang & gudang.', features: ['Semua Fitur Pro', 'Sistem Multi-Cabang (HO)', 'Manajemen Multi-Gudang', 'Dedicated Support 24/7'] }
    ],
    fnb: [
        { id: 'trial', name: 'Starter Trial', price: 'Rp 0', duration: '14 Hari', desc: 'Validasi kesesuaian sistem dengan alur dapur Anda.', features: ['POS Kasir F&B', 'Manajemen Menu & Kategori', 'Hold/Simpan Pesanan', 'Tanpa Kartu Kredit'] },
        { id: 'basic', name: 'F&B Basic', price: '59k', duration: '/Bulan', desc: 'Sistem operasional efisien untuk kedai atau coffee shop.', features: ['Manajemen Layout Meja', 'Cetak Tiket Dapur (Kitchen)', 'Pajak & Service Charge', 'Struk Thermal Bluetooth'] },
        { id: 'pro', name: 'F&B Pro', price: '169k', duration: '/Bulan', desc: 'Untuk restoran sibuk dengan kontrol bahan baku ketat.', features: ['Semua Fitur Basic', 'Resep Bahan Baku (BOM)', 'Split Bill & Gabung Meja', 'Smart Attendance & Shift'] },
        { id: 'premium', name: 'F&B Enterprise', price: '349k', duration: '/Bulan', desc: 'Skalabilitas franchise dengan analitik terpusat.', features: ['Semua Fitur Pro', 'Manajemen Franchise/Cabang', 'Self-Order QR Menu', 'Dedicated Support 24/7'] }
    ],
    jasa: [
        { id: 'trial', name: 'Starter Trial', price: 'Rp 0', duration: '14 Hari', desc: 'Coba modul layanan untuk bengkel, salon, atau laundry.', features: ['POS Layanan Jasa', 'Database Pelanggan Dasar', 'Penerimaan Pesanan', 'Tanpa Kartu Kredit'] },
        { id: 'basic', name: 'Service Basic', price: '49k', duration: '/Bulan', desc: 'Sistem tracking pesanan yang rapi untuk bisnis jasa kecil.', features: ['Tracking Status Pesanan', 'Cetak Nota / Resi Barcode', 'Manajemen Layanan & Tarif', 'Laporan Pendapatan'] },
        { id: 'pro', name: 'Service Pro', price: '159k', duration: '/Bulan', desc: 'Sistem otomatisasi notifikasi dan komisi karyawan.', features: ['Semua Fitur Basic', 'Bagi Hasil / Komisi Karyawan', 'Notif WA Otomatis (Selesai)', 'Smart Attendance & Shift'] },
        { id: 'premium', name: 'Service Enterprise', price: '329k', duration: '/Bulan', desc: 'Manajemen loyalitas pelanggan dan penjadwalan tingkat lanjut.', features: ['Semua Fitur Pro', 'Sistem Booking & Reservasi', 'Sistem Poin Loyalty', 'Manajemen Multi-Cabang'] }
    ]
};
</script>

<template>
    <div class="min-h-screen bg-[#F8FAFC] font-sans text-slate-900 selection:bg-indigo-100 selection:text-indigo-600 overflow-x-hidden">
        
        <nav :class="['fixed top-0 w-full z-50 transition-all duration-300', isScrolled ? 'bg-white/85 backdrop-blur-xl border-b border-slate-200 py-4 shadow-sm' : 'bg-transparent py-6']">
            <div class="max-w-7xl mx-auto px-6 flex items-center justify-between">
                <div class="flex flex-col">
                    <div class="font-black text-2xl md:text-3xl tracking-tighter">NEXA<span class="text-indigo-600">POS</span></div>
                </div>
                
                <div class="hidden md:flex items-center gap-8">
                    <a href="#industri" class="text-xs font-black uppercase tracking-widest text-slate-500 hover:text-indigo-600 transition-colors">Solusi</a>
                    <a href="#features" class="text-xs font-black uppercase tracking-widest text-slate-500 hover:text-indigo-600 transition-colors">Teknologi</a>
                    <a href="#pricing" class="text-xs font-black uppercase tracking-widest text-slate-500 hover:text-indigo-600 transition-colors">Harga</a>
                    <router-link to="/login" class="text-xs font-black uppercase tracking-widest text-slate-800 hover:text-indigo-600 transition-colors ml-2">Masuk</router-link>
                    <a href="#pricing" class="bg-slate-900 text-white text-[10px] font-black uppercase tracking-[0.2em] px-7 py-3.5 rounded-full hover:bg-indigo-600 shadow-xl hover:shadow-indigo-200 transition-all active:scale-95">Mulai Trial</a>
                </div>

                <button @click="isMenuOpen = !isMenuOpen" class="md:hidden p-2 text-slate-600 hover:text-indigo-600 transition-colors">
                    <svg v-if="!isMenuOpen" xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h16" /></svg>
                    <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
                </button>
            </div>

            <div v-show="isMenuOpen" class="md:hidden absolute top-full left-0 w-full bg-white border-b border-slate-100 shadow-2xl flex flex-col py-4 px-6 gap-4">
                <a @click="isMenuOpen = false" href="#industri" class="text-sm font-black uppercase tracking-widest text-slate-600 py-2 border-b border-slate-50">Solusi Industri</a>
                <a @click="isMenuOpen = false" href="#features" class="text-sm font-black uppercase tracking-widest text-slate-600 py-2 border-b border-slate-50">Teknologi Utama</a>
                <a @click="isMenuOpen = false" href="#pricing" class="text-sm font-black uppercase tracking-widest text-slate-600 py-2 border-b border-slate-50">Daftar Harga</a>
                <router-link @click="isMenuOpen = false" to="/login" class="text-sm font-black uppercase tracking-widest text-slate-600 py-2">Portal Log In</router-link>
                <a @click="isMenuOpen = false" href="#pricing" class="w-full text-center bg-indigo-600 text-white text-xs font-black uppercase tracking-[0.2em] px-6 py-4 rounded-xl mt-2">Mulai Trial Gratis</a>
            </div>
        </nav>

        <section class="relative pt-24 md:pt-44 pb-20 md:pb-28 overflow-hidden">
            <div class="absolute -top-40 -left-40 w-[50rem] h-[50rem] bg-indigo-200/40 rounded-full blur-3xl opacity-60 pointer-events-none"></div>
            <div class="absolute top-20 -right-20 w-[40rem] h-[40rem] bg-blue-200/40 rounded-full blur-3xl opacity-60 pointer-events-none"></div>

            <div class="max-w-7xl mx-auto px-6 relative z-10 flex flex-col lg:flex-row items-center gap-16 lg:gap-10">
                
                <div class="flex-1 text-center lg:text-left">
                    <div class="inline-flex items-center gap-2 px-5 py-2.5 bg-white border border-slate-200 rounded-full mb-8 shadow-sm">
                        <span class="w-2.5 h-2.5 rounded-full bg-indigo-500 animate-pulse"></span>
                        <span class="text-[9px] md:text-[10px] font-black text-slate-600 uppercase tracking-[0.2em]">Infrastruktur SaaS Skala Enterprise</span>
                    </div>
                    
                    <h1 class="text-5xl md:text-6xl lg:text-7xl font-black text-slate-900 tracking-tighter mb-8 leading-[1.05]">
                        Orkestrasi Bisnis <br class="hidden lg:block"/> Terpadu Dalam <br/>
                        <span class="text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-blue-500 italic">Satu Ekosistem.</span>
                    </h1>
                    
                    <p class="max-w-xl mx-auto lg:mx-0 text-slate-500 font-bold text-base md:text-lg mb-10 leading-relaxed">
                        Tinggalkan cara manual. Kendalikan transaksi kasir, inventaris, hingga analitik performa cabang secara real-time dari satu dashboard operasional.
                    </p>
                    
                    <div class="flex flex-col sm:flex-row items-center justify-center lg:justify-start gap-5">
                        <a href="#pricing" class="w-full sm:w-auto bg-indigo-600 text-white px-8 md:px-10 py-4 md:py-5 rounded-2xl font-black text-xs md:text-sm uppercase tracking-widest shadow-xl shadow-indigo-200 hover:bg-slate-900 transition-all hover:-translate-y-1 active:scale-95 flex items-center justify-center gap-2">
                            Lihat Paket Harga
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M5 12h14"/><path d="m12 5 7 7-7 7"/></svg>
                        </a>
                        <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest mt-2 sm:mt-0 sm:ml-2">Proses Setup Cepat &lt; 2 Menit</span>
                    </div>
                </div>

                <div class="flex-1 w-full max-w-lg lg:max-w-none relative lg:mt-0">
                    <svg class="absolute -top-10 -left-10 w-32 h-32 text-slate-200" fill="currentColor" viewBox="0 0 100 100"><pattern id="dots" x="0" y="0" width="20" height="20" patternUnits="userSpaceOnUse"><circle cx="2" cy="2" r="2"></circle></pattern><rect width="100" height="100" fill="url(#dots)"></rect></svg>
                    
                    <div class="relative bg-white p-8 rounded-[40px] shadow-2xl shadow-indigo-900/10 border border-white z-10 transform hover:-translate-y-2 transition-transform duration-500">
                        <div class="flex items-center justify-between mb-8 pb-6 border-b border-slate-100">
                            <div class="flex items-center gap-4">
                                <div class="w-14 h-14 rounded-2xl bg-gradient-to-br from-indigo-500 to-indigo-700 flex items-center justify-center text-white shadow-lg shadow-indigo-200">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-7 h-7" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/></svg>
                                </div>
                                <div>
                                    <div class="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1">Live Revenue (All Branches)</div>
                                    <div class="text-3xl font-black text-slate-800 tracking-tighter">Rp 24.580.000</div>
                                </div>
                            </div>
                            <span class="bg-emerald-50 text-emerald-600 px-3 py-1 rounded-lg text-[10px] font-black uppercase tracking-widest flex items-center gap-1">
                                <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M5 10l7-7m0 0l7 7m-7-7v18"/></svg> 12.5%
                            </span>
                        </div>
                        
                        <div class="space-y-4">
                            <div class="flex justify-between items-center text-xs font-bold text-slate-500 mb-2"><span>Aktivitas Sistem</span><span>Status</span></div>
                            <div class="flex items-center justify-between p-3 bg-slate-50 rounded-xl">
                                <div class="flex items-center gap-3"><div class="w-2 h-2 rounded-full bg-emerald-500"></div><span class="text-sm font-black text-slate-700">Sinkronisasi Multi-Gudang</span></div>
                                <span class="text-[10px] font-bold text-slate-400 uppercase">Selesai</span>
                            </div>
                            <div class="flex items-center justify-between p-3 bg-slate-50 rounded-xl">
                                <div class="flex items-center gap-3"><div class="w-2 h-2 rounded-full bg-blue-500 animate-pulse"></div><span class="text-sm font-black text-slate-700">Backup Cloud Server</span></div>
                                <span class="text-[10px] font-bold text-slate-400 uppercase">Proses</span>
                            </div>
                        </div>
                    </div>
                </div>

            </div>
        </section>

        <section id="industri" class="py-24 bg-slate-900 text-white relative">
            <div class="max-w-7xl mx-auto px-6">
                <div class="text-center mb-16">
                    <h2 class="text-3xl md:text-5xl font-black tracking-tighter mb-5">Arsitektur Modular <br/> <span class="text-indigo-400">Sesuai Model Bisnis</span></h2>
                    <p class="text-slate-400 font-bold text-sm uppercase tracking-widest max-w-2xl mx-auto">Kami tidak menggunakan sistem "Pukul Rata". Setiap industri mendapatkan modul spesifik yang otomatis aktif saat proses pendaftaran.</p>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-3 gap-6 lg:gap-8">
                    <div v-for="ind in industries" :key="ind.title" @click="activePricingTab = ind.id" onclick="document.getElementById('pricing').scrollIntoView()" class="bg-slate-800/40 p-8 rounded-[32px] border border-slate-700 hover:border-indigo-500 hover:bg-slate-800 transition-all flex flex-col items-start gap-6 group cursor-pointer">
                        <div class="w-16 h-16 rounded-2xl bg-slate-700/50 text-indigo-400 flex items-center justify-center group-hover:bg-indigo-500 group-hover:text-white transition-all duration-300 group-hover:scale-110">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                <path :d="ind.icon" />
                            </svg>
                        </div>
                        <div>
                            <h4 class="font-black text-lg uppercase tracking-wide text-white mb-2">{{ ind.title }}</h4>
                            <p class="text-xs font-bold text-slate-400 uppercase tracking-widest">{{ ind.desc }}</p>
                        </div>
                        <div class="mt-auto pt-6 w-full flex items-center justify-between text-indigo-400 text-xs font-black uppercase tracking-widest group-hover:text-indigo-300">
                            Lihat Paket Harga <svg class="w-4 h-4 transform group-hover:translate-x-2 transition-transform" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M17 8l4 4m0 0l-4 4m4-4H3"/></svg>
                        </div>
                    </div>
                </div>
            </div>
        </section>

        <section id="features" class="py-28 bg-white relative border-t border-slate-100">
            <div class="max-w-7xl mx-auto px-6">
                <div class="text-center mb-20">
                    <h2 class="text-3xl md:text-5xl font-black text-slate-900 tracking-tighter mb-5">Teknologi <span class="text-indigo-600">Terdepan</span></h2>
                    <p class="text-slate-400 font-bold text-sm uppercase tracking-widest">Infrastruktur kuat di belakang operasional harian Anda.</p>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
                    <div v-for="feat in features" :key="feat.title" class="bg-slate-50/50 p-8 rounded-[32px] border border-slate-100 hover:shadow-2xl hover:shadow-indigo-900/5 transition-all duration-300 group hover:-translate-y-2">
                        <div class="w-14 h-14 bg-white border border-slate-200 text-indigo-600 rounded-2xl flex items-center justify-center mb-6 group-hover:bg-indigo-600 group-hover:border-indigo-600 group-hover:text-white transition-all duration-300 shadow-sm" v-html="feat.icon"></div>
                        <h3 class="text-lg font-black text-slate-800 mb-3 tracking-tight">{{ feat.title }}</h3>
                        <p class="text-slate-500 font-bold text-sm leading-relaxed">{{ feat.desc }}</p>
                    </div>
                </div>
            </div>
        </section>

        <!-- 🚀 PRICING SECTION DENGAN TAB INDUSTRI -->
        <section id="pricing" class="py-28 bg-slate-50 border-t border-slate-200">
            <div class="max-w-7xl mx-auto px-6">
                <div class="text-center mb-12">
                    <h2 class="text-3xl md:text-5xl font-black text-slate-900 tracking-tighter mb-5">Investasi <span class="text-indigo-600">Transparan</span></h2>
                    <p class="text-slate-400 font-bold text-sm uppercase tracking-widest">Pilih industri Anda, dan temukan plan yang dirancang khusus untuk skala Anda.</p>
                </div>

                <!-- Tab Pemilihan Industri -->
                <div class="w-full flex justify-center mb-16 px-4">
                    <div class="bg-white p-1.5 rounded-[20px] border border-slate-200 flex w-full max-w-xl shadow-sm">
                        <button v-for="ind in industries" :key="ind.id" 
                                @click="activePricingTab = ind.id" 
                                :class="activePricingTab === ind.id ? 'bg-slate-900 text-white shadow-md' : 'text-slate-500 hover:text-slate-800 hover:bg-slate-50'" 
                                class="flex-1 px-2 py-3 rounded-xl font-black text-[9px] uppercase tracking-widest transition-all text-center">
                            {{ ind.title }}
                        </button>
                    </div>
                </div>

                <!-- Grid Kartu Harga Dinamis -->
                <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                    <div v-for="plan in pricingPlans[activePricingTab]" :key="plan.id" 
                         class="bg-white p-8 rounded-[32px] border-2 transition-all flex flex-col relative"
                         :class="plan.id === 'pro' ? 'border-indigo-500 shadow-2xl shadow-indigo-100 lg:scale-105 z-10' : 'border-slate-100 hover:border-slate-300'">
                        
                        <div v-if="plan.id === 'pro'" class="absolute -top-4 left-1/2 -translate-x-1/2 bg-gradient-to-r from-indigo-500 to-indigo-600 text-white px-5 py-1.5 rounded-full text-[10px] font-black uppercase tracking-widest shadow-lg">Rekomendasi Skala UMKM</div>

                        <div class="mb-6 mt-4">
                            <h3 class="font-black text-sm uppercase tracking-widest text-slate-400 mb-4">{{ plan.name }}</h3>
                            <div class="flex items-end gap-1">
                                <span class="text-4xl lg:text-5xl font-black text-slate-900 tracking-tighter">{{ plan.price }}</span>
                                <span class="text-slate-400 font-bold text-[10px] uppercase tracking-widest mb-1.5">{{ plan.duration }}</span>
                            </div>
                            <p class="text-slate-500 text-xs font-bold mt-5 h-12 leading-relaxed">{{ plan.desc }}</p>
                        </div>

                        <ul class="space-y-4 mb-10 flex-1 border-t border-slate-100 pt-6">
                            <li v-for="feat in plan.features" :key="feat" class="flex items-start gap-3 text-xs font-bold text-slate-700 leading-tight">
                                <svg class="w-4 h-4 text-emerald-500 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/></svg> 
                                {{ feat }}
                            </li>
                        </ul>

                        <!-- TOMBOL DAFTAR (Kirim Parameter URL ke halaman Register) -->
                        <router-link :to="`/register?industry=${activePricingTab}&plan=${plan.id}`" 
                                     :class="plan.id === 'pro' ? 'bg-indigo-600 text-white hover:bg-slate-900 shadow-lg shadow-indigo-200' : 'bg-slate-50 text-slate-700 border border-slate-200 hover:bg-slate-100 hover:border-slate-300'"
                                     class="block w-full text-center py-4 rounded-[16px] font-black text-[10px] uppercase tracking-widest transition-all active:scale-95">
                            {{ plan.id === 'trial' ? 'Mulai Eksplorasi' : 'Pilih ' + plan.name }}
                        </router-link>
                    </div>
                </div>
            </div>
        </section>

        <section class="py-24 bg-[#F8FAFC]">
            <div class="max-w-5xl mx-auto px-6">
                <div class="bg-gradient-to-br from-indigo-900 via-slate-900 to-black rounded-[40px] p-12 md:p-20 text-center relative overflow-hidden shadow-2xl">
                    <div class="absolute top-0 right-0 w-80 h-80 bg-indigo-500/20 rounded-full blur-3xl pointer-events-none"></div>
                    <div class="absolute bottom-0 left-0 w-80 h-80 bg-blue-500/20 rounded-full blur-3xl pointer-events-none"></div>
                    
                    <div class="relative z-10">
                        <h2 class="text-4xl md:text-6xl font-black text-white tracking-tighter mb-6 leading-tight">
                            Siap Elevasi Bisnis <br/> Hari Ini?
                        </h2>
                        
                        <a href="#pricing" class="inline-flex items-center gap-3 bg-white text-slate-900 px-8 md:px-12 py-5 rounded-full font-black text-xs md:text-sm uppercase tracking-widest hover:bg-indigo-50 transition-all active:scale-95 shadow-xl mt-4">
                            Mulai Perjalanan Anda
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M5 12h14"/><path d="m12 5 7 7-7 7"/></svg>
                        </a>
                        
                        <div class="mt-10 flex flex-col sm:flex-row items-center justify-center gap-3 sm:gap-6 text-indigo-200 text-[10px] font-black uppercase tracking-widest">
                            <span class="flex items-center gap-1.5"><svg class="w-4 h-4 text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path d="M20 6L9 17l-5-5"/></svg> Gratis Registrasi</span>
                            <span class="hidden sm:inline opacity-30">•</span>
                            <span class="flex items-center gap-1.5"><svg class="w-4 h-4 text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path d="M20 6L9 17l-5-5"/></svg> Tanpa Kewajiban Kartu Kredit</span>
                        </div>
                    </div>
                </div>
            </div>
        </section>

        <footer class="py-12 bg-white border-t border-slate-100">
            <div class="max-w-7xl mx-auto px-6 flex flex-col md:flex-row items-center justify-between gap-6 text-center md:text-left">
                <div>
                    <div class="font-black text-2xl tracking-tighter text-slate-800">NEXA<span class="text-indigo-600">POS</span></div>
                    <p class="text-[9px] font-black text-slate-400 uppercase tracking-[0.2em] mt-1.5">Enterprise Operations System</p>
                </div>
                
                <div class="text-[9px] font-black text-slate-400 uppercase tracking-[0.2em]">
                    &copy; 2026 Developed by Arif Juliawan
                </div>

                <div class="flex gap-8 justify-center">
                    <a href="#" class="text-[10px] font-black text-slate-400 hover:text-indigo-600 uppercase tracking-widest transition-colors">Term of Service</a>
                    <a href="#" class="text-[10px] font-black text-slate-400 hover:text-indigo-600 uppercase tracking-widest transition-colors">Privacy Policy</a>
                </div>
            </div>
        </footer>

    </div>
</template>

<style scoped>
html { scroll-behavior: smooth; }
</style>