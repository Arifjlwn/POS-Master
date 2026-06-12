<script setup>
defineProps({
    topTenants: Array,
    formatRupiah: Function
});
</script>

<template>
    <div class="bg-[#131B2E] border border-slate-800 rounded-[24px] overflow-hidden shadow-2xl">
        <div class="p-5 border-b border-slate-800 bg-[#162035]/30">
            <h3 class="text-xs font-black uppercase tracking-widest text-indigo-400">Top Performing Tenant (Whale Categories)</h3>
            <p class="text-[10px] text-slate-500 font-bold uppercase tracking-wider mt-0.5">Peringkat 5 Besar Gerai Cabang Paling Gacor Berdasarkan Akumulasi Nilai GMV Kasir</p>
        </div>
        
        <div class="hidden lg:block">
            <div class="overflow-x-auto">
                <table class="w-full text-left border-collapse">
                    <thead>
                        <tr class="border-b border-slate-800 bg-[#1a243d]/20 text-[10px] font-black text-slate-400 uppercase tracking-widest">
                            <th class="p-4 pl-6">Peringkat</th>
                            <th class="p-4">Informasi Mitra Gerai</th>
                            <th class="p-4">ID Sistem (UID)</th>
                            <th class="p-4">Nama Pemilik Toko</th>
                            <th class="p-4 text-right pr-6">Total Akumulasi GMV Kasir</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-slate-800 text-xs text-slate-300">
                        <tr v-for="(tenant, idx) in topTenants" :key="tenant.store_id" class="hover:bg-[#1a243d]/20 transition-all">
                            <td class="p-4 pl-6 font-black font-mono text-slate-500">
                                <span :class="{
                                    'text-amber-400 bg-amber-500/10 border-amber-500/20 px-2 py-0.5 border rounded': idx === 0,
                                    'text-slate-300 bg-slate-500/10 border-slate-500/20 px-2 py-0.5 border rounded': idx === 1,
                                    'text-amber-600 bg-amber-700/10 border-amber-700/20 px-2 py-0.5 border rounded': idx === 2
                                }">
                                    0{{ idx + 1 }}
                                </span>
                            </td>
                            <td class="p-4 font-black text-white text-sm">{{ tenant.nama_toko }}</td>
                            <td class="p-4 font-mono text-slate-500 select-all">{{ tenant.public_id || '-' }}</td>
                            <td class="p-4 font-bold text-slate-400">{{ tenant.owner_name }}</td>
                            <td class="p-4 text-right pr-6 font-mono font-black text-emerald-400">
                                {{ formatRupiah(tenant.total_gmv) }}
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 p-4 lg:hidden">
            <div v-for="(tenant, idx) in topTenants" :key="tenant.store_id" class="bg-[#182237] border border-slate-800 rounded-2xl p-4 flex flex-col justify-between space-y-3">
                <div class="flex items-center justify-between border-b border-slate-800/60 pb-2">
                    <div class="flex items-center gap-2">
                        <span class="text-[10px] font-black font-mono bg-slate-800 text-slate-400 px-2 py-0.5 rounded border border-slate-700">
                            RANK 0{{ idx + 1 }}
                        </span>
                        <h4 class="text-sm font-black text-white truncate max-w-[150px]">{{ tenant.nama_toko }}</h4>
                    </div>
                    <span class="font-mono font-black text-xs text-emerald-400">
                        {{ formatRupiah(tenant.total_gmv) }}
                    </span>
                </div>
                <div class="text-[11px] text-slate-400 space-y-1">
                    <p class="truncate">Owner: <span class="font-bold text-slate-200">{{ tenant.owner_name }}</span></p>
                    <p class="text-[10px] font-mono text-slate-500 select-all truncate">UID: {{ tenant.public_id || '-' }}</p>
                </div>
            </div>
        </div>

        <div v-if="!topTenants || topTenants.length === 0" class="p-8 text-center text-slate-600 font-bold uppercase tracking-widest text-xs">
            Radar Belum Mendeteksi Volume Transaksi Kasir
        </div>
    </div>
</template>