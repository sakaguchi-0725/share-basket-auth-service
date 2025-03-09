package persistence_test

import (
	"fmt"
	"log"
	"os"
	"share-basket-auth-service/core/config"
	"share-basket-auth-service/infrastructure/database"
	"sync"
	"testing"

	migrate "github.com/rubenv/sql-migrate"
	"gorm.io/gorm"
)

const migrationsPath = "../../migrations"

var (
	testDB *gorm.DB
	once   sync.Once
	cfg    = config.DBConfig{
		Host:     os.Getenv("TEST_DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("TEST_DB_USER"),
		Password: os.Getenv("TEST_DB_PASSWORD"),
		DBName:   os.Getenv("TEST_DB_NAME"),
	}
)

func TestMain(m *testing.M) {
	setupTestDB()
	code := m.Run()
	teardownTestDB()
	os.Exit(code)
}

func setupTestDB() {
	once.Do(func() {
		var err error

		testDB, err = database.New(cfg)
		if err != nil {
			log.Fatalf("initialize test database: %v", err)
		}

		db, err := testDB.DB()
		if err != nil {
			log.Fatalf("failed to get test database: %v", err)
		}

		if err := database.Migrate(db, migrationsPath); err != nil {
			log.Fatalf("migrate test database: %v", err)
		}

	})
}

func clearTestData() {
	tables := []string{"users"}

	for _, table := range tables {
		err := testDB.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE;", table)).Error
		if err != nil {
			log.Fatalf("failed to clear test data from %s: %v", table, err)
		}
	}

	fmt.Println("Test data cleared successfully")
}

func teardownTestDB() {
	db, err := testDB.DB()
	if err != nil {
		log.Fatalf("get test database: %v", err)
	}

	defer db.Close()

	migrations := &migrate.FileMigrationSource{
		Dir: migrationsPath,
	}

	n, err := migrate.ExecMax(db, "postgres", migrations, migrate.Down, -1)
	if err != nil {
		log.Fatalf("Failed to apply down migration: %v", err)
	}

	fmt.Printf("Rolled back %d migration(s)\n", n)
}
