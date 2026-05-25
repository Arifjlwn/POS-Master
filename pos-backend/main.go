package main

import (
	"log"
	"net/http"
	"os"
	"pos-backend/src/core/config"
	"pos-backend/controllers/auth"
	"pos-backend/src/core/middlewares"
	"pos-backend/models"
	"time"

	// 🚀 ALIAS DIBAWAH INI KUNCI BIAR GAK COLLISION / BENTROK KONTOL
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
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
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

	// -- Rute Terproteksi (Butuh Karcis JWT) --
	api := r.Group("/api")
	api.Use(middlewares.RequireAuth)
	{
		// --- RUTE SETUP TOKO (Global - Masih inline) ---
		api.GET("/me", auth.GetMe)
		
		// --- UPDATE INTEGRASI RUTE SETUP TOKO ---
		api.POST("/setup", func(c *gin.Context) {
			userIDRaw, exists := c.Get("user_id")
			if !exists {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak valid"})
				return
			}
			userID := uint(userIDRaw.(float64))

			var input struct {
				NamaToko      string   `json:"nama_toko" binding:"required"`
				Business_type string   `json:"business_type" binding:"required"` 
				AlamatJalan   string   `json:"alamat_toko"`                     
				Provinsi      string   `json:"provinsi"`
				Kota          string   `json:"kota"`
				Kecamatan     string   `json:"kecamatan"`
				Kelurahan     string   `json:"kelurahan"`
				KodePos       string   `json:"kode_pos"`
				Telepon       string   `json:"telepon"`
				FiturOpsional []string `json:"fitur_aktif"`                    
			}
			
			if err := c.ShouldBindJSON(&input); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Format data salah, pastikan form terisi lengkap"})
				return
			}

			var fiturString string
			for i, fitur := range input.FiturOpsional {
				if i > 0 {
					fiturString += ","
				}
				fiturString += fitur
			}

			alamatLengkap := input.AlamatJalan + ", Kel. " + input.Kelurahan + ", Kec. " + input.Kecamatan + ", " + input.Kota + ", " + input.Provinsi + " " + input.KodePos

			newStore := models.Store{
				NamaToko:     input.NamaToko,
				BusinessType: input.Business_type,
				Telepon:      input.Telepon,
				FiturAktif:   fiturString,
				Alamat:       alamatLengkap, 
				Provinsi:     input.Provinsi, 
				Kota:         input.Kota,
				Kecamatan:    input.Kecamatan,
				Kelurahan:    input.Kelurahan,
				KodePos:      input.KodePos,
			}
			
			if err := src.DB.Create(&newStore).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan infrastruktur toko baru"})
				return
			}

			if err := src.DB.Model(&models.User{}).Where("id = ?", userID).Update("store_id", newStore.ID).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengaitkan otoritas toko ke akun owner"})
				return
			}

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

			c.JSON(http.StatusOK, gin.H{
				"message":       "Konfigurasi sistem siap! Selamat datang di platform POS SaaS.",
				"store_id":      newStore.ID,
				"token":         tokenString,
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
			retailRepo := retailRepository.NewRetailRepo(src.DB)
			retailHandler := retailDelivery.NewRetailHandler(retailRepo)
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