package service

import (
	"cateogry-api/internal/domain"
	"cateogry-api/internal/repository"
)

type CategoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAllCategories() ([]domain.Category, error) {
	return s.repo.GetAll()
}

func (s *CategoryService) GetCategoryByID(id int) (*domain.Category, error) {
	return s.repo.GetByID(id)
}

func (s *CategoryService) CreateCategory(c domain.Category) (domain.Category, error) {
	return s.repo.Create(c)
}

func (s *CategoryService) UpdateCategory(id int, c domain.Category) (*domain.Category, error) {
	return s.repo.Update(id, c)
}

func (s *CategoryService) DeleteCategory(id int) error {
	return s.repo.Delete(id)
}
