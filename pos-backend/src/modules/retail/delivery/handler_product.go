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
	"pos-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	HargaJualTengah  float64 `form:"harga_jual_tengah"` // Ini tetep ada buat harga grosir layer 3
}

// Helper internal biar rapi narik StoreID tanpa ngerusak pemandangan handler utama bray
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
// 📦 PRODUCT HANDLERS
// ==========================================

func (h *RetailHandler) CreateProduct(c *gin.Context) {
	storeID := getStoreID(c)
	if storeID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi toko tidak valid bray"})
		return
	}

	var input ProductInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Form tidak valid: " + err.Error()})
		return
	}

	// 🛡️ SECURITY PATROL: Cegah nominal rusak atau minus masuk DB
	if input.HargaJual < 0 || input.HargaModal < 0 || input.Stok < 0 || input.HargaJualTengah < 0 || input.HargaJualBesar < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nilai harga atau stok tidak boleh bernilai minus!"})
		return
	}

	if input.SatuanDasar == "" {
		input.SatuanDasar = "PCS"
	}

	// 🔒 ANTI-IDOR SHIELD: Ambil public_id toko
	db := h.Repo.GetDB()
	var store models.Store
	if err := db.Select("public_id").First(&store, storeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data internal infrastruktur toko gagal diverifikasi"})
		return
	}

	var imageURL string
	bucketName := os.Getenv("SUPABASE_BUCKET_NAME")

	if file, errFile := c.FormFile("gambar"); errFile == nil {
		if file.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran gambar maksimal 5 MB bray"})
			return
		}
		contentType := file.Header.Get("Content-Type")
		fileSrc, err := file.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membaca berkas gambar"})
			return
		}
		defer fileSrc.Close()

		cleanName := regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(strings.ToLower(input.NamaProduk), "-")
		cleanName = strings.Trim(cleanName, "-")
		remotePath := fmt.Sprintf("stores/%s/products/%s_%d", store.PublicID, cleanName, time.Now().Unix())
		urlResult, errUpload := utils.UploadToSupabase(fileSrc, file.Filename, contentType, bucketName, remotePath)
		if errUpload != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengunggah aset gambar produk"})
			return
		}
		imageURL = urlResult
	}

	// 🚀 HARDCODE RETAIL DI SINI: Tanpa manggil input.ProductType lagi bray!
	product := models.Product{
		StoreID: storeID, SKU: input.SKU, NamaProduk: input.NamaProduk, Kategori: input.Kategori,
		ProductType: "retail",  // Hardcode tipe produk
		Estimasi:    "Standar", // Hardcode estimasi
		HargaModal:  input.HargaModal, HargaJual: input.HargaJual, Stok: input.Stok, Gambar: imageURL,
		SatuanDasar: input.SatuanDasar, SatuanBesar: input.SatuanBesar, IsiPerBesar: input.IsiPerBesar,
		HargaJualBesar: input.HargaJualBesar, IsNestedUom: input.IsNestedUom, SatuanTengah: input.SatuanTengah,
		IsiBesarKeTengah: input.IsiBesarKeTengah, IsiTengahKeDasar: input.IsiTengahKeDasar, HargaJualTengah: input.HargaJualTengah,
	}

	if err := h.Repo.CreateProductGlobal(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan produk. Barcode/SKU terdeteksi duplikat di toko Anda."})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Produk berhasil ditambahkan! 📦", "data": product})
}

func (h *RetailHandler) GetProducts(c *gin.Context) {
	storeID := getStoreID(c)
	search := c.Query("search")
	category := c.Query("category")
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat katalog data produk"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Katalog produk berhasil dimuat!", "total_items": totalItems, "data": products})
}

func (h *RetailHandler) UpdateProduct(c *gin.Context) {
	storeID := getStoreID(c)
	role := c.MustGet("role").(string)

	if role != "owner" && role != "manager" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak! Menu perubahan produk hanya untuk Owner atau Manager."})
		return
	}

	// 🚀 NANGKEP PUBLIC ID STRING (ULID), BUKAN INTEGER BRAY!
	publicID := c.Param("id")
	db := h.Repo.GetDB()

	var product models.Product
	// 🚀 CARI BERDASARKAN public_id
	if err := db.Where("public_id = ? AND store_id = ?", publicID, storeID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan di database!"})
		return
	}

	// 🛡️ ANTI OVERWRITE
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

	var store models.Store
	if err := db.Select("public_id").First(&store, storeID).Error; err == nil {
		bucketName := os.Getenv("SUPABASE_BUCKET_NAME")
		if file, errFile := c.FormFile("gambar"); errFile == nil {
			if file.Size <= 5*1024*1024 {
				contentType := file.Header.Get("Content-Type")
				if fileSrc, errOpen := file.Open(); errOpen == nil {
					// 🚀 UBAH input.NamaProduk JADI product.NamaProduk BRAY!
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data produk. Kode Barcode/SKU bentrok."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Produk berhasil diubah! ✏️", "data": product})
}

func (h *RetailHandler) DeleteProduct(c *gin.Context) {
	storeID := getStoreID(c)
	role := c.MustGet("role").(string)

	if role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Waduh, tim kasir dilarang menghapus aset produk gudang bray!"})
		return
	}

	// 🚀 NANGKEP PUBLIC ID STRING
	publicID := c.Param("id")
	db := h.Repo.GetDB()

	var product models.Product
	// 🚀 CARI BERDASARKAN public_id
	if err := db.Where("public_id = ? AND store_id = ?", publicID, storeID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan!"})
		return
	}

	// Buat hapus
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menarik daftar kategori"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

func (h *RetailHandler) ExportProducts(c *gin.Context) {
	storeID := getStoreID(c)
	products, err := h.Repo.GetAllProductsForExport(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menarik data untuk ekspor"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Berkas dokumen tidak ditemukan"})
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '|'
	_, _ = reader.Read() // Skip Header

	records, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca struktur berkas CSV"})
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

			// Penambahan kolom baru saat import mapping bray
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

			tengah := ""
			besar := ""
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

			// 🛡️ CRITICAL GUARD: Cegah query liar overwrite jika SKU kosong di file CSV bray!
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
