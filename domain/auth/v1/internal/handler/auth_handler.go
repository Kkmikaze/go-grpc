package handler

import (
	"context"
	"github.com/Kkmikaze/roketin/domain/auth/v1/usecase"
	authv1 "github.com/Kkmikaze/roketin/stubs/auth/v1"
)

type AuthHandler struct {
	authv1.UnimplementedAuthServiceServer
	UseCase usecase.AuthUseCase
}

func (h *AuthHandler) Check(ctx context.Context, in *authv1.HealthCheckRequest) (*authv1.HealthCheckResponse, error) {
	return &authv1.HealthCheckResponse{Message: "OK"}, nil
}

func NewAuthRestHandler(useCase usecase.AuthUseCase) authv1.AuthServiceServer {
	return &AuthHandler{UseCase: useCase}
}
