package delivery

import (
	"net/http"
	"pos-backend/src/modules/jasalayanan/laundry/domain"
	"pos-backend/src/modules/jasalayanan/laundry/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
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
