package http

import (
	"github.com/gin-gonic/gin"
	"github.com/nhutHao02/social-network-common-service/utils/logger"
	"github.com/nhutHao02/social-network-tweet-service/config"
	v1 "github.com/nhutHao02/social-network-tweet-service/internal/api/http/v1"
	"go.uber.org/zap"
)

type HTTPServer struct {
	Cfg *config.Config
	// handlers
	// example addressHandler at V1
}

func NewHTTPServer(cfg *config.Config) *HTTPServer {
	return &HTTPServer{Cfg: cfg}
}

func (s *HTTPServer) RunHTTPServer() error {
	r := gin.Default()
	v1.MapRoutes(r)
	logger.Info("HTTP Server server listening at" + s.Cfg.HTTPServer.Address)
	err := r.Run(s.Cfg.HTTPServer.Address)
	if err != nil {
		logger.Error("HTTP Server error", zap.Error(err))
		return err
	}
	return nil
}
