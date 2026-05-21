package delivery

import (
	"net/http"
	"strconv"
	"time"

	"pos-backend/src/modules/fnb/domain"
	"pos-backend/src/modules/fnb/repository"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	Repo repository.OrderRepository
}

func NewOrderHandler(repo repository.OrderRepository) *OrderHandler {
	return &OrderHandler{Repo: repo}
}

// POST /api/fnb/order
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	storeID := c.MustGet("store_id").(uint)
	userID := c.MustGet("user_id").(uint) // Ambil ID user dari JWT tracking

	var input struct {
		SessionID    int                   `json:"session_id"`
		TipeOrder    string                `json:"tipe_order" binding:"required"`
		NomorMeja    string                `json:"nomor_meja"`
		NamaPemesan  string                `json:"nama_pemesan"`
		MetodeBayar  string                `json:"metode_bayar" binding:"required"`
		UangDiterima int                   `json:"uang_diterima"`
		Kembalian    int                   `json:"kembalian"`
		TotalHarga   int                   `json:"total_harga" binding:"required"`
		Items        []domain.OrderItemFnB `json:"items" binding:"required,dive"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate Invoice Simple (Contoh: FNB-20260521-160000)
	invoiceCode := "FNB-" + time.Now().Format("20060102-150405")

	order := domain.OrderFnB{
		StoreID:      storeID,
		CreatedBy:    userID,
		Invoice:      invoiceCode,
		SessionID:    input.SessionID,
		TipeOrder:    input.TipeOrder,
		NomorMeja:    input.NomorMeja,
		NamaPemesan:  input.NamaPemesan,
		MetodeBayar:  input.MetodeBayar,
		UangDiterima: input.UangDiterima,
		Kembalian:    input.Kembalian,
		TotalHarga:   input.TotalHarga,
		StatusDapur:  "PENDING", // Otomatis masuk antrean dapur
		Items:        input.Items,
	}

	if err := h.Repo.CreateOrder(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses transaksi order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Order berhasil disimpan!", "invoice": invoiceCode, "data": order})
}

// GET /api/fnb/kitchen
func (h *OrderHandler) GetAntreanDapur(c *gin.Context) {
	storeID := c.MustGet("store_id").(uint)

	orders, err := h.Repo.GetKitchenQueue(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil antrean dapur"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// PUT /api/fnb/kitchen/:id/selesai
func (h *OrderHandler) SelesaikanOrderan(c *gin.Context) {
	storeID := c.MustGet("store_id").(uint)
	orderIDStr := c.Param("id")
	orderID, _ := strconv.Atoi(orderIDStr)

	if err := h.Repo.UpdateKitchenStatus(uint(orderID), storeID, "SELESAI"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui status dapur"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Orderan selesai dimasak dan siap disajikan!"})
}