package database

import (
	"fmt"
	"log"
	"share-basket-auth-service/core/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(cfg config.DBConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(newDSN(cfg)), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize db: %w", err)
	}

	log.Println("Successfully connected to the database")
	return db, nil
}

func newDSN(cfg config.DBConfig) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName,
	)
}
