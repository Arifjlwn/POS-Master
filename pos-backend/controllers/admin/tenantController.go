package admin

import (
	"net/http"
	"time"

	"pos-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TenantController struct {
	DB *gorm.DB
}

// 🏪 FIX ABADI: Menggunakan teknik flat scan join langsung ke users.id = stores.owner_id
func (t *TenantController) GetAllTenants(c *gin.Context) {
	// Kita bikin flat struct, semua field mentah dari database ditampung di sini
	// Biar GORM ga pusing nge-mapping nested object struct models
	type FlatStoreResponse struct {
		ID                 uint       `json:"id"`
		PublicID           string     `json:"public_id"`
		NamaToko           string     `json:"nama_toko"`
		Telepon            string     `json:"telepon"`
		BusinessType       string     `json:"business_type"`
		SubscriptionPlan   string     `json:"subscription_plan"`
		SubscriptionStatus string     `json:"subscription_status"`
		SubscriptionEnd    *time.Time `json:"subscription_end"`
		OwnerName          string     `json:"owner_name"`
		OwnerEmail         string     `json:"owner_email"`
	}

	var results []FlatStoreResponse

	// 🚀 EKSEKUSI LOGIKA LU BRAY: Cek stores, terus join ke tabel users berdasarkan owner_id
	// Kita gunakan COALESCE biar kalau datanya null di database, frontend dapet string aman, bukan null pointer
	err := t.DB.Table("stores").
		Select(`
			stores.id, 
			stores.public_id, 
			stores.nama_toko, 
			stores.telepon, 
			stores.business_type, 
			stores.subscription_plan, 
			stores.subscription_status, 
			stores.subscription_end, 
			COALESCE(users.name, 'Owner Tidak Terdeteksi') as owner_name, 
			COALESCE(users.email, '-') as owner_email
		`).
		Joins("LEFT JOIN users ON users.id = stores.owner_id"). // Membidik users.id secara presisi
		Order("stores.created_at desc").
		Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Sistem gagal mengintegrasikan data relasi owner!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "sukses",
		"data":   results,
	})
}

// 🔒 2. SUSPEND STORE (Membekukan Aktivitas Dagang Toko)
func (t *TenantController) SuspendTenant(c *gin.Context) {
	storeID := c.Param("id")

	// 🚀 SINKRON: Mengubah SubscriptionStatus menjadi 'suspended' sesuai field asli lu bray
	err := t.DB.Model(&models.Store{}).Where("id = ?", storeID).Update("subscription_status", "suspended").Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membekukan status langganan toko!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "sukses",
		"message": "Toko berhasil dibekukan! Akses fitur dagang otomatis terkunci.",
	})
}

// 🟢 3. ACTIVATE STORE (Memulihkan & Memperpanjang Masa Aktif Toko)
func (t *TenantController) ActivateTenant(c *gin.Context) {
	storeID := c.Param("id")
	newExpiry := time.Now().AddDate(0, 1, 0) // Tambah bonus masa aktif 30 hari kedepan

	// 🚀 SINKRON: Update status jadi 'active' dan isi pointer SubscriptionEnd
	err := t.DB.Model(&models.Store{}).Where("id = ?", storeID).Updates(map[string]interface{}{
		"subscription_status": "active",
		"subscription_end":    &newExpiry,
	}).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui masa aktif toko!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "sukses",
		"message": "Toko resmi diaktifkan kembali! Masa aktif diperpanjang hingga sebulan ke depan.",
	})
}
