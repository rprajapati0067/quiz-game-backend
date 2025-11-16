package repository

import (
    "context"

    "github.com/rprajapati0067/quiz-game-backend/internal/models"
)

type AwardRepository interface {
    List(ctx context.Context) ([]*models.Award, error)
    GetByID(ctx context.Context, id string) (*models.Award, error)
    CreateClaim(ctx context.Context, c *models.Claim) error
}
