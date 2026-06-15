import { Html5Qrcode } from 'html5-qrcode';
import Swal from 'sweetalert2';
import { computed, nextTick, onMounted, onUnmounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import api from '../../../api.js';
import { posService } from '../services/posService.js';

export function usePos() {
	const router = useRouter();

	// --- SETUP USER & ROLE SECARA AMAN ---
	const getUserInfo = () => {
		const token = localStorage.getItem('token');
		const role = localStorage.getItem('role') || 'kasir';
		let name = localStorage.getItem('name');

		if (token) {
			try {
				const payload = JSON.parse(atob(token.split('.')[1]));
				if (!name || name === 'undefined' || name === '') {
					name = payload.name || payload.username || 'Kasir Toko';
				}
				return { userId: payload.user_id, role, name };
			} catch (e) {
				return { userId: 0, role, name: 'Kasir Toko' };
			}
		}
		return { userId: 0, role, name: 'Kasir Toko' };
	};

	const currentUser = ref(getUserInfo());
	const currentSession = ref(null);
	const currentTime = ref('');
	let timer;

	const products = ref([]);
	const isLoadingProducts = ref(true);
	const cart = ref([]);
	const heldOrders = ref([]);
	const showHeldModal = ref(false);
	const payAmount = ref(0);
	const paymentMethod = ref('Cash');
	const showReceipt = ref(false);
	const showQrisModal = ref(false);
	const lastTransaction = ref(null);
	const showReceiptClosing = ref(false);
	const lastClosingData = ref(null);
	const noHpPelanggan = ref('');
	const isMobileCartOpen = ref(false);
	const searchQuery = ref('');
	const searchInput = ref(null);
	const selectedCategory = ref('');
	const categories = computed(() => {
		const allCats = products.value.map((p) => p.kategori).filter((cat) => cat && cat.trim() !== '');

		return [...new Set(allCats)];
	});
	const showScanner = ref(false);
	let html5QrCode = null;

	const storeSettings = ref({ payment_type: 'qris_static', qris_image: '', qris_name: '', is_tax_active: false, pajak_percent: 0 });

	// --- UTILS: GATEWAY RESOLVER LINK GAMBAR CLOUD ---
	const getImageUrl = (path) => {
		if (!path || path === 'null' || path === 'undefined') return null;
		if (path.startsWith('http://') || path.startsWith('https://')) return path;
		const cleanPath = path.startsWith('/') ? path : `/${path}`;
		return `${import.meta.env.VITE_API_BASE_URL || 'http://192.168.18.2:8080'}${cleanPath}`;
	};

	// --- LOGIKA KAMERA SCANNER BARCODE HP ---
	const startScanner = async () => {
		showScanner.value = true;
		setTimeout(async () => {
			try {
				html5QrCode = new Html5Qrcode('reader-kasir');
				await html5QrCode.start(
					{ facingMode: 'environment' },
					{ fps: 15, qrbox: { width: 250, height: 100 } },
					(decodedText) => {
						searchQuery.value = decodedText;
						stopScanner();
						handleBarcodeScan();
						const audio = new Audio('https://www.soundjay.com/buttons/beep-07a.mp3');
						audio.play().catch(() => {});
					},
					() => {}
				);
			} catch (err) {
				console.error(err);
				Swal.fire('Camera Error', 'Gagal mengakses kamera internal. Pastikan izin kamera aktif .', 'error');
				stopScanner();
			}
		}, 300);
	};

	const stopScanner = () => {
		if (html5QrCode) {
			html5QrCode
				.stop()
				.then(() => {
					html5QrCode.clear();
					showScanner.value = false;
				})
				.catch(() => {
					showScanner.value = false;
				});
		} else {
			showScanner.value = false;
		}
	};

	// --- LOGIKA MASTER DATA PRODUK RE-LOAD ---
	const fetchProducts = async () => {
		try {
			const response = await posService.getProducts();
			products.value = (response.data || response).map((p) => ({
				id: p.id,
				sku: p.sku || `SKU-${p.id}`,
				name: p.nama_produk || p.name,
				price: p.harga_jual || p.price,
				stock: p.stok || p.stock,
				image: p.gambar || p.image,
				satuan_dasar: p.satuan_dasar || 'PCS',
				satuan_besar: p.satuan_besar || null,
				isi_per_besar: p.isi_per_besar || 0,
				harga_jual_besar: p.harga_jual_besar || 0,
				is_nested_uom: p.is_nested_uom || false,
				satuan_tengah: p.satuan_tengah || null,
				isi_tengah_ke_dasar: p.isi_tengah_ke_dasar || 0,
				harga_jual_tengah: p.harga_jual_tengah || p.harga_jual * (p.isi_tengah_ke_dasar || 1),

				// 🚀 SUNTIKAN WAJIB BRAY: Tampung field kategori dari database Go lu!
				kategori: p.kategori || p.category || 'General',
			}));
		} catch (error) {
			console.error('Gagal memuat katalog produk:', error);
		} finally {
			isLoadingProducts.value = false;
		}
	};

	const filteredProducts = computed(() => {
		return products.value.filter((p) => {
			// 1. Filter via input Ketikan / Scan Barcode SKU bray
			const query = searchQuery.value ? searchQuery.value.toLowerCase() : '';
			const matchSearch = p.name.toLowerCase().includes(query) || (p.sku && p.sku.toLowerCase().includes(query));

			// 2. Filter via Klik Tab Kategori Vertikal/Horizontal di kiri layar
			const matchCategory = !selectedCategory.value || p.kategori === selectedCategory.value;

			// Produk lolos kurasi kalau memenuhi kedua syarat di atas bray bray bray!
			return matchSearch && matchCategory;
		});
	});

	const handleBarcodeScan = () => {
		if (!searchQuery.value) return;
		const query = String(searchQuery.value).trim().toLowerCase();
		const exactMatch = products.value.find((p) => p.sku && String(p.sku).toLowerCase() === query);

		if (exactMatch) {
			addToCart(exactMatch);
			searchQuery.value = '';
		} else if (filteredProducts.value.length === 1) {
			addToCart(filteredProducts.value[0]);
			searchQuery.value = '';
		}
		nextTick(() => {
			if (searchInput.value) searchInput.value.focus();
		});
	};

	// --- CORE ENGINE MANIPULASI KERANJANG KASIR ---
	const addToCart = (product) => {
		// 🛡️ FIX MULTIPLICITY: Paksa parse stok ke angka murni untuk mencegah String "0" tembus!
		if (Number(product.stock) <= 0) {
			return Swal.fire({
				icon: 'error',
				title: 'Stok Kosong!',
				text: `Stok item ${product.name} di laci gudang benar-benar telah habis, beb.`,
			});
		}

		const defaultUom = product.satuan_dasar;
		const existingItem = cart.value.find((item) => item.id === product.id && item.selected_uom === defaultUom);

		if (existingItem) {
			if ((existingItem.qty + 1) * existingItem.uom_multiplier <= Number(product.stock)) {
				existingItem.qty++;
			} else {
				return Swal.fire({ icon: 'warning', title: 'Stok Terbatas', text: 'Kuantitas order melebihi batas sisa stok fisik!' });
			}
		} else {
			cart.value.unshift({
				id: product.id,
				name: product.name,
				price: parseInt(product.price, 10),
				qty: 1,
				selected_uom: defaultUom,
				uom_multiplier: 1,
				has_grosir: !!product.satuan_besar,
				satuan_dasar: product.satuan_dasar,
				harga_dasar: parseInt(product.price, 10),
				is_nested: product.is_nested_uom,
				satuan_tengah: product.satuan_tengah,
				isi_tengah: product.isi_tengah_ke_dasar,
				harga_tengah: parseInt(product.harga_jual_tengah, 10),
				satuan_besar: product.satuan_besar,
				isi_besar: product.isi_per_besar,
				harga_besar: parseInt(product.harga_jual_besar, 10),
			});
		}

		if (window.innerWidth < 1024 && !isMobileCartOpen.value) {
			Swal.fire({ toast: true, position: 'top', icon: 'success', title: `${product.name} ditambahkan`, showConfirmButton: false, timer: 800, timerProgressBar: true });
		}
	};

	const toggleUom = (item) => {
		if (!item.has_grosir) return;
		const prodMaster = products.value.find((p) => p.id === item.id);
		if (!prodMaster) return;

		if (item.is_nested) {
			if (item.selected_uom === item.satuan_dasar) {
				if (item.qty * item.isi_tengah > prodMaster.stock) return Swal.fire({ icon: 'error', title: 'Stok Gagal Konversi', text: `Sisa stok gudang tidak cukup untuk dikonversi ke satuan ${item.satuan_tengah}` });
				item.selected_uom = item.satuan_tengah;
				item.uom_multiplier = item.isi_tengah;
				item.price = item.harga_tengah;
			} else if (item.selected_uom === item.satuan_tengah) {
				if (item.qty * item.isi_besar > prodMaster.stock) {
					Swal.fire({ icon: 'error', title: 'Stok Gagal Konversi', text: `Sisa stok tidak cukup untuk dikonversi ke satuan grosir ${item.satuan_besar}` });
					item.selected_uom = item.satuan_dasar;
					item.uom_multiplier = 1;
					item.price = item.harga_dasar;
					return;
				}
				item.selected_uom = item.satuan_besar;
				item.uom_multiplier = item.isi_besar;
				item.price = item.harga_besar;
			} else {
				item.selected_uom = item.satuan_dasar;
				item.uom_multiplier = 1;
				item.price = item.harga_dasar;
			}
		} else {
			if (item.selected_uom === item.satuan_dasar) {
				if (item.qty * item.isi_besar > prodMaster.stock) return Swal.fire({ icon: 'error', title: 'Stok Gagal Konversi', text: `Sisa stok tidak cukup untuk dikonversi ke satuan ${item.satuan_besar}` });
				item.selected_uom = item.satuan_besar;
				item.uom_multiplier = item.isi_besar;
				item.price = item.harga_besar;
			} else {
				item.selected_uom = item.satuan_dasar;
				item.uom_multiplier = 1;
				item.price = item.harga_dasar;
			}
		}
	};

	const decreaseQty = (item) => {
		const existingItem = cart.value.find((i) => i.id === item.id && i.selected_uom === item.selected_uom);
		if (existingItem) {
			if (existingItem.qty > 1) {
				existingItem.qty--;
			} else {
				cart.value = cart.value.filter((i) => !(i.id === item.id && i.selected_uom === item.selected_uom));
				if (cart.value.length === 0) isMobileCartOpen.value = false;
			}
		}
	};

	const increaseQty = (item) => {
		const existingItem = cart.value.find((i) => i.id === item.id && i.selected_uom === item.selected_uom);
		const prodMaster = products.value.find((p) => p.id === item.id);
		if (existingItem && prodMaster) {
			// 🛡️ FIX: Amankan pembanding stok fisik dengan Number()
			if ((existingItem.qty + 1) * existingItem.uom_multiplier <= Number(prodMaster.stock)) {
				existingItem.qty++;
			} else {
				Swal.fire({ icon: 'warning', title: 'Batas Maksimal', text: 'Kuantitas pesanan telah menyentuh batas stok fisik toko!' });
			}
		}
	};

	const validateQty = (item) => {
		const existingItem = cart.value.find((i) => i.id === item.id && i.selected_uom === item.selected_uom);
		if (!existingItem) return;

		// Ubah string input jadi angka float desimal bray
		let parsedQty = parseFloat(existingItem.qty);

		// Jika bukan angka atau ngaco, balikin ke default minimal 1 atau 0.25
		if (isNaN(parsedQty) || parsedQty <= 0) {
			existingItem.qty = ['KG', 'GRAM', 'LITER', 'ML'].includes(item.selected_uom?.toUpperCase()) ? 0.25 : 1;
			return;
		}

		// Set nilai aslinya yang udah bersih berbentuk angka float bray bray
		existingItem.qty = parsedQty;

		const prodMaster = products.value.find((p) => p.id === item.id);

		// Validasi Overstock (Batas limit stok gudang)
		if (prodMaster && existingItem.qty * existingItem.uom_multiplier > Number(prodMaster.stock)) {
			Swal.fire({ icon: 'warning', title: 'Overstock!', text: 'Angka kuantitas otomatis disesuaikan dengan limit sisa stok fisik.' });

			if (['KG', 'GRAM', 'LITER', 'ML'].includes(item.selected_uom?.toUpperCase())) {
				existingItem.qty = Number(prodMaster.stock) / existingItem.uom_multiplier;
			} else {
				existingItem.qty = Math.floor(Number(prodMaster.stock) / existingItem.uom_multiplier) || 1;
			}
		}
	};

	const clearCart = () => {
		if (cart.value.length === 0) return;
		Swal.fire({ title: 'Kosongkan Keranjang?', text: 'Seluruh daftar belanjaan kasir saat ini akan dihapus permanen !', icon: 'warning', showCancelButton: true, confirmButtonColor: '#ef4444', confirmButtonText: 'Ya, Hapus!', cancelButtonText: 'Batal', customClass: { popup: 'rounded-[28px]' } }).then((result) => {
			if (result.isConfirmed) {
				cart.value = [];
				payAmount.value = 0;
				setPaymentMethod('Cash');
				isMobileCartOpen.value = false;
			}
		});
	};

	const holdTransaction = () => {
		if (cart.value.length === 0) return;
		heldOrders.value.push({
			id: Date.now(),
			customer: `Pelanggan ${heldOrders.value.length + 1}`,
			items: [...cart.value],
			time: new Date().toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' }),
			total: totalBelanja.value,
		});
		cart.value = [];
		payAmount.value = 0;
		setPaymentMethod('Cash');
		isMobileCartOpen.value = false;
		Swal.fire({ toast: true, position: 'top-end', icon: 'info', title: 'Transaksi Berhasil Ditunda!', showConfirmButton: false, timer: 1500 });
	};

	const resumeOrder = (order) => {
		if (cart.value.length > 0) {
			Swal.fire({ title: 'Timpa Keranjang?', text: 'Ada transaksi aktif di keranjang saat ini. Tetap muat antrean tertunda?', icon: 'warning', showCancelButton: true, confirmButtonText: 'Ya, Timpa!', cancelButtonText: 'Batal', customClass: { popup: 'rounded-[28px]' } }).then((res) => {
				if (res.isConfirmed) processResume(order);
			});
		} else {
			processResume(order);
		}
	};

	const processResume = (order) => {
		cart.value = [...order.items];
		heldOrders.value = heldOrders.value.filter((h) => h.id !== order.id);
		showHeldModal.value = false;
		if (window.innerWidth < 1024) isMobileCartOpen.value = true;
	};

	// --- MATHEMATICS MATH ENGINE MATANG & AMAN IDR ---
	const subTotalBelanja = computed(() => cart.value.reduce((total, item) => total + parseFloat(item.price) * parseFloat(item.qty), 0));
	const nilaiPajak = computed(() => (!storeSettings.value || !storeSettings.value.is_tax_active ? 0 : Math.round(((storeSettings.value.pajak_percent || storeSettings.value.pajak_persen || 0) / 100) * subTotalBelanja.value)));
	const totalBelanja = computed(() => Math.round((subTotalBelanja.value + nilaiPajak.value) / 100) * 100);
	const kembalian = computed(() => payAmount.value - totalBelanja.value);

	const setPaymentMethod = (method) => {
		paymentMethod.value = method;
		payAmount.value = method !== 'Cash' ? totalBelanja.value : 0;
	};
	const setNominal = (amount) => {
		if (amount === 0) {
			payAmount.value = 0; // Kalau pencet 'C', reset murni jadi 0
		} else {
			payAmount.value += amount; // Tetap akumulatif! Klik 100 + 50 langsung jadi 150rb!
		}
	};

	const formatInputRupiah = (event) => {
		let rawValue = event.target.value.replace(/\D/g, '');
		payAmount.value = rawValue ? parseInt(rawValue, 10) : 0;
		event.target.value = payAmount.value === 0 ? '' : payAmount.value.toLocaleString('id-ID');
	};

	// --- FINAL CHECKOUT REQ DISPATCHER ---
	const isProcessingCheckout = ref(false);

	const executeCheckout = async () => {
		if (isProcessingCheckout.value || cart.value.length === 0) return;
		isProcessingCheckout.value = true;

		const payloadItems = cart.value.map((item) => ({
			product_id: item.id,
			kuantitas: parseFloat(Number(item.qty * item.uom_multiplier).toFixed(2)), // Garansi aman 2 angka di belakang koma
			uom_label: `${item.qty} ${item.selected_uom}`,
			harga_uom: parseInt(item.price, 10),
		}));

		try {
			const response = await posService.checkout({
				session_id: currentSession.value.id,
				items: payloadItems,
				nominal_bayar: Number(payAmount.value),
				metode_bayar: paymentMethod.value,
				no_hp_pelanggan: noHpPelanggan.value ? String(noHpPelanggan.value).replace(/\D/g, '') : '',
			});

			lastTransaction.value = {
				invoice: response.invoice || response.data?.invoice || `INV-${Date.now()}`,
				cart: [...cart.value],
				total: response.tagihan || response.data?.tagihan || totalBelanja.value,
				pay: payAmount.value,
				return: response.counter || response.kembali || response.data?.kembali || kembalian.value,
				method: paymentMethod.value,
				subtotal: subTotalBelanja.value,
				pajak: nilaiPajak.value,
				date: new Date().toLocaleString('id-ID', { year: '2-digit', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' }).replace(/\//g, '.'),
			};

			showQrisModal.value = false;
			isMobileCartOpen.value = false;
			showReceipt.value = true;
			cart.value = [];
			payAmount.value = 0;
			setPaymentMethod('Cash');
			await fetchProducts(); // Sinkronisasi ulang sisa stok terbaru pasca checkout

			nextTick(() => {
				if (searchInput.value) searchInput.value.focus();
			});
		} catch (error) {
			Swal.fire('Gagal Simpan', error.response?.data?.error || 'Koneksi database backend terputus.', 'error');
		} finally {
			isProcessingCheckout.value = false;
		}
	};

	const processCheckout = async () => {
		if (paymentMethod.value === 'Cash' && payAmount.value < totalBelanja.value) {
			return Swal.fire({ icon: 'error', title: 'Uang Kurang!', text: `Nominal kas kurang Rp ${(totalBelanja.value - payAmount.value).toLocaleString('id-ID')}` });
		}

		if (paymentMethod.value === 'QRIS') {
			const tipePayment = storeSettings.value?.payment_type;
			if (tipePayment === 'midtrans') {
				if (typeof window.snap === 'undefined') return Swal.fire({ icon: 'error', title: 'Gateway Error', text: `Script Midtrans Snap tidak ter-load sempurna. Hubungi admin SaaS .` });
				isProcessingCheckout.value = true;

				try {
					const payRes = await api.post('/retail/pos/midtrans-order', { total: totalBelanja.value });
					if (!payRes.data || !payRes.data.token) {
						Swal.fire('Error Backend', 'Server Go gagal menerbitkan Token Transaksi Midtrans.', 'error');
						isProcessingCheckout.value = false;
						return;
					}

					window.snap.pay(payRes.data.token, {
						onSuccess: () => {
							Swal.fire('Sukses', 'Dana QRIS Otomatis Berhasil Diverifikasi Server!', 'success');
							executeCheckout();
						},
						onPending: () => {
							Swal.fire('Menunggu', 'Menunggu pelanggan menyelesaikan transfer di aplikasi banking.', 'info');
							isProcessingCheckout.value = false;
						},
						onError: () => {
							Swal.fire('Gagal', 'Sesi pembayaran QRIS ditolak oleh Bank / Merchant.', 'error');
							isProcessingCheckout.value = false;
						},
						onClose: () => {
							isProcessingCheckout.value = false;
						},
					});
				} catch (error) {
					Swal.fire({ icon: 'error', title: 'Koneksi Putus', text: error.response?.data?.error || 'Endpoint /retail/pos/midtrans-order tidak merespons.' });
					isProcessingCheckout.value = false;
				}
			} else {
				showQrisModal.value = true;
			}
		} else {
			executeCheckout();
		}
	};

	// --- MANAGEMENT SHIFT CLOSING (CASH COUNT CONTROL) ---
	const showClosingModal = ref(false);
	const pecahan = ref({ p100k: 0, p50k: 0, p20k: 0, p10k: 0, p5k: 0, p2k: 0, p1k: 0, p500: 0, p200: 0, p100: 0, p50: 0, p25: 0 });
	const totalUangFisik = computed(() => {
		return pecahan.value.p100k * 100000 + pecahan.value.p50k * 50000 + pecahan.value.p20k * 20000 + pecahan.value.p10k * 10000 + pecahan.value.p5k * 5000 + pecahan.value.p2k * 2000 + pecahan.value.p1k * 1000 + pecahan.value.p500 * 500 + pecahan.value.p200 * 200 + pecahan.value.p100 * 100 + pecahan.value.p50 * 50 + pecahan.value.p25 * 25;
	});

	const handleClosing = async () => {
		try {
			const res = await posService.closeSession(currentSession.value.id, { total_aktual: totalUangFisik.value, pecahan: pecahan.value });
			Swal.fire({ icon: 'success', title: 'Shift Berhasil Ditutup!', text: 'Mencetak dokumen rekapitulasi serah terima kas.', customClass: { popup: 'rounded-[28px]' } });
			lastClosingData.value = res;
			showClosingModal.value = false;
			showReceiptClosing.value = true;
		} catch (error) {
			Swal.fire('Gagal Closing', error.response?.data?.error || 'Terjadi gangguan koneksi server.', 'error');
		}
	};

	const logout = () => {
		Swal.fire({ title: 'Akhiri Shift Kerja?', text: 'Wajib lakukan cash-count uang fisik laci sebelum mengunci sistem kasir !', icon: 'question', showCancelButton: true, confirmButtonColor: '#2563eb', confirmButtonText: 'Ya, Hitung Duit', cancelButtonText: 'Batal', customClass: { popup: 'rounded-[28px]' } }).then((result) => {
			if (result.isConfirmed) {
				Object.keys(pecahan.value).forEach((k) => (pecahan.value[k] = 0));
				showClosingModal.value = true;
			}
		});
	};

	// --- LIFECYCLE INITIALIZER ---
	onMounted(async () => {
		const token = localStorage.getItem('token');
		if (!token) return router.push('/login');

		try {
			const setRes = await api.get('/retail/store/settings');
			storeSettings.value = setRes.data.data;

			const res = await posService.checkSession();
			if (!res.has_session) return router.push('/retail/pos/buka-kasir');

			currentSession.value = res.session;
			await fetchProducts();

			if (searchInput.value) searchInput.value.focus();
		} catch (error) {
			if (error.response?.status === 401) router.push('/login');
		}

		timer = setInterval(() => {
			currentTime.value = new Date().toLocaleString('id-ID', { day: '2-digit', month: '2-digit', year: 'numeric', hour: '2-digit', minute: '2-digit', second: '2-digit' }).replace(/\//g, '.');
		}, 1000);
	});

	onUnmounted(() => {
		if (showScanner.value) stopScanner();
		clearInterval(timer);
	});

	return {
		currentUser,
		currentSession,
		currentTime,
		products,
		isLoadingProducts,
		cart,
		heldOrders,
		showHeldModal,
		payAmount,
		paymentMethod,
		showReceipt,
		showQrisModal,
		lastTransaction,
		showReceiptClosing,
		lastClosingData,
		isMobileCartOpen,
		searchQuery,
		searchInput,
		showScanner,
		selectedCategory,
		categories,
		pecahan,
		totalUangFisik,
		filteredProducts,
		subTotalBelanja,
		nilaiPajak,
		totalBelanja,
		kembalian,
		isProcessingCheckout,
		showClosingModal,
		noHpPelanggan,
		storeSettings,
		getImageUrl,
		startScanner,
		stopScanner,
		handleBarcodeScan,
		addToCart,
		decreaseQty,
		increaseQty,
		validateQty,
		clearCart,
		holdTransaction,
		resumeOrder,
		setPaymentMethod,
		executeCheckout,
		formatInputRupiah,
		processCheckout,
		handleClosing,
		logout,
		toggleUom,
		setNominal,
	};
}
