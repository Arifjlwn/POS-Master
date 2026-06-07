import { ref, computed, onMounted, watch } from 'vue';
import { journalService } from '../services/journalService.js';

export function useJournal() {
    // --- STATE DATA SALES & CLOSING ---
    const riwayat = ref([]);
    const riwayatClosing = ref([]); 
    const isLoading = ref(true);
    
    // 🚀 FIX TIMEZONE SHIFTING: Gunakan penanggalan lokal (WIB), bukan UTC!
    const dapatkanTanggalLokal = () => {
        const d = new Date();
        const offset = d.getTimezoneOffset();
        const localDate = new Date(d.getTime() - (offset * 60 * 1000));
        return localDate.toISOString().split('T')[0];
    };
    
    const tanggalDipilih = ref(dapatkanTanggalLokal()); 

    const activeTab = ref('sales'); // 'sales' atau 'closing'
    const selectedClosing = ref(null); // Data closing yang dipilih buat di-print
    const showClosingReceipt = ref(false); // Modal struk closing

    // --- STATE PENCARIAN ---
    const searchQuery = ref('');

    // --- STATE MODAL STRUK SALES ---
    const showReceipt = ref(false);
    const selectedTrx = ref(null);

    // --- STATE DATA USER DAN TOKO ---
    const currentUser = ref(null);
    const currentSession = ref(null);

    const loadSessionData = () => {
        try {
            const user = localStorage.getItem('user');
            if (user) currentUser.value = JSON.parse(user);

            const store = localStorage.getItem('storeData') || localStorage.getItem('store');
            if (store) currentSession.value = { store: JSON.parse(store) };
        } catch (e) {
            console.error("Gagal parse data session:", e);
        }
    };

    const formatRupiah = (angka) => {
        if (!angka && angka !== 0) return '0';
        return new Intl.NumberFormat('id-ID').format(angka);
    };

    // 🚀 AMBIL DATA SALES DARI API
    const fetchRiwayat = async () => {
        try {
            const response = await journalService.getDailyTransactions(tanggalDipilih.value);
            riwayat.value = response.data?.data || [];
        } catch (error) {
            console.error("Gagal menarik riwayat transaksi:", error);
            riwayat.value = [];
        }
    };

    // 🚀 AMBIL DATA CLOSING DARI API
    const fetchRiwayatClosing = async () => {
        try {
            const response = await journalService.getDailyClosing(tanggalDipilih.value);
            riwayatClosing.value = response.data?.data || [];
        } catch (error) {
            console.error("Gagal menarik riwayat closing:", error);
            riwayatClosing.value = [];
        }
    };

    // 🚀 AGREGASI BUNDLE LOADING AGAR SINKRON
    const fetchSemuaData = async () => {
        isLoading.value = true;
        await Promise.allSettled([
            fetchRiwayat(),
            fetchRiwayatClosing()
        ]);
        isLoading.value = false;
    };

    // 🚀 AUTO-LOAD PAS PERTAMA KALI MOUNT & PAS TANGGAL BERUBAH
    onMounted(() => {
        loadSessionData(); 
        fetchSemuaData();
    });

    watch(tanggalDipilih, () => {
        fetchSemuaData();
    });

    // 🚀 FITUR PENCARIAN REALTIME (SALES) - Proteksi Case-Sensitive Object
    const filteredRiwayat = computed(() => {
        if (!riwayat.value) return [];
        if (!searchQuery.value) return riwayat.value;
        const query = searchQuery.value.toLowerCase();
        
        return riwayat.value.filter(trx => {
            const invoiceMatch = trx.no_invoice && trx.no_invoice.toLowerCase().includes(query);
            
            // Cek properti 'User' kapital atau 'user' kecil dari backend
            const userObj = trx.User || trx.user;
            const userMatch = userObj?.name && userObj.name.toLowerCase().includes(query);
            
            return invoiceMatch || userMatch;
        });
    });

    // 🚀 BUKA MODAL STRUK SALES
    const openReceipt = (trx) => {
        selectedTrx.value = trx;
        showReceipt.value = true;
    };

    // 🚀 BUKA MODAL STRUK CLOSING
    const openClosingReceipt = (closingData) => {
        selectedClosing.value = closingData;
        showClosingReceipt.value = true;
    };

    return {
        riwayat, 
        isLoading, 
        tanggalDipilih, 
        searchQuery, 
        showReceipt, 
        selectedTrx,
        filteredRiwayat, 
        formatRupiah, 
        fetchRiwayat, 
        openReceipt,
        currentUser, 
        currentSession,

        // STATE & FUNGSI CLOSING
        activeTab,
        riwayatClosing,
        showClosingReceipt,
        selectedClosing,
        openClosingReceipt,
        fetchRiwayatClosing,
        fetchSemuaData
    };
}