package catalog

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
)

type mockRepo struct{}

func (m *mockRepo) GetAllProducts(offset, limit int, category string, priceLt float64) ([]models.Product, int64, error) {
	return nil, 0, nil
}

func (m *mockRepo) GetByCode(code string) (*models.Product, error) {
	return &models.Product{
		Code:  "PROD001",
		Price: decimal.NewFromFloat(10.99),
		Category: models.Category{
			Name: "Clothing",
		},
		Variants: []models.Variant{
			{
				Name: "Variant A",
				SKU:  "SKU001A",
				Price: func() *decimal.Decimal {
					d := decimal.NewFromFloat(11.99)
					return &d
				}(),
			},
			{
				Name:  "Variant B",
				SKU:   "SKU001B",
				Price: nil,
			},
		},
	}, nil
}

func TestHandleGetByCode(t *testing.T) {
	repo := &mockRepo{}
	handler := NewCatalogHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/catalog/PROD001", nil)
	w := httptest.NewRecorder()

	handler.HandleGetByCode(w, req)

	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", res.StatusCode)
	}

	var body map[string]interface{}
	json.NewDecoder(res.Body).Decode(&body)

	if body["code"] != "PROD001" {
		t.Errorf("wrong product code")
	}
}

func (m *mockRepo) GetAllCategories() ([]models.Category, error) {
	return []models.Category{
		{Code: "CLOTHING", Name: "Clothing"},
		{Code: "SHOES", Name: "Shoes"},
	}, nil

}

func TestHandleGetCategories(t *testing.T) {
	repo := &mockRepo{}
	handler := NewCatalogHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	w := httptest.NewRecorder()

	handler.HandleGetCategories(w, req)

	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", res.StatusCode)
	}

	var response []CategoryResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("expected 2 categories, got %d", len(response))
	}

	if response[0].Code != "CLOTHING" {
		t.Errorf("unexpected category code")
	}
}

func (m *mockRepo) CreateCategory(category *models.Category) error {
	return nil
}
