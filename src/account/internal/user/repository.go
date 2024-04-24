package user

import (
	"context"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	FindById(ctx context.Context, id string) (*User, error)
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *gorm.DB
}

func (r *repository) Create(ctx context.Context, user *User) error {
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) Update(ctx context.Context, user *User) error {
	err := r.db.WithContext(ctx).Model(user).UpdateColumns(User{
		Username: user.Username,
		Email:    user.Email,
	}).Error
	if err != nil {
		return err
	}

	err = r.db.WithContext(ctx).First(user, "id = ?", user.ID).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) FindById(ctx context.Context, id string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&User{}).Error
	if err != nil {
		return err
	}
	return nil
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
