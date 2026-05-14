<script setup>
import { ref, onMounted } from 'vue';
import api from '../../api.js';
import Sidebar from '../../components/Sidebar.vue';

const reports = ref([]);
const selectedDetail = ref(null);

const fetchReports = async () => {
    try {
        const res = await api.get('/stock-opname/history');
        reports.value = res.data.data;
    } catch (err) {
        console.error("Gagal ambil history SO", err);
    }
};

const showDetail = (report) => {
    selectedDetail.value = report;
};

// Hitung total kerugian/keuntungan dari selisih (Selisih * Harga Modal)
const calculateLoss = (details) => {
    return details.reduce((acc, item) => {
        return acc + (item.selisih * (item.product.harga_modal || 0));
    }, 0);
};

onMounted(fetchReports);
</script>

<template>
    <Sidebar>
        <div class="p-6 max-w-6xl mx-auto">
            <h1 class="text-2xl font-black text-slate-900 mb-6 uppercase tracking-tighter">📊 Laporan Selisih Stock Opname</h1>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div class="bg-white rounded-[30px] shadow-xl overflow-hidden border border-slate-100">
                    <table class="w-full text-left">
                        <thead class="bg-slate-50 text-[10px] font-black text-slate-400 uppercase tracking-widest">
                            <tr>
                                <th class="p-4">Tanggal & Catatan</th>
                                <th class="p-4 text-center">Status</th>
                                <th class="p-4"></th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-slate-50">
                            <tr v-for="r in reports" :key="r.id" @click="showDetail(r)" class="hover:bg-blue-50 cursor-pointer transition-colors">
                                <td class="p-4">
                                    <div class="font-black text-slate-700 uppercase">{{ new Date(r.created_at).toLocaleString() }}</div>
                                    <div class="text-[10px] text-slate-400 font-bold italic">{{ r.notes }}</div>
                                </td>
                                <td class="p-4 text-center">
                                    <span v-if="calculateLoss(r.details) < 0" class="bg-red-100 text-red-600 px-2 py-1 rounded-lg text-[10px] font-black">MINUS</span>
                                    <span v-else class="bg-green-100 text-green-600 px-2 py-1 rounded-lg text-[10px] font-black">NORMAL/PLUS</span>
                                </td>
                                <td class="p-4 text-blue-600 font-black text-xs">LIHAT ➔</td>
                            </tr>
                        </tbody>
                    </table>
                </div>

                <div v-if="selectedDetail" class="bg-white rounded-[30px] shadow-xl p-6 border-2 border-blue-600">
                    <div class="mb-4">
                        <h2 class="font-black text-slate-800 uppercase tracking-tight">Detail Selisih Barang</h2>
                        <p class="text-[10px] text-slate-400 font-bold">TOTAL ESTIMASI: 
                            <span :class="calculateLoss(selectedDetail.details) < 0 ? 'text-red-500' : 'text-green-500'">
                                Rp {{ Math.abs(calculateLoss(selectedDetail.details)).toLocaleString() }} 
                                ({{ calculateLoss(selectedDetail.details) < 0 ? 'Rugi' : 'Untung' }})
                            </span>
                        </p>
                    </div>

                    <div class="space-y-3">
                        <div v-for="d in selectedDetail.details" :key="d.id" class="p-4 bg-slate-50 rounded-2xl flex justify-between items-center border border-slate-100">
                            <div>
                                <div class="text-xs font-black text-slate-700 uppercase">{{ d.product.nama_produk }}</div>
                                <div class="text-[9px] font-bold text-slate-400">SIS: {{ d.system_qty }} | FISIK: {{ d.actual_qty }}</div>
                            </div>
                            <div :class="d.selisih < 0 ? 'text-red-600' : 'text-green-600'" class="font-black text-sm">
                                {{ d.selisih > 0 ? '+' : '' }}{{ d.selisih }}
                            </div>
                        </div>
                    </div>
                </div>
                <div v-else class="flex items-center justify-center bg-slate-100 rounded-[30px] border-2 border-dashed border-slate-200 text-slate-400 font-bold italic">
                    Pilih laporan di kiri untuk melihat detail selisih
                </div>
            </div>
        </div>
    </Sidebar>
</template>