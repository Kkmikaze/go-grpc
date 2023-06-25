package usecase

import (
	"github.com/Kkmikaze/roketin/domain/auth/v1/internal/repository"
)

type AuthUseCase interface {
}

type AuthUseCaseImpl struct {
	Repository repository.AuthRepository
}

func NewAuthUseCase(authRepository repository.AuthRepository) AuthUseCase {
	return &AuthUseCaseImpl{Repository: authRepository}
}
