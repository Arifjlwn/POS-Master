package auth

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"pos-backend/models"
	src "pos-backend/src/core/config"
	"pos-backend/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// 1. Struct Register
type RegisterInput struct {
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6"`
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
		Name:         input.Name,
		Email:        &input.Email,
		Password:     string(hashedPassword),
		Role:         "owner",
		IsVerified:   false,
		OTPCode:      otp,
		OTPExpired:   time.Now().Add(time.Minute * 5),
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
		Email  string `json:"email" binding:"required"` // Bisa email atau phone clean
		OTP    string `json:"otp" binding:"required"`
		Intent string `json:"intent"` // 🚀 Tangkap intent dari Vue!
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	var user models.User
	// Cari akun pake email atau nomor hp
	if err := src.DB.Where("email = ? OR no_hp = ?", input.Email, input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	// 🛑 1. CEK STATUS LOCKOUT PERMANEN
	if user.LockedUntil != nil && user.LockedUntil.Year() == 2099 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akun Anda telah DI-LOCK PERMANEN karena tindakan mencurigakan. Hubungi Tim IT!"})
		return
	}

	// 🛑 2. CEK STATUS LOCKOUT TIMEOUT (1 JAM)
	if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
		diff := time.Until(*user.LockedUntil)
		c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Akses dibekukan! Silakan tunggu %d menit lagi.", int(diff.Minutes()))})
		return
	}

	// 🛑 3. CEK KADALUARSA OTP
	if time.Now().After(user.OTPExpired) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kode OTP sudah kadaluarsa!"})
		return
	}

	// 🛑 4. KALO OTP-NYA SALAH, NAIKKIN COUNT ATTEMPTS!
	if user.OTPCode != input.OTP {
		newAttempts := user.OTPAttempts + 1
		updates := map[string]interface{}{"otp_attempts": newAttempts}

		if newAttempts >= 8 {
			permanentLock := time.Date(2099, 12, 31, 23, 59, 59, 0, time.UTC)
			updates["locked_until"] = permanentLock
			src.DB.Model(&user).Updates(updates)
			c.JSON(http.StatusForbidden, gin.H{"error": "Terlalu banyak percobaan! Akun Anda kini DI-LOCK PERMANEN."})
			return
		} else if newAttempts >= 4 {
			lockTime := time.Now().Add(time.Hour * 1)
			updates["locked_until"] = lockTime
			src.DB.Model(&user).Updates(updates)
			c.JSON(http.StatusForbidden, gin.H{"error": "Terlalu banyak percobaan salah! Akses dibekukan selama 1 jam."})
			return
		}

		src.DB.Model(&user).Updates(updates)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Kode OTP salah! Sisa percobaan: %d kali lagi.", 4-newAttempts%4)})
		return
	}

	// 🎉 KALO LOLOS / BENER
	updates := map[string]interface{}{
		"is_verified":  true,
		"otp_attempts": 0,
		"locked_until": nil,
	}

	// 🚀 JANGAN HAPUS OTP KALAU LAGI RESET PASSWORD! (Biar bisa diverifikasi lagi di step reset)
	if input.Intent != "reset-password" {
		updates["otp_code"] = ""
	}

	src.DB.Model(&user).Updates(updates)
	c.JSON(http.StatusOK, gin.H{"message": "Verifikasi sukses!"})
}

// 2. Struct Login
type LoginInput struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

// -- LOGIN DENGAN PROTEKSI HYBRID (MULTI-OUTLET SUPPORT) --
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	query := src.DB // Jangan preload Store dulu di sini, karena tokonya bisa banyak

	if strings.Contains(input.Identifier, "@") {
		query = query.Where("email = ?", input.Identifier)
	} else {
		cleanHP := utils.FormatPhoneNumber(input.Identifier)
		query = query.Where("no_hp = ?", cleanHP)
	}

	if err := query.First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Identitas tidak ditemukan! Pastikan Email atau No. WhatsApp benar."})
		return
	}

	if strings.Contains(input.Identifier, "@") && !user.IsVerified {
		c.JSON(http.StatusForbidden, gin.H{
			"error":      "Email Anda belum diverifikasi!",
			"unverified": true,
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password salah!"})
		return
	}

	// 🚀 CARI SEMUA CABANG YANG BISA DIAKSES OLEH USER INI
	var stores []models.Store
	if user.Role == "owner" {
		// JALUR BYPASS: Cari toko yang owner_id-nya adalah dia,
		// ATAU toko lawas (cabang pertama) yang ID-nya nempel di profil user.
		if user.StoreID != nil {
			src.DB.Where("owner_id = ? OR id = ?", user.ID, *user.StoreID).Find(&stores)
		} else {
			src.DB.Where("owner_id = ?", user.ID).Find(&stores)
		}
	} else {
		// Kasir/Manager tetap nempel di 1 toko
		if user.StoreID != nil {
			src.DB.Where("id = ?", *user.StoreID).Find(&stores)
		}
	}

	// 🚀 LOGIKA PEMBERIAN TOKEN
	// JIKA TOKONYA BANYAK ATAU DIA OWNER, KASIH "TEMPORARY TOKEN" BUAT MILIH TOKO
	if len(stores) > 1 || user.Role == "owner" {
		// Bikin Temporary Token (Cuma tahan 15 Menit, gak bisa buat akses POS/Produk)
		tempToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id":   user.ID,
			"role":      user.Role,
			"is_select": true, // Penanda ini token cuma buat milih toko
			"exp":       time.Now().Add(time.Minute * 15).Unix(),
		})

		tokenString, _ := tempToken.SignedString([]byte("KUNCI_RAHASIA_SUPER_KUAT_123"))

		c.JSON(http.StatusOK, gin.H{
			"message":        "Login berhasil, silakan pilih cabang.",
			"require_select": true,        // Frontend harus nangkep ini buat redirect ke halaman pilih toko
			"token":          tokenString, // Simpen sementara di memory frontend
			"role":           user.Role,
			"name":           user.Name,
			"stores":         stores, // Kirim list toko ke frontend buat ditampilin di Card
		})
		return
	}

	// JIKA TOKONYA CUMA 1 (Biasanya Kasir Baru), LANGSUNG KASIH TOKEN FINAL
	var storeID uint = 0
	var storeName, storeLogo, planType string

	if len(stores) == 1 {
		storeID = stores[0].ID
		storeName = stores[0].NamaToko
		storeLogo = stores[0].LogoURL
		planType = stores[0].SubscriptionPlan
	}

	// Tiket JWT Final (Bisa akses semua fitur)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"store_id":  storeID,
		"plan_type": planType,
		"role":      user.Role,
		"exp":       time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte("KUNCI_RAHASIA_SUPER_KUAT_123"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencetak token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":           "Login Sukses!",
		"require_select":    false,
		"token":             tokenString,
		"role":              user.Role,
		"name":              user.Name,
		"has_setup_store":   storeID != 0,
		"store_name":        storeName,
		"store_logo":        storeLogo,
		"subscription_plan": planType,
	})
}

// 🚀 SAKLAR FINAL: FUNGSI BUAT MILIH CABANG & DAPETIN TOKEN ASLI
type SelectStoreInput struct {
	StoreID uint `json:"store_id" binding:"required"`
}

func SelectStore(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak valid, silakan login ulang."})
		return
	}
	userID := uint(userIDRaw.(float64))
	userRole := c.GetString("role")

	var input SelectStoreInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pilihan cabang tidak valid"})
		return
	}

	// 🚀 PINDAH KE SINI: Tarik data User di luar blok IF biar kebaca sampai bawah!
	var user models.User
	if err := src.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data user tidak ditemukan"})
		return
	}

	var store models.Store
	if userRole == "owner" {
		// CARA AUTO-REPAIR: Cari toko berdasarkan ID yang diminta.
		if err := src.DB.First(&store, input.StoreID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cabang tidak ditemukan"})
			return
		}

		// Kalo ternyata toko ini owner_id-nya masih kosong, TAPI ada user yang punya store_id ini (yakni user lu)
		// Maka klaim toko ini sebagai milik lu!
		if store.OwnerID == 0 {
			if user.StoreID != nil && *user.StoreID == input.StoreID {
				// JALUR BYPASS SUKSES: Benerin database secara gaib!
				src.DB.Model(&store).Update("owner_id", userID)
			} else {
				c.JSON(http.StatusForbidden, gin.H{"error": "Akses Ditolak! Anda bukan pemilik cabang ini."})
				return
			}
		} else if store.OwnerID != userID {
			// Tokonya ada owner_id-nya, tapi BUKAN punya lu!
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses Ditolak! Anda bukan pemilik cabang ini."})
			return
		}

	} else {
		// Jalur Kasir
		if err := src.DB.First(&store, input.StoreID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cabang tidak ditemukan"})
			return
		}
	}

	// 🚀 CETAK TOKEN FINAL
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   userID,
		"store_id":  store.ID,
		"plan_type": store.SubscriptionPlan,
		"role":      userRole,
		"exp":       time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, _ := token.SignedString([]byte("KUNCI_RAHASIA_SUPER_KUAT_123"))

	c.JSON(http.StatusOK, gin.H{
		"message":           "Berhasil masuk ke cabang",
		"token":             tokenString,
		"store_id":          store.ID,
		"store_name":        store.NamaToko,
		"store_logo":        store.LogoURL,
		"subscription_plan": store.SubscriptionPlan,
		"role":              userRole,
		"name":              user.Name,    // SEKARANG SUDAH DIKENALI GOLANG
		"foto_url":          user.FotoURL, // SEKARANG SUDAH DIKENALI GOLANG
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
		"user_id":           user.ID,
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
	if name := c.PostForm("name"); name != "" {
		user.Name = name
	}
	if tempatLahir := c.PostForm("tempat_lahir"); tempatLahir != "" {
		user.TempatLahir = tempatLahir
	}
	if tanggalLahir := c.PostForm("tanggal_lahir"); tanggalLahir != "" {
		user.TanggalLahir = tanggalLahir
	}

	// 🚀 FORMAT OTOMATIS JADI 628xxx BIAR RAPI!
	if noHP := c.PostForm("no_hp"); noHP != "" {
		user.NoHP = utils.FormatPhoneNumber(noHP)
	}

	// Buat prefix nama file dari NIK
	nikClean := "user"
	if user.NIK != nil && *user.NIK != "" {
		nikClean = *user.NIK
	}

	// 2. Update Foto Profil
	if file, err := c.FormFile("foto"); err == nil {
		newFileName := fmt.Sprintf("%s_profil_%d%s", nikClean, time.Now().Unix(), filepath.Ext(file.Filename))
		uploadPath := filepath.Join("uploads", newFileName)
		if err := c.SaveUploadedFile(file, uploadPath); err == nil {
			if user.FotoURL != "" {
				os.Remove("." + user.FotoURL)
			} // Hapus foto lama
			user.FotoURL = "/uploads/" + newFileName
		}
	}

	// 3. Update Foto Biometrik
	if bioFile, errBio := c.FormFile("biometric_file"); errBio == nil {
		newBioName := fmt.Sprintf("%s_bio_%d%s", nikClean, time.Now().Unix(), filepath.Ext(bioFile.Filename))
		uploadBioPath := filepath.Join("uploads", newBioName)
		if err := c.SaveUploadedFile(bioFile, uploadBioPath); err == nil {
			if user.BiometricURL != "" {
				os.Remove("." + user.BiometricURL)
			} // Hapus bio lama
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

// 3. Struct Payload Reset Password
type ResetPasswordInput struct {
	Email    string `json:"email" binding:"required"` // Bisa Email atau No. WA
	Token    string `json:"token" binding:"required"` // OTP
	Password string `json:"password" binding:"required,min=6"`
}

// -- 🚀 RESET PASSWORD & GANTI BARU --
func ResetPassword(c *gin.Context) {
	var input ResetPasswordInput

	// 1. Validasi input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tidak valid atau password kurang dari 6 karakter"})
		return
	}

	// 2. Bersihkan inputan kalau ternyata itu Nomor HP (biar formatnya selalu 628xxx)
	identifierClean := input.Email
	if !strings.Contains(identifierClean, "@") {
		identifierClean = strings.Replace(identifierClean, "+", "", 1)
		if strings.HasPrefix(identifierClean, "0") {
			identifierClean = "62" + identifierClean[1:]
		}
		if strings.HasPrefix(identifierClean, "8") {
			identifierClean = "62" + identifierClean
		}
	}

	var user models.User

	// 3. Cari User berdasarkan (Email ATAU No HP) DAN OTP Code yang cocok
	err := src.DB.Where("(email = ? OR no_hp = ?) AND otp_code = ?", identifierClean, identifierClean, input.Token).First(&user).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kode OTP tidak valid atau salah!"})
		return
	}

	// 4. Cek apakah OTP sudah kadaluarsa
	if time.Now().After(user.OTPExpired) {
		// Opsional: Kosongkan OTP yang kadaluarsa
		src.DB.Model(&user).Update("otp_code", "")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kode OTP sudah kedaluwarsa, silakan minta kode baru"})
		return
	}

	// 5. Enkripsi Password Baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses keamanan password"})
		return
	}

	// 6. Update Password & Hanguskan OTP
	user.Password = string(hashedPassword)
	user.OTPCode = "" // Hanguskan OTP

	if err := src.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan password baru ke database"})
		return
	}

	// 7. Respon Sukses!
	c.JSON(http.StatusOK, gin.H{
		"message": "Password berhasil diperbarui secara otomatis. Silakan login kembali.",
	})
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
	// Cari pake email atau nomor hp yang udah di-format clean
	cleanID := utils.FormatPhoneNumber(input.Identifier)

	if err := src.DB.Where("email = ? OR no_hp = ?", input.Identifier, cleanID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Akun tidak terdaftar di sistem"})
		return
	}

	// Kembalin data asli ke Vue
	c.JSON(http.StatusOK, gin.H{
		"email": *user.Email,
		"phone": user.NoHP,
	})
}
