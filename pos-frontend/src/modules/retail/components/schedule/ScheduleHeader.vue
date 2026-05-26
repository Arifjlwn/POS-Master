<script setup>
import { defineProps, defineEmits } from 'vue';

defineProps({
    startDate: String,
    endDate: String,
    canEditSchedule: Boolean,
    isSaving: Boolean,
    searchQuery: String
});

const emit = defineEmits(['update:searchQuery', 'save-jadwal']);
</script>

<template>
    <div class="mb-8 relative group">
        <div class="bg-gradient-to-br from-slate-900 via-slate-800 to-indigo-900 rounded-[32px] p-6 md:p-10 mb-6 text-white shadow-2xl flex flex-col md:flex-row items-center justify-between relative overflow-hidden border border-slate-800 gap-6">
            <svg xmlns="http://www.w3.org/2000/svg" class="absolute -right-10 -bottom-10 w-64 h-64 text-indigo-500 opacity-10 pointer-events-none" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/><path d="M8 14h.01"/><path d="M12 14h.01"/><path d="M16 14h.01"/><path d="M8 18h.01"/><path d="M12 18h.01"/><path d="M16 18h.01"/></svg>
            
            <div class="z-10 text-center md:text-left">
                <h1 class="text-3xl md:text-4xl font-black tracking-tighter mb-1 uppercase italic">Shift <span class="text-blue-400">Management</span></h1>
                <p class="text-slate-300 font-bold text-[10px] md:text-xs uppercase tracking-[0.2em] flex items-center justify-center md:justify-start gap-2">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-blue-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
                    Periode: {{ startDate }} <span class="opacity-50">s/d</span> {{ endDate }}
                </p>
            </div>
            
            <div v-if="canEditSchedule" class="z-10 w-full md:w-auto">
                <button
                    @click="emit('save-jadwal')" 
                    :disabled="isSaving"
                    class="w-full md:w-auto bg-blue-600 hover:bg-blue-500 text-white px-8 py-4 rounded-[20px] md:rounded-[24px] font-black text-xs uppercase tracking-[0.2em] shadow-xl shadow-blue-900/50 flex items-center justify-center gap-3 transition-all active:scale-95 border border-blue-400/30 disabled:opacity-70 disabled:cursor-not-allowed"
                >
                    <template v-if="isSaving">
                        <div class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                        Menyimpan...
                    </template>
                    <template v-else>
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/><polyline points="17 21 17 13 7 13 7 21"/><polyline points="7 3 7 8 15 8"/></svg>
                        Simpan Jadwal
                    </template>
                </button>
            </div>
        </div>

        <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-5 flex items-center pointer-events-none">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-400 group-focus-within:text-indigo-600 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
            </div>
            <input 
                :value="searchQuery" 
                @input="emit('update:searchQuery', $event.target.value)" 
                type="text" 
                placeholder="Cari Nama Karyawan atau ID..." 
                class="block w-full pl-14 pr-4 py-4 bg-white rounded-2xl border-2 border-slate-100 shadow-sm focus:border-indigo-600 outline-none font-bold text-sm transition-all text-slate-700"
            >
        </div>
    </div>
</template>