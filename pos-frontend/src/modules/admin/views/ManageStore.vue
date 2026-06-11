<script setup>
import Swal from 'sweetalert2';
import { onMounted, ref, computed } from 'vue';
import api from '../../../api.js'; // Sesuaikan relative path ke api.js lu bray

// Import sub-komponen kedaulatan kita
import StoreHeader from '../components/store/StoreHeader.vue';
import StoreTable from '../components/store/StoreTable.vue';

const stores = ref([]);
const isLoading = ref(false);

// STATE CONTROL: Input pencarian dan metode sorting
const searchQuery = ref('');
const sortBy = ref('latest'); // default: latest, soonest_expiry, longest_expiry

const fetchStores = async () => {
	isLoading.value = true;
	try {
		const res = await api.get('/admin/stores');
		if (res.data.status === 'sukses') stores.value = res.data.data;
	} catch (err) {
		Swal.fire({
			icon: 'error',
			title: 'Gagal Memuat Data',
			text: err.response?.data?.error || 'Koneksi terputus.',
			confirmButtonColor: '#ef4444',
			customClass: { popup: 'rounded-[24px]' },
		});
	} finally {
		isLoading.value = false;
	}
};

// 🚀 FUNGSI UTAMA SUSPEND (YANG TADI HILANG)
const onSuspendStore = (store) => {
	Swal.fire({
		title: `Bekukan ${store.nama_toko}?`,
		text: 'Toko ini otomatis terkunci dari aktivitas kasir!',
		icon: 'warning',
		showCancelButton: true,
		confirmButtonColor: '#ef4444',
		cancelButtonColor: '#334155',
		confirmButtonText: 'Ya, Suspend!',
		customClass: { popup: 'rounded-[24px]' },
	}).then(async (result) => {
		if (result.isConfirmed) {
			try {
				const res = await api.put(`/admin/stores/${store.id}/suspend`);
				if (res.data.status === 'sukses') {
					Swal.fire({ icon: 'success', title: 'Berhasil!', timer: 1200, showConfirmButton: false });
					fetchStores();
				}
			} catch (err) {
				Swal.fire({ icon: 'error', title: 'Gagal memproses.' });
			}
		}
	});
};

// 🚀 FUNGSI UTAMA ACTIVATE (YANG TADI HILANG)
const onActivateStore = (store) => {
	Swal.fire({
		title: `Aktifkan ${store.nama_toko}?`,
		text: 'Masa aktif toko akan diperpanjang 30 hari ke depan!',
		icon: 'question',
		showCancelButton: true,
		confirmButtonColor: '#10b981',
		cancelButtonColor: '#334155',
		confirmButtonText: 'Ya, Aktifkan!',
		customClass: { popup: 'rounded-[24px]' },
	}).then(async (result) => {
		if (result.isConfirmed) {
			try {
				const res = await api.put(`/admin/stores/${store.id}/activate`);
				if (res.data.status === 'sukses') {
					Swal.fire({ icon: 'success', title: 'Berhasil!', timer: 1200, showConfirmButton: false });
					fetchStores();
				}
			} catch (err) {
				Swal.fire({ icon: 'error', title: 'Gagal memproses.' });
			}
		}
	});
};

// JURUS COMPUTED: Live Search & Sort Tanpa Beban Database
const filteredAndSortedStores = computed(() => {
	let result = [...stores.value];

	// 1. Logika Live Search
	if (searchQuery.value.trim() !== '') {
		const query = searchQuery.value.toLowerCase();
		result = result.filter((store) => {
			return store.nama_toko?.toLowerCase().includes(query) || store.business_type?.toLowerCase().includes(query) || store.owner_name?.toLowerCase().includes(query);
		});
	}

	// 2. Logika Sorting Berdasarkan Masa Aktif / Toko Terbaru
	if (sortBy.value === 'latest') {
		result.sort((a, b) => b.id - a.id);
	} else if (sortBy.value === 'soonest_expiry') {
		result.sort((a, b) => {
			if (!a.subscription_end) return 1;
			if (!b.subscription_end) return -1;
			return new Date(a.subscription_end) - new Date(b.subscription_end);
		});
	} else if (sortBy.value === 'longest_expiry') {
		result.sort((a, b) => {
			if (!a.subscription_end) return 1;
			if (!b.subscription_end) return -1;
			return new Date(b.subscription_end) - new Date(a.subscription_end);
		});
	}

	return result;
});

onMounted(() => {
	fetchStores();
});
</script>

<template>
	<div class="p-6 bg-[#0B0F19] min-h-screen text-white font-sans">
		<StoreHeader :is-loading="isLoading" @refresh="fetchStores" />

		<div class="flex flex-col md:flex-row gap-4 mb-6 relative z-10">
			<div class="flex-1 relative">
				<span class="absolute inset-y-0 left-0 flex items-center pl-4 pointer-events-none">
					<svg class="w-4 h-4 text-slate-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.602 10.602Z" />
					</svg>
				</span>
				<input v-model="searchQuery" type="text" placeholder="Cari nama ruko, tipe bisnis, atau nama pemilik..." class="w-full pl-11 pr-5 py-3.5 bg-[#131B2E] border border-slate-800 rounded-2xl text-white text-xs font-bold focus:outline-none focus:border-indigo-500 transition-all placeholder:text-slate-600 shadow-xl" />
			</div>

			<div class="relative min-w-[240px]">
				<span class="absolute inset-y-0 left-0 flex items-center pl-4 pointer-events-none">
					<svg class="w-4 h-4 text-indigo-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" d="M3 4.5h18m-18 5h18m-18 5h13.5m-13.5 5h9" />
					</svg>
				</span>

				<select v-model="sortBy" class="w-full pl-11 pr-10 py-3.5 bg-[#131B2E] border border-slate-800 rounded-2xl text-white text-xs font-black uppercase tracking-wider focus:outline-none focus:border-indigo-500 transition-all cursor-pointer appearance-none shadow-xl">
					<option value="latest" class="bg-[#131B2E]">Toko Terbaru</option>
					<option value="soonest_expiry" class="bg-[#131B2E]">Masa Aktif Segera Habis</option>
					<option value="longest_expiry" class="bg-[#131B2E]">Masa Aktif Terlama</option>
				</select>

				<span class="absolute inset-y-0 right-0 flex items-center pr-4 pointer-events-none text-slate-500">
					<svg class="w-3.5 h-3.5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="3" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" d="m19.5 8.25-7.5 7.5-7.5-7.5" />
					</svg>
				</span>
			</div>
		</div>

		<StoreTable :stores="filteredAndSortedStores" :is-loading="isLoading" @suspend="onSuspendStore" @activate="onActivateStore" />
	</div>
</template>
