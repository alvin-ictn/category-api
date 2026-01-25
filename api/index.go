package handler

import (
	"cateogry-api/internal/handler"
	"cateogry-api/internal/repository"
	"cateogry-api/internal/service"
	"net/http"
)

var mux *http.ServeMux

func init() {
	// Initialize the application components once
	repo := repository.NewInMemoryCategoryRepository()
	svc := service.NewCategoryService(repo)
	h := handler.NewCategoryHandler(svc)

	mux = http.NewServeMux()
	h.RegisterRoutes(mux)
}

// Handler is the entry point for Vercel Serverless Functions
func Handler(w http.ResponseWriter, r *http.Request) {
	mux.ServeHTTP(w, r)
}
