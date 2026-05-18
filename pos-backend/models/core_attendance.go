package models

import "time"

// Model untuk merekam log absensi harian
type Attendance struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	StoreID    uint      `gorm:"not null" json:"store_id"`
	UserID     uint      `gorm:"not null" json:"user_id"`
	Tanggal    string    `gorm:"type:date;not null" json:"tanggal"` // Format: YYYY-MM-DD
	JamMasuk   string    `gorm:"type:varchar(10)" json:"jam_masuk"` // Format: 08:00:00
	FotoMasuk  string    `gorm:"type:text" json:"foto_masuk"`       // Base64 Text
	JamPulang  string    `gorm:"type:varchar(10)" json:"jam_pulang"` // Format: 17:00:00
	FotoPulang string    `gorm:"type:text" json:"foto_pulang"`      // Base64 Text
	Status     string    `gorm:"type:varchar(20)" json:"status"`    // Hadir, Terlambat, dll
	
	// Relasi ke tabel User untuk mengambil nama dan NIK kasir
	User User `gorm:"foreignKey:UserID" json:"User"`
	
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}