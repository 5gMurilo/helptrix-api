package uploader

import (
	"errors"
	"io"
	"net/http"

	"github.com/5gMurilo/helptrix-api/adapter/auth"
	uploaderinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/uploader"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	allowedImageTypes    = []string{"profile-images", "service-images"}
	allowedContentTypes  = []string{"image/jpeg", "image/png", "image/webp"}
	maxFileSizeBytes    int64 = 5 * 1024 * 1024
)

type UploaderController struct {
	svc uploaderinterfaces.IUploaderService
}

func NewUploaderController(svc uploaderinterfaces.IUploaderService) uploaderinterfaces.IUploaderController {
	return &UploaderController{svc: svc}
}

func containsString(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}

// Upload godoc
//
//	@Summary		Upload an image
//	@Description	Uploads an image file to the storage. Supports profile-images and service-images types.
//	@Tags			image-uploader
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		BearerAuth
//	@Param			image-type	path		string	true	"Image type (profile-images | service-images)"
//	@Param			id			path		string	true	"User ID (profile-images) or service ID (service-images)"
//	@Param			image		formData	file	true	"Image file"
//	@Success		200			{object}	map[string]string
//	@Failure		400			{object}	map[string]string
//	@Failure		401			{object}	map[string]string
//	@Failure		403			{object}	map[string]string
//	@Failure		404			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/image-uploader/{image-type}/{id} [post]
func (ctrl *UploaderController) Upload(c *gin.Context) {
	payload := c.MustGet("authorization_payload").(*auth.Payload)

	imageType := c.Param("image-type")
	if !containsString(allowedImageTypes, imageType) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image type"})
		return
	}

	targetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	requesterID, err := uuid.Parse(payload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid requester id"})
		return
	}

	header, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
		return
	}

	if header.Size > maxFileSizeBytes {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file exceeds 5MB limit"})
		return
	}

	contentType := header.Header.Get("Content-Type")
	if !containsString(allowedContentTypes, contentType) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported image content type"})
		return
	}

	file, err := header.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	url, err := ctrl.svc.Upload(c.Request.Context(), imageType, requesterID, targetID, header.Filename, data, contentType)
	if err != nil {
		if errors.Is(err, utils.ErrNotOwner) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrUserNotFound) || errors.Is(err, utils.ErrServiceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrInvalidImageType) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}
