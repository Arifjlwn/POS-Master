<script setup>
defineProps({
    riwayat: Array,
    isLoading: Boolean,
    filterMode: String,
    tanggalDipilih: String,
    bulanDipilih: String,
    urutanTanggalTerbaru: Boolean,
    currentUser: Object
});

defineEmits(['update:filterMode', 'update:tanggalDipilih', 'update:bulanDipilih', 'toggle-sort', 'download-laporan', 'lihat-foto']);
</script>

<template>
    <div class="bg-white rounded-[32px] shadow-sm border border-slate-100 overflow-hidden mb-10">
        <div class="p-5 md:p-6 border-b border-slate-100 flex flex-col lg:flex-row lg:justify-between lg:items-center bg-slate-50/50 gap-4">
            <h3 class="font-black text-slate-800 text-lg flex items-center gap-2 uppercase tracking-tight">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-indigo-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/><polyline points="10 9 9 9 8 9"/></svg>
                Log Riwayat
            </h3>
            
            <div class="flex flex-wrap items-center gap-2 md:gap-3 w-full lg:w-auto">
                <select :value="filterMode" @change="$emit('update:filterMode', $event.target.value)" class="flex-1 lg:flex-none px-4 py-3 bg-white border border-slate-200 rounded-xl text-xs font-black text-slate-700 outline-none uppercase tracking-widest shadow-sm">
                    <option value="harian">Harian</option><option value="bulanan">Bulanan</option>
                </select>
                <input v-if="filterMode === 'harian'" type="date" :value="tanggalDipilih" @input="$emit('update:tanggalDipilih', $event.target.value)" class="flex-1 lg:flex-none px-4 py-3 bg-white border border-slate-200 rounded-xl text-xs font-black text-slate-700 outline-none uppercase shadow-sm cursor-pointer color-scheme-dark">
                <input v-else type="month" :value="bulanDipilih" @input="$emit('update:bulanDipilih', $event.target.value)" class="flex-1 lg:flex-none px-4 py-3 bg-white border border-slate-200 rounded-xl text-xs font-black text-slate-700 outline-none uppercase shadow-sm cursor-pointer color-scheme-dark">
                
                <button v-if="filterMode === 'bulanan'" @click="$emit('toggle-sort')" class="lg:hidden w-11 h-11 flex items-center justify-center bg-white border border-slate-200 rounded-xl text-slate-500 shadow-sm active:scale-95">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 transition-transform" :class="!urutanTanggalTerbaru ? 'rotate-180 text-indigo-600' : ''" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="m3 16 4 4 4-4"/><path d="M7 20V4"/><path d="m21 8-4-4-4 4"/><path d="M17 4v16"/></svg>
                </button>

                <button v-if="currentUser.role === 'owner'" @click="$emit('download-laporan')" class="w-full lg:w-auto mt-2 lg:mt-0 bg-slate-900 text-white px-5 py-3 rounded-xl font-black text-[10px] uppercase tracking-widest shadow-md flex items-center justify-center gap-2 hover:bg-indigo-600 transition-colors active:scale-95">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" x2="12" y1="3" y2="15"/></svg>
                    Ekspor Laporan
                </button>
            </div>
        </div>
        
        <div class="lg:hidden p-4 bg-slate-50/50">
            <div v-if="isLoading" class="py-12 text-center text-slate-400 font-black text-xs uppercase tracking-widest animate-pulse border-2 border-dashed border-slate-200 rounded-[24px]">Sinkronisasi Data...</div>
            <div v-else-if="riwayat.length === 0" class="py-12 text-center text-slate-400 font-black text-xs uppercase tracking-widest opacity-50 border-2 border-dashed border-slate-200 rounded-[24px]">Data Kosong</div>
            <div v-else class="flex flex-col gap-4">
                <div v-for="log in riwayat" :key="log.id" class="bg-white p-4 rounded-[20px] shadow-sm border border-slate-100 flex flex-col gap-3 relative overflow-hidden">
                    <div class="absolute left-0 top-0 bottom-0 w-1.5" 
                         :class="{
                            'bg-emerald-400': log.status === 'Hadir',
                            'bg-amber-400': log.status === 'Lupa Absen Pulang',
                            'bg-red-500': log.status === 'Mangkir',
                            'bg-slate-300': log.status === 'Libur (OFF)' || log.status === 'Belum Absen',
                            'bg-blue-400': log.status === 'Owner' || (!['Hadir','Lupa Absen Pulang','Mangkir','Libur (OFF)','Belum Absen'].includes(log.status))
                         }">
                    </div>

                    <div class="flex justify-between items-start pl-2">
                        <div>
                            <div class="font-black text-slate-800 uppercase text-sm tracking-tight">{{ log.User?.name }}</div>
                            <div class="text-[9px] text-slate-400 font-bold tracking-widest mt-0.5">{{ log.User?.nik || 'ID: '+log.user_id }}</div>
                        </div>
                        <div class="text-right">
                            <div v-if="filterMode === 'bulanan'" class="text-[10px] font-black text-slate-600 mb-1">{{ log.tanggal ? log.tanggal.substring(0, 10) : '-' }}</div>
                            <span class="bg-slate-50 border border-slate-100 px-2 py-0.5 rounded text-[8px] font-black uppercase tracking-widest text-indigo-600">{{ log.shift }}</span>
                        </div>
                    </div>

                    <div class="border-t border-dashed border-slate-100 my-1"></div>

                    <div class="grid grid-cols-2 gap-3 pl-2">
                        <div>
                            <p class="text-[8px] font-black text-slate-400 uppercase tracking-widest mb-1">Masuk</p>
                            <div class="flex items-center gap-2">
                                <div v-if="log.foto_masuk" @click="$emit('lihat-foto', log.foto_masuk, log.User?.name, 'Masuk', log.jam_masuk)" class="w-8 h-8 rounded-lg border border-slate-200 overflow-hidden cursor-zoom-in">
                                    <img :src="log.foto_masuk" class="w-full h-full object-cover">
                                </div>
                                <div v-else class="w-8 h-8 rounded-lg border border-slate-100 bg-slate-50 flex items-center justify-center">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3 text-slate-300" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M4 14.899A7 7 0 1 1 15.71 8h1.79a4.5 4.5 0 0 1 2.5 8.242"/><path d="M12 12v9"/><path d="m8 17 4-4 4 4"/></svg>
                                </div>
                                <span v-if="log.jam_masuk" class="font-black text-xs text-emerald-600">{{ log.jam_masuk }}</span>
                                <span v-else class="font-black text-xs text-slate-300">-</span>
                            </div>
                        </div>
                        <div>
                            <p class="text-[8px] font-black text-slate-400 uppercase tracking-widest mb-1">Pulang</p>
                            <div class="flex items-center gap-2">
                                <div v-if="log.foto_pulang" @click="$emit('lihat-foto', log.foto_pulang, log.User?.name, 'Pulang', log.jam_pulang)" class="w-8 h-8 rounded-lg border border-slate-200 overflow-hidden cursor-zoom-in">
                                    <img :src="log.foto_pulang" class="w-full h-full object-cover">
                                </div>
                                <div v-else class="w-8 h-8 rounded-lg border border-slate-100 bg-slate-50 flex items-center justify-center">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3 text-slate-300" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M4 14.899A7 7 0 1 1 15.71 8h1.79a4.5 4.5 0 0 1 2.5 8.242"/><path d="M12 12v9"/><path d="m8 17 4-4 4 4"/></svg>
                                </div>
                                <span v-if="log.jam_pulang" class="font-black text-xs text-amber-600">{{ log.jam_pulang }}</span>
                                <span v-else class="font-black text-xs text-slate-300">-</span>
                            </div>
                        </div>
                    </div>
                    
                    <div class="pl-2 mt-1">
                        <span v-if="log.status === 'Hadir'" class="inline-flex items-center gap-1 bg-emerald-50 text-emerald-700 font-black px-2.5 py-1 rounded text-[9px] uppercase tracking-widest">Hadir</span>
                        <span v-else-if="log.status === 'Lupa Absen Pulang'" class="inline-flex items-center gap-1 bg-amber-50 text-amber-700 font-black px-2.5 py-1 rounded text-[9px] uppercase tracking-widest">Lupa Pulang</span>
                        <span v-else-if="log.status === 'Mangkir'" class="inline-flex items-center gap-1 bg-red-50 text-red-700 font-black px-2.5 py-1 rounded text-[9px] uppercase tracking-widest">Mangkir</span>
                        <span v-else-if="log.status === 'Libur (OFF)'" class="inline-flex items-center gap-1 bg-slate-100 text-slate-500 font-black px-2.5 py-1 rounded text-[9px] uppercase tracking-widest">OFF</span>
                        <span v-else-if="log.status === 'Belum Absen'" class="inline-flex items-center gap-1 bg-blue-50 text-blue-600 font-black px-2.5 py-1 rounded text-[9px] uppercase tracking-widest">Belum Absen</span>
                        <span v-else class="inline-flex items-center gap-1 bg-slate-100 text-slate-600 font-black px-2.5 py-1 rounded text-[9px] uppercase tracking-widest">{{ log.status }}</span>
                    </div>
                </div>
            </div>
        </div>

        <div class="hidden lg:block overflow-x-auto custom-scrollbar">
            <table class="w-full text-left whitespace-nowrap">
                <thead class="bg-white border-b border-slate-100">
                    <tr class="text-[9px] font-black text-slate-400 uppercase tracking-[0.2em]">
                        <th class="px-6 py-5">Karyawan</th>
                        <th v-if="filterMode === 'bulanan'" @click="$emit('toggle-sort')" class="px-6 py-5 text-center cursor-pointer select-none hover:bg-slate-50 transition-colors group">
                            <div class="flex items-center justify-center gap-2">
                                Tanggal
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3 text-indigo-600 transition-transform" :class="!urutanTanggalTerbaru ? 'rotate-180' : ''" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="m6 9 6 6 6-6"/></svg>
                            </div>
                        </th>
                        <th class="px-6 py-5 text-center">Shift</th>
                        <th class="px-6 py-5 text-center">Foto Masuk</th>
                        <th class="px-6 py-5 text-center">Jam Masuk</th>
                        <th class="px-6 py-5 text-center">Foto Pulang</th>
                        <th class="px-6 py-5 text-center">Jam Pulang</th>
                        <th class="px-6 py-5 text-center">Status</th>
                    </tr>
                </thead>
                <tbody class="divide-y divide-slate-50">
                    <tr v-if="isLoading">
                        <td :colspan="filterMode === 'bulanan' ? 8 : 7" class="px-6 py-16 text-center text-slate-400 font-black text-xs uppercase tracking-widest animate-pulse">Sinkronisasi Data...</td>
                    </tr>
                    <tr v-else-if="riwayat.length === 0">
                        <td :colspan="filterMode === 'bulanan' ? 8 : 7" class="px-6 py-16 text-center text-slate-400 font-black text-xs uppercase tracking-widest opacity-50">Data Riwayat Kosong</td>
                    </tr>
                    <tr v-else v-for="log in riwayat" :key="log.id" class="hover:bg-indigo-50/30 transition-colors group">
                        <td class="px-6 py-5">
                            <div class="font-black text-slate-800 uppercase text-xs">{{ log.User?.name }}</div>
                            <div class="text-[9px] text-slate-400 font-bold tracking-widest mt-0.5">{{ log.User?.nik || 'ID: '+log.user_id }}</div>
                        </td>
                        
                        <td v-if="filterMode === 'bulanan'" class="px-6 py-5 text-center font-mono text-[11px] font-black text-slate-600">
                            {{ log.tanggal ? log.tanggal.substring(0, 10) : '-' }}
                        </td>

                        <td class="px-6 py-5 text-center font-black text-[9px] text-indigo-600 uppercase tracking-widest">
                            <span class="bg-slate-50 border border-slate-100 px-3 py-1.5 rounded-lg">{{ log.shift }}</span>
                        </td>

                        <td class="px-6 py-5 text-center">
                            <div v-if="log.foto_masuk" @click="$emit('lihat-foto', log.foto_masuk, log.User?.name, 'Masuk', log.jam_masuk)" class="w-12 h-12 rounded-[14px] mx-auto border-2 border-slate-200 shadow-sm overflow-hidden cursor-zoom-in hover:border-indigo-500 transition-colors">
                                <img :src="log.foto_masuk" class="w-full h-full object-cover">
                            </div>
                            <span v-else class="text-slate-300 font-black text-xs">-</span>
                        </td>

                        <td class="px-6 py-5 text-center">
                            <span v-if="log.jam_masuk" class="bg-emerald-50 text-emerald-700 border border-emerald-100 font-black px-3 py-1.5 rounded-lg text-[10px]">{{ log.jam_masuk }}</span>
                            <span v-else class="text-slate-300 font-black text-xs">-</span>
                        </td>

                        <td class="px-6 py-5 text-center">
                            <div v-if="log.foto_pulang" @click="$emit('lihat-foto', log.foto_pulang, log.User?.name, 'Pulang', log.jam_pulang)" class="w-12 h-12 rounded-[14px] mx-auto border-2 border-slate-200 shadow-sm overflow-hidden cursor-zoom-in hover:border-indigo-500 transition-colors">
                                <img :src="log.foto_pulang" class="w-full h-full object-cover">
                            </div>
                            <span v-else class="text-slate-300 font-black text-xs">-</span>
                        </td>

                        <td class="px-6 py-5 text-center">
                            <span v-if="log.jam_pulang" class="bg-amber-50 text-amber-700 border border-amber-100 font-black px-3 py-1.5 rounded-lg text-[10px]">{{ log.jam_pulang }}</span>
                            <span v-else class="text-slate-300 font-black text-xs">-</span>
                        </td>
                        
                        <td class="px-6 py-5 text-center">
                            <span v-if="log.status === 'Hadir'" class="inline-flex items-center gap-1.5 bg-emerald-50 text-emerald-700 border border-emerald-200 font-black px-3 py-1.5 rounded-full text-[9px] uppercase tracking-widest shadow-sm">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg> Hadir
                            </span>
                            <span v-else-if="log.status === 'Lupa Absen Pulang'" class="inline-flex items-center gap-1.5 bg-amber-50 text-amber-700 border border-amber-200 font-black px-3 py-1.5 rounded-full text-[9px] uppercase tracking-widest shadow-sm">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg> Lupa Pulang
                            </span>
                            <span v-else-if="log.status === 'Mangkir'" class="inline-flex items-center gap-1.5 bg-red-50 text-red-700 border border-red-200 font-black px-3 py-1.5 rounded-full text-[9px] uppercase tracking-widest shadow-sm">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polygon points="7.86 2 16.14 2 22 7.86 22 16.14 16.14 22 7.86 22 2 16.14 2 7.86 7.86 2"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg> Mangkir
                            </span>
                            <span v-else-if="log.status === 'Libur (OFF)'" class="inline-flex items-center gap-1.5 bg-slate-50 text-slate-500 border border-slate-200 font-black px-3 py-1.5 rounded-full text-[9px] uppercase tracking-widest shadow-sm">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/><path d="m9 16 2 2 4-4"/></svg> OFF
                            </span>
                            <span v-else-if="log.status === 'Belum Absen'" class="inline-flex items-center gap-1.5 bg-blue-50 text-blue-700 border border-blue-200 font-black px-3 py-1.5 rounded-full text-[9px] uppercase tracking-widest shadow-sm">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg> Belum Absen
                            </span>
                            <span v-else class="bg-slate-100 text-slate-600 font-black px-3 py-1.5 rounded-full text-[9px] uppercase tracking-widest shadow-sm border border-slate-200">
                                {{ log.status }}
                            </span>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 5px; height: 5px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #94a3b8; }

input[type="date"]::-webkit-calendar-picker-indicator,
input[type="month"]::-webkit-calendar-picker-indicator {
    cursor: pointer;
    opacity: 0.6;
    transition: 0.2s;
}
input[type="date"]::-webkit-calendar-picker-indicator:hover,
input[type="month"]::-webkit-calendar-picker-indicator:hover { opacity: 1; }

.color-scheme-dark::-webkit-calendar-picker-indicator {
    filter: invert(1);
    opacity: 0.8;
}
</style>