<script setup>
import { ref, onMounted, computed } from 'vue';
import Swal from 'sweetalert2';
import api from '../../../api.js';

// IMPORT SUB-KOMPONEN KEDAULATAN BARU
import TenantStats from '../components/store/TenantStats.vue';
import TenantList from '../components/store/TenantList.vue';

const isLoading = ref(false);
const planCounts = ref({ trial: 0, basic: 0, pro: 0, premium: 0 });
const stores = ref([]);

const searchQuery = ref('');
const sortBy = ref('latest');

const fetchTenantHubData = async () => {
    isLoading.value = true;
    try {
        const res = await api.get('/admin/subscription-overview');
        if (res.data.status === 'sukses') {
            planCounts.value = res.data.counts;
            stores.value = res.data.stores;
        }
    } catch (err) {
        Swal.fire({
            icon: 'error',
            title: 'Koneksi Terputus',
            text: err.response?.data?.error || 'Gagal memuat telemetri super admin bray.',
            confirmButtonColor: '#ef4444',
            customClass: { popup: 'rounded-[24px]' },
        });
    } finally {
        isLoading.value = false;
    }
};

const onSuspendStore = (store) => {
    Swal.fire({
        title: `Bekukan ${store.nama_toko}?`,
        text: 'Cabang ini otomatis terkunci dari aktivitas operasional!',
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
					Swal.fire({ icon: 'success', title: 'Tenant Suspended!', timer: 1200, showConfirmButton: false });
					fetchTenantHubData();
				}
			} catch (err) {
				Swal.fire({ icon: 'error', title: 'Gagal mengeksekusi perintah.' });
			}
		}
	});
};

const onActivateStore = (store) => {
    Swal.fire({
        title: `Aktifkan ${store.nama_toko}?`,
        text: 'Sistem akan memperpanjang masa aktif lisensi ruko +30 Hari ke depan!',
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
					Swal.fire({ icon: 'success', title: 'Tenant Aktif Kembali!', timer: 1200, showConfirmButton: false });
					fetchTenantHubData();
				}
			} catch (err) {
				Swal.fire({ icon: 'error', title: 'Gagal memproses.' });
			}
		}
	});
};

const filteredAndSortedStores = computed(() => {
    let result = [...stores.value];

    if (searchQuery.value.trim() !== '') {
        const query = searchQuery.value.toLowerCase();
        result = result.filter((store) => {
            return store.nama_toko?.toLowerCase().includes(query) || 
                   store.business_type?.toLowerCase().includes(query) || 
                   store.owner_name?.toLowerCase().includes(query);
        });
    }

    if (sortBy.value === 'latest') {
        result.sort((a, b) => b.id - a.id);
    } else if (sortBy.value === 'soonest_expiry') {
        result.sort((a, b) => a.sisa_hari - b.sisa_hari);
    } else if (sortBy.value === 'longest_expiry') {
        result.sort((a, b) => b.sisa_hari - a.sisa_hari);
    }

    return result;
});

onMounted(() => fetchTenantHubData());
</script>

<template>
    <div class="p-4 sm:p-6 lg:p-8 bg-[#0B0F19] min-h-screen text-white font-sans overflow-x-hidden">
        
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-8">
            <div>
                <div class="inline-flex items-center gap-2 px-3 py-1 bg-indigo-500/10 border border-indigo-500/20 text-indigo-400 rounded-xl text-[10px] font-black uppercase tracking-widest mb-2">Central Mission Control</div>
                <h1 class="text-xl sm:text-2xl font-black text-white tracking-tight uppercase">Unified Tenant & Subscription Hub</h1>
                <p class="text-slate-500 font-bold text-[9px] sm:text-[10px] uppercase tracking-widest mt-0.5">Kendali Mutlak Operasional & Financing Ekosistem Tenant POS SaaS</p>
            </div>
            
            <button @click="fetchTenantHubData" :disabled="isLoading" class="w-full sm:w-auto px-5 py-3 bg-[#131B2E] hover:bg-[#1a243d] disabled:opacity-50 border border-slate-800 rounded-xl font-black text-[10px] uppercase tracking-widest transition-all flex items-center justify-center gap-2 shadow-lg">
                <svg :class="{ 'animate-spin': isLoading }" class="w-4 h-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0 3.181 3.183a8.25 8.25 0 0 0 13.803-3.7M4.031 9.865a8.25 8.25 0 0 1 13.803-3.7l3.181 3.182m0-4.991v4.99" />
                </svg>
                <span>{{ isLoading ? 'Synchronizing...' : 'Refresh Data' }}</span>
            </button>
        </div>

        <TenantStats :plan-counts="planCounts" />

        <div class="flex flex-col md:flex-row gap-4 mb-6 relative z-10">
            <div class="flex-1 relative">
                <span class="absolute inset-y-0 left-0 flex items-center pl-4 pointer-events-none">
                    <svg class="w-4 h-4 text-slate-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.602 10.602Z" />
                    </svg>
                </span>
                <input v-model="searchQuery" type="text" placeholder="Cari nama gerai ruko, tipe bisnis, atau nama pemilik..." class="w-full pl-11 pr-5 py-3.5 bg-[#131B2E] border border-slate-800 rounded-2xl text-white text-xs font-bold focus:outline-none focus:border-indigo-500 transition-all placeholder:text-slate-600 shadow-xl" />
            </div>

            <div class="relative w-full md:w-64">
                <select v-model="sortBy" class="w-full pl-5 pr-10 py-3.5 bg-[#131B2E] border border-slate-800 rounded-2xl text-white text-xs font-black uppercase tracking-wider focus:outline-none focus:border-indigo-500 transition-all cursor-pointer appearance-none shadow-xl">
                    <option value="latest">Tenant Baru Bergabung</option>
                    <option value="soonest_expiry">Lisensi Segera Expired</option>
                    <option value="longest_expiry">Lisensi Masa Aktif Terlama</option>
                </select>
                <span class="absolute inset-y-0 right-0 flex items-center pr-4 pointer-events-none text-slate-500">
                    <svg class="w-3.5 h-3.5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="3" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="m19.5 8.25-7.5 7.5-7.5-7.5" /></svg>
                </span>
            </div>
        </div>

        <TenantList 
            :stores="filteredAndSortedStores" 
            :is-loading="isLoading" 
            @suspend="onSuspendStore" 
            @activate="onActivateStore" 
        />

    </div>
</template>