package main

import (
	"log"
	"net/http"
	"os"
	"strings"
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
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan src default system")
	}

	src.ConnectDatabase()
	r := gin.Default()

	// 🚀 SANGAR: CORS Tight Guard dinamis berbasis ENV untuk mengunci serangan XSS di Production
	// 🚀 SANGAR: CORS Tight Guard dinamis berbasis ENV
    r.Use(cors.New(cors.Config{
        AllowOriginFunc: func(origin string) bool {
            // JIKA MODE DEVELOPMENT, JANGAN KUNCI GERBANG , BIAR HP BISA MASUK LANCAR
            if os.Getenv("APP_ENV") == "development" {
                return true
            }

            allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
            if allowedOrigins == "" { return true } // Fallback development local
            for _, allowed := range strings.Split(allowedOrigins, ",") {
                if origin == strings.TrimSpace(allowed) { return true }
            }
            return false
        },
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders: []string{"Origin", "Content-Type", "Authorization", "X-Requested-With", "Accept"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

	r.Static("/uploads", "./uploads")
	r.Static("/public", "./public")

	// Payload Size Global Limiter (Max 5 MB upload protection)
	r.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 5*1024*1024)
		c.Next()
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "sukses", "message": "Halo Bos! Server Go Berhasil Menyala !"})
	})

	// ==========================================
	// -- RUTE API SAAS PUBLIC (GLOBAL & AUTH) --
	// ==========================================
	
	r.POST("/api/register", auth.Register)
	r.POST("/api/verify-otp", auth.VerifyOTP)
	r.POST("/api/login", auth.Login)
	r.POST("/api/auth/send-otp-wa", auth.SendOTPWhatsApp)
	r.POST("/api/reset-password", auth.ResetPassword)
	r.POST("/api/auth/check-account", auth.CheckAccount)

	// Instansiasi Modul Retail Layer
	retailRepo := retailRepository.NewRetailRepo(src.DB)
	retailHandler := retailDelivery.NewRetailHandler(retailRepo)
	r.POST("/api/retail/midtrans/webhook", retailHandler.MidtransWebhook)

	// ==========================================
	// -- RUTE TERPROTEKSI (MIDDLEWARE GATEWAY) --
	// ==========================================
	
	api := r.Group("/api")
	api.Use(middlewares.RequireAuth)
	{
		api.GET("/me", auth.GetMe)
		api.PUT("/profile", auth.UpdateProfile)
		api.PUT("/password", auth.UpdatePassword)
		api.POST("/auth/select-store", auth.SelectStore)
		api.POST("/setup", auth.SetupTokoBaru)

		// 🛒 Modul Bisnis: RETAIL MULTI-TENANT
        retailAPI := api.Group("/retail")
        {
            // Rute inventory bawaan
            retailDelivery.RegisterRetailInventoryRoutes(retailAPI, retailHandler)

            // 🛡️ SUNTIKAN SAKTI SINKRONISASI POS CHECKOUT:
            // Alamat penuh rute ini otomatis menjadi: POST /api/retail/pos/checkout!
            retailAPI.POST("/pos/checkout", retailHandler.CreateTransaction)
            
            // Sekalian tambahkan rute history dan closing laci yang ada di handler baru kita kemarin!
            retailAPI.GET("/pos/transactions", retailHandler.GetTransactions)
            retailAPI.GET("/pos/daily-closing", retailHandler.GetDailyClosing)
        }

		// 🧺 Modul Bisnis: LAUNDRY ECOSYSTEM
		laundryAPI := api.Group("/laundry")
		{
			laundryRepo := laundryRepository.NewLaundryRepo(src.DB)
			laundryHandler := &laundryDelivery.LaundryHandler{Repo: laundryRepo}
			laundryDelivery.RegisterLaundryRoutes(laundryAPI, laundryHandler)
		}

		// 🍔 Modul Bisnis: FOOD & BEVERAGES (FnB)
		fnbAPI := api.Group("/fnb") // 🚀 FIX: Redundansi middleware RequireAuth dicabut dari sini
		{
			fnbMenuRepo := fnbRepository.NewMenuRepo(src.DB)
			fnbMenuHandler := &fnbDelivery.MenuHandler{Repo: fnbMenuRepo}
			fnbOrderRepo := fnbRepository.NewOrderRepo(src.DB)
			fnbOrderHandler := fnbDelivery.NewOrderHandler(fnbOrderRepo)
			fnbDelivery.RegisterFnBRoutes(fnbAPI, fnbMenuHandler, fnbOrderHandler)
		}
	}

	port := os.Getenv("PORT")
	if port == "" { port = "8080" }

	log.Println("Server berjalan di port: " + port)
	r.Run("0.0.0.0:" + port)
}