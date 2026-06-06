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

// 🚀 FIX SECURITY FORMULA: Ambil modal awal dari session database murni bray, anti-rumus tebakan mundur bray!
const modalAwal = computed(() => {
  return props.lastClosingData?.session?.modal_awal || 
         props.lastClosingData?.session?.ModalAwal || 
         props.currentSession?.modal_awal || 
         props.currentSession?.ModalAwal || 0;
});

// 🛡️ SECURITY LAYER INPUT PECAHAN: Blokir mutlak angka minus, desimal, dan huruf eksponensial bray!
const sanitizePecahanInput = (event, key) => {
  let cleanVal = event.target.value.replace(/\D/g, ""); // Murni angka bulat positif, buang minus bray!
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
          
          <div v-else class="font-black text-[12px] uppercase mb-1 italic tracking-tighter">{{ storeData?.nama_toko || storeData?.NamaToko || "NEXA POS STORE" }}</div>
          <p class="text-[9px] font-black uppercase tracking-tight opacity-100 leading-tight px-1">
            {{ storeData?.alamat || "JAKARTA, INDONESIA" }}<br />
            {{ storeData?.kelurahan || "" }} {{ storeData?.kecamatan || "" }}<br />
            {{ storeData?.kota || "" }} {{ storeData?.kode_pos || "" }}
          </p>
        </div>

        <div class="text-[9px] uppercase font-mono font-black mb-2 leading-tight space-y-0.5">
          <div>TANGGAL : {{ formatDate(lastClosingData?.end_time || lastClosingData?.created_at) }}</div>
          <div class="text-center font-black my-3 tracking-[0.15em] border-y border-black py-1 uppercase">SLIP PENJUALAN<br />TUTUP SHIFT KASIR</div>
          <div class="flex justify-between"><span>KASIR: {{ currentUser?.name?.split(" ")[0] || "OPERATOR" }}</span><span>STASIUN: POS-{{ currentSession?.station_number || currentSession?.StationNumber || "01" }}</span></div>
        </div>

        <div class="border-b border-black border-dashed mt-2 mb-1"></div>
        <div class="font-black font-mono text-[9px] uppercase tracking-wider">OMSET PRODUK DAGANGAN</div>
        <div class="border-b border-black border-dashed mb-1 mt-1"></div>

        <div class="text-[9px] uppercase font-mono font-black flex flex-col gap-0.5">
          <div class="flex justify-between"><span>Penjualan Bersih (Net)</span><span>: {{ formatRupiah(lastClosingData?.net_sales || lastClosingData?.net_sales_cash) }}</span></div>
          <div class="flex justify-between"><span>Pajak Resto/PPN</span><span>: {{ formatRupiah(lastClosingData?.total_tax || lastClosingData?.tax) }}</span></div>
          <div class="flex justify-between mt-1 pt-1 border-t border-black border-dotted font-black text-[10px]"><span>TOTAL OMSET BRUTO</span><span>: {{ formatRupiah((lastClosingData?.net_sales || 0) + (lastClosingData?.total_tax || 0)) }}</span></div>
        </div>

        <div class="border-b border-black border-dashed mt-3 mb-1"></div>
        <div class="font-black font-mono text-[9px] uppercase tracking-wider">RINCIAN METODE TERIMA DANA</div>
        <div class="border-b border-black border-dashed mb-1 mt-1"></div>

        <div class="text-[9px] uppercase font-mono font-black flex flex-col gap-0.5 mb-2">
          <div class="flex justify-between"><span>Total Pendapatan Tunai</span><span>: {{ formatRupiah(lastClosingData?.sales_cash) }}</span></div>
          <div class="flex justify-between"><span>Total QRIS / EDC Non-Tunai</span><span>: {{ formatRupiah(lastClosingData?.sales_non_tunai || lastClosingData?.sales_non_cash) }}</span></div>
        </div>

        <div class="border-b border-black border-dashed mt-3 mb-1"></div>
        <div class="font-black font-mono text-[9px] uppercase tracking-wider">REKONSILIASI AUDIT LACI</div>
        <div class="border-b border-black border-dashed mb-1 mt-1"></div>

        <div class="text-[9px] uppercase font-mono font-black flex flex-col gap-0.5 mb-4">
          <div class="flex justify-between"><span>Modal Awal Laci (Float)</span><span>: {{ formatRupiah(modalAwal) }}</span></div>
          <div class="flex justify-between"><span>Expected Cash (Sistem)</span><span>: {{ formatRupiah(lastClosingData?.total_expected) }}</span></div>
          <div class="flex justify-between font-black"><span>Fisik Aktual Kasir</span><span>: {{ formatRupiah(lastClosingData?.total_actual) }}</span></div>
          
          <div class="flex justify-between mt-1 pt-1 border-t border-black border-dotted font-black text-[10px]">
            <span>Variance / Selisih Kas</span>
            <span :class="(lastClosingData?.selisih || 0) === 0 ? 'text-black' : 'text-black font-black underline'">
              : {{ formatRupiah(lastClosingData?.selisih) }} {{ (lastClosingData?.selisih || 0) === 0 ? '(BALANCED)' : '' }}
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
        <button @click="emit('close')" class="w-full bg-slate-100 text-slate-500 py-3 md:py-3.5 rounded-xl md:rounded-2xl font-black text-[9px] md:text-[10px] uppercase tracking-widest hover:bg-slate-200 active:scale-95 transition-all">Tutup</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }

/* FIX PRINTER SCALING DESIGN SLIP REKAP BRAY */
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