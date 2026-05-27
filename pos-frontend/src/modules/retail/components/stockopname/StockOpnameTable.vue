<script setup>
defineProps({ 
    cartSO: Array,
    soStep: String,
    isOwner: Boolean,
    activeTab: { type: String, default: 'SO' } // 🚀 TERIMA PROP TAB AKTIF
});
defineEmits(['remove']);

</script>

<template>
    <div class="bg-white rounded-[40px] shadow-xl shadow-slate-200/50 border border-slate-100 overflow-hidden">
        <div class="overflow-x-auto custom-scrollbar">
            <table class="w-full text-left whitespace-nowrap border-collapse min-w-[600px]">
                
                <thead>
                    <tr class="bg-slate-50/50 text-[10px] font-black text-slate-400 uppercase tracking-[0.2em] border-b border-slate-100">
                        <th class="p-6">Informasi Barang</th>
                        
                        <template v-if="activeTab === 'SO'">
                            <th v-if="isOwner || soStep === 'REVIEW'" class="p-6 text-center bg-indigo-50/50 text-indigo-600">Stok Sistem</th>
                            <th class="p-6 text-center w-36">Fisik (Real)</th>
                            <th v-if="isOwner || soStep === 'REVIEW'" class="p-6 text-center">Selisih</th>
                        </template>

                        <template v-else>
                            <th class="p-6 text-center w-36">Fisik Ketemu</th>
                            <th class="p-6 text-center w-64">Bukti Keterangan & Foto</th>
                        </template>

                        <th class="p-6 text-center w-16">Aksi</th>
                    </tr>
                </thead>

                <tbody class="divide-y divide-slate-50">
                    <tr v-if="cartSO.length === 0">
                        <td :colspan="activeTab === 'SO' ? (isOwner || soStep === 'REVIEW' ? 5 : 3) : 4" class="p-20 text-center">
                            <div class="flex flex-col items-center opacity-30">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-20 h-20 mb-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" /></svg>
                                <p class="text-xs font-black text-slate-600 uppercase tracking-[0.3em]">Belum ada barang dihitung</p>
                            </div>
                        </td>
                    </tr>

                    <template v-for="(item, index) in cartSO" :key="index">
                        <tr class="hover:bg-slate-50/50 transition-colors group">
                            <td class="p-6">
                                <div class="font-black text-slate-800 text-sm uppercase tracking-tight">{{ item.nama_produk }}</div>
                                <div class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-1">SKU: {{ item.sku || '-' }}</div>
                            </td>
                            
                            <template v-if="activeTab === 'SO'">
                                <td v-if="isOwner || soStep === 'REVIEW'" class="p-6 text-center bg-indigo-50/30">
                                    <div class="w-12 h-12 mx-auto rounded-xl bg-white flex items-center justify-center font-black text-slate-500 text-lg border border-slate-200 shadow-sm">
                                        {{ item.system_qty }}
                                    </div>
                                </td>
                                <td class="p-6 text-center">
                                    <div class="relative w-full mx-auto">
                                        <input v-model.number="item.actual_qty" type="number" min="0" class="w-full p-3 bg-white border-2 border-slate-200 focus:border-indigo-600 rounded-xl text-center font-black text-indigo-700 text-lg outline-none transition-all shadow-inner">
                                    </div>
                                </td>
                                <td v-if="isOwner || soStep === 'REVIEW'" class="p-6 text-center">
                                    <div v-if="(item.actual_qty - item.system_qty) > 0" class="inline-flex flex-col items-center justify-center w-12 h-12 rounded-full bg-emerald-50 text-emerald-600 font-black border border-emerald-200">+{{ item.actual_qty - item.system_qty }}</div>
                                    <div v-else-if="(item.actual_qty - item.system_qty) < 0" class="inline-flex flex-col items-center justify-center w-12 h-12 rounded-full bg-rose-50 text-rose-600 font-black border border-rose-200">{{ item.actual_qty - item.system_qty }}</div>
                                    <div v-else class="inline-flex flex-col items-center justify-center w-12 h-12 rounded-full bg-slate-50 text-slate-400 font-black border border-slate-200">✓</div>
                                </td>
                            </template>

                            <template v-else>
                                <td class="p-6 text-center">
                                    <div class="relative w-full mx-auto">
                                        <input v-model.number="item.actual_qty" type="number" min="1" class="w-full p-3 bg-amber-50 border-2 border-amber-200 focus:border-amber-500 rounded-xl text-center font-black text-amber-700 text-lg outline-none transition-all shadow-inner" placeholder="Qty">
                                    </div>
                                </td>
                                <td class="p-6">
                                    <div class="flex flex-col gap-3">
                                        <input v-model="item.alasan" type="text" placeholder="Contoh: Ketemu di gudang belakang..." class="w-full p-3 bg-white border border-slate-200 focus:border-amber-500 rounded-xl outline-none font-bold text-xs text-slate-700 shadow-sm">
                                    </div>
                                </td>
                            </template>

                            <td class="p-6 text-center">
                                <button @click="$emit('remove', index)" class="p-3 bg-slate-50 hover:bg-red-50 rounded-2xl transition-all border border-transparent"><svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-300 hover:text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h14" /></svg></button>
                            </td>
                        </tr>
                    </template>
                </tbody>
            </table>
        </div>
    </div>
</template>