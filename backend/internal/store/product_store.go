// backend/internal/store/store.go
package store

import (
	"context"
	"database/sql"
	"fmt"

	productsdb "db200/internal/db/products"
)

// ProductStore - хранилище ТОЛЬКО для продуктов
type ProductStore struct {
	queries *productsdb.Queries
}

// NewProductStore создает новый ProductStore
func NewProductStore(db *sql.DB) *ProductStore {
	return &ProductStore{
		queries: productsdb.New(db),
	}
}

// Create создает новый продукт
func (s *ProductStore) Create(ctx context.Context, params productsdb.CreateProductParams) (productsdb.Product, error) {
	product, err := s.queries.CreateProduct(ctx, params)
	if err != nil {
		return product, err
		//return product, fmt.Errorf("create product: %w", err)
	}
	return product, nil
}

// Create создает новый продукт
func (s *ProductStore) Get(ctx context.Context, id int32) (productsdb.Product, error) {
	product, err := s.queries.GetProductByID(ctx, id)
	if err != nil {

		return product, err
	}
	return product, nil
}

// List возвращает список продуктов с пагинацией
func (s *ProductStore) List(ctx context.Context, limit, offset int32) ([]productsdb.Product, error) {
	// Валидация пагинации
	if limit < 0 || offset < 0 {
		return nil, fmt.Errorf("store: invalid pagination: limit=%d, offset=%d", limit, offset)
	}
	if limit > 100 { // защита от слишком больших лимитов
		limit = 100
	}

	products, err := s.queries.ListProducts(ctx, productsdb.ListProductsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("store: list products: %w", err)
	}

	return products, nil
}

// UpdatePrice обновляет цену продукта
func (s *ProductStore) UpdatePrice(ctx context.Context, id, priceCents int32) (int64, error) {
	// Валидация
	if id <= 0 {
		return 0, fmt.Errorf("store: invalid product id: %d", id)
	}
	if priceCents < 0 {
		return 0, fmt.Errorf("store: price cannot be negative: %d", priceCents)
	}

	rows, err := s.queries.UpdateProductPrice(ctx, productsdb.UpdateProductPriceParams{
		ID:         id,
		PriceCents: priceCents,
	})
	if err != nil {
		return 0, fmt.Errorf("store: update product price %d: %w", id, err)
	}

	return rows, nil
}

// Delete удаляет продукт
func (s *ProductStore) Delete(ctx context.Context, id int32) (int64, error) {
	// Валидация
	if id <= 0 {
		return 0, fmt.Errorf("store: invalid product id: %d", id)
	}

	rows, err := s.queries.DeleteProduct(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("store: delete product %d: %w", id, err)
	}

	return rows, nil
}
