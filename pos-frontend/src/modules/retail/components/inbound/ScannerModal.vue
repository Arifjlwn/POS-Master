<script setup>
defineProps({ show: Boolean });
defineEmits(['close']);
</script>

<template>
    <div v-if="show" class="fixed inset-0 bg-slate-900/90 flex items-center justify-center z-[150] p-4 backdrop-blur-sm no-print">
        <div class="bg-white rounded-[32px] shadow-2xl w-full max-w-md overflow-hidden flex flex-col border-[6px] border-slate-800">
            <div class="p-6 border-b border-slate-50 flex justify-between items-center bg-slate-50/50">
                <h2 class="text-lg font-black text-slate-800 uppercase tracking-widest flex items-center gap-2">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm14 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z" /></svg>
                    Scan Barcode
                </h2>
                <button @click="$emit('close')" class="p-2 rounded-xl bg-slate-100 text-slate-400 hover:text-red-500 hover:bg-red-50 transition-all">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                </button>
            </div>
            <div class="p-6 bg-black relative">
                <div id="reader-lpb" class="w-full rounded-2xl overflow-hidden border-2 border-slate-700"></div>
                <p class="text-[10px] font-bold text-white/50 text-center mt-4 uppercase tracking-widest animate-pulse">Arahkan kamera ke barcode</p>
            </div>
        </div>
    </div>
</template>