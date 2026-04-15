package reviewinterfaces

import "github.com/gin-gonic/gin"

type IReviewController interface {
	Create(c *gin.Context)
	ListBusiness(c *gin.Context)
	ListHelper(c *gin.Context)
}
