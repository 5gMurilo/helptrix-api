package adapterhttp

import (
	"net/http"

	"github.com/5gMurilo/helptrix-api/adapter/auth"
	"github.com/5gMurilo/helptrix-api/adapter/http/middleware"
	authifaces "github.com/5gMurilo/helptrix-api/core/interfaces/auth"
	categoryinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/category"
	otpinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/otp"
	helperinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/helper"
	proposalinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/proposal"
	reviewinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/review"
	serviceinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/service"
	uploaderinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/uploader"
	userinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/user"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(
	maker *auth.PasetoMaker,
	authCtrl authifaces.IAuthController,
	userCtrl userinterfaces.IUserController,
	categoryCtrl categoryinterfaces.ICategoryController,
	svcCtrl serviceinterfaces.IServiceController,
	proposalCtrl proposalinterfaces.IProposalController,
	otpCtrl otpinterfaces.IOtpController,
	uploaderCtrl uploaderinterfaces.IUploaderController,
	helperCtrl helperinterfaces.IHelperController,
	reviewCtrl reviewinterfaces.IReviewController,
) *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/swagger/index.html")
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", authCtrl.Register)
		authGroup.POST("/login", authCtrl.Login)
	}

	otpGroup := router.Group("/otp")
	{
		otpGroup.POST("/send", otpCtrl.Send)
		otpGroup.POST("/confirm", otpCtrl.Confirm)
	}

	router.GET("/category", categoryCtrl.List)

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware(maker))
	{
		userGroup := protected.Group("/user")
		{
			userGroup.GET("/profile/:id", userCtrl.GetProfile)
			userGroup.PUT("/profile/:id", userCtrl.UpdateProfile)
			userGroup.DELETE("/profile/:id", userCtrl.DeleteProfile)
		}

		svcGroup := protected.Group("/service")
		{
			svcGroup.POST("", svcCtrl.Create)
			svcGroup.GET("", svcCtrl.List)
			svcGroup.GET("/:id", svcCtrl.GetByID)
			svcGroup.PUT("/:id", svcCtrl.Update)
			svcGroup.DELETE("/:id", svcCtrl.Delete)
		}

		proposalGroup := protected.Group("/proposal")
		{
			proposalGroup.POST("", proposalCtrl.Create)
			proposalGroup.GET("", proposalCtrl.List)
			proposalGroup.GET("/:id", proposalCtrl.GetByID)
			proposalGroup.PATCH("/:id/status", proposalCtrl.UpdateStatus)
		}

		protected.POST("/image-uploader/:image-type/:id", uploaderCtrl.Upload)

		protected.GET("/helper", helperCtrl.List)

		reviewGroup := protected.Group("/review")
		{
			reviewGroup.POST("", reviewCtrl.Create)
			reviewGroup.GET("/business", reviewCtrl.ListBusiness)
			reviewGroup.GET("/helper", reviewCtrl.ListHelper)
		}
	}

	return router
}
