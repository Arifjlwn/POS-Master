package repository

import (
	"pos-backend/src/modules/fnb/domain"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(order *domain.OrderFnB) error
	GetKitchenQueue(storeID uint) ([]domain.OrderFnB, error)
	UpdateKitchenStatus(orderID uint, storeID uint, status string) error
}

type orderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) OrderRepository {
	return &orderRepo{db}
}

// CreateOrder menyimpan nota beserta item rinciannya (GORM otomatis handle transaction via association)
func (r *orderRepo) CreateOrder(order *domain.OrderFnB) error {
	return r.db.Create(order).Error
}

// GetKitchenQueue mengambil orderan yang status dapurnya belum selesai
func (r *orderRepo) GetKitchenQueue(storeID uint) ([]domain.OrderFnB, error) {
	var orders []domain.OrderFnB
	err := r.db.Preload("Items").
		Where("store_id = ? AND status_dapur != ?", storeID, "SELESAI").
		Order("created_at asc").
		Find(&orders).Error
	return orders, err
}

// UpdateKitchenStatus mengubah status antrean dapur (MASAK/SELESAI)
func (r *orderRepo) UpdateKitchenStatus(orderID uint, storeID uint, status string) error {
	return r.db.Model(&domain.OrderFnB{}).
		Where("id = ? AND store_id = ?", orderID, storeID).
		Update("status_dapur", status).Error
}