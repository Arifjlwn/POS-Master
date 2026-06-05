package auth

import (
	"net/http"
	"os"
	"strings"
	"time"

	"pos-backend/models"
	src "pos-backend/src/core/config"
	"pos-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func SetupTokoBaru(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists { c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak valid"}); return }
	userID := uint(userIDRaw.(float64))
	
	var input struct {
		NamaToko string `json:"nama_toko" binding:"required"`; Telepon string `json:"telepon"`
		BusinessType string `json:"business_type" binding:"required"`; Industry string `json:"industry"`
		Plan string `json:"plan"`; AlamatJalan string `json:"alamat_toko"`
		Provinsi string `json:"provinsi"`; Kota string `json:"kota"`
		Kecamatan string `json:"kecamatan"`; Kelurahan string `json:"kelurahan"`; KodePos string `json:"kode_pos"`
	}
	if err := c.ShouldBindJSON(&input); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Format data salah, pastikan form terisi lengkap"}); return }
	
	var user models.User
	if err := src.DB.First(&user, userID).Error; err != nil { c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"}); return }
	if user.Role != "owner" { c.JSON(http.StatusForbidden, gin.H{"error": "Akses Ilegal! Hanya tingkat Pemilik (Owner) yang dapat mendirikan infrastruktur cabang baru."}); return }
	
	subPlan := strings.ToLower(input.Plan); if subPlan == "" { subPlan = "trial" }
	subIndustry := input.Industry; if subIndustry == "" { subIndustry = "retail" }
	
	var subEnd time.Time; subStatus := "inactive"
	if subPlan == "trial" { subStatus = "active"; subEnd = time.Now().AddDate(0, 0, 14) } else { subEnd = time.Now() }
	
	fiturAktif := `["kasir"]`
	if subPlan == "pro" { 
		fiturAktif = `["kasir", "absensi", "export_excel"]` 
	} else if subPlan == "premium" || subPlan == "trial" { 
		fiturAktif = `["kasir", "absensi", "export_excel", "multi_gudang", "ai_analyst", "whatsapp"]` 
	}
	
	newStore := models.Store{PublicID: utils.GenerateULID(), OwnerID: userID, NamaToko: input.NamaToko, Telepon: input.Telepon, BusinessType: input.BusinessType, Industry: subIndustry, SubscriptionPlan: subPlan, SubscriptionStatus: subStatus, SubscriptionEnd: &subEnd, FiturAktif: fiturAktif, Alamat: input.AlamatJalan, Provinsi: input.Provinsi, Kota: input.Kota, Kecamatan: input.Kecamatan, Kelurahan: input.Kelurahan, KodePos: input.KodePos}
	
	errTx := src.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newStore).Error; err != nil { return err }
		if err := tx.Model(&models.User{}).Where("id = ?", userID).Update("store_id", newStore.ID).Error; err != nil { return err }
		return nil
	})
	if errTx != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan infrastruktur cabang baru. Hubungi tim teknis."}); return }
	
	jwtSecret := os.Getenv("JWT_SECRET")
	
	// 🚀 FIX MUTLAK: is_select WAJIB diset false agar token ini BEBAS menembus whitelist middleware bray!
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   userID,
		"public_id": user.PublicID,
		"store_id":  newStore.ID,
		"plan_type": newStore.SubscriptionPlan,
		"role":      "owner",
		"is_select": false,
		"exp":       time.Now().Add(time.Hour * 72).Unix(),
	})
	
	tokenString, err := newToken.SignedString([]byte(jwtSecret))
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Toko sukses dibuat, namun gagal memperbarui token akses"}); return }
	
	c.JSON(http.StatusOK, gin.H{
		"message":           "Konfigurasi sistem siap! Selamat datang di platform POS SaaS.",
		"store_id":          newStore.ID,
		"token":             tokenString,
		"store_name":        newStore.NamaToko,
		"subscription_plan": newStore.SubscriptionPlan,
		"data":              gin.H{"nama_toko": newStore.NamaToko, "business_type": newStore.BusinessType},
	})
}