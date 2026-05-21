package domain

import "time"

type Menu struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    StoreID     uint      `gorm:"index;not null" json:"store_id"`
    NamaProduk  string    `gorm:"type:varchar(150);not null" json:"nama_produk"`
    Kategori    string    `gorm:"index;type:varchar(50)" json:"kategori"`
    HargaJual   float64   `gorm:"type:decimal(10,2);not null" json:"harga_jual"`
    HargaModal  float64   `gorm:"type:decimal(10,2);default:0" json:"harga_modal"`
    Stok        int       `gorm:"default:0" json:"stok"`
    Gambar      string    `gorm:"type:varchar(255)" json:"gambar"`
    IsAvailable bool      `gorm:"default:true" json:"is_available"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

func (Menu) TableName() string { return "fnb_menus" }