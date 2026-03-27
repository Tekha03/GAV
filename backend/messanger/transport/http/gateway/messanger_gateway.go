// gateway/messenger_gateway.go
package gateway

import (
	pb "api/chat_gen/chat"
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	grpc "google.golang.org/grpc"
)

func RegisterHandlers(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
  return pb.RegisterChatServiceHandler(ctx, mux, conn)
}

func NewHTTPServer(grpcAddr string) *http.Server {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{}))
	conn, err := grpc.DialContext(ctx, grpcAddr, grpc.WithInsecure()) // prod: TLS
	if err != nil {
		panic(fmt.Sprintf("dial grpc %s: %v", grpcAddr, err))
	}

	if err := RegisterHandlers(ctx, mux, conn); err != nil {
		panic(fmt.Sprintf("register handlers: %v", err))
	}

	r := chi.NewRouter()
	r.Use(corsMiddleware()) 
	r.Mount("/v1", mux)

	return &http.Server{Addr: ":8080", Handler: r}
}

func corsMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
	  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")  // мобильные domains
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if r.Method == "OPTIONS" { w.WriteHeader(http.StatusOK); return }
		next.ServeHTTP(w, r)
	  })
	}
  }