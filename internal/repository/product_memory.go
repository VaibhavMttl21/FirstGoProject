package repository

import (
	"context"
	"firstgoproject/internal/model"
	"firstgoproject/internal/shared"
	"sync"
)

type InMemoryProductRepository struct {
	mu       sync.RWMutex
	products map[string]model.Product
}

func NewInMemoryProductRepository() *InMemoryProductRepository {
	return &InMemoryProductRepository{
		products: make(map[string]model.Product),
	}
}

func (r *InMemoryProductRepository) Create(ctx context.Context, product model.Product) (model.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	select {
	case <-ctx.Done():
		return model.Product{}, ctx.Err()
	default:
	}
	_,exists := r.products[product.ID]
	if exists{
		return model.Product{},shared.ErrProductAlreadyexists
	}
	r.products[product.ID] = product

	return product, nil
}

func (r *InMemoryProductRepository) GetByID(ctx context.Context, id string) (model.Product, error) {
	select {
	case <-ctx.Done():
		return model.Product{}, ctx.Err()
	default:
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	product,exists := r.products[id]
	if !exists{
		return model.Product{},shared.ErrProductNotFound
	}

	return product,nil
}

func (r *InMemoryProductRepository) GetAll(ctx context.Context)([]model.Product,error){
	products := make([]model.Product,0)
	select {
	case <-ctx.Done():
		return products, ctx.Err()
	default:
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	for _,product := range r.products{
		products = append(products,product)
	}

	return products,nil
}

func ( r *InMemoryProductRepository) Update(ctx context.Context, product model.Product)(model.Product,error){
	select {
	case <-ctx.Done():
		return model.Product{}, ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	id := product.ID

	_,exists := r.products[id]
	if(!exists){
		return model.Product{},shared.ErrProductNotFound
	}

	r.products[id] = product

	return r.products[id],nil
}

func (r *InMemoryProductRepository) Delete(ctx context.Context, id string) error{
	select {
	case <-ctx.Done():
		return  ctx.Err()
	default:
	}

	r.mu.Lock()
    defer r.mu.Unlock()

	_,exists := r.products[id]
	if(!exists){
		return shared.ErrProductNotFound
	}

	delete(r.products,id)
	return nil
}