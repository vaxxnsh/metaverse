package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vaxxnsh/metaverse/api/internal/domain"
	"github.com/vaxxnsh/metaverse/api/internal/lib/response"
	"github.com/vaxxnsh/metaverse/api/internal/middleware"
	"github.com/vaxxnsh/metaverse/api/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type patchUserMetadataRequest struct {
	AvatarId string `json:"avatarId" binding:"required"`
}

func (h *UserHandler) PatchMetadata(c *gin.Context) {
	userId := c.GetString(middleware.UserIDKey)

	var req patchUserMetadataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_PAYLOAD", struct{}{})
		return
	}

	user, err := h.userService.PatchMetadata(c.Request.Context(), userId, req.AvatarId)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			response.SendError(c, http.StatusNotFound, "NOT_FOUND", "user not found", struct{}{})
			return
		}
		if errors.Is(err, domain.ErrAvatarNotFound) {
			response.SendError(c, http.StatusNotFound, "NOT_FOUND", "avatar not found", struct{}{})
			return
		}
		response.SendError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), struct{}{})
		return
	}

	response.SendSuccess(c, user, http.StatusOK)
}
