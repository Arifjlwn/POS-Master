import { ref, watch, nextTick } from 'vue';
import { reportService } from '../services/reportService.js';
import Chart from 'chart.js/auto';

export function useDashboard() {
    const reportData = ref(null);
    const isLoading = ref(true);
    const storeName = ref(localStorage.getItem('storeName') || 'POS UMKM');

    // Canvas DOM References
    const lineChartCanvas = ref(null);
    const pieChartCanvas = ref(null);

    // Chart Instances
    let lineChartInstance = null;
    let pieChartInstance = null;

    // Filter Dates Initial State
    const today = new Date().toISOString().split('T')[0];
    const lastWeek = new Date();
    lastWeek.setDate(lastWeek.getDate() - 6);
    const startDate = ref(lastWeek.toISOString().split('T')[0]);
    const endDate = ref(today);

    // --- UTILITIES ---
    const formatRupiah = (angka) => {
        return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(angka || 0);
    };

    const setQuickFilter = (days) => {
        const end = new Date();
        const start = new Date();
        start.setDate(end.getDate() - days);
        
        endDate.value = end.toISOString().split('T')[0];
        startDate.value = start.toISOString().split('T')[0];
    };

    // --- CHART GRAPHICS INJECTOR ---
    const renderLineChart = (grafikData) => {
        if (!lineChartCanvas.value || !grafikData) return;
        if (lineChartInstance) lineChartInstance.destroy();

        lineChartInstance = new Chart(lineChartCanvas.value, {
            type: 'line',
            data: {
                labels: grafikData.map(d => {
                    const date = new Date(d.tanggal);
                    return date.toLocaleDateString('id-ID', { day: 'numeric', month: 'short' });
                }),
                datasets: [
                    {
                        label: 'Omzet (Rp)',
                        data: grafikData.map(d => d.omzet),
                        borderColor: '#4f46e5',
                        backgroundColor: 'rgba(79, 70, 229, 0.06)',
                        borderWidth: 3,
                        tension: 0.35,
                        fill: true,
                        pointBackgroundColor: '#ffffff',
                        pointBorderColor: '#4f46e5',
                        pointBorderWidth: 2,
                        pointRadius: 4
                    },
                    {
                        label: 'Laba Bersih (Rp)',
                        data: grafikData.map(d => d.laba || 0),
                        borderColor: '#10b981',
                        backgroundColor: 'transparent',
                        borderWidth: 2,
                        borderDash: [5, 5],
                        tension: 0.35,
                        pointRadius: 0
                    },
                    {
                        label: 'Kerugian Retur (Rp)',
                        data: grafikData.map(d => d.retur_loss || 0),
                        borderColor: '#e11d48',
                        backgroundColor: 'rgba(225, 29, 72, 0.1)',
                        borderWidth: 2,
                        tension: 0.35,
                        fill: true,
                        pointRadius: 0
                    }
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                interaction: { mode: 'index', intersect: false },
                plugins: {
                    legend: { display: false },
                    tooltip: {
                        backgroundColor: 'rgba(15, 23, 42, 0.9)',
                        padding: 10,
                        cornerRadius: 8,
                        callbacks: {
                            label: (ctx) => ` ${ctx.dataset.label}: ${formatRupiah(ctx.parsed.y)}`
                        }
                    }
                },
                scales: {
                    y: {
                        beginAtZero: true,
                        grid: { color: '#f1f5f9', borderDash: [4, 4] },
                        ticks: {
                            font: { size: 10, weight: '600' },
                            color: '#94a3b8',
                            callback: (v) => v >= 1000000 ? 'Rp ' + (v / 1000000).toFixed(1) + 'M' : 'Rp ' + v.toLocaleString('id-ID')
                        }
                    },
                    x: { grid: { display: false }, ticks: { font: { size: 10, weight: '600' }, color: '#64748b' } }
                }
            }
        });
    };

    const renderPieChart = (bestSellers) => {
        if (!pieChartCanvas.value || !bestSellers || bestSellers.length === 0) return;
        if (pieChartInstance) pieChartInstance.destroy();

        const top5 = bestSellers.slice(0, 5);

        pieChartInstance = new Chart(pieChartCanvas.value, {
            type: 'doughnut',
            data: {
                labels: top5.map(item => item.nama_produk),
                datasets: [{
                    data: top5.map(item => item.qty_terjual),
                    backgroundColor: ['#4f46e5', '#3b82f6', '#0ea5e9', '#10b981', '#f59e0b'],
                    borderWidth: 0,
                    hoverOffset: 4
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                cutout: '75%',
                plugins: {
                    legend: { display: false },
                    tooltip: {
                        backgroundColor: 'rgba(15, 23, 42, 0.9)',
                        callbacks: { label: (ctx) => ` Terjual: ${ctx.parsed} Pcs` }
                    }
                }
            }
        });
    };

    // --- CORE FETCH DATA ---
    const fetchData = async () => {
        isLoading.value = true;
        try {
            const res = await reportService.getDashboardAnalytics(startDate.value, endDate.value);
            reportData.value = res.data;
            isLoading.value = false;

            // Pastikan Canvas ter-render dulu di DOM baru gambar chart
            nextTick(() => {
                if (reportData.value) {
                    renderLineChart(reportData.value.grafik_penjualan);
                    renderPieChart(reportData.value.best_sellers);
                }
            });
        } catch (error) {
            console.error("Gagal tarik data dashboard:", error);
            isLoading.value = false;
        }
    };

    // Auto hit database saat filter tanggal digeser owner
    watch([startDate, endDate], fetchData);

    return {
        reportData,
        isLoading,
        storeName,
        startDate,
        endDate,
        lineChartCanvas,
        pieChartCanvas,
        formatRupiah,
        setQuickFilter,
        fetchData
    };
}