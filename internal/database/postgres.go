package database

import (
	"database/sql"
	"fmt"

	"github.com/emiliosheinz/rinha-de-backend-2025-go/internal/config"
	_ "github.com/lib/pq"
)

func ConnectPostgres() (*sql.DB, error) {
	connectionString := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		config.PostgresDatabaseHostname,
		config.PostgresDatabasePort,
		config.PostgresDatabaseName,
		config.PostgresDatabaseUser,
		config.PostgresDatabasePassword,
	)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(0)

	fmt.Println("Database connection established successfully")
	return db, nil
}
