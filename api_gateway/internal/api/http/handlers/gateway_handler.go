package handlers

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type GatewayHandler struct {
	OrderURL   string
	PaymentURL string
	Client     *http.Client
}
type OrderCreateRequest struct {
	UserID int64   `json:"user_id" example:"42"`
	Amount float64 `json:"amount" example:"123.45"`
}

// swagger:model DepositRequest
type DepositRequest struct {
	Amount float64 `json:"amount" example:"50.5"`
}

// ProxyOrderCreate godoc
// @Summary      Create a new order
// @Description  Proxy to Order Service to create a new order
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        order  body   handlers.OrderCreateRequest  true  "Order payload"
// @Router       /orders [post]
func (h *GatewayHandler) ProxyOrderCreate(c *gin.Context) {
	h.reverseProxy(c, http.MethodPost, h.OrderURL+"/orders")
}

// ProxyOrderGet godoc
// @Summary      Get order by ID
// @Description  Proxy to Order Service to fetch an order by ID
// @Tags         orders
// @Produce      json
// @Param        id     path    int     true  "Order ID"
// @Router       /orders/{id} [get]
func (h *GatewayHandler) ProxyOrderGet(c *gin.Context) {
	id := c.Param("id")
	h.reverseProxy(c, http.MethodGet, h.OrderURL+"/orders/"+id)
}

// ProxyOrderList godoc
// @Summary      List orders for a user
// @Description  Proxy to Order Service to list all orders of a given user
// @Tags         orders
// @Produce      json
// @Param        user_id  query   int     true  "User ID"
// @Router       /orders [get]
func (h *GatewayHandler) ProxyOrderList(c *gin.Context) {
	h.reverseProxy(c, http.MethodGet, h.OrderURL+"/orders?user_id="+c.Query("user_id"))
}

// ProxyAccountCreate godoc
// @Summary      Create a new account
// @Description  Proxy to Payment Service to create a new account
// @Tags         accounts
// @Produce      json
// @Param        user_id  path    int  true  "User ID"
// @Router       /accounts/{user_id} [post]
func (h *GatewayHandler) ProxyAccountCreate(c *gin.Context) {
	uid := c.Param("user_id")
	h.reverseProxy(c, http.MethodPost, h.PaymentURL+"/accounts/"+uid)
}

// ProxyAccountDeposit godoc
// @Summary      Deposit to account
// @Description  Proxy to Payment Service to deposit funds into a user's account
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        user_id  path    int  true  "User ID"
// @Param        amount   body    handlers.DepositRequest  true  "Amount to deposit"
// @Router       /accounts/{user_id}/deposit [post]
func (h *GatewayHandler) ProxyAccountDeposit(c *gin.Context) {
	userId := c.Param("user_id")
	h.reverseProxy(c, http.MethodPost, h.PaymentURL+"/accounts/"+userId+"/deposit")
}

// ProxyAccountBalance godoc
// @Summary      Get account balance
// @Description  Proxy to Payment Service to retrieve the balance of a user's account
// @Tags         accounts
// @Produce      json
// @Param        user_id  path    int  true  "User ID"
// @Router       /accounts/{user_id}/balance [get]
func (h *GatewayHandler) ProxyAccountBalance(c *gin.Context) {
	uid := c.Param("user_id")
	h.reverseProxy(c, http.MethodGet, h.PaymentURL+"/accounts/"+uid+"/balance")
}

func (h *GatewayHandler) reverseProxy(c *gin.Context, method, target string) {
	req, err := http.NewRequest(method, target, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req.Header = c.Request.Header.Clone()

	resp, err := h.Client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		c.Writer.Header()[k] = v
	}
	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}
