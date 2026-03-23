package decrypt

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

func (h *Handler) Decrypt(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.IndentedJSON(http.StatusBadRequest, "decrypt: slug can't be empty")
		return
	}

	content, err := h.service.Decrypt(slug)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, fmt.Sprintf("decrypt: failed to decrypt slug '%s': %s", slug, err.Error()))
		return
	}

	c.String(http.StatusOK, content)
}
