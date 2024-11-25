// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package internal

import (
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	"github.com/nhutHao02/social-network-tweet-service/config"
	"github.com/nhutHao02/social-network-tweet-service/internal/api"
	"github.com/nhutHao02/social-network-tweet-service/internal/api/http"
	"github.com/nhutHao02/social-network-tweet-service/internal/api/http/v1"
	"github.com/nhutHao02/social-network-tweet-service/internal/application/imp"
	"github.com/nhutHao02/social-network-tweet-service/internal/infrastructure/tweet"
	"github.com/nhutHao02/social-network-tweet-service/pkg/redis"
	"github.com/nhutHao02/social-network-tweet-service/pkg/websocket"
	"github.com/nhutHao02/social-network-user-service/pkg/grpc"
)

// Injectors from wire.go:

func InitializeServer(cfg *config.Config, db *sqlx.DB, rdb *redis.RedisClient, userClient grpc.UserServiceClient, commentWS *websocket.Socket) *api.Server {
	tweetQueryRepository := tweet.NewTweetQueryRepository(db)
	tweetCommandRepository := tweet.NewTweetCommandRepository(db, tweetQueryRepository)
	tweetService := imp.NewTweetService(tweetQueryRepository, tweetCommandRepository, userClient, commentWS)
	tweetHandler := v1.NewTweetHandler(tweetService, userClient)
	httpServer := http.NewHTTPServer(cfg, tweetHandler)
	server := api.NewSerVer(httpServer)
	return server
}

// wire.go:

var serverSet = wire.NewSet(api.NewSerVer)

var itemServerSet = wire.NewSet(http.NewHTTPServer)

var httpHandlerSet = wire.NewSet(v1.NewTweetHandler)

var serviceSet = wire.NewSet(imp.NewTweetService)

var repositorySet = wire.NewSet(tweet.NewTweetCommandRepository, tweet.NewTweetQueryRepository)
