package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vaxxnsh/metaverse/api/internal/service"
)

type Handler struct {
	service service.Service
}

func NewHandler(s service.Service) *Handler {
	return &Handler{service: s}
}

type createUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req createUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	user, err := h.service.CreateUser(c.Request.Context(), req.Email, req.Name)
	if err != nil {
		switch err {
		case service.ErrInvalidEmail, service.ErrInvalidName:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case service.ErrUserExists:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         user.ID,
		"email":      user.Email,
		"name":       user.Name,
		"created_at": user.CreatedAt,
	})
}
