<script setup>
import { ref, onMounted, computed, onUnmounted, nextTick, watch } from 'vue';
import SidebarLaundry from '../components/SidebarLaundry.vue';
import api from '../../../../api.js';
import Swal from 'sweetalert2';
import Chart from 'chart.js/auto';

const activeTab = ref('ringkasan'); 
const selectedPeriod = ref('bulan_ini');
const searchQuery = ref('');
const isLoading = ref(false);
let pollingInterval = null;

const stats = ref({ omset: 0, total_order: 0, avg_transaction: 0, tunai: 0, qris: 0, debit: 0, piutang: 0 });
const riwayat = ref([]);
const showBuktiModal = ref(false);
const selectedBuktiUrl = ref('');

const chartCanvas = ref(null);
let omsetChart = null;

const formatRupiah = (angka) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(angka || 0);
const formatDate = (dateStr) => {
    if (!dateStr) return '-';
    const d = new Date(dateStr);
    if (isNaN(d.getTime())) return '-';
    return new Intl.DateTimeFormat('id-ID', { day: '2-digit', month: 'short', year: 'numeric' }).format(d);
};

// 🚀 FUNGSI UTAMA: EXPORT PDF / PRINT LAPORAN
const exportPDF = () => {
    // Kasih info toast kecil biar kelihatan profesional
    Swal.fire({
        toast: true,
        position: 'top-end',
        icon: 'info',
        title: 'Mempersiapkan dokumen PDF...',
        showConfirmButton: false,
        timer: 1000
    });
    
    // Tunggu animasi sebentar langsung buka jendela cetak/save PDF browser
    setTimeout(() => {
        window.print();
    }, 1000);
};

const renderGrafik = () => {
    if (!chartCanvas.value) return;
    if (omsetChart) omsetChart.destroy();

    const dataOmsetPerHari = {};
    riwayat.value.forEach(trx => {
        if (trx.status_bayar === 'LUNAS') {
            const tgl = formatDate(trx.created_at);
            dataOmsetPerHari[tgl] = (dataOmsetPerHari[tgl] || 0) + trx.total_harga;
        }
    });

    let labels = Object.keys(dataOmsetPerHari).reverse();
    let dataValues = Object.values(dataOmsetPerHari).reverse();

    if (labels.length === 0) {
        labels = ['Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu', 'Minggu'];
        dataValues = [0, 0, 0, 0, 0, 0, 0];
    }

    omsetChart = new Chart(chartCanvas.value, {
        type: 'bar',
        data: {
            labels: labels,
            datasets: [{
                label: 'Pendapatan (Rp)',
                data: dataValues,
                backgroundColor: '#4f46e5', 
                borderRadius: 8, 
                barThickness: 30, 
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: { display: false },
                tooltip: {
                    callbacks: { label: (context) => formatRupiah(context.raw) }
                }
            },
            scales: {
                y: { 
                    beginAtZero: true, 
                    grid: { color: '#f1f5f9', drawBorder: false }, 
                    ticks: { display: false } 
                },
                x: { 
                    grid: { display: false, drawBorder: false },
                    ticks: { font: { weight: 'bold', family: 'sans-serif' }, color: '#94a3b8' }
                }
            }
        }
    });
};

const fetchReportData = async (isBackground = false) => {
    if (!isBackground) isLoading.value = true;
    try {
        const response = await api.get(`/laundry/report?period=${selectedPeriod.value}`);
        stats.value = {
            omset: response.data.ringkasan.total_omset || 0,
            total_order: response.data.ringkasan.total_order || 0,
            avg_transaction: response.data.ringkasan.rata_rata || 0,
            tunai: response.data.ringkasan.tunai || 0,
            qris: response.data.ringkasan.qris || 0,
            debit: response.data.ringkasan.debit || 0,
            piutang: response.data.ringkasan.piutang || 0
        };
        riwayat.value = response.data.transaksi || [];
        
        if (activeTab.value === 'ringkasan') {
            await nextTick();
            renderGrafik();
        }
    } catch (error) {
        console.error("Gagal sinkronisasi data");
    } finally {
        if (!isBackground) isLoading.value = false;
    }
};

watch(activeTab, async (newTab) => {
    if (newTab === 'ringkasan') {
        await nextTick();
        renderGrafik();
    }
});

onMounted(() => {
    fetchReportData();
    pollingInterval = setInterval(() => fetchReportData(true), 10000); 
});

onUnmounted(() => { if (pollingInterval) clearInterval(pollingInterval); });

const filteredRiwayat = computed(() => {
    if (!searchQuery.value) return riwayat.value;
    const q = searchQuery.value.toLowerCase();
    return riwayat.value.filter(t => t.no_invoice.toLowerCase().includes(q) || t.pelanggan.toLowerCase().includes(q));
});

const vipCustomers = computed(() => {
    const customerMap = {};
    riwayat.value.forEach(trx => {
        if (trx.status_bayar === 'LUNAS') {
            if (!customerMap[trx.pelanggan]) customerMap[trx.pelanggan] = { nama: trx.pelanggan, total_belanja: 0, kunjungan: 0 };
            customerMap[trx.pelanggan].total_belanja += trx.total_harga;
            customerMap[trx.pelanggan].kunjungan += 1;
        }
    });
    return Object.values(customerMap).sort((a, b) => b.total_belanja - a.total_belanja).slice(0, 5);
});

const bukaBuktiTransfer = (url) => { selectedBuktiUrl.value = url; showBuktiModal.value = true; };
</script>

<template>
    <SidebarLaundry>
        <div id="print-area" class="flex-1 flex flex-col h-full bg-slate-50/50 overflow-hidden relative">
            
            <div class="p-5 lg:p-8 shrink-0 bg-white border-b border-slate-200 flex flex-col gap-5 md:gap-6 z-10 shadow-sm relative header-section">
                <div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
                    <div class="flex items-center gap-4">
                        <div class="w-12 h-12 bg-indigo-50 border border-indigo-100 rounded-2xl flex items-center justify-center shrink-0 shadow-inner logo-box">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-indigo-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M3 3v18h18"/><path d="m19 9-5 5-4-4-3 3"/></svg>
                        </div>
                        <div>
                            <h1 class="text-xl md:text-2xl font-black tracking-tighter uppercase text-slate-800 leading-tight report-title">Analitik Enterprise</h1>
                            <p class="text-[10px] font-black text-slate-400 mt-0.5 uppercase tracking-[0.2em] report-sub">Pusat Kontrol Owner</p>
                        </div>
                    </div>
                    
                    <button @click="exportPDF" class="w-full md:w-auto bg-slate-900 hover:bg-indigo-600 text-white px-6 py-3 rounded-xl font-black text-[10px] uppercase tracking-[0.2em] flex items-center justify-center gap-2 shadow-xl shadow-slate-200 transition-all btn-export">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
                        Export PDF / Cetak
                    </button>
                </div>

                <div class="flex flex-col xl:flex-row gap-4 items-start xl:items-center justify-between border-t border-slate-100 pt-5 action-bar">
                    <div class="flex gap-2 w-full xl:w-auto overflow-x-auto hide-scrollbar pb-2 xl:pb-0 filter-tabs">
                        <button @click="activeTab = 'ringkasan'" :class="activeTab === 'ringkasan' ? 'bg-indigo-600 text-white shadow-md' : 'bg-slate-100 text-slate-500'" class="whitespace-nowrap px-6 py-2.5 rounded-xl text-[10px] font-black uppercase tracking-widest transition-all">Ringkasan</button>
                        <button @click="activeTab = 'keuangan'" :class="activeTab === 'keuangan' ? 'bg-indigo-600 text-white shadow-md' : 'bg-slate-100 text-slate-500'" class="whitespace-nowrap px-6 py-2.5 rounded-xl text-[10px] font-black uppercase tracking-widest transition-all flex items-center gap-2">Keuangan <span v-if="stats.piutang > 0" class="w-2 h-2 rounded-full bg-rose-400 animate-pulse"></span></button>
                        <button @click="activeTab = 'analitik'" :class="activeTab === 'analitik' ? 'bg-indigo-600 text-white shadow-md' : 'bg-slate-100 text-slate-500'" class="whitespace-nowrap px-6 py-2.5 rounded-xl text-[10px] font-black uppercase tracking-widest transition-all">VIP CRM</button>
                    </div>
                    <div class="w-full xl:w-auto relative group period-select">
                        <select v-model="selectedPeriod" @change="fetchReportData" class="w-full xl:w-56 px-4 py-3 bg-slate-50 border-2 border-slate-200 rounded-xl font-bold text-[11px] uppercase tracking-widest text-slate-700 outline-none focus:border-indigo-500 transition-all cursor-pointer">
                            <option value="minggu_ini">Data Minggu Ini</option>
                            <option value="bulan_ini">Data Bulan Ini</option>
                            <option value="tahun_ini">Data Tahun Ini</option>
                        </select>
                    </div>
                </div>
            </div>

            <div class="flex-1 overflow-y-auto custom-scrollbar p-5 lg:p-8 pb-24 content-section">
                <div v-if="isLoading && riwayat.length === 0" class="flex flex-col items-center justify-center py-32 text-slate-400 loading-box">
                    <div class="w-14 h-14 border-4 border-indigo-100 border-t-indigo-600 rounded-full animate-spin mb-6"></div>
                    <p class="font-black text-xs uppercase tracking-[0.2em] animate-pulse">Sinkronisasi Server...</p>
                </div>

                <template v-else>
                    <div v-if="activeTab === 'ringkasan'" class="space-y-6 lg:space-y-8 animate-[fadeIn_0.3s_ease-out] tab-content-block">
                        <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-5 lg:gap-6 stats-grid">
                            <div class="bg-gradient-to-br from-indigo-600 via-indigo-700 to-slate-900 p-6 lg:p-8 rounded-[24px] lg:rounded-[32px] text-white shadow-xl shadow-indigo-200/50 card-stat">
                                <h3 class="text-[10px] lg:text-[11px] font-black uppercase tracking-[0.2em] text-indigo-200 mb-2">Total Estimasi Omset</h3>
                                <p class="text-3xl lg:text-4xl font-black tracking-tighter text-white">{{ formatRupiah(stats.omset) }}</p>
                            </div>
                            <div class="bg-white p-6 lg:p-8 rounded-[24px] lg:rounded-[32px] border border-slate-200 shadow-sm card-stat">
                                <h3 class="text-[10px] lg:text-[11px] font-black text-slate-400 uppercase tracking-widest mb-2">Volume Transaksi</h3>
                                <p class="text-3xl lg:text-4xl font-black text-slate-800 tracking-tighter">{{ stats.total_order }} <span class="text-sm font-bold text-slate-400">Nota</span></p>
                            </div>
                            <div class="bg-white p-6 lg:p-8 rounded-[24px] lg:rounded-[32px] border border-slate-200 shadow-sm sm:col-span-2 xl:col-span-1 card-stat">
                                <h3 class="text-[10px] lg:text-[11px] font-black text-slate-400 uppercase tracking-widest mb-2">Rata-rata / Nota (ARPU)</h3>
                                <p class="text-3xl lg:text-4xl font-black text-slate-800 tracking-tighter">{{ formatRupiah(stats.avg_transaction) }}</p>
                            </div>
                        </div>

                        <div class="bg-white p-6 lg:p-8 rounded-[24px] lg:rounded-[32px] border border-slate-200 shadow-sm flex flex-col chart-container-card">
                            <div class="flex justify-between items-center mb-6">
                                <h3 class="text-xs lg:text-sm font-black text-slate-800 uppercase tracking-widest flex items-center gap-3">
                                    <div class="w-2 h-6 bg-indigo-600 rounded-full btn-decor"></div> Tren Pendapatan Omset
                                </h3>
                            </div>
                            <div class="relative w-full h-64 lg:h-72 canvas-box">
                                <canvas ref="chartCanvas"></canvas>
                            </div>
                        </div>
                    </div>

                    <div v-if="activeTab === 'keuangan'" class="animate-[fadeIn_0.3s_ease-out] space-y-6 lg:space-y-8 tab-content-block">
                        <div class="grid grid-cols-2 xl:grid-cols-4 gap-4 lg:gap-6 money-grid">
                            <div class="bg-white p-5 rounded-[24px] border border-slate-200 shadow-sm card-money"><h3 class="text-[10px] font-black text-slate-400 uppercase mb-2">Tunai</h3><p class="text-xl lg:text-2xl font-black text-slate-800">{{ formatRupiah(stats.tunai) }}</p></div>
                            <div class="bg-white p-5 rounded-[24px] border border-emerald-100 shadow-sm card-money"><h3 class="text-[10px] font-black text-emerald-500 uppercase mb-2">Transfer / QRIS</h3><p class="text-xl lg:text-2xl font-black text-emerald-600">{{ formatRupiah(stats.qris) }}</p></div>
                            <div class="bg-white p-5 rounded-[24px] border border-slate-200 shadow-sm card-money"><h3 class="text-[10px] font-black text-slate-400 uppercase mb-2">Mesin EDC</h3><p class="text-xl lg:text-2xl font-black text-slate-800">{{ formatRupiah(stats.debit) }}</p></div>
                            <div class="bg-rose-50 p-5 rounded-[24px] border border-rose-200 shadow-sm card-money"><h3 class="text-[10px] font-black text-rose-500 uppercase mb-2 flex items-center gap-2"><div class="w-2 h-2 bg-rose-500 rounded-full animate-pulse blink-dot"></div> Piutang Berjalan</h3><p class="text-xl lg:text-2xl font-black text-rose-600">{{ formatRupiah(stats.piutang) }}</p></div>
                        </div>

                        <div class="bg-white border border-slate-200 rounded-[32px] shadow-sm overflow-hidden flex flex-col table-card-section">
                            <div class="p-6 border-b border-slate-100 flex justify-between items-center bg-slate-50/50 table-filter-bar">
                                <h2 class="text-xs font-black text-slate-800 uppercase tracking-[0.2em]">Buku Besar</h2>
                                <input v-model="searchQuery" type="text" placeholder="Cari Resi..." class="px-4 py-2 bg-white border-2 border-slate-200 rounded-xl font-bold text-xs outline-none search-input-ledger">
                            </div>
                            <div class="overflow-x-auto">
                                <table class="w-full text-left whitespace-nowrap report-main-table">
                                    <thead>
                                        <tr class="bg-white text-[9px] uppercase tracking-[0.2em] text-slate-400 border-b-2 border-slate-100">
                                            <th class="p-5 font-black text-center">Data Invoice</th>
                                            <th class="p-5 font-black text-center">Pelanggan</th>
                                            <th class="p-5 font-black text-center">Status Pembayaran</th>
                                            <th class="p-5 font-black text-center">Nominal</th>
                                            <th class="p-5 font-black text-center audit-header">Audit</th>
                                        </tr>
                                    </thead>
                                    <tbody class="text-slate-700">
                                        <tr v-for="trx in filteredRiwayat" :key="trx.id" class="border-b border-slate-50 hover:bg-slate-50 table-row-data">
                                            <td class="p-5">
                                                <p class="font-black text-xs text-slate-800 uppercase">{{ trx.no_invoice }}</p>
                                                <p class="text-[9px] font-bold text-slate-400 mt-1">{{ formatDate(trx.created_at) }}</p>
                                            </td>
                                            <td class="p-5 font-bold text-xs text-slate-700 uppercase text-center">{{ trx.pelanggan }}</td>
                                            <td class="p-5 text-center">
                                                <span v-if="trx.status_bayar === 'BELUM_LUNAS'" class="text-[9px] font-black px-3 py-1 rounded-lg uppercase bg-rose-100 text-rose-600 border border-rose-200">PIUTANG</span>
                                                <span v-else class="text-[9px] font-black px-3 py-1 rounded-lg uppercase shadow-sm border" :class="trx.metode_bayar === 'QRIS' ? 'bg-emerald-50 text-emerald-600 border-emerald-200' : 'bg-slate-100 text-slate-700 border-slate-200'">Lunas • {{ trx.metode_bayar }}</span>
                                            </td>
                                            <td class="p-5 text-center font-black text-xs tracking-tight text-slate-800">{{ formatRupiah(trx.total_harga) }}</td>
                                            <td class="p-5 text-center audit-cell">
                                                <button v-if="trx.metode_bayar === 'QRIS' && trx.bukti_transfer" @click="bukaBuktiTransfer(trx.bukti_transfer)" class="bg-white border-2 border-slate-200 text-slate-700 px-3 py-1.5 rounded-lg text-[9px] font-black uppercase hover:border-indigo-500 hover:text-indigo-600">Bukti</button>
                                                <span v-else class="text-slate-300 font-bold text-[10px] uppercase">-</span>
                                            </td>
                                        </tr>
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    </div>

                    <div v-if="activeTab === 'analitik'" class="animate-[fadeIn_0.3s_ease-out] tab-content-block">
                        <div class="bg-white border border-slate-200 rounded-[32px] p-6 lg:p-8 shadow-sm vip-section-card">
                            <h2 class="text-sm font-black text-slate-800 uppercase tracking-widest mb-6">Top 5 Pelanggan VIP</h2>
                            <div class="space-y-4">
                                <div v-for="(vip, index) in vipCustomers" :key="index" class="flex items-center justify-between p-4 bg-slate-50 rounded-2xl border border-slate-100 vip-row">
                                    <div class="flex items-center gap-4">
                                        <div class="w-8 h-8 rounded-full bg-slate-200 flex items-center justify-center font-black text-xs text-slate-600">{{ index + 1 }}</div>
                                        <div><p class="font-black text-xs uppercase text-slate-800">{{ vip.nama }}</p><p class="text-[10px] font-bold text-slate-400">{{ vip.kunjungan }}x Kunjungan</p></div>
                                    </div>
                                    <p class="font-black text-sm text-indigo-600">{{ formatRupiah(vip.total_belanja) }}</p>
                                </div>
                            </div>
                        </div>
                    </div>
                </template>
            </div>
        </div>
        
        <Teleport to="body">
            <div v-if="showBuktiModal" class="fixed inset-0 z-[9999] flex items-center justify-center bg-slate-900/90 backdrop-blur-sm p-4 animate-[fadeIn_0.2s_ease-out]" @click.self="showBuktiModal = false">
                <div class="bg-white rounded-[32px] overflow-hidden flex flex-col w-full max-w-sm">
                    <div class="bg-slate-900 p-4 flex justify-between items-center text-white"><h3 class="font-black text-[10px] uppercase tracking-[0.2em]">Audit Bukti</h3><button @click="showBuktiModal = false" class="text-slate-400 hover:text-white">X</button></div>
                    <img :src="selectedBuktiUrl" class="w-full max-h-[60vh] object-contain p-4 bg-slate-100" />
                </div>
            </div>
        </Teleport>
    </SidebarLaundry>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { height: 6px; width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
.hide-scrollbar::-webkit-scrollbar { display: none; } 
@keyframes fadeIn { from { opacity: 0; transform: translateY(15px); } to { opacity: 1; transform: translateY(0); } }

/* =========================================================================
   🚀 SAKTI CSS: KHUSUS MODE CETAK / SAVE AS PDF (ANTI KERTAS KOSONG)
   ========================================================================= */
@media print {
    /* 1. Atur Kertas A4 & Hilangkan margin berlebih */
    @page { size: A4 portrait; margin: 1cm; }

    /* 2. PAKSA SEMUA CONTAINER BUKA Gembok OVERFLOW! (Ini penyakit utamanya beb) */
    body, html, #app, :deep(#app) {
        height: auto !important;
        overflow: visible !important;
        background: white !important;
    }

    /* Buka semua gembok scroll dari Tailwind */
    #print-area, 
    .content-section,
    :deep(.flex-1), 
    :deep(.overflow-y-auto), 
    :deep(.overflow-hidden) {
        display: block !important;
        height: auto !important;
        overflow: visible !important;
        position: static !important;
        width: 100% !important;
        padding: 0 !important;
    }

    /* 3. Sembunyikan elemen navigasi yang bikin semak */
    :deep(aside), 
    .btn-export, 
    .action-bar, 
    .table-filter-bar, 
    .audit-header, 
    .audit-cell,
    .logo-box,
    .btn-decor {
        display: none !important;
    }

    /* 4. Rapikan Grid ke Mode Cetak */
    .stats-grid {
        display: grid !important;
        grid-template-columns: repeat(3, 1fr) !important;
        gap: 15px !important;
    }

    .money-grid {
        display: grid !important;
        grid-template-columns: repeat(4, 1fr) !important;
        gap: 15px !important;
    }

    /* 5. Paksa warna render PDF (Biar gradasinya ikut kecetak) */
    .card-stat, .card-money, .chart-container-card, .table-card-section {
        border: 1px solid #e2e8f0 !important;
        border-radius: 12px !important;
        page-break-inside: avoid !important; /* Biar tabel gak kepotong di tengah halaman */
        -webkit-print-color-adjust: exact !important;
        print-color-adjust: exact !important;
        box-shadow: none !important;
        margin-bottom: 20px !important;
    }

    .bg-gradient-to-br {
        background: #4f46e5 !important;
        color: white !important;
    }

    .report-title {
        color: #0f172a !important;
        font-size: 20px !important;
        margin-bottom: 20px !important;
    }

    /* 6. Fix Tabel */
    .report-main-table {
        width: 100% !important;
        border-collapse: collapse !important;
    }
    .table-row-data {
        page-break-inside: avoid !important;
    }
}
</style>