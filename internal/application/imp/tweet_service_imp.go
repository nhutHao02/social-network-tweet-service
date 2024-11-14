package imp

import "github.com/nhutHao02/social-network-tweet-service/internal/application"

type tweetService struct {
}

func NewTweetService() application.TweetService {
	return &tweetService{}
}
