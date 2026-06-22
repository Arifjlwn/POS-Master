<script setup>
defineProps({
	reports: { type: Array, required: true },
	expandedRows: { type: Array, required: true },
	isLoading: { type: Boolean, required: true },
	formatRupiah: { type: Function, required: true },
	formatDate: { type: Function, required: true },
});
const emit = defineEmits(['toggle-row', 'print-faktur']);
</script>

<template>
	<div class="bg-white border border-slate-200 rounded-[28px] overflow-hidden shadow-sm">
		<div v-if="isLoading" class="p-10 flex flex-col items-center justify-center text-slate-400">
			<svg class="animate-spin w-8 h-8 mb-3 text-indigo-600" fill="none" viewBox="0 0 24 24">
				<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
				<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
			</svg>
			<p class="text-[10px] font-black uppercase tracking-widest animate-pulse">Menyinkronkan Data...</p>
		</div>

		<div v-else-if="reports.length === 0" class="p-16 flex flex-col items-center text-center opacity-50">
			<svg xmlns="http://www.w3.org/2000/svg" class="w-16 h-16 mb-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" /></svg>
			<p class="text-sm font-black text-slate-500 uppercase tracking-widest">Tidak ada data penerimaan</p>
		</div>

		<div v-else class="w-full">
			<table class="w-full text-left border-collapse block md:table">
				<thead class="hidden md:table-header-group">
					<tr class="bg-slate-50 border-b border-slate-200">
						<th class="p-4 text-[10px] font-black text-slate-400 uppercase tracking-widest w-12 text-center">Detil</th>
						<th class="p-4 text-[10px] font-black text-slate-400 uppercase tracking-widest">Waktu & Faktur</th>
						<th class="p-4 text-[10px] font-black text-slate-400 uppercase tracking-widest">Supplier</th>
						<th class="p-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-center">Total Item</th>
						<th class="p-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-right">Nilai Tagihan</th>
						<th class="p-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-center">Aksi</th>
					</tr>
				</thead>

				<tbody v-for="lpb in reports" :key="lpb.id" class="block md:table-row-group">
					<tr class="block md:table-row border-b border-slate-200 md:border-slate-100 hover:bg-slate-50/50 transition-colors group p-4 md:p-0">
						<td class="flex items-center justify-between md:table-cell md:p-4 md:text-center border-b md:border-none border-slate-100 pb-3 mb-3 md:pb-0 md:mb-0">
							<div class="md:hidden text-left">
								<div class="font-black text-xs text-slate-800 uppercase">{{ lpb.no_faktur }}</div>
								<div class="text-[9px] font-bold text-slate-400 mt-0.5">{{ formatDate(lpb.created_at) }}</div>
							</div>

							<button @click="emit('toggle-row', lpb.id)" class="w-8 h-8 md:w-7 md:h-7 bg-white border border-slate-200 rounded-lg flex items-center justify-center text-slate-400 hover:text-indigo-600 hover:border-indigo-200 transition-all shadow-sm">
								<svg xmlns="http://www.w3.org/2000/svg" :class="expandedRows.includes(lpb.id) ? 'rotate-180 text-indigo-600' : ''" class="w-4 h-4 transition-transform" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" /></svg>
							</button>
						</td>

						<td class="hidden md:table-cell p-4">
							<div class="font-black text-xs text-slate-800 uppercase">{{ lpb.no_faktur }}</div>
							<div class="text-[9px] font-bold text-slate-400 mt-1">{{ formatDate(lpb.created_at) }}</div>
						</td>

						<td class="flex justify-between items-center md:table-cell py-2 md:p-4">
							<span class="text-[10px] font-black text-slate-400 uppercase tracking-widest md:hidden">Supplier</span>
							<span class="font-bold text-xs text-slate-700 uppercase">{{ lpb.nama_supplier }}</span>
						</td>

						<td class="flex justify-between items-center md:table-cell py-2 md:p-4 md:text-center">
							<span class="text-[10px] font-black text-slate-400 uppercase tracking-widest md:hidden">Total Item</span>
							<span class="bg-sky-50 text-sky-600 px-2.5 py-1 rounded-md text-[10px] font-black border border-sky-100">{{ lpb.total_item }} PCS</span>
						</td>

						<td class="flex justify-between items-center md:table-cell py-2 md:p-4 md:text-right">
							<span class="text-[10px] font-black text-slate-400 uppercase tracking-widest md:hidden">Nilai Tagihan</span>
							<span class="font-black text-sm text-slate-900">{{ formatRupiah(lpb.total_modal) }}</span>
						</td>

						<td class="flex justify-between items-center md:table-cell py-3 md:p-4 md:text-center mt-2 md:mt-0 border-t md:border-none border-slate-100">
							<span class="text-[10px] font-black text-slate-400 uppercase tracking-widest md:hidden">Cetak Laporan</span>
							<button @click="emit('print-faktur', lpb)" class="p-2 bg-slate-900 hover:bg-indigo-600 text-white rounded-lg transition-colors shadow-sm active:scale-95" title="Cetak Faktur LPB">
								<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z" /></svg>
							</button>
						</td>
					</tr>

					<tr v-show="expandedRows.includes(lpb.id)" class="block md:table-row bg-slate-50 border-b border-slate-200">
						<td colspan="6" class="block md:table-cell p-3 md:p-5">
							<div class="bg-white rounded-xl border border-slate-200 p-3 md:p-4 shadow-inner overflow-x-auto custom-scrollbar">
								<h4 class="text-[10px] font-black text-indigo-600 uppercase tracking-widest mb-3 flex items-center gap-2">
									<svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" /></svg>
									Rincian Barang Diterima
								</h4>
								<table class="w-full text-left min-w-[500px]">
									<thead>
										<tr class="border-b border-slate-100">
											<th class="py-2 text-[9px] font-black text-slate-400 uppercase">Nama Produk</th>
											<th class="py-2 text-[9px] font-black text-slate-400 uppercase text-center">Qty Terima</th>
											<th class="py-2 text-[9px] font-black text-slate-400 uppercase text-right">Modal/Pcs</th>
											<th class="py-2 text-[9px] font-black text-slate-400 uppercase text-right">Subtotal</th>
										</tr>
									</thead>
									<tbody>
										<tr v-for="item in lpb.items" :key="item.id" class="border-b border-slate-50 last:border-0">
											<td class="py-2.5 text-xs font-bold text-slate-700 uppercase">{{ item.nama_produk || 'Item Tidak Diketahui' }}</td>
											<td class="py-2.5 text-xs font-black text-slate-900 text-center">{{ item.qty }}</td>
											<td class="py-2.5 text-xs font-bold text-slate-500 text-right">{{ formatRupiah(item.harga_modal) }}</td>
											<td class="py-2.5 text-xs font-black text-slate-900 text-right">{{ formatRupiah(item.sub_total) }}</td>
										</tr>
									</tbody>
								</table>
							</div>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
	height: 6px;
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
