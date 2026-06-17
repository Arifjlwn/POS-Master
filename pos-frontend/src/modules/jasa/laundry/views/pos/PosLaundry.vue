<script setup>
import { computed } from 'vue';
import SidebarLaundry from '../../components/SidebarLaundry.vue';
import PosCameraModal from '../../components/pos/PosCameraModal.vue';
import PosCartSidebar from '../../components/pos/PosCartSidebar.vue';
import PosPerfumeModal from '../../components/pos/PosPerfumeModal.vue';
import PosQrisModal from '../../components/pos/PosQrisModal.vue';
import PosServiceCard from '../../components/pos/PosServiceCard.vue';
import ReceiptModal from '../../components/pos/ReceiptModal.vue'; // 🚀 IMPORT MODAL STRUK BARU BRAY
import { usePosLaundry } from '../../composables/usePosLaundry.js';

const posEngine = usePosLaundry();
const cart = computed(() => posEngine.cart.value || []);
const isPrintModalOpen = computed(() => !!posEngine.printData.value); // Buka otomatis kalau data struk keisi bray
</script>

<template>
	<SidebarLaundry class="hide-on-print">
		<div class="flex h-screen bg-[#F8FAFC] font-sans overflow-hidden w-full">
			<div class="flex-1 flex flex-col h-full overflow-hidden relative w-full lg:w-auto">
				<div class="p-4 md:p-6 shrink-0 bg-white shadow-sm z-10 flex flex-col gap-4 border-b border-slate-200">
					<div class="bg-slate-900 rounded-2xl p-4 md:p-5 flex flex-col sm:flex-row justify-between sm:items-center text-white shadow-lg gap-4">
						<div class="flex items-center gap-3">
							<div class="bg-white/10 p-2 rounded-xl border border-white/5">
								<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" /></svg>
							</div>
							<div>
								<h1 class="text-sm md:text-base font-black tracking-wider uppercase leading-none">POS Kasir Laundry</h1>
								<span class="text-[8px] font-bold bg-white/10 px-2 py-0.5 rounded border border-white/5 uppercase mt-1.5 inline-block">Sistem Kasir Pintar</span>
							</div>
						</div>
						<button @click="posEngine.showPerfumeControlModal.value = true" class="bg-white/10 hover:bg-white/20 px-4 py-2 rounded-xl border border-white/10 text-[10px] font-black uppercase tracking-widest transition-all flex items-center justify-center gap-2">
							<svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
								<circle cx="12" cy="12" r="10" />
								<path d="M12 8a4 4 0 1 0 0 8 4 4 0 0 0 0-8z" />
							</svg>
							Saklar Stok Parfum
						</button>
					</div>
					<div class="relative w-full group">
						<div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
							<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-400 group-focus-within:text-slate-800 transition-colors" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
								<circle cx="11" cy="11" r="8" />
								<line x1="21" y1="21" x2="16.65" y2="16.65" />
							</svg>
						</div>
						<input v-model="posEngine.searchQuery.value" type="text" placeholder="Cari paket layanan cuci..." class="w-full pl-11 pr-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-slate-800 outline-none font-bold text-xs text-slate-800 transition-all placeholder:text-slate-400 shadow-inner" />
					</div>
				</div>

				<div class="flex-1 overflow-y-auto custom-scrollbar p-4 md:p-6 pb-28 lg:pb-6">
					<div v-if="posEngine.isLoading.value" class="flex flex-col items-center justify-center h-40 text-slate-400">
						<div class="w-7 h-7 border-4 border-slate-200 border-t-slate-800 rounded-full animate-spin mb-3"></div>
						<p class="font-black text-[9px] uppercase tracking-widest animate-pulse">Memuat Layanan...</p>
					</div>
					<div v-else class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-4 md:gap-5">
						<PosServiceCard v-for="service in posEngine.filteredServices.value" :key="service.id" :service="service" :formatRupiah="posEngine.formatRupiah" @click="posEngine.addToCart(service)" />
					</div>
				</div>
			</div>

			<div v-if="!posEngine.isCartOpen.value && cart.length > 0" class="fixed bottom-5 left-4 right-4 z-40 lg:hidden">
				<button @click="posEngine.isCartOpen.value = true" class="w-full bg-slate-900 text-white rounded-xl p-4 flex justify-between items-center shadow-xl active:scale-95 transition-transform border border-slate-800">
					<div class="flex items-center gap-3">
						<div class="bg-white text-slate-900 font-black text-xs w-6 h-6 rounded-full flex items-center justify-center">{{ cart.length }}</div>
						<span class="text-[10px] font-black tracking-widest uppercase">Lihat Nota Timbangan</span>
					</div>
					<span class="text-sm font-black text-white">{{ posEngine.formatRupiah(posEngine.totalTagihan.value) }}</span>
				</button>
			</div>

			<PosCartSidebar
				v-model:isCartOpen="posEngine.isCartOpen.value"
				v-model:customerName="posEngine.customerName.value"
				v-model:customerPhone="posEngine.customerPhone.value"
				v-model:estimasiSelesai="posEngine.estimasiSelesai.value"
				:cart="cart"
				:customerResults="posEngine.customerResults.value"
				:showCustomerDropdown="posEngine.showCustomerDropdown.value"
				:availablePerfumes="posEngine.availablePerfumes.value"
				:paymentMethod="posEngine.paymentMethod.value"
				:mainPaymentGroup="posEngine.mainPaymentGroup.value"
				:totalTagihan="posEngine.totalTagihan.value"
				:formattedUangBayar="posEngine.formattedUangBayar.value"
				:kembalian="posEngine.kembalian.value"
				:isSubmitting="posEngine.isSubmitting.value"
				:storeInfo="posEngine.storeInfo.value"
				:photoData="posEngine.photoData.value"
				:formatRupiah="posEngine.formatRupiah"
				@search-customer="posEngine.searchCustomer"
				@select-customer="posEngine.selectCustomer"
				@close-dropdown="posEngine.closeCustomerDropdown"
				@clear-cart="posEngine.clearCart"
				@remove-item="posEngine.removeCartItem"
				@perfume-change="posEngine.handleCartPerfumeChange"
				@update-berat="posEngine.updateBerat"
				@open-camera="posEngine.openCamera"
				@remove-photo="posEngine.photoData.value = null"
				@set-payment="posEngine.setPaymentMethod"
				@set-nominal="posEngine.setNominalCash"
				@checkout="posEngine.processCheckout" />
		</div>

		<PosCameraModal :isCameraOpen="posEngine.isCameraOpen.value" :cameraTarget="posEngine.cameraTarget.value" @close="posEngine.closeCamera" @jepret="posEngine.takePhoto" @video-bind="posEngine.bindVideoStreaming" />
		<canvas
			:ref="
				(el) => {
					if (el) posEngine.canvasItemRef.value = el;
				}
			"
			class="hidden"></canvas>
		<canvas
			:ref="
				(el) => {
					if (el) posEngine.canvasQrisRef.value = el;
				}
			"
			class="hidden"></canvas>

		<PosPerfumeModal :showModal="posEngine.showPerfumeControlModal.value" :perfumes="posEngine.perfumes.value" :formatRupiah="posEngine.formatRupiah" @close="posEngine.showPerfumeControlModal.value = false" @toggle-status="posEngine.togglePerfumeStatus" />
		<PosQrisModal :showModal="posEngine.showQrisModal.value" :totalTagihan="posEngine.totalTagihan.value" :qrisStoreUrl="posEngine.qrisStoreUrl.value" :buktiTransferData="posEngine.buktiTransferData.value" :formatRupiah="posEngine.formatRupiah" @close="posEngine.cancelQris" @open-camera="posEngine.openCamera('QRIS')" @remove-photo="posEngine.buktiTransferData.value = null" @confirm="posEngine.confirmQris" />

		<ReceiptModal :show="isPrintModalOpen" :invoiceData="posEngine.printData.value" :storeData="posEngine.storeInfo.value" @close="posEngine.printData.value = null" />
	</SidebarLaundry>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
	height: 4px;
	width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
	background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
	background: #475569;
	border-radius: 10px;
}
input[type='number']::-webkit-inner-spin-button,
input[type='number']::-webkit-outer-spin-button {
	-webkit-appearance: none;
	margin: 0;
}
input[type='number'] {
	-moz-appearance: textfield;
}
</style>
