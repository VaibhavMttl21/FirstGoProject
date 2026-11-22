package repository

import (
	"context"
	"errors"
	"firstgoproject/internal/model"
	"firstgoproject/internal/shared"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type GormProductRepository struct {
	db *gorm.DB
}

func NewGormProductRepository(db *gorm.DB) *GormProductRepository {
	return &GormProductRepository{
		db: db,
	}
}

func (r *GormProductRepository) Create(
	ctx context.Context,
	product model.Product,
) (model.Product, error) {

	if err := ctx.Err(); err != nil {
		return model.Product{}, err
	}

	err := r.db.WithContext(ctx).Create(&product).Error

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.Product{}, shared.ErrProductAlreadyexists
		}

		return model.Product{}, err
	}

	return product, nil
}

func (r *GormProductRepository) GetByID(ctx context.Context, id string) (model.Product, error) {
	select {
	case <-ctx.Done():
		return model.Product{}, ctx.Err()
	default:
	}

	var product model.Product

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&product).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Product{}, shared.ErrProductNotFound
		}
		return model.Product{}, err
	}

	return product, nil
}

func (r *GormProductRepository) GetAll(ctx context.Context) ([]model.Product, error) {
	products := make([]model.Product, 0)

	select {
	case <-ctx.Done():
		return products, ctx.Err()
	default:
	}

	err := r.db.WithContext(ctx).Find(&products).Error

	if err != nil {
		return products, err
	}

	return products, nil
}

func (r *GormProductRepository) Update(ctx context.Context, product model.Product) (model.Product, error) {
	select {
	case <-ctx.Done():
		return model.Product{}, ctx.Err()
	default:
	}

	result := r.db.WithContext(ctx).
		Model(&model.Product{}).
		Where("id = ?", product.ID).
		Select("*").
		Updates(&product)

	if result.Error != nil {
		return model.Product{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.Product{}, shared.ErrProductNotFound
	}

	return product, nil
}

func (r *GormProductRepository) Delete(ctx context.Context, id string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.Product{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return shared.ErrProductNotFound
	}

	return nil
}
