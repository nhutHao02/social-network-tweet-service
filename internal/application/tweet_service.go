package application

import (
	"context"

	"github.com/nhutHao02/social-network-tweet-service/internal/domain/model"
)

type TweetService interface {
	GetTweetByUserID(ctx context.Context, req model.GetTweetByUserReq, token string) ([]model.GetTweetByUserRes, uint64, error)
}
