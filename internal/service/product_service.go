package service

import (
	"errors"

	custom_errors "github.com/abdelmounim-dev/go-tshirt/internal/errors"
	"github.com/abdelmounim-dev/go-tshirt/internal/models"
	"github.com/abdelmounim-dev/go-tshirt/internal/repository"
	"gorm.io/gorm"
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
	if len(p.Variants) == 0 {
		return &custom_errors.ValidationError{Message: "product must have at least one variant"}
	}
	return s.repo.Create(p)
}

func (s *ProductService) GetAll() ([]models.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductService) GetByID(id uint) (*models.Product, error) {
	product, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &custom_errors.NotFoundError{Message: "product not found"}
		}
		return nil, err
	}
	return product, nil
}

func (s *ProductService) Update(p *models.Product) error {
	if p.Name == "" {
		return &custom_errors.ValidationError{Message: "product name cannot be empty"}
	}
	if p.Price <= 0 {
		return &custom_errors.ValidationError{Message: "product price must be greater than zero"}
	}
	if len(p.Variants) == 0 {
		return &custom_errors.ValidationError{Message: "product must have at least one variant"}
	}
	return s.repo.Update(p)
}

func (s *ProductService) Delete(id uint) error {
	return s.repo.Delete(id)
}
