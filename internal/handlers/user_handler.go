package handlers

import (
    "context"

    user "github.com/rprajapati0067/quiz-game-backend/rpc/user"

    "github.com/rprajapati0067/quiz-game-backend/internal/service"
)

type UserHandler struct {
    user.UnimplementedUserServiceServer
    svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
    return &UserHandler{svc: svc}
}

func (h *UserHandler) Me(ctx context.Context, req *user.MeRequest) (*user.MeResponse, error) {
    // TODO: extract user id from auth (e.g. JWT). For now this is just a stub.
    return &user.MeResponse{}, nil
}
