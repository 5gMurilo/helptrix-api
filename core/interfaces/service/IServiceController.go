package serviceinterfaces

import "github.com/gin-gonic/gin"

type IServiceController interface {
	Create(c *gin.Context)
	List(c *gin.Context)
	GetByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}
