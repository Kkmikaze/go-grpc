package handler

import (
	"context"
	"github.com/Kkmikaze/roketin/common"
	"github.com/Kkmikaze/roketin/domain/movie/v1/internal/entity"
	"github.com/Kkmikaze/roketin/domain/movie/v1/internal/validator"
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

func (h *MovieHandler) CreateMovie(ctx context.Context, in *moviev1.CreateMovieRequest) (*moviev1.Response, error) {
	if err := common.ValidateRequest(&validator.MovieValidator{
		Title:       in.Title,
		Description: in.Description,
		Duration:    in.Duration,
		Artist:      in.Artist,
		Genre:       in.Genre,
		VideoUrl:    in.VideoUrl,
	}); err != nil {
		return nil, err
	}

	payload := entity.Movie{
		Title:       in.Title,
		Description: in.Description,
		Duration:    in.Duration,
		Artist:      in.Artist,
		Genre:       in.Genre,
		VideoUrl:    in.VideoUrl,
	}
	err := h.UseCase.CreateMovie(ctx, &payload)
	if err != nil {
		return nil, err
	}

	return &moviev1.Response{
		Status:  true,
		Message: "Create Movie Success",
	}, nil
}

func (h *MovieHandler) GetMovie(ctx context.Context, in *moviev1.GetMovieRequest) (*moviev1.GetMovieResponse, error) {
	row, err := h.UseCase.GetMovie(ctx, in)
	if err != nil {
		return nil, err
	}

	return &moviev1.GetMovieResponse{
		Status:  true,
		Message: "Get Movie Successfully.",
		Data:    row,
	}, nil
}

func (h *MovieHandler) GetMovieByID(ctx context.Context, in *moviev1.ParamID) (*moviev1.GetMovieByIDResponse, error) {
	row, err := h.UseCase.GetMovieByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &moviev1.GetMovieByIDResponse{
		Status:  true,
		Message: "Get Movie By ID Successfully",
		Data:    row,
	}, nil
}

func (h *MovieHandler) UpdateMovie(ctx context.Context, in *moviev1.MovieData) (*moviev1.UpdateMovieResponse, error) {
	if err := common.ValidateRequest(&validator.MovieValidator{
		Title:       in.Title,
		Description: in.Description,
		Duration:    in.Duration,
		Artist:      in.Artist,
		Genre:       in.Genre,
		VideoUrl:    in.VideoUrl,
	}); err != nil {
		return nil, err
	}
	payload := entity.Movie{
		Id:          in.Id,
		Title:       in.Title,
		Description: in.Description,
		Duration:    in.Duration,
		Artist:      in.Artist,
		Genre:       in.Genre,
		VideoUrl:    in.VideoUrl,
	}

	row, err := h.UseCase.UpdateMovie(ctx, &payload)
	if err != nil {
		return nil, err
	}

	return &moviev1.UpdateMovieResponse{
		Status:  true,
		Message: "Update Movie Success",
		Data:    row,
	}, nil
}

func (h *MovieHandler) DeleteMovie(ctx context.Context, in *moviev1.ParamID) (*moviev1.Response, error) {
	err := h.UseCase.DeleteMovie(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &moviev1.Response{
		Status:  true,
		Message: "Delete Movie Success",
	}, nil
}

func NewMovieRestHandler(useCase usecase.MovieUseCase, authConn authv1.AuthServiceClient) moviev1.MovieServiceServer {
	return &MovieHandler{UseCase: useCase, AuthService: authConn}
}
