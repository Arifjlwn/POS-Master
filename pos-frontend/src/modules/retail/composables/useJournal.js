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

    onMounted(() => fetchRiwayat());
    watch(tanggalDipilih, () => fetchRiwayat());

    // 🚀 FITUR PENCARIAN REALTIME (Tanpa perlu nembak API ulang)
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

    // 🚀 LOGIC INJECT WINDOW PRINT THERMAL STRUK
    const printReceipt = () => {
        const printContent = document.getElementById('print-area').innerHTML;
        const printWindow = window.open('', '_blank', 'width=300,height=600');
        
        printWindow.document.write(`
            <html>
                <head>
                    <title>Cetak Struk Transaksi ${selectedTrx.value?.no_invoice}</title>
                    <style>
                        body { font-family: 'Courier New', Courier, monospace; width: 58mm; margin: 0; padding: 0; font-size: 11px; color: #000; }
                        .text-center { text-align: center; }
                        .font-black { font-weight: 900; }
                        .font-bold { font-weight: bold; }
                        .flex { display: flex; }
                        .justify-between { justify-content: space-between; }
                        .uppercase { text-transform: uppercase; }
                        .border-y { border-top: 1px dashed #000; border-bottom: 1px dashed #000; padding: 5px 0; }
                        .border-b { border-bottom: 1px dashed #000; margin-bottom: 5px; padding-bottom: 5px; }
                        .border-t { border-top: 1px dashed #000; margin-top: 5px; padding-top: 5px; }
                        .truncate { white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
                        .w-full { width: 100%; }
                        .pl-4 { padding-left: 10px; }
                        p { margin: 2px 0; }
                    </style>
                </head>
                <body onload="window.print(); window.close();">
                    \${printContent}
                </body>
            </html>
        `);
        printWindow.document.close();
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
        printReceipt
    };
}