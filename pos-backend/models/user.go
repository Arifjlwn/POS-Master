package models

import "time"

type User struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    // UBAH BARIS INI: Tambahkan bintang (*) dan hapus 'not null'
    StoreID   *uint     `json:"store_id"` 
    
    Name      string    `gorm:"type:varchar(100);not null" json:"name"`
    Email     *string   `gorm:"type:varchar(100);uniqueIndex" json:"email"`
    NIK       *string   `gorm:"type:varchar(20);uniqueIndex" json:"nik"`
    Password  string    `gorm:"type:varchar(255);not null" json:"-"`
    Role      string    `gorm:"type:varchar(20);default:'kasir'" json:"role"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`

    Store     Store     `gorm:"foreignKey:StoreID" json:"store"`
}