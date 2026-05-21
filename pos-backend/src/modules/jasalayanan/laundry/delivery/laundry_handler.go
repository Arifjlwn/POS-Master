package delivery

import (
	"encoding/base64"
	"net/http"
	"os"
	"path/filepath"
	"pos-backend/models"
	"pos-backend/src/modules/jasalayanan/laundry/domain"
	"pos-backend/src/modules/jasalayanan/laundry/repository"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LaundryHandler struct {
	Repo repository.LaundryRepository
}

// FIX: Parameter diubah dari repository.Repository menjadi repository.LaundryRepository
func NewLaundryHandler(repo repository.LaundryRepository) *LaundryHandler {
	return &LaundryHandler{Repo: repo}
}

// Helper utility converter base64
func simpanGambarBase64(base64Data string, folder string, filename string) (string, error) {
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
	if err := os.WriteFile(filePath, decodedData, 0644); err != nil {
		return "", err
	}
	return strings.ReplaceAll(filePath, "\\", "/"), nil
}

type KasirInput struct {
	Name     string `json:"name" binding:"required"`
	NoHP     string `json:"no_hp"`
	Password string `json:"password"`
}

func (h *LaundryHandler) GetKasirList(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

	kasirList, err := h.Repo.GetKasirByStoreID(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat data kasir"})
		return
	}

	var response []map[string]interface{}
	for _, k := range kasirList {
		email := ""
		if k.Email != nil {
			email = *k.Email
		}
		response = append(response, map[string]interface{}{
			"id":    k.ID,
			"name":  k.Name,
			"email": email,
			"no_hp": k.NoHP,
		})
	}
	c.JSON(http.StatusOK, response)
}

func (h *LaundryHandler) CreateKasir(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

	var input KasirInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data input tidak lengkap"})
		return
	}

	store, err := h.Repo.GetStoreByID(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengonfirmasi infrastruktur toko"})
		return
	}

	namaKasirBersih := strings.ToLower(strings.ReplaceAll(input.Name, " ", ""))
	namaTokoBersih := strings.ToLower(strings.ReplaceAll(store.NamaToko, " ", ""))
	emailDummy := "kasir." + namaKasirBersih + "@" + namaTokoBersih + ".com"

	passToHash := "kasir123"
	if input.Password != "" {
		passToHash = input.Password
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(passToHash), bcrypt.DefaultCost)

	newKasir := models.User{
		Name:       input.Name,
		Email:      &emailDummy,
		Password:   string(hashedPassword),
		Role:       "kasir",
		StoreID:    &storeID,
		NoHP:       input.NoHP,
		IsVerified: true,
	}

	if err := h.Repo.CreateKasir(&newKasir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendaftarkan kasir. Pastikan nama tidak duplikat."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Kasir berhasil ditambahkan!", "email": emailDummy})
}

func (h *LaundryHandler) DeleteKasir(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))
	kasirID, _ := strconv.Atoi(c.Param("id"))

	if err := h.Repo.DeleteKasir(uint(kasirID), storeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data kasir"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Karyawan berhasil dihapus dari sistem!"})
}

type LaundryItemInput struct {
	ProductID   uint    `json:"product_id" binding:"required"`
	BeratKg     float64 `json:"berat_kg" binding:"required"`
	HargaPerKg  float64 `json:"harga_per_kg" binding:"required"`
	SubTotal    float64 `json:"sub_total" binding:"required"`
	NamaParfum  string  `json:"nama_parfum"`
	HargaParfum float64 `json:"harga_parfum"`
}

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

func (h *LaundryHandler) ProsesCheckoutLaundry(c *gin.Context) {
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

	invoiceCode := "INV/LD/" + time.Now().Format("20060102") + "/" + time.Now().Format("150405")

	var buktiPath, fotoBarangPath string
	if input.PaymentMethod == "QRIS" && input.BuktiTransferBase64 != "" {
		buktiPath, _ = simpanGambarBase64(input.BuktiTransferBase64, "public/uploads/qris", strings.ReplaceAll(invoiceCode, "/", "")+".jpg")
	}
	if input.FotoBarangBase64 != "" {
		fotoBarangPath, _ = simpanGambarBase64(input.FotoBarangBase64, "public/uploads/items", strings.ReplaceAll(invoiceCode, "/", "")+".jpg")
	}

	db := h.Repo.GetDB()
	tx := db.Begin()

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

	if err := h.Repo.CreateTransactionTx(tx, &newTx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat invoice induk"})
		return
	}

	for _, item := range input.Items {
		laundryDetail := domain.TransactionLaundryDetail{
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
			NamaParfum:    item.NamaParfum,
			HargaParfum:   item.HargaParfum,
		}

		if err := h.Repo.CreateLaundryDetailTx(tx, &laundryDetail); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan rincian cucian"})
			return
		}
	}

	existingCustomer, err := h.Repo.FindCustomerByPhone(storeID, input.CustomerPhone)
	if err != nil {
		newCustomer := models.Customer{
			StoreID:    storeID,
			Nama:       input.CustomerName,
			NoWhatsapp: input.CustomerPhone,
		}
		tx.Create(&newCustomer)
	} else {
		existingCustomer.Nama = input.CustomerName
		tx.Save(existingCustomer)
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

func (h *LaundryHandler) AmbilDaftarLayananLaundry(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

	listJasa, err := h.Repo.GetLayananLaundry(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil katalog layanan laundry"})
		return
	}
	c.JSON(http.StatusOK, listJasa)
}

func (h *LaundryHandler) CariPelanggan(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))
	keyword := c.Query("q")

	customers, err := h.Repo.SearchCustomers(storeID, keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencari data pelanggan"})
		return
	}
	c.JSON(http.StatusOK, customers)
}

type PelunasanInput struct {
	MetodeBayar         string `json:"metode_bayar" binding:"required"`
	BuktiTransferBase64 string `json:"bukti_transfer_base64"`
}

func (h *LaundryHandler) LunasiTransaksi(c *gin.Context) {
	trxIDStr := c.Param("id")
	trxID, _ := strconv.Atoi(trxIDStr)
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

	var input PelunasanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Metode pembayaran wajib dipilih"})
		return
	}

	trx, err := h.Repo.GetTransactionByID(uint(trxID), storeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaksi tidak ditemukan"})
		return
	}

	if trx.StatusBayar == "LUNAS" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaksi ini sudah lunas sebelumnya"})
		return
	}

	if input.MetodeBayar == "QRIS" && input.BuktiTransferBase64 != "" {
		buktiPath, _ := simpanGambarBase64(input.BuktiTransferBase64, "public/uploads/qris", strings.ReplaceAll(trx.NoInvoice, "/", "")+"_lunas.jpg")
		trx.BuktiTransfer = buktiPath
	}

	trx.StatusBayar = "LUNAS"
	trx.MetodeBayar = input.MetodeBayar

	if err := h.Repo.UpdateTransaction(trx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal melunasi transaksi"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tagihan berhasil dilunasi!", "data": trx})
}

type InputLayanan struct {
	NamaProduk  string  `json:"nama_produk" binding:"required"`
	HargaJual   float64 `json:"harga_jual" binding:"required"`
	SatuanDasar string  `json:"satuan_dasar" binding:"required"` 
	Estimasi    string  `json:"estimasi"`                    
}

func (h *LaundryHandler) TambahLayananLaundry(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

	var input InputLayanan
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input layanan tidak valid"})
		return
	}

	newLayanan := models.Product{
		StoreID:     storeID,
		NamaProduk:  input.NamaProduk,
		Kategori:    "JASA_LAUNDRY",
		HargaJual:   input.HargaJual,
		SatuanDasar: input.SatuanDasar,
		Estimasi:    input.Estimasi,
		Stok:        0, 
	}

	if err := h.Repo.CreateLayanan(&newLayanan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan layanan baru"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Layanan berhasil ditambahkan", "data": newLayanan})
}

func (h *LaundryHandler) HapusLayananLaundry(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))
	productID, _ := strconv.Atoi(c.Param("id"))

	if err := h.Repo.DeleteLayanan(uint(productID), storeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus layanan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Layanan berhasil deleted"})
}

func (h *LaundryHandler) EditLayananLaundry(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))
	productID, _ := strconv.Atoi(c.Param("id"))

	var input InputLayanan
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input layanan tidak valid"})
		return
	}

	layanan, err := h.Repo.GetLayananByID(uint(productID), storeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Layanan tidak ditemukan"})
		return
	}

	layanan.NamaProduk = input.NamaProduk
	layanan.HargaJual = input.HargaJual
	layanan.SatuanDasar = input.SatuanDasar
	layanan.Estimasi = input.Estimasi

	if err := h.Repo.UpdateLayanan(layanan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui layanan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Layanan berhasil diperbarui!"})
}

type PerfumeInput struct {
	Nama  string  `json:"nama" binding:"required"`
	Harga float64 `json:"harga"`
}

func (h *LaundryHandler) GetPerfumes(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

	perfumes, err := h.Repo.GetPerfumesByStoreID(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data parfum"})
		return
	}
	c.JSON(http.StatusOK, perfumes)
}

func (h *LaundryHandler) CreatePerfume(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

	var input PerfumeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data input tidak valid"})
		return
	}

	newPerfume := domain.Perfume{
		StoreID: storeID,
		Nama:    input.Nama,
		Harga:   input.Harga,
		Status:  "Tersedia",
	}

	if err := h.Repo.CreatePerfume(&newPerfume); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan varian parfum baru"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Parfum berhasil ditambahkan!", "data": newPerfume})
}

func (h *LaundryHandler) DeletePerfume(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))
	perfumeID, _ := strconv.Atoi(c.Param("id"))

	if err := h.Repo.DeletePerfume(uint(perfumeID), storeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus parfum"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Varian parfum berhasil dihapus!"})
}

type TransactionReportResponse struct {
	models.Transaction
	Invoice       string    `json:"invoice"`
	Pelanggan     string    `json:"pelanggan"`
	Whatsapp      string    `json:"whatsapp"`
	Layanan       string    `json:"layanan"`
	BeratKg       float64   `json:"berat_kg"`
	SubTotal      float64   `json:"sub_total"`
	EstimasiWaktu time.Time `json:"estimasi_waktu"`
}

func (h *LaundryHandler) GetLaporan(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

	transactions, err := h.Repo.GetAllTransactions(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data laporan"})
		return
	}

	var reportData []TransactionReportResponse
	var tunai, qris, debit, piutang, omset float64
	totalOrder := len(transactions)

	for _, trx := range transactions {
		if trx.StatusBayar == "BELUM_LUNAS" {
			piutang += trx.TotalHarga
		} else {
			omset += trx.TotalHarga
			switch trx.MetodeBayar {
			case "TUNAI":
				tunai += trx.TotalHarga
			case "QRIS":
				qris += trx.TotalHarga
			case "DEBIT":
				debit += trx.TotalHarga
			}
		}

		detail, err := h.Repo.GetLaundryDetailByTxID(trx.ID)
		var namaPelanggan, noWhatsapp string
		var productID uint
		var beratKg, subTotalDetail float64
		var estimasiWaktu time.Time

		if err == nil {
			namaPelanggan = detail.NamaPelanggan
			noWhatsapp = detail.NoWhatsapp
			productID = detail.ProductID
			beratKg = detail.BeratKg
			subTotalDetail = detail.SubTotal
			estimasiWaktu = detail.EstimasiWaktu
		}

		layananName := "Paket Laundry"
		if productID > 0 {
			if prod, err := h.Repo.GetProductByIDSimple(productID); err == nil {
				layananName = prod.NamaProduk
			}
		}

		reportData = append(reportData, TransactionReportResponse{
			Transaction:   trx,
			Invoice:       trx.NoInvoice,
			Pelanggan:     namaPelanggan,
			Whatsapp:      noWhatsapp,
			Layanan:       layananName,
			BeratKg:       beratKg,
			SubTotal:      subTotalDetail,
			EstimasiWaktu: estimasiWaktu,
		})
	}

	avg := 0.0
	if totalOrder > 0 {
		avg = omset / float64(totalOrder)
	}

	c.JSON(http.StatusOK, gin.H{
		"ringkasan": gin.H{
			"total_omset": omset,
			"total_order": totalOrder,
			"rata_rata":   avg,
			"tunai":       tunai,
			"qris":        qris,
			"debit":       debit,
			"piutang":     piutang,
		},
		"transaksi": reportData,
	})
}

type UpdateSettingInput struct {
	NamaToko      string `json:"nama_toko" binding:"required"`
	Telepon       string `json:"telepon" binding:"required"`
	Alamat        string `json:"alamat" binding:"required"`
	PaymentType   string `json:"payment_type" binding:"required"`
	QrisBase64    string `json:"qris_base64"`
	ReceiptFooter string `json:"receipt_footer"`
}

func (h *LaundryHandler) GetSettingToko(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

	store, err := h.Repo.GetStoreByID(storeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data toko tidak ditemukan"})
		return
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	if store.QrisImage != "" {
		store.QrisImage = baseURL + "/" + store.QrisImage
	}

	c.JSON(http.StatusOK, gin.H{
		"id":             store.ID,
		"nama_toko":      store.NamaToko,
		"business_type":  store.BusinessType,
		"telepon":        store.Telepon,
		"alamat":         store.Alamat,
		"payment_type":   store.PaymentType,
		"qris_image":     store.QrisImage,
		"receipt_footer": store.ReceiptFooter,
	})
}

func (h *LaundryHandler) UpdateSettingToko(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

	var input UpdateSettingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data pengaturan tidak valid"})
		return
	}

	db := h.Repo.GetDB()
	tx := db.Begin()

	store, err := h.Repo.GetStoreByID(storeID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Toko tidak ditemukan"})
		return
	}

	if input.PaymentType == "PRIBADI" && input.QrisBase64 != "" && !strings.HasPrefix(input.QrisBase64, "http") {
		// Memakai fungsi utility simpanGambarBase64 yang sudah ada di atas file ini
		qrisPath, err := simpanGambarBase64(input.QrisBase64, "public/uploads/qris_toko", "store_qris.jpg")
		if err == nil {
			store.QrisImage = qrisPath
		}
	}

	store.NamaToko = input.NamaToko
	store.Telepon = input.Telepon
	store.Alamat = input.Alamat
	store.PaymentType = input.PaymentType
	store.ReceiptFooter = input.ReceiptFooter

	if err := h.Repo.UpdateStoreTx(tx, store); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui profil toko"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"status": "sukses", "message": "Pengaturan toko berhasil diperbarui!"})
}

func (h *LaundryHandler) AmbilDataTracking(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

	results, err := h.Repo.GetTrackingCucian(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data tracking"})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h *LaundryHandler) UpdateStatusCucian(c *gin.Context) {
	trxIDStr := c.Param("id")
	trxID, _ := strconv.Atoi(trxIDStr)

	var input struct {
		StatusPesanan string `json:"status_pesanan" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format status tidak valid"})
		return
	}

	if err := h.Repo.UpdateStatusDetailCucian(uint(trxID), input.StatusPesanan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB Error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status berhasil diupdate!"})
}