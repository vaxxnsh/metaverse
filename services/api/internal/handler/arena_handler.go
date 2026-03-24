package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vaxxnsh/metaverse/api/internal/domain"
	"github.com/vaxxnsh/metaverse/api/internal/lib/response"
	"github.com/vaxxnsh/metaverse/api/internal/service"
)

type ArenaHandler struct {
	arenaService service.ArenaService
}

func NewArenaHandler(arenaService service.ArenaService) *ArenaHandler {
	return &ArenaHandler{arenaService: arenaService}
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
