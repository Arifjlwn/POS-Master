package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"pos-backend/models"
	src "pos-backend/src/core/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	// Cek apakah ada karcis di kantong (Header Auth)
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak ! Anda belum login."})
		c.Abort() //Langsung Usir, Jangan kasih masuk
		return
	}

	// Format Karcis biasanya diawali "Bearer ", jadi potong sisa kode token aja
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	// Cek Keaslian Karcis pakai kunci rahasia saat login
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metode token tidak valid")
		}
		return []byte("KUNCI_RAHASIA_SUPER_KUAT_123"), nil // harus Sama Persis dengan di authController
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau sudah hangus !"})
		c.Abort()
		return
	}

	// Kalau karcis asli, bongkar isi nya (Ambil ID User, ID TOKO, dan Role)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Simpan data ini kedalam "Context" agar bisa dibaca oleh controller nanti
		c.Set("user_id", claims["user_id"])
		c.Set("store_id", claims["store_id"])
		c.Set("role", claims["role"])

		c.Next() // Silahkan Masuk !
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Gagal membaca data dari token"})
		c.Abort()
		return
	}
}

// --- SATPAM RUANG VIP (KHUSUS OWNER) ---
func RequireOwner(c *gin.Context) {
	// Ambil data role yang udah ditaruh di kantong (Context) sama satpam RequireAuth
	roleRaw, exists := c.Get("role")
	
	// Kalau datanya ga ada, atau rolenya BUKAN owner, langsung tendang!
	if !exists || roleRaw.(string) != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses Ditolak! Fitur ini khusus untuk Bos (Owner) 😎"})
		c.Abort() // Usir! Jangan kasih lanjut ke controller
		return
	}

	// Kalau dia beneran owner, silakan lewat Bosku!
	c.Next()
}

// ========================================================
// 🚀 SATPAM LAPIS 3: SAAS PLAN GATING (SISTEM KASTA LEVEL)
// ========================================================
func RequireSaaSLevel(minLevel int) gin.HandlerFunc {
	return func(c *gin.Context) {
		storeIDRaw, exists := c.Get("store_id")
		if !exists || storeIDRaw == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak! Data toko tidak ditemukan di sesi Anda."})
			c.Abort()
			return
		}

		var storeID uint
		if val, ok := storeIDRaw.(float64); ok {
			storeID = uint(val)
		} else if val, ok := storeIDRaw.(uint); ok {
			storeID = val
		}

		// 1. Tarik Data Kasta Langsung dari Database (Real-time Check)
		var store models.Store
		if err := src.DB.Select("subscription_plan", "subscription_status").First(&store, storeID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memverifikasi status langganan toko."})
			c.Abort()
			return
		}

		// 2. Cek Apakah Langganan Mati / Nunggak Bayar
		if store.SubscriptionStatus != "active" {
			c.JSON(http.StatusPaymentRequired, gin.H{"error": "Akses dihentikan! Masa berlangganan Anda sudah berakhir. Silakan perpanjang."})
			c.Abort()
			return
		}

		// 3. Konversi Nama Paket ke Level Angka
		currentLevel := 1 // Kasta Sudra (Basic)
		plan := strings.ToLower(store.SubscriptionPlan)
		
		if plan == "premium" || plan == "enterprise" || plan == "trial" {
			currentLevel = 3 // Kasta Dewa
		} else if plan == "pro" {
			currentLevel = 2 // Kasta Kesatria
		}

		// 4. Proses Eksekusi Gembok
		if currentLevel < minLevel {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Akses API Ditolak 🛑! Fitur ini membutuhkan upgrade paket langganan.",
			})
			c.Abort() // Tendang balik ke laut!
			return
		}

		// Kasta mencukupi, silakan masuk ke Controller!
		c.Next()
	}
}