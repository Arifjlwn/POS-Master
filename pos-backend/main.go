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

func main () {
	// Inisialisasi Database
	config.ConnectDatabase()

	// Setup router baru dari GIN
	r:= gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Izinkan semua ruko (Frontend) ngakses ini
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With", "Accept"}, // Wajib ada Authorization untuk JWT
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Membuat endpoint API sederhana (Route GET)
	r.GET("/ping", func(c *gin.Context) {
		// Mengembalikan response dalam format JSON
		c.JSON(http.StatusOK, gin.H{
			"status" : "sukses",
			"message" : "Halo Bos ! Server Go Berhasil Menyala !",
		})
	})

	// -- Rute API SAAS --
	// Endpoint untuk registrasi UMKM (Method POST)
	r.POST("/api/register", controllers.Register)
	r.POST("/api/login", controllers.Login)

	// -- Rute Terproteksi (Butuh Karcis JWT) --
	// Semua rute didalam grup "api" ini akan dicegat dulu oleh middlewares.RequireAuth
	api := r.Group("/api")
	api.Use(middlewares.RequireAuth)
	{
		// Contoh: Rute untuk melihat profil sendiri
		api.GET("/me", func(c *gin.Context) {
			// Mengambil data yang dibongkar satpam tadi
			userID, _ := c.Get("user_id")
			storeID, _ := c.Get("store_id")
			role, _ := c.Get("role")

			c.JSON(http.StatusOK, gin.H{
				"message": "Ini adalah area rahasia",
				"user_id": userID,
				"store_id": storeID,
				"role": role,
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

		// Ruter Laporan (Dashboard)
		api.GET("/report/dashboard", controllers.GetDashboardReport)
	}

	// Menyalakan server di port 8080
	r.Run(":8080")
}