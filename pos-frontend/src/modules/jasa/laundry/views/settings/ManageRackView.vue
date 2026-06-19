<script setup>
import Swal from 'sweetalert2';
import { computed, onMounted, ref } from 'vue';
import api from '../../../../../api.js';
import SidebarLaundry from '../../components/SidebarLaundry.vue';

const racks = ref([]);
const isLoading = ref(true);
const isGenerating = ref(false);

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

// 🚀 COMPUTED ENGINE: Kelompokkan rak berdasarkan nama Zona/Lemari bray!
const groupedRacks = computed(() => {
	const groups = {};
	racks.value.forEach((rack) => {
		const zoneName = rack.zona || 'RAK UTAMA';
		if (!groups[zoneName]) {
			groups[zoneName] = [];
		}
		groups[zoneName].push(rack);
	});
	return groups;
});

const handleToggleStatus = (rack) => {
	const targetStatus = rack.status === 'TERSEDIA' ? 'RUSAK' : 'TERSEDIA';

	Swal.fire({
		title: `Ubah Status Rak ${rack.nama_rak}?`,
		text: `Rak di zona ${rack.zona} ini akan di-set menjadi ${targetStatus}.`,
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
			title: `ZONA ${rack.zona} - RAK ${rack.nama_rak}`,
			html: htmlContent,
			showCloseButton: true,
			showConfirmButton: true,
			confirmButtonText: 'Tutup Preview',
			confirmButtonColor: '#0f172a',
			customClass: { popup: 'rounded-[32px] p-6' },
		});
	} else {
		handleToggleStatus(rack);
	}
};

const handleEditRackConfig = () => {
	Swal.fire({
		title: 'Tambah Zona / Lemari Baru',
		html: `
            <div class="text-left font-sans space-y-4 mt-2">
                <div class="p-3 bg-indigo-50 border border-indigo-200 rounded-xl flex gap-3 text-indigo-800 shadow-inner">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                    <p class="text-[10px] font-bold leading-relaxed">
                        <b>CLUSTER INJECTION:</b> Masukkan nama identitas lemari baru Anda. Sistem akan meng-generate matrix terpisah tanpa mengganggu rak lama.
                    </p>
                </div>
                
                <div class="space-y-3 font-mono">
                    <div class="flex flex-col gap-1.5">
                        <label class="text-[10px] font-black text-slate-700 uppercase tracking-widest">Nama Lemari / Zona Rak:</label>
                        <input id="swal-input-zona" type="text" placeholder="Contoh: RAK BESI B, LEMARI DEPAN" class="p-3 bg-slate-50 border border-slate-200 rounded-xl outline-none font-black text-xs uppercase focus:border-indigo-600 transition-all text-slate-800" />
                    </div>
                    <div class="flex flex-col gap-1.5">
                        <label class="text-[10px] font-black text-slate-700 uppercase tracking-widest">Jumlah Baris Kebawah (A - Z):</label>
                        <input id="swal-input-baris" type="number" min="1" max="26" value="5" class="p-3 bg-slate-50 border border-slate-200 rounded-xl outline-none font-black text-sm focus:border-indigo-600 transition-all text-slate-800" />
                    </div>
                    <div class="flex flex-col gap-1.5">
                        <label class="text-[10px] font-black text-slate-700 uppercase tracking-widest">Jumlah Kolom Kekanan (1 - 50):</label>
                        <input id="swal-input-kolom" type="number" min="1" max="50" value="5" class="p-3 bg-slate-50 border border-slate-200 rounded-xl outline-none font-black text-sm focus:border-indigo-600 transition-all text-slate-800" />
                    </div>
                </div>
            </div>
        `,
		showCancelButton: true,
		confirmButtonText: 'Generate Zona',
		cancelButtonText: 'Batal',
		confirmButtonColor: '#4f46e5',
		cancelButtonColor: '#94a3b8',
		customClass: { popup: 'rounded-[32px] p-6', confirmButton: 'uppercase font-black text-xs px-4 py-2.5 rounded-xl', cancelButton: 'uppercase font-black text-xs px-4 py-2.5 rounded-xl' },
		preConfirm: () => {
			const zona = document.getElementById('swal-input-zona').value.trim().toUpperCase();
			const baris = document.getElementById('swal-input-baris').value;
			const kolom = document.getElementById('swal-input-kolom').value;

			if (!zona) {
				Swal.showValidationMessage('Nama Lemari / Zona wajib diisi bray!');
				return false;
			}
			if (!baris || !kolom || baris < 1 || kolom < 1) {
				Swal.showValidationMessage('Jumlah baris dan kolom tidak valid!');
				return false;
			}
			return { zona, baris: parseInt(baris), kolom: parseInt(kolom) };
		},
	}).then(async (result) => {
		if (result.isConfirmed) {
			isGenerating.value = true;
			try {
				const res = await api.post('/laundry/racks/setup', {
					zona: result.value.zona,
					jumlah_baris: result.value.baris,
					jumlah_kolom: result.value.kolom,
				});

				if (res.data.status === 'sukses') {
					Swal.fire({ icon: 'success', title: 'Zona Baru Berhasil Ditambah!', text: res.data.message, customClass: { popup: 'rounded-[28px]' }, confirmButtonColor: '#0f172a' });
					fetchRacksData();
				}
			} catch (err) {
				Swal.fire('Gagal Ekspansi', err.response?.data?.error || 'Sistem backend menolak konfigurasi.', 'error');
			} finally {
				isGenerating.value = false;
			}
		}
	});
};

// 🚀 ENGINE BARU: EDIT NAMA ZONA BRAY!
const handleEditZona = async (oldZonaName) => {
	Swal.fire({
		title: 'Ubah Nama Kluster',
		input: 'text',
		inputValue: oldZonaName,
		inputPlaceholder: 'Masukkan nama lemari/kluster baru...',
		showCancelButton: true,
		confirmButtonText: 'Simpan Perubahan',
		cancelButtonText: 'Batal',
		confirmButtonColor: '#4f46e5',
		customClass: { popup: 'rounded-[32px]', input: 'uppercase font-black text-sm text-center' },
		inputValidator: (value) => {
			if (!value) return 'Nama kluster tidak boleh kosong!';
			if (value.trim().toUpperCase() === oldZonaName) return 'Nama kluster belum diubah!';
		},
	}).then(async (result) => {
		if (result.isConfirmed) {
			try {
				const res = await api.put('/laundry/racks/zona', {
					old_zona: oldZonaName,
					new_zona: result.value.trim().toUpperCase(),
				});
				if (res.data.status === 'sukses') {
					Swal.fire({ toast: true, position: 'top-end', icon: 'success', title: 'Nama kluster berhasil diperbarui!', showConfirmButton: false, timer: 1500 });
					fetchRacksData();
				}
			} catch (err) {
				Swal.fire('Gagal Edit', err.response?.data?.error || 'Terjadi kesalahan sistem backend.', 'error');
			}
		}
	});
};

// 🚀 ENGINE BARU: HAPUS/MUSNAHKAN ZONA BRAY!
const handleDeleteZona = async (zoneName, racksInZone) => {
	// 🛡️ SECURITY KASTA TINGGI: Cek apakah ada baju nyangkut di dalam zona ini!
	const hasItems = racksInZone.some((r) => r.isi_cucian > 0);
	if (hasItems) {
		return Swal.fire('Akses Ditolak!', `Tidak bisa memusnahkan kluster <b>${zoneName}</b> karena masih ada slot laci yang berisi cucian pelanggan!`, 'error');
	}

	Swal.fire({
		title: 'Musnahkan Lemari?',
		html: `Seluruh slot matrix di dalam kluster <b>${zoneName}</b> akan dihapus permanen dari sistem POS. Tindakan ini tidak dapat dibatalkan!`,
		icon: 'warning',
		showCancelButton: true,
		confirmButtonColor: '#e11d48',
		confirmButtonText: 'Ya, Musnahkan!',
		cancelButtonText: 'Batal',
		customClass: { popup: 'rounded-[32px]' },
	}).then(async (result) => {
		if (result.isConfirmed) {
			try {
				// Tembak request delete dengan mengirimkan nama zona di payload body
				const res = await api.delete('/laundry/racks/zona', {
					data: { zona: zoneName },
				});
				if (res.data.status === 'sukses') {
					Swal.fire('Dihapus!', `Kluster ${zoneName} telah dimusnahkan.`, 'success');
					fetchRacksData();
				}
			} catch (err) {
				Swal.fire('Gagal Hapus', err.response?.data?.error || 'Terjadi kesalahan sistem backend.', 'error');
			}
		}
	});
};

onMounted(() => {
	fetchRacksData();
});
</script>

<template>
	<SidebarLaundry>
		<div class="p-4 md:p-6 bg-[#F8FAFC] min-h-screen font-sans">
			<div class="flex flex-col lg:flex-row justify-between lg:items-center gap-4 bg-white p-5 rounded-3xl shadow-sm border border-slate-200 mb-6 relative overflow-hidden">
				<div>
					<div class="flex items-center gap-2">
						<span class="bg-indigo-100 text-indigo-700 font-black text-[9px] uppercase tracking-widest px-2.5 py-1 rounded-md border border-indigo-200">Laundry Engine</span>
					</div>
					<h1 class="text-xl font-black text-slate-900 uppercase tracking-tight mt-1">Tata Letak Laci Rak Baju</h1>
					<p class="text-xs font-bold text-slate-400 mt-0.5 max-w-xl">Sistem Multi-Zona Cerdas. Memetakan kluster penempatan pakaian berdasarkan tata letak ruko fisik.</p>
				</div>

				<div class="flex flex-col sm:flex-row items-start sm:items-center gap-4 relative z-10">
					<div class="flex flex-wrap items-center gap-2 font-mono print:hidden bg-slate-50 p-1.5 rounded-2xl border border-slate-200">
						<div class="bg-emerald-50 text-emerald-700 font-bold text-[10px] px-3 py-1.5 rounded-xl flex items-center gap-1.5 border border-emerald-200">
							<span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
							KOSONG: {{ racks.filter((r) => r.status === 'TERSEDIA' && r.isi_cucian === 0).length }}
						</div>
						<div class="bg-amber-50 text-amber-700 font-bold text-[10px] px-3 py-1.5 rounded-xl flex items-center gap-1.5 border border-amber-200">
							<span class="w-2 h-2 rounded-full bg-amber-500"></span>
							TERISI: {{ racks.filter((r) => r.isi_cucian > 0).length }}
						</div>
						<div class="bg-rose-50 text-rose-700 font-bold text-[10px] px-3 py-1.5 rounded-xl flex items-center gap-1.5 border border-rose-200">
							<span class="w-2 h-2 rounded-full bg-rose-500"></span>
							RUSAK: {{ racks.filter((r) => r.status === 'RUSAK').length }}
						</div>
					</div>

					<button @click="handleEditRackConfig" :disabled="isLoading || isGenerating" class="shrink-0 bg-white hover:bg-indigo-50 text-indigo-600 border-2 border-indigo-100 hover:border-indigo-300 py-2.5 px-4 rounded-xl font-black text-[10px] uppercase tracking-widest active:scale-95 transition-all flex items-center gap-2 disabled:opacity-50">
						<svg v-if="isGenerating" class="animate-spin h-4 w-4 text-indigo-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
							<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
							<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
						</svg>
						<svg v-else xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
							<path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" />
						</svg>
						Tambah Lemari Baru
					</button>
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

				<div v-else class="space-y-10">
					<div v-for="(zoneRacks, zoneName) in groupedRacks" :key="zoneName" class="border-b border-slate-100 pb-8 last:border-none">
						<div class="flex items-center justify-between mb-4">
							<h2 class="text-xs font-black text-slate-800 uppercase tracking-[0.2em] flex items-center gap-2 bg-slate-50 px-4 py-2.5 rounded-xl border border-slate-200 inline-flex shadow-sm">
								<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
									<path stroke-linecap="round" stroke-linejoin="round" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
								</svg>
								Kluster: {{ zoneName }}
							</h2>

							<div class="flex items-center gap-1.5 bg-white p-1 rounded-xl border border-slate-100 shadow-sm">
								<button @click="handleEditZona(zoneName)" class="p-2 text-indigo-500 hover:bg-indigo-50 hover:text-indigo-600 rounded-lg transition-colors border border-transparent hover:border-indigo-100 active:scale-95" title="Ubah Nama Lemari">
									<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
										<path stroke-linecap="round" stroke-linejoin="round" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
									</svg>
								</button>
								<button @click="handleDeleteZona(zoneName, zoneRacks)" class="p-2 text-rose-400 hover:bg-rose-50 hover:text-rose-600 rounded-lg transition-colors border border-transparent hover:border-rose-100 active:scale-95" title="Musnahkan Lemari">
									<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
										<path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
									</svg>
								</button>
							</div>
						</div>

						<div class="grid grid-cols-2 sm:grid-cols-4 md:grid-cols-6 lg:grid-cols-8 xl:grid-cols-10 gap-3 md:gap-4">
							<div
								v-for="rack in zoneRacks"
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
