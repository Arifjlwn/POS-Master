import { ref, computed, onMounted, watch, onBeforeUnmount, nextTick } from 'vue';
import { useRouter } from 'vue-router';
import { productService } from '../services/productService.js';
import Swal from 'sweetalert2';
import { Html5Qrcode } from "html5-qrcode";

export function useMasterProduk() {
    const router = useRouter();

    // Data States
    const products = ref([]);
    const categories = ref([]);
    const isLoading = ref(true);
    const isSubmitting = ref(false);
    
    // Smart Mixing Stok States
    const stok_dalam_karton = ref(null);
    const eceran_tambahan = ref(null);

    // Filter, Search & Pagination States
    const searchQuery = ref('');
    const selectedCategory = ref('');
    const currentPage = ref(1);
    const limitPerPage = ref(10); 
    const totalPages = ref(1);

    // Modal & Form States
    const showFormModal = ref(false);
    const isEditing = ref(false);
    const editId = ref(null);
    const fileInput = ref(null);
    const imagePreview = ref(null);

    const form = ref({
        name: '', sku: '', category: '', 
        cost_price: 0, price: 0, stock: 0, image: null,
        satuan_dasar: 'PCS', has_satuan_besar: false, satuan_besar: '', 
        isi_per_besar: null, harga_beli_besar: null, harga_jual_besar: null,
        harga_eceran_tampil: 0, qty_eceran_tampil: 1,

        // 🚀 AMUNISI KALKULATOR SILUMAN
        input_kg: null,
        is_nested_uom: false,
        satuan_tengah: '',
        isi_besar_ke_tengah: null,
        isi_tengah_ke_dasar: null
    });

    // Barcode Scanner Kamera States
    const showScanner = ref(false);
    let html5QrCode = null;

    // --- UTILS ---
    const getImageUrl = (path) => path ? `${import.meta.env.VITE_API_BASE_URL}${path}` : null;

    // --- PENGAMBILAN DATA (API CONSUMER) ---
    const fetchCategories = async () => {
        try {
            const res = await productService.getCategories();
            categories.value = res.data.data;
        } catch (error) {}
    };

    const fetchProducts = async (page = 1) => {
        isLoading.value = true;
        currentPage.value = page;
        try {
            const res = await productService.getProducts({
                page: page, limit: limitPerPage.value,
                search: searchQuery.value, category: selectedCategory.value
            });
            products.value = res.data.data; 
            totalPages.value = Math.ceil(res.data.total_items / limitPerPage.value) || 1;
        } catch (error) {} 
        finally { isLoading.value = false; }
    };

    // 🚀 KALKULATOR SILUMAN: GULA (KG -> GRAM) & MINYAK (LITER -> ML)
    watch(
        () => form.value.input_kg, // Variabel input_kg dipake bareng buat LITER juga biar irit
        (val) => {
            if ((form.value.satuan_dasar === 'GRAM' || form.value.satuan_dasar === 'ML') && val > 0) {
                form.value.isi_per_besar = val * 1000;
            }
        }
    );

    // 🚀 KALKULATOR SILUMAN: ROKOK / KEMASAN BERLAPIS (SLOP -> BUNGKUS -> BATANG)
    watch(
        () => [form.value.isi_besar_ke_tengah, form.value.isi_tengah_ke_dasar, form.value.is_nested_uom],
        ([besarTengah, tengahDasar, isNested]) => {
            if (isNested && besarTengah > 0 && tengahDasar > 0) {
                form.value.isi_per_besar = besarTengah * tengahDasar;
            }
        }
    );

    // KALKULATOR GROSIR
    watch(
        () => [form.value.harga_beli_besar, form.value.isi_per_besar, form.value.has_satuan_besar, stok_dalam_karton.value, eceran_tambahan.value],
        ([hargaBesar, isiPerBesar, hasSatuanBesar, jmlKarton, jmlEceran]) => {
            if (hasSatuanBesar && hargaBesar > 0 && isiPerBesar > 0) {
                form.value.cost_price = Math.round(hargaBesar / isiPerBesar);
            }
            if (hasSatuanBesar && isiPerBesar > 0) {
                const karton = parseInt(jmlKarton) || 0;
                const eceran = parseInt(jmlEceran) || 0;
                const isi = parseInt(isiPerBesar) || 0;
                form.value.stock = (karton * isi) + eceran;
            } else if (!hasSatuanBesar) {
                stok_dalam_karton.value = null;
                eceran_tambahan.value = null;
            }
        }
    );

    // KALKULATOR PINTAR ECERAN WARUNG MADURA
    watch(
        () => [form.value.harga_eceran_tampil, form.value.qty_eceran_tampil],
        ([harga, qty]) => {
            if (qty > 0 && harga >= 0) {
                form.value.price = harga / qty;
            }
        }
    );

    // 🚀 AUTO SWITCH QTY PATOKAN JIKA PILIH GRAM / ML
    watch(
        () => form.value.satuan_dasar,
        (newVal) => {
            if (newVal === 'GRAM' || newVal === 'ML') {
                form.value.qty_eceran_tampil = 1000; // Otomatis nanya harga per 1000 Gram / 1000 ML
            } else {
                form.value.qty_eceran_tampil = 1;
            }
            form.value.harga_eceran_tampil = form.value.price * form.value.qty_eceran_tampil;
        }
    );

    // Debounce watcher untuk fitur ketik live-search
    let searchTimer;
    watch([searchQuery, selectedCategory], () => {
        clearTimeout(searchTimer);
        searchTimer = setTimeout(() => {
            fetchProducts(1); 
        }, 500);
    });

    // --- LOGIKA SCANNER BARCODE PRODUK ---
    const startScanner = async () => {
        showScanner.value = true;
        setTimeout(async () => {
            try {
                html5QrCode = new Html5Qrcode("reader");
                await html5QrCode.start(
                    { facingMode: "environment" }, 
                    { fps: 10, qrbox: { width: 250, height: 100 } }, 
                    (decodedText) => {
                        form.value.sku = decodedText; 
                        stopScanner();
                        const audio = new Audio('https://www.soundjay.com/buttons/beep-07a.mp3');
                        audio.play().catch(()=>{}); 
                    },
                    () => {} 
                );
            } catch (err) {
                console.error(err);
                Swal.fire('Oops!', 'Gagal menyalakan kamera scanner.', 'error');
                stopScanner();
            }
        }, 300);
    };

    const stopScanner = () => {
        if (html5QrCode) {
            html5QrCode.stop().then(() => {
                html5QrCode.clear();
                showScanner.value = false;
            }).catch(() => {
                showScanner.value = false;
            });
        } else {
            showScanner.value = false;
        }
    };

    // --- MANIPULASI DATA (CRUD ACTION HANDLERS) ---
    const changePage = (newPage) => { if (newPage >= 1 && newPage <= totalPages.value) fetchProducts(newPage); };

    const onFileChange = (e) => {
        const file = e.target.files[0];
        if (file) {
            form.value.image = file;
            imagePreview.value = URL.createObjectURL(file);
        }
    };

    const openAddModal = () => {
        isEditing.value = false;
        editId.value = null;
        form.value = { 
            name: '', sku: '', category: '', cost_price: 0, price: 0, stock: 0, image: null,
            satuan_dasar: 'PCS', has_satuan_besar: false, satuan_besar: '', isi_per_besar: null, 
            harga_beli_besar: null, harga_jual_besar: null,
            harga_eceran_tampil: 0, qty_eceran_tampil: 1,
            input_kg: null, is_nested_uom: false, satuan_tengah: '', isi_besar_ke_tengah: null, isi_tengah_ke_dasar: null // 🚀 RESET
        };
        imagePreview.value = null;
        stok_dalam_karton.value = null;
        eceran_tambahan.value = null;
        showFormModal.value = true;
    };

    const editProduct = (product) => {
        isEditing.value = true;
        editId.value = product.id;
        
        const isBulkType = product.satuan_dasar === 'GRAM' || product.satuan_dasar === 'ML';
        const acuanQty = isBulkType ? 1000 : 1;

        form.value = {
            name: product.nama_produk || '',
            sku: product.sku || '',
            category: product.kategori || '',
            cost_price: product.harga_modal || 0,
            price: product.harga_jual || 0,
            stock: product.stok || 0,
            image: null,
            satuan_dasar: product.satuan_dasar || 'PCS',
            has_satuan_besar: !!product.satuan_besar,
            satuan_besar: product.satuan_besar || '',
            isi_per_besar: product.isi_per_besar || null,
            harga_beli_besar: null,
            harga_jual_besar: product.harga_jual_besar || null,
            qty_eceran_tampil: acuanQty,
            harga_eceran_tampil: (product.harga_jual || 0) * acuanQty,
            
            // 🚀 REVERSE ENGINEER: Kalau GRAM, kita bagi 1000 lagi pas mau nampilin ke kotak KG
            input_kg: (isBulkType && product.isi_per_besar) ? (product.isi_per_besar / 1000) : null,
            is_nested_uom: false, satuan_tengah: '', isi_besar_ke_tengah: null, isi_tengah_ke_dasar: null
        };
        
        imagePreview.value = getImageUrl(product.gambar);
        showFormModal.value = true;
    };

    const submitProduct = async () => {
        if (form.value.has_satuan_besar && (!form.value.satuan_besar || !form.value.isi_per_besar || !form.value.harga_jual_besar)) {
            return Swal.fire('Data Kurang!', 'Lengkapi Nama, Isi, dan Harga Jual untuk satuan besar!', 'warning');
        }

        isSubmitting.value = true;
        const formData = new FormData();
        formData.append('nama_produk', form.value.name);
        formData.append('sku', form.value.sku);
        formData.append('kategori', form.value.category);
        formData.append('harga_modal', Number(form.value.cost_price));
        formData.append('harga_jual', Number(form.value.price));
        formData.append('stok', Number(form.value.stock));
        formData.append('satuan_dasar', form.value.satuan_dasar);
        formData.append('satuan_besar', form.value.has_satuan_besar ? form.value.satuan_besar : '');
        formData.append('isi_per_besar', form.value.has_satuan_besar ? Number(form.value.isi_per_besar) : 0);
        formData.append('harga_jual_besar', form.value.has_satuan_besar ? Number(form.value.harga_jual_besar) : 0);

        // 🚀 INI YANG KETINGGALAN BOSKU! DATA 3 LAPISNYA HARUS DIKIRIM
        formData.append('is_nested_uom', form.value.is_nested_uom);
        if (form.value.is_nested_uom) {
            formData.append('satuan_tengah', form.value.satuan_tengah);
            formData.append('isi_besar_ke_tengah', Number(form.value.isi_besar_ke_tengah));
            formData.append('isi_tengah_ke_dasar', Number(form.value.isi_tengah_ke_dasar));
        }

        if (form.value.image) formData.append('gambar', form.value.image);

        try {
            if (isEditing.value) await productService.updateProduct(editId.value, formData);
            else await productService.createProduct(formData);

            Swal.fire({ icon: 'success', title: 'Berhasil!', text: 'Sip, Data tersimpan!', timer: 2000, showConfirmButton: false });
            showFormModal.value = false;
            fetchProducts(currentPage.value);
            fetchCategories();
        } catch (error) {
            Swal.fire('Gagal!', error.response?.data?.error || error.message, 'error');
        } finally {
            isSubmitting.value = false;
        }
    };

    const deleteProduct = async (id) => {
        const result = await Swal.fire({
            title: 'Hapus Produk?',
            text: "Data yang dihapus tidak bisa dikembalikan!",
            icon: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#ef4444',
            confirmButtonText: 'Ya, Hapus!'
        });

        if (result.isConfirmed) {
            try {
                await productService.deleteProduct(id);
                Swal.fire('Terhapus!', 'Produk telah dihapus dari sistem.', 'success');
                
                if (products.value.length === 1 && currentPage.value > 1) {
                    fetchProducts(currentPage.value - 1);
                } else {
                    fetchProducts(currentPage.value);
                }
                fetchCategories();
            } catch (error) {
                Swal.fire('Gagal!', 'Gagal menghapus produk.', 'error');
            }
        }
    };

    // --- IMPOR & EKSPOR DATA CSV CORE ---
    const triggerImport = () => { if (fileInput.value) fileInput.value.click(); };

    const handleImport = async (event) => {
        const file = event.target.files[0];
        if (!file) return;

        const formData = new FormData();
        formData.append('file', file);
        Swal.fire({ title: 'Mengimpor Data...', allowOutsideClick: false, didOpen: () => Swal.showLoading() });

        try {
            await productService.importCSV(formData);
            Swal.fire('Berhasil!', 'Import CSV Sukses!', 'success');
            fetchProducts(1);
            fetchCategories();
        } catch (error) {
            Swal.fire('Gagal!', 'Gagal import file CSV.', 'error');
        } finally {
            event.target.value = '';
        }
    };

    const exportCSV = async () => {
        const token = localStorage.getItem('token');
        Swal.fire({ title: 'Menyiapkan File...', allowOutsideClick: false, didOpen: () => Swal.showLoading() });
        
        try {
            const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/api/retail/products/export`, {
                method: 'GET',
                headers: { 'Authorization': `Bearer ${token}` }
            });
            if (!response.ok) throw new Error("Gagal ekspor");

            const blob = await response.blob();
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = 'master_produk.csv';
            document.body.appendChild(a);
            a.click();
            a.remove();
            Swal.close();
        } catch (error) {
            Swal.fire('Gagal!', 'Gagal mengunduh file CSV.', 'error');
        }
    };

    const visiblePages = computed(() => {
        let pages = [];
        let start = Math.max(1, currentPage.value - 2);
        let end = Math.min(totalPages.value, currentPage.value + 2);

        if (end - start < 4) {
            if (start === 1) end = Math.min(totalPages.value, start + 4);
            else if (end === totalPages.value) start = Math.max(1, end - 4);
        }

        for (let i = start; i <= end; i++) pages.push(i);
        return pages;
    });

    onMounted(() => { fetchProducts(1); fetchCategories(); });
    onBeforeUnmount(() => { if (showScanner.value) stopScanner(); })

    return {
        products, categories, isLoading, isSubmitting, stok_dalam_karton, eceran_tambahan,
        searchQuery, selectedCategory, currentPage, totalPages, showFormModal, isEditing,
        fileInput, imagePreview, form, showScanner, visiblePages,
        getImageUrl, changePage, onFileChange, openAddModal, editProduct, submitProduct,
        deleteProduct, triggerImport, handleImport, exportCSV, startScanner, stopScanner
    };
}