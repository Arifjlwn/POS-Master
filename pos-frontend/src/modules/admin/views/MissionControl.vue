<script setup>
import { inject, computed } from 'vue';
import MetricCard from '../components/widgets/MetricCard.vue';

// Import widget modular terisolasi
import RecentActivityTable from '../components/widgets/RecentActivityTable.vue';
import ServerTelemetryCard from '../components/widgets/ServerTelemetryCard.vue';

// 🚀 JURUS INJECT: Ambil data realtime terpusat dari AdminLayout.vue bray!
const { telemetryData, isLoadingTelemetry, refreshTelemetry } = inject('globalTelemetry');

// Pemetaan data metrik tenant ruko retail terpusat (Computed reactive)
const metrics = computed(() => {
	const stats = telemetryData.value?.stats || {};
	return {
		total_tenants: stats.total_stores || 0,
		active_tenants: stats.active_stores || 0,
		pending_tenants: 0, // Fallback default 0 bray
		suspended_tenants: stats.suspended_stores || 0,
	};
});

// Pemetaan data kesehatan server riil dari runtime Go bray
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

// 🚀 STREAM LIVE AUDIT LOGS: Membaca aliran data log aktivitas riil dari database
const recentActivities = computed(() => {
	return telemetryData.value?.recent_activities || [];
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

		<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 sm:gap-6 relative z-10 mb-6 sm:mb-8">
			<MetricCard title="Total Tenant Terdaftar" :value="isLoading && metrics.total_tenants === 0 ? '...' : metrics.total_tenants">
				<template #icon>
					<svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
						<path stroke-linecap="round" stroke-linejoin="round" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
					</svg>
				</template>
			</MetricCard>

			<MetricCard title="Tenant Aktif (Premium)" :value="isLoading && metrics.active_tenants === 0 ? '...' : metrics.active_tenants" icon-color="bg-emerald-500/10 text-emerald-400 border-emerald-500/20">
				<template #icon>
					<svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
						<path stroke-linecap="round" stroke-linejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
					</svg>
				</template>
			</MetricCard>

			<MetricCard title="Aktivasi Tertunda" :value="isLoading && metrics.pending_tenants === 0 ? '...' : metrics.pending_tenants" icon-color="bg-amber-500/10 text-amber-400 border-amber-500/20">
				<template #icon>
					<svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
						<path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
					</svg>
				</template>
			</MetricCard>

			<MetricCard title="Tenant Dibekukan" :value="isLoading && metrics.suspended_tenants === 0 ? '...' : metrics.suspended_tenants" icon-color="bg-rose-500/10 text-rose-400 border-rose-500/20">
				<template #icon>
					<svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
						<path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
					</svg>
				</template>
			</MetricCard>
		</div>

		<div class="grid grid-cols-1 lg:grid-cols-3 gap-4 sm:gap-6 relative z-10 items-stretch">
			<div class="lg:col-span-2 flex flex-col">
				<RecentActivityTable :activities="recentActivities" class="flex-1" />
			</div>

			<div class="flex flex-col">
				<ServerTelemetryCard :health="serverHealth" class="flex-1" />
			</div>
		</div>

		<div class="mt-6 bg-[#131B2E] border border-slate-800 rounded-2xl p-4 text-center relative z-10 shadow-md">
			<p class="text-slate-500 font-black text-[9px] sm:text-[10px] uppercase tracking-widest leading-relaxed">Sistem Isolasi Multi-Tenant Berjalan Normal. Siap Membantai Skala 1.000 Tenant !</p>
		</div>
	</div>
</template>
