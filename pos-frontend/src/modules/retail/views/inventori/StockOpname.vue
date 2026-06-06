<script setup>
import { useStockOpname } from '../../composables/useStockOpname.js';
import Sidebar from '../../components/Sidebar.vue';
import StockOpnameHeader from '../../components/stockopname/StockOpnameHeader.vue';
import StockOpnameSearch from '../../components/stockopname/StockOpnameSearch.vue';
import StockOpnameTable from '../../components/stockopname/StockOpnameTable.vue';
import StockOpnameSummary from '../../components/stockopname/StockOpnameSummary.vue';

const {
    activeTab, soStep, notes, searchQuery, products, cartSO, cartKlaim, isSubmitting, showScanner, isOwner, // ◄ Pastikan di sini isOwner polos tanpa typo 
    isKlaimEligible, daysLeftKlaim, isSOLockedThisMonth, startScanner, stopScanner, searchProduct, addToCart, removeItem, 
    proceedToReview, backToCounting, submitSOFinal, submitKlaimTemuan, hitungTotalFisik
} = useStockOpname();
</script>

<template>
    <Sidebar>
        <div class="p-4 md:p-8 lg:p-10 max-w-7xl mx-auto font-sans bg-[#f8fafc] min-h-screen">
            
            <!-- PANEL ATAS: NAVIGASI TAB AUDIT VUE -->
            <StockOpnameHeader 
                v-model:activeTab="activeTab" 
                :isOwner="isOwner" 
                :isKlaimEligible="isKlaimEligible"
                :daysLeftKlaim="daysLeftKlaim"
            />

            <div class="grid grid-cols-1 xl:grid-cols-3 gap-8 mt-6">
                
                <!-- KIRI: PANEL INPUT & DAFTAR ITEM BARANG -->
                <div class="xl:col-span-2 space-y-6">
                    
                    <!-- ALERT BANNER: SISTEM KUNCI MANDIRI ANTI DUPLIKASI DATA DATA -->
                    <div v-if="activeTab === 'SO' && isSOLockedThisMonth" class="bg-indigo-50 border-2 border-indigo-200 p-6 rounded-2xl flex flex-col sm:flex-row items-center gap-4 text-center sm:text-left animate-fade-in">
                        <div class="w-12 h-12 bg-white rounded-full flex items-center justify-center shadow-sm shrink-0">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-indigo-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/></svg>
                        </div>
                        <div>
                            <h3 class="font-black text-indigo-700 uppercase tracking-widest text-sm">Sistem Terkunci</h3>
                            <p class="text-xs font-bold text-indigo-500 mt-1">Audit bulanan sudah selesai diproses ! Gerbang ditutup untuk mencegah manipulasi variance laci, silakan gunakan fitur KLAIM BARANG jika ada temuan stock opname positif.</p>
                        </div>
                    </div>

                    <!-- BLOK UTAMA TRANSAKSI PENCATATAN KARTU STOK -->
                    <div :class="{'opacity-40 pointer-events-none select-none grayscale-[30%]': activeTab === 'SO' && isSOLockedThisMonth}" class="space-y-6 transition-all duration-300">
                        
                        <!-- 🛡️ FIX ROW SEARCH: Izinkan pencarian tetap tampil di tab Klaim  biar kasir bisa masukin barang temuan! -->
                        <StockOpnameSearch 
                          v-if="soStep === 'COUNTING'"
                          v-model:searchQuery="searchQuery"
                          :disabled="activeTab === 'SO' && isSOLockedThisMonth"
                          :products="products" 
                          @search="searchProduct(false)" 
                          @scan="startScanner" 
                          @add="addToCart" 
                        />

                        <!-- TABLE RENDER: REKAP DAFTAR FISIK BARANG -->
                        <StockOpnameTable 
                            :cartSO="activeTab === 'SO' ? cartSO : cartKlaim"
                            :soStep="activeTab === 'SO' ? soStep : 'REVIEW'" 
                            :isOwner="isOwner"
                            :activeTab="activeTab"
                            :hitungTotalFisik="hitungTotalFisik" 
                            @remove="removeItem"
                        />
                    </div>
                </div>

                <!-- KANAN: PANEL SUMMARY KETERANGAN & TOMBOL SUBMIT FINAL -->
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
                            
                            @back-step="soStep = 'COUNTING'" 
                        />
                    </div>
                </div>
            </div>

            <!-- ========================================================================= -->
            <!-- 🛡️ FIX EXTRA SECURITY MODAL CAMERA SCANNER (ANTI CAMERA LEAK RESOURCE)       -->
            <!-- ========================================================================= -->
            <div v-if="showScanner" class="fixed inset-0 bg-slate-950/95 flex flex-col items-center justify-center z-[200] p-4 backdrop-blur-md animate-fade-in">
                <div class="bg-white rounded-[32px] p-6 w-full max-w-md text-center shadow-2xl relative border-4 border-indigo-600">
                    
                    <!-- TOMBOL DARURAT CLOSE SCANNER  -->
                    <button @click="stopScanner" class="absolute top-4 right-4 bg-slate-100 hover:bg-rose-500 text-slate-500 hover:text-white p-2 rounded-xl transition-colors active:scale-95 shadow-sm" title="Matikan Kamera">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
                    </button>

                    <div class="w-16 h-16 bg-indigo-50 text-indigo-600 rounded-full flex items-center justify-center mx-auto mb-3 border border-indigo-100">
                        <svg class="w-8 h-8" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" /><path stroke-linecap="round" stroke-linejoin="round" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
                    </div>

                    <h3 class="font-black text-slate-800 text-lg uppercase tracking-tight">Scan Barcode Produk</h3>
                    <p class="text-xs text-slate-400 font-medium mb-4">Posisikan kode batang SKU tepat di tengah kotak kamera</p>

                    <!-- CONTAINER TARGET STREAM KAMERA PERANGKAT  -->
                    <div id="reader" class="w-full aspect-square bg-slate-900 rounded-2xl overflow-hidden border-2 border-dashed border-indigo-300 shadow-inner flex items-center justify-center text-slate-500 text-xs font-mono">
                        <div class="flex flex-col items-center gap-2 animate-pulse">
                            <div class="w-6 h-6 border-2 border-indigo-400 border-t-transparent rounded-full animate-spin"></div>
                            <span>Menghubungkan Kamera...</span>
                        </div>
                    </div>
                </div>
            </div>

        </div>
    </Sidebar>
</template>

<style scoped>
.animate-fade-in { animation: fadeIn 0.2s cubic-bezier(0.16, 1, 0.3, 1); }
@keyframes fadeIn { from { opacity: 0; transform: translateY(8px); } to { opacity: 1; transform: translateY(0); } }
</style>