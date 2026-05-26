<script setup>
import { defineProps, defineEmits } from 'vue';

defineProps({
    isMobileCartOpen: Boolean,
    cart: Array,
    heldOrders: Array,
    paymentMethod: String,
    totalBelanja: Number,
    payAmount: Number,
    kembalian: Number,
    isProcessingCheckout: Boolean
});

// 🚀 Daftarkan jembatan event untuk memicu fungsi di usePos.js
const emit = defineEmits([
    'update:isMobileCartOpen', 
    'show-held', 
    'hold-order', 
    'clear-cart',
    'decrease-qty', 
    'increase-qty', 
    'validate-qty', 
    'set-payment',
    'format-rupiah', 
    'checkout'
]);
</script>

<template>
    <div :class="isMobileCartOpen ? 'fixed inset-0 z-[100] bg-slate-900/40 backdrop-blur-sm flex justify-end transition-all' : 'hidden lg:flex lg:relative w-4/12 xl:w-3/12'">
        <div :class="isMobileCartOpen ? 'w-[85%] sm:w-[360px] h-full animate-slide-in-right' : 'w-full h-full'" 
             class="bg-white lg:rounded-[32px] shadow-2xl border-l lg:border border-slate-100 flex flex-col shrink-0 overflow-hidden">
            
            <div class="lg:hidden p-4 bg-indigo-900 text-white flex justify-between items-center shrink-0">
                <h2 class="font-black tracking-widest uppercase text-sm flex items-center gap-2">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z" /></svg>
                    Keranjang
                </h2>
                <button @click="emit('update:isMobileCartOpen', false)" class="bg-white/20 p-2 rounded-xl text-white hover:bg-rose-500 transition-colors">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
                </button>
            </div>

            <div class="p-4 md:p-5 border-b border-slate-100 bg-slate-50/80 hidden lg:flex justify-between items-center shrink-0">
                <h2 class="text-sm md:text-base font-black text-slate-800 flex items-center gap-2 uppercase tracking-widest">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z" /></svg>
                    Keranjang
                </h2>
                <div class="flex gap-1.5 md:gap-2">
                    <button @click="emit('show-held')" class="p-2 bg-amber-50 hover:bg-amber-100 text-amber-600 rounded-xl transition-colors relative" title="Lihat Pesanan Tertunda">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 md:w-5 md:h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                        <span class="absolute -top-1 -right-1 flex h-4 w-4 items-center justify-center rounded-full bg-rose-500 text-[8px] font-black text-white shadow-md">{{ heldOrders.length }}</span>
                    </button>
                    <button @click="emit('hold-order')" :disabled="cart.length === 0" class="p-2 bg-indigo-50 hover:bg-indigo-100 text-indigo-600 rounded-xl transition-colors disabled:opacity-50" title="Hold Pesanan Ini">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 md:w-5 md:h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M10 9v6m4-6v6m7-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                    </button>
                    <button @click="emit('clear-cart')" :disabled="cart.length === 0" class="p-2 bg-rose-50 hover:bg-rose-100 text-rose-600 rounded-xl transition-colors disabled:opacity-50" title="Kosongkan Keranjang">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 md:w-5 md:h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h14" /></svg>
                    </button>
                </div>
            </div>

            <div class="p-3 md:p-4 flex-1 overflow-y-auto bg-white custom-scrollbar min-h-0">
                <div v-if="cart.length > 0" class="flex gap-2 lg:hidden mb-4">
                    <button @click="emit('hold-order')" class="flex-1 py-2 bg-amber-50 text-amber-600 rounded-lg text-[9px] font-black uppercase tracking-widest border border-amber-100">Hold Order</button>
                    <button @click="emit('clear-cart')" class="flex-1 py-2 bg-rose-50 text-rose-600 rounded-lg text-[9px] font-black uppercase tracking-widest border border-rose-100">Bersihkan</button>
                </div>

                <div v-if="cart.length === 0" class="h-full flex flex-col items-center justify-center opacity-30 py-12">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-12 h-12 md:w-16 md:h-16 mb-4 text-slate-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" /></svg>
                    <p class="text-slate-600 font-black text-xs md:text-sm uppercase tracking-widest">Keranjang Kosong</p>
                </div>

                <div v-for="item in cart" :key="item.id" class="flex flex-col mb-3 p-3 bg-slate-50/50 rounded-xl md:rounded-2xl border border-slate-100 shadow-sm hover:border-indigo-200 transition-colors">
                    <div class="flex justify-between items-start mb-2">
                        <h3 class="font-bold text-[10px] md:text-[11px] text-slate-800 leading-tight pr-2 line-clamp-2 uppercase">{{ item.name }}</h3>
                        <div class="font-black text-[11px] md:text-xs text-indigo-700 whitespace-nowrap">Rp {{ (item.price * item.qty).toLocaleString('id-ID') }}</div>
                    </div>
                    <div class="flex justify-between items-center">
                        <p class="text-[9px] md:text-[10px] font-bold text-slate-400">@ Rp {{ item.price.toLocaleString('id-ID') }}</p>
                        <div class="flex items-center bg-white rounded-lg md:rounded-xl p-1 border border-slate-200 shadow-sm">
                            <button @click="emit('decrease-qty', item)" class="w-6 h-6 md:w-7 md:h-7 flex items-center justify-center rounded-md md:rounded-lg text-slate-400 hover:bg-rose-50 hover:text-rose-600 font-black transition-colors">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M20 12H4" /></svg>
                            </button>
                            <input type="number" v-model.number="item.qty" @change="emit('validate-qty', item)" class="w-8 md:w-10 text-center text-xs font-black text-slate-800 bg-transparent border-none focus:ring-0 p-0 [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none">
                            <button @click="emit('increase-qty', item)" class="w-6 h-6 md:w-7 md:h-7 flex items-center justify-center rounded-md md:rounded-lg text-slate-400 hover:bg-indigo-50 hover:text-indigo-600 font-black transition-colors">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" /></svg>
                            </button>
                        </div>
                    </div>
                </div>
            </div>

            <div class="p-3 md:p-5 bg-white border-t border-slate-100 shadow-[0_-10px_20px_-10px_rgba(0,0,0,0.05)] shrink-0 z-10 lg:rounded-b-[32px]">
                <div class="mb-3 md:mb-4">
                    <span class="font-black text-[8px] md:text-[9px] text-slate-400 block mb-1.5 md:mb-2 uppercase tracking-widest text-center">Metode Pembayaran</span>
                    <div class="grid grid-cols-3 gap-1.5 md:gap-2">
                        <button v-for="method in ['Cash', 'QRIS', 'Debit']" :key="method" @click="emit('set-payment', method)" :class="paymentMethod === method ? 'bg-indigo-600 text-white shadow-md shadow-indigo-200 border-indigo-600' : 'bg-slate-50 text-slate-500 border-slate-200 hover:bg-slate-100'" class="py-2.5 rounded-lg md:rounded-xl font-black text-[9px] md:text-[10px] uppercase transition-all flex flex-col items-center gap-1.5 border">
                            <svg v-if="method === 'Cash'" xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 md:w-5 md:h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z" /></svg>
                            <svg v-else-if="method === 'QRIS'" xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 md:w-5 md:h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm14 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z" /></svg>
                            <svg v-else-if="method === 'Debit'" xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 md:w-5 md:h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z" /></svg>
                            {{ method === 'Cash' ? 'Tunai' : method }}
                        </button>
                    </div>
                </div>

                <div class="space-y-3 md:space-y-4 mb-4 md:mb-5">
                    <div class="flex justify-between items-end border-b border-dashed border-slate-200 pb-1.5 md:pb-2">
                        <span class="font-black text-[9px] md:text-[10px] text-slate-400 uppercase tracking-widest">Total Tagihan</span>
                        <span class="text-2xl md:text-3xl font-black text-indigo-800 leading-none tracking-tighter">Rp {{ totalBelanja.toLocaleString('id-ID') }}</span>
                    </div>

                    <div class="flex justify-between items-center bg-slate-50 p-2 md:p-2.5 rounded-xl md:rounded-2xl border-2 border-slate-100 focus-within:border-indigo-500 transition-all">
                        <span class="font-black text-[9px] md:text-[10px] text-slate-600 uppercase tracking-widest pl-2">Bayar</span>
                        <div class="relative flex-1 ml-2 md:ml-4">
                            <span class="absolute left-2.5 md:left-3 top-1/2 -translate-y-1/2 text-slate-400 font-bold text-xs md:text-sm italic">Rp</span>
                            <input
                                type="text"
                                :value="payAmount === 0 ? '' : payAmount.toLocaleString('id-ID')"
                                @input="emit('format-rupiah', $event)"
                                :disabled="paymentMethod !== 'Cash'"
                                :class="paymentMethod !== 'Cash' ? 'bg-slate-200/50 text-slate-400 cursor-not-allowed border-transparent' : 'bg-white text-slate-900 border-slate-200 shadow-sm'"
                                class="w-full text-right text-base md:text-lg font-black rounded-lg md:rounded-xl py-1.5 md:py-2 pl-8 pr-3 transition-all outline-none border"
                                placeholder="0">
                        </div>
                    </div>

                    <div class="flex justify-between items-center px-1">
                        <span class="font-black text-[9px] md:text-[10px] text-slate-400 uppercase tracking-widest">Kembali</span>
                        <span class="text-lg md:text-xl font-black" :class="kembalian >= 0 ? 'text-emerald-500' : 'text-rose-500'">
                            Rp {{ kembalian.toLocaleString('id-ID') }}
                        </span>
                    </div>
                </div>

                <button @click="emit('checkout')" :disabled="cart.length === 0 || payAmount < totalBelanja || isProcessingCheckout"
                    class="w-full bg-emerald-500 hover:bg-emerald-600 text-white font-black py-3 md:py-4 px-4 rounded-xl md:rounded-2xl transition-all flex justify-center items-center gap-2 md:gap-3 disabled:opacity-50 disabled:cursor-not-allowed shadow-xl shadow-emerald-200 hover:shadow-emerald-300 active:scale-95 text-xs md:text-sm uppercase tracking-widest">
                    <template v-if="isProcessingCheckout">
                        <div class="w-4 h-4 md:w-5 md:h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                        Memproses...
                    </template>
                    <template v-else>
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 md:w-5 md:h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" /></svg>
                        Proses Bayar
                    </template>
                </button>
            </div>
        </div>
    </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
@media (min-width: 768px) { .custom-scrollbar::-webkit-scrollbar { width: 6px; } }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #94a3b8; }
</style>