<script setup>
  defineProps({
    isMobileCartOpen: Boolean,
    cart: Array,
    heldOrders: Array,
    paymentMethod: String,

    subTotal: { type: Number, default: 0 },
    pajak: { type: Number, default: 0 },
    isTaxActive: { type: Boolean, default: false },
    pajakPersen: { type: Number, default: 0 },

    totalBelanja: { type: Number, default: 0 },
    payAmount: { type: Number, default: 0 },
    kembalian: { type: Number, default: 0 },
    isProcessingCheckout: Boolean,
  });

  const emit = defineEmits([
    'update:isMobileCartOpen',
    'show-held',
    'hold-order',
    'clear-cart',
    'decrease-qty',
    'increase-qty',
    'validate-qty',
    'set-nominal',
    'set-payment',
    'format-rupiah',
    'checkout',
    'toggle-uom',
  ]);

  // 🛡️ SECURITY LAYER FOR FRONTEND INPUT: Buang paksa karakter aneh di kuantitas barang bray!
  const filterQtyKeyboard = (event, item) => {
    let cleanVal = event.target.value.replace(/\D/g, ""); // Murni angka bulat positif, buang minus & desimal!
    let parsed = parseInt(cleanVal, 10);
    
    // Paksa minimal order adalah 1 PCS bray
    item.qty = parsed && parsed > 0 ? parsed : 1;
  };
</script>

<template>
  <div
    :class="[
      isMobileCartOpen
        ? 'fixed inset-0 z-[100] bg-slate-900/50 backdrop-blur-sm flex justify-end transition-opacity duration-300'
        : 'hidden lg:flex flex-col shrink-0 w-[300px] xl:w-[360px] 2xl:w-[400px]'
    ]"
  >
    <div
      :class="[
        isMobileCartOpen
          ? 'w-[85%] sm:w-[340px] h-full max-h-screen translate-x-0 transition-transform duration-300 ease-out shadow-[-10px_0_30px_rgba(0,0,0,0.1)]'
          : 'w-full h-full'
      ]"
      class="bg-white lg:rounded-[24px] lg:border border-slate-200 flex flex-col h-full overflow-hidden shadow-lg"
    >
      <div class="lg:hidden p-3 bg-indigo-900 text-white flex justify-between items-center shrink-0">
        <h2 class="font-black tracking-widest uppercase text-xs flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z"/></svg>
          Keranjang
        </h2>
        <button @click="emit('update:isMobileCartOpen', false)" class="bg-white/20 p-1.5 rounded-lg text-white hover:bg-rose-500 transition-colors">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/></svg>
        </button>
      </div>

      <div class="p-3 border-b border-slate-100 bg-slate-50/80 hidden lg:flex justify-between items-center shrink-0">
        <h2 class="text-xs xl:text-sm font-black text-slate-800 flex items-center gap-2 uppercase tracking-widest">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 xl:w-5 xl:h-5 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z"/></svg>
          Keranjang
        </h2>
        <div class="flex gap-1 xl:gap-1.5">
          <button @click="emit('show-held')" class="p-1.5 xl:p-2 bg-amber-50 hover:bg-amber-100 text-amber-600 rounded-lg xl:rounded-xl transition-colors relative" title="Lihat Pesanan Tertunda">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 xl:w-4 xl:h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
            <span v-if="heldOrders.length > 0" class="absolute -top-1 -right-1 flex h-3.5 w-3.5 xl:h-4 xl:w-4 items-center justify-center rounded-full bg-rose-500 text-[7px] xl:text-[8px] font-black text-white shadow-sm">{{ heldOrders.length }}</span>
          </button>
          <button @click="emit('hold-order')" :disabled="cart.length === 0" class="p-1.5 xl:p-2 bg-indigo-50 hover:bg-indigo-100 text-indigo-600 rounded-lg xl:rounded-xl transition-colors disabled:opacity-50" title="Hold Pesanan Ini">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 xl:w-4 xl:h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M10 9v6m4-6v6m7-3a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
          </button>
          <button @click="emit('clear-cart')" :disabled="cart.length === 0" class="p-1.5 xl:p-2 bg-rose-50 hover:bg-rose-100 text-rose-600 rounded-lg xl:rounded-xl transition-colors disabled:opacity-50" title="Kosongkan Keranjang">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 xl:w-4 xl:h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h14"/></svg>
          </button>
        </div>
      </div>

      <div class="p-2 xl:p-3 flex-1 overflow-y-auto bg-white custom-scrollbar min-h-[120px] relative">
        <div v-if="cart.length > 0" class="flex gap-2 lg:hidden mb-3">
          <button @click="emit('hold-order')" class="flex-1 py-2 bg-amber-50 text-amber-600 rounded-lg text-[9px] font-black uppercase tracking-widest border border-amber-100 active:scale-95 transition-all">Hold Order</button>
          <button @click="emit('clear-cart')" class="flex-1 py-2 bg-rose-50 text-rose-600 rounded-lg text-[9px] font-black uppercase tracking-widest border border-rose-100 active:scale-95 transition-all">Bersihkan</button>
        </div>

        <div v-if="cart.length === 0" class="h-full flex flex-col items-center justify-center opacity-40 py-10">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-10 h-10 xl:w-12 xl:h-12 mb-2 text-slate-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z"/></svg>
          <p class="text-slate-600 font-black text-[10px] xl:text-[11px] uppercase tracking-widest text-center">Keranjang Kosong</p>
        </div>

        <div v-for="item in cart" :key="item.id + '_' + item.selected_uom" class="flex flex-col mb-2 p-2 xl:p-3 bg-slate-50/50 rounded-xl border border-slate-100 shadow-sm hover:border-indigo-200 transition-all relative group">
          <div class="flex justify-between items-start mb-1.5">
            <div class="flex-1 pr-1 xl:pr-2">
              <h3 class="font-bold text-[10px] xl:text-[11px] text-slate-800 leading-tight line-clamp-2 uppercase">{{ item.name }}</h3>
              <div class="mt-1 flex flex-wrap items-center gap-1 xl:gap-1.5">
                <span
                  class="px-1.5 py-0.5 rounded text-[7px] xl:text-[8px] font-black uppercase tracking-widest border"
                  :class="{
                    'bg-indigo-50 border-indigo-200 text-indigo-600': item.selected_uom === item.satuan_dasar,
                    'bg-emerald-50 border-emerald-200 text-emerald-600': item.selected_uom === item.satuan_tengah,
                    'bg-amber-100 border-amber-300 text-amber-700': item.selected_uom === item.satuan_besar,
                  }"
                >
                  {{ item.selected_uom || 'PCS' }}
                </span>
                <button v-if="item.has_grosir" @click="emit('toggle-uom', item)" class="px-1.5 py-0.5 bg-slate-200 hover:bg-slate-300 text-slate-600 rounded text-[7px] xl:text-[8px] font-black uppercase tracking-widest flex items-center justify-center transition-colors">
                  Ubah
                </button>
              </div>
            </div>
            <div class="font-black text-[11px] xl:text-xs text-indigo-700 whitespace-nowrap text-right">
              Rp {{ (item.price * item.qty).toLocaleString('id-ID') }}
            </div>
          </div>

          <div class="flex justify-between items-center">
            <p class="text-[9px] xl:text-[10px] font-bold text-slate-400">@ Rp {{ item.price.toLocaleString('id-ID') }}</p>
            <div class="flex items-center bg-white rounded-lg p-0.5 border border-slate-200 shadow-sm">
              <button @click="emit('decrease-qty', item)" class="w-6 h-6 xl:w-7 xl:h-7 flex items-center justify-center rounded text-slate-400 hover:bg-rose-50 hover:text-rose-600 font-black transition-colors">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3 xl:w-4 xl:h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M20 12H4"/></svg>
              </button>
              
              <input 
                type="text" 
                :value="item.qty" 
                @input="filterQtyKeyboard($event, item)"
                @change="emit('validate-qty', item)" 
                class="w-8 xl:w-10 text-center text-[11px] xl:text-xs font-black text-slate-800 bg-transparent border-none focus:ring-0 p-0" 
              />
              
              <button @click="emit('increase-qty', item)" class="w-6 h-6 xl:w-7 xl:h-7 flex items-center justify-center rounded text-slate-400 hover:bg-indigo-50 hover:text-indigo-600 font-black transition-colors">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3 xl:w-4 xl:h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/></svg>
              </button>
            </div>
          </div>
        </div>
      </div>

      <div class="p-2 md:p-3 bg-white border-t border-slate-100 shadow-[0_-5px_15px_-5px_rgba(0,0,0,0.05)] shrink-0 z-10 lg:rounded-b-[24px]">
        
        <div class="mb-1.5 md:mb-2">
          <div class="grid grid-cols-3 gap-1.5">
            <button
              v-for="method in ['Cash', 'QRIS', 'Debit']"
              :key="method"
              @click="emit('set-payment', method)"
              :class="
                paymentMethod === method
                  ? 'bg-indigo-600 text-white shadow-md shadow-indigo-200 border-indigo-600'
                  : 'bg-slate-50 text-slate-500 border-slate-200 hover:bg-slate-100'
              "
              class="py-1.5 rounded-lg font-black text-[9px] uppercase transition-all flex justify-center items-center border active:scale-95"
            >
              {{ method === 'Cash' ? 'Tunai' : method }}
            </button>
          </div>
        </div>

        <div class="space-y-0.5 bg-slate-50 p-1.5 md:p-2 rounded-xl border border-slate-100 mb-1.5 md:mb-2">
          <div class="flex justify-between items-end">
            <span class="font-bold text-[9px] text-slate-500 uppercase tracking-widest">Subtotal</span>
            <span class="text-xs font-bold text-slate-700 leading-none">Rp {{ subTotal ? subTotal.toLocaleString('id-ID') : '0' }}</span>
          </div>

          <div v-if="isTaxActive" class="flex justify-between items-end pb-0.5 border-b border-dashed border-slate-200">
            <span class="font-bold text-[9px] text-slate-500 uppercase tracking-widest">Pajak ({{ pajakPersen }}%)</span>
            <span class="text-xs font-bold text-slate-700 leading-none">Rp {{ pajak ? pajak.toLocaleString('id-ID') : '0' }}</span>
          </div>
          
          <div class="flex justify-between items-end pt-0.5">
            <span class="font-black text-[10px] text-slate-600 uppercase tracking-widest">Total</span>
            <span class="text-lg md:text-xl font-black text-indigo-800 leading-none tracking-tighter">Rp {{ totalBelanja ? totalBelanja.toLocaleString('id-ID') : '0' }}</span>
          </div>
        </div>

        <div v-show="paymentMethod === 'QRIS'" class="bg-indigo-50 p-1.5 rounded-xl border border-indigo-100 flex items-center justify-between animate-[fadeInUp_0.2s_ease-out] mb-1.5">
            <div class="flex items-center gap-2">
                <div class="p-1 bg-indigo-100 text-indigo-600 rounded">
                  <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm14 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z"/></svg>
                </div>
                <div class="text-[8px] font-black text-indigo-800 uppercase tracking-widest leading-none">Uang Pas<br><span class="text-indigo-500 font-bold">Otomatis</span></div>
            </div>
            <div class="text-right">
                <div class="text-[7px] font-black text-indigo-400 uppercase tracking-widest">Validasi</div>
                <div class="text-[8px] font-black text-indigo-600 uppercase tracking-widest">Midtrans/Statis</div>
            </div>
        </div>

        <div v-show="paymentMethod === 'Debit'" class="bg-slate-100 p-1.5 rounded-xl border border-slate-200 flex items-center justify-between animate-[fadeInUp_0.2s_ease-out] mb-1.5">
            <div class="flex items-center gap-2">
                <div class="p-1 bg-slate-200 text-slate-500 rounded">
                  <svg xmlns="http://www.w3.org/2000/xl" class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z"/></svg>
                </div>
                <div class="text-[8px] font-black text-slate-600 uppercase tracking-widest leading-none">Uang Pas<br><span class="text-slate-400 font-bold">Otomatis</span></div>
            </div>
            <div class="text-right"><div class="text-[8px] font-black text-slate-400 uppercase tracking-widest">EDC External</div></div>
        </div>

        <div v-show="paymentMethod === 'Cash'" class="animate-[fadeInUp_0.2s_ease-out] mb-1.5">
          <div class="flex overflow-x-auto gap-1 custom-scrollbar pb-1.5">
            <button v-for="nom in [100000, 50000, 20000, 10000, 5000, 1000, 500, 100]" :key="nom" @click="emit('set-nominal', nom)" class="shrink-0 px-3 py-1.5 bg-white border border-slate-200 hover:border-indigo-400 hover:text-indigo-600 rounded-lg font-black text-[9px] text-slate-600 uppercase transition-all shadow-sm active:scale-95">
              {{ nom >= 1000 ? nom / 1000 + ' Rb' : nom }}
            </button>
          </div>
        </div>

        <div class="relative z-20 flex justify-between items-center bg-slate-50 p-1.5 rounded-xl border border-slate-200 focus-within:border-indigo-500 transition-all shadow-sm mb-1.5">
          <span class="font-black text-[10px] text-slate-600 uppercase tracking-widest pl-1">Bayar</span>
          <div class="flex-1 flex gap-1 ml-2">
            <div class="relative flex-1">
              <span class="absolute left-2 top-1/2 -translate-y-1/2 text-slate-400 font-bold text-xs italic">Rp</span>
              <input
                type="text"
                :value="paymentMethod === 'Cash' ? (payAmount === 0 ? '' : payAmount.toLocaleString('id-ID')) : (totalBelanja ? totalBelanja.toLocaleString('id-ID') : '0')"
                @input="emit('format-rupiah', $event)"
                :disabled="paymentMethod !== 'Cash'"
                :class="paymentMethod !== 'Cash' ? 'bg-slate-200/50 text-slate-500 cursor-not-allowed border-transparent shadow-inner' : 'bg-white text-slate-900 border-slate-300 shadow-sm focus:border-indigo-500'"
                class="w-full text-right text-sm font-black rounded-lg py-1.5 pl-7 pr-2 transition-all outline-none border"
                placeholder="0"
              />
            </div>
            
            <button v-show="paymentMethod === 'Cash'" @click="emit('set-nominal', 0)" class="w-8 h-[34px] bg-rose-50 border border-rose-200 hover:bg-rose-100 text-rose-600 rounded-lg flex items-center justify-center font-black transition-colors active:scale-95" title="Clear Nominal">
              C
            </button>
          </div>
        </div>

        <div class="flex justify-between items-center px-1 mb-1.5">
          <span class="font-black text-[9px] text-slate-400 uppercase tracking-widest">Kembali</span>
          <span class="text-xs md:text-sm font-black" :class="kembalian >= 0 && paymentMethod === 'Cash' ? 'text-emerald-500' : 'text-red-500'">
            Rp {{ paymentMethod === 'Cash' ? (kembalian ? kembalian.toLocaleString('id-ID') : '0') : '0' }}
          </span>
        </div>

        <button
          @click="emit('checkout')"
          :disabled="cart.length === 0 || (paymentMethod === 'Cash' && payAmount < totalBelanja) || isProcessingCheckout"
          class="w-full text-white font-black py-2 md:py-2.5 rounded-[12px] transition-all flex justify-center items-center gap-2 disabled:opacity-50 shadow-md text-xs uppercase tracking-[0.1em]"
          :class="paymentMethod === 'QRIS' ? 'bg-indigo-600 hover:bg-indigo-700' : 'bg-emerald-500 hover:bg-emerald-600'"
        >
          <template v-if="isProcessingCheckout">
            <div class="w-3.5 h-3.5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
            Memproses...
          </template>
          <template v-else>
            {{ paymentMethod === 'QRIS' ? 'Generate QRIS' : 'Proses Bayar' }}
          </template>
        </button>

      </div>
    </div>
  </div>
</template>