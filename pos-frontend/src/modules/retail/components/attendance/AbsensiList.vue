<script setup>
defineProps({ karyawan: Array, currentUser: Object });
defineEmits(['mulai-absen']);

// 🚀 SUNTIKAN URL SANITIZER: Biar foto profil terbaca bener dari backend/cloud bray
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

const formatImageUrl = (path) => {
    if (!path || path === 'null' || path === 'undefined' || path === '') return '';
    if (path.startsWith('data:image/')) return path;
    if (path.startsWith('http://') || path.startsWith('https://')) return path;
    
    const cleanPath = path.startsWith('/') ? path : `/${path}`;
    return `${API_BASE_URL}${cleanPath}`;
};
</script>

<template>
    <div class="mb-10">
        <div class="flex items-center gap-3 mb-6">
            <div class="w-10 h-10 bg-indigo-100 text-indigo-600 rounded-xl flex items-center justify-center shadow-inner">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
            </div>
            <h2 class="text-xl font-black text-slate-800 uppercase tracking-tighter italic">
                Panel Absensi
            </h2>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4 md:gap-6">
            <div v-for="user in karyawan" :key="user.public_id || user.id"
                 class="bg-white rounded-[24px] p-5 shadow-sm border-2 transition-all duration-300 relative overflow-hidden group flex flex-col"
                 :class="String(user.public_id || user.id) === String(currentUser.user_id) ? 'border-indigo-500 shadow-indigo-100 ring-4 ring-indigo-50' : 'border-slate-100 opacity-70 hover:opacity-100'">
                
                <div class="flex items-center gap-4 mb-4 pb-4 border-b border-slate-100 relative z-10">
    <div v-if="user.foto_url" class="w-12 h-12 md:w-14 md:h-14 rounded-[16px] border-2 border-white shadow-md overflow-hidden shrink-0">
        <img :src="formatImageUrl(user.foto_url)" class="w-full h-full object-cover" />
    </div>
    
    <div v-else class="w-12 h-12 md:w-14 md:h-14 rounded-[16px] flex items-center justify-center shrink-0 border-2 border-white shadow-md bg-gradient-to-br from-slate-800 to-slate-900 text-white font-black text-lg uppercase">
        {{ user.name.substring(0, 2) }}
    </div>

    <div class="flex-1 min-w-0">
        <h3 class="font-black text-base md:text-lg leading-tight text-slate-800 uppercase truncate">{{ user.name }}</h3>
        <div class="text-[9px] font-black px-2 py-1 mt-1.5 rounded bg-slate-100 text-slate-500 uppercase tracking-widest inline-flex items-center gap-1 border border-slate-200">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
            {{ user.role }} • {{ user.nik || 'ID:'+(user.public_id || user.id) }}
        </div>
    </div>
</div>

                <div class="flex gap-2 md:gap-3 relative z-10 mt-auto">
                    <div v-if="user.shift_hari_ini === 'OFF'" class="w-full py-3 md:py-3.5 text-center bg-slate-100 text-slate-400 font-black rounded-[16px] text-[10px] md:text-[11px] uppercase tracking-widest border-2 border-slate-200 flex items-center justify-center gap-2 shadow-inner">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/><path d="m9 16 2 2 4-4"/></svg>
                        HARI INI LIBUR (OFF)
                    </div>
                    
                    <template v-else>
                        <button v-if="!user.sudah_masuk" @click="$emit('mulai-absen', user.public_id || user.id, user.name, 'Masuk')"
                            :disabled="String(user.public_id || user.id) !== String(currentUser.user_id)"
                            class="flex-1 py-3 md:py-3.5 rounded-[16px] font-black text-[10px] md:text-[11px] uppercase tracking-widest transition-all border-2 flex items-center justify-center gap-2 bg-emerald-50 text-emerald-700 border-emerald-200 hover:bg-emerald-500 hover:text-white hover:border-emerald-500 disabled:opacity-40 disabled:cursor-not-allowed shadow-sm">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/><polyline points="10 17 15 12 10 7"/><line x1="15" y1="12" x2="3" y2="12"/></svg>
                            Absen Masuk
                        </button>
                        
                        <button v-else-if="!user.sudah_pulang" @click="$emit('mulai-absen', user.public_id || user.id, user.name, 'Pulang')"
                            :disabled="String(user.public_id || user.id) !== String(currentUser.user_id)"
                            class="flex-1 py-3 md:py-3.5 rounded-[16px] font-black text-[10px] md:text-[11px] uppercase tracking-widest transition-all border-2 flex items-center justify-center gap-2 bg-amber-50 text-amber-700 border-amber-200 hover:bg-amber-500 hover:text-white hover:border-amber-500 disabled:opacity-40 disabled:cursor-not-allowed shadow-sm">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
                            Absen Pulang
                        </button>

                        <div v-else class="w-full py-3 md:py-3.5 text-center bg-slate-50 text-slate-400 font-black rounded-[16px] text-[10px] md:text-[11px] uppercase tracking-widest border-2 border-slate-200 flex items-center justify-center gap-2">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-300" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
                            Selesai Hari Ini
                        </div>
                    </template>
                </div>
            </div>
        </div>
    </div>
</template>