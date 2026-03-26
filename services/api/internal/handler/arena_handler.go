package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vaxxnsh/shared/domain"
	"github.com/vaxxnsh/metaverse/api/internal/lib/response"
	"github.com/vaxxnsh/metaverse/api/internal/service"
)

type ArenaHandler struct {
	arenaService service.ArenaService
}

func NewArenaHandler(arenaService service.ArenaService) *ArenaHandler {
	return &ArenaHandler{arenaService: arenaService}
}

type addSpaceElementRequest struct {
	ElementId string `json:"elementId" binding:"required"`
	SpaceId   string `json:"spaceId"   binding:"required"`
	X         int32  `json:"x"         binding:"required"`
	Y         int32  `json:"y"         binding:"required"`
}

func (h *ArenaHandler) AddElementToSpace(c *gin.Context) {
	var req addSpaceElementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_PAYLOAD", struct{}{})
		return
	}

	se, err := h.arenaService.AddElementToSpace(c.Request.Context(), req.SpaceId, req.ElementId, req.X, req.Y)
	if err != nil {
		switch err {
		case domain.ErrInvalidSpaceID, domain.ErrInvalidElementID:
			response.SendError(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), struct{}{})
		default:
			response.SendError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), struct{}{})
		}
		return
	}

	response.SendSuccess(c, gin.H{"id": se.ElementID}, http.StatusCreated)
}

type deleteSpaceElementRequest struct {
	SpaceId string `json:"spaceId" binding:"required"`
	X       int32  `json:"x"`
	Y       int32  `json:"y"`
}

func (h *ArenaHandler) DeleteSpaceElement(c *gin.Context) {
	var req deleteSpaceElementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_PAYLOAD", struct{}{})
		return
	}

	if err := h.arenaService.DeleteSpaceElement(c.Request.Context(), req.SpaceId, req.X, req.Y); err != nil {
		switch err {
		case domain.ErrInvalidSpaceID:
			response.SendError(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), struct{}{})
		default:
			response.SendError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), struct{}{})
		}
		return
	}

	response.SendSuccess(c, gin.H{}, http.StatusOK)
}

func (h *ArenaHandler) GetAllElements(c *gin.Context) {
	elements, err := h.arenaService.GetAllElements(c.Request.Context())
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), struct{}{})
		return
	}

	response.SendSuccess(c, gin.H{"elements": elements}, http.StatusOK)
}

func (h *ArenaHandler) GetSpace(c *gin.Context) {
	spaceId := c.Param("spaceId")

	space, err := h.arenaService.GetSpaceWithElements(c.Request.Context(), spaceId)
	if err != nil {
		switch err {
		case domain.ErrInvalidSpaceID, domain.ErrSpaceNotFound:
			response.SendError(c, http.StatusNotFound, "NOT_FOUND", err.Error(), struct{}{})
		default:
			response.SendError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error(), struct{}{})
		}
		return
	}

	response.SendSuccess(c, space, http.StatusOK)
}
