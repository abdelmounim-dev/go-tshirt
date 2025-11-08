package service

import (
	custom_errors "github.com/abdelmounim-dev/go-tshirt/internal/errors"
	"github.com/abdelmounim-dev/go-tshirt/internal/models"
	"github.com/abdelmounim-dev/go-tshirt/internal/repository"
)

type ProductServiceInterface interface {
	Create(p *models.Product) error
	GetAll() ([]models.Product, error)
	GetByID(id uint) (*models.Product, error)
	Update(p *models.Product) error
	Delete(id uint) error
}

type ProductService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductServiceInterface {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(p *models.Product) error {
	if p.Name == "" {
		return &custom_errors.ValidationError{Message: "product name cannot be empty"}
	}
	if p.Price <= 0 {
		return &custom_errors.ValidationError{Message: "product price must be greater than zero"}
	}
	return s.repo.Create(p)
}

func (s *ProductService) GetAll() ([]models.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductService) GetByID(id uint) (*models.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) Update(p *models.Product) error {
	if p.Name == "" {
		return &custom_errors.ValidationError{Message: "product name cannot be empty"}
	}
	if p.Price <= 0 {
		return &custom_errors.ValidationError{Message: "product price must be greater than zero"}
	}
	return s.repo.Update(p)
}

func (s *ProductService) Delete(id uint) error {
	return s.repo.Delete(id)
}
