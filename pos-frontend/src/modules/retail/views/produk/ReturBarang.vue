<script setup>
import { useReturBarang } from '../../composables/useReturBarang.js';
import Sidebar from '../../components/Sidebar.vue';

// Import Modular Components
import ReturHeader from '../../components/return/ReturHeader.vue';
import ReturForm from '../../components/return/ReturForm.vue';
import ReturCart from '../../components/return/ReturCart.vue';

const {
    cart, isSubmitting, lastSubmittedReturn, searchProductQuery, 
    isDropdownOpen, selectedProduct, isScannerOpen, cameras, 
    selectedCamera, form, alasanOptions, filteredProducts, 
    startScanner, stopScanner, switchCamera, selectProduct, 
    clearSelectedProduct, addToCart, removeFromCart, 
    submitBatchReturn, getBadgeClass
} = useReturBarang();
</script>

<template>
    <Sidebar>
        <div class="p-4 md:p-8 lg:p-10 max-w-7xl mx-auto font-sans bg-[#f8fafc] min-h-screen print:bg-white print:p-0">
            
            <ReturHeader />

            <div class="grid grid-cols-1 xl:grid-cols-2 gap-8 print:hidden">
                <div class="xl:col-span-1">
                    <ReturForm 
                        v-model:searchProductQuery="searchProductQuery"
                        :isDropdownOpen="isDropdownOpen"
                        :filteredProducts="filteredProducts"
                        :selectedProduct="selectedProduct"
                        :form="form"
                        :alasanOptions="alasanOptions"
                        @start-scanner="startScanner"
                        @select-product="selectProduct"
                        @clear-product="clearSelectedProduct"
                        @add-to-cart="addToCart"
                    />
                </div>

                <div class="xl:col-span-1 flex flex-col h-full">
                    <ReturCart 
                        :cart="cart"
                        :isSubmitting="isSubmitting"
                        :getBadgeClass="getBadgeClass"
                        @remove="removeFromCart"
                        @submit="submitBatchReturn"
                    />
                </div>
            </div>

            <div v-if="lastSubmittedReturn" id="printable-area" class="hidden print:block print:p-0 print:text-black print:absolute print:inset-0 print:w-full">
                <div class="text-center border-b-[3px] border-slate-800 pb-6 mb-6">
                    <h1 class="text-2xl font-black uppercase tracking-widest text-slate-900">{{ lastSubmittedReturn.storeName }}</h1>
                    <p class="text-sm font-medium text-slate-600 mt-1">BERITA ACARA PEMUSNAHAN / RETUR BARANG</p>
                </div>
                <div class="flex justify-between items-end mb-8">
                    <div>
                        <table class="text-xs font-bold text-slate-700">
                            <tbody><tr><td class="pb-2 pr-4 text-slate-500 uppercase tracking-widest text-[9px]">No. Dokumen</td><td class="pb-2">: {{ lastSubmittedReturn.return_no }}</td></tr>
                            <tr><td class="pb-2 pr-4 text-slate-500 uppercase tracking-widest text-[9px]">Tanggal</td><td class="pb-2">: {{ new Date(lastSubmittedReturn.created_at).toLocaleDateString('id-ID', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' }) }}</td></tr>
                            <tr><td class="pr-4 text-slate-500 uppercase tracking-widest text-[9px]">Operator</td><td>: {{ lastSubmittedReturn.user.name }}</td></tr>
                            </tbody>
                        </table>
                    </div>
                    <div class="text-right">
                        <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1">Total Kuantitas</p>
                        <p class="text-3xl font-black text-slate-900 tracking-tighter">{{ lastSubmittedReturn.total_qty }}</p>
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
                        <tr v-for="(item, index) in lastSubmittedReturn.items" :key="index">
                            <td class="py-3 px-2 text-xs font-bold text-slate-600 text-center">{{ index + 1 }}</td>
                            <td class="py-3 px-2 text-[10px] font-bold text-slate-500 uppercase tracking-wider">{{ item.sku || '-' }}</td>
                            <td class="py-3 px-2 text-xs font-black text-slate-800 uppercase">{{ item.nama_produk }}</td>
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
                        <p class="text-sm font-black text-slate-800 uppercase underline">{{ lastSubmittedReturn.user.name }}</p>
                        <p class="text-[9px] font-bold text-slate-400 mt-1">Staff / Kasir</p>
                    </div>
                    <div class="text-center">
                        <p class="text-[10px] font-bold text-slate-500 uppercase tracking-widest mb-16">Mengetahui,</p>
                        <p class="text-sm font-black text-slate-800 uppercase underline">...................................</p>
                        <p class="text-[9px] font-bold text-slate-400 mt-1">Manager / Owner</p>
                    </div>
                </div>
            </div>

            <div v-if="isScannerOpen" class="fixed inset-0 bg-slate-900/90 flex items-center justify-center z-[100] p-4 backdrop-blur-sm print:hidden">
                <div class="bg-white rounded-[32px] overflow-hidden w-full max-w-sm shadow-2xl border-[8px] border-slate-800 flex flex-col">
                    <div class="p-5 border-b border-slate-100 flex justify-between items-center bg-white shrink-0">
                        <div>
                            <h3 class="font-black text-slate-800 uppercase tracking-tighter text-lg italic">Scan Barcode</h3>
                            <p class="text-[9px] text-rose-500 font-black uppercase tracking-widest mt-0.5">Retur / Buang Barang</p>
                        </div>
                        <button @click="stopScanner" class="w-10 h-10 rounded-xl bg-slate-100 text-slate-400 hover:text-red-500 hover:bg-red-50 transition-all flex items-center justify-center">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                        </button>
                    </div>
                    <div v-if="cameras.length > 0" class="p-3 bg-slate-50 border-b border-slate-100 flex gap-2 items-center">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-400 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M14.5 4h-5L7 7H4a2 2 0 0 0-2 2v9a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2h-3l-2.5-3z"/><circle cx="12" cy="13" r="3"/></svg>
                        <select v-model="selectedCamera" @change="switchCamera" class="w-full bg-white border border-slate-200 text-[10px] font-black uppercase tracking-widest text-slate-600 rounded-lg p-2 outline-none cursor-pointer">
                            <option v-for="cam in cameras" :key="cam.id" :value="cam.id">{{ cam.label || `Kamera ${cam.id.substring(0, 5)}...` }}</option>
                        </select>
                    </div>
                    <div class="relative bg-black w-full aspect-square flex items-center justify-center">
                        <div id="reader" class="w-full h-full object-cover"></div>
                        <div class="absolute inset-0 border-[12px] border-black/40 pointer-events-none z-10"></div>
                        <div class="absolute inset-x-8 inset-y-16 border-2 border-white/30 rounded-[20px] pointer-events-none z-20">
                            <div class="absolute top-0 left-0 w-8 h-8 border-t-4 border-l-4 border-rose-500 rounded-tl-[18px]"></div>
                            <div class="absolute top-0 right-0 w-8 h-8 border-t-4 border-r-4 border-rose-500 rounded-tr-[18px]"></div>
                            <div class="absolute bottom-0 left-0 w-8 h-8 border-b-4 border-l-4 border-rose-500 rounded-bl-[18px]"></div>
                            <div class="absolute bottom-0 right-0 w-8 h-8 border-b-4 border-r-4 border-rose-500 rounded-br-[18px]"></div>
                            <div class="w-full h-0.5 bg-rose-500 absolute top-0 animate-[scan_2s_infinite] shadow-[0_0_8px_#f43f5e]"></div>
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

input[type=number]::-webkit-inner-spin-button, 
input[type=number]::-webkit-outer-spin-button { -webkit-appearance: none; margin: 0; }
input[type=number] { -moz-appearance: textfield; }
select { -webkit-appearance: none; -moz-appearance: none; appearance: none; }

@keyframes scan { 0% { top: 0%; opacity: 0; } 10% { opacity: 1; } 90% { opacity: 1; } 100% { top: 100%; opacity: 0; } }

:deep(#reader) { border: none !important; }
:deep(#reader video) { object-fit: cover !important; }

@media print {
    /* 1. Hapus semua background biar hemat tinta dan bersih */
    * {
        background: transparent !important;
        color: black !important;
        box-shadow: none !important;
        text-shadow: none !important;
    }

    /* 2. Sembunyikan elemen aplikasi yang gak mau dicetak */
    .no-print, header, nav, .sidebar, aside, .bg-gradient-to-br, button {
        display: none !important;
    }

    /* 3. Pastikan area cetak lu tampil dengan ukuran normal */
    #printable-area {
        display: block !important;
        position: static !important; /* Jangan absolute biar dia ngikutin alur kertas */
        width: 100% !important;
        visibility: visible !important;
    }

    /* 4. Pastikan isi dari area cetak lu kelihatan */
    #printable-area * {
        visibility: visible !important;
    }

    /* 5. Paksa page break biar rapi */
    .print-break {
        page-break-after: always;
    }
}
</style>