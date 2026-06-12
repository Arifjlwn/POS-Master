package admin

import (
	"net/http"
	"os"
	"time"

	"pos-backend/models"
	"pos-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AdminLogin handles 🔐 OTO-1: LOGIN HANDLER PUSAT SUPER ADMIN
func (a *AuthController) AdminLogin(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data input login tidak valid !"})
		return
	}

	var user models.User
	if err := a.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kredensial login admin salah atau tidak ditemukan!"})
		return
	}

	if user.Role != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "AKSES DITOLAK: Akun Anda tidak memegang kedaulatan Root Admin!"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kredensial login admin salah atau tidak ditemukan!"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"email":   user.Email,
		"exp":     time.Now().Add(12 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengonfigurasi token kedaulatan pusat!"})
		return
	}

	c.Set("user_id", float64(user.ID))
	c.Set("email", *user.Email)

	utils.RecordAdminAction(c, "Otorisasi Super Admin", "-", "Berhasil masuk ke platform Mission Control")

	c.JSON(http.StatusOK, gin.H{
		"status": "sukses",
		"token":  tokenString,
		"user": gin.H{
			"name":  user.Name,
			"role":  user.Role,
			"email": user.Email,
		},
	})
}