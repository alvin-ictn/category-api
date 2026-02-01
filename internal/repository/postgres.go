package repository

import (
	"cateogry-api/internal/domain"
	"database/sql"
	"errors"
)

type PostgresCategoryRepository struct {
	db *sql.DB
}

func NewPostgresCategoryRepository(db *sql.DB) *PostgresCategoryRepository {
	return &PostgresCategoryRepository{db: db}
}

func (r *PostgresCategoryRepository) GetAll() ([]domain.Category, error) {
	query := "SELECT id, name, description, created_at, updated_at FROM categories"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *PostgresCategoryRepository) GetByID(id int) (*domain.Category, error) {
	query := "SELECT id, name, description, created_at, updated_at FROM categories WHERE id = $1"
	var c domain.Category
	err := r.db.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("category not found")
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *PostgresCategoryRepository) Create(category domain.Category) (domain.Category, error) {
	query := "INSERT INTO categories (name, description, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id"
	err := r.db.QueryRow(query, category.Name, category.Description, category.CreatedAt, category.UpdatedAt).Scan(&category.ID)
	if err != nil {
		return domain.Category{}, err
	}
	return category, nil
}

func (r *PostgresCategoryRepository) Update(id int, category domain.Category) (*domain.Category, error) {
	query := "UPDATE categories SET name = $1, description = $2, updated_at = $3 WHERE id = $4 RETURNING id, name, description, created_at, updated_at"
	var c domain.Category
	err := r.db.QueryRow(query, category.Name, category.Description, category.UpdatedAt, id).Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("category not found")
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *PostgresCategoryRepository) Delete(id int) error {
	query := "DELETE FROM categories WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("category not found")
	}
	return nil
}
