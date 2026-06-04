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

	var storeID uint
	storeIDRaw := c.MustGet("store_id")
	if val, ok := storeIDRaw.(float64); ok { storeID = uint(val) } else if val, ok := storeIDRaw.(uint); ok { storeID = val }

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

	var report struct {
		TotalOmzet         float64 `json:"total_omzet"`
		TotalLaba          float64 `json:"total_laba"`
		JumlahTransaksi    int64   `json:"jumlah_transaksi"`
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

	soQty, soLoss, _ := h.Repo.GetDashboardSOSummary(storeID, start, end)
	report.TotalSOQty = soQty
	report.TotalSOLoss = soLoss

	klaimQty, klaimValue, _ := h.Repo.GetDashboardKlaimSummary(storeID, start, end)
	report.TotalKlaimQty = klaimQty
	report.TotalKlaimValue = klaimValue

	report.NetSOQty = soQty - klaimQty
	if report.NetSOQty < 0 { report.NetSOQty = 0 }

	report.NetSOLoss = soLoss - klaimValue
	if report.NetSOLoss < 0 { report.NetSOLoss = 0 }

	lowStock, _ := h.Repo.GetLowStockProducts(storeID, 10)
	bestSellers, _ := h.Repo.GetTopBestSellers(storeID, start, end)

	type GrafikData struct {
		Tanggal   string  `json:"tanggal"`
		Omzet     float64 `json:"omzet"`
		Laba      float64 `json:"laba"`
		ReturLoss float64 `json:"retur_loss"`
	}

	// 🚀 SANGAR: Sikat semua data harian pake single query agregat via repo bray (Anti N+1 loop)
	grafikPenjualan, err := h.Repo.GetAggregatedDailySales(storeID, start, end)
	if err != nil || len(grafikPenjualan) == 0 {
		grafikPenjualan = make([]map[string]interface{}, 0)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"summary":          report,
			"grafik_penjualan": grafikPenjualan,
			"best_sellers":     bestSellers,
			"low_stock":        lowStock,
		},
	})
}