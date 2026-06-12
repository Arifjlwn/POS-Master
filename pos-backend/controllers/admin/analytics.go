package admin

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AnalyticsController struct {
	DB *gorm.DB
}

func (a *AnalyticsController) GetSaaSTelemetry(c *gin.Context) {
	now := time.Now()
	awalBulan := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// =========================================================================
	// 1. ESTIMASI MRR & ARR BERDASARKAN TIER TOKO AKTIF
	// =========================================================================
	type PlanAgregasi struct {
		Industry         string `gorm:"column:industry"`
		SubscriptionPlan string `gorm:"column:subscription_plan"`
		Total            int64  `gorm:"column:total"`
	}

	var grupPlan []PlanAgregasi

	errGrup := a.DB.Table("stores").
		Select("LOWER(industry) as industry, LOWER(subscription_plan) as subscription_plan, COUNT(*) as total").
		Where("subscription_status = ?", "active").
		Group("LOWER(industry), LOWER(subscription_plan)").
		Scan(&grupPlan).Error

	if errGrup != nil {
		fmt.Println("LOG ERROR AGREGASI PLAN:", errGrup.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "gagal", "error": "Gagal agregasi modul bisnis"})
		return
	}

	var mrr int64
	var totalTrial, totalBasic, totalPro, totalPremium int64

	for _, grup := range grupPlan {
		ind := strings.TrimSpace(grup.Industry)
		plan := strings.TrimSpace(grup.SubscriptionPlan)
		jumlah := grup.Total

		switch plan {
		case "trial":
			totalTrial += jumlah
		case "basic":
			totalBasic += jumlah
		case "pro":
			totalPro += jumlah
		case "premium":
			totalPremium += jumlah
		}

		var hargaPaket int64 = 0
		switch ind {
		case "retail":
			switch plan {
			case "basic":
				hargaPaket = 49000
			case "pro":
				hargaPaket = 149000
			case "premium":
				hargaPaket = 299000
			}
		case "fnb":
			switch plan {
			case "basic":
				hargaPaket = 59000
			case "pro":
				hargaPaket = 169000
			case "premium":
				hargaPaket = 349000
			}
		case "jasa":
			switch plan {
			case "basic":
				hargaPaket = 49000
			case "pro":
				hargaPaket = 159000
			case "premium":
				hargaPaket = 329000
			}
		}
		mrr += (jumlah * hargaPaket)
	}

	arr := mrr * 12

	// =========================================================================
	// 2. GROWTH TENANT REAL MATRIX BULAN BERJALAN
	// =========================================================================
	var totalTenant, tenantBulanIni int64
	a.DB.Table("stores").Count(&totalTenant)
	a.DB.Table("stores").Where("created_at >= ?", awalBulan).Count(&tenantBulanIni)

	// =========================================================================
	// 3. CHURN RATE MONITORING
	// =========================================================================
	var churnedTenants int64
	a.DB.Table("stores").Where("subscription_status = ? OR subscription_end < ?", "suspended", now).Count(&churnedTenants)

	// =========================================================================
	// 4. TOP TENANT WHALE CATEGORIES (TRANSAKSI AUTO-LUNAS)
	// =========================================================================
	type TopTenantSummary struct {
		StoreID   uint    `json:"store_id"`
		NamaToko  string  `json:"nama_toko"`
		OwnerName string  `json:"owner_name"`
		PublicID  string  `json:"public_id"`
		TotalGmv  float64 `json:"total_gmv"`
	}
	var topTenants []TopTenantSummary

	errTop := a.DB.Table("transactions").
		Select("transactions.store_id, stores.nama_toko, stores.public_id as public_id, COALESCE(users.name, 'Owner') as owner_name, COALESCE(SUM(transactions.grand_total), 0) as total_gmv").
		Joins("JOIN stores ON stores.id = transactions.store_id").
		Joins("LEFT JOIN users ON users.id = stores.owner_id").
		Group("transactions.store_id, stores.nama_toko, stores.public_id, users.name").
		Order("total_gmv desc").
		Limit(5).
		Scan(&topTenants).Error

	if errTop != nil {
		fmt.Println("LOG INFO: Data transaksi masih kosong:", errTop.Error())
		topTenants = []TopTenantSummary{}
	}

	// =========================================================================
	// 5. TREND HISTORIS BULANAN 100% DINAMIS
	// =========================================================================
	labels := make([]string, 6)
	tenantGrowthData := make([]int, 6)
	transVolumeData := make([]int, 6)

	bulanIndo := []string{"Jan", "Feb", "Mar", "Apr", "Mei", "Jun", "Jul", "Ags", "Sep", "Okt", "Nov", "Des"}
	currentYear, currentMonth, _ := now.Date()

	for i := 5; i >= 0; i-- {
		target := time.Date(currentYear, currentMonth-time.Month(i), 1, 0, 0, 0, 0, now.Location())
		startOfMonth := target
		endOfMonth := target.AddDate(0, 1, 0).Add(-time.Nanosecond)

		labels[5-i] = fmt.Sprintf("%s %02d", bulanIndo[target.Month()-1], target.Year()%100)

		var tCount int64
		a.DB.Table("stores").Where("created_at >= ? AND created_at <= ?", startOfMonth, endOfMonth).Count(&tCount)
		tenantGrowthData[5-i] = int(tCount)

		var trCount int64
		a.DB.Table("transactions").Where("created_at >= ? AND created_at <= ?", startOfMonth, endOfMonth).Count(&trCount)
		transVolumeData[5-i] = int(trCount)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "sukses",
		"data": gin.H{
			"mrr":                   mrr,
			"arr":                   arr,
			"total_tenants":         totalTenant,
			"new_tenants":           tenantBulanIni,
			"churned_tenants":       churnedTenants,
			"top_tenants":           topTenants,
			"monthly_growth_labels": labels,
			"monthly_tenant_growth": tenantGrowthData,
			"monthly_trans_growth":  transVolumeData,
			"counts": gin.H{
				"trial":   totalTrial,
				"basic":   totalBasic,
				"pro":     totalPro,
				"premium": totalPremium,
			},
		},
	})
}