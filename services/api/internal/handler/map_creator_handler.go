package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vaxxnsh/metaverse/api/internal/lib/response"
	"github.com/vaxxnsh/metaverse/api/internal/repository"
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

type createAvatarRequest struct {
	ImageUrl string `json:"imageUrl" binding:"required"`
	Name     string `json:"name"`
}

func (h *MapCreatorHandler) CreateAvatar(c *gin.Context) {
	var req createAvatarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_PAYLOAD", struct{}{})
		return
	}

	avatar, err := h.mapCreatorService.CreateAvatar(c.Request.Context(), req.ImageUrl, req.Name)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), struct{}{})
		return
	}

	response.SendSuccess(c, gin.H{"id": avatar.ID}, http.StatusCreated)
}

type defaultElementRequest struct {
	ElementId string `json:"elementId" binding:"required"`
	X         int32  `json:"x"`
	Y         int32  `json:"y"`
}

type createMapRequest struct {
	Name            string                 `json:"name"      binding:"required"`
	Thumbnail       string                 `json:"thumbnail"`
	Width           int32                  `json:"width"     binding:"required"`
	Height          int32                  `json:"height"    binding:"required"`
	DefaultElements []defaultElementRequest `json:"defaultElements"`
}

func (h *MapCreatorHandler) CreateMap(c *gin.Context) {
	var req createMapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_PAYLOAD", struct{}{})
		return
	}

	defaultElements := make([]repository.DefaultElement, 0, len(req.DefaultElements))
	for _, e := range req.DefaultElements {
		defaultElements = append(defaultElements, repository.DefaultElement{
			ElementId: e.ElementId,
			X:         e.X,
			Y:         e.Y,
		})
	}

	m, err := h.mapCreatorService.CreateMap(c.Request.Context(), req.Name, req.Width, req.Height, defaultElements)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), struct{}{})
		return
	}

	response.SendSuccess(c, gin.H{"id": m.ID}, http.StatusCreated)
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
