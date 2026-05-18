package laundry

import (
	"net/http"
	"os"
	"pos-backend/config"
	"pos-backend/models"
	"strings"

	"github.com/gin-gonic/gin"
)

// Struct payload untuk update setting
type UpdateSettingInput struct {
	NamaToko      string `json:"nama_toko" binding:"required"`
	Telepon       string `json:"telepon" binding:"required"`
	Alamat        string `json:"alamat" binding:"required"`
	PaymentType   string `json:"payment_type" binding:"required"`
	QrisBase64    string `json:"qris_base64"` // Opsional, diisi kalau owner upload QRIS baru
	ReceiptFooter string `json:"receipt_footer"`
}

// 🚀 1. FUNGSI AMBIL DATA SETTING TOKO
func GetSettingToko(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

	var store models.Store
	if err := config.DB.First(&store, storeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data toko tidak ditemukan"})
		return
	}

	// Masukkan domain URL penuh untuk gambar QRIS agar Vue bisa nampilin gambarnya
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	if store.QrisImage != "" {
		store.QrisImage = baseURL + "/" + store.QrisImage
	}

	c.JSON(http.StatusOK, store)
}

// 🚀 2. FUNGSI UPDATE DATA SETTING TOKO
func UpdateSettingToko(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

	var input UpdateSettingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data pengaturan tidak valid"})
		return
	}

	var store models.Store
	if err := config.DB.First(&store, storeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Toko tidak ditemukan"})
		return
	}

	// Jika ada upload QRIS barcode baru berbentuk Base64, kita simpan!
	if input.PaymentType == "PRIBADI" && input.QrisBase64 != "" && !strings.HasPrefix(input.QrisBase64, "http") {
		// Manfaatkan fungsi SimpanGambarBase64 yang udah kita buat di kasir_laundry.go kemarin
		qrisPath, err := SimpanGambarBase64(input.QrisBase64, "public/uploads/qris_toko", "store_"+c.Param("store_id")+"_qris.jpg")
		if err == nil {
			store.QrisImage = qrisPath
		}
	}

	// Update field lainnya
	store.NamaToko = input.NamaToko
	store.Telepon = input.Telepon
	store.Alamat = input.Alamat
	store.PaymentType = input.PaymentType
	store.ReceiptFooter = input.ReceiptFooter

	if err := config.DB.Save(&store).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui pengaturan toko"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "sukses", "message": "Pengaturan toko berhasil diperbarui!"})
}