<script setup>
import Swal from 'sweetalert2';
import { onMounted, ref, watch } from 'vue';
import api from '../../../api.js'; // Sesuaikan relative path ke api.js lu

const logs = ref([]);
const isLoading = ref(false);
const searchQuery = ref('');

const fetchAuditLogs = async () => {
	isLoading.value = true;
	try {
		const res = await api.get('/admin/audit-logs', {
			params: { search: searchQuery.value },
		});
		if (res.data.status === 'sukses') {
			logs.value = res.data.data;
		}
	} catch (err) {
		Swal.fire({
			icon: 'error',
			title: 'Gagal Memuat Log',
			text: err.response?.data?.error || 'Koneksi ke core server terputus.',
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
	return (
		d.toLocaleDateString('id-ID', {
			day: 'numeric',
			month: 'short',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit',
			second: '2-digit',
		}) + ' WIB'
	);
};

let searchTimeout;
watch(searchQuery, () => {
	clearTimeout(searchTimeout);
	searchTimeout = setTimeout(() => {
		fetchAuditLogs();
	}, 400); // 400ms delay
});

onMounted(() => {
	fetchAuditLogs();
});
</script>

<template>
	<div class="p-4 sm:p-6 bg-[#0B0F19] min-h-screen text-white font-sans overflow-x-hidden">
		<div class="flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
			<div class="w-full md:w-auto">
				<div class="inline-flex items-center gap-2 px-3 py-1 bg-indigo-500/10 border border-indigo-500/20 text-indigo-400 rounded-xl text-[10px] font-black uppercase tracking-widest mb-2">Security & Compliance</div>
				<h1 class="text-xl sm:text-2xl font-black text-white tracking-tight break-words">MASTER AUDIT LOGS</h1>
				<p class="text-slate-500 font-bold text-[9px] sm:text-[10px] uppercase tracking-widest mt-0.5 leading-relaxed">Rekam Jejak Digital Seluruh Aktivitas Ekosistem ARZURA POS</p>
			</div>

			<button @click="fetchAuditLogs" :disabled="isLoading" class="w-full md:w-auto justify-center px-4 py-2.5 bg-[#131B2E] border border-slate-800 hover:border-slate-700 rounded-xl font-bold text-xs uppercase tracking-widest transition-all flex items-center gap-2 shadow-md shrink-0">
				<svg :class="{ 'animate-spin': isLoading }" class="w-4 h-4 text-indigo-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
					<path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182m-4.105 0h4.992" />
				</svg>
				Refresh Log
			</button>
		</div>

		<div class="mb-6 relative">
			<span class="absolute inset-y-0 left-0 flex items-center pl-4 pointer-events-none">
				<svg class="w-4 h-4 text-slate-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
					<path stroke-linecap="round" stroke-linejoin="round" d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.602 10.602Z" />
				</svg>
			</span>
			<input v-model="searchQuery" type="text" placeholder="Cari jejak log aktor, jenis tindakan, atau detail..." class="w-full pl-11 pr-5 py-3.5 bg-[#131B2E] border border-slate-800 rounded-2xl text-white text-xs font-bold focus:outline-none focus:border-indigo-500 transition-all placeholder:text-slate-600 shadow-xl" />
		</div>

		<div class="hidden md:block bg-[#131B2E] border border-slate-800 rounded-[24px] overflow-hidden shadow-2xl">
			<div class="overflow-x-auto">
				<table class="w-full text-left border-collapse">
					<thead>
						<tr class="border-b border-slate-800 bg-[#1a243d]/50 text-[10px] font-black text-slate-400 uppercase tracking-widest">
							<th class="p-5">Waktu Kejadian</th>
							<th class="p-5">Aktor / Pengguna</th>
							<th class="p-5">Tindakan (Action)</th>
							<th class="p-5">Detail Deskripsi</th>
							<th class="p-5">Alamat IP</th>
							<th class="p-5">Perangkat (User Agent)</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-slate-800 text-xs">
						<tr v-for="log in logs" :key="log.id" class="hover:bg-[#1a243d]/20 transition-all text-slate-300">
							<td class="p-5 font-mono text-indigo-400 font-bold whitespace-nowrap">
								{{ formatDate(log.created_at) }}
							</td>
							<td class="p-5">
								<div class="font-black text-white">{{ log.user_name }}</div>
								<div class="text-[10px] text-slate-500 mt-0.5">{{ log.user_email }}</div>
							</td>
							<td class="p-5">
								<span class="px-2.5 py-1 bg-indigo-500/10 border border-indigo-500/20 text-indigo-300 rounded-lg font-black uppercase tracking-wide text-[10px]">
									{{ log.action }}
								</span>
							</td>
							<td class="p-5 font-medium max-w-xs break-words">
								{{ log.details }}
							</td>
							<td class="p-5 font-mono text-slate-400">
								{{ log.ip_address }}
							</td>
							<td class="p-5 text-slate-500 truncate max-w-[150px]" :title="log.user_agent">
								{{ log.user_agent }}
							</td>
						</tr>
					</tbody>
				</table>
			</div>
		</div>

		<div class="block md:hidden space-y-4">
			<div v-for="log in logs" :key="log.id" class="bg-[#131B2E] border border-slate-800 rounded-2xl p-4 shadow-xl space-y-3 relative overflow-hidden">
				<div class="flex justify-between items-center border-b border-slate-800/60 pb-2">
					<span class="text-[10px] font-mono text-indigo-400 font-black">
						{{ formatDate(log.created_at) }}
					</span>
					<span class="text-[10px] font-mono text-slate-500">IP: {{ log.ip_address }}</span>
				</div>

				<div class="flex flex-col gap-2">
					<div>
						<div class="text-xs font-black text-white">{{ log.user_name }}</div>
						<div class="text-[10px] text-slate-500">{{ log.user_email }}</div>
					</div>
					<div>
						<span class="inline-block px-2 py-0.5 bg-indigo-500/10 border border-indigo-500/20 text-indigo-300 rounded-md font-black uppercase tracking-wide text-[9px]">
							{{ log.action }}
						</span>
					</div>
				</div>

				<div class="bg-[#0B0F19]/60 border border-slate-800/40 rounded-xl p-3 text-[11px] text-slate-300 leading-relaxed font-medium break-words">
					{{ log.details }}
				</div>

				<div class="text-[9px] text-slate-600 font-medium truncate" :title="log.user_agent">UA: {{ log.user_agent }}</div>
			</div>
		</div>

		<div v-if="logs.length === 0 && !isLoading" class="p-10 bg-[#131B2E] border border-slate-800 rounded-[24px] text-center shadow-xl">
			<svg class="w-10 h-10 text-slate-600 mx-auto mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
				<path stroke-linecap="round" stroke-linejoin="round" d="M12 6v6h4.5m4.5 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" />
			</svg>
			<p class="text-xs text-slate-500 font-bold uppercase tracking-wider">Tidak ada aliran jejak log audit yang terdeteksi .</p>
		</div>
	</div>
</template>
