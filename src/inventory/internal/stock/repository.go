package stock

import (
	"context"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, stock *Stock) error
	Update(ctx context.Context, stock *Stock) error
	FindByProductId(ctx context.Context, id string) (*Stock, error)
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *gorm.DB
}

func (r *repository) Create(ctx context.Context, stock *Stock) error {
	err := r.db.WithContext(ctx).Create(stock).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) Update(ctx context.Context, stock *Stock) error {
	err := r.db.WithContext(ctx).Model(stock).Where("id = ?", stock.ID).UpdateColumns(Stock{
		Quantity: stock.Quantity,
	}).Error
	if err != nil {
		return err
	}

	err = r.db.WithContext(ctx).First(stock, "id = ?", stock.ID).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) FindByProductId(ctx context.Context, id string) (*Stock, error) {
	var stock Stock
	err := r.db.WithContext(ctx).First(&stock, "product_id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &stock, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Where("product_id = ?", id).Delete(&Stock{}).Error
	if err != nil {
		return err
	}
	return nil
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
