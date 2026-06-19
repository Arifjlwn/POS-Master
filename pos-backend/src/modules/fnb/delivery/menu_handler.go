package delivery

import (
	"net/http"
	"strconv"

	"pos-backend/src/modules/fnb/domain"
	"pos-backend/src/modules/fnb/repository"

	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	Repo repository.MenuRepository
}

// 1. TAMBAH PRODUK BARU (HANYA OWNER)
func (h *MenuHandler) CreateProduct(c *gin.Context) {
	role, _ := c.Get("role")
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

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

	product := domain.Menu{
		StoreID:     storeID,
		NamaProduk:  input.Nama,
		HargaJual:   input.Harga,
		Kategori:    input.Kategori,
		Stok:        input.Stok,
		Gambar:      input.Gambar,
		IsAvailable: true,
	}

	if err := h.Repo.Create(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan produk"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Menu baru berhasil dirilis!", "data": product})
}

// 2. MATI-NYALAKAN MENU (OWNER, KASIR, KITCHEN BISA AKSES)
func (h *MenuHandler) ToggleAvailability(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, _ := strconv.Atoi(productIDStr)
	
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

	// Cari produk berdasarkan ID dan pastiin milik toko yang sama lewat Repo
	product, err := h.Repo.GetByID(uint(productID), storeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu tidak ditemukan atau bukan milik toko Anda!"})
		return
	}

	// Balik statusnya
	product.IsAvailable = !product.IsAvailable

	if err := h.Repo.Update(product); err != nil {
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
func (h *MenuHandler) GetProducts(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))
	role, _ := c.Get("role")

	// JIKA AKSES DARI QR MEJA / KASIR BIASA, HANYA TAMPILKAN YANG AVAILABLE
	onlyAvailable := false
	if role != "owner" && role != "kitchen" && role != "kasir" {
		onlyAvailable = true
	}

	products, err := h.Repo.GetAll(storeID, onlyAvailable)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data produk"})
		return
	}
	
	c.JSON(http.StatusOK, products)
}

// 4. UPDATE PRODUCT
func (h *MenuHandler) UpdateProduct(c *gin.Context) {
	role, _ := c.Get("role")

	if role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya owner yang boleh edit menu"})
		return
	}

	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))
	
	productIDStr := c.Param("id")
	productID, _ := strconv.Atoi(productIDStr)

	// Cari produk sesuai toko lewat Repo
	product, err := h.Repo.GetByID(uint(productID), storeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	var input struct {
		Nama     string  `json:"nama"`
		Harga    float64 `json:"harga"`
		Kategori string  `json:"kategori"`
		Gambar   string  `json:"gambar"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product.NamaProduk = input.Nama
	product.HargaJual = input.Harga
	product.Kategori = input.Kategori
	product.Gambar = input.Gambar

	if err := h.Repo.Update(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update produk"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Menu berhasil diperbarui",
		"data":    product,
	})
}

// 5. DELETE PRODUCT
func (h *MenuHandler) DeleteProduct(c *gin.Context) {
	role, _ := c.Get("role")

	if role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya owner yang boleh hapus menu"})
		return
	}

	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))
	
	productIDStr := c.Param("id")
	productID, _ := strconv.Atoi(productIDStr)

	// Cek apakah produk ada lewat Repo
	_, err := h.Repo.GetByID(uint(productID), storeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	// Hapus lewat Repo
	if err := h.Repo.Delete(uint(productID), storeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hapus produk"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu berhasil dihapus"})
}