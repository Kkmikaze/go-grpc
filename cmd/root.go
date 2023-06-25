// Package cmd is described Main applications for this project.
package cmd

import (
	"context"
	"fmt"
	authFms "github.com/Kkmikaze/roketin/domain/auth/v1"
	movieFms "github.com/Kkmikaze/roketin/domain/movie/v1"
	authv1 "github.com/Kkmikaze/roketin/stubs/auth/v1"
	"google.golang.org/protobuf/proto"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Kkmikaze/roketin/configs"
	"github.com/Kkmikaze/roketin/constants"
	"github.com/Kkmikaze/roketin/interceptors"
	"github.com/Kkmikaze/roketin/internal/gateway"
	"github.com/Kkmikaze/roketin/pkg/orm"
	"github.com/Kkmikaze/roketin/pkg/rpcclient"
	"github.com/Kkmikaze/roketin/pkg/rpcserver"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/encoding/protojson"
	"gorm.io/gorm"
)

var (
	authport  string
	movieport string
	gwPort    string
	rootCmd   = &cobra.Command{
		Use:   "service",
		Short: "Running the gRPC service",
		Long:  "Used to run gRPC Service including rpc server, rpc client and gateway",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logrus.New()
			ctx := context.Background()

			// Auth Service
			dbAuth, err := orm.NewMySQL(ctx,
				configs.Configs.DBDNSAuth,
				&orm.ConfigConnProvider{
					ConnMaxLifetime: constants.ConnMaxLifeTime,
					ConnMaxIdleTime: constants.ConnMaxIdleTime,
					MaxOpenConns:    constants.MaxOpenConns,
					MaxIdleConns:    constants.MaxIdleConns,
				}, &gorm.Config{})
			if err != nil {
				panic(err)
			}

			authServer := rpcserver.NewRPCServer(authport,
				"tcp",
				false,
				[]grpc.UnaryServerInterceptor{interceptors.ServerLogUnaryInterceptor(), interceptors.RecoveryUnaryInterceptor()},
				[]grpc.StreamServerInterceptor{interceptors.ServerLogStreamInterceptor()},
			)
			defer authServer.StopListener()

			authFms.RegisterAuthServiceServer(dbAuth, authServer.Server)
			grpc_health_v1.RegisterHealthServer(authServer.Server, health.NewServer())
			logger.Infoln("Serving Auth gRPC on", authport)
			if err := authServer.Run(); err != nil {
				logger.Fatalln("Failed to listen auth grpc server err", err)
			}

			authClient, err := rpcclient.NewRPCClient(ctx,
				authport,
				false,
				[]grpc.UnaryClientInterceptor{interceptors.ClientLogUnaryInterceptor()},
				[]grpc.StreamClientInterceptor{interceptors.ClientLogStreamInterceptor()},
				grpc.WithBlock(),
			)
			if err != nil {
				logger.Fatalln("failed to dial auth server:", err)
				log.Fatalln("Failed to dial auth server:", err)
			}

			// Movie Service
			dbMovie, err := orm.NewMySQL(ctx,
				configs.Configs.DBDNSMovie,
				&orm.ConfigConnProvider{
					ConnMaxLifetime: constants.ConnMaxLifeTime,
					ConnMaxIdleTime: constants.ConnMaxIdleTime,
					MaxOpenConns:    constants.MaxOpenConns,
					MaxIdleConns:    constants.MaxIdleConns,
				}, &gorm.Config{})
			if err != nil {
				panic(err)
			}

			movieServer := rpcserver.NewRPCServer(movieport,
				"tcp",
				false,
				[]grpc.UnaryServerInterceptor{interceptors.ServerLogUnaryInterceptor(), interceptors.RecoveryUnaryInterceptor()},
				[]grpc.StreamServerInterceptor{interceptors.ServerLogStreamInterceptor()},
			)
			defer movieServer.StopListener()

			authClientConn := authv1.NewAuthServiceClient(authClient)
			movieFms.RegisterMovieServiceServer(dbMovie, movieServer.Server, authClientConn)
			grpc_health_v1.RegisterHealthServer(movieServer.Server, health.NewServer())
			logger.Infoln("Serving Movie gRPC on", movieport)
			if err := movieServer.Run(); err != nil {
				logger.Fatalln("Failed to listen movie grpc server")
			}

			movieClient, err := rpcclient.NewRPCClient(ctx,
				movieport,
				false,
				[]grpc.UnaryClientInterceptor{interceptors.ClientLogUnaryInterceptor()},
				[]grpc.StreamClientInterceptor{interceptors.ClientLogStreamInterceptor()},
				grpc.WithBlock(),
			)
			if err != nil {
				logger.Fatalln("failed to dial movie server:", err)
				log.Fatalln("Failed to dial movie server:", err)
			}

			rpcGateway := gateway.NewGateway(gwPort,
				runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
					MarshalOptions: protojson.MarshalOptions{
						UseProtoNames:   true,
						EmitUnpopulated: true,
					},
					UnmarshalOptions: protojson.UnmarshalOptions{
						DiscardUnknown: true,
					},
				}),
				runtime.WithErrorHandler(gateway.ExceptionHandler),
				runtime.WithForwardResponseOption(func(ctx context.Context, w http.ResponseWriter, message proto.Message) error {
					md, ok := runtime.ServerMetadataFromContext(ctx)
					if !ok {
						return nil
					}

					// set http status code
					if vals := md.HeaderMD.Get("x-http-code"); len(vals) > 0 {
						code, err := strconv.Atoi(vals[0])
						if err != nil {
							return err
						}
						w.WriteHeader(code)
						delete(md.HeaderMD, "x-http-code")
						delete(w.Header(), "Grpc-Metadata-X-Http-Code")
					}

					return nil
				}),
			)
			authFms.RegisterAuthServiceHandler(ctx, rpcGateway.ServeMux, authClient)
			movieFms.RegisterMovieServiceHandler(ctx, rpcGateway.ServeMux, movieClient)
			logger.Infoln("Serving gRPC-Gateway on", gwPort)
			if err := rpcGateway.Run(ctx); err != nil {
				logger.Fatalln("Failed to listen grpc server")
			}

			authServer.Terminate(ctx)
			movieServer.Terminate(ctx)
		},
	}
)

func Execute() {
	rootCmd.Flags().StringVarP(&authport, "authport", "a", "", "define auth rpc server port")
	rootCmd.Flags().StringVarP(&movieport, "movieport", "i", "", "define movie rpc server port")
	rootCmd.Flags().StringVarP(&gwPort, "gwport", "g", "", "define gateway port")
	rootCmd.MarkFlagsRequiredTogether("authport", "movieport", "gwport")

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
