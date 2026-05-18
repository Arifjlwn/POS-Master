package laundry

import (
	"net/http"
	"pos-backend/config"
	"os"

	"github.com/gin-gonic/gin"
)

type TransaksiReport struct {
	ID            uint    `json:"id"`
	NoInvoice     string  `json:"no_invoice"`
	Tanggal       string  `json:"tanggal"`
	Pelanggan     string  `json:"pelanggan"`
	Total         float64 `json:"total"`
	MetodeBayar   string  `json:"metode_bayar"`
	StatusBayar   string  `json:"status_bayar"`
	BuktiTransfer string  `json:"bukti_transfer"`
}

func AmbilLaporan(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

	var listTransaksi []TransaksiReport

	// Join tabel transactions dan details biar dapet nama pelanggan
	query := config.DB.Table("transactions").
		Select("transactions.id, transactions.no_invoice, transactions.created_at as tanggal, transactions.total_harga as total, transactions.metode_bayar, transactions.status_bayar, transactions.bukti_transfer, MAX(transaction_laundry_details.nama_pelanggan) as pelanggan").
		Joins("left join transaction_laundry_details on transaction_laundry_details.transaction_id = transactions.id").
		Where("transactions.store_id = ?", storeID).
		Group("transactions.id").
		Order("transactions.created_at desc")

	if err := query.Scan(&listTransaksi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data laporan"})
		return
	}

	// Hitung Ringkasan (Bisa di-filter harian nanti, ini totalan dulu)
	var totalTunai, totalQRIS, totalDebit, totalPiutang float64
	
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	for i, t := range listTransaksi {
		if t.StatusBayar == "LUNAS" {
			if t.MetodeBayar == "TUNAI" {
				totalTunai += t.Total
			} else if t.MetodeBayar == "QRIS" {
				totalQRIS += t.Total
			} else if t.MetodeBayar == "DEBIT" {
				totalDebit += t.Total
			}
		} else {
			totalPiutang += t.Total // Yang ngutang / NANTI_AJA
		}

		// Sisipkan domain penuh ke path gambar kalau ada bukti
		if t.BuktiTransfer != "" {
			listTransaksi[i].BuktiTransfer = baseURL + "/" + t.BuktiTransfer
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"ringkasan": gin.H{
			"tunai":   totalTunai,
			"qris":    totalQRIS,
			"debit":   totalDebit,
			"piutang": totalPiutang,
		},
		"transaksi": listTransaksi,
	})
}