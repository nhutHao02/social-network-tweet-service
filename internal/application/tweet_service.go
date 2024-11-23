package application

import (
	"context"

	"github.com/nhutHao02/social-network-tweet-service/internal/domain/model"
)

type TweetService interface {
	GetTweetByUserID(ctx context.Context, req model.GetTweetByUserReq, token string) ([]model.GetTweetByUserRes, uint64, error)
	PostTweet(ctx context.Context, req model.PostTweetReq) (bool, error)
	GetTweets(ctx context.Context, req model.GetTweetsReq) ([]model.GetTweetsRes, uint64, error)
	GetLoveTweetsByUserID(ctx context.Context, req model.GetLoveTweetsByUserIDReq) ([]model.GetTweetsRes, uint64, error)
}
