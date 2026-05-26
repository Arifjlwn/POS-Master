<script setup>
import { defineProps, defineEmits } from 'vue';

defineProps({
    filteredKaryawan: Array,
    API_BASE_URL: String,
    formatDate: Function,
    hitungMasaKerja: Function
});

const emit = defineEmits(['edit', 'delete']);
</script>

<template>
    <div class="lg:hidden flex flex-col gap-4">
        <div v-for="user in filteredKaryawan" :key="user.id" class="bg-white p-5 rounded-[24px] shadow-sm border border-slate-100 flex flex-col gap-4 relative overflow-hidden">
            <div class="flex items-center gap-4">
                <img v-if="user.foto_url" :src="API_BASE_URL + user.foto_url" class="w-14 h-14 rounded-2xl object-cover border-2 border-slate-100 shadow-sm shrink-0">
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
                        <span class="text-[9px] font-black text-blue-600 tracking-widest">{{ user.nik }}</span>
                    </div>
                </div>
            </div>

            <div class="grid grid-cols-2 gap-3 pt-3 border-t border-dashed border-slate-100">
                <div>
                    <p class="text-[8px] font-black text-slate-400 uppercase tracking-widest mb-1">Tgl Bergabung</p>
                    <p class="font-bold text-xs text-slate-700">{{ formatDate(user.created_at) }}</p>
                </div>
                <div>
                    <p class="text-[8px] font-black text-slate-400 uppercase tracking-widest mb-1">Masa Kerja</p>
                    <p class="font-black text-xs text-indigo-600">{{ hitungMasaKerja(user.created_at) }}</p>
                </div>
                <div class="col-span-2">
                    <p class="text-[8px] font-black text-slate-400 uppercase tracking-widest mb-1">Kontak / Ttl</p>
                    <p class="font-bold text-xs text-slate-600">{{ user.no_hp || '-' }} • {{ user.tempat_lahir || '-' }}, {{ user.tanggal_lahir || '-' }}</p>
                </div>
            </div>

            <div class="flex gap-2 mt-2" v-if="user.role !== 'owner'">
                <button @click="emit('edit', user)" class="flex-1 bg-slate-100 text-slate-500 hover:bg-blue-600 hover:text-white py-3 rounded-xl transition-colors font-black text-[10px] uppercase tracking-widest flex items-center justify-center gap-2">Edit</button>
                <button @click="emit('delete', user.id)" class="flex-1 bg-red-50 text-red-500 hover:bg-red-500 hover:text-white py-3 rounded-xl transition-colors font-black text-[10px] uppercase tracking-widest flex items-center justify-center gap-2">Pecat</button>
            </div>
        </div>
    </div>
</template>