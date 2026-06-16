package delivery

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"pos-backend/models"
	src "pos-backend/src/core/config"

	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

// =====================================================================
// MIDTRANS TRANSAKSI KASIR RETAIL (UANG MASUK REKENING MERCHANT/TENANT)
// =====================================================================

type PosMidtransReq struct {
	Total float64 `json:"total" binding:"required"`
}

func (h *RetailHandler) CreatePosMidtransOrder(c *gin.Context) {
	storeIDRaw, exists := c.Get("store_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak valid"})
		return
	}
	var storeID uint
	switch v := storeIDRaw.(type) {
	case float64:
		storeID = uint(v)
	case uint:
		storeID = v
	case int:
		storeID = uint(v)
	}

	var input PosMidtransReq
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data nominal tidak valid"})
		return
	}

	var store models.Store
	if err := src.DB.First(&store, storeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Toko tidak ditemukan!"})
		return
	}
	if store.PaymentType != "midtrans" || store.MidtransServerKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Toko belum mengatur Midtrans Server Key!"})
		return
	}

	var s snap.Client
	env := midtrans.Sandbox
	if os.Getenv("APP_ENV") == "production" {
		env = midtrans.Production
	}

	s.New(store.MidtransServerKey, env)
	orderID := fmt.Sprintf("POS-STR%d-%d", storeID, time.Now().Unix())

	req := &snap.Request{TransactionDetails: midtrans.TransactionDetails{OrderID: orderID, GrossAmt: int64(input.Total)}}
	snapResp, err := s.CreateTransaction(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.GetMessage()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": snapResp.Token, "order_id": orderID})
}
