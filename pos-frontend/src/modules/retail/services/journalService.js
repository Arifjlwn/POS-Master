import api from '../../../api.js';

export const journalService = {
    // 🚀 Ambil jurnal transaksi harian (Sales) berdasarkan filter tanggal
    async getDailyTransactions(tanggal) {
        return await api.get('/retail/transactions', {
            params: { tanggal: tanggal }
        });
    },

    // 🚀 Ambil riwayat tutup shift (Closing) berdasarkan filter tanggal
    async getDailyClosing(tanggal) {
        return await api.get('/retail/journal/closing', {
            // Catatan: Kalo Golang lu nangkepnya pake kata "date", ganti key "tanggal" di bawah jadi "date: tanggal"
            params: { tanggal: tanggal } 
        });
    }
};