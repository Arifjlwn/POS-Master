<script setup>
const props = defineProps({
    show: Boolean,
    invoiceData: Object,
    storeData: Object,       
    cashierName: String,     
    stationNumber: String    
});

const emit = defineEmits(['close']);

const formatRupiah = (angka) => {
    if (typeof angka !== 'number' || isNaN(angka)) return '0';
    return new Intl.NumberFormat('id-ID').format(angka);
};

const triggerPrint = () => {
    window.print();
};
</script>

<template>
    <div v-if="show && invoiceData" class="fixed inset-0 bg-slate-950/90 flex items-center justify-center z-[200] p-4 backdrop-blur-sm print:static print:bg-white print:p-0 print:block">
        
        <div class="bg-white p-6 md:p-8 rounded-[32px] md:rounded-[40px] shadow-2xl w-full max-w-sm flex flex-col max-h-[90vh] border-[6px] md:border-[8px] border-slate-800 print:border-none print:shadow-none print:max-h-none print:max-w-none print:rounded-none print:m-0 print:p-0">
            
            <div class="overflow-y-auto custom-scrollbar bg-white p-2 mx-auto print:overflow-visible print:mx-0 print:px-3 print:py-2" id="print-area" style="width: 58mm;">
                
                <div class="text-center mb-4 font-mono leading-none">
                    <h2 class="font-black text-sm uppercase tracking-tighter mb-1 italic">
                        {{ storeData?.NamaToko || storeData?.nama_toko || 'NEXA POS STORE' }}
                    </h2>
                    <p class="text-[8px] font-bold uppercase tracking-widest opacity-80">
                        {{ storeData?.Alamat || storeData?.alamat || 'JAKARTA, INDONESIA' }}
                    </p>
                    <p v-if="storeData?.Telepon || storeData?.telepon" class="text-[8px] font-bold uppercase tracking-widest opacity-80 mt-1">
                        WA: {{ storeData?.Telepon || storeData?.telepon }}
                    </p>
                </div>
                
                <div class="border-y border-black py-1.5 text-center font-black mb-3 font-mono text-[9px] tracking-[0.2em] uppercase bg-slate-100">
                    {{ invoiceData.created_at ? 'Invoice Reprint' : 'Struk Belanja' }}
                </div>
                
                <div class="mb-3 text-[8px] font-bold font-mono uppercase space-y-0.5">
                    <div class="flex justify-between">
                        <span>WAKTU:</span>
                        <span>{{ invoiceData.date || (invoiceData.created_at ? new Date(invoiceData.created_at).toLocaleString('id-ID') : '-') }}</span>
                    </div>
                    <div class="flex justify-between">
                        <span>KASIR:</span>
                        <span>{{ cashierName || 'KASIR' }} / POS-{{ stationNumber || '01' }}</span>
                    </div>
                </div>
                
                <div class="border-b border-black border-dashed mb-2"></div>
                
                <div v-for="item in (invoiceData.cart || invoiceData.details)" :key="item.id" class="mb-2 font-bold font-mono text-[9px] leading-tight uppercase">
                    <div class="truncate w-full pr-2">{{ item.name || item.product?.nama_produk || 'Item Belanja' }}</div>
                    <div class="flex justify-between pl-2 text-[8px] mt-0.5">
                        <span>{{ item.qty || item.kuantitas }} x {{ formatRupiah(item.price || item.harga_satuan) }}</span>
                        <span class="font-black text-[9px]">{{ formatRupiah(item.sub_total || (item.price * item.qty)) }}</span>
                    </div>
                </div>
                
                <div class="border-t border-black border-dashed mt-2 pt-2"></div>
                
                <div class="flex justify-between font-black text-[11px] mb-2 font-mono uppercase italic">
                    <span>TOTAL:</span>
                    <span>Rp{{ formatRupiah(invoiceData.total !== undefined ? invoiceData.total : invoiceData.total_harga) }}</span>
                </div>
                
                <div class="border-b border-black border-dashed mb-2"></div>
                
                <div class="flex justify-between mb-1 font-bold font-mono text-[8px] uppercase">
                    <span>BAYAR ({{ invoiceData.method || invoiceData.metode_bayar || 'CASH' }}):</span>
                    <span>Rp{{ formatRupiah(invoiceData.pay !== undefined ? invoiceData.pay : invoiceData.nominal_bayar) }}</span>
                </div>
                
                <div class="flex justify-between font-black font-mono text-[9px] uppercase italic text-black">
                    <span>KEMBALI:</span>
                    <span>Rp{{ formatRupiah(invoiceData.total !== undefined ? invoiceData['return'] : invoiceData.kembalian) }}</span>
                </div>
                
                <div class="mt-5 text-[7px] font-bold text-center border-t border-black border-dashed pt-2 font-mono uppercase space-y-1">
                    <p class="font-black">INV: {{ invoiceData.invoice || invoiceData.no_invoice }}</p>
                </div>
                <div class="text-center mt-4 font-black font-mono text-[8px] border-2 border-black p-1.5 uppercase">
                    Terima Kasih!<br>Barang tidak dapat ditukar.
                </div>
            </div>
            
            <div class="mt-4 md:mt-6 flex flex-col gap-2 md:gap-3 shrink-0 print:hidden">
                <button @click="triggerPrint" class="w-full bg-indigo-600 text-white py-3 md:py-4 rounded-xl md:rounded-2xl font-black text-[9px] md:text-[10px] uppercase tracking-[0.2em] shadow-xl shadow-indigo-200 flex items-center justify-center gap-2 active:scale-95 transition-all">
                    Cetak Struk
                </button>
                <button @click="emit('close')" class="w-full bg-slate-100 text-slate-500 py-3 md:py-3.5 rounded-xl md:rounded-2xl font-black text-[9px] md:text-[10px] uppercase tracking-widest hover:bg-slate-200 transition-all">Tutup</button>
            </div>
        </div>
    </div>
</template>