package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"pos-backend/controllers/auth"
	src "pos-backend/src/core/config"
	"pos-backend/src/core/middlewares"

	fnbDelivery "pos-backend/src/modules/fnb/delivery"
	fnbRepository "pos-backend/src/modules/fnb/repository"

	laundryDelivery "pos-backend/src/modules/jasalayanan/laundry/delivery"
	laundryRepository "pos-backend/src/modules/jasalayanan/laundry/repository"

	retailDelivery "pos-backend/src/modules/retail/delivery"
	retailRepository "pos-backend/src/modules/retail/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan src default system")
	}

	src.ConnectDatabase()
	r := gin.Default()

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

	retailRepo := retailRepository.NewRetailRepo(src.DB)
	retailHandler := retailDelivery.NewRetailHandler(retailRepo)
	r.POST("/api/retail/midtrans/webhook", retailHandler.MidtransWebhook)

	// -- Rute Terproteksi (Butuh Karcis JWT) --
	api := r.Group("/api")
	api.Use(middlewares.RequireAuth)
	{
		// --- RUTE SETUP TOKO & AKUN ---
		api.GET("/me", auth.GetMe)
		api.PUT("/profile", auth.UpdateProfile)
		api.PUT("/password", auth.UpdatePassword)
		api.POST("/auth/select-store", auth.SelectStore)

		// 🚀 INI DIA HASILNYA! BERSIH BANGET KAN?
		api.POST("/setup", auth.SetupTokoBaru)

		// ==========================================
		// 🛒 RUTE KHUSUS RETAIL (VERSI BERSIH MODULAR)
		// ==========================================
		retailAPI := api.Group("/retail")
		{
			retailDelivery.RegisterRetailInventoryRoutes(retailAPI, retailHandler)
		}

		// ==========================================
		// 🧺 RUTE KHUSUS LAUNDRY (MODULAR LAYER - OPSI B)
		// ==========================================
		laundryAPI := api.Group("/laundry")
		laundryRepo := laundryRepository.NewLaundryRepo(src.DB)
		laundryHandler := &laundryDelivery.LaundryHandler{Repo: laundryRepo}
		laundryDelivery.RegisterLaundryRoutes(laundryAPI, laundryHandler)

		// ==========================================
		// 🍔 RUTE KHUSUS FOOD AND BEVERAGES
		// ==========================================
		fnbAPI := api.Group("/fnb")
		fnbAPI.Use(middlewares.RequireAuth)
		fnbMenuRepo := fnbRepository.NewMenuRepo(src.DB)
		fnbMenuHandler := &fnbDelivery.MenuHandler{Repo: fnbMenuRepo}
		fnbOrderRepo := fnbRepository.NewOrderRepo(src.DB)
		fnbOrderHandler := fnbDelivery.NewOrderHandler(fnbOrderRepo)
		fnbDelivery.RegisterFnBRoutes(fnbAPI, fnbMenuHandler, fnbOrderHandler)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server berjalan di port: " + port)
	r.Run("0.0.0.0:" + port)
}
