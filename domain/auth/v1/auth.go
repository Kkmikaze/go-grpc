package fsm

import (
	"context"
	"github.com/Kkmikaze/roketin/domain/auth/v1/internal/handler"
	"github.com/Kkmikaze/roketin/domain/auth/v1/internal/repository"
	"github.com/Kkmikaze/roketin/domain/auth/v1/usecase"
	"github.com/Kkmikaze/roketin/pkg/orm"
	authv1 "github.com/Kkmikaze/roketin/stubs/auth/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func RegisterAuthServiceHandler(ctx context.Context, sv *runtime.ServeMux, conn *grpc.ClientConn) {
	err := authv1.RegisterAuthServiceHandler(ctx, sv, conn)
	if err != nil {
		panic(err)
	}
}

func RegisterAuthServiceServer(db *orm.Provider, sv *grpc.Server) {
	r := repository.NewAuthRepository(db)
	u := usecase.NewAuthUseCase(r)
	srv := handler.NewAuthRestHandler(u)
	authv1.RegisterAuthServiceServer(sv, srv)
}
