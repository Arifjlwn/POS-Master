package main

import (
	"net/http"
	"pos-backend/config"
	"pos-backend/controllers"
	"pos-backend/middlewares"
	"pos-backend/models"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	// Inisialisasi Database
	config.ConnectDatabase()

	// Setup router baru dari GIN
	r := gin.Default()
	
	// --- PERBAIKAN CORS ---
	r.Use(cors.New(cors.Config{
		// Wajib spesifik alamat Vue-nya, jangan pakai "*" karena ada AllowCredentials
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174"}, 
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Static("/uploads", "./uploads")

	// Membuat endpoint API sederhana (Route GET)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "sukses",
			"message": "Halo Bos ! Server Go Berhasil Menyala !",
		})
	})

	// -- Rute API SAAS --
	r.POST("/api/register", controllers.Register)
	r.POST("/api/login", controllers.Login)

	// -- Rute Terproteksi (Butuh Karcis JWT) --
	api := r.Group("/api")
	api.Use(middlewares.RequireAuth)
	{
		// Rute untuk melihat profil sendiri
		api.GET("/me", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			storeID, _ := c.Get("store_id")
			role, _ := c.Get("role")

			c.JSON(http.StatusOK, gin.H{
				"message":  "Ini adalah area rahasia",
				"user_id":  userID,
				"store_id": storeID,
				"role":     role,
			})
		})

		// -- Rute Produk (CRUD) --
		api.POST("/products", controllers.CreateProduct)
		api.GET("/products", controllers.GetProducts)
		api.PUT("/products/:id", controllers.UpdateProduct)
		api.DELETE("/products/:id", controllers.DeleteProduct)
		api.GET("/products/export", controllers.ExportProducts)
		api.GET("/categories", controllers.GetCategories)

		// -- Rute Karyawan --
		api.POST("/employees", controllers.CreateEmployee)
		api.GET("/employees", controllers.GetEmployees)

		// --RUTE ABSENSI--
		// --- MODULE ABSENSI ---
		api.POST("/attendance", controllers.StoreAttendance)
		api.GET("/attendance", controllers.GetAttendance)
		api.GET("/attendance/export", controllers.ExportAttendance)

		// Rute Session
		// 🚀 --- RUTE SESSION KASIR (BUKA/CEK KASIR) ---
        api.POST("/pos/open-session", controllers.OpenSession)
        api.GET("/pos/check-session", controllers.CheckSessionStatus)
		
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

	// Menyalakan server di port 8080
	r.Run(":8080")
}