package laundry

import (
	"net/http"
	"pos-backend/config"
	"pos-backend/models"

	"github.com/gin-gonic/gin"
)

type PerfumeInput struct {
	Nama  string  `json:"nama" binding:"required"`
	Harga float64 `json:"harga"`
}

// 🚀 1. AMBIL SEMUA PARFUM DI TOKO INI
func GetPerfumes(c *gin.Context) {
	storeID, _ := c.Get("store_id")

	var perfumes []models.Perfume
	if err := config.DB.Where("store_id = ?", storeID).Find(&perfumes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data parfum"})
		return
	}

	c.JSON(http.StatusOK, perfumes)
}

// 🚀 2. OWNER TAMBAH PARFUM BARU
func CreatePerfume(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

	var input PerfumeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data input tidak valid"})
		return
	}

	newPerfume := models.Perfume{
		StoreID: storeID,
		Nama:    input.Nama,
		Harga:   input.Harga,
		Status:  "Tersedia",
	}

	if err := config.DB.Create(&newPerfume).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan varian parfum baru"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Parfum berhasil ditambahkan!", "data": newPerfume})
}

// 🚀 3. OWNER HAPUS VARIANT PARFUM
func DeletePerfume(c *gin.Context) {
	storeID, _ := c.Get("store_id")
	perfumeID := c.Param("id")

	if err := config.DB.Where("id = ? AND store_id = ?", perfumeID, storeID).Delete(&models.Perfume{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus parfum"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Varian parfum berhasil dihapus!"})
}