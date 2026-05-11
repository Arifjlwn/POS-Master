<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import api from '../api';
import Sidebar from '../components/Sidebar.vue';

const reportData = ref(null);
const isLoading = ref(true);
const errorMessage = ref('');
const router = useRouter();

const formatRupiah = (angka) => {
    return new Intl.NumberFormat('id-ID').format(angka);
};

const handleLogout = () => {
    localStorage.removeItem('token');
    router.push('/login'); // <-- Kembali ke login tanpa kedip
};

onMounted(async () => {
    try {
        const response = await api.get('/report/dashboard');
        reportData.value = response.data.data;
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
        <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
            <div v-if="isLoading" class="flex flex-col justify-center items-center h-64 gap-4">
                <div class="animate-spin rounded-full h-12 w-12 border-b-4 border-blue-600"></div>
                <p class="text-gray-500 font-bold animate-pulse">Menarik data dari server...</p>
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
                        <h1 class="text-3xl font-black text-gray-800 tracking-tight">Ringkasan Hari Ini</h1>
                        <p class="text-gray-500 text-sm mt-1 font-medium">Pantau performa penjualan dan pergerakan stok toko secara real-time.</p>
                    </div>
                    <div class="bg-white px-4 py-2 rounded-lg shadow-sm border border-gray-200 text-sm font-bold text-gray-600">
                        📅 {{ new Date().toLocaleDateString('id-ID', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' }) }}
                    </div>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
                    <div class="bg-white rounded-2xl p-6 shadow-sm border border-gray-100 flex items-center gap-5 hover:-translate-y-1 hover:shadow-lg transition-all duration-300">
                        <div class="w-16 h-16 rounded-2xl bg-gradient-to-br from-blue-100 to-blue-200 text-blue-600 flex items-center justify-center text-3xl shadow-inner border border-blue-100">💰</div>
                        <div>
                            <p class="text-xs font-black text-gray-400 uppercase tracking-widest mb-1">Total Omzet</p>
                            <p class="text-2xl lg:text-3xl font-black text-gray-800 leading-none">Rp {{ formatRupiah(reportData?.summary?.omzet_hari_ini || 0) }}</p>
                        </div>
                    </div>

                    <div class="bg-white rounded-2xl p-6 shadow-sm border border-gray-100 flex items-center gap-5 hover:-translate-y-1 hover:shadow-lg transition-all duration-300">
                        <div class="w-16 h-16 rounded-2xl bg-gradient-to-br from-emerald-100 to-emerald-200 text-emerald-600 flex items-center justify-center text-3xl shadow-inner border border-emerald-100">🧾</div>
                        <div>
                            <p class="text-xs font-black text-gray-400 uppercase tracking-widest mb-1">Total Transaksi</p>
                            <p class="text-2xl lg:text-3xl font-black text-gray-800 leading-none">{{ reportData?.summary?.jumlah_transaksi || 0 }} <span class="text-sm font-bold text-gray-400">Struk</span></p>
                        </div>
                    </div>

                    <div class="bg-white rounded-2xl p-6 shadow-sm border border-gray-100 flex items-center gap-5 hover:-translate-y-1 hover:shadow-lg transition-all duration-300">
                        <div class="w-16 h-16 rounded-2xl bg-gradient-to-br from-orange-100 to-orange-200 text-orange-600 flex items-center justify-center text-3xl shadow-inner border border-orange-100">📦</div>
                        <div>
                            <p class="text-xs font-black text-gray-400 uppercase tracking-widest mb-1">Barang Terjual</p>
                            <p class="text-2xl lg:text-3xl font-black text-gray-800 leading-none">{{ reportData?.summary?.total_produk_terjual || 0 }} <span class="text-sm font-bold text-gray-400">Pcs</span></p>
                        </div>
                    </div>
                </div>

                <div class="bg-white rounded-2xl shadow-sm border border-gray-200 overflow-hidden">
                    <div class="p-5 border-b border-gray-200 flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4 bg-gray-50">
                        <div>
                            <h2 class="text-lg font-black text-gray-800 flex items-center gap-2"><span class="text-xl">⚠️</span> Peringatan Stok Menipis</h2>
                        </div>
                    </div>
                    <div class="overflow-x-auto">
                        <table v-if="reportData?.low_stock?.length > 0" class="w-full text-left border-collapse">
                            <thead>
                                <tr class="bg-white text-gray-400 text-xs uppercase tracking-widest border-b-2 border-gray-100">
                                    <th class="p-4 font-black">Kode SKU</th>
                                    <th class="p-4 font-black">Nama Barang</th>
                                    <th class="p-4 font-black text-center">Sisa Stok</th>
                                </tr>
                            </thead>
                            <tbody class="divide-y divide-gray-100">
                                <tr v-for="item in reportData.low_stock" :key="item.id" class="hover:bg-blue-50/50 transition-colors">
                                    <td class="p-4 text-sm font-mono font-bold text-gray-500">{{ item.sku || `SKU-${item.id}` }}</td>
                                    <td class="p-4 text-sm font-bold text-gray-800">{{ item.nama_produk }}</td>
                                    <td class="p-4 text-center">
                                        <div class="inline-flex items-center justify-center w-10 h-10 rounded-full bg-red-50 border-2 border-red-100 text-lg font-black text-red-600">{{ item.stok }}</div>
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                        <div v-else class="py-16 text-center">
                            <div class="w-20 h-20 bg-green-50 text-green-500 rounded-full flex items-center justify-center text-4xl mb-4 border-4 border-green-100 shadow-sm mx-auto">✅</div>
                            <p class="text-gray-800 font-black text-xl mb-1">Aman Terkendali Bos!</p>
                        </div>
                    </div>
                </div>
            </div>
        </main>
    </Sidebar>
</template>