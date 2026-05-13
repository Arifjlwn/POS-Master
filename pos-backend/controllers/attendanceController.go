package controllers

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"pos-backend/config"
	"pos-backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

type AbsenInput struct {
	UserID uint   `json:"user_id" binding:"required"`
	Jenis  string `json:"jenis" binding:"required"` // "Masuk" atau "Pulang"
	Foto   string `json:"foto" binding:"required"`  // Teks Base64 dari Vue
}

// 📸 FUNGSI REKAM ABSENSI (MASUK & PULANG)
func StoreAttendance(c *gin.Context) {
	storeIDRaw, exists := c.Get("store_id")
	if !exists || storeIDRaw == nil {
        c.JSON(http.StatusForbidden, gin.H{"error": "Toko tidak terdeteksi! Pastikan akun sudah terhubung ke toko."})
        return
    }
    
    // Pastikan konversi tipe data aman
    var storeID uint
    if val, ok := storeIDRaw.(float64); ok {
        storeID = uint(val)
    } else if val, ok := storeIDRaw.(uint); ok {
        storeID = val
    }

	var input AbsenInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid!"})
		return
	}

	// 🚀 SET TIMEZONE JAKARTA (WIB)
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	today := now.Format("2006-01-02")
	nowTime := now.Format("15:04:05")

	var attendance models.Attendance

	if input.Jenis == "Masuk" {
		// 🔍 Cek apakah sudah absen masuk hari ini?
		if err := config.DB.Where("user_id = ? AND tanggal = ?", input.UserID, today).First(&attendance).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Anda sudah melakukan Absen Masuk hari ini!"})
			return
		}

		// Simpan absen masuk
		absen := models.Attendance{
			StoreID:   storeID,
			UserID:    input.UserID,
			Tanggal:   today,
			JamMasuk:  nowTime,
			FotoMasuk: input.Foto,
			Status:    "Hadir",
		}

		if err := config.DB.Create(&absen).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan absensi masuk!"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Absen Masuk Berhasil! Selamat Bekerja."})

	} else if input.Jenis == "Pulang" {
		// 🔍 Cari record absen masuk hari ini untuk di-update (UPSERT Logic)
		if err := config.DB.Where("user_id = ? AND tanggal = ?", input.UserID, today).First(&attendance).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Anda belum melakukan Absen Masuk hari ini!"})
			return
		}

		// Cek apakah sudah absen pulang sebelumnya?
		if attendance.JamPulang != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Anda sudah melakukan Absen Pulang hari ini!"})
			return
		}

		// Update record yang ada (Isi jam pulang dan foto pulang)
		attendance.JamPulang = nowTime
		attendance.FotoPulang = input.Foto

		if err := config.DB.Save(&attendance).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan absen pulang!"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Absen Pulang Berhasil! Hati-hati di jalan."})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Jenis absen tidak dikenali!"})
	}
}

// 📋 FUNGSI TARIK DATA ABSENSI (HARIAN & BULANAN)
func GetAttendance(c *gin.Context) {
    storeIDRaw, _ := c.Get("store_id")
    storeID := uint(storeIDRaw.(float64))

    tanggal := c.Query("tanggal") 
    bulan := c.Query("bulan")     
    tahun := c.Query("tahun")     

    var riwayat []models.Attendance
    query := config.DB.Preload("User").Where("store_id = ?", storeID)

    // 🚀 PERBAIKAN LOGIKA FILTER
    if tanggal != "" {
        // Mode Harian: Pastikan formatnya YYYY-MM-DD
        query = query.Where("tanggal = ?", tanggal)
    } else if bulan != "" && tahun != "" {
        // Mode Bulanan
        prefixBulan := fmt.Sprintf("%s-%s-%%", tahun, bulan)
        query = query.Where("tanggal::text LIKE ?", prefixBulan)
    } else {
        // DEFAULT: Hari ini (WIB)
        loc, _ := time.LoadLocation("Asia/Jakarta")
        // Paksa format yang sama persis dengan yang disimpan saat StoreAttendance
        today := time.Now().In(loc).Format("2006-01-02")
        query = query.Where("tanggal = ?", today)
    }

    // 🔍 DEBUG: Tambahkan ini biar Mas bisa lihat di terminal VS Code
    // fmt.Println("DEBUG: Mencari absen tanggal:", tanggal)

    if err := query.Order("tanggal DESC, jam_masuk DESC").Find(&riwayat).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menarik data log absensi"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": riwayat})
}

// 📊 FUNGSI EXPORT LAPORAN ABSENSI KE CSV
func ExportAttendance(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

	bulan := c.Query("bulan")
	tahun := c.Query("tahun")

	if bulan == "" || tahun == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bulan dan tahun harus diisi!"})
		return
	}

	prefixBulan := fmt.Sprintf("%s-%s-%%", tahun, bulan)

	var riwayat []models.Attendance
	// Tambahkan ::text untuk filter ekspor juga
	config.DB.Preload("User").
		Where("store_id = ? AND tanggal::text LIKE ?", storeID, prefixBulan).
		Order("tanggal ASC").
		Find(&riwayat)

	b := &bytes.Buffer{}
	w := csv.NewWriter(b)
	w.Write([]string{"Tanggal", "NIK", "Nama Karyawan", "Jam Masuk", "Jam Pulang", "Status"})

	for _, log := range riwayat {
		nik := "-"
		if log.User.NIK != nil {
			nik = *log.User.NIK
		}

		jamPulang := log.JamPulang
		if jamPulang == "" {
			jamPulang = "Belum Pulang"
		}

		w.Write([]string{
			log.Tanggal,
			nik,
			log.User.Name,
			log.JamMasuk,
			jamPulang,
			log.Status,
		})
	}
	w.Flush()

	filename := fmt.Sprintf("Laporan_Absensi_%s_%s.csv", bulan, tahun)
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "text/csv", b.Bytes())
}