package main

import (
	"log"

	"github.com/5gMurilo/helptrix-api/adapter/db"
	"github.com/5gMurilo/helptrix-api/adapter/db/seeder"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("warning: .env not loaded, using environment variables")
	}

	gormDB, err := db.Connect()
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close(gormDB)

	if err := seeder.SeedCategories(gormDB); err != nil {
		log.Fatalf("seed: %v", err)
	}

	log.Println("categories seed completed")
}
