package admin

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"pos-backend/models"
	"pos-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DashboardController struct {
	DB *gorm.DB
}

func (d *DashboardController) GetTelemetryStats(c *gin.Context) {
	// --- PART 1: TENANT SUBSCRIPTION STATS ---
	var activeStores, trialStores, suspendedStores, archivedStores int64
	d.DB.Model(&models.Store{}).
		Where("subscription_status = ? AND LOWER(subscription_plan) != ?", "active", "trial").
		Count(&activeStores)
	d.DB.Model(&models.Store{}).
		Where("LOWER(subscription_plan) = ?", "trial").
		Count(&trialStores)
	d.DB.Model(&models.Store{}).Where("subscription_status = ?", "suspended").Count(&suspendedStores)
	d.DB.Model(&models.Store{}).Where("subscription_status = ?", "archived").Count(&archivedStores)

	// --- PART 2: SERVER HEALTH MONITORING ---
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	ramUsageMB := m.Alloc / 1024 / 1024
	cpuLoad := runtime.NumGoroutine()

	startPing := time.Now()
	sqlDB, err := d.DB.DB()
	dbStatus := "Mati / Putus"
	latency := int64(0)
	if err == nil && sqlDB.Ping() == nil {
		dbStatus = "Online & Stabil"
		latency = time.Since(startPing).Milliseconds()
	}

	// --- PART 3: LIVE ENGINE MATRIX (MODE ALL TIME) ---
	var totalSalesAllTime float64
	var activeCashiersCount int64

	errTrans := d.DB.Model(&models.Transaction{}).
		Select("COALESCE(SUM(total_harga), 0)").
		Row().Scan(&totalSalesAllTime)

	if errTrans != nil {
		fmt.Println("❌ [DEBUG ERROR] Gagal Hitung Total Harga Transaksi:", errTrans)
	}

	errSession := d.DB.Model(&models.CashierSession{}).Where("LOWER(status) = ?", "open").Count(&activeCashiersCount)
	if errSession != nil {
		fmt.Println("❌ [DEBUG ERROR] Gagal Hitung Sesi Kasir:", errSession)
	}

	totalTransactions := fmt.Sprintf("Rp %s", utils.FormatRupiah(int64(totalSalesAllTime)))

	// --- PART 4: 🗺️ LIVE TENANT MAP NODES (MURNI DARI DATABASE) ---
	type LiveMapNode struct {
		ID            uint    `json:"id"`
		NamaToko      string  `json:"nama_toko"`
		OwnerName     string  `json:"owner_name"`
		Plan          string  `json:"plan"`
		Status        string  `json:"status"`
		Latitude      float64 `json:"latitude"`
		Longitude     float64 `json:"longitude"`
		ActiveCashier int64   `json:"active_cashier"`
		OmzetToday    string  `json:"omzet_today"`
		ColorCode     string  `json:"color_code"`
	}

	var mapNodes []LiveMapNode

	type StoreData struct {
		ID                 uint
		NamaToko           string
		SubscriptionPlan   string
		SubscriptionStatus string
		OwnerName          string
		Latitude           float64
		Longitude          float64
	}

	var storesData []StoreData
	d.DB.Table("stores").
		Select("stores.id, stores.nama_toko, stores.subscription_plan, stores.subscription_status, stores.latitude, stores.longitude, COALESCE(users.name, 'Owner Unknown') as owner_name").
		Joins("LEFT JOIN users ON users.id = stores.owner_id").
		Scan(&storesData)

	todayStart := time.Now().Truncate(24 * time.Hour)

	for _, s := range storesData {
		// 🚀 FILTER SAKTI ANTI-HARDCODE: Cuma toko yang punya koordinat valid yang masuk Radar!
		if s.Latitude != 0 && s.Longitude != 0 {

			// 🎨 TENTUKAN WARNA RADAR BERDASARKAN STATUS
			color := "#64748b" // Abu-abu (Offline/Archived)
			if s.SubscriptionStatus == "active" {
				if strings.ToLower(s.SubscriptionPlan) == "trial" {
					color = "#f59e0b" // Kuning (Trial)
				} else {
					color = "#10b981" // Hijau (Paid/Aman)
				}
			} else if s.SubscriptionStatus == "suspended" {
				color = "#ef4444" // Merah (Dibekukan)
			}

			// 💰 HITUNG OMZET HARI INI
			var omzet float64
			d.DB.Table("transactions").
				Where("store_id = ? AND created_at >= ?", s.ID, todayStart).
				Select("COALESCE(SUM(total_harga), 0)").
				Row().Scan(&omzet)

			// 🧑‍💻 HITUNG KASIR STANDBY
			var kasirAktif int64
			d.DB.Table("cashier_sessions").
				Where("store_id = ? AND LOWER(status) = ?", s.ID, "open").
				Count(&kasirAktif)

			// 📍 MASUKIN DATA MURNI 100% DARI DB (Ga ada tebak-tebakan lagi bray)
			mapNodes = append(mapNodes, LiveMapNode{
				ID:            s.ID,
				NamaToko:      s.NamaToko,
				OwnerName:     s.OwnerName,
				Plan:          strings.ToUpper(s.SubscriptionPlan),
				Status:        strings.ToUpper(s.SubscriptionStatus),
				Latitude:      s.Latitude,
				Longitude:     s.Longitude,
				ActiveCashier: kasirAktif,
				OmzetToday:    fmt.Sprintf("Rp %s", utils.FormatRupiah(int64(omzet))),
				ColorCode:     color,
			})
		}
	}

	// --- PART 5: JSON RETURN PUSAT ---
	c.JSON(http.StatusOK, gin.H{
		"status": "sukses",
		"data": gin.H{
			"active_stores":    activeStores,
			"trial_stores":     trialStores,
			"suspended_stores": suspendedStores,
			"archived_stores":  archivedStores,
		},
		"live_stats": gin.H{
			"users_online":             1,
			"total_transactions_today": totalTransactions,
			"active_cashiers":          activeCashiersCount,
			"open_shifts":              0,
		},
		"server_health": gin.H{
			"cpu_usage": cpuLoad,
			"ram_usage": ramUsageMB,
			"db_status": dbStatus,
			"latency":   fmt.Sprintf("%dms", latency),
		},
		"live_map_nodes": mapNodes, // ◄ DATA PETA DIKIRIM KE VUE BRAY!
	})
}
