<script setup>
import Swal from 'sweetalert2';
import { onMounted, ref } from 'vue';
import api from '../../../../../api.js';
import SidebarLaundry from '../../components/SidebarLaundry.vue';

const racks = ref([]);
const isLoading = ref(true);

const fetchRacksData = async () => {
	isLoading.value = true;
	try {
		const res = await api.get('/laundry/racks');
		racks.value = res.data.data || [];
	} catch (err) {
		console.error(err);
	} finally {
		isLoading.value = false;
	}
};

const handleToggleStatus = (rack) => {
	const targetStatus = rack.status === 'TERSEDIA' ? 'RUSAK' : 'TERSEDIA';

	Swal.fire({
		title: `Ubah Status Rak ${rack.nama_rak}?`,
		text: `Rak ini akan di-set menjadi ${targetStatus}. Algoritma POS kasir otomatis me-redirect cucian keluar dari rak rusak.`,
		icon: 'warning',
		showCancelButton: true,
		confirmButtonColor: rack.status === 'TERSEDIA' ? '#ef4444' : '#10b981',
		confirmButtonText: `Ya, Set ${targetStatus}`,
		cancelButtonText: 'Batal',
		customClass: { popup: 'rounded-[28px]' },
	}).then(async (result) => {
		if (result.isConfirmed) {
			try {
				const res = await api.put(`/laundry/racks/${rack.id}/status`, {
					status: targetStatus,
				});
				if (res.data.status === 'sukses') {
					Swal.fire({ toast: true, position: 'top-end', icon: 'success', title: `Status Rak ${rack.nama_rak} diperbarui!`, showConfirmButton: false, timer: 1500 });
					rack.status = targetStatus;
				}
			} catch (e) {
				Swal.fire('Gagal Sinkron', 'Gagal merubah status fisik slot laci.', 'error');
			}
		}
	});
};

// 🚀 ENGINE PREVIEW RAK KASTA TERTINGGI BRAY!
const handleRackClick = (rack) => {
	if (rack.isi_cucian > 0) {
		let htmlContent = `<div class="text-left font-sans space-y-3 max-h-72 overflow-y-auto custom-scrollbar pr-1 mt-4">`;

		rack.detail_cucian.forEach((item, index) => {
			htmlContent += `
        <div class="bg-slate-50 p-4 rounded-2xl border border-slate-200 shadow-sm relative">
            <div class="absolute -left-1 -top-1 w-6 h-6 bg-slate-800 text-white rounded-lg flex items-center justify-center font-black text-[10px] shadow-md">${index + 1}</div>
            <div class="ml-3">
                <div class="flex justify-between items-start mb-1">
                    <div class="font-black text-sm text-slate-800 uppercase tracking-tight line-clamp-1">${item.product?.nama_produk || 'PAKET LAUNDRY'}</div>
                    <div class="text-[9px] font-black px-2 py-1 rounded bg-indigo-100 text-indigo-700 uppercase">${item.status_cucian}</div>
                </div>
                
                <div class="text-[10px] font-bold text-slate-500 mb-1 flex items-center gap-1 uppercase">
                    <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" /></svg>
                    ${item.nama_pelanggan}
                </div>
                
                <div class="text-[10px] font-bold text-slate-400 mb-2 flex items-center gap-1">
                    <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" /></svg>
                    +${item.no_whatsapp}
                </div>
                <div class="flex justify-between items-center bg-white p-2 rounded-xl border border-slate-100">
                    <span class="text-[10px] font-black uppercase text-slate-400">Muatan: <span class="text-slate-800">${item.berat_kg} KG</span></span>
                    <span class="text-xs font-black text-slate-800">Rp ${new Intl.NumberFormat('id-ID').format(item.sub_total)}</span>
                </div>
            </div>
        </div>
    `;
		});

		htmlContent += `</div>`;

		Swal.fire({
			title: `DETAIL RAK ${rack.nama_rak}`,
			html: htmlContent,
			showCloseButton: true,
			showConfirmButton: true,
			confirmButtonText: 'Tutup Preview',
			confirmButtonColor: '#0f172a',
			customClass: { popup: 'rounded-[32px] p-6' },
		});
	} else {
		// Kalau raknya bener-bener kosong, trigger fungsi ubah status rusak/tersedia.
		handleToggleStatus(rack);
	}
};

onMounted(() => {
	fetchRacksData();
});
</script>

<template>
	<SidebarLaundry>
		<div class="p-4 md:p-6 bg-[#F8FAFC] min-h-screen font-sans">
			<div class="flex flex-col sm:flex-row justify-between sm:items-center gap-4 bg-white p-5 rounded-3xl shadow-sm border border-slate-200 mb-6">
				<div>
					<div class="flex items-center gap-2">
						<span class="bg-indigo-100 text-indigo-700 font-black text-[9px] uppercase tracking-widest px-2.5 py-1 rounded-md border border-indigo-200">Laundry Engine</span>
					</div>
					<h1 class="text-xl font-black text-slate-900 uppercase tracking-tight mt-1">Tata Letak Laci Rak Baju</h1>
					<p class="text-xs font-bold text-slate-400 mt-0.5">Pantau kapasitas muatan, detail pesanan, dan atur status operasional rak secara *real-time*.</p>
				</div>

				<div class="flex flex-wrap items-center gap-2 font-mono print:hidden">
					<div class="bg-emerald-50 text-emerald-700 font-bold text-[10px] px-3 py-2 rounded-xl flex items-center gap-1.5 border border-emerald-200">
						<span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
						KOSONG: {{ racks.filter((r) => r.status === 'TERSEDIA' && r.isi_cucian === 0).length }}
					</div>
					<div class="bg-amber-50 text-amber-700 font-bold text-[10px] px-3 py-2 rounded-xl flex items-center gap-1.5 border border-amber-200">
						<span class="w-2 h-2 rounded-full bg-amber-500"></span>
						TERISI: {{ racks.filter((r) => r.isi_cucian > 0).length }}
					</div>
					<div class="bg-rose-50 text-rose-700 font-bold text-[10px] px-3 py-2 rounded-xl flex items-center gap-1.5 border border-rose-200">
						<span class="w-2 h-2 rounded-full bg-rose-500"></span>
						RUSAK: {{ racks.filter((r) => r.status === 'RUSAK').length }}
					</div>
				</div>
			</div>

			<div class="bg-white p-6 rounded-[32px] shadow-sm border border-slate-200 min-h-[50vh]">
				<div v-if="isLoading" class="flex flex-col items-center justify-center h-60 text-slate-400">
					<div class="w-8 h-8 border-4 border-slate-200 border-t-indigo-600 rounded-full animate-spin mb-3"></div>
					<p class="font-black text-[10px] uppercase tracking-widest animate-pulse">Memuat Topologi Denah Rak...</p>
				</div>

				<div v-else-if="racks.length === 0" class="flex flex-col items-center justify-center h-60 text-center">
					<div class="w-14 h-14 bg-slate-50 text-slate-400 border-2 border-dashed border-slate-200 rounded-2xl flex items-center justify-center mb-3">
						<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
							<path stroke-linecap="round" stroke-linejoin="round" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
						</svg>
					</div>
					<h3 class="text-sm font-black text-slate-800 uppercase tracking-wider">Sistem Rak Belum Di-Setup</h3>
					<p class="text-xs font-bold text-slate-400 max-w-xs mt-1 leading-normal">Silakan masuk ke halaman utama POS laundry untuk memicu setup sistem rak.</p>
				</div>

				<div v-else class="grid grid-cols-2 sm:grid-cols-4 md:grid-cols-6 lg:grid-cols-8 xl:grid-cols-10 gap-3 md:gap-4 select-none">
					<div
						v-for="rack in racks"
						:key="rack.id"
						@click="handleRackClick(rack)"
						:class="[
							'p-4 rounded-2xl border-2 cursor-pointer transition-all duration-300 active:scale-95 flex flex-col items-center justify-center text-center font-mono relative overflow-hidden group min-h-[90px]',
							rack.status === 'RUSAK' ? 'bg-rose-50 border-rose-300 text-rose-700 shadow-inner' : rack.isi_cucian > 0 ? 'bg-amber-50 border-amber-300 hover:border-amber-500 text-amber-900 shadow-md ring-2 ring-amber-100/50' : 'bg-slate-50 border-slate-200 hover:border-indigo-400 hover:bg-white text-slate-900 shadow-sm',
						]">
						<span class="absolute left-1.5 top-1.5 text-[8px] font-bold opacity-60 leading-none group-hover:text-indigo-600 transition-colors">R{{ rack.baris }}C{{ rack.kolom }}</span>

						<div v-if="rack.isi_cucian > 0" class="absolute right-1.5 top-1.5 w-5 h-5 bg-amber-500 text-white rounded-full flex items-center justify-center text-[10px] font-black shadow-sm group-hover:scale-110 transition-transform">
							{{ rack.isi_cucian }}
						</div>

						<span class="text-xl font-black tracking-tighter leading-none block mt-2 z-10">{{ rack.nama_rak }}</span>

						<span :class="['text-[8px] font-black uppercase tracking-widest px-1.5 py-0.5 rounded border mt-1.5 block leading-none z-10 transition-colors', rack.status === 'RUSAK' ? 'bg-rose-600 border-rose-700 text-white' : rack.isi_cucian > 0 ? 'bg-amber-200/50 border-amber-300 text-amber-700' : 'bg-indigo-50 border-indigo-100 text-indigo-600']">
							{{ rack.status === 'RUSAK' ? 'RUSAK' : rack.isi_cucian > 0 ? 'TERISI' : 'KOSONG' }}
						</span>
					</div>
				</div>
			</div>
		</div>
	</SidebarLaundry>
</template>

<style scoped>
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
