<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import SidebarFnB from './SidebarFnB.vue';
import api from '../../api.js'; 
import Swal from 'sweetalert2';

const activeOrders = ref([]);
const isLoading = ref(false);
let pollingInterval = null;

const fetchAntreanDapur = async (isBackground = false) => {
    if (!isBackground) isLoading.value = true;
    try {
        const response = await api.get('/fnb/kitchen');
        activeOrders.value = response.data || [];
    } catch (error) {
        console.error("Gagal menarik data", error);
    } finally {
        if (!isBackground) isLoading.value = false;
    }
};

onMounted(() => {
    fetchAntreanDapur();
    pollingInterval = setInterval(() => fetchAntreanDapur(true), 5000);
});

onUnmounted(() => clearInterval(pollingInterval));

const selesaikanOrder = (order) => {
    Swal.fire({
        title: 'Pesanan Selesai?',
        html: `Makanan untuk <b class="text-indigo-600">${order.meja || order.invoice}</b> siap?`,
        icon: 'question',
        showCancelButton: true,
        confirmButtonColor: '#10b981',
        cancelButtonText: 'Batal',
        confirmButtonText: 'Ya, Siap Saji!',
        customClass: { popup: 'rounded-[24px] font-sans' }
    }).then(async (result) => {
        if (result.isConfirmed) {
            try {
                await api.put(`/fnb/kitchen/${order.id}/selesai`);
                fetchAntreanDapur(true);
                Swal.fire({ toast: true, position: 'top-end', icon: 'success', title: 'Kasir dinotifikasi!', showConfirmButton: false, timer: 2000, customClass: { popup: 'rounded-xl' }});
            } catch (error) {
                Swal.fire('Error!', 'Gagal update status.', 'error');
            }
        }
    });
};
</script>

<template>
    <SidebarFnB>
        <div class="flex-1 flex flex-col h-full bg-[#0F172A] overflow-hidden relative font-sans text-slate-200">
            
            <div class="p-4 md:p-6 bg-[#1E293B] border-b border-slate-800 flex justify-between items-center z-10 shadow-md shrink-0">
                <div class="flex items-center gap-3 md:gap-4">
                    <div class="w-10 h-10 md:w-12 md:h-12 bg-orange-500/20 border border-orange-500/50 rounded-xl flex items-center justify-center shrink-0 shadow-inner">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 md:w-6 md:h-6 text-orange-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M6 13.87A4 4 0 0 1 7.41 6a5.11 5.11 0 0 1 1.05-1.54 5 5 0 0 1 7.08 0A5.11 5.11 0 0 1 16.59 6 4 4 0 0 1 18 13.87V21H6Z"/><line x1="6" y1="17" x2="18" y2="17"/></svg>
                    </div>
                    <div>
                        <h1 class="text-base sm:text-xl md:text-2xl font-black tracking-tighter uppercase text-white leading-none">Kitchen Monitor</h1>
                        <p class="text-[9px] md:text-[10px] font-black text-slate-400 mt-1 uppercase tracking-widest">KDS Realtime System</p>
                    </div>
                </div>
                
                <div class="bg-slate-900 px-3 py-2 md:px-5 md:py-2.5 rounded-xl border border-slate-700 flex items-center gap-3 shadow-inner">
                    <div v-if="isLoading && activeOrders.length === 0" class="w-3 h-3 md:w-4 md:h-4 rounded-full border-2 border-slate-600 border-t-orange-500 animate-spin"></div>
                    <div v-else class="w-2.5 h-2.5 bg-emerald-500 rounded-full shadow-[0_0_8px_rgba(16,185,129,0.8)] animate-pulse hidden sm:block"></div>
                    
                    <div class="text-right">
                        <p class="text-[8px] md:text-[10px] font-black text-slate-400 uppercase tracking-widest mb-0.5">Antrean</p>
                        <p class="text-sm md:text-xl font-black text-white leading-none">{{ activeOrders.length }} <span class="hidden sm:inline-block text-[10px] text-orange-400">Order</span></p>
                    </div>
                </div>
            </div>

            <div class="flex-1 overflow-y-auto custom-scrollbar p-4 md:p-6 pb-20">
                <div v-if="activeOrders.length === 0" class="h-full flex flex-col items-center justify-center text-slate-600">
                    <div class="w-20 h-20 md:w-24 md:h-24 bg-slate-800 rounded-full flex items-center justify-center mb-5 shadow-inner">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-10 h-10 md:w-12 md:h-12 text-slate-500" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M12 2v20M17 5H9.5a4.5 4.5 0 0 0 0 9H15a4.5 4.5 0 0 1 0 9H6.5"/></svg>
                    </div>
                    <h2 class="font-black text-lg md:text-xl text-slate-500 uppercase tracking-widest mb-2">Dapur Clear!</h2>
                    <p class="text-[10px] md:text-xs font-bold text-slate-600 uppercase tracking-widest text-center">Tarik nafas dulu bosku ☕</p>
                </div>

                <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-4 gap-4 md:gap-6 items-start">
                    <div v-for="order in activeOrders" :key="order.id" class="bg-white rounded-[20px] md:rounded-[24px] overflow-hidden shadow-xl flex flex-col animate-[fadeIn_0.3s_ease-out] border-4" :class="order.tipe === 'DINE_IN' ? 'border-indigo-600' : 'border-rose-500'">
                        
                        <div class="p-3.5 md:p-4 flex justify-between items-start" :class="order.tipe === 'DINE_IN' ? 'bg-indigo-600' : 'bg-rose-500'">
                            <div>
                                <div class="flex items-center gap-2 mb-1">
                                    <span class="bg-black/20 text-white text-[9px] font-black px-2 py-1 rounded uppercase tracking-widest">#{{ order.invoice.split('/').pop() || order.id }}</span>
                                    <span class="text-[9px] md:text-[10px] font-black text-white/80 uppercase flex items-center gap-1"><svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>{{ order.waktu_pesan }}</span>
                                </div>
                                <h3 class="font-black text-lg md:text-xl text-white uppercase tracking-tight leading-none mt-1">{{ order.tipe === 'DINE_IN' ? (order.meja || 'Makan Sini') : 'Bungkus' }}</h3>
                            </div>
                        </div>

                        <div class="p-3.5 md:p-4 flex-1 bg-slate-50 space-y-3">
                            <div v-for="(item, index) in order.items" :key="index" class="flex gap-3 pb-3 border-b border-slate-200 border-dashed last:border-0 last:pb-0">
                                <div class="w-8 h-8 md:w-10 md:h-10 rounded-xl bg-slate-800 flex items-center justify-center shrink-0 font-black text-sm md:text-base text-white shadow-sm">{{ item.qty }}<span class="text-[9px] ml-0.5 text-slate-400">x</span></div>
                                <div class="flex-1 pt-0.5">
                                    <p class="font-black text-xs md:text-sm text-slate-900 uppercase leading-tight">{{ item.nama }}</p>
                                    <div v-if="item.notes" class="mt-2 bg-amber-100 border-l-4 border-amber-500 px-2.5 md:px-3 py-1.5 md:py-2 rounded-r-lg shadow-sm">
                                        <p class="text-[10px] md:text-xs font-black text-amber-900 flex items-start gap-1.5 uppercase leading-snug"><svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 shrink-0 mt-0.5 text-amber-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>{{ item.notes }}</p>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <div class="p-3.5 md:p-4 bg-white border-t border-slate-100">
                            <button @click="selesaikanOrder(order)" class="w-full bg-emerald-500 hover:bg-emerald-600 text-white py-3.5 rounded-xl font-black text-[10px] md:text-xs uppercase tracking-[0.2em] shadow-lg shadow-emerald-500/30 transition-all active:scale-95 flex items-center justify-center gap-2">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
                                Siap Saji
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </SidebarFnB>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { height: 6px; width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: rgba(0,0,0,0.1); }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #334155; border-radius: 10px; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
</style>