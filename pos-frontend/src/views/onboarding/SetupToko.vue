<script setup>
import Swal from 'sweetalert2';
import { onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import api from '../../api.js';

// 🚀 IMPORT REUSABLE COMPONENTS KASTA TERTINGGI
import BusinessIdentity from './components/BusinessIdentity.vue';
import RegionSelector from './components/RegionSelector.vue';
import ResumePaymentState from './components/ResumePaymentState.vue';

const route = useRoute();
const router = useRouter();

const isResumingPayment = ref(false);
const isLoading = ref(false);

const pendingIndustry = localStorage.getItem('pendingIndustry') || 'retail';
const pendingPlan = localStorage.getItem('pendingPlan') || 'trial';

const form = ref({
	nama_toko: '',
	kategori_bisnis: 'Retail',
	detail_bisnis: '',
	telepon: '',
	alamat: '',
	provinsi: '',
	kota: '',
	kecamatan: '',
	kelurahan: '',
	kode_pos: '',
	latitude: 0,
	longitude: 0,
});

onMounted(async () => {
	// 1. Injeksi SDK Midtrans
	if (!document.getElementById('midtrans-script-owner')) {
		const midtransEnv = import.meta.env.VITE_MIDTRANS_ENV || 'sandbox';
		const clientKey = import.meta.env.VITE_MIDTRANS_CLIENT_KEY;
		const snapUrl = midtransEnv === 'production' ? 'https://app.midtrans.com/snap/snap.js' : 'https://app.sandbox.midtrans.com/snap/snap.js';

		const script = document.createElement('script');
		script.id = 'midtrans-script-owner';
		script.src = snapUrl;
		script.setAttribute('data-client-key', clientKey);
		document.head.appendChild(script);
	}

	// 2. Set Kategori Bisnis awal dari local storage
	if (pendingIndustry === 'fnb') form.value.kategori_bisnis = 'F&B';
	else if (pendingIndustry === 'jasa') form.value.kategori_bisnis = 'Jasa';
	else form.value.kategori_bisnis = 'Retail';

	// 3. RESUME CONTROLLER TRIGGER
	const resumeStoreId = route.query.resume_store_id;
	if (resumeStoreId) {
		isResumingPayment.value = true;
		await handleResumePendingStore(resumeStoreId);
	}
});

const handleResumePendingStore = async (storeId) => {
	try {
		Swal.fire({
			title: 'Menghubungkan Server...',
			text: 'Membuka gerbang billing pembayaran, mohon tunggu sebentar .',
			allowOutsideClick: false,
			didOpen: () => {
				Swal.showLoading();
			},
		});

		const res = await api.post('/re-trigger-payment', { store_id: Number(storeId) });
		Swal.close();

		if (res.data.store_name) form.value.nama_toko = res.data.store_name;

		window.snap.pay(res.data.snap_token, {
			onSuccess: function () {
				Swal.fire({
					icon: 'success',
					title: 'Pembayaran Berhasil',
					text: 'Gerai cabang Anda resmi diaktifkan global!',
					timer: 2000,
					showConfirmButton: false,
					customClass: { popup: 'rounded-[32px]' },
				}).then(() => {
					localStorage.removeItem('pendingIndustry');
					localStorage.removeItem('pendingPlan');
					localStorage.removeItem('temp_stores');
					window.location.href = '/select-store';
				});
			},
			onPending: function () {
				Swal.fire('Menunggu Pembayaran', 'Segera selesaikan transaksi invoice Anda .', 'info').then(() => {
					window.location.href = '/select-store';
				});
			},
			onError: function () {
				Swal.fire('Gagal', 'Sistem pembayaran mendeteksi anomali perbankan.', 'error').then(() => {
					window.location.href = '/select-store';
				});
			},
			onClose: function () {
				Swal.fire({
					title: 'Aktivasi Tertunda',
					text: 'Konfigurasi Toko Anda tersimpan aman. Lu bisa selesaikan pembayaran kapan pun lewat menu Pilih Toko.',
					icon: 'warning',
					confirmButtonColor: '#4f46e5',
					customClass: { popup: 'rounded-[32px]' },
				}).then(() => {
					window.location.href = '/select-store';
				});
			},
		});
	} catch (error) {
		Swal.close();
		isResumingPayment.value = false;
		Swal.fire('Error', 'Infrastruktur billing delay. Coba beberapa saat lagi.', 'error');
	}
};

const bukaSnapMidtrans = (snapToken) => {
	window.snap.pay(snapToken, {
		onSuccess: function () {
			Swal.fire({
				icon: 'success',
				title: 'Pembayaran Berhasil',
				text: 'Infrastruktur premium Anda telah aktif sepenuhnya.',
				timer: 2000,
				showConfirmButton: false,
				allowOutsideClick: false,
				customClass: { popup: 'rounded-[32px]' },
			}).then(() => {
				localStorage.removeItem('temp_stores');
				localStorage.removeItem('pendingIndustry');
				localStorage.removeItem('pendingPlan');
				window.location.href = '/select-store';
			});
		},
		onPending: function () {
			Swal.fire('Menunggu Pembayaran', 'Segera selesaikan transaksi Anda sebelum invoice kedaluwarsa.', 'info').then(() => {
				window.location.href = '/select-store';
			});
		},
		onError: function () {
			Swal.fire('Pembayaran Gagal', 'Terjadi kesalahan sistem perbankan.', 'error').then(() => {
				window.location.href = '/select-store';
			});
		},
		onClose: function () {
			Swal.fire({
				title: 'Aktivasi Ditunda',
				text: 'Konfigurasi Toko Anda tersimpan aman. Selesaikan aktivasi kapan pun melalui halaman Pilih Toko.',
				icon: 'warning',
				confirmButtonColor: '#4f46e5',
				customClass: { popup: 'rounded-[32px]' },
			}).then(() => {
				window.location.href = '/select-store';
			});
		},
	});
};

const submit = async () => {
	if (!form.value.kelurahan) {
		return Swal.fire('Data Kurang', 'Harap lengkapi pilihan Kelurahan atau Desa terlebih dahulu.', 'warning');
	}
	isLoading.value = true;

	try {
		const finalTipeBisnis = String(form.value.detail_bisnis || form.value.kategori_bisnis).toLowerCase();
		const existingOwnerToken = localStorage.getItem('token');

		const payload = {
			nama_toko: form.value.nama_toko,
			telepon: `62${form.value.telepon}`,
			business_type: finalTipeBisnis,
			industry: pendingIndustry,
			plan: pendingPlan,
			alamat_toko: form.value.alamat,
			provinsi: form.value.provinsi,
			kota: form.value.kota,
			kecamatan: form.value.kecamatan,
			kelurahan: form.value.kelurahan,
			kode_pos: String(form.value.kode_pos),
			latitude: parseFloat(form.value.latitude) || 0,
			longitude: parseFloat(form.value.longitude) || 0,
		};

		const response = await api.post('/setup', payload);
		const tokenTerupdate = response.data?.token;

		if (tokenTerupdate) {
			api.defaults.headers.common['Authorization'] = `Bearer ${tokenTerupdate}`;
			localStorage.setItem('token', tokenTerupdate);
			localStorage.setItem('store_id', response.data.store_id);
			localStorage.setItem('storeName', response.data.store_name || 'POS UMKM');
			localStorage.setItem('subscriptionPlan', response.data.subscription_plan || 'basic');
			localStorage.setItem('role', 'owner');

			const oldStoresRaw = localStorage.getItem('temp_stores');
			let currentStores = oldStoresRaw ? JSON.parse(oldStoresRaw) : [];

			currentStores.push({
				id: response.data.store_id,
				nama_toko: form.value.nama_toko,
				industry: pendingIndustry,
				subscription_plan: pendingPlan,
				subscription_status: response.data.subscription_status || 'pending',
				kota: form.value.kota || 'Lokasi Belum Diatur',
			});
			localStorage.setItem('temp_stores', JSON.stringify(currentStores));
		}

		if (['basic', 'pro', 'premium'].includes(pendingPlan)) {
			isLoading.value = false;
			Swal.fire({
				title: 'Menyiapkan Pembayaran',
				text: 'Menghubungkan ke gerbang aman Midtrans...',
				allowOutsideClick: false,
				didOpen: () => {
					Swal.showLoading();
				},
			});

			try {
				const activeToken = tokenTerupdate || existingOwnerToken;
				if (!activeToken) throw new Error('Sesi token owner tidak ditemukan di sistem. Harap login ulang.');

				const payRes = await api.post('/retail/subscription/upgrade', { plan_name: pendingPlan }, { headers: { Authorization: `Bearer ${activeToken}` } });
				Swal.close();
				bukaSnapMidtrans(payRes.data.token);
			} catch (err) {
				Swal.close();
				if (err.response?.status === 402 && (err.response.data?.token || err.response.data?.snap_token)) {
					bukaSnapMidtrans(err.response.data.token || err.response.data.snap_token);
					return;
				}
				Swal.fire('Error Gateway', err.response?.data?.error || err.message, 'error').then(() => {
					window.location.href = '/retail/account';
				});
			}
		} else {
			localStorage.removeItem('pendingIndustry');
			localStorage.removeItem('pendingPlan');
			await Swal.fire({ icon: 'success', title: 'Infrastruktur Ready !', text: 'Selamat menikmati fasilitas Free Trial selama 14 hari.', confirmButtonColor: '#4f46e5', customClass: { popup: 'rounded-[32px]' } });

			const kat = form.value.kategori_bisnis;
			const det = (form.value.detail_bisnis || '').toLowerCase();

			if (kat === 'Retail' || kat === 'Lainnya') window.location.href = '/retail/dashboard';
			else if (kat === 'F&B') window.location.href = '/fnb/dashboard';
			else if (kat === 'Jasa') {
				if (det.includes('laundry')) window.location.href = '/laundry/laporan';
				else if (det.includes('bengkel') || det.includes('otomotif')) window.location.href = '/bengkel/dashboard';
				else if (det.includes('barbershop') || det.includes('salon')) window.location.href = '/salon/dashboard';
				else if (det.includes('cuci')) window.location.href = '/cuci-kendaraan/dashboard';
				else window.location.href = '/retail/dashboard';
			}
		}
	} catch (error) {
		Swal.fire({ icon: 'error', title: 'Gagal Setup Toko', text: error.response?.data?.error || error.message, confirmButtonColor: '#ef4444' });
	} finally {
		isLoading.value = false;
	}
};
</script>

<template>
	<div class="min-h-screen bg-[#F8FAFC] flex flex-col justify-center py-10 md:py-16 sm:px-6 lg:px-8 font-sans relative overflow-hidden">
		<div class="absolute -top-24 -left-24 w-[30rem] h-[30rem] bg-indigo-200/40 rounded-full blur-3xl pointer-events-none"></div>
		<div class="absolute -bottom-24 -right-24 w-[30rem] h-[30rem] bg-blue-200/40 rounded-full blur-3xl pointer-events-none"></div>

		<ResumePaymentState v-if="isResumingPayment" />

		<div v-else class="w-full flex flex-col items-center">
			<div class="sm:mx-auto sm:w-full sm:max-w-2xl text-center relative z-10 px-4">
				<div class="w-20 h-20 bg-gradient-to-br from-indigo-600 to-blue-600 rounded-[24px] flex items-center justify-center mx-auto shadow-2xl shadow-indigo-200 mb-6 transform -rotate-6 transition-transform hover:rotate-0 duration-500 border-4 border-white">
					<svg xmlns="http://www.w3.org/2000/svg" class="w-10 h-10 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
						<path d="M3 9h18v10a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V9Z" />
						<path d="m3 9 2.45-4.9A2 2 0 0 1 7.24 3h9.52a2 2 0 0 1 1.8 1.1L21 9" />
						<path d="M12 3v6" />
					</svg>
				</div>
				<h2 class="text-3xl md:text-4xl font-black text-slate-900 tracking-tighter">
					Setup Infrastruktur
					<span class="text-indigo-600">Bisnis</span>
				</h2>
				<p class="mt-3 text-slate-400 font-bold text-[10px] md:text-xs uppercase tracking-[0.2em]">Konfigurasi Sistem Menyesuaikan Alur Kerja Anda</p>
				<div class="mt-4 inline-flex items-center gap-2 bg-indigo-100 text-indigo-700 px-4 py-2 rounded-full font-black text-[10px] uppercase tracking-widest shadow-sm">Paket Aktif: {{ pendingIndustry }} - {{ pendingPlan }}</div>
			</div>

			<div class="mt-8 sm:mx-auto sm:w-full sm:max-w-3xl px-4 relative z-10 w-full">
				<div class="bg-white/90 backdrop-blur-xl p-6 md:p-10 shadow-2xl shadow-slate-200/50 rounded-[32px] md:rounded-[40px] border border-white">
					<form @submit.prevent="submit" class="flex flex-col gap-10">
						<BusinessIdentity :formData="form" />

						<RegionSelector :formData="form" />

						<div class="pt-6 mt-2 border-t border-slate-100">
							<button type="submit" :disabled="isLoading" class="btn-submit">
								<template v-if="!isLoading">
									Luncurkan Bisnis Sekarang
									<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 ml-2" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
										<path d="M5 12h14" />
										<path d="m12 5 7 7-7 7" />
									</svg>
								</template>
								<template v-else>
									<div class="w-5 h-5 border-4 border-white/30 border-t-white rounded-full animate-spin mr-3"></div>
									MENGKONFIGURASI SISTEM...
								</template>
							</button>
						</div>
					</form>
				</div>
			</div>
		</div>
		<p class="mt-10 mb-6 text-center text-[9px] font-black text-slate-400 uppercase tracking-[0.3em]">ARZURA POS Operations &copy; 2026</p>
	</div>
</template>

<style>
/* CSS global buat form components dipindah ke sini / file css terpisah  */
.input-modern {
	@apply block w-full px-5 py-4 bg-white border-2 border-slate-200 rounded-2xl focus:bg-white focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 outline-none font-black text-slate-800 transition-all placeholder:text-slate-300 placeholder:font-bold shadow-sm;
}
.btn-submit {
	@apply w-full flex items-center justify-center py-5 md:py-6 px-6 rounded-[24px] shadow-2xl shadow-indigo-200/50 text-xs md:text-sm font-black text-white bg-indigo-600 hover:bg-slate-900 transition-all active:scale-95 disabled:opacity-50 uppercase tracking-[0.2em];
}
input[type='number']::-webkit-inner-spin-button,
input[type='number']::-webkit-outer-spin-button {
	-webkit-appearance: none;
	margin: 0;
}
input[type='number'] {
	-moz-appearance: textfield;
}
@keyframes fadeIn {
	from {
		opacity: 0;
		transform: translateY(-5px);
	}
	to {
		opacity: 1;
		transform: translateY(0);
	}
}
.animate-\[fadeIn_0\.3s_ease-out\] {
	animation: fadeIn 0.3s ease-out;
}
</style>
