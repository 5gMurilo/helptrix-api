//	@title			Helptrix API
//	@version		1.0
//	@description	REST API for the Helptrix helper marketplace platform.

//	@host		localhost:8080
//	@BasePath	/

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization

package main

import (
	"log"
	"os"

	_ "github.com/5gMurilo/helptrix-api/docs"
	"github.com/joho/godotenv"

	"github.com/5gMurilo/helptrix-api/adapter/auth"
	"github.com/5gMurilo/helptrix-api/adapter/db"
	"github.com/5gMurilo/helptrix-api/adapter/db/repository"
	adapterhttp "github.com/5gMurilo/helptrix-api/adapter/http"
	"github.com/5gMurilo/helptrix-api/core/domain"
	authmodule "github.com/5gMurilo/helptrix-api/modules/auth"
	categorymodule "github.com/5gMurilo/helptrix-api/modules/category"
	proposalmodule "github.com/5gMurilo/helptrix-api/modules/proposal"
	servicemodule "github.com/5gMurilo/helptrix-api/modules/service"
	usermodule "github.com/5gMurilo/helptrix-api/modules/user"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found or could not be loaded, relying on environment variables")
	}

	gormDB, err := db.Connect()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close(gormDB)

	if err := gormDB.AutoMigrate(
		&domain.Category{},
		&domain.User{},
		&domain.Address{},
		&domain.UserCategory{},
		&domain.Service{},
		&domain.Proposal{},
	); err != nil {
		log.Fatalf("failed to run database migrations: %v", err)
	}

	maker, err := auth.NewPasetoMaker(os.Getenv("PASETO_SYMMETRIC_KEY"))
	if err != nil {
		log.Fatalf("failed to create paseto maker: %v", err)
	}

	authRepo := repository.NewAuthRepository(gormDB)
	authSvc := authmodule.NewAuthService(authRepo, maker)
	authCtrl := authmodule.NewAuthController(authSvc)

	userRepo := repository.NewUserRepository(gormDB)
	userSvc := usermodule.NewUserService(userRepo)
	userCtrl := usermodule.NewUserController(userSvc)

	categoryRepo := repository.NewCategoryRepository(gormDB)
	categorySvc := categorymodule.NewCategoryService(categoryRepo)
	categoryCtrl := categorymodule.NewCategoryController(categorySvc)

	svcRepo := repository.NewServiceRepository(gormDB)
	svcSvc := servicemodule.NewServiceService(svcRepo)
	svcCtrl := servicemodule.NewServiceController(svcSvc)

	proposalRepo := repository.NewProposalRepository(gormDB)
	proposalSvc := proposalmodule.NewProposalService(proposalRepo)
	proposalCtrl := proposalmodule.NewProposalController(proposalSvc)

	router := adapterhttp.NewRouter(maker, authCtrl, userCtrl, categoryCtrl, svcCtrl, proposalCtrl)

	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}

	log.Printf("starting server on port %s", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
