package proposal

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/5gMurilo/helptrix-api/adapter/auth"
	"github.com/5gMurilo/helptrix-api/core/domain"
	proposalinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/proposal"
	"github.com/5gMurilo/helptrix-api/core/logger"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProposalController struct {
	service proposalinterfaces.IProposalService
}

func NewProposalController(service proposalinterfaces.IProposalService) proposalinterfaces.IProposalController {
	return &ProposalController{service: service}
}

// PostProposal godoc
//
//	@Summary		Create a proposal
//	@Description	Creates a new service proposal from a business user directed at a helper. Only business users can create proposals.
//	@Tags			proposal
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		domain.CreateProposalRequestDTO	true	"Create proposal request"
//	@Success		201		{object}	domain.ProposalResponseDTO
//	@Failure		400		{object}	map[string]string
//	@Failure		403		{object}	map[string]string
//	@Failure		409		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/proposal [post]
func (ctrl *ProposalController) Create(c *gin.Context) {
	payload := c.MustGet("authorization_payload").(*auth.Payload)
	log := logger.Get().With(
		slog.String("layer", "controller"),
		slog.String("module", "proposal"),
		slog.String("operation", "Create"),
		slog.String("user_id", payload.UserID),
	)

	if payload.UserType != utils.UserTypeBusiness {
		log.Warn("proposal creation forbidden: only business users can create proposals", slog.String("user_type", payload.UserType))
		c.JSON(http.StatusForbidden, gin.H{"error": "only business users can create proposals"})
		return
	}

	var dto domain.CreateProposalRequestDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		log.Warn("invalid request body", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		log.Error("failed to parse user ID from token", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	response, err := ctrl.service.Create(dto, userID)
	if err != nil {
		if errors.Is(err, utils.ErrProposalAlreadyPendingForHelper) {
			log.Warn("proposal creation conflict: pending proposal already exists", slog.String("error", err.Error()))
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		log.Error("proposal creation failed with unexpected error", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetProposalByID godoc
//
//	@Summary		Get proposal by ID
//	@Description	Returns a proposal by ID. Only participants (business owner or helper) can access it.
//	@Tags			proposal
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Proposal ID (UUID)"
//	@Success		200	{object}	domain.ProposalResponseDTO
//	@Failure		400	{object}	map[string]string
//	@Failure		403	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/proposal/{id} [get]
func (ctrl *ProposalController) GetByID(c *gin.Context) {
	payload := c.MustGet("authorization_payload").(*auth.Payload)
	log := logger.Get().With(
		slog.String("layer", "controller"),
		slog.String("module", "proposal"),
		slog.String("operation", "GetByID"),
		slog.String("user_id", payload.UserID),
	)

	proposalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Warn("invalid proposal ID", slog.String("id_param", c.Param("id")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid proposal id"})
		return
	}

	requesterID, err := uuid.Parse(payload.UserID)
	if err != nil {
		log.Error("failed to parse requester ID from token", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	response, err := ctrl.service.GetByID(proposalID, requesterID)
	if err != nil {
		if errors.Is(err, utils.ErrProposalNotFound) {
			log.Warn("proposal not found", slog.String("proposal_id", proposalID.String()))
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrNotProposalParticipant) {
			log.Warn("proposal access denied: requester is not a participant", slog.String("proposal_id", proposalID.String()))
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		log.Error("failed to get proposal", slog.String("error", err.Error()), slog.String("proposal_id", proposalID.String()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// PatchProposalStatus godoc
//
//	@Summary		Update proposal status
//	@Description	Updates the status of a proposal following the allowed state machine transitions. Helpers manage most transitions; both parties can cancel.
//	@Tags			proposal
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string									true	"Proposal ID (UUID)"
//	@Param			body	body		domain.UpdateProposalStatusRequestDTO	true	"Update status request"
//	@Success		200		{object}	domain.ProposalResponseDTO
//	@Failure		400		{object}	map[string]string
//	@Failure		403		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		422		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/proposal/{id}/status [patch]
func (ctrl *ProposalController) UpdateStatus(c *gin.Context) {
	payload := c.MustGet("authorization_payload").(*auth.Payload)
	log := logger.Get().With(
		slog.String("layer", "controller"),
		slog.String("module", "proposal"),
		slog.String("operation", "UpdateStatus"),
		slog.String("user_id", payload.UserID),
	)

	proposalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Warn("invalid proposal ID", slog.String("id_param", c.Param("id")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid proposal id"})
		return
	}

	var dto domain.UpdateProposalStatusRequestDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		log.Warn("invalid request body", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	requesterID, err := uuid.Parse(payload.UserID)
	if err != nil {
		log.Error("failed to parse requester ID from token", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	response, err := ctrl.service.UpdateStatus(proposalID, dto, requesterID, payload.UserType)
	if err != nil {
		if errors.Is(err, utils.ErrProposalNotFound) {
			log.Warn("status update failed: proposal not found", slog.String("proposal_id", proposalID.String()))
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrProposalUnauthorized) {
			log.Warn("status update forbidden", slog.String("proposal_id", proposalID.String()), slog.String("error", err.Error()))
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrProposalInvalidStatus) || errors.Is(err, utils.ErrProposalFinished) {
			log.Warn("status update rejected: invalid transition", slog.String("proposal_id", proposalID.String()), slog.String("error", err.Error()))
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		log.Error("status update failed with unexpected error", slog.String("error", err.Error()), slog.String("proposal_id", proposalID.String()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ListProposals godoc
//
//	@Summary		List proposals
//	@Description	Lists proposals for the authenticated user. Business users see proposals they created; helpers see proposals directed at them. Supports optional status filter.
//	@Tags			proposal
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			status	query		string	false	"Filter by proposal status (pending, accepted, refused, in progress, cancelled, finished)"
//	@Success		200		{object}	[]domain.ProposalResponseDTO
//	@Failure		500		{object}	map[string]string
//	@Router			/proposal [get]
func (ctrl *ProposalController) List(c *gin.Context) {
	payload := c.MustGet("authorization_payload").(*auth.Payload)
	log := logger.Get().With(
		slog.String("layer", "controller"),
		slog.String("module", "proposal"),
		slog.String("operation", "List"),
		slog.String("user_id", payload.UserID),
	)

	statusFilter := c.Query("status")

	requesterID, err := uuid.Parse(payload.UserID)
	if err != nil {
		log.Error("failed to parse requester ID from token", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	p := parsePagination(c)

	response, err := ctrl.service.List(requesterID, payload.UserType, statusFilter, p)
	if err != nil {
		log.Error("failed to list proposals", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func parsePagination(c *gin.Context) domain.PaginationParams {
	page := 1
	pageSize := domain.DefaultPageSize
	if v, err := strconv.Atoi(c.Query("page")); err == nil && v > 0 {
		page = v
	}
	if v, err := strconv.Atoi(c.Query("page_size")); err == nil && v > 0 {
		pageSize = v
	}
	return domain.PaginationParams{Page: page, PageSize: pageSize}
}
