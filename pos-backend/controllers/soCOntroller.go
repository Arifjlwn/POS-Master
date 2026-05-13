package controllers

import (
	"net/http"
	"pos-backend/config"
	"pos-backend/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Struct untuk menangkap input dari Frontend (Vue)
type StockOpnameInput struct {
	Notes string `json:"notes"`
	Items []struct {
		ProductID uint `json:"product_id"`
		ActualQty int  `json:"actual_qty"`
	} `json:"items"`
}

func CreateStockOpname(c *gin.Context) {
	// 1. Ambil Identitas User & Store dari Middleware Auth
	storeIDRaw, _ := c.Get("store_id")
	userIDRaw, _ := c.Get("user_id")

	// Konversi tipe data interface{} ke uint
	storeID := uint(storeIDRaw.(float64))
	userID := uint(userIDRaw.(float64))

	var input StockOpnameInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid!"})
		return
	}

	if len(input.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Daftar barang opname tidak boleh kosong!"})
		return
	}

	// 2. Mulai Transaksi Database
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		
		// Buat Header Stock Opname
		so := models.StockOpname{
			StoreID:   storeID,
			UserID:    userID,
			Notes:     input.Notes,
			CreatedAt: time.Now(),
		}

		if err := tx.Create(&so).Error; err != nil {
			return err
		}

		// 3. Looping Barang yang di-Opname
		for _, item := range input.Items {
			var product models.Product
			
			// Ambil data stok sistem saat ini
			if err := tx.First(&product, item.ProductID).Error; err != nil {
				return err // Produk tidak ditemukan, batalkan semua
			}

			// Hitung Selisih (Fisik - Sistem)
			selisih := item.ActualQty - product.Stok

			// Buat Detail Riwayat SO
			detail := models.StockOpnameDetail{
				OpnameID:  so.ID,
				ProductID: item.ProductID,
				SystemQty: product.Stok,
				ActualQty: item.ActualQty,
				Selisih:   selisih,
			}

			if err := tx.Create(&detail).Error; err != nil {
				return err
			}

			// 🚀 4. SINKRONISASI STOK MASTER
			// Kita ubah stok di tabel products menjadi angka riil hasil hitungan karyawan
			product.Stok = item.ActualQty
			if err := tx.Save(&product).Error; err != nil {
				return err
			}
		}

		return nil
	})

	// 5. Response Akhir
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses Stock Opname: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Stock Opname berhasil disimpan. Stok master telah diperbarui!",
	})
}

// Fitur Tambahan: Get History SO (Hanya untuk Owner)
func GetStockOpnameHistory(c *gin.Context) {
	storeID, _ := c.Get("store_id")
	
	var history []models.StockOpname
	config.DB.Preload("Details.Product").
		Where("store_id = ?", storeID).
		Order("created_at DESC").
		Find(&history)

	c.JSON(http.StatusOK, gin.H{"data": history})
}