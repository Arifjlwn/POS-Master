<script setup>
import { ref, defineProps, defineEmits } from 'vue';

const props = defineProps({
    show: Boolean,
    isEditing: Boolean,
    isSubmitting: Boolean,
    form: Object,
    categories: Array,
    imagePreview: String,
    stokDalamKarton: Number,
    eceranTambahan: Number
});

const emit = defineEmits([
    'close', 
    'submit', 
    'start-scanner', 
    'file-change',
    'update:stokDalamKarton',
    'update:eceranTambahan'
]);

const imageInput = ref(null);

const triggerImageUpload = () => {
    imageInput.value.click();
};
</script>

<template>
    <div v-if="show" class="fixed inset-0 bg-slate-900/80 flex items-center justify-center z-[100] p-4 backdrop-blur-sm no-print">
        <div class="bg-white rounded-[32px] shadow-2xl w-full max-w-xl flex flex-col max-h-[90vh] overflow-hidden border border-slate-100">
            
            <div class="p-6 border-b border-slate-50 bg-slate-50/50 flex justify-between items-center shrink-0">
                <h2 class="text-xl font-black text-slate-800 uppercase italic">{{ isEditing ? 'Edit Data Produk' : 'Registrasi Produk Baru' }}</h2>
                <button @click="emit('close')" class="p-2 rounded-xl bg-white text-slate-400 hover:text-red-500 hover:bg-red-50 transition-all border border-slate-100 shadow-sm">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                </button>
            </div>

            <div class="p-6 md:p-8 overflow-y-auto custom-scrollbar">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
                    
                    <div class="md:col-span-2 flex items-center gap-5 p-4 rounded-[24px] bg-slate-50 border border-slate-100 mb-2">
                        <div @click="triggerImageUpload" class="w-16 h-16 rounded-[18px] border-2 border-dashed border-slate-300 flex items-center justify-center bg-white cursor-pointer overflow-hidden shadow-inner shrink-0 hover:border-blue-400 transition-colors">
                            <img v-if="imagePreview" :src="imagePreview" class="w-full h-full object-cover">
                            <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-slate-300" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14.5 4h-5L7 7H4a2 2 0 0 0-2 2v9a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2h-3l-2.5-3z"/><circle cx="12" cy="13" r="3"/></svg>
                        </div>
                        <div class="flex-1">
                            <input v-model="form.name" type="text" placeholder="NAMA BARANG..." class="w-full px-4 py-3 rounded-xl border border-slate-200 focus:border-blue-600 outline-none font-black text-sm uppercase mb-2 text-slate-800 transition-colors">
                            <div class="flex gap-2">
                                <input v-model="form.sku" type="text" placeholder="BARCODE / SKU..." class="w-full px-4 py-2 rounded-xl border border-slate-200 focus:border-blue-600 outline-none font-bold text-xs uppercase text-slate-800 transition-colors">
                                <button @click.prevent="emit('start-scanner')" class="px-4 bg-blue-50 text-blue-600 hover:bg-blue-600 hover:text-white rounded-xl border border-blue-100 transition-colors shadow-sm">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" /><path stroke-linecap="round" stroke-linejoin="round" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
                                </button>
                            </div>
                        </div>
                        <input type="file" ref="imageInput" @change="(e) => emit('file-change', e)" accept="image/*" class="hidden">
                    </div>

                    <div class="md:col-span-2">
                        <label class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1.5 block">Kategori</label>
                        <input list="kategori-list" v-model="form.category" placeholder="Pilih / Ketik Kategori Baru..." class="w-full px-4 py-3.5 rounded-xl border border-slate-200 focus:border-blue-600 outline-none font-bold text-sm bg-white uppercase transition-all text-slate-800">
                        <datalist id="kategori-list"><option v-for="cat in categories" :key="cat" :value="cat"></option></datalist>
                    </div>

                    <div class="md:col-span-2 p-5 bg-slate-900 rounded-[28px] text-white shadow-xl mt-2 mb-2">
                        <div class="flex items-center gap-2 mb-4">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-blue-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="m7.5 4.27 9 5.15"/><path d="M21 8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16Z"/><path d="M12 22V12"/></svg>
                            <h4 class="font-black text-[10px] uppercase tracking-[0.2em]">Konversi & Satuan Jual</h4>
                        </div>

                        <div class="grid grid-cols-2 gap-4">
                            <div>
                                <label class="text-[8px] font-black text-slate-500 uppercase block mb-1">Satuan Dasar Jual</label>
                                <select v-model="form.satuan_dasar" class="w-full p-3.5 bg-slate-800 border border-slate-700 rounded-xl outline-none font-black text-xs uppercase cursor-pointer text-white">
                                    <option value="PCS">PCS</option>
                                    <option value="KG">KG</option>
                                    <option value="GRAM">GRAM</option>
                                    <option value="PACK">PACK</option>
                                    <option value="BOX">BOX</option>
                                    <option value="LITER">LITER</option>
                                    <option value="BOTOL">BOTOL</option>
                                </select>
                            </div>
                            
                            <div>
                                <label class="text-[8px] font-black text-slate-500 uppercase block mb-1">Beli Dalam Karton/Kemasan Besar?</label>
                                <div @click="form.has_satuan_besar = !form.has_satuan_besar" class="w-full p-3.5 bg-slate-800 border border-slate-700 rounded-xl font-black text-[10px] uppercase cursor-pointer flex items-center justify-between transition-colors hover:border-blue-500">
                                    {{ form.has_satuan_besar ? 'YA (AKTIF)' : 'TIDAK (HANYA PCS)' }}
                                    <div :class="form.has_satuan_besar ? 'bg-blue-500' : 'bg-slate-600'" class="w-2 h-2 rounded-full shadow-[0_0_8px_rgba(59,130,246,0.5)] transition-colors"></div>
                                </div>
                            </div>

                            <div v-if="form.has_satuan_besar" class="col-span-2 grid grid-cols-3 gap-3 pt-4 border-t border-slate-800 mt-2">
                                <div>
                                    <label class="text-[8px] font-black text-blue-400 uppercase block mb-1">Sebutannya Apa?</label>
                                    <input v-model="form.satuan_besar" type="text" placeholder="KARTON / BOX" class="w-full p-3 bg-slate-800 border border-blue-900 focus:border-blue-500 rounded-xl outline-none font-black text-xs uppercase text-white transition-all">
                                </div>
                                <div>
                                    <label class="text-[8px] font-black text-blue-400 uppercase block mb-1">1 {{ form.satuan_besar || 'KEMASAN' }} Isi Berapa {{ form.satuan_dasar }}?</label>
                                    <input v-model.number="form.isi_per_besar" type="number" class="w-full p-3 bg-slate-800 border border-blue-900 focus:border-blue-500 rounded-xl outline-none font-black text-xs text-white transition-all">
                                </div>
                                <div>
                                    <label class="text-[8px] font-black text-amber-400 uppercase block mb-1">Harga Beli 1 {{ form.satuan_besar || 'KEMASAN' }}</label>
                                    <input v-model.number="form.harga_beli_besar" type="number" placeholder="Rp" class="w-full p-3 bg-amber-900/20 border border-amber-900 focus:border-amber-500 rounded-xl outline-none font-black text-xs text-amber-400 transition-all">
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="p-5 border rounded-[28px] transition-all duration-300" :class="form.has_satuan_besar ? 'bg-slate-100/80 border-transparent' : 'bg-white border-slate-200'">
                        <label class="text-[9px] font-black uppercase tracking-widest mb-3 block" :class="form.has_satuan_besar ? 'text-indigo-500' : 'text-slate-400'">Harga Modal Dasar ({{ form.satuan_dasar }})</label>
                        <div class="relative">
                            <span class="absolute inset-y-0 left-0 pl-4 flex items-center text-xs font-black" :class="form.has_satuan_besar ? 'text-indigo-400' : 'text-slate-400'">Rp</span>
                            <input v-model.number="form.cost_price" type="number" :disabled="form.has_satuan_besar" :class="form.has_satuan_besar ? 'text-indigo-600 font-black bg-slate-200/40 cursor-not-allowed border-transparent shadow-none' : 'text-slate-800 font-black bg-white border-slate-200 focus:border-blue-600 outline-none shadow-inner'" class="w-full pl-12 pr-4 py-4 rounded-2xl text-lg border-2">
                        </div>
                        <div v-if="form.has_satuan_besar" class="mt-2.5 flex items-center justify-center gap-1.5">
                            <span class="text-[8px] font-black text-indigo-500 uppercase tracking-widest italic">* Terkunci dari kalkulator grosir</span>
                        </div>
                    </div>

                    <div class="p-5 bg-blue-50 border border-blue-100 rounded-[28px] shadow-sm">
                        <label class="text-[9px] font-black text-blue-400 uppercase tracking-widest mb-3 block">Harga Jual Eceran</label>
                        <div class="relative">
                            <span class="absolute inset-y-0 left-0 pl-4 flex items-center text-xs font-black text-blue-400">Rp</span>
                            <input v-model.number="form.price" type="number" class="w-full pl-12 pr-4 py-4 rounded-2xl bg-white border border-blue-200 focus:border-blue-600 outline-none font-black text-lg text-blue-700 shadow-inner transition-all">
                        </div>
                        <p class="text-[9px] font-black text-blue-500 mt-2 uppercase tracking-widest text-center">Profit: Rp {{ form.price - form.cost_price }} / {{ form.satuan_dasar }}</p>
                    </div>

                    <div class="md:col-span-2 grid grid-cols-1 gap-4 bg-slate-50 p-5 rounded-[24px] border border-slate-200">
                        <div v-if="form.has_satuan_besar" class="grid grid-cols-1 sm:grid-cols-3 gap-4 items-center">
                            <div>
                                <label class="text-[9px] font-black text-indigo-600 uppercase tracking-widest mb-2 block">Jumlah {{ form.satuan_besar || 'KARTON' }}</label>
                                <input :value="stokDalamKarton" @input="emit('update:stokDalamKarton', $event.target.value)" type="number" min="0" placeholder="0" class="w-full px-4 py-3.5 rounded-xl bg-white border border-slate-200 focus:border-indigo-600 outline-none font-black text-lg text-indigo-600 transition-all shadow-sm">
                            </div>
                            <div>
                                <label class="text-[9px] font-black text-amber-600 uppercase tracking-widest mb-2 block">+ Lebih Eceran ({{ form.satuan_dasar }})</label>
                                <input :value="eceranTambahan" @input="emit('update:eceranTambahan', $event.target.value)" type="number" min="0" placeholder="0" class="w-full px-4 py-3.5 rounded-xl bg-white border border-slate-200 focus:border-amber-600 outline-none font-black text-lg text-amber-600 transition-all shadow-sm">
                            </div>
                            <div>
                                <label class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-2 block">Total Stok Akhir</label>
                                <input :value="form.stock" type="number" disabled class="w-full px-4 py-3.5 rounded-xl bg-slate-100 border border-transparent font-black text-xl text-slate-500 text-center">
                            </div>
                        </div>
                        <div v-else>
                            <label class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-2 block text-center">Stok Awal Fisik (dalam {{ form.satuan_dasar }})</label>
                            <input v-model.number="form.stock" type="number" min="0" placeholder="0" class="w-full px-4 py-4 rounded-2xl bg-white border-2 border-slate-200 focus:border-blue-600 outline-none font-black text-center text-2xl text-slate-800 transition-all">
                        </div>
                    </div>

                </div>
            </div>

            <div class="p-6 bg-slate-50 border-t border-slate-100 shrink-0">
                <button @click="emit('submit')" :disabled="isSubmitting" class="w-full py-5 font-black text-xs uppercase tracking-[0.2em] bg-blue-600 text-white rounded-[24px] shadow-xl shadow-blue-200 hover:bg-slate-900 transition-all active:scale-95 flex items-center justify-center gap-3 disabled:opacity-50">
                    <template v-if="isSubmitting">
                        <div class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                        Menyimpan Data...
                    </template>
                    <template v-else>
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/><polyline points="17 21 17 13 7 13 7 21"/><polyline points="7 3 7 8 15 8"/></svg>
                        Simpan Perubahan Produk
                    </template>
                </button>
            </div>
        </div>
    </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { height: 6px; width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #94a3b8; }
input[type=number]::-webkit-inner-spin-button, 
input[type=number]::-webkit-outer-spin-button { -webkit-appearance: none; margin: 0; }
input[type=number] { -moz-appearance: textfield; }
</style>