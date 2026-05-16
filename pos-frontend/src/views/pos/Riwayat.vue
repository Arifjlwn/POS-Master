<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import api from '../../api';
import Sidebar from '../../components/Sidebar.vue';

// --- STATE DATA ---
const riwayat = ref([]);
const isLoading = ref(true);
const tanggalDipilih = ref(new Date().toISOString().split('T')[0]); 

// --- STATE PENCARIAN ---
const searchQuery = ref('');

// --- STATE MODAL STRUK ---
const showReceipt = ref(false);
const selectedTrx = ref(null);

const formatRupiah = (angka) => {
    return new Intl.NumberFormat('id-ID').format(angka);
};

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

onMounted(() => fetchRiwayat());
watch(tanggalDipilih, () => fetchRiwayat());

// 🚀 FITUR PENCARIAN REALTIME (Tanpa perlu fetch API ulang)
const filteredRiwayat = computed(() => {
    if (!searchQuery.value) return riwayat.value;
    const query = searchQuery.value.toLowerCase();
    return riwayat.value.filter(trx => 
        (trx.no_invoice && trx.no_invoice.toLowerCase().includes(query)) ||
        (trx.User?.name && trx.User.name.toLowerCase().includes(query))
    );
});

const openReceipt = (trx) => {
    selectedTrx.value = trx;
    showReceipt.value = true;
};

const printReceipt = () => {
    const printContent = document.getElementById('print-area').innerHTML;
    const printWindow = window.open('', '_blank', 'width=300,height=600');
    
    printWindow.document.write(`
        <html>
            <head>
                <title>Cetak Struk Transaksi ${selectedTrx.value?.no_invoice}</title>
                <style>
                    body { font-family: 'Courier New', Courier, monospace; width: 58mm; margin: 0; padding: 0; font-size: 11px; color: #000; }
                    .text-center { text-align: center; }
                    .font-black { font-weight: 900; }
                    .font-bold { font-weight: bold; }
                    .flex { display: flex; }
                    .justify-between { justify-content: space-between; }
                    .uppercase { text-transform: uppercase; }
                    .border-y { border-top: 1px dashed #000; border-bottom: 1px dashed #000; padding: 5px 0; }
                    .border-b { border-bottom: 1px dashed #000; margin-bottom: 5px; padding-bottom: 5px; }
                    .border-t { border-top: 1px dashed #000; margin-top: 5px; padding-top: 5px; }
                    .truncate { white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
                    .w-full { width: 100%; }
                    .pl-4 { padding-left: 10px; }
                    p { margin: 2px 0; }
                </style>
            </head>
            <body onload="window.print(); window.close();">
                ${printContent}
            </body>
        </html>
    `);
    printWindow.document.close();
};
</script>

<template>
    <Sidebar>
        <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 font-sans bg-[#f8fafc] min-h-screen">
            
            <div class="bg-gradient-to-br from-slate-900 via-slate-800 to-blue-900 rounded-[32px] p-6 md:p-10 mb-6 text-white shadow-2xl flex flex-col md:flex-row md:items-center justify-between relative overflow-hidden gap-6">
                <div class="z-10 text-center md:text-left">
                    <h1 class="text-3xl md:text-4xl font-black tracking-tighter mb-1 uppercase italic">Daily <span class="text-blue-400">Journal</span></h1>
                    <p class="text-slate-300 font-bold text-[10px] md:text-xs uppercase tracking-[0.2em]">Audit, Reprint & Track Store Sales</p>
                </div>

                <div class="z-10 flex items-center justify-center gap-3 bg-white/10 backdrop-blur-md p-2 pl-4 rounded-2xl border border-white/20 shadow-inner w-full md:w-auto">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-blue-300" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
                    <input type="date" v-model="tanggalDipilih" class="bg-transparent border-none text-sm md:text-base font-black text-white focus:ring-0 cursor-pointer outline-none uppercase tracking-tighter w-full md:w-auto [color-scheme:dark]">
                </div>
            </div>

            <div class="mb-6 relative group">
                <div class="absolute inset-y-0 left-0 pl-5 flex items-center pointer-events-none">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-400 group-focus-within:text-blue-600 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
                </div>
                <input v-model="searchQuery" type="text" placeholder="Cari Nomor Invoice atau Nama Kasir..." class="block w-full pl-14 pr-4 py-4 bg-white rounded-2xl border-2 border-slate-100 shadow-sm focus:border-blue-600 outline-none font-bold text-sm transition-all text-slate-700">
            </div>

            <div v-if="isLoading" class="flex flex-col items-center justify-center p-16 bg-white rounded-[32px] border-2 border-dashed border-slate-200">
                <div class="w-10 h-10 border-4 border-blue-200 border-t-blue-600 rounded-full animate-spin mb-4"></div>
                <p class="text-slate-400 font-black text-xs uppercase tracking-widest animate-pulse">Menarik Data Transaksi...</p>
            </div>
            
            <div v-else-if="filteredRiwayat.length === 0" class="flex flex-col items-center justify-center p-16 bg-white/50 rounded-[32px] border-2 border-dashed border-slate-300">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-20 h-20 text-slate-300 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" /></svg>
                <p class="text-slate-400 font-black text-sm uppercase tracking-widest">Tidak ada transaksi ditemukan</p>
            </div>

            <div v-else class="md:hidden grid grid-cols-1 gap-4">
                <div v-for="trx in filteredRiwayat" :key="trx.id" class="bg-white p-5 rounded-[24px] shadow-sm border border-slate-100 flex flex-col gap-3 relative overflow-hidden">
                    <div class="absolute left-0 top-0 bottom-0 w-1.5" :class="trx.metode_bayar === 'Cash' ? 'bg-emerald-400' : 'bg-blue-500'"></div>
                    
                    <div class="flex justify-between items-start">
                        <div>
                            <div class="font-mono font-black text-sm text-slate-800 tracking-tighter">{{ trx.no_invoice }}</div>
                            <div class="text-[10px] font-bold text-slate-400 uppercase mt-0.5">
                                {{ new Date(trx.created_at).toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' }) }} • OPR: {{ trx.User?.name.split(' ')[0] || 'KASIR' }}
                            </div>
                        </div>
                        <span class="px-2.5 py-1 rounded-lg text-[9px] font-black uppercase tracking-widest border" :class="trx.metode_bayar === 'Cash' ? 'border-emerald-100 bg-emerald-50 text-emerald-600' : 'border-blue-100 bg-blue-50 text-blue-600'">
                            {{ trx.metode_bayar || 'CASH' }}
                        </span>
                    </div>

                    <div class="border-t border-dashed border-slate-200 my-1"></div>

                    <div class="flex justify-between items-end">
                        <div>
                            <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest">Total Transaksi</p>
                            <p class="font-black text-slate-800 text-lg">Rp{{ formatRupiah(trx.total_harga) }}</p>
                        </div>
                        <button @click="openReceipt(trx)" class="flex items-center justify-center w-10 h-10 bg-slate-100 text-slate-600 rounded-xl hover:bg-slate-900 hover:text-white transition-colors active:scale-95 shadow-sm">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M6 9V2h12v7"/><path d="M6 18H4a2 2 0 0 1-2-2v-5a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v5a2 2 0 0 1-2 2h-2"/><rect width="12" height="8" x="6" y="14"/></svg>
                        </button>
                    </div>
                </div>
            </div>

            <div v-if="filteredRiwayat.length > 0" class="hidden md:block w-full bg-white rounded-[32px] shadow-sm border border-slate-100 overflow-hidden transition-all">
                <table class="w-full text-left border-collapse whitespace-nowrap">
                    <thead>
                        <tr class="bg-slate-50 border-b border-slate-100">
                            <th class="px-6 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">No. Invoice</th>
                            <th class="px-6 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">Waktu</th>
                            <th class="px-6 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">Kasir</th>
                            <th class="px-6 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">Metode</th>
                            <th class="px-6 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest text-right">Total Nominal</th>
                            <th class="px-6 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest text-center">Aksi</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-slate-50">
                        <tr v-for="trx in filteredRiwayat" :key="trx.id" class="hover:bg-blue-50/40 transition-all group">
                            <td class="px-6 py-4 font-mono font-black text-sm text-slate-700 tracking-tighter">{{ trx.no_invoice }}</td>
                            <td class="px-6 py-4 text-xs font-bold text-slate-500 uppercase">
                                {{ new Date(trx.created_at).toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' }) }}
                            </td>
                            <td class="px-6 py-4 text-xs font-black text-slate-700 uppercase">{{ trx.User?.name || 'KASIR' }}</td>
                            <td class="px-6 py-4">
                                <span class="px-3 py-1.5 rounded-lg text-[10px] font-black uppercase tracking-widest border" :class="trx.metode_bayar === 'Cash' ? 'border-emerald-100 bg-emerald-50 text-emerald-600' : 'border-blue-100 bg-blue-50 text-blue-600'">
                                    {{ trx.metode_bayar || 'CASH' }}
                                </span>
                            </td>
                            <td class="px-6 py-4 text-right font-black text-slate-800 text-sm">Rp{{ formatRupiah(trx.total_harga) }}</td>
                            <td class="px-6 py-4 text-center">
                                <button @click="openReceipt(trx)" class="inline-flex items-center justify-center gap-2 bg-slate-50 text-slate-500 border border-slate-200 px-4 py-2.5 rounded-xl text-[10px] font-black uppercase tracking-widest hover:bg-slate-900 hover:text-white hover:border-slate-900 transition-all active:scale-95 shadow-sm">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M6 9V2h12v7"/><path d="M6 18H4a2 2 0 0 1-2-2v-5a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v5a2 2 0 0 1-2 2h-2"/><rect width="12" height="8" x="6" y="14"/></svg>
                                    Cetak Struk
                                </button>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </main>

        <div v-if="showReceipt" class="fixed inset-0 bg-slate-950/90 flex items-center justify-center z-[200] p-4 backdrop-blur-sm no-print">
            <div class="bg-white rounded-[40px] p-6 md:p-8 max-w-sm w-full shadow-2xl border-[8px] border-slate-800 flex flex-col max-h-[90vh]">
                
                <div class="overflow-y-auto custom-scrollbar bg-white p-2 mx-auto" id="print-area" style="width: 58mm;">
                    <div class="text-center mb-4 font-mono leading-none">
                        <h2 class="font-black text-sm uppercase tracking-tighter mb-1 italic">
                            {{ selectedTrx.Store?.nama_toko || 'ARZU STORE' }}
                        </h2>
                        <p class="text-[8px] font-bold uppercase tracking-widest opacity-80">
                            {{ selectedTrx.Store?.alamat || 'JAKARTA, INDONESIA' }}
                        </p>
                        <p class="text-[8px] font-bold uppercase tracking-widest opacity-80 mt-1">
                            WA: {{ selectedTrx.Store?.no_hp || '-' }}
                        </p>
                    </div>

                    <div class="border-y border-black py-1.5 text-center font-black mb-3 font-mono text-[9px] tracking-[0.2em] uppercase bg-slate-100">
                        Invoice Reprint
                    </div>

                    <div class="mb-3 text-[8px] font-bold font-mono uppercase space-y-0.5">
                        <div class="flex justify-between"><span>DATE:</span><span>{{ new Date(selectedTrx.created_at).toLocaleString('id-ID') }}</span></div>
                        <div class="flex justify-between"><span>CASHIER:</span><span>{{ selectedTrx.User?.name.split(' ')[0] || 'KASIR' }}</span></div>
                    </div>

                    <div class="border-b border-black border-dashed mb-2"></div>

                    <div v-for="item in selectedTrx.details" :key="item.id" class="mb-2 font-bold font-mono text-[9px] leading-tight uppercase">
                        <div class="truncate w-full pr-2">{{ item.Product?.nama_produk || item.product?.nama_produk || 'Produk' }}</div>
                        <div class="flex justify-between pl-2 text-[8px] mt-0.5">
                            <span>{{ item.kuantitas }} x {{ formatRupiah(item.harga_satuan) }}</span>
                            <span class="font-black text-[9px]">{{ formatRupiah(item.sub_total) }}</span>
                        </div>
                    </div>

                    <div class="border-t border-black border-dashed mt-2 pt-2"></div>

                    <div class="flex justify-between font-black text-[11px] mb-2 font-mono uppercase italic">
                        <span>TOTAL:</span>
                        <span>Rp{{ formatRupiah(selectedTrx.total_harga) }}</span>
                    </div>

                    <div class="border-b border-black border-dashed mb-2"></div>

                    <div class="flex justify-between mb-1 font-bold font-mono text-[8px] uppercase">
                        <span>PAID ({{ selectedTrx.metode_bayar || 'CASH' }}):</span>
                        <span>Rp{{ formatRupiah(selectedTrx.nominal_bayar) }}</span>
                    </div>
                    <div class="flex justify-between font-black font-mono text-[9px] uppercase italic text-black">
                        <span>CHANGE:</span>
                        <span>Rp{{ formatRupiah(selectedTrx.kembalian) }}</span>
                    </div>

                    <div class="mt-5 text-[7px] font-bold text-center border-t border-black border-dashed pt-2 font-mono uppercase space-y-1">
                        <p>SUBTOTAL: Rp{{ formatRupiah(selectedTrx.sub_total) }} | TAX: Rp{{ formatRupiah(selectedTrx.pajak) }}</p>
                        <p class="font-black">INV: {{ selectedTrx.no_invoice }}</p>
                        <p>STN: {{ selectedTrx.station_number || '01' }}</p>
                    </div>

                    <div class="text-center mt-4 font-black font-mono text-[8px] border-2 border-black p-1.5 uppercase">
                        Terima Kasih!<br>Barang tidak dapat ditukar.
                    </div>
                </div>

                <div class="mt-6 flex flex-col gap-2 md:gap-3 shrink-0 no-print">
                    <button @click="printReceipt" class="w-full bg-blue-600 text-white py-3 md:py-4 rounded-xl md:rounded-2xl font-black text-[9px] md:text-[10px] uppercase tracking-[0.2em] shadow-xl shadow-blue-200 flex items-center justify-center gap-2 active:scale-95 transition-all">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 md:w-5 md:h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M6 9V2h12v7"/><path d="M6 18H4a2 2 0 0 1-2-2v-5a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v5a2 2 0 0 1-2 2h-2"/><rect width="12" height="8" x="6" y="14"/></svg>
                        Cetak Ulang Struk
                    </button>
                    <button @click="showReceipt = false" class="w-full bg-slate-100 text-slate-500 py-3 md:py-3.5 rounded-xl md:rounded-2xl font-black text-[9px] md:text-[10px] uppercase tracking-widest hover:bg-slate-200 transition-all">Tutup</button>
                </div>
            </div>
        </div>
    </Sidebar>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 5px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 20px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #94a3b8; }

/* Styling Calendar Icon pada Input Date biar warnanya nyatu */
input[type="date"]::-webkit-calendar-picker-indicator {
    filter: invert(1);
    cursor: pointer;
    opacity: 0.8;
}
input[type="date"]::-webkit-calendar-picker-indicator:hover {
    opacity: 1;
}

@media print {
    body * { visibility: hidden; }
    #print-area, #print-area * { visibility: visible; }
    #print-area {
        position: absolute;
        left: 0;
        top: 0;
        width: 58mm; 
        margin: 0;
        padding: 0;
    }
    @page { margin: 0; }
}
</style>