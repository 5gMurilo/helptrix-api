package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {

	if os.Getenv("DB_URL") != "" {
		return gorm.Open(postgres.Open(os.Getenv("DB_URL")), &gorm.Config{})
	}
    
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func Close(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		return
	}

	if err := sqlDB.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "error closing database connection: %v\n", err)
	}
}
