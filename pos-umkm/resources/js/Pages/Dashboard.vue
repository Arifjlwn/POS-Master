<script setup>
<<<<<<< HEAD
import { ref, onMounted } from 'vue';
import api from '../api.js'; // Panggil si kurir Axios

const reportData = ref(null);
const isLoading = ref(true);
const errorMessage = ref('');

// Fungsi untuk format angka ke Rupiah
const formatRupiah = (angka) => {
    return new Intl.NumberFormat('id-ID').format(angka);
};

// Fungsi Logout (Hapus tiket dari brankas, balik ke login)
const handleLogout = () => {
    localStorage.removeItem('token');
    window.location.href = '/login';
};

// onMounted akan berjalan otomatis saat halaman pertama kali dibuka
onMounted(async () => {
    try {
        // Tembak API Golang untuk minta data laporan
        const response = await api.get('/report/dashboard');
        reportData.value = response.data.data;
    } catch (error) {
        if (error.response && error.response.status === 403) {
            errorMessage.value = "Halaman ini khusus Bos (Owner). Kasir dilarang masuk! 🛑";
        } else if (error.response && error.response.status === 401) {
            // Kalau token kadaluarsa atau belum login
            handleLogout();
        } else {
            errorMessage.value = "Gagal menarik data dari server Golang.";
        }
    } finally {
        isLoading.value = false;
    }
});
</script>

<template>
    <div class="dashboard-container">
        <header class="header">
            <h2>📊 Dashboard POS UMKM</h2>
            <button @click="handleLogout" class="btn-logout">Logout</button>
        </header>

        <div v-if="isLoading" class="loading">Memuat data laporan...</div>

        <div v-else-if="errorMessage" class="error">
            {{ errorMessage }}
        </div>

        <div v-else class="content">
            <div class="summary-cards">
                <div class="card bg-blue">
                    <h3>Omzet Hari Ini</h3>
                    <p>Rp {{ formatRupiah(reportData?.summary?.omzet_hari_ini || 0) }}</p>
                </div>
                <div class="card bg-green">
                    <h3>Total Transaksi</h3>
                    <p>{{ reportData?.summary?.jumlah_transaksi || 0 }} Struk</p>
                </div>
                <div class="card bg-orange">
                    <h3>Barang Terjual</h3>
                    <p>{{ reportData?.summary?.total_produk_terjual || 0 }} Pcs</p>
                </div>
            </div>

            <div class="low-stock-section">
                <h3>⚠️ Peringatan Stok Menipis (< 10)</h3>
                <table v-if="reportData?.low_stock?.length > 0">
                    <thead>
                        <tr>
                            <th>SKU</th>
                            <th>Nama Barang</th>
                            <th>Sisa Stok</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="item in reportData.low_stock" :key="item.id">
                            <td>{{ item.sku }}</td>
                            <td>{{ item.nama_produk }}</td>
                            <td class="text-red"><strong>{{ item.stok }}</strong></td>
                        </tr>
                    </tbody>
                </table>
                <p v-else class="safe-stock">Aman Bos! Belum ada barang yang perlu di-restock.</p>
            </div>
        </div>
    </div>
=======
import AppLayout from '@/Layouts/AppLayout.vue';
import { Head } from '@inertiajs/vue3';
import { ref, onMounted } from 'vue';
import {
  Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend, ArcElement, BarElement
} from 'chart.js';
import { Line, Doughnut } from 'vue-chartjs';

// Registrasi komponen Chart.js
ChartJS.register( CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend, ArcElement, BarElement );

const props = defineProps({
    chartData: Object
});

// Konfigurasi Chart Tren Penjualan (Line Chart)
const salesChartData = {
    labels: props.chartData.labels,
    datasets: [{
        label: 'Omzet Penjualan (Rp)',
        backgroundColor: 'rgba(37, 99, 235, 0.2)', // Blue
        borderColor: '#2563eb',
        borderWidth: 3,
        pointBackgroundColor: '#ffffff',
        pointBorderColor: '#2563eb',
        pointBorderWidth: 2,
        pointRadius: 5,
        fill: true,
        data: props.chartData.sales,
        tension: 0.4 // Bikin garisnya melengkung halus (bezier)
    }]
};

const salesChartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
        legend: { display: false },
        tooltip: {
            callbacks: {
                label: function(context) {
                    let label = context.dataset.label || '';
                    if (label) label += ': ';
                    if (context.parsed.y !== null) {
                        label += new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(context.parsed.y);
                    }
                    return label;
                }
            }
        }
    },
    scales: {
        y: { beginAtZero: true, grid: { borderDash: [5, 5] }, ticks: { callback: (value) => 'Rp ' + (value/1000) + 'k' } },
        x: { grid: { display: false } }
    }
};

// Konfigurasi Chart Top Produk (Doughnut Chart)
const productChartData = {
    labels: props.chartData.topProducts,
    datasets: [{
        backgroundColor: ['#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6'],
        borderWidth: 0,
        data: props.chartData.topQty,
    }]
};

const productChartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    cutout: '70%',
    plugins: {
        legend: { position: 'bottom', labels: { usePointStyle: true, padding: 20 } }
    }
};

// Summary Cards (Dummy untuk visual)
const stats = [
    { name: 'Omzet Hari Ini', value: 'Rp 3.800.000', trend: '+15%', trendUp: true, icon: '💰', color: 'bg-green-500' },
    { name: 'Transaksi Selesai', value: '84', trend: '+5%', trendUp: true, icon: '🧾', color: 'bg-blue-500' },
    { name: 'Stok Kritis', value: '12', trend: 'Peringatan', trendUp: false, icon: '⚠️', color: 'bg-red-500' },
];
</script>

<template>
    <Head title="Dashboard Owner" />

    <AppLayout>
        <div class="p-6 md:p-8 max-w-7xl mx-auto">

            <div class="flex flex-col md:flex-row md:items-center justify-between mb-8 gap-4">
                <div>
                    <h1 class="text-3xl font-black text-gray-900 tracking-tight">Kinerja Toko 🚀</h1>
                    <p class="text-gray-500 font-medium mt-1">Pantauan analitik <span class="text-blue-600 font-bold">{{ $page.props.auth.store?.nama_toko }}</span>.</p>
                </div>
                <div class="bg-white border border-gray-200 shadow-sm px-4 py-2 rounded-xl text-sm font-bold text-gray-600 flex items-center gap-2">
                    <span>📅</span> 7 Hari Terakhir
                </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
                <div v-for="item in stats" :key="item.name" class="bg-white p-6 rounded-3xl shadow-sm border border-gray-100 relative overflow-hidden group">
                    <div class="absolute -right-4 -top-4 w-24 h-24 rounded-full opacity-10 transition-transform group-hover:scale-150 duration-500" :class="item.color"></div>
                    <div class="flex justify-between items-start mb-4 relative z-10">
                        <div class="text-xs font-black text-gray-400 uppercase tracking-widest">{{ item.name }}</div>
                        <div class="w-10 h-10 rounded-xl flex items-center justify-center text-xl shadow-inner bg-gray-50" :class="item.color.replace('bg-', 'text-')">
                            {{ item.icon }}
                        </div>
                    </div>
                    <div class="text-3xl font-black text-gray-900 relative z-10">{{ item.value }}</div>
                    <div class="mt-2 text-xs font-bold flex items-center gap-1 relative z-10" :class="item.trendUp ? 'text-green-500' : 'text-red-500'">
                        <span>{{ item.trendUp ? '↑' : '↓' }}</span> {{ item.trend }}
                    </div>
                </div>
            </div>

            <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">

                <div class="lg:col-span-2 bg-white p-6 rounded-3xl shadow-sm border border-gray-100">
                    <div class="flex justify-between items-center mb-6">
                        <h3 class="font-black text-gray-800 text-lg">📈 Tren Pendapatan</h3>
                        <button class="text-blue-600 hover:text-blue-800 text-sm font-bold">Lihat Detail</button>
                    </div>
                    <div class="h-[300px] w-full">
                        <Line :data="salesChartData" :options="salesChartOptions" />
                    </div>
                </div>

                <div class="bg-white p-6 rounded-3xl shadow-sm border border-gray-100 flex flex-col">
                    <h3 class="font-black text-gray-800 text-lg mb-2">🔥 Fast Moving Products</h3>
                    <p class="text-xs text-gray-400 font-medium mb-6">Berdasarkan volume (qty) terjual.</p>
                    <div class="flex-1 min-h-[250px] relative flex justify-center items-center">
                        <Doughnut :data="productChartData" :options="productChartOptions" />
                        <div class="absolute inset-0 flex flex-col items-center justify-center pointer-events-none mt-[-20px]">
                            <span class="text-3xl font-black text-gray-900">Top 5</span>
                        </div>
                    </div>
                </div>

            </div>

        </div>
    </AppLayout>
>>>>>>> fe5108ac44404475b1c66a1cb99fadaf9190e970
</template>

<style scoped>
/* CSS Sederhana untuk mempercantik Dashboard */
.dashboard-container {
    max-width: 900px;
    margin: 40px auto;
    font-family: sans-serif;
}
.header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
}
.btn-logout {
    background-color: #dc3545;
    color: white;
    border: none;
    padding: 8px 16px;
    border-radius: 4px;
    cursor: pointer;
}
.summary-cards {
    display: flex;
    gap: 20px;
    margin-bottom: 30px;
}
.card {
    flex: 1;
    padding: 20px;
    border-radius: 8px;
    color: white;
    text-align: center;
}
.card h3 { margin: 0 0 10px 0; font-size: 16px; font-weight: normal; }
.card p { margin: 0; font-size: 24px; font-weight: bold; }
.bg-blue { background-color: #007bff; }
.bg-green { background-color: #28a745; }
.bg-orange { background-color: #fd7e14; }

.low-stock-section { margin-top: 30px; }
table { width: 100%; border-collapse: collapse; margin-top: 10px; }
table, th, td { border: 1px solid #ddd; }
th, td { padding: 12px; text-align: left; }
th { background-color: #f4f4f4; }
.text-red { color: red; }
.safe-stock { color: #28a745; font-style: italic; }
.error { background: #ffe6e6; color: red; padding: 15px; border-radius: 5px; }
</style>
