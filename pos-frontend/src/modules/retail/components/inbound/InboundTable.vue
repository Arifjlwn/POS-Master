<script setup>
import { defineProps, defineEmits } from 'vue';

defineProps({
    cartLPB: Array,
    isOwner: Boolean,
    hitungTotalStok: Function,
    hitungModalPerPcs: Function
});

defineEmits(['remove']);
</script>

<template>
    <table class="w-full text-left whitespace-nowrap border-collapse">
        <thead>
            <tr class="bg-slate-50/50 text-[10px] font-black text-slate-400 uppercase tracking-[0.2em] border-b border-slate-100">
                <th class="p-6">Nama Barang</th>
                <th class="p-6 text-center">Qty Karton</th>
                <th class="p-6 text-center">Qty Eceran</th> <th class="p-6 text-center bg-blue-50/30 text-blue-600">Total Pcs</th>
                <th class="p-6 text-right">Harga Total Nota</th>
                <th v-if="isOwner" class="p-6 text-right">Modal Baru (PCS)</th>
                <th class="p-6 text-center">Aksi</th>
            </tr>
        </thead>
        <tbody class="divide-y divide-slate-50">
            <tr v-if="cartLPB.length === 0">
                <td :colspan="isOwner ? 7 : 6" class="p-16 text-center text-slate-400 font-black text-[10px] uppercase tracking-[0.3em] opacity-50">Belum ada barang di scan</td>
            </tr>
            <tr v-for="(item, index) in cartLPB" :key="item.product_id" 
                :class="(isOwner && hitungModalPerPcs(item) > item.harga_jual_saat_ini) ? 'bg-red-50/50' : 'hover:bg-slate-50/50'">
                
                <td class="p-6">
                    <div class="font-black text-slate-800 text-sm uppercase tracking-tight">{{ item.nama_produk }}</div>
                    <div v-if="isOwner && hitungModalPerPcs(item) > item.harga_jual_saat_ini" class="mt-1 text-[8px] font-black text-red-600 animate-pulse uppercase">⚠️ Potensi Rugi!</div>
                </td>

                <td class="p-6 text-center">
                    <input v-model.number="item.qty_karton" type="number" min="0" class="w-16 p-3 bg-white border-2 border-slate-100 rounded-xl text-center font-black text-slate-700 outline-none focus:border-indigo-500">
                </td>

                <td class="p-6 text-center">
                    <input v-model.number="item.qty_eceran" type="number" min="0" class="w-16 p-3 bg-white border-2 border-slate-100 rounded-xl text-center font-black text-blue-600 outline-none focus:border-blue-500">
                </td>

                <td class="p-6 text-center bg-blue-50/20">
                    <div class="flex flex-col items-center">
                        <span class="text-lg font-black text-blue-700 leading-none">{{ hitungTotalStok(item) }}</span>
                        <span class="text-[8px] font-black text-blue-400 uppercase mt-1 tracking-widest">{{ item.satuan_dasar }}</span>
                    </div>
                </td>

                <td class="p-6 min-w-[180px]"> 
                    <div class="relative w-full">
                        <span class="absolute inset-y-0 left-0 pl-4 flex items-center text-xs font-black text-slate-400 italic">Rp</span>
                        <input v-model.number="item.harga_beli_input" type="number" min="0" class="w-full pl-11 pr-4 py-3 bg-white border-2 border-slate-100 rounded-xl text-right font-black text-slate-800 outline-none focus:border-blue-500 transition-all text-sm shadow-inner">
                    </div>
                </td>

                <td v-if="isOwner" class="p-6 text-right">
                    <div class="text-base font-black tracking-tight text-emerald-600" :class="{ 'text-red-600': hitungModalPerPcs(item) > item.harga_jual_saat_ini }">
                        Rp {{ hitungModalPerPcs(item).toLocaleString('id-ID') }}
                    </div>
                </td>

                <td class="p-6 text-center">
                    <button @click="$emit('remove', index)" class="p-3 bg-slate-50 hover:bg-red-50 rounded-2xl transition-all group border border-transparent hover:border-red-100">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-300 group-hover:text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h14" /></svg>
                    </button>
                </td>
            </tr>
        </tbody>
    </table>
</template>