<script setup>
import { onMounted } from 'vue';
import SidebarLaundry from '../../components/SidebarLaundry.vue';
import CatalogFormCard from '../../components/produk/CatalogFormCard.vue';
import CatalogListDisplay from '../../components/produk/CatalogListDisplay.vue';
import { useLaundryCatalog } from '../../composables/useLaundryCatalog.js';

const {
    activeTab, isLoading, searchQuery, isEditing, showForm,
    formJasa, formParfum, filteredItems, formatRupiah,
    loadAllData, switchTab, triggerEdit, cancelForm, handleSave, handleConfirmDelete
} = useLaundryCatalog();

onMounted(() => {
    loadAllData();
});
</script>

<template>
    <SidebarLaundry>
        <div class="flex-1 flex flex-col h-full bg-[#F8FAFC] overflow-hidden relative">
            <div class="p-4 sm:p-6 md:p-8 shrink-0 bg-white border-b border-slate-100 flex flex-col z-10 shadow-sm relative gap-5">
                <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                    <div class="flex items-center gap-3 md:gap-4">
                        <div class="w-10 h-10 md:w-12 md:h-12 bg-indigo-50 border border-indigo-100 rounded-2xl flex items-center justify-center shrink-0">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 md:w-6 md:h-6 text-indigo-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/><polyline points="3.27 6.96 12 12.01 20.73 6.96"/><line x1="12" y1="22.08" x2="12" y2="12"/></svg>
                        </div>
                        <div>
                            <h1 class="text-lg md:text-2xl font-black tracking-tighter uppercase text-slate-800 leading-tight">Master Katalog</h1>
                            <p class="text-[9px] md:text-[10px] font-black text-slate-400 mt-1 uppercase tracking-widest">Katalog Paket Cuci & Add-On Parfum Premium</p>
                        </div>
                    </div>
                    
                    <button @click="showForm ? cancelForm() : (showForm = true)" :class="showForm ? 'bg-rose-500 hover:bg-rose-600 shadow-rose-200' : 'bg-indigo-600 hover:bg-indigo-700 shadow-indigo-200'" class="w-full sm:w-auto text-white px-5 py-3.5 rounded-xl font-black text-xs uppercase tracking-widest transition-all active:scale-95 shadow-lg flex items-center justify-center gap-2">
                        <svg v-if="!showForm" xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
                        <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                        {{ showForm ? 'Batal Form' : (activeTab === 'jasa' ? 'Tambah Jasa' : 'Tambah Parfum') }}
                    </button>
                </div>

                <div class="flex bg-slate-100 p-1 rounded-xl self-start w-full sm:w-auto">
                    <button @click="switchTab('jasa')" :class="activeTab === 'jasa' ? 'bg-white text-indigo-600 shadow-sm' : 'text-slate-500 hover:text-slate-800'" class="flex-1 sm:flex-initial px-4 md:px-5 py-2.5 rounded-lg text-[11px] font-black uppercase tracking-wider transition-all flex items-center justify-center gap-2">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M6 2 3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4Z"/><path d="M3 6h18"/><path d="M16 10a4 4 0 0 1-8 0"/></svg>
                        Paket Cucian
                    </button>
                    <button @click="switchTab('parfum')" :class="activeTab === 'parfum' ? 'bg-white text-indigo-600 shadow-sm' : 'text-slate-500 hover:text-slate-800'" class="flex-1 sm:flex-initial px-4 md:px-5 py-2.5 rounded-lg text-[11px] font-black uppercase tracking-wider transition-all flex items-center justify-center gap-2">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M12 22c5.523 0 10-4.477 10-10S17.523 2 12 2 2 6.477 2 12s4.477 10 10 10z"/><path d="M12 8a4 4 0 1 0 0 8 4 4 0 0 0 0-8z"/></svg>
                        Parfum Premium
                    </button>
                </div>

                <div class="relative w-full group max-w-2xl">
                    <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-400 group-focus-within:text-indigo-500 transition-colors" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
                    </div>
                    <input v-model="searchQuery" type="text" :placeholder="activeTab === 'jasa' ? 'Cari nama paket cuci...' : 'Cari nama parfum premium...'" class="w-full pl-12 pr-4 py-3 bg-slate-50 border-2 border-slate-200 rounded-xl focus:bg-white focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 outline-none font-bold text-sm text-slate-800 transition-all placeholder:text-slate-300 shadow-inner">
                </div>
            </div>

            <div class="flex-1 overflow-y-auto custom-scrollbar p-4 sm:p-6 md:p-8">
                <CatalogFormCard 
                    v-if="showForm" 
                    :activeTab="activeTab" 
                    :isEditing="isEditing" 
                    :formJasa="formJasa" 
                    :formParfum="formParfum" 
                    @save="handleSave" 
                    @cancel="cancelForm" 
                />

                <div v-if="isLoading" class="flex flex-col items-center justify-center py-20 text-slate-400">
                    <div class="w-10 h-10 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin mb-4"></div>
                    <p class="font-black text-xs uppercase tracking-[0.2em] animate-pulse">Memuat Katalog...</p>
                </div>
                
                <CatalogListDisplay 
                    v-else 
                    :activeTab="activeTab" 
                    :items="filteredItems" 
                    :formatRupiah="formatRupiah" 
                    @edit="triggerEdit" 
                    @delete="handleConfirmDelete" 
                />
            </div>
        </div>
    </SidebarLaundry>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { height: 6px; width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #94a3b8; }
input[type=number]::-webkit-inner-spin-button, input[type=number]::-webkit-outer-spin-button { -webkit-appearance: none; margin: 0; }
input[type=number] { -moz-appearance: textfield; }
@keyframes fadeInDown { from { opacity: 0; transform: translateY(-15px); } to { opacity: 1; transform: translateY(0); } }
</style>