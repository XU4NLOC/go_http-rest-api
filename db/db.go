package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect(ctx context.Context) error {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		getEnv("DB_SSLMODE", "disable"),
	)

	connection, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("open postgres connection: %w", err)
	}

	connection.SetMaxOpenConns(getEnvInt("DB_MAX_OPEN_CONNS", 20))
	connection.SetMaxIdleConns(getEnvInt("DB_MAX_IDLE_CONNS", 5))
	connection.SetConnMaxLifetime(30 * time.Minute)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := connection.PingContext(pingCtx); err != nil {
		_ = connection.Close()
		return fmt.Errorf("ping postgres: %w", err)
	}

	DB = connection

	if err := createTables(); err != nil {
		_ = DB.Close()
		DB = nil
		return err
	}

	return nil
}

func Ping(ctx context.Context) error {
	if DB == nil {
		return fmt.Errorf("database is not initialized")
	}
	return DB.PingContext(ctx)
}

func Close() error {
	if DB == nil {
		return nil
	}
	return DB.Close()
}

func createTables() error {
	accountTable := `
	CREATE TABLE IF NOT EXISTS accounts (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		type TEXT NOT NULL CHECK (type IN ('cash', 'bank', 'credit')),
		created_at TIMESTAMP DEFAULT NOW()
	);`

	transactionTable := `
	CREATE TABLE IF NOT EXISTS transactions (
		id SERIAL PRIMARY KEY,
		account_id INT REFERENCES accounts(id) ON DELETE CASCADE,
		amount NUMERIC(12,2) NOT NULL,
		type TEXT NOT NULL CHECK (type IN ('income', 'expense')),
		description TEXT,
		created_at TIMESTAMP DEFAULT NOW()
	);`

	if _, err := DB.Exec(accountTable); err != nil {
		return fmt.Errorf("create accounts table: %w", err)
	}
	if _, err := DB.Exec(transactionTable); err != nil {
		return fmt.Errorf("create transactions table: %w", err)
	}

	return nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
