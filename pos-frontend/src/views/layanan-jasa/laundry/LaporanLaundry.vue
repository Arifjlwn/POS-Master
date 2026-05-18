<script setup>
import { ref, onMounted, computed } from 'vue';
import SidebarLaundry from './SidebarLaundry.vue';
import api from '../../../api.js'; 
import Swal from 'sweetalert2';

const rekap = ref({ tunai: 0, qris: 0, debit: 0, piutang: 0 });
const riwayat = ref([]);
const isLoading = ref(false);
const searchQuery = ref('');

// State Modal Bukti QRIS
const showBuktiModal = ref(false);
const selectedBuktiUrl = ref('');

const formatRupiah = (angka) => {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(angka || 0);
};

const formatDate = (dateStr) => {
    const d = new Date(dateStr);
    return new Intl.DateTimeFormat('id-ID', { day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' }).format(d);
};

const fetchLaporan = async () => {
    isLoading.value = true;
    try {
        const response = await api.get('/laundry/report');
        rekap.value = response.data.ringkasan;
        riwayat.value = response.data.transaksi;
    } catch (error) {
        Swal.fire('Gagal!', 'Gagal mengambil data laporan.', 'error');
    } finally {
        isLoading.value = false;
    }
};

onMounted(() => fetchLaporan());

const filteredRiwayat = computed(() => {
    if (!searchQuery.value) return riwayat.value;
    return riwayat.value.filter(t => 
        t.no_invoice.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
        t.pelanggan.toLowerCase().includes(searchQuery.value.toLowerCase())
    );
});

const bukaBuktiTransfer = (url) => {
    selectedBuktiUrl.value = url;
    showBuktiModal.value = true;
};
</script>

<template>
    <SidebarLaundry>
        <div class="flex-1 flex flex-col h-full bg-[#F8FAFC] overflow-hidden relative">
            
            <div class="p-5 md:p-8 shrink-0 bg-white border-b border-slate-100 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 z-10 shadow-sm">
                <div class="flex items-center gap-4">
                    <div class="w-12 h-12 bg-indigo-50 border border-indigo-100 rounded-2xl flex items-center justify-center shrink-0">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-indigo-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"/><line x1="3" y1="9" x2="21" y2="9"/><line x1="9" y1="21" x2="9" y2="9"/></svg>
                    </div>
                    <div>
                        <h1 class="text-xl md:text-2xl font-black tracking-tighter uppercase text-slate-800 leading-tight">Laporan Keuangan</h1>
                        <p class="text-[9px] md:text-[10px] font-black text-slate-400 mt-1 uppercase tracking-widest">Dashboard Monitoring Owner</p>
                    </div>
                </div>
            </div>

            <div class="flex-1 overflow-y-auto custom-scrollbar p-5 md:p-8">
                
                <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
                    <div class="bg-white p-5 rounded-3xl border border-slate-200 shadow-sm relative overflow-hidden group">
                        <div class="absolute -right-4 -bottom-4 opacity-[0.03] group-hover:opacity-10 transition-opacity"><svg xmlns="http://www.w3.org/2000/svg" class="w-32 h-32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="6" width="20" height="12" rx="2"/><circle cx="12" cy="12" r="2"/><path d="M6 12h.01M18 12h.01"/></svg></div>
                        <h3 class="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1">Kas Laci (Tunai)</h3>
                        <p class="text-2xl font-black text-slate-800 tracking-tighter">{{ formatRupiah(rekap.tunai) }}</p>
                    </div>
                    <div class="bg-white p-5 rounded-3xl border border-emerald-100 shadow-sm relative overflow-hidden group">
                        <div class="absolute -right-4 -bottom-4 opacity-[0.03] group-hover:opacity-10 text-emerald-500 transition-opacity"><svg xmlns="http://www.w3.org/2000/svg" class="w-32 h-32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"/><rect x="7" y="7" width="3" height="3"/><rect x="14" y="7" width="3" height="3"/><rect x="7" y="14" width="3" height="3"/><rect x="14" y="14" width="3" height="3"/></svg></div>
                        <h3 class="text-[10px] font-black text-emerald-500 uppercase tracking-widest mb-1">Masuk Rekening (QRIS)</h3>
                        <p class="text-2xl font-black text-emerald-600 tracking-tighter">{{ formatRupiah(rekap.qris) }}</p>
                    </div>
                    <div class="bg-white p-5 rounded-3xl border border-sky-100 shadow-sm relative overflow-hidden group">
                        <div class="absolute -right-4 -bottom-4 opacity-[0.03] group-hover:opacity-10 text-sky-500 transition-opacity"><svg xmlns="http://www.w3.org/2000/svg" class="w-32 h-32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="1" y="4" width="22" height="16" rx="2" ry="2"/><line x1="1" y1="10" x2="23" y2="10"/></svg></div>
                        <h3 class="text-[10px] font-black text-sky-500 uppercase tracking-widest mb-1">Masuk Mesin EDC (Debit)</h3>
                        <p class="text-2xl font-black text-sky-600 tracking-tighter">{{ formatRupiah(rekap.debit) }}</p>
                    </div>
                    <div class="bg-rose-50 p-5 rounded-3xl border border-rose-100 shadow-sm relative overflow-hidden group">
                        <div class="absolute -right-4 -bottom-4 opacity-[0.05] group-hover:opacity-10 text-rose-500 transition-opacity"><svg xmlns="http://www.w3.org/2000/svg" class="w-32 h-32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/><line x1="9" y1="12" x2="15" y2="12"/></svg></div>
                        <h3 class="text-[10px] font-black text-rose-500 uppercase tracking-widest mb-1">Piutang (Belum Bayar)</h3>
                        <p class="text-2xl font-black text-rose-600 tracking-tighter">{{ formatRupiah(rekap.piutang) }}</p>
                    </div>
                </div>

                <div class="bg-white border border-slate-200 rounded-3xl shadow-sm overflow-hidden flex flex-col h-fit">
                    <div class="p-5 border-b border-slate-100 flex flex-col sm:flex-row justify-between items-center gap-4">
                        <h2 class="text-sm font-black text-slate-800 uppercase tracking-widest">Riwayat Transaksi</h2>
                        <div class="relative w-full sm:w-72">
                            <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none"><svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg></div>
                            <input v-model="searchQuery" type="text" placeholder="Cari Resi / Pelanggan..." class="w-full pl-10 pr-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-indigo-500 outline-none font-bold text-xs text-slate-800 transition-all placeholder:text-slate-400">
                        </div>
                    </div>

                    <div class="overflow-x-auto">
                        <table class="w-full text-left border-collapse">
                            <thead>
                                <tr class="bg-slate-50 text-[10px] uppercase tracking-widest text-slate-400 border-b border-slate-200">
                                    <th class="p-4 font-black">Invoice & Waktu</th>
                                    <th class="p-4 font-black">Pelanggan</th>
                                    <th class="p-4 font-black">Pembayaran</th>
                                    <th class="p-4 font-black">Total</th>
                                    <th class="p-4 font-black text-center">Audit Bukti</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr v-if="isLoading" class="border-b border-slate-100">
                                    <td colspan="5" class="p-8 text-center text-slate-400 font-bold text-xs animate-pulse uppercase tracking-widest">Memuat Data Laporan...</td>
                                </tr>
                                <tr v-else-if="filteredRiwayat.length === 0" class="border-b border-slate-100">
                                    <td colspan="5" class="p-8 text-center text-slate-400 font-bold text-xs uppercase tracking-widest">Belum Ada Transaksi Ditemukan</td>
                                </tr>
                                <tr v-else v-for="trx in filteredRiwayat" :key="trx.id" class="border-b border-slate-100 hover:bg-slate-50/50 transition-colors">
                                    <td class="p-4">
                                        <p class="font-black text-xs text-slate-800 uppercase">{{ trx.no_invoice }}</p>
                                        <p class="text-[10px] font-bold text-slate-400 mt-1">{{ formatDate(trx.tanggal) }}</p>
                                    </td>
                                    <td class="p-4 font-bold text-xs text-slate-700 uppercase">{{ trx.pelanggan }}</td>
                                    <td class="p-4">
                                        <div class="flex items-center gap-2">
                                            <span class="text-[9px] font-black px-2 py-1 rounded-md uppercase tracking-wider"
                                                  :class="trx.metode_bayar === 'QRIS' ? 'bg-emerald-100 text-emerald-600' : (trx.metode_bayar === 'TUNAI' ? 'bg-slate-200 text-slate-700' : 'bg-sky-100 text-sky-600')">
                                                {{ trx.metode_bayar }}
                                            </span>
                                            <svg v-if="trx.status_bayar === 'LUNAS'" xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-emerald-500" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
                                            <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-rose-500" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
                                        </div>
                                    </td>
                                    <td class="p-4 font-black text-sm text-indigo-600 tracking-tight">{{ formatRupiah(trx.total) }}</td>
                                    <td class="p-4 text-center">
                                        <button v-if="trx.metode_bayar === 'QRIS' && trx.bukti_transfer" @click="bukaBuktiTransfer(trx.bukti_transfer)" class="inline-flex items-center gap-1.5 bg-emerald-50 hover:bg-emerald-500 text-emerald-600 hover:text-white border border-emerald-200 px-3 py-1.5 rounded-lg font-black text-[10px] uppercase tracking-widest transition-all active:scale-95 group">
                                            <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/></svg>
                                            Cek Bukti
                                        </button>
                                        <span v-else class="text-[10px] font-bold text-slate-300">-</span>
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </div>

            </div>
        </div>

        <Teleport to="body">
            <div v-if="showBuktiModal" class="fixed inset-0 z-[9999] flex items-center justify-center bg-slate-900/90 backdrop-blur-sm p-4 animate-[fadeIn_0.2s_ease-out]" @click.self="showBuktiModal = false">
                <div class="bg-white rounded-3xl shadow-2xl w-full max-w-sm overflow-hidden flex flex-col">
                    <div class="bg-slate-800 p-4 flex justify-between items-center text-white shrink-0">
                        <h3 class="font-black text-xs uppercase tracking-widest">Audit Bukti Transfer</h3>
                        <button @click="showBuktiModal = false" class="text-slate-400 hover:text-white bg-slate-700 hover:bg-rose-500 p-1.5 rounded-full transition-colors"><svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg></button>
                    </div>
                    <div class="p-6 bg-slate-100 flex justify-center items-center">
                        <img :src="selectedBuktiUrl" alt="Bukti Transfer" class="max-w-full max-h-[60vh] object-contain rounded-xl shadow-md border-4 border-white" />
                    </div>
                </div>
            </div>
        </Teleport>
    </SidebarLaundry>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { height: 6px; width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
@keyframes fadeIn { from { opacity: 0; transform: scale(0.95); } to { opacity: 1; transform: scale(1); } }
</style>