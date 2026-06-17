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

// 🚀 FIXED 08 TELEPON ENGINE
const formatPhoneTo08 = (phone) => {
	if (!phone) return '';
	let raw = String(phone).replace(/\D/g, '');
	if (raw.startsWith('62')) {
		raw = '0' + raw.substring(2);
	} else if (!raw.startsWith('0')) {
		raw = '0' + raw;
	}
	return raw;
};

// 🚀 FIXED DATE TIME ENGINE STANDARD KASIR (Hari, DD/MM/YYYY HH:MM)
const parseFormatTanggal = (rawDate) => {
	if (!rawDate) return '--:--';
	try {
		let str = String(rawDate).trim();
		// Bersihkan pemisah T bawaan database Go
		str = str.replace('T', ' ').replace(',', ' ');

		// Atasi string format tanggal picker murni
		if (str.length === 10 && str.includes('-')) {
			str += 'T00:00:00Z';
		}

		const d = new Date(str);
		if (isNaN(d.getTime())) {
			const now = new Date();
			const namaHari = ['Minggu', 'Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu'];
			return `${namaHari[now.getDay()]}, ${String(now.getDate()).padStart(2, '0')}/${String(now.getMonth() + 1).padStart(2, '0')}/${now.getFullYear()} ${String(now.getHours()).padStart(2, '0')}:${String(now.getMinutes()).padStart(2, '0')}`;
		}

		const namaHari = ['Minggu', 'Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu'];
		const hari = namaHari[d.getDay()];
		const tanggal = String(d.getDate()).padStart(2, '0');
		const bulan = String(d.getMonth() + 1).padStart(2, '0');
		const tahun = d.getFullYear();
		const jam = String(d.getHours()).padStart(2, '0');
		const menit = String(d.getMinutes()).padStart(2, '0');

		return `${hari}, ${tanggal}/${bulan}/${tahun}`;
	} catch (e) {
		return String(rawDate);
	}
};

const waktuDisplay = computed(() => {
	if (!props.invoiceData) return '--:--';
	return parseFormatTanggal(props.invoiceData.tanggal);
});

const estimasiDisplay = computed(() => {
	if (!props.invoiceData?.estimasi) return null;
	return parseFormatTanggal(props.invoiceData.estimasi);
});

// SUB-TOTAL DISPLAY ENGINE
const subTotalDisplay = computed(() => {
	const items = props.invoiceData?.items || [];
	if (items.length > 0) {
		return items.reduce((sum, item) => {
			const harga = item.harga || 0;
			const berat = item.berat || 0;
			const chargeParfum = item.harga_parfum || 0;
			return sum + Number(harga) * Number(berat) + Number(chargeParfum);
		}, 0);
	}
	return 0;
});

const grandTotal = computed(() => {
	if (props.invoiceData?.total !== undefined) return Number(props.invoiceData.total);
	return subTotalDisplay.value;
});

const kembalianDisplay = computed(() => {
	if (!props.invoiceData) return 0;
	return Number(props.invoiceData.kembali || 0);
});

const formatKemasan = (item) => {
	return `${item.berat || 0} ${item.satuan_dasar || 'KG'}`;
};

const triggerPrint = () => {
	window.print();
};
</script>

<template>
	<div v-if="show && invoiceData" class="fixed inset-0 bg-slate-950/90 flex items-center justify-center z-[200] p-4 backdrop-blur-sm print:absolute print:inset-0 print:bg-white print:p-0 print:block print:z-[9999] print-modal-overlay">
		<div class="bg-white p-6 md:p-8 rounded-[32px] md:rounded-[40px] shadow-2xl w-full max-w-sm flex flex-col max-h-[90vh] border-[6px] md:border-[8px] border-slate-800 print:border-none print:shadow-none print:max-h-none print:max-w-none print:rounded-none print:m-0 print:p-0 print-modal-content">
			<div class="overflow-y-auto custom-scrollbar bg-white p-2 mx-auto print:overflow-visible print:mx-0 print:px-0 print:py-0 select-none" id="print-area">
				<div class="text-center mb-3 font-mono leading-none">
					<div v-if="invoiceData?.logo_url && invoiceData.logo_url !== ''" class="mb-1.5 mt-0">
						<img :src="invoiceData.logo_url.startsWith('http') ? invoiceData.logo_url : API_BASE_URL + invoiceData.logo_url" class="w-12 h-12 object-contain mx-auto grayscale contrast-200 brightness-90" alt="Logo Toko" />
					</div>
					<div v-else>
						<h2 class="font-black text-sm md:text-base uppercase tracking-tighter mb-1 italic text-black">{{ invoiceData?.toko_nama || 'LAUNDRY POS' }}</h2>
					</div>

					<p class="text-[9px] font-black uppercase tracking-tight text-black leading-tight px-1 whitespace-pre-line mt-1">
						{{ invoiceData?.toko_alamat || 'JAKARTA, INDONESIA' }}
						<span v-if="invoiceData?.toko_telepon">
							<br />
							TELP: {{ formatPhoneTo08(invoiceData.toko_telepon) }}
						</span>
					</p>
				</div>

				<div class="border-y-2 border-black border-dashed py-1 text-center font-black mb-2 font-mono text-xs tracking-[0.1em] uppercase text-black">STRUK NOTA TIMBANG</div>

				<div class="mb-3 text-[10px] font-black font-mono uppercase space-y-1 text-black leading-tight">
					<div class="flex justify-between items-start gap-1">
						<span class="shrink-0">WAKTU DATANG:</span>
						<span class="text-right">{{ waktuDisplay }}</span>
					</div>
					<div class="flex justify-between items-start gap-1">
						<span class="shrink-0">KASIR:</span>
						<span class="text-right">{{ cashierName || 'Admin' }} / POS-01</span>
					</div>
					<div class="flex justify-between items-start gap-1">
						<span class="shrink-0">PELANGGAN:</span>
						<span class="text-right">{{ invoiceData.pelanggan || 'UMUM' }}</span>
					</div>
					<div v-if="invoiceData.pelanggan_phone" class="flex justify-between items-start gap-1">
						<span class="shrink-0">NO HP PELANGGAN:</span>
						<span class="text-right">{{ formatPhoneTo08(invoiceData.pelanggan_phone) }}</span>
					</div>
					<div class="flex justify-between items-start gap-1">
						<span class="shrink-0">NO. RESI:</span>
						<span class="font-black text-right tracking-tighter text-[9px] break-all max-w-[55%] ml-2">{{ invoiceData.invoice || 'INV-TEMP' }}</span>
					</div>
					<div v-if="estimasiDisplay" class="flex justify-between items-start gap-1">
						<span class="shrink-0">EST SELESAI:</span>
						<span class="text-right">{{ estimasiDisplay }}</span>
					</div>
				</div>

				<div class="border-b-2 border-black border-dashed mb-2"></div>

				<div v-for="(item, index) in invoiceData.items" :key="index" class="mb-2 font-mono leading-tight uppercase break-inside-avoid text-black">
					<div class="w-full font-black text-[11px] mb-1 break-words">{{ item.nama_produk || 'JASA LAUNDRY' }}</div>
					<div class="flex justify-between pl-2 text-[10px] font-black">
						<span>
							{{ formatKemasan(item) }}
							<span class="lowercase">x {{ formatRupiah(item.harga) }}</span>
						</span>
						<span class="text-[11px]">{{ formatRupiah(item.harga * item.berat + item.harga_parfum) }}</span>
					</div>
					<div class="text-[9px] pl-2 mt-1 font-black text-black">* AROMA: {{ item.nama_parfum || 'PARFUM STANDAR TOKO' }}</div>
				</div>

				<div class="border-t-2 border-black border-dashed mt-2 pt-1.5"></div>

				<div class="flex justify-between font-black text-[10px] mb-1 font-mono uppercase text-black">
					<span>SUBTOTAL:</span>
					<span>{{ formatRupiah(subTotalDisplay) }}</span>
				</div>
				<div class="flex justify-between font-black text-xs mb-1.5 font-mono uppercase border-t border-black pt-1 mt-1 text-black">
					<span>TOTAL TAGIHAN:</span>
					<span>{{ formatRupiah(grandTotal) }}</span>
				</div>
				<div class="border-b-2 border-black border-dashed mb-2"></div>

				<div class="flex justify-between mb-1 font-black font-mono text-[10px] uppercase text-black">
					<span>BAYAR ({{ invoiceData.metode || 'TUNAI' }}):</span>
					<span>{{ formatRupiah(invoiceData.bayar) }}</span>
				</div>
				<div class="flex justify-between font-black font-mono text-[10px] uppercase text-black">
					<span>KEMBALI:</span>
					<span>{{ kembalianDisplay > 0 ? formatRupiah(kembalianDisplay) : '0' }}</span>
				</div>

				<div class="text-center mt-5 font-black font-mono text-[10px] border-2 border-black p-2 uppercase leading-tight whitespace-pre-line text-black">
					{{ invoiceData?.toko_footer || 'TERIMA KASIH ATAS KUNJUNGAN ANDA!' }}
				</div>

				<div class="hidden print:block font-mono text-xs leading-none select-none text-white">
					<br />
					&nbsp; .
				</div>
			</div>

			<div class="mt-4 md:mt-6 flex flex-col gap-2 md:gap-3 shrink-0 print:hidden">
				<button @click="triggerPrint" class="w-full bg-slate-900 hover:bg-slate-800 text-white py-3 md:py-4 rounded-xl md:rounded-2xl font-black text-xs uppercase tracking-[0.2em] flex items-center justify-center gap-2 active:scale-95 transition-all">Cetak Struk</button>
				<button @click="emit('close')" class="w-full bg-slate-100 text-slate-500 py-3 md:py-3.5 rounded-xl md:rounded-2xl font-black text-xs uppercase tracking-widest hover:bg-slate-200 transition-all">Tutup</button>
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
</style>

<style>
@media print {
	@page {
		margin: 0mm !important;
		size: v-bind("storeData?.printer_width || storeData?.PrinterWidth || '58mm'") auto !important;
	}

	body * {
		visibility: hidden !important;
	}

	#print-area,
	#print-area * {
		visibility: visible !important;
		color: #000000 !important;
		font-family: 'Courier New', Courier, monospace !important;
		font-weight: 900 !important;
	}

	#print-area {
		position: absolute !important;
		left: 0 !important;
		top: 0 !important;
		width: v-bind("storeData?.printer_width || storeData?.PrinterWidth || '58mm'") !important;

		/* 🚀 FIX NENDANG & KIRI-KANAN BERNAFAS: Atas 0mm (mepet maksimal), Kiri-Kanan 4mm (gak kepotong paper) */
		padding: 0mm 4mm 5mm 4mm !important;
		margin: 0 !important;
		background: #ffffff !important;
		box-sizing: border-box !important;
	}

	.print\:hidden,
	button {
		display: none !important;
		visibility: hidden !important;
	}

	html,
	body {
		background: #ffffff !important;
		height: auto !important;
		overflow: visible !important;
	}
}
</style>
