<script setup>
import Swal from 'sweetalert2';
import { computed, ref, watch } from 'vue';
import api from '../../../../../api.js';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

const props = defineProps({
	show: Boolean,
	invoiceData: Object,
	storeData: Object,
	cashierName: String,
	stationNumber: String,
});

const emit = defineEmits(['close', 'refresh-data']);

// 🚀 STATE UNTUK PINDAH RAK
const isEditingRack = ref(false);
const allRacks = ref([]);
const availableRacks = ref([]);
const selectedRackId = ref('');
const isSubmittingRack = ref(false);

const formatRupiah = (angka) => {
	if (angka === 0) return '0';
	if (typeof angka !== 'number' || isNaN(angka)) return '0';
	return new Intl.NumberFormat('id-ID').format(angka);
};

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

const parseFormatTanggal = (rawDate) => {
	if (!rawDate) return '--:--';
	try {
		let str = String(rawDate).trim();
		str = str.replace('T', ' ').replace(',', ' ');
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

// 🚀 COMPUTED: Ambil info Rak Otomatis
const currentRack = computed(() => {
	const items = props.invoiceData?.items || [];
	if (items.length > 0 && items[0].nomor_rak) {
		return items[0].nomor_rak;
	}
	return '-';
});

// 🚀 COMPUTED BARU: Deteksi isi cucian yang ada di dalam rak aktif saat ini
const currentRackDetails = computed(() => {
	const rack = allRacks.value.find((r) => r.nama_rak === currentRack.value);
	return rack ? rack.detail_cucian : [];
});

// 🚀 ENGINE: Fetch data rak kosong pas modal terbuka
watch(
	() => props.show,
	async (newVal) => {
		if (newVal) {
			isEditingRack.value = false;

			try {
				const res = await api.get('/laundry/racks');

				allRacks.value = res.data.data || [];

				availableRacks.value = allRacks.value.filter((r) => r.status === 'TERSEDIA' && r.isi_cucian === 0);
			} catch (err) {
				console.error('Gagal memuat list master denah rak:', err);
			}
		}
	}
);

// 🚀 ENGINE: Eksekusi Grid Klik Langsung Pindah Rak
const handlePilihRakGrid = async (id) => {
	if (!id || isSubmittingRack.value) return;

	selectedRackId.value = id;

	// 🎯 KITA PAKAI NOMOR INVOICE SEKARANG BRAY, BUKAN ID PRODUK!
	const invoiceCode = props.invoiceData?.invoice;
	if (!invoiceCode) {
		return Swal.fire('Error', 'Data Nomor Invoice tidak ditemukan!', 'error');
	}

	isSubmittingRack.value = true;
	try {
		// Bypass URL param :id pakai kata 'update', lempar data asli ke dalam body JSON
		const res = await api.put(`/laundry/transactions/update/pindah-rak`, {
			new_rack_id: parseInt(id),
			invoice: invoiceCode, // 🚀 Injeksi Invoice ke Golang!
		});

		if (res.data.status === 'sukses') {
			const newRackName = availableRacks.value.find((r) => r.id === parseInt(id))?.nama_rak;

			// Ubah reaktif nama rak di struk layar kasir
			if (props.invoiceData.items && props.invoiceData.items.length > 0) {
				props.invoiceData.items.forEach((item) => {
					item.nomor_rak = newRackName;
				});
			}

			isEditingRack.value = false;
			Swal.fire({
				toast: true,
				position: 'top-end',
				icon: 'success',
				title: 'Evakuasi Rak Berhasil!',
				showConfirmButton: false,
				timer: 1500,
			});
		}
	} catch (err) {
		Swal.fire('Gagal Pindah', err.response?.data?.error || 'Koneksi terputus', 'error');
	} finally {
		isSubmittingRack.value = false;
	}
};

const triggerPrint = () => {
	window.print();
};
</script>

<template>
	<div v-if="show && invoiceData" class="fixed inset-0 bg-slate-950/90 flex items-center justify-center z-[200] p-4 backdrop-blur-sm print:absolute print:inset-0 print:bg-white print:p-0 print:block print:z-[9999] print-modal-overlay">
		<div class="absolute top-4 md:top-8 left-1/2 -translate-x-1/2 w-full max-w-sm z-20 print:hidden transition-all">
			<div class="bg-indigo-600 text-white rounded-[24px] p-4 shadow-2xl border-4 border-indigo-900 relative overflow-hidden">
				<div class="absolute -right-4 -top-4 w-24 h-24 bg-white/10 rounded-full blur-xl pointer-events-none"></div>

				<div v-if="!isEditingRack" class="relative z-10 animate-fade-in">
					<div class="flex items-center justify-between">
						<div class="flex flex-col">
							<span class="text-[10px] font-black uppercase tracking-widest text-indigo-200">Alokasi Rak Otomatis</span>
							<span class="text-3xl font-black tracking-tighter leading-none mt-1 drop-shadow-sm">{{ currentRack }}</span>
						</div>
						<button @click="isEditingRack = true" class="bg-white/20 hover:bg-white/30 text-white p-2.5 rounded-xl transition-all active:scale-95 flex items-center justify-center backdrop-blur-md shadow-sm border border-white/10" title="Ubah Rak Fisik">
							<svg xmlns="http://www.w3.org/2000/xl" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
								<path stroke-linecap="round" stroke-linejoin="round" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
							</svg>
						</button>
					</div>

					<div class="mt-4 pt-3 border-t border-indigo-500/50">
						<span class="text-[9px] font-black uppercase tracking-widest text-indigo-300 mb-2 block flex items-center gap-1.5">
							<svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
								<path stroke-linecap="round" stroke-linejoin="round" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
							</svg>
							Penghuni Rak Saat Ini:
						</span>

						<div class="max-h-24 overflow-y-auto custom-scrollbar pr-1 space-y-1.5">
							<div v-if="currentRackDetails.length > 0" v-for="(item, idx) in currentRackDetails" :key="idx" class="bg-indigo-800/60 p-2 rounded-lg border border-indigo-700/50 flex justify-between items-center shadow-inner">
								<div class="flex flex-col min-w-0 pr-2">
									<span class="text-[10px] font-bold text-indigo-50 line-clamp-1 uppercase">{{ item.product?.nama_produk || 'PAKET LAUNDRY' }}</span>
									<span class="text-[8px] font-medium text-indigo-300 line-clamp-1 mt-0.5 uppercase">PELANGGAN: {{ item.nama_pelanggan }}</span>
								</div>
								<span class="text-[9px] font-black bg-indigo-900 px-2 py-0.5 rounded text-indigo-300 shrink-0 border border-indigo-800">{{ item.berat_kg }} KG</span>
							</div>

							<div v-else class="text-center py-3 text-[10px] font-black uppercase tracking-widest text-indigo-300/60 border border-dashed border-indigo-400/40 rounded-xl bg-indigo-700/30">Tidak ada item apapun</div>
						</div>
					</div>
				</div>

				<div v-else class="flex flex-col gap-3 relative z-10 animate-fade-in w-full">
					<div class="flex justify-between items-center mb-1">
						<span class="text-[10px] font-black uppercase tracking-widest text-indigo-200">Pilih Denah Rak Kosong</span>
						<button @click="isEditingRack = false" class="text-white hover:text-rose-300 transition-colors p-1 bg-white/10 rounded-lg active:scale-95">
							<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
								<path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
							</svg>
						</button>
					</div>

					<div v-if="isSubmittingRack" class="flex flex-col items-center justify-center py-8">
						<svg class="animate-spin h-6 w-6 text-white mb-2" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
							<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
							<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
						</svg>
						<span class="text-[10px] font-black tracking-widest uppercase text-indigo-200 animate-pulse">Menyinkronkan Rak...</span>
					</div>

					<div v-else-if="availableRacks.length > 0" class="grid grid-cols-4 sm:grid-cols-5 gap-2 max-h-[220px] overflow-y-auto custom-scrollbar pr-1">
						<button v-for="rack in availableRacks" :key="rack.id" @click="handlePilihRakGrid(rack.id)" class="bg-indigo-700/50 hover:bg-white text-indigo-100 hover:text-indigo-900 border border-indigo-500 hover:border-white rounded-xl p-2 flex flex-col items-center justify-center transition-all active:scale-95 group shadow-sm">
							<span class="text-[8px] font-bold opacity-70 group-hover:opacity-100 transition-opacity">R{{ rack.baris }}C{{ rack.kolom }}</span>
							<span class="font-black text-sm mt-0.5 tracking-tighter">{{ rack.nama_rak }}</span>
						</button>
					</div>

					<div v-else class="text-center bg-white/10 border border-white/20 rounded-xl p-4 mt-2">
						<svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 mx-auto mb-2 text-indigo-300" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
							<path stroke-linecap="round" stroke-linejoin="round" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
						</svg>
						<p class="text-[10px] font-black uppercase tracking-widest text-indigo-200">Semua Rak Sedang Penuh / Rusak</p>
					</div>
				</div>
			</div>
		</div>

		<div class="bg-white p-6 md:p-8 rounded-[32px] md:rounded-[40px] shadow-2xl w-full max-w-sm flex flex-col max-h-[65vh] mt-40 border-[6px] md:border-[8px] border-slate-800 print:border-none print:shadow-none print:max-h-none print:max-w-none print:rounded-none print:m-0 print:p-0 print-modal-content relative z-10">
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
						<span class="shrink-0">WAKTU:</span>
						<span class="text-right">{{ waktuDisplay }}</span>
					</div>
					<div class="flex justify-between items-start gap-1">
						<span class="shrink-0">KASIR:</span>
						<span class="text-right">{{ cashierName || 'Admin' }}</span>
					</div>
					<div class="flex justify-between items-start gap-1">
						<span class="shrink-0">PELANGGAN:</span>
						<span class="text-right">{{ invoiceData.pelanggan || 'UMUM' }}</span>
					</div>
					<div v-if="invoiceData.pelanggan_phone" class="flex justify-between items-start gap-1">
						<span class="shrink-0">NO.HP:</span>
						<span class="text-right">{{ formatPhoneTo08(invoiceData.pelanggan_phone) }}</span>
					</div>
					<div class="flex justify-between items-start gap-1">
						<span class="shrink-0">NO.INV:</span>
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
					<div class="text-[9px] pl-2 mt-1 font-black text-black">* AROMA: {{ item.nama_parfum || 'PARFUM STANDAR' }}</div>
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

				<div class="mt-4 mb-4 flex items-center gap-2 overflow-hidden opacity-80 text-black print:mt-0">
					<span class="font-mono text-xs font-black tracking-widest leading-none">-</span>
					<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
						<path stroke-linecap="round" stroke-linejoin="round" d="M14.121 14.121L19 19m-7-7l3-3m-3 3L7.05 9.05a2.828 2.828 0 114-4l4 4a2.828 2.828 0 11-4 4zm0 0l-3 3m3-3l3 3m-3-3l-4 4a2.828 2.828 0 11-4-4l4-4a2.828 2.828 0 114 4z" />
					</svg>
					<span class="font-mono text-xs font-black tracking-widest leading-none w-full border-b border-dashed border-black"></span>
				</div>

				<div class="text-center font-mono leading-none text-black break-inside-avoid">
					<h2 class="font-black text-[10px] uppercase tracking-tighter mb-1">{{ invoiceData?.toko_nama || 'LAUNDRY POS' }}</h2>
					<div class="border-2 border-black p-2 my-2 bg-black text-white">
						<div class="text-[9px] font-bold tracking-widest mb-1 opacity-90">LOKASI RAK BAJU:</div>
						<div class="text-3xl font-black tracking-tighter">{{ currentRack }}</div>
					</div>
					<div class="flex justify-between items-center text-[10px] font-black uppercase text-left mt-2">
						<span class="w-1/2 truncate pr-1">P: {{ invoiceData.pelanggan || 'UMUM' }}</span>
						<span class="w-1/2 text-right truncate">#{{ invoiceData.invoice ? invoiceData.invoice.split('/').pop() : 'INV' }}</span>
					</div>
					<div class="text-left text-[9px] font-black mt-1">ITEM: {{ invoiceData.items ? invoiceData.items.length : 0 }} BARANG</div>
				</div>

				<div class="hidden print:block font-mono text-xs leading-none select-none text-white">
					<br />
					&nbsp; .
				</div>
			</div>

			<div class="mt-4 md:mt-6 flex flex-col gap-2 md:gap-3 shrink-0 print:hidden relative z-10">
				<button @click="triggerPrint" class="w-full bg-slate-900 hover:bg-slate-800 text-white py-3 md:py-4 rounded-xl md:rounded-2xl font-black text-xs uppercase tracking-[0.2em] flex items-center justify-center gap-2 active:scale-95 transition-all">
					<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
						<path stroke-linecap="round" stroke-linejoin="round" d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z" />
					</svg>
					Cetak Struk Ganda
				</button>
				<button @click="emit('close')" class="w-full bg-slate-100 text-slate-500 py-3 md:py-3.5 rounded-xl md:rounded-2xl font-black text-xs uppercase tracking-widest hover:bg-slate-200 transition-all">Selesai & Tutup</button>
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
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
	background: #94a3b8;
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
		padding: 0mm 2mm 2mm 3mm !important;
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
