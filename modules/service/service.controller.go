package service

import (
	"errors"
	"net/http"

	"github.com/5gMurilo/helptrix-api/adapter/auth"
	"github.com/5gMurilo/helptrix-api/core/domain"
	serviceinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/service"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ServiceController struct {
	svc serviceinterfaces.IServiceService
}

func NewServiceController(svc serviceinterfaces.IServiceService) serviceinterfaces.IServiceController {
	return &ServiceController{svc: svc}
}

// PostService godoc
//
//	@Summary		Create a service
//	@Description	Creates a new service offering for the authenticated helper user
//	@Tags			service
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		domain.CreateServiceRequestDTO	true	"Create service request"
//	@Success		201		{object}	domain.ServiceResponseDTO
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		403		{object}	map[string]string
//	@Failure		409		{object}	map[string]string
//	@Failure		422		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/service [post]
func (ctrl *ServiceController) Create(c *gin.Context) {
	payload := c.MustGet("authorization_payload").(*auth.Payload)

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var dto domain.CreateServiceRequestDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := ctrl.svc.Create(userID, payload.UserType, dto)
	if err != nil {
		if errors.Is(err, utils.ErrHelperOnly) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrCategoryNotAssignedToUser) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrServiceNameNotUnique) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrInvalidValueFormat) ||
			errors.Is(err, utils.ErrValueNotPositive) ||
			errors.Is(err, utils.ErrInvalidStartTimeFormat) ||
			errors.Is(err, utils.ErrInvalidEndTimeFormat) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetServices godoc
//
//	@Summary		List services
//	@Description	Lists all services for the authenticated helper user
//	@Tags			service
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}		domain.ServiceResponseDTO
//	@Failure		401	{object}	map[string]string
//	@Failure		403	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/service [get]
func (ctrl *ServiceController) List(c *gin.Context) {
	payload := c.MustGet("authorization_payload").(*auth.Payload)

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	response, err := ctrl.svc.List(userID, payload.UserType)
	if err != nil {
		if errors.Is(err, utils.ErrHelperOnly) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetServiceByID godoc
//
//	@Summary		Get service by ID
//	@Description	Returns a specific service by its ID for the authenticated helper user
//	@Tags			service
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Service ID"
//	@Success		200	{object}	domain.ServiceResponseDTO
//	@Failure		400	{object}	map[string]string
//	@Failure		401	{object}	map[string]string
//	@Failure		403	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/service/{id} [get]
func (ctrl *ServiceController) GetByID(c *gin.Context) {
	payload := c.MustGet("authorization_payload").(*auth.Payload)

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	serviceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service id"})
		return
	}

	response, err := ctrl.svc.GetByID(serviceID, userID, payload.UserType)
	if err != nil {
		if errors.Is(err, utils.ErrHelperOnly) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrServiceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// PutService godoc
//
//	@Summary		Update a service
//	@Description	Updates fields of an existing service for the authenticated helper user
//	@Tags			service
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string							true	"Service ID"
//	@Param			body	body		domain.UpdateServiceRequestDTO	true	"Update service request"
//	@Success		200		{object}	domain.ServiceResponseDTO
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		403		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		409		{object}	map[string]string
//	@Failure		422		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/service/{id} [put]
func (ctrl *ServiceController) Update(c *gin.Context) {
	payload := c.MustGet("authorization_payload").(*auth.Payload)

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	serviceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service id"})
		return
	}

	var dto domain.UpdateServiceRequestDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := ctrl.svc.Update(serviceID, userID, payload.UserType, dto)
	if err != nil {
		if errors.Is(err, utils.ErrHelperOnly) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrServiceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrServiceNameNotUnique) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrCategoryNotAssignedToUser) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrInvalidValueFormat) ||
			errors.Is(err, utils.ErrValueNotPositive) ||
			errors.Is(err, utils.ErrInvalidStartTimeFormat) ||
			errors.Is(err, utils.ErrInvalidEndTimeFormat) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteService godoc
//
//	@Summary		Delete a service
//	@Description	Soft-deletes a service for the authenticated helper user
//	@Tags			service
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path	string	true	"Service ID"
//	@Success		204
//	@Failure		400	{object}	map[string]string
//	@Failure		401	{object}	map[string]string
//	@Failure		403	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/service/{id} [delete]
func (ctrl *ServiceController) Delete(c *gin.Context) {
	payload := c.MustGet("authorization_payload").(*auth.Payload)

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	serviceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service id"})
		return
	}

	if err := ctrl.svc.Delete(serviceID, userID, payload.UserType); err != nil {
		if errors.Is(err, utils.ErrHelperOnly) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrServiceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}
