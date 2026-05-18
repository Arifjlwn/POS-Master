package laundry

import (
	"net/http"
	"pos-backend/config"
	"pos-backend/models"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type KasirInput struct {
	Name     string `json:"name" binding:"required"`
	NoHP     string `json:"no_hp"`
	Password string `json:"password"`
}

// 🚀 1. AMBIL SEMUA DATA KASIR DI TOKO INI
func GetKasirList(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	
	var kasirList []models.User
	// Cari semua user yang role-nya kasir dan kerja di toko ini
	if err := config.DB.Where("store_id = ? AND role = ?", storeIDRaw, "kasir").Find(&kasirList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat data kasir"})
		return
	}

	// Format data biar enak dibaca Vue
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

// 🚀 2. TAMBAH KASIR BARU (Otomatis buatin Email Dummy)
func CreateKasir(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

	var input KasirInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data input tidak lengkap"})
		return
	}

	// Ambil nama toko buat bikin domain email
	var store models.Store
	config.DB.First(&store, storeID)

	// Format Email: kasir.nama@namatoko.com
	namaKasirBersih := strings.ToLower(strings.ReplaceAll(input.Name, " ", ""))
	namaTokoBersih := strings.ToLower(strings.ReplaceAll(store.NamaToko, " ", ""))
	emailDummy := "kasir." + namaKasirBersih + "@" + namaTokoBersih + ".com"

	// Tentukan Password (Kalau dikosongin, pakai 'kasir123')
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
		IsVerified: true, // 🚀 Wajib True biar lolos OTP
	}

	if err := config.DB.Create(&newKasir).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendaftarkan kasir. Pastikan nama tidak duplikat."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Kasir berhasil ditambahkan!",
		"email":   emailDummy,
	})
}

// 🚀 3. HAPUS KASIR (Pecat Karyawan)
func DeleteKasir(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	kasirID := c.Param("id") // Ambil ID dari URL

	// Hapus kasir yang ID-nya cocok dan pastikan dia beneran kasir di toko ini (biar aman dari hacker)
	if err := config.DB.Where("id = ? AND store_id = ? AND role = ?", kasirID, storeIDRaw, "kasir").Delete(&models.User{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data kasir"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Karyawan berhasil dihapus dari sistem!"})
}