package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"pos-backend/models"
	src "pos-backend/src/core/config"
	"pos-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/gorm"
)

func SetupTokoBaru(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak valid"})
		return
	}
	userID := uint(userIDRaw.(float64))

	var input struct {
		NamaToko     string `json:"nama_toko" binding:"required"`
		Telepon      string `json:"telepon"`
		BusinessType string `json:"business_type" binding:"required"`
		Industry     string `json:"industry"`
		Plan         string `json:"plan"`
		AlamatJalan  string `json:"alamat_toko"`
		Provinsi     string `json:"provinsi"`
		Kota         string `json:"kota"`
		Kecamatan    string `json:"kecamatan"`
		Kelurahan    string `json:"kelurahan"`
		KodePos      string `json:"kode_pos"`
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
	if user.Role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses Ilegal! Hanya tingkat Pemilik (Owner) yang dapat mendirikan infrastruktur cabang baru."})
		return
	}

	subPlan := strings.ToLower(input.Plan)
	if subPlan == "" {
		subPlan = "trial"
	}
	subIndustry := input.Industry
	if subIndustry == "" {
		subIndustry = "retail"
	}

	var subEnd time.Time
	var subStatus string
	var fiturAktif string

	switch subPlan {
	case "trial":
		subStatus = "active"
		subEnd = time.Now().AddDate(0, 0, 14) // 14 Hari Uji Coba Gratis bray
		fiturAktif = `["kasir", "absensi", "export_excel", "multi_gudang", "ai_analyst", "whatsapp"]`

	case "pro":
		subStatus = "pending" // 🔒 Gembok 'pending' sebelum settlement Midtrans
		subEnd = time.Now()   // Nanti di-update otomatis lewat webhook sukses
		fiturAktif = `["kasir", "absensi", "export_excel"]`

	case "premium":
		subStatus = "pending" // 🔒 Gembok 'pending' sebelum settlement Midtrans
		subEnd = time.Now()
		fiturAktif = `["kasir", "absensi", "export_excel", "multi_gudang", "whatsapp"]`

	default: // Fallback ke paket "basic" jika kosong atau tidak cocok bray
		subPlan = "basic"
		subStatus = "pending"
		subEnd = time.Now()
		fiturAktif = `["kasir"]`
	}

	newStore := models.Store{
		PublicID:           utils.GenerateULID(),
		OwnerID:            userID,
		NamaToko:           input.NamaToko,
		Telepon:            input.Telepon,
		BusinessType:       input.BusinessType,
		Industry:           subIndustry,
		SubscriptionPlan:   subPlan,
		SubscriptionStatus: subStatus,
		SubscriptionEnd:    &subEnd,
		FiturAktif:         fiturAktif,
		Alamat:             input.AlamatJalan,
		Provinsi:           input.Provinsi,
		Kota:               input.Kota,
		Kecamatan:          input.Kecamatan,
		Kelurahan:          input.Kelurahan,
		KodePos:            input.KodePos,
	}

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

	jwtSecret := os.Getenv("JWT_SECRET")

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
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Toko sukses dibuat, namun gagal memperbarui token akses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":             "Konfigurasi sistem siap! Selamat datang di platform POS SaaS.",
		"store_id":            newStore.ID,
		"token":               tokenString,
		"store_name":          newStore.NamaToko,
		"subscription_plan":   newStore.SubscriptionPlan,
		"subscription_status": newStore.SubscriptionStatus, // Kirim balik ke Vue buat validasi routing frontend
		"data":                gin.H{"nama_toko": newStore.NamaToko, "business_type": newStore.BusinessType},
	})
}

func ReTriggerPaymentHandler(c *gin.Context) {
	// 1. Ambil payload JSON data store_id dari Vue
	var input struct {
		StoreID uint `json:"store_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID gerai tidak valid"})
		return
	}

	// 2. Query data ruko/toko dari DB PostgreSQL
	var store models.Store
	if err := src.DB.First(&store, input.StoreID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Infrastruktur data toko tidak ditemukan"})
		return
	}

	// 3. Hitung harga/nominal invoice berdasarkan tipe plan toko tersebut bray
	var finalPrice int64 = 0
	targetPlan := strings.ToLower(strings.TrimSpace(store.SubscriptionPlan))

	switch targetPlan {
	case "basic":
		finalPrice = 49000
	case "pro":
		finalPrice = 149000
	case "premium":
		finalPrice = 299000
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tier tingkatan paket toko tidak terdaftar"})
		return
	}

	// 4. Inisialisasi Environment Gateway Midtrans bray
	// Variabel 'env' lu panggil di sini untuk menentukan core routing server Midtrans!
	var env midtrans.EnvironmentType = midtrans.Sandbox
	if os.Getenv("APP_ENV") == "production" {
		env = midtrans.Production
	}

	// Suntik kredensial server key ke SDK Midtrans Go global instance bray
	midtrans.ServerKey = os.Getenv("MIDTRANS_SERVER_KEY")
	midtrans.Environment = env

	// 5. Racik kode Order ID unik ber-prefix "UPGRADE" agar otomatis lolos sensor Webhook lu bray!
	planCode := strings.ReplaceAll(strings.ToUpper(targetPlan), " ", "")
	orderID := fmt.Sprintf("UPGRADE-TOKO-%d-%s-%d", store.ID, planCode, time.Now().Unix())

	// 6. Buat payload Request Snap untuk dikirim ke API Midtrans
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{OrderID: orderID, GrossAmt: finalPrice},
		Items:              &[]midtrans.ItemDetails{{ID: "SUB-" + planCode, Price: finalPrice, Qty: 1, Name: "Aktivasi Langganan " + store.SubscriptionPlan}},
	}

	// 7. Tembak server Midtrans buat nge-generate snapToken riil bray!
	snapResp, errSnap := snap.CreateTransaction(req)
	if errSnap != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal bernegosiasi ulang dengan Midtrans Gateway"})
		return
	}

	// 8. Return response sukses 200 OK ke SetupToko.vue lu bray!
	c.JSON(http.StatusOK, gin.H{
		"message":    "Sesi billing pembayaran berhasil di-restore!",
		"store_id":   store.ID,
		"store_name": store.NamaToko,
		"snap_token": snapResp.Token, // ◄ Token ril siap meledakkan modal di Vue lu bray!
		"order_id":   orderID,
	})
}
