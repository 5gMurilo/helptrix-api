package service

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/5gMurilo/helptrix-api/adapter/auth"
	"github.com/5gMurilo/helptrix-api/core/domain"
	serviceinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/service"
	"github.com/5gMurilo/helptrix-api/core/logger"
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
	log := logger.Get().With(
		slog.String("layer", "controller"),
		slog.String("module", "service"),
		slog.String("operation", "Create"),
		slog.String("user_id", payload.UserID),
	)

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		log.Error("failed to parse user ID from token", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var dto domain.CreateServiceRequestDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		log.Warn("invalid request body", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := ctrl.svc.Create(userID, payload.UserType, dto)
	if err != nil {
		if errors.Is(err, utils.ErrHelperOnly) {
			log.Warn("service creation forbidden: only helpers allowed", slog.String("user_type", payload.UserType))
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrCategoryNotAssignedToUser) {
			log.Warn("service creation rejected: category not assigned to user", slog.String("error", err.Error()))
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrServiceNameNotUnique) {
			log.Warn("service creation conflict: name already in use", slog.String("error", err.Error()))
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrInvalidValueFormat) ||
			errors.Is(err, utils.ErrValueNotPositive) ||
			errors.Is(err, utils.ErrInvalidStartTimeFormat) ||
			errors.Is(err, utils.ErrInvalidEndTimeFormat) {
			log.Warn("service creation rejected: invalid field value", slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Error("service creation failed with unexpected error", slog.String("error", err.Error()))
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
	log := logger.Get().With(
		slog.String("layer", "controller"),
		slog.String("module", "service"),
		slog.String("operation", "List"),
		slog.String("user_id", payload.UserID),
	)

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		log.Error("failed to parse user ID from token", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	response, err := ctrl.svc.List(userID, payload.UserType)
	if err != nil {
		if errors.Is(err, utils.ErrHelperOnly) {
			log.Warn("service list forbidden: only helpers allowed", slog.String("user_type", payload.UserType))
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		log.Error("service list failed with unexpected error", slog.String("error", err.Error()))
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
	log := logger.Get().With(
		slog.String("layer", "controller"),
		slog.String("module", "service"),
		slog.String("operation", "GetByID"),
		slog.String("user_id", payload.UserID),
	)

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		log.Error("failed to parse user ID from token", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	serviceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Warn("invalid service ID", slog.String("id_param", c.Param("id")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service id"})
		return
	}

	response, err := ctrl.svc.GetByID(serviceID, userID, payload.UserType)
	if err != nil {
		if errors.Is(err, utils.ErrHelperOnly) {
			log.Warn("service get forbidden: only helpers allowed", slog.String("user_type", payload.UserType))
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrServiceNotFound) {
			log.Warn("service not found", slog.String("service_id", serviceID.String()))
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Error("service get failed with unexpected error", slog.String("error", err.Error()), slog.String("service_id", serviceID.String()))
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
	log := logger.Get().With(
		slog.String("layer", "controller"),
		slog.String("module", "service"),
		slog.String("operation", "Update"),
		slog.String("user_id", payload.UserID),
	)

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		log.Error("failed to parse user ID from token", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	serviceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Warn("invalid service ID", slog.String("id_param", c.Param("id")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service id"})
		return
	}

	var dto domain.UpdateServiceRequestDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		log.Warn("invalid request body", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := ctrl.svc.Update(serviceID, userID, payload.UserType, dto)
	if err != nil {
		if errors.Is(err, utils.ErrHelperOnly) {
			log.Warn("service update forbidden: only helpers allowed", slog.String("user_type", payload.UserType))
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrServiceNotFound) {
			log.Warn("service update failed: service not found", slog.String("service_id", serviceID.String()))
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrServiceNameNotUnique) {
			log.Warn("service update conflict: name already in use", slog.String("error", err.Error()))
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrCategoryNotAssignedToUser) {
			log.Warn("service update rejected: category not assigned to user", slog.String("error", err.Error()))
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrInvalidValueFormat) ||
			errors.Is(err, utils.ErrValueNotPositive) ||
			errors.Is(err, utils.ErrInvalidStartTimeFormat) ||
			errors.Is(err, utils.ErrInvalidEndTimeFormat) {
			log.Warn("service update rejected: invalid field value", slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Error("service update failed with unexpected error", slog.String("error", err.Error()), slog.String("service_id", serviceID.String()))
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
	log := logger.Get().With(
		slog.String("layer", "controller"),
		slog.String("module", "service"),
		slog.String("operation", "Delete"),
		slog.String("user_id", payload.UserID),
	)

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		log.Error("failed to parse user ID from token", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	serviceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Warn("invalid service ID", slog.String("id_param", c.Param("id")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service id"})
		return
	}

	if err := ctrl.svc.Delete(serviceID, userID, payload.UserType); err != nil {
		if errors.Is(err, utils.ErrHelperOnly) {
			log.Warn("service delete forbidden: only helpers allowed", slog.String("user_type", payload.UserType))
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrServiceNotFound) {
			log.Warn("service delete failed: service not found", slog.String("service_id", serviceID.String()))
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Error("service delete failed with unexpected error", slog.String("error", err.Error()), slog.String("service_id", serviceID.String()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}
