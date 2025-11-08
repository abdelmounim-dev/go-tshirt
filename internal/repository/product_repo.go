package repository

import (
	"go-tshirt/internal/models"
	"gorm.io/gorm"
)

type ProductRepo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) AutoMigrate() error {
	return r.db.AutoMigrate(&models.Product{})
}

func (r *ProductRepo) Create(p *models.Product) error {
	return r.db.Create(p).Error
}

func (r *ProductRepo) GetAll() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Find(&products).Error
	return products, err
}

func (r *ProductRepo) GetByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, id).Error
	return &product, err
}

func (r *ProductRepo) Update(p *models.Product) error {
	return r.db.Save(p).Error
}

func (r *ProductRepo) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}
