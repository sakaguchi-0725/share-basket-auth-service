package main

import (
	"context"
	"log"
	"share-basket-auth-service/core/config"
	"share-basket-auth-service/infrastructure/database"
	"share-basket-auth-service/presentation/api/middleware"
	"share-basket-auth-service/presentation/api/server"
	"share-basket-auth-service/registry"
)

func main() {
	cfg := config.Load()

	db, err := database.New(cfg.DB)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	registry, err := registry.InjectAPI(context.Background(), db, cfg.AWS)

	s := server.New(":8080")
	s.Use(middleware.ErrorMiddleware())
	s.POST("/signup", registry.SignUpHandler())
	s.POST("/signup/confirm", registry.SignUpConfirmHandler())

	if err := s.Run(); err != nil {
		log.Fatalf("Server could not be run: %v", err)
	}
}
