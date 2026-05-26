<script setup>
import { usePenerimaanBarang } from '../../composables/usePenerimaanBarang.js';
import Sidebar from '../../components/Sidebar.vue';
import InboundHeader from '../../components/inbound/InboundHeader.vue';
import SupplierForm from '../../components/inbound/SupplierForm.vue';
import ScannerModal from '../../components/inbound/ScannerModal.vue';
import InboundSearch from '../../components/inbound/InboundSearch.vue';
import InboundTable from '../../components/inbound/InboundTable.vue';

const {
    currentUser, isOwner, supplierName, noFaktur, searchQuery, products, cartLPB, isSubmitting, showScanner,
    startScanner, stopScanner, searchProduct, addToCart, hitungTotalStok, hitungModalPerPcs, removeItem, submitLPB
} = usePenerimaanBarang();
</script>

<template>
    <Sidebar>
        <div class="p-4 md:p-8 lg:p-10 max-w-7xl mx-auto font-sans bg-[#f8fafc] min-h-screen">
            
            <InboundHeader :currentUser="currentUser" />

            <div class="grid grid-cols-1 xl:grid-cols-4 gap-8">
                
                <div class="xl:col-span-1 space-y-6">
                    <SupplierForm 
                        v-model:supplierName="supplierName" 
                        v-model:noFaktur="noFaktur" 
                    />
                </div>

                <div class="xl:col-span-3 space-y-6">
                    <div class="bg-white p-6 rounded-[32px] shadow-sm border border-slate-100 relative">
                        <InboundSearch 
                            v-model:searchQuery="searchQuery"
                            @search="searchProduct(false)"
                            @scan="startScanner"
                        />
                        
                        <div v-if="products.length && searchQuery.length > 0" class="absolute left-6 right-6 mt-3 bg-white border border-slate-100 rounded-[28px] shadow-2xl z-[100] overflow-hidden">
                            <div v-for="p in products" :key="p.id" @click="addToCart(p)" class="p-5 hover:bg-blue-50 cursor-pointer border-b last:border-0 flex justify-between items-center group transition-all">
                                <div class="flex items-center gap-4">
                                    <div class="w-10 h-10 bg-slate-100 rounded-xl flex items-center justify-center text-xl shadow-inner group-hover:bg-blue-100 transition-colors">📦</div>
                                    <div class="font-black text-slate-800 uppercase text-sm group-hover:text-blue-600 transition-colors">{{ p.nama_produk }}</div>
                                </div>
                                <span class="text-[9px] bg-slate-100 text-slate-500 px-3 py-1.5 rounded-lg font-black border border-slate-200">STOK: {{ p.stok }}</span>
                            </div>
                        </div>
                    </div>

                    <div class="bg-white rounded-[40px] shadow-xl shadow-slate-200/50 border border-slate-100 overflow-hidden">
                        <div class="overflow-x-auto custom-scrollbar">
                            <InboundTable 
            :cartLPB="cartLPB"
            :isOwner="isOwner"
            :hitungTotalStok="hitungTotalStok"
            :hitungModalPerPcs="hitungModalPerPcs"
            @remove="removeItem"
        />
                        </div>
                    </div>

                    <button @click="submitLPB" :disabled="isSubmitting" class="w-full bg-blue-600 hover:bg-slate-900 text-white py-6 rounded-[30px] font-black text-xs uppercase tracking-[0.3em] shadow-2xl transition-all active:scale-[0.98] disabled:opacity-50 flex items-center justify-center gap-4">
                        <template v-if="isSubmitting"><div class="w-5 h-5 border-4 border-white/30 border-t-white rounded-full animate-spin"></div>Sinkronisasi Stok...</template>
                        <template v-else>Posting Penerimaan Barang</template>
                    </button>
                </div>
            </div>

            <ScannerModal :show="showScanner" @close="stopScanner" />

        </div>
    </Sidebar>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { height: 8px; width: 8px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #e2e8f0; border-radius: 20px; border: 2px solid transparent; background-clip: content-box; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #cbd5e1; background-clip: content-box; }
input[type=number]::-webkit-inner-spin-button, input[type=number]::-webkit-outer-spin-button { -webkit-appearance: none; margin: 0; }
</style>