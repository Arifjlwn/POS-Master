<script setup>
import { ref, onMounted, watch } from 'vue';
import api from '../api';
import Sidebar from '../components/Sidebar.vue';

// State Data
const riwayat = ref([]);
const isLoading = ref(true);
const tanggalDipilih = ref(new Date().toISOString().split('T')[0]); // Default hari ini

// State Modal Struk
const showReceipt = ref(false);
const selectedTrx = ref(null);

const formatRupiah = (angka) => {
    return new Intl.NumberFormat('id-ID').format(angka);
};

// Fungsi Tarik Data dari Golang
const fetchRiwayat = async () => {
    isLoading.value = true;
    try {
        const response = await api.get('/transactions', {
            params: { tanggal: tanggalDipilih.value }
        });
        riwayat.value = response.data.data || [];
    } catch (error) {
        console.error("Gagal menarik riwayat transaksi:", error);
    } finally {
        isLoading.value = false;
    }
};

onMounted(() => {
    fetchRiwayat();
});

// Otomatis tarik data ulang kalau tanggalnya diganti bos
watch(tanggalDipilih, () => {
    fetchRiwayat();
});

// Buka Struk
const openReceipt = (trx) => {
    selectedTrx.value = trx;
    showReceipt.value = true;
};

// Print
const printReceipt = () => {
    window.print();
};
</script>

<template>
    <Sidebar>
        <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 font-sans">
            <div class="mb-8 flex flex-col md:flex-row md:justify-between md:items-end gap-4">
                <div>
                    <h1 class="text-3xl font-black text-gray-800 tracking-tight">Riwayat Transaksi</h1>
                    <p class="text-gray-500 text-sm mt-1 font-medium">Lihat daftar transaksi, cek detail struk, atau cetak ulang nota.</p>
                </div>
                <div class="flex items-center gap-3">
                    <span class="text-sm font-bold text-gray-500">Pilih Tanggal:</span>
                    <input type="date" v-model="tanggalDipilih" class="px-4 py-2 bg-white border border-gray-200 rounded-xl text-sm font-bold text-gray-700 shadow-sm focus:ring-blue-500 focus:border-blue-500 cursor-pointer outline-none transition-all">
                </div>
            </div>

            <div class="bg-white rounded-3xl shadow-sm border border-gray-200 overflow-hidden">
                <div class="p-5 border-b border-gray-100 flex items-center gap-2 bg-gray-50/50">
                    <span class="text-xl">📜</span>
                    <h2 class="text-lg font-black text-gray-800">Daftar Transaksi</h2>
                </div>

                <div class="overflow-x-auto">
                    <table class="w-full text-left whitespace-nowrap">
                        <thead class="bg-white border-b border-gray-100">
                            <tr>
                                <th class="px-6 py-4 text-xs font-black text-gray-400 uppercase tracking-widest">No. Invoice</th>
                                <th class="px-6 py-4 text-xs font-black text-gray-400 uppercase tracking-widest">Waktu</th>
                                <th class="px-6 py-4 text-xs font-black text-gray-400 uppercase tracking-widest">Kasir</th>
                                <th class="px-6 py-4 text-xs font-black text-gray-400 uppercase tracking-widest text-right">Total Transaksi</th>
                                <th class="px-6 py-4 text-xs font-black text-gray-400 uppercase tracking-widest text-center">Aksi</th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-gray-50">
                            <tr v-if="isLoading">
                                <td colspan="5" class="px-6 py-12 text-center text-gray-400 font-medium">Sedang memuat data transaksi...</td>
                            </tr>
                            <tr v-else-if="riwayat.length === 0">
                                <td colspan="5" class="px-6 py-12 text-center text-gray-400 font-medium italic">Tidak ada transaksi pada tanggal ini.</td>
                            </tr>
                            <tr v-for="trx in riwayat" :key="trx.id" class="hover:bg-blue-50/30 transition-colors">
                                <td class="px-6 py-4">
                                    <div class="font-mono font-bold text-blue-700">{{ trx.no_invoice }}</div>
                                </td>
                                <td class="px-6 py-4 text-sm font-bold text-gray-600">
                                    {{ new Date(trx.created_at).toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' }) }} WIB
                                </td>
                                <td class="px-6 py-4">
                                    <span class="bg-gray-100 text-gray-700 text-xs font-bold px-2.5 py-1 rounded-md">{{ trx.User?.name || 'Kasir' }}</span>
                                </td>
                                <td class="px-6 py-4 text-right">
                                    <div class="font-black text-gray-800">Rp {{ formatRupiah(trx.total_harga) }}</div>
                                </td>
                                <td class="px-6 py-4 text-center">
                                    <button @click="openReceipt(trx)" class="bg-blue-50 text-blue-600 hover:bg-blue-600 hover:text-white px-4 py-2 rounded-xl text-xs font-black transition-all border border-blue-100 active:scale-95 shadow-sm">
                                        👁️ Lihat Struk
                                    </button>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </main>

        <div v-if="showReceipt" class="fixed inset-0 bg-gray-900/80 flex items-center justify-center z-50 p-4 backdrop-blur-sm no-print">
            <div class="bg-gray-200 p-4 rounded-2xl shadow-2xl w-full max-w-sm overflow-hidden border-t-8 border-gray-800 flex flex-col max-h-[90vh]">
                
                <div class="overflow-y-auto custom-scrollbar bg-white p-4 mx-auto" id="print-area" style="width: 58mm;">
                    <div class="text-center mb-3 font-mono leading-tight">
                        <h2 class="font-black text-sm mb-1">INDO UMKM</h2>
                        <p class="font-medium text-[11px]">JL. KEBON KOSONG NO 56 F</p>
                    </div>

                    <div class="text-center my-2 font-bold tracking-widest border-y border-black py-1 font-mono text-[11px]">
                        <p>S T R U K   B E L A N J A</p>
                    </div>

                    <div class="mb-2 text-[10px] font-bold font-mono">
                        <p>{{ new Date(selectedTrx.created_at).toLocaleString('id-ID') }} / KASIR: {{ selectedTrx.User?.name.split(' ')[0] || 'KASIR' }}</p>
                    </div>

                    <p class="border-b border-dashed border-black mb-2"></p>

                    <div v-for="item in selectedTrx.details" :key="item.id" class="mb-1.5 font-bold font-mono text-[11px] leading-tight uppercase">
                        <div class="truncate w-full pr-2">{{ item.product?.nama_produk || 'Produk' }}</div>
                        <div class="flex justify-between pl-4 text-[10px]">
                            <span>{{ item.kuantitas }} x {{ formatRupiah(item.harga_satuan) }}</span>
                            <span>{{ formatRupiah(item.sub_total) }}</span>
                        </div>
                    </div>

                    <p class="border-t border-dashed border-black mt-2 pt-2"></p>

                    <div class="flex justify-between font-black text-xs mb-2 font-mono uppercase">
                        <span>TOTAL BELANJA :</span>
                        <span>{{ formatRupiah(selectedTrx.total_harga) }}</span>
                    </div>

                    <p class="border-b border-dashed border-black mb-2"></p>

                    <div class="flex justify-between mb-1 font-bold font-mono text-[11px] uppercase">
                        <span>TUNAI/BAYAR :</span>
                        <span>{{ formatRupiah(selectedTrx.nominal_bayar) }}</span>
                    </div>
                    <div class="flex justify-between mb-2 font-bold font-mono text-[11px] uppercase">
                        <span>KEMBALIAN :</span>
                        <span>{{ formatRupiah(selectedTrx.kembalian) }}</span>
                    </div>

                    <div class="mt-4 text-[9px] font-medium text-center border-t border-dashed border-black pt-2 font-mono uppercase">
                        <p>SUBTOTAL: {{ formatRupiah(selectedTrx.sub_total) }} | PAJAK: {{ formatRupiah(selectedTrx.pajak) }}</p>
                        <p class="mt-1">TRX-ID: {{ selectedTrx.no_invoice }}</p>
                    </div>

                    <div class="text-center mt-4 font-bold font-mono text-[11px]">
                        <p>=== TERIMA KASIH ===</p>
                        <p>BARANG YANG SUDAH DIBELI</p>
                        <p>TIDAK DAPAT DITUKAR/DIKEMBALIKAN</p>
                    </div>
                </div>

                <div class="mt-4 flex gap-3 shrink-0">
                    <button @click="printReceipt" class="flex-1 bg-gray-800 text-white py-3 rounded-xl font-black hover:bg-gray-900 transition-colors shadow-md text-sm">🖨️ CETAK NOTA</button>
                    <button @click="showReceipt = false" class="flex-1 bg-white border border-gray-300 py-3 rounded-xl font-black text-gray-800 hover:bg-gray-50 transition-colors text-sm">TUTUP</button>
                </div>
            </div>
        </div>
    </Sidebar>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 5px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #e2e8f0; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #cbd5e1; }

/* CSS KHUSUS PRINT */
@media print {
    body * { visibility: hidden; }
    #print-area, #print-area * { visibility: visible; }
    #print-area { position: absolute; left: 0; top: 0; width: 58mm; padding: 0; margin: 0; }
    @page { margin: 0; }
    .no-print { display: none !important; }
}
</style>