<script setup>
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import { onMounted, onUnmounted, ref, watch } from 'vue';

const props = defineProps({
	nodes: {
		type: Array,
		default: () => [],
	},
});

const mapContainer = ref(null);
let map = null;
let markersLayer = null;

onMounted(() => {
	// 1. Inisialisasi Peta (Center awal di titik tengah Indonesia )
	map = L.map(mapContainer.value).setView([-2.5489, 118.0149], 5);

	// 2. Suntik Tema CartoDB Dark Matter (Radar Militer Kasta Tertinggi)
	L.tileLayer('https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png', {
		attribution: '&copy; OpenStreetMap &copy; CartoDB',
		subdomains: 'abcd',
		maxZoom: 19,
	}).addTo(map);

	markersLayer = L.layerGroup().addTo(map);
	renderRadarPins();
});

// Kalau ada data sinkronisasi baru dari backend, render ulang pin-nya!
watch(
	() => props.nodes,
	() => {
		renderRadarPins();
	},
	{ deep: true }
);

const renderRadarPins = () => {
	if (!markersLayer || !map) return;
	markersLayer.clearLayers(); // Bersihkan titik lama

	props.nodes.forEach((node) => {
		// Kita pakai CircleMarker biar mirip titik radar menyala, bukan icon pin standar
		const marker = L.circleMarker([node.latitude, node.longitude], {
			radius: 6,
			fillColor: node.color_code,
			color: node.color_code,
			weight: 2,
			opacity: 0.9,
			fillOpacity: 0.6,
		});

		// Racik HTML untuk Tooltip/Popup pas titiknya diklik
		const popupContent = `
            <div style="min-width: 160px; font-family: sans-serif;">
                <h4 style="margin: 0 0 4px 0; color: #fff; font-weight: 900; font-size: 13px; text-transform: uppercase;">${node.nama_toko}</h4>
                <div style="font-size: 10px; color: #94a3b8; margin-bottom: 10px; font-weight: bold; text-transform: uppercase; letter-spacing: 0.05em;">
                    ${node.owner_name} • <span style="color: ${node.color_code};">${node.plan} (${node.status})</span>
                </div>
                
                <div style="display: flex; justify-content: space-between; border-top: 1px dashed #334155; padding-top: 8px; margin-bottom: 4px;">
                    <span style="font-size: 10px; color: #64748b; font-weight: bold; text-transform: uppercase;">Kasir Aktif:</span>
                    <span style="font-size: 11px; color: #fff; font-weight: 900;">${node.active_cashier} Sesi</span>
                </div>
                
                <div style="display: flex; justify-content: space-between;">
                    <span style="font-size: 10px; color: #64748b; font-weight: bold; text-transform: uppercase;">Omzet Hari Ini:</span>
                    <span style="font-size: 11px; color: #10b981; font-weight: 900;">${node.omzet_today}</span>
                </div>
            </div>
        `;

		// Bind popup dengan custom class biar bisa kita styling jadi dark mode
		marker.bindPopup(popupContent, {
			className: 'arzura-dark-popup',
			closeButton: false, // Hilangin tombol silang biar clean
		});

		markersLayer.addLayer(marker);
	});
};

onUnmounted(() => {
	if (map) map.remove(); // Bersihin memory biar browser gak bocor
});
</script>

<template>
	<div class="bg-[#131B2E] border border-slate-800 rounded-[24px] sm:rounded-[28px] p-4 sm:p-5 shadow-xl flex flex-col h-full w-full overflow-hidden relative z-10">
		<div class="flex items-center justify-between mb-4 shrink-0">
			<div>
				<h3 class="text-xs sm:text-sm font-black text-white uppercase tracking-wide flex items-center gap-2">
					<svg class="w-4 h-4 text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
						<path stroke-linecap="round" stroke-linejoin="round" d="M3.75 13.5l10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75z" />
					</svg>
					Live Tenant Radar
				</h3>
				<p class="text-[9px] sm:text-[10px] text-slate-500 font-bold uppercase tracking-wider mt-1">Pemetaan Operasional Armada Ruko Realtime</p>
			</div>
			<span class="w-2 h-2 rounded-full bg-emerald-500 animate-ping shrink-0"></span>
		</div>

		<div ref="mapContainer" class="flex-1 w-full rounded-2xl z-0 border border-slate-800 overflow-hidden" style="min-height: 350px"></div>
	</div>
</template>

<style>
/* 🚀 JURUS SAKTI CSS: Ubah popup putih bawaan Leaflet jadi Dark Mode Premium */
.arzura-dark-popup .leaflet-popup-content-wrapper {
	background-color: #0b0f19 !important;
	border: 1px solid #1e293b !important;
	border-radius: 16px !important;
	box-shadow:
		0 20px 25px -5px rgb(0 0 0 / 0.5),
		0 8px 10px -6px rgb(0 0 0 / 0.5) !important;
}
.arzura-dark-popup .leaflet-popup-tip {
	background-color: #1e293b !important;
}
.arzura-dark-popup .leaflet-popup-content {
	margin: 14px 16px !important;
	line-height: 1.4 !important;
}
</style>
