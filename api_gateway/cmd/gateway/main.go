package main

import (
	_ "api_gateway/docs"
	gatewayHttp "api_gateway/internal/api/http"
	"api_gateway/internal/config"
	files "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"log"
	_ "net/http/pprof"
)

// @title API Gateway
// @version 1.0
// @description API Gateway for Orders and Payments Services
// @host localhost:8083
// @BasePath /
func main() {
	cfg, err := config.MustLoad()
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	r := gatewayHttp.NewGatewayRouter(cfg)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))

	log.Printf("API Gateway listening on %s", cfg.GatewayAddr)
	if err := r.Run(cfg.GatewayAddr); err != nil {
		log.Fatalf("gateway server error: %v", err)
	}
	log.Printf("server shutdown!")
}
