package service

import (
    "context"

    "github.com/google/uuid"

    "github.com/rprajapati0067/quiz-game-backend/internal/models"
    "github.com/rprajapati0067/quiz-game-backend/internal/repository"
)

type AuthService interface {
    Signup(ctx context.Context, name, phone, email string) (*models.User, error)
    Login(ctx context.Context, phone string) (*models.User, error)
}

type authService struct {
    users repository.UserRepository
}

func NewAuthService(users repository.UserRepository) AuthService {
    return &authService{users: users}
}

func (s *authService) Signup(ctx context.Context, name, phone, email string) (*models.User, error) {
    u := &models.User{
        ID:       uuid.NewString(),
        Name:     name,
        Phone:    phone,
        Email:    email,
        Verified: false,
        Blocked:  false,
        Points:   0,
    }
    if err := s.users.CreateUser(ctx, u); err != nil {
        return nil, err
    }
    return u, nil
}

func (s *authService) Login(ctx context.Context, phone string) (*models.User, error) {
    // TODO: implement GetByPhone properly in repo
    return nil, nil
}
