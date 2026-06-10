import { computed, ref, watch, nextTick, onUnmounted } from 'vue';
import { reportService } from '../services/reportService.js';
import Chart from 'chart.js/auto';

export function useDashboard() {
	const reportData = ref(null);
	const isLoading = ref(true);
	const errorMessage = ref('');
	const storeName = ref(localStorage.getItem('storeName') || 'POS UMKM');

	const lineChartCanvas = ref(null);
	const pieChartCanvas = ref(null);

	let lineChartInstance = null;
	let pieChartInstance = null;

	// =====================================
	// DATE UTILS
	// =====================================

	const getLocalDate = (d) => {
		const year = d.getFullYear();
		const month = String(d.getMonth() + 1).padStart(2, '0');
		const day = String(d.getDate()).padStart(2, '0');
		return `${year}-${month}-${day}`;
	};

	const today = new Date();
	const lastWeek = new Date();
	lastWeek.setDate(lastWeek.getDate() - 6);

	const startDate = ref(getLocalDate(lastWeek));
	const endDate = ref(getLocalDate(today));

	// =====================================
	// FORMAT UTILS
	// =====================================

	const formatRupiah = (angka) => {
		return new Intl.NumberFormat('id-ID', {
			style: 'currency',
			currency: 'IDR',
			maximumFractionDigits: 0,
		}).format(angka || 0);
	};

	const setQuickFilter = (days) => {
		const end = new Date();
		const start = new Date();
		start.setDate(end.getDate() - days);

		startDate.value = getLocalDate(start);
		endDate.value = getLocalDate(end);
	};

	// =====================================
	// COMPUTED ANALYTICS (FIX MISMATCH PAYLOAD)
	// =====================================

	// 🚀 FIX PROFIT MARGIN: Hitung margin murni berdasarkan Laba Netto baru dari backend Go !
	const profitMargin = computed(() => {
		if (!reportData.value?.summary) return 0;
		const totalOmzet = reportData.value.summary.total_omzet || 0;

		// Ambil laba netto murni, fallback ke gross atau 0
		const totalLabaNetto = reportData.value.summary.total_laba_netto ?? reportData.value.summary.total_laba_gross ?? reportData.value.summary.total_laba ?? 0;

		if (totalOmzet === 0) return 0;
		return ((totalLabaNetto / totalOmzet) * 100).toFixed(1);
	});

	// 🚀 FIX STRUK PER HARI: Langsung baca kalkulasi solid dari level database via Backend payload !
	const strukPerHari = computed(() => {
		if (!reportData.value?.summary) return 0;
		const val = reportData.value.summary.struk_per_hari || 0;
		return Number(val).toFixed(1);
	});

	// =====================================
	// LINE CHART (TREND ANALYTICS)
	// =====================================

	const renderLineChart = (grafikData = []) => {
		if (!lineChartCanvas.value) return;

		// SAFE MEMORY DESTRUCTION: Hancurkan instance lama biar gak stack dan flicker pas di-hover !
		if (lineChartInstance) {
			lineChartInstance.destroy();
			lineChartInstance = null;
		}

		if (!grafikData.length) return;

		lineChartInstance = new Chart(lineChartCanvas.value, {
			type: 'line',
			data: {
				labels: grafikData.map((d) =>
					new Date(d.tanggal).toLocaleDateString('id-ID', {
						day: 'numeric',
						month: 'short',
					})
				),
				datasets: [
					{
						label: 'Omzet (Rp)',
						data: grafikData.map((item) => item.omzet || 0),
						borderColor: '#4f46e5',
						backgroundColor: 'rgba(79,70,229,0.06)',
						borderWidth: 3,
						tension: 0.35,
						fill: true,
						pointBackgroundColor: '#ffffff',
						pointBorderColor: '#4f46e5',
						pointBorderWidth: 2,
						pointRadius: 4,
						pointHoverRadius: 6,
					},
					{
						label: 'Laba Bersih (Rp)',
						data: grafikData.map((item) => item.laba || 0),
						borderColor: '#10b981',
						backgroundColor: 'transparent',
						borderWidth: 2,
						borderDash: [5, 5],
						tension: 0.35,
						pointRadius: 0,
						pointHoverRadius: 5,
						pointHoverBackgroundColor: '#10b981',
					},
					{
						label: 'Kerugian Retur (Rp)',
						data: grafikData.map((item) => item.retur_loss || 0),
						borderColor: '#e11d48',
						backgroundColor: 'rgba(225,29,72,0.06)',
						borderWidth: 2,
						tension: 0.35,
						fill: true,
						pointRadius: 0,
					},
				],
			},
			options: {
				responsive: true,
				maintainAspectRatio: false,
				interaction: {
					mode: 'index',
					intersect: false,
				},
				plugins: {
					legend: { display: false },
					tooltip: {
						backgroundColor: 'rgba(15,23,42,0.95)',
						padding: 12,
						cornerRadius: 12,
						callbacks: {
							label: (ctx) => `${ctx.dataset.label}: ${formatRupiah(ctx.parsed.y)}`,
						},
					},
				},
				scales: {
					y: {
						beginAtZero: true,
						grid: {
							color: '#f1f5f9',
							borderDash: [4, 4],
						},
						ticks: {
							color: '#94a3b8',
							callback: (value) => {
								if (Math.abs(value) >= 1000000) {
									return `Rp ${(value / 1000000).toFixed(1)}M`;
								}
								return `Rp ${value.toLocaleString('id-ID')}`;
							},
						},
					},
					x: {
						grid: { display: false },
						ticks: { color: '#64748b' },
					},
				},
			},
		});
	};

	// =====================================
	// PIE CHART (TOP 5 LEADERBOARD)
	// =====================================

	const renderPieChart = (bestSellers = []) => {
		if (!pieChartCanvas.value) return;

		if (pieChartInstance) {
			pieChartInstance.destroy();
			pieChartInstance = null;
		}

		if (!bestSellers.length) return;
		const top5 = bestSellers.slice(0, 5);

		pieChartInstance = new Chart(pieChartCanvas.value, {
			type: 'doughnut',
			data: {
				labels: top5.map((item) => item.nama_produk),
				datasets: [
					{
						data: top5.map((item) => item.qty_terjual),
						backgroundColor: ['#4f46e5', '#3b82f6', '#0ea5e9', '#10b981', '#f59e0b'],
						borderWidth: 0,
						hoverOffset: 6,
					},
				],
			},
			options: {
				responsive: true,
				maintainAspectRatio: false,
				cutout: '75%',
				plugins: {
					legend: { display: false },
					tooltip: {
						backgroundColor: 'rgba(15,23,42,0.95)',
						padding: 10,
						cornerRadius: 10,
						callbacks: {
							label: (ctx) => {
								const item = top5[ctx.dataIndex];
								return ` Terjual: ${ctx.parsed} ${item.satuan_dasar || 'Pcs'}`;
							},
						},
					},
				},
			},
		});
	};

	// =====================================
	// FETCH DATA ENGINE
	// =====================================

	const fetchData = async () => {
		try {
			isLoading.value = true;
			errorMessage.value = '';

			const res = await reportService.getDashboardAnalytics(startDate.value, endDate.value);

			// Simpan data array terstruktur murni
			reportData.value = res.data;

			isLoading.value = false;
			await nextTick();

			if (reportData.value) {
				renderLineChart(reportData.value.grafik_penjualan || []);
				renderPieChart(reportData.value.best_sellers || []);
			}
		} catch (error) {
			console.error('Gagal tarik data dashboard:', error);
			errorMessage.value = 'Gagal memuat data laporan. Silakan coba lagi .';
			isLoading.value = false;
		}
	};

	// Tracking filter realtime
	watch([startDate, endDate], () => {
		fetchData();
	});

	onUnmounted(() => {
		if (lineChartInstance) lineChartInstance.destroy();
		if (pieChartInstance) pieChartInstance.destroy();
	});

	return {
		reportData,
		isLoading,
		errorMessage,
		storeName,
		startDate,
		endDate,
		lineChartCanvas,
		pieChartCanvas,
		profitMargin,
		strukPerHari,
		formatRupiah,
		setQuickFilter,
		fetchData,
	};
}
