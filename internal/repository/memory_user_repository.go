package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/rprajapati0067/quiz-game-backend/internal/models"
)

type MemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]*models.User
	byPhone map[string]*models.User
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users:   make(map[string]*models.User),
		byPhone: make(map[string]*models.User),
	}
}

func (r *MemoryUserRepository) CreateUser(ctx context.Context, u *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[u.ID]; exists {
		return errors.New("user already exists")
	}
	if _, exists := r.byPhone[u.Phone]; exists {
		return errors.New("phone number already registered")
	}

	r.users[u.ID] = u
	r.byPhone[u.Phone] = u
	return nil
}

func (r *MemoryUserRepository) GetByPhone(ctx context.Context, phone string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.byPhone[phone]
	if !exists {
		return nil, nil
	}
	
	// Return a copy to avoid race conditions
	u := *user
	return &u, nil
}

func (r *MemoryUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, nil
	}
	
	// Return a copy to avoid race conditions
	u := *user
	return &u, nil
}

func (r *MemoryUserRepository) Update(ctx context.Context, u *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[u.ID]; !exists {
		return errors.New("user not found")
	}

	r.users[u.ID] = u
	r.byPhone[u.Phone] = u
	return nil
}

