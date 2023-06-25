package repository

import (
	"context"
	"github.com/Kkmikaze/roketin/domain/movie/v1/internal/entity"
	"github.com/Kkmikaze/roketin/pkg/orm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MovieRepository interface {
	CreateMovie(ctx context.Context, payload *entity.Movie) error
	GetMovie(ctx context.Context, query *orm.QueryBuilder) ([]*entity.Movie, int64, error)
	GetMovieByID(ctx context.Context, id uint64) (*entity.Movie, error)
	UpdateMovie(ctx context.Context, e *entity.Movie) (*entity.Movie, error)
	DeleteMovie(ctx context.Context, id uint64) error
}

type MovieRepositoryImpl struct {
	Provider *orm.Provider
}

func (r *MovieRepositoryImpl) CreateMovie(ctx context.Context, payload *entity.Movie) error {
	if err := r.Provider.WithContext(ctx).Create(payload).Error; err != nil {
		return status.Error(codes.Internal, "Internal Server Error")
	}

	return nil
}

func (r *MovieRepositoryImpl) GetMovie(ctx context.Context, query *orm.QueryBuilder) ([]*entity.Movie, int64, error) {
	var count int64
	var rows []*entity.Movie

	statement := `*`
	result := r.Provider.WithContext(ctx).Model(&entity.Movie{}).Select(statement)

	if query.Search != "" {
		keyword := "%" + query.Search + "%"
		result = result.Where("title LIKE ? OR description LIKE ? OR artist LIKE ? OR genre LIKE ?", keyword, keyword, keyword, keyword)
	}

	if err := result.Count(&count).Error; err != nil {
		return nil, 0, status.Error(codes.Internal, "Internal Server Error")
	}

	if query.Page > 0 {
		result = result.Limit(query.ItemPerPage).Offset((query.Page - 1) * query.ItemPerPage)
	}

	result = result.Order("created_at desc")

	if err := result.Find(&rows).Error; err != nil {
		return nil, 0, status.Error(codes.Internal, "Internal Server Error")
	}

	return rows, count, nil
}

func (r *MovieRepositoryImpl) GetMovieByID(ctx context.Context, id uint64) (*entity.Movie, error) {
	var row entity.Movie
	if err := r.Provider.WithContext(ctx).Where("id = ?", id).First(&row).Error; err != nil {
		if row.Id == 0 {
			return nil, status.Error(codes.NotFound, "Movie not found")
		}
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	return &row, nil
}

func (r *MovieRepositoryImpl) UpdateMovie(ctx context.Context, e *entity.Movie) (*entity.Movie, error) {
	if err := r.Provider.WithContext(ctx).Updates(e).Error; err != nil {
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	return e, nil
}

func (r *MovieRepositoryImpl) DeleteMovie(ctx context.Context, id uint64) error {
	if err := r.Provider.WithContext(ctx).Delete(&entity.Movie{Id: id}).Error; err != nil {
		return status.Error(codes.Internal, "Internal Server Error")
	}

	return nil
}

func NewMovieRepository(conn *orm.Provider) MovieRepository {
	return &MovieRepositoryImpl{Provider: conn}
}
