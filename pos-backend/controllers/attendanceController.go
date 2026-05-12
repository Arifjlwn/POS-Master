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
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

	var input AbsenInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid!"})
		return
	}

	// Ambil tanggal & waktu server saat tombol ditekan
	today := time.Now().Format("2006-01-02")
	nowTime := time.Now().Format("15:04:05")

	if input.Jenis == "Masuk" {
		// Cek apakah hari ini sudah absen masuk?
		var existing models.Attendance
		if err := config.DB.Where("user_id = ? AND tanggal = ?", input.UserID, today).First(&existing).Error; err == nil {
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan absensi!"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Absen Masuk Berhasil!"})

	} else if input.Jenis == "Pulang" {
		// Cari absen masuknya hari ini
		var absen models.Attendance
		if err := config.DB.Where("user_id = ? AND tanggal = ?", input.UserID, today).First(&absen).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Anda belum Absen Masuk hari ini!"})
			return
		}

		// Cek apakah sudah absen pulang sebelumnya?
		if absen.JamPulang != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Anda sudah melakukan Absen Pulang hari ini!"})
			return
		}

		// Update data dengan jam pulang & foto pulang
		absen.JamPulang = nowTime
		absen.FotoPulang = input.Foto
		
		if err := config.DB.Save(&absen).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan absen pulang!"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Absen Pulang Berhasil!"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Jenis absen tidak dikenali!"})
	}
}

// 📋 FUNGSI TARIK DATA ABSENSI BERDASARKAN TANGGAL
// 📋 FUNGSI TARIK DATA ABSENSI (HARIAN & BULANAN)
// 📋 FUNGSI TARIK DATA ABSENSI (HARIAN & BULANAN)
func GetAttendance(c *gin.Context) {
    storeIDRaw, _ := c.Get("store_id")
    storeID := uint(storeIDRaw.(float64))
    
    // Ambil parameter dari URL (dikirim oleh Vue)
    tanggal := c.Query("tanggal") // Untuk Harian (YYYY-MM-DD)
    bulan := c.Query("bulan")     // Untuk Bulanan (MM)
    tahun := c.Query("tahun")     // Untuk Bulanan (YYYY)

    var riwayat []models.Attendance
    // Awali Query dengan Preload User biar Nama Karyawan muncul
    query := config.DB.Preload("User").Where("store_id = ?", storeID)

    // 🚀 LOGIKA FILTER DINAMIS
    if tanggal != "" {
        // Mode Harian
        query = query.Where("tanggal = ?", tanggal)
    } else if bulan != "" && tahun != "" {
    // Di Postgres, kolom 'tanggal' harus diubah jadi TEXT dulu baru bisa pakai LIKE
    // Kita pakai cast(tanggal as text)
    prefixBulan := fmt.Sprintf("%s-%s-%%", tahun, bulan)
    query = query.Where("tanggal::text LIKE ?", prefixBulan) // Tambahkan ::text
    } else {
        // Default: Ambil Hari ini
        today := time.Now().Format("2006-01-02")
        query = query.Where("tanggal = ?", today)
    }

    // Urutkan dari yang terbaru (Descending)
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

	bulan := c.Query("bulan") // Format: "05"
	tahun := c.Query("tahun") // Format: "2026"

	if bulan == "" || tahun == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bulan dan tahun harus diisi!"})
		return
	}

	// Filter LIKE berdasarkan YYYY-MM
	prefixBulan := fmt.Sprintf("%s-%s-%%", tahun, bulan)

	var riwayat []models.Attendance
	config.DB.Preload("User").Where("store_id = ? AND tanggal LIKE ?", storeID, prefixBulan).Order("tanggal ASC").Find(&riwayat)

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
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "text/csv", b.Bytes())
}