import { ref, computed, onMounted, watch } from 'vue';
import { journalService } from '../services/journalService.js';

export function useJournal() {
    // --- STATE DATA SALES ---
    const riwayat = ref([]);
    const isLoading = ref(true);
    const tanggalDipilih = ref(new Date().toISOString().split('T')[0]); 

    // --- 🚀 STATE DATA CLOSING (BARU) ---
    const riwayatClosing = ref([]); // Nyimpen data closing
    const activeTab = ref('sales'); // 'sales' atau 'closing'
    const selectedClosing = ref(null); // Data closing yang dipilih buat di-print
    const showClosingReceipt = ref(false); // Modal struk closing

    // --- STATE PENCARIAN ---
    const searchQuery = ref('');

    // --- STATE MODAL STRUK SALES ---
    const showReceipt = ref(false);
    const selectedTrx = ref(null);

    // 🚀 AMBIL DATA USER DAN TOKO DARI LOCAL STORAGE
    const currentUser = ref(null);
    const currentSession = ref(null);

    const loadSessionData = () => {
        try {
            const user = localStorage.getItem('user');
            if (user) currentUser.value = JSON.parse(user);

            const store = localStorage.getItem('storeData') || localStorage.getItem('store');
            if (store) currentSession.value = { store: JSON.parse(store) };
        } catch (e) {
            console.error("Gagal parse data session", e);
        }
    };

    const formatRupiah = (angka) => {
        if (!angka && angka !== 0) return '0';
        return new Intl.NumberFormat('id-ID').format(angka);
    };

    // 🚀 AMBIL DATA SALES DARI API
    const fetchRiwayat = async () => {
        isLoading.value = true;
        try {
            const response = await journalService.getDailyTransactions(tanggalDipilih.value);
            riwayat.value = response.data.data || [];
        } catch (error) {
            console.error("Gagal menarik riwayat transaksi:", error);
        } finally {
            isLoading.value = false;
        }
    };

    // 🚀 AMBIL DATA CLOSING DARI API (FUNGSI BARU)
    const fetchRiwayatClosing = async () => {
        try {
            // Asumsi di journalService.js lu bikin fungsi getDailyClosing(date)
            // Kalo belum ada di service-nya, jangan lupa lu buat ya!
            const response = await journalService.getDailyClosing(tanggalDipilih.value);
            riwayatClosing.value = response.data.data || [];
        } catch (error) {
            console.error("Gagal menarik riwayat closing:", error);
        }
    };

    // 🚀 AUTO-LOAD PAS TANGGAL BERUBAH ATAU PERTAMA KALI MOUNT
    onMounted(() => {
        loadSessionData(); 
        fetchRiwayat();
        fetchRiwayatClosing(); // Langsung tarik data closing juga
    });

    watch(tanggalDipilih, () => {
        fetchRiwayat();
        fetchRiwayatClosing(); // Kalo tanggal ganti, update dua-duanya
    });

    // 🚀 FITUR PENCARIAN REALTIME (SALES)
    const filteredRiwayat = computed(() => {
        if (!searchQuery.value) return riwayat.value;
        const query = searchQuery.value.toLowerCase();
        return riwayat.value.filter(trx => 
            (trx.no_invoice && trx.no_invoice.toLowerCase().includes(query)) ||
            (trx.User?.name && trx.User.name.toLowerCase().includes(query))
        );
    });

    // 🚀 BUKA MODAL STRUK SALES
    const openReceipt = (trx) => {
        selectedTrx.value = trx;
        showReceipt.value = true;
    };

    // 🚀 BUKA MODAL STRUK CLOSING (FUNGSI BARU)
    const openClosingReceipt = (closingData) => {
        selectedClosing.value = closingData;
        showClosingReceipt.value = true;
    };

    return {
        riwayat, isLoading, tanggalDipilih, searchQuery, showReceipt, selectedTrx,
        filteredRiwayat, formatRupiah, fetchRiwayat, openReceipt,
        currentUser, currentSession,

        // 🚀 PASTIKAN SEMUA STATE & FUNGSI CLOSING DI-EXPORT!
        activeTab,
        riwayatClosing,
        showClosingReceipt,
        selectedClosing,
        openClosingReceipt,
        fetchRiwayatClosing
    };
}