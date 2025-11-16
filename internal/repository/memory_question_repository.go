package repository

import (
	"context"
	"sync"

	"github.com/rprajapati0067/quiz-game-backend/internal/models"
)

type MemoryQuestionRepository struct {
	mu        sync.RWMutex
	questions map[string]*models.Question
}

func NewMemoryQuestionRepository() *MemoryQuestionRepository {
	return &MemoryQuestionRepository{
		questions: make(map[string]*models.Question),
	}
}

func (r *MemoryQuestionRepository) Create(ctx context.Context, q *models.Question) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.questions[q.ID] = q
	return nil
}

func (r *MemoryQuestionRepository) ListBySlot(ctx context.Context, slot int32) ([]*models.Question, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*models.Question
	for _, q := range r.questions {
		if q.Slot == slot {
			// Return a copy to avoid race conditions
			qCopy := *q
			result = append(result, &qCopy)
		}
	}
	return result, nil
}

func (r *MemoryQuestionRepository) GetByID(ctx context.Context, id string) (*models.Question, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	question, exists := r.questions[id]
	if !exists {
		return nil, nil
	}
	
	// Return a copy to avoid race conditions
	q := *question
	return &q, nil
}

