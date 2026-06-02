package models

import "time"

type ShiftClosing struct {
    ID              uint      `gorm:"primaryKey" json:"id"`
    StoreID         uint      `gorm:"not null" json:"store_id"`
    SessionID       uint      `gorm:"not null" json:"session_id"`
    UserID          uint      `gorm:"not null" json:"user_id"`
    
    // Summary Data (Ini yang lu print di struk!)
    StartTime       time.Time `json:"start_time"`
    EndTime         time.Time `json:"end_time"`
    NetSales        float64   `gorm:"type:decimal(15,2)" json:"net_sales"`
    TotalTax        float64   `gorm:"type:decimal(15,2)" json:"total_tax"`
    SalesCash       float64   `gorm:"type:decimal(15,2)" json:"sales_cash"`
    SalesNonTunai   float64   `gorm:"type:decimal(15,2)" json:"sales_non_tunai"`
    
    // Cash Count (Hasil hitung fisik laci)
    TotalExpected   float64   `gorm:"type:decimal(15,2)" json:"total_expected"`
    TotalActual     float64   `gorm:"type:decimal(15,2)" json:"total_actual"`
    Selisih         float64   `gorm:"type:decimal(15,2)" json:"selisih"`

    // Relasi biar bisa ditarik namanya
    User            User           `gorm:"foreignKey:UserID" json:"User"`
    Session         CashierSession `gorm:"foreignKey:SessionID" json:"session"`
    Store           Store          `gorm:"foreignKey:StoreID" json:"Store"` // 🚀 INI BARU DITAMBAHIN
    CreatedAt       time.Time      `json:"created_at"`
}