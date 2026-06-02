<script setup>
defineProps({
    show: Boolean,
    showReceiptClosing: Boolean,
    pecahan: Object,
    totalUangFisik: Number,
    lastClosingData: Object,
    currentSession: Object,
    currentUser: Object,
    storeLogo: String // 🚀 TAMBAHIN INI BIAR LOGO BISA MASUK
});

const emit = defineEmits(['close', 'process-closing', 'print-closing', 'finish-closing']);
</script>

<template>
    <div v-if="show" class="fixed inset-0 bg-slate-900/90 flex items-center justify-center z-[160] p-4 backdrop-blur-sm no-print">
        <div class="bg-white rounded-[24px] md:rounded-[32px] w-full max-w-2xl overflow-hidden shadow-2xl border-[6px] md:border-[8px] border-slate-800 flex flex-col max-h-[95vh] md:max-h-[90vh]">
            <div class="bg-slate-800 p-4 md:p-6 text-center shrink-0 relative">
                <button @click="emit('close')" class="absolute top-4 right-4 text-slate-400 hover:text-white">✕</button>
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
                    <span class="text-2xl font-black text-indigo-700">Rp {{ totalUangFisik.toLocaleString('id-ID') }}</span>
                </div>
                <button @click="emit('process-closing')" class="w-full bg-indigo-600 text-white py-3.5 rounded-xl font-black text-xs tracking-widest flex items-center justify-center gap-2 uppercase shadow-lg shadow-indigo-200">
                    Proses Tutup Shift
                </button>
            </div>
        </div>
    </div>

    <div v-if="showReceiptClosing" class="fixed inset-0 bg-slate-900/90 flex items-center justify-center z-[200] p-4 backdrop-blur-sm">
        <div class="bg-white p-6 rounded-[32px] shadow-2xl w-full max-w-sm overflow-hidden border-t-8 border-slate-800">
            <div id="print-closing" class="text-left font-mono text-[10px] leading-relaxed uppercase text-black bg-white p-2 mx-auto" :style="{ width: currentSession?.store?.printer_width + 'mm' || '58mm' }">
                
                <div class="text-center mb-4">
                    <img v-if="storeLogo" :src="`http://localhost:8080${storeLogo}`" alt="Logo Toko" class="h-10 mx-auto mb-2 grayscale" />
                    <h2 class="font-black text-xs uppercase">{{ currentSession?.store?.nama_toko || 'ARZU STORE' }}</h2>
                    <div class="w-full h-[1px] bg-black my-1"></div>
                    <p class="font-bold text-[8px] tracking-widest">CLOSING REPORT</p>
                    <p class="font-bold text-[8px]">POS #{{ currentSession?.station_number }}</p>
                </div>
        
                <div class="mb-3">
                    <div class="flex justify-between font-bold">
                        <span>START</span>
                        <span>{{ lastClosingData?.start_time }}</span></div>
                    <div class="flex justify-between font-bold">
                        <span>END</span>
                        <span>{{ lastClosingData?.end_time }}</span></div>
                    <div class="flex justify-between font-bold">
                        <span>CASHIER</span>
                        <span>{{ currentUser?.name?.split(' ')[0] }}</span>
                    </div>
                </div>
        
                <div class="border-b border-dashed border-black mb-3"></div>
        
                <div class="mb-3">
                    <div class="flex justify-between font-bold text-[9px] mt-1">
                        <span>- NET SALES</span>
                        <span class="font-black">{{ lastClosingData?.net_sales?.toLocaleString('id-ID') || 0 }}</span>
                    </div>
                    <div class="flex justify-between font-bold text-[9px]">
                        <span>- TAX</span>
                        <span>{{ lastClosingData?.total_tax?.toLocaleString('id-ID') || 0 }}</span>
                    </div>
                    <div class="flex justify-between font-bold text-[9px]">
                        <span>TOTAL TRANSACTION</span>
                        <span>{{ ((lastClosingData?.net_sales || 0) + (lastClosingData?.total_tax || 0)).toLocaleString('id-ID') }}</span>
                    </div>
                </div>

                <div class="border-b border-black mb-3"></div>
                <div class="flex justify-between font-bold text-[9px]">
                    <span>DRAWER SUMMARY :</span>    
                </div>
        
                <div class="mb-3 space-y-0.5">
                    <div class="flex justify-between font-bold">
                        <span>CASH FLOAT</span>
                        <span>{{ currentSession?.modal_awal?.toLocaleString('id-ID') || 0 }}</span>
                    </div>
                    <div class="flex justify-between font-bold">
                        <span>CASH SALES</span>
                        <span>{{ lastClosingData?.sales_cash?.toLocaleString('id-ID') || 0 }}</span>
                    </div>
                    <div class="flex justify-between font-bold">
                        <span>NON-CASH</span>
                        <span>{{ lastClosingData?.sales_non_tunai?.toLocaleString('id-ID') || 0 }}</span>
                    </div>
                </div>
        
                <div class="border-b border-black mb-2"></div>
        
                <div class="mb-4">
                    <div class="flex justify-between font-black text-[10px] mb-1">
                        <span>EXPECTED CASH</span>
                        <span>{{ lastClosingData?.total_expected?.toLocaleString('id-ID') || 0 }}</span>
                    </div>
                    <div class="flex justify-between font-black text-[10px]">
                        <span>ACTUAL CASH</span>
                        <span>{{ lastClosingData?.total_actual?.toLocaleString('id-ID') || 0 }}</span>
                    </div>
                    <div class="flex justify-between font-black text-[11px] mt-2 border-t border-black pt-1">
                        <span>VARIANCE</span>
                        <span>{{ lastClosingData?.selisih?.toLocaleString('id-ID') || 0 }}</span>
                    </div>
                </div>
        
                <div class="text-center font-bold text-[8px] mt-6">
                    <p>=== END OF SHIFT ===</p>
                    <div class="mt-8 border-t border-black w-3/4 mx-auto"></div>
                    <p class="mt-1">{{ currentUser?.name }}</p>
                </div>
            </div>
            
            <div class="mt-6 flex gap-2 no-print">
                <button @click="emit('print-closing')" class="flex-1 bg-slate-900 text-white py-3 rounded-xl font-black text-[10px] uppercase tracking-widest shadow-lg">Print Struk</button>
                <button @click="emit('finish-closing')" class="flex-1 bg-slate-200 text-slate-700 py-3 rounded-xl font-black text-[10px] uppercase tracking-widest">Selesai</button>
            </div>
        </div>
    </div>
</template>