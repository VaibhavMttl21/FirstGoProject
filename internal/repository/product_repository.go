package repository

import (
	"context"
	"firstgoproject/internal/model"
)


type ProductRepository interface{
	Create(ctx context.Context, product model.Product)(model.Product,error)
	GetByID(ctx context.Context, id string)(model.Product,error)
	GetAll(ctx context.Context)([]model.Product,error)
	Update(ctx context.Context, product model.Product)(model.Product,error)
	Delete(ctx context.Context, id string)(error)
}