package category

import (
	"log/slog"
	"net/http"

	categoryinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/category"
	"github.com/5gMurilo/helptrix-api/core/logger"
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
	log := logger.Get().With(
		slog.String("layer", "controller"),
		slog.String("module", "category"),
		slog.String("operation", "List"),
	)

	list, err := ctrl.svc.List()
	if err != nil {
		log.Error("failed to list categories", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, list)
}
