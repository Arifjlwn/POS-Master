import api from '../../../api.js';

export const reportService = {
    // 🚀 Tarik ringkasan metrik keuangan dan grafik dari backend
    async getDashboardAnalytics(startDate, endDate) {
        const response = await api.get('/retail/report/dashboard', {
            params: { start_date: startDate, end_date: endDate }
        });
        return response.data;
    }
};