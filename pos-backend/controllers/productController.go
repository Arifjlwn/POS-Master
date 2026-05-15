package controllers

import (
    "bytes"
    "encoding/csv"
    "fmt"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "time"

    "pos-backend/config"
    "pos-backend/models"

    "github.com/gin-gonic/gin"
)

// Struct untuk menangkap data dari Frontend
type ProductInput struct {
    SKU          *string `form:"sku"`
    NamaProduk   string  `form:"nama_produk" binding:"required"`
    Kategori     string  `form:"kategori"`
    HargaModal   float64 `form:"harga_modal"`
    HargaJual    float64 `form:"harga_jual" binding:"required"`
    Stok         int     `form:"stok"`
    // 🚀 TAMBAHAN FIELD SATUAN
    SatuanDasar  string  `form:"satuan_dasar"`
    SatuanBesar  string  `form:"satuan_besar"`
    IsiPerBesar  int     `form:"isi_per_besar"`
}

// Fungsi Tambah Produk + Upload Gambar
func CreateProduct(c *gin.Context) {
    // 1. Ambil ID Toko dari token JWT yang lolos satpam
    storeIDraw, exists := c.Get("store_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak, ID Toko tidak ditemukan!"})
        return
    }
    storeID := uint(storeIDraw.(float64))

    // 2. Tangkap JSON/FormData dari Frontend
    var input ProductInput
    if err := c.ShouldBind(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Periksa kembali isian form Anda: " + err.Error()})
        return
    }

    // Default Satuan Dasar jika kosong
    if input.SatuanDasar == "" {
        input.SatuanDasar = "PCS"
    }

    // 3. Logika Upload Gambar
    var imagePath string
    file, errFile := c.FormFile("gambar")

    if errFile == nil {
        // Buat folder otomatis kalo belom ada
        folderPath := "uploads/products"
        os.MkdirAll(folderPath, os.ModePerm)

        // Bikin nama file unik (Waktu + Nama Asli)
        fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
        savePath := filepath.Join(folderPath, fileName)

        // Simpan gambar ke folder
        if err := c.SaveUploadedFile(file, savePath); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan gambar produk"})
            return
        }

        // Simpan alamat URL nya di database
        imagePath = "/" + savePath
    }

    // 4. Rakit data produknya
    product := models.Product{
        StoreID:      storeID,
        SKU:          input.SKU,
        NamaProduk:   input.NamaProduk,
        Kategori:     input.Kategori,
        HargaModal:   input.HargaModal,
        HargaJual:    input.HargaJual,
        Stok:         input.Stok,
        Gambar:       imagePath,
        // 🚀 MASUKKAN DATA SATUAN KE MODEL
        SatuanDasar:  input.SatuanDasar,
        SatuanBesar:  input.SatuanBesar,
        IsiPerBesar:  input.IsiPerBesar,
    }

    // 5. Simpan ke database
    if err := config.DB.Create(&product).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan produk. Barcode mungkin sudah dipakai barang lain."})
        return
    }

    // 6. Beri balasan sukses
    c.JSON(http.StatusCreated, gin.H{
        "message": "Produk berhasil ditambahkan! 📦",
        "data":    product,
    })
}

// Fungsi Lihat Daftar Produk
func GetProducts(c *gin.Context) {
    // 1. Ambil ID Toko dari Satpam JWT
    storeIDRaw, exists := c.Get("store_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses Ditolak!"})
        return
    }
    storeID := uint(storeIDRaw.(float64))

    // 2. Tangkap parameter dari URL Vue (query string)
    search := c.Query("search")     // Untuk cari nama/sku
    category := c.Query("category") // Untuk filter kategori

    var products []models.Product
    
    // 3. Bangun Query Dasar (Wajib milik toko yang login)
    query := config.DB.Where("store_id = ?", storeID)

    // 🚀 LOGIKA PENCARIAN (Berdasarkan Nama atau SKU)
    if search != "" {
        // Kita pakai ILIKE biar "aqua" atau "AQUA" tetep ketemu (Case Insensitive)
        searchTerm := "%" + search + "%"
        query = query.Where("(nama_produk ILIKE ? OR sku ILIKE ?)", searchTerm, searchTerm)
    }

    // 🚀 LOGIKA FILTER KATEGORI (Exact Match)
    if category != "" {
        query = query.Where("kategori = ?", category)
    }

    // 4. Eksekusi ke database (Urutkan dari yang terbaru ditambah)
    if err := query.Order("id DESC").Find(&products).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data produk"})
        return
    }

    // 5. Kirim balik ke Vue
    c.JSON(http.StatusOK, gin.H{
        "message": "Katalog produk berhasil dimuat!",
        "total":   len(products),
        "data":    products, // Ini yang dibaca oleh products.value di Vue
    })
}

// Fungsi Ubah Produk (Update)
func UpdateProduct(c *gin.Context) {
    // 1. Cek ID Toko dari Satpam JWT
    storeID, _ := c.Get("store_id")
    role, _ := c.Get("role")

    // Logika RBAC
    if role != "owner" {
        c.JSON(http.StatusForbidden, gin.H{"error": "Hentikan! Cuma Owner yang boleh ubah harga/data barang."})
        return
    }
    
    // 2. Tangkap ID Produk dari ujung URL (Contoh: /api/products/1)
    productID := c.Param("id") 
    var product models.Product
    
    // 3. Cari produknya. Syarat Wajib: ID Produk harus cocok DAN ID Toko harus cocok!
    if err := config.DB.Where("id = ? AND store_id = ?", productID, storeID).First(&product).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan atau bukan milik toko Anda!"})
        return
    }

    // 4. Tangkap data baru dari Frontend (Karena dikirim sebagai FormData, kita bisa ambil manual atau pakai ShouldBind)
    // Supaya lebih aman menangani update parsial di FormData, kita tangkap manual nilai yang masuk
    
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

    // 🚀 UPDATE DATA SATUAN
    if satuanDasar := c.PostForm("satuan_dasar"); satuanDasar != "" {
        product.SatuanDasar = satuanDasar
    }
    
    product.SatuanBesar = c.PostForm("satuan_besar") // Bisa kosong
    
    if isiPerBesar, err := strconv.Atoi(c.PostForm("isi_per_besar")); err == nil {
        product.IsiPerBesar = isiPerBesar
    } else {
        product.IsiPerBesar = 0 // Reset ke 0 jika kosong
    }

    // 5. Cek apakah ada file foto baru yang diupload
    file, errFile := c.FormFile("gambar")
    if errFile == nil {
        folderPath := "uploads/products"
        os.MkdirAll(folderPath, os.ModePerm)
        fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
        savePath := filepath.Join(folderPath, fileName)

        if err := c.SaveUploadedFile(file, savePath); err == nil {
            // Hapus gambar lama jika ada
            if product.Gambar != "" {
                os.Remove("." + product.Gambar)
            }
            product.Gambar = "/" + savePath
        }
    }

    // 6. Simpan kembali ke database
    if err := config.DB.Save(&product).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update produk. Barcode mungkin bentrok."})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Produk berhasil diubah! ✏️", 
        "data": product,
    })
}

// Fungsi Hapus Produk (Delete)
func DeleteProduct(c *gin.Context) {
    storeID, _ := c.Get("store_id")
    role, _ := c.Get("role")

    // Logika RBAC
    if role != "owner" {
        c.JSON(http.StatusForbidden, gin.H{"error": "Waduh, Kasir dilarang hapus barang dari sistem!"})
        return
    }

    productID := c.Param("id")
    var product models.Product
    
    // Pastikan produk yang mau dihapus itu beneran ada dan milik dia
    if err := config.DB.Where("id = ? AND store_id = ?", productID, storeID).First(&product).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan atau bukan milik toko Anda!"})
        return
    }

    // Hapus gambar fisiknya dari harddisk kalau ada
    if product.Gambar != "" {
        os.Remove("." + product.Gambar)
    }

    // Hapus dari muka bumi (database)
    config.DB.Delete(&product)
    
    c.JSON(http.StatusOK, gin.H{"message": "Barang berhasil dihapus dari gudang! 🗑️"})
}

// Fungsi Ambil Daftar Kategori Unik
func GetCategories(c *gin.Context) {
    // Ambil ID Toko
    storeIDRaw, exists := c.Get("store_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses Ditolak!"})
        return
    }
    storeID := uint(storeIDRaw.(float64))

    var categories []string

    // Minta GORM ambil kolom "kategori" yang unik (tidak dobel) dan tidak kosong
    config.DB.Model(&models.Product{}).
        Where("store_id = ? AND kategori IS NOT NULL AND kategori != ''", storeID).
        Distinct("kategori").
        Pluck("kategori", &categories)

    // Kirim Array kategori ke Vue
    c.JSON(http.StatusOK, gin.H{
        "data": categories,
    })
}

// Fungsi Ekspor CSV
func ExportProducts(c *gin.Context) {
    // Ambil ID Toko
    storeIDRaw, exists := c.Get("store_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses Ditolak !"})
        return
    }
    storeID := uint(storeIDRaw.(float64))

    var products []models.Product
    if err := config.DB.Where("store_id = ?", storeID).Order("id DESC").Find(&products).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data produk"})
        return
    }

    // Siapkan CSV
    b := &bytes.Buffer{}
    w := csv.NewWriter(b)

    // Tulis Header Kolom
    w.Write([]string{"SKU", "Nama Produk", "Kategori", "Harga Modal", "Harga Jual", "Stok", "Satuan Dasar", "Satuan Besar", "Isi Per Besar"})

    // Looping isi data produk
    for _, p := range products {
        sku := ""
        if p.SKU != nil {
            sku = *p.SKU
        }

        w.Write([]string{
            sku,
            p.NamaProduk,
            p.Kategori,
            fmt.Sprintf("%.0f", p.HargaModal),
            fmt.Sprintf("%.0f", p.HargaJual),
            fmt.Sprintf("%d", p.Stok),
            p.SatuanDasar,
            p.SatuanBesar,
            fmt.Sprintf("%d", p.IsiPerBesar),
        })
    }
    w.Flush()

    // Paksa Browser untuk download file
    c.Header("Content-Description", "File Transfer")
    c.Header("Content-Disposition", "attachment; filename = katalog_produk_pos.csv")
    c.Data(http.StatusOK, "text/csv", b.Bytes())
}