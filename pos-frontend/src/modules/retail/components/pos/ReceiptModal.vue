<script setup>

import { computed } from 'vue';

  const API_BASE_URL =
    import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

  const props = defineProps({
    show: Boolean,
    invoiceData: Object,
    storeData: Object,
    cashierName: String,
    stationNumber: String,
  });

  const emit = defineEmits(['close']);

  const formatRupiah = (angka) => {
    if (typeof angka !== 'number' || isNaN(angka)) return '0';
    return new Intl.NumberFormat('id-ID').format(angka);
  };

  // 🚀 1. TAMPILKAN SUBTOTAL (TOTAL SEBELUM PAJAK) - PALING AKURAT
const subTotalDisplay = computed(() => {
    // Cek field asli dari executeCheckout (subtotal) atau dari Golang (sub_total)
    if (props.invoiceData.subtotal !== undefined) return Number(props.invoiceData.subtotal);
    if (props.invoiceData.sub_total !== undefined) return Number(props.invoiceData.sub_total);

    // JURUS PAMUNGKAS: Hitung langsung dari harga barang x qty! 
    // Ini jaminan 100% akurat dan gak bakal kecampur sama total akhir.
    const items = props.invoiceData.cart || props.invoiceData.details || [];
    if (items.length > 0) {
        return items.reduce((sum, item) => {
            const harga = item.price || item.harga_satuan || 0;
            const qty = item.qty || item.kuantitas || 0;
            return sum + (harga * qty);
        }, 0);
    }

    return 0; // Fallback paling aman
});

// 🚀 2. TAMPILKAN PAJAK
const totalPajak = computed(() => {
    // Kalo dari riwayat dan udah ada field pajaknya
    if (props.invoiceData.pajak !== undefined) return Number(props.invoiceData.pajak);
    
    // Kalo kasir live, hitung dari subTotalDisplay yang udah fix bener
    if (!props.storeData?.is_tax_active) return 0;
    const persen = props.storeData?.pajak_persen || 0;
    
    return subTotalDisplay.value * (persen / 100);
});

// 🚀 3. TAMPILKAN TOTAL AKHIR
const grandTotal = computed(() => {
    // Langsung tembak total tagihan akhir dari Golang / object checkout
    if (props.invoiceData.total !== undefined) return Number(props.invoiceData.total);
    if (props.invoiceData.total_harga !== undefined) return Number(props.invoiceData.total_harga);

    // Fallback buat jaga-jaga
    return subTotalDisplay.value + totalPajak.value;
});

  const formatWA = (waNumber) => {
    if (!waNumber) return '-';
    const cleaned = waNumber.replace(/\D/g, '');
    let localNumber = cleaned;
    if (localNumber.startsWith('62')) {
      localNumber = '0' + localNumber.slice(2);
    }
    const match = localNumber.match(/^(\d{4})(\d{4})(\d{4,5})$/);
    if (match) {
      return `${match[1]}-${match[2]}-${match[3]}`;
    }
    return localNumber;
  };

  // FUNGSI SILUMAN PEMBACA KEMASAN (VERSI UPDATE UNTUK usePos.js)
  const formatKemasan = (item) => {
    // 1. Kalo ini dari history laporan yang udah ada detail_notes dari Golang
    if (item.detail_notes && item.detail_notes !== 'Transaksi Retail Toko') {
      return item.detail_notes;
    }

    // 2. BACA LANGSUNG DARI KERANJANG KASIR (selected_uom)
    // Pas kasir nge-toggle grosir di CartSidebar, dia ngubah "item.selected_uom"
    const satuanPilihan =
      item.selected_uom || item.satuan_terpilih || item.satuan || item.kemasan;
    const qtyPilihan = item.qty || item.kuantitas || 0;

    if (satuanPilihan) {
      return `${qtyPilihan} ${satuanPilihan}`;
    }

    // 3. Fallback (Kalau gak ada datanya sama sekali)
    const product = item.product || item;
    return `${qtyPilihan} ${product.satuan_dasar || 'pcs'}`;
  };

  const triggerPrint = () => {
    window.print();
  };
</script>

<template>
  <div
    v-if="show && invoiceData"
    class="fixed inset-0 bg-slate-950/90 flex items-center justify-center z-[200] p-4 backdrop-blur-sm print:static print:bg-white print:p-0 print:block"
  >
    <div
      class="bg-white p-6 md:p-8 rounded-[32px] md:rounded-[40px] shadow-2xl w-full max-w-sm flex flex-col max-h-[90vh] border-[6px] md:border-[8px] border-slate-800 print:border-none print:shadow-none print:max-h-none print:max-w-none print:rounded-none print:m-0 print:p-0"
    >
      <div
        class="overflow-y-auto custom-scrollbar bg-white p-2 mx-auto print:overflow-visible print:mx-0 print:px-3 print:py-2"
        id="print-area"
        :style="{  width: storeData?.printer_width || '58mm' }"
      >
        <div class="text-center mb-4 font-mono leading-none">
          <div v-if="storeData?.logo_url && storeData.logo_url !== ''">
            <img
              :src="API_BASE_URL + storeData.logo_url"
              class="w-20 h-20 object-contain mx-auto grayscale contrast-150"
              alt="Logo Toko"
            />
          </div>
          <h2
            v-else
            class="font-black text-sm uppercase tracking-tighter mb-1 italic"
          >
            {{ storeData?.nama_toko || 'NEXA POS STORE' }}
          </h2>
          <p
            class="text-[9px] font-black uppercase tracking-widest opacity-100 leading-tight px-1"
          >
            {{ storeData?.alamat || 'JAKARTA, INDONESIA' }}<br />
            {{ storeData?.kelurahan || 'KELURAHAN' }},
            {{ storeData?.kecamatan || 'KECAMATAN' }}<br />
            {{ storeData?.kota || 'KOTA' }},
            {{ storeData?.provinsi || 'PROVINSI' }}
            {{ storeData?.kode_pos || 'KODE POS' }}
          </p>
          <!-- <p v-if="storeData?.Telepon || storeData?.telepon" class="text-[9px] font-black uppercase tracking-widest border-slate-600 pt-1 border-dotted inline-block">
                        WA: {{ formatWA(storeData?.Telepon || storeData?.telepon) }}
                    </p> -->
        </div>

        <div
          class="border-y border-black py-1.5 text-center font-black mb-3 font-mono text-[9px] tracking-[0.2em] uppercase"
        >
          {{ invoiceData.created_at ? 'Invoice Reprint' : 'Struk Belanja' }}
        </div>

        <div class="mb-3 text-[8px] font-black font-mono uppercase space-y-0.5">
          <div class="flex justify-between">
            <span>WAKTU:</span>
            <span>
              {{
                invoiceData.date ||
                (invoiceData.created_at
                  ? new Date(invoiceData.created_at)
                      .toLocaleString('id-ID', {
                        year: '2-digit',
                        month: '2-digit',
                        day: '2-digit',
                        hour: '2-digit',
                        minute: '2-digit',
                        hour12: false,
                      })
                      .replace('.', ':')
                  : '-')
              }}
            </span>
          </div>
          <div class="flex justify-between">
            <span>KASIR:</span>
            <span
              >{{ cashierName || 'KASIR' }} / POS-{{
                stationNumber || '01'
              }}</span
            >
          </div>
          <div class="flex justify-between">
            <span>Inv:</span>
            <span>{{ invoiceData.invoice || invoiceData.no_invoice }}</span>
          </div>
        </div>

        <div class="border-b border-black border-dashed mb-2"></div>

        <div
          v-for="item in invoiceData.cart || invoiceData.details"
          :key="item.id"
          class="mb-2 font-bold font-mono text-[9px] leading-tight uppercase"
        >
          <div class="truncate w-full pr-2">
            {{ item.name || item.product?.nama_produk || 'Item Belanja' }}
          </div>

          <div class="flex justify-between pl-2 text-[8px] mt-0.5">
            <span class="text-black-600">
              {{ formatKemasan(item) }}
              <span class="lowercase">
                x {{ formatRupiah(item.price || item.harga_satuan) }}</span
              >
            </span>
            <span class="font-black text-[9px]">{{
              formatRupiah(item.sub_total || item.price * item.qty)
            }}</span>
          </div>
        </div>

        <div class="border-t border-black border-dashed mt-2 pt-2"></div>

        <div class="flex justify-between font-bold text-[10px] mb-1 font-mono uppercase italic">
            <span>SUBTOTAL:</span>
            <span>{{ formatRupiah(subTotalDisplay) }}</span>
        </div>

        <div v-if="totalPajak > 0" class="flex justify-between font-bold text-[10px] mb-1 font-mono uppercase italic">
            <span>PAJAK:</span>
            <span>{{ formatRupiah(totalPajak) }}</span>
        </div>

        <div class="flex justify-between font-black text-[11px] mb-2 font-mono uppercase italic border-t border-black pt-1 mt-1">
            <span>TOTAL BELANJA:</span>
            <span>{{ formatRupiah(grandTotal) }}</span>
        </div>

        <div class="border-b border-black border-dashed mb-2"></div>

        <div
          class="flex justify-between mb-1 font-bold font-mono text-[10px] uppercase"
        >
          <span
            >BAYAR ({{
              invoiceData.method || invoiceData.metode_bayar || 'CASH'
            }}):</span
          >
          <span>{{
            formatRupiah(
              invoiceData.pay !== undefined
                ? invoiceData.pay
                : invoiceData.nominal_bayar
            )
          }}</span>
        </div>

        <div
          class="flex justify-between font-black font-mono text-[10px] uppercase italic text-black"
        >
          <span>KEMBALI:</span>
          <span
            >{{
              formatRupiah(
                invoiceData.total !== undefined
                  ? invoiceData['return']
                  : invoiceData.kembalian
              )
            }}</span
          >
        </div>

        <div
          class="text-center mt-4 font-black font-mono text-[10px] border-2 border-black p-1.5 uppercase leading-tight"
        >
          {{ storeData?.receipt_footer || 'TERIMA KASIH ATAS KUNJUNGAN ANDA!' }}
        </div>
      </div>

      <div
        class="mt-4 md:mt-6 flex flex-col gap-2 md:gap-3 shrink-0 print:hidden"
      >
        <button
          @click="triggerPrint"
          class="w-full bg-indigo-600 text-white py-3 md:py-4 rounded-xl md:rounded-2xl font-black text-[9px] md:text-[10px] uppercase tracking-[0.2em] shadow-xl shadow-indigo-200 flex items-center justify-center gap-2 active:scale-95 transition-all"
        >
          Cetak Struk
        </button>
        <button
          @click="emit('close')"
          class="w-full bg-slate-100 text-slate-500 py-3 md:py-3.5 rounded-xl md:rounded-2xl font-black text-[9px] md:text-[10px] uppercase tracking-widest hover:bg-slate-200 transition-all"
        >
          Tutup
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
@media print {
  #print-area {
    /* 🚀 Kunci ukuran kertas saat diprint berdasarkan database */
    width: v-bind("storeData?.printer_width || '58mm'") !important;
    margin: 0 !important;
    padding: 0 !important;
  }
}
</style>