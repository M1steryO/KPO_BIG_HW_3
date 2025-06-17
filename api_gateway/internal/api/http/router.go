package http

import (
	"api_gateway/internal/api/http/handlers"
	"api_gateway/internal/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewGatewayRouter(cfg *config.Config) *gin.Engine {
	h := &handlers.GatewayHandler{
		OrderURL:   cfg.OrderServiceURL,
		PaymentURL: cfg.PaymentServiceURL,
		Client:     &http.Client{},
	}

	r := gin.Default()

	r.POST("/orders", h.ProxyOrderCreate)
	r.GET("/orders/:id", h.ProxyOrderGet)
	r.GET("/orders", h.ProxyOrderList)

	r.POST("/accounts/:user_id", h.ProxyAccountCreate)
	r.POST("/accounts/:user_id/deposit", h.ProxyAccountDeposit)
	r.GET("/accounts/:user_id/balance", h.ProxyAccountBalance)

	return r
}
