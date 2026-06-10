<script setup>
// 🚀 Menggunakan prop agar data bersifat dinamis dan dapat dikendalikan dari halaman utama (reactive)
defineProps({
    activities: {
        type: Array,
        required: true
    }
});
</script>

<template>
    <div class="bg-[#131B2E] border border-slate-800 rounded-[24px] sm:rounded-[28px] p-5 sm:p-6 shadow-xl flex flex-col justify-between h-full w-full overflow-hidden">
        <div>
            <div class="flex items-center justify-between mb-4 border-b border-slate-800 pb-4">
                <div class="min-w-0 flex-1 pr-2">
                    <h3 class="text-xs sm:text-sm font-black text-white uppercase tracking-wide truncate">Live Audit Logs Activity</h3>
                    <p class="text-[9px] sm:text-[10px] text-slate-500 font-bold uppercase tracking-wider mt-0.5 truncate">Realtime Global Event Stream Detector</p>
                </div>
                <span class="w-2 h-2 rounded-full bg-indigo-500 animate-ping shrink-0"></span>
            </div>

            <div class="space-y-3 max-h-[260px] overflow-y-auto pr-1 custom-scrollbar">
                <div 
                    v-for="log in activities" 
                    :key="log.id" 
                    class="flex items-start gap-2.5 p-2.5 sm:p-3 bg-[#0B0F19]/60 border border-slate-800/40 rounded-xl transition-all hover:bg-[#0B0F19]"
                >
                    <span class="text-[9px] sm:text-[10px] font-black text-indigo-400 bg-indigo-500/10 px-1.5 py-0.5 rounded-md mt-0.5 shrink-0">
                        {{ log.time }}
                    </span>
                    <p class="text-[11px] sm:text-xs text-slate-300 font-medium leading-relaxed flex-1 break-words overflow-hidden">
                        <span v-if="log.type === 'warning'" class="text-amber-400 font-bold shrink-0">[WARN] </span>
                        <span v-else-if="log.type === 'danger'" class="text-rose-400 font-bold shrink-0">[CRIT] </span>
                        {{ log.event }}
                    </p>
                </div>
            </div>
        </div>
    </div>
</template>