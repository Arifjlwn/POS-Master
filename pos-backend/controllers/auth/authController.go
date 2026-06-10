package auth

import (
	crand "crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"pos-backend/models"
	src "pos-backend/src/core/config"
	"pos-backend/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func generateOTP() string {
	n, err := crand.Int(crand.Reader, big.NewInt(1000000))
	if err != nil {
		return "000000"
	}
	return fmt.Sprintf("%06d", n.Int64())
}

type RegisterInput struct {
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6"`
	TempatLahir  string `json:"tempat_lahir"`
	TanggalLahir string `json:"tanggal_lahir"`
	NoHP         string `json:"no_hp"`
}
type LoginInput struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}
type SelectStoreInput struct {
	StoreID uint `json:"store_id" binding:"required"`
}
type ResetPasswordInput struct {
	Email    string `json:"email" binding:"required"`
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var existingUser models.User
	email := strings.ToLower(strings.TrimSpace(input.Email))
	cleanPhone := utils.FormatPhoneNumber(input.NoHP)

	// 1. Cek apakah email atau nomor WA sudah pernah di-input sebelumnya 
	if err := src.DB.Where("LOWER(email) = ? OR no_hp = ?", email, cleanPhone).First(&existingUser).Error; err == nil {

		// 🔒 JALUR A: Jika akun sudah terverifikasi murni (Aktif), langsung blokir keras!
		if existingUser.IsVerified {
			c.JSON(http.StatusConflict, gin.H{"error": "Email atau Nomor WhatsApp sudah terdaftar aktif di sistem !"})
			return
		}

		// 🔓 JALUR B: Jika akun masih pending (IsVerified == false, efek klik batal)
		// Kita lakukan OVERWRITE / UPDATE data lamanya biar ga memicu error duplicate key  bantai!
		hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if errHash != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengamankan sandi akun"})
			return
		}

		errUpdate := src.DB.Model(&existingUser).Updates(map[string]interface{}{
			"name":          input.Name,
			"password":      string(hashedPassword),
			"tempat_lahir":  input.TempatLahir,
			"tanggal_lahir": input.TanggalLahir,
			"otp_code":      "",  // Pastikan tetep bersih kosong sebelum dipilih
			"otp_attempts":  0,   // Reset hitungan gagal
			"locked_until":  nil, // Bebaskan gembok brute-force kalau ada
		}).Error

		if errUpdate != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui konfigurasi data pendaftaran tertunda "})
			return
		}

		// Lempar status 200 OK, suruh frontend lari ke halaman Pilih OTP lagi !
		c.JSON(http.StatusOK, gin.H{"message": "Melanjutkan aktivasi akun Anda yang tertunda.", "email": email, "phone": cleanPhone})
		return
	}

	// 2. JALUR C: Pendaftaran Murni Baru (Belum pernah terdata sama sekali di DB )
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengamankan sandi akun"})
		return
	}

	user := models.User{
		PublicID:     utils.GenerateULID(),
		Name:         input.Name,
		Email:        &email,
		Password:     string(hashedPassword),
		Role:         "owner",
		IsVerified:   false,
		OTPCode:      "",
		OTPExpired:   time.Now(),
		TempatLahir:  input.TempatLahir,
		TanggalLahir: input.TanggalLahir,
		NoHP:         cleanPhone,
	}

	if err := src.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat infrastruktur tenant"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Pendaftaran berhasil! Silakan tentukan metode aktivasi akun Anda.", "email": input.Email, "phone": cleanPhone})
}

func SendOTPEmailEndpoint(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email target wajib diisi dengan benar!"})
		return
	}

	var user models.User
	emailClean := strings.ToLower(strings.TrimSpace(input.Email))
	if err := src.DB.Where("email = ?", emailClean).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Alamat email tidak terdaftar di sistem kami!"})
		return
	}

	otp := generateOTP()
	if err := src.DB.Model(&user).Updates(map[string]interface{}{"otp_code": otp, "otp_expired": time.Now().Add(time.Minute * 5)}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui token keamanan"})
		return
	}

	// Menjalankan pengiriman email pihak ketiga secara asynchronous agar response API tetap instan
	go utils.SendOTPEmail(emailClean, otp)

	c.JSON(http.StatusOK, gin.H{"message": "Kode OTP sukses dikirim ke Email Anda! Silakan cek kotak masuk atau spam.", "email": emailClean})
}

func VerifyOTP(c *gin.Context) {
	var input struct {
		Email  string `json:"email" binding:"required"`
		OTP    string `json:"otp" binding:"required"`
		Intent string `json:"intent"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format parameter data tidak valid"})
		return
	}

	var user models.User
	cleanIdentifier := strings.ToLower(strings.TrimSpace(input.Email))

	// FIX QUERY: Presisi pencarian data berdasarkan tipe input identitas
	query := src.DB
	if strings.Contains(cleanIdentifier, "@") {
		query = query.Where("email = ?", cleanIdentifier)
	} else {
		query = query.Where("no_hp = ?", utils.FormatPhoneNumber(cleanIdentifier))
	}

	if err := query.First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Identitas pengguna tidak ditemukan"})
		return
	}

	if user.LockedUntil != nil && user.LockedUntil.Year() == 2099 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akun dibekukan permanen. Hubungi IT Operations pusat."})
		return
	}
	if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
		diff := time.Until(*user.LockedUntil)
		c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Akses terkunci! Silakan coba kembali dalam %d menit.", int(diff.Minutes()))})
		return
	}
	if time.Now().After(user.OTPExpired) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kode OTP telah kedaluwarsa!"})
		return
	}

	if user.OTPCode == "" || user.OTPCode != input.OTP {
		newAttempts := user.OTPAttempts + 1
		updates := map[string]interface{}{"otp_attempts": newAttempts}
		if newAttempts >= 8 {
			updates["locked_until"] = time.Date(2099, 12, 31, 23, 59, 59, 0, time.UTC)
			src.DB.Model(&user).Updates(updates)
			c.JSON(http.StatusForbidden, gin.H{"error": "Batas percobaan terlampaui! Keamanan mendeteksi anomali, akun dikunci permanen."})
			return
		}
		if newAttempts >= 4 {
			updates["locked_until"] = time.Now().Add(time.Hour * 1)
			src.DB.Model(&user).Updates(updates)
			c.JSON(http.StatusForbidden, gin.H{"error": "Terlalu banyak kegagalan! Akses ditangguhkan selama 1 hour."})
			return
		}
		src.DB.Model(&user).Updates(updates)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Kode OTP tidak cocok! Sisa percobaan sebelum dikunci: %d kali.", 4-newAttempts%4)})
		return
	}

	updates := map[string]interface{}{"is_verified": true, "otp_attempts": 0, "locked_until": nil, "otp_code": ""}
	src.DB.Model(&user).Updates(updates)
	c.JSON(http.StatusOK, gin.H{"message": "Verifikasi identitas sukses!"})
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	query := src.DB
	if strings.Contains(input.Identifier, "@") {
		query = query.Where("email = ?", strings.ToLower(strings.TrimSpace(input.Identifier)))
	} else {
		query = query.Where("no_hp = ?", utils.FormatPhoneNumber(input.Identifier))
	}
	if err := query.First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kredensial salah! Pastikan identitas dan kata sandi Anda sesuai."})
		return
	}
	if !user.IsVerified {
		c.JSON(http.StatusForbidden, gin.H{"error": "Infrastruktur akun Anda belum diverifikasi sepenuhnya!", "unverified": true})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kredensial salah! Pastikan identitas dan kata sandi Anda sesuai."})
		return
	}

	var stores []models.Store
	if user.Role == "owner" {
		if user.StoreID != nil {
			src.DB.Where("owner_id = ? OR id = ?", user.ID, *user.StoreID).Find(&stores)
		} else {
			src.DB.Where("owner_id = ?", user.ID).Find(&stores)
		}
	} else {
		if user.StoreID != nil {
			src.DB.Where("id = ?", *user.StoreID).Find(&stores)
		}
	}
	jwtSecret := os.Getenv("JWT_SECRET")

	if len(stores) > 1 || user.Role == "owner" {
		tempToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": user.ID, "public_id": user.PublicID, "role": user.Role, "is_select": true, "exp": time.Now().Add(time.Minute * 15).Unix()})
		tokenString, err := tempToken.SignedString([]byte(jwtSecret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menginisialisasi gerbang sesi"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Otentikasi berhasil, silakan tentukan cabang operasional.", "require_select": true, "token": tokenString, "role": user.Role, "name": user.Name, "stores": stores})
		return
	}

	var storeID uint = 0
	var storeName, storeLogo, planType string
	if len(stores) == 1 {
		storeID = stores[0].ID
		storeName = stores[0].NamaToko
		storeLogo = stores[0].LogoURL
		planType = stores[0].SubscriptionPlan
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": user.ID, "public_id": user.PublicID, "store_id": storeID, "plan_type": planType, "role": user.Role, "is_select": false, "exp": time.Now().Add(time.Hour * 72).Unix()})
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal merilis token otoritas sesi"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Login Berhasil!", "require_select": false, "token": tokenString, "role": user.Role, "name": user.Name, "has_setup_store": storeID != 0, "store_name": storeName, "store_logo": storeLogo, "subscription_plan": planType})
}

func SelectStore(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi kedaluwarsa, silakan masuk kembali."})
		return
	}
	userID := uint(userIDRaw.(float64))
	userRole := c.GetString("role")
	var input SelectStoreInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cabang sasaran tidak valid"})
		return
	}
	var user models.User
	if err := src.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data user tidak sinkron"})
		return
	}
	var store models.Store
	if err := src.DB.First(&store, input.StoreID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Infrastruktur cabang tidak ditemukan"})
		return
	}

	if userRole == "owner" {
		if store.OwnerID != user.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Pelanggaran Otoritas! Anda tidak memiliki hak izin akses pada gerai cabang ini."})
			return
		}
	} else {
		if user.StoreID == nil || *user.StoreID != store.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak! Akun staf Anda tidak terdaftar di cabang gerai ini."})
			return
		}
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": user.ID, "public_id": user.PublicID, "store_id": store.ID, "plan_type": store.SubscriptionPlan, "role": user.Role, "is_select": false, "exp": time.Now().Add(time.Hour * 72).Unix()})
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menerbitkan token sesi cabang"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Berhasil beralih ke gerai cabang.", "token": tokenString, "store_id": store.ID, "store_name": store.NamaToko, "store_logo": store.LogoURL, "subscription_plan": store.SubscriptionPlan, "role": user.Role, "name": user.Name, "foto_url": user.FotoURL})
}

func GetMe(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := uint(userIDRaw.(float64))

	var user models.User
	if err := src.DB.Preload("Store").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan!"})
		return
	}

	// 🚀 STEP PENYELAMAT: Query list seluruh toko milik owner ini dari database
	var listStores []models.Store
	if user.Role == "owner" {
		// Kalau dia owner, tarik semua toko yang owner_id-nya adalah dia 
		src.DB.Where("owner_id = ?", user.ID).Find(&listStores)
	} else {
		// Kalau staf/kasir, masukin aja ruko tempat dia bekerja biar gak kosong 
		if user.StoreID != nil {
			src.DB.Where("id = ?", *user.StoreID).Find(&listStores)
		}
	}

	// Balikin semua field asli milik lu, tapi kita selipin array "stores" di paling bawah!
	c.JSON(http.StatusOK, gin.H{
		"user_id":           user.ID,
		"public_id":         user.PublicID,
		"name":              user.Name,
		"no_hp":             user.NoHP,
		"nik":               user.NIK,
		"role":              user.Role,
		"is_verified":       user.IsVerified,
		"foto_url":          user.FotoURL,
		"biometric_url":     user.BiometricURL,
		"tempat_lahir":      user.TempatLahir,
		"tanggal_lahir":     user.TanggalLahir,
		"store_name":        user.Store.NamaToko,
		"store_logo":        user.Store.LogoURL,
		"business_type":     user.Store.BusinessType,
		"subscription_plan": user.Store.SubscriptionPlan,
		"fitur_aktif":       user.Store.FiturAktif,

		// 🔒 SEKAT SAKTI: Ini yang ditunggu-tunggu sama SelectStore.vue lu !
		"stores": listStores,
	})
}

func UpdateProfile(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := uint(userIDRaw.(float64))
	var user models.User
	if err := src.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan!"})
		return
	}
	if name := c.PostForm("name"); name != "" {
		user.Name = name
	}
	if tempatLahir := c.PostForm("tempat_lahir"); tempatLahir != "" {
		user.TempatLahir = tempatLahir
	}
	if tanggalLahir := c.PostForm("tanggal_lahir"); tanggalLahir != "" {
		user.TanggalLahir = tanggalLahir
	}

	bucketName := os.Getenv("SUPABASE_BUCKET_NAME")
	if fileHeader, err := c.FormFile("foto"); err == nil {
		if fileHeader.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran foto profil maksimal 5 MB"})
			return
		}
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membuka file foto"})
			return
		}
		defer file.Close()
		remotePath := fmt.Sprintf("users/%s/profile", user.PublicID)
		publicURL, err := utils.UploadToSupabase(file, fileHeader.Filename, fileHeader.Header.Get("Content-Type"), bucketName, remotePath)
		if err != nil {
			log.Printf("upload foto gagal: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal upload foto profil"})
			return
		}
		user.FotoURL = publicURL
	}
	if bioHeader, err := c.FormFile("biometric_file"); err == nil {
		if bioHeader.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran file biometrik maksimal 5 MB"})
			return
		}
		bioFile, err := bioHeader.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membuka file biometrik"})
			return
		}
		defer bioFile.Close()
		remotePath := fmt.Sprintf("users/%s/biometric", user.PublicID)
		publicURL, err := utils.UploadToSupabase(bioFile, bioHeader.Filename, bioHeader.Header.Get("Content-Type"), bucketName, remotePath)
		if err != nil {
			log.Printf("upload biometrik gagal: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal upload file biometrik"})
			return
		}
		user.BiometricURL = publicURL
	}
	if err := src.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan perubahan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Profil berhasil diperbarui!", "data": gin.H{"name": user.Name, "foto_url": user.FotoURL}})
}

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

func SendOTPWhatsApp(c *gin.Context) {
	var input struct {
		Phone string `json:"phone" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nomor WhatsApp wajib diisi!"})
		return
	}
	phoneClean := input.Phone
	phoneClean = strings.Replace(phoneClean, "+", "", 1)
	if strings.HasPrefix(phoneClean, "0") {
		phoneClean = "62" + phoneClean[1:]
	}
	if strings.HasPrefix(phoneClean, "8") {
		phoneClean = "62" + phoneClean
	}
	var user models.User
	if err := src.DB.Where("no_hp = ?", phoneClean).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nomor WhatsApp tidak terdaftar di sistem kami!"})
		return
	}
	otp := generateOTP()
	if err := src.DB.Model(&user).Updates(map[string]interface{}{"otp_code": otp, "otp_expired": time.Now().Add(time.Minute * 3)}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengonfigurasi token verifikasi baru"})
		return
	}
	message := fmt.Sprintf("Halo Bos %s!\n\nKode OTP Verifikasi Akun ARZURA POS Anda adalah: *%s*\n\nKode ini rahasia dan berlaku selama 3 menit. Jangan bagikan kode ini kepada siapapun demi keamanan infrastruktur bisnis Anda. 😎", user.Name, otp)
	utils.SendSystemWhatsApp(phoneClean, message)
	c.JSON(http.StatusOK, gin.H{"message": "Kode OTP berhasil dikirim ke WhatsApp Anda! Silakan cek chat masuk.", "phone": phoneClean})
}

func ResetPassword(c *gin.Context) {
	var input ResetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tidak valid atau password kurang dari 6 karakter"})
		return
	}
	identifierClean := input.Email
	query := src.DB
	if strings.Contains(identifierClean, "@") {
		query = query.Where("email = ?", strings.ToLower(strings.TrimSpace(identifierClean)))
	} else {
		identifierClean = strings.Replace(identifierClean, "+", "", 1)
		if strings.HasPrefix(identifierClean, "0") {
			identifierClean = "62" + identifierClean[1:]
		}
		if strings.HasPrefix(identifierClean, "8") {
			identifierClean = "62" + identifierClean
		}
		query = query.Where("no_hp = ?", identifierClean)
	}

	var user models.User
	if err := query.First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Identitas pengguna tidak ditemukan"})
		return
	}

	// 🚀 FIX SECURITY: Blokir brute-force OTP pada gerbang ResetPassword
	if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
		diff := time.Until(*user.LockedUntil)
		c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Akses terkunci! Silakan coba kembali dalam %d menit.", int(diff.Minutes()))})
		return
	}

	if time.Now().After(user.OTPExpired) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kode OTP sudah kedaluwarsa, silakan minta kode baru"})
		return
	}

	if user.OTPCode == "" || user.OTPCode != input.Token {
		newAttempts := user.OTPAttempts + 1
		updates := map[string]interface{}{"otp_attempts": newAttempts}
		if newAttempts >= 5 {
			updates["locked_until"] = time.Now().Add(time.Hour * 2) // Kunci 2 jam jika brute-force token reset sandi
			src.DB.Model(&user).Updates(updates)
			c.JSON(http.StatusForbidden, gin.H{"error": "Terlalu banyak kegagalan verifikasi! Akses pemulihan ditangguhkan selama 2 jam."})
			return
		}
		src.DB.Model(&user).Updates(updates)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Kode verifikasi OTP tidak cocok! Sisa percobaan: %d kali.", 5-newAttempts)})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses keamanan password"})
		return
	}

	// Sukses memperbarui data, reset counter percobaan ke nol
	if err := src.DB.Model(&user).Updates(map[string]interface{}{
		"password":     string(hashedPassword),
		"otp_code":     "",
		"otp_attempts": 0,
		"locked_until": nil,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan password baru ke database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil diperbarui secara otomatis. Silakan login kembali."})
}

func CheckAccount(c *gin.Context) {
	var input struct {
		Identifier string `json:"identifier" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input wajib diisi"})
		return
	}
	var user models.User
	cleanID := utils.FormatPhoneNumber(input.Identifier)
	if err := src.DB.Where("email = ? OR no_hp = ?", input.Identifier, cleanID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Akun tidak terdaftar di sistem"})
		return
	}
	email := ""
	if user.Email != nil {
		email = *user.Email
	}
	c.JSON(http.StatusOK, gin.H{"email": email, "phone": user.NoHP})
}
