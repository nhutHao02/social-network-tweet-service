package v1

import "github.com/nhutHao02/social-network-tweet-service/internal/application"

type TweetHandler struct {
	tweerService application.TweetService
}

func NewTweetHandler(tweerService application.TweetService) *TweetHandler {
	return &TweetHandler{tweerService: tweerService}
}
