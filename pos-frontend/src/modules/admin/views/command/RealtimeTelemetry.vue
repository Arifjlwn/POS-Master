<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import api from '../../../../api.js';

const stats = ref({ tenant_online: 0, kasir_online: 0, total_stores: 0 });
let intervalId = null;

const fetchStats = async () => {
	try {
		const res = await api.get('/admin/command/telemetry');
		if (res.data.status === 'sukses') stats.value = res.data.data;
	} catch (e) {
		console.error('Gagal sync telemetry');
	}
};

onMounted(() => {
	fetchStats();
	intervalId = setInterval(fetchStats, 5000); // Polling setiap 5 detik
});

onUnmounted(() => clearInterval(intervalId));
</script>

<template>
	<div class="p-6 space-y-6">
		<h2 class="text-white font-black text-lg uppercase tracking-widest">Realtime Telemetry</h2>

		<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
			<div class="bg-[#131B2E] border border-slate-800 p-6 rounded-2xl">
				<p class="text-slate-500 font-bold text-[10px] uppercase">Tenant Online</p>
				<h3 class="text-4xl font-black text-emerald-400 mt-2">{{ stats.tenant_online }}</h3>
			</div>
			<div class="bg-[#131B2E] border border-slate-800 p-6 rounded-2xl">
				<p class="text-slate-500 font-bold text-[10px] uppercase">Kasir Online</p>
				<h3 class="text-4xl font-black text-indigo-400 mt-2">{{ stats.kasir_online }}</h3>
			</div>
			<div class="bg-[#131B2E] border border-slate-800 p-6 rounded-2xl">
				<p class="text-slate-500 font-bold text-[10px] uppercase">Total Terdaftar</p>
				<h3 class="text-4xl font-black text-white mt-2">{{ stats.total_stores }}</h3>
			</div>
		</div>
	</div>
</template>
