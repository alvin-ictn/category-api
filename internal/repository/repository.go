package repository

import "cateogry-api/internal/domain"

type CategoryRepository interface {
	GetAll() ([]domain.Category, error)
	GetByID(id int) (*domain.Category, error)
	Create(category domain.Category) (domain.Category, error)
	Update(id int, category domain.Category) (*domain.Category, error)
	Delete(id int) error
}
