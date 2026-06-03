package auth

import (
	"net/http"
	"strings"
	"time"

	"pos-backend/models"
	src "pos-backend/src/core/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// SetupTokoBaru menangani proses pembuatan cabang/infrastruktur baru (Khusus Owner)
func SetupTokoBaru(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak valid"})
		return
	}
	userID := uint(userIDRaw.(float64))

	// 1. TANGKEP DATA DARI VUE TERMASUK INDUSTRY & PLAN SAAS
	var input struct {
		NamaToko      string `json:"nama_toko" binding:"required"`
		Telepon       string `json:"telepon"`
		Business_type string `json:"business_type" binding:"required"`
		Industry      string `json:"industry"` // Titipan Landing Page: retail/fnb/jasa
		Plan          string `json:"plan"`     // Titipan Landing Page: trial/basic/pro/premium
		AlamatJalan   string `json:"alamat_toko"`
		Provinsi      string `json:"provinsi"`
		Kota          string `json:"kota"`
		Kecamatan     string `json:"kecamatan"`
		Kelurahan     string `json:"kelurahan"`
		KodePos       string `json:"kode_pos"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data salah, pastikan form terisi lengkap"})
		return
	}

	var user models.User
	if err := src.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	// 🚀 LOGIKA ENTERPRISE MULTI-OUTLET:
	// Hanya ROLE OWNER yang boleh mendirikan cabang baru. Kasir DILARANG KERAS!
	if user.Role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Akses Ilegal! Hanya tingkat Pemilik (Owner) yang dapat mendirikan infrastruktur cabang baru.",
		})
		return
	}

	// 2. LOGIKA SUBSCRIPTION (SAAS)
	subPlan := strings.ToLower(input.Plan)
	if subPlan == "" {
		subPlan = "trial"
	}

	subIndustry := input.Industry
	if subIndustry == "" {
		subIndustry = "retail"
	}

	var subEnd time.Time
	subStatus := "inactive" // Default berbayar adalah inactive (Nunggak)

	if subPlan == "trial" {
		subStatus = "active"                  // Trial langsung aktif gratis
		subEnd = time.Now().AddDate(0, 0, 14) // Trial 14 Hari
	} else {
		// Kalau berbayar, set expired-nya detik ini juga (biar ke-lock sebelum bayar via Midtrans)
		subEnd = time.Now()
	}

	// SETTING FITUR BERDASARKAN KASTA
	fiturAktif := `["kasir"]`
	if subPlan == "pro" {
		fiturAktif = `["kasir", "absensi", "export_excel"]`
	} else if subPlan == "premium" || subPlan == "trial" {
		fiturAktif = `["kasir", "absensi", "export_excel", "multi_gudang", "ai_analyst", "whatsapp"]`
	}

	// SUSUN MODEL DATABASE TOKO BARU
	newStore := models.Store{
		OwnerID:            userID,
		NamaToko:           input.NamaToko,
		Telepon:            input.Telepon,
		BusinessType:       input.Business_type,
		Industry:           subIndustry,
		SubscriptionPlan:   subPlan,
		SubscriptionStatus: subStatus,
		SubscriptionEnd:    subEnd,
		FiturAktif:         fiturAktif,
		Alamat:             input.AlamatJalan,
		Provinsi:           input.Provinsi,
		Kota:               input.Kota,
		Kecamatan:          input.Kecamatan,
		Kelurahan:          input.Kelurahan,
		KodePos:            input.KodePos,
	}

	// 5. TRANSAKSI DATABASE (Simpan Toko + Update ID di User)
	errTx := src.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newStore).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.User{}).Where("id = ?", userID).Update("store_id", newStore.ID).Error; err != nil {
			return err
		}
		return nil
	})

	if errTx != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan infrastruktur cabang baru. Hubungi tim teknis."})
		return
	}

	// 6. TERBITIN KARTU AKSES (JWT) BARU YANG UDAH ADA STORE ID-NYA
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   userID,
		"store_id":  newStore.ID,
		"plan_type": newStore.SubscriptionPlan,
		"role":      "owner",
		"exp":       time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := newToken.SignedString([]byte("KUNCI_RAHASIA_SUPER_KUAT_123"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Toko sukses dibuat, namun gagal memperbarui token akses"})
		return
	}

	// 7. RESPONSE KEMBALI KE VUE
	c.JSON(http.StatusOK, gin.H{
		"message":  "Konfigurasi sistem siap! Selamat datang di platform POS SaaS.",
		"store_id": newStore.ID,
		"token":    tokenString,
		"data": gin.H{
			"nama_toko":     newStore.NamaToko,
			"business_type": newStore.BusinessType,
		},
	})
}
