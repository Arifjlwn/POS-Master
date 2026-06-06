import { Html5Qrcode } from 'html5-qrcode';
import Swal from 'sweetalert2';
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue';
import { productService } from '../services/productService.js';

export function useMasterProduk() {
	const products = ref([]),
		categories = ref([]),
		isLoading = ref(true),
		isSubmitting = ref(false);
	const stok_dalam_karton = ref(null),
		eceran_tambahan = ref(null),
		searchQuery = ref(''),
		selectedCategory = ref('');
	const currentPage = ref(1),
		limitPerPage = ref(10),
		totalPages = ref(1);
	const showFormModal = ref(false),
		isEditing = ref(false),
		editId = ref(null),
		fileInput = ref(null),
		imagePreview = ref(null);
	const showScanner = ref(false);
	let html5QrCode = null,
		searchTimer = null,
		scannerTimeout = null;

	const form = ref({
		name: '',
		sku: '',
		category: '',
		product_type: 'retail',
		estimasi: 'Standar',
		cost_price: 0,
		price: 0,
		stock: 0,
		image: null,
		satuan_dasar: 'PCS',
		has_satuan_besar: false,
		satuan_besar: '',
		isi_per_besar: null,
		harga_beli_besar: null,
		harga_jual_besar: null,
		harga_eceran_tampil: 0,
		qty_eceran_tampil: 1,
		input_kg: null,
		is_nested_uom: false,
		satuan_tengah: '',
		isi_besar_ke_tengah: null,
		isi_tengah_ke_dasar: null,
	});

	const getImageUrl = (p) => (p ? (p.startsWith('http://') || p.startsWith('https://') ? p : `${import.meta.env.VITE_API_BASE_URL}${p}`) : null);

	const fetchCategories = async () => {
		try {
			const res = await productService.getCategories();
			categories.value = res.data.data;
		} catch (e) {}
	};
	const fetchProducts = async (page = 1) => {
		isLoading.value = true;
		currentPage.value = page;
		try {
			const res = await productService.getProducts({ page, limit: limitPerPage.value, search: searchQuery.value, category: selectedCategory.value });
			products.value = res.data.data;
			totalPages.value = Math.ceil(res.data.total_items / limitPerPage.value) || 1;
		} catch (e) {
		} finally {
			isLoading.value = false;
		}
	};

	watch(
		() => form.value.input_kg,
		(v) => {
			if ((form.value.satuan_dasar === 'GRAM' || form.value.satuan_dasar === 'ML') && v > 0) form.value.isi_per_besar = v * 1000;
		}
	);
	watch(
		() => [form.value.isi_besar_ke_tengah, form.value.isi_tengah_ke_dasar, form.value.is_nested_uom],
		([bt, td, isN]) => {
			if (isN && bt > 0 && td > 0) form.value.isi_per_besar = bt * td;
		}
	);
	watch(
		() => [form.value.harga_beli_besar, form.value.isi_per_besar, form.value.has_satuan_besar, stok_dalam_karton.value, eceran_tambahan.value],
		([hb, ipb, hsb, jk, je]) => {
			if (hsb && hb > 0 && ipb > 0) form.value.cost_price = Math.round(hb / ipb);
			if (hsb && ipb > 0) {
				form.value.stock = (parseInt(jk) || 0) * (parseInt(ipb) || 0) + (parseInt(je) || 0);
			} else if (!hsb) {
				stok_dalam_karton.value = null;
				eceran_tambahan.value = null;
			}
		}
	);
	watch(
		() => [form.value.harga_eceran_tampil, form.value.qty_eceran_tampil],
		([h, q]) => {
			if (q > 0 && h >= 0) form.value.price = h / q;
		}
	);
	watch(
		() => form.value.satuan_dasar,
		(nV) => {
			form.value.qty_eceran_tampil = nV === 'GRAM' || nV === 'ML' ? 1000 : 1;
			form.value.harga_eceran_tampil = form.value.price * form.value.qty_eceran_tampil;
		}
	);
	watch([searchQuery, selectedCategory], () => {
		clearTimeout(searchTimer);
		searchTimer = setTimeout(() => {
			fetchProducts(1);
		}, 500);
	});

	const startScanner = async () => {
		showScanner.value = true;
		scannerTimeout = setTimeout(async () => {
			try {
				html5QrCode = new Html5Qrcode('reader');
				await html5QrCode.start(
					{ facingMode: 'environment' },
					{ fps: 10, qrbox: { width: 250, height: 100 } },
					(txt) => {
						form.value.sku = txt;
						stopScanner();
						new Audio('https://www.soundjay.com/buttons/beep-07a.mp3').play().catch(() => {});
					},
					() => {}
				);
			} catch (err) {
				Swal.fire('Oops!', 'Gagal menyalakan kamera scanner.', 'error');
				stopScanner();
			}
		}, 300);
	};

	const stopScanner = () => {
		clearTimeout(scannerTimeout);
		if (html5QrCode && html5QrCode.isScanning) {
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

	const changePage = (p) => {
		if (p >= 1 && p <= totalPages.value) fetchProducts(p);
	};
	const onFileChange = (e) => {
		const file = e.target.files[0];
		if (file) {
			if (imagePreview.value && imagePreview.value.startsWith('blob:')) URL.revokeObjectURL(imagePreview.value); // Clean memory leak
			form.value.image = file;
			imagePreview.value = URL.createObjectURL(file);
		}
	};

	const openAddModal = () => {
		isEditing.value = false;
		editId.value = null;
		imagePreview.value = null;
		stok_dalam_karton.value = null;
		eceran_tambahan.value = null;
		form.value = {
			name: '',
			sku: '',
			category: '',
			// 🚀 TAMBAHAN BARU
			product_type: 'retail',
			estimasi: 'Standar',
			// ---
			cost_price: 0,
			price: 0,
			stock: 0,
			image: null,
			satuan_dasar: 'PCS',
			has_satuan_besar: false,
			satuan_besar: '',
			isi_per_besar: null,
			harga_beli_besar: null,
			harga_jual_besar: null,
			harga_eceran_tampil: 0,
			qty_eceran_tampil: 1,
			input_kg: null,
			is_nested_uom: false,
			satuan_tengah: '',
			isi_besar_ke_tengah: null,
			isi_tengah_ke_dasar: null,
			// 🚀 TAMBAHAN BARU
			harga_jual_tengah: 0,
		};
		showFormModal.value = true;
	};

	const editProduct = (product) => {
		isEditing.value = true;
		editId.value = product.public_id;
		const isBulk = product.satuan_dasar === 'GRAM' || product.satuan_dasar === 'ML',
			acuan = isBulk ? 1000 : 1;

		form.value = {
			name: product.nama_produk || '',
			sku: product.sku || '',
			category: product.kategori || '',
			// 🚀 SINKRONISASI FIELD BARU DARI DATABASE
			product_type: product.product_type || 'retail',
			estimasi: product.estimasi || 'Standar',
			// ---
			cost_price: product.harga_modal || 0,
			price: product.harga_jual || 0,
			stock: product.stok || 0,
			image: null,
			satuan_dasar: product.satuan_dasar || 'PCS',
			has_satuan_besar: !!product.satuan_besar,
			satuan_besar: product.satuan_besar || '',
			isi_per_besar: product.isi_per_besar || null,
			harga_beli_besar: product.satuan_besar ? (product.harga_modal || 0) * (product.isi_per_besar || 1) : null,
			harga_jual_besar: product.harga_jual_besar || null,
			qty_eceran_tampil: acuan,
			harga_eceran_tampil: (product.harga_jual || 0) * acuan,
			is_nested_uom: [true, 'true', 1].includes(product.is_nested_uom),
			satuan_tengah: product.satuan_tengah || '',
			isi_besar_ke_tengah: product.isi_besar_ke_tengah || null,
			isi_tengah_ke_dasar: product.isi_tengah_ke_dasar || null,
			// 🚀 SINKRONISASI HARGA JUAL TENGAH
			harga_jual_tengah: product.harga_jual_tengah || 0,
			input_kg: isBulk && product.isi_per_besar ? product.isi_per_besar / 1000 : null,
		};

		if (product.satuan_besar && product.isi_per_besar > 0) {
			stok_dalam_karton.value = Math.floor((product.stok || 0) / product.isi_per_besar);
			eceran_tambahan.value = (product.stok || 0) % product.isi_per_besar;
		} else {
			stok_dalam_karton.value = null;
			eceran_tambahan.value = null;
		}
		imagePreview.value = getImageUrl(product.gambar);
		showFormModal.value = true;
	};

	const submitProduct = async () => {
		if (form.value.has_satuan_besar && (!form.value.satuan_besar || !form.value.isi_per_besar || !form.value.harga_jual_besar)) {
			return Swal.fire('Data Kurang!', 'Lengkapi detail satuan besar!', 'warning');
		}
		isSubmitting.value = true;
		const fData = new FormData();
		fData.append('nama_produk', form.value.name);
		// Pastikan SKU gak "undefined"
		fData.append('sku', form.value.sku || '');
		fData.append('kategori', form.value.category || '');

		// 🚀 INJECT FIELD BARU KE FORM DATA
		fData.append('product_type', form.value.product_type || 'retail');
		fData.append('estimasi', form.value.estimasi || 'Standar');

		fData.append('harga_modal', Number(form.value.cost_price));
		fData.append('harga_jual', Number(form.value.price));
		fData.append('satuan_dasar', form.value.satuan_dasar);

		fData.append('satuan_besar', form.value.has_satuan_besar ? form.value.satuan_besar : '');
		fData.append('isi_per_besar', form.value.has_satuan_besar ? Number(form.value.isi_per_besar) : 0);
		fData.append('harga_jual_besar', form.value.has_satuan_besar ? Number(form.value.harga_jual_besar) : 0);

		fData.append('is_nested_uom', form.value.is_nested_uom);
		if (form.value.is_nested_uom) {
			fData.append('satuan_tengah', form.value.satuan_tengah);
			fData.append('isi_besar_ke_tengah', Number(form.value.isi_besar_ke_tengah));
			fData.append('isi_tengah_ke_dasar', Number(form.value.isi_tengah_ke_dasar));
			// 🚀 JANGAN LUPA HARGA JUAL TENGAH
			fData.append('harga_jual_tengah', Number(form.value.harga_jual_tengah));
		}

		if (form.value.image) fData.append('gambar', form.value.image);
		if (!isEditing.value) fData.append('stok', Number(form.value.stock));

		try {
			if (isEditing.value) {
				await productService.updateProduct(editId.value, fData);
			} else {
				await productService.createProduct(fData);
			}
			Swal.fire({ icon: 'success', title: 'Berhasil!', text: 'Data tersimpan!', timer: 1500, showConfirmButton: false });
			showFormModal.value = false;
			fetchProducts(currentPage.value);
			fetchCategories();
		} catch (err) {
			// Tangkep error dari backend Go lu
			Swal.fire('Gagal!', err.response?.data?.error || err.message, 'error');
		} finally {
			isSubmitting.value = false;
		}
	};

	const deleteProduct = async (publicId) => {
		// 🔐 SECURITY FIXED: Parameter diganti eksplisit menerima publicId (ULID)
		const res = await Swal.fire({ title: 'Hapus Produk?', text: 'Data tidak bisa dikembalikan!', icon: 'warning', showCancelButton: true, confirmButtonColor: '#ef4444', confirmButtonText: 'Ya, Hapus!' });
		if (res.isConfirmed) {
			try {
				await productService.deleteProduct(publicId);
				Swal.fire('Terhapus!', 'Produk terhapus.', 'success');
				fetchProducts(products.value.length === 1 && currentPage.value > 1 ? currentPage.value - 1 : currentPage.value);
				fetchCategories();
			} catch (e) {
				Swal.fire('Gagal!', 'Gagal menghapus produk.', 'error');
			}
		}
	};

	const triggerImport = () => {
		if (fileInput.value) fileInput.value.click();
	};
	const handleImport = async (e) => {
		const file = e.target.files[0];
		if (!file) return;
		const fData = new FormData();
		fData.append('file', file);
		Swal.fire({ title: 'Mengimpor Data...', allowOutsideClick: false, didOpen: () => Swal.showLoading() });
		try {
			await productService.importCSV(fData);
			Swal.fire('Berhasil!', 'Import Sukses!', 'success');
			fetchProducts(1);
			fetchCategories();
		} catch (err) {
			Swal.fire('Gagal!', 'Gagal import file CSV.', 'error');
		} finally {
			e.target.value = '';
		}
	};

	const exportCSV = async () => {
		Swal.fire({ title: 'Menyiapkan File...', allowOutsideClick: false, didOpen: () => Swal.showLoading() });
		try {
			// 🔐 SECURITY FIXED: Gunakan Axios instance/service agar interseptor auth 401 & token otomatis terpasang sinkron
			const res = await productService.exportProducts();
			const url = window.URL.createObjectURL(new Blob([res.data])),
				a = document.createElement('a');
			a.href = url;
			a.download = 'master_produk.csv';
			document.body.appendChild(a);
			a.click();
			a.remove();
			URL.revokeObjectURL(url);
			Swal.close();
		} catch (e) {
			Swal.fire('Gagal!', 'Gagal mengunduh file CSV.', 'error');
		}
	};

	const visiblePages = computed(() => {
		let p = [],
			s = Math.max(1, currentPage.value - 2),
			e = Math.min(totalPages.value, currentPage.value + 2);
		if (e - s < 4) {
			if (s === 1) e = Math.min(totalPages.value, s + 4);
			else if (e === totalPages.value) s = Math.max(1, e - 4);
		}
		for (let i = s; i <= e; i++) p.push(i);
		return p;
	});

	onMounted(() => {
		fetchProducts(1);
		fetchCategories();
	});
	onBeforeUnmount(() => {
		clearTimeout(scannerTimeout);
		if (showScanner.value) stopScanner();
		if (imagePreview.value && imagePreview.value.startsWith('blob:')) URL.revokeObjectURL(imagePreview.value);
	});

	return {
		products,
		categories,
		isLoading,
		isSubmitting,
		stok_dalam_karton,
		eceran_tambahan,
		searchQuery,
		selectedCategory,
		currentPage,
		totalPages,
		showFormModal,
		isEditing,
		fileInput,
		imagePreview,
		form,
		showScanner,
		visiblePages,
		getImageUrl,
		changePage,
		onFileChange,
		openAddModal,
		editProduct,
		submitProduct,
		deleteProduct,
		triggerImport,
		handleImport,
		exportCSV,
		startScanner,
		stopScanner,
	};
}
