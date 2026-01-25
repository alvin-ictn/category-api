package repository

import (
	"cateogry-api/internal/domain"
	"errors"
	"time"
)

type InMemoryCategoryRepository struct {
	categories []domain.Category
}

func NewInMemoryCategoryRepository() *InMemoryCategoryRepository {
	return &InMemoryCategoryRepository{
		categories: []domain.Category{
			{ID: 1, Name: "Electronics", Description: "Electronic devices and accessories", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: 2, Name: "Books", Description: "Books and literature", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: 3, Name: "Clothing", Description: "Clothing and accessories", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
	}
}

func (r *InMemoryCategoryRepository) GetAll() ([]domain.Category, error) {
	return r.categories, nil
}

func (r *InMemoryCategoryRepository) GetByID(id int) (*domain.Category, error) {
	for _, c := range r.categories {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, errors.New("category not found")
}

func (r *InMemoryCategoryRepository) Create(c domain.Category) (domain.Category, error) {
	maxID := 0
	for _, cat := range r.categories {
		if cat.ID > maxID {
			maxID = cat.ID
		}
	}
	c.ID = maxID + 1
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	r.categories = append(r.categories, c)
	return c, nil
}

func (r *InMemoryCategoryRepository) Update(id int, u domain.Category) (*domain.Category, error) {
	for i, c := range r.categories {
		if c.ID == id {
			r.categories[i].Name = u.Name
			r.categories[i].Description = u.Description
			r.categories[i].UpdatedAt = time.Now()
			return &r.categories[i], nil
		}
	}
	return nil, errors.New("category not found")
}

func (r *InMemoryCategoryRepository) Delete(id int) error {
	for i, c := range r.categories {
		if c.ID == id {
			r.categories = append(r.categories[:i], r.categories[i+1:]...)
			return nil
		}
	}
	return errors.New("category not found")
}
