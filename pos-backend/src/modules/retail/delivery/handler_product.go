package delivery

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"pos-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ===========================
// 📦 PRODUCT HANDLERS
// ===========================
type ProductInput struct {
	SKU            *string `form:"sku"`
	NamaProduk     string  `form:"nama_produk" binding:"required"`
	Kategori       string  `form:"kategori"`
	HargaModal     float64 `form:"harga_modal"`
	HargaJual      float64 `form:"harga_jual" binding:"required"`
	Stok           int     `form:"stok"`
	SatuanDasar    string  `form:"satuan_dasar"`
	SatuanBesar    string  `form:"satuan_besar"`
	IsiPerBesar    int     `form:"isi_per_besar"`
	HargaJualBesar float64 `form:"harga_jual_besar"`

	// 🚀 TAMBAHAN: Nangkep Form 3 Lapis (Rokok/Renteng)
	IsNestedUom      bool   `form:"is_nested_uom"`
	SatuanTengah     string `form:"satuan_tengah"`
	IsiBesarKeTengah int    `form:"isi_besar_ke_tengah"`
	IsiTengahKeDasar int    `form:"isi_tengah_ke_dasar"`
}

func (h *RetailHandler) CreateProduct(c *gin.Context) {
	storeID := uint(c.MustGet("store_id").(float64))

	var input ProductInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Periksa kembali isian form Anda: " + err.Error()})
		return
	}
	if input.SatuanDasar == "" {
		input.SatuanDasar = "PCS"
	}

	var imagePath string
	file, errFile := c.FormFile("gambar")
	if errFile == nil {
		folderPath := "uploads/products"
		os.MkdirAll(folderPath, os.ModePerm)
		fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
		savePath := filepath.Join(folderPath, fileName)

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan gambar produk"})
			return
		}
		imagePath = "/" + savePath
	}

	product := models.Product{
		StoreID:        storeID,
		SKU:            input.SKU,
		NamaProduk:     input.NamaProduk,
		Kategori:       input.Kategori,
		HargaModal:     input.HargaModal,
		HargaJual:      input.HargaJual,
		Stok:           input.Stok,
		Gambar:         imagePath,
		SatuanDasar:    input.SatuanDasar,
		SatuanBesar:    input.SatuanBesar,
		IsiPerBesar:    input.IsiPerBesar,
		HargaJualBesar: input.HargaJualBesar,

		// 🚀 TAMBAHAN: Masukin data 3 Lapis ke Database
		IsNestedUom:      input.IsNestedUom,
		SatuanTengah:     input.SatuanTengah,
		IsiBesarKeTengah: input.IsiBesarKeTengah,
		IsiTengahKeDasar: input.IsiTengahKeDasar,
	}

	if err := h.Repo.CreateProductGlobal(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan produk. Barcode mungkin sudah dipakai barang lain."})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Produk berhasil ditambahkan! 📦", "data": product})
}

func (h *RetailHandler) GetProducts(c *gin.Context) {
	storeID := uint(c.MustGet("store_id").(float64))
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data produk"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Katalog produk berhasil dimuat!",
		"total_items": totalItems,
		"data":        products,
	})
}

func (h *RetailHandler) UpdateProduct(c *gin.Context) {
	storeID := uint(c.MustGet("store_id").(float64))
	role := c.MustGet("role").(string)

	if role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hentikan! Cuma Owner yang boleh ubah harga/data barang."})
		return
	}

	productID, _ := strconv.Atoi(c.Param("id"))
	product, err := h.Repo.GetProductByIDSimple(uint(productID), storeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan atau bukan milik toko Anda!"})
		return
	}

	product.NamaProduk = c.PostForm("nama_produk")
	if sku := c.PostForm("sku"); sku != "" {
		product.SKU = &sku
	} else {
		product.SKU = nil
	}
	product.Kategori = c.PostForm("kategori")

	if hargaModal, err := strconv.ParseFloat(c.PostForm("harga_modal"), 64); err == nil {
		product.HargaModal = hargaModal
	}
	if hargaJual, err := strconv.ParseFloat(c.PostForm("harga_jual"), 64); err == nil {
		product.HargaJual = hargaJual
	}
	if stok, err := strconv.Atoi(c.PostForm("stok")); err == nil {
		product.Stok = stok
	}
	if satuanDasar := c.PostForm("satuan_dasar"); satuanDasar != "" {
		product.SatuanDasar = satuanDasar
	}
	product.SatuanBesar = c.PostForm("satuan_besar")
	if isiPerBesar, err := strconv.Atoi(c.PostForm("isi_per_besar")); err == nil {
		product.IsiPerBesar = isiPerBesar
	} else {
		product.IsiPerBesar = 0
	}

	// Nangkep harga_jual_besar
	if hargaJualBesar, err := strconv.ParseFloat(c.PostForm("harga_jual_besar"), 64); err == nil {
		product.HargaJualBesar = hargaJualBesar
	} else {
		product.HargaJualBesar = 0
	}

	// 🚀 TAMBAHAN: Nangkep data kemasan tengah pas di-Edit
	product.IsNestedUom = c.PostForm("is_nested_uom") == "true"
	product.SatuanTengah = c.PostForm("satuan_tengah")
	if ibt, err := strconv.Atoi(c.PostForm("isi_besar_ke_tengah")); err == nil {
		product.IsiBesarKeTengah = ibt
	} else {
		product.IsiBesarKeTengah = 0
	}
	if itd, err := strconv.Atoi(c.PostForm("isi_tengah_ke_dasar")); err == nil {
		product.IsiTengahKeDasar = itd
	} else {
		product.IsiTengahKeDasar = 0
	}

	file, errFile := c.FormFile("gambar")
	if errFile == nil {
		folderPath := "uploads/products"
		os.MkdirAll(folderPath, os.ModePerm)
		fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
		savePath := filepath.Join(folderPath, fileName)

		if err := c.SaveUploadedFile(file, savePath); err == nil {
			if product.Gambar != "" {
				os.Remove("." + product.Gambar)
			}
			product.Gambar = "/" + savePath
		}
	}

	db := h.Repo.GetDB()
	if err := h.Repo.SaveProduct(db, product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update produk. Barcode mungkin bentrok."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Produk berhasil diubah! ✏️", "data": product})
}

func (h *RetailHandler) DeleteProduct(c *gin.Context) {
	storeID := uint(c.MustGet("store_id").(float64))
	role := c.MustGet("role").(string)

	if role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Waduh, Kasir dilarang hapus barang dari sistem!"})
		return
	}

	productID, _ := strconv.Atoi(c.Param("id"))
	product, err := h.Repo.GetProductByIDSimple(uint(productID), storeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan atau bukan milik toko Anda!"})
		return
	}

	if product.Gambar != "" {
		os.Remove("." + product.Gambar)
	}
	if err := h.Repo.DeleteProductGlobal(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus produk"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Barang berhasil dihapus dari gudang! 🗑️"})
}

func (h *RetailHandler) GetCategories(c *gin.Context) {
	storeID := uint(c.MustGet("store_id").(float64))
	categories, err := h.Repo.GetDistinctCategories(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil kategori"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

func (h *RetailHandler) ExportProducts(c *gin.Context) {
	storeID := uint(c.MustGet("store_id").(float64))
	products, err := h.Repo.GetAllProductsForExport(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data produk"})
		return
	}

	b := &bytes.Buffer{}
	w := csv.NewWriter(b)
	w.Comma = '|'

	// 🚀 1. TAMBAHIN HEADER "Harga Jual Besar"
	w.Write([]string{"SKU", "Nama Produk", "Kategori", "Harga Modal", "Harga Jual", "Stok", "Satuan Terkecil", "Satuan Tengah", "Satuan Besar", "Isi Per Besar", "Harga Jual Besar"})

	for _, p := range products {
		sku := ""
		if p.SKU != nil {
			sku = *p.SKU
		}
		w.Write([]string{
			sku, p.NamaProduk, p.Kategori,
			fmt.Sprintf("%.0f", p.HargaModal), fmt.Sprintf("%.0f", p.HargaJual),
			fmt.Sprintf("%d", p.Stok), p.SatuanDasar, p.SatuanTengah, p.SatuanBesar, fmt.Sprintf("%d", p.IsiPerBesar),
			fmt.Sprintf("%.0f", p.HargaJualBesar), // 🚀 2. MASUKIN DATANYA KE BARIS CSV
		})
	}
	w.Flush()

	c.Header("Content-Disposition", "attachment; filename = katalog_produk_pos.csv")
	c.Data(http.StatusOK, "text/csv", b.Bytes())
}

func (h *RetailHandler) ImportProducts(c *gin.Context) {
	storeID := uint(c.MustGet("store_id").(float64))
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File tidak ditemukan"})
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '|'
	_, _ = reader.Read() // Skip header

	records, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca isi CSV"})
		return
	}

	db := h.Repo.GetDB()
	err = db.Transaction(func(tx *gorm.DB) error {
		for _, row := range records {
			sku := row[0]
			nama := row[1]
			kategori := row[2]
			modal, _ := strconv.ParseFloat(row[3], 64)
			jual, _ := strconv.ParseFloat(row[4], 64)
			stok, _ := strconv.Atoi(row[5])
			dasar := row[6]
			besar := row[7]
			isi, _ := strconv.Atoi(row[8])

			// 🚀 3. BACA KOLOM BARU (Pake pengaman if len(row) > 9 biar ga error kalo upload CSV lama)
			var jualBesar float64
			if len(row) > 9 {
				jualBesar, _ = strconv.ParseFloat(row[9], 64)
			}

			if nama == "" {
				continue
			}

			var product models.Product
			res := tx.Where("sku = ? AND store_id = ?", sku, storeID).First(&product)

			if res.Error == nil {
				// UPDATE BARANG YANG UDAH ADA
				product.NamaProduk = nama
				product.Kategori = kategori
				product.HargaModal = modal
				product.HargaJual = jual
				product.Stok = stok
				product.SatuanDasar = dasar
				product.SatuanBesar = besar
				product.IsiPerBesar = isi
				product.HargaJualBesar = jualBesar // 🚀 4. SIMPAN HARGA GROSIR

				tx.Save(&product)
			} else {
				// BIKIN BARANG BARU
				newProduct := models.Product{
					StoreID: storeID, SKU: &sku, NamaProduk: nama, Kategori: kategori,
					HargaModal: modal, HargaJual: jual, Stok: stok, SatuanDasar: dasar,
					SatuanBesar: besar, IsiPerBesar: isi,
					HargaJualBesar: jualBesar, // 🚀 5. SIMPAN HARGA GROSIR
				}
				tx.Create(&newProduct)
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal impor: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Berhasil mengimpor " + strconv.Itoa(len(records)) + " produk"})
}