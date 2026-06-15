<script setup>
defineProps({
    activeTab: { type: String, required: true },
    items: { type: Array, required: true },
    formatRupiah: { type: Function, required: true }
});

const emit = defineEmits(['edit', 'delete']);
</script>

<template>
    <div v-if="items.length === 0" class="flex flex-col items-center justify-center py-20 text-slate-400 bg-white rounded-3xl border border-slate-100 border-dashed border-2">
        <div class="w-14 h-14 bg-slate-50 rounded-full flex items-center justify-center mb-4 border border-slate-100">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-slate-300" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
        </div>
        <p class="font-black text-xs uppercase tracking-widest text-slate-500 mb-1">Data Belum Tersedia</p>
        <p class="text-[11px] font-bold">Silahkan tambahkan item katalog baru di atas</p>
    </div>

    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4 md:gap-5">
        <template v-if="activeTab === 'jasa'">
            <div v-for="item in items" :key="item.id" class="bg-white rounded-[24px] border-2 border-slate-100 p-5 md:p-6 relative group hover:border-indigo-500 transition-all shadow-sm hover:shadow-xl hover:shadow-indigo-100/40 overflow-hidden flex flex-col justify-between h-40 md:h-48">
                <div class="absolute top-4 right-4 flex gap-2 sm:opacity-0 group-hover:opacity-100 transition-all z-20">
                    <button @click.stop="emit('edit', item)" class="w-8 h-8 bg-white border border-slate-100 text-slate-400 rounded-full flex items-center justify-center hover:bg-sky-500 hover:text-white hover:border-sky-500 shadow-sm transition-colors">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
                    </button>
                    <button @click.stop="emit('delete', item.id, item.nama_produk)" class="w-8 h-8 bg-white border border-slate-100 text-slate-400 rounded-full flex items-center justify-center hover:bg-rose-500 hover:text-white hover:border-rose-500 shadow-sm transition-colors">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/></svg>
                    </button>
                </div>
                <div class="relative z-10">
                    <div class="mb-2.5">
                        <span class="text-[9px] font-black bg-indigo-50 text-indigo-600 px-2.5 py-1 rounded-lg uppercase tracking-widest border border-indigo-100 inline-block">Est: {{ item.estimasi || 'Standar' }}</span>
                    </div>
                    <h3 class="text-sm font-black text-slate-800 uppercase leading-tight line-clamp-2 pr-14">{{ item.nama_produk }}</h3>
                </div>
                <div class="relative z-10 pt-4 border-t border-slate-100 mt-auto flex items-end justify-between">
                    <p class="text-base md:text-xl font-black text-emerald-500 leading-none">{{ formatRupiah(item.harga_jual) }}</p>
                    <span class="text-[9px] md:text-[10px] text-slate-400 font-bold uppercase tracking-widest">/ {{ item.satuan_dasar }}</span>
                </div>
            </div>
        </template>

        <template v-else>
            <div v-for="perfume in items" :key="perfume.id" class="bg-white rounded-[24px] border-2 border-slate-100 p-5 md:p-6 relative group hover:border-pink-500 transition-all shadow-sm hover:shadow-xl hover:shadow-pink-100/40 overflow-hidden flex flex-col justify-between h-40 md:h-48">
                <div class="absolute top-4 right-4 flex gap-2 sm:opacity-0 group-hover:opacity-100 transition-all z-20">
                    <button @click.stop="emit('delete', perfume.id, perfume.nama)" class="w-8 h-8 bg-white border border-slate-100 text-slate-400 rounded-full flex items-center justify-center hover:bg-rose-500 hover:text-white hover:border-rose-500 shadow-sm transition-colors">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/></svg>
                    </button>
                </div>
                <div class="relative z-10">
                    <div class="mb-2.5">
                        <span :class="perfume.status === 'Tersedia' ? 'bg-emerald-50 text-emerald-600 border-emerald-100' : 'bg-rose-50 text-rose-600 border-rose-100'" class="text-[9px] font-black px-2.5 py-1 rounded-lg uppercase tracking-widest border inline-block">
                            {{ perfume.status || 'Tersedia' }}
                        </span>
                    </div>
                    <div class="flex items-center gap-2">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-pink-500 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M12 22c5.523 0 10-4.477 10-10S17.523 2 12 2 2 6.477 2 12s4.477 10 10 10z"/><path d="M12 8a4 4 0 1 0 0 8 4 4 0 0 0 0-8z"/></svg>
                        <h3 class="text-sm font-black text-slate-800 uppercase leading-tight line-clamp-2 pr-10 capitalize">{{ perfume.nama }}</h3>
                    </div>
                </div>
                <div class="relative z-10 pt-4 border-t border-slate-100 mt-auto flex items-end justify-between">
                    <p class="text-base md:text-xl font-black text-pink-500 leading-none">
                        {{ perfume.harga > 0 ? '+' + formatRupiah(perfume.harga) : 'Gratis Add-On' }}
                    </p>
                </div>
            </div>
        </template>
    </div>
</template>