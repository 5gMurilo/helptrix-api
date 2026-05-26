package review

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/5gMurilo/helptrix-api/adapter/auth"
	"github.com/5gMurilo/helptrix-api/core/domain"
	reviewinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/review"
	"github.com/5gMurilo/helptrix-api/core/logger"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReviewController struct {
	svc reviewinterfaces.IReviewService
}

func NewReviewController(svc reviewinterfaces.IReviewService) reviewinterfaces.IReviewController {
	return &ReviewController{svc: svc}
}

// CreateReview godoc
//
//	@Summary		Create review
//	@Description	Business user creates a review for a helper after service completion
//	@Tags			review
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body	domain.CreateReviewRequestDTO	true	"Create review request"
//	@Success		201
//	@Failure		400	{object}	map[string]string
//	@Failure		403	{object}	map[string]string
//	@Failure		409	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/review [post]
func (ctrl *ReviewController) Create(c *gin.Context) {
	payload := c.MustGet("authorization_payload").(*auth.Payload)
	log := logger.Get().With(
		slog.String("layer", "controller"),
		slog.String("module", "review"),
		slog.String("operation", "Create"),
		slog.String("user_id", payload.UserID),
	)

	if payload.UserType != utils.UserTypeBusiness {
		log.Warn("review creation forbidden: only business users can create reviews", slog.String("user_type", payload.UserType))
		c.JSON(http.StatusForbidden, gin.H{"error": "only business users can create reviews"})
		return
	}

	businessID, err := uuid.Parse(payload.UserID)
	if err != nil {
		log.Error("failed to parse user ID from token", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var dto domain.CreateReviewRequestDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		log.Warn("invalid request body", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.svc.CreateReview(businessID, dto); err != nil {
		if errors.Is(err, utils.ErrCannotReviewSelf) {
			log.Warn("review creation forbidden: cannot review self")
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrProposalNotFound) ||
			errors.Is(err, utils.ErrReviewProposalMismatch) ||
			errors.Is(err, utils.ErrProposalNotFinished) {
			log.Warn("review creation forbidden: proposal constraint violation", slog.String("error", err.Error()))
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrReviewAlreadyExists) {
			log.Warn("review creation conflict: review already exists", slog.String("error", err.Error()))
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		log.Error("review creation failed with unexpected error", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "created"})
}

// ListBusinessReviews godoc
//
//	@Summary		List reviews made by business
//	@Description	Returns all reviews created by the authenticated business user
//	@Tags			review
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}	domain.ReviewListResponseDTO
//	@Failure		403	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/review/business [get]
func (ctrl *ReviewController) ListBusiness(c *gin.Context) {
	payload := c.MustGet("authorization_payload").(*auth.Payload)
	log := logger.Get().With(
		slog.String("layer", "controller"),
		slog.String("module", "review"),
		slog.String("operation", "ListBusiness"),
		slog.String("user_id", payload.UserID),
	)

	if payload.UserType != utils.UserTypeBusiness {
		log.Warn("list business reviews forbidden: only business users can access this endpoint", slog.String("user_type", payload.UserType))
		c.JSON(http.StatusForbidden, gin.H{"error": "only business users can access this endpoint"})
		return
	}

	businessID, err := uuid.Parse(payload.UserID)
	if err != nil {
		log.Error("failed to parse user ID from token", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	reviews, err := ctrl.svc.ListBusinessReviews(businessID, parsePagination(c))
	if err != nil {
		log.Error("failed to list business reviews", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// ListHelperReviews godoc
//
//	@Summary		List reviews received by helper
//	@Description	Returns all reviews received by the authenticated helper user
//	@Tags			review
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}	domain.ReviewListResponseDTO
//	@Failure		403	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/review/helper [get]
func (ctrl *ReviewController) ListHelper(c *gin.Context) {
	payload := c.MustGet("authorization_payload").(*auth.Payload)
	log := logger.Get().With(
		slog.String("layer", "controller"),
		slog.String("module", "review"),
		slog.String("operation", "ListHelper"),
		slog.String("user_id", payload.UserID),
	)

	if payload.UserType != utils.UserTypeHelper {
		log.Warn("list helper reviews forbidden: only helper users can access this endpoint", slog.String("user_type", payload.UserType))
		c.JSON(http.StatusForbidden, gin.H{"error": "only helper users can access this endpoint"})
		return
	}

	helperID, err := uuid.Parse(payload.UserID)
	if err != nil {
		log.Error("failed to parse user ID from token", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	reviews, err := ctrl.svc.ListHelperReviews(helperID, parsePagination(c))
	if err != nil {
		log.Error("failed to list helper reviews", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, reviews)
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
