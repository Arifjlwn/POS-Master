<script setup>
import { ref } from 'vue';
import SidebarFnB from './SidebarFnB.vue';

// Dummy Data Laporan
const stats = ref({ omset: 4500000, total_order: 124, avg_transaction: 36200 });
const bestSellers = ref([
    { rank: 1, nama: 'Nasi Goreng Arzu Spesial', terjual: 45, omset: 1125000 },
    { rank: 2, nama: 'Kopi Susu Gula Aren Dewa', terjual: 38, omset: 684000 },
    { rank: 3, nama: 'Es Teh Manis Jumbo', terjual: 30, omset: 180000 }
]);

const formatRupiah = (angka) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(angka || 0);
</script>

<template>
    <SidebarFnB>
        <div class="flex-1 flex flex-col h-full bg-slate-50 overflow-hidden relative font-sans">
            <div class="p-5 md:p-8 shrink-0 bg-white border-b border-slate-200 flex flex-col md:flex-row justify-between items-start md:items-center gap-4 z-10 shadow-sm">
                <div class="flex items-center gap-4">
                    <div class="w-12 h-12 bg-indigo-50 border border-indigo-100 rounded-2xl flex items-center justify-center shrink-0">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-indigo-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"/><line x1="3" y1="9" x2="21" y2="9"/><line x1="9" y1="21" x2="9" y2="9"/></svg>
                    </div>
                    <div>
                        <h1 class="text-xl md:text-2xl font-black tracking-tighter uppercase text-slate-800">Analitik Resto</h1>
                        <p class="text-[10px] font-black text-slate-400 mt-1 uppercase tracking-widest">Performa Bisnis F&B</p>
                    </div>
                </div>
                <button class="bg-white border-2 border-slate-200 hover:border-indigo-500 hover:text-indigo-600 text-slate-600 px-6 py-3 rounded-xl font-black text-[10px] uppercase tracking-widest transition-all shadow-sm">
                    Export PDF
                </button>
            </div>

            <div class="flex-1 overflow-y-auto custom-scrollbar p-5 md:p-8 space-y-6 md:space-y-8">
                <div class="grid grid-cols-1 md:grid-cols-3 gap-5 lg:gap-6">
                    <div class="bg-gradient-to-br from-slate-900 to-indigo-900 p-6 md:p-8 rounded-3xl text-white shadow-xl shadow-slate-200">
                        <h3 class="text-[10px] font-black uppercase tracking-widest text-slate-400 mb-2">Total Omset Hari Ini</h3>
                        <p class="text-3xl md:text-4xl font-black tracking-tighter">{{ formatRupiah(stats.omset) }}</p>
                    </div>
                    <div class="bg-white border border-slate-200 p-6 md:p-8 rounded-3xl shadow-sm">
                        <h3 class="text-[10px] font-black uppercase tracking-widest text-slate-400 mb-2">Total Order Berhasil</h3>
                        <p class="text-3xl md:text-4xl font-black tracking-tighter text-slate-800">{{ stats.total_order }} <span class="text-sm text-slate-400">Tiket</span></p>
                    </div>
                    <div class="bg-white border border-slate-200 p-6 md:p-8 rounded-3xl shadow-sm">
                        <h3 class="text-[10px] font-black uppercase tracking-widest text-slate-400 mb-2">Rata-rata Penjualan / Meja</h3>
                        <p class="text-3xl md:text-4xl font-black tracking-tighter text-slate-800">{{ formatRupiah(stats.avg_transaction) }}</p>
                    </div>
                </div>

                <div class="grid grid-cols-1 xl:grid-cols-2 gap-6 md:gap-8">
                    <div class="bg-white border border-slate-200 rounded-3xl shadow-sm p-6 md:p-8">
                        <h2 class="text-xs font-black text-slate-800 uppercase tracking-widest mb-6 flex items-center gap-2">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-amber-500" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/></svg>
                            Menu Paling Laris (Top 3)
                        </h2>
                        <div class="space-y-4">
                            <div v-for="item in bestSellers" :key="item.rank" class="flex items-center justify-between p-4 bg-slate-50 rounded-2xl border border-slate-100">
                                <div class="flex items-center gap-4">
                                    <div class="w-10 h-10 rounded-xl flex items-center justify-center font-black text-lg" :class="item.rank === 1 ? 'bg-amber-100 text-amber-600' : 'bg-slate-200 text-slate-600'">#{{ item.rank }}</div>
                                    <div>
                                        <p class="font-black text-xs uppercase text-slate-800">{{ item.nama }}</p>
                                        <p class="text-[10px] font-bold text-slate-500 mt-0.5">Terjual: {{ item.terjual }} Porsi</p>
                                    </div>
                                </div>
                                <p class="font-black text-sm text-indigo-600 hidden sm:block">{{ formatRupiah(item.omset) }}</p>
                            </div>
                        </div>
                    </div>

                    <div class="bg-white border border-slate-200 rounded-3xl shadow-sm p-6 md:p-8 flex flex-col justify-center">
                        <h2 class="text-xs font-black text-slate-800 uppercase tracking-widest mb-6 text-center">Rasio Tipe Pesanan</h2>
                        <div class="flex justify-center gap-8 md:gap-16">
                            <div class="text-center">
                                <div class="w-24 h-24 rounded-full border-8 border-indigo-500 flex items-center justify-center mb-3 text-2xl font-black text-indigo-600 shadow-inner">65%</div>
                                <p class="text-[10px] font-black uppercase tracking-widest text-slate-800">Dine In</p>
                                <p class="text-xs font-bold text-slate-400">80 Order</p>
                            </div>
                            <div class="text-center">
                                <div class="w-24 h-24 rounded-full border-8 border-rose-400 flex items-center justify-center mb-3 text-2xl font-black text-rose-500 shadow-inner">35%</div>
                                <p class="text-[10px] font-black uppercase tracking-widest text-slate-800">Take Away</p>
                                <p class="text-xs font-bold text-slate-400">44 Order</p>
                            </div>
                        </div>
                    </div>
                </div>

            </div>
        </div>
    </SidebarFnB>
</template>
<style scoped>
.custom-scrollbar::-webkit-scrollbar { height: 6px; width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
</style>