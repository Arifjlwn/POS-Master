package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"pos-backend/config"
	"pos-backend/models"
	"strconv" // 🚀 Tambahkan ini untuk Atoi
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Fungsi Menambahkan Karyawan (Create)
func CreateEmployee(c *gin.Context) {
	// 1. Cek role wajib Owner
	storeIDRaw, _ := c.Get("store_id")
	role, _ := c.Get("role")

	if role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya Owner yang bisa mendaftarkan karyawan baru!"})
		return
	}

	storeID := uint(storeIDRaw.(float64))

	// 2. Tangkap inputan dari Multipart Form (karena ada upload foto)
	name := c.PostForm("name")
	password := c.PostForm("password")
	tempatLahir := c.PostForm("tempat_lahir")
	tanggalLahir := c.PostForm("tanggal_lahir")
	noHP := c.PostForm("no_hp")

	if name == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama dan Password wajib diisi!"})
		return
	}

	// 3. LOGIKA GENERATE NIK OTOMATIS
	currentYear := time.Now().Format("2006")
	var lastEmployee models.User
	var newNIK string

	err := config.DB.Where("store_id = ? AND role = ? AND nik LIKE ?", storeID, "kasir", currentYear+"%").
		Order("nik desc").
		First(&lastEmployee).Error

	if err != nil {
		newNIK = currentYear + "0001"
	} else {
		lastNIK := *lastEmployee.NIK
		lastSequenceStr := lastNIK[4:]
		lastSequence, _ := strconv.Atoi(lastSequenceStr)
		newSequence := lastSequence + 1
		newNIK = fmt.Sprintf("%s%04d", currentYear, newSequence)
	}

	// 4. Hash Password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// 5. Handle Foto (Wajib ada di pendaftaran pertama sesuai request Mas)
	file, err := c.FormFile("foto")
	var fotoURL string
	if err == nil {
		newFileName := fmt.Sprintf("%s_%d%s", newNIK, time.Now().Unix(), filepath.Ext(file.Filename))
		uploadPath := filepath.Join("uploads", newFileName)
		c.SaveUploadedFile(file, uploadPath)
		fotoURL = "/uploads/" + newFileName
	}

	// 6. Simpan ke database
	employee := models.User{
		StoreID:      &storeID,
		Name:         name,
		NIK:          &newNIK,
		Password:     string(hashedPassword),
		Role:         "kasir",
		TempatLahir:  tempatLahir,
		TanggalLahir: tanggalLahir,
		NoHP:         noHP,
		FotoURL:      fotoURL,
	}

	if err := config.DB.Create(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan ke database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Karyawan baru berhasil didaftarkan! 🤝",
		"data": gin.H{
			"nama": employee.Name,
			"nik":  newNIK,
		},
	})
}

// Fungsi Lihat Daftar Karyawan
func GetEmployees(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	// role, _ := c.Get("role")

	// if role != "owner" {
	// 	c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak!"})
	// 	return
	// }

	storeID := uint(storeIDRaw.(float64))
	var employees []models.User

	if err := config.DB.Where("store_id = ?", storeID).Find(&employees).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
        return
    }

	c.JSON(http.StatusOK, gin.H{"data": employees})
}

// Fungsi Update Karyawan (Edit)
func UpdateEmployee(c *gin.Context) {
	// Cek Role Owner
	role, _ := c.Get("role")
	if role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya owner yang bisa edit data tim!"})
		return
	}

	id := c.Param("id")
	var employee models.User

	if err := config.DB.First(&employee, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Karyawan tidak ditemukan!"})
		return
	}

	// Update data teks
	employee.Name = c.PostForm("name")
	employee.TempatLahir = c.PostForm("tempat_lahir")
	employee.TanggalLahir = c.PostForm("tanggal_lahir")
	employee.NoHP = c.PostForm("no_hp")

	// Update Password jika diisi
	password := c.PostForm("password")
	if password != "" {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		employee.Password = string(hashed)
	}

	// Handle Update Foto
	file, err := c.FormFile("foto")
	if err == nil {
		newFileName := fmt.Sprintf("%s_%d%s", *employee.NIK, time.Now().Unix(), filepath.Ext(file.Filename))
		uploadPath := filepath.Join("uploads", newFileName)

		if err := c.SaveUploadedFile(file, uploadPath); err == nil {
			// Hapus foto lama biar gak menuh-menuhin PC i5 Mas Arif
			if employee.FotoURL != "" {
				os.Remove("." + employee.FotoURL)
			}
			employee.FotoURL = "/uploads/" + newFileName
		}
	}

	config.DB.Save(&employee)

	c.JSON(http.StatusOK, gin.H{
		"message": "Data berhasil diperbarui!",
		"data":    employee,
	})
}