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
    // Auto-scroll ke bagian detail kalau di HP
    if (window.innerWidth < 1024) {
        setTimeout(() => {
            document.getElementById('print-area').scrollIntoView({ behavior: 'smooth' });
        }, 100);
    }
};

// Hitung total kerugian/keuntungan dari selisih (Selisih * Harga Modal)
const calculateLoss = (details) => {
    return details.reduce((acc, item) => {
        return acc + (item.selisih * (item.product.harga_modal || 0));
    }, 0);
};

// Format Tanggal ala Indonesia
const formatDate = (dateStr) => {
    return new Intl.DateTimeFormat('id-ID', { 
        weekday: 'long', year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit' 
    }).format(new Date(dateStr));
};

// Fungsi Print Sakti
const printReport = () => {
    window.print();
};

onMounted(fetchReports);
</script>

<template>
    <Sidebar>
        <div class="p-4 md:p-8 lg:p-10 max-w-7xl mx-auto font-sans bg-[#f8fafc] min-h-screen">
            
            <div class="no-print bg-gradient-to-br from-slate-900 via-slate-800 to-indigo-900 rounded-[40px] p-8 md:p-10 mb-8 text-white shadow-2xl flex flex-col md:flex-row items-center justify-between relative overflow-hidden border border-white/10">
                <div class="absolute -right-10 -top-10 opacity-10">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-64 h-64" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1"><path stroke-linecap="round" stroke-linejoin="round" d="M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" /></svg>
                </div>
                
                <div class="z-10 text-center md:text-left">
                    <h1 class="text-3xl md:text-4xl font-black tracking-tighter mb-2 uppercase italic leading-none">Audit <span class="text-indigo-400">Report</span></h1>
                    <p class="text-slate-300 font-bold text-[10px] uppercase tracking-[0.2em] flex items-center justify-center md:justify-start gap-2">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" /></svg>
                        Laporan Selisih Stock Opname
                    </p>
                </div>
            </div>

            <div class="grid grid-cols-1 lg:grid-cols-12 gap-8">
                <div class="lg:col-span-5 no-print space-y-6">
                    <div class="bg-white rounded-[32px] shadow-xl shadow-slate-200/50 border border-slate-100 overflow-hidden">
                        <div class="p-6 bg-slate-50/50 border-b border-slate-100">
                            <h3 class="font-black text-slate-800 text-xs uppercase tracking-widest flex items-center gap-2">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                                Riwayat Audit
                            </h3>
                        </div>
                        <div class="overflow-x-auto custom-scrollbar max-h-[600px] overflow-y-auto">
                            <table class="w-full text-left">
                                <tbody class="divide-y divide-slate-50">
                                    <tr v-if="reports.length === 0">
                                        <td class="p-10 text-center text-slate-400 font-bold text-xs uppercase tracking-widest italic">Belum ada riwayat SO</td>
                                    </tr>
                                    <tr v-for="r in reports" :key="r.id" @click="showDetail(r)" 
                                        class="cursor-pointer transition-all group"
                                        :class="selectedDetail?.id === r.id ? 'bg-indigo-50 border-l-4 border-l-indigo-600' : 'hover:bg-slate-50 border-l-4 border-l-transparent'">
                                        
                                        <td class="p-5">
                                            <div class="font-black text-slate-800 text-xs uppercase tracking-tight group-hover:text-indigo-600 transition-colors">{{ formatDate(r.created_at) }}</div>
                                            <div class="text-[10px] text-slate-400 font-bold mt-1 line-clamp-1">{{ r.notes }}</div>
                                            
                                            <div class="mt-3 flex items-center gap-2">
                                                <span v-if="calculateLoss(r.details) < 0" class="bg-red-100 text-red-600 px-2.5 py-1 rounded-md text-[9px] font-black uppercase tracking-widest shadow-sm">
                                                    RUGI / MINUS
                                                </span>
                                                <span v-else-if="calculateLoss(r.details) > 0" class="bg-emerald-100 text-emerald-600 px-2.5 py-1 rounded-md text-[9px] font-black uppercase tracking-widest shadow-sm">
                                                    SURPLUS / PLUS
                                                </span>
                                                <span v-else class="bg-slate-100 text-slate-600 px-2.5 py-1 rounded-md text-[9px] font-black uppercase tracking-widest shadow-sm">
                                                    BALANCE
                                                </span>
                                            </div>
                                        </td>
                                        <td class="p-5 text-right">
                                            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-300 group-hover:text-indigo-600 transition-colors inline-block" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" /></svg>
                                        </td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>

                <div class="lg:col-span-7" id="print-area">
                    <div v-if="selectedDetail" class="bg-white rounded-[32px] print:rounded-none shadow-xl shadow-slate-200/50 print:shadow-none p-6 md:p-8 border-2 border-indigo-100 print:border-none relative">
                        
                        <div class="flex flex-col sm:flex-row sm:items-start justify-between gap-4 border-b-2 border-slate-800 pb-6 mb-6">
                            <div class="text-center">
                                <h1 class="text-2xl font-black text-slate-900 uppercase tracking-[0.2em] leading-tight print:text-center">Laporan Stock Opname</h1>
                                <p class="text-xs font-bold text-slate-500 mt-2 uppercase tracking-widest print:text-center">{{ formatDate(selectedDetail.created_at) }}</p>
                            </div>

                            <button @click="printReport" class="no-print shrink-0 p-3 bg-indigo-50 hover:bg-indigo-600 text-indigo-600 hover:text-white rounded-2xl transition-all shadow-sm flex items-center justify-center gap-2" title="Print A4">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z" /></svg>
                                <span class="text-[10px] font-black uppercase tracking-widest">Print Dokumen</span>
                            </button>
                        </div>

                        <div class="mb-8 p-6 rounded-2xl border-2 border-dashed" :class="calculateLoss(selectedDetail.details) < 0 ? 'bg-red-50 border-red-200' : 'bg-emerald-50 border-emerald-200'">
                            <h3 class="text-[10px] font-black text-slate-500 uppercase tracking-[0.2em] mb-2 text-center">Estimasi Valuasi Selisih (Berdasarkan Modal)</h3>
                            <div class="text-3xl font-black text-center" :class="calculateLoss(selectedDetail.details) < 0 ? 'text-red-600' : 'text-emerald-600'">
                                {{ calculateLoss(selectedDetail.details) < 0 ? '-' : '+' }} Rp {{ Math.abs(calculateLoss(selectedDetail.details)).toLocaleString('id-ID') }}
                            </div>
                            <p class="text-[10px] font-bold text-center mt-2 uppercase tracking-widest" :class="calculateLoss(selectedDetail.details) < 0 ? 'text-red-400' : 'text-emerald-400'">
                                Catatan: {{ selectedDetail.notes }}
                            </p>
                        </div>

                        <h4 class="font-black text-slate-800 text-[10px] uppercase tracking-[0.2em] mb-4 flex items-center gap-2">
                            <span class="w-1.5 h-4 bg-indigo-600 rounded-full"></span> Rincian Item Dihitung
                        </h4>
                        
                        <div class="overflow-x-auto">
                            <table class="w-full text-left border-collapse">
                                <thead>
                                    <tr class="bg-slate-100/50 text-[9px] font-black text-slate-500 uppercase tracking-widest border-b-2 border-slate-200">
                                        <th class="p-3">Nama Barang</th>
                                        <th class="p-3 text-center">Sistem</th>
                                        <th class="p-3 text-center">Fisik</th>
                                        <th class="p-3 text-center">Selisih</th>
                                    </tr>
                                </thead>
                                <tbody class="divide-y divide-slate-100">
                                    <tr v-for="d in selectedDetail.details" :key="d.id" class="text-xs font-bold text-slate-700">
                                        <td class="p-3 uppercase">{{ d.product.nama_produk }}</td>
                                        <td class="p-3 text-center text-slate-400">{{ d.system_qty }}</td>
                                        <td class="p-3 text-center text-indigo-600 font-black">{{ d.actual_qty }}</td>
                                        <td class="p-3 text-center font-black" :class="d.selisih < 0 ? 'text-red-600' : d.selisih > 0 ? 'text-emerald-600' : 'text-slate-300'">
                                            {{ d.selisih > 0 ? '+' : '' }}{{ d.selisih }}
                                        </td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>

                        <div class="hidden print:flex justify-end mt-16 pt-8">
                            <div class="text-center w-48">
                                <p class="text-[10px] font-bold uppercase tracking-widest mb-16">Mengetahui,</p>
                                <div class="border-b border-slate-800"></div>
                                <p class="text-[10px] font-bold uppercase tracking-widest mt-2">Owner / Auditor</p>
                            </div>
                        </div>

                    </div>
                    
                    <div v-else class="no-print flex flex-col items-center justify-center h-full min-h-[400px] bg-slate-50 rounded-[32px] border-2 border-dashed border-slate-200 opacity-60">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-20 h-20 text-slate-300 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" /></svg>
                        <p class="text-xs font-black text-slate-400 uppercase tracking-widest text-center px-6">Pilih salah satu riwayat di sebelah kiri <br> untuk melihat detail audit.</p>
                    </div>
                </div>
            </div>
        </div>
    </Sidebar>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { height: 6px; width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #94a3b8; }

/* 🚀 CSS SAKTI UNTUK MENCETAK A4 KERTAS */
@media print {
    /* Sembunyikan semua elemen di luar area print */
    body * {
        visibility: hidden;
    }
    /* Sembunyikan elemen yang punya class no-print */
    .no-print {
        display: none !important;
    }
    /* Tampilkan area print dan isinya */
    #print-area, #print-area * {
        visibility: visible;
    }
    /* Posisikan area print ke ujung kiri atas kertas */
    #print-area {
        position: absolute;
        left: 0;
        top: 0;
        width: 100%;
        margin: 0;
        padding: 0;
    }
    /* Seting kertas A4 */
    @page {
        size: A4 portrait;
        margin: 1.5cm;
    }
}
</style>