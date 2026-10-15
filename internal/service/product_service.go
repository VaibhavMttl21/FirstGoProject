package service

import (
	"context"
	"firstgoproject/internal/model"
	"firstgoproject/internal/repository"
	"firstgoproject/internal/shared"
	"strings"
	"time"
)

type ProductService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func validateProduct(product model.Product) error {
	if strings.TrimSpace(product.Name) == "" {
		return shared.ErrProductIsEmpty
	}
	if product.PriceCents < 0 {
		return shared.ErrProductPriceNegative
	}
	if product.StockQuantity < 0 {
		return shared.ErrStockIsNegative
	}
	currency := product.Currency

	if len(currency) != 3 {
		return shared.ErrCurrencyInvalid
	}

	for _, ch := range currency {
		if ch < 'A' || ch > 'Z' {
			return shared.ErrCurrencyInvalid
		}
	}
	return nil
}

func (s *ProductService) Create(ctx context.Context, product model.Product) (model.Product, error) {
	if err := ctx.Err(); err != nil {
		return model.Product{}, err
	}
	err := validateProduct(product)
	if err != nil {
		return model.Product{}, err
	}
	id, err := newProductID()
	if err != nil {
		return model.Product{}, err
	}
	product.ID = id
	now := time.Now().UTC()
	product.CreatedAt = now
	product.UpdatedAt = now

	return s.repo.Create(ctx, product)
}

func (s *ProductService) GetByID(ctx context.Context, id string) (model.Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return model.Product{}, err
	}
	return product, nil
}

func (s *ProductService) GetAll(ctx context.Context) ([]model.Product, error) {
	return s.repo.GetAll(ctx)
}

func (s *ProductService) Update(ctx context.Context, product model.Product) (model.Product, error) {
	if err := ctx.Err(); err != nil {
		return model.Product{}, err
	}
	err := validateProduct(product)
	if err != nil {
		return model.Product{}, err
	}
	existing, err := s.repo.GetByID(ctx, product.ID)
	if err != nil {
		return model.Product{}, err
	}

	product.CreatedAt = existing.CreatedAt
	product.UpdatedAt = time.Now().UTC()

	return s.repo.Update(ctx, product)
}

func (s *ProductService) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
