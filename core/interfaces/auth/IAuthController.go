package authinterfaces

import "github.com/gin-gonic/gin"

type IAuthController interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}
