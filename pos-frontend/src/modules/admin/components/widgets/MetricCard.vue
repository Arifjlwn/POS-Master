<script setup>
// 🚀 Clean Architecture Props
defineProps({
    title: {
        type: String,
        required: true
    },
    value: {
        type: [String, Number],
        required: true
    },
    iconColor: {
        type: String,
        default: 'bg-indigo-500/10 text-indigo-400 border-indigo-500/20'
    }
});
</script>

<template>
    <div class="bg-[#131B2E] border border-slate-800 rounded-[28px] p-6 shadow-xl transition-all duration-300 hover:border-slate-700 hover:-translate-y-1 flex items-center justify-between group">
        <div class="space-y-1.5">
            <span class="text-[10px] font-black text-slate-500 uppercase tracking-widest block">
                {{ title }}
            </span>
            <h3 class="text-2xl md:text-3xl font-black text-white tracking-tight group-hover:text-indigo-400 transition-colors">
                {{ value }}
            </h3>
        </div>
        
        <div :class="iconColor" class="w-12 h-12 rounded-2xl flex items-center justify-center shrink-0 shadow-inner border transition-transform duration-500 group-hover:scale-110">
            <slot name="icon"></slot>
        </div>
    </div>
</template>