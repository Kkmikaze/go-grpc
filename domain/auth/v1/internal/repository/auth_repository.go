package repository

import (
	"github.com/Kkmikaze/roketin/pkg/orm"
)

type AuthRepository interface {
}

type AuthRepositoryImpl struct {
	Provider *orm.Provider
}

func NewAuthRepository(conn *orm.Provider) AuthRepository {
	return &AuthRepositoryImpl{Provider: conn}
}
