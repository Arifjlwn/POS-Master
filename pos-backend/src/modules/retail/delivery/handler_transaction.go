package delivery

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"pos-backend/models"
	"pos-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ====================================
// 🛒 POS TRANSACTION LOGIC HANDLERS
// ====================================

type CartItem struct {
	ProductID uint    `json:"product_id" binding:"required"`
	Kuantitas int     `json:"kuantitas" binding:"required,gt=0"`
	UomLabel  string  `json:"uom_label"`
	HargaUom  float64 `json:"harga_uom"`
}

type TransactionInput struct {
	Items         []CartItem `json:"items" binding:"required,gt=0"`
	NominalBayar  float64    `json:"nominal_bayar" binding:"required"`
	MetodeBayar   string     `json:"metode_bayar"`
	NoHPPelanggan string     `json:"no_hp_pelanggan"`
}

func (h *RetailHandler) CreateTransaction(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id"); userIDRaw, _ := c.Get("user_id")
	var storeID, userID uint
	switch v := storeIDRaw.(type) { case float64: storeID = uint(v); case uint: storeID = v; case int: storeID = uint(v) }
	switch v := userIDRaw.(type) { case float64: userID = uint(v); case uint: userID = v; case int: userID = uint(v) }

	var input TransactionInput
	if err := c.ShouldBindJSON(&input); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Format keranjang tidak sesuai!"}); return }

	formatRupiah := func(amount float64) string {
		str := fmt.Sprintf("%.0f", amount); var result string
		for i, n := len(str)-1, 0; i >= 0; i-- {
			result = string(str[i]) + result; n++
			if n%3 == 0 && i > 0 { result = "." + result }
		}
		return "Rp " + result
	}

	var savedTransaction models.Transaction
	db := h.Repo.GetDB()

	err := db.Transaction(func(tx *gorm.DB) error {
		activeSession, err := h.Repo.GetActiveSession(tx, userID, storeID)
		if err != nil { return fmt.Errorf("session kasir tidak ditemukan, silakan buka kasir dulu bray") }

		var store models.Store
		if err := tx.Select("id", "nama_toko", "business_type", "pajak_persen", "wa_token", "receipt_footer").First(&store, storeID).Error; err != nil {
			return fmt.Errorf("data toko tidak ditemukan")
		}

		tipeBisnis := "RETAIL"; statusPesanan := "SELESAI"
		if store.BusinessType == "Jasa - Laundry" { tipeBisnis = "LAUNDRY"; statusPesanan = "ANTRI" } 
		if store.BusinessType == "Kuliner - F&B" { tipeBisnis = "FNB"; statusPesanan = "PROSES" }

		var subTotal float64; var details []models.TransactionDetail; rincianBarangWA := ""

		for _, item := range input.Items {
			product, err := h.Repo.GetProductByID(tx, item.ProductID, storeID)
			if err != nil { return err }
			if product.Stok < item.Kuantitas { return fmt.Errorf("Stok %s habis! Sisa: %d", product.NamaProduk, product.Stok) }

			product.Stok -= item.Kuantitas
			if err := h.Repo.SaveProduct(tx, product); err != nil { return err }

			itemSubTotal := float64(item.Kuantitas) * product.HargaJual
			rincianDisplay := item.UomLabel; hargaSatuanDisplay := item.HargaUom

			if hargaSatuanDisplay > 0 {
				qtyOriginal := 1; fmt.Sscanf(rincianDisplay, "%d", &qtyOriginal)
				if qtyOriginal > 0 { itemSubTotal = float64(qtyOriginal) * hargaSatuanDisplay }
			}

			subTotal += itemSubTotal
			rincianBarangWA += fmt.Sprintf("▪️ *%s*\n   %s x %s = *%s*\n", product.NamaProduk, rincianDisplay, formatRupiah(hargaSatuanDisplay), formatRupiah(itemSubTotal))

			// 🚀 HYBRID OPTIMAL: Detail bersih tanpa PublicID ULID, hemat index space DB!
			details = append(details, models.TransactionDetail{
				ProductID:   product.ID,
				HargaSatuan: hargaSatuanDisplay,
				Kuantitas:   item.Kuantitas,
				SubTotal:    itemSubTotal,
				ItemType:    "PRODUCT",
				DetailNotes: rincianDisplay,
			})
		}

		pajak := (store.PajakPersen / 100.0) * subTotal
		rawTotal := subTotal + pajak
		roundedTotal := math.Round(rawTotal/100) * 100
		pembulatan := roundedTotal - rawTotal

		if input.NominalBayar < roundedTotal { return fmt.Errorf("Uang pelanggan kurang! Tagihan: %s", formatRupiah(roundedTotal)) }
		kembalian := input.NominalBayar - roundedTotal
		noInvoice := fmt.Sprintf("INV-%s", time.Now().Format("20060102150405"))

		savedTransaction = models.Transaction{
			PublicID:      utils.GenerateULID(), // Header tetep wajib dipasang bray buat top-level masking API eksternal
			SessionID:     activeSession.ID,
			StoreID:       storeID,
			UserID:        userID,
			NoInvoice:     noInvoice,
			SubTotal:      subTotal,
			Pajak:         pajak,
			Pembulatan:    pembulatan,
			TotalHarga:    roundedTotal,
			MetodeBayar:   input.MetodeBayar,
			StatusBayar:   "LUNAS",
			TipeBisnis:    tipeBisnis,
			StatusPesanan: statusPesanan,
			NominalBayar:  input.NominalBayar,
			Kembalian:     kembalian,
			Details:       details,
		}

		if input.NoHPPelanggan != "" && store.WaToken != "" {
			pesanNota := fmt.Sprintf(
				`🏪 *%s*
========================
Halo Bosku! Terima kasih sudah berbelanja. Berikut rincian transaksi Anda:

🧾 *No. Invoice:* %s
📅 *Tanggal:* %s

*Rincian Pesanan:*
%s========================
💰 *Subtotal:* %s
⚖️ *Pajak/Biaya:* %s
========================
✅ *TOTAL BAYAR: %s*
💳 *Metode:* %s
💵 *Tunai:* %s
💸 *Kembali:* %s
========================
				%s`,
				store.NamaToko, noInvoice, time.Now().Format("02 Jan 2006, 15:04 WIB"), rincianBarangWA,
				formatRupiah(subTotal), formatRupiah(pajak), formatRupiah(roundedTotal), input.MetodeBayar,
				formatRupiah(input.NominalBayar), formatRupiah(kembalian), store.ReceiptFooter,
			)
			utils.SendWhatsAppFonnte(store.WaToken, input.NoHPPelanggan, pesanNota)
		}

		return h.Repo.CreateTransactionTx(tx, &savedTransaction)
	})

	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaksi berhasil! Struk siap dicetak.",
		"invoice": savedTransaction.NoInvoice,
		"tagihan": savedTransaction.TotalHarga,
		"kembali": savedTransaction.Kembalian,
	})
}

func (h *RetailHandler) GetTransactions(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	var storeID uint
	switch v := storeIDRaw.(type) { case float64: storeID = uint(v); case uint: storeID = v; case int: storeID = uint(v) }

	tanggal := c.Query("tanggal")
	if tanggal == "" { tanggal = time.Now().Format("2006-01-02") }

	parsedDate, err := time.ParseInLocation("2006-01-02", tanggal, time.Local)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal tidak valid"}); return }

	startOfDay := parsedDate; endOfDay := startOfDay.Add(24 * time.Hour)
	transactions, err := h.Repo.GetTransactionsByRange(storeID, startOfDay, endOfDay)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menarik riwayat transaksi"}); return }

	c.JSON(http.StatusOK, gin.H{"message": "Riwayat transaksi berhasil ditarik!", "data": transactions})
}

func (h *RetailHandler) GetDailyClosing(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	var storeID uint
	switch v := storeIDRaw.(type) { case float64: storeID = uint(v); case uint: storeID = v; case int: storeID = uint(v) }

	tanggal := c.Query("tanggal")
	if tanggal == "" { tanggal = time.Now().Format("2006-01-02") }

	parsedDate, err := time.ParseInLocation("2006-01-02", tanggal, time.Local)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal tidak valid"}); return }

	startOfDay := parsedDate; endOfDay := startOfDay.Add(24 * time.Hour)
	closings, err := h.Repo.GetClosingByRange(storeID, startOfDay, endOfDay)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menarik riwayat closing"}); return }

	c.JSON(http.StatusOK, gin.H{"message": "Riwayat closing berhasil ditarik!", "data": closings})
}