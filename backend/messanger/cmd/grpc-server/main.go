// cmd/grpc-server/main.go
package main

import (
	pb "api/chat_gen"
	"log"
	"messanger/container"
	gr "messanger/transport/grpc"
	"net"
	"os"

	"google.golang.org/grpc"
)

func main() {
    postgresDSN := os.Getenv("POSTGRES_DSN")
    redisAddr := os.Getenv("REDIS_ADDR")
    socialNetworkAddr := os.Getenv("SOCIAL_NETWORK_ADDR")

    container, err := container.NewHybridContainer(
        postgresDSN,
        redisAddr,
        socialNetworkAddr,
    )
    if err != nil {
        log.Fatal(err)
    }
    
    lis, err := net.Listen("tcp", ":9090")
    if err != nil {
        log.Fatal(err)
    }
    
    server := grpc.NewServer()
    pb.RegisterChatServiceServer(server, gr.NewServer(container.ChatService()))
    
    log.Println("gRPC server on :9090")
    if err := server.Serve(lis); err != nil {
        log.Fatal(err)
    }
}
