package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vaxxnsh/metaverse/api/internal/lib/response"
	"github.com/vaxxnsh/metaverse/api/internal/service"
)

type MetadataHandler struct {
	metadataService service.MetaDataService
}

func NewMetadataHandler(metadataService service.MetaDataService) *MetadataHandler {
	return &MetadataHandler{metadataService: metadataService}
}

type bulkAvatarsRequest struct {
	UserIds []string `json:"userIds" binding:"required"`
}

func (h *MetadataHandler) GetBulkUserAvatars(c *gin.Context) {
	var req bulkAvatarsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_PAYLOAD", struct{}{})
		return
	}

	avatars, err := h.metadataService.GetBulkUserAvatars(c.Request.Context(), req.UserIds)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), struct{}{})
		return
	}

	response.SendSuccess(c, gin.H{"avatars": avatars}, http.StatusOK)
}
