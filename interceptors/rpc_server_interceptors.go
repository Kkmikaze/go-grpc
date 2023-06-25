package interceptors

import (
	"context"
	"log"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ServerLogUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		log.Printf(info.FullMethod)
		log.Printf("Received a request: %v\n", req)

		headers, ok := metadata.FromIncomingContext(ctx)

		if ok {
			log.Printf("Received headers: %v\n", headers)
		}

		return handler(ctx, req)
	}
}

func ServerLogStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		log.Printf("Received a request: %v\n", info)

		headers, ok := metadata.FromIncomingContext(ss.Context())

		if ok {
			log.Printf("Received headers: %v\n", headers)
		}

		return handler(srv, ss)
	}
}

func RecoveryUnaryInterceptor() grpc.UnaryServerInterceptor {
	return recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(func(p any) (err error) {
		return status.Errorf(codes.Unknown, "Internal Server Error")
	}))
}