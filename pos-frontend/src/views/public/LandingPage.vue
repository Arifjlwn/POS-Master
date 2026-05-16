<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();
const isMenuOpen = ref(false);
const isScrolled = ref(false);

// Menangani efek blur navigasi saat di-scroll
onMounted(() => {
    window.addEventListener('scroll', () => {
        isScrolled.value = window.scrollY > 20;
    });
});

const features = [
    {
        title: "Cloud POS Multi-Industri",
        desc: "Kasir pintar yang menyesuaikan bisnis Anda. Dari scan barcode minimarket, manajemen meja resto, hingga hitung berat laundry.",
        icon: `<svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="16" height="20" x="4" y="2" rx="2" ry="2"/><line x1="12" y1="18" x2="12.01" y2="18"/></svg>`
    },
    {
        title: "Sistem Presensi Adaptif",
        desc: "Metode absensi yang fleksibel untuk tim Anda. Pilih menggunakan Face AI, PIN, atau integrasi mesin absen sesuai kebutuhan.",
        icon: `<svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M22 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>`
    },
    {
        title: "Smart Inventory",
        desc: "Stok otomatis terpotong saat transaksi. Dilengkapi fitur resep bahan baku (F&B) dan manajemen multi-satuan barang.",
        icon: `<svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m7.5 4.27 9 5.15"/><path d="M21 8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16Z"/><path d="m3.3 7 8.7 5 8.7-5"/><path d="M12 22V12"/></svg>`
    },
    {
        title: "Jadwal & Shift Dinamis",
        desc: "Atur rotasi shift karyawan dengan mudah melalui matriks visual. Terintegrasi langsung dengan sistem penggajian.",
        icon: `<svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>`
    }
];

const industries = [
    { title: "Retail & Kelontong", desc: "Minimarket, Butik, Elektronik", icon: "M16 11V7a4 4 0 0 0-8 0v4M5 9h14l1 12H4L5 9z" },
    { title: "Food & Beverage", desc: "Cafe, Restoran, Foodcourt", icon: "M18 8h1a4 4 0 0 1 0 8h-1M2 8h16v9a4 4 0 0 1-4 4H6a4 4 0 0 1-4-4V8z M6 1v3 M10 1v3 M14 1v3" },
    { title: "Layanan & Jasa", desc: "Laundry, Barbershop, Bengkel", icon: "M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z M3.27 6.96L12 12.01l8.73-5.05" }
];
</script>

<template>
    <div class="min-h-screen bg-[#F8FAFC] font-sans text-slate-900 selection:bg-indigo-100 selection:text-indigo-600 overflow-x-hidden">
        
        <nav :class="['fixed top-0 w-full z-50 transition-all duration-300', isScrolled ? 'bg-white/80 backdrop-blur-xl border-b border-slate-200 py-3 shadow-sm' : 'bg-transparent py-5']">
            <div class="max-w-7xl mx-auto px-6 flex items-center justify-between">
                <div class="flex flex-col">
                    <div class="font-black text-2xl md:text-3xl tracking-tighter">POS<span class="text-indigo-600">UMKM</span></div>
                </div>
                
                <div class="hidden md:flex items-center gap-6">
                    <a href="#industri" class="text-xs font-black uppercase tracking-widest text-slate-500 hover:text-indigo-600 transition-colors">Tipe Bisnis</a>
                    <a href="#features" class="text-xs font-black uppercase tracking-widest text-slate-500 hover:text-indigo-600 transition-colors">Fitur</a>
                    <a href="#pricing" class="text-xs font-black uppercase tracking-widest text-slate-500 hover:text-indigo-600 transition-colors">Harga</a>
                    <router-link to="/login" class="text-xs font-black uppercase tracking-widest text-slate-800 hover:text-indigo-600 transition-colors ml-4">Masuk</router-link>
                    <router-link to="/register" class="bg-indigo-600 text-white text-[10px] font-black uppercase tracking-[0.2em] px-6 py-3.5 rounded-2xl hover:bg-slate-900 shadow-xl shadow-indigo-200 transition-all active:scale-95">Mulai Trial Gratis</router-link>
                </div>

                <button @click="isMenuOpen = !isMenuOpen" class="md:hidden p-2 text-slate-600 hover:text-indigo-600 transition-colors">
                    <svg v-if="!isMenuOpen" xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h16" /></svg>
                    <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
                </button>
            </div>

            <div v-show="isMenuOpen" class="md:hidden absolute top-full left-0 w-full bg-white border-b border-slate-100 shadow-xl flex flex-col py-4 px-6 gap-4 animate-[fadeIn_0.2s_ease-out]">
                <a @click="isMenuOpen = false" href="#industri" class="text-sm font-black uppercase tracking-widest text-slate-600 py-2">Tipe Bisnis</a>
                <a @click="isMenuOpen = false" href="#features" class="text-sm font-black uppercase tracking-widest text-slate-600 py-2">Fitur Utama</a>
                <a @click="isMenuOpen = false" href="#pricing" class="text-sm font-black uppercase tracking-widest text-slate-600 py-2">Harga</a>
                <router-link @click="isMenuOpen = false" to="/login" class="text-sm font-black uppercase tracking-widest text-slate-600 py-2">Masuk Log</router-link>
                <router-link @click="isMenuOpen = false" to="/register" class="w-full text-center bg-indigo-600 text-white text-xs font-black uppercase tracking-[0.2em] px-6 py-4 rounded-xl mt-2">Mulai Trial Gratis</router-link>
            </div>
        </nav>

        <section class="relative pt-32 md:pt-40 pb-20 md:pb-24 overflow-hidden">
            <div class="absolute -top-40 -left-40 w-[40rem] h-[40rem] bg-indigo-200/40 rounded-full blur-3xl opacity-60 pointer-events-none"></div>
            <div class="absolute top-20 -right-20 w-[30rem] h-[30rem] bg-blue-200/40 rounded-full blur-3xl opacity-60 pointer-events-none"></div>

            <div class="max-w-7xl mx-auto px-6 relative z-10 flex flex-col lg:flex-row items-center gap-12 lg:gap-8">
                
                <div class="flex-1 text-center lg:text-left">
                    <div class="inline-flex items-center gap-2 px-4 py-2 bg-white border border-slate-200 rounded-full mb-6 shadow-sm">
                        <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
                        <span class="text-[9px] md:text-[10px] font-black text-slate-600 uppercase tracking-[0.2em]">100% Cocok Untuk Retail, F&B, dan Jasa</span>
                    </div>
                    
                    <h1 class="text-5xl md:text-6xl lg:text-7xl font-black text-slate-900 tracking-tighter mb-6 leading-[1.05]">
                        Sistem Pintar <br class="hidden lg:block"/> Untuk Bisnis <br/>
                        <span class="text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-blue-500 italic">Makin Bersinar.</span>
                    </h1>
                    
                    <p class="max-w-xl mx-auto lg:mx-0 text-slate-500 font-bold text-base md:text-lg mb-10 leading-relaxed">
                        Apapun jenis bisnis Anda, kendalikan kasir, kelola stok bahan/barang, dan pantau kehadiran tim dari satu dashboard eksklusif.
                    </p>
                    
                    <div class="flex flex-col sm:flex-row items-center justify-center lg:justify-start gap-4">
                        <router-link to="/register" class="w-full sm:w-auto bg-indigo-600 text-white px-8 md:px-10 py-4 md:py-5 rounded-[20px] font-black text-xs md:text-sm uppercase tracking-widest shadow-xl shadow-indigo-200 hover:bg-slate-900 transition-all hover:-translate-y-1 active:scale-95 flex items-center justify-center gap-2">
                            Coba Gratis 14 Hari
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="5" y1="12" x2="19" y2="12"/><polyline points="12 5 19 12 12 19"/></svg>
                        </router-link>
                        <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest mt-2 sm:mt-0 sm:ml-2">Tanpa Kartu Kredit</span>
                    </div>
                </div>

                <div class="flex-1 w-full max-w-lg lg:max-w-none relative animate-[float_6s_ease-in-out_infinite] mt-10 lg:mt-0">
                    <svg class="absolute -top-10 -left-10 w-24 h-24 text-slate-200" fill="currentColor" viewBox="0 0 100 100"><pattern id="dots" x="0" y="0" width="20" height="20" patternUnits="userSpaceOnUse"><circle cx="2" cy="2" r="2"></circle></pattern><rect width="100" height="100" fill="url(#dots)"></rect></svg>
                    
                    <div class="relative bg-white p-6 rounded-[32px] md:rounded-[40px] shadow-2xl shadow-slate-200/50 border-[8px] border-slate-50 z-10 rotate-2 hover:rotate-0 transition-transform duration-500">
                        <div class="flex items-center gap-3 mb-6">
                            <div class="w-12 h-12 rounded-xl bg-gradient-to-br from-emerald-400 to-emerald-600 flex items-center justify-center text-white shadow-md">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/></svg>
                            </div>
                            <div>
                                <div class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Pendapatan Hari Ini</div>
                                <div class="text-2xl font-black text-slate-800 tracking-tighter">Rp 4.580.000</div>
                            </div>
                        </div>
                        <div class="space-y-3">
                            <div class="h-4 bg-slate-100 rounded-full w-full"></div>
                            <div class="h-4 bg-slate-100 rounded-full w-5/6"></div>
                            <div class="h-4 bg-slate-100 rounded-full w-4/6"></div>
                        </div>
                        <div class="mt-6 pt-6 border-t border-slate-100 flex justify-between items-center">
                            <div class="flex -space-x-2">
                                <div class="w-8 h-8 rounded-full bg-blue-100 border-2 border-white flex items-center justify-center text-[10px] font-black">AJ</div>
                                <div class="w-8 h-8 rounded-full bg-amber-100 border-2 border-white flex items-center justify-center text-[10px] font-black">BD</div>
                            </div>
                            <span class="text-[9px] font-black bg-indigo-50 text-indigo-600 px-3 py-1.5 rounded-lg uppercase tracking-widest border border-indigo-100">Sistem Aktif</span>
                        </div>
                    </div>
                </div>

            </div>
        </section>

        <section id="industri" class="py-20 bg-slate-900 text-white relative">
            <div class="max-w-7xl mx-auto px-6">
                <div class="text-center mb-12">
                    <h2 class="text-2xl md:text-3xl font-black tracking-tighter mb-4">Satu Sistem, <span class="text-blue-400">Beragam Industri</span></h2>
                    <p class="text-slate-400 font-bold text-xs uppercase tracking-widest max-w-2xl mx-auto">Kami mengerti setiap bisnis punya alur berbeda. Aplikasi akan menyesuaikan fitur secara otomatis saat proses setup.</p>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
                    <div v-for="ind in industries" :key="ind.title" class="bg-slate-800/50 p-6 rounded-[24px] border border-slate-700 hover:border-indigo-500 transition-all flex items-center gap-4 group cursor-default">
                        <div class="w-12 h-12 rounded-xl bg-slate-700 text-indigo-400 flex items-center justify-center group-hover:bg-indigo-600 group-hover:text-white transition-colors">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                <path :d="ind.icon" />
                            </svg>
                        </div>
                        <div>
                            <h4 class="font-black text-sm uppercase tracking-wider text-slate-100">{{ ind.title }}</h4>
                            <p class="text-[10px] font-bold text-slate-400 mt-1 uppercase tracking-widest">{{ ind.desc }}</p>
                        </div>
                    </div>
                </div>
            </div>
        </section>

        <section id="features" class="py-24 bg-white relative border-t border-slate-100">
            <div class="max-w-7xl mx-auto px-6">
                <div class="text-center mb-16">
                    <h2 class="text-3xl md:text-5xl font-black text-slate-900 tracking-tighter mb-4">Fitur Kelas <span class="text-indigo-600">Enterprise</span></h2>
                    <p class="text-slate-400 font-bold text-sm uppercase tracking-widest">Tanpa Ribet. Semua Yang Anda Butuhkan Ada Di Sini.</p>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                    <div v-for="feat in features" :key="feat.title" class="bg-slate-50 p-8 rounded-[32px] border border-slate-100 shadow-sm hover:shadow-2xl hover:shadow-indigo-100 transition-all duration-300 group hover:-translate-y-2">
                        <div class="w-14 h-14 bg-white border-2 border-slate-100 text-indigo-600 rounded-2xl flex items-center justify-center mb-6 group-hover:bg-indigo-600 group-hover:border-indigo-600 group-hover:text-white transition-all duration-300 shadow-sm group-hover:rotate-6" v-html="feat.icon"></div>
                        <h3 class="text-lg font-black text-slate-800 mb-3 tracking-tight">{{ feat.title }}</h3>
                        <p class="text-slate-500 font-medium text-sm leading-relaxed">{{ feat.desc }}</p>
                    </div>
                </div>
            </div>
        </section>

        <section id="pricing" class="py-24 bg-slate-50 border-t border-slate-200">
            <div class="max-w-7xl mx-auto px-6">
                <div class="text-center mb-16">
                    <h2 class="text-3xl md:text-5xl font-black text-slate-900 tracking-tighter mb-4">Harga <span class="text-indigo-600">Transparan</span></h2>
                    <p class="text-slate-400 font-bold text-sm uppercase tracking-widest">Mulai dengan gratis, bayar saat bisnis Anda sudah siap bertumbuh.</p>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-2 gap-8 max-w-4xl mx-auto">
                    <div class="bg-white p-8 md:p-10 rounded-[40px] shadow-sm border-2 border-slate-100 hover:border-indigo-200 transition-all">
                        <div class="inline-block px-4 py-1.5 bg-slate-100 text-slate-600 rounded-full mb-6 font-black text-[10px] uppercase tracking-widest">Paket Starter</div>
                        <div class="mb-6">
                            <span class="text-5xl font-black text-slate-900 tracking-tighter">Rp 0</span>
                            <span class="text-slate-400 font-bold text-sm uppercase tracking-widest"> / 14 Hari</span>
                        </div>
                        <p class="text-slate-500 text-sm font-medium mb-8">Eksplorasi seluruh fitur tanpa batasan dan tanpa komitmen. Cocok untuk mencoba kecocokan sistem.</p>
                        <ul class="space-y-4 mb-10">
                            <li class="flex items-center gap-3 text-sm font-bold text-slate-700"><svg class="w-5 h-5 text-emerald-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/></svg> Full Akses Modul POS & Kasir</li>
                            <li class="flex items-center gap-3 text-sm font-bold text-slate-700"><svg class="w-5 h-5 text-emerald-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/></svg> Unlimited Transaksi</li>
                            <li class="flex items-center gap-3 text-sm font-bold text-slate-700"><svg class="w-5 h-5 text-emerald-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/></svg> Tanpa Kartu Kredit</li>
                        </ul>
                        <router-link to="/register" class="block w-full text-center bg-slate-100 text-slate-600 hover:bg-slate-200 py-4 rounded-[20px] font-black text-xs uppercase tracking-widest transition-colors">Coba Gratis</router-link>
                    </div>

                    <div class="bg-slate-900 p-8 md:p-10 rounded-[40px] shadow-2xl shadow-indigo-900/50 border-4 border-indigo-500 relative overflow-hidden transform md:-translate-y-4">
                        <div class="absolute top-6 right-6 bg-amber-400 text-amber-900 px-3 py-1 rounded-full text-[9px] font-black uppercase tracking-widest shadow-sm">Paling Laris</div>
                        <div class="inline-block px-4 py-1.5 bg-indigo-500/20 text-indigo-300 rounded-full mb-6 font-black text-[10px] uppercase tracking-widest border border-indigo-500/30">Paket Pro Bisnis</div>
                        <div class="mb-6">
                            <span class="text-5xl font-black text-white tracking-tighter">149k</span>
                            <span class="text-slate-400 font-bold text-sm uppercase tracking-widest"> / Bulan</span>
                        </div>
                        <p class="text-slate-300 text-sm font-medium mb-8">Solusi lengkap untuk operasional bisnis UMKM sehari-hari dengan dukungan prioritas.</p>
                        <ul class="space-y-4 mb-10">
                            <li class="flex items-center gap-3 text-sm font-bold text-white"><svg class="w-5 h-5 text-indigo-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/></svg> Semua Fitur Starter</li>
                            <li class="flex items-center gap-3 text-sm font-bold text-white"><svg class="w-5 h-5 text-indigo-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/></svg> Modul Absensi Face AI & Shift</li>
                            <li class="flex items-center gap-3 text-sm font-bold text-white"><svg class="w-5 h-5 text-indigo-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/></svg> Laporan Keuangan Export Excel</li>
                        </ul>
                        <router-link to="/register" class="block w-full text-center bg-indigo-600 text-white hover:bg-indigo-500 py-4 rounded-[20px] font-black text-xs uppercase tracking-widest transition-colors shadow-lg">Langganan Sekarang</router-link>
                    </div>
                </div>
            </div>
        </section>

        <section class="py-24 bg-[#F8FAFC]">
            <div class="max-w-5xl mx-auto px-6">
                <div class="bg-gradient-to-br from-indigo-900 to-blue-950 rounded-[40px] p-10 md:p-16 text-center relative overflow-hidden shadow-2xl">
                    <div class="absolute top-0 right-0 w-64 h-64 bg-indigo-500/20 rounded-full blur-3xl pointer-events-none"></div>
                    <div class="absolute bottom-0 left-0 w-64 h-64 bg-blue-500/20 rounded-full blur-3xl pointer-events-none"></div>
                    
                    <div class="relative z-10">
                        <h2 class="text-3xl md:text-5xl font-black text-white tracking-tighter mb-6 leading-tight">
                            Siap Transformasi Bisnis <br/> Hari Ini?
                        </h2>
                        
                        <router-link to="/register" class="inline-flex items-center gap-3 bg-white text-slate-900 px-8 md:px-12 py-4 md:py-5 rounded-[24px] font-black text-xs md:text-sm uppercase tracking-widest hover:bg-indigo-50 transition-all active:scale-95 shadow-xl mt-4">
                            Mulai Trial 14 Hari
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M5 12h14"/><path d="m12 5 7 7-7 7"/></svg>
                        </router-link>
                        
                        <div class="mt-8 flex flex-col sm:flex-row items-center justify-center gap-2 sm:gap-4 text-indigo-300 text-[10px] font-black uppercase tracking-widest">
                            <span class="flex items-center gap-1"><svg class="w-3 h-3 text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path d="M20 6L9 17l-5-5"/></svg> Gratis Registrasi</span>
                            <span class="hidden sm:inline">•</span>
                            <span class="flex items-center gap-1"><svg class="w-3 h-3 text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path d="M20 6L9 17l-5-5"/></svg> Tanpa Komitmen</span>
                        </div>
                    </div>
                </div>
            </div>
        </section>

        <footer class="py-10 bg-white border-t border-slate-100">
            <div class="max-w-7xl mx-auto px-6 flex flex-col md:flex-row items-center justify-between gap-6 text-center md:text-left">
                <div>
                    <div class="font-black text-xl tracking-tighter text-slate-800">POS<span class="text-indigo-600">UMKM</span></div>
                    <p class="text-[9px] font-black text-slate-400 uppercase tracking-[0.2em] mt-1">Sistem Cerdas Bisnis Modern</p>
                </div>
                
                <div class="text-[9px] font-black text-slate-300 uppercase tracking-[0.2em]">
                    &copy; 2026 Crafted by Arif Juliawan
                </div>

                <div class="flex gap-6 justify-center">
                    <a href="#" class="text-[10px] font-black text-slate-400 hover:text-indigo-600 uppercase tracking-widest transition-colors">Panduan</a>
                    <a href="#" class="text-[10px] font-black text-slate-400 hover:text-indigo-600 uppercase tracking-widest transition-colors">Privasi</a>
                </div>
            </div>
        </footer>

    </div>
</template>

<style scoped>
/* Smooth Scroll Default Behavior */
html {
    scroll-behavior: smooth;
}

/* Animasi Floating Card Mulus */
@keyframes float {
    0% { transform: translateY(0px); }
    50% { transform: translateY(-15px); }
    100% { transform: translateY(0px); }
}

/* Animasi Fade In Dropdown Mobile */
@keyframes fadeIn {
    from { opacity: 0; transform: translateY(-10px); }
    to { opacity: 1; transform: translateY(0); }
}
</style>