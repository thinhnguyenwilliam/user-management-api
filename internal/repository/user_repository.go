// user-management-api/internal/repository/user_repository.go
package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
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

func (r *InMemoryUserRepository) Create(ctx context.Context, user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.UUID]; exists {
		return errors.New("user already exists")
	}

	r.users[user.UUID] = user
	return nil
}

func (r *InMemoryUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id.String()]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}
