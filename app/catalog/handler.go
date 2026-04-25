package catalog

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/mytheresa/go-hiring-challenge/app/api"
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

type VariantResponse struct {
	Name  string  `json:"name"`
	SKU   string  `json:"sku"`
	Price float64 `json:"price"`
}

type ProductDetailResponse struct {
	Code     string            `json:"code"`
	Price    float64           `json:"price"`
	Category string            `json:"category"`
	Variants []VariantResponse `json:"variants"`
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

	limit := 10
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

	res, total, err := h.repo.GetAllProducts(offset, limit, category, priceLt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

func (h *CatalogHandler) HandleGetByCode(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path // e.g. /catalog/PROD001
	parts := strings.Split(path, "/")

	if len(parts) < 3 {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	code := parts[2]

	product, err := h.repo.GetByCode(code)
	if err != nil {
		api.ErrorResponse(w, http.StatusNotFound, "product not found")
		return
	}

	variants := make([]VariantResponse, len(product.Variants))

	for i, v := range product.Variants {
		price := product.Price.InexactFloat64()

		if v.Price != nil {
			price = v.Price.InexactFloat64()
		}

		variants[i] = VariantResponse{
			Name:  v.Name,
			SKU:   v.SKU,
			Price: price,
		}
	}

	response := ProductDetailResponse{
		Code:     product.Code,
		Price:    product.Price.InexactFloat64(),
		Category: product.Category.Name,
		Variants: variants,
	}

	api.OKResponse(w, response)

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *CatalogHandler) HandleGetProducts(w http.ResponseWriter, r *http.Request) {

	offset := 0
	limit := 10

	// Parse offset
	if val := r.URL.Query().Get("offset"); val != "" {
		if v, err := strconv.Atoi(val); err == nil && v >= 0 {
			offset = v
		}
	}

	// Parse limit
	if val := r.URL.Query().Get("limit"); val != "" {
		if v, err := strconv.Atoi(val); err == nil {
			if v < 1 {
				v = 1
			}
			if v > 100 {
				v = 100
			}
			limit = v
		}
	}

	// Filters
	category := r.URL.Query().Get("category")

	priceLt := 0.0
	if val := r.URL.Query().Get("price_lt"); val != "" {
		if v, err := strconv.ParseFloat(val, 64); err == nil {
			priceLt = v
		}
	}

	products, total, err := h.repo.GetAllProducts(offset, limit, category, priceLt)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, "failed to fetch products")
		return
	}

	response := map[string]interface{}{
		"total":    total,
		"products": products,
	}

	api.OKResponse(w, response)
}
