package repository

import (
    "context"

    "github.com/rprajapati0067/quiz-game-backend/internal/models"
)

type QuestionRepository interface {
    Create(ctx context.Context, q *models.Question) error
    ListBySlot(ctx context.Context, slot int32) ([]*models.Question, error)
    GetByID(ctx context.Context, id string) (*models.Question, error)
}
