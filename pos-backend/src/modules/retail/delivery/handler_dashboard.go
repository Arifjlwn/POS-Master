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
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak! Laporan keuangan eksklusif untuk Owner."})
		return
	}

	var storeID uint
	storeIDRaw := c.MustGet("store_id")
	if val, ok := storeIDRaw.(float64); ok {
		storeID = uint(val)
	} else if val, ok := storeIDRaw.(uint); ok {
		storeID = val
	}

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	now := time.Now()
	location := now.Location()

	// 1. Kalibrasi Waktu Mulai (00:00:00)
	start, errStart := time.ParseInLocation("2006-01-02", startDateStr, location)
	if startDateStr == "" || errStart != nil {
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	} else {
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, location)
	}

	// 2. 🚀 FIX TIME WINDOW TRAP: Kunci tepat di jam 23:59:59 hari itu juga ! Anti-lewat hari!
	end, errEnd := time.ParseInLocation("2006-01-02", endDateStr, location)
	if endDateStr == "" || errEnd != nil {
		end = time.Date(start.Year(), start.Month(), start.Day(), 23, 59, 59, 999999999, location)
	} else {
		end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, location)
	}

	var report struct {
		TotalOmzet         float64 `json:"total_omzet"`
		TotalLabaGross     float64 `json:"total_laba_gross"` // Kita perjelas ini laba kotor penjualan
		TotalLabaNetto     float64 `json:"total_laba_netto"` // 🚀 INDIKATOR SAKTI: Laba bersih riil Owner!
		JumlahTransaksi    int64   `json:"jumlah_transaksi"`
		StrukPerHari       float64 `json:"struk_per_hari"`
		TotalProdukTerjual float64 `json:"total_produk_terjual"`
		AvgTransaksi       float64 `json:"avg_transaksi"`
		TotalReturQty      float64 `json:"total_retur_qty"`
		TotalReturLoss     float64 `json:"total_retur_loss"`
		TotalSOQty         float64 `json:"total_so_qty"`
		TotalSOLoss        float64 `json:"total_so_loss"`
		TotalKlaimQty      float64 `json:"total_klaim_qty"`
		TotalKlaimValue    float64 `json:"total_klaim_value"`
		NetSOQty           float64 `json:"net_so_qty"`
		NetSOLoss          float64 `json:"net_so_loss"`
	}

	// Tarik data dasar omzet & produk terjual
	omzet, qty, _ := h.Repo.GetDashboardSummary(storeID, start, end)
	report.TotalOmzet = omzet
	report.TotalProdukTerjual = qty

	// Tarik profit kotor awal penjualan
	labaGross, _ := h.Repo.GetDashboardLaba(storeID, start, end)
	report.TotalLabaGross = labaGross

	// Hitung total transaksi invoice via GORM
	db := h.Repo.GetDB()
	db.Model(&models.Transaction{}).Where("store_id = ? AND created_at BETWEEN ? AND ?", storeID, start, end).Count(&report.JumlahTransaksi)

	// Hitung rata-rata nilai belanja per struk
	if report.JumlahTransaksi > 0 {
		report.AvgTransaksi = report.TotalOmzet / float64(report.JumlahTransaksi)
	}

	// Hitung pembagi rata-rata transaksi harian
	jumlahHari := int(end.Sub(start).Hours()/24) + 1 // +1 biar hari yang sama dihitung 1 hari penuh
	if jumlahHari <= 0 {
		jumlahHari = 1
	}
	report.StrukPerHari = float64(report.JumlahTransaksi) / float64(jumlahHari)

	// Tarik rangkuman pembuangan waste & retur rusak
	returQty, returLoss, _ := h.Repo.GetDashboardReturSummary(storeID, start, end)
	report.TotalReturQty = returQty
	report.TotalReturLoss = returLoss

	// Tarik kerugian awal audit laci (SO)
	soQty, soLoss, _ := h.Repo.GetDashboardSOSummary(storeID, start, end)
	report.TotalSOQty = soQty
	report.TotalSOLoss = soLoss

	// Tarik klaim temuan barang
	klaimQty, klaimValue, _ := h.Repo.GetDashboardKlaimSummary(storeID, start, end)
	report.TotalKlaimQty = klaimQty
	report.TotalKlaimValue = klaimValue

	// Hitung selisih bersih kebocoran barang audit
	report.NetSOQty = soQty - klaimQty
	if report.NetSOQty < 0 {
		report.NetSOQty = 0
	}

	report.NetSOLoss = soLoss - klaimValue
	if report.NetSOLoss < 0 {
		report.NetSOLoss = 0
	}

	// 🚀 RUMUS EKONOMI MUTLAK: Laba Bersih = Laba Kotor - Kerugian Retur/Waste - Kerugian Bersih SO
	report.TotalLabaNetto = report.TotalLabaGross - report.TotalReturLoss - report.NetSOLoss

	// Pull data produk limit menipis & produk terlaris
	lowStock, _ := h.Repo.GetLowStockProducts(storeID, 10)
	bestSellers, _ := h.Repo.GetTopBestSellers(storeID, start, end)

	type GrafikData struct {
		Tanggal   string  `json:"tanggal"`
		Omzet     float64 `json:"omzet"`
		Laba      float64 `json:"laba"`
		ReturLoss float64 `json:"retur_loss"`
	}

	// Sikat data grafik penjualan harian agregat
	grafikPenjualan, err := h.Repo.GetAggregatedDailySales(storeID, start, end)

	if err != nil || len(grafikPenjualan) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"summary":          report,
				"grafik_penjualan": []GrafikData{},
				"best_sellers":     bestSellers,
				"low_stock":        lowStock,
				"role_info":        role,
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"summary":          report,
			"grafik_penjualan": grafikPenjualan,
			"best_sellers":     bestSellers,
			"low_stock":        lowStock,
			"role_info":        role,
		},
	})
}
