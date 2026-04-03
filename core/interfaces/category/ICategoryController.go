package categoryinterfaces

import "github.com/gin-gonic/gin"

type ICategoryController interface {
	List(c *gin.Context)
}
