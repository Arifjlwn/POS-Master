<script setup>
import { useStockOpname } from '../../composables/useStockOpname.js';
import Sidebar from '../../components/Sidebar.vue';
import StockOpnameHeader from '../../components/stockopname/StockOpnameHeader.vue';
import StockOpnameSearch from '../../components/stockopname/StockOpnameSearch.vue';
import StockOpnameTable from '../../components/stockopname/StockOpnameTable.vue';
import StockOpnameSummary from '../../components/stockopname/StockOpnameSummary.vue';

const {
    activeTab, soStep, notes, searchQuery, products, cartSO, cartKlaim, isSubmitting, showScanner, isOwner,
    isKlaimEligible, daysLeftKlaim,isSOLockedThisMonth,startScanner, stopScanner, searchProduct, addToCart, removeItem, 
    proceedToReview, backToCounting, submitSOFinal, submitKlaimTemuan, hitungTotalFisik
} = useStockOpname();
</script>

<template>
    <Sidebar>
        <div class="p-4 md:p-8 lg:p-10 max-w-7xl mx-auto font-sans bg-[#f8fafc] min-h-screen">
            
            <StockOpnameHeader 
                v-model:activeTab="activeTab" 
                :isOwner="isOwner" 
                :isKlaimEligible="isKlaimEligible"
                :daysLeftKlaim="daysLeftKlaim"
            />

            <div class="grid grid-cols-1 xl:grid-cols-3 gap-8">
                
                <div class="xl:col-span-2 space-y-6">
                    
                    <div v-if="activeTab === 'SO' && isSOLockedThisMonth" class="bg-indigo-50 border-2 border-indigo-200 p-6 rounded-2xl flex flex-col sm:flex-row items-center gap-4 text-center sm:text-left">
                        <div class="w-12 h-12 bg-white rounded-full flex items-center justify-center shadow-sm shrink-0">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-indigo-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/></svg>
                        </div>
                        <div>
                            <h3 class="font-black text-indigo-700 uppercase tracking-widest text-sm">Sistem Terkunci</h3>
                            <p class="text-xs font-bold text-indigo-500 mt-1">Audit (Stock Opname) sudah dilakukan bulan ini. Silakan kembali bulan depan, atau gunakan menu KLAIM BARANG untuk koreksi stok positif.</p>
                        </div>
                    </div>

                    <div :class="{'opacity-40 pointer-events-none select-none grayscale-[30%]': activeTab === 'SO' && isSOLockedThisMonth}" class="space-y-6 transition-all duration-300">
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
                            :activeTab="activeTab"
                            :hitungTotalFisik="hitungTotalFisik" @remove="removeItem"
                        />
                    </div>
                </div>

                <div class="xl:col-span-1">
                    <div :class="{'opacity-40 pointer-events-none select-none grayscale-[30%]': activeTab === 'SO' && isSOLockedThisMonth}" class="transition-all duration-300">
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
            </div>

            <div v-if="showScanner" class="fixed inset-0 bg-slate-900/90 flex items-center justify-center z-[150] p-4 backdrop-blur-sm">
                </div>

        </div>
    </Sidebar>
</template>

<style scoped>
.animate-fade-in { animation: fadeIn 0.3s ease-in-out; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(5px); } to { opacity: 1; transform: translateY(0); } }
</style>