package usecase

import (
	"context"
	"github.com/Kkmikaze/roketin/domain/movie/v1/internal/entity"
	"github.com/Kkmikaze/roketin/domain/movie/v1/internal/repository"
	"github.com/Kkmikaze/roketin/pkg/orm"
	moviev1 "github.com/Kkmikaze/roketin/stubs/movie/v1"
)

type MovieUseCase interface {
	CreateMovie(ctx context.Context, e *entity.Movie) error
	GetMovie(ctx context.Context, query *moviev1.GetMovieRequest) (*moviev1.GetMovieData, error)
	GetMovieByID(ctx context.Context, id uint64) (*moviev1.MovieData, error)
	UpdateMovie(ctx context.Context, e *entity.Movie) (*moviev1.MovieData, error)
	DeleteMovie(ctx context.Context, id uint64) error
}

type MovieUseCaseImpl struct {
	Repository repository.MovieRepository
}

func (u *MovieUseCaseImpl) CreateMovie(ctx context.Context, e *entity.Movie) error {

	if err := u.Repository.CreateMovie(ctx, e); err != nil {
		return err
	}

	return nil
}

func (u *MovieUseCaseImpl) GetMovie(ctx context.Context, query *moviev1.GetMovieRequest) (*moviev1.GetMovieData, error) {
	var results []*moviev1.MovieData

	builder := &orm.QueryBuilder{
		Search:      query.Search,
		Page:        int(query.Page),
		ItemPerPage: int(query.ItemPerPage),
	}

	rows, total, err := u.Repository.GetMovie(ctx, builder)
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		results = append(results, &moviev1.MovieData{
			Id:          row.Id,
			Title:       row.Title,
			Description: row.Description,
			Duration:    row.Duration,
			Artist:      row.Artist,
			Genre:       row.Genre,
			VideoUrl:    row.VideoUrl,
		})
	}

	return &moviev1.GetMovieData{
		Items: results,
		Total: uint64(total),
	}, nil
}

func (u *MovieUseCaseImpl) GetMovieByID(ctx context.Context, id uint64) (*moviev1.MovieData, error) {
	row, err := u.Repository.GetMovieByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &moviev1.MovieData{
		Id:          row.Id,
		Title:       row.Title,
		Description: row.Description,
		Duration:    row.Duration,
		Artist:      row.Artist,
		Genre:       row.Genre,
		VideoUrl:    row.VideoUrl,
	}, nil
}

func (u *MovieUseCaseImpl) UpdateMovie(ctx context.Context, e *entity.Movie) (*moviev1.MovieData, error) {
	_, err := u.Repository.GetMovieByID(ctx, e.Id)
	if err != nil {
		return nil, err
	}

	row, err := u.Repository.UpdateMovie(ctx, e)
	if err != nil {
		return nil, err
	}
	return &moviev1.MovieData{
		Id:          row.Id,
		Title:       row.Title,
		Description: row.Description,
		Duration:    row.Duration,
		Artist:      row.Artist,
		Genre:       row.Genre,
		VideoUrl:    row.VideoUrl,
	}, nil
}

func (u *MovieUseCaseImpl) DeleteMovie(ctx context.Context, id uint64) error {
	_, err := u.Repository.GetMovieByID(ctx, id)
	if err != nil {
		return err
	}

	if err = u.Repository.DeleteMovie(ctx, id); err != nil {
		return err
	}

	return nil
}

func NewMovieUseCase(movieRepository repository.MovieRepository) MovieUseCase {
	return &MovieUseCaseImpl{Repository: movieRepository}
}
