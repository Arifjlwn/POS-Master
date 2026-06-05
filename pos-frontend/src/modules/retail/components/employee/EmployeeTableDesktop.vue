<script setup>
defineProps({
    filteredKaryawan: Array,
    API_BASE_URL: String,
    formatDate: Function,
    hitungMasaKerja: Function
});

const emit = defineEmits(['edit', 'delete']);

const formatNoHP = (phone) => {
    if (!phone) return 'Belum Diatur';
    let str = String(phone);
    if (str.startsWith('62')) return '0' + str.slice(2);
    if (str.startsWith('+62')) return '0' + str.slice(3);
    return str;
};
</script>

<template>
    <div class="hidden lg:block bg-white rounded-[32px] shadow-sm border border-slate-100 overflow-hidden">
        <div class="overflow-x-auto custom-scrollbar">
            <table class="w-full text-left whitespace-nowrap">
                <thead class="bg-slate-50/80 border-b border-slate-100">
                    <tr>
                        <th class="px-6 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">Employee Profile</th>
                        <th class="px-6 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">Akses Login (WA)</th>
                        <th class="px-6 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">Tenure & Date Joined</th>
                        <th class="px-6 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">Personal Data (NIK)</th>
                        <th class="px-6 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest text-right">Actions</th>
                    </tr>
                </thead>
                <tbody class="divide-y divide-slate-50">
                    <tr v-for="user in filteredKaryawan" :key="user.id" class="hover:bg-slate-50/50 transition-colors group">
                        <td class="px-6 py-5">
                            <div class="flex items-center gap-4">
                                <img v-if="user.foto_url" :src="(user.foto_url.startsWith('http://') || user.foto_url.startsWith('https://')) ? user.foto_url : API_BASE_URL + user.foto_url" class="w-14 h-14 rounded-[14px] object-cover border-2 border-slate-100 shadow-sm">
                                <div v-else class="w-14 h-14 rounded-[14px] bg-slate-100 flex items-center justify-center text-slate-400 border-2 border-slate-200">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
                                </div>
                                <div>
                                    <div class="text-slate-800 font-black text-sm uppercase tracking-tight">{{ user.name }}</div>
                                    <div class="mt-1.5">
                                        <span v-if="user.role === 'owner'" class="inline-flex items-center gap-1 bg-slate-900 text-white font-black px-2.5 py-0.5 rounded text-[9px] uppercase tracking-widest shadow-sm">Owner</span>
                                        <span v-else-if="user.role === 'manager'" class="inline-flex items-center gap-1 bg-purple-50 text-purple-600 border border-purple-100 font-black px-2.5 py-0.5 rounded text-[9px] uppercase tracking-widest">Manager</span>
                                        <span v-else-if="user.role === 'supervisor'" class="inline-flex items-center gap-1 bg-blue-50 text-blue-600 border border-blue-100 font-black px-2.5 py-0.5 rounded text-[9px] uppercase tracking-widest">Supervisor</span>
                                        <span v-else class="inline-flex items-center gap-1 bg-emerald-50 text-emerald-600 border border-emerald-100 font-black px-2.5 py-0.5 rounded text-[9px] uppercase tracking-widest">Staff</span>
                                    </div>
                                </div>
                            </div>
                        </td>
                        
                        <td class="px-6 py-5">
                            <div v-if="user.role !== 'owner'" class="inline-flex items-center gap-2 bg-emerald-50 px-3 py-2 rounded-xl border border-emerald-200">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-emerald-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" /></svg>
                                <span class="text-emerald-700 font-black tracking-widest text-xs">{{ formatNoHP(user.no_hp) }}</span>
                            </div>
                            <div v-else class="inline-flex items-center gap-2 bg-slate-50 px-3 py-2 rounded-xl border border-slate-200">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" /></svg>
                                <span class="text-slate-500 font-black tracking-widest text-[10px] uppercase">Login Via Email</span>
                            </div>
                        </td>

                        <td class="px-6 py-5">
                            <div class="font-black text-indigo-600 text-sm">{{ hitungMasaKerja(user.created_at) }}</div>
                            <div class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-1">Sejak: {{ formatDate(user.created_at) }}</div>
                        </td>
                        
                        <td class="px-6 py-5">
                            <div class="text-slate-600 font-black text-xs tracking-widest">NIK: {{ user.nik || '-' }}</div>
                            <div class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-1">{{ user.tempat_lahir || '-' }}, {{ user.tanggal_lahir || '-' }}</div>
                        </td>
                        
                        <td class="px-6 py-5 text-right">
                            <div class="flex items-center justify-end gap-2">
                                <template v-if="user.role !== 'owner'">
                                    <button @click="emit('edit', user)" class="bg-slate-50 border border-slate-200 text-slate-500 hover:bg-blue-600 hover:text-white p-2.5 rounded-xl transition-colors shadow-sm" title="Edit">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M17 3a2.828 2.828 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5L17 3z"/></svg>
                                    </button>
                                    <button @click="emit('delete', user.public_id || user.id)" class="bg-red-50 text-red-500 hover:bg-red-500 hover:text-white p-2.5 rounded-xl transition-colors shadow-sm" title="Pecat">
    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/></svg>
</button>
                                </template>
                                <span v-else class="text-[9px] font-black text-slate-300 uppercase tracking-widest bg-slate-50 border border-slate-100 px-3 py-2 rounded-xl shadow-sm">Protected</span>
                            </div>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
</template>