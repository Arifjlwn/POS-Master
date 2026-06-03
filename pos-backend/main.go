package main

import (
	"log"
	"net/http"
	"os"
	"pos-backend/controllers/auth"
	"pos-backend/models"
	src "pos-backend/src/core/config"
	"pos-backend/src/core/middlewares"
	"strings"
	"time"

	"gorm.io/gorm"

	// 🚀 ALIAS DIBAWAH INI KUNCI BIAR GAK COLLISION
	fnbDelivery "pos-backend/src/modules/fnb/delivery"
	fnbRepository "pos-backend/src/modules/fnb/repository"

	laundryDelivery "pos-backend/src/modules/jasalayanan/laundry/delivery"
	laundryRepository "pos-backend/src/modules/jasalayanan/laundry/repository"

	retailDelivery "pos-backend/src/modules/retail/delivery"
	retailRepository "pos-backend/src/modules/retail/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func main() {
	// GO Baca file .env difolder ini
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan src default system")
	}

	// Inisialisasi Database (Pake package src/core/config)
	src.ConnectDatabase()

	// Setup router baru dari GIN
	r := gin.Default()

	// --- PERBAIKAN CORS ---
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			if origin == "http://localhost:5173" || origin == "http://localhost:5174" {
				return true
			}
			return true
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Static("/uploads", "./uploads")
	r.Static("/public", "./public")

	r.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 5*1024*1024)
		c.Next()
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "sukses",
			"message": "Halo Bos ! Server Go Berhasil Menyala !",
		})
	})

	// ==========================================
	// -- RUTE API SAAS (GLOBAL & AUTH) --
	// ==========================================
	r.POST("/api/register", auth.Register)
	r.POST("/api/verify-otp", auth.VerifyOTP)
	r.POST("/api/login", auth.Login)
	r.POST("/api/auth/send-otp-wa", auth.SendOTPWhatsApp)
	r.POST("/api/reset-password", auth.ResetPassword)
	r.POST("/api/auth/check-account", auth.CheckAccount)

	// 1. TAMBAHKAN INI: Inisialisasi Retail Handler di sini (di luar rute terproteksi)
	retailRepo := retailRepository.NewRetailRepo(src.DB)
	retailHandler := retailDelivery.NewRetailHandler(retailRepo)
	r.POST("/api/retail/midtrans/webhook", retailHandler.MidtransWebhook)

	// -- Rute Terproteksi (Butuh Karcis JWT) --
	api := r.Group("/api")
	api.Use(middlewares.RequireAuth)
	{
		// --- RUTE SETUP TOKO (Global - Masih inline) ---
		api.GET("/me", auth.GetMe)
		api.PUT("/profile", auth.UpdateProfile)
		api.PUT("/password", auth.UpdatePassword)
		api.POST("/auth/select-store", auth.SelectStore)

		// --- UPDATE INTEGRASI RUTE SETUP TOKO ---
		api.POST("/setup", func(c *gin.Context) {
			userIDRaw, exists := c.Get("user_id")
			if !exists {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak valid"})
				return
			}
			userID := uint(userIDRaw.(float64))

			// 🚀 1. TANGKEP DATA DARI VUE TERMASUK INDUSTRY & PLAN SAAS
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
			if user.StoreID != nil {
				c.JSON(http.StatusConflict, gin.H{"error": "Sistem mendeteksi Anda sudah memiliki toko yang terdaftar."})
				return
			}

			// 🚀 2. LOGIKA SUBSCRIPTION (SAAS)
			subPlan := strings.ToLower(input.Plan)
			if subPlan == "" {
				subPlan = "trial"
			}

			subIndustry := input.Industry
			if subIndustry == "" {
				subIndustry = "retail"
			}

			var subEnd time.Time
			subStatus := "inactive" // 🚀 Default berbayar adalah inactive (Nunggak)

			if subPlan == "trial" {
				subStatus = "active"                  // Trial langsung aktif gratis
				subEnd = time.Now().AddDate(0, 0, 14) // Trial 14 Hari
			} else {
				// Kalau berbayar, set expired-nya detik ini juga (biar ke-lock sebelum bayar)
				subEnd = time.Now()
			}

			// 🚀 SETTING FITUR
			fiturAktif := `["kasir"]`
			if subPlan == "pro" {
				fiturAktif = `["kasir", "absensi", "export_excel"]`
			} else if subPlan == "premium" || subPlan == "trial" {
				fiturAktif = `["kasir", "absensi", "export_excel", "multi_gudang", "ai_analyst", "whatsapp"]`
			}

			// 🚀 SUSUN MODEL DATABASE TOKO BARU
			newStore := models.Store{
				OwnerID:            userID,
				NamaToko:           input.NamaToko,
				Telepon:            input.Telepon,
				BusinessType:       input.Business_type,
				Industry:           subIndustry,
				SubscriptionPlan:   subPlan,
				SubscriptionStatus: subStatus, // 🚀 Gunakan status dinamis hasil filter di atas!
				SubscriptionEnd:    subEnd,
				FiturAktif:         fiturAktif,
				Alamat:             input.AlamatJalan,
				Provinsi:           input.Provinsi,
				Kota:               input.Kota,
				Kecamatan:          input.Kecamatan,
				Kelurahan:          input.Kelurahan,
				KodePos:            input.KodePos,
			}

			// 🚀 5. TRANSAKSI DATABASE (Simpan Toko + Update ID di User)
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
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan infrastruktur toko baru. Hubungi tim teknis."})
				return
			}

			// 🚀 6. TERBITIN KARTU AKSES (JWT) BARU YANG UDAH ADA STORE ID-NYA
			newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"user_id":  userID,
				"store_id": newStore.ID,
				"role":     "owner",
				"exp":      time.Now().Add(time.Hour * 72).Unix(),
			})

			tokenString, err := newToken.SignedString([]byte("KUNCI_RAHASIA_SUPER_KUAT_123"))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Toko sukses dibuat, namun gagal memperbarui token akses"})
				return
			}

			// 🚀 7. RESPONSE KEMBALI KE VUE
			c.JSON(http.StatusOK, gin.H{
				"message":  "Konfigurasi sistem siap! Selamat datang di platform POS SaaS.",
				"store_id": newStore.ID,
				"token":    tokenString,
				"data": gin.H{
					"nama_toko":     newStore.NamaToko,
					"business_type": newStore.BusinessType,
				},
			})
		})

		// ==========================================
		// 🛒 RUTE KHUSUS RETAIL (VERSI BERSIH MODULAR)
		// ==========================================
		retailAPI := api.Group("/retail")
		{
			// 🚀 KUNCI UTAMA: WIRING MODUL MODULAR RETAIL YANG BARU
			// Ini bakal nge-handle: Schedule, Session Kasir, Checkout POS, Absensi, Retur, LPB, dan Dashboard SO!
			// retailRepo := retailRepository.NewRetailRepo(src.DB)
			// retailHandler := retailDelivery.NewRetailHandler(retailRepo)
			retailDelivery.RegisterRetailInventoryRoutes(retailAPI, retailHandler)
		}

		// ==========================================
		// 🧺 RUTE KHUSUS LAUNDRY (MODULAR LAYER - OPSI B)
		// ==========================================
		laundryAPI := api.Group("/laundry")

		// Inisialisasi Dependency Pake Nama Alias Repository & Delivery Khusus Laundry
		laundryRepo := laundryRepository.NewLaundryRepo(src.DB)
		laundryHandler := &laundryDelivery.LaundryHandler{Repo: laundryRepo}

		// Daftarkan rute laundry
		laundryDelivery.RegisterLaundryRoutes(laundryAPI, laundryHandler)

		// ==========================================
		//      RUTE KHUSUS FOOD AND BEVERAGES
		// ==========================================
		fnbAPI := api.Group("/fnb")
		fnbAPI.Use(middlewares.RequireAuth)

		// Init Master Menu Pake Alias Fnb
		fnbMenuRepo := fnbRepository.NewMenuRepo(src.DB)
		fnbMenuHandler := &fnbDelivery.MenuHandler{Repo: fnbMenuRepo}

		// Init Transaksi & Dapur Pake Alias Fnb
		fnbOrderRepo := fnbRepository.NewOrderRepo(src.DB)
		fnbOrderHandler := fnbDelivery.NewOrderHandler(fnbOrderRepo)

		// Daftarkan ke rute fnb
		fnbDelivery.RegisterFnBRoutes(fnbAPI, fnbMenuHandler, fnbOrderHandler)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server berjalan di port: " + port)
	r.Run("0.0.0.0:" + port)
}
