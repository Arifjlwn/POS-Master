package fnb

import (
	"fmt"
	"net/http"
	"pos-backend/config"
	"pos-backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

// Payload dari Vue KasirFnB.vue
type ItemOrder struct {
	ProductID uint    `json:"product_id" binding:"required"`
	Qty       int     `json:"qty" binding:"required"`
	SubTotal  float64 `json:"sub_total" binding:"required"`
	Notes     string  `json:"notes"`
}

type CheckoutPayload struct {
	SessionID  uint        `json:"session_id" binding:"required"` // 🚀 Wajib bawa Session ID aktif kasir
	TipeOrder  string      `json:"tipe_order" binding:"required"` // DINE_IN / TAKE_AWAY
	NomorMeja  string      `json:"nomor_meja"`
	SubTotal   float64     `json:"sub_total" binding:"required"`
	Pajak      float64     `json:"pajak"`
	TotalHarga float64     `json:"total_harga" binding:"required"`
	Items      []ItemOrder `json:"items" binding:"required"`
}

// 🚀 FUNGSI UTAMA KASIR F&B: 1 PAGE 1 BACKEND INTERACTION
func CreateOrder(c *gin.Context) {
	// Ambil metadata toko & user dari JWT Middleware
	storeIDRaw, _ := c.Get("store_id")
	userIDRaw, _ := c.Get("user_id")
	
	storeID := uint(storeIDRaw.(float64))
	userID := uint(userIDRaw.(float64))

	var payload CheckoutPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data pesanan fnb tidak valid: " + err.Error()})
		return
	}

	// Memulai DB Transaction demi keamanan data integrity
	tx := config.DB.Begin()

	// 1. Generate Nomor Invoice Khusus Resto
	// Format contoh: INV/FNB/20260519/0001 (bisa disesuaikan seleramu beb)
	invoiceNumber := fmt.Sprintf("INV/FNB/%s/%d", time.Now().Format("20060102"), time.Now().UnixNano()%10000)

	// 2. Petakan ke Struktur Tabel Core Transaction Global
	newTransaction := models.Transaction{
		SessionID:     payload.SessionID, // Diisi dari session_id kasir aktif
		StoreID:       storeID,
		UserID:        userID,
		NoInvoice:     invoiceNumber,
		SubTotal:      payload.SubTotal,
		Pajak:         payload.Pajak,
		Pembulatan:    0,
		TotalHarga:    payload.TotalHarga,
		MetodeBayar:   "Cash",       // Default, nanti bisa di-update pas pelunasan ambil makanan
		StatusBayar:   "BELUM_LUNAS", // Mengikuti alur pesanan F&B berjalan
		StatusPesanan: payload.TipeOrder, // Menyimpan tipe order (DINE_IN/TAKE_AWAY)
		NominalBayar:  0,            // Belum bayar
		Kembalian:     0,
	}

	// Eksekusi insert ke database core
	if err := tx.Create(&newTransaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan core transaction bisnis: " + err.Error()})
		return
	}

	// 3. Simpan Rincian Menu ke Tabel Modul Spesifik (FnBDetail)
	for _, item := range payload.Items {
		var produk models.Product
		if err := tx.First(&produk, item.ProductID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Produk/Menu tidak ditemukan di database"})
			return
		}

		fnbItem := models.FnBDetail{
			TransactionID: newTransaction.ID,
			ProductID:     item.ProductID,
			Qty:           item.Qty,
			HargaSatuan:   produk.HargaJual, // Ambil nominal tervalidasi dari server
			SubTotal:      item.SubTotal,
			Notes:         item.Notes,       // Menyimpan catatan koki kustom
			StatusDapur:   "PROSES",         // Otomatis berjejer di monitor dapur koki
		}

		if err := tx.Create(&fnbItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menginput rincian item ke dapur"})
			return
		}
	}

	// Jika sukses tanpa kendala, sahkan transaksi!
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Pesanan berhasil diteruskan ke Monitor Dapur!",
		"invoice": invoiceNumber,
	})
}