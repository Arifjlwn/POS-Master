package delivery

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"strconv"

	"pos-backend/models"
	"pos-backend/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// =======================
// 🤝 EMPLOYEE HANDLERS
// =======================
func (h *RetailHandler) CreateEmployee(c *gin.Context) {
	roleOwner, _ := c.Get("role")
	if roleOwner != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya Owner yang bisa mendaftarkan karyawan baru!"})
		return
	}

	storeID := uint(c.MustGet("store_id").(float64))
	name := c.PostForm("name")
	password := c.PostForm("password")
	tempatLahir := c.PostForm("tempat_lahir")
	tanggalLahir := c.PostForm("tanggal_lahir")
	noHP := c.PostForm("no_hp")
	inputRole := c.PostForm("role")

	// 🚀 PERBAIKAN: Nomor HP sekarang WAJIB DIISI!
	if name == "" || password == "" || noHP == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama, Nomor HP, dan Password wajib diisi!"})
		return
	}
	if inputRole == "" {
		inputRole = "kasir"
	}

	// 🚀 FORMAT NOMOR HP JADI 628xxx SEBELUM DISIMPAN
	formattedHP := utils.FormatPhoneNumber(noHP)

	currentYear := time.Now().Format("2006")
	var newNIK string

	lastEmployee, err := h.Repo.GetLastEmployeeNIK(storeID, currentYear)
	if err != nil {
		// 🚀 JIKA BELUM ADA, MULAI DARI TAHUN + 0001 (Contoh: 20260001)
		newNIK = fmt.Sprintf("%s0001", currentYear)
	} else {
		lastNIK := *lastEmployee.NIK

		// 🚀 AMBIL 4 DIGIT TERAKHIR UNTUK URUTAN (Bukan 3 digit lagi)
		if len(lastNIK) >= 4 {
			lastSequence, _ := strconv.Atoi(lastNIK[len(lastNIK)-4:])
			newNIK = fmt.Sprintf("%s%04d", currentYear, lastSequence+1)
		} else {
			newNIK = fmt.Sprintf("%s0001", currentYear)
		}
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	file, err := c.FormFile("foto")
	var fotoURL string
	if err == nil {
		newFileName := fmt.Sprintf("%s_%d%s", newNIK, time.Now().Unix(), filepath.Ext(file.Filename))
		uploadPath := filepath.Join("uploads", newFileName)
		c.SaveUploadedFile(file, uploadPath)
		fotoURL = "/uploads/" + newFileName
	}

	employee := models.User{
		StoreID:      &storeID,
		Name:         name,
		NIK:          &newNIK,
		Password:     string(hashedPassword),
		Role:         inputRole,
		TempatLahir:  tempatLahir,
		TanggalLahir: tanggalLahir,
		NoHP:         formattedHP, // 🚀 MASUKIN NOMOR HP YANG UDAH DIBERSIHKAN
		FotoURL:      fotoURL,
	}

	if err := h.Repo.CreateEmployee(&employee); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan ke database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Karyawan baru berhasil didaftarkan! 🤝",
		"data":    gin.H{"nama": employee.Name, "nik": newNIK, "jabatan": employee.Role},
	})
}

func (h *RetailHandler) GetEmployees(c *gin.Context) {
	storeID := uint(c.MustGet("store_id").(float64))
	employees, err := h.Repo.GetAllEmployees(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": employees})
}

func (h *RetailHandler) UpdateEmployee(c *gin.Context) {
	roleOwner, _ := c.Get("role")
	if roleOwner != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya owner yang bisa edit data tim!"})
		return
	}

	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	employee, err := h.Repo.GetEmployeeByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Karyawan tidak ditemukan!"})
		return
	}

	employee.Name = c.PostForm("name")
	employee.TempatLahir = c.PostForm("tempat_lahir")
	employee.TanggalLahir = c.PostForm("tanggal_lahir")

	// 🚀 FORMAT NOMOR HP JUGA SAAT UPDATE
	if inputHP := c.PostForm("no_hp"); inputHP != "" {
		employee.NoHP = utils.FormatPhoneNumber(inputHP)
	}

	if newRole := c.PostForm("role"); newRole != "" {
		employee.Role = newRole
	}
	if password := c.PostForm("password"); password != "" {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		employee.Password = string(hashed)
	}

	nikClean := "karyawan"
	if employee.NIK != nil {
		nikClean = *employee.NIK
	}

	if file, err := c.FormFile("foto"); err == nil {
		newFileName := fmt.Sprintf("%s_%d%s", nikClean, time.Now().Unix(), filepath.Ext(file.Filename))
		uploadPath := filepath.Join("uploads", newFileName)
		if err := c.SaveUploadedFile(file, uploadPath); err == nil {
			if employee.FotoURL != "" {
				os.Remove("." + employee.FotoURL)
			}
			employee.FotoURL = "/uploads/" + newFileName
		}
	}

	if bioFile, errBio := c.FormFile("biometric_file"); errBio == nil {
		newBioName := fmt.Sprintf("%s_bio_%d%s", nikClean, time.Now().Unix(), filepath.Ext(bioFile.Filename))
		uploadBioPath := filepath.Join("uploads", newBioName)
		if err := c.SaveUploadedFile(bioFile, uploadBioPath); err == nil {
			if employee.BiometricURL != "" {
				os.Remove("." + employee.BiometricURL)
			}
			employee.BiometricURL = "/uploads/" + newBioName
		}
	}

	if err := h.Repo.SaveEmployee(employee); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan perubahan ke database"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil diperbarui! 💾", "data": employee})
}

// ========================
// 📅 SCHEDULE HANDLERS
// ========================
type ScheduleItem struct {
	UserID    uint   `json:"user_id" binding:"required"`
	Tanggal   string `json:"tanggal" binding:"required"`
	ShiftType string `json:"shift_type" binding:"required"`
}

type BulkScheduleInput struct {
	Schedules []ScheduleItem `json:"schedules" binding:"required"`
}

func (h *RetailHandler) SaveSchedules(c *gin.Context) {
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

	db := h.Repo.GetDB()
	tx := db.Begin()

	for _, item := range input.Schedules {
		jamMasuk, jamPulang := "-", "-"
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

		existing, err := h.Repo.GetScheduleByDate(tx, item.UserID, item.Tanggal)
		if err == nil {
			existing.ShiftType = item.ShiftType
			existing.JamMasukJadwal = jamMasuk
			existing.JamPulangJadwal = jamPulang
			if err := h.Repo.SaveScheduleTx(tx, existing); err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui jadwal lama!"})
				return
			}
		} else {
			newSchedule := models.Schedule{
				StoreID:         storeID,
				UserID:          item.UserID,
				Tanggal:         item.Tanggal,
				ShiftType:       item.ShiftType,
				JamMasukJadwal:  jamMasuk,
				JamPulangJadwal: jamPulang,
			}
			if err := h.Repo.CreateScheduleTx(tx, &newSchedule); err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan jadwal baru!"})
				return
			}
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Jadwal mingguan berhasil disimpan! 🚀"})
}

func (h *RetailHandler) GetSchedules(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	var storeID uint
	if val, ok := storeIDRaw.(float64); ok {
		storeID = uint(val)
	} else if val, ok := storeIDRaw.(uint); ok {
		storeID = val
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	listJadwal, err := h.Repo.GetSchedulesRange(storeID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menarik data jadwal"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": listJadwal})
}

// ================================
// 🚀 ATTENDANCE HANDLERS
// ================================

type AbsenInput struct {
	UserID uint   `json:"user_id" binding:"required"`
	Jenis  string `json:"jenis" binding:"required"`
	Foto   string `json:"foto" binding:"required"`
}

func (h *RetailHandler) StoreAttendance(c *gin.Context) {
	storeIDRaw, exists := c.Get("store_id")
	if !exists || storeIDRaw == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Toko tidak terdeteksi! Pastikan akun sudah terhubung."})
		return
	}

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

	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	today := now.Format("2006-01-02")
	nowTime := now.Format("15:04:05")

	attendance, err := h.Repo.GetAttendanceToday(input.UserID, today)

	if input.Jenis == "Masuk" {
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Anda sudah melakukan Absen Masuk hari ini!"})
			return
		}

		absen := models.Attendance{
			StoreID:   storeID,
			UserID:    input.UserID,
			Tanggal:   today,
			JamMasuk:  nowTime,
			FotoMasuk: input.Foto,
			Status:    "Hadir",
		}

		if err := h.Repo.CreateAttendance(&absen); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan absensi masuk!"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Absen Masuk Berhasil! Selamat Bekerja."})

	} else if input.Jenis == "Pulang" {
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Anda belum melakukan Absen Masuk hari ini!"})
			return
		}

		if attendance.JamPulang != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Anda sudah melakukan Absen Pulang hari ini!"})
			return
		}

		attendance.JamPulang = nowTime
		attendance.FotoPulang = input.Foto

		if err := h.Repo.SaveAttendance(attendance); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan absen pulang!"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Absen Pulang Berhasil! Hati-hati di jalan."})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Jenis absen tidak dikenali!"})
	}
}

func (h *RetailHandler) GetAttendance(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	var storeID uint
	if val, ok := storeIDRaw.(float64); ok {
		storeID = uint(val)
	} else if val, ok := storeIDRaw.(uint); ok {
		storeID = val
	}

	tanggal := c.Query("tanggal")
	bulan := c.Query("bulan")
	tahun := c.Query("tahun")

	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	todayStr := now.Format("2006-01-02")

	var prefixBulan string
	if tanggal == "" && bulan != "" && tahun != "" {
		prefixBulan = fmt.Sprintf("%s-%s-%%", tahun, bulan)
	}

	riwayat, err := h.Repo.GetAttendanceReport(storeID, tanggal, prefixBulan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menarik data log absensi"})
		return
	}

	db := h.Repo.GetDB()
	for i := 0; i < len(riwayat); i++ {
		if riwayat[i].JamMasuk != "" && riwayat[i].JamPulang != "" {
			riwayat[i].Status = "Hadir"
		} else if riwayat[i].JamMasuk != "" && riwayat[i].JamPulang == "" {
			if riwayat[i].Tanggal < todayStr {
				riwayat[i].Status = "Lupa Absen Pulang"
				db.Model(&riwayat[i]).Update("status", "Lupa Absen Pulang")
			} else {
				riwayat[i].Status = "Hadir"
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": riwayat})
}

func (h *RetailHandler) ExportAttendance(c *gin.Context) {
	storeID := uint(c.MustGet("store_id").(float64))
	bulan := c.Query("bulan")
	tahun := c.Query("tahun")

	if bulan == "" || tahun == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bulan dan tahun harus diisi!"})
		return
	}

	prefixBulan := fmt.Sprintf("%s-%s-%%", tahun, bulan)
	riwayat, err := h.Repo.GetAttendanceReport(storeID, "", prefixBulan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses ekspor laporan"})
		return
	}

	b := &bytes.Buffer{}
	w := csv.NewWriter(b)
	w.Write([]string{"Tanggal", "NIK", "Nama Karyawan", "Jam Masuk", "Jam Pulang", "Status"})

	for _, logData := range riwayat {
		nik := "-"
		if logData.User.NIK != nil {
			nik = *logData.User.NIK
		}

		jamPulang := logData.JamPulang
		if jamPulang == "" {
			jamPulang = "Belum Pulang"
		}

		w.Write([]string{
			logData.Tanggal,
			nik,
			logData.User.Name,
			logData.JamMasuk,
			jamPulang,
			logData.Status,
		})
	}
	w.Flush()

	filename := fmt.Sprintf("Laporan_Absensi_%s_%s.csv", bulan, tahun)
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "text/csv", b.Bytes())
}