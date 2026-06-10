<script setup>
import { computed } from 'vue';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

const props = defineProps({
	show: Boolean,
	invoiceData: Object,
	storeData: Object,
	cashierName: String,
	stationNumber: String,
});

const emit = defineEmits(['close']);

const formatRupiah = (angka) => {
	if (angka === 0) return '0';
	if (typeof angka !== 'number' || isNaN(angka)) return '0';
	return new Intl.NumberFormat('id-ID').format(angka);
};

// 🚀 1. SUB-TOTAL DISPLAY ENGINE: Kebal dari jebakan null/undefined data
const subTotalDisplay = computed(() => {
	if (props.invoiceData?.subtotal !== undefined) return Number(props.invoiceData.subtotal);
	if (props.invoiceData?.sub_total !== undefined) return Number(props.invoiceData.sub_total);

	const items = props.invoiceData?.cart || props.invoiceData?.details || [];
	if (items.length > 0) {
		return items.reduce((sum, item) => {
			const harga = item.price || item.harga_satuan || 0;
			const qty = item.qty || item.kuantitas || 0;
			return sum + Number(harga) * Number(qty);
		}, 0);
	}
	return 0;
});

// 🚀 2. TAX CALCULATOR: Akurat membaca status pajak SaaS merchant side
const totalPajak = computed(() => {
	if (props.invoiceData?.pajak !== undefined) return Number(props.invoiceData.pajak);
	const isTaxActive = props.storeData?.is_tax_active || props.storeData?.IsTaxActive || false;
	if (!isTaxActive) return 0;
	const persen = props.storeData?.pajak_persen || props.storeData?.pajak_percent || 0;
	return subTotalDisplay.value * (Number(persen) / 100);
});

// 🚀 3. GRAND TOTAL: Hasil final klop antara server database vs client render
const grandTotal = computed(() => {
	if (props.invoiceData?.total !== undefined) return Number(props.invoiceData.total);
	if (props.invoiceData?.tagihan !== undefined) return Number(props.invoiceData.tagihan);
	if (props.invoiceData?.total_harga !== undefined) return Number(props.invoiceData.total_harga);
	return subTotalDisplay.value + totalPajak.value;
});

// 🚀 4. NOMINAL KEMBALIAN KAS: Anti-Bug Falsy Value (Uang Pas tercetak sempurna )
const kembalianDisplay = computed(() => {
	if (!props.invoiceData) return 0;

	// Tampung semua kemungkinan variasi nama field dari backend Go lu
	const values = [props.invoiceData.return, props.invoiceData['return'], props.invoiceData.kembalian, props.invoiceData.kembali, props.invoiceData.counter];

	for (let val of values) {
		if (val !== undefined && val !== null) {
			return Number(val);
		}
	}

	// Fallback hitungan manual jika backend ga ngirim data kembalian
	const nominalBayar = props.invoiceData.pay ?? props.invoiceData.nominal_bayar ?? 0;
	const hasil = Number(nominalBayar) - grandTotal.value;
	return hasil > 0 ? hasil : 0;
});

// 🚀 WAKTU DISPLAY ENGINE: Sniper Regex (Garansi Aman di POS & Riwayat)
const waktuDisplay = computed(() => {
	if (!props.invoiceData) return '--:--';

	let rawDate = props.invoiceData.created_at || props.invoiceData.date || props.invoiceData.date_time;

	if (!rawDate) return '--:--';

	// Ubah ke string biar seragam cara nge-ceknya
	const strDate = String(rawDate);

	// 🎯 1. KASUS KHUSUS REPRINT RIWAYAT (Database Go UTC)
	// Tandanya: ada field 'created_at', atau datanya ditarik dari backend harian
	if (props.invoiceData.created_at) {
		try {
			let tempDate = strDate;
			// Jika string polos dari DB belum ada penanda UTC, paksa inject 'Z' biar JS tahu itu UTC
			if (!tempDate.endsWith('Z') && !tempDate.includes('+')) {
				tempDate = tempDate.replace(' ', 'T') + 'Z';
			}
			const d = new Date(tempDate);
			if (!isNaN(d.getTime())) {
				// Tambah 7 jam otomatis ke WIB
				return new Intl.DateTimeFormat('id-ID', {
					hour: '2-digit',
					minute: '2-digit',
					hour12: false,
					timeZone: 'Asia/Jakarta',
				})
					.format(d)
					.replace('.', ':');
			}
		} catch (e) {
			console.error(e);
		}
	}

	// 🎯 2. KASUS LIVE POS KASIR (Waktu Lokal Akurat)
	// Kita langsung nyomot pola angka "HH:MM" dari string tanpa lewat proses "new Date()" yang bikin jam bergeser/bocor T-Z
	const match = strDate.match(/(\d{2}):(\d{2})/);
	if (match && match[0]) {
		return match[0]; // Langsung keluarin angka jam mentahnya, misal "18:49"
	}

	return strDate; // Fallback terakhir tampilkan apa adanya
});

// FUNGSI SILUMAN PEMBACA KEMASAN UOM MULTI-UNIT
const formatKemasan = (item) => {
	if (item.detail_notes && item.detail_notes !== 'Transaksi Retail Toko') return item.detail_notes;
	if (item.uom_label) return item.uom_label; // Jika backend udah ngirim label jadi, pakai langsung!

	const satuanPilihan = item.selected_uom || item.satuan_terpilih || item.satuan || item.kemasan;
	const qtyPilihan = item.qty || item.kuantitas || 0;
	if (satuanPilihan) return `${qtyPilihan} ${satuanPilihan}`;

	const product = item.product || item;
	return `${qtyPilihan} ${product.satuan_dasar || 'PCS'}`;
};

const triggerPrint = () => {
	window.print();
};
</script>

<template>
	<div v-if="show && invoiceData" class="fixed inset-0 bg-slate-950/90 flex items-center justify-center z-[200] p-4 backdrop-blur-sm print:static print:bg-white print:p-0 print:block">
		<div class="bg-white p-6 md:p-8 rounded-[32px] md:rounded-[40px] shadow-2xl w-full max-w-sm flex flex-col max-h-[90vh] border-[6px] md:border-[8px] border-slate-800 print:border-none print:shadow-none print:max-h-none print:max-w-none print:rounded-none print:m-0 print:p-0">
			<div class="overflow-y-auto custom-scrollbar bg-white p-2 mx-auto print:overflow-visible print:mx-0 print:px-3 print:py-2 select-none" id="print-area">
				<div class="text-center mb-4 font-mono leading-none">
					<div v-if="storeData?.logo_url && storeData.logo_url !== ''">
						<img :src="storeData.logo_url.startsWith('http://') || storeData.logo_url.startsWith('https://') ? storeData.logo_url : API_BASE_URL + storeData.logo_url" class="w-16 h-16 object-contain mx-auto grayscale contrast-200 brightness-90" alt="Logo Toko" />
					</div>

					<h2 v-else class="font-black text-xs md:text-sm uppercase tracking-tighter mb-1 italic">{{ storeData?.nama_toko || storeData?.NamaToko || 'ARZURA POS STORE' }}</h2>
					<p class="text-[9px] font-black uppercase tracking-tight opacity-100 leading-tight px-1">
						{{ storeData?.alamat || 'JAKARTA, INDONESIA' }}
						<br />
						{{ storeData?.kelurahan || 'KELURAHAN' }}, {{ storeData?.kecamatan || 'KECAMATAN' }}
						<br />
						{{ storeData?.kota || 'KOTA' }}, {{ storeData?.provinsi || 'PROVINSI' }} {{ storeData?.kode_pos || 'KODE POS' }}
					</p>
				</div>

				<div class="border-y border-black py-1 text-center font-black mb-3 font-mono text-[9px] tracking-[0.2em] uppercase">
					{{ invoiceData.created_at ? 'Invoice Reprint' : 'Struk Belanja' }}
				</div>

				<div class="mb-3 text-[8px] font-black font-mono uppercase space-y-0.5">
					<div class="flex justify-between">
						<span>WAKTU:</span>
						<span>{{ waktuDisplay }}</span>
					</div>

					<div class="flex justify-between">
						<span>KASIR:</span>
						<span>{{ cashierName || 'KASIR' }} / POS-{{ stationNumber || '01' }}</span>
					</div>
					<div class="flex justify-between">
						<span class="truncate pr-2">INV:</span>
						<span class="font-bold shrink-0">{{ invoiceData.invoice || invoiceData.no_invoice || 'INV-TEMP' }}</span>
					</div>
				</div>

				<div class="border-b border-black border-dashed mb-2"></div>

				<div v-for="item in invoiceData.cart || invoiceData.details" :key="item.id" class="mb-2 font-bold font-mono text-[9px] leading-tight uppercase break-inside-avoid">
					<div class="truncate w-full pr-1">{{ item.name || item.product_name || item.product?.nama_produk || 'Item Belanja' }}</div>
					<div class="flex justify-between pl-2 text-[8px] mt-0.5">
						<span>
							{{ formatKemasan(item) }}
							<span class="lowercase">x {{ formatRupiah(item.price || item.harga_uom || item.harga_satuan) }}</span>
						</span>
						<span class="font-black text-[9px]">{{ formatRupiah(item.sub_total || item.subtotal_item || item.price * item.qty) }}</span>
					</div>
				</div>

				<div class="border-t border-black border-dashed mt-2 pt-1.5"></div>
				<div class="flex justify-between font-bold text-[9px] mb-0.5 font-mono uppercase italic">
					<span>SUBTOTAL:</span>
					<span>{{ formatRupiah(subTotalDisplay) }}</span>
				</div>
				<div v-if="totalPajak >= 0" class="flex justify-between font-bold text-[9px] mb-0.5 font-mono uppercase italic">
					<span>PAJAK:</span>
					<span>{{ formatRupiah(totalPajak) }} %</span>
				</div>
				<div class="flex justify-between font-black text-[10px] mb-1.5 font-mono uppercase italic border-t border-black pt-1 mt-1">
					<span>TOTAL BELANJA:</span>
					<span>{{ formatRupiah(grandTotal) }}</span>
				</div>
				<div class="border-b border-black border-dashed mb-2"></div>

				<div class="flex justify-between mb-0.5 font-bold font-mono text-[9px] uppercase">
					<span>BAYAR ({{ invoiceData.method || invoiceData.metode_bayar || 'CASH' }}):</span>
					<span>{{ formatRupiah(invoiceData.pay !== undefined ? invoiceData.pay : (invoiceData.nominal_bayar ?? grandTotal)) }}</span>
				</div>

				<div class="flex justify-between font-black font-mono text-[9px] uppercase italic text-black">
					<span>KEMBALI:</span>
					<span>{{ kembalianDisplay > 0 ? formatRupiah(kembalianDisplay) : '0' }}</span>
				</div>

				<div class="text-center mt-4 font-black font-mono text-[9px] border border-black p-1.5 uppercase leading-tight whitespace-pre-line">
					{{ storeData?.receipt_footer || storeData?.ReceiptFooter || 'TERIMA KASIH ATAS KUNJUNGAN ANDA!' }}
				</div>

				<div class="hidden print:block font-mono text-[9px] leading-none select-none text-white">
					<br />
					&nbsp; .
				</div>
			</div>

			<div class="mt-4 md:mt-6 flex flex-col gap-2 md:gap-3 shrink-0 print:hidden">
				<button @click="triggerPrint" class="w-full bg-indigo-600 text-white py-3 md:py-4 rounded-xl md:rounded-2xl font-black text-[9px] md:text-[10px] uppercase tracking-[0.2em] shadow-xl shadow-indigo-200 flex items-center justify-center gap-2 active:scale-95 transition-all">Cetak Struk</button>
				<button @click="emit('close')" class="w-full bg-slate-100 text-slate-500 py-3 md:py-3.5 rounded-xl md:rounded-2xl font-black text-[9px] md:text-[10px] uppercase tracking-widest hover:bg-slate-200 transition-all">Tutup</button>
			</div>
		</div>
	</div>
</template>

<style scoped>
#print-area {
	width: v-bind("storeData?.printer_width || storeData?.PrinterWidth || '58mm'") !important;
	max-width: 100% !important;
}

.custom-scrollbar::-webkit-scrollbar {
	width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
	background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
	background: #cbd5e1;
	border-radius: 10px;
}

/* ==========================================
   🖨️ PREMIUM MULTI-TENANT DYNAMIC PRINT ENGINE 
   ========================================== */
@media print {
	@page {
		/* 🚀 FORCE ATAS DEMPET: Bersihkan margin halaman virtual browser tanpa toleransi */
		margin: 0mm 0mm 0mm 0mm !important;
		size: v-bind("storeData?.printer_width || storeData?.PrinterWidth || '58mm'") auto !important;
	}

	html,
	body {
		background: #ffffff !important;
		margin: 0 !important;
		padding: 0 !important;
		height: auto !important;
		max-height: none !important;
		box-sizing: border-box !important;
	}

	#print-area {
		width: v-bind("storeData?.printer_width || storeData?.PrinterWidth || '58mm'") !important;
		margin: 0 auto !important;

		/* 🚀 FIX CRITICAL MARGIN ATAS: 
           - Padding atas dikunci mati ke 0mm biar konten langsung nempel ke bibir printer !
           - Padding kanan-kiri 4mm (safe zone scale 100%)
           - Padding bawah 5mm (penahan dorongan titik putih html) */
		padding: 0mm 4mm 5mm 4mm !important;

		box-sizing: border-box !important;
		background: #ffffff !important;
		display: block !important;
	}

	#print-area * {
		box-sizing: border-box !important;
	}
}
</style>
