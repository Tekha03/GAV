package main

import (
    "log"
    "net"
    "messanger/chat/config"
    "messanger/chat/service"
    "messanger/chat/transport/grpc"
    pb "api"
    "google.golang.org/grpc"
)

func main() {
    cfg := config.Load()
	chatServie := service.NewService()
    
    lis, err := net.Listen("tcp", cfg.GRPC.Port)
    if err != nil {
        log.Fatalf("listen: %v", err)
    }
    
    server := grpc.NewServer(svc)
    pb.RegisterChatServiceServer(grpc.NewServer(), server)
    
    log.Printf("gRPC server on %s", cfg.GRPC.Port)
    log.Fatal(grpc.Serve(lis))
}
