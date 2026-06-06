<script setup>
defineProps({
    notes: String,
    cartLength: Number,
    isSubmitting: Boolean,
    isOwner: Boolean, 
    soStep: String,
    activeTab: String
});

// 🚀 FIX EMITS: Daftarkan emit 'back-step' untuk memicu pembalikan state di parent
defineEmits(['update:notes', 'submit', 'proceed', 'back-step']);
</script>

<template>
    <div class="space-y-6 select-none font-sans">
        
        <div v-if="activeTab === 'SO'" class="bg-white p-5 md:p-6 rounded-[32px] shadow-sm border border-slate-100">
            <h3 class="font-black text-slate-800 text-[10px] uppercase tracking-[0.2em] mb-4 flex items-center gap-2">
                <span class="w-1.5 h-4 bg-indigo-600 rounded-full"></span> Catatan Laporan Audit
            </h3>
            <textarea 
                :value="notes" 
                @input="$emit('update:notes', $event.target.value)" 
                rows="3" 
                placeholder="Masukkan keterangan stock opname hari ini ..." 
                class="w-full p-4 bg-slate-50 border-2 border-slate-100 focus:border-indigo-500 focus:bg-white rounded-2xl outline-none font-bold text-xs md:text-sm transition-all text-slate-700 resize-none shadow-inner"
            ></textarea>
            
            <div class="mt-4 p-4 bg-indigo-50 rounded-2xl border border-indigo-100/70 flex items-center justify-between">
                <h4 class="text-[9px] font-black text-indigo-600 uppercase tracking-widest leading-none">Total SKU Dihitung</h4>
                <p class="text-xl font-black text-indigo-900 leading-none">{{ cartLength }} <span class="text-xs font-bold text-indigo-500">Item</span></p>
            </div>
        </div>

        <template v-if="activeTab === 'KLAIM'">
            <button 
                @click="$emit('submit')" 
                :disabled="isSubmitting" 
                class="w-full bg-amber-500 hover:bg-amber-600 disabled:bg-amber-300 text-white py-5 rounded-[30px] font-black text-xs uppercase tracking-[0.2em] shadow-xl shadow-amber-100 transition-all active:scale-[0.97] flex justify-center items-center gap-2 disabled:cursor-not-allowed"
            >
                <div v-if="isSubmitting" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                <span>{{ isSubmitting ? 'MEMPROSES KLAIM...' : 'KIRIM KLAIM BARANG TEMUAN' }}</span>
            </button>
        </template>

        <template v-else>
            <div v-if="soStep === 'REVIEW'" class="flex flex-col gap-3 animate-fade-in">
                <button 
                    @click="$emit('submit')" 
                    :disabled="isSubmitting" 
                    class="w-full bg-slate-900 hover:bg-slate-800 disabled:bg-slate-400 text-white py-5 rounded-[30px] font-black text-xs uppercase tracking-[0.2em] shadow-xl shadow-slate-200/50 transition-all active:scale-[0.97] flex justify-center items-center gap-2 disabled:cursor-not-allowed"
                >
                    <div v-if="isSubmitting" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                    <template v-else>
                        <span>{{ isOwner ? 'FINALISASI & PERBARUI STOK MASTER' : 'KIRIM BERKAS PENGAJUAN AUDIT' }}</span>
                    </template>
                </button>
                
                <button @click="$emit('back-step')" class="w-full text-slate-400 hover:text-indigo-600 py-2 font-black text-[9px] uppercase tracking-widest transition-colors flex items-center justify-center gap-1.5 active:scale-95">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M10 19l-7-7m0 0l7-7m-7 7h18" /></svg>
                    Kembali Edit Angka Fisik Rak
                </button>
            </div>

            <div v-else class="flex flex-col gap-4">
                <button 
                    @click="$emit('proceed')" 
                    class="w-full bg-indigo-600 hover:bg-indigo-700 text-white py-5 rounded-[30px] font-black text-xs md:text-sm uppercase tracking-[0.15em] shadow-xl shadow-indigo-100 transition-all active:scale-[0.97]"
                >
                    CEK SELISIH (REVIEW MODE)
                </button>
                <div class="p-4 bg-indigo-50/50 border border-indigo-100/50 rounded-2xl">
                    <p class="text-[10px] font-bold text-indigo-700 text-center uppercase tracking-widest leading-relaxed">Selesaikan hitung fisik rak gudang secara jujur sebelum masuk ke menu review selisih.</p>
                </div>
            </div>
        </template>
    </div>
</template>

<style scoped>
.animate-fade-in { animation: fadeIn 0.2s cubic-bezier(0.16, 1, 0.3, 1); }
@keyframes fadeIn { from { opacity: 0; transform: translateY(5px); } to { opacity: 1; transform: translateY(0); } }
</style>