<script setup>
defineProps({
    detail: Object,
    isOwner: Boolean,
    isApproving: Boolean,
    formatDate: Function,
    calculateLoss: Function,
    reportType: { type: String, default: 'SO' } // 🚀 TERIMA PROP TYPE DI SINI
});
defineEmits(['print', 'approve']);

// Fungsi Format Angka Rupiah & Ribuan
const formatNumber = (val) => {
    if (val === null || val === undefined) return '0';
    return Number(val).toLocaleString('id-ID');
};
</script>

<template>
    <div id="print-area">
        <div v-if="detail" class="bg-white rounded-[32px] print:rounded-none shadow-xl shadow-slate-200/50 print:shadow-none p-6 md:p-8 border-2 border-indigo-100 print:border-none relative">
            
            <div class="flex flex-col sm:flex-row sm:items-start justify-between gap-4 border-b-2 border-slate-800 pb-6 mb-6">
                <div class="text-center sm:text-left">
                    <h1 class="text-2xl font-black text-slate-900 uppercase tracking-[0.2em] leading-tight print:text-center">
                        {{ reportType === 'SO' ? 'Laporan Stock Opname' : 'Laporan Klaim Barang' }}
                    </h1>
                    <p class="text-xs font-bold text-slate-500 mt-2 uppercase tracking-widest print:text-center">{{ formatDate(detail.created_at) }}</p>
                    
                    <div class="mt-3 no-print">
                        <span v-if="detail.status === 'PENDING_APPROVAL'" class="inline-flex items-center gap-1.5 bg-amber-100 text-amber-700 px-3 py-1 rounded-lg text-[10px] font-black uppercase tracking-widest border border-amber-200">
                            <span class="w-1.5 h-1.5 rounded-full bg-amber-500 animate-pulse"></span> Menunggu Persetujuan
                        </span>
                        <span v-else-if="detail.status === 'APPROVED'" class="inline-flex items-center gap-1.5 bg-emerald-100 text-emerald-700 px-3 py-1 rounded-lg text-[10px] font-black uppercase tracking-widest border border-emerald-200">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" /></svg> Disetujui
                        </span>
                    </div>
                </div>

                <div class="no-print flex gap-2">
                    <button v-if="isOwner && detail.status === 'PENDING_APPROVAL'" 
                            @click="$emit('approve', detail.id)" 
                            :disabled="isApproving"
                            class="p-3 px-4 bg-emerald-500 hover:bg-emerald-600 text-white rounded-2xl transition-all shadow-md flex items-center justify-center gap-2">
                        <div v-if="isApproving" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                        <template v-else>
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" /></svg>
                            <span class="text-[10px] font-black uppercase tracking-widest">SETUJUI</span>
                        </template>
                    </button>
                </div>
            </div>

            <div class="mb-8 p-6 rounded-2xl border-2 border-dashed" :class="calculateLoss(detail.details) < 0 ? 'bg-red-50 border-red-200' : 'bg-emerald-50 border-emerald-200'">
                <h3 class="text-[10px] font-black text-slate-500 uppercase tracking-[0.2em] mb-2 text-center">
                    {{ reportType === 'SO' ? 'Estimasi Valuasi Selisih' : 'Total Valuasi Barang Temuan' }}
                </h3>
                <div class="text-3xl font-black text-center" :class="calculateLoss(detail.details) < 0 ? 'text-red-600' : 'text-emerald-600'">
                    {{ calculateLoss(detail.details) < 0 ? '-' : '+' }} Rp {{ formatNumber(Math.abs(calculateLoss(detail.details))) }}
                </div>
                <p class="text-[10px] font-bold text-center mt-2 uppercase tracking-widest" :class="calculateLoss(detail.details) < 0 ? 'text-red-400' : 'text-emerald-400'">
                    Catatan: {{ detail.notes }}
                </p>
            </div>

            <h4 class="font-black text-slate-800 text-[10px] uppercase tracking-[0.2em] mb-4 flex items-center gap-2">
                <span class="w-1.5 h-4 bg-indigo-600 rounded-full"></span> Rincian Item Laporan
            </h4>
            
            <div class="overflow-x-auto">
                <table class="w-full text-left border-collapse min-w-[600px]">
                    
                    <thead>
                        <tr class="bg-slate-100/50 text-[9px] font-black text-slate-500 uppercase tracking-widest border-b-2 border-slate-200">
                            <th class="p-3 w-[40%]">Nama Barang</th>
                            <th class="p-3 text-center">Stok Sistem</th>
                            <th class="p-3 text-center">{{ reportType === 'SO' ? 'Fisik' : 'Fisik Ketemu' }}</th>
                            <th class="p-3 text-center">{{ reportType === 'SO' ? 'Selisih' : 'Total Stok' }}</th>
                        </tr>
                    </thead>

                    <tbody class="divide-y divide-slate-100">
                        <tr v-for="d in detail.details" :key="d.id" class="text-xs font-bold text-slate-700">
                            
                            <td class="p-3">
                                <div class="uppercase">{{ d.product?.nama_produk || 'Produk Terhapus' }}</div>
                                <div v-if="d.alasan" class="text-[9px] text-amber-600 mt-1 flex items-center gap-1 font-bold">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                                    {{ d.alasan }}
                                </div>
                            </td>
                            
                            <td class="p-3 text-center text-slate-400">
                                {{ reportType === 'SO' ? formatNumber(d.system_qty) : formatNumber(d.product?.stok ?? 0) }} 
                                <span class="text-[8px] ml-1">{{ d.product?.satuan_dasar || 'PCS' }}</span>
                            </td>
                            
                            <td class="p-3 text-center text-indigo-600 font-black">
                                <span v-if="reportType === 'KLAIM'" class="text-amber-500">+</span>
                                {{ reportType === 'SO' ? formatNumber(d.actual_qty) : formatNumber(d.qty) }}
                                <span class="text-[8px] ml-1">{{ d.product?.satuan_dasar || 'PCS' }}</span>
                            </td>
                            
                            <td class="p-3 text-center font-black" :class="reportType === 'SO' ? (d.selisih < 0 ? 'text-red-600' : d.selisih > 0 ? 'text-emerald-600' : 'text-slate-300') : 'text-emerald-600'">
                                <template v-if="reportType === 'SO'">
                                    {{ d.selisih > 0 ? '+' : '' }}{{ formatNumber(d.selisih) }}
                                </template>
                                <template v-else>
                                    {{ formatNumber((d.product?.stok ?? 0) + d.qty) }}
                                </template>
                                <span class="text-[8px] ml-1 text-slate-400">{{ d.product?.satuan_dasar || 'PCS' }}</span>
                            </td>

                        </tr>
                    </tbody>

                </table>
            </div>
        </div>
        
        <div v-else class="no-print flex flex-col items-center justify-center h-full min-h-[400px] bg-slate-50 rounded-[32px] border-2 border-dashed border-slate-200 opacity-60">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-20 h-20 text-slate-300 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" /></svg>
            <p class="text-xs font-black text-slate-400 uppercase tracking-widest text-center px-6">Pilih salah satu riwayat di sebelah kiri <br> untuk melihat detail audit.</p>
        </div>
    </div>
</template>