<script setup>
import { computed } from 'vue';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

const props = defineProps({
    show: Boolean,
    showReceiptClosing: Boolean,
    pecahan: Object,
    totalUangFisik: Number,
    lastClosingData: Object,
    currentSession: Object,
    currentUser: Object,
    storeData: Object 
});

const emit = defineEmits(['close', 'process-closing', 'print-closing', 'finish-closing']);

// 🚀 FUNGSI FORMAT RUPIAH
const formatRupiah = (angka) => {
    if (!angka && angka !== 0) return '0';
    return new Intl.NumberFormat('id-ID').format(angka);
};

// 🚀 FUNGSI SAKTI FORMAT TANGGAL JADI STANDAR RITEL (DD.MM.YYYY HH:MM:SS)
const formatDate = (dateString) => {
    if (!dateString) return '-';
    const date = new Date(dateString);
    if (isNaN(date.getTime())) return String(dateString); 
    
    const pad = (n) => n.toString().padStart(2, '0');
    const d = pad(date.getDate());
    const m = pad(date.getMonth() + 1);
    const y = date.getFullYear();
    const h = pad(date.getHours());
    const min = pad(date.getMinutes());
    const s = pad(date.getSeconds());
    
    return `${d}.${m}.${y} ${h}:${min}:${s}`;
};

// 🚀 PERHITUNGAN MODAL AWAL (FLOAT) JAGA-JAGA KALAU BACKEND LUPA NGIRIM
const modalAwal = computed(() => {
    if (props.lastClosingData?.session?.modal_awal) return props.lastClosingData.session.modal_awal;
    const expected = props.lastClosingData?.total_expected || 0;
    const cash = props.lastClosingData?.sales_cash || 0;
    return expected - cash;
});

const triggerPrint = () => {
    window.print();
};
</script>

<template>
    <div v-if="show" class="fixed inset-0 bg-slate-900/90 flex items-center justify-center z-[160] p-4 backdrop-blur-sm no-print">
        <div class="bg-white rounded-[24px] md:rounded-[32px] w-full max-w-2xl overflow-hidden shadow-2xl border-[6px] md:border-[8px] border-slate-800 flex flex-col max-h-[95vh] md:max-h-[90vh]">
            <div class="bg-slate-800 p-4 md:p-6 text-center shrink-0 relative">
                <button @click="emit('close')" class="absolute top-4 right-4 text-slate-400 hover:text-white">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
                </button>
                <h2 class="text-white font-black text-xl uppercase italic">Cash Count</h2>
                <p class="text-slate-400 text-[9px] font-bold uppercase tracking-widest mt-1">Hitung Uang Fisik Laci Kasir</p>
            </div>

            <div class="p-4 md:p-6 grid grid-cols-2 md:grid-cols-3 gap-2 md:gap-3 overflow-y-auto custom-scrollbar flex-1 bg-slate-50">
                <div v-for="denon in [
                    { label: '100.000', key: 'p100k' }, { label: '50.000', key: 'p50k' }, { label: '20.000', key: 'p20k' },
                    { label: '10.000', key: 'p10k' }, { label: '5.000', key: 'p5k' }, { label: '2.000', key: 'p2k' },
                    { label: '1.000', key: 'p1k' }, { label: '500', key: 'p500' }, { label: '200', key: 'p200' },
                    { label: '100', key: 'p100' }, { label: '50', key: 'p50' }
                ]" :key="denon.key" class="bg-white p-2 md:p-3 rounded-xl border-2 border-slate-200 focus-within:border-indigo-500">
                    <label class="text-[8px] md:text-[9px] font-black text-slate-400 block mb-1">Pecahan {{ denon.label }}</label>
                    <input type="number" v-model.number="pecahan[denon.key]" class="w-full bg-transparent text-lg font-black text-slate-800 outline-none" placeholder="0" min="0">
                </div>
            </div>

            <div class="p-4 md:p-6 bg-white border-t border-slate-200 shrink-0">
                <div class="flex justify-between items-center mb-4">
                    <span class="text-slate-500 text-[9px] font-black uppercase">Total Fisik:</span>
                    <span class="text-2xl font-black text-indigo-700">Rp {{ formatRupiah(totalUangFisik) }}</span>
                </div>
                <button @click="emit('process-closing')" class="w-full bg-indigo-600 text-white py-3.5 rounded-xl font-black text-xs tracking-widest flex items-center justify-center gap-2 uppercase shadow-lg shadow-indigo-200 active:scale-95 transition-all">
                    Proses Tutup Shift
                </button>
            </div>
        </div>
    </div>

    <div v-if="showReceiptClosing" class="fixed inset-0 bg-slate-950/90 flex items-center justify-center z-[200] p-4 backdrop-blur-sm print:static print:bg-white print:p-0 print:block">
        <div class="bg-white p-6 md:p-8 rounded-[32px] md:rounded-[40px] shadow-2xl w-full max-w-sm flex flex-col max-h-[90vh] border-[6px] md:border-[8px] border-slate-800 print:border-none print:shadow-none print:max-h-none print:max-w-none print:rounded-none print:m-0 print:p-0">
            
            <div class="overflow-y-auto custom-scrollbar bg-white p-2 mx-auto print:overflow-visible print:mx-0 print:px-3 print:py-2" id="print-area" :style="{ width: storeData?.printer_width || '58mm' }">
                
                <div class="text-center mb-4 font-mono leading-none">
                    <div v-if="storeData?.logo_url && storeData.logo_url !== ''">
                        <img :src="API_BASE_URL + storeData.logo_url" class="w-20 h-20 object-contain mx-auto grayscale contrast-150 mb-2" alt="Logo Toko" />
                    </div>
                    <div v-else class="font-black text-[12px] uppercase mb-1">
                        {{ storeData?.nama_toko || 'NEXA POS STORE' }}
                    </div>
                    <p class="text-[9px] font-black uppercase tracking-widest opacity-100 leading-tight px-1">
                        {{ storeData?.alamat || 'JAKARTA, INDONESIA' }}<br />
                        {{ storeData?.kelurahan || '' }} {{ storeData?.kecamatan || '' }}<br />
                        {{ storeData?.kota || '' }} {{ storeData?.kode_pos || '' }}
                    </p>
                </div>
        
                <div class="text-[9px] uppercase font-mono font-black mb-2 leading-tight">
                    <div>TANGGAL : {{ formatDate(lastClosingData?.end_time || lastClosingData?.created_at) }}</div>
                    <div class="text-center font-bold my-2 tracking-widest">
                        SLIP PENJUALAN<br/>TUTUP SHIFT
                    </div>
                    <div>Nama: {{ currentUser?.name?.split(' ')[0] || '1' }} || Station : {{ currentSession?.station_number || '01' }}</div>
                </div>
        
                <div class="border-b border-black border-dashed mt-2 mb-1"></div>
                <div class="font-black font-mono text-[9px] uppercase tracking-widest">BARANG DAGANGAN</div>
                <div class="border-b border-black border-dashed mb-1 mt-1"></div>
        
                <div class="text-[9px] uppercase font-mono font-black flex flex-col gap-0.5">
                    <div class="flex justify-between">
                        <span>Penjualan (Net)</span>
                        <span>: {{ formatRupiah(lastClosingData?.net_sales) }}</span>
                    </div>
                    <div class="flex justify-between">
                        <span>PPN</span>
                        <span>: {{ formatRupiah(lastClosingData?.total_tax) }}</span>
                    </div>
                    <div class="flex justify-between mt-1">
                        <span>TOTAL TRX</span>
                        <span>: {{ formatRupiah((lastClosingData?.net_sales || 0) + (lastClosingData?.total_tax || 0)) }}</span>
                    </div>
                </div>

                <div class="border-b border-black border-dashed mt-3 mb-1"></div>
                <div class="font-black font-mono text-[9px] uppercase tracking-widest">Rincian Transaksi</div>
                <div class="border-b border-black border-dashed mb-1 mt-1"></div>
                
                <div class="text-[9px] uppercase font-mono font-black flex flex-col gap-0.5 mb-2">
                    <div class="flex justify-between">
                        <span>Uang Tunai</span>
                        <span>: {{ formatRupiah(lastClosingData?.sales_cash) }}</span>
                    </div>
                    <div class="flex justify-between">
                        <span>Non Tunai</span>
                        <span>: {{ formatRupiah(lastClosingData?.sales_non_tunai) }}</span>
                    </div>
                </div>

                <div class="border-b border-black border-dashed mt-3 mb-1"></div>
                <div class="font-black font-mono text-[9px] uppercase tracking-widest">REKONSILIASI LACI</div>
                <div class="border-b border-black border-dashed mb-1 mt-1"></div>
        
                <div class="text-[9px] uppercase font-mono font-black flex flex-col gap-0.5 mb-4">
                    <div class="flex justify-between">
                        <span>Modal Awal (Float)</span>
                        <span>: {{ formatRupiah(modalAwal) }}</span>
                    </div>
                    <div class="flex justify-between">
                        <span>Expected Cash</span>
                        <span>: {{ formatRupiah(lastClosingData?.total_expected) }}</span>
                    </div>
                    <div class="flex justify-between">
                        <span>Fisik Aktual</span>
                        <span>: {{ formatRupiah(lastClosingData?.total_actual) }}</span>
                    </div>
                    <div class="flex justify-between mt-1 pt-1 border-t border-black border-dotted">
                        <span>Variance</span>
                        <span>: {{ formatRupiah(lastClosingData?.selisih) }}</span>
                    </div>
                </div>
        
                <div class="text-center font-black font-mono text-[9px] mt-6 pb-2 leading-tight uppercase">
                    <div class="mb-8">TANDA TANGAN</div>
                    <div class="w-2/3 mx-auto border-t border-black pt-1">{{ currentUser?.name || 'KASIR' }}</div>
                    <div class="text-[8px] mt-1">( K A S I R )</div>
                </div>
            </div>
            
            <div class="mt-4 md:mt-6 flex flex-col gap-2 md:gap-3 shrink-0 print:hidden">
                <button @click="triggerPrint" class="w-full bg-indigo-600 text-white py-3 md:py-4 rounded-xl md:rounded-2xl font-black text-[9px] md:text-[10px] uppercase tracking-[0.2em] shadow-xl shadow-indigo-200 active:scale-95 transition-all">
                    Cetak Struk
                </button>
                <button @click="emit('close')" class="w-full bg-slate-100 text-slate-500 py-3 md:py-3.5 rounded-xl md:rounded-2xl font-black text-[9px] md:text-[10px] uppercase tracking-widest hover:bg-slate-200 active:scale-95 transition-all">
                    Tutup
                </button>
            </div>
        </div>
    </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #94a3b8; }

@media print {
  #print-area {
    width: v-bind("storeData?.printer_width || '58mm'") !important;
    margin: 0 !important;
    padding: 0 !important;
  }
}
</style>