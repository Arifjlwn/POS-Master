package delivery

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"pos-backend/models"
	"pos-backend/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ==========================================
// 🤝 MANAGEMENT TIM & KARYAWAN HANDLERS
// ==========================================

func (h *RetailHandler) CreateEmployee(c *gin.Context) {
    roleOwner, _ := c.Get("role")
    if roleOwner != "owner" {
        c.JSON(http.StatusForbidden, gin.H{"error": "Hanya Owner yang bisa mendaftarkan karyawan baru!"})
        return
    }

    storeID := getStoreID(c)

    name := c.PostForm("name")
    password := c.PostForm("password")
    tempatLahir := c.PostForm("tempat_lahir")
    tanggalLahir := c.PostForm("tanggal_lahir")
    noHP := c.PostForm("no_hp")
    inputRole := c.PostForm("role")

    if name == "" || password == "" || noHP == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Nama, Nomor HP, dan Password wajib diisi!"})
        return
    }
    if inputRole == "" {
        inputRole = "kasir"
    }

    formattedHP := utils.FormatPhoneNumber(noHP)
    currentYear := time.Now().Format("2006")
    var newNIK string

    lastEmployee, err := h.Repo.GetLastEmployeeNIK(storeID, currentYear)
    if err != nil {
        newNIK = fmt.Sprintf("%s0001", currentYear)
    } else {
        lastNIK := *lastEmployee.NIK
        if len(lastNIK) >= 4 {
            lastSequence, _ := strconv.Atoi(lastNIK[len(lastNIK)-4:])
            newNIK = fmt.Sprintf("%s%04d", currentYear, lastSequence+1)
        } else {
            newNIK = fmt.Sprintf("%s0001", currentYear)
        }
    }

    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    var fotoURL string
    bucketName := os.Getenv("SUPABASE_BUCKET_NAME")

    db := h.Repo.GetDB()
    var store models.Store
    if err := db.Select("public_id").First(&store, storeID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Verifikasi keamanan toko gagal"})
        return
    }

    if file, err := c.FormFile("foto"); err == nil {
        if file.Size > 5*1024*1024 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran foto maksimal 5 MB"})
            return
        }
        contentType := file.Header.Get("Content-Type")
        fileSrc, _ := file.Open()
        defer fileSrc.Close()

        remotePath := fmt.Sprintf("stores/%s/employees/avatar_%s_%d", store.PublicID, newNIK, time.Now().Unix())
        urlResult, errUpload := utils.UploadToSupabase(fileSrc, file.Filename, contentType, bucketName, remotePath)
        if errUpload == nil {
            fotoURL = urlResult
        }
    }

    employee := models.User{
        PublicID:     utils.GenerateULID(),
        StoreID:      &storeID,
        Name:         name,
        NIK:          &newNIK,
        Password:     string(hashedPassword),
        Role:         inputRole,
        TempatLahir:  tempatLahir,
        TanggalLahir: tanggalLahir,
        NoHP:         formattedHP,
        FotoURL:      fotoURL,
		IsVerified:   true,
    }

    if err := h.Repo.CreateEmployee(&employee); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan data ke cloud database"})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"message": "Karyawan baru berhasil didaftarkan!", "data": gin.H{"nama": employee.Name, "nik": newNIK, "jabatan": employee.Role}})
}

func (h *RetailHandler) GetEmployees(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	var storeID uint
	switch v := storeIDRaw.(type) {
	case float64:
		storeID = uint(v)
	case uint:
		storeID = v
	case int:
		storeID = uint(v)
	}

	employees, err := h.Repo.GetAllEmployees(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data karyawan"})
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

    storeID := getStoreID(c)
    publicID := c.Param("id") 

    db := h.Repo.GetDB()
    var employee models.User
    
    if err := db.Where("public_id = ? AND store_id = ?", publicID, storeID).First(&employee).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Karyawan tidak ditemukan!"})
        return
    }

    employee.Name = c.PostForm("name")
    employee.TempatLahir = c.PostForm("tempat_lahir")
    employee.TanggalLahir = c.PostForm("tanggal_lahir")

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
    bucketName := os.Getenv("SUPABASE_BUCKET_NAME")

    var store models.Store
    if err := db.Select("public_id").First(&store, storeID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memverifikasi keamanan penyimpanan data."})
        return
    }

    if file, err := c.FormFile("foto"); err == nil {
        contentType := file.Header.Get("Content-Type")
        fileSrc, _ := file.Open()
        defer fileSrc.Close()

        remotePath := fmt.Sprintf("stores/%s/employees/avatar_%s_%d", store.PublicID, nikClean, time.Now().Unix())
        if urlResult, errUpload := utils.UploadToSupabase(fileSrc, file.Filename, contentType, bucketName, remotePath); errUpload == nil {
            employee.FotoURL = urlResult
        }
    }

    if bioFile, errBio := c.FormFile("biometric_file"); errBio == nil {
        contentType := bioFile.Header.Get("Content-Type")
        fileSrc, _ := bioFile.Open()
        defer fileSrc.Close()

        remotePath := fmt.Sprintf("stores/%s/employees/biometric_%s_%d", store.PublicID, nikClean, time.Now().Unix())
        if urlResult, errUpload := utils.UploadToSupabase(fileSrc, bioFile.Filename, contentType, bucketName, remotePath); errUpload == nil {
            employee.BiometricURL = urlResult
        }
    }

    if err := h.Repo.SaveEmployee(&employee); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan perubahan"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Data berhasil diperbarui!", "data": employee})
}

func (h *RetailHandler) DeleteEmployee(c *gin.Context) {
    roleOwner, _ := c.Get("role")
    if roleOwner != "owner" {
        c.JSON(http.StatusForbidden, gin.H{"error": "Hanya owner yang berhak memecat karyawan!"})
        return
    }

    storeID := getStoreID(c)
    publicID := c.Param("id")

    db := h.Repo.GetDB()
    var employee models.User
    
    if err := db.Where("public_id = ? AND store_id = ?", publicID, storeID).First(&employee).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Karyawan tidak ditemukan!"})
        return
    }

    if employee.Role == "owner" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Data Owner tidak bisa dihapus!"})
        return
    }

    if err := db.Delete(&employee).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data karyawan."})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Karyawan berhasil diberhentikan dari sistem."})
}

// ==========================================
// 📅 SHIFT & JADWAL KARYAWAN HANDLERS
// ==========================================

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
	switch v := storeIDRaw.(type) {
	case float64:
		storeID = uint(v)
	case uint:
		storeID = v
	case int:
		storeID = uint(v)
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
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update jadwal"})
				return
			}
		} else {
			newSchedule := models.Schedule{
				StoreID: storeID, UserID: item.UserID, Tanggal: item.Tanggal, ShiftType: item.ShiftType, JamMasukJadwal: jamMasuk, JamPulangJadwal: jamPulang,
			}
			if err := h.Repo.CreateScheduleTx(tx, &newSchedule); err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan jadwal baru"})
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
	switch v := storeIDRaw.(type) {
	case float64:
		storeID = uint(v)
	case uint:
		storeID = v
	case int:
		storeID = uint(v)
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

// ==========================================
// 🚀 PRESENSI & ATTENDANCE LOG HANDLERS
// ==========================================

type AbsenInput struct {
	UserID uint   `json:"user_id" binding:"required"`
	Jenis  string `json:"jenis" binding:"required"`
	Foto   string `json:"foto" binding:"required"`
}

func (h *RetailHandler) StoreAttendance(c *gin.Context) {
    storeID := getStoreID(c)
    
    userIDRaw, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak valid!"})
        return
    }
    var userID uint
    switch v := userIDRaw.(type) {
    case float64:
        userID = uint(v)
    case uint:
        userID = v
    case int:
        userID = uint(v)
    }

    var input AbsenInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid!"})
        return
    }

    // PAKSA PAKAI ID DARI TOKEN, TOLAK INPUT JSON DARI FRONTEND
    input.UserID = userID

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
            StoreID: storeID, UserID: input.UserID, Tanggal: today, JamMasuk: nowTime, FotoMasuk: input.Foto, Status: "Hadir",
        }
        if err := h.Repo.CreateAttendance(&absen); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal absensi masuk!"})
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
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal absen pulang!"})
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
	switch v := storeIDRaw.(type) {
	case float64:
		storeID = uint(v)
	case uint:
		storeID = v
	case int:
		storeID = uint(v)
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menarik log absensi"})
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
		w.Write([]string{logData.Tanggal, nik, logData.User.Name, logData.JamMasuk, jamPulang, logData.Status})
	}
	w.Flush()

	filename := fmt.Sprintf("Laporan_Absensi_%s_%s.csv", bulan, tahun)
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "text/csv", b.Bytes())
}
