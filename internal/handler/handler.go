package handler

import (
	"cateogry-api/internal/domain"
	"cateogry-api/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(s *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: s}
}

func (h *CategoryHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/categories", h.categoriesHandler)
	mux.HandleFunc("/categories/", h.categoryHandler)
}

func (h *CategoryHandler) categoriesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getCategories(w, r)
	case http.MethodPost:
		h.createCategory(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *CategoryHandler) categoryHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getCategoryByID(w, r, id)
	case http.MethodPut:
		h.updateCategory(w, r, id)
	case http.MethodDelete:
		h.deleteCategory(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetCategories godoc
//
//	@Summary		Show all categories
//	@Description	Get all categories
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	domain.Category
//	@Router			/categories [get]
func (h *CategoryHandler) getCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// CreateCategory godoc
//
//	@Summary		Create a new category
//	@Description	Create a new category
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			category	body		domain.Category	true	"Category Data"
//	@Success		201			{object}	domain.Category
//	@Router			/categories [post]
func (h *CategoryHandler) createCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory domain.Category
	err := json.NewDecoder(r.Body).Decode(&newCategory)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	createdCategory, err := h.service.CreateCategory(newCategory)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdCategory)
}

// GetCategoryByID godoc
//
//	@Summary		Get a category by ID
//	@Description	Get a category by ID
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Category ID"
//	@Success		200	{object}	domain.Category
//	@Failure		404	{string}	string	"Category not found"
//	@Router			/categories/{id} [get]
func (h *CategoryHandler) getCategoryByID(w http.ResponseWriter, r *http.Request, id int) {
	category, err := h.service.GetCategoryByID(id)
	if err != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) updateCategory(w http.ResponseWriter, r *http.Request, id int) {
	var updatedData domain.Category
	err := json.NewDecoder(r.Body).Decode(&updatedData)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	updatedCategory, err := h.service.UpdateCategory(id, updatedData)
	if err != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedCategory)
}

func (h *CategoryHandler) deleteCategory(w http.ResponseWriter, r *http.Request, id int) {
	err := h.service.DeleteCategory(id)
	if err != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
