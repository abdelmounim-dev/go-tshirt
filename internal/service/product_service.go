package service

import (
	"go-tshirt/internal/models"
	"go-tshirt/internal/repository"
)

type ProductService struct {
	repo *repository.ProductRepo
}

func NewProductService(repo *repository.ProductRepo) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(p *models.Product) error {
	return s.repo.Create(p)
}

func (s *ProductService) GetAll() ([]models.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductService) GetByID(id uint) (*models.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) Update(p *models.Product) error {
	return s.repo.Update(p)
}

func (s *ProductService) Delete(id uint) error {
	return s.repo.Delete(id)
}
