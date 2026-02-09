package service

import (
	"cateogry-api/internal/domain"
	"cateogry-api/internal/repository"
	"time"
)

type TransactionService struct {
	repo *repository.TransactionRepository
}

func NewTransactionService(repo *repository.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(items []domain.CheckoutItem) (*domain.Transaction, error) {
	return s.repo.CreateTransaction(items)
}

func (s *TransactionService) GetDailyReport(date time.Time) (domain.DailyReport, error) {
	return s.repo.GetDailyReport(date)
}

func (s *TransactionService) GetReport(startDate, endDate time.Time) (domain.DailyReport, error) {
	return s.repo.GetReport(startDate, endDate)
}
