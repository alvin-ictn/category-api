package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type HealthHandler struct {
	db *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

type HealthResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// CheckHealth godoc
// @Summary      Health Check
// @Description  Check database connection status
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  HealthResponse
// @Failure      503  {object}  HealthResponse
// @Router       /health [get]
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	err := h.db.Ping()
	response := HealthResponse{Status: "healthy"}
	statusCode := http.StatusOK

	if err != nil {
		response.Status = "unhealthy"
		response.Error = err.Error()
		statusCode = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
