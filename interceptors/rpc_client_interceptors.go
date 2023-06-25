package interceptors

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ClientLogUnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context,
		method string,
		req,
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		log.Printf("%s was invoked with %v\n", method, req)

		headers, ok := metadata.FromOutgoingContext(ctx)

		if ok {
			log.Printf("Sending headers: %v\n", headers)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func ClientLogStreamInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		log.Printf("%s was invoked with %v\n", method, desc)

		headers, ok := metadata.FromOutgoingContext(ctx)

		if ok {
			log.Printf("Sending headers: %v\n", headers)
		}

		return streamer(ctx, desc, cc, method, opts...)
	}
}