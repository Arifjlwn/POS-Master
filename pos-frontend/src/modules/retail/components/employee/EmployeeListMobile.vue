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
    <div class="lg:hidden flex flex-col gap-4">
        <div v-for="user in filteredKaryawan" :key="user.id" class="bg-white p-5 rounded-[24px] shadow-sm border border-slate-100 flex flex-col gap-4 relative overflow-hidden">
            
            <div class="flex items-center gap-4">
                <img v-if="user.foto_url" :src="(user.foto_url.startsWith('http://') || user.foto_url.startsWith('https://')) ? user.foto_url : API_BASE_URL + user.foto_url" class="w-14 h-14 rounded-2xl object-cover border-2 border-slate-100 shadow-sm shrink-0">
                <div v-else class="w-14 h-14 rounded-2xl bg-slate-100 flex items-center justify-center text-slate-400 border-2 border-slate-200 shrink-0">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
                </div>
                
                <div class="flex-1 min-w-0">
                    <h3 class="font-black text-base text-slate-800 uppercase truncate">{{ user.name }}</h3>
                    <div class="flex items-center gap-1.5 mt-1">
                        <span v-if="user.role === 'owner'" class="inline-flex items-center gap-1 bg-slate-900 text-white font-black px-2 py-0.5 rounded text-[8px] uppercase tracking-widest">Owner</span>
                        <span v-else-if="user.role === 'manager'" class="inline-flex items-center gap-1 bg-purple-50 text-purple-600 border border-purple-100 font-black px-2 py-0.5 rounded text-[8px] uppercase tracking-widest">Manager</span>
                        <span v-else-if="user.role === 'supervisor'" class="inline-flex items-center gap-1 bg-blue-50 text-blue-600 border border-blue-100 font-black px-2 py-0.5 rounded text-[8px] uppercase tracking-widest">Supervisor</span>
                        <span v-else class="inline-flex items-center gap-1 bg-emerald-50 text-emerald-600 border border-emerald-100 font-black px-2 py-0.5 rounded text-[8px] uppercase tracking-widest">Staff</span>
                    </div>
                </div>
            </div>

            <div class="grid grid-cols-2 gap-3 pt-3 border-t border-dashed border-slate-100">
                <div v-if="user.role !== 'owner'" class="col-span-2 bg-emerald-50 rounded-xl p-3 flex items-center gap-3 border border-emerald-100">
                    <div class="p-1.5 bg-emerald-100 rounded-lg text-emerald-600"><svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" /></svg></div>
                    <div>
                        <p class="text-[8px] font-black text-emerald-600 uppercase tracking-widest mb-0.5">Akses Login (WA)</p>
                        <p class="font-black text-sm text-emerald-700 tracking-widest">{{ formatNoHP(user.no_hp) }}</p>
                    </div>
                </div>
                
                <div v-else class="col-span-2 bg-slate-50 rounded-xl p-3 flex items-center gap-3 border border-slate-200">
                    <div class="p-1.5 bg-slate-200 rounded-lg text-slate-500"><svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" /></svg></div>
                    <div>
                        <p class="text-[8px] font-black text-slate-400 uppercase tracking-widest mb-0.5">Akses Login (Email)</p>
                        <p class="font-black text-xs text-slate-600 tracking-widest">{{ user.email || 'Email Owner Terlindungi' }}</p>
                    </div>
                </div>

                <div class="bg-slate-50 p-3 rounded-xl border border-slate-100">
                    <p class="text-[8px] font-black text-slate-400 uppercase tracking-widest mb-1">Tgl Bergabung</p>
                    <p class="font-bold text-xs text-slate-700">{{ formatDate(user.created_at) }}</p>
                </div>
                <div class="bg-slate-50 p-3 rounded-xl border border-slate-100">
                    <p class="text-[8px] font-black text-slate-400 uppercase tracking-widest mb-1">Masa Kerja</p>
                    <p class="font-black text-xs text-indigo-600">{{ hitungMasaKerja(user.created_at) }}</p>
                </div>

                <div class="col-span-2 pt-1">
                    <p class="text-[8px] font-black text-slate-400 uppercase tracking-widest mb-1">Data Personal (NIK)</p>
                    <p class="font-bold text-xs text-slate-600"><span class="text-indigo-600 font-black">{{ user.nik || '-' }}</span> • {{ user.tempat_lahir || '-' }}, {{ user.tanggal_lahir || '-' }}</p>
                </div>
            </div>

            <div class="flex gap-2 mt-2 pt-3 border-t border-dashed border-slate-100" v-if="user.role !== 'owner'">
                <button @click="emit('edit', user)" class="flex-1 bg-slate-100 text-slate-500 hover:bg-blue-600 hover:text-white py-3 rounded-xl transition-colors font-black text-[10px] uppercase tracking-widest flex items-center justify-center gap-2">Edit</button>
                <button @click="emit('delete', user.id)" class="flex-1 bg-red-50 text-red-500 hover:bg-red-500 hover:text-white py-3 rounded-xl transition-colors font-black text-[10px] uppercase tracking-widest flex items-center justify-center gap-2">Pecat</button>
            </div>
        </div>
    </div>
</template>