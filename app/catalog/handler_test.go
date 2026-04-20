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
