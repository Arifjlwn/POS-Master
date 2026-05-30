<script setup>
import { ref } from 'vue';

const props = defineProps({ form: Object, qrisPreview: String });
const emit = defineEmits(['update-file']);
const qrisInput = ref(null);

const onQrisSelect = (e) => {
    const file = e.target.files[0];
    if (file) emit('update-file', 'qris', file, URL.createObjectURL(file));
};
</script>

<template>
    <div class="space-y-6 animate-fade-in-up">
        <div class="bg-indigo-50 p-4 rounded-2xl border border-indigo-100 flex items-start gap-4">
            <div class="p-2 bg-indigo-100 text-indigo-600 rounded-xl shrink-0"><svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg></div>
            <div>
                <h4 class="font-black text-indigo-800 text-xs uppercase tracking-widest mb-1">Mode Pembayaran QRIS Statis</h4>
                <p class="text-xs font-medium text-indigo-600 leading-relaxed">Pembeli akan melakukan scan gambar QRIS di bawah ini, dan kasir akan melakukan validasi mutasi secara manual sebelum menyelesaikan transaksi.</p>
            </div>
        </div>

        <div class="flex flex-col md:flex-row gap-6">
            <div class="w-full md:w-1/2">
                <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Nama Pemilik Rekening (A.N)</label>
                <input v-model="form.qris_name" type="text" placeholder="Misal: Budi Santoso" class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:border-blue-600 outline-none font-bold text-sm mb-6">
            </div>

            <div class="w-full md:w-1/2 flex flex-col items-center p-6 bg-slate-50 rounded-3xl border-2 border-slate-100">
                <label class="text-[10px] font-black text-slate-500 uppercase tracking-widest mb-4">Upload Barcode QRIS</label>
                <div @click="$refs.qrisInput.click()" class="w-48 h-48 bg-white rounded-3xl shadow-sm flex items-center justify-center cursor-pointer overflow-hidden border-4 border-slate-200 group relative hover:border-emerald-400 transition-all">
                    <img v-if="qrisPreview" :src="qrisPreview" class="w-full h-full object-contain p-2">
                    <div v-else class="text-center text-slate-300 group-hover:text-emerald-500 transition-colors">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-10 h-10 mx-auto mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm14 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z" /></svg>
                        <p class="text-[10px] font-black uppercase tracking-widest">Pilih Gambar</p>
                    </div>
                </div>
                <input type="file" ref="qrisInput" @change="onQrisSelect" class="hidden" accept="image/*">
            </div>
        </div>
    </div>
</template>