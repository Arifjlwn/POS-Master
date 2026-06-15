package delivery

import (
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"pos-backend/models"
	"pos-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ====================================
// 🛒 POS TRANSACTION LOGIC HANDLERS
// ====================================

type CartItem struct {
	ProductID uint `json:"product_id" binding:"required"`
	// 🚀 AMAN 1: Ubah ke float64 untuk mendukung desimal 0.25 / 0.5 bray
	Kuantitas float64 `json:"kuantitas" binding:"required,gt=0"`
	UomLabel  string  `json:"uom_label"`
	HargaUom  float64 `json:"harga_uom"`
}

type TransactionInput struct {
	Items         []CartItem `json:"items" binding:"required,gt=0"`
	NominalBayar  float64    `json:"nominal_pay" binding:"required"` // Disesuaikan atau disamakan nominal_bayar
	MetodeBayar   string     `json:"metode_bayar"`
	NoHPPelanggan string     `json:"no_hp_pelanggan"`
}

func (h *RetailHandler) CreateTransaction(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	userIDRaw, _ := c.Get("user_id")
	var storeID, userID uint
	switch v := storeIDRaw.(type) {
	case float64:
		storeID = uint(v)
	case uint:
		storeID = v
	case int:
		storeID = uint(v)
	}
	switch v := userIDRaw.(type) {
	case float64:
		userID = uint(v)
	case uint:
		userID = v
	case int:
		userID = uint(v)
	}

	var input TransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format keranjang tidak sesuai !"})
		return
	}

	formatRupiah := func(amount float64) string {
		str := fmt.Sprintf("%.0f", amount)
		var result string
		for i, n := len(str)-1, 0; i >= 0; i-- {
			result = string(str[i]) + result
			n++
			if n%3 == 0 && i > 0 {
				result = "." + result
			}
		}
		return "Rp " + result
	}

	var savedTransaction models.Transaction
	var store models.Store
	rincianBarangWA := ""
	db := h.Repo.GetDB()

	err := db.Transaction(func(tx *gorm.DB) error {
		activeSession, err := h.Repo.GetActiveSession(tx, userID, storeID)
		if err != nil {
			return fmt.Errorf("sesi laci kasir belum dibuka, isi modal awal dulu ")
		}

		if err := tx.Select("id", "nama_toko", "business_type", "pajak_persen", "wa_token", "receipt_footer").First(&store, storeID).Error; err != nil {
			return fmt.Errorf("data identitas toko gagal diverifikasi")
		}

		tipeBisnis := "RETAIL"
		statusPesanan := "SELESAI"
		if store.BusinessType == "Jasa - Laundry" {
			tipeBisnis = "LAUNDRY"
			statusPesanan = "ANTRI"
		}
		if store.BusinessType == "Kuliner - F&B" {
			tipeBisnis = "FNB"
			statusPesanan = "PROSES"
		}

		var subTotal float64
		var details []models.TransactionDetail

		for _, item := range input.Items {
			var product models.Product

			// 🛡️ LOCK ROW CONCURRENCY
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&product, "id = ? AND store_id = ?", item.ProductID, storeID).Error; err != nil {
				return fmt.Errorf("produk ID %d tidak ditemukan di rak toko", item.ProductID)
			}

			// 🚀 AMAN 2: Paksa casting product.Stok ke float64 jika di model masih int, mencegah crash mismatched types bray!
			if float64(product.Stok) < item.Kuantitas {
				return fmt.Errorf("stok %s habis ! Sisa fisik: %.2f", product.NamaProduk, float64(product.Stok))
			}

			// Potong stok produk (Casting kembali agar tipe data sinkron)
			product.Stok -= int(item.Kuantitas) // Jika models.Product.Stok masih int, kalau nanti udah dimigrasi ke float64, buang int() nya bray
			if err := tx.Save(&product).Error; err != nil {
				return err
			}

			// 🛡️ ANTI TWEAKING HARGA FRONTEND
			cleanLabel := strings.ToUpper(item.UomLabel)
			var hargaSatuanResmi float64 = product.HargaJual

			if cleanLabel != "" {
				if product.SatuanBesar != "" && strings.Contains(cleanLabel, strings.ToUpper(product.SatuanBesar)) {
					hargaSatuanResmi = product.HargaJualBesar
				} else if product.SatuanTengah != "" && strings.Contains(cleanLabel, strings.ToUpper(product.SatuanTengah)) {
					hargaSatuanResmi = product.HargaJualTengah
				}
			}

			// 🚀 AMAN 3: Ubah %d menjadi %f agar mengekstrak nilai koma murni (cth: "0.5") dengan akurat bray!
			qtyOriginal := 1.0
			fmt.Sscanf(item.UomLabel, "%f", &qtyOriginal)
			if qtyOriginal <= 0 {
				qtyOriginal = 1.0
			}

			itemSubTotal := qtyOriginal * hargaSatuanResmi
			subTotal += itemSubTotal

			rincianBarangWA += fmt.Sprintf("▪️ *%s*\n   %s x %s = *%s*\n", product.NamaProduk, item.UomLabel, formatRupiah(hargaSatuanResmi), formatRupiah(itemSubTotal))

			skuSnapshot := ""
			if product.SKU != nil {
				skuSnapshot = *product.SKU
			}

			// Masukkan ke detail item transaksi
			details = append(details, models.TransactionDetail{
				ProductID:          product.ID,
				NamaProdukSnapshot: product.NamaProduk,
				SKUProductSnapshot: skuSnapshot,
				HargaSatuan:        hargaSatuanResmi,
				Kuantitas:          item.Kuantitas, // Auto-sinkron karena models.TransactionDetail sudah float64 bray
				ItemType:           "PRODUCT",
				DetailNotes:        item.UomLabel,
				SubTotal:           itemSubTotal,
			})
		}

		pajak := (store.PajakPersen / 100.0) * subTotal
		rawTotal := subTotal + pajak
		roundedTotal := math.Round(rawTotal/100) * 100
		pembulatan := roundedTotal - rawTotal

		if input.NominalBayar < roundedTotal {
			return fmt.Errorf("uang kas kurang! Total wajib bayar: %s", formatRupiah(roundedTotal))
		}

		kembalian := input.NominalBayar - roundedTotal
		noInvoice := fmt.Sprintf("INV-%s", time.Now().Format("20060102150405"))

		savedTransaction = models.Transaction{
			PublicID:      utils.GenerateULID(),
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

		return tx.Create(&savedTransaction).Error
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 🚀 ASYNCHRONOUS BACKGROUND THREAD WA BROADCASTER
	if input.NoHPPelanggan != "" && store.WaToken != "" {
		go func(token, hp, namaToko, inv, rincian string, sub, paj, tot, nom, kemb float64, method, footer string) {
			pesanNota := fmt.Sprintf(
				`🏪 *%s*
========================
Halo Bosku! Terima kasih sudah berbelanja. Berikut rincian kuitansi digital Anda:

🧾 *No. Invoice:* %s
📅 *Tanggal:* %s

*Rincian Pesanan:*
%s========================
💰 *Subtotal:* %s
⚖️ *Pajak Resto:* %s
========================
✅ *TOTAL BELANJA: %s*
💳 *Metode:* %s
💵 *Tunai Kas:* %s
💸 *Kembalian:* %s
========================
%s`,
				namaToko, inv, time.Now().Format("02 Jan 2006, 15:04 WIB"), rincian,
				formatRupiah(sub), formatRupiah(paj), formatRupiah(tot), method,
				formatRupiah(nom), formatRupiah(kemb), footer,
			)
			utils.SendWhatsAppFonnte(token, hp, pesanNota)
		}(store.WaToken, input.NoHPPelanggan, store.NamaToko, savedTransaction.NoInvoice, rincianBarangWA, savedTransaction.SubTotal, savedTransaction.Pajak, savedTransaction.TotalHarga, input.NominalBayar, savedTransaction.Kembalian, input.MetodeBayar, store.ReceiptFooter)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Transaksi sukses ! Dokumen kuitansi siap dicetak.",
		"invoice":    savedTransaction.NoInvoice,
		"no_invoice": savedTransaction.NoInvoice,
		"tagihan":    savedTransaction.TotalHarga,
		"kembali":    savedTransaction.Kembalian,
	})
}

// =========================================================================
// 🛒 RIWAYAT LOG REPORT LIST DATA TRANS HANDLERS - FIX TIME WINDOW SINKRON
// =========================================================================

func (h *RetailHandler) GetTransactions(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	var storeID uint
	switch v := storeIDRaw.(type) {
	case float64:
		storeID = uint(v)
	case uint:
		storeID = v
	case int:
		storeID = uint(v)
	}

	tanggal := c.Query("tanggal")
	if tanggal == "" {
		tanggal = time.Now().Format("2006-01-02")
	}

	var startOfDay, endOfDay time.Time

	// 🚀 DETEKSI OTOMATIS: Apakah inputnya Bulanan (YYYY-MM) atau Harian (YYYY-MM-DD)
	if len(tanggal) == 7 { // Format: 2026-06 (Bulanan)
		parsedMonth, err := time.ParseInLocation("2006-01", tanggal, time.Local)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format bulanan tidak valid!"})
			return
		}
		// Batas awal: Tanggal 1 jam 00:00:00
		startOfDay = time.Date(parsedMonth.Year(), parsedMonth.Month(), 1, 0, 0, 0, 0, parsedMonth.Location())
		// Batas akhir: Tanggal 1 bulan berikutnya dikurangi 1 detik (Otomatis dapat tanggal terakhir bulan ini)
		endOfDay = startOfDay.AddDate(0, 1, 0).Add(-time.Second)

	} else { // Format: 2026-06-09 (Harian)
		parsedDate, err := time.ParseInLocation("2006-01-02", tanggal, time.Local)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format rujukan tanggal tidak valid !"})
			return
		}
		startOfDay = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, parsedDate.Location())
		endOfDay = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 23, 59, 59, 999999999, parsedDate.Location())
	}

	transactions, err := h.Repo.GetTransactionsByRange(storeID, startOfDay, endOfDay)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menarik riwayat log laporan transaksi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Daftar riwayat log transaksi berhasil ditarik!",
		"data":    transactions,
	})
}

func (h *RetailHandler) GetDailyClosing(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	var storeID uint
	switch v := storeIDRaw.(type) {
	case float64:
		storeID = uint(v)
	case uint:
		storeID = v
	case int:
		storeID = uint(v)
	}

	tanggal := c.Query("tanggal")
	if tanggal == "" {
		tanggal = time.Now().Format("2006-01-02")
	}

	var startOfDay, endOfDay time.Time

	// 🚀 DETEKSI OTOMATIS BULANAN ATAU HARIAN UNTUK AUDIT CLOSING
	if len(tanggal) == 7 { // Format: 2026-06 (Bulanan)
		parsedMonth, err := time.ParseInLocation("2006-01", tanggal, time.Local)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format bulan audit tidak valid"})
			return
		}
		startOfDay = time.Date(parsedMonth.Year(), parsedMonth.Month(), 1, 0, 0, 0, 0, parsedMonth.Location())
		endOfDay = startOfDay.AddDate(0, 1, 0).Add(-time.Second)

	} else { // Format: 2026-06-09 (Harian)
		parsedDate, err := time.ParseInLocation("2006-01-02", tanggal, time.Local)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal rujukan audit closing tidak valid"})
			return
		}
		startOfDay = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, parsedDate.Location())
		endOfDay = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 23, 59, 59, 999999999, parsedDate.Location())
	}

	closings, err := h.Repo.GetClosingByRange(storeID, startOfDay, endOfDay)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengumpulkan berkas rekapitulasi harian shift"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berkas data rekapitulasi closing berhasil dikumpulkan!",
		"data":    closings,
	})
}
