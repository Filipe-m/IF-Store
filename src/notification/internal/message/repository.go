package message

import (
	"context"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, message *Message) error
}

type repository struct {
	db *gorm.DB
}

func (r *repository) Create(ctx context.Context, message *Message) error {
	err := r.db.WithContext(ctx).Create(message).Error
	if err != nil {
		return err
	}
	return nil
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
