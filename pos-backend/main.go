package main

import (
	"log"
	"net/http"
	"os"
	"pos-backend/config"
	"pos-backend/controllers/auth"   // 🚀 IMPORT FOLDER AUTH
	"pos-backend/controllers/retail" // 🚀 IMPORT FOLDER RETAIL
	"pos-backend/controllers/laundry" // 🚀 IMPORT FOLDER RETAIL
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
		
		// --- UPDATE INTEGRASI RUTE SETUP TOKO (Di dalam main.go) ---
		api.POST("/setup", func(c *gin.Context) {
			userIDRaw, exists := c.Get("user_id")
			if !exists {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak valid"})
				return
			}
			userID := uint(userIDRaw.(float64))

			// Struct penangkap disesuaikan 100% dengan state Vue kamu
			var input struct {
				NamaToko     string   `json:"nama_toko" binding:"required"`
				TipeBisnis   string   `json:"tipe_bisnis" binding:"required"` // Gabungan kategori & spesifikasi
				AlamatJalan  string   `json:"alamat_toko"`                    // Menerima form.alamat_jalan
				Provinsi     string   `json:"provinsi"`
				Kota         string   `json:"kota"`
				Kecamatan    string   `json:"kecamatan"`
				Kelurahan    string   `json:"kelurahan"`
				KodePos      string   `json:"kode_pos"`
				Telepon      string   `json:"telepon"`
				FiturOpsional []string `json:"fitur_aktif"`                    // Menerima array form.fitur_opsional
			}
			
			if err := c.ShouldBindJSON(&input); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Format data salah, pastikan form terisi lengkap"})
				return
			}

			// Merangkai array modul premium menjadi string teks dipisah koma
			var fiturString string
			for i, fitur := range input.FiturOpsional {
				if i > 0 {
					fiturString += ","
				}
				fiturString += fitur
			}

			// Menjahit teks alamat lengkap sekali jadi untuk kebutuhan invoice cetak nota cepat
			alamatLengkap := input.AlamatJalan + ", Kel. " + input.Kelurahan + ", Kec. " + input.Kecamatan + ", " + input.Kota + ", " + input.Provinsi + " " + input.KodePos

			newStore := models.Store{
				NamaToko:     input.NamaToko,
				BusinessType: input.TipeBisnis,
				Telepon:      input.Telepon,
				FiturAktif:   fiturString,
				Alamat:       alamatLengkap, // Masuk alamat gabungan
				Provinsi:     input.Provinsi, // Masuk pecahan data area operasional
				Kota:         input.Kota,
				Kecamatan:    input.Kecamatan,
				Kelurahan:    input.Kelurahan,
				KodePos:      input.KodePos,
			}
			
			if err := config.DB.Create(&newStore).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan infrastruktur toko baru"})
				return
			}

			if err := config.DB.Model(&models.User{}).Where("id = ?", userID).Update("store_id", newStore.ID).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengaitkan otoritas toko ke akun owner"})
				return
			}

			// Generate ulang JWT token bawaan kamu
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
				"message":  "Konfigurasi sistem siap! Selamat datang di platform POS SaaS.",
				"store_id": newStore.ID,
				"token":    tokenString,
				"data": gin.H{
					"nama_toko":   newStore.NamaToko,
					"tipe_bisnis": newStore.BusinessType, 
				},
			})
		})

		// ==========================================
		// 🛒 RUTE KHUSUS RETAIL (URL nambah /retail)
		// ==========================================
		retailAPI := api.Group("/retail")
		{
			// -- Rute Produk (CRUD) --
			retailAPI.POST("/products", retail.CreateProduct)
			retailAPI.GET("/products", retail.GetProducts)
			retailAPI.PUT("/products/:id", retail.UpdateProduct)
			retailAPI.DELETE("/products/:id", retail.DeleteProduct)
			retailAPI.GET("/products/export", retail.ExportProducts)
			retailAPI.POST("/products/import", retail.ImportProducts)
			retailAPI.GET("/categories", retail.GetCategories)

			// -- Rute Karyawan --
			retailAPI.POST("/employees", retail.CreateEmployee)
			retailAPI.GET("/employees", retail.GetEmployees)
			retailAPI.PUT("/employees/:id", retail.UpdateEmployee)

			retailAPI.POST("/schedules/bulk", retail.SaveSchedules)
			retailAPI.GET("/schedules", retail.GetSchedules)

			// -- Rute Absensi --
			retailAPI.POST("/attendance", retail.StoreAttendance)
			retailAPI.GET("/attendance", retail.GetAttendance)
			retailAPI.GET("/attendance/export", retail.ExportAttendance)

			// --- RUTE RETUR, LPB, STOK OPNAME ---
			retailAPI.POST("/returns", retail.CreateReturn)
			retailAPI.GET("/returns", retail.GetReturns)
			retailAPI.POST("/purchases", retail.CreateLPB)
			retailAPI.POST("/stock-opname", retail.CreateStockOpname)
			retailAPI.GET("/stock-opname/history", retail.GetStockOpnameHistory)

			// --- RUTE KASIR (POS) ---
			retailAPI.POST("/pos/open-session", retail.OpenSession)
			retailAPI.GET("/pos/check-session", retail.CheckSessionStatus)
			retailAPI.POST("/pos/close-session/:id", retail.CloseSession)
			retailAPI.POST("/checkout", retail.CreateTransaction)
			retailAPI.GET("/transactions", retail.GetTransactions)

			// Rute Laporan
			retailAPI.GET("/report/dashboard", retail.GetDashboardReport)
		}

		// ==========================================
		// 🧺 RUTE KHUSUS LAUNDRY 
		// ==========================================
		laundryAPI := api.Group("/laundry")
		{
			laundryAPI.GET("/services", laundry.AmbilDaftarLayananLaundry) 
			laundryAPI.POST("/services", laundry.TambahLayananLaundry)
			laundryAPI.PUT("/services/:id", laundry.EditLayananLaundry)
			laundryAPI.DELETE("/services/:id", laundry.HapusLayananLaundry)
			laundryAPI.POST("/checkout", laundry.ProsesCheckoutLaundry)
			laundryAPI.GET("/tracking", laundry.AmbilDataTracking) 
			laundryAPI.PUT("/tracking/:id/status", laundry.UpdateStatusCucian)
			laundryAPI.GET("/customers/search", laundry.CariPelanggan)
			laundryAPI.GET("/report", laundry.AmbilLaporan)
			laundryAPI.GET("/setting", laundry.GetSettingToko)
			laundryAPI.PUT("/setting", laundry.UpdateSettingToko)
			laundryAPI.GET("/kasir", middlewares.RequireAuth, middlewares.RequireOwner, laundry.GetKasirList)
			laundryAPI.POST("/kasir", middlewares.RequireAuth, middlewares.RequireOwner, laundry.CreateKasir)
			laundryAPI.DELETE("/kasir/:id", middlewares.RequireAuth, middlewares.RequireOwner, laundry.DeleteKasir)
			laundryAPI.GET("/perfumes", middlewares.RequireAuth, laundry.GetPerfumes)
			laundryAPI.POST("/perfumes", middlewares.RequireAuth, middlewares.RequireOwner, laundry.CreatePerfume)
			laundryAPI.DELETE("/perfumes/:id", middlewares.RequireAuth, middlewares.RequireOwner, laundry.DeletePerfume)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server berjalan di port: " + port)
	r.Run("0.0.0.0:" + port)
}