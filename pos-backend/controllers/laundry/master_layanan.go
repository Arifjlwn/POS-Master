package laundry

import (
	"net/http"
	"pos-backend/config"
	"pos-backend/models"

	"github.com/gin-gonic/gin"
)

type InputLayanan struct {
	NamaProduk  string  `json:"nama_produk" binding:"required"`
	HargaJual   float64 `json:"harga_jual" binding:"required"`
	SatuanDasar string  `json:"satuan_dasar" binding:"required"` // KG atau PCS
	Estimasi    string  `json:"estimasi"`                        // 1 Hari, Kilat, dll
}

// 🚀 FUNGSI TAMBAH PAKET LAUNDRY
func TambahLayananLaundry(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

	var input InputLayanan
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input layanan tidak valid"})
		return
	}

	// Kita jadikan produk dengan settingan khusus JASA (Stok 0, SKU nil)
	newLayanan := models.Product{
		StoreID:     storeID,
		NamaProduk:  input.NamaProduk,
		Kategori:    "JASA_LAUNDRY",
		HargaJual:   input.HargaJual,
		SatuanDasar: input.SatuanDasar,
		Estimasi:    input.Estimasi,
		Stok:        0, // Jasa nggak butuh stok
	}

	if err := config.DB.Create(&newLayanan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan layanan baru"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Layanan berhasil ditambahkan",
		"data":    newLayanan,
	})
}

// 🚀 FUNGSI HAPUS PAKET LAUNDRY
func HapusLayananLaundry(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))
	productID := c.Param("id")

	if err := config.DB.Where("id = ? AND store_id = ?", productID, storeID).Delete(&models.Product{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus layanan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Layanan berhasil dihapus"})
}

// 🚀 FUNGSI EDIT PAKET LAUNDRY
func EditLayananLaundry(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))
	productID := c.Param("id")

	var input InputLayanan
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input layanan tidak valid"})
		return
	}

	// Cari layanannya dulu, beneran ada gak di toko ini?
	var layanan models.Product
	if err := config.DB.Where("id = ? AND store_id = ?", productID, storeID).First(&layanan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Layanan tidak ditemukan"})
		return
	}

	// Update data yang baru
	layanan.NamaProduk = input.NamaProduk
	layanan.HargaJual = input.HargaJual
	layanan.SatuanDasar = input.SatuanDasar
	layanan.Estimasi = input.Estimasi

	if err := config.DB.Save(&layanan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui layanan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Layanan berhasil diperbarui!"})
}