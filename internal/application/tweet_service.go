package application

import (
	"context"

	"github.com/gorilla/websocket"
	"github.com/nhutHao02/social-network-tweet-service/internal/domain/model"
)

type TweetService interface {
	GetTweetByUserID(ctx context.Context, req model.GetTweetByUserReq, token string) ([]model.GetTweetByUserRes, uint64, error)
	PostTweet(ctx context.Context, req model.PostTweetReq) (bool, error)
	GetTweets(ctx context.Context, req model.GetTweetsReq) ([]model.GetTweetsRes, uint64, error)
	GetActionTweetsByUserID(ctx context.Context, req model.GetActionTweetsByUserIDReq) ([]model.GetTweetsRes, uint64, error)
	ActionTweetsByUserID(ctx context.Context, req model.ActionTweetReq) (bool, error)
	DeleteActionTweetsByUserID(ctx context.Context, req model.ActionTweetReq) (bool, error)
	CommentWebSocket(ctx context.Context, conn *websocket.Conn, req model.CommentWSReq)
	GetTweetComments(ctx context.Context, req model.TweetCommentReq) ([]model.TweetCommentRes, uint64, error)
}
