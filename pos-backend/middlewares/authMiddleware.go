package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"pos-backend/models"
	src "pos-backend/src/core/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak! Anda belum login."})
		c.Abort()
		return
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Sistem keamanan server mengalami gangguan teknis internal."})
		c.Abort()
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metode token tidak valid")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau sudah hangus!"})
		c.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Gagal membaca data dari token"})
		c.Abort()
		return
	}

	isSelectToken := false
	if isSelectRaw, exists := claims["is_select"]; exists {
		isSelectToken = isSelectRaw.(bool)
	}

	if isSelectToken {
		currentPath := c.Request.URL.Path

		isAllowed := currentPath == "/api/auth/select-store" ||
			currentPath == "/api/setup" ||
			currentPath == "/api/retail/subscription/upgrade" ||
			currentPath == "/api/me"

		if !isAllowed {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Akses dibatasi! Selesaikan konfigurasi infrastruktur gerai atau tentukan cabang operasional Anda terlebih dahulu.",
				"code":  "REQUIRE_STORE_SELECTION",
			})
			c.Abort()
			return
		}
	}

	c.Set("user_id", claims["user_id"])
	c.Set("public_id", claims["public_id"])
	c.Set("role", claims["role"])

	if storeID, exists := claims["store_id"]; exists && storeID != nil {
		c.Set("store_id", storeID)
	} else {
		c.Set("store_id", float64(0))
	}

	if planType, exists := claims["plan_type"]; exists && planType != nil {
		c.Set("plan_type", planType)
	} else {
		c.Set("plan_type", "")
	}

	c.Next()
}

func RequireOwner(c *gin.Context) {
	roleRaw, exists := c.Get("role")
	if !exists || roleRaw.(string) != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses Ditolak! Fitur ini khusus untuk tingkat Pemilik (Owner)"})
		c.Abort()
		return
	}
	c.Next()
}

func RequireSaaSLevel(minLevel int) gin.HandlerFunc {
	return func(c *gin.Context) {
		storeIDRaw, exists := c.Get("store_id")
		if !exists || storeIDRaw == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak! Klaim data toko tidak ditemukan di dalam session token."})
			c.Abort()
			return
		}

		var storeID uint
		if val, ok := storeIDRaw.(float64); ok {
			storeID = uint(val)
		} else if val, ok := storeIDRaw.(uint); ok {
			storeID = val
		}

		if storeID == 0 {
			c.Next()
			return
		}

		var store models.Store
		if err := src.DB.Select("id", "subscription_plan", "subscription_status", "subscription_end").First(&store, storeID).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error":     fmt.Sprintf("Gagal memverifikasi lisensi cabang. Detail database: %v", err.Error()),
				"code":      "STORE_RECORD_NOT_FOUND",
				"target_id": storeID,
			})
			c.Abort()
			return
		}

		if store.SubscriptionEnd != nil && time.Now().After(*store.SubscriptionEnd) {
			if store.SubscriptionStatus == "active" {
				src.DB.Model(&store).Update("subscription_status", "inactive")
			}
			c.JSON(http.StatusPaymentRequired, gin.H{"error": "Masa aktif toko telah habis. Silakan perpanjang langganan.", "code": "SUBSCRIPTION_EXPIRED"})
			c.Abort()
			return
		}

		if store.SubscriptionStatus != "active" {
			c.JSON(http.StatusPaymentRequired, gin.H{"error": "Akses dihentikan! Status lisensi toko Anda tidak aktif di server."})
			c.Abort()
			return
		}

		// 🚀 BYPASS DEWA: Jika plan berstatus premium atau trial, berikan hak akses mutlak tanpa batas level !
		plan := strings.ToLower(strings.TrimSpace(store.SubscriptionPlan))
		if plan == "premium" || plan == "trial" || strings.Contains(plan, "premium") {
			c.Next()
			return
		}

		currentLevel := 1
		if plan == "pro" || strings.Contains(plan, "pro") {
			currentLevel = 2
		} else if plan == "basic" || strings.Contains(plan, "basic") {
			currentLevel = 1
		}

		if currentLevel < minLevel {
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses API Ditolak! Fitur ini membutuhkan upgrade paket langganan.", "code": "UPGRADE_REQUIRED"})
			c.Abort()
			return
		}

		c.Next()
	}
}
