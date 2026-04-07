package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"taoshop/internal/service"
)

type OrderHandler struct {
	orders *service.OrderService
}

func NewOrderHandler(orders *service.OrderService) *OrderHandler {
	return &OrderHandler{orders: orders}
}

type createOrderRequest struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int   `json:"quantity" binding:"required,min=1,max=10"`
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req createOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	userID := c.GetInt64("user_id")
	order, err := h.orders.Create(c.Request.Context(), userID, req.ProductID, req.Quantity)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrInsufficientStock) {
			status = http.StatusBadRequest
		}

		c.JSON(status, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "purchase success",
		"order":   order,
	})
}

func (h *OrderHandler) Mine(c *gin.Context) {
	userID := c.GetInt64("user_id")
	orders, err := h.orders.ListByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": orders})
}
