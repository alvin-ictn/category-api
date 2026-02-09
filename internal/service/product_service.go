package service

import (
	"cateogry-api/internal/domain"
	"cateogry-api/internal/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAll(name string) ([]domain.Product, error) {
	return s.repo.GetAll(name)
}

func (s *ProductService) GetByID(id int) (*domain.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) Create(product domain.Product) (domain.Product, error) {
	return s.repo.Create(product)
}

func (s *ProductService) Update(id int, product domain.Product) (*domain.Product, error) {
	return s.repo.Update(id, product)
}

func (s *ProductService) Delete(id int) error {
	return s.repo.Delete(id)
}
