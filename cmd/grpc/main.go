package main

import (
	"context"
	"log"
	"share-basket-auth-service/core/config"
	"share-basket-auth-service/infrastructure/database"
	"share-basket-auth-service/presentation/grpc/server"
	"share-basket-auth-service/registry"
)

func main() {
	cfg := config.Load()

	db, err := database.New(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	registry, err := registry.InjectGRPC(context.Background(), db, cfg.AWS)
	if err != nil {
		log.Fatalf("failed to create registry: %v", err)
	}

	s := server.New(":8081")
	s.MapServices(registry.VerifyTokenHandler())
	s.Run()
}
