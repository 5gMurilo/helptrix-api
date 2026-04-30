package auth

import (
	"errors"
	"net/http"

	"github.com/5gMurilo/helptrix-api/core/domain"
	authinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/auth"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service authinterfaces.IAuthService
}

func NewAuthController(service authinterfaces.IAuthService) authinterfaces.IAuthController {
	return &AuthController{service: service}
}

// Register godoc
//
//	@Summary		Register a new user
//	@Description	Creates a helper or business user with address and categories
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.RegisterRequestDTO	true	"Register request"
//	@Success		201		{object}	domain.RegisterResponseDTO
//	@Failure		400		{object}	map[string]string
//	@Failure		409		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/auth/register [post]
func (ctrl *AuthController) Register(c *gin.Context) {
	var dto domain.RegisterRequestDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := ctrl.service.Register(dto)
	if err != nil {
		if errors.Is(err, utils.ErrUserAlreadyRegistered) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error" + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Login godoc
//
//	@Summary		Authenticate user
//	@Description	Authenticates a user by email and password, returns a Paseto token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.LoginRequestDTO	true	"Login request"
//	@Success		201		{object}	domain.LoginResponseDTO
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Router			/auth/login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	var dto domain.LoginRequestDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := ctrl.service.Login(dto)
	if err != nil {
		if errors.Is(err, utils.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, response)
}
