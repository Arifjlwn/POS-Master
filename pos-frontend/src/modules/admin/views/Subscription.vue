<script setup>
import { ref, onMounted } from 'vue';
import Swal from 'sweetalert2';
import api from '../../../api.js';

const isLoading = ref(false);
const planCounts = ref({ trial: 0, basic: 0, pro: 0, premium: 0 });
const storeBillingList = ref([]);

const fetchBillingData = async () => {
	isLoading.value = true;
	try {
		const res = await api.get('/admin/subscription-overview');
		if (res.data.status === 'sukses') {
			planCounts.value = res.data.counts;
			storeBillingList.value = res.data.stores;
		}
	} catch (err) {
		Swal.fire({
			icon: 'error',
			title: 'Gagal Sinkronisasi Billing',
			text: 'Terputus dari server internal.',
			confirmButtonColor: '#ef4444',
			customClass: { popup: 'rounded-[24px]' },
		});
	} finally {
		isLoading.value = false;
	}
};

const formatDate = (dateStr) => {
	if (!dateStr) return '-';
	const d = new Date(dateStr);
	return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' });
};

onMounted(() => fetchBillingData());
</script>

<template>
	<div class="p-4 sm:p-6 bg-[#0B0F19] min-h-screen text-white font-sans overflow-x-hidden">
		<div class="mb-6">
			<div class="inline-flex items-center gap-2 px-3 py-1 bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 rounded-xl text-[10px] font-black uppercase tracking-widest mb-2">Financing & Analytics</div>
			<h1 class="text-xl sm:text-2xl font-black text-white tracking-tight">SUBSCRIPTION CONTROL HUB</h1>
			<p class="text-slate-500 font-bold text-[9px] sm:text-[10px] uppercase tracking-widest mt-0.5">Manajemen Akumulasi Paket Lisensi Tenant Multi-Ruko</p>
		</div>

		<div class="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
			<div class="bg-[#131B2E] border border-slate-800 p-4 rounded-2xl shadow-xl flex flex-col justify-between">
				<span class="text-[10px] text-slate-500 font-black uppercase tracking-wider">Masa Trial</span>
				<h2 class="text-2xl font-black text-indigo-400 mt-2">
					{{ planCounts.trial }}
					<span class="text-xs text-slate-500 font-bold">Ruko</span>
				</h2>
			</div>
			<div class="bg-[#131B2E] border border-slate-800 p-4 rounded-2xl shadow-xl flex flex-col justify-between">
				<span class="text-[10px] text-slate-500 font-black uppercase tracking-wider">Paket Basic</span>
				<h2 class="text-2xl font-black text-amber-400 mt-2">
					{{ planCounts.basic }}
					<span class="text-xs text-slate-500 font-bold">Ruko</span>
				</h2>
			</div>
			<div class="bg-[#131B2E] border border-slate-800 p-4 rounded-2xl shadow-xl flex flex-col justify-between">
				<span class="text-[10px] text-slate-500 font-black uppercase tracking-wider">Paket Pro</span>
				<h2 class="text-2xl font-black text-emerald-400 mt-2">
					{{ planCounts.pro }}
					<span class="text-xs text-slate-500 font-bold">Ruko</span>
				</h2>
			</div>
			<div class="bg-[#131B2E] border border-slate-800 p-4 rounded-2xl shadow-xl flex flex-col justify-between">
				<span class="text-[10px] text-slate-500 font-black uppercase tracking-wider">Paket Premium</span>
				<h2 class="text-2xl font-black text-purple-400 mt-2">
					{{ planCounts.premium }}
					<span class="text-xs text-slate-500 font-bold">Ruko</span>
				</h2>
			</div>
		</div>

		<div class="hidden md:block bg-[#131B2E] border border-slate-800 rounded-[24px] overflow-hidden shadow-2xl">
			<div class="overflow-x-auto">
				<table class="w-full text-left border-collapse">
					<thead>
						<tr class="border-b border-slate-800 bg-[#1a243d]/50 text-[10px] font-black text-slate-400 uppercase tracking-widest">
							<th class="p-5">Nama Mitra Ruko</th>
							<th class="p-5">Nama Owner</th>
							<th class="p-5">Tier Paket</th>
							<th class="p-5">Masa Berlaku</th>
							<th class="p-5">Sisa Hari</th>
							<th class="p-5">Status Layanan</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-slate-800 text-xs text-slate-300">
						<tr v-for="store in storeBillingList" :key="store.id" class="hover:bg-[#1a243d]/20 transition-all">
							<td class="p-5 font-black text-white">{{ store.nama_toko }}</td>
							<td class="p-5 font-bold text-slate-400">{{ store.owner_name }}</td>
							<td class="p-5">
								<span class="px-2 py-0.5 bg-indigo-500/10 text-indigo-400 border border-indigo-500/20 rounded font-black uppercase text-[10px]">
									{{ store.subscription_plan }}
								</span>
							</td>
							<td class="p-5 font-mono">{{ formatDate(store.subscription_end) }}</td>
							<td class="p-5 font-bold" :class="store.sisa_hari <= 3 ? 'text-red-400' : 'text-slate-400'">{{ store.sisa_hari }} Hari Lagi</td>
							<td class="p-5">
								<span :class="store.subscription_status === 'active' ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20' : 'bg-red-500/10 text-red-400 border-red-500/20'" class="px-2 py-0.5 border rounded font-black text-[10px] uppercase">
									{{ store.subscription_status }}
								</span>
							</td>
						</tr>
					</tbody>
				</table>
			</div>
		</div>

		<div class="block md:hidden space-y-4">
			<div v-for="store in storeBillingList" :key="store.id" class="bg-[#131B2E] border border-slate-800 rounded-2xl p-4 shadow-xl space-y-3">
				<div class="flex justify-between items-center border-b border-slate-800 pb-2">
					<h4 class="text-xs font-black text-white">{{ store.nama_toko }}</h4>
					<span :class="store.subscription_status === 'active' ? 'text-emerald-400' : 'text-red-400'" class="text-[9px] font-black uppercase">● {{ store.subscription_status }}</span>
				</div>
				<div class="text-[11px] space-y-1.5 text-slate-400">
					<p>
						Owner:
						<span class="text-slate-200 font-bold">{{ store.owner_name }}</span>
					</p>
					<p>
						Paket:
						<span class="text-indigo-400 font-black uppercase text-[10px]">{{ store.subscription_plan }}</span>
					</p>
					<p>
						Habis Tanggal:
						<span class="text-slate-300 font-mono">{{ formatDate(store.subscription_end) }}</span>
					</p>
					<p>
						Sisa Waktu:
						<span :class="store.sisa_hari <= 3 ? 'text-red-400 font-black' : 'text-slate-300 font-bold'">{{ store.sisa_hari }} Hari Lagi</span>
					</p>
				</div>
			</div>
		</div>
	</div>
</template>
