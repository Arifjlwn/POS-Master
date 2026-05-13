<script setup>
import { ref, onMounted, nextTick, watch } from 'vue';
import api from '../api';
import Sidebar from '../components/Sidebar.vue';
import Chart from 'chart.js/auto';

const reportData = ref(null);
const isLoading = ref(true);
const chartCanvas = ref(null);
let salesChart = null;

const storeName = ref(localStorage.getItem('storeName') || 'Toko UMKM');

// --- FILTER TANGGAL ---
const today = new Date().toISOString().split('T')[0];
const startDate = ref(today);
const endDate = ref(today);

const formatRupiah = (angka) => {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(angka);
};

const renderChart = (grafikData) => {
    if (!chartCanvas.value || !grafikData) return;
    if (salesChart) salesChart.destroy();

    salesChart = new Chart(chartCanvas.value, {
        type: 'line',
        data: {
            labels: grafikData.map(d => d.tanggal),
            datasets: [
                {
                    label: 'Omzet Penjualan',
                    data: grafikData.map(d => d.omzet),
                    borderColor: '#2563eb', // Biru
                    backgroundColor: 'rgba(37, 99, 235, 0.05)',
                    borderWidth: 4,
                    tension: 0.4,
                    fill: true,
                    pointRadius: 4
                },
                // 🚀 Jika Mas Arif mau tambahkan garis laba di grafik nantinya:
                {
                    label: 'Laba Bersih',
                    data: grafikData.map(d => d.laba || 0), 
                    borderColor: '#10b981', // Hijau
                    backgroundColor: 'transparent',
                    borderWidth: 3,
                    borderDash: [5, 5],
                    tension: 0.4,
                    pointRadius: 4
                } 
                
            ]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: { legend: { display: false } },
            scales: {
                y: { beginAtZero: true, grid: { borderDash: [5, 5] }, ticks: { callback: (v) => 'Rp ' + v.toLocaleString() } },
                x: { grid: { display: false } }
            }
        }
    });
};

const fetchData = async () => {
    isLoading.value = true;
    try {
        const response = await api.get('/report/dashboard', {
            params: { start_date: startDate.value, end_date: endDate.value }
        });
        reportData.value = response.data.data;
        
        nextTick(() => {
            renderChart(reportData.value.grafik_penjualan);
        });
    } catch (error) {
        console.error("Gagal tarik data dashboard:", error);
    } finally {
        isLoading.value = false;
    }
};

onMounted(fetchData);
watch([startDate, endDate], fetchData);
</script>

<template>
    <Sidebar>
        <div class="p-6 md:p-10 max-w-7xl mx-auto font-sans bg-slate-50 min-h-screen">
            
            <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-6 mb-10">
                <div>
                    <h1 class="text-4xl font-black text-slate-900 tracking-tighter">Business Insight</h1>
                    <p class="text-slate-500 font-bold uppercase text-[10px] tracking-[0.3em]">
                        {{ storeName }} - ANALYTICS CENTER
                    </p>
                </div>

                <div class="flex items-center gap-3 bg-white p-2 rounded-[24px] shadow-xl shadow-slate-200/50 border border-white">
                    <input type="date" v-model="startDate" class="text-xs font-black border-none focus:ring-0 bg-slate-50 rounded-xl px-4 py-2">
                    <span class="text-slate-300 font-bold">to</span>
                    <input type="date" v-model="endDate" class="text-xs font-black border-none focus:ring-0 bg-slate-50 rounded-xl px-4 py-2">
                </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-10">
    <div class="bg-white p-8 rounded-[35px] border border-white shadow-xl shadow-slate-200/60 relative">
        <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Total Omzet</p>
        <p class="text-2xl font-black text-slate-900">
            {{ formatRupiah(reportData?.summary?.total_omzet || 0) }}
        </p>
    </div>
    
    <div class="bg-emerald-500 p-8 rounded-[35px] shadow-xl shadow-emerald-100 text-white">
        <p class="text-[10px] font-black text-emerald-100 uppercase tracking-widest mb-2">Estimasi Laba</p>
        <p class="text-2xl font-black">
            {{ formatRupiah(reportData?.summary?.total_laba || 0) }}
        </p>
    </div>

    <div class="bg-white p-8 rounded-[35px] border border-white shadow-xl shadow-slate-200/60">
        <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Barang Terjual</p>
        <p class="text-2xl font-black text-slate-900">
            {{ reportData?.summary?.total_produk_terjual || 0 }} <span class="text-sm font-bold text-slate-300 italic">Pcs</span>
        </p>
    </div>

    <div class="bg-blue-600 p-8 rounded-[35px] shadow-xl shadow-blue-100 text-white">
        <p class="text-[10px] font-black text-blue-200 uppercase tracking-widest mb-2">Avg. Per Struk</p>
        <p class="text-2xl font-black">
            {{ formatRupiah(reportData?.summary?.avg_transaksi || 0) }}
        </p>
    </div>
</div>

            <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
                <div class="lg:col-span-2 bg-white rounded-[45px] p-10 shadow-xl shadow-slate-200/50 border border-white">
                    <div class="flex items-center justify-between mb-8">
                        <h3 class="text-xl font-black text-slate-800 tracking-tight">Performa Penjualan</h3>
                        <div class="flex gap-2">
                            <span class="w-3 h-3 rounded-full bg-blue-600 animate-pulse"></span>
                            <span class="text-[10px] font-black text-slate-400 uppercase tracking-tighter">Realtime Tracking</span>
                        </div>
                    </div>
                    <div class="h-[400px]">
                        <canvas ref="chartCanvas"></canvas>
                    </div>
                </div>

                <div class="bg-white rounded-[45px] p-10 shadow-xl shadow-slate-200/50 border border-white flex flex-col">
                    <h3 class="text-xl font-black text-slate-800 tracking-tight mb-8">Top Selling 🏆</h3>
                    <div class="space-y-8 flex-1">
                        <div v-for="(item, index) in reportData?.best_sellers" :key="index" class="flex items-center gap-4 group">
                            <div class="w-12 h-12 rounded-2xl bg-slate-50 flex items-center justify-center font-black text-slate-300 group-hover:bg-blue-600 group-hover:text-white transition-all duration-300">
                                0{{ index + 1 }}
                            </div>
                            <div class="flex-1">
                                <p class="text-sm font-black text-slate-800 uppercase leading-none mb-1">{{ item.nama_produk }}</p>
                                <p class="text-[10px] font-bold text-slate-400 uppercase">{{ item.qty_terjual }} Pcs Terjual</p>
                            </div>
                            <p class="text-sm font-black text-blue-600">{{ formatRupiah(item.total_omzet) }}</p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </Sidebar>
</template>