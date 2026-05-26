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
	"context"
	"log/slog"
	"os"

	_ "github.com/5gMurilo/helptrix-api/docs"
	"github.com/joho/godotenv"

	"github.com/5gMurilo/helptrix-api/adapter/auth"
	"github.com/5gMurilo/helptrix-api/adapter/db"
	"github.com/5gMurilo/helptrix-api/adapter/db/repository"
	"github.com/5gMurilo/helptrix-api/adapter/email"
	adapterhttp "github.com/5gMurilo/helptrix-api/adapter/http"
	adapterstorage "github.com/5gMurilo/helptrix-api/adapter/storage"
	"github.com/5gMurilo/helptrix-api/core/domain"
	uploaderinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/uploader"
	"github.com/5gMurilo/helptrix-api/core/logger"
	authmodule "github.com/5gMurilo/helptrix-api/modules/auth"
	categorymodule "github.com/5gMurilo/helptrix-api/modules/category"
	otpmodule "github.com/5gMurilo/helptrix-api/modules/otp"
	proposalmodule "github.com/5gMurilo/helptrix-api/modules/proposal"
	servicemodule "github.com/5gMurilo/helptrix-api/modules/service"
	uploadermodule "github.com/5gMurilo/helptrix-api/modules/uploader"
	uploaderstrategies "github.com/5gMurilo/helptrix-api/modules/uploader/strategies"
	helpermodule "github.com/5gMurilo/helptrix-api/modules/helper"
	reviewmodule "github.com/5gMurilo/helptrix-api/modules/review"
	usermodule "github.com/5gMurilo/helptrix-api/modules/user"
)

func main() {
	logger.Init()
	log := logger.Get()

	if err := godotenv.Load(); err != nil {
		log.Warn("env file not found, relying on environment variables")
	}

	gormDB, err := db.Connect()
	if err != nil {
		log.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close(gormDB)

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("failed to get underlying sql.DB: %v", err)
	}
	if _, err := sqlDB.Exec("CREATE EXTENSION IF NOT EXISTS pg_trgm"); err != nil {
		log.Printf("warning: could not create pg_trgm extension: %v", err)
	}

	if err := gormDB.AutoMigrate(
		&domain.Category{},
		&domain.User{},
		&domain.Address{},
		&domain.UserCategory{},
		&domain.Service{},
		&domain.Proposal{},
		&domain.OTP{},
		&domain.Review{},
	); err != nil {
		log.Error("failed to run database migrations", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if _, err := sqlDB.Exec("CREATE INDEX IF NOT EXISTS idx_users_name_trgm ON users USING GIN (name gin_trgm_ops)"); err != nil {
		log.Printf("warning: could not create trigram index on users.name: %v", err)
	}

	maker, err := auth.NewPasetoMaker(os.Getenv("PASETO_SYMMETRIC_KEY"))
	if err != nil {
		log.Error("failed to create paseto maker", slog.String("error", err.Error()))
		os.Exit(1)
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

	emailSender := email.NewResendEmailSender()
	otpRepo := repository.NewOtpRepository(gormDB)
	otpSvc := otpmodule.NewOtpService(otpRepo, emailSender)
	otpCtrl := otpmodule.NewOtpController(otpSvc)

	storageClient, err := adapterstorage.NewFirebaseStorageClient(context.Background())
	if err != nil {
		log.Error("failed to create firebase storage client", slog.String("error", err.Error()))
		os.Exit(1)
	}

	bucketName := os.Getenv("FIREBASE_STORAGE_BUCKET")

	strategies := map[string]uploaderinterfaces.IImageUploadStrategy{
		"profile-images": uploaderstrategies.NewProfileImageStrategy(storageClient, userRepo, bucketName),
		"service-images": uploaderstrategies.NewServiceImageStrategy(storageClient, svcRepo, bucketName),
	}

	uploaderSvc := uploadermodule.NewUploaderService(strategies)
	uploaderCtrl := uploadermodule.NewUploaderController(uploaderSvc)

	helperRepo := repository.NewHelperRepository(gormDB)
	helperSvc := helpermodule.NewHelperService(helperRepo)
	helperCtrl := helpermodule.NewHelperController(helperSvc)

	reviewRepo := repository.NewReviewRepository(gormDB)
	reviewSvc := reviewmodule.NewReviewService(reviewRepo)
	reviewCtrl := reviewmodule.NewReviewController(reviewSvc)

	router := adapterhttp.NewRouter(maker, authCtrl, userCtrl, categoryCtrl, svcCtrl, proposalCtrl, otpCtrl, uploaderCtrl, helperCtrl, reviewCtrl)

	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}

	log.Info("starting server", slog.String("port", port))

	if err := router.Run(":" + port); err != nil {
		log.Error("failed to start server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
