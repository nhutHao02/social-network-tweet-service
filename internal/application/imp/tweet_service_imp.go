package imp

import (
	"context"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/nhutHao02/social-network-common-service/utils/logger"
	"github.com/nhutHao02/social-network-tweet-service/internal/application"
	"github.com/nhutHao02/social-network-tweet-service/internal/domain/interface/tweet"
	"github.com/nhutHao02/social-network-tweet-service/internal/domain/model"
	"github.com/nhutHao02/social-network-tweet-service/pkg/constants"
	ws "github.com/nhutHao02/social-network-tweet-service/pkg/websocket"
	grpcUser "github.com/nhutHao02/social-network-user-service/pkg/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

type tweetService struct {
	queryRepo   tweet.TweetQueryRepository
	commandRepo tweet.TweetCommandRepository
	userClient  grpcUser.UserServiceClient
	commentWS   *ws.Socket
}

// CommentWebSocket implements application.TweetService.
func (t *tweetService) CommentWebSocket(ctx context.Context, conn *websocket.Conn, req model.CommentWSReq) {
	var (
		roomWSID        = strconv.Itoa(int(req.TweetID)) + "tweet"
		userWSID        = strconv.Itoa(int(req.UserID))
		incomingMessage = model.IncomingMessageWSReq{}
	)
	// Add connection
	t.commentWS.AddConnection(roomWSID, userWSID, conn)

	// Listen for connection close event
	defer t.commentWS.RemoveConnection(roomWSID, userWSID, conn)

	// Get User Info
	// Create context with metadata
	md := metadata.Pairs("authorization", constants.BearerString+req.Token)
	ctxx := metadata.NewOutgoingContext(ctx, md)

	userRes, err := t.userClient.GetUserInfo(ctxx, &grpcUser.GetUserRequest{UserID: req.UserID})
	if err != nil {
		logger.Error("tweetService-CommentWebSocket: Error get UserInfo, call grpcUser to server error", zap.Error(err))
	}

	// Handle incoming messages
	for {
		if err := conn.ReadJSON(&incomingMessage); err != nil {
			logger.Error("TweetHandler-TweetCommentWebSocketHandler: Error reading message", zap.Error(err))
			t.commentWS.RemoveConnection(roomWSID, userWSID, conn)
			break
		}

		// Save comment to db
		params := map[string]interface{}{
			"UserID":      req.UserID,
			"TweetID":     req.TweetID,
			"Description": incomingMessage.Content,
		}
		outGoingMessage, err := t.commandRepo.PostComment(ctx, params)
		if err != nil {
			logger.Warn("tweetService-CommentWebSocket: Error saving comment to DB and ignore broadcast", zap.Error(err))
			continue
		}

		// pass user info to outgoingmessage
		outGoingMessage.From = &model.UserInfo{
			ID:       int(userRes.Id),
			Email:    userRes.Email,
			FullName: &userRes.FullName,
			UrlAvt:   &userRes.UrlAvt,
		}

		// Broadcast message to the room
		t.commentWS.Broadcast(roomWSID, userWSID, outGoingMessage)
	}
}

// DeleteActionTweetsByUserID implements application.TweetService.
func (t *tweetService) DeleteActionTweetsByUserID(ctx context.Context, req model.ActionTweetReq) (bool, error) {
	success, err := t.commandRepo.DeleteActionTweetsByUserID(ctx, req)
	if err != nil {
		return success, err
	}
	return success, nil
}

// ActionTweetsByUserID implements application.TweetService.
func (t *tweetService) ActionTweetsByUserID(ctx context.Context, req model.ActionTweetReq) (bool, error) {
	success, err := t.commandRepo.ActionTweetsByUserID(ctx, req)
	if err != nil {
		return success, err
	}
	return success, nil
}

func getUserInfoFromUserClient(ctx context.Context, userClient grpcUser.UserServiceClient, token string, res *[]model.GetTweetsRes) {
	// Get User Info
	// Create context with metadata
	md := metadata.Pairs("authorization", constants.BearerString+token)
	ctxx := metadata.NewOutgoingContext(ctx, md)

	// pass user info to res
	for index, item := range *res {
		userRes, err := userClient.GetUserInfo(ctxx, &grpcUser.GetUserRequest{UserID: int64(item.UserID)})
		if err != nil {
			logger.Error("token: call grpcUser to server error", zap.Error(err))
		}
		(*res)[index].UserInfor = &model.UserInfo{
			ID:       int(userRes.Id),
			Email:    userRes.Email,
			FullName: &userRes.FullName,
			UrlAvt:   &userRes.UrlAvt,
		}
	}
}

// GetActionTweetsByUserID implements application.TweetService.
func (t *tweetService) GetActionTweetsByUserID(ctx context.Context, req model.GetActionTweetsByUserIDReq) ([]model.GetTweetsRes, uint64, error) {
	res, total, err := t.queryRepo.GetActionTweetsByUserID(ctx, req)
	if err != nil {
		return res, total, err
	}

	getUserInfoFromUserClient(ctx, t.userClient, req.Token, &res)

	return res, total, nil

}

// GetTweets implements application.TweetService.
func (t *tweetService) GetTweets(ctx context.Context, req model.GetTweetsReq) ([]model.GetTweetsRes, uint64, error) {
	res, total, err := t.queryRepo.GetTweets(ctx, req)
	if err != nil {
		return res, total, err
	}

	getUserInfoFromUserClient(ctx, t.userClient, req.Token, &res)

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
	md := metadata.Pairs("authorization", constants.BearerString+token)
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
	commentWS *ws.Socket,
) application.TweetService {
	return &tweetService{queryRepo: queryRepo, commandRepo: commandRepo, userClient: userClient, commentWS: commentWS}
}
