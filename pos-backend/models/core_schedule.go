package models

import "time"

type Schedule struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	PublicID string `gorm:"size:26;uniqueIndex" json:"public_id"`

	StoreID uint `gorm:"not null;index:idx_schedule_store_date" json:"store_id"`
	UserID  uint `gorm:"not null;index" json:"user_id"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	// Format YYYY-MM-DD
	Tanggal string `gorm:"type:date;not null;index:idx_schedule_store_date" json:"tanggal"`

	// SHIFT_1, SHIFT_2, MIDDLE, OFF
	ShiftType string `gorm:"type:varchar(20);not null" json:"shift_type"`

	JamMasukJadwal  string `gorm:"type:varchar(5)" json:"jam_masuk_jadwal"`
	JamPulangJadwal string `gorm:"type:varchar(5)" json:"jam_pulang_jadwal"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Store Store `gorm:"foreignKey:StoreID" json:"-"`
}