<script setup>
import L from 'leaflet';
import { nextTick, onMounted, ref, watch } from 'vue';

const props = defineProps({
	formData: {
		type: Object,
		required: true,
	},
});

const BASE_REGIONAL_API = import.meta.env.VITE_BASE_REGIONAL_API;

const regIds = ref({ provinsi: '', kota: '', kecamatan: '', kelurahan: '' });
const listProvinsi = ref([]);
const listKota = ref([]);
const listKecamatan = ref([]);
const listKelurahan = ref([]);

// State Peta Leaflet
const mapContainer = ref(null);
let map = null;
let marker = null;

// Inisialisasi Peta Awal (Default pusat Jakarta / Indonesia)
const initMap = () => {
	if (!mapContainer.value) return;

	const initLat = parseFloat(props.form.latitude) && parseFloat(props.form.latitude) !== 0 ? parseFloat(props.form.latitude) : -6.2;
	const initLng = parseFloat(props.form.longitude) && parseFloat(props.form.longitude) !== 0 ? parseFloat(props.form.longitude) : 106.816666;

	map = L.map(mapContainer.value, {
		center: [initLat, initLng],
		zoom: 14,
		zoomControl: false,
	});

	// Tambah tombol zoom di kanan bawah biar rapi
	L.control.zoom({ position: 'bottomright' }).addTo(map);

	// Pake tile layer OpenStreetMap gratisan kasta tertinggi
	L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
		attribution: '&copy; OpenStreetMap contributors',
	}).addTo(map);

	// Bikin Marker yang bisa digeser-geser (draggable: true)
	marker = L.marker([initLat, initLng], { draggable: true }).addTo(map);

	// Event pas marker selesai digeser manual sama owner
	marker.on('dragend', function (event) {
		const position = marker.getLatLng();
		props.formData.latitude = position.lat;
		props.formData.longitude = position.lng;
	});
};

// Fungsi Geocoding sederhana memanfaatkan nominatim gratisan untuk tracking wilayah
const seachLocationByAddress = async () => {
	const alamatLengkap = `${props.formData.alamat || ''}, ${props.formData.kelurahan || ''}, ${props.formData.kecamatan || ''}, ${props.formData.kota || ''}, ${props.formData.provinsi || ''}`;
	if (!props.formData.kelurahan) return;

	try {
		const response = await fetch(`https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(alamatLengkap)}&limit=1`);
		const data = await response.json();

		if (data && data.length > 0) {
			const lat = parseFloat(data[0].lat);
			const lon = parseFloat(data[0].lon);

			props.formData.latitude = lat;
			props.formData.longitude = lon;

			// Geser kamera peta dan marker ke lokasi baru hasil pencarian alamat
			if (map && marker) {
				map.setView([lat, lon], 16);
				marker.setLatLng([lat, lon]);
			}
		}
	} catch (error) {
		console.error('Gagal tracking koordinat alamat:', error);
	}
};

const loadProvinsi = async () => {
	try {
		const res = await fetch(`${BASE_REGIONAL_API}/provinsi.json`);
		const data = await res.json();
		listProvinsi.value = data.map((item) => ({ id: item.id, name: item.nama }));
	} catch (e) {
		console.error('Gagal memuat data provinsi:', e);
	}
};

watch(
	() => regIds.value.provinsi,
	async (newId) => {
		regIds.value.kota = '';
		listKota.value = [];
		regIds.value.kecamatan = '';
		listKecamatan.value = [];
		regIds.value.kelurahan = '';
		listKelurahan.value = [];
		props.formData.provinsi = listProvinsi.value.find((p) => p.id === newId)?.name || '';
		if (newId) {
			const res = await fetch(`${BASE_REGIONAL_API}/kabupaten/${newId}.json`);
			listKota.value = (await res.json()).map((item) => ({ id: item.id, name: item.nama }));
		}
	}
);

watch(
	() => regIds.value.kota,
	async (newId) => {
		regIds.value.kecamatan = '';
		listKecamatan.value = [];
		regIds.value.kelurahan = '';
		listKelurahan.value = [];
		props.formData.kota = listKota.value.find((p) => p.id === newId)?.name || '';
		if (newId) {
			const res = await fetch(`${BASE_REGIONAL_API}/kecamatan/${newId}.json`);
			listKecamatan.value = (await res.json()).map((item) => ({ id: item.id, name: item.nama }));
		}
	}
);

watch(
	() => regIds.value.kecamatan,
	async (newId) => {
		regIds.value.kelurahan = '';
		listKelurahan.value = [];
		props.formData.kecamatan = listKecamatan.value.find((p) => p.id === newId)?.name || '';
		if (newId) {
			const res = await fetch(`${BASE_REGIONAL_API}/kelurahan/${newId}.json`);
			listKelurahan.value = (await res.json()).map((item) => ({ id: item.id, name: item.nama }));
		}
	}
);

watch(
	() => regIds.value.kelurahan,
	(newId) => {
		props.formData.kelurahan = listKelurahan.value.find((p) => p.id === newId)?.name || '';
		if (newId) {
			// Pas kelurahan dipilih, lari cari koordinat kasarnya dulu buat penempatan peta awal bray
			setTimeout(() => {
				seachLocationByAddress();
			}, 500);
		}
	}
);

onMounted(() => {
	loadProvinsi();
	nextTick(() => {
		initMap();
	});
});
</script>

<template>
	<div class="flex flex-col gap-6 pt-4 border-t border-slate-100">
		<div class="flex items-center gap-3 border-b border-slate-100 pb-3">
			<div class="w-8 h-8 rounded-full bg-emerald-50 flex items-center justify-center text-emerald-600 font-black text-xs border border-emerald-100 shadow-sm">2</div>
			<h3 class="text-lg font-black text-slate-800 uppercase tracking-tight">Lokasi Operasional Cabang</h3>
		</div>

		<div class="grid grid-cols-1 md:grid-cols-2 gap-5 md:gap-6">
			<div class="md:col-span-2">
				<label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block">Detail Alamat Ruko (Jalan, No, Blok)</label>
				<textarea v-model="formData.alamat" @blur="seachLocationByAddress" rows="2" required class="block w-full px-5 py-4 bg-white border-2 border-slate-200 rounded-2xl focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 outline-none font-black text-slate-800 transition-all placeholder:text-slate-300 placeholder:font-bold shadow-sm resize-none uppercase" placeholder="Contoh: Ruko Frankfurt Blok C No. 12, Jl. Boulevard Raya..."></textarea>
			</div>

			<div class="relative">
				<label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block">Provinsi</label>
				<div class="relative">
					<select v-model="regIds.provinsi" required class="block w-full px-5 py-4 bg-white border-2 border-slate-200 rounded-2xl focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 outline-none font-black text-slate-800 transition-all shadow-sm cursor-pointer appearance-none">
						<option value="" disabled selected hidden>Pilih Provinsi...</option>
						<option v-for="prov in listProvinsi" :key="prov.id" :value="prov.id">{{ prov.name }}</option>
					</select>
					<div class="absolute inset-y-0 right-5 flex items-center pointer-events-none text-slate-400">
						<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
							<path d="m6 9 6 6 6-6" />
						</svg>
					</div>
				</div>
			</div>

			<div class="relative">
				<label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block">Kota / Kabupaten</label>
				<div class="relative">
					<select v-model="regIds.kota" :disabled="!regIds.provinsi" required class="block w-full px-5 py-4 bg-white border-2 border-slate-200 rounded-2xl focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 outline-none font-black text-slate-800 transition-all shadow-sm cursor-pointer appearance-none disabled:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-60">
						<option value="" disabled selected hidden>{{ regIds.provinsi ? 'Pilih Kota...' : 'Pilih Provinsi Terlebih Dahulu' }}</option>
						<option v-for="kota in listKota" :key="kota.id" :value="kota.id">{{ kota.name }}</option>
					</select>
					<div class="absolute inset-y-0 right-5 flex items-center pointer-events-none text-slate-400">
						<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
							<path d="m6 9 6 6 6-6" />
						</svg>
					</div>
				</div>
			</div>

			<div class="relative">
				<label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block">Kecamatan</label>
				<div class="relative">
					<select v-model="regIds.kecamatan" :disabled="!regIds.kota" required class="block w-full px-5 py-4 bg-white border-2 border-slate-200 rounded-2xl focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 outline-none font-black text-slate-800 transition-all shadow-sm cursor-pointer appearance-none disabled:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-60">
						<option value="" disabled selected hidden>{{ regIds.kota ? 'Pilih Kecamatan...' : 'Pilih Kota Terlebih Dahulu' }}</option>
						<option v-for="kec in listKecamatan" :key="kec.id" :value="kec.id">{{ kec.name }}</option>
					</select>
					<div class="absolute inset-y-0 right-5 flex items-center pointer-events-none text-slate-400">
						<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
							<path d="m6 9 6 6 6-6" />
						</svg>
					</div>
				</div>
			</div>

			<div class="grid grid-cols-2 gap-3">
				<div class="relative">
					<label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block">Desa / Kel</label>
					<div class="relative">
						<select v-model="regIds.kelurahan" :disabled="!regIds.kecamatan" required class="block w-full px-4 py-4 bg-white border-2 border-slate-200 rounded-2xl focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 outline-none font-black text-slate-800 transition-all shadow-sm cursor-pointer appearance-none disabled:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-60">
							<option value="" disabled selected hidden>Pilih Kelurahan...</option>
							<option v-for="kel in listKelurahan" :key="kel.id" :value="kel.id" class="text-xs">{{ kel.name }}</option>
						</select>
						<div class="absolute inset-y-0 right-4 flex items-center pointer-events-none text-slate-400">
							<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
								<path d="m6 9 6 6 6-6" />
							</svg>
						</div>
					</div>
				</div>

				<div>
					<label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block">Kode Pos</label>
					<input v-model="formData.kode_pos" :disabled="!regIds.kelurahan" type="number" required class="w-full px-4 py-4 bg-white border-2 border-slate-200 rounded-2xl focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 outline-none font-black text-slate-800 transition-all text-sm placeholder:text-slate-300 disabled:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-60" placeholder="Kode Pos" />
				</div>
			</div>

			<div class="md:col-span-2">
				<label class="text-[10px] font-black text-indigo-600 uppercase tracking-widest ml-1 mb-2 flex items-center gap-1">
					<svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
						<path stroke-linecap="round" stroke-linejoin="round" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
						<path stroke-linecap="round" stroke-linejoin="round" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
					</svg>
					Pinpoint Alamat Akurat (Geser Pin Sesuai Posisi Ruko Anda)
				</label>

				<div ref="mapContainer" class="w-full h-64 sm:h-72 bg-slate-100 rounded-2xl border-2 border-slate-200 shadow-inner z-10 overflow-hidden"></div>

				<div class="mt-2 text-[10px] font-bold text-slate-400 text-right uppercase tracking-wider">
					Satelit Log:
					<span class="text-slate-600 font-mono font-black">{{ props.formData.latitude?.toFixed(6) || 0 }}, {{ props.formData.longitude?.toFixed(6) || 0 }}</span>
				</div>
			</div>
		</div>
	</div>
</template>
