package controllers

import (
	"net/http"
	"pos-backend/config"
	"pos-backend/models"

	"github.com/gin-gonic/gin"
)

// Input struct untuk submit jadwal per baris
type ScheduleItem struct {
	UserID    uint   `json:"user_id" binding:"required"`
	Tanggal   string `json:"tanggal" binding:"required"`    // Format: YYYY-MM-DD
	ShiftType string `json:"shift_type" binding:"required"` // 'Shift 1', 'Shift 2', 'Middle', 'OFF'
}

type BulkScheduleInput struct {
	Schedules []ScheduleItem `json:"schedules" binding:"required"`
}

// 📅 1. FUNGSI SUBMIT JADWAL (BULK INSERT / UPSERT)
func SaveSchedules(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	var storeID uint
	if val, ok := storeIDRaw.(float64); ok {
		storeID = uint(val)
	} else if val, ok := storeIDRaw.(uint); ok {
		storeID = val
	}

	var input BulkScheduleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data jadwal tidak valid!"})
		return
	}

	tx := config.DB.Begin() // Gunakan transaction biar kalau satu gagal, semua di-rollback

	for _, item := range input.Schedules {
		// Tentukan jam kerja default berdasarkan shift biar karyawan ga perlu input manual
		jamMasuk := "-"
		jamPulang := "-"
		if item.ShiftType == "Shift 1" {
			jamMasuk = "07:00"
			jamPulang = "15:00"
		} else if item.ShiftType == "Shift 2" {
			jamMasuk = "15:00"
			jamPulang = "23:00"
		} else if item.ShiftType == "Middle" {
			jamMasuk = "11:00"
			jamPulang = "19:00"
		}

		var existing models.Schedule
		// Cek apakah user sudah punya jadwal di tanggal tersebut? Kalau ada, kita UPDATE (Upsert logic ala TSM)
		err := tx.Where("user_id = ? AND tanggal = ?", item.UserID, item.Tanggal).First(&existing).Error
		
		if err == nil {
			// Update jadwal yang sudah ada
			existing.ShiftType = item.ShiftType
			existing.JamMasukJadwal = jamMasuk
			existing.JamPulangJadwal = jamPulang
			if err := tx.Save(&existing).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui jadwal lama!"})
				return
			}
		} else {
			// Buat jadwal baru jika belum ada
			newSchedule := models.Schedule{
				StoreID:         storeID,
				UserID:          item.UserID,
				Tanggal:         item.Tanggal,
				ShiftType:       item.ShiftType,
				JamMasukJadwal:  jamMasuk,
				JamPulangJadwal: jamPulang,
			}
			if err := tx.Create(&newSchedule).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan jadwal baru!"})
				return
			}
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Jadwal mingguan berhasil disimpan! 🚀"})
}

// 🔍 2. FUNGSI AMBIL DATA JADWAL (DENGAN FILTER RANGE TANGGAL)
func GetSchedules(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	var storeID uint
	if val, ok := storeIDRaw.(float64); ok {
		storeID = uint(val)
	} else if val, ok := storeIDRaw.(uint); ok {
		storeID = val
	}

	startDate := c.Query("start_date") // Contoh: 2026-05-11
	endDate := c.Query("end_date")     // Contoh: 2026-05-17

	var listJadwal []models.Schedule
	query := config.DB.Preload("User").Where("store_id = ?", storeID)

	// Filter berdasarkan range mingguan jika start dan end date diisi dari Vue
	if startDate != "" && endDate != "" {
		query = query.Where("tanggal BETWEEN ? AND ?", startDate, endDate)
	}

	if err := query.Order("tanggal ASC").Find(&listJadwal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menarik data jadwal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": listJadwal})
}