<script setup>
import { ref, onMounted } from 'vue';
import Swal from 'sweetalert2';
import api from '../../../api.js';

// IMPORT SUB-COMPONENTS BRAY
import SaaSRevenueCards from '../components/analytics/SaaSRevenueCards.vue';
import SaaSChartGrid from '../components/analytics/SaaSChartGrid.vue';
import WhaleTenantTable from '../components/analytics/WhaleTenantTable.vue';

const isLoading = ref(false);

const telemetry = ref({
    mrr: 0,
    arr: 0,
    total_tenants: 0,
    new_tenants: 0,
    churned_tenants: 0,
    top_tenants: [],
    monthly_growth_labels: [], 
    monthly_tenant_growth: [], 
    monthly_trans_growth: []   
});

const formatRupiah = (value) => {
    return new Intl.NumberFormat('id-ID', {
        style: 'currency',
        currency: 'IDR',
        minimumFractionDigits: 0
    }).format(value);
};

const fetchAnalyticsData = async () => {
    isLoading.value = true;
    try {
        const res = await api.get('/admin/analytics/telemetry');
        if (res.data.status === 'sukses') {
            telemetry.value = res.data.data;
        }
    } catch (err) {
        Swal.fire({
            icon: 'error',
            title: 'Koneksi Gagal',
            text: err.response?.data?.error || 'Gagal memetakan telemetri finansial dari server pusat.',
            confirmButtonColor: '#ef4444',
            customClass: { popup: 'rounded-[24px]' },
        });
    } finally {
        isLoading.value = false;
    }
};

onMounted(() => fetchAnalyticsData());
</script>

<template>
    <div class="p-4 sm:p-6 lg:p-8 bg-[#0B0F19] min-h-screen text-white font-sans overflow-x-hidden">
        
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-8">
            <div>
                <div class="inline-flex items-center gap-2 px-3 py-1 bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 rounded-xl text-[10px] font-black uppercase tracking-widest mb-2">Financing & Analytics Center</div>
                <h1 class="text-xl sm:text-2xl font-black text-white tracking-tight uppercase">SaaS Financial Telemetry</h1>
                <p class="text-slate-500 font-bold text-[9px] sm:text-[10px] uppercase tracking-widest mt-0.5">Analisis Akumulasi Pendapatan Bersih, Retensi, Dan Skalabilitas Platform</p>
            </div>
            
            <button @click="fetchAnalyticsData" :disabled="isLoading" class="w-full sm:w-auto px-5 py-3 bg-[#131B2E] hover:bg-[#1a243d] disabled:opacity-50 border border-slate-800 rounded-xl font-black text-[10px] uppercase tracking-widest transition-all flex items-center justify-center gap-2 shadow-lg">
                <svg :class="{ 'animate-spin': isLoading }" class="w-4 h-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0 3.181 3.183a8.25 8.25 0 0 0 13.803-3.7M4.031 9.865a8.25 8.25 0 0 1 13.803-3.7l3.181 3.182m0-4.991v4.99" />
                </svg>
                <span>{{ isLoading ? 'Recalculating...' : 'Refresh Metrics' }}</span>
            </button>
        </div>

        <SaaSRevenueCards 
            :mrr="telemetry.mrr" 
            :arr="telemetry.arr" 
            :total-tenants="telemetry.total_tenants" 
            :new-tenants="telemetry.new_tenants" 
            :formatRupiah="formatRupiah" 
        />

        <SaaSChartGrid :telemetry="telemetry" />

        <WhaleTenantTable 
            :top-tenants="telemetry.top_tenants" 
            :formatRupiah="formatRupiah" 
        />

    </div>
</template>

<style scoped>
::-webkit-scrollbar {
    width: 6px;
    height: 6px;
}
::-webkit-scrollbar-track {
    background: #0B0F19;
}
::-webkit-scrollbar-thumb {
    background: #1e293b;
    border-radius: 10px;
}
::-webkit-scrollbar-thumb:hover {
    background: #334155;
}
</style>