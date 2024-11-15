//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	"github.com/nhutHao02/social-network-tweet-service/config"
	"github.com/nhutHao02/social-network-tweet-service/internal/api"
	"github.com/nhutHao02/social-network-tweet-service/internal/api/http"
	v1 "github.com/nhutHao02/social-network-tweet-service/internal/api/http/v1"
	"github.com/nhutHao02/social-network-tweet-service/internal/application/imp"
	"github.com/nhutHao02/social-network-tweet-service/internal/infrastructure/tweet"
	"github.com/nhutHao02/social-network-tweet-service/pkg/redis"
	grpcUser "github.com/nhutHao02/social-network-user-service/pkg/grpc"
)

var serverSet = wire.NewSet(
	api.NewSerVer,
)

var itemServerSet = wire.NewSet(
	http.NewHTTPServer,
)

var httpHandlerSet = wire.NewSet(
	v1.NewTweetHandler,
)

var serviceSet = wire.NewSet(
	imp.NewTweetService,
)

var repositorySet = wire.NewSet(
	tweet.NewTweetCommandRepository,
	tweet.NewTweetQueryRepository,
)

func InitializeServer(
	cfg *config.Config,
	db *sqlx.DB,
	rdb *redis.RedisClient,
	userClient grpcUser.UserServiceClient,
) *api.Server {
	wire.Build(serverSet, itemServerSet, httpHandlerSet, serviceSet, repositorySet)
	return &api.Server{}
}
