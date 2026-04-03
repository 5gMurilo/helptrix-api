package userinterfaces

import "github.com/gin-gonic/gin"

type IUserController interface {
	GetProfile(c *gin.Context)
	UpdateProfile(c *gin.Context)
	DeleteProfile(c *gin.Context)
}
