package controllers

import (
	"net/http"
	"pos-backend/config"
	"pos-backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

func GetDashboardReport(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	role, _ := c.Get("role")

	// 1. Gembok Keamanan
	if role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak! Laporan keuangan cuma untuk Owner."})
		return
	}

	storeID := uint(storeIDRaw.(float64))

	// 🚀 2. TANGKAP FILTER DARI VUE
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

	// --- LOGIKA SUMMARY (BERDASARKAN RANGE TANGGAL) ---
	var report struct {
		TotalOmzet         float64 `json:"total_omzet"`
		JumlahTransaksi    int64   `json:"jumlah_transaksi"`
		TotalProdukTerjual float64 `json:"total_produk_terjual"`
		AvgTransaksi       float64 `json:"avg_transaksi"`
	}

	config.DB.Model(&models.Transaction{}).
		Where("store_id = ? AND created_at BETWEEN ? AND ?", storeID, start, end).
		Select("COALESCE(SUM(total_harga), 0) as total_omzet, COUNT(id) as jumlah_transaksi").
		Scan(&report)

	if report.JumlahTransaksi > 0 {
		report.AvgTransaksi = report.TotalOmzet / float64(report.JumlahTransaksi)
	}

	// --- LOGIKA PRODUK TERJUAL ---
	config.DB.Table("transaction_details").
		Joins("JOIN transactions ON transactions.id = transaction_details.transaction_id").
		Where("transactions.store_id = ? AND transactions.created_at BETWEEN ? AND ?", storeID, start, end).
		Select("COALESCE(SUM(transaction_details.kuantitas), 0)").
		Row().Scan(&report.TotalProdukTerjual)

	// --- LOGIKA STOK MENIPIS (lowStockProducts) ---
	var lowStockProducts []models.Product
	config.DB.Where("store_id = ? AND stok < ?", storeID, 10).Find(&lowStockProducts)

	// --- LOGIKA GRAFIK (grafikPenjualan) ---
	type GrafikData struct {
		Tanggal string  `json:"tanggal"`
		Omzet   float64 `json:"omzet"`
	}
	var grafikPenjualan []GrafikData

	days := int(end.Sub(start).Hours() / 24)
	if days <= 0 { days = 1 }
	if days > 31 { days = 31 }

	for i := 0; i < days; i++ {
		tgl := start.AddDate(0, 0, i)
		tglEnd := tgl.Add(24 * time.Hour)

		var dailyOmzet float64
		config.DB.Model(&models.Transaction{}).
			Where("store_id = ? AND created_at BETWEEN ? AND ?", storeID, tgl, tglEnd).
			Select("COALESCE(SUM(total_harga), 0)").
			Row().Scan(&dailyOmzet)

		grafikPenjualan = append(grafikPenjualan, GrafikData{
			Tanggal: tgl.Format("02 Jan"),
			Omzet:   dailyOmzet,
		})
	}

	// --- LOGIKA TOP 5 BEST SELLER (bestSellers) ---
	type BestSeller struct {
		NamaProduk string  `json:"nama_produk"`
		SKU        string  `json:"sku"`
		QtyTerjual int     `json:"qty_terjual"`
		TotalOmzet float64 `json:"total_omzet"`
	}
	var bestSellers []BestSeller

	config.DB.Table("transaction_details").
		Select("products.nama_produk, products.sku, SUM(transaction_details.kuantitas) as qty_terjual, SUM(transaction_details.sub_total) as total_omzet").
		Joins("JOIN transactions ON transactions.id = transaction_details.transaction_id").
		Joins("JOIN products ON products.id = transaction_details.product_id").
		Where("transactions.store_id = ? AND transactions.created_at BETWEEN ? AND ?", storeID, start, end).
		Group("products.nama_produk, products.sku").
		Order("qty_terjual DESC").
		Limit(5).
		Scan(&bestSellers)

	// 3. KIRIM JSON (Satu Nama, Gak Dobel)
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"summary":          report,
			"grafik_penjualan": grafikPenjualan,
			"best_sellers":     bestSellers,
			"low_stock":        lowStockProducts,
		},
	})
}