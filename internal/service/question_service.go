package service

import (
    "context"

    "github.com/google/uuid"

    "github.com/rprajapati0067/quiz-game-backend/internal/models"
    "github.com/rprajapati0067/quiz-game-backend/internal/repository"
)

type QuestionService interface {
    Create(ctx context.Context, text string, options []string, correctIndex, slot int32, createdBy string) (*models.Question, error)
    ListBySlot(ctx context.Context, slot int32) ([]*models.Question, error)
}

type questionService struct {
    repo repository.QuestionRepository
}

func NewQuestionService(repo repository.QuestionRepository) QuestionService {
    return &questionService{repo: repo}
}

func (s *questionService) Create(ctx context.Context, text string, options []string, correctIndex, slot int32, createdBy string) (*models.Question, error) {
    q := &models.Question{
        ID:           uuid.NewString(),
        Text:         text,
        Options:      options,
        CorrectIndex: correctIndex,
        Slot:         slot,
        CreatedBy:    createdBy,
    }
    if err := s.repo.Create(ctx, q); err != nil {
        return nil, err
    }
    return q, nil
}

func (s *questionService) ListBySlot(ctx context.Context, slot int32) ([]*models.Question, error) {
    return s.repo.ListBySlot(ctx, slot)
}
