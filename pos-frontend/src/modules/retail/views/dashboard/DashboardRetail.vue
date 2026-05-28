<script setup>
import { onMounted } from 'vue';
import { useDashboard } from '../../composables/useDashboard.js';
import Sidebar from '../../components/Sidebar.vue';

const {
    reportData,
    isLoading,
    storeName,
    startDate,
    endDate,
    lineChartCanvas,
    pieChartCanvas,
    profitMargin,
    formatRupiah,
    setQuickFilter,
    fetchData
} = useDashboard();

onMounted(fetchData);

// Format number utility
const formatNumber = (val) => {
    return Number(val).toLocaleString('id-ID');
};
</script>

<template>
    <Sidebar>
        <div class="p-4 md:p-6 lg:p-8 max-w-[1600px] mx-auto font-sans bg-[#F8FAFC] min-h-screen">
            
            <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-4 mb-6">
                <div>
                    <h1 class="text-2xl md:text-3xl font-black text-slate-900 tracking-tight uppercase">{{ storeName }}</h1>
                    <div class="flex items-center gap-2 mt-1">
                        <span class="w-2 h-2 rounded-full bg-indigo-500 animate-pulse"></span>
                        <p class="text-slate-500 font-bold uppercase text-[10px] tracking-widest">Executive Analytics Dashboard</p>
                    </div>
                </div>

                <div class="flex flex-col sm:flex-row items-start sm:items-center gap-3">
                    <div class="flex bg-slate-200/50 p-1 rounded-xl">
                        <button @click="setQuickFilter(0)" class="px-4 py-1.5 rounded-lg text-[10px] font-black uppercase tracking-wider transition-all" :class="startDate === endDate ? 'bg-white text-indigo-600 shadow-sm' : 'text-slate-500 hover:text-slate-700'">Hari Ini</button>
                        <button @click="setQuickFilter(6)" class="px-4 py-1.5 rounded-lg text-[10px] font-black uppercase tracking-wider transition-all" :class="startDate !== endDate && (new Date(endDate) - new Date(startDate)) / (1000 * 3600 * 24) === 6 ? 'bg-white text-indigo-600 shadow-sm' : 'text-slate-500 hover:text-slate-700'">7 Hari</button>
                        <button @click="setQuickFilter(29)" class="px-4 py-1.5 rounded-lg text-[10px] font-black uppercase tracking-wider transition-all" :class="startDate !== endDate && (new Date(endDate) - new Date(startDate)) / (1000 * 3600 * 24) === 29 ? 'bg-white text-indigo-600 shadow-sm' : 'text-slate-500 hover:text-slate-700'">30 Hari</button>
                    </div>
                    <div class="flex items-center bg-white p-1.5 rounded-xl shadow-sm border border-slate-200">
                        <input type="date" v-model="startDate" class="text-xs font-bold text-slate-700 border-none focus:ring-0 bg-transparent cursor-pointer py-1 px-2 outline-none">
                        <span class="text-[10px] text-slate-300 font-black px-1">-</span>
                        <input type="date" v-model="endDate" class="text-xs font-bold text-slate-700 border-none focus:ring-0 bg-transparent cursor-pointer py-1 px-2 outline-none">
                    </div>
                    
                    <!-- 🚀 TOMBOL REFRESH SAKTI -->
                    <button @click="fetchData" :disabled="isLoading" class="p-2.5 bg-white text-slate-500 rounded-xl shadow-sm border border-slate-200 hover:bg-indigo-50 hover:border-indigo-200 hover:text-indigo-600 transition-all disabled:opacity-50 group" title="Refresh Data Analytics">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" :class="{'animate-spin text-indigo-600': isLoading, 'group-hover:rotate-180 transition-transform duration-500': !isLoading}" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                        </svg>
                    </button>
                </div>
            </div>

            <div v-if="isLoading" class="py-32 flex flex-col items-center justify-center bg-white rounded-[24px] border border-slate-200 shadow-sm">
                <div class="w-10 h-10 border-4 border-indigo-100 border-t-indigo-600 rounded-full animate-spin mb-4"></div>
                <p class="text-slate-400 font-bold text-xs uppercase tracking-widest animate-pulse">Menghubungkan Data Lapangan...</p>
            </div>

            <div v-else>
                <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-5 mb-8">
                    
                    <div class="bg-white p-6 rounded-[24px] border border-slate-200 shadow-sm flex flex-col justify-between group hover:border-indigo-300 transition-all">
                        <div class="flex justify-between items-start mb-4">
                            <p class="text-[10px] font-black text-slate-400 uppercase tracking-wider">Total Omzet</p>
                            <div class="w-8 h-8 rounded-xl bg-indigo-50 text-indigo-600 flex items-center justify-center group-hover:scale-110 transition-transform">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><rect width="20" height="14" x="2" y="5" rx="2"/><line x1="2" y1="10" x2="22" y2="10"/></svg>
                            </div>
                        </div>
                        <div>
                            <p class="text-2xl font-black text-slate-900 tracking-tight">{{ formatRupiah(reportData?.summary?.total_omzet) }}</p>
                            <p class="text-[10px] font-bold text-slate-400 mt-1">Bruto Penjualan Toko</p>
                        </div>
                    </div>
                    
                    <div class="bg-white p-6 rounded-[24px] border border-slate-200 shadow-sm flex flex-col justify-between group hover:border-emerald-300 transition-all relative overflow-hidden">
                        <div class="absolute right-0 top-0 w-20 h-20 bg-gradient-to-br from-emerald-50 to-transparent opacity-60 rounded-bl-[100%]"></div>
                        <div class="flex justify-between items-start mb-4 relative z-10">
                            <p class="text-[10px] font-black text-slate-400 uppercase tracking-wider">Estimasi Laba</p>
                            <div class="w-8 h-8 rounded-xl bg-emerald-50 text-emerald-600 flex items-center justify-center group-hover:scale-110 transition-transform">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>
                            </div>
                        </div>
                        <div class="relative z-10 flex items-end justify-between">
                            <div>
                                <p class="text-2xl font-black text-emerald-600 tracking-tight">{{ formatRupiah(reportData?.summary?.total_laba) }}</p>
                                <p class="text-[10px] font-bold text-slate-400 mt-1">Netto Bersih (Sebelum Loss)</p>
                            </div>
                            <div class="bg-emerald-100 text-emerald-700 px-2 py-1 rounded-lg text-[10px] font-black tracking-widest border border-emerald-200">
                                {{ profitMargin }}% MARGIN
                            </div>
                        </div>
                    </div>

                    <div class="bg-slate-900 p-6 rounded-[24px] shadow-lg flex flex-col justify-between group relative overflow-hidden">
                        <div class="absolute right-0 top-0 w-24 h-24 bg-white/5 rounded-bl-[100%]"></div>
                        <div class="flex justify-between items-start mb-4 relative z-10">
                            <p class="text-[10px] font-black text-slate-400 uppercase tracking-wider">Avg. Per Struk</p>
                            <div class="w-8 h-8 rounded-xl bg-white/10 text-white flex items-center justify-center group-hover:scale-110 transition-transform">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><rect width="16" height="20" x="4" y="2" rx="2" ry="2"/><line x1="8" y1="6" x2="16" y2="6"/><line x1="8" y1="10" x2="16" y2="10"/></svg>
                            </div>
                        </div>
                        <div class="relative z-10">
                            <p class="text-2xl font-black text-white tracking-tight">{{ formatRupiah(reportData?.summary?.avg_transaksi) }}</p>
                            <p class="text-[10px] font-bold text-slate-400 mt-1">Daya Beli Pelanggan</p>
                        </div>
                    </div>

                    <div class="bg-white p-6 rounded-[24px] border border-slate-200 shadow-sm flex flex-col justify-between group hover:border-amber-300 transition-all">
                        <div class="flex justify-between items-start mb-4">
                            <p class="text-[10px] font-black text-slate-400 uppercase tracking-wider">Barang Terjual</p>
                            <div class="w-8 h-8 rounded-xl bg-amber-50 text-amber-600 flex items-center justify-center group-hover:scale-110 transition-transform">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="m7.5 4.27 9 5.15"/><path d="M21 8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16Z"/><path d="m3.3 7 8.7 5 8.7-5"/><path d="M12 22V12"/></svg>
                            </div>
                        </div>
                        <div>
                            <p class="text-2xl font-black text-slate-900 tracking-tight">{{ formatNumber(reportData?.summary?.total_produk_terjual) }} <span class="text-[12px] font-bold text-slate-400">Total</span></p>
                            <p class="text-[10px] font-bold text-slate-400 mt-1">Produk Keluar Kasir</p>
                        </div>
                    </div>

                    <div class="bg-white p-6 rounded-[24px] border border-slate-200 shadow-sm flex flex-col justify-between group hover:border-rose-300 transition-all relative overflow-hidden">
                        <div class="absolute right-0 top-0 w-16 h-16 bg-gradient-to-br from-rose-50 to-transparent opacity-80 rounded-bl-[100%]"></div>
                        <div class="flex justify-between items-start mb-4 relative z-10">
                            <p class="text-[10px] font-black text-slate-400 uppercase tracking-wider">Waste & Retur</p>
                            <div class="w-8 h-8 rounded-xl bg-rose-50 text-rose-600 flex items-center justify-center group-hover:scale-110 transition-transform">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/></svg>
                            </div>
                        </div>
                        <div class="relative z-10">
                            <p class="text-2xl font-black text-rose-600 tracking-tight">{{ formatRupiah(reportData?.summary?.total_retur_loss) }}</p>
                            <p class="text-[10px] font-bold text-slate-400 mt-1">Total: <span class="text-rose-500">{{ formatNumber(reportData?.summary?.total_retur_qty) }}</span> Dibuang/Rusak</p>
                        </div>
                    </div>

                    <div class="bg-white p-6 rounded-[24px] border border-slate-200 shadow-sm flex flex-col justify-between group hover:border-purple-300 transition-all relative overflow-hidden">
                        <div class="absolute right-0 top-0 w-16 h-16 bg-gradient-to-br from-purple-50 to-transparent opacity-80 rounded-bl-[100%]"></div>
                        <div class="flex justify-between items-start mb-4 relative z-10">
                            <p class="text-[10px] font-black text-slate-400 uppercase tracking-wider">Selisih SO (Kehilangan)</p>
                            <div class="w-8 h-8 rounded-xl bg-purple-50 text-purple-600 flex items-center justify-center group-hover:scale-110 transition-transform">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
                            </div>
                        </div>
                        <div class="relative z-10">
                            <p class="text-2xl font-black text-purple-600 tracking-tight">{{ formatRupiah(reportData?.summary?.total_so_loss) }}</p>
                            <p class="text-[10px] font-bold text-slate-400 mt-1">Total: <span class="text-purple-500">{{ formatNumber(reportData?.summary?.total_so_qty) }}</span> Hilang/Minus</p>
                        </div>
                    </div>
                </div>

                <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-6">
                    <div class="lg:col-span-2 bg-white rounded-[24px] p-6 border border-slate-200 shadow-sm flex flex-col">
                        <div class="flex items-center justify-between mb-6">
                            <div>
                                <h3 class="text-sm font-black text-slate-800 tracking-tight uppercase">Analitik Pendapatan</h3>
                                <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mt-0.5">Tren Omzet vs Laba</p>
                            </div>
                            <div class="flex gap-3">
                                <div class="flex items-center gap-1.5"><span class="w-2.5 h-2.5 rounded-sm bg-[#4f46e5]"></span><span class="text-[9px] font-black text-slate-600 uppercase">Omzet</span></div>
                                <div class="flex items-center gap-1.5"><span class="w-2.5 h-2.5 rounded-sm border-2 border-dashed border-[#10b981]"></span><span class="text-[9px] font-black text-slate-600 uppercase">Laba</span></div>
                                <div class="flex items-center gap-1.5"><span class="w-2.5 h-2.5 rounded-sm bg-[#e11d48]"></span><span class="text-[9px] font-black text-slate-600 uppercase">Retur</span></div>
                            </div>
                        </div>
                        <div class="flex-1 min-h-[250px] w-full relative">
                            <canvas ref="lineChartCanvas"></canvas>
                        </div>
                    </div>

                    <div class="bg-white rounded-[24px] p-6 border border-slate-200 shadow-sm flex flex-col">
                        <div>
                            <h3 class="text-sm font-black text-slate-800 tracking-tight uppercase">Komposisi Top 5</h3>
                            <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mt-0.5">Berdasarkan Qty Terjual</p>
                        </div>
                        
                        <div v-if="!reportData?.best_sellers || reportData.best_sellers.length === 0" class="flex-1 flex flex-col items-center justify-center text-slate-300 py-10">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-10 h-10 mb-2 opacity-50" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><circle cx="12" cy="12" r="10"/><path d="M12 2a10 10 0 0 1 10 10"/></svg>
                            <p class="text-[10px] font-bold uppercase tracking-widest">Data Kosong</p>
                        </div>

                        <div v-else class="flex-1 flex flex-col justify-center mt-6">
                            <div class="h-[180px] relative w-full mb-4 flex justify-center">
                                <canvas ref="pieChartCanvas"></canvas>
                                <div class="absolute inset-0 flex flex-col items-center justify-center pointer-events-none">
                                    <span class="text-2xl font-black text-slate-800">{{ formatNumber(reportData.summary.total_produk_terjual) }}</span>
                                    <span class="text-[8px] font-bold text-slate-400 uppercase tracking-widest">Total Keseluruhan</span>
                                </div>
                            </div>
                            <div class="space-y-1.5 mt-auto max-h-[110px] overflow-y-auto custom-scrollbar">
                                <div v-for="(item, index) in reportData.best_sellers.slice(0, 5)" :key="index" class="flex items-center justify-between text-[10px] font-bold text-slate-600 py-0.5 border-b border-slate-50">
                                    <div class="flex items-center gap-2 truncate pr-2">
                                        <span class="w-2 h-2 rounded-full shrink-0" :style="{ backgroundColor: ['#4f46e5', '#3b82f6', '#0ea5e9', '#10b981', '#f59e0b'][index] }"></span>
                                        <span class="truncate uppercase">{{ item.nama_produk }}</span>
                                    </div>
                                    <span class="font-black text-slate-900 shrink-0">{{ formatNumber(item.qty_terjual) }} {{ item.satuan_dasar || 'Pcs' }}</span>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="bg-white rounded-[24px] border border-slate-200 shadow-sm overflow-hidden">
                    <div class="p-5 border-b border-slate-100 flex items-center justify-between bg-slate-50/50">
                        <h3 class="text-xs font-black text-slate-800 uppercase tracking-widest flex items-center gap-2">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-amber-500" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M6 9H4.5a2.5 2.5 0 0 1 0-5H6"/><path d="M18 9h1.5a2.5 2.5 0 0 0 0-5H18"/><path d="M4 22h16"/><path d="M18 2H6v7a6 6 0 0 0 12 0V2Z"/></svg>
                            Peringkat Barang Terlaris (Product Leaderboard)
                        </h3>
                    </div>
                    
                    <div v-if="!reportData?.best_sellers || reportData.best_sellers.length === 0" class="p-10 text-center">
                        <p class="text-xs font-bold text-slate-400 uppercase tracking-widest">Belum ada transaksi di rentang tanggal ini.</p>
                    </div>

                    <div v-else class="overflow-x-auto">
                        <table class="w-full text-left whitespace-nowrap">
                            <thead class="bg-white border-b border-slate-100">
                                <tr class="text-[9px] font-black text-slate-400 uppercase tracking-[0.2em]">
                                    <th class="px-6 py-4 w-16 text-center">Rank</th>
                                    <th class="px-6 py-4">Nama Produk / SKU</th>
                                    <th class="px-6 py-4 text-center">Qty Terjual</th>
                                    <th class="px-6 py-4 text-right">Kontribusi Omzet</th>
                                </tr>
                            </thead>
                            <tbody class="divide-y divide-slate-50">
                                <tr v-for="(item, index) in reportData.best_sellers" :key="index" class="hover:bg-slate-50/50 transition-colors">
                                    <td class="px-6 py-3 text-center">
                                        <span class="inline-flex items-center justify-center w-6 h-6 rounded-md text-[10px] font-black" 
                                                :class="index === 0 ? 'bg-amber-100 text-amber-600' : (index === 1 ? 'bg-slate-200 text-slate-600' : (index === 2 ? 'bg-orange-100 text-orange-700' : 'text-slate-400'))">
                                            #{{ index + 1 }}
                                        </span>
                                    </td>
                                    <td class="px-6 py-3">
                                        <div class="font-black text-slate-800 text-xs uppercase">{{ item.nama_produk }}</div>
                                        <div class="text-[9px] font-bold text-slate-400 tracking-widest mt-0.5">{{ item.sku || 'SKU-NA' }}</div>
                                    </td>
                                    <td class="px-6 py-3 text-center">
                                        <span class="font-black text-slate-700 text-xs bg-slate-100 px-2.5 py-1 rounded-md">
                                            {{ formatNumber(item.qty_terjual) }} {{ item.satuan_dasar || 'Pcs' }}
                                        </span>
                                    </td>
                                    <td class="px-6 py-3 text-right">
                                        <span class="font-black text-indigo-600 text-sm">{{ formatRupiah(item.total_omzet) }}</span>
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </div>

            </div>
        </div>
    </Sidebar>
</template>

<style scoped>
input[type="date"]::-webkit-calendar-picker-indicator {
    cursor: pointer;
    opacity: 0.5;
    transition: 0.2s;
}
input[type="date"]::-webkit-calendar-picker-indicator:hover {
    opacity: 1;
}
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
</style>