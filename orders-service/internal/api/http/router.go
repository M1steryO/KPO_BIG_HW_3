package http

import (
	"github.com/gin-gonic/gin"
	"orders-service/internal/api/http/handlers"
	"orders-service/internal/storage"
)

func NewOrderRouter(store storage.OrderStorage) *gin.Engine {
	r := gin.Default()
	h := handlers.NewHandler(store)

	grp := r.Group("/orders")
	grp.POST("", h.Create)
	grp.GET("/:id", h.Get)
	grp.GET("", h.List)

	return r
}
