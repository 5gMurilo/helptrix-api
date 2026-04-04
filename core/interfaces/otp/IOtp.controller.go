package otpinterfaces

import "github.com/gin-gonic/gin"

type IOtpController interface {
	Send(c *gin.Context)
	Confirm(c *gin.Context)
}
