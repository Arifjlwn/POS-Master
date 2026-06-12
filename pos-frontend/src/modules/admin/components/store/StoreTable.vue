<script setup>
const emit = defineEmits(['suspend', 'activate']);

defineProps({
	stores: Array,
	isLoading: Boolean,
});

const formatDate = (dateStr) => {
	if (!dateStr) return 'Selamanya / Batas Kustom';
	const d = new Date(dateStr);
	return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' });
};
</script>

<template>
	<div>
		<div class="hidden md:block bg-[#131B2E] border border-slate-800 rounded-[24px] overflow-hidden shadow-2xl">
			<div class="overflow-x-auto">
				<table class="w-full text-left border-collapse">
					<thead>
						<tr class="border-b border-slate-800 bg-[#1a243d]/50 text-[10px] font-black text-slate-400 uppercase tracking-widest">
							<th class="p-5">Informasi Toko</th>
							<th class="p-5">Pemilik (Owner ID)</th>
							<th class="p-5">Paket Langganan</th>
							<th class="p-5">Masa Aktif</th>
							<th class="p-5">Status</th>
							<th class="p-5 text-center">Aksi Operasional</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-slate-800 text-sm">
						<tr v-for="item in stores" :key="item.id" class="hover:bg-[#1a243d]/20 transition-all">
							<td class="p-5">
								<div class="font-black text-white flex items-center gap-1.5">
									{{ item.nama_toko }}
								</div>
								<div class="text-xs text-slate-400 mt-1 flex items-center gap-2">
									<span>{{ item.business_type }}</span>
									<span class="text-slate-600">•</span>
									<span class="flex items-center gap-1">
										<svg class="w-3 h-3 text-slate-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
											<path stroke-linecap="round" stroke-linejoin="round" d="M2.25 6.622a2.25 2.25 0 0 1 2.25-2.25h1.378c.513 0 .96.344 1.087.835l.383 1.437c.127.476-.053.983-.453 1.28l-.753.559a11.356 11.356 0 0 0 5.445 5.446l.559-.753c.297-.4.805-.58 1.28-.453l1.437.383c.49.128.835.574.835 1.087V17.25a2.25 2.25 0 0 1-2.25 2.25h-.377a10.518 10.518 0 0 1-9.516-9.516V6.622Z" />
										</svg>
										{{ item.telepon || '-' }}
									</span>
								</div>
							</td>

							<td class="p-5">
								<div class="font-bold text-slate-200">{{ item.owner_name }}</div>
								<div class="text-xs text-slate-500 mt-0.5">{{ item.owner_email }}</div>
								<div class="text-[10px] text-slate-600 font-mono mt-1 select-all">UID: {{ item.public_id }}</div>
							</td>

							<td class="p-5">
								<span class="px-2.5 py-1 bg-indigo-500/10 border border-indigo-500/20 Triang-indigo-400 rounded-lg text-xs font-bold uppercase tracking-wide">
									{{ item.subscription_plan }}
								</span>
							</td>
							<td class="p-5 font-medium text-slate-300">
								{{ formatDate(item.subscription_end) }}
							</td>
							<td class="p-5">
								<span
									:class="{
										'bg-emerald-500/10 border-emerald-500/20 text-emerald-400': item.subscription_status === 'active',
										'bg-red-500/10 border-red-500/20 text-red-400': item.subscription_status === 'suspended',
										'bg-amber-500/10 border-amber-500/20 text-amber-400': item.subscription_status === 'pending',
									}"
									class="px-2.5 py-1 border rounded-lg text-xs font-black uppercase tracking-wider">
									{{ item.subscription_status }}
								</span>
							</td>

							<td class="p-5">
								<div class="flex justify-center gap-2">
									<button v-if="item.subscription_status === 'active'" @click="emit('suspend', item)" class="px-3 py-2 bg-red-600 hover:bg-red-700 active:scale-95 text-white font-black text-xs uppercase tracking-widest rounded-xl transition-all flex items-center gap-1.5">Suspend</button>
									<button v-else @click="emit('activate', item)" class="px-3 py-2 bg-emerald-600 hover:bg-emerald-700 active:scale-95 text-white font-black text-xs uppercase tracking-widest rounded-xl transition-all flex items-center gap-1.5">Activate</button>
								</div>
							</td>
						</tr>
					</tbody>
				</table>
			</div>
		</div>

		<div class="block md:hidden space-y-4">
			<div v-for="item in stores" :key="item.id" class="bg-[#131B2E] border border-slate-800 rounded-2xl p-4 shadow-xl space-y-4 relative overflow-hidden">
				<div class="flex justify-between items-start gap-2 border-b border-slate-800 pb-3">
					<div>
						<h4 class="text-sm font-black text-white tracking-tight">{{ item.nama_toko }}</h4>
						<span class="text-[10px] px-2 py-0.5 bg-slate-800 rounded text-slate-400 font-bold uppercase tracking-wider mt-1 inline-block">
							{{ item.business_type }}
						</span>
					</div>
					<span
						:class="{
							'bg-emerald-500/10 border-emerald-500/20 text-emerald-400': item.subscription_status === 'active',
							'bg-red-500/10 border-red-500/20 text-red-400': item.subscription_status === 'suspended',
							'bg-amber-500/10 border-amber-500/20 text-amber-400': item.subscription_status === 'pending',
						}"
						class="px-2 py-0.5 border rounded-lg text-[9px] font-black uppercase tracking-wider">
						{{ item.subscription_status }}
					</span>
				</div>

				<div class="grid grid-cols-2 gap-2 text-xs">
					<div>
						<span class="text-[10px] text-slate-500 uppercase font-black tracking-wider block">Pemilik</span>
						<span class="font-bold text-slate-300 block truncate">{{ item.owner_name }}</span>
						<span class="text-[10px] text-slate-500 block truncate">{{ item.owner_email }}</span>
					</div>
					<div>
						<span class="text-[10px] text-slate-500 uppercase font-black tracking-wider block">Masa Aktif</span>
						<span class="font-bold text-indigo-400 block">{{ formatDate(item.subscription_end) }}</span>
						<span class="text-[9px] px-1.5 py-0.5 bg-indigo-500/10 text-indigo-300 rounded font-black uppercase tracking-wide inline-block mt-0.5">Plan: {{ item.subscription_plan }}</span>
					</div>
				</div>

				<div class="flex flex-col sm:flex-row justify-between items-stretch sm:items-center gap-3 pt-2 border-t border-slate-800/60 text-[10px]">
					<div class="font-mono text-slate-500 truncate select-all">UID: {{ item.public_id }}</div>

					<button v-if="item.subscription_status === 'active'" @click="emit('suspend', item)" class="w-full sm:w-auto px-4 py-2.5 bg-red-600 hover:bg-red-700 active:scale-95 text-white font-black text-center text-xs uppercase tracking-widest rounded-xl transition-all shadow-md">Suspend Toko</button>
					<button v-else @click="emit('activate', item)" class="w-full sm:w-auto px-4 py-2.5 bg-emerald-600 hover:bg-emerald-700 active:scale-95 text-white font-black text-center text-xs uppercase tracking-widest rounded-xl transition-all shadow-md">Aktivasi Toko</button>
				</div>
			</div>
		</div>

		<div v-if="stores.length === 0 && !isLoading" class="p-10 bg-[#131B2E] border border-slate-800 rounded-[24px] text-center shadow-xl">
			<p class="text-xs text-slate-500 font-bold uppercase tracking-wider">Tidak ada data ruko retail yang terdeteksi di server pusat bray.</p>
		</div>
	</div>
</template>
