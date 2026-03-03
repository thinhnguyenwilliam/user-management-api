// user-management-api/internal/repository/user_repository.go
package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/thinhnguyenwilliam/user-management-api/internal/models"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
}

type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[int64]*models.User
}

func NewUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[int64]*models.User),
	}
}

func (r *InMemoryUserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *InMemoryUserRepository) Create(ctx context.Context, user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; exists {
		return errors.New("user already exists")
	}

	r.users[user.ID] = user
	return nil
}
