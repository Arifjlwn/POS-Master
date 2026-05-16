package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"pos-backend/config" // Sesuaikan nama module-mu
	"pos-backend/models" // Sesuaikan nama module-mu

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 🚀 STRUCT BARU UNTUK NERIMA KERANJANG RETUR
type ReturnItem struct {
	ProductID uint   `json:"product_id" binding:"required"`
	Qty       int    `json:"qty" binding:"required,gt=0"`
	Alasan    string `json:"alasan" binding:"required"`
	Catatan   string `json:"catatan"`
}

type ReturnInputBatch struct {
	Items []ReturnItem `json:"items" binding:"required,min=1"`
}

func CreateReturn(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	userIDRaw, _ := c.Get("user_id")
	storeID := uint(storeIDRaw.(float64))
	userID := uint(userIDRaw.(float64))

	var input ReturnInputBatch
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data keranjang retur tidak valid!"})
		return
	}

	// 🚀 BIKIN NOMOR RETUR OTOMATIS (Contoh: RET-160526-12345)
	returnNo := fmt.Sprintf("RET-%s-%d", time.Now().Format("060102150405"), userID)

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var newReturns []models.ProductReturn

	// 🚀 LOOPING SEMUA BARANG DI KERANJANG
	for _, item := range input.Items {
		var product models.Product
		if err := tx.Where("id = ? AND store_id = ?", item.ProductID, storeID).First(&product).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Ada produk yang tidak ditemukan di toko ini!"})
			return
		}

		if product.Stok < item.Qty {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Stok %s tidak mencukupi (Sisa: %d)!", product.NamaProduk, product.Stok)})
			return
		}

		// Potong stok
		if err := tx.Model(&product).Update("stok", gorm.Expr("stok - ?", item.Qty)).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memotong stok produk!"})
			return
		}

		// Siapkan data log
		newReturns = append(newReturns, models.ProductReturn{
			ReturnNo:  returnNo, // Pakai nomor retur yang sama untuk 1 keranjang
			StoreID:   storeID,
			ProductID: item.ProductID,
			UserID:    userID,
			Qty:       item.Qty,
			Alasan:    item.Alasan,
			Catatan:   item.Catatan,
		})
	}

	// Simpan semua log retur sekaligus
	if err := tx.Create(&newReturns).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencatat log retur batch!"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message":   "Berita Acara Retur berhasil diproses!",
		"return_no": returnNo, // Kasih tau nomor returnya ke Vue
	})
}

// GetReturns buat ambil history 
func GetReturns(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	var returns []models.ProductReturn
	var totalItems int64

	query := config.DB.Model(&models.ProductReturn{}).Where("store_id = ?", storeID).
		Preload("Product").
		Preload("User")

	query.Count(&totalItems)

	if pageStr != "" && limitStr != "" {
		page, _ := strconv.Atoi(pageStr)
		limit, _ := strconv.Atoi(limitStr)
		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	}

	// Order by terbaru
	if err := query.Order("created_at DESC").Find(&returns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data retur"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Data retur berhasil dimuat!",
		"total_items": totalItems,
		"data":        returns,
	})
}