package fnb

import (
	"net/http"
	"pos-backend/config"
	"pos-backend/models"

	"github.com/gin-gonic/gin"
)

// Response cetakan agar format JSON pas dengan layar DapurFnB.vue
type DapurItemResponse struct {
	ID    uint   `json:"id"`
	Nama  string `json:"nama"`
	Qty   int    `json:"qty"`
	Notes string `json:"notes"`
}

type DapurOrderResponse struct {
	ID         uint                `json:"id"`
	Invoice    string              `json:"invoice"`
	Tipe       string              `json:"tipe"`
	Meja       string              `json:"meja"`
	WaktuPesan string              `json:"waktu_pesan"`
	Status     string              `json:"status"`
	Items      []DapurItemResponse `json:"items"`
}

// 🚀 FUNGSI TARIK DATA ANTREAN DAPUR
func GetAntreanDapur(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

	var details []models.FnBDetail

	// Join tabel FnBDetail dengan Transaction dan Product buat narik info komplit
	if err := config.DB.Preload("Transaction").Preload("Product").
		Joins("JOIN transactions ON transactions.id = fn_b_details.transaction_id").
		Where("transactions.store_id = ? AND fn_b_details.status_dapur = ?", storeID, "PROSES").
		Order("transactions.created_at ASC").
		Find(&details).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil antrean dapur"})
		return
	}

	// Grouping menu berdasarkan ID Transaksi biar jadi 1 Tiket KOT utuh
	orderMap := make(map[uint]*DapurOrderResponse)

	for _, d := range details {
		trxID := d.TransactionID
		if _, exists := orderMap[trxID]; !exists {
			orderMap[trxID] = &DapurOrderResponse{
				ID:         trxID,
				Invoice:    d.Transaction.NoInvoice,
				Tipe:       d.Transaction.StatusPesanan, // Nampung TipeOrder (DINE_IN/TAKE_AWAY)
				Meja:       "Meja", // Nanti bisa dikaitkan dengan field meja jika ada di core
				WaktuPesan: d.Transaction.CreatedAt.Format("15:04"),
				Status:     "PROSES",
				Items:      []DapurItemResponse{},
			}
		}

		orderMap[trxID].Items = append(orderMap[trxID].Items, DapurItemResponse{
			ID:    d.ID,
			Nama:  d.Product.NamaProduk,
			Qty:   d.Qty,
			Notes: d.Notes, // Catatan koki tampil di sini
		})
	}

	// Convert hasil grouping ke array
	var result []DapurOrderResponse
	for _, order := range orderMap {
		result = append(result, *order)
	}

	c.JSON(http.StatusOK, result)
}

// 🚀 FUNGSI SELASAIKAN PESANAN (Ubah Status PROSES -> SELESAI)
func SelesaikanOrderan(c *gin.Context) {
	trxID := c.Param("id")

	// Update seluruh item yang ada di Transaksi tersebut jadi SELESAI
	if err := config.DB.Model(&models.FnBDetail{}).
		Where("transaction_id = ?", trxID).
		Update("status_dapur", "SELESAI").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyelesaikan pesanan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pesanan siap dihidangkan!"})
}