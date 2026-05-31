<script setup>
defineProps({
    currentUser: Object,
    currentSession: Object,
    currentTime: String
});

const emit = defineEmits(['go-dashboard', 'logout']);
</script>

<template>
    <div class="p-2 md:p-4 shrink-0 pb-0 no-print">
        <div class="bg-indigo-900 rounded-[20px] md:rounded-[24px] shadow-lg flex flex-col md:flex-row overflow-hidden border border-indigo-800">
            
            <div @click="emit('go-dashboard')" class="p-3 md:p-4 bg-indigo-950 md:w-64 flex flex-col justify-center items-center cursor-pointer hover:bg-indigo-900 transition-colors border-b md:border-b-0 md:border-r border-indigo-800 shrink-0">
                <div class="font-black text-xl md:text-2xl tracking-tight text-white flex items-center gap-2">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 md:w-8 md:h-8 text-indigo-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" /></svg>
                    NEXA POS
                </div>
                <div class="text-center text-[9px] md:text-[10px] uppercase mt-1 md:mt-2 tracking-widest font-black px-3 py-1 rounded-lg text-amber-400 bg-amber-400/10 hover:bg-amber-400/20 transition-colors">
                    MENU UTAMA
                </div>
            </div>

            <div class="flex-1 flex flex-col min-w-0">
                <div class="bg-rose-600 text-white text-[9px] md:text-[10px] font-black py-1.5 text-center tracking-[0.2em] overflow-hidden whitespace-nowrap shadow-inner flex w-full">
                    <div class="animate-marquee min-w-full">
                        INFO KASIR: PASTIKAN SCAN BARCODE DENGAN BENAR | SELALU TAWARKAN PRODUK PROMO KEPADA KONSUMEN SEBELUM CHECKOUT
                    </div>
                </div>

                <div class="flex flex-wrap md:flex-nowrap justify-between items-center px-3 md:px-5 py-2 md:py-3 flex-1 bg-indigo-900 gap-2 md:gap-4">
                    <div class="flex items-center gap-2 md:gap-4 text-xs md:text-sm font-medium text-indigo-100 min-w-0">
                        <div class="flex items-center gap-1.5 md:gap-2 bg-indigo-950/50 px-2.5 py-1.5 rounded-lg md:rounded-xl border border-indigo-700/50 truncate max-w-[120px] md:max-w-[200px]">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 md:w-4 md:h-4 text-indigo-400 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" /></svg>
                            <span class="font-black text-white uppercase tracking-wider text-[10px] md:text-xs truncate" :title="currentUser?.name">{{ currentUser?.name ? currentUser.name.split(' ')[0] : 'KASIR' }}</span>
                        </div>
                        
                        <div class="flex items-center gap-3 sm:gap-5 border-l border-indigo-800 pl-3 sm:pl-5">
                            <span class="flex items-center gap-1 sm:gap-2 opacity-90 text-[10px] uppercase tracking-widest font-black whitespace-nowrap">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 md:w-4 md:h-4 text-indigo-400 hidden sm:block" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                                <span>Shift:</span><strong class="text-white text-xs ml-1">1</strong>
                            </span>
                            <span class="flex items-center gap-1 sm:gap-2 opacity-90 text-[10px] uppercase tracking-widest font-black whitespace-nowrap">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 md:w-4 md:h-4 text-indigo-400 hidden sm:block" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" /></svg>
                                <span>POS:</span><strong class="text-white text-xs ml-1">{{ currentSession?.station_number || '01' }}</strong>
                            </span>
                        </div>
                    </div>

                    <div class="flex items-center gap-2 md:gap-3 shrink-0">
                        <div class="font-mono text-xs lg:text-sm font-black text-amber-400 tracking-wider bg-indigo-950/50 px-2 md:px-3 py-1.5 rounded-lg md:rounded-xl border border-indigo-800/50 hidden lg:block">
                            {{ currentTime }}
                        </div>
                        <button @click="emit('logout')" class="bg-rose-500 hover:bg-rose-600 text-white p-2 px-3 md:px-4 md:py-2 rounded-lg md:rounded-xl transition-all shadow-lg shadow-rose-900/50 active:scale-95 flex items-center gap-2 border border-rose-400 whitespace-nowrap">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 md:w-5 md:h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" /></svg>
                            <span class="text-[9px] md:text-[10px] font-black uppercase tracking-widest">Tutup Shift</span>
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.animate-marquee {
    display: inline-block;
    padding-left: 100%;
    animation: marquee 20s linear infinite;
}
@keyframes marquee {
    0% { transform: translateX(0); }
    100% { transform: translateX(-100%); }
}
</style>