package repository

import (
	"github.com/Kkmikaze/roketin/pkg/orm"
)

type MovieRepository interface {
}

type MovieRepositoryImpl struct {
	Provider *orm.Provider
}

func NewMovieRepository(conn *orm.Provider) MovieRepository {
	return &MovieRepositoryImpl{Provider: conn}
}
