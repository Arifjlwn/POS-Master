package models

import "time"

type Attendance struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	PublicID string `gorm:"size:26;uniqueIndex" json:"public_id"`

	StoreID uint `gorm:"not null;index:idx_attendance_store_date" json:"store_id"`
	UserID  uint `gorm:"not null;index" json:"user_id"`

	// Format: YYYY-MM-DD
	Tanggal string `gorm:"type:date;not null;index:idx_attendance_store_date" json:"tanggal"`

	// Format: 08:00:00
	JamMasuk string `gorm:"type:varchar(10)" json:"jam_masuk"`

	// URL Supabase Storage
	FotoMasuk string `gorm:"type:text" json:"foto_masuk"`

	// Format: 17:00:00
	JamPulang string `gorm:"type:varchar(10)" json:"jam_pulang"`

	// URL Supabase Storage
	FotoPulang string `gorm:"type:text" json:"foto_pulang"`

	// HADIR, TERLAMBAT, IZIN, SAKIT, ALFA
	Status string `gorm:"type:varchar(20);default:'HADIR'" json:"status"`

	// Opsional buat payroll nanti
	Catatan string `gorm:"type:text" json:"catatan"`

	User  User  `gorm:"foreignKey:UserID" json:"user"`
	Store Store `gorm:"foreignKey:StoreID" json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}