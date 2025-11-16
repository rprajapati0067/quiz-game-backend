package handlers

import (
    "context"

    question "github.com/rprajapati0067/quiz-game-backend/rpc/question"

    "github.com/rprajapati0067/quiz-game-backend/internal/service"
)

type QuestionHandler struct {
    question.UnimplementedQuestionServiceServer
    svc service.QuestionService
}

func NewQuestionHandler(svc service.QuestionService) *QuestionHandler {
    return &QuestionHandler{svc: svc}
}

func (h *QuestionHandler) CreateQuestion(ctx context.Context, req *question.CreateQuestionRequest) (*question.CreateQuestionResponse, error) {
    q, err := h.svc.Create(ctx, req.Text, req.Options, req.CorrectIndex, req.Slot, "admin")
    if err != nil {
        return nil, err
    }
    return &question.CreateQuestionResponse{
        Question: &question.Question{
            Id:           q.ID,
            Text:         q.Text,
            Options:      q.Options,
            CorrectIndex: q.CorrectIndex,
            Slot:         q.Slot,
        },
    }, nil
}

func (h *QuestionHandler) ListQuestions(ctx context.Context, req *question.ListQuestionsRequest) (*question.ListQuestionsResponse, error) {
    qs, err := h.svc.ListBySlot(ctx, req.Slot)
    if err != nil {
        return nil, err
    }
    res := &question.ListQuestionsResponse{}
    for _, q := range qs {
        res.Questions = append(res.Questions, &question.Question{
            Id:           q.ID,
            Text:         q.Text,
            Options:      q.Options,
            CorrectIndex: q.CorrectIndex,
            Slot:         q.Slot,
        })
    }
    return res, nil
}

// SubmitAnswer is left as TODO for you to implement.
