<script setup>
import { ref, onMounted, nextTick } from 'vue';
import { useRouter } from 'vue-router';
import api from '../api';
import Sidebar from '../components/Sidebar.vue';
import Chart from 'chart.js/auto'; // 🚀 Import library grafik

const reportData = ref(null);
const isLoading = ref(true);
const errorMessage = ref('');
const router = useRouter();
const chartCanvas = ref(null);
let salesChart = null;

const formatRupiah = (angka) => {
    return new Intl.NumberFormat('id-ID').format(angka);
};

const handleLogout = () => {
    localStorage.removeItem('token');
    router.push('/login');
};

// Fungsi Render Grafik
const renderChart = (grafikData) => {
    if (!chartCanvas.value || !grafikData || grafikData.length === 0) return;
    
    if (salesChart) salesChart.destroy(); // Hapus grafik lama kalau ada

    salesChart = new Chart(chartCanvas.value, {
        type: 'line',
        data: {
            labels: grafikData.map(d => d.tanggal),
            datasets: [{
                label: 'Omzet Harian (Rp)',
                data: grafikData.map(d => d.omzet),
                borderColor: '#2563eb', // Warna Biru Tailwind
                backgroundColor: 'rgba(37, 99, 235, 0.1)', // Biru transparan
                borderWidth: 3,
                pointBackgroundColor: '#ffffff',
                pointBorderColor: '#2563eb',
                pointBorderWidth: 2,
                pointRadius: 4,
                tension: 0.4, // Bikin garisnya melengkung smooth
                fill: true
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: { display: false } // Sembunyikan legenda biar bersih
            },
            scales: {
                y: { beginAtZero: true, grid: { borderDash: [5, 5] } },
                x: { grid: { display: false } }
            }
        }
    });
};

onMounted(async () => {
    try {
        const response = await api.get('/report/dashboard');
        reportData.value = response.data.data;
        
        // Render grafik setelah data selesai ditarik dan UI siap
        nextTick(() => {
            renderChart(reportData.value.grafik_7_hari);
        });
        
    } catch (error) {
        if (error.response && error.response.status === 403) {
            errorMessage.value = "Halaman ini khusus Bos (Owner). Kasir dilarang masuk! 🛑";
        } else if (error.response && error.response.status === 401) {
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
    <Sidebar>
        <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 font-sans">
            <div v-if="isLoading" class="flex flex-col justify-center items-center h-64 gap-4">
                <div class="animate-spin rounded-full h-12 w-12 border-b-4 border-blue-600"></div>
                <p class="text-gray-500 font-bold animate-pulse">Menyiapkan Laporan Toko...</p>
            </div>

            <div v-else-if="errorMessage" class="bg-red-50 border-l-8 border-red-500 p-6 rounded-xl shadow-sm flex items-center gap-4">
                <div class="text-4xl">🛑</div>
                <div>
                    <h3 class="text-red-800 font-black text-lg">Akses Ditolak!</h3>
                    <p class="text-red-600 font-medium">{{ errorMessage }}</p>
                </div>
            </div>

            <div v-else>
                <div class="mb-8 flex flex-col md:flex-row md:justify-between md:items-end gap-4">
                    <div>
                        <h1 class="text-3xl font-black text-gray-800 tracking-tight">Dashboard & Laporan</h1>
                        <p class="text-gray-500 text-sm mt-1 font-medium">Pantau performa penjualan, grafik, dan pergerakan stok toko secara real-time.</p>
                    </div>
                    <div class="bg-white px-4 py-2 rounded-lg shadow-sm border border-gray-200 text-sm font-bold text-gray-600">
                        📅 {{ new Date().toLocaleDateString('id-ID', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' }) }}
                    </div>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
                    <div class="bg-white rounded-3xl p-6 shadow-sm border border-gray-100 flex items-center gap-5 hover:-translate-y-1 hover:shadow-md transition-all duration-300">
                        <div class="w-16 h-16 rounded-2xl bg-gradient-to-br from-blue-100 to-blue-200 text-blue-600 flex items-center justify-center text-3xl shadow-inner border border-blue-100">💰</div>
                        <div>
                            <p class="text-xs font-black text-gray-400 uppercase tracking-widest mb-1">Omzet Hari Ini</p>
                            <p class="text-2xl lg:text-3xl font-black text-gray-800 leading-none">Rp {{ formatRupiah(reportData?.summary?.omzet_hari_ini || 0) }}</p>
                        </div>
                    </div>

                    <div class="bg-white rounded-3xl p-6 shadow-sm border border-gray-100 flex items-center gap-5 hover:-translate-y-1 hover:shadow-md transition-all duration-300">
                        <div class="w-16 h-16 rounded-2xl bg-gradient-to-br from-emerald-100 to-emerald-200 text-emerald-600 flex items-center justify-center text-3xl shadow-inner border border-emerald-100">🧾</div>
                        <div>
                            <p class="text-xs font-black text-gray-400 uppercase tracking-widest mb-1">Total Transaksi</p>
                            <p class="text-2xl lg:text-3xl font-black text-gray-800 leading-none">{{ reportData?.summary?.jumlah_transaksi || 0 }} <span class="text-sm font-bold text-gray-400">Struk</span></p>
                        </div>
                    </div>

                    <div class="bg-white rounded-3xl p-6 shadow-sm border border-gray-100 flex items-center gap-5 hover:-translate-y-1 hover:shadow-md transition-all duration-300">
                        <div class="w-16 h-16 rounded-2xl bg-gradient-to-br from-orange-100 to-orange-200 text-orange-600 flex items-center justify-center text-3xl shadow-inner border border-orange-100">📦</div>
                        <div>
                            <p class="text-xs font-black text-gray-400 uppercase tracking-widest mb-1">Barang Terjual</p>
                            <p class="text-2xl lg:text-3xl font-black text-gray-800 leading-none">{{ reportData?.summary?.total_produk_terjual || 0 }} <span class="text-sm font-bold text-gray-400">Pcs</span></p>
                        </div>
                    </div>
                </div>

                <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
                    
                    <div class="bg-white rounded-3xl shadow-sm border border-gray-200 overflow-hidden flex flex-col">
                        <div class="p-5 border-b border-gray-100 flex items-center gap-2 bg-gray-50/50">
                            <span class="text-xl">📈</span>
                            <h2 class="text-lg font-black text-gray-800">Grafik Penjualan 7 Hari Terakhir</h2>
                        </div>
                        <div class="p-6 flex-1 relative min-h-[300px]">
                            <canvas ref="chartCanvas"></canvas>
                        </div>
                    </div>

                    <div class="bg-white rounded-3xl shadow-sm border border-gray-200 overflow-hidden flex flex-col">
                        <div class="p-5 border-b border-gray-100 flex items-center gap-2 bg-yellow-50">
                            <span class="text-xl">🏆</span>
                            <h2 class="text-lg font-black text-gray-800">Top 5 Produk Terlaris Bulan Ini</h2>
                        </div>
                        <div class="overflow-x-auto">
                            <table class="w-full text-left">
                                <thead class="bg-white border-b border-gray-100">
                                    <tr>
                                        <th class="px-5 py-4 text-xs font-black text-gray-400 uppercase tracking-widest">Produk</th>
                                        <th class="px-5 py-4 text-xs font-black text-gray-400 uppercase tracking-widest text-center">Terjual</th>
                                        <th class="px-5 py-4 text-xs font-black text-gray-400 uppercase tracking-widest text-right">Omzet</th>
                                    </tr>
                                </thead>
                                <tbody class="divide-y divide-gray-50">
                                    <tr v-if="!reportData?.best_sellers || reportData.best_sellers.length === 0">
                                        <td colspan="3" class="px-5 py-8 text-center text-gray-400 font-medium italic">Belum ada data penjualan bulan ini.</td>
                                    </tr>
                                    <tr v-for="(item, index) in reportData?.best_sellers" :key="index" class="hover:bg-gray-50 transition-colors">
                                        <td class="px-5 py-4">
                                            <div class="flex items-center gap-3">
                                                <div class="w-6 h-6 rounded-full bg-yellow-100 text-yellow-700 flex items-center justify-center font-black text-xs">{{ index + 1 }}</div>
                                                <div>
                                                    <div class="font-bold text-gray-800 line-clamp-1">{{ item.nama_produk }}</div>
                                                    <div class="text-[10px] font-mono font-bold text-gray-400">{{ item.sku }}</div>
                                                </div>
                                            </div>
                                        </td>
                                        <td class="px-5 py-4 text-center">
                                            <span class="bg-blue-50 text-blue-700 font-black px-2 py-1 rounded-lg text-sm">{{ item.qty_terjual }}</span>
                                        </td>
                                        <td class="px-5 py-4 text-right font-black text-gray-800">
                                            Rp {{ formatRupiah(item.total_omzet) }}
                                        </td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>

                <div class="bg-white rounded-3xl shadow-sm border border-gray-200 overflow-hidden">
                    <div class="p-5 border-b border-gray-100 flex items-center gap-2 bg-red-50">
                        <span class="text-xl">⚠️</span>
                        <h2 class="text-lg font-black text-gray-800">Peringatan Stok Menipis (Di Bawah 10 Pcs)</h2>
                    </div>
                    <div class="overflow-x-auto">
                        <table v-if="reportData?.low_stock?.length > 0" class="w-full text-left">
                            <thead class="bg-white border-b border-gray-100">
                                <tr>
                                    <th class="px-6 py-4 text-xs font-black text-gray-400 uppercase tracking-widest">Kode SKU / PLU</th>
                                    <th class="px-6 py-4 text-xs font-black text-gray-400 uppercase tracking-widest">Nama Barang</th>
                                    <th class="px-6 py-4 text-xs font-black text-gray-400 uppercase tracking-widest text-center">Sisa Stok</th>
                                </tr>
                            </thead>
                            <tbody class="divide-y divide-gray-50">
                                <tr v-for="item in reportData.low_stock" :key="item.id" class="hover:bg-red-50/30 transition-colors">
                                    <td class="px-6 py-4 text-sm font-mono font-bold text-gray-500">{{ item.sku || `SKU-${item.id}` }}</td>
                                    <td class="px-6 py-4 text-sm font-bold text-gray-800">{{ item.nama_produk }}</td>
                                    <td class="px-6 py-4 text-center">
                                        <span class="bg-red-100 text-red-700 font-black px-3 py-1.5 rounded-xl text-sm animate-pulse">{{ item.stok }} Pcs</span>
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                        <div v-else class="py-12 text-center">
                            <div class="w-16 h-16 bg-green-50 text-green-500 rounded-full flex items-center justify-center text-3xl mb-3 mx-auto">✅</div>
                            <p class="text-gray-800 font-black text-lg">Aman Terkendali Bos!</p>
                            <p class="text-gray-500 text-sm mt-1">Belum ada barang yang perlu di-restock.</p>
                        </div>
                    </div>
                </div>

            </div>
        </main>
    </Sidebar>
</template>