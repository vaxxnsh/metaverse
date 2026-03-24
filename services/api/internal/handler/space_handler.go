package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vaxxnsh/metaverse/api/internal/lib/response"
	"github.com/vaxxnsh/metaverse/api/internal/middleware"
	"github.com/vaxxnsh/metaverse/api/internal/service"
)

type SpaceHandler struct {
	spaceService service.SpaceService
}

func NewSpaceHandler(spaceService service.SpaceService) *SpaceHandler {
	return &SpaceHandler{spaceService: spaceService}
}

type createSpaceRequest struct {
	Name   string `json:"name"   binding:"required"`
	Width  int32  `json:"width"  binding:"required"`
	Height int32  `json:"height" binding:"required"`
	MapId  string `json:"mapId"`
}

func (h *SpaceHandler) CreateSpace(c *gin.Context) {
	userId := c.GetString(middleware.UserIDKey)

	var req createSpaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_PAYLOAD", struct{}{})
		return
	}

	space, err := h.spaceService.CreateSpace(c.Request.Context(), userId, req.Name, req.Width, req.Height, req.MapId)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), struct{}{})
		return
	}

	response.SendSuccess(c, gin.H{"space": space}, http.StatusCreated)
}

func (h *SpaceHandler) DeleteSpace(c *gin.Context) {
	spaceId := c.Param("spaceId")

	if err := h.spaceService.DeleteSpace(c.Request.Context(), spaceId); err != nil {
		response.SendError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), struct{}{})
		return
	}

	response.SendSuccess(c, gin.H{}, http.StatusOK)
}

func (h *SpaceHandler) GetMySpaces(c *gin.Context) {
	userId := c.GetString(middleware.UserIDKey)

	spaces, err := h.spaceService.GetSpacesByCreator(c.Request.Context(), userId)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), struct{}{})
		return
	}

	response.SendSuccess(c, gin.H{"spaces": spaces}, http.StatusOK)
}
