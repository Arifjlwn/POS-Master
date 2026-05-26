<script setup>
defineProps({
    isOpen: Boolean,
    document: Object,
    storeName: String
});
defineEmits(['close', 'print']);
</script>

<template>
    <div v-if="isOpen" class="fixed inset-0 z-[100] flex items-center justify-center p-4 print:p-0 print:relative print:block print:z-auto">
        <div @click="$emit('close')" class="absolute inset-0 bg-slate-900/80 backdrop-blur-sm print:hidden transition-all"></div>
        
        <div class="relative bg-white w-full max-w-4xl max-h-[90vh] overflow-y-auto rounded-[32px] shadow-2xl flex flex-col print:shadow-none print:rounded-none print:max-h-none print:h-auto print:overflow-visible custom-scrollbar">
            
            <div class="sticky top-0 bg-white/90 backdrop-blur-md px-6 py-4 border-b border-slate-100 flex justify-between items-center z-10 print:hidden shrink-0">
                <h3 class="font-black text-slate-800 text-lg uppercase tracking-tight">Detail Dokumen</h3>
                <div class="flex items-center gap-3">
                    <button @click="$emit('print')" class="bg-indigo-600 hover:bg-slate-900 text-white px-4 py-2.5 rounded-xl font-black text-[10px] uppercase tracking-widest transition-all flex items-center gap-2 shadow-md">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="6 9 6 2 18 2 18 9"/><path d="M6 18H4a2 2 0 0 1-2-2v-5a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v5a2 2 0 0 1-2 2h-2"/><rect width="12" height="8" x="6" y="14"/></svg>
                        Cetak PDF
                    </button>
                    <button @click="$emit('close')" class="w-10 h-10 rounded-xl bg-slate-100 text-slate-400 hover:text-red-500 hover:bg-red-50 transition-all flex items-center justify-center">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                    </button>
                </div>
            </div>

            <div id="printable-area" class="p-8 md:p-12 print:p-0 print:text-black">
                
                <div class="text-center border-b-[3px] border-slate-800 pb-6 mb-6">
                    <h1 class="text-2xl font-black uppercase tracking-widest text-slate-900">{{ storeName }}</h1>
                    <p class="text-sm font-medium text-slate-600 mt-1">BERITA ACARA PEMUSNAHAN / RETUR BARANG</p>
                </div>

                <div class="flex justify-between items-end mb-8">
                    <div>
                        <table class="text-xs font-bold text-slate-700">
                            <tr><td class="pb-2 pr-4 text-slate-500 uppercase tracking-widest text-[9px]">No. Dokumen</td><td class="pb-2">: {{ document?.return_no }}</td></tr>
                            <tr><td class="pb-2 pr-4 text-slate-500 uppercase tracking-widest text-[9px]">Tanggal</td><td class="pb-2">: {{ document ? new Date(document.created_at).toLocaleDateString('id-ID', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' }) : '' }}</td></tr>
                            <tr><td class="pr-4 text-slate-500 uppercase tracking-widest text-[9px]">Operator</td><td>: {{ document?.user.name }}</td></tr>
                        </table>
                    </div>
                    <div class="text-right">
                        <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1">Total Kuantitas</p>
                        <p class="text-3xl font-black text-slate-900 tracking-tighter">{{ document?.total_qty }}</p>
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
                        <tr v-for="(item, index) in document?.items" :key="item.id">
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
                        <p class="text-sm font-black text-slate-800 uppercase underline">{{ document?.user.name }}</p>
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
</template>