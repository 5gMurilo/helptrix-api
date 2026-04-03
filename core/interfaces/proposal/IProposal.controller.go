package proposalinterfaces

import "github.com/gin-gonic/gin"

type IProposalController interface {
	Create(c *gin.Context)
	GetByID(c *gin.Context)
	UpdateStatus(c *gin.Context)
	List(c *gin.Context)
}
