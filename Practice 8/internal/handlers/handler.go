package handlers

import (
	"net/http"

	"practice-8/internal/service"
	"practice-8/pkg/modules"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.UserService
}

func NewHandler(svc service.UserService) *Handler {
	return &Handler{service: svc}
}

func (h *Handler) Register(c *gin.Context) {
	var req modules.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.service.Register(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *Handler) Login(c *gin.Context) {
	var req modules.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, user, err := h.service.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, modules.AuthResponse{
		SessionID: token,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
	})
}

func (h *Handler) GetRate(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")

	if from == "" || to == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from and to parameters required"})
		return
	}

	rate, err := h.service.GetRate(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"from": from, "to": to, "rate": rate})
}
