package encoder

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Request struct {
	Content string `json:"content"`
}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Encode(ctx *gin.Context) {
	var request Request

	if err := ctx.BindJSON(&request); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, fmt.Sprintf("failed to map JSON body to request: %s", err.Error()))
		return
	}

	slug, err := h.service.Encode(ctx, request.Content)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, fmt.Sprintf("failed to encode content '%s': %s", request.Content, err.Error()))
		return
	}

	ctx.IndentedJSON(http.StatusOK, slug)
}
