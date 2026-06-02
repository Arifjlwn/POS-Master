package delivery

import (
	"net/http"
	"time"

	"pos-backend/models"

	"github.com/gin-gonic/gin"
)

// =================================
// 📊 DASHBOARD ANALYTICS HANDLERS
// =================================

func (h *RetailHandler) GetDashboardReport(c *gin.Context) {
	role := c.MustGet("role").(string)
	if role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak! Laporan keuangan cuma untuk Owner."})
		return
	}

	storeID := uint(c.MustGet("store_id").(float64))
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	now := time.Now()
	location := now.Location()

	start, _ := time.ParseInLocation("2006-01-02", startDateStr, location)
	if startDateStr == "" {
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	}

	end, _ := time.ParseInLocation("2006-01-02", endDateStr, location)
	if endDateStr == "" {
		end = start.Add(24 * time.Hour)
	} else {
		end = end.Add(24 * time.Hour)
	}

	// 🚀 STRUCT DI-UPDATE BUAT NAMPUNG DATA KLAIM & FINAL NETTO
	var report struct {
		TotalOmzet         float64 `json:"total_omzet"`
		TotalLaba          float64 `json:"total_laba"`
		JumlahTransaksi    int64   `json:"jumlah_transaksi"`
		TotalProdukTerjual float64 `json:"total_produk_terjual"`
		AvgTransaksi       float64 `json:"avg_transaksi"`
		TotalReturQty      float64 `json:"total_retur_qty"`
		TotalReturLoss     float64 `json:"total_retur_loss"`

		// Data SO Asli
		TotalSOQty  float64 `json:"total_so_qty"`
		TotalSOLoss float64 `json:"total_so_loss"`

		// Data Klaim Barang Nyempil
		TotalKlaimQty   float64 `json:"total_klaim_qty"`
		TotalKlaimValue float64 `json:"total_klaim_value"`

		// 🚀 HASIL AKHIR (FINAL) SETELAH DIKURANGI KLAIM
		NetSOQty  float64 `json:"net_so_qty"`
		NetSOLoss float64 `json:"net_so_loss"`
	}

	omzet, qty, _ := h.Repo.GetDashboardSummary(storeID, start, end)
	report.TotalOmzet = omzet
	report.TotalProdukTerjual = qty

	laba, _ := h.Repo.GetDashboardLaba(storeID, start, end)
	report.TotalLaba = laba

	db := h.Repo.GetDB()
	db.Model(&models.Transaction{}).Where("store_id = ? AND created_at BETWEEN ? AND ?", storeID, start, end).Count(&report.JumlahTransaksi)

	if report.JumlahTransaksi > 0 {
		report.AvgTransaksi = report.TotalOmzet / float64(report.JumlahTransaksi)
	}

	returQty, returLoss, _ := h.Repo.GetDashboardReturSummary(storeID, start, end)
	report.TotalReturQty = returQty
	report.TotalReturLoss = returLoss

	// 🚀 1. TARIK DATA SO AWAL
	soQty, soLoss, _ := h.Repo.GetDashboardSOSummary(storeID, start, end)
	report.TotalSOQty = soQty
	report.TotalSOLoss = soLoss

	// 🚀 2. TARIK DATA KLAIM YANG UDAH DI-APPROVE
	klaimQty, klaimValue, _ := h.Repo.GetDashboardKlaimSummary(storeID, start, end)
	report.TotalKlaimQty = klaimQty
	report.TotalKlaimValue = klaimValue

	// 🚀 3. HITUNG HASIL FINAL (NETTO RUGI)
	report.NetSOQty = soQty - klaimQty
	if report.NetSOQty < 0 {
		report.NetSOQty = 0
	} // Biar gak minus kalau anomali

	report.NetSOLoss = soLoss - klaimValue
	if report.NetSOLoss < 0 {
		report.NetSOLoss = 0
	} // Biar gak minus kalau anomali

	lowStock, _ := h.Repo.GetLowStockProducts(storeID, 10)

	type GrafikData struct {
		Tanggal   string  `json:"tanggal"`
		Omzet     float64 `json:"omzet"`
		Laba      float64 `json:"laba"`
		ReturLoss float64 `json:"retur_loss"`
	}
	var grafikPenjualan []GrafikData

	days := int(end.Sub(start).Hours() / 24)
	if days <= 0 {
		days = 1
	}
	if days > 31 {
		days = 31
	}

	for i := 0; i < days; i++ {
		tgl := start.AddDate(0, 0, i)
		tglEnd := tgl.Add(24 * time.Hour)

		dailyOmzet, dailyLaba, dailyReturLoss, _ := h.Repo.GetDailySalesReport(storeID, tgl, tglEnd)
		grafikPenjualan = append(grafikPenjualan, GrafikData{
			Tanggal:   tgl.Format("02 Jan"),
			Omzet:     dailyOmzet,
			Laba:      dailyLaba,
			ReturLoss: dailyReturLoss,
		})
	}

	bestSellers, _ := h.Repo.GetTopBestSellers(storeID, start, end)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"summary":          report,
			"grafik_penjualan": grafikPenjualan,
			"best_sellers":     bestSellers,
			"low_stock":        lowStock,
		},
	})
}