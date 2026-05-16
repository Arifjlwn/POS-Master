package main

import (
	"log"
	"net/http"
	"os"
	"pos-backend/config"
	"pos-backend/controllers"
	"pos-backend/middlewares"
	"pos-backend/models"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func main() {
	// GO Baca file .env difolder ini
	err := godotenv.Load()
		if err != nil {
			log.Println("Peringatan: File .env tidak ditemukan, menggunakan config default system")
		}

	// Inisialisasi Database
	config.ConnectDatabase()

	// Setup router baru dari GIN
	r := gin.Default()
	
	// --- PERBAIKAN CORS ---
	r.Use(cors.New(cors.Config{
        // 🚀 Mengizinkan localhost DAN semua IP Local Network (192.168.x.x) secara dinamis
        AllowOriginFunc: func(origin string) bool {
            // Tetap izinkan localhost untuk development di PC
            if origin == "http://localhost:5173" || origin == "http://localhost:5174" {
                return true
            }
            // Izinkan otomatis semua perangkat HP yang terhubung satu WiFi/Hotspot (192.168.xx.xx)
            return true 
        },
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With", "Accept"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    r.Static("/uploads", "./uploads")

	r.Use(func(c *gin.Context) {
    // Batasi 5MB secara manual
    c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 5*1024*1024)
    c.Next()
})

	// Membuat endpoint API sederhana (Route GET)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "sukses",
			"message": "Halo Bos ! Server Go Berhasil Menyala !",
		})
	})

	// -- Rute API SAAS --
	r.POST("/api/register", controllers.Register)
	r.POST("/api/verify-otp", controllers.VerifyOTP)
	r.POST("/api/login", controllers.Login)

	// -- Rute Terproteksi (Butuh Karcis JWT) --
	api := r.Group("/api")
	api.Use(middlewares.RequireAuth)
	{
		// -- Rute Produk (CRUD) --
		api.POST("/products", controllers.CreateProduct)
		api.GET("/products", controllers.GetProducts)
		api.PUT("/products/:id", controllers.UpdateProduct)
		api.DELETE("/products/:id", controllers.DeleteProduct)
		api.GET("/products/export", controllers.ExportProducts)
		api.POST("/products/import", controllers.ImportProducts)
		api.GET("/categories", controllers.GetCategories)

		// -- Rute Karyawan --
		api.POST("/employees", controllers.CreateEmployee)
		api.GET("/employees", controllers.GetEmployees)
		api.PUT("/employees/:id", controllers.UpdateEmployee)
		api.GET("/me", controllers.GetMe)

		// -- RUTE TOKO SHIFT MANAGEMENT (TSM)
		api.POST("/schedules/bulk", controllers.SaveSchedules)
		api.GET("/schedules", controllers.GetSchedules)

		// --RUTE ABSENSI--
		api.POST("/attendance", controllers.StoreAttendance)
		api.GET("/attendance", controllers.GetAttendance)
		api.GET("/attendance/export", controllers.ExportAttendance)

		//--- RUTE RETUR & BARANG RUSAK ---
		api.POST("/returns", controllers.CreateReturn)
		api.GET("/returns", controllers.GetReturns)

		// Rute LPB
		api.POST("/purchases", controllers.CreateLPB)

		// Rute Stock Opname
		api.POST("/stock-opname", controllers.CreateStockOpname)
		api.GET("/stock-opname/history", controllers.GetStockOpnameHistory)

		// Rute Session
		// 🚀 --- RUTE SESSION KASIR (BUKA/CEK KASIR) ---
        api.POST("/pos/open-session", controllers.OpenSession)
        api.GET("/pos/check-session", controllers.CheckSessionStatus)

		// --- RUTE CLOSING ---
		api.POST("/pos/close-session/:id", controllers.CloseSession)
		
		// Rute Transaksi (Mesin Kasir)
		api.POST("/checkout", controllers.CreateTransaction)
		api.GET("/transactions", controllers.GetTransactions)
		

		// Rute Laporan (Dashboard)
		api.GET("/report/dashboard", controllers.GetDashboardReport)

		// --- RUTE SETUP TOKO ---
		api.POST("/setup-toko", func(c *gin.Context) {
			// 1. Ambil User ID dari Satpam JWT (Ingat, dari JWT itu bentuknya float64, harus diconvert)
			userIDRaw, exists := c.Get("user_id")
			if !exists {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak valid"})
				return
			}
			userID := uint(userIDRaw.(float64))

			// 2. Tangkap JSON dari Frontend Vue
			var input struct {
				NamaToko   string `json:"nama_toko" binding:"required"`
				TipeBisnis string `json:"tipe_bisnis" binding:"required"`
				AlamatToko string `json:"alamat_toko"`
				Telepon    string `json:"telepon"`
			}
			
			if err := c.ShouldBindJSON(&input); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Format data salah, pastikan kirim JSON"})
				return
			}

			// 3. Masukkan ke tabel Stores di Supabase
			newStore := models.Store{
				NamaToko:     input.NamaToko,
				BusinessType: input.TipeBisnis,
				Alamat:       input.AlamatToko,
				Telepon:      input.Telepon,
			}
			
			if err := config.DB.Create(&newStore).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat toko"})
				return
			}

			// 4. Update Akun Bos! (Kaitkan ID User dengan ID Toko yang baru dibuat)
			if err := config.DB.Model(&models.User{}).Where("id = ?", userID).Update("store_id", newStore.ID).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengaitkan toko ke akun Anda"})
				return
			}

			// 🚀 OPERASI SENYAP DIMULAI DI SINI: Cetak ulang tiket VIP!
			newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"user_id":  userID,
				"store_id": newStore.ID, // 👈 Sekarang udah terisi ID Toko yang baru dibuat!
				"role":     "owner",
				"exp":      time.Now().Add(time.Hour * 72).Unix(),
			})

			tokenString, err := newToken.SignedString([]byte("KUNCI_RAHASIA_SUPER_KUAT_123"))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Toko dibuat, tapi gagal cetak tiket baru"})
				return
			}

			// 5. Berhasil!
			c.JSON(http.StatusOK, gin.H{
				"message": "Toko berhasil dibuat! Selamat berbisnis.",
				"store_id": newStore.ID,
				"token": tokenString,
			})
		})
	}

	// Menyalakan server dari port .env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server berjalan dan terbuka untuk lokal network di port: " + port)
	r.Run("0.0.0.0:" + port)
}