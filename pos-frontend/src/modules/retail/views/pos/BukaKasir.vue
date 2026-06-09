<script setup>
import Swal from 'sweetalert2';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import api from '../../../../api.js';
import { useBukaKasir } from '../../composables/useBukaKasir.js';

const router = useRouter();
const { role, name, storeName, stationNumber, modalAwal, loading, checkExistingSession, handleInputModal, handleBukaKasir } = useBukaKasir();

// STATE UTAMA MODAL & SINKRONISASI KUOTA TERMINAL
const showUpgradeModal = ref(false);
const quotaTerminal = ref(1);
const loadingPayment = ref(false);

// WRAPPER FUNGSI BUKA KASIR DENGAN VALIDASI ENKAPSULASI
const submitBukaKasir = async () => {
	if (!stationNumber.value) {
		return Swal.fire({
			icon: 'warning',
			title: 'Stasiun Belum Dipilih',
			text: 'Silakan pilih nomor pos stasiun kasir terlebih dahulu.',
			confirmButtonColor: '#3b82f6',
			customClass: { popup: 'rounded-[32px]' },
		});
	}

	try {
		await handleBukaKasir();
	} catch (error) {
		if (error.response && error.response.status === 403 && error.response.data.error_code === 'QUOTA_FULL') {
			showUpgradeModal.value = true;
		} else {
			Swal.fire({
				icon: 'error',
				title: 'Gagal Membuka Sesi',
				text: error.response?.data?.error || 'Terjadi kesalahan sistem.',
				confirmButtonColor: '#3b82f6',
				customClass: { popup: 'rounded-[32px]' },
			});
		}
	}
};

// ALUR AMAN MIDTRANS SNAP PEMBELIAN TERMINAL BARU
const beliLisensiTambahan = async () => {
	loadingPayment.value = true;
	try {
		const response = await api.post('/retail/subscription/upgrade', {
			plan_name: 'Terminal Tambahan',
		});

		const snapToken = response.data.token;

		window.snap.pay(snapToken, {
			onSuccess: function (result) {
				Swal.fire({
					icon: 'success',
					title: 'Pembayaran Berhasil!',
					text: 'Kuota stasiun mesin kasir Anda otomatis telah bertambah.',
					confirmButtonColor: '#4f46e5',
					customClass: { popup: 'rounded-[32px]' },
				}).then(() => {
					showUpgradeModal.value = false;
					loadingPayment.value = false;
					handleBukaKasir();
				});
			},
			onPending: function (result) {
				Swal.fire({
					icon: 'info',
					title: 'Menunggu Pembayaran',
					text: 'Silakan selesaikan tagihan Anda pada aplikasi e-wallet / banking.',
					confirmButtonColor: '#4f46e5',
					customClass: { popup: 'rounded-[32px]' },
				});
				loadingPayment.value = false;
			},
			onError: function (result) {
				Swal.fire({ icon: 'error', title: 'Pembayaran Gagal', text: 'Sesi transaksi dibatalkan oleh sistem.', confirmButtonColor: '#ef4444', customClass: { popup: 'rounded-[32px]' } });
				loadingPayment.value = false;
			},
			onClose: function () {
				loadingPayment.value = false;
			},
		});
	} catch (error) {
		Swal.fire({
			icon: 'error',
			title: 'Koneksi Terputus',
			text: error.response?.data?.error || 'Gagal menghubungi server pembayaran Midtrans.',
			confirmButtonColor: '#ef4444',
			customClass: { popup: 'rounded-[32px]' },
		});
		loadingPayment.value = false;
	}
};

onMounted(async () => {
	checkExistingSession();

	try {
		const res = await api.get('/retail/store/settings');
		if (res.data && res.data.data) {
			quotaTerminal.value = res.data.data.quota_terminal || 1;
		}
	} catch (error) {
		console.error('Gagal menarik pengaturan toko:', error);
	}

	const script = document.createElement('script');
	const isSandbox = import.meta.env.VITE_MIDTRANS_ENV === 'sandbox';
	script.src = isSandbox ? 'https://app.sandbox.midtrans.com/snap/snap.js' : 'https://app.midtrans.com/snap/snap.js';

	const clientKey = import.meta.env.VITE_MIDTRANS_CLIENT_KEY || 'SB-Mid-client-fallback';
	script.setAttribute('data-client-key', clientKey);
	document.head.appendChild(script);
});
</script>

<template>
	<div class="min-h-screen bg-slate-950 flex items-center justify-center p-4 sm:p-6 relative overflow-hidden font-sans selection:bg-blue-100">
		<div class="absolute top-0 left-0 w-full h-full opacity-10 pointer-events-none">
			<div class="absolute -top-24 -left-24 w-72 sm:w-96 h-72 sm:h-96 bg-blue-600 rounded-full blur-[80px] sm:blur-[120px]"></div>
			<div class="absolute -bottom-24 -right-24 w-72 sm:w-96 h-72 sm:h-96 bg-indigo-600 rounded-full blur-[80px] sm:blur-[120px]"></div>
		</div>

		<div class="w-full max-w-lg relative z-10">
			<div class="bg-white rounded-[36px] sm:rounded-[48px] p-5 sm:p-8 md:p-12 shadow-2xl border-4 sm:border-[12px] border-slate-900/5 relative overflow-hidden">
				<div class="text-center mb-6 sm:mb-10">
					<div class="w-16 h-16 sm:w-20 sm:h-20 bg-slate-900 rounded-[22px] sm:rounded-[28px] flex items-center justify-center mx-auto mb-4 sm:mb-6 shadow-xl shadow-blue-500/10 transform -rotate-3 hover:rotate-0 transition-all duration-500">
						<svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8 sm:w-10 sm:h-10 text-blue-500" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
							<rect width="18" height="12" x="3" y="4" rx="2" ry="2" />
							<line x1="2" x2="22" y1="20" y2="20" />
							<line x1="12" x2="12" y1="16" y2="20" />
						</svg>
					</div>
					<h1 class="text-2xl sm:text-3xl font-black text-slate-900 tracking-tighter uppercase">
						Point of
						<span class="text-blue-600">Sale</span>
					</h1>
					<p class="text-slate-400 font-bold text-[9px] sm:text-[10px] uppercase tracking-[0.3em] sm:tracking-[0.4em] mt-1 break-words px-2">
						{{ storeName }}
					</p>
				</div>

				<div class="space-y-6 sm:space-y-8">
					<div class="bg-slate-50 p-4 sm:p-5 rounded-[24px] sm:rounded-[32px] border border-slate-100 flex items-center justify-between">
						<div class="flex items-center gap-3 sm:gap-4">
							<div class="w-10 h-10 sm:w-12 sm:h-12 rounded-xl sm:rounded-2xl bg-white border border-slate-200 flex items-center justify-center text-base sm:text-lg shadow-sm shrink-0">👤</div>
							<div class="min-w-0">
								<label class="text-[8px] sm:text-[9px] font-black text-slate-400 uppercase tracking-widest block truncate">Logged Operator</label>
								<div class="text-xs sm:text-sm font-black text-slate-800 uppercase flex items-center gap-2 truncate">
									<span class="w-1.5 h-1.5 sm:w-2 sm:h-2 bg-green-500 rounded-full animate-pulse shrink-0"></span>
									<span class="truncate">{{ name }}</span>
								</div>
							</div>
						</div>
					</div>

					<div class="space-y-2.5 sm:space-y-3">
						<label class="text-[9px] sm:text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 italic block">Select Device Station</label>
						<div class="grid grid-cols-3 gap-2 sm:gap-3">
							<button
								v-for="n in Array.from({ length: Math.max(3, quotaTerminal) }, (_, i) => String(i + 1).padStart(2, '0'))"
								:key="n"
								@click="stationNumber = n"
								:class="stationNumber === n ? 'bg-slate-900 text-white shadow-lg shadow-slate-900/20 scale-[1.02] border-slate-900' : 'bg-slate-50 text-slate-400 grayscale border-transparent hover:bg-slate-100'"
								class="flex flex-col items-center gap-1.5 sm:gap-2 py-3 sm:py-4 rounded-[20px] sm:rounded-[28px] border-2 transition-all duration-300 outline-none focus:ring-2 focus:ring-offset-2 focus:ring-slate-900">
								<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 sm:w-5 sm:h-5 transition-transform duration-300 group-hover:scale-110" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
									<rect width="18" height="12" x="3" y="4" rx="2" ry="2" />
									<line x1="2" x2="22" y1="20" y2="20" />
									<line x1="12" x2="12" y1="16" y2="20" />
								</svg>
								<span class="text-[9px] sm:text-[10px] font-black tracking-tighter">POS {{ n }}</span>
							</button>
						</div>
					</div>

					<div class="space-y-2.5 sm:space-y-3">
						<label class="text-[9px] sm:text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 block">Floating Capital (Modal Awal)</label>
						<div class="relative group">
							<div class="absolute left-4 sm:left-6 top-1/2 -translate-y-1/2 flex flex-col items-center select-none pointer-events-none z-10">
								<span class="text-[8px] sm:text-[10px] font-black text-blue-400 uppercase leading-none">IDR</span>
								<span class="text-base sm:text-xl font-black text-blue-600">Rp</span>
							</div>
							<input type="text" placeholder="0" :value="modalAwal === 0 ? '' : modalAwal.toLocaleString('id-ID')" @input="handleInputModal" class="w-full bg-blue-50/30 border-2 border-blue-100 p-5 sm:p-8 pl-14 sm:pl-24 rounded-[24px] sm:rounded-[36px] font-black text-xl sm:text-4xl text-slate-900 focus:border-blue-600 focus:bg-white focus:ring-4 focus:ring-blue-600/5 outline-none transition-all placeholder:text-slate-200" />
						</div>
					</div>

					<div class="pt-2 sm:pt-4">
						<button @click="submitBukaKasir" :disabled="loading || loadingPayment" class="w-full bg-blue-600 hover:bg-slate-900 text-white p-4 sm:p-6 rounded-[24px] sm:rounded-[32px] font-black text-xs sm:text-sm uppercase tracking-[0.15em] sm:tracking-[0.2em] shadow-xl shadow-blue-600/10 hover:shadow-none transition-all duration-300 transform active:scale-[0.98] disabled:opacity-50 flex items-center justify-center gap-3 sm:gap-4 outline-none focus:ring-4 focus:ring-blue-600/20">
							<span>{{ loading || loadingPayment ? 'Accessing Server...' : 'Initialize Session' }}</span>
							<svg v-if="!loading && !loadingPayment" xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 sm:w-5 sm:h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
								<path d="M5 12h14" />
								<path d="m12 5 7 7-7 7" />
							</svg>
						</button>

						<div class="mt-6 sm:mt-8 text-center" v-if="role === 'owner'">
							<button @click="router.push('/retail/dashboard')" class="text-[9px] sm:text-[10px] font-black text-slate-400 hover:text-blue-600 uppercase tracking-[0.25em] sm:tracking-[0.3em] transition-colors flex items-center justify-center gap-2 mx-auto outline-none">
								<svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
									<path d="m15 18-6-6 6-6" />
								</svg>
								Back to Dashboard
							</button>
						</div>
					</div>
				</div>
			</div>
		</div>

		<Transition name="modal-bounce">
			<div v-if="showUpgradeModal" class="fixed inset-0 bg-slate-950/80 backdrop-blur-sm flex items-center justify-center z-[100] px-4 selection:bg-indigo-100">
				<div class="bg-white p-6 sm:p-8 rounded-[32px] sm:rounded-[40px] shadow-2xl max-w-sm w-full text-center border-4 sm:border-8 border-slate-100 relative overflow-hidden cubic-bezier-modal">
					<div class="absolute top-0 left-0 w-full h-2 bg-gradient-to-r from-red-500 to-rose-500"></div>

					<div class="w-16 h-16 sm:w-20 sm:h-20 bg-rose-50 text-rose-500 rounded-full flex items-center justify-center mx-auto mb-4 sm:mb-6">
						<svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8 sm:w-10 sm:h-10" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
							<path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
						</svg>
					</div>

					<h3 class="text-xl sm:text-2xl font-black text-slate-900 tracking-tight">Kuota Penuh!</h3>
					<p class="text-[11px] sm:text-xs font-bold text-slate-500 mt-2.5 sm:mt-3 leading-relaxed px-1">Lisensi toko Anda saat ini hanya mengizinkan {{ quotaTerminal }} Kasir beroperasi sekaligus. Karyawan lain sedang menggunakan laci aktif di stasiun kerja lain.</p>

					<div v-if="role === 'owner'" class="mt-6 sm:mt-8 bg-indigo-50/50 border-2 border-indigo-100 p-4 sm:p-5 rounded-[20px] sm:rounded-[24px]">
						<h4 class="font-black text-indigo-900 text-xs sm:text-sm uppercase tracking-wider">Upgrade System</h4>
						<p class="text-[10px] sm:text-[11px] font-bold text-indigo-600 mt-1.5 leading-relaxed">Buka terminal tambahan sekaligus untuk mempercepat antrean pelanggan retail Anda.</p>
						<button @click="beliLisensiTambahan" :disabled="loadingPayment" class="mt-4 sm:mt-5 w-full bg-indigo-600 hover:bg-indigo-700 text-white py-3 sm:py-4 rounded-xl sm:rounded-2xl font-black text-[10px] sm:text-xs uppercase tracking-widest transition-all shadow-lg shadow-indigo-600/10 transform active:scale-[0.97] disabled:opacity-50">
							{{ loadingPayment ? 'Memproses...' : 'Beli Kuota (Rp 50.000)' }}
						</button>
					</div>

					<div v-else class="mt-6 sm:mt-8 bg-amber-50 border-2 border-amber-100 p-4 sm:p-5 rounded-[20px] sm:rounded-[24px]">
						<h4 class="font-black text-amber-900 text-xs sm:text-sm uppercase tracking-wider">Akses Terbatas</h4>
						<p class="text-[10px] sm:text-[11px] font-bold text-amber-700 mt-1.5">
							Silakan hubungi
							<b>Owner</b>
							untuk melakukan penambahan lisensi station operasional baru.
						</p>
					</div>

					<button @click="showUpgradeModal = false" class="mt-5 sm:mt-6 text-[9px] sm:text-[10px] font-black text-slate-400 hover:text-slate-600 uppercase tracking-widest transition-colors outline-none">Tutup & Batal</button>
				</div>
			</div>
		</Transition>
	</div>
</template>

<style scoped>
/* Transisi Kurva Elastis Premium */
.cubic-bezier-modal {
	transition-timing-function: cubic-bezier(0.34, 1.56, 0.64, 1);
}

/* Animasi Vue Transition Modal (Bounce & Fade Smooth) */
.modal-bounce-enter-from {
	opacity: 0;
}
.modal-bounce-enter-from .cubic-bezier-modal {
	opacity: 0;
	transform: scale(0.9) translateY(20px);
}
.modal-bounce-enter-active,
.modal-bounce-leave-active {
	transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}
.modal-bounce-leave-to {
	opacity: 0;
}
.modal-bounce-leave-to .cubic-bezier-modal {
	opacity: 0;
	transform: scale(0.95) translateY(10px);
}
</style>
