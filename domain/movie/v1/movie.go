package fsm

import (
	"context"
	"github.com/Kkmikaze/roketin/domain/movie/v1/internal/handler"
	"github.com/Kkmikaze/roketin/domain/movie/v1/internal/repository"
	"github.com/Kkmikaze/roketin/domain/movie/v1/usecase"
	"github.com/Kkmikaze/roketin/pkg/orm"
	authv1 "github.com/Kkmikaze/roketin/stubs/auth/v1"
	moviev1 "github.com/Kkmikaze/roketin/stubs/movie/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func RegisterMovieServiceHandler(ctx context.Context, sv *runtime.ServeMux, conn *grpc.ClientConn) {
	err := moviev1.RegisterMovieServiceHandler(ctx, sv, conn)
	if err != nil {
		panic(err)
	}
}

func RegisterMovieServiceServer(db *orm.Provider, sv *grpc.Server, authClient authv1.AuthServiceClient) {
	r := repository.NewMovieRepository(db)
	u := usecase.NewMovieUseCase(r)
	srv := handler.NewMovieRestHandler(u, authClient)
	moviev1.RegisterMovieServiceServer(sv, srv)
}
