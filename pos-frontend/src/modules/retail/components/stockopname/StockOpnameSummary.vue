<script setup>
defineProps({
    notes: String,
    cartLength: Number,
    isSubmitting: Boolean,
    isOwner: Boolean, 
    soStep: String,
    activeTab: String
});
defineEmits(['update:notes', 'submit', 'proceed', 'back']);
</script>

<template>
    <div class="space-y-6">
        <div v-if="activeTab === 'SO'" class="bg-white p-8 rounded-[32px] shadow-xl shadow-slate-200/50 border border-slate-100">
            <h3 class="font-black text-slate-800 text-[10px] uppercase tracking-[0.2em] mb-6 flex items-center gap-3">
                <span class="w-1.5 h-6 bg-indigo-600 rounded-full"></span> Catatan Audit
            </h3>
            <textarea :value="notes" @input="$emit('update:notes', $event.target.value)" rows="3" placeholder="Keterangan SO..." class="w-full p-4 bg-slate-50 border-2 border-slate-100 focus:border-indigo-500 rounded-2xl outline-none font-bold text-sm transition-all text-slate-700 resize-none shadow-inner"></textarea>
            
            <div class="mt-6 p-4 bg-indigo-50 rounded-2xl border border-indigo-100">
                <h4 class="text-[9px] font-black text-indigo-600 uppercase tracking-widest mb-1">Total SKU Dihitung</h4>
                <p class="text-2xl font-black text-indigo-900">{{ cartLength }} <span class="text-sm">Item</span></p>
            </div>
        </div>

        <template v-if="activeTab === 'KLAIM'">
            <button @click="$emit('submit')" :disabled="isSubmitting" class="w-full bg-amber-500 hover:bg-amber-600 text-white py-6 rounded-[30px] font-black text-xs uppercase tracking-[0.2em] shadow-xl shadow-amber-200 transition-all active:scale-[0.98]">
                {{ isSubmitting ? 'MENGIRIM...' : 'KIRIM KLAIM BARANG' }}
            </button>
        </template>

        <template v-else>
            <div v-if="soStep === 'REVIEW'" class="flex flex-col gap-4">
                <button @click="$emit('submit')" :disabled="isSubmitting" class="w-full bg-slate-900 hover:bg-slate-800 text-white py-6 rounded-[30px] font-black text-[10px] sm:text-xs uppercase tracking-[0.2em] shadow-xl shadow-slate-200 transition-all active:scale-[0.98] flex justify-center items-center gap-2">
                    <div v-if="isSubmitting" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                    <span v-else>KIRIM PENGAJUAN AUDIT</span>
                </button>
                <button @click="$emit('back')" class="w-full text-slate-400 py-3 font-black text-[9px] uppercase tracking-widest hover:text-slate-600 transition-all">
                    Kembali Edit Data
                </button>
            </div>

            <div v-else class="flex flex-col gap-4">
                <button @click="$emit('proceed')" class="w-full bg-indigo-600 hover:bg-indigo-700 text-white py-6 rounded-[30px] font-black text-xs sm:text-sm uppercase tracking-[0.2em] shadow-xl shadow-indigo-200 transition-all active:scale-[0.98]">
                    CEK SELISIH (REVIEW MODE)
                </button>
                <div class="p-4 bg-indigo-50 border border-indigo-100 rounded-2xl">
                    <p class="text-[10px] font-bold text-indigo-800 text-center uppercase tracking-widest">Selesaikan hitung fisik sebelum review.</p>
                </div>
            </div>
        </template>
    </div>
</template>