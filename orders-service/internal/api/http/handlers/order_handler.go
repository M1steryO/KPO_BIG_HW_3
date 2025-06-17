package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"orders-service/internal/services"
	"orders-service/internal/storage"
	"strconv"
	"time"
)

type Handler struct {
	svc *services.OrderService
}

func NewHandler(store storage.OrderStorage) *Handler {
	return &Handler{svc: services.NewOrderService(store)}
}

func (h *Handler) Create(c *gin.Context) {
	var req struct {
		UserID int64   `json:"user_id" binding:"required"`
		Amount float64 `json:"amount" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	order, err := h.svc.Create(req.UserID, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, order)
}

func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	order, err := h.svc.Get(id)
	if err != nil {
		if errors.Is(err, storage.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

	}
	updatedAtTime := time.Unix(order.CreatedAt, 0).Local()
	formattedTime := updatedAtTime.Format(time.RFC3339)
	c.JSON(http.StatusOK, gin.H{
		"id":         order.ID,
		"status":     order.Status,
		"amount":     order.Amount,
		"user_id":    order.UserID,
		"message":    order.Message,
		"created_at": formattedTime,
	})
}

func (h *Handler) List(c *gin.Context) {
	userId := c.Query("user_id")
	id, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}
	orders, err := h.svc.List(id)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("user with id %d haven't got any orders", id)})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}
