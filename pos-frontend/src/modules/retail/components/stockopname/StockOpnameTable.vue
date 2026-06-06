<script setup>
const props = defineProps({ 
    cartSO: Array,
    soStep: String,
    isOwner: Boolean,
    activeTab: { type: String, default: 'SO' },
    hitungTotalFisik: Function 
});
defineEmits(['remove']);

const formatNumber = (val) => {
    if (val === null || val === undefined || val === '') return '0';
    return Number(val).toLocaleString('id-ID');
};

// 🛡️ SECURITY LAYER INPUT: Paksa buang semua karakter non-angka bulat secara real-time  !
const enforceOnlyPureDigits = (event, item, key) => {
    let cleanVal = event.target.value.replace(/\D/g, ""); // Buang minus, desimal, huruf e 
    item[key] = cleanVal ? parseInt(cleanVal, 10) : 0;
};
</script>

<template>
    <div class="bg-white rounded-[32px] shadow-xl shadow-slate-200/50 border border-slate-100 overflow-hidden select-none">
        <div class="overflow-x-auto custom-scrollbar">
            <table class="w-full text-left whitespace-nowrap border-collapse min-w-[700px]">
                
                <thead>
                    <tr class="bg-slate-50/50 text-[10px] font-black text-slate-400 uppercase tracking-[0.2em] border-b border-slate-100">
                        <th class="p-4 pl-6 w-1/4">Informasi Barang</th>
                        
                        <template v-if="activeTab === 'SO'">
                            <th v-if="isOwner || soStep === 'REVIEW'" class="p-4 text-center bg-indigo-50/50 text-indigo-600">Stok Sistem</th>
                            <th class="p-4 text-center">Fisik Aktual (Kemasan)</th>
                            <th class="p-4 text-center bg-blue-50/30 text-blue-600">Total Fisik</th>
                            <th v-if="isOwner || soStep === 'REVIEW'" class="p-4 text-center">Selisih</th>
                        </template>

                        <template v-else>
                            <th class="p-4 text-center">Kuantitas Koreksi (Kemasan)</th>
                            <th class="p-4 text-center bg-amber-50/30 text-amber-600">Total Temuan</th>
                            <th class="p-4 text-center w-1/4">Catatan Otorisasi Lokasi</th>
                        </template>

                        <th class="p-4 text-center w-12 pr-6">Aksi</th>
                    </tr>
                </thead>

                <tbody class="divide-y divide-slate-50 font-sans">
                    <tr v-if="!cartSO || cartSO.length === 0">
                        <td :colspan="activeTab === 'SO' ? (isOwner || soStep === 'REVIEW' ? 6 : 4) : 4" class="p-16 text-center">
                            <div class="flex flex-col items-center opacity-30">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-16 h-16 mb-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" /></svg>
                                <p class="text-[10px] font-black text-slate-600 uppercase tracking-[0.3em]">Belum ada data barang dalam daftar hitung</p>
                            </div>
                        </td>
                    </tr>

                    <template v-else v-for="(item, index) in cartSO" :key="item.product_id + '_' + index">
                        <tr class="hover:bg-slate-50/50 transition-colors group">
                            
                            <td class="p-4 pl-6">
                                <div class="font-black text-slate-800 text-xs uppercase tracking-tight line-clamp-1" :title="item.nama_produk">{{ item.nama_produk }}</div>
                                <div class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-0.5">SKU: {{ item.sku || '-' }}</div>
                                
                                <div class="mt-2 flex gap-1">
                                    <span v-if="item.is_nested" class="px-2 py-0.5 bg-indigo-50 text-indigo-600 rounded text-[8px] font-black uppercase tracking-widest border border-indigo-100">3 Lapis (Grosir)</span>
                                    <span v-else-if="item.has_satuan_besar" class="px-2 py-0.5 bg-emerald-50 text-emerald-600 rounded text-[8px] font-black uppercase tracking-widest border border-emerald-100">2 Lapis (Medium)</span>
                                    <span v-else class="px-2 py-0.5 bg-slate-100 text-slate-500 rounded text-[8px] font-black uppercase tracking-widest border border-slate-200">1 Lapis (Dasar)</span>
                                </div>
                            </td>
                            
                            <template v-if="activeTab === 'SO'">
                                <td v-if="isOwner || soStep === 'REVIEW'" class="p-4 text-center bg-indigo-50/30">
                                    <div class="inline-flex items-center gap-1.5 bg-white px-3 py-1.5 rounded-lg border border-slate-200 shadow-sm">
                                        <span class="font-black text-slate-600 text-xs md:text-sm">{{ formatNumber(item.system_qty) }}</span>
                                        <span class="text-[8px] font-black text-slate-400 uppercase mt-0.5">{{ item.satuan_dasar }}</span>
                                    </div>
                                </td>

                                <td class="p-4 text-center">
                                    <div class="flex items-center justify-center gap-1.5">
                                        <div v-if="item.has_satuan_besar" class="flex items-center bg-white border-2 border-slate-200 focus-within:border-indigo-500 rounded-lg overflow-hidden transition-all shadow-inner w-[95px]">
                                            <input :value="item.qty_besar || ''" @input="enforceOnlyPureDigits($event, item, 'qty_besar')" type="text" inputmode="numeric" class="w-10 p-1.5 text-center font-black text-indigo-700 text-xs md:text-sm outline-none bg-transparent">
                                            <span class="text-[8px] font-black text-slate-400 uppercase bg-slate-50 w-full p-2 border-l border-slate-200 truncate" :title="item.satuan_besar">{{ item.satuan_besar || 'BSR' }}</span>
                                        </div>

                                        <span v-if="item.is_nested" class="text-slate-300 font-bold text-xs">+</span>

                                        <div v-if="item.is_nested" class="flex items-center bg-white border-2 border-slate-200 focus-within:border-emerald-500 rounded-lg overflow-hidden transition-all shadow-inner w-[95px]">
                                            <input :value="item.qty_tengah || ''" @input="enforceOnlyPureDigits($event, item, 'qty_tengah')" type="text" inputmode="numeric" class="w-10 p-1.5 text-center font-black text-emerald-600 text-xs md:text-sm outline-none bg-transparent">
                                            <span class="text-[8px] font-black text-emerald-400 uppercase bg-slate-50 w-full p-2 border-l border-slate-200 truncate" :title="item.satuan_tengah">{{ item.satuan_tengah || 'TGH' }}</span>
                                        </div>

                                        <span v-if="item.has_satuan_besar" class="text-slate-300 font-bold text-xs">+</span>

                                        <div class="flex items-center bg-white border-2 border-slate-200 focus-within:border-blue-500 rounded-lg overflow-hidden transition-all shadow-inner w-[95px]">
                                            <input :value="item.qty_dasar || ''" @input="enforceOnlyPureDigits($event, item, 'qty_dasar')" type="text" inputmode="numeric" class="w-10 p-1.5 text-center font-black text-blue-600 text-xs md:text-sm outline-none bg-transparent">
                                            <span class="text-[8px] font-black text-blue-400 uppercase bg-slate-50 w-full p-2 border-l border-slate-200 truncate" :title="item.satuan_dasar">{{ item.satuan_dasar || 'DSR' }}</span>
                                        </div>
                                    </div>
                                </td>

                                <td class="p-4 text-center bg-blue-50/20">
                                    <div class="flex flex-col items-center">
                                        <span class="text-sm md:text-base font-black text-blue-700 leading-none">{{ formatNumber(hitungTotalFisik(item)) }}</span>
                                        <span class="text-[8px] font-black text-blue-400 uppercase mt-1 tracking-widest">{{ item.satuan_dasar }}</span>
                                    </div>
                                </td>

                                <td v-if="isOwner || soStep === 'REVIEW'" class="p-4 text-center">
                                    <div class="inline-flex items-center justify-center px-2.5 py-1 rounded-lg border font-black text-xs min-w-[55px]"
                                         :class="(hitungTotalFisik(item) - item.system_qty) > 0 ? 'bg-emerald-50 text-emerald-600 border-emerald-200' : ((hitungTotalFisik(item) - item.system_qty) < 0 ? 'bg-rose-50 text-rose-600 border-rose-200' : 'bg-slate-50 text-slate-400 border-slate-200')">
                                        <template v-if="(hitungTotalFisik(item) - item.system_qty) > 0">+</template>
                                        {{ (hitungTotalFisik(item) - item.system_qty) === 0 ? '✓ OK' : formatNumber(hitungTotalFisik(item) - item.system_qty) }}
                                    </div>
                                </td>
                            </template>

                            <template v-else>
                                <td class="p-4 text-center">
                                    <div class="flex items-center justify-center gap-1.5">
                                        <div v-if="item.has_satuan_besar" class="flex items-center bg-amber-50 border-2 border-amber-100 focus-within:border-amber-400 rounded-lg overflow-hidden transition-all shadow-inner w-[95px]">
                                            <input :value="item.qty_besar || ''" @input="enforceOnlyPureDigits($event, item, 'qty_besar')" type="text" inputmode="numeric" class="w-10 p-1.5 text-center font-black text-amber-700 text-xs md:text-sm outline-none bg-transparent">
                                            <span class="text-[8px] font-black text-amber-500 uppercase bg-amber-100/30 w-full p-2 border-l border-amber-200/40 truncate">{{ item.satuan_besar || 'BSR' }}</span>
                                        </div>

                                        <span v-if="item.is_nested" class="text-amber-300 font-bold text-xs">+</span>

                                        <div v-if="item.is_nested" class="flex items-center bg-amber-50 border-2 border-amber-100 focus-within:border-amber-400 rounded-lg overflow-hidden transition-all shadow-inner w-[95px]">
                                            <input :value="item.qty_tengah || ''" @input="enforceOnlyPureDigits($event, item, 'qty_tengah')" type="text" inputmode="numeric" class="w-10 p-1.5 text-center font-black text-amber-700 text-xs md:text-sm outline-none bg-transparent">
                                            <span class="text-[8px] font-black text-amber-500 uppercase bg-amber-100/30 w-full p-2 border-l border-amber-200/40 truncate">{{ item.satuan_tengah || 'TGH' }}</span>
                                        </div>

                                        <span v-if="item.has_satuan_besar" class="text-amber-300 font-bold text-xs">+</span>

                                        <div class="flex items-center bg-amber-50 border-2 border-amber-100 focus-within:border-amber-400 rounded-lg overflow-hidden transition-all shadow-inner w-[95px]">
                                            <input :value="item.qty_dasar || ''" @input="enforceOnlyPureDigits($event, item, 'qty_dasar')" type="text" inputmode="numeric" class="w-10 p-1.5 text-center font-black text-amber-700 text-xs md:text-sm outline-none bg-transparent">
                                            <span class="text-[8px] font-black text-amber-500 uppercase bg-amber-100/30 w-full p-2 border-l border-amber-200/40 truncate">{{ item.satuan_dasar || 'DSR' }}</span>
                                        </div>
                                    </div>
                                </td>

                                <td class="p-4 text-center bg-amber-50/10">
                                    <div class="flex flex-col items-center">
                                        <span class="text-sm md:text-base font-black leading-none" :class="hitungTotalFisik(item) > item.max_klaim ? 'text-rose-600 animate-pulse' : 'text-amber-600'">
                                            {{ formatNumber(hitungTotalFisik(item)) }}
                                        </span>
                                        <span class="text-[8px] font-black uppercase mt-1.5 tracking-wider px-2 py-0.5 rounded whitespace-nowrap"
                                              :class="hitungTotalFisik(item) > item.max_klaim ? 'bg-rose-100 text-rose-600' : 'bg-amber-50 text-amber-600 border border-amber-100'">
                                            LIMIT MAKS: {{ formatNumber(item.max_klaim) }} {{ item.satuan_dasar }}
                                        </span>
                                    </div>
                                </td>

                                <td class="p-4">
                                    <input v-model="item.alasan" :disabled="hitungTotalFisik(item) === 0" type="text" placeholder="Contoh: Ketemu nyempil di kolong rak 4 B ..." class="w-full p-2.5 bg-white border-2 border-slate-200 focus:border-amber-500 rounded-xl outline-none font-bold text-[10px] text-slate-700 shadow-sm disabled:opacity-50 disabled:bg-slate-50 disabled:cursor-not-allowed">
                                </td>
                            </template>
                            
                            <td class="p-4 text-center pr-6">
                                <button v-if="activeTab === 'SO' && soStep === 'COUNTING'" @click="$emit('remove', index)" class="p-2 bg-slate-50 hover:bg-rose-50 rounded-lg transition-all border border-transparent text-slate-300 hover:text-rose-600 active:scale-95 shadow-sm" title="Hapus Dari Daftar">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h14" /></svg>
                                </button>
                                <span v-else class="text-[9px] font-black text-slate-300 tracking-wider bg-slate-100 px-2 py-1 rounded border uppercase select-none">Kunci</span>
                            </td>
                        </tr>
                    </template>
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