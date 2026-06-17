<script setup>
const props = defineProps({
	isCartOpen: Boolean,
	cart: Array,
	customerName: String,
	customerPhone: String,
	estimasiSelesai: String,
	customerResults: Array,
	showCustomerDropdown: Boolean,
	availablePerfumes: Array,
	paymentMethod: String,
	mainPaymentGroup: String,
	totalTagihan: Number,
	formattedUangBayar: String,
	kembalian: Number,
	isSubmitting: Boolean,
	storeInfo: Object,
	photoData: String,
	formatRupiah: Function,
});

const emit = defineEmits(['update:isCartOpen', 'update:customerName', 'update:customerPhone', 'update:estimasiSelesai', 'search-customer', 'select-customer', 'close-dropdown', 'clear-cart', 'remove-item', 'perfume-change', 'update-berat', 'open-camera', 'remove-photo', 'set-payment', 'set-nominal', 'checkout', 'format-no-hp']);
</script>

<template>
	<div v-if="isCartOpen" @click="emit('update:isCartOpen', false)" class="fixed inset-0 bg-slate-900/40 z-40 lg:hidden backdrop-blur-sm transition-opacity"></div>

	<div :class="isCartOpen ? 'translate-x-0' : 'translate-x-full lg:translate-x-0'" class="fixed inset-y-0 right-0 z-50 w-full sm:w-[380px] lg:w-[330px] xl:w-[380px] bg-white shadow-2xl lg:static lg:shadow-none lg:border-l border-slate-200 flex flex-col transition-transform duration-300 ease-in-out shrink-0 h-full overflow-hidden">
		<div class="bg-white text-slate-800 p-4 shrink-0 border-b border-slate-200 flex justify-between items-center">
			<h2 class="text-xs font-black uppercase tracking-widest flex items-center gap-2 text-slate-800">
				<div class="bg-slate-100 text-slate-800 p-1.5 rounded-lg border border-slate-200 shadow-sm">
					<svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
						<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
						<polyline points="14 2 14 8 20 8" />
					</svg>
				</div>
				Rincian Nota Timbang
			</h2>
			<button @click="emit('clear-cart')" class="text-[9px] font-black uppercase text-rose-600 hover:bg-rose-50 px-3 py-2 rounded-xl transition-colors">Batal</button>
		</div>

		<div class="flex-1 overflow-y-auto custom-scrollbar bg-slate-50/50 flex flex-col">
			<div class="p-4 border-b border-slate-200 bg-white flex flex-col gap-4 shadow-sm relative z-10 shrink-0">
				<div class="relative">
					<label class="block text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1.5 ml-0.5">Nama Pelanggan / Cari Database</label>
					<input
						:value="customerName"
						@input="
							emit('update:customerName', $event.target.value);
							emit('search-customer');
						"
						@blur="emit('close-dropdown')"
						@focus="emit('search-customer')"
						type="text"
						class="w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-slate-800 outline-none font-bold text-slate-800 text-xs placeholder:text-slate-300 shadow-inner"
						placeholder="Ketik nama pelanggan..." />
					<div v-if="showCustomerDropdown" class="absolute left-0 right-0 top-[105%] z-50 bg-white border border-slate-200 rounded-xl shadow-2xl overflow-hidden">
						<div class="max-h-40 overflow-y-auto custom-scrollbar">
							<div v-for="cust in customerResults" :key="cust.id" @click="emit('select-customer', cust)" class="p-3 hover:bg-slate-50 cursor-pointer border-b border-slate-100 flex flex-col transition-colors">
								<span class="font-black text-xs text-slate-800">{{ cust.nama }}</span>
								<span class="text-[9px] font-bold text-slate-400 mt-0.5">WhatsApp: +{{ cust.no_whatsapp }}</span>
							</div>
						</div>
					</div>
				</div>

				<div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
					<div>
						<label class="block text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1.5 ml-0.5">No. WhatsApp</label>
						<div class="flex items-center bg-slate-50 border border-slate-200 rounded-xl focus-within:bg-white focus-within:border-slate-800 transition-all overflow-hidden shadow-inner">
							<div class="pl-3 pr-2 py-3 bg-slate-100 border-r border-slate-200 shrink-0"><span class="text-slate-500 font-black text-xs">+62</span></div>
							<input
								:value="customerPhone"
								@input="
									emit('update:customerPhone', $event.target.value);
									emit('format-no-hp');
								"
								type="number"
								class="w-full px-3 py-3 bg-transparent outline-none font-bold text-slate-800 text-xs" />
						</div>
					</div>
					<div>
						<label class="block text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1.5 ml-0.5">Estimasi Selesai</label>
						<input :value="estimasiSelesai" @input="emit('update:estimasiSelesai', $event.target.value)" type="date" class="w-full px-3 py-2.5 bg-slate-50 border border-slate-200 rounded-xl outline-none font-bold text-xs text-slate-700 focus:bg-white focus:border-slate-800 transition-all shadow-inner" />
					</div>
				</div>

				<div class="mt-1 border border-dashed border-slate-200 rounded-xl p-3 bg-slate-50 flex flex-col items-center justify-center relative overflow-hidden min-h-[110px]">
					<label class="block text-[9px] font-black text-slate-400 uppercase tracking-widest mb-3 absolute top-3 left-3 z-10">Bukti Dokumentasi Barang</label>
					<button v-if="!photoData" @click="emit('open-camera', 'ITEM')" class="flex flex-col items-center text-slate-400 hover:text-slate-800 transition-colors mt-4">
						<svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 mb-1.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
							<path stroke-linecap="round" stroke-linejoin="round" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
							<circle cx="12" cy="13" r="3" />
						</svg>
						<span class="text-[9px] font-black uppercase tracking-widest">Ambil Foto Masuk</span>
					</button>
					<div v-else class="w-full h-full relative z-20 mt-4 rounded-lg overflow-hidden border border-slate-200">
						<img :src="photoData" class="w-full h-28 object-cover" />
						<button @click="emit('remove-photo')" class="absolute top-2 right-2 bg-rose-600 text-white w-5 h-5 rounded-full flex items-center justify-center shadow-lg">
							<svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
								<line x1="18" y1="6" x2="6" y2="18" />
								<line x1="6" y1="6" x2="18" y2="18" />
							</svg>
						</button>
					</div>
				</div>
			</div>

			<div class="p-4 flex flex-col gap-3">
				<div v-if="cart.length === 0" class="flex flex-col items-center justify-center text-center opacity-40 py-12">
					<p class="text-slate-500 font-black text-[10px] uppercase tracking-widest">Belum Ada Item Ditimbang</p>
				</div>
				<div v-else v-for="(item, index) in cart" :key="index" class="bg-white p-4 rounded-2xl border border-slate-200 relative flex flex-col gap-3 shadow-sm animate-[fadeIn_0.2s_ease-out]">
					<button @click="emit('remove-item', index)" class="absolute -right-1.5 -top-1.5 w-5 h-5 bg-white text-rose-600 border border-slate-200 rounded-full flex items-center justify-center shadow-sm z-10">
						<svg xmlns="http://www.w3.org/2000/svg" class="w-2.5 h-2.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
							<line x1="18" y1="6" x2="6" y2="18" />
							<line x1="6" y1="6" x2="18" y2="18" />
						</svg>
					</button>
					<div class="flex justify-between items-start">
						<div class="pr-2">
							<h4 class="text-xs font-black text-slate-800 uppercase leading-tight">{{ item.nama_produk }}</h4>
							<p class="text-[9px] font-bold text-slate-400 mt-0.5 uppercase tracking-wider">@ {{ formatRupiah(item.harga) }} /{{ item.satuan_dasar }}</p>
						</div>
						<span class="text-xs font-black text-slate-800 bg-slate-50 px-2 py-0.5 rounded border border-slate-200 shrink-0">{{ formatRupiah(item.harga * parseFloat(item.berat) + item.harga_parfum) }}</span>
					</div>
					<div class="flex flex-col gap-1 bg-slate-50 p-2 rounded-xl border border-slate-200">
						<label class="text-[8px] font-black text-slate-500 uppercase tracking-widest ml-0.5">Varian Aroma Parfum</label>
						<select @change="emit('perfume-change', index, $event)" class="w-full text-[10px] font-bold bg-white text-slate-700 border border-slate-200 rounded-lg px-2 py-1 outline-none cursor-pointer focus:border-slate-800 transition-colors">
							<option value="default">Parfum Standar Toko (Bawaan - Gratis)</option>
							<option v-for="perfume in availablePerfumes" :key="perfume.id" :value="perfume.id">{{ perfume.nama }} {{ perfume.harga > 0 ? `(+ ${formatRupiah(perfume.harga)})` : '(Gratis)' }}</option>
						</select>
					</div>
					<div class="flex items-center justify-between bg-slate-50 rounded-xl p-1.5 border border-slate-200">
						<span class="text-[9px] font-black text-slate-500 uppercase tracking-widest pl-1.5">{{ item.satuan_dasar === 'KG' ? 'Berat Timbangan:' : 'Jumlah Item:' }}</span>
						<div class="flex items-center bg-white rounded-lg border border-slate-200 overflow-hidden shadow-sm">
							<button @click="emit('update-berat', index, -0.5)" class="w-8 h-8 flex items-center justify-center text-slate-500 active:bg-slate-100">
								<svg xmlns="http://www.w3.org/2000/svg" class="w-2.5 h-2.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="5" y1="12" x2="19" y2="12" /></svg>
							</button>
							<input v-model="item.berat" type="number" step="0.1" min="0.1" class="w-10 text-center font-black text-xs text-slate-800 outline-none bg-transparent" />
							<button @click="emit('update-berat', index, 0.5)" class="w-8 h-8 flex items-center justify-center text-slate-500 active:bg-slate-100 border-l border-slate-100">
								<svg xmlns="http://www.w3.org/2000/svg" class="w-2.5 h-2.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
									<line x1="12" y1="5" x2="12" y2="19" />
									<line x1="5" y1="12" x2="19" y2="12" />
								</svg>
							</button>
							<span class="text-[9px] font-black text-slate-400 bg-slate-50 h-8 flex items-center px-2 border-l border-slate-200 uppercase">{{ item.satuan_dasar }}</span>
						</div>
					</div>
				</div>
			</div>
		</div>

		<div class="bg-white border-t border-slate-200 shrink-0 p-4 shadow-[0_-10px_40px_-15px_rgba(0,0,0,0.05)] z-20">
			<div class="mb-3 grid grid-cols-2 gap-2">
				<button @click="emit('set-payment', 'TUNAI')" :class="mainPaymentGroup === 'Cash' ? 'bg-slate-900 border-slate-900 text-white font-black' : 'bg-white border-slate-200 text-slate-500'" class="py-2.5 rounded-xl border flex items-center justify-center gap-1.5 text-[10px] font-black uppercase tracking-wider transition-all active:scale-95">
					<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
					Tunai
				</button>
				<button @click="emit('set-payment', 'QRIS')" :class="mainPaymentGroup === 'Non-Cash' ? 'bg-slate-900 border-slate-900 text-white font-black' : 'bg-white border-slate-200 text-slate-500'" class="py-2.5 rounded-xl border flex items-center justify-center gap-1.5 text-[10px] font-black uppercase tracking-wider transition-all active:scale-95">
					<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z" /></svg>
					Non-Tunai
				</button>
			</div>

			<div v-show="mainPaymentGroup === 'Non-Cash'" class="mb-3 bg-slate-50 p-1 rounded-xl border border-slate-200/60 flex gap-1 animate-[fadeIn_0.15s_ease-out]">
				<button @click="emit('set-payment', 'QRIS')" :class="paymentMethod === 'QRIS' ? 'bg-white text-slate-900 border-slate-200 font-black shadow-sm' : 'text-slate-400 font-bold'" class="flex-1 py-1.5 rounded-lg text-[9px] uppercase tracking-widest border border-transparent transition-all flex items-center justify-center gap-1">QRIS</button>
				<button @click="emit('set-payment', 'DEBIT')" :class="paymentMethod === 'DEBIT' ? 'bg-white text-slate-900 border-slate-200 font-black shadow-sm' : 'text-slate-400 font-bold'" class="flex-1 py-1.5 rounded-lg text-[9px] uppercase tracking-widest border border-transparent transition-all flex items-center justify-center gap-1">Debit Bank</button>
				<button @click="emit('set-payment', 'PAYLATER')" :class="paymentMethod === 'PAYLATER' ? 'bg-amber-600 text-white border-amber-600 font-black shadow-sm' : 'text-amber-600/60 font-bold'" class="flex-1 py-1.5 rounded-lg text-[9px] uppercase tracking-widest border border-transparent transition-all flex items-center justify-center gap-1">Bayar Nanti</button>
			</div>

			<div class="space-y-2 mb-3 bg-slate-50 p-3 rounded-xl border border-slate-200 shadow-inner">
				<div class="flex justify-between items-center">
					<span class="text-[9px] font-black text-slate-400 uppercase tracking-widest">GRAND TOTAL</span>
					<span class="text-lg font-black text-slate-900 tracking-tighter">{{ formatRupiah(totalTagihan) }}</span>
				</div>

				<div class="pt-2 border-t border-slate-200/60 animate-[fadeIn_0.2s_ease-out]" v-if="paymentMethod === 'TUNAI'">
					<div class="flex overflow-x-auto gap-1 mb-2.5 custom-scrollbar pb-1">
						<button v-for="nom in [100000, 50000, 20000, 10000, 5000, 2000, 1000]" :key="nom" type="button" @click="emit('set-nominal', nom)" class="shrink-0 px-2.5 py-1 bg-white border border-slate-200 hover:border-slate-800 rounded-lg font-black text-[9px] text-slate-600 transition-all active:scale-95">{{ nom / 1000 }}Rb</button>
					</div>

					<div class="flex justify-between items-center mb-2.5">
						<span class="text-[9px] font-black text-slate-500 uppercase tracking-widest">Uang Cash</span>
						<div class="flex items-center gap-1.5">
							<div class="relative">
								<span class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 font-bold text-[10px] italic">Rp</span>
								<input :value="formattedUangBayar" @input="emit('set-nominal', $event.target.value)" type="text" class="w-28 bg-white border border-slate-200 rounded-xl pl-8 pr-3 py-1.5 text-right font-black text-xs text-slate-800 outline-none focus:border-slate-800 shadow-sm" placeholder="0" />
							</div>
							<button type="button" @click="emit('set-nominal', 'BACKSPACE')" class="w-8 h-[30px] bg-white border border-slate-200 hover:bg-slate-100 text-slate-700 rounded-xl flex items-center justify-center transition-colors active:scale-95 shadow-sm">
								<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2M3 12l6.414 6.414A2 2 0 0010.828 19H19a2 2 0 002-2V7a2 2 0 00-2-2h-8.172a2 2 0 00-1.414.586L3 12z" /></svg>
							</button>
							<button type="button" @click="emit('set-nominal', 0)" class="w-8 h-[30px] bg-rose-50 border border-rose-200 hover:bg-rose-100 text-rose-600 rounded-xl flex items-center justify-center font-black text-xs transition-colors active:scale-95 shadow-sm">C</button>
						</div>
					</div>
					<div class="flex justify-between items-center bg-white p-2 rounded-xl border border-slate-100">
						<span class="text-[9px] font-black text-slate-400 uppercase tracking-widest">Uang Kembali</span>
						<span class="text-xs font-black tracking-tight" :class="kembalian < 0 ? 'text-rose-600' : 'text-emerald-600'">{{ formatRupiah(kembalian) }}</span>
					</div>
				</div>

				<div v-show="paymentMethod === 'QRIS'" class="pt-2 border-t border-slate-200/60 flex items-center justify-between text-[9px] font-black text-slate-400 uppercase tracking-widest animate-[fadeIn_0.2s_ease-out]">
					<span>Satelit Routing:</span>
					<span :class="(storeInfo.payment_type || 'qris_static') === 'midtrans' ? 'text-indigo-600' : 'text-slate-600'" class="font-black">
						{{ (storeInfo.payment_type || 'qris_static') === 'midtrans' ? 'AUTOMATIC MIDTRANS GATEWAY' : 'MANUAL STATIC QRIS' }}
					</span>
				</div>
			</div>

			<button @click="emit('checkout')" :disabled="cart.length === 0 || isSubmitting" class="w-full bg-slate-900 hover:bg-slate-800 text-white py-3.5 rounded-xl font-black text-xs uppercase tracking-[0.2em] transition-all active:scale-95 disabled:opacity-40 shadow-md flex items-center justify-center gap-2">
				<template v-if="isSubmitting">
					<svg class="animate-spin w-4 h-4 text-white" fill="none" viewBox="0 0 24 24">
						<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
						<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
					</svg>
					<span>Memproses...</span>
				</template>
				<template v-else-if="paymentMethod === 'QRIS' && (storeInfo.payment_type || 'qris_static') === 'midtrans'">
					<svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24">
						<rect width="20" height="14" x="2" y="5" rx="2" />
						<path stroke-linecap="round" stroke-linejoin="round" d="M2 10h20M7 15h.01M11 15h2" />
					</svg>
					<span>Bayar via Midtrans Snap</span>
				</template>
				<template v-else>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
					<span>Proses & Kirim Nota WA</span>
				</template>
			</button>
		</div>
	</div>
</template>
