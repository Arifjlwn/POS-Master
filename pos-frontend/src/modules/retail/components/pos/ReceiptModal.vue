<script setup>
defineProps({
    show: Boolean, transaction: Object, user: Object, session: Object
});
const emit = defineEmits(['close']);
const print = () => window.print();
</script>

<template>
    <div v-if="show" class="fixed inset-0 bg-slate-950/90 flex items-center justify-center z-[200] p-4 backdrop-blur-sm no-print">
        <div class="bg-white p-6 md:p-8 rounded-[32px] shadow-2xl w-full max-w-sm flex flex-col max-h-[90vh] border-[8px] border-slate-800">
            <div class="overflow-y-auto bg-white p-2 mx-auto shadow-inner" id="print-area" style="width: 58mm; font-family: monospace;">
                <div class="text-center mb-4 leading-none text-xs">
                    <h2 class="font-black uppercase italic">{{ session?.store?.nama_toko || 'ARZU STORE JKT' }}</h2>
                    <p class="text-[8px] opacity-70 mt-1">{{ session?.store?.alamat_toko || 'JAKARTA 2 BRANCH, INDONESIA' }}</p>
                </div>
                <div class="border-y border-black py-1 text-center font-black text-[9px] uppercase tracking-wider mb-2 bg-slate-50">Struk Belanja</div>
                <div class="text-[8px] space-y-0.5 mb-2">
                    <div class="flex justify-between"><span>WAKTU:</span><span>{{ transaction?.date }}</span></div>
                    <div class="flex justify-between"><span>KASIR:</span><span>{{ user?.name?.split(' ')[0] }}</span></div>
                </div>
                <div class="border-b border-black border-dashed mb-2"></div>
                <div v-for="item in transaction?.cart" :key="item.id" class="text-[9px] mb-2 leading-tight">
                    <div class="truncate uppercase font-bold">{{ item.name }}</div>
                    <div class="flex justify-between text-[8px] text-slate-600 pl-1">
                        <span>{{ item.qty }} x {{ item.price.toLocaleString('id-ID') }}</span>
                        <span class="font-black text-black">Rp{{ (item.price * item.qty).toLocaleString('id-ID') }}</span>
                    </div>
                </div>
                <div class="border-t border-black border-dashed pt-1.5 flex justify-between font-black text-xs">
                    <span>TOTAL:</span><span>Rp{{ transaction?.total?.toLocaleString('id-ID') }}</span>
                </div>
                <div class="text-center mt-4 text-[8px] border border-black p-1 uppercase font-black">Terima Kasih!</div>
            </div>
            <div class="mt-4 flex flex-col gap-2 no-print">
                <button @click="print" class="w-full bg-indigo-600 text-white py-3.5 rounded-xl font-black text-xs uppercase tracking-widest shadow-md">Cetak Struk 58mm</button>
                <button @click="emit('close')" class="w-full bg-slate-100 text-slate-500 py-3 rounded-xl font-bold text-xs uppercase">Tutup Jendela</button>
            </div>
        </div>
    </div>
</template>