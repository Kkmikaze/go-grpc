package usecase

import (
	"github.com/Kkmikaze/roketin/domain/movie/v1/internal/repository"
)

type MovieUseCase interface {
}

type MovieUseCaseImpl struct {
	Repository repository.MovieRepository
}

func NewMovieUseCase(movieRepository repository.MovieRepository) MovieUseCase {
	return &MovieUseCaseImpl{Repository: movieRepository}
}
