package service

import (
    "context"

    "github.com/rprajapati0067/quiz-game-backend/internal/models"
    "github.com/rprajapati0067/quiz-game-backend/internal/repository"
)

type UserService interface {
    GetByID(ctx context.Context, id string) (*models.User, error)
}

type userService struct {
    users repository.UserRepository
}

func NewUserService(users repository.UserRepository) UserService {
    return &userService{users: users}
}

func (s *userService) GetByID(ctx context.Context, id string) (*models.User, error) {
    return s.users.GetByID(ctx, id)
}
