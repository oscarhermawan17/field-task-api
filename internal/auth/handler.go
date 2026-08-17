package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
	secret  []byte
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string `json:"accessToken"`
	User        User   `json:"user"`
}

func NewHandler(service *Service, secret []byte) *Handler {
	return &Handler{service: service, secret: secret}
}

func (h *Handler) Login(c *gin.Context) {
	var request loginRequest
	if err := c.ShouldBindJSON(&request); err != nil || strings.TrimSpace(request.Email) == "" || request.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "email and password are required"})
		return
	}

	user, err := h.service.Authenticate(request.Email, request.Password)
	if errors.Is(err, ErrInvalidCredentials) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid email or password"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "unable to sign in"})
		return
	}

	accessToken, err := CreateAccessToken(user.ID, h.secret, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "unable to create access token"})
		return
	}

	c.JSON(http.StatusOK, loginResponse{AccessToken: accessToken, User: user})
}

func (h *Handler) Me(c *gin.Context) {
	userID := c.GetString("authUserID")
	user, ok := h.service.UserByID(userID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "user session is no longer valid"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
