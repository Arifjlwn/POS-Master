package controllers

import (
	"net/http"
	"pos-backend/config"
	"pos-backend/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// 1. Struct Register Simpel
type RegisterInput struct {
	Name string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// -- REGISTER --
func Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Acak Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal Mengenkripsi password"})
		return
	}

	// Bikin Akun Owner
	user := models.User{
		Name: input.Name,
		Email: &input.Email,
		Password: string(hashedPassword),
		Role: "owner",
	}

	// Simpan ke DB
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal daftar, email mungkin sudah dipakai"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Akun berhasil dibuat! Silakan login untuk setup toko Anda. 🚀",
		"name": user.Name,
	})
}

// 2. Struct Login
type LoginInput struct {
	Identifier string `json:"identifier" binding:"required"` //Bisa email atau nik
	Password string `json:"password" binding:"required"`
}

//--- LOGIN ---
func Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	// Deteksi login
	if strings.Contains(input.Identifier, "@") {
    // 🚀 Tambahkan Preload("Store") di sini
    if err := config.DB.Preload("Store").Where("email = ?", input.Identifier).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Email tidak terdaftar !"})
        return
    }
} else {
    // 🚀 Dan di sini juga
    if err := config.DB.Preload("Store").Where("nik = ?", input.Identifier).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "NIK tidak ditemukan !"})
        return
    }
}

	// Cocokan password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password Salah !"})
		return
	}

	// Ekstrak ID Toko
	var storeID uint = 0
	if user.StoreID != nil {
		storeID = *user.StoreID
	}

	// Pembuatan tiket JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"store_id": storeID,
		"role": user.Role,
		"exp": time.Now().Add(time.Hour*72).Unix(),
	})

	tokenString, err := token.SignedString([]byte("KUNCI_RAHASIA_SUPER_KUAT_123"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencetak tiket masuk"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login Sukses ! Selamat Datang Bos.",
		"token": tokenString,
		"role": user.Role,
		"name": user.Name,
		"has_setup_store": storeID != 0, // Indikator buat Frontend Vue: false = suruh ke SetupToko, true = ke Dashboard
		"store_name": user.Store.NamaToko,
	})
}