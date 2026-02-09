package handler

import (
	"cateogry-api/internal/domain"
	"cateogry-api/internal/service"
	"encoding/json"
	"net/http"
	"time"
)

type TransactionHandler struct {
	service *service.TransactionService
}

func NewTransactionHandler(service *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/checkout", h.handleCheckout)
	mux.HandleFunc("/report/hari-ini", h.handleDailyReport)
	mux.HandleFunc("/report", h.handleReport)
}

func (h *TransactionHandler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req domain.CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	println("TransactionHandler.handleCheckout", req.Items)
	transaction, err := h.service.Checkout(req.Items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	println("TransactionHandler.handleCheckout", transaction)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

func (h *TransactionHandler) handleDailyReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Assuming today's date
	report, err := h.service.GetDailyReport(time.Now())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

func (h *TransactionHandler) handleReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	if startDateStr == "" || endDateStr == "" {
		http.Error(w, "start_date and end_date are required", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		http.Error(w, "Invalid start_date format (YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		http.Error(w, "Invalid end_date format (YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	// Add 1 day to end_date to include the full day if user inputs only date part
	// But since time.Parse returns 00:00:00, we check specifically if we need to include end date
	// Usually reports assume inclusive end date or exclusive
	// Instruction example: start_date=2026-01-01&end_date=2026-02-01. Usually implies up to 2026-02-01.
	// If we assume exclusive upper bound: [start, end)

	report, err := h.service.GetReport(startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
