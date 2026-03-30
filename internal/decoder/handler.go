package decoder

import (
	"errors"
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
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "slug can't be empty"})
		return
	}

	content, err := h.service.Decode(ctx, slug)
	if errors.Is(err, ErrNotFound) {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("no data under slug %s", slug)})
		return
	}

	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("something went wrong: %s", err.Error())})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"content": content})
}
