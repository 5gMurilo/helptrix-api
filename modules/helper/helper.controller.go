package helper

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/5gMurilo/helptrix-api/adapter/auth"
	"github.com/5gMurilo/helptrix-api/core/domain"
	helperinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/helper"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"github.com/gin-gonic/gin"
)

type HelperController struct {
	svc helperinterfaces.IHelperService
}

func NewHelperController(svc helperinterfaces.IHelperService) helperinterfaces.IHelperController {
	return &HelperController{svc: svc}
}

// List godoc
//
//	@Summary		List and search helpers
//	@Description	Returns a paginated list of helper-type users. Supports optional filtering by name (substring) and category ID. Restricted to business users.
//	@Tags			helper
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			name			query		string	false	"Substring match on helper name (case-insensitive)"
//	@Param			category_id		query		int		false	"Filter helpers by assigned category ID"
//	@Param			page			query		int		false	"Page number (default: 1)"
//	@Param			page_size		query		int		false	"Results per page (default: 20)"
//	@Success		200	{object}	domain.HelperListResponseDTO
//	@Failure		401	{object}	map[string]string
//	@Failure		403	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/helper [get]
func (ctrl *HelperController) List(c *gin.Context) {
	payload := c.MustGet("authorization_payload").(*auth.Payload)

	params := domain.HelperSearchParams{
		Name:     c.Query("name"),
		Page:     1,
		PageSize: 20,
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if v, err := strconv.Atoi(pageStr); err == nil && v > 0 {
			params.Page = v
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if v, err := strconv.Atoi(pageSizeStr); err == nil && v > 0 {
			params.PageSize = v
		}
	}

	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if v, err := strconv.ParseUint(categoryIDStr, 10, 64); err == nil && v > 0 {
			categoryID := uint(v)
			params.CategoryID = &categoryID
		}
	}

	result, err := ctrl.svc.Search(payload.UserType, params)
	if err != nil {
		if errors.Is(err, utils.ErrBusinessOnly) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, result)
}
