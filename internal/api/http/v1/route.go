package v1

import (
	_ "github.com/nhutHao02/social-network-tweet-service/docs"

	"github.com/gin-gonic/gin"
	"github.com/nhutHao02/social-network-common-service/middleware"
	"github.com/nhutHao02/social-network-common-service/utils/logger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
			vTweet.POST("", tweetHandler.PostTweet)
			vTweet.GET("/all", tweetHandler.GetAllTweets)
			vTweet.GET("/love", tweetHandler.GetLoveTweetsByUserID)
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
