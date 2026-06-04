package models

import "time"

type Store struct {
	// ==========================================
	// PRIMARY IDENTIFIER
	// ==========================================
	ID       uint   `gorm:"primaryKey" json:"id"`
	PublicID string `gorm:"size:26;uniqueIndex;not null" json:"public_id"`

	// ==========================================
	// OWNER
	// ==========================================
	OwnerID uint `gorm:"not null;index" json:"owner_id"`

	// ==========================================
	// INFORMASI TOKO
	// ==========================================
	NamaToko string `gorm:"type:varchar(100);not null;index" json:"nama_toko"`
	Telepon  string `gorm:"type:varchar(20);index" json:"telepon"`

	BusinessType string `gorm:"type:varchar(100);not null" json:"business_type"`

	Industry string `gorm:"type:varchar(50);default:'retail';index" json:"industry"`

	// ==========================================
	// LANGGANAN
	// ==========================================
	SubscriptionPlan string `gorm:"type:varchar(50);default:'trial'" json:"subscription_plan"`

	SubscriptionStatus string `gorm:"type:varchar(20);default:'active';index" json:"subscription_status"`

	SubscriptionEnd *time.Time `json:"subscription_end"`

	FiturAktif string `gorm:"type:text" json:"fitur_aktif"`

	QuotaTerminal int `gorm:"default:1;not null" json:"quota_terminal"`

	// ==========================================
	// ALAMAT TOKO
	// ==========================================
	Alamat    string `gorm:"type:text" json:"alamat"`
	Provinsi  string `gorm:"type:varchar(100)" json:"provinsi"`
	Kota      string `gorm:"type:varchar(100)" json:"kota"`
	Kecamatan string `gorm:"type:varchar(100)" json:"kecamatan"`
	Kelurahan string `gorm:"type:varchar(100)" json:"kelurahan"`
	KodePos   string `gorm:"type:varchar(10)" json:"kode_pos"`

	// ==========================================
	// BRANDING
	// ==========================================
	LogoURL string `gorm:"type:text" json:"logo_url"`

	IsTaxActive bool `gorm:"default:false" json:"is_tax_active"`

	PajakPersen float64 `gorm:"type:decimal(5,2);default:0.00" json:"pajak_persen"`

	ReceiptFooter string `gorm:"type:text" json:"receipt_footer"`

	// ==========================================
	// PAYMENT
	// ==========================================
	PaymentType string `gorm:"type:varchar(50);default:'qris_static';index" json:"payment_type"`

	QrisImage string `gorm:"type:text" json:"qris_image"`
	QrisName  string `gorm:"type:varchar(100)" json:"qris_name"`

	// ==========================================
	// SECRET CONFIG
	// ==========================================
	MidtransServerKey string `gorm:"type:varchar(255)" json:"-"`

	MidtransClientKey string `gorm:"type:varchar(255)" json:"-"`

	WaToken string `gorm:"type:varchar(255)" json:"-"`

	// ==========================================
	// PRINTER
	// ==========================================
	PrinterWidth string `gorm:"type:varchar(20);default:'58mm'" json:"printer_width"`

	PrinterType string `gorm:"type:varchar(20);default:'bluetooth'" json:"printer_type"`

	// ==========================================
	// RELATIONS
	// ==========================================
	Users []User `gorm:"foreignKey:StoreID" json:"-"`

	// ==========================================
	// TIMESTAMP
	// ==========================================
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}