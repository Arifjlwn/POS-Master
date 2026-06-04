package delivery

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"pos-backend/models"
	"pos-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ===========================
// 📦 PRODUCT HANDLERS
// ===========================

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
}

func (h *RetailHandler) CreateProduct(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	var storeID uint
	switch v := storeIDRaw.(type) { case float64: storeID = uint(v); case uint: storeID = v; case int: storeID = uint(v) }

	var input ProductInput
	if err := c.ShouldBind(&input); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Form tidak valid: " + err.Error()}); return }
	if input.SatuanDasar == "" { input.SatuanDasar = "PCS" }

	var imageURL string
	bucketName := os.Getenv("SUPABASE_BUCKET_NAME")

	if file, errFile := c.FormFile("gambar"); errFile == nil {
		if file.Size > 5*1024*1024 { c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran gambar maksimal 5 MB"}); return }
		contentType := file.Header.Get("Content-Type")
		fileSrc, err := file.Open(); if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membaca file gambar"}); return }
		defer fileSrc.Close()

		remotePath := fmt.Sprintf("stores/%d/products/%d_img", storeID, time.Now().Unix())
		urlResult, errUpload := utils.UploadToSupabase(fileSrc, file.Filename, contentType, bucketName, remotePath)
		if errUpload != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal upload gambar ke cloud storage"}); return }
		imageURL = urlResult
	}

	product := models.Product{
		StoreID: storeID, SKU: input.SKU, NamaProduk: input.NamaProduk, Kategori: input.Kategori,
		HargaModal: input.HargaModal, HargaJual: input.HargaJual, Stok: input.Stok, Gambar: imageURL,
		SatuanDasar: input.SatuanDasar, SatuanBesar: input.SatuanBesar, IsiPerBesar: input.IsiPerBesar,
		HargaJualBesar: input.HargaJualBesar, IsNestedUom: input.IsNestedUom, SatuanTengah: input.SatuanTengah,
		IsiBesarKeTengah: input.IsiBesarKeTengah, IsiTengahKeDasar: input.IsiTengahKeDasar,
	}

	if err := h.Repo.CreateProductGlobal(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan produk. Barcode mungkin duplikat."})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Produk berhasil ditambahkan! 📦", "data": product})
}

func (h *RetailHandler) GetProducts(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	var storeID uint
	switch v := storeIDRaw.(type) { case float64: storeID = uint(v); case uint: storeID = v; case int: storeID = uint(v) }

	search := c.Query("search"); category := c.Query("category")
	pageStr := c.Query("page"); limitStr := c.Query("limit")
	usePagination := false; limit, offset := 10, 0

	if pageStr != "" || limitStr != "" {
		usePagination = true
		if pageStr == "" { pageStr = "1" }; if limitStr == "" { limitStr = "10" }
		p, _ := strconv.Atoi(pageStr); l, _ := strconv.Atoi(limitStr)
		limit = l; offset = (p - 1) * limit
	}

	products, totalItems, err := h.Repo.GetProductsCatalog(storeID, search, category, limit, offset, usePagination)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data produk"}); return }

	c.JSON(http.StatusOK, gin.H{"message": "Katalog produk berhasil dimuat!", "total_items": totalItems, "data": products})
}

func (h *RetailHandler) UpdateProduct(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id"); role := c.MustGet("role").(string)
	var storeID uint
	switch v := storeIDRaw.(type) { case float64: storeID = uint(v); case uint: storeID = v; case int: storeID = uint(v) }

	if role != "owner" { c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak! Fitur ini khusus Owner."}); return }

	productID, _ := strconv.Atoi(c.Param("id"))
	product, err := h.Repo.GetProductByIDSimple(uint(productID), storeID)
	if err != nil { c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan!"}); return }

	product.NamaProduk = c.PostForm("nama_produk")
	if sku := c.PostForm("sku"); sku != "" { product.SKU = &sku } else { product.SKU = nil }
	product.Kategori = c.PostForm("kategori")

	if hModal, err := strconv.ParseFloat(c.PostForm("harga_modal"), 64); err == nil { product.HargaModal = hModal }
	if hJual, err := strconv.ParseFloat(c.PostForm("harga_jual"), 64); err == nil { product.HargaJual = hJual }
	if stok, err := strconv.Atoi(c.PostForm("stok")); err == nil { product.Stok = stok }
	if sDasar := c.PostForm("satuan_dasar"); sDasar != "" { product.SatuanDasar = sDasar }
	product.SatuanBesar = c.PostForm("satuan_besar")
	if iPerBesar, err := strconv.Atoi(c.PostForm("isi_per_besar")); err == nil { product.IsiPerBesar = iPerBesar } else { product.IsiPerBesar = 0 }
	if hJualBesar, err := strconv.ParseFloat(c.PostForm("harga_jual_besar"), 64); err == nil { product.HargaJualBesar = hJualBesar } else { product.HargaJualBesar = 0 }

	product.IsNestedUom = c.PostForm("is_nested_uom") == "true"
	product.SatuanTengah = c.PostForm("satuan_tengah")
	if ibt, err := strconv.Atoi(c.PostForm("isi_besar_ke_tengah")); err == nil { product.IsiBesarKeTengah = ibt } else { product.IsiBesarKeTengah = 0 }
	if itd, err := strconv.Atoi(c.PostForm("isi_tengah_ke_dasar")); err == nil { product.IsiTengahKeDasar = itd } else { product.IsiTengahKeDasar = 0 }

	bucketName := os.Getenv("SUPABASE_BUCKET_NAME")
	if file, errFile := c.FormFile("gambar"); errFile == nil {
		if file.Size > 5*1024*1024 { c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran gambar maksimal 5 MB"}); return }
		contentType := file.Header.Get("Content-Type")
		fileSrc, err := file.Open(); if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membaca file gambar"}); return }
		defer fileSrc.Close()

		remotePath := fmt.Sprintf("stores/%d/products/%d_img", storeID, time.Now().Unix())
		urlResult, errUpload := utils.UploadToSupabase(fileSrc, file.Filename, contentType, bucketName, remotePath)
		if errUpload == nil { product.Gambar = urlResult }
	}

	db := h.Repo.GetDB()
	if err := h.Repo.SaveProduct(db, product); err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update produk. Barcode mungkin bentrok."}); return }
	c.JSON(http.StatusOK, gin.H{"message": "Produk berhasil diubah! ✏️", "data": product})
}

func (h *RetailHandler) DeleteProduct(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id"); role := c.MustGet("role").(string)
	var storeID uint
	switch v := storeIDRaw.(type) { case float64: storeID = uint(v); case uint: storeID = v; case int: storeID = uint(v) }

	if role != "owner" { c.JSON(http.StatusForbidden, gin.H{"error": "Waduh, Kasir dilarang hapus barang!"}); return }

	productID, _ := strconv.Atoi(c.Param("id"))
	product, err := h.Repo.GetProductByIDSimple(uint(productID), storeID)
	if err != nil { c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan!"}); return }

	if err := h.Repo.DeleteProductGlobal(product); err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus produk"}); return }
	c.JSON(http.StatusOK, gin.H{"message": "Barang berhasil dihapus dari gudang! 🗑️"})
}

func (h *RetailHandler) GetCategories(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	var storeID uint
	switch v := storeIDRaw.(type) { case float64: storeID = uint(v); case uint: storeID = v; case int: storeID = uint(v) }

	categories, err := h.Repo.GetDistinctCategories(storeID)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil kategori"}); return }
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

func (h *RetailHandler) ExportProducts(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	var storeID uint
	switch v := storeIDRaw.(type) { case float64: storeID = uint(v); case uint: storeID = v; case int: storeID = uint(v) }

	products, err := h.Repo.GetAllProductsForExport(storeID)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data produk"}); return }

	b := &bytes.Buffer{}; w := csv.NewWriter(b); w.Comma = '|'
	w.Write([]string{"SKU", "Nama Produk", "Kategori", "Harga Modal", "Harga Jual", "Stok", "Satuan Terkecil", "Satuan Tengah", "Satuan Besar", "Isi Per Besar", "Harga Jual Besar"})

	for _, p := range products {
		sku := ""
		if p.SKU != nil { sku = *p.SKU }
		w.Write([]string{
			sku, p.NamaProduk, p.Kategori, fmt.Sprintf("%.0f", p.HargaModal), fmt.Sprintf("%.0f", p.HargaJual),
			fmt.Sprintf("%d", p.Stok), p.SatuanDasar, p.SatuanTengah, p.SatuanBesar, fmt.Sprintf("%d", p.IsiPerBesar), fmt.Sprintf("%.0f", p.HargaJualBesar),
		})
	}
	w.Flush()

	c.Header("Content-Disposition", "attachment; filename=katalog_produk_pos.csv")
	c.Data(http.StatusOK, "text/csv", b.Bytes())
}

func (h *RetailHandler) ImportProducts(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	var storeID uint
	switch v := storeIDRaw.(type) { case float64: storeID = uint(v); case uint: storeID = v; case int: storeID = uint(v) }

	file, _, err := c.Request.FormFile("file"); if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "File tidak ditemukan"}); return }
	defer file.Close()

	reader := csv.NewReader(file); reader.Comma = '|'; _, _ = reader.Read()
	records, err := reader.ReadAll(); if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca isi CSV"}); return }

	db := h.Repo.GetDB()
	err = db.Transaction(func(tx *gorm.DB) error {
		for _, row := range records {
			if len(row) < 7 { continue }
			sku := row[0]; nama := row[1]; kategori := row[2]
			modal, _ := strconv.ParseFloat(row[3], 64); jual, _ := strconv.ParseFloat(row[4], 64); stok, _ := strconv.Atoi(row[5])
			dasar := row[6]

			tengah := ""; besar := ""; isi := 0; var jualBesar float64
			if len(row) > 7 { tengah = row[7] }
			if len(row) > 8 { besar = row[8] }
			if len(row) > 9 { isi, _ = strconv.Atoi(row[9]) }
			if len(row) > 10 { jualBesar, _ = strconv.ParseFloat(row[10], 64) }

			if nama == "" { continue }

			var product models.Product
			res := tx.Where("sku = ? AND store_id = ?", sku, storeID).First(&product)

			if res.Error == nil {
				product.NamaProduk = nama; product.Kategori = kategori; product.HargaModal = modal; product.HargaJual = jual; product.Stok = stok
				product.SatuanDasar = dasar; product.SatuanTengah = tengah; product.SatuanBesar = besar; product.IsiPerBesar = isi; product.HargaJualBesar = jualBesar
				tx.Save(&product)
			} else {
				newProduct := models.Product{
					StoreID: storeID, SKU: &sku, NamaProduk: nama, Kategori: kategori, HargaModal: modal, HargaJual: jual, Stok: stok,
					SatuanDasar: dasar, SatuanTengah: tengah, SatuanBesar: besar, IsiPerBesar: isi, HargaJualBesar: jualBesar,
				}
				tx.Create(&newProduct)
			}
		}
		return nil
	})

	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal impor: " + err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"message": "Berhasil mengimpor " + strconv.Itoa(len(records)) + " produk"})
}