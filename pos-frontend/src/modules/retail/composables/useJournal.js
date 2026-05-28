import { ref, computed, onMounted, watch } from 'vue';
import { journalService } from '../services/journalService.js';

export function useJournal() {
    // --- STATE DATA ---
    const riwayat = ref([]);
    const isLoading = ref(true);
    const tanggalDipilih = ref(new Date().toISOString().split('T')[0]); 

    // --- STATE PENCARIAN ---
    const searchQuery = ref('');

    // --- STATE MODAL STRUK ---
    const showReceipt = ref(false);
    const selectedTrx = ref(null);

    // 🚀 AMBIL DATA USER DAN TOKO DARI LOCAL STORAGE
    const currentUser = ref(null);
    const currentSession = ref(null);

    const loadSessionData = () => {
        try {
            const user = localStorage.getItem('user');
            if (user) currentUser.value = JSON.parse(user);

            // Coba ambil data store/toko (tergantung lu nyimpennya gimana pas login)
            const store = localStorage.getItem('storeData') || localStorage.getItem('store');
            if (store) currentSession.value = { store: JSON.parse(store) };
        } catch (e) {
            console.error("Gagal parse data session", e);
        }
    };

    const formatRupiah = (angka) => {
        return new Intl.NumberFormat('id-ID').format(angka);
    };

    // 🚀 Ambil data dari Service API
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

    onMounted(() => {
        loadSessionData(); // Panggil fungsi load profil
        fetchRiwayat();
    });
    watch(tanggalDipilih, () => fetchRiwayat());

    // 🚀 FITUR PENCARIAN REALTIME
    const filteredRiwayat = computed(() => {
        if (!searchQuery.value) return riwayat.value;
        const query = searchQuery.value.toLowerCase();
        return riwayat.value.filter(trx => 
            (trx.no_invoice && trx.no_invoice.toLowerCase().includes(query)) ||
            (trx.User?.name && trx.User.name.toLowerCase().includes(query))
        );
    });

    const openReceipt = (trx) => {
        selectedTrx.value = trx;
        showReceipt.value = true;
    };

    // (HAPUS FUNGSI printReceipt KARENA UDAH DI-HANDLE SAMA RECEIPTMODAL KITA)

    return {
        riwayat, isLoading, tanggalDipilih, searchQuery, showReceipt, selectedTrx,
        filteredRiwayat, formatRupiah, fetchRiwayat, openReceipt,
        currentUser, currentSession // 🚀 EXPORT 2 VARIABEL INI BIAR BISA DIPAKE DI VUE
    };
}