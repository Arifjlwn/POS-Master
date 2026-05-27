<script setup>
import { useStockOpnameReport } from '../../composables/useStockOpnameReport.js';
import Sidebar from '../../components/Sidebar.vue';
import ReportList from '../../components/stockopname/report/ReportList.vue';
import ReportDetail from '../../components/stockopname/report/ReportDetail.vue';

// 🚀 Inject reportType ke view
const { 
    reportType, reports, isLoading, selectedDetail, isOwner, isApproving,
    showDetail, calculateLoss, formatDate, approveAudit
} = useStockOpnameReport();
</script>

<template>
    <Sidebar>
        <div class="p-4 md:p-8 lg:p-10 max-w-7xl mx-auto font-sans bg-[#f8fafc] min-h-screen">
            
            <div class="no-print bg-gradient-to-br from-slate-900 via-slate-800 to-indigo-900 rounded-[40px] p-8 md:p-10 mb-6 text-white shadow-2xl flex flex-col md:flex-row items-center justify-between relative overflow-hidden border border-white/10">
                <div class="z-10 text-center md:text-left">
                    <h1 class="text-3xl md:text-4xl font-black tracking-tighter mb-2 uppercase italic leading-none">Audit & <span class="text-indigo-400">Claims</span></h1>
                    <p class="text-slate-300 font-bold text-[10px] uppercase tracking-[0.2em]">Pusat Persetujuan Stok & Barang Temuan</p>
                </div>
            </div>

            <div class="no-print flex justify-start mb-6">
                <div class="bg-white p-1.5 rounded-2xl shadow-sm border border-slate-100 flex gap-2">
                    <button @click="reportType = 'SO'" :class="reportType === 'SO' ? 'bg-indigo-600 text-white shadow-md' : 'text-slate-500'" class="px-5 py-2 rounded-xl font-black text-[10px] uppercase tracking-widest transition-all">
                        Laporan Audit Reguler
                    </button>
                    <button @click="reportType = 'KLAIM'" :class="reportType === 'KLAIM' ? 'bg-amber-500 text-white shadow-md' : 'text-slate-500'" class="px-5 py-2 rounded-xl font-black text-[10px] uppercase tracking-widest transition-all">
                        Persetujuan Klaim Nyempil
                    </button>
                </div>
            </div>

            <div class="grid grid-cols-1 lg:grid-cols-12 gap-8">
                <div class="lg:col-span-5 no-print space-y-6">
                    <div v-if="isLoading" class="p-10 text-center text-slate-400 font-bold text-xs animate-pulse">Memuat Data...</div>
                    <ReportList 
                        v-else
                        :reports="reports" 
                        :selectedDetail="selectedDetail" 
                        :formatDate="formatDate"
                        :calculateLoss="calculateLoss"
                        @select="showDetail" 
                    />
                </div>

                <div class="lg:col-span-7">
                    <ReportDetail 
                        :detail="selectedDetail"
                        :isOwner="isOwner"
                        :isApproving="isApproving"
                        :formatDate="formatDate"
                        :calculateLoss="calculateLoss"
                        :reportType="reportType" @approve="approveAudit"
                    />
                </div>
            </div>
        </div>
    </Sidebar>
</template>