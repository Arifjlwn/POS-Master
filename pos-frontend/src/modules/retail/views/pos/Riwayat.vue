<script setup>
import { useJournal } from '../../composables/useJournal.js';
import Sidebar from '../../components/Sidebar.vue';
// 🚀 PERUBAHAN 1: Import file komponen ReceiptModal lu di sini beb!
import ReceiptModal from '../../components/pos/ReceiptModal.vue';

const {
    isLoading,
    tanggalDipilih,
    searchQuery,
    showReceipt,
    selectedTrx,
    filteredRiwayat,
    formatRupiah,
    openReceipt,
    currentUser,
    currentSession
} = useJournal();
</script>

<template>
    <Sidebar class="print:bg-white print:h-auto ">
        <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 font-sans bg-[#f8fafc] min-h-screen print:hidden">
            
            <div class="bg-gradient-to-br from-slate-900 via-slate-800 to-blue-900 rounded-[32px] p-6 md:p-10 mb-6 text-white shadow-2xl flex flex-col md:flex-row md:items-center justify-between relative overflow-hidden gap-6">
                <div class="z-10 text-center md:text-left">
                    <h1 class="text-3xl md:text-4xl font-black tracking-tighter mb-1 uppercase italic">Daily <span class="text-blue-400">Journal</span></h1>
                    <p class="text-slate-300 font-bold text-[10px] md:text-xs uppercase tracking-[0.2em]">Audit, Reprint & Track Store Sales</p>
                </div>

                <div class="z-10 flex items-center justify-center gap-3 bg-white/10 backdrop-blur-md p-2 pl-4 rounded-2xl border border-white/20 shadow-inner w-full md:w-auto">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-blue-300" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
                    <input type="date" v-model="tanggalDipilih" class="bg-transparent border-none text-sm md:text-base font-black text-white focus:ring-0 cursor-pointer outline-none uppercase tracking-tighter w-full md:w-auto [color-scheme:dark]">
                </div>
            </div>

            <div class="mb-6 relative group">
                <div class="absolute inset-y-0 left-0 pl-5 flex items-center pointer-events-none">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-400 group-focus-within:text-blue-600 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
                </div>
                <input v-model="searchQuery" type="text" placeholder="Cari Nomor Invoice atau Nama Kasir..." class="block w-full pl-14 pr-4 py-4 bg-white rounded-2xl border-2 border-slate-100 shadow-sm focus:border-blue-600 outline-none font-bold text-sm transition-all text-slate-700">
            </div>

            <div class="md:hidden grid grid-cols-1 gap-4">
                <div v-for="trx in filteredRiwayat" :key="trx.id" class="bg-white p-5 rounded-[24px] shadow-sm border border-slate-100 flex flex-col gap-3 relative overflow-hidden">
                    <div class="absolute left-0 top-0 bottom-0 w-1.5" :class="trx.metode_bayar === 'Cash' ? 'bg-emerald-400' : 'bg-blue-500'"></div>
                    
                    <div class="flex justify-between items-start">
                        <div>
                            <div class="font-mono font-black text-sm text-slate-800 tracking-tighter">{{ trx.no_invoice }}</div>
                            <div class="text-[10px] font-bold text-slate-400 uppercase mt-0.5">
                                {{ new Date(trx.created_at).toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' }) }} • OPR: {{ trx.User?.name ? trx.User.name.split(' ')[0] : 'KASIR' }}
                            </div>
                        </div>
                        <span class="px-2.5 py-1 rounded-lg text-[9px] font-black uppercase tracking-widest border" :class="trx.metode_bayar === 'Cash' ? 'border-emerald-100 bg-emerald-50 text-emerald-600' : 'border-blue-100 bg-blue-50 text-blue-600'">
                            {{ trx.metode_bayar || 'CASH' }}
                        </span>
                    </div>

                    <div class="border-t border-dashed border-slate-200 my-1"></div>

                    <div class="flex justify-between items-end">
                        <div>
                            <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest">Total Transaksi</p>
                            <p class="font-black text-slate-800 text-lg">Rp{{ formatRupiah(trx.total_harga) }}</p>
                        </div>
                        <button @click="openReceipt(trx)" class="flex items-center justify-center w-10 h-10 bg-slate-100 text-slate-600 rounded-xl hover:bg-slate-900 hover:text-white transition-colors active:scale-95 shadow-sm">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M6 9V2h12v7"/><path d="M6 18H4a2 2 0 0 1-2-2v-5a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v5a2 2 0 0 1-2 2h-2"/><rect width="12" height="8" x="6" y="14"/></svg>
                        </button>
                    </div>
                </div>
            </div>

            <div v-if="filteredRiwayat.length > 0" class="hidden md:block w-full bg-white rounded-[32px] shadow-sm border border-slate-100 overflow-hidden transition-all">
                <table class="w-full text-left border-collapse whitespace-nowrap">
                    <thead>
                        <tr class="bg-slate-50 border-b border-slate-100">
                            <th class="px-6 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">No. Invoice</th>
                            <th class="px-6 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">Waktu</th>
                            <th class="px-6 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">Kasir</th>
                            <th class="px-6 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">Metode</th>
                            <th class="px-6 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest text-right">Total Nominal</th>
                            <th class="px-6 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest text-center">Aksi</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-slate-50">
                        <tr v-for="trx in filteredRiwayat" :key="trx.id" class="hover:bg-blue-50/40 transition-all group">
                            <td class="px-6 py-4 font-mono font-black text-sm text-slate-700 tracking-tighter">{{ trx.no_invoice }}</td>
                            <td class="px-6 py-4 text-xs font-bold text-slate-500 uppercase">
                                {{ new Date(trx.created_at).toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' }) }}
                            </td>
                            <td class="px-6 py-4 text-xs font-black text-slate-700 uppercase">{{ trx.User?.name || 'KASIR' }}</td>
                            <td class="px-6 py-4">
                                <span class="px-3 py-1.5 rounded-lg text-[10px] font-black uppercase tracking-widest border" :class="trx.metode_bayar === 'Cash' ? 'border-emerald-100 bg-emerald-50 text-emerald-600' : 'border-blue-100 bg-blue-50 text-blue-600'">
                                    {{ trx.metode_bayar || 'CASH' }}
                                </span>
                            </td>
                            <td class="px-6 py-4 text-right font-black text-slate-800 text-sm">Rp{{ formatRupiah(trx.total_harga) }}</td>
                            <td class="px-6 py-4 text-center">
                                <button @click="openReceipt(trx)" class="inline-flex items-center justify-center gap-2 bg-slate-50 text-slate-500 border border-slate-200 px-4 py-2.5 rounded-xl text-[10px] font-black uppercase tracking-widest hover:bg-slate-900 hover:text-white hover:border-slate-900 transition-all active:scale-95 shadow-sm">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M6 9V2h12v7"/><path d="M6 18H4a2 2 0 0 1-2-2v-5a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v5a2 2 0 0 1-2 2h-2"/><rect width="12" height="8" x="6" y="14"/></svg>
                                    Cetak Struk
                                </button>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </main>
    </Sidebar>
    <ReceiptModal 
        :show="showReceipt" 
        :invoiceData="selectedTrx" 
        :storeData="selectedTrx?.Store || currentSession?.store || currentSession?.Store"
        :cashierName="selectedTrx?.User?.name ? selectedTrx.User.name.split(' ')[0] : (currentUser?.name ? currentUser.name.split(' ')[0] : 'KASIR')"
        :stationNumber="'01'"
        @close="showReceipt = false" 
    />
</template>

<style scoped>
/* 🚀 TAMBAHKAN CSS INI BIAR LAYAR BERSIH PAS NGEPRINT */
@media print {
    @page { 
        margin: 0; 
    }
    body { 
        background: white; 
        -webkit-print-color-adjust: exact; 
    }
}

.custom-scrollbar::-webkit-scrollbar { width: 5px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 20px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #94a3b8; }

input[type="date"]::-webkit-calendar-picker-indicator {
    filter: invert(1);
    cursor: pointer;
    opacity: 0.8;
}
input[type="date"]::-webkit-calendar-picker-indicator:hover {
    opacity: 1;
}
</style>