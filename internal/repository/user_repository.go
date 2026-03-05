// user-management-api/internal/repository/user_repository.go
package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/thinhnguyenwilliam/user-management-api/internal/models"
)

type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]*models.User
}

func NewUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*models.User),
	}
}

func (r *InMemoryUserRepository) FindByEmail(
	ctx context.Context,
	email string,
) (*models.User, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, errors.New("user not found")
}

func (r *InMemoryUserRepository) Create(ctx context.Context, user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.UUID.String()]; exists {
		return errors.New("user already exists")
	}

	r.users[user.UUID.String()] = user
	return nil
}
