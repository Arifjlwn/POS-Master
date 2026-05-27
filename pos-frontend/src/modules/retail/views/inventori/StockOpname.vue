<script setup>
import { useStockOpname } from '../../composables/useStockOpname.js';
import Sidebar from '../../components/Sidebar.vue';
import StockOpnameHeader from '../../components/stockopname/StockOpnameHeader.vue';
import StockOpnameSearch from '../../components/stockopname/StockOpnameSearch.vue';
import StockOpnameTable from '../../components/stockopname/StockOpnameTable.vue';
import StockOpnameSummary from '../../components/stockopname/StockOpnameSummary.vue';

const {
    activeTab, soStep, notes, searchQuery, products, cartSO, cartKlaim, isSubmitting, showScanner, isOwner,
    startScanner, stopScanner, searchProduct, addToCart, removeItem, 
    proceedToReview, submitSOFinal, submitKlaimTemuan
} = useStockOpname();
</script>

<template>
    <Sidebar>
        <div class="p-4 md:p-8 lg:p-10 max-w-7xl mx-auto font-sans bg-[#f8fafc] min-h-screen">
            
            <StockOpnameHeader 
                v-model:activeTab="activeTab" 
                :isOwner="isOwner" 
                :isKlaimEligible = "isKlaimEligible"
                :daysLeftKlaim = "daysLeftKlaim"
                
            />

            <div class="grid grid-cols-1 xl:grid-cols-3 gap-8">
                
                <div class="xl:col-span-2 space-y-6">
                    <div v-if="activeTab === 'KLAIM'" class="bg-amber-50 border-2 border-amber-200 p-6 rounded-[32px] animate-fade-in">
                        <h3 class="font-black text-amber-800 text-sm uppercase tracking-widest mb-2">Klaim Barang Nyempil</h3>
                        <p class="text-xs text-amber-700 font-medium">Gunakan form ini HANYA jika menemukan barang <b class="font-black">PLUS</b> setelah audit selesai. Membutuhkan Approval Owner.</p>
                    </div>

                    <StockOpnameSearch 
                        v-if="soStep === 'COUNTING' || activeTab === 'KLAIM'"
                        v-model:searchQuery="searchQuery" 
                        :products="products" 
                        @search="searchProduct(false)" 
                        @scan="startScanner" 
                        @add="addToCart" 
                    />

                    <StockOpnameTable 
                        :cartSO="activeTab === 'SO' ? cartSO : cartKlaim"
                        :soStep="activeTab === 'SO' ? soStep : 'REVIEW'" 
                        :isOwner="isOwner"
                        @remove="removeItem"
                    />
                </div>

                <div class="xl:col-span-1">
                    <StockOpnameSummary 
                        v-model:notes="notes"
                        :cartLength="activeTab === 'SO' ? cartSO.length : cartKlaim.length"
                        :isSubmitting="isSubmitting"
                        :isOwner="isOwner"
                        :soStep="soStep"
                        :activeTab="activeTab"
                        @proceed="proceedToReview"
                        @submit="activeTab === 'SO' ? submitSOFinal() : submitKlaimTemuan()"
                    />
                </div>
            </div>

            <div v-if="showScanner" class="fixed inset-0 bg-slate-900/90 flex items-center justify-center z-[150] p-4 backdrop-blur-sm">
                <div class="bg-white rounded-[32px] shadow-2xl w-full max-w-md overflow-hidden flex flex-col">
                    <div class="p-6 border-b border-slate-50 flex justify-between items-center bg-slate-50/50">
                        <h2 class="text-lg font-black text-slate-800 uppercase tracking-widest">Scan Barcode</h2>
                        <button @click="stopScanner" class="p-2 rounded-xl bg-slate-100 text-slate-400 hover:text-red-500 transition-all"><svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg></button>
                    </div>
                    <div class="p-6 bg-black relative"><div id="reader-so" class="w-full rounded-2xl overflow-hidden border-2 border-slate-700"></div></div>
                </div>
            </div>

        </div>
    </Sidebar>
</template>

<style scoped>
.animate-fade-in { animation: fadeIn 0.3s ease-in-out; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(5px); } to { opacity: 1; transform: translateY(0); } }
</style>