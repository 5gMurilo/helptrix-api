package category

import (
	"net/http"

	categoryinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/category"
	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	svc categoryinterfaces.ICategoryService
}

func NewCategoryController(svc categoryinterfaces.ICategoryService) categoryinterfaces.ICategoryController {
	return &CategoryController{svc: svc}
}

// GetCategories godoc
//
//	@Summary		List categories
//	@Description	Returns all categories available on the platform (public catalog)
//	@Tags			category
//	@Produce		json
//	@Success		200	{array}		domain.CategoryListItemResponseDTO
//	@Failure		500	{object}	map[string]string
//	@Router			/category [get]
func (ctrl *CategoryController) List(c *gin.Context) {
	list, err := ctrl.svc.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, list)
}
