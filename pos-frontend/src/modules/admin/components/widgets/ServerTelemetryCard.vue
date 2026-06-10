<script setup>
defineProps({
    health: {
        type: Object,
        required: true
    }
});
</script>

<template>
    <div class="bg-[#131B2E] border border-slate-800 rounded-[28px] p-6 shadow-xl flex flex-col justify-between h-full">
        <div>
            <div class="mb-4 border-b border-slate-800 pb-4">
                <h3 class="text-sm font-black text-white uppercase tracking-wide">Core Engine Telemetry</h3>
                <p class="text-[10px] text-slate-500 font-bold uppercase tracking-wider mt-0.5">Cloud Database & Server Performance</p>
            </div>

            <div class="space-y-5">
                <div class="space-y-2">
                    <div class="flex justify-between items-center text-[11px] font-bold uppercase tracking-wider text-slate-400">
                        <span>Core CPU Allocation</span>
                        <span class="text-indigo-400 font-black">{{ health.cpu_usage }}%</span>
                    </div>
                    <div class="w-full h-2 bg-[#0B0F19] rounded-full overflow-hidden border border-slate-800">
                        <div class="h-full bg-indigo-500 rounded-full transition-all duration-500" :style="{ width: health.cpu_usage + '%' }"></div>
                    </div>
                </div>

                <div class="space-y-2">
                    <div class="flex justify-between items-center text-[11px] font-bold uppercase tracking-wider text-slate-400">
                        <span>Memory Heap Usage</span>
                        <span class="text-purple-400 font-black">{{ health.ram_usage }}%</span>
                    </div>
                    <div class="w-full h-2 bg-[#0B0F19] rounded-full overflow-hidden border border-slate-800">
                        <div class="h-full bg-purple-500 rounded-full transition-all duration-500" :style="{ width: health.ram_usage + '%' }"></div>
                    </div>
                </div>

                <div class="grid grid-cols-2 gap-3 pt-3 border-t border-slate-800/60">
                    <div class="p-3 bg-[#0B0F19]/50 rounded-xl border border-slate-800/40 text-center">
                        <span class="text-[9px] font-bold text-slate-500 uppercase tracking-widest block">Supabase GORM</span>
                        <span class="text-xs font-black text-emerald-400 uppercase tracking-wider mt-1 block">{{ health.db_status }}</span>
                    </div>
                    <div class="p-3 bg-[#0B0F19]/50 rounded-xl border border-slate-800/40 text-center">
                        <span class="text-[9px] font-bold text-slate-500 uppercase tracking-widest block">Core Latency</span>
                        <span class="text-xs font-black text-indigo-400 uppercase tracking-wider mt-1 block">{{ health.latency }}</span>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>