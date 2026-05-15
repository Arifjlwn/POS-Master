<script setup>
import { ref, onMounted, watch } from 'vue';
import { useRouter } from 'vue-router';
import api from '../../api.js';
import Sidebar from '../../components/Sidebar.vue';
import Swal from 'sweetalert2';

const router = useRouter();

// State Data
const products = ref([]);
const categories = ref([]);
const isLoading = ref(true);
const isSubmitting = ref(false);

const fetchCategories = async () => {
    try {
        const response = await api.get('/categories');
        categories.value = response.data.data;
    } catch (error) {
        console.error("Gagal ambil kategori:", error);
    }
};

// State Filter & Search
const searchQuery = ref('');
const selectedCategory = ref('');
const currentPage = ref(1);

// State Form Modal
const showFormModal = ref(false);
const isEditing = ref(false);
const editId = ref(null);
const fileInput = ref(null);
const imagePreview = ref(null);

const form = ref({
    name: '',
    sku: '',
    category: '',
    cost_price: 0,
    price: 0,
    stock: 0,
    image: null,
    satuan_dasar: 'PCS',
    has_satuan_besar: false,
    satuan_besar: '',
    isi_per_besar: null,
    harga_beli_besar: null 
});

// 🚀 WATCHER KALKULATOR PINTAR: Ngitung otomatis Harga Modal!
watch(
    () => [form.value.harga_beli_besar, form.value.isi_per_besar, form.value.has_satuan_besar],
    ([hargaBesar, isiPerBesar, hasSatuanBesar]) => {
        if (hasSatuanBesar && hargaBesar > 0 && isiPerBesar > 0) {
            form.value.cost_price = Math.round(hargaBesar / isiPerBesar);
        }
    }
);

// Helper untuk nampilin gambar dari Golang
const getImageUrl = (path) => {
    return path ? `${import.meta.env.VITE_API_BASE_URL}${path}` : null;
};

// --- FUNGSI TARIK DATA (API GOLANG) ---
const fetchProducts = async (page = 1) => {
    isLoading.value = true;
    try {
        const response = await api.get('/products', {
            params: {
                page: page,
                search: searchQuery.value,
                category: selectedCategory.value
            }
        });
        products.value = response.data.data; 
    } catch (error) {
        console.error("Gagal ambil produk:", error);
    } finally {
        isLoading.value = false;
    }
};

onMounted(() =>{
     fetchProducts();
     fetchCategories();
});

let searchTimer;
watch([searchQuery, selectedCategory], () => {
    clearTimeout(searchTimer);
    searchTimer = setTimeout(() => {
        fetchProducts(1);
    }, 500);
});

// --- FUNGSI CRUD & UPLOAD FOTO ---
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
        satuan_dasar: 'PCS', has_satuan_besar: false, satuan_besar: '', isi_per_besar: null, harga_beli_besar: null
    };
    imagePreview.value = null;
    showFormModal.value = true;
};

const editProduct = (product) => {
    isEditing.value = true;
    editId.value = product.id;
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
        harga_beli_besar: null // Dikosongkan saat edit biar nggak over-ride harga yang udah ada
    };
    imagePreview.value = getImageUrl(product.gambar);
    showFormModal.value = true;
};

const submitProduct = async () => {
    if (form.value.has_satuan_besar && (!form.value.satuan_besar || !form.value.isi_per_besar)) {
        return Swal.fire('Data Kurang!', 'Lengkapi detail satuan besar beserta isinya!', 'warning');
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

    if (form.value.image) {
        formData.append('gambar', form.value.image);
    }

    try {
        const token = localStorage.getItem('token');
        const url = isEditing.value 
            ? `${import.meta.env.VITE_API_BASE_URL}/api/products/${editId.value}` 
            : `${import.meta.env.VITE_API_BASE_URL}/api/products`;
            
        const method = isEditing.value ? 'PUT' : 'POST';

        const response = await fetch(url, {
            method: method,
            headers: { 'Authorization': `Bearer ${token}` },
            body: formData
        });

        const data = await response.json();
        if (!response.ok) throw new Error(data.error || 'Gagal menyimpan produk dari server');

        Swal.fire({
            icon: 'success',
            title: 'Berhasil!',
            text: `Produk berhasil ${isEditing.value ? 'diperbarui' : 'ditambahkan'}!`,
            timer: 2000,
            showConfirmButton: false
        });
        
        showFormModal.value = false;
        fetchProducts(currentPage.value);
        fetchCategories();
    } catch (error) {
        Swal.fire('Gagal!', error.message, 'error');
        console.error("Detail Error:", error);
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
        confirmButtonColor: '#d33',
        confirmButtonText: 'Ya, Hapus!'
    });

    if (result.isConfirmed) {
        try {
            await api.delete(`/products/${id}`);
            Swal.fire('Terhapus!', 'Produk telah dihapus dari sistem.', 'success');
            fetchProducts(currentPage.value);
        } catch (error) {
            Swal.fire('Gagal!', 'Gagal menghapus produk.', 'error');
        }
    }
};

const triggerImport = () => fileInput.value.click();

const handleImport = async (event) => {
    const file = event.target.files[0];
    if (!file) return;

    const formData = new FormData();
    formData.append('file', file);

    Swal.fire({ title: 'Mengimpor Data...', allowOutsideClick: false, didOpen: () => Swal.showLoading() });

    try {
        await api.post('/products/import', formData);
        Swal.fire('Berhasil!', 'Import CSV Sukses!', 'success');
        fetchProducts();
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
        const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/api/products/export`, {
            method: 'GET',
            headers: { 'Authorization': `Bearer ${token}` }
        });
        
        if (!response.ok) throw new Error("Gagal ekspor");

        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = 'katalog_produk.csv';
        document.body.appendChild(a);
        a.click();
        a.remove();
        Swal.close();
    } catch (error) {
        Swal.fire('Gagal!', 'Gagal mengunduh file CSV.', 'error');
    }
};
</script>

<template>
    <Sidebar>
        <div class="p-4 md:p-8 lg:p-10 w-full max-w-7xl mx-auto font-sans">
            <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-6 mb-8">
                <div>
                    <h1 class="text-3xl font-black text-slate-900 tracking-tighter uppercase italic flex items-center gap-3">
                        <span class="p-3 bg-blue-600 text-white rounded-2xl shadow-lg shadow-blue-200">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/><polyline points="3.27 6.96 12 12.01 20.73 6.96"/><line x1="12" y1="22.08" x2="12" y2="12"/></svg>
                        </span>
                        Master <span class="text-blue-600">Inventory</span>
                    </h1>
                    <p class="text-xs text-slate-400 font-black uppercase tracking-[0.2em] mt-2 ml-16">Katalog Produk & Harga</p>
                </div>

                <div class="flex flex-wrap items-center gap-3">
                    <button @click="exportCSV" class="bg-emerald-50 hover:bg-emerald-500 text-emerald-600 hover:text-white px-5 py-3 rounded-[16px] font-black transition-all shadow-sm items-center gap-2 flex text-[10px] uppercase tracking-widest border border-emerald-100">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" x2="12" y1="3" y2="15"/></svg>
                        Ekspor
                    </button>

                    <button @click="triggerImport" class="bg-amber-50 hover:bg-amber-500 text-amber-600 hover:text-white px-5 py-3 rounded-[16px] font-black transition-all shadow-sm items-center gap-2 flex text-[10px] uppercase tracking-widest border border-amber-100">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" x2="12" y1="3" y2="15"/></svg>
                        Import
                    </button>
                    <input type="file" ref="fileInput" class="hidden" accept=".csv" @change="handleImport">

                    <button @click="openAddModal" class="bg-blue-600 hover:bg-slate-900 text-white px-6 py-3 rounded-[16px] font-black transition-all shadow-xl shadow-blue-200 items-center gap-2 flex text-[10px] uppercase tracking-widest active:scale-95 ml-2">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
                        Tambah Barang
                    </button>
                </div>
            </div>

            <div class="flex flex-col sm:flex-row gap-4 mb-6">
                <div class="relative flex-1">
                    <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-300" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
                    </div>
                    <input v-model="searchQuery" type="text" placeholder="Cari nama barang atau barcode..." class="block w-full pl-12 pr-4 py-3.5 bg-white rounded-2xl border border-slate-100 shadow-sm focus:border-blue-500 focus:ring-4 focus:ring-blue-500/10 text-sm font-bold outline-none transition-all placeholder:text-slate-300">
                </div>
                <div class="w-full sm:w-64 shrink-0 relative">
                    <select v-model="selectedCategory" class="block w-full pl-4 pr-10 py-3.5 bg-white rounded-2xl border border-slate-100 shadow-sm focus:border-blue-500 focus:ring-4 focus:ring-blue-500/10 text-sm font-bold text-slate-700 cursor-pointer outline-none appearance-none transition-all">
                        <option value="">SEMUA KATEGORI</option>
                        <option v-for="cat in categories" :key="cat" :value="cat">{{ cat.toUpperCase() }}</option>
                    </select>
                    <div class="absolute inset-y-0 right-0 pr-4 flex items-center pointer-events-none">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="6 9 12 15 18 9"/></svg>
                    </div>
                </div>
            </div>

            <div class="w-full bg-white rounded-[32px] shadow-sm border border-slate-100 overflow-hidden">
                <div class="overflow-x-auto custom-scrollbar">
                    <table class="w-full text-left border-collapse whitespace-nowrap min-w-[800px]">
                        <thead>
                            <tr class="bg-slate-50/50 text-slate-400 border-b border-slate-100">
                                <th class="p-5 font-black uppercase tracking-widest text-[10px]">Produk & Detail</th>
                                <th class="p-5 font-black uppercase tracking-widest text-[10px]">Kategori</th>
                                <th class="p-5 font-black uppercase tracking-widest text-[10px] text-right">Modal / Dasar</th>
                                <th class="p-5 font-black uppercase tracking-widest text-[10px] text-right">Harga Jual</th>
                                <th class="p-5 font-black uppercase tracking-widest text-[10px] text-center">Stok</th>
                                <th class="p-5 font-black uppercase tracking-widest text-[10px] text-center">Aksi</th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-slate-50">
                            <tr v-if="isLoading">
                                <td colspan="6" class="p-16 text-center text-slate-400 font-black text-xs uppercase tracking-widest animate-pulse">Memuat Katalog Produk...</td>
                            </tr>
                            <tr v-else-if="products.length === 0">
                                <td colspan="6" class="p-16 text-center text-slate-400 font-black text-xs uppercase tracking-widest opacity-50">Produk Tidak Ditemukan</td>
                            </tr>
                            <tr v-else v-for="product in products" :key="product.id" class="hover:bg-slate-50/80 transition-colors group">
                                <td class="p-5 flex items-center gap-4">
                                    <div class="w-14 h-14 rounded-[16px] border-2 border-white shadow-sm bg-slate-50 flex items-center justify-center text-slate-300 overflow-hidden shrink-0 group-hover:border-blue-100 transition-colors">
                                        <img v-if="product.gambar" :src="getImageUrl(product.gambar)" class="w-full h-full object-cover">
                                        <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/></svg>
                                    </div>
                                    <div>
                                        <div class="font-black text-slate-800 text-sm uppercase">{{ product.nama_produk }}</div>
                                        <div class="flex items-center gap-2 mt-1">
                                            <span class="text-[9px] font-black bg-slate-100 text-slate-400 px-2 py-0.5 rounded uppercase tracking-widest">{{ product.sku || 'NO-SKU' }}</span>
                                            <span v-if="product.satuan_besar" class="text-[9px] font-black bg-purple-50 text-purple-600 px-2 py-0.5 rounded uppercase tracking-widest border border-purple-100">{{ product.satuan_besar }} / {{ product.isi_per_besar }}</span>
                                        </div>
                                    </div>
                                </td>
                                <td class="p-5">
                                    <span class="bg-blue-50 text-blue-600 border border-blue-100 px-3 py-1.5 rounded-lg font-black text-[9px] uppercase tracking-widest">{{ product.kategori || 'General' }}</span>
                                </td>
                                <td class="p-5 text-slate-400 font-bold text-right text-xs">
                                    Rp {{ (product.harga_modal || 0).toLocaleString('id-ID') }}
                                    <div class="text-[8px] font-black uppercase mt-0.5 tracking-widest">/ {{ product.satuan_dasar || 'PCS' }}</div>
                                </td>
                                <td class="p-5 text-emerald-600 font-black text-right text-sm">
                                    Rp {{ (product.harga_jual || 0).toLocaleString('id-ID') }}
                                </td>
                                <td class="p-5 text-center">
                                    <span class="px-3 py-1.5 text-[10px] rounded-lg font-black tracking-widest shadow-sm inline-flex items-center gap-1.5" :class="(product.stok || 0) > 10 ? 'bg-emerald-50 text-emerald-600 border border-emerald-100' : 'bg-red-50 text-red-600 border border-red-100'">
                                        {{ product.stok || 0 }} <span class="uppercase text-[8px]">{{ product.satuan_dasar || 'PCS' }}</span>
                                    </span>
                                </td>
                                <td class="p-5 text-center">
                                    <div class="flex justify-center gap-2">
                                        <button @click="editProduct(product)" class="bg-slate-100 text-slate-400 p-2.5 rounded-xl hover:bg-blue-600 hover:text-white transition-colors" title="Edit">
                                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M17 3a2.828 2.828 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5L17 3z"/></svg>
                                        </button>
                                        <button @click="deleteProduct(product.id)" class="bg-slate-100 text-slate-400 p-2.5 rounded-xl hover:bg-red-500 hover:text-white transition-colors" title="Hapus">
                                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/></svg>
                                        </button>
                                    </div>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>

            <div v-if="showFormModal" class="fixed inset-0 bg-slate-900/80 flex items-center justify-center z-[100] p-4 backdrop-blur-sm">
                <div class="bg-white rounded-[32px] shadow-2xl w-full max-w-2xl flex flex-col max-h-[90vh] overflow-hidden border border-slate-100">
                    
                    <div class="p-6 border-b border-slate-50 bg-slate-50/50 flex justify-between items-center shrink-0">
                        <h2 class="text-xl font-black text-slate-800 uppercase italic">{{ isEditing ? 'Edit Data Produk' : 'Registrasi Produk Baru' }}</h2>
                        <button @click="showFormModal = false" class="p-2 rounded-xl bg-white text-slate-400 hover:text-red-500 hover:bg-red-50 transition-all border border-slate-100 shadow-sm">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                        </button>
                    </div>

                    <div class="p-6 md:p-8 overflow-y-auto custom-scrollbar">
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
                            
                            <div class="md:col-span-2 flex items-center gap-5 p-4 rounded-[24px] bg-slate-50 border border-slate-100 mb-2">
                                <div @click="$refs.fileInput.click()" class="w-16 h-16 rounded-[18px] border-2 border-dashed border-slate-300 flex items-center justify-center bg-white cursor-pointer overflow-hidden shadow-inner">
                                    <img v-if="imagePreview" :src="imagePreview" class="w-full h-full object-cover">
                                    <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-slate-300" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14.5 4h-5L7 7H4a2 2 0 0 0-2 2v9a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2h-3l-2.5-3z"/><circle cx="12" cy="13" r="3"/></svg>
                                </div>
                                <div class="flex-1">
                                    <input v-model="form.name" type="text" placeholder="NAMA BARANG..." class="w-full px-4 py-3 rounded-xl border border-slate-200 focus:border-blue-600 outline-none font-black text-sm uppercase mb-2 text-slate-800">
                                    <input v-model="form.sku" type="text" placeholder="BARCODE / SKU..." class="w-full px-4 py-2 rounded-xl border border-slate-200 focus:border-blue-600 outline-none font-bold text-xs uppercase text-slate-800">
                                </div>
                                <input type="file" ref="fileInput" @change="onFileChange" accept="image/*" class="hidden">
                            </div>

                            <div class="md:col-span-2">
                                <label class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1.5 block">Kategori</label>
                                <input list="kategori-list" v-model="form.category" placeholder="Pilih / Ketik Kategori Baru..." class="w-full px-4 py-3.5 rounded-xl border border-slate-200 focus:border-blue-600 focus:ring-4 focus:ring-blue-600/10 outline-none font-bold text-sm bg-white uppercase transition-all text-slate-800">
                                <datalist id="kategori-list"><option v-for="cat in categories" :key="cat" :value="cat"></option></datalist>
                            </div>

                            <div class="md:col-span-2 p-5 bg-slate-900 rounded-[28px] text-white shadow-xl mt-2 mb-2">
                                <div class="flex items-center gap-2 mb-4">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-blue-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="m7.5 4.27 9 5.15"/><path d="M21 8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16Z"/><path d="M12 22V12"/></svg>
                                    <h4 class="font-black text-[10px] uppercase tracking-[0.2em]">Konversi & Satuan</h4>
                                </div>

                                <div class="grid grid-cols-2 gap-4">
                                    <div>
                                        <label class="text-[8px] font-black text-slate-500 uppercase block mb-1">Satuan Dasar</label>
                                        <select v-model="form.satuan_dasar" class="w-full p-3.5 bg-slate-800 border border-slate-700 rounded-xl outline-none font-black text-xs uppercase cursor-pointer text-white">
                                            <option value="PCS">PCS</option><option value="KG">KG</option><option value="GRAM">GRAM</option><option value="PACK">PACK</option><option value="LITER">LITER</option>
                                        </select>
                                    </div>
                                    <div>
                                        <label class="text-[8px] font-black text-slate-500 uppercase block mb-1">Punya Grosir?</label>
                                        <div @click="form.has_satuan_besar = !form.has_satuan_besar" class="w-full p-3.5 bg-slate-800 border border-slate-700 rounded-xl font-black text-[10px] uppercase cursor-pointer flex items-center justify-between transition-colors hover:border-blue-500">
                                            {{ form.has_satuan_besar ? 'YA (AKTIF)' : 'TIDAK' }}
                                            <div :class="form.has_satuan_besar ? 'bg-blue-500' : 'bg-slate-600'" class="w-2 h-2 rounded-full shadow-[0_0_8px_rgba(59,130,246,0.5)] transition-colors"></div>
                                        </div>
                                    </div>

                                    <div v-if="form.has_satuan_besar" class="col-span-2 grid grid-cols-3 gap-3 pt-4 border-t border-slate-800 mt-2 animate-[fadeIn_0.3s_ease-out]">
                                        <div class="col-span-1">
                                            <label class="text-[8px] font-black text-blue-400 uppercase block mb-1">Nama Grosir</label>
                                            <input v-model="form.satuan_besar" type="text" placeholder="KARTON/BOX" class="w-full p-3 bg-slate-800 border border-blue-900 focus:border-blue-500 rounded-xl outline-none font-black text-xs uppercase text-white placeholder:text-slate-600 transition-all">
                                        </div>
                                        <div class="col-span-1">
                                            <label class="text-[8px] font-black text-blue-400 uppercase block mb-1">Isi per {{ form.satuan_besar || '...' }}</label>
                                            <input v-model.number="form.isi_per_besar" type="number" class="w-full p-3 bg-slate-800 border border-blue-900 focus:border-blue-500 rounded-xl outline-none font-black text-xs text-white transition-all">
                                        </div>
                                        <div class="col-span-1">
                                            <label class="text-[8px] font-black text-amber-400 uppercase block mb-1">Harga 1 {{ form.satuan_besar || '...' }}</label>
                                            <input v-model.number="form.harga_beli_besar" type="number" placeholder="Rp" class="w-full p-3 bg-amber-900/20 border border-amber-900 focus:border-amber-500 rounded-xl outline-none font-black text-xs text-amber-400 transition-all placeholder:text-amber-900/50">
                                        </div>
                                    </div>
                                </div>
                            </div>

                            <div class="p-5 bg-white border border-slate-200 rounded-[28px] shadow-sm">
                                <label class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-3 block">Harga Modal Dasar</label>
                                <div class="relative">
                                    <span class="absolute inset-y-0 left-0 pl-4 flex items-center text-xs font-black text-slate-400">Rp</span>
                                    <input v-model.number="form.cost_price" type="number" class="w-full pl-12 pr-4 py-4 rounded-2xl bg-slate-50 border border-slate-100 outline-none font-black text-lg transition-all" :class="form.has_satuan_besar ? 'text-blue-600 ring-2 ring-blue-500/10' : 'text-slate-800'">
                                </div>
                                <p v-if="form.has_satuan_besar" class="text-[8px] font-black text-blue-500 mt-2 uppercase tracking-widest italic text-center">*Dihitung otomatis dari kalkulator grosir di atas</p>
                            </div>

                            <div class="p-5 bg-blue-50 border border-blue-100 rounded-[28px] shadow-sm">
                                <label class="text-[9px] font-black text-blue-400 uppercase tracking-widest mb-3 block">Harga Jual Eceran</label>
                                <div class="relative">
                                    <span class="absolute inset-y-0 left-0 pl-4 flex items-center text-xs font-black text-blue-400">Rp</span>
                                    <input v-model.number="form.price" type="number" class="w-full pl-12 pr-4 py-4 rounded-2xl bg-white border border-blue-200 focus:border-blue-600 outline-none font-black text-lg text-blue-700 shadow-inner transition-all">
                                </div>
                                <p class="text-[9px] font-black text-blue-500 mt-2 uppercase tracking-widest text-center">Profit: Rp {{ form.price - form.cost_price }} / {{ form.satuan_dasar }}</p>
                            </div>

                            <div class="md:col-span-2">
                                <label class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-2 block text-center">Stok Awal Fisik (dalam {{ form.satuan_dasar }})</label>
                                <input v-model.number="form.stock" type="number" class="w-full px-4 py-4 rounded-2xl bg-slate-50 border border-slate-200 focus:border-blue-600 focus:ring-4 focus:ring-blue-600/10 outline-none font-black text-center text-2xl text-slate-800 transition-all">
                            </div>

                        </div>
                    </div>

                    <div class="p-6 bg-slate-50 border-t border-slate-100 shrink-0">
                        <button @click="submitProduct" :disabled="isSubmitting" class="w-full py-5 font-black text-xs uppercase tracking-[0.2em] bg-blue-600 text-white rounded-[24px] shadow-xl shadow-blue-200 hover:bg-slate-900 transition-all active:scale-95 flex items-center justify-center gap-3 disabled:opacity-50">
                            <template v-if="isSubmitting">
                                <div class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                                Menyimpan Data...
                            </template>
                            <template v-else>
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/><polyline points="17 21 17 13 7 13 7 21"/><polyline points="7 3 7 8 15 8"/></svg>
                                Simpan Perubahan Produk
                            </template>
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </Sidebar>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { height: 6px; width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #94a3b8; }

input[type=number]::-webkit-inner-spin-button, 
input[type=number]::-webkit-outer-spin-button { 
  -webkit-appearance: none; 
  margin: 0; 
}
input[type=number] {
  -moz-appearance: textfield;
}

@keyframes fadeIn {
    from { opacity: 0; transform: translateY(-5px); }
    to { opacity: 1; transform: translateY(0); }
}
</style>