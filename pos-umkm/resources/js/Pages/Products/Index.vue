<script setup>
import { ref, watch } from 'vue';
import { Head, useForm, router, Link } from '@inertiajs/vue3';
import AppLayout from '@/Layouts/AppLayout.vue';

const props = defineProps({
    products: Object,
    filters: Object,
    categories: Array,
});

// fileInput untuk Upload
const fileInput = ref(null);

const triggerImport = () => {
    fileInput.value.click();
};

const handleImport = (event) => {
    const file = event.target.files[0];
    if (!file) return;

    if (!confirm('Yakin ingin import data CSV ini? \n\nCatatan: Jika SKU/Barcode sudah ada, sistem akan otomatis mengupdate harga dan stoknya.')) {
        event.target.value = ''; // Reset input
        return;
    }

    router.post('/products/import', {
        file: file,
    }, {
        forceFormData: true, // Wajib true agar file bisa terkirim
        preserveScroll: true,
        onSuccess: () => {
            alert('Yeay! Berhasil mengimport data produk.');
            event.target.value = '';
        },
        onError: () => {
            alert('Gagal import data. Pastikan format file adalah CSV.');
            event.target.value = '';
        }
    });
};

const isEditing = ref(false);
const editId = ref(null);
// Variabel baru untuk mengontrol muncul/hilangnya Form Modal
const showFormModal = ref(false);

// Fungsi Seacrh & Filter
const searchQuery = ref(props.filters.search || '');
const selectedCategory = ref(props.filters.category || '');

let searchTimeOut;
watch([searchQuery, selectedCategory], () => {
    clearTimeout(searchTimeOut);

    searchTimeOut = setTimeout(() => {
        router.get('/products', {
            search: searchQuery.value,
            category: selectedCategory.value
        }, {
            preserveState: true,
            preserveScroll: true,
            replace: true
        });
    }, 300);
});

const form = useForm({
    name: '',
    sku: '',
    category: '',
    cost_price: '',
    price: '',
    stock: '',
    image: null,
});

// Fungsi untuk membuka form modal kosong (Tambah Data)
const openAddModal = () => {
    isEditing.value = false;
    editId.value = null;
    form.reset();
    showFormModal.value = true;
};

const editProduct = (product) => {
    isEditing.value = true;
    editId.value = product.id;
    form.name = product.name;
    form.sku = product.sku || '';
    form.category = product.category || '';
    form.cost_price = product.cost_price;
    form.price = product.price;
    form.stock = product.stock;
    form.image = null;
    showFormModal.value = true; // Buka modal saat klik edit
};

const cancelForm = () => {
    showFormModal.value = false; // Tutup modal
    setTimeout(() => {
        isEditing.value = false;
        editId.value = null;
        form.reset();
    }, 300); // Jeda sedikit biar animasi tutup modalnya mulus
};

const deleteProduct = (id) => {
    if (confirm('Yakin ingin menghapus produk ini?')) {
        router.delete(`/products/${id}`, {
            preserveScroll: true,
            onError: (errors) => {
                if (errors.message) {
                    alert(errors.message);
                }
            }
        });
    }
};

const submitProduct = () => {
    if (isEditing.value) {
        router.post(`/products/${editId.value}`, {
            _method: 'put',
            name: form.name,
            sku: form.sku,
            category: form.category,
            cost_price: form.cost_price,
            price: form.price,
            stock: form.stock,
            image: form.image
        }, {
            preserveScroll: true,
            onSuccess: () => {
                cancelForm();
                alert('Data Produk Berhasil Diperbarui!');
            }
        });
    } else {
        form.post('/products', {
            preserveScroll: true,
            onSuccess: () => {
                cancelForm();
                alert('Produk Baru Berhasil Ditambahkan!');
            },
        });
    }
};
</script>

<template>
    <Head title="Master Stock" />

    <AppLayout>
        <div class="p-4 md:p-6 w-full max-w-full mx-auto">

            <div class="flex justify-between items-center mb-6">
                <div>
                    <h1 class="text-2xl font-black text-gray-800 tracking-tight">📦 Master Produk</h1>
                    <p class="text-sm text-gray-500 mt-1">Kelola daftar barang, harga, dan stok toko Anda.</p>
                </div>

                <div class="hidden sm:flex gap-2">
                    <!-- Tombol Export (Langsung hit ke endpoint export) -->
                    <a href="/products/export" target="_blank" class="bg-green-100 hover:bg-green-200 text-green-700 px-4 py-2.5 rounded-xl font-bold transition-all shadow-sm items-center gap-2 border border-green-200 flex">
                        📤 Export CSV
                    </a>

                    <!-- Tombol Import -->
                    <button @click="triggerImport" class="bg-orange-100 hover:bg-orange-200 text-orange-700 px-4 py-2.5 rounded-xl font-bold transition-all shadow-sm items-center gap-2 border border-orange-200 flex">
                        📥 Import CSV
                    </button>
                    <!-- Input File Tersembunyi -->
                    <input type="file" ref="fileInput" class="hidden" accept=".csv" @change="handleImport">

                    <!-- Tombol Tambah Barang yang sudah ada -->
                    <button @click="openAddModal" class="bg-blue-600 hover:bg-blue-700 text-white px-5 py-2.5 rounded-xl font-bold transition-all shadow-md items-center gap-2 flex ml-2">
                        <span class="text-xl leading-none">+</span> Tambah Barang
                    </button>
                </div>
            </div>

            <div class="flex flex-col sm:flex-row gap-4 mb-4">
                <div class="relative flex-1">
                    <span class="absolute inset-y-0 left-0 pl-4 flex items-center text-gray-400">
                        🔍
                    </span>
                    <input
                        type="text"
                        v-model="searchQuery"
                        placeholder="Cari nama barang atau scan barcode..."
                        class="block w-full pl-11 pr-4 py-2.5 rounded-xl border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 text-sm font-medium"
                    >
                </div>

                <div class="w-full sm:w-64 shrink-0 relative">
                    <span class="absolute inset-y-0 left-0 pl-3 flex items-center text-gray-400">
                        📁
                    </span>
                    <select
                        v-model="selectedCategory"
                        class="block w-full pl-10 pr-10 py-2.5 rounded-xl border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 text-sm font-medium text-gray-700 cursor-pointer appearance-none bg-white"
                    >
                        <option value="">Semua Kategori</option>
                        <option v-for="cat in categories" :key="cat" :value="cat">
                            {{ cat }}
                        </option>
                    </select>
                </div>
            </div>

            <div class="w-full bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
                <div class="overflow-x-auto w-full">
                    <table class="w-full text-left border-collapse whitespace-nowrap min-w-[800px]">
                        <thead>
                            <tr class="bg-gray-50 text-gray-700 border-b">
                                <th class="p-4 font-bold uppercase tracking-wider text-xs">Produk & Barcode</th>
                                <th class="p-4 font-bold uppercase tracking-wider text-xs">Kategori</th>
                                <th class="p-4 font-bold uppercase tracking-wider text-xs text-right">Modal (HPP)</th>
                                <th class="p-4 font-bold uppercase tracking-wider text-xs text-right">Harga Jual</th>
                                <th class="p-4 font-bold uppercase tracking-wider text-xs text-center">Stok</th>
                                <th class="p-4 font-bold uppercase tracking-wider text-xs text-center">Aksi</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="product in products.data" :key="product.id" class="border-b border-gray-50 hover:bg-blue-50/30 transition-colors">
                                <td class="p-4 flex items-center gap-4">
                                    <img :src="product.image || 'https://placehold.co/100x100?text=No+Image'" class="w-12 h-12 object-contain p-1 rounded-lg border border-gray-200 bg-white shadow-sm shrink-0">
                                    <div>
                                        <div class="font-bold text-gray-800">{{ product.name }}</div>
                                        <div class="text-xs text-gray-400 font-mono mt-0.5 flex items-center gap-1">
                                            <span class="text-[10px]">📟</span> {{ product.sku || 'Belum ada Barcode' }}
                                        </div>
                                    </div>
                                </td>
                                <td class="p-4 text-sm text-gray-600">
                                    <span v-if="product.category" class="bg-gray-100 px-3 py-1 rounded-md font-medium">{{ product.category }}</span>
                                    <span v-else class="text-gray-400 italic">-</span>
                                </td>
                                <td class="p-4 text-gray-500 font-medium text-right text-sm">Rp {{ product.cost_price?.toLocaleString('id-ID') || 0 }}</td>
                                <td class="p-4 text-blue-600 font-black text-right text-sm">Rp {{ product.price.toLocaleString('id-ID') }}</td>
                                <td class="p-4 text-center">
                                    <span class="px-3 py-1 text-xs rounded-full font-bold shadow-sm" :class="product.stock > 10 ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'">
                                        {{ product.stock }}
                                    </span>
                                </td>
                                <td class="p-4">
                                    <div class="flex justify-center gap-2">
                                        <button @click="editProduct(product)" class="bg-yellow-100 text-yellow-700 px-4 py-1.5 rounded-lg shadow-sm hover:bg-yellow-200 text-xs font-bold transition">Edit</button>
                                        <button @click="deleteProduct(product.id)" class="bg-red-100 text-red-700 px-4 py-1.5 rounded-lg shadow-sm hover:bg-red-200 text-xs font-bold transition">Hapus</button>
                                    </div>
                                </td>
                            </tr>
                            <tr v-if="products.data.length === 0">
                                <td colspan="6" class="p-12 text-center">
                                    <div class="text-4xl mb-3">📦</div>
                                    <div class="text-gray-500 font-medium">Belum ada data produk.</div>
                                    <div class="text-sm text-gray-400 mt-1">Klik tombol tambah di pojok kanan bawah untuk memulai.</div>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                    <div v-if="products.data.length > 0" class="p-4 border-t border-gray-100 flex flex-col sm:flex-row items-center justify-between gap-4 bg-gray-50">
                    <div class="text-sm text-gray-500 font-medium">
                        Menampilkan <span class="font-bold text-gray-800">{{ products.from }}</span> sampai <span class="font-bold text-gray-800">{{ products.to }}</span> dari <span class="font-bold text-blue-600">{{ products.total }}</span> data
                    </div>
                    <div class="flex gap-1 flex-wrap justify-center">
                        <Link
                            v-for="(link, index) in products.links"
                            :key="index"
                            :href="link.url || '#'"
                            class="px-3.5 py-1.5 rounded-lg text-sm font-bold transition-colors shadow-sm"
                            :class="[
                                link.active ? 'bg-blue-600 text-white border border-blue-600' : 'bg-white text-gray-600 border border-gray-200 hover:bg-gray-100',
                                !link.url ? 'opacity-50 cursor-not-allowed bg-gray-50' : ''
                            ]"
                            v-html="link.label"
                        ></Link>
                    </div>
                </div>
                </div>
            </div>

            <div v-if="showFormModal" class="fixed inset-0 bg-black/60 flex items-center justify-center z-50 p-4 backdrop-blur-sm transition-opacity">
                <div class="bg-white rounded-2xl shadow-2xl w-full max-w-2xl overflow-hidden flex flex-col max-h-[90vh]">

                    <div class="p-5 border-b border-gray-100 flex justify-between items-center bg-gray-50 shrink-0" :class="isEditing ? 'border-t-4 border-t-yellow-400' : 'border-t-4 border-t-blue-600'">
                        <h2 class="text-xl font-black text-gray-800">
                            {{ isEditing ? '✏️ Edit Data Produk' : '➕ Tambah Produk Baru' }}
                        </h2>
                        <button @click="cancelForm" class="text-gray-400 hover:text-red-500 transition-colors p-1 rounded-lg hover:bg-red-50">
                            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M6 18L18 6M6 6l12 12" /></svg>
                        </button>
                    </div>

                    <div class="p-6 overflow-y-auto flex-1">
                        <form @submit.prevent="submitProduct" class="space-y-5">
                            <div>
                                <label class="block text-sm font-bold text-gray-700 mb-1">Nama Barang <span class="text-red-500">*</span></label>
                                <input type="text" v-model="form.name" class="block w-full rounded-xl border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500" required>
                            </div>

                            <div class="flex flex-col sm:flex-row gap-5">
                                <div class="w-full sm:w-1/2">
                                    <label class="block text-sm font-bold text-gray-700 mb-1">SKU / Barcode</label>
                                    <div class="relative">
                                        <span class="absolute inset-y-0 left-0 pl-3 flex items-center text-gray-400">📟</span>
                                        <input type="text" v-model="form.sku" placeholder="Scan barcode..." class="block w-full pl-9 rounded-xl border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500">
                                    </div>
                                </div>
                                <div class="w-full sm:w-1/2">
                                    <label class="block text-sm font-bold text-gray-700 mb-1">Kategori</label>
                                    <input type="text" v-model="form.category" placeholder="Cth: Minuman, Snack" class="block w-full rounded-xl border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500">
                                </div>
                            </div>

                            <div class="flex flex-col sm:flex-row gap-5">
                                <div class="w-full sm:w-1/2">
                                    <label class="block text-sm font-bold text-gray-700 mb-1">Harga Modal (HPP) <span class="text-red-500">*</span></label>
                                    <div class="relative">
                                        <span class="absolute inset-y-0 left-0 pl-3 flex items-center text-gray-500 font-bold">Rp</span>
                                        <input type="number" v-model="form.cost_price" class="block w-full pl-10 rounded-xl border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500" required>
                                    </div>
                                </div>
                                <div class="w-full sm:w-1/2">
                                    <label class="block text-sm font-bold text-gray-700 mb-1">Harga Jual <span class="text-red-500">*</span></label>
                                    <div class="relative">
                                        <span class="absolute inset-y-0 left-0 pl-3 flex items-center text-blue-600 font-bold">Rp</span>
                                        <input type="number" v-model="form.price" class="block w-full pl-10 rounded-xl border-blue-200 bg-blue-50 shadow-sm focus:border-blue-500 focus:ring-blue-500" required>
                                    </div>
                                </div>
                            </div>

                            <div class="flex flex-col sm:flex-row gap-5">
                                <div class="w-full sm:w-1/3">
                                    <label class="block text-sm font-bold text-gray-700 mb-1">Sisa Stok <span class="text-red-500">*</span></label>
                                    <input type="number" v-model="form.stock" class="block w-full rounded-xl border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 text-center font-bold" required>
                                </div>
                                <div class="w-full sm:w-2/3">
                                    <label class="block text-sm font-bold text-gray-700 mb-1">Foto {{ isEditing ? '(Kosongkan jika tidak diganti)' : '' }}</label>
                                    <input type="file" @input="form.image = $event.target.files[0]" class="block w-full text-sm text-gray-500 file:mr-4 file:py-2.5 file:px-4 file:rounded-xl file:border-0 file:text-sm file:font-bold file:bg-gray-100 file:text-gray-700 hover:file:bg-gray-200 border border-gray-200 rounded-xl" accept="image/*">
                                </div>
                            </div>
                        </form>
                    </div>

                    <div class="p-5 border-t border-gray-100 bg-white flex gap-3 shrink-0">
                        <button type="button" @click="cancelForm" class="flex-1 bg-gray-100 text-gray-700 font-bold py-3 px-4 rounded-xl hover:bg-gray-200 transition-colors">
                            Batal
                        </button>
                        <button @click="submitProduct" :disabled="form.processing" class="flex-1 bg-blue-600 text-white font-bold py-3 px-4 rounded-xl hover:bg-blue-700 transition-colors shadow-md disabled:opacity-50">
                            {{ form.processing ? 'Menyimpan...' : 'Simpan Data' }}
                        </button>
                    </div>

                </div>
            </div>

        </div>
    </AppLayout>
</template>
