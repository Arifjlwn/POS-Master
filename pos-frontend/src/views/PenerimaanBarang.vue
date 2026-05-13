<script setup>
import { ref, computed } from 'vue';
import Sidebar from '../components/Sidebar.vue';
import api from '../api.js';
import Swal from 'sweetalert2';

// --- STATE ---
const supplierName = ref('');
const noFaktur = ref('');
const searchQuery = ref('');
const products = ref([]); // Data produk dari master untuk pencarian
const cartLPB = ref([]); // Barang yang sedang di-input
const isSubmitting = ref(false);

// --- SEARCH LOGIC ---
const searchProduct = async () => {
    if (searchQuery.value.length < 2) return;
    try {
        const res = await api.get(`/products?search=${searchQuery.value}`);
        products.value = res.data.data;
    } catch (err) { console.error(err); }
};

const addToCart = (product) => {
    const existing = cartLPB.value.find(item => item.product_id === product.id);
    if (existing) {
        existing.qty_masuk++;
    } else {
        cartLPB.value.push({
            product_id: product.id,
            nama_produk: product.nama_produk,
            qty_masuk: 1,
            harga_beli: product.harga_beli || 0 // Default harga beli dari master
        });
    }
    searchQuery.value = '';
    products.value = [];
};

const removeItem = (index) => cartLPB.value.splice(index, 1);

// --- SUBMIT LPB ---
const submitLPB = async () => {
    if (!supplierName.value || cartLPB.value.length === 0) {
        return Swal.fire('Oops!', 'Lengkapi data supplier dan barang!', 'warning');
    }

    const result = await Swal.fire({
        title: 'Konfirmasi LPB',
        text: "Apakah data barang masuk sudah sesuai?",
        icon: 'question',
        showCancelButton: true,
        confirmButtonText: 'Ya, Simpan!'
    });

    if (!result.isConfirmed) return;

    isSubmitting.value = true;
    try {
        await api.post('/purchases', {
            supplier_name: supplierName.value,
            no_faktur: noFaktur.value,
            items: cartLPB.value
        });

        Swal.fire('Berhasil!', 'Stok barang telah diperbarui.', 'success');
        // Reset Form
        supplierName.value = '';
        noFaktur.value = '';
        cartLPB.value = [];
    } catch (err) {
        Swal.fire('Gagal', err.response?.data?.error || 'Terjadi kesalahan', 'error');
    } finally {
        isSubmitting.value = false;
    }
};
</script>

<template>
    <Sidebar>
        <div class="p-6 md:p-8 max-w-6xl mx-auto font-sans">
            <h1 class="text-2xl font-black text-gray-800 mb-6 uppercase tracking-tight">📦 Penerimaan Barang (LPB)</h1>

            <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
                <div class="lg:col-span-1 space-y-6">
                    <div class="bg-white p-6 rounded-3xl shadow-sm border border-gray-100">
                        <label class="block text-xs font-black text-gray-400 uppercase tracking-widest mb-2">Informasi Faktur</label>
                        <input v-model="supplierName" type="text" placeholder="Nama Supplier" class="w-full p-4 mb-4 bg-gray-50 border-2 border-transparent focus:border-blue-500 rounded-2xl outline-none font-bold transition-all">
                        <input v-model="noFaktur" type="text" placeholder="Nomor Faktur/Surat Jalan" class="w-full p-4 bg-gray-50 border-2 border-transparent focus:border-blue-500 rounded-2xl outline-none font-bold transition-all">
                    </div>

                    <div class="bg-blue-600 p-6 rounded-3xl shadow-lg text-white">
                        <div class="text-xs font-black opacity-60 uppercase mb-1">Total Item Masuk</div>
                        <div class="text-4xl font-black">{{ cartLPB.length }} Jenis</div>
                    </div>
                </div>

                <div class="lg:col-span-2 space-y-6">
                    <div class="bg-white p-6 rounded-3xl shadow-sm border border-gray-100 relative">
                        <label class="block text-xs font-black text-gray-400 uppercase tracking-widest mb-2">Cari Produk</label>
                        <input v-model="searchQuery" @input="searchProduct" type="text" placeholder="Scan Barcode atau Ketik Nama Barang..." class="w-full p-4 bg-gray-50 border-2 border-transparent focus:border-blue-500 rounded-2xl outline-none font-bold transition-all">
                        
                        <div v-if="products.length" class="absolute left-6 right-6 mt-2 bg-white border border-gray-100 rounded-2xl shadow-2xl z-50 overflow-hidden">
                            <div v-for="p in products" :key="p.id" @click="addToCart(p)" class="p-4 hover:bg-blue-50 cursor-pointer border-b border-gray-50 last:border-0 flex justify-between">
                                <span class="font-bold uppercase text-gray-700">{{ p.nama_produk }}</span>
                                <span class="text-xs font-black text-blue-600 bg-blue-50 px-2 py-1 rounded">STOK: {{ p.stok }}</span>
                            </div>
                        </div>
                    </div>

                    <div class="bg-white rounded-3xl shadow-sm border border-gray-100 overflow-hidden">
                        <table class="w-full text-left border-collapse">
                            <thead>
                                <tr class="bg-gray-50 text-[10px] font-black text-gray-400 uppercase tracking-widest">
                                    <th class="p-5">Nama Barang</th>
                                    <th class="p-5 text-center">Qty Masuk</th>
                                    <th class="p-5 text-right">Harga Beli (Satuan)</th>
                                    <th class="p-5 text-center">Aksi</th>
                                </tr>
                            </thead>
                            <tbody class="divide-y divide-gray-50">
                                <tr v-if="cartLPB.length === 0">
                                    <td colspan="4" class="p-10 text-center text-gray-400 italic font-medium">Belum ada barang dipilih.</td>
                                </tr>
                                <tr v-for="(item, index) in cartLPB" :key="item.product_id" class="hover:bg-gray-50/50">
                                    <td class="p-5 font-black text-gray-700 uppercase text-sm">{{ item.nama_produk }}</td>
                                    <td class="p-5">
                                        <input v-model.number="item.qty_masuk" type="number" class="w-20 mx-auto block p-2 bg-gray-100 rounded-lg text-center font-black text-blue-600">
                                    </td>
                                    <td class="p-5">
                                        <input v-model.number="item.harga_beli" type="number" class="w-32 ml-auto block p-2 bg-gray-100 rounded-lg text-right font-black text-green-600">
                                    </td>
                                    <td class="p-5 text-center">
                                        <button @click="removeItem(index)" class="text-red-400 hover:text-red-600 transition-colors">✕</button>
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </div>

                    <button @click="submitLPB" :disabled="isSubmitting" class="w-full bg-blue-600 hover:bg-blue-700 text-white py-5 rounded-3xl font-black text-xl shadow-xl transition-all active:scale-95 disabled:opacity-50">
                        {{ isSubmitting ? 'MEMPROSES...' : '🔥 SIMPAN PENERIMAAN BARANG' }}
                    </button>
                </div>
            </div>
        </div>
    </Sidebar>
</template>