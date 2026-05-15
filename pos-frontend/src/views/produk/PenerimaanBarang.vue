<script setup>
import { ref, computed } from 'vue';
import Sidebar from '../../components/Sidebar.vue';
import api from '../../api.js';
import Swal from 'sweetalert2';

// --- STATE USER & ROLE ---
const getPayloadFromToken = () => {
    const token = localStorage.getItem('token');
    const role = localStorage.getItem('role') || 'staff';
    if (!token) return { role: 'staff' };
    try {
        return { role: role.toLowerCase() };
    } catch (e) { return { role: 'staff' }; }
};
const currentUser = ref(getPayloadFromToken());

// 🚀 CEK APAKAH USER BOLEH LIHAT HARGA?
const canSeePrice = computed(() => {
    return currentUser.value.role === 'owner';
});

// --- STATE ---
const supplierName = ref('');
const noFaktur = ref('');
const searchQuery = ref('');
const products = ref([]); 
const cartLPB = ref([]); 
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
            harga_jual_saat_ini: product.harga_jual || 0,
            // 🚀 HARGA MODAL ASLI (Cuma dipake kalau dia Owner)
            harga_modal_database: product.harga_modal || 0,
            harga_beli_input: product.satuan_besar ? (product.harga_modal * product.isi_per_besar) : product.harga_modal, 
            satuan_dasar: product.satuan_dasar || 'PCS',
            satuan_besar: product.satuan_besar || null,
            isi_per_besar: product.isi_per_besar || 1,
            has_satuan_besar: !!product.satuan_besar && product.isi_per_besar > 1,
            jenis_satuan: product.satuan_besar ? 'BESAR' : 'DASAR' 
        });
    }
    searchQuery.value = '';
    products.value = [];
};

const hitungTotalStok = (item) => {
    return item.jenis_satuan === 'BESAR' ? (item.qty_masuk * item.isi_per_besar) : item.qty_masuk;
};

const hitungModalPerPcs = (item) => {
    // Kalau Karyawan, kita pake harga modal lama dari DB biar nggak ngerusak data
    if (!canSeePrice.value) return item.harga_modal_database;
    
    if (item.jenis_satuan === 'BESAR') {
        return Math.round(item.harga_beli_input / item.isi_per_besar);
    }
    return item.harga_beli_input;
};

const removeItem = (index) => cartLPB.value.splice(index, 1);

const submitLPB = async () => {
    if (!supplierName.value || cartLPB.value.length === 0) {
        return Swal.fire('Oops!', 'Lengkapi data!', 'warning');
    }

    const adaYangRugi = canSeePrice.value && cartLPB.value.some(item => hitungModalPerPcs(item) >= item.harga_jual_saat_ini);
    
    const result = await Swal.fire({
        title: 'Konfirmasi LPB',
        text: adaYangRugi ? "⚠️ Modal baru rugi! Lanjut?" : "Simpan penerimaan barang?",
        icon: 'question',
        showCancelButton: true,
        confirmButtonText: 'Ya, Simpan!'
    });

    if (!result.isConfirmed) return;

    isSubmitting.value = true;
    try {
        const payloadItems = cartLPB.value.map(item => ({
            product_id: item.product_id,
            qty_masuk: hitungTotalStok(item),
            // Kalau bukan owner, kirim harga_modal_database (biar modal nggak berubah)
            harga_beli: canSeePrice.value ? hitungModalPerPcs(item) : item.harga_modal_database
        }));

        await api.post('/purchases', {
            supplier_name: supplierName.value,
            no_faktur: noFaktur.value,
            items: payloadItems
        });

        Swal.fire('Berhasil!', 'Stok telah diperbarui.', 'success');
        supplierName.value = ''; noFaktur.value = ''; cartLPB.value = [];
    } catch (err) {
        Swal.fire('Gagal!', 'Error system', 'error');
    } finally {
        isSubmitting.value = false;
    }
};
</script>

<template>
    <Sidebar>
        <div class="p-4 md:p-8 lg:p-10 max-w-7xl mx-auto font-sans bg-[#f8fafc] min-h-screen">
            
            <div class="bg-gradient-to-br from-slate-900 to-indigo-900 rounded-[40px] p-8 md:p-10 mb-10 text-white shadow-2xl flex flex-col md:flex-row items-center justify-between relative overflow-hidden border border-white/10">
                <div class="z-10">
                    <h1 class="text-3xl font-black tracking-tighter mb-2 uppercase italic leading-none">Inbound <span class="text-blue-400">Logistics</span></h1>
                    <p class="text-blue-200 font-bold text-[10px] uppercase tracking-[0.2em] flex items-center gap-2">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" /></svg>
                        Penerimaan Stok Barang
                    </p>
                </div>
                <div class="z-10 mt-4 md:mt-0 px-4 py-2 bg-white/10 rounded-2xl border border-white/20 backdrop-blur-md">
                    <span class="text-[9px] font-black uppercase tracking-widest text-blue-300">Mode Akses: {{ currentUser.role }}</span>
                </div>
            </div>

            <div class="grid grid-cols-1 xl:grid-cols-4 gap-8">
                <div class="xl:col-span-1 space-y-6">
                    <div class="bg-white p-8 rounded-[32px] shadow-xl shadow-slate-200/50 border border-slate-100">
                        <h3 class="font-black text-slate-800 text-[10px] uppercase tracking-[0.2em] mb-6 flex items-center gap-2">Data Pengiriman</h3>
                        <div class="space-y-4">
                            <input v-model="supplierName" type="text" placeholder="Nama Supplier..." class="w-full p-4 bg-slate-50 border-2 border-slate-100 rounded-2xl outline-none font-bold text-sm focus:border-blue-500">
                            <input v-model="noFaktur" type="text" placeholder="No Faktur..." class="w-full p-4 bg-slate-50 border-2 border-slate-100 rounded-2xl outline-none font-bold text-sm focus:border-blue-500 uppercase">
                        </div>
                    </div>
                </div>

                <div class="xl:col-span-3 space-y-6">
                    <div class="bg-white p-6 rounded-[32px] shadow-sm border border-slate-100 relative">
                        <input v-model="searchQuery" @input="searchProduct" type="text" placeholder="Scan Barcode / Cari Produk..." class="w-full p-5 bg-slate-50 border-2 border-transparent rounded-[24px] outline-none font-bold text-sm focus:border-blue-500 focus:bg-white transition-all">
                        
                        <div v-if="products.length" class="absolute left-6 right-6 mt-3 bg-white border border-slate-100 rounded-[28px] shadow-2xl z-[100] overflow-hidden">
                            <div v-for="p in products" :key="p.id" @click="addToCart(p)" class="p-5 hover:bg-blue-50 cursor-pointer border-b last:border-0 flex justify-between items-center group transition-all">
                                <div class="flex items-center gap-4">
                                    <div class="w-10 h-10 bg-slate-100 rounded-xl flex items-center justify-center text-xl shadow-inner group-hover:bg-blue-100 transition-colors">📦</div>
                                    <div class="font-black text-slate-800 uppercase text-sm group-hover:text-blue-600 transition-colors">{{ p.nama_produk }}</div>
                                </div>
                                <span class="text-[9px] bg-slate-100 text-slate-500 px-3 py-1.5 rounded-lg font-black border border-slate-200">STOK: {{ p.stok }}</span>
                            </div>
                        </div>
                    </div>

                    <div class="bg-white rounded-[40px] shadow-xl shadow-slate-200/50 border border-slate-100 overflow-hidden">
                        <div class="overflow-x-auto custom-scrollbar">
                            <table class="w-full text-left whitespace-nowrap border-collapse">
                                <thead>
                                    <tr class="bg-slate-50/50 text-[10px] font-black text-slate-400 uppercase tracking-[0.2em] border-b border-slate-100">
                                        <th class="p-6">Informasi Barang</th>
                                        <th class="p-6 text-center">Qty Masuk</th>
                                        <th class="p-6 text-center bg-blue-50/30 text-blue-600">Total Pcs</th>
                                        <th v-if="canSeePrice" class="p-6 text-right">Harga Nota</th>
                                        <th v-if="canSeePrice" class="p-6 text-right">Modal Baru</th>
                                        <th class="p-6 text-center">Aksi</th>
                                    </tr>
                                </thead>
                                <tbody class="divide-y divide-slate-50">
                                    <tr v-for="(item, index) in cartLPB" :key="item.product_id" 
                                        :class="(canSeePrice && hitungModalPerPcs(item) > item.harga_jual_saat_ini) ? 'bg-red-50/50' : 'hover:bg-slate-50/50'">
                                        
                                        <td class="p-6">
                                            <div class="font-black text-slate-800 text-sm uppercase tracking-tight">{{ item.nama_produk }}</div>
                                            <div v-if="canSeePrice && hitungModalPerPcs(item) > item.harga_jual_saat_ini" class="mt-1 text-[8px] font-black text-red-600 animate-pulse uppercase">⚠️ Potensi Rugi!</div>
                                        </td>

                                        <td class="p-6 text-center">
                                            <div class="flex items-center justify-center gap-2">
                                                <input v-model.number="item.qty_masuk" type="number" class="w-16 p-3 bg-white border-2 border-slate-100 rounded-xl text-center font-black text-blue-600 outline-none focus:border-blue-500">
                                                <select v-model="item.jenis_satuan" class="p-3 bg-white border-2 border-slate-100 rounded-xl font-black text-[10px] uppercase cursor-pointer">
                                                    <option value="DASAR">{{ item.satuan_dasar }}</option>
                                                    <option v-if="item.has_satuan_besar" value="BESAR">{{ item.satuan_besar }}</option>
                                                </select>
                                            </div>
                                        </td>

                                        <td class="p-6 text-center bg-blue-50/20">
                                            <div class="flex flex-col items-center">
                                                <span class="text-lg font-black text-blue-700 leading-none">{{ hitungTotalStok(item) }}</span>
                                                <span class="text-[8px] font-black text-blue-400 uppercase mt-1 tracking-widest">{{ item.satuan_dasar }}</span>
                                            </div>
                                        </td>

                                        <td v-if="canSeePrice" class="p-6">
                                            <div class="relative">
                                                <span class="absolute inset-y-0 left-0 pl-3 flex items-center text-[10px] font-black text-slate-300 italic">Rp</span>
                                                <input v-model.number="item.harga_beli_input" type="number" class="w-full pl-8 pr-3 py-3 bg-white border-2 border-slate-100 rounded-xl text-right font-black text-slate-700 outline-none focus:border-blue-500 transition-all">
                                            </div>
                                        </td>

                                        <td v-if="canSeePrice" class="p-6 text-right">
                                            <div class="text-base font-black tracking-tight text-emerald-600">
                                                Rp {{ hitungModalPerPcs(item).toLocaleString('id-ID') }}
                                            </div>
                                        </td>

                                        <td class="p-6 text-center">
                                            <button @click="removeItem(index)" class="p-3 bg-slate-50 hover:bg-red-50 rounded-2xl transition-all group">
                                                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-300 group-hover:text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h14" /></svg>
                                            </button>
                                        </td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                    </div>

                    <button @click="submitLPB" :disabled="isSubmitting" class="w-full bg-blue-600 hover:bg-slate-900 text-white py-6 rounded-[30px] font-black text-xs uppercase tracking-[0.3em] shadow-2xl transition-all active:scale-[0.98] disabled:opacity-50 flex items-center justify-center gap-4">
                        <template v-if="isSubmitting"><div class="w-5 h-5 border-4 border-white/30 border-t-white rounded-full animate-spin"></div>Sinkronisasi Stok...</template>
                        <template v-else>Posting Penerimaan Barang</template>
                    </button>
                </div>
            </div>
        </div>
    </Sidebar>
</template>