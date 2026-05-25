import api from '../../../api.js';

export const journalService = {
    // 🚀 Ambil jurnal transaksi harian berdasarkan filter tanggal
    async getDailyTransactions(tanggal) {
        return await api.get('/retail/transactions', {
            params: { tanggal }
        });
    }
};