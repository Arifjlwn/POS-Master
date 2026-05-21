<script setup>
defineProps({
    totalMenu: Number,
    totalAvailable: Number,
    totalUnavailable: Number
})
</script>

<template>
    <div
        class="grid grid-cols-3 gap-2 mb-5 md:gap-4 lg:gap-5"
    >
        <!-- TOTAL -->
        <div
            class="mobile-stat p-3 md:p-5 border rounded-[20px] md:rounded-[28px] bg-white/70 backdrop-blur-xl border-white/40 shadow-[0_10px_40px_rgba(15,23,42,0.06)]"
        >
            <p
                class="text-[8px] md:text-[11px] font-black uppercase tracking-[0.18em] md:tracking-[0.25em] text-slate-400"
            >
                Total Menu
            </p>

            <h2
                class="mt-1 md:mt-2 text-[22px] md:text-4xl font-black text-slate-800 leading-none"
            >
                {{ totalMenu }}
            </h2>
        </div>

        <!-- READY -->
        <div
            class="mobile-stat p-3 md:p-5 border rounded-[20px] md:rounded-[28px] bg-white/70 backdrop-blur-xl border-white/40 shadow-[0_10px_40px_rgba(15,23,42,0.06)]"
        >
            <p
                class="text-[8px] md:text-[11px] font-black uppercase tracking-[0.18em] md:tracking-[0.25em] text-slate-400"
            >
                Menu Ready
            </p>

            <h2
                class="mt-1 md:mt-2 text-[22px] md:text-4xl font-black text-emerald-600 leading-none"
            >
                {{ totalAvailable }}
            </h2>
        </div>

        <!-- HABIS -->
        <div
            class="mobile-stat p-3 md:p-5 border rounded-[20px] md:rounded-[28px] bg-white/70 backdrop-blur-xl border-white/40 shadow-[0_10px_40px_rgba(15,23,42,0.06)]"
        >
            <p
                class="text-[8px] md:text-[11px] font-black uppercase tracking-[0.18em] md:tracking-[0.25em] text-slate-400"
            >
                Menu Habis
            </p>

            <h2
                class="mt-1 md:mt-2 text-[22px] md:text-4xl font-black text-rose-600 leading-none"
            >
                {{ totalUnavailable }}
            </h2>
        </div>
    </div>
</template>