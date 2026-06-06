<script setup>
import { computed, ref } from 'vue';

const props = defineProps({
    detail: Object,
    isOwner: Boolean,
    isApproving: Boolean,
    formatDate: Function,
    calculateLoss: Function,
    reportType: { type: String, default: 'SO' }
});

const emit = defineEmits(['print', 'approve']);

// Ambil BASE URL secara dinamis bray, biar ga hardcoded localhost terus pas production
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/';

const formatNumber = (val) => {
    if (val === null || val === undefined) return '0';
    return Number(val).toLocaleString('id-ID');
};

const soLoss = computed(() => {
    return props.detail?.so ? props.calculateLoss(props.detail.so.details, 'SO') : 0;
});

const klaimRecovered = computed(() => {
    return props.detail?.klaim ? props.calculateLoss(props.detail.klaim.details, 'KLAIM') : 0;
});

const finalNetto = computed(() => {
    return soLoss.value + klaimRecovered.value;
});

const isSettled = computed(() => {
    if (props.detail?.klaim) return props.detail.klaim.status === 'APPROVED';
    return props.detail?.so?.status === 'APPROVED';
});

const triggerPrint = () => {
    window.print();
};

const getMinusQty = (productId) => {
    const soItem = props.detail?.so?.details?.find(item => item.product_id === productId);
    return soItem ? soItem.selisih : 0;
};

// STATE UNTUK FILE UPLOAD BERITA ACARA (PDF)
const pdfFile = ref(null);
const pdfFileName = ref(''); // State tambahan buat indikator UX nama file bray

const handleFileUpload = (e) => {
    const file = e.target.files[0];
    if (file) {
        pdfFile.value = file;
        pdfFileName.value = file.name;
    }
};

// FIX LOGIC EMIT: Pastikan dibungkus dengan objek atau dipastikan urutannya pas
const handleApprove = (id, type) => {
    emit('approve', id, type, pdfFile.value);
    // Reset form sesudah emit berhasil dilempar
    pdfFile.value = null; 
    pdfFileName.value = '';
};

const pdfUrl = computed(() => {
    if (props.detail?.klaim?.bukti_bar) return props.detail.klaim.bukti_bar;
    if (props.detail?.so?.bukti_bar) return props.detail.so.bukti_bar;
    return null;
});

const isPdfModalOpen = ref(false);
</script>

<template>
    <div>
        <div v-if="detail && detail.so" id="print-area" class="bg-white rounded-[32px] print:rounded-none shadow-xl shadow-slate-200/50 print:shadow-none p-6 md:p-8 border-2 border-indigo-100 print:border-none relative print:p-0 print:w-full overflow-hidden">

            <div class="relative z-10">
                <div v-if="isSettled" class="watermark-print">
                    APPROVED
                </div>
                
                <div class="flex flex-col sm:flex-row sm:items-start justify-between gap-4 border-b-2 border-slate-800 print:border-slate-300 pb-6 mb-6">
                    <div class="text-center sm:text-left print:w-full print:text-center">
                        <h1 class="text-2xl print:text-xl font-black text-slate-900 uppercase tracking-[0.2em] leading-tight">
                            Laporan Rekonsiliasi Audit
                        </h1>
                        <p class="text-xs print:text-[10px] font-bold text-slate-500 mt-2 uppercase tracking-widest">{{ formatDate(detail.so.created_at) }}</p>
                        
                        <div class="mt-3 flex gap-2 flex-wrap justify-center sm:justify-start print:justify-center">
                            <span class="inline-flex items-center gap-1.5 px-3 py-1 print:px-2 print:py-0.5 rounded-lg text-[10px] font-black uppercase tracking-widest border"
                                :class="detail.so.status === 'APPROVED' ? 'bg-emerald-100 text-emerald-700 border-emerald-200' : 'bg-amber-100 text-amber-700 border-amber-200'">
                                SO: {{ detail.so.status === 'APPROVED' ? 'Disetujui' : 'Menunggu' }}
                            </span>
                            
                            <span v-if="detail.klaim" class="inline-flex items-center gap-1.5 px-3 py-1 print:px-2 print:py-0.5 rounded-lg text-[10px] font-black uppercase tracking-widest border"
                                :class="detail.klaim.status === 'APPROVED' ? 'bg-emerald-100 text-emerald-700 border-emerald-200' : 'bg-amber-100 text-amber-700 border-amber-200'">
                                KLAIM: {{ detail.klaim.status === 'APPROVED' ? 'Disetujui' : 'Menunggu' }}
                            </span>
                        </div>
                    </div>

                    <div class="no-print flex flex-col gap-2 shrink-0 w-full sm:w-auto">
                        <button @click="triggerPrint" class="w-full p-2.5 px-4 bg-slate-800 hover:bg-slate-900 text-white rounded-xl transition-all shadow-md text-[10px] font-black uppercase tracking-widest flex items-center justify-center gap-2">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z" /></svg>
                            CETAK LAPORAN
                        </button>

                        <button v-if="pdfUrl" 
                           @click="isPdfModalOpen = true"
                           class="w-full p-2.5 px-4 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl transition-all shadow-md text-[10px] font-black uppercase tracking-widest flex items-center justify-center gap-2">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" /><path stroke-linecap="round" stroke-linejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" /></svg>
                            PREVIEW BUKTI (PDF)
                        </button>

                        <div v-if="isOwner && !isSettled" class="bg-indigo-50 border border-indigo-100 p-3 rounded-xl flex flex-col gap-2 mt-2">
                            <label class="text-[9px] font-black text-indigo-800 uppercase tracking-widest text-center">Upload PDF TTD (Berita Acara)</label>
                            <input type="file" accept=".pdf" @change="handleFileUpload" 
                                class="text-[9px] font-bold text-slate-500 file:mr-2 file:py-1.5 file:px-3 file:rounded-lg file:border-0 file:text-[9px] file:font-black file:bg-indigo-600 file:text-white hover:file:bg-indigo-700 cursor-pointer w-full bg-white rounded-lg border border-slate-200">
                            <p v-if="pdfFileName" class="text-[9px] text-indigo-600 font-bold text-center truncate max-w-[180px]">📂 {{ pdfFileName }}</p>
                        </div>

                        <button v-if="isOwner && detail.so.status === 'PENDING_APPROVAL'" 
                                @click="handleApprove(detail.so.id, 'SO')" 
                                :disabled="isApproving || !pdfFile"
                                :class="!pdfFile ? 'bg-slate-300 text-slate-500 cursor-not-allowed shadow-none' : 'bg-indigo-500 hover:bg-indigo-600 text-white shadow-md'"
                                class="w-full p-2.5 px-4 rounded-xl transition-all text-[10px] font-black uppercase tracking-widest mt-1">
                            {{ !pdfFile ? '🔒 UPLOAD PDF DULU' : 'SETUJUI AUDIT AWAL' }}
                        </button>
                        
                        <button v-if="isOwner && detail.klaim && detail.klaim.status === 'PENDING_APPROVAL'" 
                                @click="handleApprove(detail.klaim.id, 'KLAIM')" 
                                :disabled="isApproving || !pdfFile"
                                :class="!pdfFile ? 'bg-slate-300 text-slate-500 cursor-not-allowed shadow-none' : 'bg-amber-500 hover:bg-amber-600 text-white shadow-md'"
                                class="w-full p-2.5 px-4 rounded-xl transition-all text-[10px] font-black uppercase tracking-widest mt-1">
                            {{ !pdfFile ? '🔒 UPLOAD PDF DULU' : 'SETUJUI KLAIM TEMUAN' }}
                        </button>
                    </div>
                </div>

                <div class="mb-8 p-6 print:p-4 rounded-2xl print:rounded-xl border-2 border-solid print:border-black" :class="finalNetto < 0 ? 'bg-red-50 border-red-200' : 'bg-emerald-50 border-emerald-200'">
                    <div class="flex flex-col md:flex-row print:flex-row justify-between items-center print:items-center gap-4">
                        <div class="text-center md:text-left print:text-left">
                            <p class="text-[10px] font-black uppercase tracking-widest" :class="finalNetto < 0 ? 'text-red-500' : 'text-emerald-600'">Rugi SO Awal : - Rp {{ formatNumber(Math.abs(soLoss)) }}</p>
                            <p class="text-[10px] font-black uppercase tracking-widest mt-1 text-amber-600">Klaim : + Rp {{ formatNumber(klaimRecovered) }}</p>
                        </div>
                        <div class="text-center md:text-right print:text-right">
                            <h3 class="text-[10px] font-black text-slate-500 uppercase tracking-[0.2em] mb-1">Final Netto Loss</h3>
                            <div class="text-3xl print:text-2xl font-black" :class="finalNetto < 0 ? 'text-red-600' : 'text-emerald-600'">
                                {{ finalNetto < 0 ? '-' : '' }} Rp {{ formatNumber(Math.abs(finalNetto)) }}
                            </div>
                        </div>
                    </div>
                </div>

                <h4 class="font-black text-slate-800 text-[10px] uppercase tracking-[0.2em] mb-3 flex items-center gap-2">
                    <span class="w-1.5 h-4 bg-indigo-600 rounded-full"></span> 1. Rincian Selisih Audit (SO)
                </h4>
                <div class="overflow-x-auto mb-8 print:mb-6 border border-slate-200 print:border-none rounded-xl print:rounded-none">
                    <table class="w-full text-left border-collapse min-w-[600px] print:min-w-full print-table">
                        <thead class="bg-slate-50 border-b-2 border-slate-200 print:border-slate-300">
                            <tr class="text-[9px] font-black text-slate-500 uppercase tracking-widest">
                                <th class="p-3 print:py-2 print:px-3 w-[40%] border-r border-transparent print:border-slate-200">Nama Barang</th>
                                <th class="p-3 print:py-2 print:px-3 text-center border-r border-transparent print:border-slate-200">Stok Sistem</th>
                                <th class="p-3 print:py-2 print:px-3 text-center border-r border-transparent print:border-slate-200">Hasil Fisik</th>
                                <th class="p-3 print:py-2 print:px-3 text-center">Selisih</th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-slate-100 print:divide-slate-200 border-b border-slate-200">
                            <tr v-for="d in detail.so.details" :key="d.id" class="text-xs font-bold text-slate-700">
                                <td class="p-3 print:py-2 print:px-3 border-r border-transparent print:border-slate-200">
                                    <div class="uppercase">{{ d.product?.nama_produk || 'Produk Terhapus' }}</div>
                                    <div v-if="d.alasan" class="text-[9px] text-amber-600 mt-0.5">Note: {{ d.alasan }}</div>
                                </td>
                                <td class="p-3 print:py-2 print:px-3 text-center text-slate-500 border-r border-transparent print:border-slate-200">
                                    {{ formatNumber(d.system_qty) }} <span class="text-[8px] uppercase ml-0.5">{{ d.product?.satuan_dasar || 'PCS' }}</span>
                                </td>
                                <td class="p-3 print:py-2 print:px-3 text-center text-indigo-600 border-r border-transparent print:border-slate-200">
                                    {{ formatNumber(d.actual_qty) }} <span class="text-[8px] uppercase ml-0.5">{{ d.product?.satuan_dasar || 'PCS' }}</span>
                                </td>
                                <td class="p-3 print:py-2 print:px-3 text-center" :class="d.selisih < 0 ? 'text-red-600' : d.selisih > 0 ? 'text-emerald-600' : 'text-slate-500'">
                                    {{ d.selisih > 0 ? '+' : '' }}{{ formatNumber(d.selisih) }} <span class="text-[8px] uppercase ml-0.5">{{ d.product?.satuan_dasar || 'PCS' }}</span>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>

                <template v-if="detail.klaim">
                    <h4 class="font-black text-slate-800 text-[10px] uppercase tracking-[0.2em] mb-3 flex items-center gap-2">
                        <span class="w-1.5 h-4 bg-amber-500 rounded-full"></span> 2. Rincian Klaim Penemuan Barang
                    </h4>
                    <div class="overflow-x-auto border border-amber-200 print:border-none rounded-xl print:rounded-none">
                        <table class="w-full text-left border-collapse min-w-[700px] print:min-w-full print-table">
                            <thead class="bg-amber-50 border-b-2 border-amber-200 print:border-slate-300">
                                <tr class="text-[9px] font-black text-amber-700 uppercase tracking-widest">
                                    <th class="p-3 print:py-2 print:px-3 w-[30%] border-r border-transparent print:border-slate-200">Nama Barang</th>
                                    <th class="p-3 print:py-2 print:px-3 text-center border-r border-transparent print:border-slate-200">Qty Minus (SO)</th>
                                    <th class="p-3 print:py-2 print:px-3 text-center border-r border-transparent print:border-slate-200">Klaim Ketemu</th>
                                    <th class="p-3 print:py-2 print:px-3 text-center border-r border-transparent print:border-slate-200">Selisih Akhir</th>
                                    <th class="p-3 print:py-2 print:px-3 text-center">Lokasi / Alasan Ditemukan</th>
                                </tr>
                            </thead>
                            <tbody class="divide-y divide-amber-100 print:divide-slate-200 border-b border-slate-200">
                                <tr v-for="d in detail.klaim.details" :key="d.id" class="text-xs font-bold text-slate-700">
                                    <td class="p-3 print:py-2 print:px-3 uppercase border-r border-transparent print:border-slate-200">{{ d.product?.nama_produk || 'Produk Terhapus' }}</td>
                                    <td class="p-3 print:py-2 print:px-3 text-center text-red-600 border-r border-transparent print:border-slate-200">
                                        {{ formatNumber(getMinusQty(d.product_id)) }} <span class="text-[8px] uppercase ml-0.5">{{ d.product?.satuan_dasar || 'PCS' }}</span>
                                    </td>
                                    <td class="p-3 print:py-2 print:px-3 text-center text-emerald-600 font-black border-r border-transparent print:border-slate-200">
                                        +{{ formatNumber(d.qty) }} <span class="text-[8px] uppercase ml-0.5">{{ d.product?.satuan_dasar || 'PCS' }}</span>
                                    </td>
                                    <td class="p-3 print:py-2 print:px-3 text-center font-black border-r border-transparent print:border-slate-200"
                                        :class="(getMinusQty(d.product_id) + d.qty) < 0 ? 'text-red-600' : 'text-slate-500'">
                                        {{ formatNumber(getMinusQty(d.product_id) + d.qty) }} <span class="text-[8px] uppercase ml-0.5">{{ d.product?.satuan_dasar || 'PCS' }}</span>
                                    </td>
                                    <td class="p-3 print:py-2 print:px-3 text-center text-slate-500 text-[10px]">{{ d.alasan || '-' }}</td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </template>

                <div v-if="!isSettled" class="print-avoid-break">
                    <h4 class="font-black text-slate-800 print:text-black text-[10px] uppercase tracking-[0.2em] mb-4 mt-10 flex items-center gap-2">
                        <span class="w-1.5 h-4 bg-slate-800 print:bg-black rounded-full"></span> 3. Otorisasi & Pengesahan Laporan
                    </h4>
                    
                    <div class="grid grid-cols-2 md:grid-cols-3 gap-6 text-center border-t-2 border-solid border-slate-300 print:border-black pt-8 mt-2 pb-4">
                        <div class="flex flex-col items-center justify-between min-h-[100px]">
                            <span class="text-[10px] font-black uppercase tracking-widest text-slate-500 print:text-black">Dihitung Oleh,</span>
                            <div class="w-32 border-b-2 border-slate-800 print:border-black mt-16"></div>
                            <span class="text-[9px] font-bold text-slate-400 print:text-black mt-1.5">Tim Audit / Kasir</span>
                        </div>
                        
                        <div class="flex flex-col items-center justify-between min-h-[100px]">
                            <span class="text-[10px] font-black uppercase tracking-widest text-slate-500 print:text-black">Diperiksa Oleh,</span>
                            <div class="w-32 border-b-2 border-slate-800 print:border-black mt-16"></div>
                            <span class="text-[9px] font-bold text-slate-400 print:text-black mt-1.5">Kepala Toko</span>
                        </div>
                        
                        <div class="flex flex-col items-center justify-between min-h-[100px] col-span-2 md:col-span-1">
                            <span class="text-[10px] font-black uppercase tracking-widest text-slate-500 print:text-black">Disetujui Oleh,</span>
                            <div class="w-32 border-b-2 border-slate-800 print:border-black mt-16"></div>
                            <span class="text-[9px] font-bold text-slate-400 print:text-black mt-1.5">Manajemen / Owner</span>
                        </div>
                    </div>
                </div>

            </div>
        </div>
        
        <div v-else class="no-print flex flex-col items-center justify-center h-full min-h-[400px] bg-slate-50 rounded-[32px] border-2 border-dashed border-slate-200 opacity-60">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-20 h-20 text-slate-300 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" /></svg>
            <p class="text-xs font-black text-slate-400 uppercase tracking-widest text-center px-6">Pilih salah satu riwayat di sebelah kiri <br> untuk melihat detail rekonsiliasi audit.</p>
        </div>
    </div>

    <div v-if="isPdfModalOpen" class="fixed inset-0 z-[9999] flex items-center justify-center bg-slate-900/80 backdrop-blur-sm p-4 md:p-8 no-print">
        <div class="bg-white rounded-3xl shadow-2xl w-full max-w-5xl h-full max-h-[90vh] flex flex-col overflow-hidden relative animate-fade-in-up">
            
            <div class="flex justify-between items-center p-4 md:px-6 border-b-2 border-slate-200 bg-slate-50/50">
                <div class="flex items-center gap-3">
                    <div class="p-2 bg-indigo-100 text-indigo-600 rounded-lg">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" /></svg>
                    </div>
                    <div>
                        <h3 class="font-black text-slate-800 uppercase tracking-widest text-sm">Dokumen Berita Acara</h3>
                        <p class="text-[9px] font-bold text-slate-400 mt-0.5">Bukti fisik rekonsiliasi yang sah</p>
                    </div>
                </div>
                
                <div class="flex items-center gap-2">
                    <a :href="API_BASE_URL + pdfUrl" download="Berita_Acara_Audit.pdf" target="_blank" class="px-4 py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white rounded-xl text-[10px] font-black uppercase tracking-widest flex items-center gap-2 transition-all shadow-md">
                        <svg xmlns="http://www.w3.org/2000/xl" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" /></svg>
                        UNDUH
                    </a>
                    <button @click="isPdfModalOpen = false" class="p-2.5 bg-slate-200 hover:bg-rose-100 text-slate-500 hover:text-rose-600 rounded-xl transition-all">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
                    </button>
                </div>
            </div>

            <div class="flex-1 w-full bg-slate-300 relative">
                <div class="absolute inset-0 flex items-center justify-center -z-10 text-slate-400 font-bold text-xs">
                    Memuat dokumen...
                </div>
                <iframe :src="API_BASE_URL + pdfUrl" class="w-full h-full border-0 relative z-10 shadow-inner"></iframe>
            </div>
            
        </div>
    </div>
</template>