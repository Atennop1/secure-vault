package encrypt

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

func (h *Handler) Encrypt(c *gin.Context) {
	var request Request

	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusBadRequest, fmt.Sprintf("failed to map JSON body to request: %s", err.Error()))
		return
	}

	slug, err := h.service.Encrypt(request.Content)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, fmt.Sprintf("failed to encrypt content '%s': %s", request.Content, err.Error()))
		return
	}

	c.IndentedJSON(http.StatusOK, slug)
}
