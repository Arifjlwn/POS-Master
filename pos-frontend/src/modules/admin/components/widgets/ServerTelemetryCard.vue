<script setup>
defineProps({
    health: {
        type: Object,
        required: true
    }
});
</script>

<template>
    <div class="bg-[#131B2E] border border-slate-800 rounded-[24px] sm:rounded-[28px] p-5 sm:p-6 shadow-xl flex flex-col justify-between h-full w-full overflow-hidden">
        <div>
            <div class="mb-4 border-b border-slate-800 pb-4">
                <h3 class="text-xs sm:text-sm font-black text-white uppercase tracking-wide truncate">Core Engine Telemetry</h3>
                <p class="text-[9px] sm:text-[10px] text-slate-500 font-bold uppercase tracking-wider mt-0.5 truncate">Cloud Database & Server Performance</p>
            </div>

            <div class="space-y-4 sm:space-y-5">
                <div class="space-y-1.5 sm:space-y-2">
                    <div class="flex justify-between items-center text-[10px] sm:text-[11px] font-bold uppercase tracking-wider text-slate-400">
                        <span class="truncate">Core CPU Allocation</span>
                        <span class="text-indigo-400 font-black ml-1">{{ health.cpu_usage }}%</span>
                    </div>
                    <div class="w-full h-2 bg-[#0B0F19] rounded-full overflow-hidden border border-slate-800">
                        <div class="h-full bg-indigo-500 rounded-full transition-all duration-500" :style="{ width: health.cpu_usage + '%' }"></div>
                    </div>
                </div>

                <div class="space-y-1.5 sm:space-y-2">
                    <div class="flex justify-between items-center text-[10px] sm:text-[11px] font-bold uppercase tracking-wider text-slate-400">
                        <span class="truncate">Memory Heap Usage</span>
                        <span class="text-purple-400 font-black ml-1">{{ health.ram_usage }}%</span>
                    </div>
                    <div class="w-full h-2 bg-[#0B0F19] rounded-full overflow-hidden border border-slate-800">
                        <div class="h-full bg-purple-500 rounded-full transition-all duration-500" :style="{ width: health.ram_usage + '%' }"></div>
                    </div>
                </div>

                <div class="grid grid-cols-2 gap-2 sm:gap-3 pt-3 border-t border-slate-800/60">
                    <div class="p-2 sm:p-3 bg-[#0B0F19]/50 rounded-xl border border-slate-800/40 text-center min-w-0">
                        <span class="text-[8px] sm:text-[9px] font-bold text-slate-500 uppercase tracking-widest block truncate">Supabase GORM</span>
                        <span class="text-[11px] sm:text-xs font-black text-emerald-400 uppercase tracking-wider mt-1 block truncate">{{ health.db_status }}</span>
                    </div>
                    <div class="p-2 sm:p-3 bg-[#0B0F19]/50 rounded-xl border border-slate-800/40 text-center min-w-0">
                        <span class="text-[8px] sm:text-[9px] font-bold text-slate-500 uppercase tracking-widest block truncate">Core Latency</span>
                        <span class="text-[11px] sm:text-xs font-black text-indigo-400 uppercase tracking-wider mt-1 block truncate">{{ health.latency }}</span>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>