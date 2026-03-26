package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vaxxnsh/shared/domain"
	"github.com/vaxxnsh/metaverse/api/internal/lib/response"
	"github.com/vaxxnsh/metaverse/api/internal/middleware"
	"github.com/vaxxnsh/metaverse/api/internal/service"
)

type AdminHandler struct {
	adminService service.AdminService
}

func NewAdminHandler(adminService service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

type patchAdminMetadataRequest struct {
	AvatarId string `json:"avatarId" binding:"required"`
}

func (h *AdminHandler) PatchMetadata(c *gin.Context) {
	adminId := c.GetString(middleware.AdminIDKey)

	var req patchAdminMetadataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_PAYLOAD", struct{}{})
		return
	}

	admin, err := h.adminService.PatchMetadata(c.Request.Context(), adminId, req.AvatarId)
	if err != nil {
		if errors.Is(err, domain.ErrAdminNotFound) {
			response.SendError(c, http.StatusNotFound, "NOT_FOUND", "admin not found", struct{}{})
			return
		}
		if errors.Is(err, domain.ErrAvatarNotFound) {
			response.SendError(c, http.StatusNotFound, "NOT_FOUND", "avatar not found", struct{}{})
			return
		}
		response.SendError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), struct{}{})
		return
	}

	response.SendSuccess(c, admin, http.StatusOK)
}
