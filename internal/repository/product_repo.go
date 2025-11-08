package repository

import (
	"github.com/abdelmounim-dev/go-tshirt/internal/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	AutoMigrate() error
	Create(p *models.Product) error
	GetAll() ([]models.Product, error)
	GetByID(id uint) (*models.Product, error)
	Update(p *models.Product) error
	Delete(id uint) error
}

type ProductRepo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) ProductRepository {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) AutoMigrate() error {
	return r.db.AutoMigrate(&models.Product{}, &models.ProductVariant{})
}

func (r *ProductRepo) Create(p *models.Product) error {
	return r.db.Create(p).Error
}

func (r *ProductRepo) GetAll() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Variants").Find(&products).Error
	return products, err
}

func (r *ProductRepo) GetByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Variants").First(&product, id).Error
	return &product, err
}

func (r *ProductRepo) Update(p *models.Product) error {
	return r.db.Save(p).Error
}

func (r *ProductRepo) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}
