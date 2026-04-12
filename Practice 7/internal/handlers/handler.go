package handlers

import (
	"net/http"
	"strconv"

	"practice-7/internal/usecase"
	"practice-7/pkg/modules"
	"practice-7/pkg/utils"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase *usecase.UserUsecase
}

func NewHandler(usecase *usecase.UserUsecase) *Handler {
	return &Handler{usecase: usecase}
}

// POST /register
func (h *Handler) Register(c *gin.Context) {
	var req modules.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.usecase.Register(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// POST /login
func (h *Handler) Login(c *gin.Context) {
	var req modules.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, user, err := h.usecase.Login(req)
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

// GET /me
func (h *Handler) GetMe(c *gin.Context) {
	// Get user ID from JWT claims (set by JWT middleware)
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userClaims, ok := claims.(*utils.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	user, err := h.usecase.GetUserByID(userClaims.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// PATCH /users/promote/:id
func (h *Handler) PromoteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.usecase.PromoteUserToAdmin(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user promoted to admin"})
}
