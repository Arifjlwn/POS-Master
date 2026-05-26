<script setup>
import { defineProps, defineEmits } from 'vue';

defineProps({ searchQuery: String });
const emit = defineEmits(['update:searchQuery', 'add-new']);
</script>

<template>
    <div>
        <div class="bg-gradient-to-br from-slate-900 via-slate-800 to-indigo-900 rounded-[32px] p-6 md:p-10 mb-6 text-white shadow-2xl flex flex-col md:flex-row items-center justify-between relative overflow-hidden border border-slate-800 gap-6">
            <svg xmlns="http://www.w3.org/2000/svg" class="absolute -right-10 -bottom-10 w-64 h-64 text-indigo-500 opacity-10 pointer-events-none" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><line x1="19" y1="8" x2="19" y2="14"/><line x1="22" y1="11" x2="16" y2="11"/></svg>
            
            <div class="z-10 text-center md:text-left">
                <h1 class="text-3xl md:text-4xl font-black tracking-tighter mb-1 uppercase italic">Human <span class="text-blue-400">Resources</span></h1>
                <p class="text-slate-300 font-bold text-[10px] md:text-xs uppercase tracking-[0.2em] flex items-center justify-center md:justify-start gap-2">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-blue-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M17 21v2a4 4 0 0 1-4 4H5a4 4 0 0 1-4-4v-2a4 4 0 0 1 4-4h1v-1a4 4 0 0 1 4-4h2a4 4 0 0 1 4 4v1h1a4 4 0 0 1 4 4v2z"/><path d="M9 11v-1a2 2 0 0 1 2-2h2a2 2 0 0 1 2 2v1"/></svg>
                    Employee Directory & Access
                </p>
            </div>
            
            <div class="z-10 w-full md:w-auto">
                <button @click="emit('add-new')" class="w-full md:w-auto bg-blue-600 hover:bg-blue-500 text-white px-8 py-4 rounded-[20px] md:rounded-[24px] font-black text-xs uppercase tracking-[0.2em] shadow-xl shadow-blue-900/50 flex items-center justify-center gap-3 transition-all active:scale-95 border border-blue-400/30">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><line x1="19" y1="8" x2="19" y2="14"/><line x1="22" y1="11" x2="16" y2="11"/></svg>
                    Add Employee
                </button>
            </div>
        </div>

        <div class="mb-8 relative group">
            <div class="absolute inset-y-0 left-0 pl-5 flex items-center pointer-events-none">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-400 group-focus-within:text-blue-600 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
            </div>
            <input :value="searchQuery" @input="emit('update:searchQuery', $event.target.value)" type="text" placeholder="Cari Nama, NIK, atau Jabatan..." class="block w-full pl-14 pr-4 py-4 bg-white rounded-2xl border-2 border-slate-100 shadow-sm focus:border-blue-600 outline-none font-bold text-sm transition-all text-slate-700">
        </div>
    </div>
</template>