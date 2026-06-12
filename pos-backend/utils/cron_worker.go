package utils

import (
	"fmt"
	"log"
	"time"

	"pos-backend/models"
	src "pos-backend/src/core/config"
)

// InitSubscriptionCronJob menyalakan robot pengecekan masa aktif ruko di latar belakang
func InitSubscriptionCronJob() {
	log.Println("🤖 [CRON WORKER] Robot Algoritma Auto-Suspend diaktifkan!")

	// 1. Eksekusi sapu bersih sekali saat server Go baru aja di-restart/dinyalain
	executeAutoSuspend()

	// 2. Buat alarm Ticker. Di sini gua set jalan SETIAP 1 JAM.
	// (Kalau lu mau ngetes cepet, ganti aja jadi '1 * time.Minute' )
	ticker := time.NewTicker(1 * time.Hour)

	// 3. Jalankan di Goroutine agar tidak memblokir jalannya web server API Gin lu!
	go func() {
		for {
			<-ticker.C
			executeAutoSuspend()
		}
	}()
}

func executeAutoSuspend() {
	var expiredStores []models.Store
	now := time.Now()

	// 🔍 Misi: Cari semua toko yang statusnya 'active' TAPI masa aktifnya sudah kadaluwarsa (< waktu sekarang)
	// (Aman dari error NULL: Kalau toko trialnya unlimited / null, dia ga akan ikut ke-query)
	err := src.DB.Where("subscription_status = ? AND subscription_end < ?", "active", now).Find(&expiredStores).Error
	if err != nil {
		log.Println("❌ [CRON ERROR] Gagal menarik data toko kedaluwarsa:", err)
		return
	}

	// Kalau ga ada ruko yang expired, robot diam / tidur lagi
	if len(expiredStores) == 0 {
		return
	}

	log.Printf("🤖 [CRON WORKER] Menemukan %d toko kedaluwarsa! Mengeksekusi pembekuan masal...\n", len(expiredStores))

	for _, store := range expiredStores {
		// 1. Update status di database jadi 'suspended'
		errUpdate := src.DB.Model(&store).Update("subscription_status", "suspended").Error
		if errUpdate != nil {
			log.Printf("❌ [CRON ERROR] Gagal membekukan toko ID %d: %v\n", store.ID, errUpdate)
			continue
		}

		// 2. Tembak Log ke Mission Control pakai fungsi khusus Robot kita tadi!
		details := fmt.Sprintf("Masa aktif habis pada %s. Sistem membekukan ruko secara otomatis.", store.SubscriptionEnd.Format("02 Jan 2006 15:04 WIB"))
		RecordWorkerLog("Suspend Otomatis", store.PublicID, details)

		log.Printf("✅ [CRON WORKER] Toko %s (%s) resmi dibekukan.", store.NamaToko, store.PublicID)
	}
}
