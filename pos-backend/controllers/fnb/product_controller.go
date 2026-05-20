package fnb

import (
	"net/http"
	"pos-backend/config"
	"pos-backend/models"
	"github.com/gin-gonic/gin"
)

// 1. TAMBAH PRODUK BARU (HANYA OWNER)
func CreateProduct(c *gin.Context) {
	// Ambil data dari Auth Middleware kamu
	role, _ := c.Get("role")
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

	// Proteksi tingkat tinggi: Kalau bukan owner, tendang!
	if role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak! Hanya Owner yang boleh menambah menu."})
		return
	}

	var input struct {
		Nama     string  `json:"nama" binding:"required"`
		Harga    float64 `json:"harga" binding:"required"`
		Kategori string  `json:"kategori" binding:"required"`
		Stok     int     `json:"stok"`
		Gambar   string  `json:"gambar"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := models.Product{
		StoreID:     storeID,
		NamaProduk:  input.Nama,
		HargaJual:   input.Harga,
		Kategori:    input.Kategori,
		Stok:        input.Stok,
		Gambar:      input.Gambar,
		IsAvailable: true, // Default aktif pas dibuat
	}

	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan produk"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Menu baru berhasil dirilis!", "data": product})
}

// 2. MATI-NYALAKAN MENU (OWNER, KASIR, KITCHEN BISA AKSES)
func ToggleAvailability(c *gin.Context) {
	productID := c.Param("id")
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

	var product models.Product
	// Cari produk berdasarkan ID dan pastiin milik toko yang sama
	if err := config.DB.Where("id = ? AND store_id = ?", productID, storeID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu tidak ditemukan atau bukan milik toko Anda!"})
		return
	}

	// 🚀 Balik statusnya (kalau true jadi false, kalau false jadi true)
	product.IsAvailable = !product.IsAvailable

	if err := config.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengubah status menu"})
		return
	}

	statusTeks := "Dihidupkan (Pelanggan bisa order)"
	if !product.IsAvailable {
		statusTeks = "Dimatikan (Menu disembunyikan dari QR Meja)"
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Status menu berhasil diperbarui!",
		"nama_menu":    product.NamaProduk,
		"is_available": product.IsAvailable,
		"keterangan":   statusTeks,
	})
}

// 3. AMBIL DATA PRODUK (UNTUK KASIR & QR MEJA SELF-SERVICE)
func GetProducts(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64)) // 👈 Convert dengan aman
	role, _ := c.Get("role")

	var products []models.Product
	query := config.DB.Where("store_id = ?", storeID)

	// 🚀 JIKA AKSES DARI QR MEJA / KASIR BIASA, HANYA TAMPILKAN YANG AVAILABLE
	if role != "owner" && role != "kitchen" && role != "kasir" {
		query = query.Where("is_available = ?", true)
	}

	query.Find(&products)
	c.JSON(http.StatusOK, products)
}