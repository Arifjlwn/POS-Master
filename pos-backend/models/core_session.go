package models

import (
	"time"
)

type CashierSession struct {
    ID            uint       `gorm:"primaryKey" json:"id"`
    StoreID       uint       `gorm:"not null" json:"store_id"`
    UserID        uint       `gorm:"not null" json:"user_id"`
    StationNumber string     `gorm:"size:10;not null" json:"station_number"`
    ModalAwal     float64    `gorm:"type:decimal(15,2);not null" json:"modal_awal"`
    TotalMasuk    float64    `gorm:"type:decimal(15,2);default:0" json:"total_masuk"`
    TotalKeluar   float64    `gorm:"type:decimal(15,2);default:0" json:"total_keluar"`
    
    // 🚀 TAMBAHKAN DUA INI UNTUK CLOSING
    TotalAktual   float64    `gorm:"type:decimal(15,2);default:0" json:"total_aktual"`
    Selisih       float64    `gorm:"type:decimal(15,2);default:0" json:"selisih"`
    
    StartTime     time.Time  `json:"start_time"`
    EndTime       *time.Time `json:"end_time"`
    Status        string     `gorm:"size:20;default:'open'" json:"status"`
    
    User  User  `gorm:"foreignKey:UserID" json:"user"`
    Store Store `gorm:"foreignKey:StoreID" json:"store"`
}