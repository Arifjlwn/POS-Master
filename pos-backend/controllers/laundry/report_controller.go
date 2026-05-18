package laundry

import (
	"net/http"
	"pos-backend/config"
	"pos-backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

// 🚀 STRUCT KUSTOM: Jembatan penghubung antara Tabel Transaksi & Detail
type TransactionReportResponse struct {
	models.Transaction
	Invoice   string  `json:"invoice"`
	Pelanggan string  `json:"pelanggan"`
	Whatsapp  string  `json:"whatsapp"`
	Layanan   string  `json:"layanan"`
	BeratKg   float64 `json:"berat_kg"`
	SubTotal  float64 `json:"sub_total"`
	EstimasiWaktu time.Time `json:"estimasi_waktu"`
}

// 🚀 FUNGSI AMBIL SEMUA DATA LAPORAN & STATUS CUCIAN
func GetLaporan(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

	var transactions []models.Transaction
	if err := config.DB.Where("store_id = ?", storeID).Order("created_at desc").Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data laporan"})
		return
	}

	var reportData []TransactionReportResponse
	var tunai, qris, debit, piutang, omset float64
	totalOrder := len(transactions)

	for _, trx := range transactions {
		// 1. Hitung-hitungan Duit & Omset
		if trx.StatusBayar == "BELUM_LUNAS" {
			piutang += trx.TotalHarga
		} else {
			omset += trx.TotalHarga
			switch trx.MetodeBayar {
			case "TUNAI":
				tunai += trx.TotalHarga
			case "QRIS":
				qris += trx.TotalHarga
			case "DEBIT":
				debit += trx.TotalHarga
			}
		}

		// 2. 🚀 CARI NAMA PELANGGAN DI TABEL DETAIL
		var detail models.TransactionLaundryDetail
		config.DB.Where("transaction_id = ?", trx.ID).First(&detail)

		// 3. 🚀 CARI NAMA LAYANAN (Misal: Cuci Komplit, Setrika)
		var product models.Product
		layananName := "Paket Laundry"
		if err := config.DB.Where("id = ?", detail.ProductID).First(&product).Error; err == nil {
			layananName = product.NamaProduk
		}

		// 4. BUNGKUS JADI SATU KEMASAN UNTUK VUE
		reportData = append(reportData, TransactionReportResponse{
			Transaction: trx,
			Invoice:     trx.NoInvoice,
			Pelanggan:   detail.NamaPelanggan,
			Whatsapp:    detail.NoWhatsapp,
			Layanan:     layananName,
			BeratKg:     detail.BeratKg,
			SubTotal:    detail.SubTotal,
			EstimasiWaktu: detail.EstimasiWaktu,
		})
	}

	avg := 0.0
	if totalOrder > 0 {
		avg = omset / float64(totalOrder)
	}

	c.JSON(http.StatusOK, gin.H{
		"ringkasan": gin.H{
			"total_omset": omset,
			"total_order": totalOrder,
			"rata_rata":   avg,
			"tunai":       tunai,
			"qris":        qris,
			"debit":       debit,
			"piutang":     piutang,
		},
		"transaksi": reportData, // 👈 Kirim bungkusan yang udah lengkap ada namanya!
	})
}