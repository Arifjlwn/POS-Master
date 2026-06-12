package admin

import (
	"net/http"
	"time"

	"pos-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SubscriptionController struct {
	DB *gorm.DB
}

// GetSubscriptionOverview handles MONITORING MASTER PAKET LANGGANAN TENANT
func (s *SubscriptionController) GetSubscriptionOverview(c *gin.Context) {
	type StorePlanSummary struct {
		ID                 uint       `json:"id"`
		NamaToko           string     `json:"nama_toko"`
		OwnerName          string     `json:"owner_name"`
		SubscriptionPlan   string     `json:"subscription_plan"`
		SubscriptionStatus string     `json:"subscription_status"`
		SubscriptionEnd    *time.Time `json:"subscription_end"`
		SisaHari           int        `json:"sisa_hari"`
	}

	var rawStores []models.Store
	if err := s.DB.Order("subscription_end asc").Find(&rawStores).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memetakan telemetri billing"})
		return
	}

	var summaries []StorePlanSummary
	var totalTrial, totalBasic, totalPro, totalPremium int64

	s.DB.Model(&models.Store{}).Where("subscription_plan = ?", "trial").Count(&totalTrial)
	s.DB.Model(&models.Store{}).Where("subscription_plan = ?", "basic").Count(&totalBasic)
	s.DB.Model(&models.Store{}).Where("subscription_plan = ?", "pro").Count(&totalPro)
	s.DB.Model(&models.Store{}).Where("subscription_plan = ?", "premium").Count(&totalPremium)

	now := time.Now()
	for _, sStore := range rawStores {
		var ownerName string
		var user models.User
		if err := s.DB.Select("name").Where("id = ?", sStore.OwnerID).First(&user).Error; err == nil {
			ownerName = user.Name
		} else {
			ownerName = "Tidak Terdeteksi"
		}

		sisaHari := 0
		if sStore.SubscriptionEnd != nil {
			diff := sStore.SubscriptionEnd.Sub(now)
			sisaHari = int(diff.Hours() / 24)
		}

		summaries = append(summaries, StorePlanSummary{
			ID:                 sStore.ID,
			NamaToko:           sStore.NamaToko,
			OwnerName:          ownerName,
			SubscriptionPlan:   sStore.SubscriptionPlan,
			SubscriptionStatus: sStore.SubscriptionStatus,
			SubscriptionEnd:    sStore.SubscriptionEnd,
			SisaHari:           sisaHari,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "sukses",
		"counts": gin.H{
			"trial":   totalTrial,
			"basic":   totalBasic,
			"pro":     totalPro,
			"premium": totalPremium,
		},
		"stores": summaries,
	})
}