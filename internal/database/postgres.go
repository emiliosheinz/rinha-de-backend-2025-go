package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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
	
	var db *sql.DB
	var err error
	maxRetries := 10
	retryDelay := 2 * time.Second
	
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", connectionString)
		if err != nil {
			log.Printf("Attempt %d/%d: Failed to open database connection: %v", i+1, maxRetries, err)
			time.Sleep(retryDelay)
			continue
		}
		
		err = db.Ping()
		if err != nil {
			log.Printf("Attempt %d/%d: Failed to ping database: %v", i+1, maxRetries, err)
			db.Close()
			time.Sleep(retryDelay)
			continue
		}
		
		break
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %v", maxRetries, err)
	}

	db.SetMaxOpenConns(64)
	db.SetMaxIdleConns(32)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(30 * time.Second)

	fmt.Println("Database connection established successfully")
	return db, nil
}
