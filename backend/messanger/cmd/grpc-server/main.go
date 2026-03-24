// cmd/grpc-server/main.go
package main

import (
    "log"
    "net"
    "messanger/container"
    pb "api/chat_gen" 
    "google.golang.org/grpc"
	gr "messanger/transport/grpc"
)

func main() {
    container, err := container.NewHybridContainer(
        "postgres://postgres:password@localhost:5433/messanger?sslmode=disable",
        "localhost:6379",
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
