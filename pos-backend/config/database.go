package config

import (
	"log"
	"os"
	"pos-backend/models" // Pastikan path ini benar sesuai struktur folder Mas

	"gorm.io/driver/postgres" // Ganti dari mysql ke postgres
	"gorm.io/gorm"
)

// Variabel Global untuk menampung Koneksi DB
var DB *gorm.DB

func ConnectDatabase() {
	// MASUKKAN URI DARI SUPABASE DI SINI
	// Pastikan [PASSWORD] sudah diganti dengan password asli Mas tanpa kurung siku
	dsn := os.Getenv("DB_URL")
	
	// Gunakan postgres.Open bukan mysql.Open
	database, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

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
		&models.Attendance{},
		&models.CashierSession{},
		&models.Purchase{},
		&models.PurchaseDetail{},
		&models.StockOpname{},
		&models.StockOpnameDetail{},
		&models.Schedule{},
	)
	
	if err != nil {
		log.Fatal("Gagal Melakukan Migrasi Database! Error: ", err)
	}
	log.Println("✅ Tabel database berhasil di-generate di Cloud Supabase!")

	DB = database
}