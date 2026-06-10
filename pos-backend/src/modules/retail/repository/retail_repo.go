package repository

import (
	"math"
	"pos-backend/models"
	"pos-backend/src/modules/retail/domain"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RetailRepository interface {
	GetDB() *gorm.DB
	GetProductByID(tx *gorm.DB, id uint, storeID uint) (*models.Product, error)
	SaveProduct(tx *gorm.DB, product *models.Product) error
	UpdateProductStokExpr(tx *gorm.DB, id uint, storeID uint, qty int) error

	// Stock Opname
	CreateStockOpname(tx *gorm.DB, so *domain.StockOpname) error
	CreateStockOpnameDetail(tx *gorm.DB, detail *domain.StockOpnameDetail) error
	GetStockOpnameHistory(storeID uint) ([]domain.StockOpname, error)
	CheckStockOpnameThisMonth(storeID uint, currentMonth int, currentYear int) bool

	// Return
	CreateProductReturns(tx *gorm.DB, returns []domain.ProductReturn) error
	GetReturnsHistory(storeID uint, limit int, offset int) ([]domain.ProductReturn, int64, error)

	// Purchase / LPB
	CreatePurchase(tx *gorm.DB, purchase *domain.Purchase) error
	CreatePurchaseDetail(tx *gorm.DB, detail *domain.PurchaseDetail) error
	CreatePurchaseWithMovingAverage(db *gorm.DB, purchase *domain.Purchase) error

	// Absensi Karyawan
	GetAttendanceToday(userID uint, tanggal string) (*models.Attendance, error)
	CreateAttendance(attendance *models.Attendance) error
	SaveAttendance(attendance *models.Attendance) error
	GetAttendanceReport(storeID uint, filterTanggal string, prefixBulan string) ([]models.Attendance, error)
	GetSchedulesForMangkir(storeID uint, tanggal string) ([]models.Schedule, error)

	// HR / Data Karyawan
	GetLastEmployeeNIK(storeID uint, currentYear string) (*models.User, error)
	CreateEmployee(user *models.User) error
	GetAllEmployees(storeID uint) ([]models.User, error)
	GetEmployeeByID(id uint, storeID uint) (*models.User, error)
	SaveEmployee(user *models.User) error

	// Master Produk & Dashboard Report
	CreateProductGlobal(product *models.Product) error
	GetProductsCatalog(storeID uint, search string, category string, limit int, offset int, usePagination bool) ([]models.Product, int64, error)
	GetProductByIDSimple(id uint, storeID uint) (*models.Product, error)
	DeleteProductGlobal(product *models.Product) error
	GetDistinctCategories(storeID uint) ([]string, error)
	GetAllProductsForExport(storeID uint) ([]models.Product, error)
	GetDashboardSummary(storeID uint, start time.Time, end time.Time) (float64, float64, error)
	GetDashboardLaba(storeID uint, start time.Time, end time.Time) (float64, error)
	GetDashboardReturSummary(storeID uint, start time.Time, end time.Time) (float64, float64, error)
	GetDashboardSOSummary(storeID uint, start time.Time, end time.Time) (float64, float64, error)
	GetDashboardKlaimSummary(storeID uint, start time.Time, end time.Time) (float64, float64, error)
	GetLowStockProducts(storeID uint, limitStock int) ([]models.Product, error)
	GetDailySalesReport(storeID uint, tgl time.Time, tglEnd time.Time) (float64, float64, float64, error)
	GetAggregatedDailySales(storeID uint, start time.Time, end time.Time) ([]map[string]interface{}, error) // 🚀 SANGAR: Senjata Pemungkas Anti N+1 Loop
	GetTopBestSellers(storeID uint, start time.Time, end time.Time) ([]domain.BestSeller, error)

	// Jadwal Karyawan
	GetScheduleByDate(tx *gorm.DB, userID uint, tanggal string) (*models.Schedule, error)
	SaveScheduleTx(tx *gorm.DB, schedule *models.Schedule) error
	CreateScheduleTx(tx *gorm.DB, schedule *models.Schedule) error
	GetSchedulesRange(storeID uint, start string, end string) ([]models.Schedule, error)

	// Cashier Session Shift
	GetActiveSession(tx *gorm.DB, userID uint, storeID uint) (*models.CashierSession, error)
	CreateSession(session *models.CashierSession) error
	GetSessionByIDPreloaded(id uint) (*models.CashierSession, error)
	SaveSession(session *models.CashierSession) error
	GetSalesMethodSummary(sessionID string, method string) (float64, error)
	GetSalesTotalAndTax(sessionID string) (float64, float64, error)

	// Transaksi POS Checkout
	GetStoreByIDSimple(tx *gorm.DB, storeID uint) (*models.Store, error)
	CreateTransactionTx(tx *gorm.DB, transaction *models.Transaction) error
	GetTransactionsByRange(storeID uint, start time.Time, end time.Time) ([]models.Transaction, error)
	GetClosingByRange(storeID uint, startOfDay time.Time, endOfDay time.Time) ([]models.ShiftClosing, error)
}

type retailRepo struct{ db *gorm.DB }

func NewRetailRepo(db *gorm.DB) RetailRepository { return &retailRepo{db} }
func (r *retailRepo) GetDB() *gorm.DB            { return r.db }

func (r *retailRepo) CreatePurchase(tx *gorm.DB, p *domain.Purchase) error { return tx.Create(p).Error }
func (r *retailRepo) CreatePurchaseDetail(tx *gorm.DB, d *domain.PurchaseDetail) error {
	return tx.Create(d).Error
}

func (r *retailRepo) CreatePurchaseWithMovingAverage(db *gorm.DB, purchase *domain.Purchase) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(purchase).Error; err != nil {
			return err
		}
		for _, detail := range purchase.Details {
			var product models.Product
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("id = ? AND store_id = ?", detail.ProductID, purchase.StoreID).
				First(&product).Error; err != nil {
				return err
			}
			stokLama := float64(product.Stok)
			hargaModalLama := product.HargaModal
			stokBaruMasuk := float64(detail.QtyMasuk)
			hargaModalBaruMasuk := detail.HargaModal
			totalAsetLama := stokLama * hargaModalLama
			totalAsetBaru := stokBaruMasuk * hargaModalBaruMasuk
			totalStokAkhir := stokLama + stokBaruMasuk
			var hargaPokokRataRata float64
			if totalStokAkhir > 0 {
				rawAverage := (totalAsetLama + totalAsetBaru) / totalStokAkhir
				hargaPokokRataRata = math.Ceil(rawAverage/100) * 100
			} else {
				hargaPokokRataRata = hargaModalLama
			}
			if err := tx.Model(&product).Updates(map[string]interface{}{"stok": int(totalStokAkhir), "harga_modal": hargaPokokRataRata}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *retailRepo) GetProductByID(tx *gorm.DB, id uint, storeID uint) (*models.Product, error) {
	var p models.Product
	err := tx.Where("id = ? AND store_id = ?", id, storeID).First(&p).Error
	return &p, err
}
func (r *retailRepo) SaveProduct(tx *gorm.DB, p *models.Product) error {
	// 🔒 ANTI-SABOTASE: Memastikan cuma bisa nge-update barang di toko lu sendiri
	return tx.Model(&models.Product{}).
		Where("id = ? AND store_id = ?", p.ID, p.StoreID).
		Updates(p).Error
}
func (r *retailRepo) UpdateProductStokExpr(tx *gorm.DB, id uint, storeID uint, qty int) error {
	return tx.Model(&models.Product{}).
		Where("id = ? AND store_id = ?", id, storeID). // 🔒 AMAN DARI SABOTASE
		Update("stok", gorm.Expr("stok - ?", qty)).Error
}

func (r *retailRepo) CreateStockOpname(tx *gorm.DB, so *domain.StockOpname) error {
	return tx.Create(so).Error
}
func (r *retailRepo) CheckStockOpnameThisMonth(storeID uint, currentMonth int, currentYear int) bool {
	var count int64
	r.db.Model(&domain.StockOpname{}).Where("store_id = ? AND EXTRACT(MONTH FROM created_at) = ? AND EXTRACT(YEAR FROM created_at) = ?", storeID, currentMonth, currentYear).Count(&count)
	return count > 0
}
func (r *retailRepo) CreateStockOpnameDetail(tx *gorm.DB, detail *domain.StockOpnameDetail) error {
	return tx.Create(detail).Error
}
func (r *retailRepo) GetStockOpnameHistory(storeID uint) ([]domain.StockOpname, error) {
	var history []domain.StockOpname
	err := r.db.Preload("Details.Product").Where("store_id = ?", storeID).Order("created_at DESC").Find(&history).Error
	return history, err
}

func (r *retailRepo) CreateProductReturns(tx *gorm.DB, returns []domain.ProductReturn) error {
	return tx.Create(&returns).Error
}
func (r *retailRepo) GetReturnsHistory(storeID uint, limit int, offset int) ([]domain.ProductReturn, int64, error) {
	var list []domain.ProductReturn
	var total int64
	query := r.db.Model(&domain.ProductReturn{}).Where("store_id = ?", storeID).Preload("Product").Preload("User")
	query.Count(&total)
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	err := query.Order("created_at DESC").Find(&list).Error
	return list, total, err
}

func (r *retailRepo) GetAttendanceToday(userID uint, tanggal string) (*models.Attendance, error) {
	var a models.Attendance
	err := r.db.Where("user_id = ? AND tanggal = ?", userID, tanggal).First(&a).Error
	return &a, err
}
func (r *retailRepo) CreateAttendance(a *models.Attendance) error { return r.db.Create(a).Error }
func (r *retailRepo) SaveAttendance(a *models.Attendance) error   { return r.db.Save(a).Error }
func (r *retailRepo) GetAttendanceReport(storeID uint, filterTanggal string, prefixBulan string) ([]models.Attendance, error) {
	var list []models.Attendance
	query := r.db.Preload("User").Where("store_id = ?", storeID)
	if filterTanggal != "" {
		query = query.Where("tanggal = ?", filterTanggal)
	} else if prefixBulan != "" {
		query = query.Where("tanggal::text LIKE ?", prefixBulan)
	}
	err := query.Order("tanggal DESC, jam_masuk DESC").Find(&list).Error
	return list, err
}
func (r *retailRepo) GetSchedulesForMangkir(storeID uint, tanggal string) ([]models.Schedule, error) {
	var list []models.Schedule
	err := r.db.Where("store_id = ? AND tanggal = ? AND shift_type != ?", storeID, tanggal, "OFF").Find(&list).Error
	return list, err
}

func (r *retailRepo) GetLastEmployeeNIK(storeID uint, currentYear string) (*models.User, error) {
	var user models.User
	err := r.db.Where("store_id = ? AND role != ? AND nik LIKE ?", storeID, "owner", currentYear+"%").Order("nik desc").First(&user).Error
	return &user, err
}
func (r *retailRepo) CreateEmployee(u *models.User) error { return r.db.Create(u).Error }
func (r *retailRepo) GetAllEmployees(storeID uint) ([]models.User, error) {
	var list []models.User
	err := r.db.Where("store_id = ?", storeID).Find(&list).Error
	return list, err
}
func (r *retailRepo) GetEmployeeByID(id uint, storeID uint) (*models.User, error) {
	var u models.User
	err := r.db.Where("id = ? AND store_id = ?", id, storeID).First(&u).Error // 🔒 AMAN
	return &u, err
}
func (r *retailRepo) SaveEmployee(u *models.User) error { return r.db.Save(u).Error }

func (r *retailRepo) CreateProductGlobal(p *models.Product) error { return r.db.Create(p).Error }
func (r *retailRepo) GetProductsCatalog(storeID uint, search string, category string, limit int, offset int, usePagination bool) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64
	query := r.db.Model(&models.Product{}).Where("store_id = ?", storeID)
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("(nama_produk ILIKE ? OR sku ILIKE ?)", searchTerm, searchTerm)
	}
	if category != "" {
		query = query.Where("kategori = ?", category)
	}
	query.Count(&total)
	if usePagination {
		query = query.Limit(limit).Offset(offset)
	}
	err := query.Order("id DESC").Find(&products).Error
	return products, total, err
}
func (r *retailRepo) GetProductByIDSimple(id uint, storeID uint) (*models.Product, error) {
	var p models.Product
	err := r.db.Where("id = ? AND store_id = ?", id, storeID).First(&p).Error
	return &p, err
}
func (r *retailRepo) DeleteProductGlobal(p *models.Product) error { return r.db.Delete(p).Error }
func (r *retailRepo) GetDistinctCategories(storeID uint) ([]string, error) {
	var categories []string
	err := r.db.Model(&models.Product{}).Where("store_id = ? AND kategori IS NOT NULL AND kategori != ''", storeID).Distinct("kategori").Pluck("kategori", &categories).Error
	return categories, err
}
func (r *retailRepo) GetAllProductsForExport(storeID uint) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Where("store_id = ?", storeID).Order("id DESC").Find(&products).Error
	return products, err
}

func (r *retailRepo) GetDashboardSummary(storeID uint, start time.Time, end time.Time) (float64, float64, error) {
	var res struct {
		Omzet float64
		Qty   float64
	}
	err := r.db.Table("transaction_details").Joins("JOIN transactions ON transactions.id = transaction_details.transaction_id").Where("transactions.store_id = ? AND transactions.created_at BETWEEN ? AND ?", storeID, start, end).Select("COALESCE(SUM(transaction_details.sub_total), 0) as omzet, COALESCE(SUM(transaction_details.kuantitas), 0) as qty").Scan(&res).Error
	return res.Omzet, res.Qty, err
}
func (r *retailRepo) GetDashboardLaba(storeID uint, start time.Time, end time.Time) (float64, error) {
	var totalLaba float64
	err := r.db.Table("transaction_details").Select("COALESCE(SUM(transaction_details.sub_total - (COALESCE(products.harga_modal, 0) * transaction_details.kuantitas)), 0)").Joins("JOIN transactions ON transactions.id = transaction_details.transaction_id").Joins("LEFT JOIN products ON products.id = transaction_details.product_id").Where("transactions.store_id = ? AND transactions.created_at BETWEEN ? AND ?", storeID, start, end).Row().Scan(&totalLaba)
	return totalLaba, err
}
func (r *retailRepo) GetDashboardReturSummary(storeID uint, start time.Time, end time.Time) (float64, float64, error) {
	var res struct {
		Qty  float64
		Loss float64
	}
	err := r.db.Table("retail_product_returns").Select("COALESCE(SUM(retail_product_returns.qty), 0) as qty, COALESCE(SUM(retail_product_returns.qty * COALESCE(products.harga_modal, 0)), 0) as loss").Joins("LEFT JOIN products ON products.id = retail_product_returns.product_id").Where("retail_product_returns.store_id = ? AND retail_product_returns.created_at BETWEEN ? AND ?", storeID, start, end).Scan(&res).Error
	return res.Qty, res.Loss, err
}
func (r *retailRepo) GetDashboardKlaimSummary(storeID uint, start time.Time, end time.Time) (float64, float64, error) {
	var res struct {
		Qty            float64
		RecoveredValue float64
	}
	err := r.db.Table("retail_stock_adjustment_details").Select("COALESCE(SUM(retail_stock_adjustment_details.qty), 0) as qty, COALESCE(SUM(retail_stock_adjustment_details.qty * COALESCE(products.harga_modal, 0)), 0) as recovered_value").Joins("JOIN retail_stock_adjustments ON retail_stock_adjustments.id = retail_stock_adjustment_details.adjustment_id").Joins("LEFT JOIN products ON products.id = retail_stock_adjustment_details.product_id").Where("retail_stock_adjustments.store_id = ? AND retail_stock_adjustments.status = 'APPROVED' AND retail_stock_adjustments.created_at BETWEEN ? AND ?", storeID, start, end).Scan(&res).Error
	return res.Qty, res.RecoveredValue, err
}
func (r *retailRepo) GetDashboardSOSummary(storeID uint, start time.Time, end time.Time) (float64, float64, error) {
	var res struct {
		Qty  float64
		Loss float64
	}
	err := r.db.Table("retail_stock_opname_details").Select("COALESCE(SUM(ABS(retail_stock_opname_details.selisih)), 0) as qty, COALESCE(SUM(ABS(retail_stock_opname_details.selisih) * COALESCE(products.harga_modal, 0)), 0) as loss").Joins("JOIN retail_stock_opnames ON retail_stock_opnames.id = retail_stock_opname_details.opname_id").Joins("LEFT JOIN products ON products.id = retail_stock_opname_details.product_id").Where("retail_stock_opnames.store_id = ? AND retail_stock_opname_details.selisih < 0 AND retail_stock_opnames.created_at BETWEEN ? AND ?", storeID, start, end).Scan(&res).Error
	return res.Qty, res.Loss, err
}
func (r *retailRepo) GetLowStockProducts(storeID uint, limitStock int) ([]models.Product, error) {
	var list []models.Product
	err := r.db.Where("store_id = ? AND stok < ?", storeID, limitStock).Find(&list).Error
	return list, err
}
func (r *retailRepo) GetDailySalesReport(storeID uint, tgl time.Time, tglEnd time.Time) (float64, float64, float64, error) {
	var sales struct {
		Omzet float64
		Laba  float64
	}
	r.db.Table("transaction_details").Select("COALESCE(SUM(transaction_details.sub_total), 0) as omzet, COALESCE(SUM(transaction_details.sub_total - (COALESCE(products.harga_modal, 0) * transaction_details.kuantitas)), 0) as laba").Joins("JOIN transactions ON transactions.id = transaction_details.transaction_id").Joins("LEFT JOIN products ON products.id = transaction_details.product_id").Where("transactions.store_id = ? AND transactions.created_at BETWEEN ? AND ?", storeID, tgl, tglEnd).Scan(&sales)
	var returLoss float64
	r.db.Table("retail_product_returns").Select("COALESCE(SUM(retail_product_returns.qty * COALESCE(products.harga_modal, 0)), 0)").Joins("LEFT JOIN products ON products.id = retail_product_returns.product_id").Where("retail_product_returns.store_id = ? AND retail_product_returns.created_at BETWEEN ? AND ?", storeID, tgl, tglEnd).Row().Scan(&returLoss)
	return sales.Omzet, sales.Laba, returLoss, nil
}

// 🚀 SANGAR V2: Senjata Utama Dashboard Tingkat Pro - FIX FILTER TANGGAL BERGESER
func (r *retailRepo) GetAggregatedDailySales(storeID uint, start time.Time, end time.Time) ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	// Kita paksa Postgres convert date_val pake TO_CHAR agar output tanggalnya string bersih "YYYY-MM-DD" !
	query := `
WITH dates AS (
    SELECT generate_series(
        DATE(?),
        DATE(?),
        INTERVAL '1 day'
    )::date AS date_val
),

sales_data AS (
    SELECT
        DATE(t.created_at) as date_val,
        COALESCE(SUM(td.sub_total), 0) as omzet,
        COALESCE(
            SUM(
                td.sub_total -
                (COALESCE(p.harga_modal, 0) * td.kuantitas)
            ),
            0
        ) as laba
    FROM transaction_details td
    JOIN transactions t
        ON t.id = td.transaction_id
    LEFT JOIN products p
        ON p.id = td.product_id
    WHERE
        t.store_id = ?
        AND t.created_at BETWEEN ? AND ?
    GROUP BY DATE(t.created_at)
),

return_data AS (
    SELECT
        DATE(r.created_at) as date_val,
        COALESCE(
            SUM(
                r.qty * COALESCE(p2.harga_modal, 0)
            ),
            0
        ) as retur_loss
    FROM retail_product_returns r
    LEFT JOIN products p2
        ON p2.id = r.product_id
    WHERE
        r.store_id = ?
        AND r.created_at BETWEEN ? AND ?
    GROUP BY DATE(r.created_at)
)

SELECT
    TO_CHAR(d.date_val, 'YYYY-MM-DD') as tanggal,
    COALESCE(s.omzet, 0) as omzet,
    COALESCE(s.laba, 0) as laba,
    COALESCE(r.retur_loss, 0) as retur_loss
FROM dates d
LEFT JOIN sales_data s
    ON s.date_val = d.date_val
LEFT JOIN return_data r
    ON r.date_val = d.date_val
ORDER BY d.date_val ASC;
`

	// 🚀 FIX CRITICAL BOUNDARY: Jangan kurangi 'end' dengan 24 jam !
	// Karena 'end' dari handler udah dikunci di jam 23:59:59 pada hari terakhir yang dipilih owner.
	err := r.db.Raw(
		query,
		start, // Generator Tanggal Mulai
		end,   // Generator Tanggal Akhir (Tetap utuh hari terakhir !)

		storeID,
		start,
		end,

		storeID,
		start,
		end,
	).Scan(&result).Error

	return result, err
}

func (r *retailRepo) GetTopBestSellers(storeID uint, start time.Time, end time.Time) ([]domain.BestSeller, error) {
	var list []domain.BestSeller
	err := r.db.Table("transaction_details").Select("products.nama_produk, products.sku, products.satuan_dasar, SUM(transaction_details.kuantitas) as qty_terjual, SUM(transaction_details.sub_total) as total_omzet").Joins("JOIN transactions ON transactions.id = transaction_details.transaction_id").Joins("JOIN products ON products.id = transaction_details.product_id").Where("transactions.store_id = ? AND transactions.created_at BETWEEN ? AND ?", storeID, start, end).Group("products.nama_produk, products.sku, products.satuan_dasar").Order("qty_terjual DESC").Limit(5).Scan(&list).Error
	return list, err
}

func (r *retailRepo) GetScheduleByDate(tx *gorm.DB, userID uint, tanggal string) (*models.Schedule, error) {
	var s models.Schedule
	err := tx.Where("user_id = ? AND tanggal = ?", userID, tanggal).First(&s).Error
	return &s, err
}
func (r *retailRepo) SaveScheduleTx(tx *gorm.DB, s *models.Schedule) error { return tx.Save(s).Error }
func (r *retailRepo) CreateScheduleTx(tx *gorm.DB, s *models.Schedule) error {
	return tx.Create(s).Error
}
func (r *retailRepo) GetSchedulesRange(storeID uint, start string, end string) ([]models.Schedule, error) {
	var list []models.Schedule
	query := r.db.Preload("User").Where("store_id = ?", storeID)
	if start != "" && end != "" {
		query = query.Where("tanggal BETWEEN ? AND ?", start, end)
	}
	err := query.Order("tanggal ASC").Find(&list).Error
	return list, err
}

func (r *retailRepo) GetActiveSession(tx *gorm.DB, userID uint, storeID uint) (*models.CashierSession, error) {
	var s models.CashierSession
	err := tx.Where("user_id = ? AND store_id = ? AND status = ?", userID, storeID, "open").First(&s).Error
	return &s, err
}
func (r *retailRepo) CreateSession(s *models.CashierSession) error { return r.db.Create(s).Error }
func (r *retailRepo) GetSessionByIDPreloaded(id uint) (*models.CashierSession, error) {
	var s models.CashierSession
	err := r.db.Preload("Store").First(&s, id).Error
	return &s, err
}
func (r *retailRepo) SaveSession(s *models.CashierSession) error { return r.db.Save(s).Error }
func (r *retailRepo) GetSalesMethodSummary(sessionID string, method string) (float64, error) {
	var sum float64
	err := r.db.Table("transactions").Select("COALESCE(SUM(total_harga), 0)").Where("session_id = ? AND metode_bayar = ?", sessionID, method).Scan(&sum).Error
	return sum, err
}
func (r *retailRepo) GetSalesTotalAndTax(sessionID string) (float64, float64, error) {
	var res struct {
		Gross float64
		Tax   float64
	}
	// 🛠️ Fix:
	err := r.db.Table("transactions").Select("COALESCE(SUM(total_harga), 0) as gross, COALESCE(SUM(pajak), 0) as tax").Where("session_id = ?", sessionID).Scan(&res).Error
	return res.Gross, res.Tax, err
}

func (r *retailRepo) GetStoreByIDSimple(tx *gorm.DB, storeID uint) (*models.Store, error) {
	var s models.Store
	err := tx.First(&s, storeID).Error
	return &s, err
}
func (r *retailRepo) CreateTransactionTx(tx *gorm.DB, t *models.Transaction) error {
	return tx.Create(t).Error
}
func (r *retailRepo) GetTransactionsByRange(storeID uint, start time.Time, end time.Time) ([]models.Transaction, error) {
	var list []models.Transaction
	err := r.db.Preload("User").Preload("Store").Preload("Details").Preload("Details.Product").Where("store_id = ? AND created_at BETWEEN ? AND ?", storeID, start, end).Order("created_at DESC").Find(&list).Error
	return list, err
}
func (r *retailRepo) GetClosingByRange(storeID uint, startOfDay, endOfDay time.Time) ([]models.ShiftClosing, error) {
	var closings []models.ShiftClosing
	err := r.db.Preload("User").Preload("Session").Preload("Store").Where("store_id = ? AND created_at BETWEEN ? AND ?", storeID, startOfDay, endOfDay).Order("created_at DESC").Find(&closings).Error
	return closings, err
}
