package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MapRoutes(
	router *gin.Engine,
) {

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}