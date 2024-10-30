package http

import (
	"github.com/gin-gonic/gin"
	"github.com/nhutHao02/social-network-tweet-service/config"
	v1 "github.com/nhutHao02/social-network-tweet-service/internal/api/http/v1"
)

type HTTPServer struct {
	cfg *config.Config
	// handlers
	// example addressHandler at V1
}

func NewHTTPServer(cfg *config.Config) *HTTPServer {
	return &HTTPServer{cfg: cfg}
}

func (s *HTTPServer) RunHTTPServer() {
	r := gin.Default()
	v1.MapRoutes(r)
	r.Run(s.cfg.HTTPServer.Address)
}
