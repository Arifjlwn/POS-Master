<script setup>
import { ref, computed } from 'vue';
import Sidebar from '../components/Sidebar.vue';
import api from '../api.js';
import Swal from 'sweetalert2';

const notes = ref('Stock Opname Reguler');
const searchQuery = ref('');
const products = ref([]);
const cartSO = ref([]);
const isSubmitting = ref(false);

// --- SEARCH LOGIC ---
const searchProduct = async () => {
    if (searchQuery.value.length < 2) return;
    try {
        // Kita hanya ambil nama dan ID, jangan tampilkan stok sistem di sini!
        const res = await api.get(`/products?search=${searchQuery.value}`);
        products.value = res.data.data;
    } catch (err) { console.error(err); }
};

const addToCart = (product) => {
    const existing = cartSO.value.find(item => item.product_id === product.id);
    if (existing) {
        // Jika scan lagi, tambah 1 (khas barcode scanner)
        existing.actual_qty++;
    } else {
        cartSO.value.unshift({
            product_id: product.id,
            nama_produk: product.nama_produk,
            actual_qty: 1
        });
    }
    searchQuery.value = '';
    products.value = [];
};

const removeItem = (index) => cartSO.value.splice(index, 1);

// --- SUBMIT STOCK OPNAME ---
const submitSO = async () => {
    if (cartSO.value.length === 0) {
        return Swal.fire('Kosong!', 'Input barang yang dihitung dulu Mas.', 'warning');
    }

    const result = await Swal.fire({
        title: 'Selesai Hitung?',
        text: "Pastikan jumlah fisik sudah benar. Stok sistem akan langsung disesuaikan!",
        icon: 'warning',
        showCancelButton: true,
        confirmButtonText: 'Ya, Finalisasi!'
    });

    if (!result.isConfirmed) return;

    isSubmitting.value = true;
    try {
        await api.post('/stock-opname', {
            notes: notes.value,
            items: cartSO.value
        });

        Swal.fire('Berhasil!', 'Stok master telah disesuaikan dengan fisik.', 'success');
        cartSO.value = [];
        notes.value = 'Stock Opname Reguler';
    } catch (err) {
        Swal.fire('Gagal', err.response?.data?.error || 'Terjadi kesalahan', 'error');
    } finally {
        isSubmitting.value = false;
    }
};
</script>

<template>
    <Sidebar>
        <div class="p-6 md:p-10 max-w-5xl mx-auto font-sans">
            <div class="flex items-center justify-between mb-8">
                <div>
                    <h1 class="text-3xl font-black text-slate-900 tracking-tighter uppercase">🔍 Stock Opname</h1>
                    <p class="text-slate-400 font-bold text-[10px] tracking-widest">PENYESUAIAN STOK FISIK VS SISTEM</p>
                </div>
                <div class="bg-amber-100 text-amber-700 px-4 py-2 rounded-2xl text-xs font-black italic">
                    MODE: BLIND COUNTING 🔒
                </div>
            </div>

            <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
                <div class="lg:col-span-2 space-y-6">
                    <div class="bg-white p-6 rounded-[30px] shadow-xl shadow-slate-200/50 border border-white relative">
                        <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-3">Scan Barcode / Cari Barang</label>
                        <input v-model="searchQuery" @input="searchProduct" type="text" placeholder="Masukkan Barcode..." class="w-full p-4 bg-slate-50 border-2 border-transparent focus:border-blue-600 rounded-2xl outline-none font-bold transition-all">
                        
                        <div v-if="products.length" class="absolute left-6 right-6 mt-2 bg-white border border-slate-100 rounded-2xl shadow-2xl z-50 overflow-hidden">
                            <div v-for="p in products" :key="p.id" @click="addToCart(p)" class="p-4 hover:bg-blue-50 cursor-pointer border-b border-slate-50 last:border-0 flex justify-between items-center">
                                <span class="font-black text-slate-700 uppercase">{{ p.nama_produk }}</span>
                                <span class="text-[10px] font-black text-blue-600">PILIH</span>
                            </div>
                        </div>
                    </div>

                    <div class="bg-white rounded-[30px] shadow-xl shadow-slate-200/50 border border-white overflow-hidden">
                        <table class="w-full text-left">
                            <thead class="bg-slate-50 text-[10px] font-black text-slate-400 uppercase tracking-widest">
                                <tr>
                                    <th class="p-5">Barang</th>
                                    <th class="p-5 text-center">Jumlah Fisik (Real)</th>
                                    <th class="p-5 text-center">Aksi</th>
                                </tr>
                            </thead>
                            <tbody class="divide-y divide-slate-50">
                                <tr v-for="(item, index) in cartSO" :key="index" class="hover:bg-slate-50/50 transition-colors">
                                    <td class="p-5 font-black text-slate-700 uppercase">{{ item.nama_produk }}</td>
                                    <td class="p-5">
                                        <input v-model.number="item.actual_qty" type="number" class="w-24 mx-auto block p-3 bg-slate-100 rounded-xl text-center font-black text-blue-600 focus:ring-2 focus:ring-blue-600 outline-none">
                                    </td>
                                    <td class="p-5 text-center">
                                        <button @click="removeItem(index)" class="text-slate-300 hover:text-red-500 transition-colors">✕</button>
                                    </td>
                                </tr>
                                <tr v-if="cartSO.length === 0">
                                    <td colspan="3" class="p-10 text-center text-slate-300 font-bold italic">Belum ada barang yang di-scan.</td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </div>

                <div class="space-y-6">
                    <div class="bg-white p-6 rounded-[30px] shadow-xl shadow-slate-200/50 border border-white">
                        <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-3">Catatan SO</label>
                        <textarea v-model="notes" rows="3" class="w-full p-4 bg-slate-50 border-2 border-transparent focus:border-blue-600 rounded-2xl outline-none font-bold transition-all text-sm"></textarea>
                    </div>

                    <button @click="submitSO" :disabled="isSubmitting" class="w-full bg-slate-900 hover:bg-blue-600 text-white py-6 rounded-[30px] font-black text-lg shadow-xl shadow-slate-200 transition-all active:scale-95 disabled:opacity-50">
                        {{ isSubmitting ? 'PROSES...' : '🚀 FINALISASI SO' }}
                    </button>
                    <p class="text-[9px] text-center text-slate-400 font-bold uppercase leading-relaxed px-4">
                        Perhatian: Finalisasi SO akan langsung mengubah stok di sistem sesuai dengan jumlah fisik yang Anda input.
                    </p>
                </div>
            </div>
        </div>
    </Sidebar>
</template>