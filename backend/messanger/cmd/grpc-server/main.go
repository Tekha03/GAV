package main

import (
    chatv1 "github.com/Tekha03/GAV/api_gav/gen/chat/v1"
	"log"
	"messanger/internal/config"
	"messanger/storage/container"
	gr "messanger/transport/grpc"
	"messanger/transport/http/gateway"
	"net"

	"google.golang.org/grpc"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatal("load config:", err)
    }

    container, err := container.NewHybridContainer(cfg.PostgresDSN, cfg.RedisAddr, cfg.SocialNetworkAddr)
    if err != nil { log.Fatal(err) }

    grpcLis, err := net.Listen("tcp", cfg.GRPCAddr)
    if err != nil { log.Fatal(err) }
    grpcServer := grpc.NewServer()
    chatv1.RegisterChatServiceServer(grpcServer, gr.NewServer(container.ChatService()))
    go func() { log.Printf("gRPC on %s", cfg.GRPCAddr); grpcServer.Serve(grpcLis) }()

    httpServer := gateway.NewHTTPServer(cfg.GRPCAddr)
    log.Printf("HTTP gateway on :8080")
    log.Fatal(httpServer.ListenAndServe())
}
