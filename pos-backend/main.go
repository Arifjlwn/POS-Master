package main

import (
	"time"
	"net/http"
	"pos-backend/config"
	"pos-backend/controllers"
	"pos-backend/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
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

		// -- Rute Karyawan --
		api.POST("/employees", controllers.CreateEmployee)
		api.GET("/employees", controllers.GetEmployees)

		// Rute Transaksi (Mesin Kasir)
		api.POST("/checkout", controllers.CreateTransaction)

		// Rute Laporan (Dashboard)
		api.GET("/report/dashboard", controllers.GetDashboardReport)

		// --- PERBAIKAN RUTE SETUP TOKO ---
		// Gunakan "api" bukan "protected"
		api.POST("/setup-toko", func(c *gin.Context) {
			var input struct {
				NamaToko   string `json:"nama_toko"`
				TipeBisnis string `json:"tipe_bisnis"`
				AlamatToko string `json:"alamat_toko"`
				Telepon    string `json:"telepon"`
			}
			
			if err := c.ShouldBindJSON(&input); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Format data salah, pastikan kirim JSON"})
				return
			}

			// (Opsional) Di sini nanti kodenya untuk save ke tabel Stores pakai db.Create(&store)
			// Untuk sekarang kita kasih respon sukses dulu
			c.JSON(http.StatusOK, gin.H{
				"message": "Toko berhasil dibuat!",
				"data":    input,
			})
		})
	}

	// Menyalakan server di port 8080
	r.Run(":8080")
}