package src

import (
	"log"
	"os"
	"time"
	"pos-backend/models"
	fnbDomain "pos-backend/src/modules/fnb/domain"
	laundryDomain "pos-backend/src/modules/jasalayanan/laundry/domain"
	retailDomain "pos-backend/src/modules/retail/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("DB_URL")
	database, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn, PreferSimpleProtocol: true}), &gorm.Config{})
	if err != nil { log.Fatal("Gagal Menyambung ke Supabase! Error: ", err) }
	log.Println("✅ Berhasil terhubung ke Database Supabase (Postgres)!")

	sqlDB, err := database.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
		log.Println("⚡ Connection Pool berhasil dikonfigurasi: Max Open 100, Max Idle 10!")
	}

	err = database.AutoMigrate(
		// 🌐 CORE GLOBAL (Pondasi Utama & SDM - Urutan Wajib Di-lock Duluan Bray)
		&models.Store{}, &models.User{}, &models.CashierSession{}, &models.Product{}, &models.Customer{}, &models.Attendance{}, &models.Schedule{}, &models.ShiftClosing{}, &models.Transaction{},

		// 🛍️ RETAIL & INVENTORY
		&models.TransactionDetail{}, &retailDomain.StockOpname{}, &retailDomain.StockOpnameDetail{}, &retailDomain.StockAdjustment{}, &retailDomain.StockAdjustmentDetail{}, &retailDomain.ProductReturn{}, &retailDomain.Purchase{}, &retailDomain.PurchaseDetail{},

		// 🧺 LAYANAN & JASA (LAUNDRY)
		&laundryDomain.TransactionLaundryDetail{}, &laundryDomain.Perfume{},

		// 🍔 FOOD & BEVERAGES
		&fnbDomain.Menu{}, &fnbDomain.OrderFnB{}, &fnbDomain.OrderItemFnB{},
	)
	if err != nil { log.Fatal("Gagal Melakukan Migrasi Database! Error: ", err) }

	log.Println("✅ Tabel database berhasil di-generate di Cloud Supabase!")
	DB = database
}