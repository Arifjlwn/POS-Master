<script setup>
import { ref, onMounted, watch } from 'vue';
import api from '../../api';
import Sidebar from '../../components/Sidebar.vue';

// --- STATE DATA ---
const riwayat = ref([]);
const isLoading = ref(true);
const tanggalDipilih = ref(new Date().toISOString().split('T')[0]); 

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
                <title>Cetak Struk</title>
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
        <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 font-sans">
            
            <div class="mb-8 flex flex-col md:flex-row md:items-center justify-between gap-6 bg-white p-6 rounded-[32px] border border-slate-100 shadow-sm">
                <div>
                    <h1 class="text-3xl font-black text-slate-800 tracking-tighter uppercase italic">Transaction <span class="text-blue-600">History</span></h1>
                    <p class="text-slate-400 text-[10px] font-black uppercase tracking-[0.2em] mt-1">Audit, Reprint, and Track Store Sales</p>
                </div>
                <div class="flex items-center gap-4 bg-slate-50 p-2 rounded-2xl border border-slate-100">
                    <div class="pl-3">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
                    </div>
                    <input type="date" v-model="tanggalDipilih" class="bg-transparent border-none text-sm font-black text-slate-700 focus:ring-0 cursor-pointer outline-none uppercase tracking-tighter">
                </div>
            </div>

            <div class="bg-white rounded-[40px] shadow-sm border border-slate-100 overflow-hidden transition-all">
                <div class="p-6 border-b border-slate-50 flex items-center gap-3 bg-slate-50/30">
                    <div class="w-10 h-10 bg-blue-600 rounded-xl flex items-center justify-center shadow-lg shadow-blue-200">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/><polyline points="10 9 9 9 8 9"/></svg>
                    </div>
                    <h2 class="text-lg font-black text-slate-800 uppercase tracking-tighter italic">Daily Journal</h2>
                </div>

                <div class="overflow-x-auto">
                    <table class="w-full text-left border-collapse">
                        <thead>
                            <tr class="bg-white border-b border-slate-100">
                                <th class="px-8 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">Inv. Number</th>
                                <th class="px-8 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">Timestamp</th>
                                <th class="px-8 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">Operator</th>
                                <th class="px-8 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest">Method</th>
                                <th class="px-8 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest text-right">Total</th>
                                <th class="px-8 py-5 text-[10px] font-black text-slate-400 uppercase tracking-widest text-center">Action</th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-slate-50">
                            <tr v-if="isLoading">
                                <td colspan="6" class="px-8 py-20 text-center font-black text-slate-300 uppercase tracking-widest text-xs">Loading Journal...</td>
                            </tr>
                            <tr v-else-if="riwayat.length === 0">
                                <td colspan="6" class="px-8 py-20 text-center font-black text-slate-300 uppercase tracking-widest text-xs italic">No Transactions Found</td>
                            </tr>
                            <tr v-for="trx in riwayat" :key="trx.id" class="hover:bg-slate-50/80 transition-all">
                                <td class="px-8 py-5 font-mono font-black text-sm text-blue-600 tracking-tighter">{{ trx.no_invoice }}</td>
                                <td class="px-8 py-5 text-[11px] font-bold text-slate-500 uppercase italic">
                                    {{ new Date(trx.created_at).toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' }) }}
                                </td>
                                <td class="px-8 py-5 text-[11px] font-black text-slate-700 uppercase">{{ trx.User?.name || 'KASIR' }}</td>
                                <td class="px-8 py-5">
                                    <span class="px-3 py-1 rounded-full text-[9px] font-black uppercase tracking-widest border border-blue-100 bg-blue-50 text-blue-600">
                                        {{ trx.metode_bayar || 'Cash' }}
                                    </span>
                                </td>
                                <td class="px-8 py-5 text-right font-black text-slate-800 text-sm">Rp{{ formatRupiah(trx.total_harga) }}</td>
                                <td class="px-8 py-5 text-center">
                                    <button @click="openReceipt(trx)" class="inline-flex items-center gap-2 bg-slate-900 text-white px-5 py-2.5 rounded-2xl text-[10px] font-black uppercase tracking-widest hover:bg-blue-600 transition-all active:scale-95 shadow-sm">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M6 9V2h12v7"/><path d="M6 18H4a2 2 0 0 1-2-2v-5a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v5a2 2 0 0 1-2 2h-2"/><rect width="12" height="8" x="6" y="14"/></svg>
                                        Reprint
                                    </button>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </main>

        <div v-if="showReceipt" class="fixed inset-0 bg-slate-900/90 flex items-center justify-center z-[100] p-4 backdrop-blur-sm no-print">
            <div class="bg-white rounded-[48px] p-8 max-w-sm w-full shadow-2xl border-[12px] border-slate-950/5 flex flex-col max-h-[90vh]">
                
                <div class="overflow-y-auto custom-scrollbar bg-white p-2 mx-auto" id="print-area" style="width: 58mm;">
                    <div class="text-center mb-5 font-mono leading-none">
                        <h2 class="font-black text-base uppercase tracking-tighter mb-1 italic">
                            {{ selectedTrx.Store?.nama_toko || 'NAMA TOKO' }}
                        </h2>
                        <p class="text-[8px] font-bold uppercase tracking-widest opacity-60">
                            {{ selectedTrx.Store?.alamat || 'ALAMAT BELUM DISET' }}
                        </p>
                        <p class="text-[7px] font-medium opacity-40 mt-1">
                            WA: {{ selectedTrx.Store?.no_hp || '-' }}
                        </p>
                    </div>

                    <div class="border-y border-black py-2 text-center font-black mb-3 font-mono text-[10px] tracking-[0.2em] uppercase bg-slate-50">
                        Invoice Reprint
                    </div>

                    <div class="mb-3 text-[9px] font-bold font-mono uppercase space-y-0.5">
                        <div class="flex justify-between"><span>DATE:</span><span>{{ new Date(selectedTrx.created_at).toLocaleString('id-ID') }}</span></div>
                        <div class="flex justify-between"><span>CASHIER:</span><span>{{ selectedTrx.User?.name.split(' ')[0] || 'KASIR' }}</span></div>
                    </div>

                    <div class="border-b border-black border-dashed mb-3"></div>

                    <div v-for="item in selectedTrx.details" :key="item.id" class="mb-3 font-bold font-mono text-[10px] leading-tight uppercase">
                        <div class="truncate w-full pr-2">{{ item.Product?.nama_produk || item.product?.nama_produk || 'Produk' }}</div>
                        <div class="flex justify-between pl-4 text-[9px] mt-0.5">
                            <span>{{ item.kuantitas }} x {{ formatRupiah(item.harga_satuan) }}</span>
                            <span class="font-black text-[10px]">{{ formatRupiah(item.sub_total) }}</span>
                        </div>
                    </div>

                    <div class="border-t border-black border-dashed mt-3 pt-3"></div>

                    <div class="flex justify-between font-black text-xs mb-3 font-mono uppercase italic">
                        <span>GRAND TOTAL :</span>
                        <span>{{ formatRupiah(selectedTrx.total_harga) }}</span>
                    </div>

                    <div class="border-b border-black border-dashed mb-3"></div>

                    <div class="flex justify-between mb-1 font-bold font-mono text-[9px] uppercase">
                        <span>PAID ({{ selectedTrx.metode_bayar || 'CASH' }}):</span>
                        <span>{{ formatRupiah(selectedTrx.nominal_bayar) }}</span>
                    </div>
                    <div class="flex justify-between font-black font-mono text-[10px] uppercase italic text-blue-800">
                        <span>CHANGE:</span>
                        <span>{{ formatRupiah(selectedTrx.kembalian) }}</span>
                    </div>

                    <div class="mt-6 text-[8px] font-bold text-center border-t border-black border-dashed pt-3 font-mono uppercase space-y-1">
                        <p>SUBTOTAL: {{ formatRupiah(selectedTrx.sub_total) }} | TAX: {{ formatRupiah(selectedTrx.pajak) }}</p>
                        <p class="font-black">INV-ID: {{ selectedTrx.no_invoice }}</p>
                        <p>STATION: {{ selectedTrx.station_number || '01' }}</p>
                    </div>

                    <div class="text-center mt-8 font-black font-mono text-[10px] border-2 border-black p-2 rotate-1 uppercase">
                        Terima Kasih Telah Berbelanja!
                    </div>
                </div>

                <div class="mt-8 flex flex-col gap-3 shrink-0">
                    <button @click="printReceipt" class="w-full bg-blue-600 text-white py-5 rounded-3xl font-black text-[10px] uppercase tracking-[0.2em] shadow-xl shadow-blue-200 flex items-center justify-center gap-3">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><path d="M6 9V2h12v7"/><path d="M6 18H4a2 2 0 0 1-2-2v-5a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v5a2 2 0 0 1-2 2h-2"/><rect width="12" height="8" x="6" y="14"/></svg>
                        Reprint Nota
                    </button>
                    <button @click="showReceipt = false" class="w-full bg-slate-100 text-slate-400 py-4 rounded-3xl font-black text-[10px] uppercase tracking-widest hover:text-slate-600 transition-all">Dismiss</button>
                </div>
            </div>
        </div>
    </Sidebar>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 5px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #e2e8f0; border-radius: 20px; }

@media (max-width: 1024px) {
    .overflow-x-auto {
        -webkit-overflow-scrolling: touch;
    }
}

@media print {
    /* Sembunyikan SEMUA elemen di halaman */
    body * {
        visibility: hidden;
    }
    /* Tampilkan hanya area struk dan semua isinya */
    #print-area, #print-area * {
        visibility: visible;
    }
    /* Atur posisi struk ke pojok kiri atas kertas */
    #print-area {
        position: absolute;
        left: 0;
        top: 0;
        width: 58mm; /* Sesuaikan lebar printer thermal */
        margin: 0;
        padding: 0;
    }
    /* Hilangkan header/footer default browser (tanggal, url) */
    @page {
        margin: 0;
    }
}
</style>