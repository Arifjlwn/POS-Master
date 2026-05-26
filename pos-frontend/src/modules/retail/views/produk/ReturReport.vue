<script setup>
import { useReturReport } from '../../composables/useReturReport.js';
import Sidebar from '../../components/Sidebar.vue';
import ReturReportModal from '../../components/return/report/ReturReportModal.vue';

const {
    isLoading, searchQuery, isModalOpen, selectedDocument, user,
    filteredDocuments, openDetail, closeDetail, printDocument
} = useReturReport();
</script>

<template>
    <Sidebar>
        <div class="p-4 md:p-8 lg:p-10 max-w-7xl mx-auto font-sans bg-[#f8fafc] min-h-screen no-print">
            
            <div class="bg-gradient-to-br from-indigo-900 via-slate-800 to-slate-900 rounded-[32px] p-6 md:p-10 mb-8 text-white shadow-2xl flex flex-col md:flex-row items-center justify-between relative overflow-hidden border border-slate-800 gap-6">
                <svg xmlns="http://www.w3.org/2000/svg" class="absolute -right-10 -bottom-10 w-64 h-64 text-indigo-500 opacity-10 pointer-events-none" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/><polyline points="10 9 9 9 8 9"/></svg>
                
                <div class="z-10 text-center md:text-left">
                    <h1 class="text-3xl md:text-4xl font-black tracking-tighter mb-1 uppercase italic">Laporan <span class="text-indigo-400">Waste & Retur</span></h1>
                    <p class="text-slate-300 font-bold text-[10px] md:text-xs uppercase tracking-[0.2em] flex items-center justify-center md:justify-start gap-2">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-indigo-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>
                        Riwayat Berita Acara Pemusnahan Barang
                    </p>
                </div>
            </div>

            <div class="mb-6 relative group">
                <div class="absolute inset-y-0 left-0 pl-5 flex items-center pointer-events-none">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-400 group-focus-within:text-indigo-600 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
                </div>
                <input v-model="searchQuery" type="text" placeholder="Cari Nomor Dokumen (RET-...) atau Nama Barang..." class="block w-full pl-14 pr-4 py-4 bg-white rounded-2xl border-2 border-slate-100 shadow-sm focus:border-indigo-600 outline-none font-bold text-sm transition-all text-slate-700">
            </div>

            <div v-if="isLoading" class="py-20 flex flex-col items-center justify-center bg-white rounded-[32px] border border-slate-100">
                <div class="w-12 h-12 border-4 border-slate-200 border-t-indigo-600 rounded-full animate-spin mb-4"></div>
                <p class="text-slate-400 font-black text-xs uppercase tracking-widest animate-pulse">Menyusun Laporan...</p>
            </div>

            <div v-else-if="filteredDocuments.length === 0" class="flex flex-col items-center justify-center py-20 bg-white/50 rounded-[32px] border-2 border-dashed border-slate-200">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-20 h-20 text-slate-300 mb-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/><polyline points="10 9 9 9 8 9"/></svg>
                <p class="text-slate-400 font-black text-sm uppercase tracking-widest">Dokumen Tidak Ditemukan</p>
            </div>

            <div v-else class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
                <div v-for="doc in filteredDocuments" :key="doc.return_no" class="bg-white rounded-[24px] p-6 shadow-sm border border-slate-100 hover:shadow-xl hover:border-indigo-100 transition-all group flex flex-col justify-between">
                    
                    <div>
                        <div class="flex items-center justify-between mb-4 border-b border-dashed border-slate-100 pb-4">
                            <div class="inline-flex items-center gap-2 bg-indigo-50 text-indigo-700 px-3 py-1.5 rounded-xl text-[10px] font-black uppercase tracking-widest">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
                                {{ doc.return_no }}
                            </div>
                            <span class="text-[9px] font-bold text-slate-400">{{ new Date(doc.created_at).toLocaleDateString('id-ID') }}</span>
                        </div>

                        <div class="flex items-center gap-4 mb-4">
                            <div class="w-12 h-12 rounded-2xl bg-rose-50 flex flex-col items-center justify-center border border-rose-100">
                                <span class="text-xs font-black text-rose-600 leading-none">{{ doc.items.length }}</span>
                                <span class="text-[8px] font-black text-rose-400 uppercase tracking-widest mt-1">Item</span>
                            </div>
                            <div class="flex-1">
                                <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest">Total Dibuang</p>
                                <p class="text-base font-black text-slate-800">{{ doc.total_qty }} <span class="text-[10px] text-slate-500 uppercase">Unit</span></p>
                            </div>
                        </div>

                        <div class="flex -space-x-2 overflow-hidden py-1 mb-4">
                            <span v-for="(item, idx) in doc.items.slice(0, 3)" :key="idx" class="inline-block px-2 py-1 bg-slate-100 text-slate-600 text-[9px] font-bold rounded-lg border border-white whitespace-nowrap truncate max-w-[120px]">
                                {{ item.product?.nama_produk || 'Item' }}
                            </span>
                            <span v-if="doc.items.length > 3" class="inline-block px-2 py-1 bg-slate-200 text-slate-600 text-[9px] font-bold rounded-lg border border-white">
                                +{{ doc.items.length - 3 }} Lainnya
                            </span>
                        </div>
                    </div>

                    <div class="flex items-center justify-between pt-4 border-t border-slate-100 mt-2">
                        <div class="flex items-center gap-2">
                            <div class="w-6 h-6 rounded-full bg-slate-800 text-white flex items-center justify-center text-[8px] font-black">{{ doc.user.name.substring(0,2).toUpperCase() }}</div>
                            <span class="text-[10px] font-bold text-slate-500 truncate max-w-[80px]">{{ doc.user.name }}</span>
                        </div>
                        <button @click="openDetail(doc)" class="bg-slate-100 hover:bg-indigo-600 text-slate-600 hover:text-white px-4 py-2 rounded-xl text-[10px] font-black uppercase tracking-widest transition-colors flex items-center gap-2">
                            Detail
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M5 12h14"/><path d="m12 5 7 7-7 7"/></svg>
                        </button>
                    </div>

                </div>
            </div>
        </div>

        <ReturReportModal 
            :isOpen="isModalOpen"
            :document="selectedDocument"
            :storeName="user.storeName"
            @close="closeDetail"
            @print="printDocument"
        />

    </Sidebar>
</template>

<style>
/* CSS Print Global (Taruh tanpa 'scoped' atau di CSS Global) */
@media print {
    body * { visibility: hidden !important; }
    #printable-area, #printable-area * { visibility: visible !important; }
    #printable-area { position: absolute !important; left: 0 !important; top: 0 !important; width: 100% !important; margin: 0 !important; padding: 0 !important; }
    body { background-color: white !important; }
    .no-print, header, nav, .sidebar { display: none !important; }
    table { page-break-inside: auto; }
    tr    { page-break-inside: avoid; page-break-after: auto; }
    @page { margin: 20mm; size: A4 portrait; }
}
</style>