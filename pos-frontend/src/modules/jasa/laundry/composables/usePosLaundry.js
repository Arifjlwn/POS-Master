import Swal from 'sweetalert2';
import { computed, onMounted, ref, watch } from 'vue';
import api from '../../../../api.js';

export function usePosLaundry() {
	const services = ref([]);
	const perfumes = ref([]);
	const cart = ref([]);
	const isLoading = ref(false);
	const isSubmitting = ref(false);

	const searchQuery = ref('');
	const customerName = ref('');
	const customerPhone = ref('');
	const estimasiSelesai = ref('');

	const customerResults = ref([]);
	const showCustomerDropdown = ref(false);

	const isCameraOpen = ref(false);
	const cameraTarget = ref('');
	const photoData = ref(null);
	const buktiTransferData = ref(null);

	const videoItemRef = ref(null);
	const canvasItemRef = ref(null);
	const videoQrisRef = ref(null);
	const canvasQrisRef = ref(null);
	const cameraStream = ref(null);

	const showQrisModal = ref(false);
	const showPerfumeControlModal = ref(false);
	const qrisStoreUrl = ref('');

	const storeInfo = ref({
		nama_toko: 'LAUNDRY POS',
		alamat: '',
		telepon: '',
		receipt_footer: 'Terima Kasih',
		payment_type: 'qris_static',
		logo_url: '',
	});

	const isCartOpen = ref(false);
	const paymentMethod = ref('TUNAI');
	const mainPaymentGroup = ref('Cash');
	const uangBayar = ref('');
	const printData = ref(null);
	const printerSize = ref('58');

	// 🚀 ENGINE BARU: Buat reset tanggal otomatis ke H+2
	const resetEstimasiSelesai = () => {
		const today = new Date();
		today.setDate(today.getDate() + 2);
		estimasiSelesai.value = today.toISOString().split('T')[0];
	};

	const formattedUangBayar = computed({
		get() {
			if (!uangBayar.value) return '';
			return new Intl.NumberFormat('id-ID').format(uangBayar.value);
		},
		set(newValue) {
			const cleanValue = String(newValue).replace(/\D/g, '');
			uangBayar.value = cleanValue ? parseInt(cleanValue, 10) : '';
		},
	});

	const formatRupiah = (angka) => {
		return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(angka || 0);
	};

	const formatNoHpCustomer = () => {
		let val = String(customerPhone.value);
		if (val.startsWith('0')) val = val.substring(1);
		if (val.startsWith('62')) val = val.substring(2);
		customerPhone.value = val;
	};

	const setPaymentMethod = (method) => {
		paymentMethod.value = method;
		mainPaymentGroup.value = method === 'TUNAI' ? 'Cash' : 'Non-Cash';
	};

	const setNominalCash = (amount) => {
		if (amount === 0) {
			uangBayar.value = '';
			return;
		}
		if (amount === 'BACKSPACE') {
			const currentStr = String(uangBayar.value || '');
			if (currentStr.length > 0) {
				uangBayar.value = currentStr.substring(0, currentStr.length - 1);
			}
			return;
		}
		if (typeof amount === 'string') {
			const cleanValue = amount.replace(/\D/g, '');
			uangBayar.value = cleanValue ? parseInt(cleanValue, 10) : '';
			return;
		}
		uangBayar.value = (uangBayar.value ? parseInt(uangBayar.value, 10) : 0) + amount;
	};

	const bindVideoStreaming = (el) => {
		if (!el) return;
		if (cameraTarget.value === 'ITEM') {
			videoItemRef.value = el;
		} else {
			videoQrisRef.value = el;
		}
	};

	watch(
		[cameraStream, videoItemRef, videoQrisRef, isCameraOpen],
		() => {
			if (!isCameraOpen.value || !cameraStream.value) return;
			const videoEl = cameraTarget.value === 'ITEM' ? videoItemRef.value : videoQrisRef.value;
			if (videoEl && videoEl.srcObject !== cameraStream.value) {
				videoEl.srcObject = cameraStream.value;
				videoEl.play().catch((e) => console.warn('Autoplay terinterupsi:', e));
			}
		},
		{ deep: true }
	);

	const openCamera = async (target) => {
		cameraTarget.value = target;
		isCameraOpen.value = true;
		try {
			cameraStream.value = await navigator.mediaDevices.getUserMedia({
				video: { facingMode: 'environment' },
				audio: false,
			});
		} catch (err) {
			Swal.fire('Oops!', 'Kamera tidak diizinkan atau hardware error.', 'error');
			isCameraOpen.value = false;
		}
	};

	const takePhoto = () => {
		const video = cameraTarget.value === 'ITEM' ? videoItemRef.value : videoQrisRef.value;
		const canvas = cameraTarget.value === 'ITEM' ? canvasItemRef.value : canvasQrisRef.value;

		if (!video || !canvas) {
			return Swal.fire('Gagal', 'Element visual hardware belum siap!', 'warning');
		}

		try {
			const ctx = canvas.getContext('2d');
			canvas.width = video.videoWidth || 640;
			canvas.height = video.videoHeight || 480;
			ctx.drawImage(video, 0, 0, canvas.width, canvas.height);
			const base64Image = canvas.toDataURL('image/jpeg', 0.8);

			if (cameraTarget.value === 'ITEM') photoData.value = base64Image;
			else buktiTransferData.value = base64Image;

			closeCamera();
		} catch (e) {
			console.error(e);
		}
	};

	const closeCamera = () => {
		isCameraOpen.value = false;
		if (cameraStream.value) {
			cameraStream.value.getTracks().forEach((track) => track.stop());
			cameraStream.value = null;
		}
	};

	const fetchServices = async () => {
		isLoading.value = true;
		try {
			const response = await api.get('/laundry/services');
			services.value = response.data.data || response.data || [];
		} catch (error) {
			Swal.fire('Gagal', 'Gagal mengambil katalog.', 'error');
		} finally {
			isLoading.value = false;
		}
	};

	const fetchPerfumes = async () => {
		try {
			const res = await api.get('/laundry/perfumes');
			const savedPerfumeStates = JSON.parse(localStorage.getItem('perfume_states')) || {};

			// 🚀 FIXED: Jaring pengaman! Tarik data dari res.data.data, kalau kosong ambil dari res.data, kalau kosong lagi jadikan array kosong []
			const rawPerfumes = res.data?.data || res.data || [];

			perfumes.value = rawPerfumes.map((p) => {
				return {
					...p,
					status: savedPerfumeStates[p.id] !== undefined ? savedPerfumeStates[p.id] : 'Tersedia',
				};
			});
		} catch (err) {
			console.error('Gagal menarik data parfum:', err);
		}
	};

	const togglePerfumeStatus = (perfume) => {
		perfume.status = perfume.status === 'Tersedia' ? 'Habis' : 'Tersedia';
		const savedPerfumeStates = JSON.parse(localStorage.getItem('perfume_states')) || {};
		savedPerfumeStates[perfume.id] = perfume.status;
		localStorage.setItem('perfume_states', JSON.stringify(savedPerfumeStates));
	};

	const fetchStoreSetting = async () => {
		try {
			const response = await api.get('/store/settings');
			const data = response.data.data || response.data;
			if (data) {
				let alamatLengkap = data.alamat || '';
				if (data.kelurahan || data.kecamatan) alamatLengkap += `\n${data.kelurahan || ''}, ${data.kecamatan || ''}`;
				if (data.kota || data.provinsi) alamatLengkap += `\n${data.kota || ''}, ${data.provinsi || ''} ${data.kode_pos || ''}`;

				storeInfo.value = {
					nama_toko: data.nama_toko || 'LAUNDRY POS',
					alamat: alamatLengkap,
					telepon: data.telepon || '',
					receipt_footer: data.receipt_footer || 'Terima Kasih',
					payment_type: data.payment_type || 'qris_static',
					logo_url: data.logo_url || '',
				};
				if (data.qris_image) qrisStoreUrl.value = data.qris_image;
			}
		} catch (error) {
			console.error(error);
		}
	};

	const searchCustomer = async () => {
		if (customerName.value.length < 2) {
			showCustomerDropdown.value = false;
			customerResults.value = [];
			return;
		}
		try {
			const response = await api.get(`/laundry/customers/search?q=${customerName.value}`);
			customerResults.value = response.data.data || response.data;
			showCustomerDropdown.value = customerResults.value.length > 0;
		} catch (error) {
			console.error(error);
		}
	};

	const selectCustomer = (cust) => {
		customerName.value = cust.nama;
		let phone = String(cust.no_whatsapp);
		if (phone.startsWith('62')) phone = phone.substring(2);
		customerPhone.value = phone;
		showCustomerDropdown.value = false;
	};

	const addToCart = (service) => {
		if (cart.value.findIndex((item) => item.id === service.id) !== -1) {
			Swal.fire({ toast: true, position: 'top-end', icon: 'info', title: 'Sudah di keranjang!', showConfirmButton: false, timer: 1500 });
		} else {
			cart.value.push({
				id: service.id,
				nama_produk: service.nama_produk,
				harga: service.harga_jual,
				berat: 1,
				satuan_dasar: service.satuan_dasar || 'KG',
				nama_parfum: 'Parfum Standar Bawaan',
				harga_parfum: 0,
			});
		}
	};

	const totalTagihan = computed(() => {
		return cart.value.reduce((acc, item) => acc + item.harga * (parseFloat(item.berat) || 0) + item.harga_parfum, 0);
	});

	const kembalian = computed(() => {
		return (Number(uangBayar.value) || 0) - totalTagihan.value;
	});

	const processCheckout = async () => {
		if (cart.value.length === 0) return Swal.fire('Oops!', 'Keranjang kosong.', 'warning');
		if (!customerName.value || !customerPhone.value) return Swal.fire('Oops!', 'Data Pelanggan WAJIB diisi!', 'warning');
		if (paymentMethod.value === 'TUNAI' && kembalian.value < 0) return Swal.fire('Oops!', 'Uang bayar kurang!', 'warning');

		if (paymentMethod.value === 'QRIS') {
			const currentPaymentType = storeInfo.value.payment_type || 'qris_static';

			if (currentPaymentType === 'qris_static' && !buktiTransferData.value) {
				showQrisModal.value = true;
				return;
			}

			if (currentPaymentType === 'midtrans' && !isSubmitting.value) {
				isSubmitting.value = true;
				try {
					// 🚀 SEKARANG SUDAH PUNYA JALUR MANDIRI KASTA TERTINGGI BRAY!
					// Kita ubah dari '/retail/pos/checkout' ke endpoint laundry baru kita:
					const midtransRes = await api.post('/laundry/midtrans-token', {
						total: parseFloat(totalTagihan.value),
					});

					// Nangkep tokennya tetep sama bray
					const snapToken = midtransRes.data?.token || midtransRes.data?.data?.token;

					if (snapToken) {
						isSubmitting.value = false;
						window.snap.pay(snapToken, {
							onSuccess: () => {
								Swal.fire('Lunas!', 'Pembayaran Midtrans berhasil divalidasi.', 'success');
								executeFinalLaundryCheckout('LUNAS');
							},
							onPending: () => Swal.fire('Tertunda', 'Menunggu penyelesaian pembayaran.', 'info'),
							onError: () => Swal.fire('Gagal!', 'Transaksi Payment Gateway dibatalkan.', 'error'),
						});
						return;
					}
				} catch (midtransErr) {
					isSubmitting.value = false;
					console.error('🚨 [FORENSIK MIDTRANS CRASH]:', midtransErr);

					if (midtransErr.response) {
						return Swal.fire('Backend Go Menolak!', midtransErr.response?.data?.error || `HTTP ${midtransErr.response.status}: Cek terminal Go lu bray!`, 'error');
					} else {
						return Swal.fire('Frontend Crash!', 'Script Midtrans Snap belum terload di index.html lu bray!', 'error');
					}
				}
			}
		}

		executeFinalLaundryCheckout(paymentMethod.value === 'PAYLATER' ? 'BELUM_LUNAS' : 'LUNAS');
	};

	const executeFinalLaundryCheckout = async (forcedStatusBayar) => {
		if (isSubmitting.value) return;
		isSubmitting.value = true;

		try {
			const payload = {
				customer_name: customerName.value,
				customer_phone: `62${customerPhone.value}`,
				estimasi_selesai: estimasiSelesai.value,
				items: cart.value.map((item) => ({
					product_id: item.id,
					berat_kg: parseFloat(item.berat),
					harga_per_kg: item.harga,
					nama_parfum: item.nama_parfum,
					harga_parfum: item.harga_parfum,
					sub_total: parseFloat(item.berat) * item.harga + item.harga_parfum,
				})),
				total_amount: totalTagihan.value,
				payment_method: paymentMethod.value,
				payment_status: forcedStatusBayar,
				foto_barang_base64: photoData.value || '',
				bukti_transfer_base64: buktiTransferData.value || '',
			};

			const response = await api.post('/laundry/checkout', payload);

			if (response.data && response.data.status === 'sukses') {
				printData.value = {
					toko_nama: storeInfo.value.nama_toko || 'LAUNDRY POS',
					toko_alamat: storeInfo.value.alamat,
					toko_telepon: storeInfo.value.telepon,
					toko_footer: storeInfo.value.receipt_footer,
					logo_url: storeInfo.value.logo_url,
					invoice: response.data.invoice_code,
					kasir: 'Admin',
					tanggal: new Date().toISOString(),
					pelanggan: customerName.value,
					pelanggan_phone: customerPhone.value,
					estimasi: estimasiSelesai.value,
					items: cart.value.map((item) => ({
						...item,
						nomor_rak: response.data.nomor_rak,
					})),
					total: totalTagihan.value,
					metode: paymentMethod.value,
					bayar: paymentMethod.value === 'TUNAI' ? uangBayar.value || totalTagihan.value : totalTagihan.value,
					kembali: kembalian.value > 0 ? kembalian.value : 0,
				};

				// 🚀 Kosongin cart dan form
				cart.value = [];
				customerName.value = '';
				customerPhone.value = '';
				uangBayar.value = '';
				paymentMethod.value = 'TUNAI';
				mainPaymentGroup.value = 'Cash';
				photoData.value = null;
				buktiTransferData.value = null;
				isCartOpen.value = false;

				// 🚀 FIXED: RESET TANGGAL ESTIMASI KEMBALI KE H+2 DEFAULT!
				resetEstimasiSelesai();
			}
		} catch (error) {
			Swal.fire('Gagal!', 'Gagal memproses transaksi final.', 'error');
		} finally {
			isSubmitting.value = false;
		}
	};

	onMounted(() => {
		fetchServices();
		fetchPerfumes();
		fetchStoreSetting();
		resetEstimasiSelesai(); // 🚀 Setup awal pas komponen load

		if (!document.getElementById('midtrans-snap-script')) {
			const script = document.createElement('script');
			script.id = 'midtrans-snap-script';
			script.src = 'https://app.sandbox.midtrans.com/snap/snap.js';
			script.setAttribute('data-client-key', 'PASANG_BEBAS_BRAY_KARENA_YANG_DIPAKAI_TETEP_TOKEN_GO');
			script.async = true;
			document.head.appendChild(script);
		}
	});

	return {
		services,
		perfumes,
		cart,
		isLoading,
		isSubmitting,
		searchQuery,
		customerName,
		customerPhone,
		estimasiSelesai,
		customerResults,
		showCustomerDropdown,
		isCameraOpen,
		cameraTarget,
		photoData,
		buktiTransferData,
		videoItemRef,
		canvasItemRef,
		videoQrisRef,
		canvasQrisRef,
		cameraStream,
		showQrisModal,
		showPerfumeControlModal,
		qrisStoreUrl,
		storeInfo,
		isCartOpen,
		paymentMethod,
		mainPaymentGroup,
		formattedUangBayar,
		printData,
		printerSize,
		formatRupiah,
		formatNoHpCustomer,
		searchCustomer,
		selectCustomer,
		// 🚀 FIXED: Hapus timeout ganda karena udah diurus di PosCartSidebar.vue (handleBlurDropdown)
		closeCustomerDropdown: () => {
			showCustomerDropdown.value = false;
		},
		filteredServices: computed(() => services.value.filter((s) => s.nama_produk.toLowerCase().includes(searchQuery.value.toLowerCase()))),
		availablePerfumes: computed(() => perfumes.value.filter((p) => p.status === 'Tersedia')),
		addToCart,
		handleCartPerfumeChange: (index, event) => {
			const selectedId = event.target.value;
			if (selectedId === 'default') {
				cart.value[index].nama_parfum = 'Parfum Standar Bawaan';
				cart.value[index].harga_parfum = 0;
			} else {
				const pObj = perfumes.value.find((p) => p.id == selectedId);
				if (pObj) {
					cart.value[index].nama_parfum = pObj.nama;
					cart.value[index].harga_parfum = pObj.harga;
				}
			}
		},
		updateBerat: (index, delta) => {
			let item = cart.value[index];
			let current = parseFloat(item.berat) || 0;
			let actualDelta = item.satuan_dasar === 'KG' ? delta : delta > 0 ? 1 : -1;
			let newVal = current + actualDelta;
			item.berat = newVal < (item.satuan_dasar === 'KG' ? 0.1 : 1) ? (item.satuan_dasar === 'KG' ? 0.1 : 1) : item.satuan_dasar === 'KG' ? Math.round(newVal * 10) / 10 : Math.round(newVal);
		},
		removeCartItem: (index) => {
			cart.value.splice(index, 1);
			if (cart.value.length === 0) isCartOpen.value = false;
		},
		clearCart: () =>
			Swal.fire({ title: 'Batalkan Cucian?', icon: 'warning', showCancelButton: true, confirmButtonColor: '#e11d48', confirmButtonText: 'Ya, Bersihkan' }).then((result) => {
				if (result.isConfirmed) {
					cart.value = [];
					isCartOpen.value = false;
					photoData.value = null;
					buktiTransferData.value = null;
				}
			}),
		totalTagihan,
		kembalian,
		openCamera,
		takePhoto,
		closeCamera,
		cancelQris: () => {
			showQrisModal.value = false;
			buktiTransferData.value = null;
			closeCamera();
		},
		confirmQris: () => {
			showQrisModal.value = false;
			processCheckout();
		},
		processCheckout,
		togglePerfumeStatus,
		setPaymentMethod,
		setNominalCash,
		bindVideoStreaming,
	};
}
