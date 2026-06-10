package utils // 🚀 Menggunakan kasta package utils lu

import (
	"log"
	"pos-backend/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedSuperAdmin buat nge-inject akun Root Admin pertama kali secara aman
func SeedSuperAdmin(db *gorm.DB) {
	adminEmail := "founder@arzura-pos.my.id"
	adminPassword := "Yourbitch123"

	var count int64
	// Cek apakah akun admin ini udah terdaftar atau belum
	db.Model(&models.User{}).Where("email = ?", adminEmail).Count(&count)
	if count > 0 {
		log.Println("[SEEDER] Akun Root Admin sudah ada, skip injection.")
		return
	}

	// Hash password menggunakan Bcrypt tingkat tinggi
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("[SEEDER] Gagal melakukan enkripsi password : %v", err)
	}

	// 🔒 HYBRID ARCHITECTURE COMPLIANCE
	admin := models.User{
		Name:     "Arif Juliawan (Founder)",
		Email:    &adminEmail,
		Password: string(hashedPassword),
		Role:     "super_admin",
		PublicID: GenerateULID(), // ◄ Langsung panggil tanpa utils. karena udah satu package !
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Printf("[SEEDER] Gagal meng-inject akun admin ke DB : %v", err)
	} else {
		log.Println("[SEEDER] SUCCESS! Akun Root Admin ARZURA POS berhasil ditanam ke database pusat! 🔥")
	}
}
