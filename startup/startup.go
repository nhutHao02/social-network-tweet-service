package startup

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nhutHao02/social-network-tweet-service/config"
)

func StartServer() {
	// load congig
	cfg := config.LoadConfig()

	// database setup

	// setup route
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// setup server
	r.Run(cfg.HTTPServer.Address)

}
