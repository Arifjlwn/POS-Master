package admin

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"pos-backend/models"
	"pos-backend/utils" // 🚀 IMPORT WAJIB: Panggil package utils lu bray!

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminController struct {
	DB *gorm.DB
}

// Struct payload input dari frontend Vue
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Struct penampung response flat tabel toko (Prioritas #2)
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

// ==========================================
// 🔐 OTO-1: LOGIN HANDLER PUSAT SUPER ADMIN
// ==========================================
func (a *AdminController) AdminLogin(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data input login tidak valid bray!"})
		return
	}

	var user models.User
	if err := a.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kredensial login admin salah atau tidak ditemukan!"})
		return
	}

	if user.Role != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "AKSES DITOLAK: Akun Anda tidak memegang kedaulatan Root Admin!"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kredensial login admin salah atau tidak ditemukan!"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"email":   user.Email,
		"exp":     time.Now().Add(12 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengonfigurasi token kedaulatan pusat!"})
		return
	}

	// 🚀 CATAT LOG: Super Admin Login
	utils.RecordAdminAction(c, "Otorisasi Super Admin", "-", "Berhasil masuk ke platform Mission Control")

	c.JSON(http.StatusOK, gin.H{
		"status": "sukses",
		"token":  tokenString,
		"user": gin.H{
			"name":  user.Name,
			"role":  user.Role,
			"email": user.Email,
		},
	})
}

// ==========================================
// 🏪 OTO-2: AMBIL DATA FLAT MITRA TOKO (STORES)
// ==========================================
func (a *AdminController) GetAllTenants(c *gin.Context) {
	var results []FlatStoreResponse

	err := a.DB.Table("stores").
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
		Joins("LEFT JOIN users ON users.id = stores.owner_id").
		Order("stores.created_at desc").
		Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyatukan relasi toko dan owner!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "sukses",
		"data":   results,
	})
}

// ==========================================
// 🔒 OTO-3: SUSPEND/BEKUKAN LANGGANAN TOKO
// ==========================================
func (a *AdminController) SuspendTenant(c *gin.Context) {
	storeID := c.Param("id")

	err := a.DB.Model(&models.Store{}).Where("id = ?", storeID).Update("subscription_status", "suspended").Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membekukan toko!"})
		return
	}

	// 🚀 CATAT LOG: Eksekusi Suspend Tenant
	utils.RecordAdminAction(c, "Pembekuan Otoritas", storeID, fmt.Sprintf("Membekukan akses kasir toko ID: %s", storeID))

	c.JSON(http.StatusOK, gin.H{
		"status":  "sukses",
		"message": "Toko resmi dibekukan! Gerbang kasir otomatis terkunci.",
	})
}

// ==========================================
// 🟢 OTO-4: AKTIFKAN & PERPANJANG MASA AKTIF TOKO
// ==========================================
func (a *AdminController) ActivateTenant(c *gin.Context) {
	storeID := c.Param("id")
	newExpiry := time.Now().AddDate(0, 1, 0)

	err := a.DB.Model(&models.Store{}).Where("id = ?", storeID).Updates(map[string]interface{}{
		"subscription_status": "active",
		"subscription_end":    &newExpiry,
	}).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengaktifkan toko!"})
		return
	}

	// 🚀 CATAT LOG: Eksekusi Activate Tenant
	utils.RecordAdminAction(c, "Aktivasi Layanan", storeID, fmt.Sprintf("Memperpanjang masa aktif toko ID: %s", storeID))

	c.JSON(http.StatusOK, gin.H{
		"status":  "sukses",
		"message": "Toko diaktifkan kembali! Masa aktif diperpanjang 30 hari ke depan.",
	})
}

// ==========================================
// 📊 OTO-5: 100% REALTIME TELEMETRI MONITOR PUSAT
// ==========================================
func (a *AdminController) GetTelemetryStats(c *gin.Context) {
	// 1. DATA METRIK TOKO (REAL DARI DB)
	var totalStores, activeStores, suspendedStores int64
	a.DB.Model(&models.Store{}).Count(&totalStores)
	a.DB.Model(&models.Store{}).Where("subscription_status = ?", "active").Count(&activeStores)
	a.DB.Model(&models.Store{}).Where("subscription_status = ?", "suspended").Count(&suspendedStores)

	// 2. KESEHATAN SERVER (REAL DARI RUNTIME GO & DB PING)
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	ramUsageMB := m.Alloc / 1024 / 1024
	cpuLoad := runtime.NumGoroutine()

	startPing := time.Now()
	sqlDB, err := a.DB.DB()
	dbStatus := "Mati / Putus"
	if err == nil && sqlDB.Ping() == nil {
		dbStatus = "Online & Stabil"
	}
	latency := time.Since(startPing).Milliseconds()

	// 3. SINKRONISASI AKTIVITAS TERBARU DARI TABEL AUDIT_LOGS
	type SqlAuditResponse struct {
		ID        uint
		Action    string
		Details   string
		UserName  string
		CreatedAt time.Time
	}
	var dbLogs []SqlAuditResponse

	err = a.DB.Table("audit_logs").
		Select(`
			audit_logs.id, 
			audit_logs.action, 
			audit_logs.details, 
			COALESCE(users.name, audit_logs.user_email) as user_name, 
			audit_logs.created_at
		`).
		Joins("LEFT JOIN users ON users.id = audit_logs.user_id").
		Order("audit_logs.created_at desc").
		Limit(5).
		Scan(&dbLogs).Error

	var finalActivities = []gin.H{}
	if err == nil {
		for _, logItem := range dbLogs {
			displayAction := logItem.Action
			if logItem.Details != "" {
				displayAction = fmt.Sprintf("%s (%s)", logItem.Action, logItem.Details)
			}

			finalActivities = append(finalActivities, gin.H{
				"id":     logItem.ID,
				"action": displayAction,
				"user":   logItem.UserName,
				"time":   logItem.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}
	}

	// 4. LEMPAR PAYLOAD REALTIME
	c.JSON(http.StatusOK, gin.H{
		"status": "sukses",
		"data": gin.H{
			"total_stores":     totalStores,
			"active_stores":    activeStores,
			"suspended_stores": suspendedStores,
		},
		"server_health": gin.H{
			"cpu_usage": cpuLoad,
			"ram_usage": ramUsageMB,
			"db_status": dbStatus,
			"latency":   fmt.Sprintf("%dms", latency),
		},
		"recent_activities": finalActivities,
	})
}
