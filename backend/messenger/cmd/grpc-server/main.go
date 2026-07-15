package main

import (
	"log"
	"messenger/internal/config"
	"messenger/internal/kafka"
	"messenger/storage/container"
	gr "messenger/transport/grpc"
	"messenger/transport/http/gateway"
	"net"

	chatv1 "api/gen/chat/v1"

	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("load config:", err)
	}

	var producer *kafka.Producer
	if cfg.KafkaEnabled {
		producer, err = kafka.NewProducer(cfg.KafkaBrokers)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Print("Kafka disabled")
	}

	container, err := container.NewHybridContainer(cfg.PostgresDSN, cfg.RedisAddr, cfg.SocialNetworkAddr, producer)
	if err != nil {
		log.Fatal(err)
	}

	grpcLis, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	chatv1.RegisterChatServiceServer(grpcServer, gr.NewServer(container.ChatService()))
	go func() {
		log.Printf("gRPC on %s", cfg.GRPCAddr)
		grpcServer.Serve(grpcLis)
	}()

	httpServer := gateway.NewHTTPServer(cfg.HTTPAddr, container.ChatService(), cfg.JWTSecret)
	log.Printf("HTTP gateway on %s", cfg.HTTPAddr)
	log.Fatal(httpServer.ListenAndServe())
}
