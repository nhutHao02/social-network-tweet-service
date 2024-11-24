package tweet

import (
	"context"

	"github.com/nhutHao02/social-network-tweet-service/internal/domain/model"
)

type TweetQueryRepository interface {
	GetTweetByUserID(ctx context.Context, req model.GetTweetByUserReq) ([]model.GetTweetByUserRes, uint64, error)
	GetTweets(ctx context.Context, req model.GetTweetsReq) ([]model.GetTweetsRes, uint64, error)
	GetActionTweetsByUserID(ctx context.Context, req model.GetActionTweetsByUserIDReq) ([]model.GetTweetsRes, uint64, error)
	ExistedTweet(ctx context.Context, tweetId int64) (bool, error)
}

type TweetCommandRepository interface {
	PostTweet(ctx context.Context, req model.PostTweetReq) (bool, error)
	ActionTweetsByUserID(ctx context.Context, req model.ActionTweetReq) (bool, error)
	DeleteActionTweetsByUserID(ctx context.Context, req model.ActionTweetReq) (bool, error)
}
