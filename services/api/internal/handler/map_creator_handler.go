package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vaxxnsh/metaverse/api/internal/lib/response"
	"github.com/vaxxnsh/metaverse/api/internal/service"
)

type MapCreatorHandler struct {
	mapCreatorService service.MapCreatorService
}

func NewMapCreatorHandler(mapCreatorService service.MapCreatorService) *MapCreatorHandler {
	return &MapCreatorHandler{mapCreatorService: mapCreatorService}
}

type createElementRequest struct {
	ImageUrl string `json:"imageUrl" binding:"required"`
	Width    int32  `json:"width"    binding:"required"`
	Height   int32  `json:"height"   binding:"required"`
	Static   bool   `json:"static"`
}

type updateElementRequest struct {
	ImageUrl string `json:"imageUrl" binding:"required"`
}

func (h *MapCreatorHandler) UpdateElement(c *gin.Context) {
	elementId := c.Param("elementId")

	var req updateElementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_PAYLOAD", struct{}{})
		return
	}

	_, err := h.mapCreatorService.UpdateElementImage(c.Request.Context(), elementId, req.ImageUrl)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), struct{}{})
		return
	}

	response.SendSuccess(c, gin.H{}, http.StatusOK)
}

func (h *MapCreatorHandler) CreateElement(c *gin.Context) {
	var req createElementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_PAYLOAD", struct{}{})
		return
	}

	element, err := h.mapCreatorService.CreateElement(c.Request.Context(), req.ImageUrl, req.Width, req.Height, req.Static)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), struct{}{})
		return
	}

	response.SendSuccess(c, gin.H{"id": element.ID}, http.StatusCreated)
}
