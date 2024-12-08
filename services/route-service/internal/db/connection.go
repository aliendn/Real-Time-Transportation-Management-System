package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // Postgres driver
)

var DB *sql.DB

func ConnectPostgres() error {
	connStr := "host=postgres-db user=postgres password=postgres dbname=routes sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to Postgres: %w", err)
	}

	return DB.Ping()
}