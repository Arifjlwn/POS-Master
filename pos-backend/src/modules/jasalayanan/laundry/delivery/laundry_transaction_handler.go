package delivery

import (
	"fmt"
	"net/http"
	"os"
	"pos-backend/src/modules/jasalayanan/laundry/domain"
	"pos-backend/src/modules/jasalayanan/laundry/usecase"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type LaundryTransactionHandler struct {
	usecase usecase.LaundryUseCase
}

func NewLaundryTransactionHandler(uc usecase.LaundryUseCase) *LaundryTransactionHandler {
	return &LaundryTransactionHandler{usecase: uc}
}

func (h *LaundryTransactionHandler) ProsesCheckoutLaundry(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	userIDRaw, _ := c.Get("user_id")
	storeID := uint(storeIDRaw.(float64))
	userID := uint(userIDRaw.(float64))

	var input domain.CheckoutLaundryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data transaksi laundry tidak valid !"})
		return
	}

	invoiceCode, fotoURL, nomorRak, err := h.usecase.ProcessCheckout(storeID, userID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       "sukses",
		"message":      "Transaksi laundry berhasil diproses!",
		"invoice_code": invoiceCode,
		"foto_url":     fotoURL,
		"nomor_rak":    nomorRak, // 🚀 LEMPAR NOMOR RAK KE FRONTEND VUE BRAY!
	})
}

func (h *LaundryTransactionHandler) LunasiTransaksi(c *gin.Context) {
	trxID, _ := strconv.Atoi(c.Param("id"))
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

	var input domain.PelunasanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Metode pembayaran wajib dipilih"})
		return
	}

	if err := h.usecase.ProcessPelunasan(uint(trxID), storeID, input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "sukses", "message": "Tagihan berhasil dilunasi!"})
}

func (h *LaundryTransactionHandler) GetLaporan(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

	summary, err := h.usecase.GetLaporanRingkasan(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengompilasi laporan audit laundry"})
		return
	}
	c.JSON(http.StatusOK, summary)
}

func (h *LaundryTransactionHandler) UpdateStatusCucian(c *gin.Context) {
	trxID, _ := strconv.Atoi(c.Param("id"))
	var input struct {
		StatusPesanan string `json:"status_pesanan" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data perubahan status tidak valid"})
		return
	}
	if err := h.usecase.UpdateStatusCucian(uint(trxID), input.StatusPesanan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Status pengerjaan cucian berhasil diperbarui!"})
}

// Struct buat nangkep total tagihan dari Vue
type MidtransLaundryInput struct {
	Total float64 `json:"total" binding:"required"`
}

func (h *LaundryTransactionHandler) GetMidtransTokenLaundry(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := extractUintID(storeIDRaw)

	var input MidtransLaundryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Total nominal tagihan wajib dilampirkan bray!"})
		return
	}

	// 🛡️ SECURITY PATCH: Bikin Order ID unik khusus laundry biar gak tabrakan sama retail
	orderID := fmt.Sprintf("LAUNDRY-%d-%d", storeID, time.Now().UnixNano())

	// 🚀 DYNAMIC ENVIRONMENT CHECK (Kasta Tertinggi!)
	// Default ke Sandbox, tapi kalau di .env tulisannya production, dia otomatis pindah ke Production!
	midtransEnv := midtrans.Sandbox
	if os.Getenv("APP_ENV") == "production" {
		midtransEnv = midtrans.Production
	}

	// 🚀 TEMBAK KE MIDTRANS CORE ENGINE
	s := snap.Client{}
	s.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtransEnv)

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(input.Total),
		},
	}

	snapTokenResponse, errMidtrans := s.CreateTransaction(req)
	if errMidtrans != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Midtrans Tolak Transaksi: %v", errMidtrans.GetMessage())})
		return
	}

	// Kembalikan tokennya ke Vue bray!
	c.JSON(http.StatusOK, gin.H{
		"status": "sukses",
		"token":  snapTokenResponse.Token,
	})
}
