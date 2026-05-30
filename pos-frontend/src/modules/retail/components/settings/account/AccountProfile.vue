<script setup>
import { ref } from 'vue';

const props = defineProps({ form: Object, fotoPreview: String });
const emit = defineEmits(['update-file', 'submit']);
const fotoInput = ref(null);

const onFotoSelect = (e) => {
    const file = e.target.files[0];
    if (file) emit('update-file', 'foto', file, URL.createObjectURL(file));
};
</script>

<template>
    <div class="space-y-6 animate-fade-in-up">
        <div class="flex flex-col md:flex-row gap-6">
            <div class="w-full md:w-1/3 flex flex-col items-center p-6 bg-slate-50 rounded-3xl border-2 border-slate-100">
                <label class="text-[10px] font-black text-slate-500 uppercase tracking-widest mb-4">Foto Profil</label>
                <div @click="$refs.fotoInput.click()" class="w-32 h-32 bg-white rounded-full shadow-sm flex items-center justify-center cursor-pointer overflow-hidden border-4 border-slate-200 group relative hover:border-blue-400 transition-colors">
                    <img v-if="fotoPreview" :src="fotoPreview" class="w-full h-full object-cover">
                    <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-8 h-8 text-slate-300 group-hover:text-blue-500 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" /></svg>
                </div>
                <input type="file" ref="fotoInput" @change="onFotoSelect" class="hidden" accept="image/*">
                <p class="text-[9px] text-slate-400 font-bold mt-4 text-center">Tap untuk ubah foto</p>
            </div>

            <div class="w-full md:w-2/3 grid grid-cols-1 md:grid-cols-2 gap-4">
                <div class="md:col-span-2">
                    <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Nama Lengkap</label>
                    <input v-model="form.name" type="text" class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:border-blue-600 outline-none font-bold text-sm">
                </div>
                <div class="md:col-span-2">
                    <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">No. WhatsApp</label>
                    <input v-model="form.no_hp" type="text" class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:border-blue-600 outline-none font-bold text-sm">
                </div>
                <div>
                    <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Tempat Lahir</label>
                    <input v-model="form.tempat_lahir" type="text" class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:border-blue-600 outline-none font-bold text-sm">
                </div>
                <div>
                    <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Tanggal Lahir</label>
                    <input v-model="form.tanggal_lahir" type="date" class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:border-blue-600 outline-none font-bold text-xs uppercase">
                </div>
            </div>
        </div>
        
        <div class="pt-4 border-t border-slate-100">
            <button @click="emit('submit')" class="w-full md:w-auto px-8 py-4 bg-blue-600 hover:bg-slate-900 text-white rounded-2xl font-black text-[10px] uppercase tracking-widest shadow-xl transition-all">
                Simpan Profil
            </button>
        </div>
    </div>
</template>