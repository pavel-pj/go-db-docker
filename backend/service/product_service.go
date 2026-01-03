// backend/service/product_service.go
package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	productsdb "db200/internal/db/products"
	"db200/internal/store"
)

type ProductService struct {
	store *store.ProductStore
}

func NewProductService(productStore *store.ProductStore) *ProductService {
	return &ProductService{
		store: productStore,
	}
}

// Ошибки объявляются в сервисе
var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
)

type CreateProductInput struct {
	Slug        string
	Title       string
	Description string
	PriceCents  int32
}

func (s *ProductService) Create(ctx context.Context, input CreateProductInput) (productsdb.Product, error) {
	// Валидация
	if input.Slug == "" {
		return productsdb.Product{}, fmt.Errorf("slug is required")
	}
	if input.PriceCents <= 0 {
		return productsdb.Product{}, fmt.Errorf("price must be positive")
	}

	product, err := s.store.Create(ctx, productsdb.CreateProductParams{
		Slug:        input.Slug,
		Title:       input.Title,
		Description: input.Description,
		PriceCents:  input.PriceCents,
	})
	if err != nil {
		return product, fmt.Errorf("create: %w", err)
	}
	return product, nil
}

func (s *ProductService) Get(ctx context.Context, id int32) (productsdb.Product, error) {

	product, err := s.store.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return product, fmt.Errorf("product %d not found: %w", id, ErrNotFound)
		}
		return product, fmt.Errorf("get product %d: %w", id, err)
	}
	return product, nil
}

func (s *ProductService) List(ctx context.Context, limit, offset int32) ([]productsdb.Product, error) {
	// Бизнес-правила для пагинации
	if limit <= 0 {
		limit = 10 // дефолтное значение
	}
	if limit > 50 {
		return nil, fmt.Errorf("service: list products: %w: limit too large %d",
			ErrInvalidInput, limit)
	}

	products, err := s.store.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("service: list products: %w", err)
	}

	// Можно добавить бизнес-логику (фильтрацию, форматирование и т.д.)
	return products, nil
}

func (s *ProductService) UpdatePrice(ctx context.Context, id, priceCents int32) error {
	// Бизнес-правила
	if id <= 0 {
		return fmt.Errorf("service: update price: %w: invalid id %d",
			ErrInvalidInput, id)
	}
	if priceCents <= 0 {
		return fmt.Errorf("service: update price: %w: price must be positive, got %d",
			ErrInvalidInput, priceCents)
	}

	rows, err := s.store.UpdatePrice(ctx, id, priceCents)
	if err != nil {
		return fmt.Errorf("service: update price: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("service: update price: %w: product %d not found",
			ErrNotFound, id)
	}

	return nil
}

func (s *ProductService) Delete(ctx context.Context, id int32) error {
	// Бизнес-правила
	if id <= 0 {
		return fmt.Errorf("service: delete product: %w: invalid id %d",
			ErrInvalidInput, id)
	}

	rows, err := s.store.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("service: delete product: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("service: delete product: %w: product %d not found",
			ErrNotFound, id)
	}

	return nil
}
