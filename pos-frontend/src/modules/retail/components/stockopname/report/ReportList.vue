<script setup>
defineProps({ 
    reports: Array, 
    selectedDetail: Object,
    formatDate: Function,
    calculateLoss: Function
});
defineEmits(['select']);
</script>

<template>
    <div class="bg-white rounded-[32px] shadow-xl shadow-slate-200/50 border border-slate-100 overflow-hidden">
        <div class="p-6 bg-slate-50/50 border-b border-slate-100">
            <h3 class="font-black text-slate-800 text-xs uppercase tracking-widest flex items-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                Riwayat Rekonsiliasi (Bulanan)
            </h3>
        </div>
        <div class="overflow-x-auto custom-scrollbar max-h-[600px] overflow-y-auto">
            <table class="w-full text-left">
                <tbody class="divide-y divide-slate-50">
                    <tr v-if="reports.length === 0">
                        <td class="p-10 text-center text-slate-400 font-bold text-xs uppercase tracking-widest italic">Belum ada riwayat Audit</td>
                    </tr>
                    <tr v-for="r in reports" :key="r.so.id" @click="$emit('select', r)" 
                        class="cursor-pointer transition-all group relative"
                        :class="selectedDetail?.so?.id === r.so.id ? 'bg-indigo-50 border-l-4 border-l-indigo-600' : 'hover:bg-slate-50 border-l-4 border-l-transparent'">
                        
                        <td class="p-5">
                            <div class="flex justify-between items-start mb-1">
                                <div class="font-black text-slate-800 text-xs uppercase tracking-tight group-hover:text-indigo-600 transition-colors">
                                    {{ formatDate(r.so.created_at) }}
                                </div>
                                <span v-if="r.so.status === 'PENDING_APPROVAL'" class="w-2 h-2 rounded-full bg-amber-500 animate-pulse" title="SO Menunggu Approval"></span>
                                <span v-else-if="r.so.status === 'APPROVED'" class="w-2 h-2 rounded-full bg-emerald-500" title="SO Disetujui"></span>
                            </div>
                            
                            <div class="text-[10px] text-slate-400 font-bold mt-1 line-clamp-1">{{ r.so.notes }}</div>
                            
                            <div class="mt-3 flex items-center gap-2">
                                <span v-if="calculateLoss(r.so.details, 'SO') < 0" class="bg-red-100 text-red-600 px-2.5 py-1 rounded-md text-[9px] font-black uppercase tracking-widest shadow-sm">
                                    RUGI / MINUS
                                </span>
                                <span v-else-if="calculateLoss(r.so.details, 'SO') > 0" class="bg-emerald-100 text-emerald-600 px-2.5 py-1 rounded-md text-[9px] font-black uppercase tracking-widest shadow-sm">
                                    SURPLUS / PLUS
                                </span>
                                <span v-else class="bg-slate-100 text-slate-600 px-2.5 py-1 rounded-md text-[9px] font-black uppercase tracking-widest shadow-sm">
                                    BALANCE
                                </span>

                                <span v-if="r.klaim" class="bg-amber-100 text-amber-600 px-2.5 py-1 rounded-md text-[9px] font-black uppercase tracking-widest shadow-sm flex items-center gap-1" title="Ada Klaim Barang Ketemu">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-2.5 h-2.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" /></svg>
                                    TERKLAIM
                                </span>
                            </div>
                        </td>
                        <td class="p-5 text-right w-10">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-300 group-hover:text-indigo-600 transition-colors inline-block" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" /></svg>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { height: 6px; width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #94a3b8; }
</style>