<script setup>
import { computed, inject } from 'vue';

import MetricCard from '../components/widgets/MetricCard.vue';
import ServerTelemetryCard from '../components/widgets/ServerTelemetryCard.vue';

// 🚀 IMPORT RADAR KASTA TERTINGGI (Pengganti RecentActivityTable)
import LiveTenantMap from '../components/widgets/LiveTenantMap.vue';

// 🚀 JURUS INJECT: Ambil data realtime terpusat dari AdminLayout.vue !
const { telemetryData, isLoadingTelemetry, refreshTelemetry } = inject('globalTelemetry');

// 🚀 FIX REALTIME MAPPING: Disamain persis sama struktur response JSON backend !
const metrics = computed(() => {
	// Backend ngereturn objek pembungkus utama bernama "data" untuk status ruko
	const dataResponse = telemetryData.value?.data || {};
	const live = telemetryData.value?.live_stats || {};

	return {
		// 1. Statistik Tenant Berdasarkan Status Kontrak (Mapping Langsung ke Properti "data")
		active_tenants: dataResponse.active_stores || 0,
		trial_tenants: dataResponse.trial_stores || 0,
		suspended_tenants: dataResponse.suspended_stores || 0,
		archived_tenants: dataResponse.archived_stores || 0,

		// 2. Realtime Telemetri Aktivitas Operasional Merchant
		users_online: live.users_online || 0,
		global_transactions: live.total_transactions_today || 'Rp 0',
		active_cashiers: live.active_cashiers || 0,
		open_shifts: live.open_shifts || 0,
	};
});

// Pemetaan data kesehatan server riil dari runtime Go
const serverHealth = computed(() => {
	return (
		telemetryData.value?.server_health || {
			cpu_usage: 0,
			ram_usage: 0,
			db_status: 'Connecting...',
			latency: '0ms',
		}
	);
});

// 🗺️ STREAM LIVE MAP NODES: Membaca titik koordinat radar dari backend
const liveMapNodes = computed(() => {
	return telemetryData.value?.live_map_nodes || [];
});

const isLoading = computed(() => isLoadingTelemetry.value);
</script>

<template>
	<div class="w-full min-h-screen bg-[#0B0F19] p-3 sm:p-5 lg:p-6 font-sans select-none text-white relative overflow-x-hidden">
		<div class="absolute w-[300px] h-[300px] sm:w-[600px] sm:h-[600px] bg-indigo-900/10 rounded-full blur-[80px] sm:blur-[140px] -top-20 sm:-top-40 -right-20 sm:-right-40 pointer-events-none"></div>

		<div class="flex flex-col lg:flex-row justify-between items-start lg:items-center mb-6 sm:mb-8 lg:mb-10 gap-4 relative z-10">
			<div class="w-full lg:w-auto">
				<div class="inline-flex items-center gap-2 px-3 py-1 bg-indigo-500/10 border border-indigo-500/20 text-indigo-400 rounded-xl text-[10px] font-black uppercase tracking-widest mb-2.5">Core Platform Control v1.0</div>
				<h1 class="text-xl sm:text-2xl md:text-3xl font-black text-white tracking-tight break-words">ARZURA MISSION CONTROL</h1>
				<p class="text-slate-500 font-bold text-[9px] sm:text-[10px] uppercase tracking-widest mt-1 leading-relaxed">Realtime Multi-Tenant Telemetry Monitoring System</p>
			</div>

			<button @click="refreshTelemetry" :disabled="isLoading" class="w-full lg:w-auto justify-center px-5 py-3 bg-[#131B2E] border border-slate-800 hover:border-slate-700 active:scale-[0.98] disabled:opacity-50 text-white rounded-2xl font-black text-[10px] uppercase tracking-widest transition-all shadow-xl flex items-center gap-2 shrink-0">
				<svg :class="{ 'animate-spin': isLoading }" class="w-3.5 h-3.5 text-indigo-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
					<path stroke-linecap="round" stroke-linejoin="round" d="M4 4v5h.582m15.356 2A8.001 8.001 0 1121.253 8H18" />
				</svg>
				Sync Telemetri
			</button>
		</div>

		<div class="mb-4 text-xs font-black uppercase tracking-widest text-slate-500 relative z-10 px-1">Subscription Matrix</div>
		<div class="grid grid-cols-2 lg:grid-cols-4 gap-4 sm:gap-6 relative z-10 mb-6 sm:mb-8">
			<MetricCard title="Tenant Aktif" :value="isLoading ? '...' : metrics.active_tenants" icon-color="bg-emerald-500/10 text-emerald-400 border-emerald-500/20">
				<template #icon>
					<svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
						<path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z" />
					</svg>
				</template>
			</MetricCard>

			<MetricCard title="Tenant Masa Trial" :value="isLoading ? '...' : metrics.trial_tenants" icon-color="bg-indigo-500/10 text-indigo-400 border-indigo-500/20">
				<template #icon>
					<svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
						<path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
					</svg>
				</template>
			</MetricCard>

			<MetricCard title="Tenant Suspended" :value="isLoading ? '...' : metrics.suspended_tenants" icon-color="bg-amber-500/10 text-amber-400 border-amber-500/20">
				<template #icon>
					<svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
						<path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
					</svg>
				</template>
			</MetricCard>

			<MetricCard title="Tenant Archived" :value="isLoading ? '...' : metrics.archived_tenants" icon-color="bg-rose-500/10 text-rose-400 border-rose-500/20">
				<template #icon>
					<svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
						<path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-10V4a1 1 0 00-1-1h-2a1 1 0 00-1 1v3M4 7h16" />
					</svg>
				</template>
			</MetricCard>
		</div>

		<div class="mb-4 text-xs font-black uppercase tracking-widest text-slate-500 relative z-10 px-1">Live Engine Matrix</div>
		<div class="grid grid-cols-1 sm:grid-cols-3 gap-4 sm:gap-6 relative z-10 mb-6 sm:mb-8">
			<MetricCard title="User / Admin Online" :value="isLoading ? '...' : metrics.users_online" icon-color="bg-cyan-500/10 text-cyan-400 border-cyan-500/20">
				<template #icon>
					<svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
						<path stroke-linecap="round" stroke-linejoin="round" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
					</svg>
				</template>
			</MetricCard>

			<MetricCard title="Volume Transaksi Global (All Time)" :value="isLoading ? '...' : metrics.global_transactions" icon-color="bg-emerald-500/10 text-emerald-400 border-emerald-500/20">
				<template #icon>
					<svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
						<path stroke-linecap="round" stroke-linejoin="round" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
					</svg>
				</template>
			</MetricCard>

			<MetricCard title="Kasir Aktif Terdaftar" :value="isLoading ? '...' : metrics.active_cashiers" icon-color="bg-purple-500/10 text-purple-400 border-purple-500/20">
				<template #icon>
					<svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
						<path stroke-linecap="round" stroke-linejoin="round" d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
					</svg>
				</template>
			</MetricCard>
		</div>

		<div class="grid grid-cols-1 lg:grid-cols-3 gap-4 sm:gap-6 relative z-10 items-stretch min-h-[400px]">
			<div class="lg:col-span-2 flex flex-col h-full">
				<LiveTenantMap :nodes="liveMapNodes" class="flex-1" />
			</div>

			<div class="flex flex-col h-full">
				<ServerTelemetryCard :health="serverHealth" class="flex-1" />
			</div>
		</div>

		<div class="mt-6 bg-[#131B2E] border border-slate-800 rounded-2xl p-4 text-center relative z-10 shadow-md">
			<p class="text-slate-500 font-black text-[9px] sm:text-[10px] uppercase tracking-widest leading-relaxed">Sistem Isolasi Multi-Tenant Berjalan Normal. Siap Membantai Skala 1.000 Tenant !</p>
		</div>
	</div>
</template>
