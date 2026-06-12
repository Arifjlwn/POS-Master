<script setup>
import { computed } from 'vue';
import apexchart from 'vue3-apexcharts';

const props = defineProps({
    telemetry: Object
});

const growthChartOptions = computed(() => ({
    chart: { type: 'bar', toolbar: { show: false }, background: 'transparent' },
    theme: { mode: 'dark' },
    colors: ['#6366f1', '#10b981'],
    plotOptions: { bar: { borderRadius: 6, columnWidth: '45%' } },
    dataLabels: { enabled: false },
    stroke: { show: true, width: 2, colors: ['transparent'] },
    xaxis: { 
        categories: props.telemetry.monthly_growth_labels?.length ? props.telemetry.monthly_growth_labels : ['Periode'],
        axisBorder: { show: false },
        axisTicks: { show: false }
    },
    grid: { borderColor: '#1e293b', strokeDashArray: 4 },
    legend: { position: 'top', fontMedium: 'bold' }
}));

const growthChartSeries = computed(() => [
    { name: 'Growth Tenant', data: props.telemetry.monthly_tenant_growth?.length ? props.telemetry.monthly_tenant_growth : [0] },
    { name: 'Growth Transaksi', data: props.telemetry.monthly_trans_growth?.length ? props.telemetry.monthly_trans_growth : [0] }
]);

const churnChartOptions = {
    chart: { type: 'donut', background: 'transparent' },
    theme: { mode: 'dark' },
    colors: ['#10b981', '#f43f5e'], 
    labels: ['Tenant Aktif', 'Tenant Churn / Suspended'],
    legend: { position: 'bottom', fontMedium: 'bold' },
    stroke: { colors: ['#131B2E'] },
    dataLabels: { enabled: false },
    plotOptions: { pie: { donut: { size: '75%' } } }
};

const churnChartSeries = computed(() => {
    const activeTenants = props.telemetry.total_tenants - props.telemetry.churned_tenants;
    return [activeTenants >= 0 ? activeTenants : 0, props.telemetry.churned_tenants || 0];
});
</script>

<template>
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
        <div class="lg:col-span-2 bg-[#131B2E] border border-slate-800 p-5 rounded-[24px] shadow-2xl">
            <div class="mb-4">
                <h3 class="text-xs font-black uppercase tracking-widest text-slate-400">Growth Index & Volume Velocity</h3>
                <p class="text-[10px] text-slate-500 font-bold uppercase tracking-wider mt-0.5">Analisis Akselerasi Pendaftaran Ruko vs Frekuensi Laci Kasir Tenant</p>
            </div>
            <div class="w-full h-64">
                <apexchart type="bar" height="100%" :options="growthChartOptions" :series="growthChartSeries" />
            </div>
        </div>

        <div class="bg-[#131B2E] border border-slate-800 p-5 rounded-[24px] shadow-2xl flex flex-col justify-between">
            <div>
                <h3 class="text-xs font-black uppercase tracking-widest text-slate-400">Churn Rate & Retention Audit</h3>
                <p class="text-[10px] text-slate-500 font-bold uppercase tracking-wider mt-0.5">Deteksi Rasio Ruko Mati / Suspend Terhadap Total Tenant Aktif</p>
            </div>
            <div class="w-full h-48 flex items-center justify-center py-2">
                <apexchart type="donut" width="100%" :options="churnChartOptions" :series="churnChartSeries" />
            </div>
            <div class="border-t border-slate-800/80 pt-3 text-center text-[10px] font-bold uppercase tracking-widest" :class="props.telemetry.churned_tenants > 3 ? 'text-rose-400 animate-pulse' : 'text-slate-500'">
                {{ props.telemetry.churned_tenants }} Tenant Terdeteksi Rusak / Beku
            </div>
        </div>
    </div>
</template>