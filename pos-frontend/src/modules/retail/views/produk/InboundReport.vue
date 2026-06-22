<script setup>
import { onMounted } from 'vue';
import InboundReportFilter from '../../components/inbound/report/InboundReportFilter.vue';
import InboundReportKPI from '../../components/inbound/report/InboundReportKPI.vue';
import InboundReportTable from '../../components/inbound/report/InboundReportTable.vue';
import Sidebar from '../../components/Sidebar.vue';
import { useInboundReport } from '../../composables/useInboundReport.js';

const { isLoading, searchQuery, startDate, endDate, expandedRows, filteredReports, kpiStats, formatRupiah, formatDate, fetchReports, toggleRow, printFaktur } = useInboundReport();

onMounted(() => {
	fetchReports();
});
</script>

<template>
	<Sidebar>
		<div class="p-4 md:p-8 lg:p-10 max-w-7xl mx-auto font-sans bg-[#f8fafc] min-h-screen">
			<div class="mb-8 flex items-center justify-between">
				<div>
					<h1 class="text-2xl lg:text-3xl font-black tracking-tight text-slate-900 uppercase">Laporan Penerimaan Barang</h1>
					<p class="text-xs font-bold text-slate-500 mt-1 uppercase tracking-widest flex items-center gap-2">Audit & Rekapitulasi Faktur Inbound</p>
				</div>
				<button @click="fetchReports" class="w-10 h-10 bg-white border border-slate-200 rounded-xl flex items-center justify-center text-slate-400 hover:text-indigo-600 hover:border-indigo-200 transition-all shadow-sm active:scale-95" title="Refresh Data">
					<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" /></svg>
				</button>
			</div>

			<InboundReportKPI :stats="kpiStats" :formatRupiah="formatRupiah" />

			<InboundReportFilter v-model:searchQuery="searchQuery" v-model:startDate="startDate" v-model:endDate="endDate" />

			<InboundReportTable :reports="filteredReports" :expandedRows="expandedRows" :isLoading="isLoading" :formatRupiah="formatRupiah" :formatDate="formatDate" @toggle-row="toggleRow" @print-faktur="printFaktur" />
		</div>
	</Sidebar>
</template>
