<script setup>
defineProps({
    stores: Array,
    isLoading: Boolean
});

const emit = defineEmits(['suspend', 'activate']);

const formatDate = (dateStr) => {
    if (!dateStr) return 'Batas Kustom';
    const d = new Date(dateStr);
    return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' });
};
</script>

<template>
    <div>
        <div class="hidden lg:block bg-[#131B2E] border border-slate-800 rounded-[24px] overflow-hidden shadow-2xl">
            <div class="overflow-x-auto">
                <table class="w-full text-left border-collapse">
                    <thead>
                        <tr class="border-b border-slate-800 bg-[#1a243d]/50 text-[10px] font-black text-slate-400 uppercase tracking-widest">
                            <th class="p-5">Informasi Ruko Retail</th>
                            <th class="p-5">ID Sistem (UID)</th> <th class="p-5">Nama Pemilik</th>
                            <th class="p-5">Kategori & Tier</th>
                            <th class="p-5">Masa Berlaku</th>
                            <th class="p-5">Sisa Waktu</th>
                            <th class="p-5 text-center">Status & Aksi Komando</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-slate-800 text-xs text-slate-300">
                        <tr v-for="store in stores" :key="store.id" class="hover:bg-[#1a243d]/20 transition-all">
                            <td class="p-5">
                                <div class="font-black text-white text-sm">{{ store.nama_toko }}</div>
                                <div class="text-[10px] text-slate-500 font-bold uppercase mt-1 flex items-center gap-2">
                                    <span>{{ store.business_type || 'Retail' }}</span>
                                    <span>•</span>
                                    <span>{{ store.telepon || '-' }}</span>
                                </div>
                            </td>
                            
                            <td class="p-5 font-mono text-[11px] text-slate-400 select-all tracking-tight">
                                {{ store.public_id || store.PublicID || '-' }}
                            </td>

                            <td class="p-5">
                                <div class="font-bold text-slate-300">{{ store.owner_name }}</div>
                                <div class="text-[10px] text-slate-500 font-mono mt-0.5">{{ store.owner_email }}</div>
                            </td>
                            <td class="p-5">
                                <span class="px-2.5 py-1 bg-indigo-500/10 text-indigo-400 border border-indigo-500/20 rounded-lg font-black uppercase text-[10px] tracking-wider">
                                    {{ store.subscription_plan }}
                                </span>
                            </td>
                            <td class="p-5 font-mono text-slate-400">{{ formatDate(store.subscription_end) }}</td>
                            <td class="p-5 font-black" :class="store.sisa_hari <= 3 ? 'text-rose-400 animate-pulse' : 'text-slate-400'">
                                {{ store.sisa_hari }} Hari Lagi
                            </td>
                            <td class="p-5">
                                <div class="flex items-center justify-center gap-3">
                                    <span :class="store.subscription_status === 'active' ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20' : 'bg-rose-500/10 text-rose-400 border-rose-500/20'" class="px-2.5 py-1 border rounded-lg font-black text-[10px] uppercase tracking-wider">
                                        {{ store.subscription_status }}
                                    </span>
                                    <button v-if="store.subscription_status === 'active'" @click="emit('suspend', store)" class="px-3 py-1.5 bg-rose-500/10 hover:bg-rose-600 text-rose-400 hover:text-white rounded-xl font-black text-[10px] uppercase tracking-wider transition-all border border-rose-500/20 active:scale-95">
                                        Suspend
                                    </button>
                                    <button v-else @click="emit('activate', store)" class="px-3 py-1.5 bg-emerald-500/10 hover:bg-emerald-600 text-emerald-400 hover:text-white rounded-xl font-black text-[10px] uppercase tracking-wider transition-all border border-emerald-500/20 active:scale-95">
                                        Activate
                                    </button>
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 lg:hidden">
            <div v-for="store in stores" :key="store.id" class="bg-[#131B2E] border border-slate-800 rounded-2xl p-5 shadow-xl flex flex-col justify-between space-y-4 relative overflow-hidden">
                <div class="flex justify-between items-start border-b border-slate-800 pb-3">
                    <div>
                        <h3 class="text-sm font-black text-white tracking-tight">{{ store.nama_toko }}</h3>
                        <p class="text-[10px] text-slate-500 font-bold uppercase tracking-wider mt-0.5">{{ store.business_type || 'Retail' }}</p>
                    </div>
                    <span :class="store.subscription_status === 'active' ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20' : 'bg-rose-500/10 text-rose-400 border-rose-500/20'" class="px-2 py-0.5 border rounded-md text-[9px] font-black uppercase tracking-wider">
                        {{ store.subscription_status }}
                    </span>
                </div>

                <div class="grid grid-cols-2 gap-y-3 gap-x-2 text-[11px] text-slate-400">
                    <div>
                        <span class="text-[9px] text-slate-500 uppercase font-black tracking-wider block">Pemilik</span>
                        <span class="font-bold text-slate-200 block truncate">{{ store.owner_name }}</span>
                        <span class="text-[10px] text-slate-500 block truncate">{{ store.owner_email }}</span>
                    </div>
                    <div>
                        <span class="text-[9px] text-slate-500 uppercase font-black tracking-wider block">Paket Lisensi</span>
                        <span class="px-1.5 py-0.5 bg-indigo-500/10 text-indigo-400 border border-indigo-500/20 rounded font-black uppercase text-[9px] inline-block mt-0.5 tracking-wider">
                            {{ store.subscription_plan }}
                        </span>
                    </div>
                    <div>
                        <span class="text-[9px] text-slate-500 uppercase font-black tracking-wider block">Habis Kontrak</span>
                        <span class="font-bold font-mono text-slate-300">{{ formatDate(store.subscription_end) }}</span>
                    </div>
                    <div>
                        <span class="text-[9px] text-slate-500 uppercase font-black tracking-wider block">Sisa Waktu</span>
                        <span :class="store.sisa_hari <= 3 ? 'text-rose-400 font-black animate-pulse' : 'text-slate-300 font-bold'">
                            {{ store.sisa_hari }} Hari Lagi
                        </span>
                    </div>
                </div>

                <div class="pt-3 border-t border-slate-800/60 flex items-center justify-between gap-4">
                    <span class="text-[9px] font-mono text-slate-600 select-all truncate max-w-[120px]">
                        UID: {{ store.public_id || store.PublicID || '-' }}
                    </span>
                    
                    <button v-if="store.subscription_status === 'active'" @click="emit('suspend', store)" class="px-4 py-2 bg-rose-600 hover:bg-rose-700 text-white rounded-xl font-black text-[10px] uppercase tracking-widest transition-all active:scale-95 shadow-md">
                        Suspend
                    </button>
                    <button v-else @click="emit('activate', store)" class="px-4 py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl font-black text-[10px] uppercase tracking-widest transition-all active:scale-95 shadow-md">
                        Activate
                    </button>
                </div>
            </div>
        </div>

        <div v-if="stores.length === 0 && !isLoading" class="p-10 bg-[#131B2E] border border-slate-800 rounded-[24px] text-center shadow-xl mt-4">
            <p class="text-xs text-slate-500 font-bold uppercase tracking-wider">Data Tenant Tidak Terdeteksi di Radar Hub</p>
        </div>
    </div>
</template>