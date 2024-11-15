package imp

import (
	"github.com/nhutHao02/social-network-tweet-service/internal/application"
	"github.com/nhutHao02/social-network-tweet-service/internal/domain/interface/tweet"
)

type tweetService struct {
	queryRepo   tweet.TweetQueryRepository
	commandRepo tweet.TweetCommandRepository
}

func NewTweetService(queryRepo tweet.TweetQueryRepository, commandRepo tweet.TweetCommandRepository) application.TweetService {
	return &tweetService{queryRepo: queryRepo, commandRepo: commandRepo}
}
