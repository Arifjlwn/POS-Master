package models

import "time"

type User struct {
	// ==========================================
	// PRIMARY IDENTIFIER
	// ==========================================
	ID       uint   `gorm:"primaryKey" json:"id"`
	PublicID string `gorm:"size:26;uniqueIndex" json:"public_id"`

	// ==========================================
	// RELASI TOKO
	// ==========================================
	StoreID *uint `gorm:"index" json:"store_id"`
	Store   Store `gorm:"foreignKey:StoreID" json:"store"`

	// ==========================================
	// DATA AKUN
	// ==========================================
	Name     string  `gorm:"type:varchar(100);not null" json:"name"`
	Email    *string `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	NoHP     string  `gorm:"type:varchar(20);uniqueIndex" json:"no_hp"`
	NIK      *string `gorm:"type:varchar(20)" json:"nik"`

	Password string `gorm:"type:varchar(255);not null" json:"-"`

	Role string `gorm:"type:varchar(20);default:'kasir';index" json:"role"`

	// ==========================================
	// DATA KARYAWAN
	// ==========================================
	TempatLahir  string `gorm:"type:varchar(100)" json:"tempat_lahir"`
	TanggalLahir string `gorm:"type:varchar(20)" json:"tanggal_lahir"`

	FotoURL      string `gorm:"type:text" json:"foto_url"`
	BiometricURL string `gorm:"type:text" json:"biometric_url"`

	// ==========================================
	// VERIFIKASI AKUN
	// ==========================================
	IsVerified bool `gorm:"default:false" json:"is_verified"`

	OTPCode     string     `json:"-"`
	OTPExpired  time.Time  `json:"-"`
	OTPAttempts int        `gorm:"default:0" json:"-"`
	LockedUntil *time.Time `json:"-"`

	// ==========================================
	// TIMESTAMP
	// ==========================================
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}