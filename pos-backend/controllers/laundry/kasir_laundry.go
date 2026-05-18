package laundry

import (
	"encoding/base64"
	"net/http"
	"os"
	"path/filepath"
	"pos-backend/config"
	"pos-backend/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Struct untuk menampung item cucian dari Vue Frontend
type LaundryItemInput struct {
	ProductID   uint    `json:"product_id" binding:"required"`
	BeratKg     float64 `json:"berat_kg" binding:"required"`
	HargaPerKg  float64 `json:"harga_per_kg" binding:"required"`
	SubTotal    float64 `json:"sub_total" binding:"required"`
	
	// 🚀 TAMBAHAN PAYLOAD ADD-ON PARFUM DARI VUE
	NamaParfum  string  `json:"nama_parfum"`
	HargaParfum float64 `json:"harga_parfum"`
}

// Struct utama untuk menerima payload checkout laundry
type CheckoutLaundryInput struct {
	CustomerName        string             `json:"customer_name" binding:"required"`
	CustomerPhone       string             `json:"customer_phone" binding:"required"`
	EstimasiSelesai     string             `json:"estimasi_selesai" binding:"required"`
	Items               []LaundryItemInput `json:"items" binding:"required"`
	TotalAmount         float64            `json:"total_amount" binding:"required"`
	PaymentMethod       string             `json:"payment_method" binding:"required"`
	PaymentStatus       string             `json:"payment_status" binding:"required"`
	FotoBarangBase64    string             `json:"foto_barang_base64"`
	BuktiTransferBase64 string             `json:"bukti_transfer_base64"`
}

// 🚀 FUNGSI SAKTI UBAH BASE64 JADI GAMBAR .JPG
func SimpanGambarBase64(base64Data string, folder string, filename string) (string, error) {
	if base64Data == "" {
		return "", nil
	}

	parts := strings.Split(base64Data, ",")
	var pureBase64 string
	if len(parts) > 1 {
		pureBase64 = parts[1]
	} else {
		pureBase64 = parts[0]
	}

	decodedData, err := base64.StdEncoding.DecodeString(pureBase64)
	if err != nil {
		return "", err
	}

	os.MkdirAll(folder, os.ModePerm)

	filePath := filepath.Join(folder, filename)
	err = os.WriteFile(filePath, decodedData, 0644)
	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(filePath, "\\", "/"), nil
}

// 🚀 1. FUNGSI UNTUK MEMPROSES TRANSAKSI MASUK
func ProsesCheckoutLaundry(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	userIDRaw, _ := c.Get("user_id")
	
	storeID := uint(storeIDRaw.(float64))
	userID := uint(userIDRaw.(float64))

	var input CheckoutLaundryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data transaksi tidak valid"})
		return
	}

	estimasiTime, err := time.Parse("2006-01-02", input.EstimasiSelesai)
	if err != nil {
		estimasiTime = time.Now().Add(time.Hour * 48) 
	}

	tx := config.DB.Begin()

	invoiceCode := "INV/LD/" + time.Now().Format("20060102") + "/" + time.Now().Format("150405")

	var buktiPath string
	if input.PaymentMethod == "QRIS" && input.BuktiTransferBase64 != "" {
		buktiPath, _ = SimpanGambarBase64(input.BuktiTransferBase64, "public/uploads/qris", strings.ReplaceAll(invoiceCode, "/", "")+".jpg")
	}

	var fotoBarangPath string
	if input.FotoBarangBase64 != "" {
		fotoBarangPath, _ = SimpanGambarBase64(input.FotoBarangBase64, "public/uploads/items", strings.ReplaceAll(invoiceCode, "/", "")+".jpg")
	}

	newTx := models.Transaction{
		SessionID:     1,
		StoreID:       storeID,
		UserID:        userID,
		NoInvoice:     invoiceCode,
		SubTotal:      input.TotalAmount,
		Pajak:         0,
		Pembulatan:    0,
		TotalHarga:    input.TotalAmount,
		MetodeBayar:   input.PaymentMethod,
		StatusBayar:   input.PaymentStatus,
		NominalBayar:  input.TotalAmount,
		Kembalian:     0,
		BuktiTransfer: buktiPath, 
	}

	if err := tx.Create(&newTx).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat invoice induk"})
		return
	}

	for _, item := range input.Items {
		laundryDetail := models.TransactionLaundryDetail{
			TransactionID: newTx.ID,
			ProductID:     item.ProductID,
			NamaPelanggan: input.CustomerName,
			NoWhatsapp:    input.CustomerPhone,
			BeratKg:       item.BeratKg,
			HargaPerKg:    item.HargaPerKg,
			SubTotal:      item.SubTotal,
			StatusCucian:  "ANTRI", 
			EstimasiWaktu: estimasiTime,
			FotoBarang:    fotoBarangPath,
			// 🚀 TANGKAP DATA PARFUM DAN SIMPAN KE DATABASE
			NamaParfum:    item.NamaParfum,
			HargaParfum:   item.HargaParfum,
		}

		if err := tx.Create(&laundryDetail).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan rincian"})
			return
		}
	}

	var existingCustomer models.Customer
	if err := tx.Where("store_id = ? AND no_whatsapp = ?", storeID, input.CustomerPhone).First(&existingCustomer).Error; err != nil {
		newCustomer := models.Customer{
			StoreID:    storeID,
			Nama:       input.CustomerName,
			NoWhatsapp: input.CustomerPhone,
		}
		tx.Create(&newCustomer)
	} else {
		existingCustomer.Nama = input.CustomerName
		tx.Save(&existingCustomer)
	}

	tx.Commit()

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	fotoPublicUrl := ""
	if fotoBarangPath != "" {
		fotoPublicUrl = baseURL + "/" + fotoBarangPath
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       "sukses",
		"message":      "Transaksi laundry berhasil disimpan!",
		"invoice_code": invoiceCode,
		"foto_url":     fotoPublicUrl,
	})
}

// 🚀 2. FUNGSI UNTUK MENAMPILKAN KATALOG JASA CUCIAN DI KASIR
func AmbilDaftarLayananLaundry(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

	var listJasa []models.Product

	if err := config.DB.Where("store_id = ?", storeID).Find(&listJasa).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil katalog layanan laundry"})
		return
	}

	c.JSON(http.StatusOK, listJasa)
}

// 🚀 3. FUNGSI PENCARIAN PELANGGAN (LIVE SEARCH)
func CariPelanggan(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))
	keyword := c.Query("q")

	var customers []models.Customer
	query := config.DB.Where("store_id = ?", storeID)

	if keyword != "" {
		query = query.Where("nama ILIKE ? OR no_whatsapp LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Order("updated_at desc").Limit(5).Find(&customers)
	c.JSON(http.StatusOK, customers)
}

// =========================================================================
// 🚀 4. FUNGSI BARU: LUNASI PIUTANG KASIR (Dipanggil dari Halaman Laporan)
// =========================================================================

type PelunasanInput struct {
	MetodeBayar string `json:"metode_bayar" binding:"required"`
}

func LunasiTransaksi(c *gin.Context) {
	trxID := c.Param("id")
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

	var input PelunasanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Metode pembayaran wajib dipilih"})
		return
	}

	var trx models.Transaction
	if err := config.DB.Where("id = ? AND store_id = ?", trxID, storeID).First(&trx).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaksi tidak ditemukan"})
		return
	}

	if trx.StatusBayar == "LUNAS" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaksi ini sudah lunas sebelumnya"})
		return
	}

	// Ubah status dan catat dia lunas pakai apa
	trx.StatusBayar = "LUNAS"
	trx.MetodeBayar = input.MetodeBayar

	if err := config.DB.Save(&trx).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal melunasi transaksi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tagihan berhasil dilunasi!", "data": trx})
}