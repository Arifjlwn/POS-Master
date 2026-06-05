<script setup>
import { computed, ref } from 'vue';

const props = defineProps({ form: Object, qrisPreview: String });
const emit = defineEmits(['update-file', 'remove-qris']);
const qrisInput = ref(null);

const onQrisSelect = (e) => {
	const file = e.target.files[0];
	if (file) {
		// VALIDASI KEAMANAN: Cegah upload file gajah yang bikin server ngos-ngosan bray!
		if (file.size > 2 * 1024 * 1024) {
			alert('Maksimal ukuran file QRIS adalah 2MB');
			return;
		}
		emit('update-file', 'qris', file, URL.createObjectURL(file));
	}
};

// HELPER UX: Generate URL Webhook dinamis sesuai domain saat ini untuk dicopy user ke Midtrans
const webhookUrl = computed(() => {
	const baseUrl = import.meta.env.VITE_API_BASE_URL || window.location.origin;
	// Bersihkan slash ganda jika ada
	return `${baseUrl}/api/retail/midtrans/webhook`.replace(/([^:]\/)\/+/g, '$1');
});

const copyWebhook = () => {
	navigator.clipboard.writeText(webhookUrl.value);
	alert('URL Webhook berhasil disalin! Silakan paste di Dashboard Midtrans Anda.');
};
</script>

<template>
	<div class="space-y-6 animate-fade-in-up">
		<div class="bg-white p-6 rounded-[24px] border border-slate-200">
			<label class="text-[10px] font-black text-slate-800 uppercase tracking-widest mb-4 block">Metode Pembayaran Digital</label>
			<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
				<div @click="form.payment_type = 'qris_static'" :class="form.payment_type === 'qris_static' ? 'bg-indigo-50 border-indigo-500 shadow-md ring-2 ring-indigo-100' : 'bg-transparent border-slate-200 hover:border-indigo-300'" class="p-4 rounded-2xl border-2 cursor-pointer transition-all flex items-center gap-4">
					<div :class="form.payment_type === 'qris_static' ? 'bg-indigo-500 text-white' : 'bg-slate-200 text-slate-400'" class="w-8 h-8 rounded-full flex items-center justify-center shrink-0 transition-colors">
						<svg v-if="form.payment_type === 'qris_static'" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="4"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" /></svg>
					</div>
					<div>
						<div class="text-sm font-black text-slate-800 uppercase tracking-tight">QRIS Statis</div>
						<div class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-1">Manual Validasi Kasir</div>
					</div>
				</div>

				<div @click="form.payment_type = 'midtrans'" :class="form.payment_type === 'midtrans' ? 'bg-blue-50 border-blue-500 shadow-md ring-2 ring-blue-100' : 'bg-transparent border-slate-200 hover:border-blue-300'" class="p-4 rounded-2xl border-2 cursor-pointer transition-all flex items-center gap-4">
					<div :class="form.payment_type === 'midtrans' ? 'bg-blue-500 text-white' : 'bg-slate-200 text-slate-400'" class="w-8 h-8 rounded-full flex items-center justify-center shrink-0 transition-colors">
						<svg v-if="form.payment_type === 'midtrans'" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="4"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" /></svg>
					</div>
					<div>
						<div class="text-sm font-black text-slate-800 uppercase tracking-tight">Midtrans Gateway</div>
						<div class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-1">Validasi Otomatis via API</div>
					</div>
				</div>
			</div>
		</div>

		<div v-if="form.payment_type === 'qris_static'" class="animate-[fadeInUp_0.2s_ease-out]">
			<div class="bg-indigo-50 p-4 rounded-2xl border border-indigo-100 flex items-start gap-4 mb-6">
				<div class="p-2 bg-indigo-100 text-indigo-600 rounded-xl shrink-0">
					<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
				</div>
				<div>
					<p class="text-[10px] font-bold text-indigo-600 leading-relaxed uppercase tracking-widest mt-1">Pembeli melakukan scan gambar QRIS. Kasir melakukan validasi mutasi secara manual sebelum menyelesaikan transaksi.</p>
				</div>
			</div>

			<div class="flex flex-col md:flex-row gap-6">
				<div class="w-full md:w-1/2">
					<label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Nama Pemilik Rekening (A.N)</label>
					<input v-model="form.qris_name" type="text" placeholder="Misal: Budi Santoso" class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:border-indigo-500 outline-none font-bold text-sm mb-6 transition-colors" />
				</div>

				<div class="w-full md:w-1/2 flex flex-col items-center p-6 bg-slate-50 rounded-3xl border-2 border-slate-100">
					<label class="text-[10px] font-black text-slate-500 uppercase tracking-widest mb-4">Upload Barcode QRIS</label>

					<div class="relative w-48 h-48 bg-white rounded-3xl shadow-sm flex items-center justify-center overflow-hidden border-4 border-slate-200 group hover:border-emerald-400 transition-all">
						<div v-if="qrisPreview" class="w-full h-full relative cursor-default">
							<img :src="qrisPreview" class="w-full h-full object-contain p-2" />
							<button @click.prevent="emit('remove-qris')" class="absolute inset-0 bg-rose-500/80 text-white flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity backdrop-blur-sm cursor-pointer z-10">
								<svg class="w-10 h-10" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
							</button>
						</div>

						<div v-else @click="$refs.qrisInput.click()" class="w-full h-full flex flex-col items-center justify-center cursor-pointer text-slate-300 group-hover:text-emerald-500 transition-colors">
							<svg xmlns="http://www.w3.org/2000/svg" class="w-10 h-10 mx-auto mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm14 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z" /></svg>
							<p class="text-[10px] font-black uppercase tracking-widest">Pilih Gambar</p>
						</div>
					</div>
					<input type="file" ref="qrisInput" @change="onQrisSelect" class="hidden" accept="image/*" />
				</div>
			</div>
		</div>

		<div v-if="form.payment_type === 'midtrans'" class="bg-blue-50/50 p-6 rounded-[24px] border border-blue-100 animate-[fadeInUp_0.2s_ease-out]">
			<div class="flex items-center gap-3 mb-6">
				<div class="p-2 bg-blue-100 text-blue-600 rounded-xl">
					<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" /></svg>
				</div>
				<div>
					<h3 class="text-sm font-black text-blue-900 uppercase tracking-widest">Kredensial API Midtrans</h3>
					<p class="text-[9px] font-bold text-blue-600 uppercase tracking-widest mt-1">Dapatkan kunci ini dari dashboard akun Midtrans Anda.</p>
				</div>
			</div>

			<div class="space-y-4">
				<div>
					<label class="text-[9px] font-black text-blue-500 uppercase tracking-widest mb-2 block">Client Key</label>
					<input v-model="form.midtrans_client_key" type="text" placeholder="SB-Mid-client-xxx..." class="w-full p-4 bg-white border border-blue-200 focus:border-blue-600 rounded-xl outline-none font-bold text-xs text-slate-800 transition-all shadow-sm" />
				</div>
				<div>
					<label class="text-[9px] font-black text-blue-500 uppercase tracking-widest mb-2 block">Server Key (Rahasia)</label>
					<input v-model="form.midtrans_server_key" type="password" placeholder="SB-Mid-server-xxx..." class="w-full p-4 bg-white border border-blue-200 focus:border-blue-600 rounded-xl outline-none font-bold text-xs text-slate-800 transition-all shadow-inner" />
				</div>
			</div>

			<div class="mt-6 pt-6 border-t border-blue-100">
				<label class="text-[9px] font-black text-blue-500 uppercase tracking-widest mb-2 block">URL Webhook Notifikasi (Wajib Diisi di Midtrans)</label>
				<div class="flex items-center gap-2">
					<input :value="webhookUrl" type="text" readonly class="w-full p-3 bg-blue-100/50 border border-blue-200 rounded-xl outline-none font-bold text-[10px] text-blue-800 cursor-not-allowed" />
					<button @click.prevent="copyWebhook" type="button" class="shrink-0 bg-blue-600 hover:bg-blue-700 text-white px-4 py-3 rounded-xl font-black text-[10px] uppercase tracking-widest shadow-md transition-all active:scale-95 flex items-center gap-2">
						<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" /></svg>
						Salin
					</button>
				</div>
				<p class="text-[9px] font-bold text-slate-400 mt-2">Buka Dashboard Midtrans &rarr; Settings &rarr; Configuration &rarr; Payment Notification URL. Lalu paste URL di atas.</p>
			</div>
		</div>
	</div>
</template>
