package main

import (
	"cateogry-api/internal/handler"
	"cateogry-api/internal/repository"
	"cateogry-api/internal/service"
	"fmt"
	"net/http"
)

func main() {
	repo := repository.NewInMemoryCategoryRepository()
	svc := service.NewCategoryService(repo)
	h := handler.NewCategoryHandler(svc)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
