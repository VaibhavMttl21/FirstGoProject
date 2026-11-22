package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPostgres() (*gorm.DB, error) {
    databaseURL := os.Getenv("DATABASE_URL")

    if databaseURL == "" {
        return nil, fmt.Errorf("DATABASE_URL is not set")
    }

    db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
    }

    sqlDB, err := db.DB()
    if err != nil {
        return nil, fmt.Errorf("failed to access database connection pool: %w", err)
    }

    if err := sqlDB.Ping(); err != nil {
        return nil, fmt.Errorf("PostgreSQL ping failed: %w", err)
    }

    return db, nil
}

// Step 2: Understand the important parts
// gorm.Open()

// Initializes GORM with the PostgreSQL driver and connection configuration.

// db.DB()

// Returns the underlying *sql.DB connection pool managed by Go's database/sql package.

// sqlDB.Ping()

// Verifies that PostgreSQL is reachable.

// Why return (*gorm.DB, error)?

// Because the caller should decide what happens if the database connection fails. We don't want to terminate the application from inside our database package.