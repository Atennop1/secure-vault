package decoder

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Decode(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if slug == "" {
		ctx.IndentedJSON(http.StatusBadRequest, "slug can't be empty")
		return
	}

	content, err := h.service.Decode(ctx, slug)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, fmt.Sprintf("failed to decode slug '%s': %s", slug, err.Error()))
		return
	}

	ctx.String(http.StatusOK, content)
}
