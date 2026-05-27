import { ref, onMounted, watch } from 'vue';
import api from '../../../api.js';
import Swal from 'sweetalert2';

export function useStockOpnameReport() {
    const reportType = ref('SO'); // 🚀 'SO' atau 'KLAIM'
    const reports = ref([]);
    const isLoading = ref(true);
    const selectedDetail = ref(null);
    const isApproving = ref(false);

    const role = localStorage.getItem('role') || 'staff';
    const isOwner = role.toLowerCase() === 'owner';

    const fetchReports = async () => {
        isLoading.value = true;
        selectedDetail.value = null; // Reset detail pas pindah tab
        try {
            // 🚀 DINAMIS: Tembak url sesuai tab yang dipilih Owner
            const endpoint = reportType.value === 'SO' 
                ? '/retail/stock-opname/history' 
                : '/retail/stock-adjustment/history';
                
            const res = await api.get(endpoint);
            reports.value = res.data.data || [];
        } catch (err) {
            console.error("Gagal ambil history", err);
        } finally {
            isLoading.value = false;
        }
    };

    // Awasi kalau Owner mindah tab, langsung otomatis tarik data baru
    watch(reportType, () => { fetchReports(); });

    const showDetail = (report) => {
        selectedDetail.value = report;
    };

    const calculateLoss = (details) => {
        if (!details) return 0;
        // Kalau SO pake .selisih, kalau KLAIM barang ketemu pake .qty
        return details.reduce((acc, item) => {
            const kuantitas = reportType.value === 'SO' ? item.selisih : item.qty;
            return acc + (kuantitas * (item.product?.harga_modal || 0));
        }, 0);
    };

    const formatDate = (dateStr) => {
        if (!dateStr) return '';
        return new Intl.DateTimeFormat('id-ID', { 
            weekday: 'long', year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit' 
        }).format(new Date(dateStr));
    };

    // 🚀 APPROVE DUA JALUR
    const approveAudit = async (id) => {
        const title = reportType.value === 'SO' ? 'Setujui Audit?' : 'Setujui Klaim Barang?';
        const text = reportType.value === 'SO' 
            ? 'Stok master akan DITIMPA sesuai hasil fisik audit.' 
            : 'Stok master akan DITAMBAH sesuai barang nyempil yang ditemukan.';

        const confirm = await Swal.fire({
            title: title,
            text: text,
            icon: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#10b981',
            confirmButtonText: 'Ya, Setujui!'
        });

        if (!confirm.isConfirmed) return;

        isApproving.value = true;
        try {
            const endpoint = reportType.value === 'SO'
                ? `/retail/stock-opname/${id}/approve`
                : `/retail/stock-adjustment/${id}/approve`;

            await api.patch(endpoint);
            Swal.fire('Berhasil!', 'Persetujuan sukses diproses.', 'success');
            await fetchReports();
        } catch (err) {
            Swal.fire('Gagal', err.response?.data?.error || 'Sistem error', 'error');
        } finally {
            isApproving.value = false;
        }
    };

    onMounted(fetchReports);

    return { 
        reportType, reports, isLoading, selectedDetail, isOwner, isApproving,
        showDetail, calculateLoss, formatDate, approveAudit 
    };
}