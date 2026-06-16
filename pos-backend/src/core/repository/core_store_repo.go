package repository

import (
	"gorm.io/gorm"
)

type CoreRepo interface {
	GetDB() *gorm.DB
}

type coreRepo struct {
	db *gorm.DB
}

func NewCoreRepo(db *gorm.DB) CoreRepo {
	return &coreRepo{db: db}
}

func (r *coreRepo) GetDB() *gorm.DB {
	return r.db
}
