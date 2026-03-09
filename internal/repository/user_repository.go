// user-management-api/internal/repository/user_repository.go
package repository

import (
	"context"

	sqlc "github.com/thinhnguyenwilliam/user-management-api/internal/db/sqlc"
	"github.com/thinhnguyenwilliam/user-management-api/internal/models"
	"github.com/thinhnguyenwilliam/user-management-api/internal/models/mapper"
)

type userRepository struct {
	q sqlc.Querier
}

func NewUserRepository(q sqlc.Querier) IUserRepository {
	return &userRepository{
		q: q,
	}
}

func (r *userRepository) Create(
	ctx context.Context,
	user *models.User,
) error {

	params := mapper.ToCreateUserParams(user)

	result, err := r.q.CreateUser(ctx, params)
	if err != nil {
		return err
	}

	mapper.MapUserFromDB(result, user)

	return nil
}

// type userRepository struct {
// 	db *gorm.DB
// }

// func NewUserRepository(db *gorm.DB) IUserRepository {
// 	return &userRepository{
// 		db: db,
// 	}
// }

// func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
// 	var user models.User

// 	err := r.db.WithContext(ctx).
// 		Where("email = ?", email).
// 		First(&user).Error

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &user, nil
// }

// func (r *userRepository) Create(ctx context.Context, user *models.User) error {
// 	return r.db.WithContext(ctx).Create(user).Error
// }

// func (r *userRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
// 	var user models.User

// 	err := r.db.WithContext(ctx).
// 		First(&user, id).Error

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &user, nil
// }

// func (r *userRepository) Delete(ctx context.Context, id int) error {
// 	return r.db.WithContext(ctx).
// 		Delete(&models.User{}, id).Error
// }
