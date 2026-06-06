<script setup>
import Swal from 'sweetalert2';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import api from '../../../../api.js';
import { usePos } from '../../composables/usePos.js';

// PASUKAN SUB-KOMPONEN EMANG JOSS!
import CartSidebar from '../../components/pos/CartSidebar.vue';
import ClosingModal from '../../components/pos/ClosingModal.vue';
import PosHeader from '../../components/pos/PosHeader.vue';
import ProductCatalog from '../../components/pos/ProductCatalog.vue';
import ReceiptModal from '../../components/pos/ReceiptModal.vue';
import WaPromptModal from '../../components/pos/WaPromptModal.vue';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';
const router = useRouter();
const storeData = ref({});
const showWaModal = ref(false);

// 🚀 SATPAM KHUSUS POS & LOAD DATA TOKO
onMounted(async () => {
	try {
		const res = await api.get('/retail/store/settings');
		storeData.value = res.data.data;
		localStorage.setItem('storeLogo', storeData.value.logo_url || '');

		// SUNTIK SCRIPT MIDTRANS DINAMIS (BACA CONFIG BACKEND DATABASE)
		if (storeData.value.payment_type === 'midtrans' && storeData.value.midtrans_client_key) {
			if (!document.getElementById('midtrans-script')) {
				const midtransEnv = import.meta.env.VITE_MIDTRANS_ENV || 'sandbox';
				const snapUrl = midtransEnv === 'production' ? 'https://app.midtrans.com/snap/snap.js' : 'https://app.sandbox.midtrans.com/snap/snap.js';

				const script = document.createElement('script');
				script.id = 'midtrans-script';
				script.src = snapUrl;
				script.setAttribute('data-client-key', storeData.value.midtrans_client_key);
				document.head.appendChild(script);
			}
		}

		// --- BLOK SATPAM GLOBAL (ANTI-KABUR TOKO KADALUWARSA) ---
		const role = localStorage.getItem('role') || 'owner';
		if (role === 'owner') {
			let isDead = false;
			if (storeData.value.subscription_status !== 'active') {
				isDead = true;
			} else if (storeData.value.subscription_end) {
				const endDate = new Date(storeData.value.subscription_end),
					today = new Date();
				const diffDays = Math.ceil((endDate.getTime() - today.getTime()) / (1000 * 60 * 60 * 24));
				if (diffDays <= 0) isDead = true;
			}
			if (isDead) {
				Swal.fire({ icon: 'error', title: 'Masa Aktif Habis!', text: 'Silakan selesaikan tagihan langganan SaaS Anda.', confirmButtonColor: '#ef4444' });
				router.push('/retail/account');
				return;
			}
		}
	} catch (e) {
		console.error('Gagal narik data toko / eksekusi satpam POS:', e);
	}
});

const {
	currentUser,
	currentSession,
	currentTime,
	products,
	isLoadingProducts,
	cart,
	heldOrders,
	showHeldModal,
	payAmount,
	paymentMethod,
	showReceipt,
	showQrisModal,
	lastTransaction,
	showReceiptClosing,
	lastClosingData,
	isMobileCartOpen,
	searchQuery,
	searchInput,
	showScanner,
	pecahan,
	filteredProducts,
	subTotalBelanja,
	nilaiPajak,
	totalBelanja,
	kembalian,
	totalUangFisik,
	showClosingModal,
	isProcessingCheckout,
	storeLogo,
	getImageUrl,
	startScanner,
	stopScanner,
	handleBarcodeScan,
	addToCart,
	toggleUom,
	decreaseQty,
	increaseQty,
	validateQty,
	clearCart,
	holdTransaction,
	resumeOrder,
	setPaymentMethod,
	executeCheckout,
	formatInputRupiah,
	processCheckout,
	handleClosing,
	logout,
	setNominal,
	noHpPelanggan,
} = usePos();

// INITIAL CHECKOUT ROUTER: Otak percabangan Midtrans vs QRIS Statis bray!
const handleInitialCheckout = async () => {
	if (cart.length === 0) return;

	if (paymentMethod.value === 'Cash' && payAmount.value < totalBelanja.value) {
		return Swal.fire({
			icon: 'error',
			title: 'Uang Kurang!',
			text: `Uang fisik kas kurang Rp ${(totalBelanja.value - payAmount.value).toLocaleString('id-ID')}`,
			confirmButtonColor: '#ef4444',
			customClass: { popup: 'rounded-[28px]' },
		});
	}

	if (paymentMethod.value === 'QRIS') {
		if (storeData.value?.payment_type === 'midtrans') {
			if (typeof window.snap === 'undefined') {
				return Swal.fire('Sistem Loading', 'Script Midtrans belum siap, pastikan koneksi internet stabil bray.', 'warning');
			}
			isProcessingCheckout.value = true;
			try {
				const payRes = await api.post('/retail/pos/midtrans-order', { total: totalBelanja.value });
				window.snap.pay(payRes.data.token, {
					onSuccess: () => {
						isProcessingCheckout.value = false;
						Swal.fire('Berhasil', 'Pembayaran QRIS Midtrans Diterima!', 'success');
						triggerCheckoutFlow();
					},
					onPending: () => {
						Swal.fire('Menunggu', 'Pelanggan sedang melakukan scan saldo.', 'info');
						isProcessingCheckout.value = false;
					},
					onError: () => {
						Swal.fire('Gagal', 'Pembayaran ditolak atau kadaluwarsa.', 'error');
						isProcessingCheckout.value = false;
					},
					onClose: () => {
						isProcessingCheckout.value = false;
					},
				});
			} catch (error) {
				Swal.fire('Error Gateway', 'Gagal memanggil API Midtrans dari Golang.', 'error');
				isProcessingCheckout.value = false;
			}
		} else {
			showQrisModal.value = true; // Buka Modal Gambar QRIS Manual bray
		}
	} else {
		triggerCheckoutFlow();
	}
};

// 🛡️ ANTI-FRAUD QRIS MANUAL: Paksa kasir sumpah konfirmasi mutasi masuk bray!
const verifyManualQrisPayment = () => {
	Swal.fire({
		title: 'Konfirmasi Uang Masuk',
		text: `Pastikan dana Rp ${totalBelanja.value.toLocaleString('id-ID')} SUDAH MASUK ke e-wallet toko sebelum mencetak struk bray!`,
		icon: 'warning',
		showCancelButton: true,
		confirmButtonColor: '#4f46e5',
		cancelButtonColor: '#cbd5e1',
		confirmButtonText: 'Ya, Sudah Masuk!',
		cancelButtonText: 'Cek Lagi',
		customClass: { popup: 'rounded-[28px]' },
	}).then((result) => {
		if (result.isConfirmed) {
			triggerCheckoutFlow();
		}
	});
};

// 🛡️ FIX SECURITY: Cabut validasi premium dari LocalStorage!
// Ambil status langsung dari object 'storeData' database aman backend bray!
const triggerCheckoutFlow = () => {
	showQrisModal.value = false;
	const planFromDb = storeData.value?.subscription_plan || 'basic';
	const isPremium = planFromDb === 'premium' || planFromDb === 'trial';
	const hasWaToken = storeData.value?.wa_token && storeData.value.wa_token !== '';

	if (isPremium && hasWaToken) {
		showWaModal.value = true;
	} else {
		proceedToFinalCheckout();
	}
};

const proceedToFinalCheckout = async () => {
	showWaModal.value = false;
	await executeCheckout();
};

const handleWaSubmit = (phone) => {
	noHpPelanggan.value = phone;
	proceedToFinalCheckout();
};

const handleWaSkip = () => {
	noHpPelanggan.value = '';
	proceedToFinalCheckout();
	showReceipt.value = true;
};

const printClosing = () => {
	window.print();
};
const finishClosing = () => router.push('/retail/pos/riwayat');
const goToRiwayat = () => router.push('/retail/pos/riwayat');
</script>

<template>
	<div class="h-[100dvh] flex flex-col bg-slate-100 font-sans overflow-hidden print:h-auto print:bg-white print:overflow-visible print:block">
		<PosHeader class="print:hidden" :currentUser="currentUser" :currentSession="currentSession" :currentTime="currentTime" @go-dashboard="goToRiwayat" @logout="logout" />

		<div class="print:hidden flex-1 flex overflow-hidden p-2 md:p-4 pt-2 md:pt-4 gap-4 relative">
			<ProductCatalog v-model:searchQuery="searchQuery" :filteredProducts="filteredProducts" :heldOrders="heldOrders" :getImageUrl="getImageUrl" @barcode-scan="handleBarcodeScan" @start-scanner="startScanner" @show-held="showHeldModal = true" @add-to-cart="addToCart" />

			<CartSidebar
				v-model:noHpPelanggan="noHpPelanggan"
				v-model:isMobileCartOpen="isMobileCartOpen"
				:cart="cart"
				:heldOrders="heldOrders"
				:paymentMethod="paymentMethod"
				:subTotal="subTotalBelanja"
				:pajak="nilaiPajak"
				:totalBelanja="totalBelanja"
				:payAmount="payAmount"
				:kembalian="kembalian"
				:isTaxActive="storeData?.is_tax_active"
				:pajakPersen="storeData?.pajak_percent"
				:isProcessingCheckout="isProcessingCheckout"
				@show-held="showHeldModal = true"
				@hold-order="holdTransaction"
				@clear-cart="clearCart"
				@decrease-qty="decreaseQty"
				@increase-qty="increaseQty"
				@validate-qty="validateQty"
				@set-payment="setPaymentMethod"
				@format-rupiah="formatInputRupiah"
				@set-nominal="($event === 0) ? payAmount = 0 : payAmount += $event"
				@checkout="handleInitialCheckout"
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
					<h2 class="text-lg font-black text-slate-800 uppercase tracking-widest flex items-center gap-2">Pesanan Tertunda</h2>
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
				<div class="p-6 bg-black"><div id="reader-kasir" class="w-full rounded-2xl overflow-hidden border-2 border-slate-700"></div></div>
			</div>
		</div>

		<div v-if="showQrisModal" class="fixed inset-0 bg-slate-900/80 flex items-center justify-center z-[150] p-4 backdrop-blur-sm print:hidden">
			<div class="bg-white p-6 md:p-8 rounded-[32px] shadow-2xl w-full max-w-sm text-center flex flex-col border-t-8 border-indigo-600">
				<h2 class="text-xl font-black text-slate-900 uppercase tracking-widest mb-1">Bayar via QRIS</h2>
				<div class="bg-white p-3 rounded-2xl border border-slate-300 w-full mb-4 flex flex-col justify-center items-center min-h-[200px]">
					<img v-if="storeData?.qris_image" :src="getImageUrl(storeData.qris_image)" alt="QRIS Toko" class="w-full aspect-square object-cover rounded-xl shadow-inner" />
					<div v-else class="text-slate-400 text-xs font-bold p-10 text-center aspect-square flex items-center justify-center">
						QRIS Belum Di-Setup.
						<br />
						Silakan upload di Pengaturan Toko.
					</div>
					<p v-if="storeData?.qris_name" class="mt-3 font-bold text-[10px] text-slate-500 uppercase tracking-widest border-t border-slate-100 pt-2 w-full">{{ storeData.qris_name }}</p>
				</div>
				<div class="bg-indigo-50 text-indigo-900 p-4 rounded-2xl mb-6 font-black text-2xl">Rp {{ totalBelanja.toLocaleString('id-ID') }}</div>
				<div class="flex gap-3">
					<button @click="showQrisModal = false" :disabled="isProcessingCheckout" class="flex-1 bg-slate-100 py-4 rounded-xl font-black text-[10px] uppercase text-slate-500">Batal</button>
					<button @click="verifyManualQrisPayment" :disabled="isProcessingCheckout" class="flex-1 bg-indigo-600 py-4 rounded-xl font-black text-[10px] uppercase text-white shadow-lg">{{ isProcessingCheckout ? 'Proses...' : 'Lunas' }}</button>
				</div>
			</div>
		</div>

		<WaPromptModal :show="showWaModal" :totalBelanja="totalBelanja" @submit="handleWaSubmit" @skip="handleWaSkip" @close="showWaModal = false" />
		<ReceiptModal :show="showReceipt" :invoiceData="lastTransaction" :storeData="currentSession?.store || currentSession?.Store" :cashierName="currentUser?.name ? currentUser.name.split(' ')[0] : 'KASIR'" :stationNumber="currentSession?.station_number || '01'" @close="showReceipt = false" />
		<ClosingModal :show="showClosingModal" :showReceiptClosing="showReceiptClosing" :pecahan="pecahan" :totalUangFisik="totalUangFisik" :lastClosingData="lastClosingData" :currentSession="currentSession" :currentUser="currentUser" :storeLogo="storeLogo" @close="showClosingModal = false" @process-closing="handleClosing" @print-closing="printClosing" @finish-closing="finishClosing" />
	</div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
	width: 4px;
}
@media (min-width: 768px) {
	.custom-scrollbar::-webkit-scrollbar {
		width: 6px;
	}
}
.custom-scrollbar::-webkit-scrollbar-track {
	background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
	background: #cbd5e1;
	border-radius: 10px;
}
@media print {
	@page {
		margin: 0;
	}
	body {
		background: white;
		-webkit-print-color-adjust: exact;
	}
}
:deep(.swal2-popup) {
	font-family: 'Inter', sans-serif !important;
	border-radius: 28px !important;
}
</style>
