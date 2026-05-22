package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Fail to open DB:", err)
	}

	log.Printf("Connect to postgresSQL")
	createTables()
}

func createTables() {
	accountTable := `
	CREATE TABLE IF NOT EXISTED accounts (
		id					SERIAL PRIMARY KEY,
		name				TEXT NOT NULL,
		type				TEXT NOT NULL CHECK (type IN ('cash', 'bank', 'credit')),
		created_at	TIMESTAMP DEFAULT NOW()
);`

	transactionTable := `
  CREATE TABLE IF NOT EXISTS transactions (
		id          SERIAL PRIMARY KEY,
    account_id  INT REFERENCES accounts(id) ON DELETE CASCADE,
		amount      NUMERIC(12,2) NOT NULL,
    type        TEXT NOT NULL CHECK (type IN ('income', 'expense')),
		description TEXT,
    created_at  TIMESTAMP DEFAULT NOW()
);`

	if _, err := DB.Exec(accountTable); err != nil {
		log.Fatal("Fail to create accounts table", err)
	}

	if _, err := DB.Exec(transactionTable); err != nil {
		log.Fatal("Fail to create transactions table", err)
	}

	log.Println("Table ready")
}
