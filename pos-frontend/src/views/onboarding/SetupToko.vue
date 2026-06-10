<script setup>
import Swal from 'sweetalert2';
import { computed, onMounted, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router'; // 🚀 FIX: Import useRouter buat handling redirect pasca resume
import api from '../../api.js';

const route = useRoute();
const router = useRouter(); // 🚀 Definisi router instance
const isExpansion = route.query.is_expansion === 'true';

// 🚀 FONDASI UTAMA: Deklarasikan state penanda resume pembayaran toko pending bray
const isResumingPayment = ref(false);

const pendingIndustry = localStorage.getItem('pendingIndustry') || 'retail';
const pendingPlan = localStorage.getItem('pendingPlan') || 'trial';

const kategoriOptions = [
	{ id: 'Retail', label: 'Retail & Barang', icon: 'M16 11V7a4 4 0 0 0-8 0v4M5 9h14l1 12H4L5 9z' },
	{ id: 'F&B', label: 'Food & Beverage', icon: 'M18 8h1a4 4 0 0 1 0 8h-1M2 8h16v9a4 4 0 0 1-4 4H6a4 4 0 0 1-4-4V8z M6 1v3 M10 1v3 M14 1v3' },
	{ id: 'Jasa', label: 'Layanan & Jasa', icon: 'M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z M3.27 6.96L12 12.01l8.73-5.05 M12 22.08V12' },
	{ id: 'Lainnya', label: 'Bisnis Lainnya', icon: 'M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z M22 6l-10 7L2 6' },
];

const detailOptions = {
	Retail: ['Minimarket / Toko Kelontong', 'Toko Pakaian / Butik', 'Elektronik / Gadget', 'Apotek / Farmasi', 'Toko Bangunan', 'Pet Shop', 'Lainnya'],
	'F&B': ['Restoran / Rumah Makan', 'Cafe / Coffee Shop', 'Bakery / Toko Roti', 'Food Court / Kaki Lima', 'Katering', 'Lainnya'],
	Jasa: ['Laundry', 'Barbershop / Salon', 'Cuci Mobil / Motor', 'Bengkel Otomotif', 'Klinik / Praktek', 'Lainnya'],
	Lainnya: ['Bisnis Umum / Lainnya'],
};

const listProvinsi = ref([]);
const listKota = ref([]);
const listKecamatan = ref([]);
const listKelurahan = ref([]);

const regIds = ref({ provinsi: '', kota: '', kecamatan: '', kelurahan: '' });

const BASE_REGIONAL_API = import.meta.env.VITE_BASE_REGIONAL_API;

const loadProvinsi = async () => {
	try {
		const res = await fetch(`${BASE_REGIONAL_API}/provinsi.json`);
		const data = await res.json();
		listProvinsi.value = data.map((item) => ({ id: item.id, name: item.nama }));
	} catch (e) {
		console.error('Gagal memuat data provinsi via Ibnux API ', e);
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
		form.value.provinsi = listProvinsi.value.find((p) => p.id === newId)?.name || '';
		if (newId) {
			try {
				const res = await fetch(`${BASE_REGIONAL_API}/kabupaten/${newId}.json`);
				const data = await res.json();
				listKota.value = data.map((item) => ({ id: item.id, name: item.nama }));
			} catch (e) {
				console.error(e);
			}
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
		form.value.kota = listKota.value.find((p) => p.id === newId)?.name || '';
		if (newId) {
			try {
				const res = await fetch(`${BASE_REGIONAL_API}/kecamatan/${newId}.json`);
				const data = await res.json();
				listKecamatan.value = data.map((item) => ({ id: item.id, name: item.nama }));
			} catch (e) {
				console.error(e);
			}
		}
	}
);

watch(
	() => regIds.value.kecamatan,
	async (newId) => {
		regIds.value.kelurahan = '';
		listKelurahan.value = [];
		form.value.kecamatan = listKecamatan.value.find((p) => p.id === newId)?.name || '';
		if (newId) {
			try {
				const res = await fetch(`${BASE_REGIONAL_API}/kelurahan/${newId}.json`);
				const data = await res.json();
				listKelurahan.value = data.map((item) => ({ id: item.id, name: item.nama }));
			} catch (e) {
				console.error(e);
			}
		}
	}
);

watch(
	() => regIds.value.kelurahan,
	(newId) => {
		form.value.kelurahan = listKelurahan.value.find((p) => p.id === newId)?.name || '';
	}
);

onMounted(async () => {
	// 1. Injeksi SDK Midtrans bawaan lu bray
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

	// 2. Load region data bawaan lu bray
	loadProvinsi();

	// 3. Set Kategori Bisnis awal dari local storage
	const pendingIndustry = localStorage.getItem('pendingIndustry');
	if (pendingIndustry === 'fnb') form.value.kategori_bisnis = 'F&B';
	else if (pendingIndustry === 'jasa') form.value.kategori_bisnis = 'Jasa';
	else form.value.kategori_bisnis = 'Retail';

	// 🚀 4. RESUME CONTROLLER TRIGGER: Deteksi kedatangan user bawa utang invoice toko bray!
	const resumeStoreId = route.query.resume_store_id;
	if (resumeStoreId) {
		isResumingPayment.value = true; // Gembok layar form utama
		await handleResumePendingStore(resumeStoreId);
	}
});

const handleResumePendingStore = async (storeId) => {
	try {
		Swal.fire({
			title: 'Menghubungkan Server...',
			text: 'Membuka gerbang billing pembayaran, mohon tunggu sebentar bray.',
			allowOutsideClick: false,
			didOpen: () => {
				Swal.showLoading();
			},
		});

		// Tembak endpoint re-trigger payment terpusat lu di Go bray
		const res = await api.post('/re-trigger-payment', { store_id: Number(storeId) });
		Swal.close();

		if (res.data.store_name) {
			form.value.nama_toko = res.data.store_name;
		}

		// 🚀 BOOM: Paksa Modal Snap Midtrans Meledak Keluar di Layar User!
		// 🚀 PERBARUI BLOK SDK SNAP DI DALAM FUNGSI handleResumePendingStore (SETUPTOKO.VUE) LU BRAY!
		window.snap.pay(res.data.snap_token, {
			onSuccess: function (result) {
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

					// Sukses aktif? Paksa hapus cache biar data toko ter-update jadi active bray!
					localStorage.removeItem('temp_stores');
					window.location.href = '/select-store';
				});
			},
			onPending: function (result) {
				Swal.fire('Menunggu Pembayaran', 'Segera selesaikan transaksi invoice Anda bray.', 'info').then(() => {
					window.location.href = '/select-store';
				});
			},
			onError: function (result) {
				Swal.fire('Gagal', 'Sistem pembayaran mendeteksi anomali perbankan.', 'error').then(() => {
					window.location.href = '/select-store';
				});
			},
			onClose: function () {
				// 🚀 FIX ABSOLUT: Begitu disilang saat resume billing, tendang balik ke Select Store!
				Swal.fire({
					title: 'Aktivasi Tertunda',
					text: 'Tenang , konfigurasi toko tersimpan aman. Lu bisa selesaikan pembayaran kapan pun lewat menu Pilih Toko.',
					icon: 'warning',
					confirmButtonColor: '#4f46e5',
					customClass: { popup: 'rounded-[32px]' },
				}).then(() => {
					window.location.href = '/select-store'; // ◄ SEKARANG DI-LOCK KE PILIH CABANG BRAY!
				});
			},
		});
	} catch (error) {
		Swal.close();
		console.error('Gagal memuat billing resume:', error);
		Swal.fire('Error', 'Infrastruktur billing delay. Coba beberapa saat lagi.', 'error');
		isResumingPayment.value = false;
	}
};

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
});

const isLoading = ref(false);

const activeKategori = computed(() => {
	return kategoriOptions.find((k) => k.id === form.value.kategori_bisnis) || kategoriOptions[0];
});

watch(
	() => form.value.kategori_bisnis,
	(newVal) => {
		form.value.detail_bisnis = detailOptions[newVal][0];
	},
	{ immediate: true }
);

const formatNoHp = () => {
	let val = String(form.value.telepon);
	if (val.startsWith('0')) val = val.substring(1);
	if (val.startsWith('62')) val = val.substring(2);
	form.value.telepon = val;
};

const bukaSnapMidtrans = (snapToken) => {
	window.snap.pay(snapToken, {
		onSuccess: function (result) {
			Swal.fire({
				icon: 'success',
				title: 'Pembayaran Berhasil',
				text: 'Infrastruktur premium Anda telah aktif sepenuhnya.',
				timer: 2000,
				showConfirmButton: false,
				allowOutsideClick: false,
				customClass: { popup: 'rounded-[32px]' },
			}).then(() => {
				// HAPUS CACHE LAMA BIAR FRONTEND DIPAKSA RELOAD DATA FRESH DARI DB BRAY!
				localStorage.removeItem('temp_stores');

				localStorage.removeItem('pendingIndustry');
				localStorage.removeItem('pendingPlan');
				window.location.href = '/select-store'; // ◄ Arahkan ke select-store biar dia nge-load ulang ruko aktifnya bray!
			});
		},
		onPending: function (result) {
			Swal.fire('Menunggu Pembayaran', 'Segera selesaikan transaksi Anda sebelum invoice kedaluwarsa.', 'info').then(() => {
				// Pending bayar? Balikkan ke select-store biar ruko pending-nya kelihatan bray!
				window.location.href = '/select-store';
			});
		},
		onError: function (result) {
			Swal.fire('Pembayaran Gagal', 'Terjadi kesalahan sistem perbankan.', 'error').then(() => {
				window.location.href = '/select-store';
			});
		},
		onClose: function () {
			// 🚀 FIX MUTLAK: Begitu disilang olehnya, putar balik haluan ke halaman Pilih Gerai bray!
			Swal.fire({
				title: 'Aktivasi Ditunda',
				text: 'Konfigurasi Toko Anda tersimpan aman. Selesaikan aktivasi kapan pun melalui halaman Pilih Toko.',
				icon: 'warning',
				confirmButtonColor: '#4f46e5',
				customClass: { popup: 'rounded-[32px]' },
			}).then(() => {
				window.location.href = '/select-store'; // ◄ Diarahkan ke select-store!
			});
		},
	});
};

const submit = async () => {
	if (!regIds.value.kelurahan) {
		return Swal.fire('Data Kurang', 'Harap lengkapi pilihan Kelurahan atau Desa terlebih dahulu.', 'warning');
	}
	isLoading.value = true;

	try {
		const finalTipeBisnis = String(form.value.detail_bisnis || form.value.kategori_bisnis).toLowerCase();
		const currentPendingIndustry = localStorage.getItem('pendingIndustry') || 'retail';
		const currentPendingPlan = (localStorage.getItem('pendingPlan') || 'trial').toLowerCase();
		const existingOwnerToken = localStorage.getItem('token');

		const payload = {
			nama_toko: form.value.nama_toko,
			telepon: `62${form.value.telepon}`,
			business_type: finalTipeBisnis,
			industry: currentPendingIndustry,
			plan: currentPendingPlan,
			alamat_toko: form.value.alamat,
			provinsi: form.value.provinsi,
			kota: form.value.kota,
			kecamatan: form.value.kecamatan,
			kelurahan: form.value.kelurahan,
			kode_pos: String(form.value.kode_pos),
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

			// ====================================================================
			// 🚀 INI DIA SUNTIKAN SAKTI NYA BRAY! (Inject Toko Baru ke Cache Lokal)
			// ====================================================================
			// 1. Ambil list ruko lama yang ada di local storage dolo bray
			const oldStoresRaw = localStorage.getItem('temp_stores');
			let currentStores = oldStoresRaw ? JSON.parse(oldStoresRaw) : [];

			// 2. Bikin objek ruko baru tiruan sesuai data yang barusan diinput user bray
			const newStoreObj = {
				id: response.data.store_id,
				nama_toko: form.value.nama_toko,
				industry: currentPendingIndustry,
				subscription_plan: currentPendingPlan,
				subscription_status: response.data.subscription_status || 'pending', // Gembok pending biar merah bray!
				kota: form.value.kota || 'Lokasi Belum Diatur',
			};

			// 3. Masukkan ruko baru ini ke dalam array ruko yang udah ada
			currentStores.push(newStoreObj);

			// 4. Timpa local storage 'temp_stores' pake data paling fresh!
			localStorage.setItem('temp_stores', JSON.stringify(currentStores));
			// ====================================================================
		}

		if (currentPendingPlan === 'basic' || currentPendingPlan === 'pro' || currentPendingPlan === 'premium') {
			isLoading.value = false;

			Swal.fire({
				title: 'Menyiapkan Pembayaran',
				text: 'Menghubungkan ke gerbang aman Midtrans untuk verifikasi lisensi...',
				allowOutsideClick: false,
				didOpen: () => {
					Swal.showLoading();
				},
			});

			try {
				const activeToken = tokenTerupdate || existingOwnerToken;

				if (!activeToken) {
					throw new Error('Sesi token owner tidak ditemukan di sistem. Harap login ulang.');
				}

				const payRes = await api.post('/retail/subscription/upgrade', { plan_name: currentPendingPlan }, { headers: { Authorization: `Bearer ${activeToken}` } });

				const snapToken = payRes.data.token;
				Swal.close();
				bukaSnapMidtrans(snapToken);
			} catch (err) {
				Swal.close();

				if (err.response && err.response.status === 402) {
					const backupSnapToken = err.response.data?.token || err.response.data?.snap_token;
					if (backupSnapToken) {
						bukaSnapMidtrans(backupSnapToken);
						return;
					}
				}

				Swal.fire('Error Gateway', err.response?.data?.error || err.message || 'Gagal memanggil gerbang pembayaran.', 'error').then(() => {
					window.location.href = '/retail/account';
				});
			}
		} else {
			localStorage.removeItem('pendingIndustry');
			localStorage.removeItem('pendingPlan');

			await Swal.fire({
				icon: 'success',
				title: 'Infrastruktur Ready bray!',
				text: 'Selamat menikmati fasilitas Free Trial selama 14 hari.',
				confirmButtonColor: '#4f46e5',
				customClass: { popup: 'rounded-[32px]' },
			});

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
		Swal.fire({
			icon: 'error',
			title: 'Gagal Setup Toko',
			text: error.response?.data?.error || error.message || 'Terjadi kesalahan sistem database.',
			confirmButtonColor: '#ef4444',
		});
	}
	{
		isLoading.value = false;
	}
};
</script>

<template>
	<div class="min-h-screen bg-[#F8FAFC] flex flex-col justify-center py-10 md:py-16 sm:px-6 lg:px-8 font-sans relative overflow-hidden">
		<div class="absolute -top-24 -left-24 w-[30rem] h-[30rem] bg-indigo-200/40 rounded-full blur-3xl pointer-events-none"></div>
		<div class="absolute -bottom-24 -right-24 w-[30rem] h-[30rem] bg-blue-200/40 rounded-full blur-3xl pointer-events-none"></div>

		<!-- 🚀 1. TAMPILAN RESUME PAYMENT (Bila user memuat utang toko pending bray) -->
		<div v-if="isResumingPayment" class="sm:mx-auto sm:w-full sm:max-w-xl px-4 relative z-10 animate-fade-in-up">
			<div class="bg-white p-8 md:p-12 text-center shadow-2xl rounded-[40px] border border-slate-100/50 flex flex-col items-center">
				<div class="w-14 h-14 border-4 border-indigo-100 border-t-indigo-600 rounded-full animate-spin mb-6"></div>
				<h3 class="font-black text-slate-900 uppercase text-base tracking-wider">Menyambungkan Ke Midtrans</h3>
				<p class="text-[11px] font-bold text-slate-400 mt-2 uppercase leading-relaxed max-w-xs mx-auto">Sistem mendeteksi transaksi tertunda. Mengembalikan data invoice gerai Anda secara realtime...</p>
			</div>
		</div>

		<!-- 🚀 2. TAMPILAN FORM UTAMA (Hanya muncul jika kondisi normal pendaftaran baru) -->
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
						<div class="flex flex-col gap-6">
							<div class="flex items-center gap-3 border-b border-slate-100 pb-3">
								<div class="w-8 h-8 rounded-full bg-indigo-50 flex items-center justify-center text-indigo-600 font-black text-xs border border-indigo-100 shadow-sm">1</div>
								<h3 class="text-lg font-black text-slate-800 uppercase tracking-tight">Identitas Bisnis</h3>
							</div>

							<div class="grid grid-cols-1 md:grid-cols-2 gap-5 md:gap-6">
								<div class="md:col-span-2">
									<label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block">Nama Brand / Toko</label>
									<input v-model="form.nama_toko" type="text" required class="input-modern text-lg" placeholder="Contoh: Indomaret, Laundry Bersih, dsb..." />
								</div>

								<div class="md:col-span-2">
									<label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block">Kategori Industri</label>
									<div class="p-4 rounded-[20px] border-2 border-indigo-200 bg-indigo-50/50 flex items-center gap-4 select-none opacity-90">
										<div class="w-12 h-12 bg-white rounded-xl shadow-sm border border-indigo-100 flex items-center justify-center">
											<svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-indigo-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
												<path :d="activeKategori.icon" />
											</svg>
										</div>
										<div class="flex flex-col">
											<span class="font-black text-slate-800 uppercase tracking-widest">{{ activeKategori.label }}</span>
											<span class="text-[10px] font-bold text-indigo-500 uppercase tracking-widest">Pilihan Dari Landing Page</span>
										</div>
										<div class="ml-auto text-slate-400" title="Kategori telah dikunci">
											<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
												<rect width="18" height="11" x="3" y="11" rx="2" ry="2" />
												<path d="M7 11V7a5 5 0 0 1 10 0v4" />
											</svg>
										</div>
									</div>
								</div>

								<div class="md:col-span-2 animate-[fadeIn_0.3s_ease-out]">
									<label class="text-[10px] font-black text-indigo-500 uppercase tracking-widest ml-1 mb-2 flex items-center gap-2">
										<svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
											<polyline points="9 18 15 12 9 6" />
										</svg>
										Spesifikasi {{ form.kategori_bisnis }}
									</label>
									<div class="relative">
										<select v-model="form.detail_bisnis" required class="input-modern bg-slate-50/50 cursor-pointer appearance-none text-indigo-900 border-indigo-100">
											<option v-for="opt in detailOptions[form.kategori_bisnis]" :key="opt" :value="opt">
												{{ opt }}
											</option>
										</select>
										<div class="absolute inset-y-0 right-5 flex items-center pointer-events-none">
											<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-indigo-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
												<path d="m6 9 6 6 6-6" />
											</svg>
										</div>
									</div>
								</div>

								<div class="md:col-span-2">
									<label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block">No. WhatsApp Bisnis</label>
									<div class="flex items-center bg-white border-2 border-slate-200 rounded-2xl focus-within:border-indigo-500 focus-within:ring-4 focus-within:ring-indigo-500/10 transition-all shadow-sm overflow-hidden">
										<div class="pl-5 pr-4 py-4 bg-slate-50 border-r border-slate-200 flex items-center justify-center select-none">
											<span class="text-slate-500 font-black text-sm">+62</span>
										</div>
										<input v-model="form.telepon" @input="formatNoHp" type="number" required class="w-full px-4 py-4 bg-transparent outline-none font-black text-slate-800 placeholder:text-slate-300 placeholder:font-bold" placeholder="81234567890" />
									</div>
								</div>
							</div>
						</div>

						<div class="flex flex-col gap-6 pt-4 border-t border-slate-100">
							<div class="flex items-center gap-3 border-b border-slate-100 pb-3">
								<div class="w-8 h-8 rounded-full bg-emerald-50 flex items-center justify-center text-emerald-600 font-black text-xs border border-emerald-100 shadow-sm">2</div>
								<h3 class="text-lg font-black text-slate-800 uppercase tracking-tight">Lokasi Operasional</h3>
							</div>

							<div class="grid grid-cols-1 md:grid-cols-2 gap-5 md:gap-6">
								<div class="md:col-span-2">
									<label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block">Detail Alamat (Jalan, No, RT/RW)</label>
									<textarea v-model="form.alamat" rows="2" required class="input-modern resize-none uppercase" placeholder="Contoh: Jl. Jendral Sudirman Kav 21, RT 01 / RW 02..."></textarea>
								</div>

								<div class="relative">
									<label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block">Provinsi</label>
									<select v-model="regIds.provinsi" required class="input-modern cursor-pointer appearance-none">
										<option value="" disabled selected hidden>Pilih Provinsi...</option>
										<option v-for="prov in listProvinsi" :key="prov.id" :value="prov.id">
											{{ prov.name }}
										</option>
									</select>
									<div class="absolute inset-y-0 right-5 mt-7 flex items-center pointer-events-none">
										<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
											<path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" />
										</svg>
									</div>
								</div>

								<div class="relative">
									<label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block">Kota / Kabupaten</label>
									<select v-model="regIds.kota" :disabled="!regIds.provinsi" required class="input-modern cursor-pointer appearance-none disabled:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-60">
										<option value="" disabled selected hidden>
											{{ regIds.provinsi ? 'Pilih Kota...' : 'Pilih Provinsi Dulu' }}
										</option>
										<option v-for="kota in listKota" :key="kota.id" :value="kota.id">
											{{ kota.name }}
										</option>
									</select>
									<div class="absolute inset-y-0 right-5 mt-7 flex items-center pointer-events-none">
										<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
											<path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" />
										</svg>
									</div>
								</div>

								<div class="relative">
									<label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block">Kecamatan</label>
									<select v-model="regIds.kecamatan" :disabled="!regIds.kota" required class="input-modern cursor-pointer appearance-none disabled:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-60">
										<option value="" disabled selected hidden>
											{{ regIds.kota ? 'Pilih Kecamatan...' : 'Pilih Kota Dulu' }}
										</option>
										<option v-for="kec in listKecamatan" :key="kec.id" :value="kec.id">
											{{ kec.name }}
										</option>
									</select>
									<div class="absolute inset-y-0 right-5 mt-7 flex items-center pointer-events-none">
										<svg xmlns="http://www.w3.org/2000/xl" class="w-4 h-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
											<path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" />
										</svg>
									</div>
								</div>

								<div class="grid grid-cols-2 gap-3">
									<div class="relative">
										<label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block flex items-center gap-1">
											<span class="truncate">Desa / Kel</span>
										</label>
										<select v-model="regIds.kelurahan" :disabled="!regIds.kecamatan" required class="w-full px-4 py-4 bg-white border-2 border-slate-200 rounded-2xl focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 outline-none font-black text-slate-800 transition-all text-sm cursor-pointer appearance-none disabled:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-60">
											<option value="" disabled selected hidden>
												{{ regIds.kecamatan ? 'Pilih Kelurahan...' : 'Pilih Kelurahan' }}
											</option>
											<option v-for="kel in listKelurahan" :key="kel.id" :value="kel.id" class="text-xs">
												{{ kel.name }}
											</option>
										</select>
										<div class="absolute inset-y-0 right-4 mt-7 flex items-center pointer-events-none">
											<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
												<path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" />
											</svg>
										</div>
									</div>

									<div>
										<label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block flex items-center gap-1">Kode Pos</label>
										<input v-model="form.kode_pos" :disabled="!regIds.kelurahan" type="number" required class="w-full px-4 py-4 bg-white border-2 border-slate-200 rounded-2xl focus-within:border-indigo-500 focus-within:ring-4 focus-within:ring-indigo-500/10 outline-none font-black text-slate-800 transition-all text-sm placeholder:text-slate-300 disabled:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-60" :placeholder="regIds.kelurahan ? 'Ketik Kode Pos...' : 'Ketik Kode Pos'" />
									</div>
								</div>
							</div>
						</div>

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

<style scoped>
.input-modern {
	@apply block w-full px-5 py-4 bg-white border-2 border-slate-200 rounded-2xl focus:bg-white focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 outline-none font-black text-slate-800 transition-all placeholder:text-slate-300 placeholder:font-bold shadow-sm;
}
.btn-submit {
	@apply w-full flex items-center justify-center py-5 md:py-6 px-6 rounded-[24px] shadow-2xl shadow-indigo-200/50 text-xs md:text-sm font-black text-white bg-indigo-600 hover:bg-slate-900 transition-all active:scale-95 disabled:opacity-50 uppercase tracking-[0.2em];
}
textarea::-webkit-scrollbar {
	width: 6px;
}
textarea::-webkit-scrollbar-track {
	background: transparent;
}
textarea::-webkit-scrollbar-thumb {
	background: #cbd5e1;
	border-radius: 10px;
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
input[type='number']::-webkit-inner-spin-button,
input[type='number']::-webkit-outer-spin-button {
	-webkit-appearance: none;
	margin: 0;
}
input[type='number'] {
	-moz-appearance: textfield;
}
</style>
