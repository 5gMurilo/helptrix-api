package user

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/5gMurilo/helptrix-api/adapter/auth"
	"github.com/5gMurilo/helptrix-api/core/domain"
	userinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/user"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	svc userinterfaces.IUserService
}

func NewUserController(svc userinterfaces.IUserService) userinterfaces.IUserController {
	return &UserController{svc: svc}
}

// GetProfile godoc
//
//	@Summary		Get user profile
//	@Description	Returns a user profile by ID, with optional filters for business users
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id				path		string	true	"User ID (UUID)"
//	@Param			category_id		query		int		false	"Filter services by category ID (business only)"
//	@Param			actuation_day	query		string	false	"Filter services by actuation day (business only)"
//	@Success		200				{object}	domain.GetProfileResponseDTO
//	@Failure		400				{object}	map[string]string
//	@Failure		401				{object}	map[string]string
//	@Failure		404				{object}	map[string]string
//	@Failure		500				{object}	map[string]string
//	@Router			/user/profile/{id} [get]
func (ctrl *UserController) GetProfile(c *gin.Context) {
	payload := c.MustGet("authorization_payload").(*auth.Payload)

	targetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	requesterID, err := uuid.Parse(payload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid requester id"})
		return
	}

	var filters domain.ProfileFilters

	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		parsed, err := strconv.ParseUint(categoryIDStr, 10, 64)
		if err == nil {
			categoryID := uint(parsed)
			filters.CategoryID = &categoryID
		}
	}

	if actuationDay := c.Query("actuation_day"); actuationDay != "" {
		filters.ActuationDays = []string{actuationDay}
	}

	_ = requesterID

	response, err := ctrl.svc.GetProfile(requesterID, payload.UserType, targetID, filters)
	if err != nil {
		if errors.Is(err, utils.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateProfile godoc
//
//	@Summary		Update user profile
//	@Description	Updates mutable fields of a user profile. Only the owner can update.
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path	string							true	"User ID (UUID)"
//	@Param			body	body	domain.UpdateProfileRequestDTO	true	"Update profile request"
//	@Success		204
//	@Failure		400	{object}	map[string]string
//	@Failure		401	{object}	map[string]string
//	@Failure		403	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		409	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/user/profile/{id} [put]
func (ctrl *UserController) UpdateProfile(c *gin.Context) {
	payload := c.MustGet("authorization_payload").(*auth.Payload)

	targetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	requesterID, err := uuid.Parse(payload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid requester id"})
		return
	}

	var dto domain.UpdateProfileRequestDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.svc.UpdateProfile(requesterID, targetID, dto); err != nil {
		if errors.Is(err, utils.ErrNotOwner) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrCategoryHasLinkedServices) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteProfile godoc
//
//	@Summary		Delete user profile
//	@Description	Soft-deletes a user profile. Only the owner can delete.
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path	string	true	"User ID (UUID)"
//	@Success		204
//	@Failure		400	{object}	map[string]string
//	@Failure		401	{object}	map[string]string
//	@Failure		403	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/user/profile/{id} [delete]
func (ctrl *UserController) DeleteProfile(c *gin.Context) {
	payload := c.MustGet("authorization_payload").(*auth.Payload)

	targetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	requesterID, err := uuid.Parse(payload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid requester id"})
		return
	}

	if err := ctrl.svc.DeleteProfile(requesterID, targetID); err != nil {
		if errors.Is(err, utils.ErrNotOwner) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}
