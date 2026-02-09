package repository

import (
	"cateogry-api/internal/domain"
	"database/sql"
	"errors"
	"time"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAll(name string) ([]domain.Product, error) {
	query := `
		SELECT p.id, p.name, p.description, p.price, p.stock, p.category_id, 
		       p.created_at, p.updated_at, c.name as category_name
		FROM products p
		JOIN categories c ON p.category_id = c.id
		WHERE p.deleted_at IS NULL
	`

	args := []interface{}{}
	if name != "" {
		query += " AND p.name ILIKE $1"
		args = append(args, "%"+name+"%")
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CategoryID, &p.CreatedAt, &p.UpdatedAt, &p.CategoryName); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepository) GetByID(id int) (*domain.Product, error) {
	query := `
		SELECT p.id, p.name, p.description, p.price, p.stock, p.category_id, 
		       p.created_at, p.updated_at, c.name as category_name
		FROM products p
		JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1 AND p.deleted_at IS NULL
	`
	var p domain.Product
	err := r.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CategoryID, &p.CreatedAt, &p.UpdatedAt, &p.CategoryName)
	if err == sql.ErrNoRows {
		return nil, errors.New("product not found")
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Create(product domain.Product) (domain.Product, error) {
	query := `
		INSERT INTO products (name, description, price, stock, category_id, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id, created_at, updated_at
	`
	now := time.Now()
	err := r.db.QueryRow(query, product.Name, product.Description, product.Price, product.Stock, product.CategoryID, now, now).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return domain.Product{}, err
	}
	return product, nil
}

func (r *ProductRepository) Update(id int, product domain.Product) (*domain.Product, error) {
	query := `
		UPDATE products 
		SET name = $1, description = $2, price = $3, stock = $4, category_id = $5, updated_at = $6 
		WHERE id = $7 AND deleted_at IS NULL
		RETURNING id, name, description, price, stock, category_id, created_at, updated_at
	`
	var p domain.Product
	err := r.db.QueryRow(query, product.Name, product.Description, product.Price, product.Stock, product.CategoryID, time.Now(), id).Scan(
		&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CategoryID, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("product not found")
	}
	if err != nil {
		return nil, err
	}
	// We need to fetch CategoryName separately or assume it hasn't changed drastically,
	// but simpler to return without it for update or do another query if strictly needed.
	// For now let's keep it simple.
	return &p, nil
}

func (r *ProductRepository) Delete(id int) error {
	query := "UPDATE products SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL"
	result, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("product not found")
	}
	return nil
}

func (r *ProductRepository) CleanUpOldDeleted(duration time.Duration) error {
	threshold := time.Now().Add(-duration)
	query := "DELETE FROM products WHERE deleted_at < $1"
	_, err := r.db.Exec(query, threshold)
	return err
}
