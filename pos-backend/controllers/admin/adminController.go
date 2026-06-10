package admin

import (
	"net/http"
	"pos-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm" // 🚀 PAKE DRIVER GORM MURNI !
)

// Kita pasang field DB di dalam struct biar independen 
type AdminController struct {
	DB *gorm.DB
}

// GetTelemetryStats menarik ringkasan performa tenant ARZURA POS secara realtime
func (a *AdminController) GetTelemetryStats(c *gin.Context) {
	var totalTenants, activeTenants, pendingTenants, suspendedTenants int64

	// 🚀 SIKAT DATABASE LEWAT POINTER LOCAL STRUCT (Bukan dari config lagi, putus siklus !)
	a.DB.Model(&models.Store{}).Count(&totalTenants)
	a.DB.Model(&models.Store{}).Where("subscription_status = ?", "active").Count(&activeTenants)
	a.DB.Model(&models.Store{}).Where("subscription_status = ?", "pending").Count(&pendingTenants)
	a.DB.Model(&models.Store{}).Where("subscription_status = ?", "suspended").Count(&suspendedTenants)

	c.JSON(http.StatusOK, gin.H{
		"stats": gin.H{
			"total_tenants":     totalTenants,
			"active_tenants":    activeTenants,
			"pending_tenants":   pendingTenants,
			"suspended_tenants": suspendedTenants,
		},
	})
}