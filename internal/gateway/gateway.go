// Package gateway is described reusable package for create gateway server
package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/Kkmikaze/roketin/api"
	"github.com/Kkmikaze/roketin/constants"
	"github.com/Kkmikaze/roketin/middleware"
	"github.com/Kkmikaze/roketin/third_party"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

type Gateway struct {
	*runtime.ServeMux
	Addr           string
	MaxHeaderBytes int
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

type ErrorResponse struct {
	Status  bool              `json:"status"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

func NewGateway(addr string, opts ...runtime.ServeMuxOption) *Gateway {
	gwMux := runtime.NewServeMux(opts...)

	return &Gateway{
		ServeMux:       gwMux,
		Addr:           addr,
		MaxHeaderBytes: constants.MaxHeaderBytes,
		ReadTimeout:    constants.ReadTimeout,
		WriteTimeout:   constants.WriteTimeout,
	}
}

func (gw *Gateway) Run(ctx context.Context) error {
	sw := swaggerUIHandler()

	fileServer := http.FileServer(http.FS(api.FS))
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.Handle("/", sw)

	gwServer := &http.Server{
		Addr: fmt.Sprintf(":%v", gw.Addr),
		Handler: middleware.CORS(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if strings.HasPrefix(request.URL.Path, "/api/v1") {
				gw.ServeMux.ServeHTTP(writer, request)
				return
			}
			mux.ServeHTTP(writer, request)
		})),
		ReadTimeout:    gw.ReadTimeout,
		WriteTimeout:   gw.WriteTimeout,
		MaxHeaderBytes: gw.MaxHeaderBytes,
	}

	go func() {
		<-ctx.Done()
		log.Println("Shutting down the http gateway server")
		if err := gwServer.Shutdown(ctx); err != nil {
			log.Fatalf("Failed to shutdown http gateway server: %v", err)
		}
	}()
	if err := gwServer.ListenAndServe(); err != nil {
		log.Fatalln("Server gRPC-Gateway exited with error:", err)
	}

	return nil
}

func ExceptionHandler(ctx context.Context, mux *runtime.ServeMux, m runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	s := status.Convert(err)
	code := runtime.HTTPStatusFromCode(s.Code())
	fallback := `{"status": false, "message": "Failed to marshall message" }`
	w.Header().Set("Content-type", m.ContentType("application/json"))
	w.WriteHeader(code)

	objectMapper := make(map[string]string)

	for _, detail := range s.Details() {
		switch t := detail.(type) {
		case *errdetails.BadRequest:
			for _, violation := range t.FieldViolations {
				objectMapper[strings.ToLower(violation.GetField())] = violation.GetDescription()
			}
		}
	}

	response := ErrorResponse{
		false,
		s.Message(),
		objectMapper,
	}

	marshal, err := json.Marshal(&response)
	if err != nil {
		_, _ = w.Write([]byte(fallback))
	}

	_, _ = w.Write(marshal)

}

func swaggerUIHandler() http.Handler {
	err := mime.AddExtensionType(".svg", "image/svg+xml")
	if err != nil {
		panic(err)
	}
	subFS, err := fs.Sub(third_party.SwaggerUI, "swagger-ui")
	if err != nil {
		panic("couldn't create sub filesystem: " + err.Error())
	}
	return http.FileServer(http.FS(subFS))
}
