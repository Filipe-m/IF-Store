package order

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"order/internal/platform"
)

type Repository interface {
	Create(ctx context.Context, order *Order) error
	FindByOrderId(ctx context.Context, id string) (*Order, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, order *Order) error
	ExistOrderByStatus(ctx context.Context, userID string, status Status) (*Order, error)
	FindActualByUserId(ctx context.Context, id string) (*Order, error)
	DeleteItem(ctx context.Context, id string, order *Order) (err error)
}

type repository struct {
	db *gorm.DB
}

func (r *repository) Create(ctx context.Context, order *Order) error {
	err := r.db.WithContext(ctx).Preload(clause.Associations).Create(order).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) Update(ctx context.Context, order *Order) (err error) {
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		err = platform.CommitOrRollback(tx, err)
	}()

	err = tx.WithContext(ctx).Save(order).Error
	if err != nil {
		return err
	}

	for i := range order.Items {
		order.Items[i].OrderID = order.ID
	}

	err = tx.WithContext(ctx).Save(order.Items).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) ExistOrderByStatus(ctx context.Context, userID string, status Status) (*Order, error) {
	var order Order
	err := r.db.WithContext(ctx).Where("user_id = ? AND status = ?", userID, status).
		Preload(clause.Associations).First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &order, nil
}

func (r *repository) FindByOrderId(ctx context.Context, id string) (*Order, error) {
	var order Order
	err := r.db.WithContext(ctx).Preload(clause.Associations).First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *repository) FindActualByUserId(ctx context.Context, id string) (*Order, error) {
	var order Order
	err := r.db.WithContext(ctx).Preload(clause.Associations).Order("created_at desc").First(&order, "user_id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *repository) Delete(ctx context.Context, id string) (err error) {
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		err = platform.CommitOrRollback(tx, err)
	}()

	err = tx.WithContext(ctx).Where("order_id = ?", id).Delete(&Item{}).Error
	if err != nil {
		return
	}

	err = tx.WithContext(ctx).Where("id = ?", id).Delete(&Order{}).Error
	if err != nil {
		return
	}

	return
}

func (r *repository) DeleteItem(ctx context.Context, id string, order *Order) (err error) {
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		err = platform.CommitOrRollback(tx, err)
	}()

	err = tx.WithContext(ctx).Save(order).Error
	if err != nil {
		return err
	}

	err = tx.WithContext(ctx).Where("id = ?", id).Delete(&Item{}).Error
	if err != nil {
		return
	}

	return
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
