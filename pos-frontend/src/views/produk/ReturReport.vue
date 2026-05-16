<script setup>
import { ref, onMounted, computed } from 'vue';
import Sidebar from '../../components/Sidebar.vue';
import api from '../../api.js';
import Swal from 'sweetalert2';

const returns = ref([]);
const isLoading = ref(true);
const searchQuery = ref('');

// State untuk Modal Detail/Cetak
const isModalOpen = ref(false);
const selectedDocument = ref(null);

const user = ref({
    storeName: localStorage.getItem('storeName') || 'POS UMKM',
    name: localStorage.getItem('name') || 'Admin'
});

const fetchReturns = async () => {
    isLoading.value = true;
    try {
        const response = await api.get('/returns');
        returns.value = response.data.data || [];
    } catch (error) {
        Swal.fire('Error', 'Gagal memuat data riwayat retur.', 'error');
    } finally {
        isLoading.value = false;
    }
};

onMounted(() => fetchReturns());

// 🚀 SIHIR GHAIB: Grouping data flat menjadi per-Nomor Dokumen
const groupedReturns = computed(() => {
    const groups = {};
    returns.value.forEach(ret => {
        // Karena di backend mungkin return_no ada yang kosong (data lama), kita kasih fallback
        const docNo = ret.return_no || `RET-OLD-${ret.id}`; 
        
        if (!groups[docNo]) {
            groups[docNo] = {
                return_no: docNo,
                created_at: ret.created_at,
                user: ret.user || { name: 'Sistem' },
                items: [],
                total_qty: 0
            };
        }
        groups[docNo].items.push(ret);
        groups[docNo].total_qty += ret.qty;
    });

    // Ubah Object jadi Array dan urutkan dari yang terbaru
    return Object.values(groups).sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
});

// 🚀 FITUR PENCARIAN REALTIME
const filteredDocuments = computed(() => {
    if (!searchQuery.value) return groupedReturns.value;
    const query = searchQuery.value.toLowerCase();
    
    return groupedReturns.value.filter(doc => 
        doc.return_no.toLowerCase().includes(query) ||
        doc.user.name.toLowerCase().includes(query) ||
        // Cari juga ke dalam item-itemnya, siapa tau nyari nama barang
        doc.items.some(item => item.product?.nama_produk?.toLowerCase().includes(query))
    );
});

// Buka Modal Detail
const openDetail = (doc) => {
    selectedDocument.value = doc;
    isModalOpen.value = true;
};

// Tutup Modal Detail
const closeDetail = () => {
    isModalOpen.value = false;
    setTimeout(() => { selectedDocument.value = null; }, 300);
};

// 🚀 FUNGSI CETAK BERITA ACARA
const printDocument = () => {
    window.print();
};
</script>

<template>
    <Sidebar>
        <div class="p-4 md:p-8 lg:p-10 max-w-7xl mx-auto font-sans bg-[#f8fafc] min-h-screen print:bg-white print:p-0 print:m-0">
            
            <div class="bg-gradient-to-br from-indigo-900 via-slate-800 to-slate-900 rounded-[32px] p-6 md:p-10 mb-8 text-white shadow-2xl flex flex-col md:flex-row items-center justify-between relative overflow-hidden border border-slate-800 gap-6 print:hidden">
                <svg xmlns="http://www.w3.org/2000/svg" class="absolute -right-10 -bottom-10 w-64 h-64 text-indigo-500 opacity-10 pointer-events-none" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/><polyline points="10 9 9 9 8 9"/></svg>
                
                <div class="z-10 text-center md:text-left">
                    <h1 class="text-3xl md:text-4xl font-black tracking-tighter mb-1 uppercase italic">Laporan <span class="text-indigo-400">Waste & Retur</span></h1>
                    <p class="text-slate-300 font-bold text-[10px] md:text-xs uppercase tracking-[0.2em] flex items-center justify-center md:justify-start gap-2">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-indigo-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>
                        Riwayat Berita Acara Pemusnahan Barang
                    </p>
                </div>
            </div>

            <div class="mb-6 relative group print:hidden">
                <div class="absolute inset-y-0 left-0 pl-5 flex items-center pointer-events-none">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-400 group-focus-within:text-indigo-600 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
                </div>
                <input v-model="searchQuery" type="text" placeholder="Cari Nomor Dokumen (RET-...) atau Nama Barang..." class="block w-full pl-14 pr-4 py-4 bg-white rounded-2xl border-2 border-slate-100 shadow-sm focus:border-indigo-600 outline-none font-bold text-sm transition-all text-slate-700">
            </div>

            <div v-if="isLoading" class="py-20 flex flex-col items-center justify-center bg-white rounded-[32px] border border-slate-100 print:hidden">
                <div class="w-12 h-12 border-4 border-slate-200 border-t-indigo-600 rounded-full animate-spin mb-4"></div>
                <p class="text-slate-400 font-black text-xs uppercase tracking-widest animate-pulse">Menyusun Laporan...</p>
            </div>

            <div v-else-if="filteredDocuments.length === 0" class="flex flex-col items-center justify-center py-20 bg-white/50 rounded-[32px] border-2 border-dashed border-slate-200 print:hidden">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-20 h-20 text-slate-300 mb-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/><polyline points="10 9 9 9 8 9"/></svg>
                <p class="text-slate-400 font-black text-sm uppercase tracking-widest">Dokumen Tidak Ditemukan</p>
            </div>

            <div v-else class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6 print:hidden">
                <div v-for="doc in filteredDocuments" :key="doc.return_no" class="bg-white rounded-[24px] p-6 shadow-sm border border-slate-100 hover:shadow-xl hover:border-indigo-100 transition-all group flex flex-col justify-between">
                    
                    <div>
                        <div class="flex items-center justify-between mb-4 border-b border-dashed border-slate-100 pb-4">
                            <div class="inline-flex items-center gap-2 bg-indigo-50 text-indigo-700 px-3 py-1.5 rounded-xl text-[10px] font-black uppercase tracking-widest">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
                                {{ doc.return_no }}
                            </div>
                            <span class="text-[9px] font-bold text-slate-400">{{ new Date(doc.created_at).toLocaleDateString('id-ID') }}</span>
                        </div>

                        <div class="flex items-center gap-4 mb-4">
                            <div class="w-12 h-12 rounded-2xl bg-rose-50 flex flex-col items-center justify-center border border-rose-100">
                                <span class="text-xs font-black text-rose-600 leading-none">{{ doc.items.length }}</span>
                                <span class="text-[8px] font-black text-rose-400 uppercase tracking-widest mt-1">Item</span>
                            </div>
                            <div class="flex-1">
                                <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest">Total Dibuang</p>
                                <p class="text-base font-black text-slate-800">{{ doc.total_qty }} <span class="text-[10px] text-slate-500 uppercase">Unit</span></p>
                            </div>
                        </div>

                        <div class="flex -space-x-2 overflow-hidden py-1 mb-4">
                            <span v-for="(item, idx) in doc.items.slice(0, 3)" :key="idx" class="inline-block px-2 py-1 bg-slate-100 text-slate-600 text-[9px] font-bold rounded-lg border border-white whitespace-nowrap truncate max-w-[120px]">
                                {{ item.product?.nama_produk || 'Item' }}
                            </span>
                            <span v-if="doc.items.length > 3" class="inline-block px-2 py-1 bg-slate-200 text-slate-600 text-[9px] font-bold rounded-lg border border-white">
                                +{{ doc.items.length - 3 }} Lainnya
                            </span>
                        </div>
                    </div>

                    <div class="flex items-center justify-between pt-4 border-t border-slate-100 mt-2">
                        <div class="flex items-center gap-2">
                            <div class="w-6 h-6 rounded-full bg-slate-800 text-white flex items-center justify-center text-[8px] font-black">{{ doc.user.name.substring(0,2).toUpperCase() }}</div>
                            <span class="text-[10px] font-bold text-slate-500 truncate max-w-[80px]">{{ doc.user.name }}</span>
                        </div>
                        <button @click="openDetail(doc)" class="bg-slate-100 hover:bg-indigo-600 text-slate-600 hover:text-white px-4 py-2 rounded-xl text-[10px] font-black uppercase tracking-widest transition-colors flex items-center gap-2">
                            Detail
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M5 12h14"/><path d="m12 5 7 7-7 7"/></svg>
                        </button>
                    </div>

                </div>
            </div>

            <div v-if="isModalOpen" class="fixed inset-0 z-[100] flex items-center justify-center p-4 print:p-0 print:relative print:block print:z-auto">
                <div @click="closeDetail" class="absolute inset-0 bg-slate-900/80 backdrop-blur-sm print:hidden transition-all"></div>
                
                <div class="relative bg-white w-full max-w-4xl max-h-[90vh] overflow-y-auto rounded-[32px] shadow-2xl flex flex-col print:shadow-none print:rounded-none print:max-h-none print:h-auto print:overflow-visible custom-scrollbar">
                    
                    <div class="sticky top-0 bg-white/90 backdrop-blur-md px-6 py-4 border-b border-slate-100 flex justify-between items-center z-10 print:hidden shrink-0">
                        <h3 class="font-black text-slate-800 text-lg uppercase tracking-tight">Detail Dokumen</h3>
                        <div class="flex items-center gap-3">
                            <button @click="printDocument" class="bg-indigo-600 hover:bg-slate-900 text-white px-4 py-2.5 rounded-xl font-black text-[10px] uppercase tracking-widest transition-all flex items-center gap-2 shadow-md">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="6 9 6 2 18 2 18 9"/><path d="M6 18H4a2 2 0 0 1-2-2v-5a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v5a2 2 0 0 1-2 2h-2"/><rect width="12" height="8" x="6" y="14"/></svg>
                                Cetak PDF
                            </button>
                            <button @click="closeDetail" class="w-10 h-10 rounded-xl bg-slate-100 text-slate-400 hover:text-red-500 hover:bg-red-50 transition-all flex items-center justify-center">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                            </button>
                        </div>
                    </div>

                    <div id="printable-area" class="p-8 md:p-12 print:p-0 print:text-black">
                        
                        <div class="text-center border-b-[3px] border-slate-800 pb-6 mb-6">
                            <h1 class="text-2xl font-black uppercase tracking-widest text-slate-900">{{ user.storeName }}</h1>
                            <p class="text-sm font-medium text-slate-600 mt-1">BERITA ACARA PEMUSNAHAN / RETUR BARANG</p>
                        </div>

                        <div class="flex justify-between items-end mb-8">
                            <div>
                                <table class="text-xs font-bold text-slate-700">
                                    <tr><td class="pb-2 pr-4 text-slate-500 uppercase tracking-widest text-[9px]">No. Dokumen</td><td class="pb-2">: {{ selectedDocument?.return_no }}</td></tr>
                                    <tr><td class="pb-2 pr-4 text-slate-500 uppercase tracking-widest text-[9px]">Tanggal</td><td class="pb-2">: {{ selectedDocument ? new Date(selectedDocument.created_at).toLocaleDateString('id-ID', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' }) : '' }}</td></tr>
                                    <tr><td class="pr-4 text-slate-500 uppercase tracking-widest text-[9px]">Operator</td><td>: {{ selectedDocument?.user.name }}</td></tr>
                                </table>
                            </div>
                            <div class="text-right">
                                <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1">Total Kuantitas</p>
                                <p class="text-3xl font-black text-slate-900 tracking-tighter">{{ selectedDocument?.total_qty }}</p>
                            </div>
                        </div>

                        <table class="w-full text-left border-collapse mb-12">
                            <thead>
                                <tr class="border-y-2 border-slate-800">
                                    <th class="py-3 px-2 text-[10px] font-black text-slate-800 uppercase tracking-widest w-12 text-center">No</th>
                                    <th class="py-3 px-2 text-[10px] font-black text-slate-800 uppercase tracking-widest">SKU</th>
                                    <th class="py-3 px-2 text-[10px] font-black text-slate-800 uppercase tracking-widest">Nama Barang</th>
                                    <th class="py-3 px-2 text-[10px] font-black text-slate-800 uppercase tracking-widest text-center">Qty</th>
                                    <th class="py-3 px-2 text-[10px] font-black text-slate-800 uppercase tracking-widest">Alasan / Klasifikasi</th>
                                </tr>
                            </thead>
                            <tbody class="divide-y divide-slate-200">
                                <tr v-for="(item, index) in selectedDocument?.items" :key="item.id">
                                    <td class="py-3 px-2 text-xs font-bold text-slate-600 text-center">{{ index + 1 }}</td>
                                    <td class="py-3 px-2 text-[10px] font-bold text-slate-500 uppercase tracking-wider">{{ item.product?.sku || '-' }}</td>
                                    <td class="py-3 px-2 text-xs font-black text-slate-800 uppercase">{{ item.product?.nama_produk || 'Produk Dihapus' }}</td>
                                    <td class="py-3 px-2 text-sm font-black text-slate-900 text-center">{{ item.qty }}</td>
                                    <td class="py-3 px-2">
                                        <div class="text-xs font-bold text-slate-700">{{ item.alasan }}</div>
                                        <div v-if="item.catatan" class="text-[10px] font-medium text-slate-500 italic mt-0.5">Catatan: {{ item.catatan }}</div>
                                    </td>
                                </tr>
                            </tbody>
                        </table>

                        <div class="grid grid-cols-2 gap-8 mt-16 pt-8 break-inside-avoid">
                            <div class="text-center">
                                <p class="text-[10px] font-bold text-slate-500 uppercase tracking-widest mb-16">Dibuat Oleh,</p>
                                <p class="text-sm font-black text-slate-800 uppercase underline">{{ selectedDocument?.user.name }}</p>
                                <p class="text-[9px] font-bold text-slate-400 mt-1">Staff / Kasir</p>
                            </div>
                            <div class="text-center">
                                <p class="text-[10px] font-bold text-slate-500 uppercase tracking-widest mb-16">Mengetahui,</p>
                                <p class="text-sm font-black text-slate-800 uppercase underline">...................................</p>
                                <p class="text-[9px] font-bold text-slate-400 mt-1">Manager / Owner</p>
                            </div>
                        </div>

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

/* 🚀 CSS SAKTI BUAT CETAK PDF / PRINTER */
@media print {
    /* Sembunyikan semua elemen kecuali modal yang lagi kebuka */
    body * {
        visibility: hidden;
    }
    
    /* Tampilkan hanya area kertas putih */
    #printable-area, #printable-area * {
        visibility: visible;
    }
    
    /* Posisikan kertas putih di pojok kiri atas pas diprint */
    #printable-area {
        position: absolute;
        left: 0;
        top: 0;
        width: 100%;
    }

    /* Hilangkan background color bawaan browser pas print */
    body {
        background-color: white !important;
    }

    /* Pastikan tabel tidak terpotong di tengah halaman */
    table { page-break-inside: auto; }
    tr    { page-break-inside: avoid; page-break-after: auto; }
    
    /* Atur margin kertas A4 */
    @page {
        margin: 20mm;
        size: A4 portrait;
    }
}
</style>