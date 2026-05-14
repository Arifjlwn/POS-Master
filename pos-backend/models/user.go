package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	StoreID   *uint     `json:"store_id"` 
	
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Email     *string   `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	NIK       *string   `gorm:"type:varchar(20);uniqueIndex" json:"nik"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	Role      string    `gorm:"type:varchar(20);default:'kasir'" json:"role"`

	// --- 🚀 FIELD BARU UNTUK KARYAWAN & ABSENSI WAJAH ---
	TempatLahir  string    `gorm:"type:varchar(100)" json:"tempat_lahir"`
	TanggalLahir string    `gorm:"type:varchar(20)" json:"tanggal_lahir"`
	NoHP         string    `gorm:"type:varchar(20)" json:"no_hp"`
	FotoURL      string    `gorm:"type:text" json:"foto_url"` 
	// ----------------------------------------------------

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Store     Store     `gorm:"foreignKey:StoreID" json:"store"`

	// FIELD Verifikasi OTP Email
	IsVerified bool `gorm:"default:false" json:"is_verified"`
	OTPCode string `json:"-"`
	OTPExpired time.Time `json:"-"`
}