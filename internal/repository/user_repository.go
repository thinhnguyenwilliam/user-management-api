// user-management-api/internal/repository/user_repository.go
package repository

import (
	"context"

	"github.com/thinhnguyenwilliam/user-management-api/internal/models"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User

	err := r.db.WithContext(ctx).
		First(&user, id).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).
		Delete(&models.User{}, id).Error
}
