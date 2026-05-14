<script setup>
import { useRouter } from 'vue-router';

const router = useRouter();

const features = [
  {
    title: "Mesin Kasir Cloud",
    desc: "Transaksi lancar tanpa hambatan, stok otomatis terpotong secara realtime.",
    icon: `<svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="8" cy="21" r="1"/><circle cx="19" cy="21" r="1"/><path d="M2.05 2.05h2l2.66 12.42a2 2 0 0 0 2 1.58h9.78a2 2 0 0 0 1.95-1.57l1.56-7.43a1 1 0 0 0-1-1.21H5.84"/></svg>`
  },
  {
    title: "Sistem Absensi", // 🚀 Dibuat general sesuai request
    desc: "Berbagai metode absensi digital yang fleksibel, mulai dari foto selfie hingga verifikasi lokasi.",
    icon: `<svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><rect width="18" height="18" x="3" y="4" rx="2"/><path d="M3 10h18"/></svg>`
  },
  {
    title: "Manajemen Inventori",
    desc: "Pantau stok masuk dan keluar secara akurat. Hindari selisih barang dengan fitur Stock Opname.",
    icon: `<svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m7.5 4.27 9 5.15"/><path d="M21 8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16Z"/><path d="m3.3 7 8.7 5 8.7-5"/><path d="M12 22V12"/></svg>`
  }
];
</script>

<template>
  <div class="min-h-screen bg-white font-sans text-slate-900 selection:bg-blue-100 selection:text-blue-600">
    
    <nav class="sticky top-0 z-50 bg-white/80 backdrop-blur-xl border-b border-slate-100 px-6 py-4">
      <div class="max-w-7xl mx-auto flex items-center justify-between">
        <div class="flex flex-col">
          <div class="font-black text-2xl tracking-tighter">POS<span class="text-blue-600">UMKM</span></div>
        </div>
        <div class="flex items-center gap-4">
          <router-link to="/login" class="text-xs font-black uppercase tracking-widest text-slate-500 hover:text-blue-600 transition-all px-4 py-2">Masuk</router-link>
          <router-link to="/register" class="bg-slate-900 text-white text-[10px] font-black uppercase tracking-[0.2em] px-6 py-3 rounded-2xl hover:bg-blue-600 shadow-xl shadow-slate-200 transition-all active:scale-95">Daftar Sekarang</router-link>
        </div>
      </div>
    </nav>

    <section class="relative pt-20 pb-32 overflow-hidden">
      <div class="absolute -top-24 -left-24 w-96 h-96 bg-blue-100/50 rounded-full blur-3xl opacity-50"></div>
      <div class="absolute top-1/2 -right-24 w-80 h-80 bg-indigo-100/50 rounded-full blur-3xl opacity-50"></div>

      <div class="max-w-7xl mx-auto px-6 text-center relative z-10">
        <div class="inline-block px-4 py-1.5 bg-blue-50 rounded-full mb-6">
          <span class="text-[10px] font-black text-blue-600 uppercase tracking-[0.3em]">🚀 Solusi Retail Masa Depan</span>
        </div>
        <h1 class="text-5xl md:text-7xl font-black text-slate-900 tracking-tighter mb-8 leading-[0.9]">
          Kelola Toko Jadi <br/> <span class="text-blue-600 italic">Lebih Gampang.</span>
        </h1>
        <p class="max-w-2xl mx-auto text-slate-400 font-bold text-lg mb-12">
          Satu aplikasi untuk kasir, absensi wajah, hingga laporan keuangan otomatis. Dirancang khusus untuk memajukan UMKM Indonesia.
        </p>
        <div class="flex flex-col sm:flex-row items-center justify-center gap-4">
          <router-link to="/register" class="w-full sm:w-auto bg-blue-600 text-white px-10 py-5 rounded-[24px] font-black uppercase tracking-widest shadow-2xl shadow-blue-200 hover:bg-slate-900 transition-all hover:-translate-y-1">Mulai Gratis Sekarang</router-link>
          <a href="#features" class="w-full sm:w-auto px-10 py-5 rounded-[24px] font-black text-slate-500 uppercase tracking-widest border-2 border-slate-100 hover:bg-slate-50 transition-all">Lihat Fitur</a>
        </div>
      </div>
    </section>

    <section id="features" class="py-24 bg-slate-50 relative">
      <div class="max-w-7xl mx-auto px-6">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div v-for="feat in features" :key="feat.title" class="bg-white p-10 rounded-[40px] border border-white shadow-sm hover:shadow-xl transition-all group">
            <div class="w-16 h-16 bg-blue-50 text-blue-600 rounded-3xl flex items-center justify-center mb-8 group-hover:bg-blue-600 group-hover:text-white transition-all duration-500" v-html="feat.icon"></div>
            <h3 class="text-xl font-black text-slate-900 mb-4 uppercase tracking-tight">{{ feat.title }}</h3>
            <p class="text-slate-400 font-bold text-sm leading-relaxed">{{ feat.desc }}</p>
          </div>
        </div>
      </div>
    </section>

    <section class="py-24">
      <div class="max-w-5xl mx-auto px-6">
        <div class="bg-slate-900 rounded-[50px] p-12 md:p-20 text-center relative overflow-hidden">
          <div class="absolute top-0 right-0 w-64 h-64 bg-blue-600/20 rounded-full blur-3xl"></div>
          <div class="relative z-10">
            <h2 class="text-3xl md:text-5xl font-black text-white tracking-tighter mb-8 leading-tight">
              Siap Bawa Tokomu <br/> Naik Kelas Hari Ini?
            </h2>
            <router-link to="/register" class="inline-block bg-white text-slate-900 px-10 py-5 rounded-[24px] font-black uppercase tracking-widest hover:bg-blue-50 transition-all">Buka Toko Gratis</router-link>
            <p class="mt-8 text-slate-500 text-[10px] font-black uppercase tracking-widest">Tanpa kartu kredit • Setup hanya 2 menit</p>
          </div>
        </div>
      </div>
    </section>

    <footer class="py-12 border-t border-slate-100">
      <div class="max-w-7xl mx-auto px-6 flex flex-col md:flex-row items-center justify-between gap-6 text-center md:text-left">
        <div>
          <div class="font-black text-xl tracking-tighter">POS<span class="text-blue-600">UMKM</span></div>
          <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mt-1 italic">&copy; 2026 Developed by Arif Juliawan</p>
        </div>
        <div class="flex gap-8">
          <a href="#" class="text-[10px] font-black text-slate-400 hover:text-blue-600 uppercase tracking-widest">Syarat & Ketentuan</a>
          <a href="#" class="text-[10px] font-black text-slate-400 hover:text-blue-600 uppercase tracking-widest">Kebijakan Privasi</a>
        </div>
      </div>
    </footer>

  </div>
</template>

<style scoped>
/* Smooth Scroll */
html {
  scroll-behavior: smooth;
}
</style>