package catalog

import (
	"encoding/json"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/models"
)

type CategoryResponse struct {
	Code string `json:"code,omitempty"`
	Name string `json:"name,omitempty"`
}

func (h *CatalogHandler) HandleGetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.repo.GetAllCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]CategoryResponse, len(categories))

	for i, c := range categories {
		response[i] = CategoryResponse{
			Code: c.Code,
			Name: c.Name,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func (h *CatalogHandler) HandleCreateCategory(w http.ResponseWriter, r *http.Request) {
	var req CategoryResponse

	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	category := models.Category{
		Code: req.Code,
		Name: req.Name,
	}

	if err := h.repo.CreateCategory(&category); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
