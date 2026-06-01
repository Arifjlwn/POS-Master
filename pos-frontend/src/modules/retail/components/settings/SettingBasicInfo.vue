<script setup>
import { ref } from 'vue';

const props = defineProps({ form: Object, logoPreview: String });
const emit = defineEmits(['update-file', 'remove-logo']); // 🚀 Tambahin emit remove-logo
const fileInput = ref(null);

const onLogoSelect = (e) => {
    const file = e.target.files[0];
    if (file) emit('update-file', 'logo', file, URL.createObjectURL(file));
};
</script>

<template>
    <div class="space-y-6 animate-fade-in-up">
        <div class="flex flex-col md:flex-row gap-6">
            <div class="w-full md:w-1/3 flex flex-col items-center p-6 bg-slate-50 rounded-3xl border-2 border-slate-100">
                <label class="text-[10px] font-black text-slate-500 uppercase tracking-widest mb-4">Logo Toko</label>
                
                <div class="relative w-32 h-32 bg-white rounded-[24px] shadow-sm flex items-center justify-center overflow-hidden border-2 border-slate-200 group hover:border-blue-400 transition-colors">
                    
                    <div v-if="logoPreview" class="w-full h-full relative cursor-default">
                        <img :src="logoPreview" class="w-full h-full object-contain p-2">
                        
                        <button @click.prevent="emit('remove-logo')" class="absolute inset-0 bg-rose-500/80 text-white flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity backdrop-blur-sm cursor-pointer z-10">
                            <svg class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/></svg>
                        </button>
                    </div>

                    <div v-else @click="$refs.fileInput.click()" class="w-full h-full flex items-center justify-center cursor-pointer">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8 text-slate-300 group-hover:text-blue-500 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" /></svg>
                    </div>

                </div>
                
                <input type="file" ref="fileInput" @change="onLogoSelect" class="hidden" accept="image/*">
                <p class="text-[9px] text-slate-400 font-bold mt-4 text-center">Format: JPG, PNG (Max 2MB)</p>
            </div>

            <div class="w-full md:w-2/3 grid grid-cols-1 md:grid-cols-2 gap-4">
                <div class="md:col-span-2">
                    <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Nama Toko</label>
                    <input v-model="form.nama_toko" type="text" class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:border-blue-600 outline-none font-bold text-sm">
                </div>
                <div>
                    <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">No. WA Bisnis</label>
                    <input v-model="form.telepon" type="text" class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:border-blue-600 outline-none font-bold text-sm">
                </div>
                <div>
                    <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Kategori Bisnis</label>
                    <input v-model="form.business_type" type="text" disabled class="w-full p-4 bg-slate-100 text-slate-500 rounded-2xl border-2 border-slate-100 outline-none font-black text-xs uppercase tracking-widest cursor-not-allowed">
                </div>
            </div>
        </div>

        <div class="border-t border-slate-100 pt-6 space-y-4">
            <h4 class="font-black text-slate-800 uppercase tracking-widest text-sm mb-2">Lokasi Operasional</h4>
            
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                <div class="md:col-span-2 lg:col-span-3">
                    <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Alamat Jalan / Patokan</label>
                    <textarea v-model="form.alamat" rows="2" class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:border-blue-600 outline-none font-bold text-sm resize-none" placeholder="Contoh: Jl. Sudirman No. 123..."></textarea>
                </div>
                
                <div>
                    <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Provinsi</label>
                    <input v-model="form.provinsi" type="text" class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:border-blue-600 outline-none font-bold text-sm">
                </div>
                <div>
                    <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Kota / Kabupaten</label>
                    <input v-model="form.kota" type="text" class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:border-blue-600 outline-none font-bold text-sm">
                </div>
                <div>
                    <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Kecamatan</label>
                    <input v-model="form.kecamatan" type="text" class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:border-blue-600 outline-none font-bold text-sm">
                </div>
                <div>
                    <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Kelurahan / Desa</label>
                    <input v-model="form.kelurahan" type="text" class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:border-blue-600 outline-none font-bold text-sm">
                </div>
                <div>
                    <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Kode Pos</label>
                    <input v-model="form.kode_pos" type="text" class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:border-blue-600 outline-none font-bold text-sm">
                </div>
            </div>
        </div>
    </div>
</template>