<script setup>
import Swal from 'sweetalert2';
import { onMounted, ref } from 'vue';
import api from '../api.js';

const waToken = ref('');
const isLoading = ref(false);
const isTesting = ref(false);
const testNumber = ref('');
const isDataLoaded = ref(false);

onMounted(async () => {
	try {
		// 🚀 KASTA GLOBAL: Langsung tembak pintu gerbang utama /store/settings terpusat bray!
		const response = await api.get('/store/settings');
		const storeData = response.data.data || response.data;

		if (storeData) {
			waToken.value = storeData.wa_token || '';
		}
	} catch (error) {
		console.error('Gagal memuat pengaturan WhatsApp toko:', error);
	} finally {
		isDataLoaded.value = true;
	}
});

const saveSettings = async () => {
	// 🛡️ SECURITY PATCH: Cegah teks masking samaran ikut ter-post jika tidak diubah oleh user!
	if (waToken.value === 'HAS_TOKEN_HIDDEN_BY_SYSTEM') {
		return Swal.fire({ icon: 'info', title: 'Tidak Ada Perubahan', text: 'Token tersimpan aman dan tidak diubah.', timer: 2000, showConfirmButton: false });
	}

	isLoading.value = true;
	try {
		const formData = new FormData();
		formData.append('wa_token', waToken.value.trim());

		// 🚀 KASTA GLOBAL: Simpan murni terpusat ke PUT /store/settings bray bray bray!
		await api.put('/store/settings', formData);

		Swal.fire({ icon: 'success', title: 'Tersimpan', text: 'Token WhatsApp berhasil diperbarui!', timer: 1500 });

		// Paksa samarkan lagi di interface setelah sukses simpan bray
		waToken.value = 'HAS_TOKEN_HIDDEN_BY_SYSTEM';
	} catch (error) {
		Swal.fire('Gagal!', error.response?.data?.error || 'Terjadi kesalahan saat menyimpan.', 'error');
	} finally {
		isLoading.value = false;
	}
};

const testConnection = async () => {
	if (!waToken.value) return Swal.fire('Oops!', 'Simpan Token WA dulu!', 'warning');
	if (!testNumber.value) return Swal.fire('Oops!', 'Masukkan nomor tujuan!', 'warning');

	isTesting.value = true;
	try {
		// 🚀 KASTA GLOBAL: Tembak rute uji coba terpusat ke POST /store/whatsapp/test bray!
		const response = await api.post('/store/whatsapp/test', {
			target: String(testNumber.value).trim(),
		});

		// Backend Go lu merespon dengan status: "success" / "sukses" bray bray bray
		if (response.data.status === 'success' || response.data.status === 'sukses' || response.data.status === true) {
			Swal.fire('Berhasil!', 'Pesan uji coba terkirim ke WhatsApp Anda!', 'success');
		} else {
			Swal.fire('Gagal', response.data.message || 'Cek kembali validitas token!', 'error');
		}
	} catch (error) {
		Swal.fire('Error', 'Gagal menghubungi server.', 'error');
	} finally {
		isTesting.value = false;
	}
};
</script>

<template>
	<div class="p-4 sm:p-6 lg:p-8 max-w-5xl mx-auto space-y-6 font-sans">
		<div class="flex flex-col md:flex-row md:items-center justify-between gap-4 bg-white p-6 rounded-[24px] shadow-sm border border-slate-100">
			<div>
				<h1 class="text-2xl font-black text-slate-800 tracking-tight">Integrasi WhatsApp Gateway</h1>
				<p class="text-[10px] font-bold text-slate-400 mt-1 uppercase tracking-widest">Kirim Nota & Laporan Otomatis ke Pelanggan</p>
			</div>
			<div class="w-12 h-12 rounded-2xl bg-emerald-50 text-emerald-600 flex items-center justify-center border border-emerald-100 shrink-0">
				<svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
					<path d="M3 21l1.65-3.8a9 9 0 1 1 3.4 2.9L3 21" />
					<path d="M9 10a.5.5 0 0 0 1 0V9a.5.5 0 0 0-1 0v1a5 5 0 0 0 5 5h1a.5.5 0 0 0 0-1h-1a.5.5 0 0 0 0 1" />
				</svg>
			</div>
		</div>

		<div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
			<div class="lg:col-span-2 space-y-6">
				<div class="bg-white p-6 rounded-[24px] shadow-sm border border-slate-100">
					<h2 class="text-xs font-black text-slate-800 uppercase tracking-widest mb-6 border-b border-slate-100 pb-4">Konfigurasi Gateway</h2>

					<form @submit.prevent="saveSettings" class="space-y-6">
						<div class="space-y-2">
							<label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Token API (Provider Fonnte)</label>
							<input
								v-model="waToken"
								:type="waToken === 'HAS_TOKEN_HIDDEN_BY_SYSTEM' ? 'text' : 'password'"
								class="w-full px-5 py-4 bg-slate-50 border-2 border-slate-100 rounded-2xl focus:bg-white focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 outline-none font-bold text-sm text-slate-800 transition-all placeholder:text-slate-300"
								:placeholder="waToken === 'HAS_TOKEN_HIDDEN_BY_SYSTEM' ? '•••••••••••••••••••• (Token Tersimpan Aman)' : 'Masukkan Token API Fonnte Anda...'"
								:disabled="isTesting"
								@focus="waToken === 'HAS_TOKEN_HIDDEN_BY_SYSTEM' ? (waToken = '') : null"
								required />
							<p class="text-[10px] font-bold text-slate-400 mt-2 ml-1">
								Belum punya akun?
								<a href="https://fonnte.com" target="_blank" class="text-indigo-600 hover:underline">Daftar di Fonnte.com</a>
							</p>
						</div>

						<div class="pt-4">
							<button type="submit" :disabled="isLoading || isTesting" class="w-full sm:w-auto px-8 py-4 bg-slate-900 text-white rounded-[20px] font-black text-xs uppercase tracking-[0.2em] shadow-lg shadow-slate-900/20 hover:bg-indigo-600 transition-all active:scale-95 disabled:opacity-50 flex items-center justify-center gap-2">
								<span v-if="!isLoading">Simpan Konfigurasi</span>
								<span v-else>Menyimpan...</span>
								<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round">
									<path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z" />
									<polyline points="17 21 17 13 7 13 7 21" />
									<polyline points="7 3 7 8 15 8" />
								</svg>
							</button>
						</div>
					</form>
				</div>

				<div class="bg-indigo-50 p-6 rounded-[24px] border border-indigo-100 relative overflow-hidden">
					<div class="absolute top-0 right-0 w-32 h-32 bg-indigo-500/10 rounded-full blur-2xl transform translate-x-10 -translate-y-10"></div>
					<h2 class="text-xs font-black text-indigo-900 uppercase tracking-widest mb-4">Uji Coba Pengiriman</h2>

					<div class="flex flex-col sm:flex-row gap-3 relative z-10">
						<input v-model="testNumber" type="number" class="flex-1 px-5 py-4 bg-white border-2 border-indigo-100 rounded-2xl focus:border-indigo-500 outline-none font-bold text-sm text-slate-800 transition-all placeholder:text-slate-300" placeholder="081234567890" :disabled="isTesting" />
						<button @click="testConnection" :disabled="isTesting" class="px-6 py-4 bg-indigo-600 text-white rounded-2xl font-black text-xs uppercase tracking-widest hover:bg-slate-900 transition-all shadow-md active:scale-95 disabled:opacity-50 shrink-0">
							<span v-if="!isTesting">Kirim Test</span>
							<span v-else>Mengirim...</span>
						</button>
					</div>
				</div>
			</div>

			<div class="bg-white p-6 rounded-[24px] shadow-sm border border-slate-100 h-fit">
				<div class="flex items-center gap-3 mb-6">
					<div class="w-8 h-8 rounded-full bg-blue-50 flex items-center justify-center text-blue-600 border border-blue-100">
						<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round">
							<circle cx="12" cy="12" r="10" />
							<path d="M12 16v-4" />
							<path d="M12 8h.01" />
						</svg>
					</div>
					<h3 class="font-black text-slate-800 text-xs uppercase tracking-widest">Cara Integrasi</h3>
				</div>

				<ul class="space-y-4 text-[11px] font-bold text-slate-500">
					<li class="flex gap-3">
						<span class="w-5 h-5 rounded-full bg-slate-100 flex items-center justify-center shrink-0 text-[10px] text-slate-800">1</span>
						<span>Buat akun di platform penyedia WhatsApp API (Fonnte).</span>
					</li>
					<li class="flex gap-3">
						<span class="w-5 h-5 rounded-full bg-slate-100 flex items-center justify-center shrink-0 text-[10px] text-slate-800">2</span>
						<span>Hubungkan nomor WhatsApp Toko Anda dengan memindai kode QR.</span>
					</li>
					<li class="flex gap-3">
						<span class="w-5 h-5 rounded-full bg-slate-100 flex items-center justify-center shrink-0 text-[10px] text-slate-800">3</span>
						<span>Salin "Token API" yang diberikan oleh provider.</span>
					</li>
					<li class="flex gap-3">
						<span class="w-5 h-5 rounded-full bg-slate-100 flex items-center justify-center shrink-0 text-[10px] text-slate-800">4</span>
						<span>Tempelkan token tersebut pada form di samping dan simpan bray bray bray!</span>
					</li>
				</ul>
			</div>
		</div>
	</div>
</template>

<style scoped>
input[type='number']::-webkit-inner-spin-button,
input[type='number']::-webkit-outer-spin-button {
	-webkit-appearance: none;
	margin: 0;
}
input[type='number'] {
	-moz-appearance: textfield;
}
</style>
