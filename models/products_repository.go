package models

import (
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetAllProducts(offset, limit int, category string, priceLt float64) ([]Product, int64, error)
	GetByCode(code string) (*Product, error)
}

type productsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) ProductRepository {
	return &productsRepository{
		db: db,
	}
}

func (r *productsRepository) GetAllProducts(offset, limit int, category string, priceLt float64) ([]Product, int64, error) {
	var products []Product
	var total int64

	// 1. Build base query
	db := r.db.Model(&Product{}).
		Preload("Variants").
		Preload("Category")

	// 2. Apply filters
	if category != "" {
		db = db.Joins("Category").Where("categories.name = ?", category)
	}

	if priceLt > 0 {
		db = db.Where("price < ?", priceLt)
	}

	// 3. Count (filtered)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 4. Fetch with pagination
	if err := db.
		Offset(offset).
		Limit(limit).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productsRepository) GetByCode(code string) (*Product, error) {
	var product Product

	if err := r.db.
		Preload("Variants").
		Preload("Category").
		Where("code = ?", code).
		First(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}
