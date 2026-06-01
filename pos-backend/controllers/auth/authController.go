package auth

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"pos-backend/src/core/config"
	"pos-backend/models"
	"pos-backend/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// 1. Struct Register
type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	TempatLahir  string `json:"tempat_lahir"`
    TanggalLahir string `json:"tanggal_lahir"`
    NoHP         string `json:"no_hp"`
}

// -- REGISTER DENGAN OTP --
func Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Cek Email Terdaftar
	var existingUser models.User
	if err := src.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email sudah terdaftar!"})
		return
	}

	// 2. Hash Password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	// 3. Generate OTP (6 Digit)
	rand.Seed(time.Now().UnixNano())
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	user := models.User{
		Name:       input.Name,
		Email:      &input.Email,
		Password:   string(hashedPassword),
		Role:       "owner",
		IsVerified: false,
		OTPCode:    otp,
		OTPExpired: time.Now().Add(time.Minute * 5),
		TempatLahir:  input.TempatLahir,
        TanggalLahir: input.TanggalLahir,
        NoHP:         input.NoHP,
	}

	if err := src.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat akun"})
		return
	}

	// 4. Kirim Email OTP (Goroutine)
	go utils.SendOTPEmail(input.Email, otp)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pendaftaran berhasil! Silakan cek email untuk kode OTP.",
		"email":   input.Email,
	})
}

// -- VERIFIKASI OTP --
func VerifyOTP(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required"`
		OTP   string `json:"otp" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	var user models.User
	if err := src.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	if time.Now().After(user.OTPExpired) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kode OTP sudah kadaluarsa!"})
		return
	}

	if user.OTPCode != input.OTP {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kode OTP salah!"})
		return
	}

	// Update Status Verified
	src.DB.Model(&user).Updates(map[string]interface{}{
		"is_verified": true,
		"otp_code":    "",
	})

	c.JSON(http.StatusOK, gin.H{"message": "Email terverifikasi! Silakan login."})
}

// 2. Struct Login
type LoginInput struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

// -- LOGIN DENGAN PROTEKSI HYBRID --
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	query := src.DB.Preload("Store")

	// 🚀 DETEKTOR SILUMAN: EMAIL (OWNER) VS NO WA (KARYAWAN)
	if strings.Contains(input.Identifier, "@") {
		// Jalur VIP: Login via Email
		query = query.Where("email = ?", input.Identifier)
	} else {
		// Jalur Operasional: Bersihkan format no HP (08xxx jadi 628xxx)
		cleanHP := utils.FormatPhoneNumber(input.Identifier)
		query = query.Where("no_hp = ?", cleanHP)
	}

	// Cari kecocokan di database
	if err := query.First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Identitas tidak ditemukan! Pastikan Email atau No. WhatsApp benar."})
		return
	}

	// 🚀 Proteksi OTP: HANYA untuk login via Email (Owner)
	if strings.Contains(input.Identifier, "@") && !user.IsVerified {
		c.JSON(http.StatusForbidden, gin.H{
			"error":      "Email Anda belum diverifikasi!",
			"unverified": true,
		})
		return
	}

	// Bandingkan Password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password salah!"})
		return
	}

	// Ambil Store ID
	var storeID uint = 0
	if user.StoreID != nil {
		storeID = *user.StoreID
	}

	// Tiket JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"store_id": storeID,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte("KUNCI_RAHASIA_SUPER_KUAT_123"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencetak token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Login Sukses!",
		"token":           tokenString,
		"role":            user.Role,
		"name":            user.Name,
		"foto_url":        user.FotoURL,
		"has_setup_store": storeID != 0,
		"store_name":      user.Store.NamaToko,
		"store_logo":      user.Store.LogoURL,
		"business_type":   user.Store.BusinessType,
		"subscription_plan": user.Store.SubscriptionPlan,
    	"fitur_aktif":       user.Store.FiturAktif,
	})
}

// -- GET ME --
func GetMe(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := uint(userIDRaw.(float64))

	var user models.User
	if err := src.DB.Preload("Store").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":       user.ID,
		"name":          user.Name,
		"no_hp":		 user.NoHP,
		"nik":           user.NIK,
		"role":          user.Role,
		"is_verified":   user.IsVerified,
		"foto_url":      user.FotoURL,
		"biometric_url": user.BiometricURL,
		"tempat_lahir":  user.TempatLahir,
		"tanggal_lahir": user.TanggalLahir,
		"store_name":    user.Store.NamaToko,
		"store_logo":    user.Store.LogoURL,
		"business_type": user.Store.BusinessType,
		"subscription_plan": user.Store.SubscriptionPlan,
    	"fitur_aktif":       user.Store.FiturAktif,
	})
}

// -- UPDATE PROFIL (Nama, WA, Foto) --
func UpdateProfile(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := uint(userIDRaw.(float64))

	var user models.User
	if err := src.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan!"})
		return
	}

	// 1. Update Data Teks
	if name := c.PostForm("name"); name != "" { user.Name = name }
	if tempatLahir := c.PostForm("tempat_lahir"); tempatLahir != "" { user.TempatLahir = tempatLahir }
	if tanggalLahir := c.PostForm("tanggal_lahir"); tanggalLahir != "" { user.TanggalLahir = tanggalLahir }
	
	// 🚀 FORMAT OTOMATIS JADI 628xxx BIAR RAPI!
	if noHP := c.PostForm("no_hp"); noHP != "" {
		user.NoHP = utils.FormatPhoneNumber(noHP)
	}

	// Buat prefix nama file dari NIK
	nikClean := "user"
	if user.NIK != nil && *user.NIK != "" { nikClean = *user.NIK }

	// 2. Update Foto Profil
	if file, err := c.FormFile("foto"); err == nil {
		newFileName := fmt.Sprintf("%s_profil_%d%s", nikClean, time.Now().Unix(), filepath.Ext(file.Filename))
		uploadPath := filepath.Join("uploads", newFileName)
		if err := c.SaveUploadedFile(file, uploadPath); err == nil {
            if user.FotoURL != "" { os.Remove("." + user.FotoURL) } // Hapus foto lama
			user.FotoURL = "/uploads/" + newFileName
		}
	}

	// 3. Update Foto Biometrik
	if bioFile, errBio := c.FormFile("biometric_file"); errBio == nil {
		newBioName := fmt.Sprintf("%s_bio_%d%s", nikClean, time.Now().Unix(), filepath.Ext(bioFile.Filename))
		uploadBioPath := filepath.Join("uploads", newBioName)
		if err := c.SaveUploadedFile(bioFile, uploadBioPath); err == nil {
            if user.BiometricURL != "" { os.Remove("." + user.BiometricURL) } // Hapus bio lama
			user.BiometricURL = "/uploads/" + newBioName
		}
	}

	if err := src.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan perubahan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profil berhasil diperbarui!", "data": user})
}

// -- UPDATE PASSWORD --
func UpdatePassword(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := uint(userIDRaw.(float64))

	var input struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak lengkap!"})
		return
	}

	var user models.User
	if err := src.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan!"})
		return
	}

	// Cek apakah password lama yang diketik itu Bener?
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password lama salah!"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	if err := src.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan password baru"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil diubah!"})
}

// ========================================================
// 🚀 FUNGSI BARU: GENERATE & KIRIM KODE OTP VIA WHATSAPP (FONNTE SISTEM)
// ========================================================
func SendOTPWhatsApp(c *gin.Context) {
	var input struct {
		Phone string `json:"phone" binding:"required"`
	}

	// 1. Validasi input request dari Vue
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nomor WhatsApp wajib diisi!"})
		return
	}

	// 🚀 FIX MUTLAK: Bersihkan nomor HP dan pastikan diawali '62'
	phoneClean := input.Phone
	// Kalau nomornya dikirim pake '+' (misal +628xxx), buang plus-nya
	phoneClean = strings.Replace(phoneClean, "+", "", 1)
	
	// Kalau nomornya diawali '0' (misal 08xxx), ganti '0'-nya jadi '62'
	if strings.HasPrefix(phoneClean, "0") {
		phoneClean = "62" + phoneClean[1:]
	}
	
	// Kalau nomornya murni langsung '8xxx' (tanpa 0 dan tanpa 62), tambahin '62' di depannya
	if strings.HasPrefix(phoneClean, "8") {
		phoneClean = "62" + phoneClean
	}

	// Sekarang phoneClean dijamin isinya bakal '6289666168123'

	// 2. Cari user berdasarkan nomor HP yang sudah berformat 628xx
	var user models.User
	if err := src.DB.Where("no_hp = ?", phoneClean).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nomor WhatsApp tidak terdaftar di sistem kami!"})
		return
	}

	// 3. Generate Kode OTP Baru (6 Digit)
	rand.Seed(time.Now().UnixNano())
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	// 4. Update Kode OTP dan Waktu Expired (3 Menit) ke Database
	err := src.DB.Model(&user).Updates(map[string]interface{}{
		"otp_code":    otp,
		"otp_expired": time.Now().Add(time.Minute * 3),
	}).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengonfigurasi token verifikasi baru"})
		return
	}

	// 5. Susun template pesan OTP
	message := fmt.Sprintf(
		"Halo Bos %s!\n\nKode OTP Verifikasi Akun NEXA POS Anda adalah: *%s*\n\nKode ini rahasia dan berlaku selama 3 menit. Jangan bagikan kode ini kepada siapapun demi keamanan infrastruktur bisnis Anda. 😎", 
		user.Name, 
		otp,
	)

	// 6. Tembak Fonnte pake Token Sistem lu lewat utils
	utils.SendSystemWhatsApp(phoneClean, message)

	c.JSON(http.StatusOK, gin.H{
		"message": "Kode OTP berhasil dikirim ke WhatsApp Anda! Silakan cek chat masuk.",
		"phone":   phoneClean,
	})
}