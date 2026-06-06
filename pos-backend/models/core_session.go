package models

import "time"

type CashierSession struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	PublicID string `gorm:"size:26;uniqueIndex" json:"public_id"` // ULID / UUID Aman

	StoreID uint `gorm:"not null;index" json:"store_id"`
	
	// FIX AUDIT FINANSIAL:UserID tetap menjadi pemilik utama laci kasir shift tersebut
	UserID uint `gorm:"not null;index" json:"user_id"`

	// 🔒 GEMBOK PROTEKSI POS MULTI-TENANT (PARTIAL UNIQUE INDEX POSTGRESQL):
	// Taktik sakti via tag GORM gorm:"uniqueIndex:idx_one_active_session_per_user..." 
	// Memastikan satu user/stasiun di toko lu GAK BISA punya 2 sesi 'OPEN' sekaligus. 
	// Kalau kasir maksa nembak via Postman/Bypass FE, database otomatis nolak (Duplicate Error)!
	StationNumber string `gorm:"size:20;uniqueIndex:idx_one_active_station_per_store,where:status = 'OPEN';not null" json:"station_number"`

	// Modal Awal Shift (Floating Capital)
	ModalAwal float64 `gorm:"type:decimal(15,2);not null;default:0" json:"modal_awal"`

	// Rekap Arus Kas Riil di POS
	TotalMasuk  float64 `gorm:"type:decimal(15,2);default:0" json:"total_masuk"`  // Penjualan Tunai dll
	TotalKeluar float64 `gorm:"type:decimal(15,2);default:0" json:"total_keluar"` // Retur / Paid Out

	// Closing Audit Data
	TotalAktual float64 `gorm:"type:decimal(15,2);default:0" json:"total_aktual"` // Uang fisik di laci kasir nyatanya 
	Selisih     float64 `gorm:"type:decimal(15,2);default:0" json:"selisih"`      // TotalAktual - (ModalAwal + TotalMasuk - TotalKeluar)
	ClosingNote string  `gorm:"type:text" json:"closing_note"`

	// Status Manajemen Sesi Laci Kasir
	// OPEN, CLOSED, CANCELLED
	Status string `gorm:"size:20;uniqueIndex:idx_one_active_session_per_user,where:status = 'OPEN';default:'OPEN'" json:"status"`

	OpenedAt time.Time  `gorm:"index" json:"opened_at"`
	ClosedAt *time.Time `gorm:"index" json:"closed_at"`

	// 🛡️ SUNTIKAN SAKTI SAKSI FORENSIK: Rekam jejak user ID siapa yang membuka dan menutup sesi 
	// (Penting banget kalau ada kasus selisih kas raib jutaan rupiah)
	OpenedBy uint `gorm:"not null;index" json:"opened_by"`
	ClosedBy *uint `gorm:"index" json:"closed_by"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations Mapping
	User   User   `gorm:"foreignKey:UserID" json:"user"`
	Store  Store  `gorm:"foreignKey:StoreID" json:"store"`
	
	// Relasi virtual tambahan untuk tracking saksi audit
	Opener User   `gorm:"foreignKey:OpenedBy" json:"-"`
	Closer *User  `gorm:"foreignKey:ClosedBy" json:"-"`

	Transactions []Transaction `gorm:"foreignKey:SessionID" json:"-"`
}