<script setup>
import L from 'leaflet';
import { nextTick, onMounted, ref, watch } from 'vue';

import 'leaflet/dist/leaflet.css';

const props = defineProps({
	form: Object,
	logoPreview: String,
});

const emit = defineEmits(['update-file', 'remove-logo']);
const fileInput = ref(null);

const onLogoSelect = (e) => {
	const file = e.target.files[0];
	if (file) {
		if (file.size > 2 * 1024 * 1024) {
			alert('Maksimal ukuran file adalah 2MB');
			return;
		}
		emit('update-file', 'logo', file, URL.createObjectURL(file));
	}
};

const removeImage = () => {
	emit('remove-logo');
	if (fileInput.value) fileInput.value.value = '';
};

// State Peta Leaflet
const mapContainer = ref(null);
let map = null;
let marker = null;

// Inisialisasi Peta
const initMap = () => {
	if (!mapContainer.value) return;

	const initLat = parseFloat(props.form.latitude) || -6.224168;
	const initLng = parseFloat(props.form.longitude) || 106.864388;

	map = L.map(mapContainer.value, {
		center: [initLat, initLng],
		zoom: 15,
		zoomControl: false,
	});

	L.control.zoom({ position: 'bottomright' }).addTo(map);

	L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
		attribution: '&copy; OpenStreetMap contributors',
	}).addTo(map);

	marker = L.marker([initLat, initLng], { draggable: isEditingLocation.value }).addTo(map);

	marker.on('dragend', function () {
		const position = marker.getLatLng();
		props.form.latitude = position.lat;
		props.form.longitude = position.lng;
	});
};

// Geocoding Pasif: Mengubah posisi peta tanpa merusak state form induk bray
const searchLocationByAddress = async () => {
	if (!props.form.kelurahan || !props.form.alamat) return;

	const alamatLengkap = `${props.form.alamat}, ${props.form.kelurahan}, ${props.form.kecamatan}, ${props.form.kota}, ${props.form.provinsi}`;

	try {
		const response = await fetch(`https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(alamatLengkap)}&limit=1`);
		const data = await response.json();

		if (data && data.length > 0) {
			const lat = parseFloat(data[0].lat);
			const lon = parseFloat(data[0].lon);

			// 🚀 AMAN: Update koordinat form lokal tanpa memicu re-render reset form induk
			props.form.latitude = lat;
			props.form.longitude = lon;

			if (map && marker) {
				map.setView([lat, lon], 16);
				marker.setLatLng([lat, lon]);
				map.invalidateSize();
			}
		}
	} catch (error) {
		console.error('Gagal tracking koordinat alamat:', error);
	}
};

// AUTO-DETECT EDIT MODE
const isEditingLocation = ref(false);
const toggleEditLocation = () => {
	isEditingLocation.value = !isEditingLocation.value;

	if (marker) {
		if (isEditingLocation.value) {
			marker.dragging.enable();
		} else {
			marker.dragging.disable();
			const resetLat = parseFloat(props.form.latitude) || -6.224168;
			const resetLng = parseFloat(props.form.longitude) || 106.864388;
			map.setView([resetLat, resetLng], 15);
			marker.setLatLng([resetLat, resetLng]);
		}
	}

	nextTick(() => {
		if (map) {
			// Beri nafas 100ms agar DOM Tailwind selesai render transisi opacity bray
			setTimeout(() => {
				map.invalidateSize();
			}, 100);
		}
	});
};

// INTEGRASI API WILAYAH (IBNUX)
const listProvinsi = ref([]);
const listKota = ref([]);
const listKecamatan = ref([]);
const listKelurahan = ref([]);

const regIds = ref({ provinsi: '', kota: '', kecamatan: '', kelurahan: '' });
const isLoading = ref({ reg: false, dist: false, vil: false });

const cleanStr = (str) => {
	return str ? String(str).trim().toUpperCase() : '';
};

const matchWilayah = (apiName, dbName) => {
	if (!apiName || !dbName) return false;
	let api = cleanStr(apiName)
		.replace(/KABUPATEN |KOTA |ADMINISTRASI /g, '')
		.replace('DAERAH KHUSUS IBUKOTA JAKARTA', 'DKI JAKARTA')
		.trim();
	let db = cleanStr(dbName)
		.replace(/KABUPATEN |KOTA |ADMINISTRASI /g, '')
		.replace('DAERAH KHUSUS IBUKOTA JAKARTA', 'DKI JAKARTA')
		.trim();
	return api === db;
};

const loadKota = async (idProvinsi, targetName = null) => {
	isLoading.value.reg = true;
	try {
		const res = await fetch(`https://ibnux.github.io/data-indonesia/kabupaten/${idProvinsi}.json`);
		const data = await res.json();
		listKota.value = data.map((item) => ({ id: item.id, name: item.nama.toUpperCase() }));

		if (targetName) {
			const found = listKota.value.find((r) => matchWilayah(r.name, targetName));
			if (found) {
				regIds.value.kota = found.id;
				props.form.kota = found.name;
				return found.id;
			}
		}
		return null;
	} catch (err) {
		console.error(err);
		return null;
	} finally {
		isLoading.value.reg = false;
	}
};

const loadKecamatan = async (idKota, targetName = null) => {
	isLoading.value.dist = true;
	try {
		const res = await fetch(`https://ibnux.github.io/data-indonesia/kecamatan/${idKota}.json`);
		const data = await res.json();
		listKecamatan.value = data.map((item) => ({ id: item.id, name: item.nama.toUpperCase() }));

		if (targetName) {
			const found = listKecamatan.value.find((d) => matchWilayah(d.name, targetName));
			if (found) {
				regIds.value.kecamatan = found.id;
				props.form.kecamatan = found.name;
				return found.id;
			}
		}
		return null;
	} catch (err) {
		console.error(err);
		return null;
	} finally {
		isLoading.value.dist = false;
	}
};

const loadKelurahan = async (idKecamatan, targetName = null) => {
	isLoading.value.vil = true;
	try {
		const res = await fetch(`https://ibnux.github.io/data-indonesia/kelurahan/${idKecamatan}.json`);
		const data = await res.json();
		listKelurahan.value = data.map((item) => ({ id: item.id, name: item.nama.toUpperCase() }));

		if (targetName) {
			const found = listKelurahan.value.find((v) => matchWilayah(v.name, targetName));
			if (found) {
				regIds.value.kelurahan = found.id;
				props.form.kelurahan = found.name;
				return found.id;
			}
		}
		return null;
	} catch (err) {
		console.error(err);
		return null;
	} finally {
		isLoading.value.vil = false;
	}
};

const initializeWilayah = async () => {
	if (!props.form.provinsi) {
		isEditingLocation.value = true;
		return;
	}
	const prov = listProvinsi.value.find((p) => matchWilayah(p.name, props.form.provinsi));
	if (!prov) {
		isEditingLocation.value = true;
		return;
	}

	regIds.value.provinsi = prov.id;
	props.form.provinsi = prov.name;

	if (!props.form.kota) {
		isEditingLocation.value = true;
		return;
	}
	const idKota = await loadKota(prov.id, props.form.kota);
	if (!idKota) {
		isEditingLocation.value = true;
		return;
	}

	if (!props.form.kecamatan) {
		isEditingLocation.value = true;
		return;
	}
	const idKecamatan = await loadKecamatan(idKota, props.form.kecamatan);
	if (!idKecamatan) {
		isEditingLocation.value = true;
		return;
	}

	if (props.form.kelurahan) {
		await loadKelurahan(idKecamatan, props.form.kelurahan);
	} else {
		isEditingLocation.value = true;
	}

	// 🚀 LOGIKA AMAN: Jika koordinat 0, geser petanya tanpa memicu loop resetting form bray
	const currentLat = parseFloat(props.form.latitude) || 0;
	const currentLng = parseFloat(props.form.longitude) || 0;

	if (currentLat === 0 || currentLng === 0) {
		await searchLocationByAddress();
	} else {
		if (map && marker) {
			map.setView([currentLat, currentLng], 16);
			marker.setLatLng([currentLat, currentLng]);
			nextTick(() => {
				if (map) map.invalidateSize();
			});
		}
	}
};

onMounted(async () => {
	try {
		const res = await fetch(`https://ibnux.github.io/data-indonesia/provinsi.json`);
		const data = await res.json();
		listProvinsi.value = data.map((item) => ({ id: item.id, name: item.nama.toUpperCase() }));
	} catch (e) {
		console.error('Gagal load provinsi dari IBNUX:', e);
	}

	nextTick(() => {
		initMap();
	});
});

watch(
	() => ({
		provinsi: props.form.provinsi,
		provincesLoaded: listProvinsi.value.length,
	}),
	async (val) => {
		if (val.provincesLoaded > 0 && !regIds.value.provinsi) {
			await initializeWilayah();
		}
	},
	{ immediate: true, deep: true }
);

const handleProvinsiChange = async () => {
	regIds.value.kota = '';
	regIds.value.kecamatan = '';
	regIds.value.kelurahan = '';
	listKota.value = [];
	listKecamatan.value = [];
	listKelurahan.value = [];
	props.form.kota = '';
	props.form.kecamatan = '';
	props.form.kelurahan = '';

	const prov = listProvinsi.value.find((p) => p.id === regIds.value.provinsi);
	if (prov) props.form.provinsi = prov.name;
	if (regIds.value.provinsi) await loadKota(regIds.value.provinsi);
};

const handleKotaChange = async () => {
	regIds.value.kecamatan = '';
	regIds.value.kelurahan = '';
	listKecamatan.value = [];
	listKelurahan.value = [];
	props.form.kecamatan = '';
	props.form.kelurahan = '';

	const kota = listKota.value.find((k) => k.id === regIds.value.kota);
	if (kota) props.form.kota = kota.name;
	if (regIds.value.kota) await loadKecamatan(regIds.value.kota);
};

const handleKecamatanChange = async () => {
	regIds.value.kelurahan = '';
	listKelurahan.value = [];
	props.form.kelurahan = '';
	const kec = listKecamatan.value.find((k) => k.id === regIds.value.kecamatan);
	if (kec) props.form.kecamatan = kec.name;
	if (regIds.value.kecamatan) {
		await loadKelurahan(regIds.value.kecamatan);
		setTimeout(() => {
			searchLocationByAddress();
		}, 500);
	}
};

const handleKelurahanChange = () => {
	const kel = listKelurahan.value.find((k) => k.id === regIds.value.kelurahan);
	if (kel) {
		props.form.kelurahan = kel.name;
		setTimeout(() => {
			searchLocationByAddress();
		}, 300);
	}
};

const forceUppercase = (field) => {
	if (props.form[field]) props.form[field] = props.form[field].toUpperCase();
};

watch(isEditingLocation, (newVal) => {
	if (newVal && map) {
		nextTick(() => {
			setTimeout(() => {
				map.invalidateSize();
				// Opsional: sekalian centering ulang ke koordinat agar tidak melenceng
				const currentLat = parseFloat(props.form.latitude) || -6.224168;
				const currentLng = parseFloat(props.form.longitude) || 106.864388;
				map.setView([currentLat, currentLng], 15);
			}, 150);
		});
	}
});
</script>

<template>
	<div class="space-y-6 animate-fade-in-up w-full">
		<div class="flex flex-col md:flex-row gap-6">
			<div class="w-full md:w-1/3 flex flex-col items-center p-6 bg-slate-50 rounded-3xl border-2 border-slate-100">
				<label class="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-4">Logo Toko</label>
				<div class="relative w-32 h-32 bg-white rounded-[24px] shadow-sm flex items-center justify-center overflow-hidden border-2 border-slate-200 group hover:border-blue-400 transition-colors">
					<div v-if="logoPreview" class="w-full h-full relative cursor-default">
						<img :src="logoPreview" class="w-full h-full object-contain p-2" />
						<button @click.prevent="removeImage" class="absolute inset-0 bg-rose-500/80 text-white flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity backdrop-blur-sm cursor-pointer z-10">
							<svg class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
						</button>
					</div>
					<div v-else @click="$refs.fileInput.click()" class="w-full h-full flex items-center justify-center cursor-pointer">
						<svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8 text-slate-300 group-hover:text-blue-500 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" /></svg>
					</div>
				</div>
				<input type="file" ref="fileInput" @change="onLogoSelect" class="hidden" accept="image/*" />
				<p class="text-[9px] text-slate-400 font-bold mt-4 text-center">Format: JPG, PNG (Max 2MB)</p>
			</div>

			<div class="w-full md:w-2/3 grid grid-cols-1 md:grid-cols-2 gap-4">
				<div class="md:col-span-2">
					<label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Nama Toko / Outlet</label>
					<input v-model="form.nama_toko" @input="forceUppercase('nama_toko')" type="text" class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:bg-white focus:border-blue-600 focus:shadow-lg focus:shadow-blue-500/10 outline-none font-bold text-sm text-slate-800 transition-all" />
				</div>
				<div>
					<label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">No. WA Bisnis</label>
					<input v-model="form.telepon" type="tel" class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:bg-white focus:border-blue-600 focus:shadow-lg focus:shadow-blue-500/10 outline-none font-bold text-sm text-slate-800 transition-all" placeholder="08xxxxxxxx" />
				</div>
				<div>
					<label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Kategori Bisnis</label>
					<div class="w-full p-4 bg-slate-100/80 rounded-2xl border-2 border-slate-100 flex items-center justify-between cursor-not-allowed">
						<span class="font-black text-xs text-slate-400 uppercase tracking-widest">{{ form.business_type || 'RETAIL' }}</span>
						<svg class="w-4 h-4 text-slate-300" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" /></svg>
					</div>
				</div>
			</div>
		</div>

		<div class="border-t border-slate-100 pt-6 space-y-4 w-full">
			<div class="flex items-center justify-between mb-2">
				<h4 class="font-black text-slate-800 uppercase tracking-widest text-sm">Lokasi Administrasi & Radar</h4>
				<button @click.prevent="toggleEditLocation" class="text-[10px] font-black uppercase tracking-widest flex items-center gap-1 transition-colors" :class="isEditingLocation ? 'text-rose-500 hover:text-rose-600' : 'text-blue-500 hover:text-blue-600'">
					<svg v-if="!isEditingLocation" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" /></svg>
					<svg v-else xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
					{{ isEditingLocation ? 'Batal Ubah' : 'Ubah Lokasi' }}
				</button>
			</div>

			<div v-if="isEditingLocation" class="bg-amber-50 border border-amber-200 text-amber-700 p-3 rounded-xl text-xs font-bold flex items-start gap-2 mb-4">
				<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" /></svg>
				<p>Mengubah lokasi akan mempengaruhi sistem radar peta area dan data zonasi cabang.</p>
			</div>

			<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-12 gap-4 w-full" :class="{ 'opacity-60 pointer-events-none': !isEditingLocation }">
				<div class="col-span-1 sm:col-span-2 md:col-span-12">
					<label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Alamat Lengkap (Jalan / Patokan)</label>
					<textarea v-model="form.alamat" @blur="searchLocationByAddress" @input="forceUppercase('alamat')" :disabled="!isEditingLocation" rows="2" class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:bg-white focus:border-blue-600 focus:shadow-lg focus:shadow-blue-500/10 outline-none font-bold text-sm text-slate-800 transition-all resize-none disabled:bg-slate-100/50" placeholder="Contoh: Jl. Raya Frankfurt Blok B No. 9..."></textarea>
				</div>

				<div class="relative w-full min-w-0 col-span-1 sm:col-span-1 md:col-span-5">
					<label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Provinsi</label>
					<select v-model="regIds.provinsi" @change="handleProvinsiChange" :disabled="!isEditingLocation" class="w-full p-4 pr-10 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:bg-white focus:border-blue-600 focus:shadow-lg focus:shadow-blue-500/10 outline-none font-bold text-sm text-slate-800 transition-all appearance-none cursor-pointer disabled:bg-slate-100/50">
						<option value="" disabled>Pilih Provinsi</option>
						<option v-for="p in listProvinsi" :key="p.id" :value="p.id">{{ p.name }}</option>
					</select>
					<div class="absolute inset-y-0 right-4 flex items-center top-[28px] pointer-events-none text-slate-400">
						<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7"></path></svg>
					</div>
				</div>

				<div class="relative w-full min-w-0 col-span-1 sm:col-span-1 md:col-span-7">
					<label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Kota / Kabupaten</label>
					<select v-model="regIds.kota" @change="handleKotaChange" :disabled="!regIds.provinsi || isLoading.reg || !isEditingLocation" class="w-full p-4 pr-10 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:bg-white focus:border-blue-600 focus:shadow-lg focus:shadow-blue-500/10 outline-none font-bold text-sm text-slate-800 transition-all appearance-none disabled:opacity-60 disabled:cursor-not-allowed cursor-pointer disabled:bg-slate-100/50">
						<option value="" disabled>{{ isLoading.reg ? 'Loading...' : 'Pilih Kota/Kab' }}</option>
						<option v-for="r in listKota" :key="r.id" :value="r.id">{{ r.name }}</option>
					</select>
					<div class="absolute inset-y-0 right-4 flex items-center top-[28px] pointer-events-none text-slate-400">
						<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7"></path></svg>
					</div>
				</div>

				<div class="relative w-full min-w-0 col-span-1 sm:col-span-1 md:col-span-4">
					<label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Kecamatan</label>
					<select v-model="regIds.kecamatan" @change="handleKecamatanChange" :disabled="!regIds.kota || isLoading.dist || !isEditingLocation" class="w-full p-4 pr-10 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:bg-white focus:border-blue-600 focus:shadow-lg focus:shadow-blue-500/10 outline-none font-bold text-sm text-slate-800 transition-all appearance-none disabled:opacity-60 disabled:cursor-not-allowed cursor-pointer disabled:bg-slate-100/50">
						<option value="" disabled>{{ isLoading.dist ? 'Loading...' : 'Pilih Kecamatan' }}</option>
						<option v-for="d in listKecamatan" :key="d.id" :value="d.id">{{ d.name }}</option>
					</select>
					<div class="absolute inset-y-0 right-4 flex items-center top-[28px] pointer-events-none text-slate-400">
						<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7"></path></svg>
					</div>
				</div>

				<div class="relative w-full min-w-0 col-span-1 sm:col-span-1 md:col-span-5">
					<label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Kelurahan / Desa</label>
					<select v-model="regIds.kelurahan" @change="handleKelurahanChange" :disabled="!regIds.kecamatan || isLoading.vil || !isEditingLocation" class="w-full p-4 pr-10 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:bg-white focus:border-blue-600 focus:shadow-lg focus:shadow-blue-500/10 outline-none font-bold text-sm text-slate-800 transition-all appearance-none disabled:opacity-60 disabled:cursor-not-allowed cursor-pointer disabled:bg-slate-100/50">
						<option value="" disabled>{{ isLoading.vil ? 'Loading...' : 'Pilih Kelurahan' }}</option>
						<option v-for="v in listKelurahan" :key="v.id" :value="v.id">{{ v.name }}</option>
					</select>
					<div class="absolute inset-y-0 right-4 flex items-center top-[28px] pointer-events-none text-slate-400">
						<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7"></path></svg>
					</div>
				</div>

				<div class="w-full min-w-0 col-span-1 sm:col-span-2 md:col-span-3">
					<label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Kode Pos</label>
					<input v-model="form.kode_pos" :disabled="!isEditingLocation" type="text" class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:bg-white focus:border-blue-600 focus:shadow-lg focus:shadow-blue-500/10 outline-none font-black text-sm text-slate-800 transition-all disabled:bg-slate-100/50" placeholder="Kode Pos" />
				</div>
			</div>

			<div class="w-full mt-6">
				<label class="text-[10px] font-black uppercase tracking-widest ml-1 mb-2 flex items-center gap-1" :class="isEditingLocation ? 'text-indigo-600' : 'text-slate-400'">
					<svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
						<path stroke-linecap="round" stroke-linejoin="round" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
					</svg>
					{{ isEditingLocation ? 'Radar Aktif: Silakan Geser Pin Toko' : 'Radar Terkunci (Klik "Ubah Lokasi" Untuk Menggeser Pin)' }}
				</label>

				<div ref="mapContainer" class="w-full h-60 bg-slate-100 rounded-2xl border-2 border-slate-200 shadow-inner z-10 overflow-hidden transition-all" :class="{ 'border-indigo-400 ring-4 ring-indigo-500/10': isEditingLocation }"></div>

				<div class="mt-2 text-[9px] font-bold text-slate-400 text-right uppercase tracking-wider">
					Radar Koordinat:
					<span class="text-slate-600 font-mono font-black">{{ form.latitude || 0 }}, {{ form.longitude || 0 }}</span>
				</div>
			</div>
		</div>
	</div>
</template>
