import { ref, onMounted } from 'vue';
import api from '../../../api.js';
import Swal from 'sweetalert2';

export function useStockOpnameReport() {
    const reports = ref([]);
    const isLoading = ref(true);
    const selectedDetail = ref(null); // Ini nanti isinya { so: {...}, klaim: {...} }
    const isApproving = ref(false);

    const role = localStorage.getItem('role') || 'staff';
    const isOwner = role.toLowerCase() === 'owner';

    const fetchReports = async () => {
        isLoading.value = true;
        selectedDetail.value = null; 
        try {
            // 🚀 Cukup 1 endpoint yang balikin Compare Result dari Golang
            const res = await api.get('/retail/stock-opname/history');
            reports.value = res.data.data || [];
        } catch (err) {
            console.error("Gagal ambil history", err);
        } finally {
            isLoading.value = false;
        }
    };

    const showDetail = (report) => {
        selectedDetail.value = report;
    };

    // 🚀 Update fungsi ngitung Valuasi biar dinamis (SO vs KLAIM)
    const calculateLoss = (details, type = 'SO') => {
        if (!details) return 0;
        return details.reduce((acc, item) => {
            const kuantitas = type === 'SO' ? item.selisih : item.qty;
            return acc + (kuantitas * (item.product?.harga_modal || 0));
        }, 0);
    };

    const formatDate = (dateStr) => {
        if (!dateStr) return '';
        return new Intl.DateTimeFormat('id-ID', { 
            weekday: 'long', year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit' 
        }).format(new Date(dateStr));
    };

    // 🚀 Update fungsi approve (Kirim ke ID spesifik SO atau KLAIM)
    // 🚀 TERIMA PARAMETER FILE DARI KOMPONEN
    const approveAudit = async (id, type = 'SO', pdfFile = null) => {
        const confirm = await Swal.fire({
            title: type === 'SO' ? 'Setujui Audit?' : 'Setujui Klaim?',
            text: 'Pastikan file PDF Berita Acara (BAR) sudah benar.',
            icon: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#10b981',
            confirmButtonText: 'Ya, Setujui & Upload!'
        });

        if (!confirm.isConfirmed) return;

        isApproving.value = true;
        try {
            const endpoint = type === 'SO'
                ? `/retail/stock-opname/${id}/approve`
                : `/retail/stock-adjustment/${id}/approve`;

            // 🚀 BIKIN FORMDATA BIAR BISA NGIRIM FILE
            let payload = new FormData();
            if (pdfFile) {
                payload.append('bukti_bar', pdfFile); // Nama parameter file-nya harus klop sama Golang
            }

            // 🚀 PERBAIKAN: HAPUS HEADERS MANUAL! 
            // Cukup lempar 'payload' aja, Axios bakal otomatis nambahin boundary multipart-nya.
            await api.patch(endpoint, payload);
            
            Swal.fire('Berhasil!', 'Persetujuan sukses & Bukti tersimpan.', 'success');
            await fetchReports();
        } catch (err) {
            Swal.fire('Gagal', err.response?.data?.error || 'Sistem error', 'error');
        } finally {
            isApproving.value = false;
        }
    };

    onMounted(fetchReports);

    return { 
        reports, isLoading, selectedDetail, isOwner, isApproving,
        showDetail, calculateLoss, formatDate, approveAudit 
    };
}