package product

import (
	"context"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, product *Product) error
	Update(ctx context.Context, product *Product) error
	FindById(ctx context.Context, id string) (*Product, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, limit, page int) ([]Product, error)
}

type repository struct {
	db *gorm.DB
}

func (r *repository) Create(ctx context.Context, product *Product) error {
	err := r.db.WithContext(ctx).Create(product).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) Update(ctx context.Context, product *Product) error {
	err := r.db.WithContext(ctx).Model(product).Where("id = ?", product.ID).UpdateColumns(Product{
		Name:        product.Name,
		Price:       product.Price,
		Description: product.Description,
	}).Error
	if err != nil {
		return err
	}

	err = r.db.WithContext(ctx).First(product, "id = ?", product.ID).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) FindById(ctx context.Context, id string) (*Product, error) {
	var product Product
	err := r.db.WithContext(ctx).First(&product, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&Product{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context, limit, page int) ([]Product, error) {
	var products []Product
	err := r.db.WithContext(ctx).Raw(`
		SELECT products.* FROM products
		INNER JOIN stocks ON products.id = stocks.product_id
		WHERE stocks.quantity > 0
		ORDER BY products.updated_at DESC
		LIMIT ? 
		OFFSET ?;
		`, limit, page).Scan(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
