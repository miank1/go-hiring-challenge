package catalog

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mytheresa/go-hiring-challenge/models"
)

type Response struct {
	Total    int64     `json:"total"`
	Products []Product `json:"products"`
}

type Product struct {
	Code     string  `json:"code"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

type CatalogHandler struct {
	repo models.ProductRepository
}

func NewCatalogHandler(r models.ProductRepository) *CatalogHandler {
	return &CatalogHandler{
		repo: r,
	}
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	offsetStr := r.URL.Query().Get("offset")

	offset := 0
	if offsetStr != "" {
		if v, err := strconv.Atoi(offsetStr); err == nil && v >= 0 {
			offset = v
		}
	}

	limitStr := r.URL.Query().Get("limit")

	limit := 10 // default
	if limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil {
			if v >= 1 && v <= 100 {
				limit = v
			}
		}
	}

	category := r.URL.Query().Get("category")
	priceLtStr := r.URL.Query().Get("price_lt")

	priceLt := 0.0
	if priceLtStr != "" {
		val, err := strconv.ParseFloat(priceLtStr, 64)
		if err == nil {
			priceLt = val
		}
	}

	// _ = category
	// _ = priceLtStr

	res, total, err := h.repo.GetAllProducts(offset, limit, category, priceLt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Map response
	products := make([]Product, len(res))
	for i, p := range res {
		products[i] = Product{
			Code:     p.Code,
			Price:    p.Price.InexactFloat64(),
			Category: p.Category.Name,
		}
	}

	// Return the products as a JSON response
	w.Header().Set("Content-Type", "application/json")

	response := Response{
		Total:    total,
		Products: products,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
