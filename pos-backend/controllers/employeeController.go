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
	roleOwner, _ := c.Get("role")

	if roleOwner != "owner" {
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
	inputRole := c.PostForm("role") // 🚀 TANGKAP ROLE DINAMIS DARI VUE

	if name == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama dan Password wajib diisi!"})
		return
	}

	// Jika dari frontend kosong, kasih default terendah ("kasir")
	if inputRole == "" {
		inputRole = "kasir"
	}

	// 3. LOGIKA GENERATE NIK OTOMATIS (DIPERBAIKI)
	currentYear := time.Now().Format("2006")
	var lastEmployee models.User
	var newNIK string

	// 🚀 FIX: Ubah dari query `role = "kasir"` menjadi `role != "owner"` 
	// Supaya manager, supervisor, dan kasir urutan nomor urut NIK-nya menyatu dalam 1 cabang toko!
	err := config.DB.Where("store_id = ? AND role != ? AND nik LIKE ?", storeID, "owner", currentYear+"%").
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

	// 5. Handle Foto
	file, err := c.FormFile("foto")
	var fotoURL string
	if err == nil {
		newFileName := fmt.Sprintf("%s_%d%s", newNIK, time.Now().Unix(), filepath.Ext(file.Filename))
		uploadPath := filepath.Join("uploads", newFileName)
		c.SaveUploadedFile(file, uploadPath)
		fotoURL = "/uploads/" + newFileName
	}

	// 6. Simpan ke database menggunakan inputRole dinamis
	employee := models.User{
		StoreID:      &storeID,
		Name:         name,
		NIK:          &newNIK,
		Password:     string(hashedPassword),
		Role:         inputRole,
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
			"nama":    employee.Name,
			"nik":     newNIK,
			"jabatan": employee.Role,
		},
	})
}

// Fungsi Lihat Daftar Karyawan
func GetEmployees(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))
	var employees []models.User

	// 🚀 SUDAH BERSIH TANPA SATPAM "owner": Diizinkan ditarik oleh manager/kasir untuk kebutuhan matriks TSM
	if err := config.DB.Where("store_id = ?", storeID).Find(&employees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": employees})
}

// Fungsi Update Karyawan (Edit)
func UpdateEmployee(c *gin.Context) {
	// 1. Cek Role Owner lewat token JWT
	roleOwner, _ := c.Get("role")
	if roleOwner != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya owner yang bisa edit data tim!"})
		return
	}

	id := c.Param("id")
	var employee models.User

	if err := config.DB.First(&employee, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Karyawan tidak ditemukan!"})
		return
	}

	// 2. Tangkap semua data perubahan form teks dari Vue
	employee.Name = c.PostForm("name")
	employee.TempatLahir = c.PostForm("tempat_lahir")
	employee.TanggalLahir = c.PostForm("tanggal_lahir")
	employee.NoHP = c.PostForm("no_hp")
	
	newRole := c.PostForm("role")
	if newRole != "" {
		employee.Role = newRole
	}

	// Update Password jika diisi oleh Owner
	password := c.PostForm("password")
	if password != "" {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		employee.Password = string(hashed)
	}

	// 3. Handle Update Foto (Diberi Pelindung Anti-Nil Pointer)
	file, err := c.FormFile("foto")
	if err == nil {
		// 🚀 AMANKAN POINTER NIK: Jika pointer nil, gunakan nama file default "karyawan"
		nikClean := "karyawan"
		if employee.NIK != nil {
			nikClean = *employee.NIK
		}

		newFileName := fmt.Sprintf("%s_%d%s", nikClean, time.Now().Unix(), filepath.Ext(file.Filename))
		uploadPath := filepath.Join("uploads", newFileName)

		if err := c.SaveUploadedFile(file, uploadPath); err == nil {
			// Hapus foto lama di PC Mas Arif biar storage gak bengkak
			if employee.FotoURL != "" {
				os.Remove("." + employee.FotoURL)
			}
			employee.FotoURL = "/uploads/" + newFileName
		}
	}

	// 🚀 4. HANDLE UPDATE FOTO BIOMETRIK (ABSENSI)
bioFile, errBio := c.FormFile("biometric_file")
if errBio == nil {
    // Amankan Pointer NIK lagi
    nikClean := "karyawan"
    if employee.NIK != nil {
        nikClean = *employee.NIK
    }

    // Kasih nama beda (tambah _bio_) biar file-nya nggak nimpa foto profil
    newBioName := fmt.Sprintf("%s_bio_%d%s", nikClean, time.Now().Unix(), filepath.Ext(bioFile.Filename))
    uploadBioPath := filepath.Join("uploads", newBioName)

    if err := c.SaveUploadedFile(bioFile, uploadBioPath); err == nil {
        // Hapus foto biometrik lama jika ada
        if employee.BiometricURL != "" {
            os.Remove("." + employee.BiometricURL)
        }
        // Simpan path baru ke database
        employee.BiometricURL = "/uploads/" + newBioName
    }
}

	// 4. Eksekusi simpan perubahan aman ke Supabase
	if err := config.DB.Save(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan perubahan ke database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data berhasil diperbarui! 💾",
		"data":    employee,
	})
}