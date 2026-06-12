<script setup>
import { onMounted, onUnmounted, provide, ref } from 'vue';
import api from '../../../../api.js'; // Sesuaikan relative path ke api.js lu bray
import Sidebar from './Sidebar.vue';
import TopNavbar from './TopNavbar.vue';

// 🚀 STATE REALTIME INFRASTRUKTUR GLOBAL (STRUKTUR SIMETRIS DENGAN BACKEND & MISSION CONTROL)
const telemetryData = ref({
	data: {
		active_stores: 0,
		trial_stores: 0,
		suspended_stores: 0,
		archived_stores: 0,
	},
	live_stats: {
		users_online: 0,
		total_transactions_today: 'Rp 0',
		active_cashiers: 0,
		open_shifts: 0,
	},
	server_health: {
		cpu_usage: 0,
		ram_usage: 0,
		db_status: 'Connecting...',
		latency: '0ms',
	},
	recent_activities: [],
	live_map_nodes: [], // ◄ 🗺️ TAMBAHAN SAKTI 1: Wadah kosong buat titik radar
});

const isSidebarOpen = ref(false);
const isLoadingTelemetry = ref(false);
let globalTelemetryInterval = null;

// Fungsi terpusat menarik data dari Go Backend (Murni Realtime)
const fetchGlobalTelemetry = async (showSilently = false) => {
	if (!showSilently) isLoadingTelemetry.value = true;
	try {
		const res = await api.get('/admin/dashboard-stats');

		// 🚀 PROSES DISTRIBUSI TELEMETRI KASTA TERTINGGI: Tangkap semua payload tanpa ada yang dibuang!
		if (res.data?.status === 'sukses') {
			telemetryData.value = {
				data: res.data.data || { active_stores: 0, trial_stores: 0, suspended_stores: 0, archived_stores: 0 },
				live_stats: res.data.live_stats || { users_online: 0, total_transactions_today: 'Rp 0', active_cashiers: 0, open_shifts: 0 },
				server_health: res.data.server_health || { cpu_usage: 0, ram_usage: 0, db_status: 'Offline', latency: '0ms' },
				recent_activities: res.data.recent_activities || [],
				live_map_nodes: res.data.live_map_nodes || [], // ◄ 🗺️ TAMBAHAN SAKTI 2: Tangkap data titik radarnya dari API!
			};
		}
	} catch (err) {
		console.error('Gagal menarik sinkronisasi telemetri global bray:', err);
	} finally {
		if (!showSilently) isLoadingTelemetry.value = false;
	}
};

// 🚀 JURUS PROVIDE: Menyebarkan data realtime ke semua anak halaman (MissionControl dll)
provide('globalTelemetry', {
	telemetryData,
	isLoadingTelemetry,
	refreshTelemetry: () => fetchGlobalTelemetry(false),
});

// Fungsi memantau perubahan ukuran layar guna menjaga konsistensi UI
const handleResize = () => {
	if (window.innerWidth < 1024) {
		isSidebarOpen.value = false;
	}
};

onMounted(() => {
	window.addEventListener('resize', handleResize);

	// Eksekusi trigger telemetri kasta tertinggi pertama kali
	fetchGlobalTelemetry(false);

	// Silent polling tiap 10 detik secara background bray
	globalTelemetryInterval = setInterval(() => {
		fetchGlobalTelemetry(true);
	}, 10000);
});

onUnmounted(() => {
	window.removeEventListener('resize', handleResize);
	if (globalTelemetryInterval) clearInterval(globalTelemetryInterval);
});

const handleCloseSidebar = () => {
	if (window.innerWidth < 1024) {
		isSidebarOpen.value = false;
	}
};

const toggleSidebar = () => {
	isSidebarOpen.value = !isSidebarOpen.value;
};
</script>

<template>
	<div class="min-h-screen bg-[#0B0F19] text-slate-300 font-sans flex overflow-hidden selection:bg-indigo-500/30">
		<Sidebar :isOpen="isSidebarOpen" @close="handleCloseSidebar" />

		<div class="flex-1 flex flex-col h-screen overflow-hidden relative">
			<TopNavbar @toggle-sidebar="toggleSidebar" :db-status="telemetryData.server_health.db_status" :latency="telemetryData.server_health.latency" />

			<main class="flex-1 overflow-x-hidden overflow-y-auto bg-[#0B0F19] p-4 lg:p-8">
				<router-view v-slot="{ Component }">
					<transition name="fade" mode="out-in">
						<component :is="Component" />
					</transition>
				</router-view>
			</main>
		</div>
	</div>
</template>

<style>
/* Animasi transisi perpindahan halaman bray */
.fade-enter-active,
.fade-leave-active {
	transition:
		opacity 0.2s ease,
		transform 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
	opacity: 0;
	transform: translateY(10px);
}
</style>
