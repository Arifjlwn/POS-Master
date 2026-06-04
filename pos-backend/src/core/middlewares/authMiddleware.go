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

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		isSelectToken := false
		if isSelectRaw, exists := claims["is_select"]; exists {
			isSelectToken = isSelectRaw.(bool)
		}

		if isSelectToken {
			currentPath := c.Request.URL.Path
			isAllowed := currentPath == "/api/auth/select-store" || currentPath == "/api/setup"

			if !isAllowed {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "GAGAL SETUP TOKO ANDA HARUS MEMILIH CABANG DULU",
					"code":  "REQUIRE_STORE_SELECTION",
				})
				c.Abort()
				return
			}
		}

		c.Set("user_id", claims["user_id"])
		c.Set("public_id", claims["public_id"]) // 🚀 AMAN: inject ke konteks untuk tracking folder Supabase
		c.Set("role", claims["role"])

		if !isSelectToken {
			if storeID, exists := claims["store_id"]; exists { c.Set("store_id", storeID) }
			if planType, exists := claims["plan_type"]; exists { c.Set("plan_type", planType) }
		}

		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Gagal membaca data dari token"})
		c.Abort()
		return
	}
}

func RequireOwner(c *gin.Context) {
	roleRaw, exists := c.Get("role")
	if !exists || roleRaw.(string) != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses Ditolak! Fitur ini khusus untuk Bos (Owner)"})
		c.Abort()
		return
	}
	c.Next()
}

func RequireSaaSLevel(minLevel int) gin.HandlerFunc {
	return func(c *gin.Context) {
		storeIDRaw, exists := c.Get("store_id")
		if !exists || storeIDRaw == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak! Data toko tidak ditemukan."})
			c.Abort()
			return
		}

		var storeID uint
		if val, ok := storeIDRaw.(float64); ok { storeID = uint(val) } else if val, ok := storeIDRaw.(uint); ok { storeID = val }

		var store models.Store
		if err := src.DB.Select("id", "subscription_plan", "subscription_status", "subscription_end").First(&store, storeID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memverifikasi status langganan toko."})
			c.Abort()
			return
		}

		if store.SubscriptionEnd != nil && time.Now().After(*store.SubscriptionEnd) {
			if store.SubscriptionStatus == "active" { src.DB.Model(&store).Update("subscription_status", "inactive") }
			c.JSON(http.StatusPaymentRequired, gin.H{"error": "Masa aktif toko telah habis. Silakan perpanjang langganan.", "code": "SUBSCRIPTION_EXPIRED"})
			c.Abort()
			return
		}

		if store.SubscriptionStatus != "active" {
			c.JSON(http.StatusPaymentRequired, gin.H{"error": "Akses dihentikan! Status toko Anda tidak aktif."})
			c.Abort()
			return
		}

		currentLevel := 1
		plan := strings.ToLower(store.SubscriptionPlan)
		if plan == "premium" || plan == "trial" { currentLevel = 3 } else if plan == "pro" { currentLevel = 2 }

		if currentLevel < minLevel {
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses API Ditolak! Fitur ini membutuhkan upgrade paket langganan.", "code": "UPGRADE_REQUIRED"})
			c.Abort()
			return
		}

		c.Next()
	}
}