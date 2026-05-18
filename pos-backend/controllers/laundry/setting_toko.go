package laundry

import (
	"net/http"
	"os"
	"pos-backend/config"
	"pos-backend/models"
	"strings"

	"github.com/gin-gonic/gin"
)

// Struct payload (Bersih dari urusan kasir)
type UpdateSettingInput struct {
	NamaToko      string `json:"nama_toko" binding:"required"`
	Telepon       string `json:"telepon" binding:"required"`
	Alamat        string `json:"alamat" binding:"required"`
	PaymentType   string `json:"payment_type" binding:"required"`
	QrisBase64    string `json:"qris_base64"`
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

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	if store.QrisImage != "" {
		store.QrisImage = baseURL + "/" + store.QrisImage
	}

	c.JSON(http.StatusOK, gin.H{
		"id":             store.ID,
		"nama_toko":      store.NamaToko,
		"business_type":  store.BusinessType,
		"telepon":        store.Telepon,
		"alamat":         store.Alamat,
		"payment_type":   store.PaymentType,
		"qris_image":     store.QrisImage,
		"receipt_footer": store.ReceiptFooter,
	})
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

	tx := config.DB.Begin()

	var store models.Store
	if err := tx.First(&store, storeID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Toko tidak ditemukan"})
		return
	}

	// Proses Simpan Gambar QRIS Base64
	if input.PaymentType == "PRIBADI" && input.QrisBase64 != "" && !strings.HasPrefix(input.QrisBase64, "http") {
		qrisPath, err := SimpanGambarBase64(input.QrisBase64, "public/uploads/qris_toko", "store_qris.jpg")
		if err == nil {
			store.QrisImage = qrisPath
		}
	}

	store.NamaToko = input.NamaToko
	store.Telepon = input.Telepon
	store.Alamat = input.Alamat
	store.PaymentType = input.PaymentType
	store.ReceiptFooter = input.ReceiptFooter

	if err := tx.Save(&store).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui profil toko"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"status": "sukses", "message": "Pengaturan toko berhasil diperbarui!"})
}