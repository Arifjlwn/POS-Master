package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"pos-backend/models" // Sesuaikan dengan nama modul go.mod Anda
	src "pos-backend/src/core/config"
)

// Catatan memori untuk menyimpan jejak digital laju permintaan berdasarkan IP
type clientLimiter struct {
	lastSeen time.Time
	count    int
}

var (
	mu      sync.Mutex
	clients = make(map[string]*clientLimiter)
)

// AdminLoginRateLimiter membatasi percobaan login maksimal 5 kali per menit demi mengamankan peladen produksi
func AdminLoginRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		mu.Lock()
		defer mu.Unlock()

		ip := c.ClientIP()
		now := time.Now()

		if client, exists := clients[ip]; exists {
			// Jika kunjungan terakhir sudah lebih dari 1 menit, reset hitungan dari awal
			if now.Sub(client.lastSeen) > 1*time.Minute {
				client.count = 1
				client.lastSeen = now
			} else {
				client.count++
				// Blokir permintaan jika melebihi batas ambang keamanan
				if client.count > 5 {
					c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
						"error": "Terlalu banyak percobaan login. Akses ditangguhkan selama 1 menit.",
					})
					return
				}
			}
		} else {
			clients[ip] = &clientLimiter{lastSeen: now, count: 1}
		}

		c.Next()
	}
}

// RequireSuperAdmin memastikan token yang masuk memegang role super_admin secara sah di DB
func RequireSuperAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Cek role dari context JWT
		role, exists := c.Get("role")
		if !exists || role.(string) != "super_admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "AKSES ILEGAL: Otoritas Root Admin Tidak Valid!"})
			c.Abort()
			return
		}

		// 2. Double-check anti-token hijacking langsung ke DB core
		userIDRaw, _ := c.Get("user_id")
		userID := uint(userIDRaw.(float64))
		
		var dbUser models.User
		if err := src.DB.Select("role").First(&dbUser, userID).Error; err != nil || dbUser.Role != "super_admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Hak akses dicabut oleh server pusat!"})
			c.Abort()
			return
		}

		c.Next()
	}
}