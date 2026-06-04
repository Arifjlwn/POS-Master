package auth

import (
	crand "crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"
	"pos-backend/models"
	src "pos-backend/src/core/config"
	"pos-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func generateOTP() string {
	n, err := crand.Int(crand.Reader, big.NewInt(1000000))
	if err != nil { return "000000" }
	return fmt.Sprintf("%06d", n.Int64())
}

type RegisterInput struct {
	Name string `json:"name" binding:"required"`; Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`; TempatLahir string `json:"tempat_lahir"`
	TanggalLahir string `json:"tanggal_lahir"`; NoHP string `json:"no_hp"`
}
type LoginInput struct { Identifier string `json:"identifier" binding:"required"`; Password string `json:"password" binding:"required"` }
type SelectStoreInput struct { StoreID uint `json:"store_id" binding:"required"` }
type ResetPasswordInput struct { Email string `json:"email" binding:"required"`; Token string `json:"token" binding:"required"`; Password string `json:"password" binding:"required,min=6"` }

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	var existingUser models.User; email := strings.ToLower(strings.TrimSpace(input.Email))
	if err := src.DB.Where("LOWER(email) = ?", email).First(&existingUser).Error; err == nil { c.JSON(http.StatusConflict, gin.H{"error": "Email sudah terdaftar!"}); return }
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengamankan password"}); return }
	otp := generateOTP()
	user := models.User{PublicID: utils.GenerateULID(), Name: input.Name, Email: &email, Password: string(hashedPassword), Role: "owner", IsVerified: false, OTPCode: otp, OTPExpired: time.Now().Add(time.Minute * 5), TempatLahir: input.TempatLahir, TanggalLahir: input.TanggalLahir, NoHP: utils.FormatPhoneNumber(input.NoHP)}
	if err := src.DB.Create(&user).Error; err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat akun tenant"}); return }
	go utils.SendOTPEmail(input.Email, otp)
	c.JSON(http.StatusCreated, gin.H{"message": "Pendaftaran berhasil! Silakan cek email untuk kode OTP.", "email": input.Email})
}

func VerifyOTP(c *gin.Context) {
	var input struct { Email string `json:"email" binding:"required"`; OTP string `json:"otp" binding:"required"`; Intent string `json:"intent"` }
	if err := c.ShouldBindJSON(&input); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"}); return }
	var user models.User
	if err := src.DB.Where("email = ? OR no_hp = ?", input.Email, input.Email).First(&user).Error; err != nil { c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"}); return }
	if user.LockedUntil != nil && user.LockedUntil.Year() == 2099 { c.JSON(http.StatusForbidden, gin.H{"error": "Akun Anda telah DI-LOCK PERMANEN karena tindakan mencurigakan. Hubungi Tim IT!"}); return }
	if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) { diff := time.Until(*user.LockedUntil); c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Akses dibekukan! Silakan tunggu %d menit lagi.", int(diff.Minutes()))}); return }
	if time.Now().After(user.OTPExpired) { c.JSON(http.StatusBadRequest, gin.H{"error": "Kode OTP sudah kadaluarsa!"}); return }
	if user.OTPCode != input.OTP {
		newAttempts := user.OTPAttempts + 1; updates := map[string]interface{}{"otp_attempts": newAttempts}
		if newAttempts >= 8 { updates["locked_until"] = time.Date(2099, 12, 31, 23, 59, 59, 0, time.UTC); src.DB.Model(&user).Updates(updates); c.JSON(http.StatusForbidden, gin.H{"error": "Terlalu banyak percobaan! Akun Anda kini DI-LOCK PERMANEN."}); return } 
		if newAttempts >= 4 { updates["locked_until"] = time.Now().Add(time.Hour * 1); src.DB.Model(&user).Updates(updates); c.JSON(http.StatusForbidden, gin.H{"error": "Terlalu banyak percobaan salah! Akses dibekukan selama 1 jam."}); return }
		src.DB.Model(&user).Updates(updates); c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Kode OTP salah! Sisa percobaan: %d kali lagi.", 4-newAttempts%4)}); return
	}
	updates := map[string]interface{}{"is_verified": true, "otp_attempts": 0, "locked_until": nil}
	if input.Intent != "reset-password" { updates["otp_code"] = "" }
	src.DB.Model(&user).Updates(updates); c.JSON(http.StatusOK, gin.H{"message": "Verifikasi sukses!"})
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	var user models.User; query := src.DB
	if strings.Contains(input.Identifier, "@") { query = query.Where("email = ?", input.Identifier) } else { query = query.Where("no_hp = ?", utils.FormatPhoneNumber(input.Identifier)) }
	if err := query.First(&user).Error; err != nil { c.JSON(http.StatusUnauthorized, gin.H{"error": "Identitas tidak ditemukan! Pastikan Email atau No. WhatsApp benar."}); return }
	if strings.Contains(input.Identifier, "@") && !user.IsVerified { c.JSON(http.StatusForbidden, gin.H{"error": "Email Anda belum diverifikasi!", "unverified": true}); return }
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil { c.JSON(http.StatusUnauthorized, gin.H{"error": "Password salah!"}); return }
	var stores []models.Store
	if user.Role == "owner" { if user.StoreID != nil { src.DB.Where("owner_id = ? OR id = ?", user.ID, *user.StoreID).Find(&stores) } else { src.DB.Where("owner_id = ?", user.ID).Find(&stores) } } else { if user.StoreID != nil { src.DB.Where("id = ?", *user.StoreID).Find(&stores) } }
	jwtSecret := os.Getenv("JWT_SECRET")
	if len(stores) > 1 || user.Role == "owner" {
		tempToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": user.ID, "public_id": user.PublicID, "role": user.Role, "is_select": true, "exp": time.Now().Add(time.Minute * 15).Unix()})
		tokenString, err := tempToken.SignedString([]byte(jwtSecret))
		if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"}); return }
		c.JSON(http.StatusOK, gin.H{"message": "Login berhasil, silakan pilih cabang.", "require_select": true, "token": tokenString, "role": user.Role, "name": user.Name, "stores": stores}); return
	}
	var storeID uint = 0; var storeName, storeLogo, planType string
	if len(stores) == 1 { storeID = stores[0].ID; storeName = stores[0].NamaToko; storeLogo = stores[0].LogoURL; planType = stores[0].SubscriptionPlan }
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": user.ID, "public_id": user.PublicID, "store_id": storeID, "plan_type": planType, "role": user.Role, "exp": time.Now().Add(time.Hour * 72).Unix()})
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"}); return }
	c.JSON(http.StatusOK, gin.H{"message": "Login Sukses!", "require_select": false, "token": tokenString, "role": user.Role, "name": user.Name, "has_setup_store": storeID != 0, "store_name": storeName, "store_logo": storeLogo, "subscription_plan": planType})
}

func SelectStore(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists { c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak valid, silakan login ulang."}); return }
	userID := uint(userIDRaw.(float64)); userRole := c.GetString("role"); var input SelectStoreInput
	if err := c.ShouldBindJSON(&input); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Pilihan cabang tidak valid"}); return }
	var user models.User
	if err := src.DB.First(&user, userID).Error; err != nil { c.JSON(http.StatusNotFound, gin.H{"error": "Data user tidak ditemukan"}); return }
	var store models.Store
	if err := src.DB.First(&store, input.StoreID).Error; err != nil { c.JSON(http.StatusNotFound, gin.H{"error": "Cabang tidak ditemukan"}); return }
	if userRole == "owner" {
		if store.OwnerID == 0 { if user.StoreID != nil && *user.StoreID == input.StoreID { src.DB.Model(&store).Update("owner_id", userID) } else { c.JSON(http.StatusForbidden, gin.H{"error": "Akses Ditolak! Anda bukan pemilik cabang ini."}); return } } else if store.OwnerID != userID { c.JSON(http.StatusForbidden, gin.H{"error": "Akses Ditolak! Anda bukan pemilik cabang ini."}); return }
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": userID, "public_id": user.PublicID, "store_id": store.ID, "plan_type": store.SubscriptionPlan, "role": userRole, "exp": time.Now().Add(time.Hour * 72).Unix()})
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"}); return }
	c.JSON(http.StatusOK, gin.H{"message": "Berhasil masuk ke cabang", "token": tokenString, "store_id": store.ID, "store_name": store.NamaToko, "store_logo": store.LogoURL, "subscription_plan": store.SubscriptionPlan, "role": userRole, "name": user.Name, "foto_url": user.FotoURL})
}

func GetMe(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id"); userID := uint(userIDRaw.(float64)); var user models.User
	if err := src.DB.Preload("Store").First(&user, userID).Error; err != nil { c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan!"}); return }
	c.JSON(http.StatusOK, gin.H{"user_id": user.ID, "public_id": user.PublicID, "name": user.Name, "no_hp": user.NoHP, "nik": user.NIK, "role": user.Role, "is_verified": user.IsVerified, "foto_url": user.FotoURL, "biometric_url": user.BiometricURL, "tempat_lahir": user.TempatLahir, "tanggal_lahir": user.TanggalLahir, "store_name": user.Store.NamaToko, "store_logo": user.Store.LogoURL, "business_type": user.Store.BusinessType, "subscription_plan": user.Store.SubscriptionPlan, "fitur_aktif": user.Store.FiturAktif})
}

func UpdateProfile(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id"); userID := uint(userIDRaw.(float64)); var user models.User
	if err := src.DB.First(&user, userID).Error; err != nil { c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan!"}); return }
	if name := c.PostForm("name"); name != "" { user.Name = name }
	if tempatLahir := c.PostForm("tempat_lahir"); tempatLahir != "" { user.TempatLahir = tempatLahir }
	if tanggalLahir := c.PostForm("tanggal_lahir"); tanggalLahir != "" { user.TanggalLahir = tanggalLahir }
	if noHP := c.PostForm("no_hp"); noHP != "" { user.NoHP = utils.FormatPhoneNumber(noHP) }
	bucketName := os.Getenv("SUPABASE_BUCKET_NAME")
	if fileHeader, err := c.FormFile("foto"); err == nil {
		if fileHeader.Size > 5*1024*1024 { c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran foto profil maksimal 5 MB"}); return }
		file, err := fileHeader.Open(); if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membuka file foto"}); return }
		defer file.Close(); remotePath := fmt.Sprintf("users/%s/profile", user.PublicID)
		publicURL, err := utils.UploadToSupabase(file, fileHeader.Filename, fileHeader.Header.Get("Content-Type"), bucketName, remotePath)
		if err != nil { log.Printf("upload foto gagal: %v", err); c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal upload foto profil"}); return }
		user.FotoURL = publicURL
	}
	if bioHeader, err := c.FormFile("biometric_file"); err == nil {
		if bioHeader.Size > 5*1024*1024 { c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran file biometrik maksimal 5 MB"}); return }
		bioFile, err := bioHeader.Open(); if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membuka file biometrik"}); return }
		defer bioFile.Close(); remotePath := fmt.Sprintf("users/%s/biometric", user.PublicID)
		publicURL, err := utils.UploadToSupabase(bioFile, bioHeader.Filename, bioHeader.Header.Get("Content-Type"), bucketName, remotePath)
		if err != nil { log.Printf("upload biometrik gagal: %v", err); c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal upload file biometrik"}); return }
		user.BiometricURL = publicURL
	}
	if err := src.DB.Save(&user).Error; err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan perubahan"}); return }
	c.JSON(http.StatusOK, gin.H{"message": "Profil berhasil diperbarui!", "data": gin.H{"name": user.Name, "foto_url": user.FotoURL}})
}

func UpdatePassword(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id"); userID := uint(userIDRaw.(float64)); var input struct { OldPassword string `json:"old_password" binding:"required"`; NewPassword string `json:"new_password" binding:"required"` }
	if err := c.ShouldBindJSON(&input); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak lengkap!"}); return }
	var user models.User
	if err := src.DB.First(&user, userID).Error; err != nil { c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan!"}); return }
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Password lama salah!"}); return }
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost); user.Password = string(hashedPassword)
	if err := src.DB.Save(&user).Error; err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan password baru"}); return }
	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil diubah!"})
}

func SendOTPWhatsApp(c *gin.Context) {
	var input struct { Phone string `json:"phone" binding:"required"` }
	if err := c.ShouldBindJSON(&input); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Nomor WhatsApp wajib diisi!"}); return }
	phoneClean := input.Phone; phoneClean = strings.Replace(phoneClean, "+", "", 1)
	if strings.HasPrefix(phoneClean, "0") { phoneClean = "62" + phoneClean[1:] }
	if strings.HasPrefix(phoneClean, "8") { phoneClean = "62" + phoneClean }
	var user models.User
	if err := src.DB.Where("no_hp = ?", phoneClean).First(&user).Error; err != nil { c.JSON(http.StatusNotFound, gin.H{"error": "Nomor WhatsApp tidak terdaftar di sistem kami!"}); return }
	otp := generateOTP()
	if err := src.DB.Model(&user).Updates(map[string]interface{}{"otp_code": otp, "otp_expired": time.Now().Add(time.Minute * 3)}).Error; err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengonfigurasi token verifikasi baru"}); return }
	message := fmt.Sprintf("Halo Bos %s!\n\nKode OTP Verifikasi Akun NEXA POS Anda adalah: *%s*\n\nKode ini rahasia dan berlaku selama 3 menit. Jangan bagikan kode ini kepada siapapun demi keamanan infrastruktur bisnis Anda. 😎", user.Name, otp)
	utils.SendSystemWhatsApp(phoneClean, message)
	c.JSON(http.StatusOK, gin.H{"message": "Kode OTP berhasil dikirim ke WhatsApp Anda! Silakan cek chat masuk.", "phone": phoneClean})
}

func ResetPassword(c *gin.Context) {
	var input ResetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Format tidak valid atau password kurang dari 6 karakter"}); return }
	identifierClean := input.Email
	if !strings.Contains(identifierClean, "@") {
		identifierClean = strings.Replace(identifierClean, "+", "", 1)
		if strings.HasPrefix(identifierClean, "0") { identifierClean = "62" + identifierClean[1:] }
		if strings.HasPrefix(identifierClean, "8") { identifierClean = "62" + identifierClean }
	}
	var user models.User
	if err := src.DB.Where("(email = ? OR no_hp = ?) AND otp_code = ?", identifierClean, identifierClean, input.Token).First(&user).Error; err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Kode OTP tidak valid atau salah!"}); return }
	if time.Now().After(user.OTPExpired) { src.DB.Model(&user).Update("otp_code", ""); c.JSON(http.StatusBadRequest, gin.H{"error": "Kode OTP sudah kedaluwarsa, silakan minta kode baru"}); return }
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses keamanan password"}); return }
	user.Password = string(hashedPassword); user.OTPCode = ""
	if err := src.DB.Save(&user).Error; err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan password baru ke database"}); return }
	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil diperbarui secara otomatis. Silakan login kembali."})
}

func CheckAccount(c *gin.Context) {
	var input struct { Identifier string `json:"identifier" binding:"required"` }
	if err := c.ShouldBindJSON(&input); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Input wajib diisi"})
		return }
	var user models.User; cleanID := utils.FormatPhoneNumber(input.Identifier)
	if err := src.DB.Where("email = ? OR no_hp = ?", input.Identifier, cleanID).First(&user).Error; err != nil { c.JSON(http.StatusNotFound, gin.H{"error": "Akun tidak terdaftar di sistem"}); return }
	email := ""; if user.Email != nil { email = *user.Email }
	c.JSON(http.StatusOK, gin.H{"email": email, "phone": user.NoHP})
}