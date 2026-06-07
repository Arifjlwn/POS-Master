<script setup>
defineProps({ cart: Array, isSubmitting: Boolean, getBadgeClass: Function });
defineEmits(['remove', 'submit']);

const formatNumber = (val) => {
    return Number(val).toLocaleString('id-ID');
};
</script>

<template>
    <div class="bg-white rounded-[32px] shadow-sm border border-slate-100 overflow-hidden flex flex-col flex-1 min-h-[350px] max-h-[75vh] xl:max-h-none xl:h-auto xl:sticky xl:top-24">
        <div class="p-5 sm:p-6 border-b border-slate-50 bg-slate-50/50 flex items-center justify-between shrink-0">
            <h3 class="font-black text-slate-800 text-sm uppercase tracking-widest flex items-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-rose-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z"/><line x1="3" y1="6" x2="21" y2="6"/><path d="M16 10a4 4 0 0 1-8 0"/></svg>
                Keranjang Retur
                <span v-if="cart.length > 0" class="ml-1.5 bg-rose-600 text-white w-6 h-6 rounded-full flex items-center justify-center text-[10px]">{{ cart.length }}</span>
            </h3>
        </div>

        <div v-if="cart.length === 0" class="flex-1 flex flex-col items-center justify-center text-center p-8">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-16 h-16 sm:w-20 sm:h-20 text-slate-200 mb-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z"/><line x1="3" y1="6" x2="21" y2="6"/><path d="M16 10a4 4 0 0 1-8 0"/></svg>
            <p class="text-slate-400 font-black text-sm uppercase tracking-widest">Keranjang Kosong</p>
            <p class="text-slate-400 font-medium text-xs mt-2 max-w-[200px]">Pilih barang di form samping untuk memproses retur.</p>
        </div>

        <div v-auto-animate class="flex-1 overflow-y-auto custom-scrollbar p-3 sm:p-4 bg-slate-50/30 flex flex-col gap-3">
            <div v-for="(item, index) in cart" :key="index" class="bg-white p-3 sm:p-4 rounded-[20px] shadow-sm border border-slate-100 flex flex-col gap-3 relative group">
                
                <div class="flex justify-between items-start pr-12">
                    <div>
                        <div class="font-black text-slate-800 text-xs sm:text-sm uppercase leading-tight line-clamp-2" :title="item.nama_produk">{{ item.nama_produk }}</div>
                        <div class="mt-1.5 flex flex-wrap gap-1">
                            <span class="inline-block px-2 py-0.5 rounded text-[8px] font-black uppercase tracking-widest border" :class="getBadgeClass(item.alasan)">{{ item.alasan }}</span>
                        </div>
                    </div>
                </div>

                <div class="bg-slate-50 border border-slate-100 p-2.5 sm:p-3 rounded-xl flex items-center justify-between">
                    <div class="flex flex-wrap gap-x-3 gap-y-1.5">
                        <span v-if="item.has_satuan_besar" class="text-[9px] sm:text-[10px] font-bold text-slate-600 uppercase">
                            <span class="font-black text-rose-600 text-[10px] sm:text-xs">{{ item.qty_besar }}</span> {{ item.satuan_besar }}
                        </span>
                        <span v-if="item.is_nested" class="text-[9px] sm:text-[10px] font-bold text-slate-600 uppercase">
                            <span class="font-black text-rose-600 text-[10px] sm:text-xs">{{ item.qty_tengah }}</span> {{ item.satuan_tengah }}
                        </span>
                        <span class="text-[9px] sm:text-[10px] font-bold text-slate-600 uppercase">
                            <span class="font-black text-rose-600 text-[10px] sm:text-xs">{{ item.qty_dasar }}</span> {{ item.satuan_dasar }}
                        </span>
                    </div>
                    
                    <div class="bg-white px-2.5 sm:px-3 py-1.5 rounded-lg border border-slate-200 shadow-sm text-center shrink-0 ml-2">
                        <div class="text-[8px] font-black text-slate-400 uppercase tracking-widest mb-0.5">Total Dibuang</div>
                        <div class="text-xs sm:text-sm font-black text-slate-800 leading-none">{{ formatNumber(item.qty) }} <span class="text-[9px] sm:text-[10px]">{{ item.satuan_dasar }}</span></div>
                    </div>
                </div>

                <p v-if="item.catatan" class="text-[9px] sm:text-[10px] text-slate-500 italic mt-0.5">"{{ item.catatan }}"</p>
                
                <button @click="$emit('remove', index)" class="absolute top-3 right-3 sm:top-4 sm:right-4 w-7 h-7 sm:w-8 sm:h-8 flex items-center justify-center rounded-lg bg-red-50 text-red-400 hover:text-white hover:bg-red-500 transition-colors opacity-100 xl:opacity-0 xl:group-hover:opacity-100 active:scale-95">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/></svg>
                </button>
            </div>
        </div>

        <div v-if="cart.length > 0" class="p-4 sm:p-5 bg-white border-t border-slate-100 shrink-0">
            <button @click="$emit('submit')" :disabled="isSubmitting" class="w-full bg-rose-600 hover:bg-slate-900 text-white py-4 sm:py-5 rounded-[20px] sm:rounded-[24px] font-black text-xs sm:text-sm uppercase tracking-[0.2em] shadow-xl shadow-rose-200 hover:shadow-slate-300 transition-all active:scale-95 disabled:opacity-50 flex items-center justify-center gap-2">
                <template v-if="isSubmitting">
                    <div class="w-4 h-4 sm:w-5 sm:h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                    MENYIMPAN...
                </template>
                <template v-else>
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 sm:w-5 sm:h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"/></svg>
                    PROSES PEMUSNAHAN
                </template>
            </button>
        </div>
    </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #94a3b8; }
</style>