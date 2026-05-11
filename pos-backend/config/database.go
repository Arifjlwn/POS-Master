package config

import (
	"log"
	"pos-backend/models" // Pastikan path ini benar sesuai struktur folder Mas

	"gorm.io/driver/postgres" // Ganti dari mysql ke postgres
	"gorm.io/gorm"
)

// Variabel Global untuk menampung Koneksi DB
var DB *gorm.DB

func ConnectDatabase() {
	// MASUKKAN URI DARI SUPABASE DI SINI
	// Pastikan [PASSWORD] sudah diganti dengan password asli Mas tanpa kurung siku
	dsn := "postgresql://postgres:Arifjlwn020700@db.pxxjqewukgpfxmbwjnni.supabase.co:5432/postgres"
	
	// Gunakan postgres.Open bukan mysql.Open
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Gagal Menyambung ke Supabase! Error: ", err)
	}

	log.Println("✅ Berhasil terhubung ke Database Supabase (Postgres)!")

	// Auto Migrate (GORM bakal otomatis ngebangun tabel di Supabase)
	err = database.AutoMigrate(
		&models.Store{},
		&models.User{},
		&models.Product{},
		&models.Transaction{},
		&models.TransactionDetail{},
	)
	
	if err != nil {
		log.Fatal("Gagal Melakukan Migrasi Database! Error: ", err)
	}
	log.Println("✅ Tabel database berhasil di-generate di Cloud Supabase!")

	DB = database
}