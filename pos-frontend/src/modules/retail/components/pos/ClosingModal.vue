<script setup>
import { computed } from "vue";

defineOptions({ name: "ClosingModal" });

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";

const props = defineProps({
  show: Boolean, 
  showReceiptClosing: Boolean, 
  pecahan: Object, 
  totalUangFisik: Number, 
  lastClosingData: Object, 
  currentSession: Object, 
  currentUser: Object, 
  storeData: Object, 
  storeLogo: String,
});

const emit = defineEmits(["close", "process-closing", "print-closing", "finish-closing"]);

const formatRupiah = (angka) => {
  if (angka === 0) return "0";
  if (!angka || isNaN(angka)) return "0";
  return new Intl.NumberFormat("id-ID").format(angka);
};

const formatDate = (dateString) => {
  if (!dateString) return "-";
  const date = new Date(dateString);
  if (isNaN(date.getTime())) return String(dateString);

  const pad = (n) => n.toString().padStart(2, "0");
  return `${pad(date.getDate())}.${pad(date.getMonth() + 1)}.${date.getFullYear()} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`;
};

// 🛡️ MAPPING DATA CORES DARI GO BACKEND SINKRON TOTAL
const modalAwal = computed(() => {
  return props.lastClosingData?.session?.modal_awal || 
         props.lastClosingData?.session?.ModalAwal || 
         props.lastClosingData?.data?.session?.modal_awal ||
         props.currentSession?.modal_awal || 
         props.currentSession?.ModalAwal || 0;
});

const netSales = computed(() => {
  return props.lastClosingData?.net_sales ?? props.lastClosingData?.NetSales ?? props.lastClosingData?.data?.net_sales ?? 0;
});

const totalTax = computed(() => {
  return props.lastClosingData?.total_tax ?? props.lastClosingData?.TotalTax ?? props.lastClosingData?.tax ?? props.lastClosingData?.data?.total_tax ?? 0;
});

const salesCash = computed(() => {
  return props.lastClosingData?.sales_cash ?? props.lastClosingData?.SalesCash ?? props.lastClosingData?.data?.sales_cash ?? 0;
});

const salesNonTunai = computed(() => {
  return props.lastClosingData?.sales_non_tunai ?? props.lastClosingData?.SalesNonTunai ?? props.lastClosingData?.sales_non_cash ?? props.lastClosingData?.data?.sales_non_tunai ?? 0;
});

const totalExpected = computed(() => {
  return props.lastClosingData?.total_expected ?? props.lastClosingData?.TotalExpected ?? props.lastClosingData?.data?.total_expected ?? 0;
});

const totalActual = computed(() => {
  return props.lastClosingData?.total_actual ?? props.lastClosingData?.TotalActual ?? props.lastClosingData?.data?.total_actual ?? 0;
});

const selisihKas = computed(() => {
  return props.lastClosingData?.selisih ?? props.lastClosingData?.Selisih ?? props.lastClosingData?.data?.selisih ?? 0;
});

const sanitizePecahanInput = (event, key) => {
  let cleanVal = event.target.value.replace(/\D/g, ""); 
  props.pecahan[key] = cleanVal ? parseInt(cleanVal, 10) : 0;
};

const triggerPrint = () => { window.print(); };
</script>

<template>
  <div v-if="show" class="fixed inset-0 bg-slate-900/90 flex items-center justify-center z-[160] p-4 backdrop-blur-sm no-print">
    <div class="bg-white rounded-[24px] md:rounded-[32px] w-full max-w-2xl overflow-hidden shadow-2xl border-[6px] md:border-[8px] border-slate-800 flex flex-col max-h-[95vh] md:max-h-[90vh]">
      
      <div class="bg-slate-800 p-4 md:p-6 text-center shrink-0 relative">
        <button @click="emit('close')" class="absolute top-4 right-4 text-slate-400 hover:text-white transition-colors">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
        </button>
        <h2 class="text-white font-black text-xl uppercase italic tracking-wider">Cash Count Drawer</h2>
        <p class="text-slate-400 text-[9px] font-bold uppercase tracking-widest mt-1">Audit Serah Terima Laci Kasir Cabang</p>
      </div>

      <div class="p-4 md:p-6 grid grid-cols-2 md:grid-cols-3 gap-2 md:gap-3 overflow-y-auto custom-scrollbar flex-1 bg-slate-50">
        <div v-for="denon in [{ label: '100.000', key: 'p100k' }, { label: '50.000', key: 'p50k' }, { label: '20.000', key: 'p20k' }, { label: '10.000', key: 'p10k' }, { label: '5.000', key: 'p5k' }, { label: '2.000', key: 'p2k' }, { label: '1.000', key: 'p1k' }, { label: '500', key: 'p500' }, { label: '200', key: 'p200' }, { label: '100', key: 'p100' }, { label: '50', key: 'p50' }]" :key="denon.key" class="bg-white p-2 md:p-3 rounded-xl border-2 border-slate-200 focus-within:border-indigo-500 shadow-sm transition-all">
          <label class="text-[8px] md:text-[9px] font-black text-slate-400 block mb-1 uppercase tracking-widest">Pecahan Rp {{ denon.label }}</label>
          <input 
            type="text" 
            inputmode="numeric"
            :value="pecahan[denon.key] || ''"
            @input="sanitizePecahanInput($event, denon.key)"
            class="w-full bg-transparent text-lg font-black text-slate-800 outline-none" 
            placeholder="0" 
          />
        </div>
      </div>

      <div class="p-4 md:p-6 bg-white border-t border-slate-200 shrink-0">
        <div class="flex justify-between items-center mb-4">
          <span class="text-slate-500 text-[9px] font-black uppercase tracking-widest">Total Fisik Aktual:</span>
          <span class="text-2xl font-black text-indigo-700">Rp {{ formatRupiah(totalUangFisik) }}</span>
        </div>
        <button @click="emit('process-closing')" class="w-full bg-indigo-600 hover:bg-indigo-700 text-white py-3.5 rounded-xl font-black text-xs tracking-widest flex items-center justify-center gap-2 uppercase shadow-lg shadow-indigo-200 active:scale-95 transition-all">Proses Tutup Shift Kasir</button>
      </div>
    </div>
  </div>

  <div v-if="showReceiptClosing" class="fixed inset-0 bg-slate-950/90 flex items-center justify-center z-[200] p-4 backdrop-blur-sm print:static print:bg-white print:p-0 print:block">
    <div class="bg-white p-6 md:p-8 rounded-[32px] md:rounded-[40px] shadow-2xl w-full max-w-sm flex flex-col max-h-[90vh] border-[6px] md:border-[8px] border-slate-800 print:border-none print:shadow-none print:max-h-none print:max-w-none print:rounded-none print:m-0 print:p-0">
      
      <div class="overflow-y-auto custom-scrollbar bg-white p-2 mx-auto print:overflow-visible print:mx-0 print:px-3 print:py-2 select-none" id="print-area">
        <div class="text-center mb-4 font-mono leading-none">
          
          <div v-if="storeData?.logo_url && storeData.logo_url !== ''" class="mb-2">
            <img :src="(storeData.logo_url.startsWith('http://') || storeData.logo_url.startsWith('https://')) ? storeData.logo_url : API_BASE_URL + storeData.logo_url" class="w-16 h-16 object-contain mx-auto grayscale contrast-200 brightness-90" alt="Logo Toko" />
          </div>
          <div v-else-if="storeLogo && storeLogo !== ''" class="mb-2">
            <img :src="(storeLogo.startsWith('http://') || storeLogo.startsWith('https://')) ? storeLogo : API_BASE_URL + storeLogo" class="w-16 h-16 object-contain mx-auto grayscale contrast-200 brightness-90" alt="Logo Toko" />
          </div>
          
          <div v-else class="font-black text-[12px] uppercase mb-1 italic tracking-tighter">{{ storeData?.nama_toko || storeData?.NamaToko || "NEXA POS STORE" }}</div>
          <p class="text-[9px] font-black uppercase tracking-tight opacity-100 leading-tight px-1">
            {{ storeData?.alamat || storeData?.Alamat || "JAKARTA, INDONESIA" }}<br />
            {{ storeData?.kelurahan || storeData?.Kelurahan || "" }} {{ storeData?.kecamatan || storeData?.Kecamatan || "" }}<br />
            {{ storeData?.kota || storeData?.Kota || "" }} {{ storeData?.kode_pos || storeData?.KodePos || "" }}
          </p>
        </div>

        <div class="text-[9px] uppercase font-mono font-black mb-2 leading-tight space-y-0.5">
          <div>TANGGAL : {{ formatDate(lastClosingData?.end_time || lastClosingData?.created_at || lastClosingData?.data?.end_time) }}</div>
          <div class="text-center font-black my-3 tracking-[0.15em] border-y border-black py-1 uppercase">SLIP PENJUALAN<br />TUTUP SHIFT KASIR</div>
          <div class="flex justify-between"><span>KASIR: {{ currentUser?.name || "OPERATOR" }}</span><span>STASIUN: POS-{{ currentSession?.station_number || currentSession?.StationNumber || "01" }}</span></div>
        </div>

        <div class="border-b border-black border-dashed mt-2 mb-1"></div>
        <div class="font-black font-mono text-[9px] uppercase tracking-wider">OMSET PRODUK DAGANGAN</div>
        <div class="border-b border-black border-dashed mb-1 mt-1"></div>

        <div class="text-[9px] uppercase font-mono font-black flex flex-col gap-0.5">
  <div class="flex items-center">
    <span class="w-[150px] shrink-0">Penjualan Bersih (Net)</span>
    <span class="mr-1">:</span>
    <span class="flex-1 text-right">{{ formatRupiah(netSales) }}</span>
  </div>
  <div class="flex items-center">
    <span class="w-[150px] shrink-0">Pajak Resto/PPN</span>
    <span class="mr-1">:</span>
    <span class="flex-1 text-right">{{ formatRupiah(totalTax) }}</span>
  </div>
  <div class="border-t border-black border-dotted my-1"></div>
  <div class="flex items-center font-black text-[10px]">
    <span class="w-[150px] shrink-0">TOTAL OMSET BRUTO</span>
    <span class="mr-1">:</span>
    <span class="flex-1 text-right">{{ formatRupiah(netSales + totalTax) }}</span>
  </div>
</div>

<div class="text-[9px] uppercase font-mono font-black flex flex-col gap-0.5 mb-2">
  <div class="flex items-center">
    <span class="w-[150px] shrink-0">Total Pendapatan Tunai</span>
    <span class="mr-1">:</span>
    <span class="flex-1 text-right">{{ formatRupiah(salesCash) }}</span>
  </div>
  <div class="flex items-center">
    <span class="w-[150px] shrink-0">Total QRIS / Non-Tunai</span>
    <span class="mr-1">:</span>
    <span class="flex-1 text-right">{{ formatRupiah(salesNonTunai) }}</span>
  </div>
</div>

<div class="text-[9px] uppercase font-mono font-black flex flex-col gap-0.5 mb-4">
  <div class="flex items-center">
    <span class="w-[150px] shrink-0">[+] Modal Awal Sesi</span>
    <span class="mr-1">:</span>
    <span class="flex-1 text-right">{{ formatRupiah(modalAwal) }}</span>
  </div>
  <div class="flex items-center">
    <span class="w-[150px] shrink-0">[+] Penjualan Tunai Shift</span>
    <span class="mr-1">:</span>
    <span class="flex-1 text-right">{{ formatRupiah(salesCash) }}</span>
  </div>
  
  <div class="border-t border-black border-dotted my-1"></div>
  
  <div class="flex items-center font-bold">
    <span class="w-[150px] shrink-0">(=) Uang Sistem Seharusnya</span>
    <span class="mr-1">:</span>
    <span class="flex-1 text-right">{{ formatRupiah(totalExpected) }}</span>
  </div>
  
  <div class="flex items-center font-bold text-indigo-700 bg-slate-100 px-0.5 rounded py-0.5 my-0.5">
    <span class="w-[150px] shrink-0">(X) Duit Fisik Aktual Kasir</span>
    <span class="mr-1">:</span>
    <span class="flex-1 text-right">{{ formatRupiah(totalActual) }}</span>
  </div>
  
  <div class="border-t border-black border-dotted my-1"></div>
  
  <div class="flex items-center font-black text-[10px]">
    <span class="w-[150px] shrink-0">Selisih / Variance Kas</span>
    <span class="mr-1">:</span>
    <span :class="selisihKas === 0 ? 'text-black' : (selisihKas > 0 ? 'text-emerald-600' : 'text-rose-600')" class="flex-1 text-right">
      {{ formatRupiah(selisihKas) }} 
      <span class="text-[7px] block sm:inline font-bold">
        {{ selisihKas === 0 ? '(BALANCED)' : (selisihKas > 0 ? '(PLUS)' : '(MINUS)') }}
      </span>
    </span>
  </div>
</div>

        <div class="text-center mt-6 pb-2 font-black font-mono text-[9px] leading-tight uppercase break-inside-avoid">
          <div class="mb-8 tracking-widest">VALIDASI OTORISASI SHIFT</div>
          <div class="w-2/3 mx-auto border-t border-black pt-1 font-black">{{ currentUser?.name || "KASIR OPERATOR" }}</div>
          <div class="text-[8px] mt-1 tracking-widest">( K A S I R )</div>
        </div>
      </div>

      <div class="mt-4 md:mt-6 flex flex-col gap-2 md:gap-3 shrink-0 print:hidden">
        <button @click="triggerPrint" class="w-full bg-indigo-600 text-white py-3 md:py-4 rounded-xl md:rounded-2xl font-black text-[9px] md:text-[10px] uppercase tracking-[0.2em] shadow-xl shadow-indigo-200 active:scale-95 transition-all">Cetak Dokumen Closing</button>
        <button @click="emit('finish-closing')" class="w-full bg-slate-100 text-slate-500 py-3 md:py-3.5 rounded-xl md:rounded-2xl font-black text-[9px] md:text-[10px] uppercase tracking-widest hover:bg-slate-200 active:scale-95 transition-all">Tutup</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }

#print-area {
  width: v-bind("storeData?.printer_width || storeData?.PrinterWidth || '58mm'") !important;
  max-width: 100% !important;
}

@media print {
  #print-area {
    width: v-bind("storeData?.printer_width || storeData?.PrinterWidth || '58mm'") !important;
    margin: 0 auto !important; 
    padding: 0 !important;
    background: white !important;
  }
}
</style>