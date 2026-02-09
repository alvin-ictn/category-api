package repository

import (
	"cateogry-api/internal/domain"
	"database/sql"
	"fmt"
	"time"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransaction(items []domain.CheckoutItem) (*domain.Transaction, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	println("TransactionRepository.CreateTransaction", items)
	totalAmount := 0
	details := make([]domain.TransactionDetail, 0)

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		if stock < item.Quantity {
			return nil, fmt.Errorf("stock for product %s is not enough. available: %d, requested: %d", productName, stock, item.Quantity)
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, domain.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName, // Optional, for response
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	for i := range details {
		details[i].TransactionID = transactionID
		// FIX: Use correctly indexed access as per instruction
		_, err = tx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)",
			transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &domain.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		CreatedAt:   time.Now(), // This is approximate, actual time is in DB
		Details:     details,
	}, nil
}

func (r *TransactionRepository) GetDailyReport(date time.Time) (domain.DailyReport, error) {
	// Start of the day
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	// End of the day
	endOfDay := startOfDay.Add(24 * time.Hour)

	return r.GetReport(startOfDay, endOfDay)
}

func (r *TransactionRepository) GetReport(startDate, endDate time.Time) (domain.DailyReport, error) {
	var report domain.DailyReport

	// 1. Total Revenue and Total Transactions
	queryRevenue := `
		SELECT COALESCE(SUM(total_amount), 0), COUNT(id)
		FROM transactions
		WHERE created_at >= $1 AND created_at < $2
	`
	err := r.db.QueryRow(queryRevenue, startDate, endDate).Scan(&report.TotalRevenue, &report.TotalTransactions)
	if err != nil {
		return report, err
	}

	// 2. Best Selling Product
	queryBestSeller := `
		SELECT p.name, COALESCE(SUM(td.quantity), 0) as qty_sold
		FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		JOIN products p ON td.product_id = p.id
		WHERE t.created_at >= $1 AND t.created_at < $2
		GROUP BY p.name
		ORDER BY qty_sold DESC
		LIMIT 1
	`
	err = r.db.QueryRow(queryBestSeller, startDate, endDate).Scan(&report.BestSellingProduct.Name, &report.BestSellingProduct.QtySold)
	if err == sql.ErrNoRows {
		// No transactions, so no best seller
		report.BestSellingProduct = domain.BestSellingProduct{Name: "-", QtySold: 0}
	} else if err != nil {
		return report, err
	}

	return report, nil
}
