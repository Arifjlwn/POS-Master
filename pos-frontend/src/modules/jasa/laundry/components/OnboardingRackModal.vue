<script setup>
import Swal from 'sweetalert2';
import { ref } from 'vue';
import api from '../../../../api.js';

const props = defineProps({
	show: Boolean,
});

const emit = defineEmits(['setup-success', 'bypass']);

const step = ref(1); // Step 1: Pertanyaan Utama, Step 2: Input Matrix Grid
const jumlahBaris = ref(5); // Default 5 baris kebawah (A-E)
const jumlahKolom = ref(5); // Default 5 kolom ke kanan (1-5)
const isSubmitting = ref(false);

const handleGenerateRack = async () => {
	if (jumlahBaris.value > 26 || jumlahKolom.value > 50) {
		return Swal.fire('Batas Limit!', 'Maksimal baris adalah 26 (A-Z) dan kolom 50.', 'warning');
	}

	isSubmitting.value = true;
	try {
		const res = await api.post('/laundry/racks/setup', {
			jumlah_baris: parseInt(jumlahBaris.value, 10),
			jumlah_kolom: parseInt(jumlahKolom.value, 10),
		});

		if (res.data.status === 'sukses') {
			Swal.fire('Sukses!', 'Susunan laci rak laundry berhasil dibuat!', 'success');
			emit('setup-success');
		}
	} catch (err) {
		Swal.fire('Gagal Setup', err.response?.data?.error || 'Sirkuit backend terputus.', 'error');
	} finally {
		isSubmitting.value = false;
	}
};
</script>

<template>
	<div v-if="show" class="fixed inset-0 bg-slate-950/80 flex items-center justify-center z-[300] p-4 backdrop-blur-sm animate-fade-in">
		<div class="bg-white p-6 md:p-8 rounded-[32px] shadow-2xl w-full max-w-md border-[6px] border-slate-900 transition-all flex flex-col relative overflow-hidden">
			<!-- STICKER LAUNDRY -->
			<div class="absolute -right-6 -top-6 bg-indigo-600 text-white font-black text-[9px] uppercase tracking-widest px-8 py-3 rotate-45 border-b-4 border-slate-900 shadow-md">Laundry</div>

			<!-- STEP 1: PERTANYAAN UTAMA -->
			<div v-if="step === 1" class="text-center py-4">
				<div class="w-16 h-16 bg-indigo-50 text-indigo-600 rounded-2xl flex items-center justify-center mx-auto mb-4 border-2 border-slate-900 shadow-md">
					<svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
						<path stroke-linecap="round" stroke-linejoin="round" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
					</svg>
				</div>
				<h2 class="text-xl font-black text-slate-900 uppercase tracking-tight leading-tight">Konfigurasi Modul Rak</h2>
				<p class="text-xs font-bold text-slate-400 mt-2 px-2 leading-relaxed">Apakah cabang ruko Anda menggunakan sistem penomoran rak untuk mengorganisir pakaian konsumen?</p>

				<div class="flex gap-4 mt-6">
					<button @click="emit('bypass')" class="flex-1 bg-slate-100 hover:bg-slate-200 text-slate-500 py-3.5 rounded-2xl font-black text-xs uppercase tracking-wider active:scale-95 transition-all border-2 border-slate-200 flex items-center justify-center gap-2">
						<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
							<path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
						</svg>
						Lewati
					</button>
					<button @click="step = 2" class="flex-1 bg-indigo-600 hover:bg-indigo-700 text-white py-3.5 rounded-2xl font-black text-xs uppercase tracking-wider shadow-lg shadow-indigo-600/30 active:scale-95 transition-all border-2 border-slate-900 flex items-center justify-center gap-2">
						<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
							<path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
						</svg>
						Gunakan Rak
					</button>
				</div>
			</div>

			<!-- STEP 2: SETUP UKURAN MATRIX RAK -->
			<div v-if="step === 2" class="flex flex-col">
				<h3 class="text-lg font-black text-slate-900 uppercase tracking-tight mb-1">Ukuran Kapasitas Rak</h3>
				<p class="text-[11px] font-bold text-slate-400 mb-5 leading-normal">Tentukan jumlah susunan kotak penyimpanan. Sistem otomatis membuat identitas rak.</p>

				<div class="space-y-4 font-mono">
					<div class="flex flex-col gap-1.5">
						<label class="text-[10px] font-black text-slate-700 uppercase tracking-widest">Jumlah Baris Rak Kebawah (A - Z):</label>
						<input v-model.number="jumlahBaris" type="number" min="1" max="26" class="p-3.5 bg-slate-50 border-2 border-slate-200 rounded-xl outline-none font-black text-sm focus:border-indigo-600 transition-all text-slate-800" />
					</div>
					<div class="flex flex-col gap-1.5">
						<label class="text-[10px] font-black text-slate-700 uppercase tracking-widest">Jumlah Baris Rak Kekanan (1, 2, 3...):</label>
						<input v-model.number="jumlahKolom" type="number" min="1" max="50" class="p-3.5 bg-slate-50 border-2 border-slate-200 rounded-xl outline-none font-black text-sm focus:border-indigo-600 transition-all text-slate-800" />
					</div>
				</div>

				<!-- PREVIEW INFO BOX -->
				<div class="bg-indigo-50 border-2 border-dashed border-indigo-200 rounded-2xl p-3.5 mt-5 text-center flex flex-col items-center justify-center gap-2">
					<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-indigo-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
						<path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
					</svg>
					<p class="text-[11px] font-bold text-indigo-900 uppercase">
						Sistem membuat total
						<span class="font-black text-sm text-indigo-600 underline px-1">{{ jumlahBaris * jumlahKolom }} slot rak</span>
						(Mulai dari A-1 s/d {{ String.fromCharCode(64 + (jumlahBaris || 1)) }}-{{ jumlahKolom }}).
					</p>
				</div>

				<div class="flex gap-3 mt-6">
					<button @click="step = 1" :disabled="isSubmitting" class="bg-slate-100 p-4 rounded-2xl text-slate-500 font-black text-xs uppercase active:scale-95 transition-all flex justify-center items-center">
						<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
							<path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" />
						</svg>
					</button>
					<button @click="handleGenerateRack" :disabled="isSubmitting" class="flex-1 bg-indigo-600 hover:bg-indigo-700 text-white py-4 rounded-2xl font-black text-xs uppercase tracking-widest shadow-md active:scale-95 transition-all border-2 border-slate-900 flex justify-center items-center gap-2">
						<!-- SVG SPINNER LOADING -->
						<svg v-if="isSubmitting" class="animate-spin h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
							<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
							<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
						</svg>

						<!-- SVG PETIR GENERATE -->
						<svg v-else xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
							<path stroke-linecap="round" stroke-linejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
						</svg>

						{{ isSubmitting ? 'Memproses Data...' : 'Generate Rak' }}
					</button>
				</div>
			</div>
		</div>
	</div>
</template>
