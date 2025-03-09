package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

func Migrate(db *sql.DB, path string) error {
	migrations := &migrate.FileMigrationSource{
		Dir: path,
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	fmt.Println("Applied", n, "migrations!")
	return nil
}
