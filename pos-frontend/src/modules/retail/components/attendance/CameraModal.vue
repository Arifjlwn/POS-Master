<script setup>
defineProps({
    show: Boolean,
    target: Object,
    currentTime: String,
    isAiLoading: Boolean,
    isSubmitting: Boolean,
    setVideoRef: Function // Menerima function binding dari composable
});

defineEmits(['close', 'capture']);
</script>

<template>
    <div v-if="show" class="fixed inset-0 bg-slate-950/95 z-[100] flex items-center justify-center p-4 backdrop-blur-xl no-print">
        <div class="bg-white rounded-[32px] md:rounded-[40px] overflow-hidden shadow-2xl w-full max-w-sm border-[8px] md:border-[10px] border-slate-900/5 flex flex-col">
            <div class="p-5 md:p-6 border-b border-slate-50 flex justify-between items-center bg-white shrink-0">
                <div>
                    <h3 class="font-black text-slate-800 uppercase tracking-tighter text-lg md:text-xl italic">Verifikasi Wajah</h3>
                    <p class="text-[9px] md:text-[10px] text-indigo-600 font-black uppercase tracking-widest mt-0.5">{{ target.jenis }} • {{ target.nama }}</p>
                </div>
                <button @click="$emit('close')" class="w-10 h-10 rounded-[14px] bg-slate-100 text-slate-400 hover:text-rose-500 hover:bg-rose-50 transition-all flex items-center justify-center">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                </button>
            </div>
            
            <div class="relative bg-slate-900 w-full aspect-[3/4] flex items-center justify-center overflow-hidden">
                <div v-if="isAiLoading" class="absolute inset-0 bg-slate-950/90 backdrop-blur-md flex flex-col items-center justify-center z-50 text-white">
                    <div class="w-10 h-10 md:w-12 md:h-12 border-4 border-indigo-500 border-t-transparent rounded-full animate-spin mb-4"></div>
                    <p class="font-black text-[9px] md:text-[10px] uppercase tracking-[0.2em] animate-pulse">Menyiapkan Engine AI...</p>
                </div>
                
                <video :ref="setVideoRef" autoplay playsinline class="w-full h-full object-cover transform scale-x-[-1]"></video>
                
                <div class="absolute inset-0 border-[12px] md:border-[16px] border-black/20 pointer-events-none"></div>
                <div class="absolute inset-x-8 inset-y-12 md:inset-x-12 md:inset-y-16 border-2 border-white/30 rounded-[32px] md:rounded-[40px] pointer-events-none">
                    <div class="absolute top-0 left-0 w-8 h-8 border-t-4 border-l-4 border-white/80 rounded-tl-[30px] md:rounded-tl-[38px]"></div>
                    <div class="absolute top-0 right-0 w-8 h-8 border-t-4 border-r-4 border-white/80 rounded-tr-[30px] md:rounded-tr-[38px]"></div>
                    <div class="absolute bottom-0 left-0 w-8 h-8 border-b-4 border-l-4 border-white/80 rounded-bl-[30px] md:rounded-bl-[38px]"></div>
                    <div class="absolute bottom-0 right-0 w-8 h-8 border-b-4 border-r-4 border-white/80 rounded-br-[30px] md:rounded-br-[38px]"></div>
                </div>

                <div class="absolute bottom-4 left-4 md:bottom-6 md:left-6 text-amber-400 font-mono drop-shadow-md">
                    <div class="font-black tracking-widest text-[10px] md:text-xs">{{ target.nama.toUpperCase() }}</div>
                    <div class="text-lg md:text-xl font-black tracking-widest mt-0.5 md:mt-1">{{ currentTime }}</div>
                    <div class="font-bold opacity-80 mt-0.5 text-[9px] md:text-[10px]">{{ new Date().toLocaleDateString('id-ID') }}</div>
                </div>
            </div>
            
            <div class="p-5 md:p-6 bg-white shrink-0">
                <button @click="$emit('capture')" :disabled="isSubmitting" class="w-full bg-indigo-600 hover:bg-slate-900 text-white py-4 md:py-5 rounded-[20px] md:rounded-[24px] font-black text-xs md:text-sm uppercase tracking-[0.2em] shadow-xl shadow-indigo-200 transition-all disabled:opacity-50 flex items-center justify-center gap-2 md:gap-3 active:scale-95">
                    <template v-if="isSubmitting">
                        <div class="w-4 h-4 md:w-5 md:h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                        Verifikasi...
                    </template>
                    <template v-else>
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 md:w-5 md:h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M14.5 4h-5L7 7H4a2 2 0 0 0-2 2v9a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2h-3l-2.5-3z"/><circle cx="12" cy="13" r="3"/></svg>
                        Scan & Kirim
                    </template>
                </button>
            </div>
        </div>
    </div>
</template>