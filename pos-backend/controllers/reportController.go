package controllers

import (
	"fmt"
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

	// --- LOGIKA SUMMARY ---
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
	}

	// 1. Ambil Omzet & Produk Terjual
	var resultSummary struct {
		Omzet float64
		Qty   float64
	}
	config.DB.Table("transaction_details").
		Joins("JOIN transactions ON transactions.id = transaction_details.transaction_id").
		Where("transactions.store_id = ? AND transactions.created_at BETWEEN ? AND ?", storeID, start, end).
		Select("COALESCE(SUM(transaction_details.sub_total), 0) as omzet, COALESCE(SUM(transaction_details.kuantitas), 0) as qty").
		Scan(&resultSummary)

	report.TotalOmzet = resultSummary.Omzet
	report.TotalProdukTerjual = resultSummary.Qty

	// 2. Ambil Total Laba
	var totalLaba float64
	err := config.DB.Table("transaction_details").
		Select("COALESCE(SUM(transaction_details.sub_total - (COALESCE(products.harga_modal, 0) * transaction_details.kuantitas)), 0)").
		Joins("JOIN transactions ON transactions.id = transaction_details.transaction_id").
		Joins("LEFT JOIN products ON products.id = transaction_details.product_id").
		Where("transactions.store_id = ? AND transactions.created_at BETWEEN ? AND ?", storeID, start, end).
		Row().Scan(&totalLaba)

	if err != nil {
		fmt.Println("Gagal hitung laba:", err)
	}
	report.TotalLaba = totalLaba

	// 3. Hitung Jumlah Transaksi
	config.DB.Model(&models.Transaction{}).
		Where("store_id = ? AND created_at BETWEEN ? AND ?", storeID, start, end).
		Count(&report.JumlahTransaksi)

	if report.JumlahTransaksi > 0 {
		report.AvgTransaksi = report.TotalOmzet / float64(report.JumlahTransaksi)
	}

	// 🚀 4. Ambil Summary Total Retur & Loss (Waste)
	var resultRetur struct {
		Qty  float64
		Loss float64
	}
	// Perhitungan Kerugian (Loss) = Qty Retur * Harga Modal Produk
	config.DB.Table("product_returns").
		Select("COALESCE(SUM(product_returns.qty), 0) as qty, COALESCE(SUM(product_returns.qty * COALESCE(products.harga_modal, 0)), 0) as loss").
		Joins("LEFT JOIN products ON products.id = product_returns.product_id").
		Where("product_returns.store_id = ? AND product_returns.created_at BETWEEN ? AND ?", storeID, start, end).
		Scan(&resultRetur)

	report.TotalReturQty = resultRetur.Qty
	report.TotalReturLoss = resultRetur.Loss

	// 🚀 5. Ambil Summary Kerugian Stock Opname (HANYA YANG SELISIH MINUS)
	var resultSO struct {
		Qty  float64
		Loss float64
	}
	
	// Kita join 3 tabel: Details (buat selisih & harga), Header (buat filter toko & tanggal), Products (buat harga_modal)
	config.DB.Table("stock_opname_details").
		Select("COALESCE(SUM(ABS(stock_opname_details.selisih)), 0) as qty, COALESCE(SUM(ABS(stock_opname_details.selisih) * COALESCE(products.harga_modal, 0)), 0) as loss").
		Joins("JOIN stock_opnames ON stock_opnames.id = stock_opname_details.opname_id").
		Joins("LEFT JOIN products ON products.id = stock_opname_details.product_id").
		Where("stock_opnames.store_id = ? AND stock_opname_details.selisih < 0 AND stock_opnames.created_at BETWEEN ? AND ?", storeID, start, end).
		Scan(&resultSO)

	report.TotalSOQty = resultSO.Qty
	report.TotalSOLoss = resultSO.Loss

	// --- LOGIKA STOK MENIPIS (lowStockProducts) ---
	var lowStockProducts []models.Product
	config.DB.Where("store_id = ? AND stok < ?", storeID, 10).Find(&lowStockProducts)

	// --- LOGIKA GRAFIK (grafikPenjualan) ---
	type GrafikData struct {
		Tanggal   string  `json:"tanggal"`
		Omzet     float64 `json:"omzet"`
		Laba      float64 `json:"laba"`
		ReturLoss float64 `json:"retur_loss"` // 🚀 Tambahan line merah untuk grafik
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

		var dailyData struct {
			Omzet     float64
			Laba      float64
			ReturLoss float64
		}

		// Ambil Omzet & Laba Harian
		config.DB.Table("transaction_details").
			Select(`
				COALESCE(SUM(transaction_details.sub_total), 0) as omzet,
				COALESCE(SUM(transaction_details.sub_total - (COALESCE(products.harga_modal, 0) * transaction_details.kuantitas)), 0) as laba
			`).
			Joins("JOIN transactions ON transactions.id = transaction_details.transaction_id").
			Joins("LEFT JOIN products ON products.id = transaction_details.product_id").
			Where("transactions.store_id = ? AND transactions.created_at BETWEEN ? AND ?", storeID, tgl, tglEnd).
			Scan(&dailyData)

		// 🚀 Ambil Kerugian Retur Harian
		config.DB.Table("product_returns").
			Select("COALESCE(SUM(product_returns.qty * COALESCE(products.harga_modal, 0)), 0)").
			Joins("LEFT JOIN products ON products.id = product_returns.product_id").
			Where("product_returns.store_id = ? AND product_returns.created_at BETWEEN ? AND ?", storeID, tgl, tglEnd).
			Row().Scan(&dailyData.ReturLoss)

		grafikPenjualan = append(grafikPenjualan, GrafikData{
			Tanggal:   tgl.Format("02 Jan"),
			Omzet:     dailyData.Omzet,
			Laba:      dailyData.Laba,
			ReturLoss: dailyData.ReturLoss, // 🚀 Masuk ke array grafik
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

	// 3. KIRIM JSON
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"summary":          report,
			"grafik_penjualan": grafikPenjualan,
			"best_sellers":     bestSellers,
			"low_stock":        lowStockProducts,
		},
	})
}