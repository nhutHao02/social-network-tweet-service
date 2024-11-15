package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/nhutHao02/social-network-common-service/middleware"
	"github.com/nhutHao02/social-network-common-service/utils/logger"
)

func MapRoutes(
	router *gin.Engine,
	tweetHandler *TweetHandler,
) {
	v1 := router.Group("/api/v1")
	{
		v1.Use(middleware.JwtAuthMiddleware(logger.GetDefaultLogger()))
		{
			vTweet := v1.Group("/tweet")
			vTweet.GET("", tweetHandler.GetTweetByUserID)

		}
	}
	// vGuest := router.Group("api/v1/guest")
	// {

	// }

}
