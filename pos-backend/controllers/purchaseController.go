package controllers

import (
	"net/http"
	"pos-backend/config"
	"pos-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PurchaseInput struct {
	SupplierName string `json:"supplier_name"`
	NoFaktur     string `json:"no_faktur"`
	Items        []struct {
		ProductID uint    `json:"product_id"`
		QtyMasuk  int     `json:"qty_masuk"`
		HargaModal float64 `json:"harga_modal"`
	} `json:"items"`
}

func CreateLPB(c *gin.Context) {
	storeID, _ := c.Get("store_id")
	userID, _ := c.Get("user_id")

	var input PurchaseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Buat Header Purchase
		purchase := models.Purchase{
			StoreID:      uint(storeID.(float64)),
			UserID:       uint(userID.(float64)),
			SupplierName: input.SupplierName,
			NoFaktur:     input.NoFaktur,
			TotalItem:    len(input.Items),
		}

		if err := tx.Create(&purchase).Error; err != nil {
			return err
		}

		// 2. Loop Items: Update Stok & Simpan Detail
		for _, item := range input.Items {
			// Update Stok di tabel Product
			var product models.Product
			if err := tx.First(&product, item.ProductID).Error; err != nil {
				return err
			}

			product.Stok += item.QtyMasuk
			// Optional: Mas Arif bisa update harga beli master barang di sini
			// product.HargaBeli = item.HargaBeli 

			if err := tx.Save(&product).Error; err != nil {
				return err
			}

			// Simpan Detail Purchase
			detail := models.PurchaseDetail{
				PurchaseID: purchase.ID,
				ProductID:  item.ProductID,
				QtyMasuk:   item.QtyMasuk,
				HargaModal:  item.HargaModal,
			}
			if err := tx.Create(&detail).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal proses LPB: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "LPB Berhasil! Stok barang sudah bertambah."})
}