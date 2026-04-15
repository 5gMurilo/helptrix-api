package review

import (
	"errors"
	"net/http"

	"github.com/5gMurilo/helptrix-api/adapter/auth"
	"github.com/5gMurilo/helptrix-api/core/domain"
	reviewinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/review"
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

	if payload.UserType != utils.UserTypeBusiness {
		c.JSON(http.StatusForbidden, gin.H{"error": "only business users can create reviews"})
		return
	}

	businessID, err := uuid.Parse(payload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var dto domain.CreateReviewRequestDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.svc.CreateReview(businessID, dto); err != nil {
		if errors.Is(err, utils.ErrCannotReviewSelf) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "no finished proposal found for this business and helper" {
			c.JSON(http.StatusForbidden, gin.H{"error": utils.ErrProposalNotFinished.Error()})
			return
		}
		if err.Error() == "business has already reviewed this helper" {
			c.JSON(http.StatusConflict, gin.H{"error": utils.ErrReviewAlreadyExists.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Status(http.StatusCreated)
}

// ListBusinessReviews godoc

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

	if payload.UserType != utils.UserTypeBusiness {
		c.JSON(http.StatusForbidden, gin.H{"error": "only business users can access this endpoint"})
		return
	}

	businessID, err := uuid.Parse(payload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	reviews, err := ctrl.svc.ListBusinessReviews(businessID)
	if err != nil {
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

	if payload.UserType != utils.UserTypeHelper {
		c.JSON(http.StatusForbidden, gin.H{"error": "only helper users can access this endpoint"})
		return
	}

	helperID, err := uuid.Parse(payload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	reviews, err := ctrl.svc.ListHelperReviews(helperID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}
