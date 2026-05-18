package models

import "time"

type Schedule struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	StoreID          uint      `gorm:"not null" json:"store_id"`
	UserID           uint      `gorm:"not null" json:"user_id"`
	User             User      `gorm:"foreignKey:UserID" json:"user,omitempty"` // Biar bisa ambil nama karyawan pas ditarik ke Vue
	Tanggal          string    `gorm:"type:date;not null" json:"tanggal"`       // Format: "YYYY-MM-DD"
	ShiftType        string    `gorm:"type:varchar(20);not null" json:"shift_type"` // 'Shift 1', 'Shift 2', 'Middle', atau 'OFF'
	JamMasukJadwal   string    `gorm:"type:varchar(5)" json:"jam_masuk_jadwal"`  // Contoh: "07:00"
	JamPulangJadwal  string    `gorm:"type:varchar(5)" json:"jam_pulang_jadwal"` // Contoh: "15:00"
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}