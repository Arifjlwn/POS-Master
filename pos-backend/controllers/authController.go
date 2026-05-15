package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"pos-backend/config"
	"pos-backend/models"
	"pos-backend/utils" // 🚀 Pastikan utils.SendOTPEmail sudah Mas buat
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
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
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

	if err := config.DB.Create(&user).Error; err != nil {
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
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
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
	config.DB.Model(&user).Updates(map[string]interface{}{
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

// -- LOGIN DENGAN PROTEKSI --
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	// 🚀 Cari User berdasarkan Email atau NIK
	query := config.DB.Preload("Store")
	if strings.Contains(input.Identifier, "@") {
		query = query.Where("email = ?", input.Identifier)
	} else {
		query = query.Where("nik = ?", input.Identifier)
	}

	if err := query.First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Identitas tidak ditemukan!"})
		return
	}

	// 🚀 Proteksi OTP: Hanya untuk login via Email
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
	})
}

// -- GET ME --
func GetMe(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := uint(userIDRaw.(float64))

	var user models.User
	if err := config.DB.Preload("Store").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":       user.ID,
		"name":          user.Name,
		"nik":           user.NIK,
		"role":          user.Role,
		"is_verified":   user.IsVerified,
		"foto_url":      user.FotoURL,
		"biometric_url": user.BiometricURL,
		"tempat_lahir":  user.TempatLahir,
		"tanggal_lahir": user.TanggalLahir,
		"store_name":    user.Store.NamaToko,
	})
}