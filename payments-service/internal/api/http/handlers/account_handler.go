package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"payments-service/internal/storage"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"payments-service/internal/services"
)

type AccountHandler struct{ svc *services.AccountService }

func NewAccountHandler(svc *services.AccountService) *AccountHandler { return &AccountHandler{svc} }

func (h *AccountHandler) CreateAccount(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.Param("user_id"), 10, 64)
	accountId, err := h.svc.Create(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": accountId})
}

func (h *AccountHandler) Deposit(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.Param("user_id"), 10, 64)
	var req struct {
		Amount float64 `json:"amount" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.Deposit(userId, req.Amount); err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			c.JSON(http.StatusNotFound,
				gin.H{"error": fmt.Sprintf("user with id %d not found", userId)})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *AccountHandler) Balance(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.Param("user_id"), 10, 64)
	acct, err := h.svc.Get(userId)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("user with id %d not found", userId)})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	updatedAtTime := time.Unix(acct.UpdatedAt, 0).Local()
	formattedTime := updatedAtTime.Format(time.RFC3339)

	c.JSON(http.StatusOK, gin.H{
		"user_id":    acct.UserID,
		"balance":    acct.Balance,
		"updated_at": formattedTime,
	})
}
