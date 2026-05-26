<script setup>
defineProps({ me: Object });
</script>

<template>
    <div class="bg-gradient-to-br from-indigo-900 via-blue-900 to-slate-900 rounded-[32px] p-6 md:p-10 mb-8 text-white shadow-2xl flex flex-col md:flex-row items-center justify-between overflow-hidden relative border border-indigo-800">
        <div class="absolute -right-10 -top-10 opacity-10 text-white pointer-events-none">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-72 h-72" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
        </div>
        
        <div class="relative z-10 text-center md:text-left mb-6 md:mb-0 w-full md:w-auto">
            <h1 class="text-3xl md:text-4xl font-black tracking-tighter mb-1 uppercase italic">Presensi <span class="text-blue-400">Tim</span></h1>
            <p class="text-indigo-200 font-bold text-[10px] md:text-xs uppercase tracking-[0.2em] flex items-center justify-center md:justify-start gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M14 9a2 2 0 0 1-2 2H6l-4 4V4c0-1.1.9-2 2-2h8a2 2 0 0 1 2 2v5Z"/><path d="M18 9h2a2 2 0 0 1 2 2v11l-4-4h-6a2 2 0 0 1-2-2v-1"/></svg>
                Absen Tepat Waktu Ya!
            </p>
        </div>
        
        <div class="relative z-10 bg-white/10 backdrop-blur-md border border-white/20 rounded-2xl p-4 md:px-8 text-center shadow-inner w-full md:w-auto">
            <div class="text-[9px] md:text-[10px] font-black text-indigo-200 uppercase tracking-[0.2em] mb-1">Status Shift Anda</div>
            <div class="font-mono text-2xl md:text-3xl font-black flex items-center justify-center gap-2 tracking-tighter"
                 :class="me.shift_hari_ini === 'OFF' ? 'text-slate-400' : (me.sudah_masuk ? (me.sudah_pulang ? 'text-slate-300' : 'text-emerald-400') : 'text-amber-400')">
                <template v-if="me.shift_hari_ini === 'OFF'">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 md:w-8 md:h-8" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/><path d="m9 16 2 2 4-4"/></svg>
                    LIBUR (OFF)
                </template>
                <template v-else>
                    <template v-if="me.sudah_masuk">
                        <template v-if="me.sudah_pulang">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 md:w-8 md:h-8" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg>
                            SELESAI
                        </template>
                        <template v-else>
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 md:w-8 md:h-8 animate-pulse" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
                            AKTIF
                        </template>
                    </template>
                    <template v-else>
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 md:w-8 md:h-8" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
                        BELUM
                    </template>
                </template>
            </div>
        </div>
    </div>
</template>