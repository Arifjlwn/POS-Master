package delivery

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"pos-backend/models"
	"pos-backend/src/modules/retail/domain"
	"pos-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause" // 🚀 FIX COMPILER: Import wajib biar row locking Retur jalan !
)

// ==========================================
// 📦 STRUCT INPUT VALIDATION
// ==========================================

type ProductInput struct {
	SKU              *string `form:"sku"`
	NamaProduk       string  `form:"nama_produk" binding:"required"`
	Kategori         string  `form:"kategori"`
	HargaModal       float64 `form:"harga_modal"`
	HargaJual        float64 `form:"harga_jual" binding:"required"`
	Stok             int     `form:"stok"`
	SatuanDasar      string  `form:"satuan_dasar"`
	SatuanBesar      string  `form:"satuan_besar"`
	IsiPerBesar      int     `form:"isi_per_besar"`
	HargaJualBesar   float64 `form:"harga_jual_besar"`
	IsNestedUom      bool    `form:"is_nested_uom"`
	SatuanTengah     string  `form:"satuan_tengah"`
	IsiBesarKeTengah int     `form:"isi_besar_ke_tengah"`
	IsiTengahKeDasar int     `form:"isi_tengah_ke_dasar"`
	HargaJualTengah  float64 `form:"harga_jual_tengah"`
}

type PurchaseInput struct {
	SupplierName string `json:"supplier_name" binding:"required"`
	NoFaktur     string `json:"no_faktur" binding:"required"`
	Items        []struct {
		ProductID  uint    `json:"product_id" binding:"required"`
		QtyMasuk   int     `json:"qty_masuk" binding:"required,gt=0"`
		HargaModal float64 `json:"harga_modal" binding:"required,min=0"` // ◄ Kunci tag json:"harga_modal"!
	} `json:"items"`
}

func getStoreID(c *gin.Context) uint {
	storeIDRaw, _ := c.Get("store_id")
	switch v := storeIDRaw.(type) {
	case float64:
		return uint(v)
	case uint:
		return v
	case int:
		return uint(v)
	}
	return 0
}

// ==========================================
// 📦 PRODUCT HANDLERS (CRUD CATALOG MODULE)
// ==========================================

func (h *RetailHandler) CreateProduct(c *gin.Context) {
	storeID := getStoreID(c)
	if storeID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi identitas toko tidak valid !"})
		return
	}

	var input ProductInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Form isian produk tidak valid: " + err.Error()})
		return
	}

	// 🛡️ SECURITY PATROL: Blokir mutlak angka minus masuk master database
	if input.HargaJual < 0 || input.HargaModal < 0 || input.Stok < 0 || input.HargaJualTengah < 0 || input.HargaJualBesar < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nilai nominal harga penjualan/stok dilarang bernilai minus !"})
		return
	}

	if input.SatuanDasar == "" {
		input.SatuanDasar = "PCS"
	}

	db := h.Repo.GetDB()
	var store models.Store
	if err := db.Select("public_id").First(&store, storeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data internal infrastruktur merchant gagal diverifikasi"})
		return
	}

	var imageURL string
	bucketName := os.Getenv("SUPABASE_BUCKET_NAME")

	if file, errFile := c.FormFile("gambar"); errFile == nil {
		if file.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran aset gambar maksimal 5 MB !"})
			return
		}
		contentType := file.Header.Get("Content-Type")
		fileSrc, err := file.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membuka file stream gambar"})
			return
		}
		defer fileSrc.Close()

		cleanName := regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(strings.ToLower(input.NamaProduk), "-")
		cleanName = strings.Trim(cleanName, "-")
		remotePath := fmt.Sprintf("stores/%s/products/%s_%d", store.PublicID, cleanName, time.Now().Unix())
		urlResult, errUpload := utils.UploadToSupabase(fileSrc, file.Filename, contentType, bucketName, remotePath)
		if errUpload != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Koneksi Supabase Storage terputus, gagal simpan gambar"})
			return
		}
		imageURL = urlResult
	}

	product := models.Product{
		StoreID: storeID, SKU: input.SKU, NamaProduk: input.NamaProduk, Kategori: input.Kategori,
		ProductType: "retail",
		Estimasi:    "Standar",
		HargaModal:  input.HargaModal, HargaJual: input.HargaJual, Stok: input.Stok, Gambar: imageURL,
		SatuanDasar: input.SatuanDasar, SatuanBesar: input.SatuanBesar, IsiPerBesar: input.IsiPerBesar,
		HargaJualBesar: input.HargaJualBesar, IsNestedUom: input.IsNestedUom, SatuanTengah: input.SatuanTengah,
		IsiBesarKeTengah: input.IsiBesarKeTengah, IsiTengahKeDasar: input.IsiTengahKeDasar, HargaJualTengah: input.HargaJualTengah,
	}

	if err := h.Repo.CreateProductGlobal(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan produk. Kode Barcode/SKU sudah terpakai di toko Anda !"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Aset produk baru berhasil ditambahkan! 📦", "data": product})
}

func (h *RetailHandler) GetProducts(c *gin.Context) {
	storeID := getStoreID(c)
	search := strings.TrimSpace(c.Query("search")) // FIX CLEANING INDEX SPACE
	category := strings.TrimSpace(c.Query("category"))
	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	usePagination := false
	limit, offset := 10, 0

	if pageStr != "" || limitStr != "" {
		usePagination = true
		if pageStr == "" {
			pageStr = "1"
		}
		if limitStr == "" {
			limitStr = "10"
		}
		p, _ := strconv.Atoi(pageStr)
		l, _ := strconv.Atoi(limitStr)
		limit = l
		offset = (p - 1) * limit
	}

	products, totalItems, err := h.Repo.GetProductsCatalog(storeID, search, category, limit, offset, usePagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat katalog berkas data produk"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Katalog produk berhasil dimuat!", "total_items": totalItems, "data": products})
}

func (h *RetailHandler) UpdateProduct(c *gin.Context) {
	storeID := getStoreID(c)
	role := strings.ToLower(c.MustGet("role").(string))

	if role != "owner" && role != "manager" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak! Menu perubahan produk khusus level Owner/Manager !"})
		return
	}

	publicID := c.Param("id")
	db := h.Repo.GetDB()

	var product models.Product
	if err := db.Where("public_id = ? AND store_id = ?", publicID, storeID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kode produk tidak terdaftar dalam database master!"})
		return
	}

	if v := c.PostForm("nama_produk"); v != "" {
		product.NamaProduk = v
	}
	if v := c.PostForm("kategori"); v != "" {
		product.Kategori = v
	}
	if sku := c.PostForm("sku"); sku != "" {
		product.SKU = &sku
	}

	if hModal, err := strconv.ParseFloat(c.PostForm("harga_modal"), 64); err == nil && hModal >= 0 {
		product.HargaModal = hModal
	}
	if hJual, err := strconv.ParseFloat(c.PostForm("harga_jual"), 64); err == nil && hJual >= 0 {
		product.HargaJual = hJual
	}
	if stok, err := strconv.Atoi(c.PostForm("stok")); err == nil && stok >= 0 {
		product.Stok = stok
	}

	if sDasar := c.PostForm("satuan_dasar"); sDasar != "" {
		product.SatuanDasar = sDasar
	}
	if sBesar := c.PostForm("satuan_besar"); sBesar != "" {
		product.SatuanBesar = sBesar
	}
	if iPerBesar, err := strconv.Atoi(c.PostForm("isi_per_besar")); err == nil {
		product.IsiPerBesar = iPerBesar
	}
	if hJualBesar, err := strconv.ParseFloat(c.PostForm("harga_jual_besar"), 64); err == nil {
		product.HargaJualBesar = hJualBesar
	}

	if isNested := c.PostForm("is_nested_uom"); isNested != "" {
		product.IsNestedUom = (isNested == "true")
	}
	if sTengah := c.PostForm("satuan_tengah"); sTengah != "" {
		product.SatuanTengah = sTengah
	}
	if ibt, err := strconv.Atoi(c.PostForm("isi_besar_ke_tengah")); err == nil {
		product.IsiBesarKeTengah = ibt
	}
	if itd, err := strconv.Atoi(c.PostForm("isi_tengah_ke_dasar")); err == nil {
		product.IsiTengahKeDasar = itd
	}
	if hJualTengah, err := strconv.ParseFloat(c.PostForm("harga_jual_tengah"), 64); err == nil {
		product.HargaJualTengah = hJualTengah
	}

	// =========================================================================
	// 🚀 FIX SYNTAX ENGINE: Tambahkan pemisah titik koma (;) setelah assignment !
	// =========================================================================
	var store models.Store
	if err := db.Select("public_id").First(&store, storeID).Error; err == nil {
		bucketName := os.Getenv("SUPABASE_BUCKET_NAME")
		if file, errFile := c.FormFile("gambar"); errFile == nil {
			if file.Size <= 5*1024*1024 {
				contentType := file.Header.Get("Content-Type")
				if fileSrc, errOpen := file.Open(); errOpen == nil {
					cleanName := regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(strings.ToLower(product.NamaProduk), "-")
					cleanName = strings.Trim(cleanName, "-")

					remotePath := fmt.Sprintf("stores/%s/products/%s_%d", store.PublicID, cleanName, time.Now().Unix())
					urlResult, errUpload := utils.UploadToSupabase(fileSrc, file.Filename, contentType, bucketName, remotePath)
					if errUpload == nil {
						product.Gambar = urlResult
					}
					fileSrc.Close()
				}
			}
		}
	}

	if err := h.Repo.SaveProduct(db, &product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update , nomor barcode/SKU sudah dimiliki produk lain!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Spesifikasi produk berhasil diubah! ✏️", "data": product})
}

func (h *RetailHandler) DeleteProduct(c *gin.Context) {
	storeID := getStoreID(c)
	role := strings.ToLower(c.MustGet("role").(string))

	if role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak! Tim kasir dilarang membumihanguskan aset gudang !"})
		return
	}

	publicID := c.Param("id")
	db := h.Repo.GetDB()

	var product models.Product
	if err := db.Where("public_id = ? AND store_id = ?", publicID, storeID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan di database"})
		return
	}

	if err := h.Repo.DeleteProductGlobal(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus produk"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Barang berhasil dihapus dari gudang! 🗑️"})
}

func (h *RetailHandler) GetCategories(c *gin.Context) {
	storeID := getStoreID(c)
	categories, err := h.Repo.GetDistinctCategories(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengumpulkan daftar kategori"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

// ==========================================
// 📦 BULK COMPLIANCE: EXPORT IMPORT CSV DATA
// ==========================================

func (h *RetailHandler) ExportProducts(c *gin.Context) {
	storeID := getStoreID(c)
	products, err := h.Repo.GetAllProductsForExport(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengunduh berkas data ekspor"})
		return
	}

	b := &bytes.Buffer{}
	w := csv.NewWriter(b)
	w.Comma = '|'
	w.Write([]string{"SKU", "Nama Produk", "Kategori", "Tipe Produk", "Estimasi", "Harga Modal", "Harga Jual", "Stok", "Satuan Terkecil", "Satuan Tengah", "Satuan Besar", "Isi Per Besar", "Harga Jual Besar", "Harga Jual Tengah"})

	for _, p := range products {
		sku := ""
		if p.SKU != nil {
			sku = *p.SKU
		}
		w.Write([]string{
			sku, p.NamaProduk, p.Kategori, p.ProductType, p.Estimasi, fmt.Sprintf("%.0f", p.HargaModal), fmt.Sprintf("%.0f", p.HargaJual),
			fmt.Sprintf("%d", p.Stok), p.SatuanDasar, p.SatuanTengah, p.SatuanBesar, fmt.Sprintf("%d", p.IsiPerBesar), fmt.Sprintf("%.0f", p.HargaJualBesar), fmt.Sprintf("%.0f", p.HargaJualTengah),
		})
	}
	w.Flush()

	c.Header("Content-Disposition", "attachment; filename=katalog_produk_pos.csv")
	c.Data(http.StatusOK, "text/csv", b.Bytes())
}

func (h *RetailHandler) ImportProducts(c *gin.Context) {
	storeID := getStoreID(c)
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Berkas file CSV tidak ditemukan"})
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '|'
	_, _ = reader.Read() // Skip Header

	records, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca struktur pembatas file CSV "})
		return
	}

	db := h.Repo.GetDB()
	err = db.Transaction(func(tx *gorm.DB) error {
		for _, row := range records {
			if len(row) < 7 {
				continue
			}
			sku := strings.TrimSpace(row[0])
			nama := strings.TrimSpace(row[1])
			kategori := strings.TrimSpace(row[2])

			pType := "retail"
			if len(row) > 3 && strings.TrimSpace(row[3]) != "" {
				pType = strings.TrimSpace(row[3])
			}

			estimasi := "Standar"
			if len(row) > 4 && strings.TrimSpace(row[4]) != "" {
				estimasi = strings.TrimSpace(row[4])
			}

			modal, _ := strconv.ParseFloat(row[5], 64)
			jual, _ := strconv.ParseFloat(row[6], 64)
			stok, _ := strconv.Atoi(row[7])
			dasar := row[8]

			tengah, besar := "", ""
			isi := 0
			var jualBesar, jualTengah float64
			if len(row) > 9 {
				tengah = row[9]
			}
			if len(row) > 10 {
				besar = row[10]
			}
			if len(row) > 11 {
				isi, _ = strconv.Atoi(row[11])
			}
			if len(row) > 12 {
				jualBesar, _ = strconv.ParseFloat(row[12], 64)
			}
			if len(row) > 13 {
				jualTengah, _ = strconv.ParseFloat(row[13], 64)
			}

			if nama == "" {
				continue
			}

			var product models.Product
			var res *gorm.DB
			if sku != "" {
				res = tx.Where("sku = ? AND store_id = ?", sku, storeID).First(&product)
			}

			if sku != "" && res != nil && res.Error == nil {
				product.NamaProduk = nama
				product.Kategori = kategori
				product.ProductType = pType
				product.Estimasi = estimasi
				if modal >= 0 {
					product.HargaModal = modal
				}
				if jual >= 0 {
					product.HargaJual = jual
				}
				if stok >= 0 {
					product.Stok = stok
				}
				product.SatuanDasar = dasar
				product.SatuanTengah = tengah
				product.SatuanBesar = besar
				product.IsiPerBesar = isi
				product.HargaJualBesar = jualBesar
				product.HargaJualTengah = jualTengah
				tx.Save(&product)
			} else {
				var skuPtr *string
				if sku != "" {
					skuPtr = &sku
				}

				newProduct := models.Product{
					StoreID: storeID, SKU: skuPtr, NamaProduk: nama, Kategori: kategori,
					ProductType: pType, Estimasi: estimasi, HargaModal: modal, HargaJual: jual, Stok: stok,
					SatuanDasar: dasar, SatuanTengah: tengah, SatuanBesar: besar, IsiPerBesar: isi,
					HargaJualBesar: jualBesar, HargaJualTengah: jualTengah,
				}
				tx.Create(&newProduct)
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses transaksi impor data: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Berhasil memproses impor " + strconv.Itoa(len(records)) + " data produk ke database"})
}

// =========================================================================
// 🚀 JALUR BARU: PURCHASE LPB FAKTUR MASUK (MOVING AVERAGE IMPLEMENTATION)
// =========================================================================

func (h *RetailHandler) CreateLPB(c *gin.Context) {
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

	var input PurchaseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data faktur penerimaan LPB tidak valid!"})
		return
	}
	if len(input.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Keranjang data barang masuk kosong !"})
		return
	}

	var totalHargaFaktur float64
	var details []domain.PurchaseDetail

	for _, item := range input.Items {
		// 🛡️ SECURITY GUARD FINANSIAL: Tolak mentah-mentah kuantitas minus/rusak dari luar  !
		if item.QtyMasuk <= 0 || item.HargaModal < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Kuantitas masuk harus lebih dari 0 dan harga modal tidak boleh minus!"})
			return
		}

		subTotalItem := float64(item.QtyMasuk) * item.HargaModal
		totalHargaFaktur += subTotalItem

		details = append(details, domain.PurchaseDetail{
			ProductID:  item.ProductID,
			QtyMasuk:   item.QtyMasuk,
			HargaModal: item.HargaModal,
			SubTotal:   subTotalItem,
		})
	}

	purchase := domain.Purchase{
		PublicID:     utils.GenerateULID(),
		StoreID:      storeID,
		UserID:       userID,
		SupplierName: input.SupplierName,
		NoFaktur:     input.NoFaktur,
		TotalItem:    len(input.Items),
		TotalHarga:   totalHargaFaktur,
		StatusBayar:  "LUNAS",
		Details:      details,
	}

	db := h.Repo.GetDB()
	if err := h.Repo.CreatePurchaseWithMovingAverage(db, &purchase); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung valuasi HPP Moving Average: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Faktur Penerimaan Barang (LPB) diproses, modal HPP Moving Average diperbarui!"})
}

func (h *RetailHandler) GetLaporanInbound(c *gin.Context) {
	// 1. Ambil StoreID pakai helper lu yang udah ada
	storeID := getStoreID(c)
	if storeID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi identitas toko tidak valid !"})
		return
	}

	// 2. Tarik datanya dari UseCase/Repo
	// (Pastikan fungsi GetLaporanInbound ini udah lu bikin di retail_repo.go / retail_usecase.go sesuai instruksi sebelumnya)
	reports, err := h.Repo.GetLaporanInbound(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Sirkuit gagal mengambil laporan penerimaan barang: " + err.Error()})
		return
	}

	// 3. 🚀 Output Kasta Tertinggi ke Frontend Vue
	c.JSON(http.StatusOK, gin.H{
		"status": "sukses",
		"data":   reports,
	})
}

// =========================================================================
// 🚀 JALUR BARU: RETUR BARANG ADAPTIF ENGINE
// =========================================================================

type ReturnItem struct {
	ProductID uint   `json:"product_id" binding:"required"`
	Qty       int    `json:"qty" binding:"required,gt=0"`
	Alasan    string `json:"alasan" binding:"required"`
	Catatan   string `json:"catatan"`
}

type ReturnInputBatch struct {
	Items []ReturnItem `json:"items" binding:"required,min=1"`
}

func (h *RetailHandler) CreateReturn(c *gin.Context) {
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

	var input ReturnInputBatch
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data keranjang retur tidak valid !"})
		return
	}

	returnNo := fmt.Sprintf("RET-%s-%d", time.Now().Format("060102150405"), userID)
	db := h.Repo.GetDB()

	err := db.Transaction(func(tx *gorm.DB) error {
		var newReturns []domain.ProductReturn
		for _, item := range input.Items {
			var product models.Product

			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&product, "id = ? AND store_id = ?", item.ProductID, storeID).Error; err != nil {
				return fmt.Errorf("produk ID %d tidak terdaftar di rak toko ", item.ProductID)
			}

			if product.Stok < item.Qty {
				return fmt.Errorf("stok %s tidak mencukupi untuk diretur (Sisa di DB: %d PCS)", product.NamaProduk, product.Stok)
			}

			if err := h.Repo.UpdateProductStokExpr(tx, item.ProductID, storeID, item.Qty); err != nil {
				return fmt.Errorf("gagal memotong saldo kuantitas kartu stok produk ID %d", item.ProductID)
			}

			newReturns = append(newReturns, domain.ProductReturn{
				PublicID:  utils.GenerateULID(),
				ReturnNo:  returnNo,
				StoreID:   storeID,
				ProductID: item.ProductID,
				UserID:    userID,
				Qty:       item.Qty,
				Alasan:    item.Alasan,
				Catatan:   item.Catatan,
			})
		}

		return h.Repo.CreateProductReturns(tx, newReturns)
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Berita Acara Retur berhasil diproses!", "return_no": returnNo})
}

func (h *RetailHandler) GetReturns(c *gin.Context) {
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

	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	limit, offset := 0, 0
	if pageStr != "" && limitStr != "" {
		p, _ := strconv.Atoi(pageStr)
		l, _ := strconv.Atoi(limitStr)
		limit = l
		offset = (p - 1) * limit
	}

	returns, totalItems, err := h.Repo.GetReturnsHistory(storeID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat arsip data lembar retur"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data retur berhasil dimuat!", "total_items": totalItems, "data": returns})
}
