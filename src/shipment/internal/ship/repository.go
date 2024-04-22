package ship

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"shipment/internal/platform"
)

type Repository interface {
	Create(ctx context.Context, shipment *Shipment) error
	FindByOrderId(ctx context.Context, id string) (*Shipment, error)
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *gorm.DB
}

func (r *repository) Create(ctx context.Context, shipment *Shipment) error {
	err := r.db.WithContext(ctx).Preload(clause.Associations).Create(shipment).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) FindByOrderId(ctx context.Context, id string) (*Shipment, error) {
	var shipment Shipment
	err := r.db.WithContext(ctx).Preload(clause.Associations).First(&shipment, "order_id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &shipment, nil
}

func (r *repository) Delete(ctx context.Context, id string) (err error) {
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		err = platform.CommitOrRollback(tx, err)
	}()

	err = tx.WithContext(ctx).Where("shipment_id = ?", id).Delete(&Item{}).Error
	if err != nil {
		return
	}

	err = tx.WithContext(ctx).Where("id = ?", id).Delete(&Shipment{}).Error
	if err != nil {
		return
	}

	return
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
