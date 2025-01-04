package v1

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nhutHao02/social-network-common-service/middleware"
	"github.com/nhutHao02/social-network-common-service/utils/logger"
	_ "github.com/nhutHao02/social-network-tweet-service/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func MapRoutes(
	router *gin.Engine,
	tweetHandler *TweetHandler,
) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	v1 := router.Group("/api/v1")
	{
		v1.GET("/comment-tweet-websocket", tweetHandler.TweetCommentWSHandler)
		v1.Use(middleware.JwtAuthMiddleware(logger.GetDefaultLogger()))
		{
			vTweet := v1.Group("/tweet")
			vTweet.GET("", tweetHandler.GetTweetByUserID)
			vTweet.POST("", tweetHandler.PostTweet)
			vTweet.GET("/all", tweetHandler.GetAllTweets)
			vTweet.GET("/tweet-action", tweetHandler.GetActionTweetsByUserID)
			vTweet.POST("/action", tweetHandler.ActionTweet)
			vTweet.DELETE("/delete-action", tweetHandler.DeleteActionTweet)
			vTweet.GET("/comment", tweetHandler.GetTweetComments)

			vSocket := v1.Group("/ws")
			vSocket.GET("/comment-tweet-ws", tweetHandler.TweetCommentWebSocketHandler)
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
