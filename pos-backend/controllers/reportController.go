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

	// 1. Gembok! Cuma Bos yang boleh lihat laporan keuangan
	if role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak! Laporan keuangan cuma untuk Owner."})
		return
	}

	storeID := uint(storeIDRaw.(float64))

	// 2. Tentukan rentang waktu "Hari Ini" (00:00:00 sampai 23:59:59)
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	// --- LOGIKA HITUNG OMZET & JUMLAH TRANSAKSI (HARI INI) ---
	var report struct {
		OmzetHariIni       float64 `json:"omzet_hari_ini"`
		JumlahTransaksi    int64   `json:"jumlah_transaksi"`
		TotalProdukTerjual float64 `json:"total_produk_terjual"`
	}

	config.DB.Model(&models.Transaction{}).
		Where("store_id = ? AND created_at BETWEEN ? AND ?", storeID, startOfDay, endOfDay).
		Select("COALESCE(SUM(total_harga), 0) as omzet_hari_ini, COUNT(id) as jumlah_transaksi").
		Scan(&report)

	// --- LOGIKA HITUNG TOTAL PRODUK TERJUAL (HARI INI) ---
	config.DB.Table("transaction_details").
		Joins("JOIN transactions ON transactions.id = transaction_details.transaction_id").
		Where("transactions.store_id = ? AND transactions.created_at BETWEEN ? AND ?", storeID, startOfDay, endOfDay).
		Select("COALESCE(SUM(transaction_details.kuantitas), 0)").
		Row().Scan(&report.TotalProdukTerjual)

	// --- LOGIKA STOK MENIPIS (ALERT) ---
	var lowStockProducts []models.Product
	config.DB.Where("store_id = ? AND stok < ?", storeID, 10).
		Find(&lowStockProducts)

	// 🚀 --- LOGIKA GRAFIK 7 HARI TERAKHIR ---
	type GrafikData struct {
		Tanggal string  `json:"tanggal"`
		Omzet   float64 `json:"omzet"`
	}
	var grafik7Hari []GrafikData

	// Looping mundur dari 6 hari yang lalu sampai hari ini
	for i := 6; i >= 0; i-- {
		targetDate := startOfDay.AddDate(0, 0, -i)
		targetDateEnd := targetDate.Add(24 * time.Hour)

		var dailyOmzet float64
		config.DB.Model(&models.Transaction{}).
			Where("store_id = ? AND created_at BETWEEN ? AND ?", storeID, targetDate, targetDateEnd).
			Select("COALESCE(SUM(total_harga), 0)").
			Row().Scan(&dailyOmzet)

		grafik7Hari = append(grafik7Hari, GrafikData{
			Tanggal: targetDate.Format("02 Jan"), // Format: "12 May"
			Omzet:   dailyOmzet,
		})
	}

	// 🚀 --- LOGIKA TOP 5 BEST SELLER (BULAN INI) ---
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

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
		Where("transactions.store_id = ? AND transactions.created_at BETWEEN ? AND ?", storeID, startOfMonth, endOfMonth).
		Group("transaction_details.product_id, products.nama_produk, products.sku").
		Order("qty_terjual DESC").
		Limit(5).
		Scan(&bestSellers)

	// 3. Kirim hasil laporan lengkap ke Frontend
	c.JSON(http.StatusOK, gin.H{
		"message": "Data laporan berhasil ditarik! 📊",
		"data": gin.H{
			"summary":       report,
			"low_stock":     lowStockProducts,
			"grafik_7_hari": grafik7Hari, // Data buat Chart.js
			"best_sellers":  bestSellers, // Data buat Tabel Ranking
		},
	})
}