package repository

import (
	"strings"
	"fmt"


	"pos-backend/models"
	"pos-backend/src/modules/jasalayanan/laundry/domain"

	"gorm.io/gorm"
)

type LaundryRepository interface {
	GetKasirByStoreID(storeID uint) ([]models.User, error)
	GetStoreByID(storeID uint) (*models.Store, error)
	CreateKasir(user *models.User) error
	DeleteKasir(id uint, storeID uint) error
	GetLayananLaundry(storeID uint) ([]models.Product, error)
	FindCustomerByPhone(storeID uint, phone string) (*models.Customer, error)
	CreateCustomer(customer *models.Customer) error
	UpdateCustomer(customer *models.Customer) error
	SearchCustomers(storeID uint, keyword string) ([]models.Customer, error)
	CreateTransactionTx(tx *gorm.DB, transaction *models.Transaction) error
	CreateLaundryDetailTx(tx *gorm.DB, detail *domain.TransactionLaundryDetail) error
	GetTransactionByID(id uint, storeID uint) (*models.Transaction, error)
	UpdateTransaction(transaction *models.Transaction) error
	GetDB() *gorm.DB
	CreateLayanan(product *models.Product) error
	DeleteLayanan(id uint, storeID uint) error
	GetLayananByID(id uint, storeID uint) (*models.Product, error)
	UpdateLayanan(product *models.Product) error
	GetPerfumesByStoreID(storeID uint) ([]domain.Perfume, error)
	CreatePerfume(perfume *domain.Perfume) error
	DeletePerfume(id uint, storeID uint) error
	GetAllTransactions(storeID uint) ([]models.Transaction, error)
	GetLaundryDetailByTxID(txID uint) (*domain.TransactionLaundryDetail, error)
	GetProductByIDSimple(id uint) (*models.Product, error)
	UpdateStoreTx(tx *gorm.DB, store *models.Store) error
	GetTrackingCucian(storeID uint) ([]domain.TrackingResponse, error)
	UpdateStatusDetailCucian(id uint, status string) error
}

type laundryRepo struct{ db *gorm.DB }

func NewLaundryRepo(db *gorm.DB) LaundryRepository { return &laundryRepo{db} }

func (r *laundryRepo) GetDB() *gorm.DB { return r.db }

func (r *laundryRepo) GetKasirByStoreID(storeID uint) ([]models.User, error) {
	var list []models.User
	err := r.db.Where("store_id = ? AND role = ?", storeID, "kasir").Find(&list).Error
	return list, err
}
func (r *laundryRepo) GetStoreByID(storeID uint) (*models.Store, error) {
	var store models.Store
	err := r.db.First(&store, storeID).Error
	return &store, err
}
func (r *laundryRepo) CreateKasir(user *models.User) error { return r.db.Create(user).Error }
func (r *laundryRepo) DeleteKasir(id uint, storeID uint) error {
	return r.db.Where("id = ? AND store_id = ? AND role = ?", id, storeID, "kasir").Delete(&models.User{}).Error
}
func (r *laundryRepo) GetLayananLaundry(storeID uint) ([]models.Product, error) {
    var products []models.Product
    
    // 🚀 FIXED KASTA TERTINGGI: Filter berdasarkan product_type, bukan kategori!
    err := r.db.Where("store_id = ? AND product_type = ? AND is_active = ?", storeID, "JASA_LAUNDRY", true).Find(&products).Error
    
    return products, err
}
func (r *laundryRepo) FindCustomerByPhone(storeID uint, phone string) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.Where("store_id = ? AND no_whatsapp = ?", storeID, phone).First(&customer).Error
	return &customer, err
}
func (r *laundryRepo) CreateCustomer(customer *models.Customer) error {
	return r.db.Create(customer).Error
}
func (r *laundryRepo) UpdateCustomer(customer *models.Customer) error {
	return r.db.Save(customer).Error
}
func (r *laundryRepo) SearchCustomers(storeID uint, keyword string) ([]models.Customer, error) {
	var customers []models.Customer
	query := r.db.Where("store_id = ?", storeID)
	if keyword != "" {
		query = query.Where("nama ILIKE ? OR no_whatsapp LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	return customers, query.Order("updated_at desc").Limit(5).Find(&customers).Error
}
func (r *laundryRepo) CreateTransactionTx(tx *gorm.DB, transaction *models.Transaction) error {
	return tx.Create(transaction).Error
}
func (r *laundryRepo) CreateLaundryDetailTx(tx *gorm.DB, detail *domain.TransactionLaundryDetail) error {
	return tx.Create(detail).Error
}
func (r *laundryRepo) GetTransactionByID(id uint, storeID uint) (*models.Transaction, error) {
	var trx models.Transaction
	err := r.db.Where("id = ? AND store_id = ?", id, storeID).First(&trx).Error
	return &trx, err
}
func (r *laundryRepo) UpdateTransaction(transaction *models.Transaction) error {
	return r.db.Save(transaction).Error
}
func (r *laundryRepo) CreateLayanan(p *models.Product) error { return r.db.Create(p).Error }
func (r *laundryRepo) DeleteLayanan(id uint, storeID uint) error {
	// 1. Coba Hard Delete dulu (kalau layanannya belum pernah dipakai transaksi, hapus aja biar database bersih)
	err := r.db.Where("id = ? AND store_id = ?", id, storeID).Delete(&models.Product{}).Error
	
	if err != nil {
		// 2. 🛡️ SECURITY KASTA TINGGI: Kalau PostgreSQL nolak karena Foreign Key (23503)
		if strings.Contains(err.Error(), "23503") || strings.Contains(strings.ToLower(err.Error()), "foreign key") {
			
			// 🚀 AUTO-FALLBACK: Kita arsipkan (Soft Delete) layanannya dengan mengubah is_active = false
			errArsip := r.db.Model(&models.Product{}).
				Where("id = ? AND store_id = ?", id, storeID).
				Update("is_active", false).Error
				
			if errArsip != nil {
				return fmt.Errorf("gagal mengarsipkan layanan: %v", errArsip)
			}
			
			// Return nil karena proses arsip sukses
			return nil
		}
		
		// Kalau error lain (misal database putus)
		return fmt.Errorf("gagal menghapus layanan dari database: %v", err)
	}
	
	return nil
}

func (r *laundryRepo) GetLayananByID(id uint, storeID uint) (*models.Product, error) {
    var p models.Product
    // 🚀 FIXED: Tambahin is_active = true
    err := r.db.Where("id = ? AND store_id = ? AND is_active = ?", id, storeID, true).First(&p).Error
    return &p, err
}
func (r *laundryRepo) UpdateLayanan(p *models.Product) error { return r.db.Save(p).Error }
func (r *laundryRepo) GetPerfumesByStoreID(storeID uint) ([]domain.Perfume, error) {
	var list []domain.Perfume
	err := r.db.Where("store_id = ?", storeID).Find(&list).Error
	return list, err
}
func (r *laundryRepo) CreatePerfume(p *domain.Perfume) error { return r.db.Create(p).Error }
func (r *laundryRepo) DeletePerfume(id uint, storeID uint) error {
	return r.db.Where("id = ? AND store_id = ?", id, storeID).Delete(&domain.Perfume{}).Error
}
func (r *laundryRepo) GetAllTransactions(storeID uint) ([]models.Transaction, error) {
	var list []models.Transaction
	err := r.db.Where("store_id = ?", storeID).Order("created_at desc").Find(&list).Error
	return list, err
}
func (r *laundryRepo) GetLaundryDetailByTxID(txID uint) (*domain.TransactionLaundryDetail, error) {
	var d domain.TransactionLaundryDetail
	err := r.db.Where("transaction_id = ?", txID).First(&d).Error
	return &d, err
}
func (r *laundryRepo) GetProductByIDSimple(id uint) (*models.Product, error) {
	var p models.Product
	err := r.db.Where("id = ?", id).First(&p).Error
	return &p, err
}
func (r *laundryRepo) UpdateStoreTx(tx *gorm.DB, store *models.Store) error {
	return tx.Save(store).Error
}

func (r *laundryRepo) GetTrackingCucian(storeID uint) ([]domain.TrackingResponse, error) {
	var results []domain.TrackingResponse
	err := r.db.Table("laundry_transaction_details").
		Select(`
            laundry_transaction_details.id, 
            transactions.no_invoice as invoice, 
            laundry_transaction_details.nama_pelanggan as pelanggan, 
            laundry_transaction_details.no_whatsapp as whatsapp, 
            products.nama_produk as layanan, 
            laundry_transaction_details.berat_kg, 
            laundry_transaction_details.sub_total, 
            laundry_transaction_details.status_cucian,
            laundry_transaction_details.status_bayar,
            laundry_transaction_details.rack_id,
            laundry_transaction_details.nomor_rak,
            laundry_transaction_details.estimasi_waktu
        `).
		Joins("left join transactions on transactions.id = laundry_transaction_details.transaction_id").
		Joins("left join products on products.id = laundry_transaction_details.product_id").
		Where("transactions.store_id = ? AND laundry_transaction_details.status_cucian != 'DIAMBIL'", storeID).
		Order("laundry_transaction_details.id desc").
		Scan(&results).Error
	return results, err
}

func (r *laundryRepo) UpdateStatusDetailCucian(id uint, status string) error {
	return r.db.Table("laundry_transaction_details").Where("id = ?", id).Update("status_cucian", status).Error
}
