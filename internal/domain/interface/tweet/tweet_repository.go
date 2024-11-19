package tweet

import (
	"context"

	"github.com/nhutHao02/social-network-tweet-service/internal/domain/model"
)

type TweetQueryRepository interface {
	GetTweetByUserID(ctx context.Context, req model.GetTweetByUserReq) ([]model.GetTweetByUserRes, uint64, error)
}

type TweetCommandRepository interface {
	PostTweet(ctx context.Context, req model.PostTweetReq) (bool, error)
}
