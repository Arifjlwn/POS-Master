package repository

import (
	"pos-backend/src/modules/fnb/domain"
	"gorm.io/gorm"
)

type MenuRepository interface {
	Create(menu *domain.Menu) error
	GetAll(storeID uint, onlyAvailable bool) ([]domain.Menu, error)
	GetByID(id uint, storeID uint) (*domain.Menu, error) // <-- INI TADI YANG KETINGGALAN
	Update(menu *domain.Menu) error
	Delete(id uint, storeID uint) error
}

type menuRepo struct {
	db *gorm.DB
}

func NewMenuRepo(db *gorm.DB) MenuRepository {
	return &menuRepo{db}
}

func (r *menuRepo) Create(m *domain.Menu) error {
	return r.db.Create(m).Error
}

func (r *menuRepo) GetAll(storeID uint, onlyAvailable bool) ([]domain.Menu, error) {
	var menus []domain.Menu
	query := r.db.Where("store_id = ?", storeID)
	if onlyAvailable {
		query = query.Where("is_available = ?", true)
	}
	err := query.Find(&menus).Error
	return menus, err
}

func (r *menuRepo) GetByID(id uint, storeID uint) (*domain.Menu, error) {
	var menu domain.Menu
	err := r.db.Where("id = ? AND store_id = ?", id, storeID).First(&menu).Error
	return &menu, err
}

func (r *menuRepo) Update(m *domain.Menu) error {
	return r.db.Save(m).Error
}

func (r *menuRepo) Delete(id uint, storeID uint) error {
	return r.db.Where("id = ? AND store_id = ?", id, storeID).Delete(&domain.Menu{}).Error
}