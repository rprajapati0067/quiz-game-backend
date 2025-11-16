package repository

import (
    "context"

    "github.com/rprajapati0067/quiz-game-backend/internal/models"
)

type UserRepository interface {
    CreateUser(ctx context.Context, u *models.User) error
    GetByPhone(ctx context.Context, phone string) (*models.User, error)
    GetByID(ctx context.Context, id string) (*models.User, error)
    Update(ctx context.Context, u *models.User) error
}
