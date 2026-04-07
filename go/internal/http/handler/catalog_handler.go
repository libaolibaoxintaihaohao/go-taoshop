package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"taoshop/internal/service"
)

type CatalogHandler struct {
	catalog *service.CatalogService
}

func NewCatalogHandler(catalog *service.CatalogService) *CatalogHandler {
	return &CatalogHandler{catalog: catalog}
}

func (h *CatalogHandler) ListProducts(c *gin.Context) {
	products, err := h.catalog.ListProducts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": products})
}
