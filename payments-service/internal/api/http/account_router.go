package http

import (
	"github.com/gin-gonic/gin"
	"payments-service/internal/api/http/handlers"
	"payments-service/internal/services"
	"payments-service/internal/storage"
)

func NewAccountRouter(store storage.AccountStorage) *gin.Engine {
	r := gin.Default()
	svc := services.NewAccountService(store)
	h := handlers.NewAccountHandler(svc)
	r.POST("/accounts/:user_id", h.CreateAccount)
	r.POST("/accounts/:user_id/deposit", h.Deposit)
	r.GET("/accounts/:user_id/balance", h.Balance)
	return r
}
