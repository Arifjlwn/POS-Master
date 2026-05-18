package laundry

import (
	"net/http"
	"pos-backend/config"

	"github.com/gin-gonic/gin"
)

// Wadah untuk nangkep hasil Join Table
type TrackingResponse struct {
	ID        uint    `json:"id"`
	Invoice   string  `json:"invoice"`
	Pelanggan string  `json:"pelanggan"`
	Whatsapp  string  `json:"whatsapp"`
	Layanan   string  `json:"layanan"`
	BeratKg   float64 `json:"berat_kg"`
	SubTotal  float64 `json:"sub_total"`
	Status    string  `json:"status"`
}

// 🚀 FUNGSI TARIK SEMUA DATA CUCIAN
func AmbilDataTracking(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

	var results []TrackingResponse

	// Join 3 Tabel Sekaligus: Detail Cucian + Transaksi + Produk
	query := config.DB.Table("transaction_laundry_details").
		Select("transaction_laundry_details.id, transactions.no_invoice as invoice, transaction_laundry_details.nama_pelanggan as pelanggan, transaction_laundry_details.no_whatsapp as whatsapp, products.nama_produk as layanan, transaction_laundry_details.berat_kg, transaction_laundry_details.sub_total, transaction_laundry_details.status_cucian as status").
		Joins("left join transactions on transactions.id = transaction_laundry_details.transaction_id").
		Joins("left join products on products.id = transaction_laundry_details.product_id").
		Where("transactions.store_id = ? AND transaction_laundry_details.status_cucian != 'DIAMBIL'", storeID).
		Order("transaction_laundry_details.id desc")

	if err := query.Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data tracking"})
		return
	}

	c.JSON(http.StatusOK, results)
}

// 🚀 FUNGSI UPDATE STATUS (ANTRI -> PROSES -> SELESAI -> DIAMBIL)
func UpdateStatusCucian(c *gin.Context) {
	trxID := c.Param("id")

	var input struct {
		StatusPesanan string `json:"status_pesanan" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format status tidak valid"})
		return
	}

	// 2. Cek update ke database
	if err := config.DB.Table("transactions").Where("id = ?", trxID).Update("status_pesanan", input.StatusPesanan).Error; err != nil {
		
		// 🚀 UBAH BARIS INI BIAR ERROR ASLINYA MUNCUL DI CONSOLE VUE KAMU!
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB Error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status berhasil diupdate!"})
}