package handlers

import (
    "context"

    auth "github.com/rprajapati0067/quiz-game-backend/rpc/auth"

    "github.com/rprajapati0067/quiz-game-backend/internal/service"
)

type AuthHandler struct {
    auth.UnimplementedAuthServiceServer
    svc service.AuthService
}

func NewAuthHandler(svc service.AuthService) *AuthHandler {
    return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Signup(ctx context.Context, req *auth.SignupRequest) (*auth.SignupResponse, error) {
    u, err := h.svc.Signup(ctx, req.Name, req.Phone, req.Email)
    if err != nil {
        return nil, err
    }
    return &auth.SignupResponse{UserId: u.ID}, nil
}

// TODO: implement Login and VerifyPhone using AuthService once repository supports it.
