package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"pos-backend/controllers/admin"
	"pos-backend/controllers/auth"
	"pos-backend/middlewares"
	src "pos-backend/src/core/config"
	"pos-backend/utils"

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
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan environment default system")
	}

	// 1. KONEKSI & INSTANSIASI LAYER (Dependency Injection Terpusat)
	src.ConnectDatabase()
	utils.SeedSuperAdmin(src.DB)
	utils.InitSubscriptionCronJob()

	// Inisialisasi Modul Retail
	retailRepo := retailRepository.NewRetailRepo(src.DB)
	retailHandler := retailDelivery.NewRetailHandler(retailRepo)

	// Inisialisasi Modul Laundry
	laundryRepo := laundryRepository.NewLaundryRepo(src.DB)

	// Inisialisasi Modul Food & Beverages (FnB)
	fnbMenuRepo := fnbRepository.NewMenuRepo(src.DB)
	fnbMenuHandler := &fnbDelivery.MenuHandler{Repo: fnbMenuRepo}
	fnbOrderRepo := fnbRepository.NewOrderRepo(src.DB)
	fnbOrderHandler := fnbDelivery.NewOrderHandler(fnbOrderRepo)

	// 2. KONFIGURASI ENGINE WEB SERVER
	r := gin.Default()

	// CORS Tight Guard dinamis berbasis ENV untuk mengunci eksploitasi di Production
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			if os.Getenv("APP_ENV") == "development" {
				return true
			}

			allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
			if allowedOrigins == "" {
				log.Println("CRITICAL ERROR: ALLOWED_ORIGINS tidak terdeteksi di lingkungan Production!")
				return false
			}

			for _, allowed := range strings.Split(allowedOrigins, ",") {
				if origin == strings.TrimSpace(allowed) {
					return true
				}
			}
			return false
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With", "Accept", "ngrok-skip-browser-warning"},
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
	r.POST("/api/auth/send-otp-email", auth.SendOTPEmailEndpoint)
	r.POST("/api/reset-password", auth.ResetPassword)
	r.POST("/api/auth/check-account", auth.CheckAccount)
	r.POST("/api/re-trigger-payment", auth.ReTriggerPaymentHandler)

	// Webhook Gateway Midtrans
	r.POST("/api/retail/midtrans/webhook", retailHandler.MidtransWebhook)

	// ==========================================
	// --       🚀 GERBANG UTAMA ADMIN         --
	// ==========================================
	// 🚀 DI DALAM FUNCTION MAIN() LU BRAY, PAS BAGIAN INSTANSIASI CONTROLLER ADMIN:
		adminAuthCtrl := &admin.AuthController{DB: src.DB}
		adminDashboardCtrl := &admin.DashboardController{DB: src.DB}
		analyticsCtrl := &admin.AnalyticsController{DB: src.DB}
		adminTenantCtrl := &admin.TenantController{DB: src.DB}
		adminAuditCtrl := &admin.AuditController{DB: src.DB}
		adminSubCtrl := &admin.SubscriptionController{DB: src.DB}

		// Bagian Group Routing Admin Lu Tinggal Colok Sesuai Modular-nya bray:
		adminRoutes := r.Group("/api/admin")
		{
			adminRoutes.POST("/login", adminAuthCtrl.AdminLogin)
			adminRoutes.GET("/dashboard-stats", adminDashboardCtrl.GetTelemetryStats)
			adminRoutes.GET("/analytics/telemetry", analyticsCtrl.GetSaaSTelemetry)
			adminRoutes.GET("/stores", adminTenantCtrl.GetAllTenants)
			adminRoutes.PUT("/stores/:id/suspend", adminTenantCtrl.SuspendTenant)
			adminRoutes.PUT("/stores/:id/activate", adminTenantCtrl.ActivateTenant)
			adminRoutes.GET("/audit-logs", adminAuditCtrl.GetDetailedAuditLogs)
			adminRoutes.GET("/subscription-overview", adminSubCtrl.GetSubscriptionOverview)
		}

	// ==========================================
	// -- RUTE TERPROTEKSI MERCHANT (STORE-LEVEL)--
	// ==========================================
	apiGroup := r.Group("/api")
	apiGroup.Use(middlewares.RequireAuth)
	{
		apiGroup.GET("/me", auth.GetMe)
		apiGroup.PUT("/profile", auth.UpdateProfile)
		apiGroup.PUT("/password", auth.UpdatePassword)
		apiGroup.POST("/auth/select-store", auth.SelectStore)
		apiGroup.POST("/setup", auth.SetupTokoBaru)

		// 🛒 Modul Bisnis: RETAIL MULTI-TENANT
		retailAPI := apiGroup.Group("/retail")
		{
			retailDelivery.RegisterRetailInventoryRoutes(retailAPI, retailHandler)
			retailAPI.POST("/pos/checkout", retailHandler.CreateTransaction)
			retailAPI.GET("/pos/transactions", retailHandler.GetTransactions)
			retailAPI.GET("/pos/daily-closing", retailHandler.GetDailyClosing)
		}

		// 🧺 Modul Bisnis: LAUNDRY ECOSYSTEM
		laundryAPI := apiGroup.Group("/laundry")
		{
			laundryDelivery.RegisterLaundryRoutes(laundryAPI, laundryRepo)
		}

		// 🍔 Modul Bisnis: FOOD & BEVERAGES (FnB)
		fnbAPI := apiGroup.Group("/fnb")
		{
			fnbDelivery.RegisterFnBRoutes(fnbAPI, fnbMenuHandler, fnbOrderHandler)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server berjalan di port: " + port)
	r.Run("0.0.0.0:" + port)
}
