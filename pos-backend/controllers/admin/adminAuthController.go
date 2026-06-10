package admin

import (
	"net/http"
	"os"
	"pos-backend/models"
	src "pos-backend/src/core/config"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// LoginAdminHandler menangani autentikasi masuk ke Mission Control
func LoginAdminHandler(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input tidak valid bray"})
		return
	}

	var user models.User
	cleanEmail := strings.ToLower(strings.TrimSpace(input.Email))

	// 1. Cari user berdasarkan email
	if err := src.DB.Where("LOWER(email) = ?", cleanEmail).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kredensial login admin salah atau tidak terdaftar!"})
		return
	}

	// 2. 🔒 PROTEKSI KETAT: Cek apakah user ini emang super_admin bray
	if user.Role != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "AKSES ILEGAL: Akun Anda tidak memiliki otoritas Root Admin!"})
		return
	}

	// 3. Verifikasi Password bray
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kredensial login admin salah!"})
		return
	}

	// 4. Rilis Token JWT Khusus Admin (Bisa disetel durasi pendek, misal 8 Jam)
	jwtSecret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"public_id": user.PublicID,
		"role":      user.Role,
		"exp":       time.Now().Add(time.Hour * 8).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal merilis token otoritas pusat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Otoritas diberikan! Selamat datang Founder.",
		"token":   tokenString,
		"role":    user.Role,
		"name":    user.Name,
	})
}