package main

import (
	"log"
	"share-basket-auth-service/core/config"
	"share-basket-auth-service/infrastructure/database"
)

func main() {
	cfg := config.Load()
	gormDB, err := database.New(cfg.DB)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	db, err := gormDB.DB()
	if err != nil {
		log.Fatalf("failed to get the database connection: %v", err)
	}

	defer db.Close()

	if err := database.Migrate(db, "migrations"); err != nil {
		log.Fatalf("failed to migrate the database: %v", err)
	}
}
