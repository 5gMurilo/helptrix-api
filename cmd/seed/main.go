package main

import (
	"log"

	"github.com/5gMurilo/helptrix-api/adapter/db"
	"github.com/5gMurilo/helptrix-api/adapter/db/repository"
	"github.com/5gMurilo/helptrix-api/adapter/db/seeder"
	"github.com/5gMurilo/helptrix-api/modules/category"
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

	categoryRepo := repository.NewCategoryRepository(gormDB)
	categorySvc := category.NewCategoryService(categoryRepo)

	if err := seeder.SeedCategories(categorySvc); err != nil {
		log.Fatalf("seed: %v", err)
	}

	log.Println("categories seed completed")
}
