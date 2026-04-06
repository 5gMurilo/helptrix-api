package uploaderinterfaces

import "github.com/gin-gonic/gin"

type IUploaderController interface {
	Upload(c *gin.Context)
}
