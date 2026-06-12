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

// GetSubscriptionOverview - UNIFIED ENDPOINT UNTUK MANAGEMENT HUB
func (s *SubscriptionController) GetSubscriptionOverview(c *gin.Context) {
	// 🚀 FLAT STRUCT SAKTI: Gabungan telemetri ruko, data owner hasil join, dan sisa hari 
	type StorePlanSummary struct {
		ID                 uint       `json:"id"`
		PublicID           string     `json:"public_id"`
		NamaToko           string     `json:"nama_toko"`
		Telepon            string     `json:"telepon"`
		BusinessType       string     `json:"business_type"`
		SubscriptionPlan   string     `json:"subscription_plan"`
		SubscriptionStatus string     `json:"subscription_status"`
		SubscriptionEnd    *time.Time `json:"subscription_end"`
		OwnerName          string     `json:"owner_name"`
		OwnerEmail         string     `json:"owner_email"`
		SisaHari           int        `json:"sisa_hari"` // ◄ DI-HITUNG LIVE DI LOOP 
	}

	// 1. Eksekusi Flat Scan Join dari tenantController dengan performa kilat ⚡
	var rawResults []StorePlanSummary
	err := s.DB.Table("stores").
		Select(`
			stores.id, 
			stores.public_id AS public_id, 
			stores.nama_toko, 
			stores.telepon, 
			stores.business_type, 
			stores.subscription_plan, 
			stores.subscription_status, 
			stores.subscription_end, 
			COALESCE(users.name, 'Owner Tidak Terdeteksi') as owner_name, 
			COALESCE(users.email, '-') as owner_email
		`).
		Joins("LEFT JOIN users ON users.id = stores.owner_id").
		Order("stores.created_at desc").
		Scan(&rawResults).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Sistem gagal mengintegrasikan data relasi owner!"})
		return
	}

	// 2. Akumulasi data statistik counter box atas (Murni 1 baris query group by sebenernya bisa, tapi ini biar tetep aman pake model lu)
	var totalTrial, totalBasic, totalPro, totalPremium int64
	s.DB.Model(&models.Store{}).Where("LOWER(subscription_plan) = ?", "trial").Count(&totalTrial)
	s.DB.Model(&models.Store{}).Where("LOWER(subscription_plan) = ?", "basic").Count(&totalBasic)
	s.DB.Model(&models.Store{}).Where("LOWER(subscription_plan) = ?", "pro").Count(&totalPro)
	s.DB.Model(&models.Store{}).Where("LOWER(subscription_plan) = ?", "premium").Count(&totalPremium)

	// 3. Hitung Sisa Hari secara forensik untuk masing-masing ruko
	now := time.Now()
	for i := range rawResults {
		if rawResults[i].SubscriptionEnd != nil {
			diff := rawResults[i].SubscriptionEnd.Sub(now)
			hari := int(diff.Hours() / 24)
			if hari < 0 {
				hari = 0 // Biar gak minus  kalau udah expired
			}
			rawResults[i].SisaHari = hari
		} else {
			rawResults[i].SisaHari = 0
		}
	}

	// 4. Lempar data terpadu ke frontend !
	c.JSON(http.StatusOK, gin.H{
		"status": "sukses",
		"counts": gin.H{
			"trial":   totalTrial,
			"basic":   totalBasic,
			"pro":     totalPro,
			"premium": totalPremium,
		},
		"stores": rawResults, // ◄ DATA INI UDAH KOMPLIT LENGKAP KAP KAP KAP!
	})
}