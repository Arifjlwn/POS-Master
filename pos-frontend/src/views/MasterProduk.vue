<script setup>
import { ref, onMounted, watch, computed } from 'vue';
import { useRouter } from 'vue-router';
import api from '../api.js';
import Sidebar from '../components/Sidebar.vue';

const router = useRouter();

// State Data
const products = ref([]);

const categories = ref([]);

const fetchCategories = async () => {
    try {
        const response = await api.get('/categories');
        categories.value = response.data.data;
    } catch (error) {
        console.error("Gagal ambil kategori:", error)
    }
};

const isLoading = ref(true);

// State Filter & Search
const searchQuery = ref('');
const selectedCategory = ref('');
const currentPage = ref(1);

// State Form Modal
const showFormModal = ref(false);
const isEditing = ref(false);
const editId = ref(null);
const fileInput = ref(null);
const imagePreview = ref(null); // Tambahan untuk preview foto

const form = ref({
    name: '',
    sku: '',
    category: '',
    cost_price: 0,
    price: 0,
    stock: 0,
    image: null,
});

// Helper untuk nampilin gambar dari Golang
const getImageUrl = (path) => {
    return path ? `http://localhost:8080${path}` : null;
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

// Nangkap file foto pas user milih gambar
const onFileChange = (e) => {
    const file = e.target.files[0];
    if (file) {
        form.value.image = file;
        imagePreview.value = URL.createObjectURL(file); // Bikin preview instan
    }
};

const openAddModal = () => {
    isEditing.value = false;
    editId.value = null;
    form.value = { name: '', sku: '', category: '', cost_price: 0, price: 0, stock: 0, image: null };
    imagePreview.value = null; // Reset preview
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
        image: null
    };
    // Tampilkan gambar lama jika ada
    imagePreview.value = getImageUrl(product.gambar);
    showFormModal.value = true;
};

const submitProduct = async () => {
    // 1. Siapkan Paketnya
    const formData = new FormData();
    formData.append('nama_produk', form.value.name);
    formData.append('sku', form.value.sku);
    formData.append('kategori', form.value.category);
    formData.append('harga_modal', Number(form.value.cost_price));
    formData.append('harga_jual', Number(form.value.price));
    formData.append('stok', Number(form.value.stock));

    if (form.value.image) {
        formData.append('gambar', form.value.image);
    }

    try {
        // 2. Ambil Karcis VIP
        const token = localStorage.getItem('token');
        
        // 3. Tentukan Tujuan (Tambah atau Edit)
        const url = isEditing.value 
            ? `http://localhost:8080/api/products/${editId.value}` 
            : 'http://localhost:8080/api/products';
            
        const method = isEditing.value ? 'PUT' : 'POST';

        // 🚀 4. KITA BYPASS api.js PAKAI NATIVE FETCH
        const response = await fetch(url, {
            method: method,
            headers: {
                'Authorization': `Bearer ${token}`
                // 🚨 PENTING: JANGAN ADA 'Content-Type' DI SINI SAMA SEKALI!
            },
            body: formData
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || 'Gagal menyimpan produk dari server');
        }

        alert(`Produk berhasil ${isEditing.value ? 'diperbarui ✏️' : 'ditambahkan 🚀'}!`);
        showFormModal.value = false;
        fetchProducts(currentPage.value);
        fetchCategories();

    } catch (error) {
        alert('Gagal menyimpan produk: ' + error.message);
        console.error("Detail Error:", error);
    }
};

const deleteProduct = async (id) => {
    if (confirm('Yakin ingin menghapus produk ini? 🗑️')) {
        try {
            await api.delete(`/products/${id}`);
            fetchProducts(currentPage.value);
        } catch (error) {
            alert('Gagal menghapus produk.');
        }
    }
};

const triggerImport = () => fileInput.value.click();

const handleImport = async (event) => {
    const file = event.target.files[0];
    if (!file) return;

    const formData = new FormData();
    formData.append('file', file);

    try {
        await api.post('/products/import', formData);
        alert('Import CSV Berhasil!');
        fetchProducts();
    } catch (error) {
        alert('Gagal import file CSV.');
    } finally {
        event.target.value = '';
    }
};

const exportCSV = async () => {
    const token = localStorage.getItem('token');
    
    // Tarik file pakai fetch biar token bisa dikirim di Header
    try {
        const response = await fetch('http://localhost:8080/api/products/export', {
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
    } catch (error) {
        alert("Gagal mengunduh file CSV.");
    }
};
</script>

<template>
    <Sidebar>
        <div class="p-4 md:p-6 w-full max-w-full mx-auto font-sans">
            <div class="flex justify-between items-center mb-6">
                <div>
                    <h1 class="text-2xl font-black text-gray-800 tracking-tight text-left">📦 Master Produk</h1>
                    <p class="text-sm text-gray-500 mt-1 font-medium">Kelola daftar barang, harga, dan stok toko Anda.</p>
                </div>

                <div class="hidden sm:flex gap-2">
                    <button @click="exportCSV" class="bg-green-100 hover:bg-green-200 text-green-700 px-4 py-2.5 rounded-xl font-bold transition-all shadow-sm items-center gap-2 border border-green-200 flex text-xs">
                        📤 Ekspor CSV
                    </button>

                    <button @click="triggerImport" class="bg-orange-100 hover:bg-orange-200 text-orange-700 px-4 py-2.5 rounded-xl font-bold transition-all shadow-sm items-center gap-2 border border-orange-200 flex text-xs">
                        📥 Import CSV
                    </button>
                    <input type="file" ref="fileInput" class="hidden" accept=".csv" @change="handleImport">

                    <button @click="openAddModal" class="bg-blue-600 hover:bg-blue-700 text-white px-5 py-2.5 rounded-xl font-bold transition-all shadow-md items-center gap-2 flex ml-2 text-xs">
                        <span class="text-xl leading-none">+</span> Tambah Barang
                    </button>
                </div>
            </div>

            <div class="flex flex-col sm:flex-row gap-4 mb-4">
                <div class="relative flex-1">
                    <span class="absolute inset-y-0 left-0 pl-4 flex items-center text-gray-400">🔍</span>
                    <input v-model="searchQuery" type="text" placeholder="Cari nama barang atau barcode..." class="block w-full pl-11 pr-4 py-2.5 rounded-xl border border-gray-200 shadow-sm focus:border-blue-500 focus:ring-blue-500 text-sm font-medium outline-none">
                </div>
                <div class="w-full sm:w-64 shrink-0">
                    <select v-model="selectedCategory" class="block w-full px-4 py-2.5 rounded-xl border border-gray-200 shadow-sm focus:border-blue-500 focus:ring-blue-500 text-sm font-medium text-gray-700 cursor-pointer bg-white outline-none">
                        <option value="">Semua Kategori</option>
                        <option v-for="cat in categories" :key="cat" :value="cat">{{ cat }}</option>
                    </select>
                </div>
            </div>

            <div class="w-full bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
                <div class="overflow-x-auto w-full">
                    <table class="w-full text-left border-collapse whitespace-nowrap min-w-[800px]">
                        <thead>
                            <tr class="bg-gray-50 text-gray-400 border-b border-gray-100">
                                <th class="p-4 font-black uppercase tracking-widest text-[10px]">Produk & Barcode</th>
                                <th class="p-4 font-black uppercase tracking-widest text-[10px]">Kategori</th>
                                <th class="p-4 font-black uppercase tracking-widest text-[10px] text-right">Modal</th>
                                <th class="p-4 font-black uppercase tracking-widest text-[10px] text-right">Harga Jual</th>
                                <th class="p-4 font-black uppercase tracking-widest text-[10px] text-center">Stok</th>
                                <th class="p-4 font-black uppercase tracking-widest text-[10px] text-center">Aksi</th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-gray-50">
                            <tr v-for="product in products" :key="product.id" class="hover:bg-blue-50/30 transition-colors group">
                                <td class="p-4 flex items-center gap-4">
                                    <div class="w-12 h-12 rounded-xl border border-gray-100 bg-gray-50 flex items-center justify-center text-xl shadow-inner overflow-hidden">
                                        <img v-if="product.gambar" :src="getImageUrl(product.gambar)" class="w-full h-full object-cover">
                                        <span v-else>📦</span>
                                    </div>
                                    <div>
                                        <div class="font-bold text-gray-800 text-sm">{{ product.nama_produk }}</div>
                                        <div class="text-[10px] text-gray-400 font-mono mt-0.5 tracking-tighter">📟 {{ product.sku || '-' }}</div>
                                    </div>
                                </td>
                                <td class="p-4">
                                    <span class="bg-gray-100 text-gray-600 px-3 py-1 rounded-lg font-bold text-[10px] uppercase tracking-wider">{{ product.kategori || 'General' }}</span>
                                </td>
                                <td class="p-4 text-gray-500 font-bold text-right text-sm">Rp {{ (product.harga_modal || 0).toLocaleString('id-ID') }}</td>
                                <td class="p-4 text-blue-700 font-black text-right text-sm">Rp {{ (product.harga_jual || 0).toLocaleString('id-ID') }}</td>
                                <td class="p-4 text-center">
                                    <span class="px-3 py-1 text-[11px] rounded-lg font-black shadow-sm" :class="(product.stok || 0) > 10 ? 'bg-green-50 text-green-600 border border-green-100' : 'bg-red-50 text-red-600 border border-red-100'">
                                        {{ product.stok || 0 }}
                                    </span>
                                </td>
                                <td class="p-4 text-center">
                                    <div class="flex justify-center gap-2">
                                        <button @click="editProduct(product)" class="bg-yellow-50 text-yellow-600 p-2 rounded-lg hover:bg-yellow-100 transition-colors" title="Edit">✏️</button>
                                        <button @click="deleteProduct(product.id)" class="bg-red-50 text-red-600 p-2 rounded-lg hover:bg-red-100 transition-colors" title="Hapus">🗑️</button>
                                    </div>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>

            <div v-if="showFormModal" class="fixed inset-0 bg-black/60 flex items-center justify-center z-[60] p-4 backdrop-blur-sm transition-all">
                <div class="bg-white rounded-3xl shadow-2xl w-full max-w-xl overflow-hidden flex flex-col transform transition-all border-t-8" :class="isEditing ? 'border-yellow-400' : 'border-blue-600'">
                    <div class="p-6">
                        <h2 class="text-2xl font-black text-gray-800 mb-6">{{ isEditing ? 'Edit Produk' : 'Produk Baru' }}</h2>
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                            
                            <div class="md:col-span-2 flex items-center gap-4 mb-2">
                                <div class="w-16 h-16 rounded-xl border-2 border-dashed border-gray-300 flex items-center justify-center bg-gray-50 overflow-hidden shrink-0">
                                    <img v-if="imagePreview" :src="imagePreview" class="w-full h-full object-cover">
                                    <span v-else class="text-2xl text-gray-400">📷</span>
                                </div>
                                <div class="flex-1">
                                    <label class="text-[10px] font-black text-gray-400 uppercase mb-1 block">Foto Produk (Opsional)</label>
                                    <input type="file" @change="onFileChange" accept="image/*" class="w-full text-xs text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-xs file:font-bold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100 cursor-pointer">
                                </div>
                            </div>

                            <div class="md:col-span-2">
                                <label class="text-[10px] font-black text-gray-400 uppercase mb-1 block">Nama Barang</label>
                                <input v-model="form.name" type="text" class="w-full px-4 py-2.5 rounded-xl border border-gray-200 focus:ring-2 focus:ring-blue-500 outline-none font-bold">
                            </div>
                            <div>
                                <label class="text-[10px] font-black text-gray-400 uppercase mb-1 block">SKU / Barcode</label>
                                <input v-model="form.sku" type="text" class="w-full px-4 py-2.5 rounded-xl border border-gray-200 focus:ring-2 focus:ring-blue-500 outline-none font-mono">
                            </div>
                            <div>
                                <label class="text-[10px] font-black text-gray-400 uppercase mb-1 block">Kategori</label>
                                <input list="kategori-list" v-model="form.category" placeholder="Pilih atau ketik baru..." class="w-full px-4 py-2.5 rounded-xl border border-gray-200 focus:ring-2 focus:ring-blue-500 outline-none font-bold bg-white">
                                <datalist id="kategori-list">
                                    <option v-for="cat in categories" :key="cat" :value="cat"></option>
                                </datalist>
                            </div>
                            <div>
                                <label class="text-[10px] font-black text-gray-400 uppercase mb-1 block">Harga Modal</label>
                                <input v-model="form.cost_price" type="number" class="w-full px-4 py-2.5 rounded-xl border border-gray-200 focus:ring-2 focus:ring-blue-500 outline-none font-bold">
                            </div>
                            <div>
                                <label class="text-[10px] font-black text-gray-400 uppercase mb-1 block">Harga Jual</label>
                                <input v-model="form.price" type="number" class="w-full px-4 py-2.5 rounded-xl border border-gray-200 focus:ring-2 focus:ring-blue-500 outline-none font-bold text-blue-600">
                            </div>
                            <div class="md:col-span-2">
                                <label class="text-[10px] font-black text-gray-400 uppercase mb-1 block">Stok Awal</label>
                                <input v-model="form.stock" type="number" class="w-full px-4 py-2.5 rounded-xl border border-gray-200 focus:ring-2 focus:ring-blue-500 outline-none font-black text-center text-lg">
                            </div>
                        </div>
                    </div>
                    <div class="p-4 bg-gray-50 flex gap-3">
                        <button @click="showFormModal = false" class="flex-1 py-3 font-bold text-gray-500 hover:bg-gray-200 rounded-2xl transition-colors">Batal</button>
                        <button @click="submitProduct" class="flex-1 py-3 font-black bg-blue-600 text-white rounded-2xl shadow-lg hover:bg-blue-700 transition-all">Simpan Perubahan</button>
                    </div>
                </div>
            </div>
        </div>
    </Sidebar>
</template>