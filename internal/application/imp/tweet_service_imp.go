package imp

import (
	"context"

	"github.com/nhutHao02/social-network-common-service/utils/logger"
	"github.com/nhutHao02/social-network-tweet-service/internal/application"
	"github.com/nhutHao02/social-network-tweet-service/internal/domain/interface/tweet"
	"github.com/nhutHao02/social-network-tweet-service/internal/domain/model"
	grpcUser "github.com/nhutHao02/social-network-user-service/pkg/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

type tweetService struct {
	queryRepo   tweet.TweetQueryRepository
	commandRepo tweet.TweetCommandRepository
	userClient  grpcUser.UserServiceClient
}

// GetTweets implements application.TweetService.
func (t *tweetService) GetTweets(ctx context.Context, req model.GetTweetsReq) ([]model.GetTweetsRes, uint64, error) {
	res, total, err := t.queryRepo.GetTweets(ctx, req)
	if err != nil {
		return res, total, err
	}

	// Get User Info
	// Create context with metadata
	md := metadata.Pairs("authorization", "Bearer "+req.Token)
	ctxx := metadata.NewOutgoingContext(ctx, md)

	// pass user info to res
	for index, item := range res {
		userRes, err := t.userClient.GetUserInfo(ctxx, &grpcUser.GetUserRequest{UserID: int64(item.UserID)})
		if err != nil {
			logger.Error("tweetService-GetTweets: call grpcUser to server error", zap.Error(err))
		}
		res[index].UserInfor = &model.UserInfo{
			ID:       int(userRes.Id),
			Email:    userRes.Email,
			FullName: &userRes.FullName,
			UrlAvt:   &userRes.UrlAvt,
		}
	}

	return res, total, nil
}

// PostTweet implements application.TweetService.
func (t *tweetService) PostTweet(ctx context.Context, req model.PostTweetReq) (bool, error) {
	success, err := t.commandRepo.PostTweet(ctx, req)
	if err != nil {
		return false, err
	}
	return success, nil

}

// GetTweetByUserID implements application.TweetService.
func (t *tweetService) GetTweetByUserID(ctx context.Context, req model.GetTweetByUserReq, token string) ([]model.GetTweetByUserRes, uint64, error) {
	res, total, err := t.queryRepo.GetTweetByUserID(ctx, req)
	if err != nil {
		return res, total, err
	}

	// Get User Info
	// Create context with metadata
	md := metadata.Pairs("authorization", "Bearer "+token)
	ctxx := metadata.NewOutgoingContext(ctx, md)

	userRes, err := t.userClient.GetUserInfo(ctxx, &grpcUser.GetUserRequest{UserID: int64(req.UserID)})
	if err != nil {
		logger.Error("tweetService-GetTweetByUserID: call grpcUser to server error", zap.Error(err))
		return res, total, err
	}

	// pass user info to res
	for index, _ := range res {
		res[index].UserInfor = &model.UserInfo{
			ID:       int(userRes.Id),
			Email:    userRes.Email,
			FullName: &userRes.FullName,
			UrlAvt:   &userRes.UrlAvt,
		}
	}

	return res, total, nil
}

func NewTweetService(
	queryRepo tweet.TweetQueryRepository,
	commandRepo tweet.TweetCommandRepository,
	userClient grpcUser.UserServiceClient,
) application.TweetService {
	return &tweetService{queryRepo: queryRepo, commandRepo: commandRepo, userClient: userClient}
}
