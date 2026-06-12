package utils

import (
	"pos-backend/models"
	src "pos-backend/src/core/config"

	"github.com/gin-gonic/gin"
)

// RecordAdminAction otomatis menangkap IP & browser admin untuk dicatat ke database
func RecordAdminAction(c *gin.Context, action string, targetULID string, details string) error {
	userIDRaw, exists := c.Get("user_id")

	var userID uint = 0 // Default 0 kalau ternyata datanya nil (misal saat login)

	// 🚀 SAFE GUARDING: Cek dulu valus-nya ada apa kagak, jangan langsung dihantam float64!
	if exists && userIDRaw != nil {
		if val, ok := userIDRaw.(float64); ok {
			userID = uint(val)
		} else if valInt, ok := userIDRaw.(uint); ok {
			userID = valInt
		}
	}

	emailRaw, _ := c.Get("email")
	userEmail := "Arif Juliawan"
	if emailRaw != nil {
		userEmail = emailRaw.(string)
	}

	logData := models.AuditLog{
		UserID:         userID,
		UserEmail:      userEmail,
		Action:         action,
		TargetPublicID: targetULID,
		Details:        details,
		IPAddress:      c.ClientIP(),
		UserAgent:      c.Request.UserAgent(),
	}

	return src.DB.Create(&logData).Error
}

func RecordSystemLog(c *gin.Context, action string, targetULID string, details string) error {
	logData := models.AuditLog{
		UserID:         0,
		UserEmail:      "SYSTEM_AUTO",
		Action:         action,
		TargetPublicID: targetULID,
		Details:        details,
		IPAddress:      c.ClientIP(),
		UserAgent:      c.Request.UserAgent(),
	}

	return src.DB.Create(&logData).Error
}

// RecordWorkerLog mencatat aktivitas dari background process tanpa butuh context web (Gin)
func RecordWorkerLog(action string, targetULID string, details string) error {
	logData := models.AuditLog{
		UserID:         0,
		UserEmail:      "SYSTEM_CRON", // ◄ Menandakan ini dieksekusi oleh Robot Auto-Pilot!
		Action:         action,
		TargetPublicID: targetULID,
		Details:        details,
		IPAddress:      "127.0.0.1",        // IP Localhost server
		UserAgent:      "ARZURA_WORKER_V1", // Identitas Robot lu 
	}

	return src.DB.Create(&logData).Error
}
