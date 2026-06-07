<script setup>
import { useStockOpnameReport } from '../../composables/useStockOpnameReport.js';
import Sidebar from '../../components/Sidebar.vue';
import ReportList from '../../components/stockopname/report/ReportList.vue';
import ReportDetail from '../../components/stockopname/report/ReportDetail.vue';

const { 
    reports, isLoading, selectedDetail, isOwner, isApproving,
    showDetail, calculateLoss, formatDate, approveAudit
} = useStockOpnameReport();
</script>

<template>
    <Sidebar>
        <div id="stock-opname-view" class="p-4 md:p-8 lg:p-10 max-w-7xl mx-auto font-sans bg-[#f8fafc] min-h-screen">
            
            <div class="no-print bg-gradient-to-br from-slate-900 via-slate-800 to-indigo-900 rounded-[40px] p-8 md:p-10 mb-6 text-white shadow-2xl flex flex-col md:flex-row items-center justify-between relative overflow-hidden border border-white/10 select-none">
                <div class="z-10 text-center md:text-left">
                    <h1 class="text-3xl md:text-4xl font-black tracking-tighter mb-2 uppercase italic leading-none">Audit & <span class="text-indigo-400">Claims</span></h1>
                    <p class="text-slate-300 font-bold text-[10px] uppercase tracking-[0.2em] flex items-center justify-center md:justify-start gap-1.5">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-indigo-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" /></svg>
                        Pusat Rekonsiliasi & Persetujuan Stok
                    </p>
                </div>
            </div>

            <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 mt-6">
                
                <div class="lg:col-span-4 no-print space-y-6">
                    
                    <div v-if="isLoading" class="bg-white border border-slate-100 rounded-[32px] p-12 text-center text-slate-400 font-black text-xs uppercase tracking-widest animate-pulse shadow-sm flex flex-col items-center gap-3">
                        <div class="w-8 h-8 border-4 border-indigo-500 border-t-transparent rounded-full animate-spin"></div>
                        <span>Sinkronisasi Data Cloud...</span>
                    </div>
                    
                    <ReportList 
                        v-else
                        :reports="reports" 
                        :selectedDetail="selectedDetail" 
                        :formatDate="formatDate"
                        :calculateLoss="calculateLoss"
                        @select="showDetail" 
                    />
                </div>

                <div class="lg:col-span-8">
                    
                    <ReportDetail 
                        v-if="selectedDetail && !isLoading"
                        :detail="selectedDetail"
                        :isOwner="isOwner"
                        :isApproving="isApproving"
                        :formatDate="formatDate"
                        :calculateLoss="calculateLoss"
                        @approve="(id, type, file) => approveAudit(id, type, file)"
                    />

                    <div v-else class="bg-white border-2 border-dashed border-slate-200 rounded-[40px] p-16 text-center h-full min-h-[400px] flex flex-col items-center justify-center gap-4 no-print transition-all">
                        <div class="w-20 h-20 bg-slate-50 text-slate-300 rounded-full flex items-center justify-center shadow-inner border border-slate-100">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-10 h-10 text-slate-300" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                            </svg>
                        </div>
                        <div class="flex flex-col gap-1">
                            <h3 class="font-black text-slate-700 uppercase tracking-wider text-sm">Review Berita Acara</h3>
                            <p class="text-xs font-bold text-slate-400 max-w-xs mx-auto leading-relaxed">Silakan klik salah satu daftar dokumen audit atau klaim barang temuan di sebelah kiri untuk memvalidasi angka selisih finansial laci!</p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </Sidebar>
</template>

<style scoped>
/* ==========================================
   📈 CORE LAYOUT PRINT RECONCILIATION ENGINE
   ========================================== */
@media print {
    /* Gembok total sirkulasi shadow DOM sidebar pas print view dipicu bray */
    :deep(aside), 
    :deep(nav), 
    :deep(header),
    :deep(#sidebar-wrapper),
    .no-print {
        display: none !important;
    }
    
    #stock-opname-view {
        padding: 0 !important;
        margin: 0 !important;
        max-width: 100% !important;
        width: 100% !important;
        background: #ffffff !important;
    }
}
</style>