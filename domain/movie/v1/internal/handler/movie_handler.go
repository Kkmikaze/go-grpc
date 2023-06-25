package handler

import (
	"context"
	"github.com/Kkmikaze/roketin/domain/movie/v1/usecase"
	authv1 "github.com/Kkmikaze/roketin/stubs/auth/v1"
	moviev1 "github.com/Kkmikaze/roketin/stubs/movie/v1"
)

type MovieHandler struct {
	moviev1.UnimplementedMovieServiceServer
	UseCase     usecase.MovieUseCase
	AuthService authv1.AuthServiceClient
}

func (h *MovieHandler) Check(ctx context.Context, in *moviev1.HealthCheckRequest) (*moviev1.HealthCheckResponse, error) {
	return &moviev1.HealthCheckResponse{Message: "OK"}, nil
}

func NewMovieRestHandler(useCase usecase.MovieUseCase, authConn authv1.AuthServiceClient) moviev1.MovieServiceServer {
	return &MovieHandler{UseCase: useCase, AuthService: authConn}
}
