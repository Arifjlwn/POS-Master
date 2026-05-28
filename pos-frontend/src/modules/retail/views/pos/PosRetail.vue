<script setup>
import { useRouter } from 'vue-router';
import { usePos } from '../../composables/usePos.js';

// 🚀 IMPORT PASUKAN SUB-KOMPONEN BARU KITA BEB!
import PosHeader from '../../components/pos/PosHeader.vue';
import ProductCatalog from '../../components/pos/ProductCatalog.vue';
import CartSidebar from '../../components/pos/CartSidebar.vue';
import ClosingModal from '../../components/pos/ClosingModal.vue';
import ReceiptModal from '../../components/pos/ReceiptModal.vue';

const router = useRouter();

const {
    currentUser, currentSession, currentTime, products, isLoadingProducts, cart, heldOrders,
    showHeldModal, payAmount, paymentMethod, showReceipt, showQrisModal, lastTransaction,
    showReceiptClosing, lastClosingData, isMobileCartOpen, searchQuery, searchInput,
    showScanner, pecahan, filteredProducts, totalBelanja, kembalian, totalUangFisik,
    showClosingModal, isProcessingCheckout,
    getImageUrl, startScanner, stopScanner, handleBarcodeScan, addToCart,toggleUom,
    decreaseQty, increaseQty, validateQty, clearCart, holdTransaction, resumeOrder,
    setPaymentMethod, executeCheckout, formatInputRupiah, processCheckout, handleClosing, logout
} = usePos();

// DOM Logic Print spesifik untuk view layer
const printClosing = () => {
    window.print(); 
};

const finishClosing = () => router.push('/retail/absensi');
const goToDashboard = () => currentUser.value.role === 'owner' ? router.push('/retail/dashboard') : router.push('/retail/absensi');
</script>

<template>
    <div class="h-[100dvh] flex flex-col bg-slate-100 font-sans overflow-hidden print:h-auto print:bg-white print:overflow-visible print:block">
        
        <PosHeader
            class="print:hidden" 
            :currentUser="currentUser" 
            :currentSession="currentSession" 
            :currentTime="currentTime"
            @go-dashboard="goToDashboard"
            @logout="logout"
        />

        <div class="flex-1 flex overflow-hidden p-2 md:p-4 pt-2 md:pt-4 gap-4 relative print:hidden">
            
            <ProductCatalog 
                v-model:searchQuery="searchQuery"
                :filteredProducts="filteredProducts"
                :heldOrders="heldOrders"
                :getImageUrl="getImageUrl"
                @barcode-scan="handleBarcodeScan"
                @start-scanner="startScanner"
                @show-held="showHeldModal = true"
                @add-to-cart="addToCart"
            />

            <CartSidebar 
                v-model:isMobileCartOpen="isMobileCartOpen"
                :cart="cart"
                :heldOrders="heldOrders"
                :paymentMethod="paymentMethod"
                :totalBelanja="totalBelanja"
                :payAmount="payAmount"
                :kembalian="kembalian"
                :isProcessingCheckout="isProcessingCheckout"
                @show-held="showHeldModal = true"
                @hold-order="holdTransaction"
                @clear-cart="clearCart"
                @decrease-qty="decreaseQty"
                @increase-qty="increaseQty"
                @validate-qty="validateQty"
                @set-payment="setPaymentMethod"
                @format-rupiah="formatInputRupiah"
                @checkout="processCheckout"
                @toggle-uom="toggleUom" />
        </div>

        <div v-if="cart.length > 0" class="lg:hidden fixed bottom-0 left-0 right-0 p-3 bg-white/90 backdrop-blur-sm border-t border-slate-200 z-40 shadow-md print:hidden">
            <button @click="isMobileCartOpen = true" class="w-full bg-indigo-600 text-white p-3 rounded-xl flex justify-between items-center active:scale-95 transition-all">
                <span class="bg-white text-indigo-600 font-black px-3 py-1 rounded-lg text-xs">{{ cart.length }} Item</span>
                <span class="font-black text-sm">Rp {{ totalBelanja.toLocaleString('id-ID') }} ➔</span>
            </button>
        </div>

        <div v-if="showHeldModal" class="fixed inset-0 bg-slate-900/80 flex items-center justify-center z-[150] p-4 backdrop-blur-sm print:hidden">
            <div class="bg-white p-6 md:p-8 rounded-[32px] shadow-2xl w-full max-w-xl max-h-[80vh] flex flex-col">
                <div class="flex justify-between items-center mb-6">
                    <h2 class="text-lg font-black text-slate-800 uppercase tracking-widest flex items-center gap-2">💾 Pesanan Tertunda</h2>
                    <button @click="showHeldModal = false" class="p-2 bg-slate-100 hover:bg-rose-100 text-slate-400 rounded-xl">✕</button>
                </div>
                <div class="flex-1 overflow-y-auto custom-scrollbar pr-2">
                    <div v-if="heldOrders.length === 0" class="text-center py-10 text-slate-400 font-bold text-xs uppercase tracking-widest italic">Tidak ada pesanan ditunda.</div>
                    <div v-for="order in heldOrders" :key="order.id" class="p-4 bg-slate-50 border border-slate-200 rounded-2xl mb-3 flex justify-between items-center group">
                        <div>
                            <p class="font-black text-sm text-slate-800 uppercase">{{ order.customer }}</p>
                            <p class="text-[10px] font-bold text-slate-400 mt-1">Jam: {{ order.time }} | {{ order.items.length }} Item</p>
                            <p class="text-indigo-600 font-black mt-1 text-sm">Rp {{ order.total.toLocaleString('id-ID') }}</p>
                        </div>
                        <button @click="resumeOrder(order)" class="bg-indigo-100 text-indigo-600 hover:bg-indigo-600 hover:text-white px-4 py-2 rounded-xl font-black text-[10px] uppercase tracking-widest transition-all">Lanjutkan</button>
                    </div>
                </div>
            </div>
        </div>

        <div v-if="showScanner" class="fixed inset-0 bg-slate-900/90 flex items-center justify-center z-[150] p-4 backdrop-blur-sm print:hidden">
            <div class="bg-white rounded-[32px] shadow-2xl w-full max-w-md overflow-hidden flex flex-col">
                <div class="p-6 border-b flex justify-between items-center bg-slate-50/50">
                    <h2 class="text-lg font-black text-slate-800 uppercase tracking-widest">Scan Barcode</h2>
                    <button @click="stopScanner" class="p-2 rounded-xl bg-slate-100 text-slate-400">✕</button>
                </div>
                <div class="p-6 bg-black">
                    <div id="reader-kasir" class="w-full rounded-2xl overflow-hidden border-2 border-slate-700"></div>
                </div>
            </div>
        </div>

        <div v-if="showQrisModal" class="fixed inset-0 bg-slate-900/80 flex items-center justify-center z-[150] p-4 backdrop-blur-sm print:hidden">
            <div class="bg-white p-6 md:p-8 rounded-[32px] shadow-2xl w-full max-w-sm text-center flex flex-col border-t-8 border-indigo-600">
                <h2 class="text-xl font-black text-slate-900 uppercase tracking-widest mb-1">Bayar via QRIS</h2>
                <div class="bg-white p-3 rounded-2xl border border-slate-300 w-full mb-4 flex justify-center items-center min-h-[200px]">
                    <img src="https://placehold.co/300x300?text=QRIS+TOKO" alt="QRIS Toko" class="w-full h-full max-h-48 object-contain mx-auto rounded-xl">
                </div>
                <div class="bg-indigo-50 text-indigo-900 p-4 rounded-2xl mb-6 font-black text-2xl">Rp {{ totalBelanja.toLocaleString('id-ID') }}</div>
                <div class="flex gap-3">
                    <button @click="showQrisModal = false" :disabled="isProcessingCheckout" class="flex-1 bg-slate-100 py-4 rounded-xl font-black text-[10px] uppercase text-slate-500">Batal</button>
                    <button @click="executeCheckout" :disabled="isProcessingCheckout" class="flex-1 bg-indigo-600 py-4 rounded-xl font-black text-[10px] uppercase text-white shadow-lg">{{ isProcessingCheckout ? 'Proses...' : 'Lunas' }}</button>
                </div>
            </div>
        </div>

        <ReceiptModal 
            :show="showReceipt"
            :invoiceData="lastTransaction"
            :storeData="currentSession?.store || currentSession?.Store"
            :cashierName="currentUser?.name ? currentUser.name.split(' ')[0] : 'KASIR'"
            :stationNumber="currentSession?.station_number || '01'"
            @close="showReceipt = false" 
        />

        <ClosingModal 
            :show="showClosingModal"
            :showReceiptClosing="showReceiptClosing"
            :pecahan="pecahan"
            :totalUangFisik="totalUangFisik"
            :lastClosingData="lastClosingData"
            :currentSession="currentSession"
            :currentUser="currentUser"
            @close="showClosingModal = false"
            @process-closing="handleClosing"
            @print-closing="printClosing"
            @finish-closing="finishClosing"
        />

    </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
@media (min-width: 768px) { .custom-scrollbar::-webkit-scrollbar { width: 6px; } }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }

.animate-marquee { display: inline-block; padding-left: 100%; animation: marquee 20s linear infinite; }
@keyframes marquee { 0% { transform: translateX(0); } 100% { transform: translateX(-100%); } }

@keyframes slide-in-right { from { transform: translateX(100%); } to { transform: translateX(0); } }
.animate-slide-in-right { animation: slide-in-right 0.3s cubic-bezier(0.25, 0.8, 0.25, 1); }

@media print {
    @page { 
        margin: 0; 
    }
    body { 
        background: white; 
        -webkit-print-color-adjust: exact; 
    }
}

:deep(.swal2-popup) { font-family: 'Inter', sans-serif !important; border-radius: 28px !important; }
</style>