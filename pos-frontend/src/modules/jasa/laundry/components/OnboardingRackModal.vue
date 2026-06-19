<script setup>
import Swal from 'sweetalert2';
import { onMounted, ref } from 'vue';
import api from '../../../../api.js';

const props = defineProps({
	show: Boolean,
});

const emit = defineEmits(['setup-success', 'bypass']);

const step = ref(1);
const zonaNama = ref('RAK UTAMA');
const jumlahBaris = ref(5);
const jumlahKolom = ref(5);
const isSubmitting = ref(false);

// 🚀 STATE UNTUK NGECEK APAKAH INI RECOVERY DARI SNOOZE 7 HARI?
const isRecoverySnooze = ref(false);

onMounted(() => {
	// Cek apakah ada sisa jejak snooze di memori walau udah kedaluwarsa
	const snoozeTime = localStorage.getItem('snooze_rack_setup');
	if (snoozeTime) {
		isRecoverySnooze.value = true;
	}
});

const handleBypass = () => {
	// 🚀 ENGINE SNOOZE 7 HARI: Catat waktu saat ini + 7 hari (dalam milidetik)
	const snoozeTime = new Date().getTime() + 7 * 24 * 60 * 60 * 1000;
	localStorage.setItem('snooze_rack_setup', snoozeTime.toString());

	// Beri notifikasi halus ke user
	Swal.fire({
		toast: true,
		position: 'top-end',
		icon: 'info',
		title: 'Setup Rak ditunda selama 7 hari.',
		showConfirmButton: false,
		timer: 2000,
	});

	emit('bypass');
};

const handleGenerateRack = async () => {
	if (!zonaNama.value.trim()) {
		return Swal.fire('Data Tidak Lengkap', 'Nama Lemari Perdana wajib diisi.', 'warning');
	}
	if (jumlahBaris.value > 26 || jumlahKolom.value > 50) {
		return Swal.fire('Batas Limit!', 'Maksimal baris adalah 26 (A-Z) dan kolom 50.', 'warning');
	}

	isSubmitting.value = true;
	try {
		const res = await api.post('/laundry/racks/setup', {
			zona: zonaNama.value.trim().toUpperCase(),
			jumlah_baris: parseInt(jumlahBaris.value, 10),
			jumlah_kolom: parseInt(jumlahKolom.value, 10),
		});

		if (res.data.status === 'sukses') {
			// Hapus dosa snooze kalau dia akhirnya mau setup bray
			localStorage.removeItem('snooze_rack_setup');

			Swal.fire({
				icon: 'success',
				title: 'Konfigurasi Berhasil!',
				text: 'Susunan laci rak laundry telah mengudara!',
				customClass: { popup: 'rounded-[28px]' },
				confirmButtonColor: '#0f172a',
			});
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
	<div v-if="show" class="fixed inset-0 bg-slate-900/60 flex items-center justify-center z-[300] p-4 backdrop-blur-md animate-fade-in">
		<div class="bg-white/95 backdrop-blur-xl p-8 md:p-10 rounded-[36px] shadow-[0_0_60px_-15px_rgba(79,70,229,0.3)] w-full max-w-md border border-white/40 ring-1 ring-slate-900/5 flex flex-col relative overflow-hidden">
			<div class="absolute -right-8 -top-8 bg-gradient-to-br from-indigo-500 to-indigo-700 text-white font-black text-[9px] uppercase tracking-[0.3em] px-10 py-4 rotate-45 shadow-lg">Laundry</div>

			<Transition name="slide-fade" mode="out-in">
				<div v-if="step === 1" class="text-center py-2" key="step1">
					<div class="w-20 h-20 bg-gradient-to-tr from-indigo-50 to-white text-indigo-600 rounded-3xl flex items-center justify-center mx-auto mb-6 shadow-inner ring-1 ring-indigo-100/50">
						<svg xmlns="http://www.w3.org/2000/svg" class="w-10 h-10 drop-shadow-sm" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
							<path stroke-linecap="round" stroke-linejoin="round" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
						</svg>
					</div>
					<h2 class="text-2xl font-black bg-clip-text text-transparent bg-gradient-to-r from-slate-900 to-slate-600 uppercase tracking-tight leading-tight mb-3">Sistem Rak Cerdas</h2>

					<p class="text-xs font-bold text-slate-500 leading-relaxed px-2">
						{{ isRecoverySnooze ? 'Apakah toko SUDAH MEMILIKI rak untuk menyimpan baju customer saat selesai?' : 'Apakah toko MEMILIKI rak untuk menyimpan baju customer saat selesai?' }}
					</p>

					<div class="flex gap-4 mt-8">
						<button @click="handleBypass" class="flex-1 bg-white hover:bg-slate-50 text-slate-500 py-3.5 rounded-2xl font-black text-xs uppercase tracking-wider active:scale-95 transition-all border border-slate-200 shadow-sm flex items-center justify-center gap-2 group">
							<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-400 group-hover:text-rose-500 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
								<path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
							</svg>
							Belum Punya
						</button>

						<button v-if="isRecoverySnooze" @click="emit('setup-success', true)" class="flex-1 bg-gradient-to-r from-indigo-600 to-indigo-700 hover:from-indigo-500 hover:to-indigo-600 text-white py-3.5 rounded-2xl font-black text-xs uppercase tracking-wider shadow-lg shadow-indigo-600/20 active:scale-95 transition-all flex items-center justify-center gap-2 group">
							<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 group-hover:translate-x-1 transition-transform" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
								<path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
							</svg>
							Sudah Punya
						</button>

						<button v-else @click="step = 2" class="flex-1 bg-gradient-to-r from-indigo-600 to-indigo-700 hover:from-indigo-500 hover:to-indigo-600 text-white py-3.5 rounded-2xl font-black text-xs uppercase tracking-wider shadow-lg shadow-indigo-600/20 active:scale-95 transition-all flex items-center justify-center gap-2 group">
							<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 group-hover:translate-x-1 transition-transform" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
								<path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
							</svg>
							Punya (Setup)
						</button>
					</div>
				</div>

				<div v-else-if="step === 2" class="flex flex-col" key="step2">
					<h3 class="text-xl font-black text-slate-900 uppercase tracking-tight mb-2">Kapasitas Penyimpanan</h3>
					<p class="text-xs font-bold text-slate-500 mb-6 leading-relaxed">Sistem akan secara presisi meng-generate topologi matrix rak Anda.</p>

					<div class="space-y-5 font-mono">
						<div class="flex flex-col gap-2">
							<label class="text-[10px] font-black text-slate-400 uppercase tracking-widest flex items-center gap-1.5">
								<span class="w-1.5 h-1.5 rounded-full bg-indigo-500"></span>
								Nama Lemari Perdana:
							</label>
							<input v-model="zonaNama" type="text" placeholder="Contoh: RAK KAYU A" class="p-4 bg-slate-50/50 border border-slate-200 rounded-2xl outline-none font-black text-sm uppercase focus:border-indigo-500 focus:bg-white focus:ring-4 focus:ring-indigo-500/10 transition-all text-slate-800 shadow-inner" />
						</div>
						<div class="flex flex-col gap-2">
							<label class="text-[10px] font-black text-slate-400 uppercase tracking-widest flex items-center gap-1.5">
								<span class="w-1.5 h-1.5 rounded-full bg-indigo-500"></span>
								Baris Kebawah (A - Z):
							</label>
							<input v-model.number="jumlahBaris" type="number" min="1" max="26" class="p-4 bg-slate-50/50 border border-slate-200 rounded-2xl outline-none font-black text-base focus:border-indigo-500 focus:bg-white focus:ring-4 focus:ring-indigo-500/10 transition-all text-slate-800 shadow-inner" />
						</div>
						<div class="flex flex-col gap-2">
							<label class="text-[10px] font-black text-slate-400 uppercase tracking-widest flex items-center gap-1.5">
								<span class="w-1.5 h-1.5 rounded-full bg-indigo-500"></span>
								Kolom Kekanan (1 - 50):
							</label>
							<input v-model.number="jumlahKolom" type="number" min="1" max="50" class="p-4 bg-slate-50/50 border border-slate-200 rounded-2xl outline-none font-black text-base focus:border-indigo-500 focus:bg-white focus:ring-4 focus:ring-indigo-500/10 transition-all text-slate-800 shadow-inner" />
						</div>
					</div>

					<div class="bg-gradient-to-br from-indigo-50 to-blue-50 border border-indigo-100 rounded-2xl p-4 mt-6 flex items-center gap-4 shadow-inner">
						<div class="w-10 h-10 rounded-full bg-indigo-100 flex items-center justify-center shrink-0">
							<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
								<path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
							</svg>
						</div>
						<p class="text-[10px] font-bold text-indigo-900/80 uppercase leading-relaxed">
							Total kapasitas aktif:
							<span class="font-black text-sm text-indigo-700 bg-white px-2 py-0.5 rounded shadow-sm mx-1">{{ jumlahBaris * jumlahKolom }} Slot</span>
							<br />
							(A-1 s/d {{ String.fromCharCode(64 + (jumlahBaris || 1)) }}-{{ jumlahKolom }})
						</p>
					</div>

					<div class="flex gap-3 mt-8">
						<button @click="step = 1" :disabled="isSubmitting" class="bg-white border border-slate-200 p-4 rounded-2xl text-slate-400 hover:text-slate-600 hover:bg-slate-50 font-black text-xs uppercase active:scale-95 transition-all flex justify-center items-center shadow-sm">
							<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
								<path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" />
							</svg>
						</button>
						<button @click="handleGenerateRack" :disabled="isSubmitting" class="flex-1 bg-slate-900 hover:bg-slate-800 text-white py-4 rounded-2xl font-black text-xs uppercase tracking-widest shadow-xl shadow-slate-900/20 active:scale-95 transition-all flex justify-center items-center gap-2">
							<svg v-if="isSubmitting" class="animate-spin h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
								<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
								<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
							</svg>
							<svg v-else xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-indigo-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
								<path stroke-linecap="round" stroke-linejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
							</svg>
							{{ isSubmitting ? 'Memproses...' : 'Eksekusi Sistem' }}
						</button>
					</div>
				</div>
			</Transition>
		</div>
	</div>
</template>

<style scoped>
.slide-fade-enter-active {
	transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1);
}
.slide-fade-leave-active {
	transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}
.slide-fade-enter-from {
	transform: translateX(20px);
	opacity: 0;
}
.slide-fade-leave-to {
	transform: translateX(-20px);
	opacity: 0;
}
.animate-fade-in {
	animation: fadeIn 0.5s ease-out;
}
@keyframes fadeIn {
	from {
		opacity: 0;
	}
	to {
		opacity: 1;
	}
}
</style>
