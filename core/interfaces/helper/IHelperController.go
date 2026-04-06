package helperinterfaces

import "github.com/gin-gonic/gin"

type IHelperController interface {
	List(c *gin.Context)
}
