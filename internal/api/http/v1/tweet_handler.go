package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	common "github.com/nhutHao02/social-network-common-service/model"
	"github.com/nhutHao02/social-network-common-service/request"
	"github.com/nhutHao02/social-network-common-service/utils/logger"
	"github.com/nhutHao02/social-network-common-service/utils/token"
	"github.com/nhutHao02/social-network-tweet-service/internal/application"
	"github.com/nhutHao02/social-network-tweet-service/internal/domain/model"
	"github.com/nhutHao02/social-network-tweet-service/pkg/constants"
	"github.com/nhutHao02/social-network-tweet-service/pkg/websocket"
	grpcUser "github.com/nhutHao02/social-network-user-service/pkg/grpc"
	"go.uber.org/zap"
)

type TweetHandler struct {
	tweerService application.TweetService
	userClient   grpcUser.UserServiceClient
}

func NewTweetHandler(tweerService application.TweetService, userClient grpcUser.UserServiceClient) *TweetHandler {
	return &TweetHandler{tweerService: tweerService, userClient: userClient}
}

// GetTweetByUserID godoc
//
//	@Summary		GetTweetByUserID
//	@Description	Get tweet by user id
//	@Tags			Tweet
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string															true	"Bearer <your_token>"
//	@Param			userID			query		int																true	"User ID"
//	@Param			limit			query		int																true	"Limit"
//	@Param			page			query		int																true	"Page"
//	@Success		200				{object}	common.PagingSuccessResponse{data=[]model.GetTweetByUserRes}	"successful"
//	@Failure		default			{object}	common.Response{data=nil}										"failure"
//	@Router			/tweet [get]
func (h *TweetHandler) GetTweetByUserID(c *gin.Context) {
	var req model.GetTweetByUserReq

	err := request.GetQueryParamsFromUrl(c, &req)
	if err != nil {
		return
	}
	token, err := token.GetTokenString(c)
	if err != nil {
		logger.Error("TweetHandler-GetTweetByUserID: get token from request error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(err.Error(), "GetTweetByUserID failure"))
		return
	}
	res, total, err := h.tweerService.GetTweetByUserID(c.Request.Context(), req, token)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), "GetTweetByUserID failure"))
		return
	}
	c.JSON(http.StatusOK, common.NewPagingSuccessResponse(res, total))
}

// PostTweet godoc
//
//	@Summary		PostTweet
//	@Description	Post new Tweet
//	@Tags			Tweet
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							true	"Bearer <your_token>"
//	@Param			body			body		model.PostTweetReq				true	"Post Tweet Request"
//	@Success		200				{object}	common.Response{data=boolean}	"successfully"
//	@Failure		default			{object}	common.Response{data=nil}		"failure"
//	@Router			/tweet [post]
func (h *TweetHandler) PostTweet(c *gin.Context) {
	var req model.PostTweetReq

	err := request.GetBodyJSON(c, &req)
	if err != nil {
		return
	}

	userId, err := token.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(err.Error(), constants.PostTweetFailure))
		return
	}
	if int(req.UserID) != userId {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(constants.InvalidUserID, constants.PostTweetFailure))
		return
	}

	success, err := h.tweerService.PostTweet(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), constants.PostTweetFailure))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(success))
}

// GetAllTweets godoc
//
//	@Summary		GetAllTweets
//	@Description	Get All Tweets
//	@Tags			Tweet
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string													true	"Bearer <your_token>"
//	@Param			limit			query		int														true	"Limit"
//	@Param			page			query		int														true	"Page"
//	@Success		200				{object}	common.PagingSuccessResponse{data=[]model.GetTweetsRes}	"successfully"
//	@Failure		default			{object}	common.Response{data=nil}								"failure"
//	@Router			/all [get]
func (h *TweetHandler) GetAllTweets(c *gin.Context) {
	paging := request.GetPaging(c)

	tokenString, err := token.GetTokenString(c)
	if err != nil {
		logger.Error("TweetHandler-GetAllTweets: get token from request error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(err.Error(), constants.GetTweetsFailure))
		return
	}

	userID, err := token.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(err.Error(), constants.PostTweetFailure))
		return
	}

	req := model.GetTweetsReq{
		UserID: int64(userID),
		Token:  tokenString,
		Page:   paging.Page,
		Limit:  paging.Limit,
	}
	res, total, err := h.tweerService.GetTweets(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), constants.GetTweetsFailure))
		return
	}

	c.JSON(http.StatusOK, common.NewPagingSuccessResponse(res, total))
}

// GetActionTweetsByUserID godoc
//
//	@Summary		GetActionTweetsByUserID
//	@Description	Get Tweets By UserID and Action
//	@Tags			Tweet
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string													true	"Bearer <your_token>"
//	@Param			userID			query		int														true	"User ID"
//	@Param			action			query		string													true	"Action"
//	@Param			limit			query		int														true	"Limit"
//	@Param			page			query		int														true	"Page"
//	@Success		200				{object}	common.PagingSuccessResponse{data=[]model.GetTweetsRes}	"successful"
//	@Failure		default			{object}	common.Response{data=nil}								"failure"
//	@Router			/tweet/love [get]
func (h *TweetHandler) GetActionTweetsByUserID(c *gin.Context) {
	var req model.GetActionTweetsByUserIDReq

	err := request.GetQueryParamsFromUrl(c, &req)
	if err != nil {
		return
	}

	if !req.Action.IsValid() {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(constants.InvalidAction, constants.GetLoveTweetsFailure))
		return
	}

	token, err := token.GetTokenString(c)
	if err != nil {
		logger.Error("TweetHandler-GetActionTweetsByUserID: get token from request error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(err.Error(), constants.GetLoveTweetsFailure))
		return
	}

	req.Token = token

	res, total, err := h.tweerService.GetActionTweetsByUserID(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), constants.GetLoveTweetsFailure))
		return
	}
	c.JSON(http.StatusOK, common.NewPagingSuccessResponse(res, total))
}

// ActionTweet godoc
//
//	@Summary		ActionTweet
//	@Description	Action with Tweet such as Love, Bookmark, Repost
//	@Tags			Tweet
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							true	"Bearer <your_token>"
//	@Param			model			body		int								true	"ActionTweetReq"
//	@Success		200				{object}	common.Response{data=boolean}	"successful"
//	@Failure		default			{object}	common.Response{data=nil}		"failure"
//	@Router			/tweet/action [post]
func (h *TweetHandler) ActionTweet(c *gin.Context) {
	var req model.ActionTweetReq

	err := request.GetBodyJSON(c, &req)
	if err != nil {
		return
	}

	if !req.Action.IsValid() {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(constants.InvalidAction, constants.ActionTweetsFailure))
		return
	}

	userID, err := token.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), constants.ActionTweetsFailure))
		return
	}
	if userID != int(req.UserID) {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(constants.InvalidUserID, constants.ActionTweetsFailure))
		return
	}

	success, err := h.tweerService.ActionTweetsByUserID(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), constants.ActionTweetsFailure))
		return
	}
	c.JSON(http.StatusOK, common.NewSuccessResponse(success))
}

// DeleteActionTweet godoc
//
//	@Summary		DeleteActionTweet
//	@Description	Delete Action to Tweet such as Love, Bookmark, Repost
//	@Tags			Tweet
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							true	"Bearer <your_token>"
//	@Param			model			body		int								true	"ActionTweetReq"
//	@Success		200				{object}	common.Response{data=boolean}	"successful"
//	@Failure		default			{object}	common.Response{data=nil}		"failure"
//	@Router			/tweet/delete-action [delete]
func (h *TweetHandler) DeleteActionTweet(c *gin.Context) {
	var req model.ActionTweetReq

	err := request.GetBodyJSON(c, &req)
	if err != nil {
		return
	}

	if !req.Action.IsValid() {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(constants.InvalidAction, constants.ActionTweetsFailure))
		return
	}

	userID, err := token.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), constants.ActionTweetsFailure))
		return
	}

	if userID != int(req.UserID) {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(constants.InvalidUserID, constants.ActionTweetsFailure))
		return
	}

	success, err := h.tweerService.DeleteActionTweetsByUserID(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), constants.ActionTweetsFailure))
		return
	}
	c.JSON(http.StatusOK, common.NewSuccessResponse(success))
}

// TweetCommentWebSocketHandler godoc
//
//	@Summary		TweetCommentWebSocketHandler
//	@Description	Establish a WebSocket connection to comment on tweets in real-time.
//	@Tags			Tweet
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string						true	"Bearer <your_token>"
//	@Param			userID			query		int							true	"User ID"
//	@Param			tweetID			query		int							true	"Tweet ID"
//	@Success		101				{string}	string						"WebSocket connection established"
//	@Failure		default			{object}	common.Response{data=nil}	"failure"
//	@Router			/ws/comment-tweet-ws [get]
func (h *TweetHandler) TweetCommentWebSocketHandler(c *gin.Context) {
	var req model.CommentWSReq

	if err := request.GetQueryParamsFromUrl(c, &req); err != nil {
		return
	}
	userID, err := token.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), constants.ConnectWSFailure))
		return
	}

	if userID != int(req.UserID) {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(constants.InvalidUserID, constants.ConnectWSFailure))
		return
	}

	token, err := token.GetTokenString(c)
	if err != nil {
		logger.Error("TweetHandler-GetActionTweetsByUserID: get token from request error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(err.Error(), constants.GetLoveTweetsFailure))
		return
	}

	req.Token = token

	// Upgrade HTTP connection to WebSocket
	conn, err := websocket.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error("Error when upgrade HTTP connection to WebSocket", zap.Error(err))
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(err.Error(), constants.ConnectWSFailure))
		return
	}

	h.tweerService.CommentWebSocket(c.Request.Context(), conn, req)
}

// GetTweetComments godoc
//
//	@Summary		GetTweetComments
//	@Description	Get Tweet Comments
//	@Tags			Tweet
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string														true	"Bearer <your_token>"
//	@Param			tweetID			query		int															true	"Tweet ID"
//	@Param			limit			query		int															true	"Limit"
//	@Param			page			query		int															true	"Page"
//	@Success		200				{object}	common.PagingSuccessResponse{data=[]model.TweetCommentRes}	"successful"
//	@Failure		default			{object}	common.Response{data=nil}									"failure"
//	@Router			/comment [get]
func (h *TweetHandler) GetTweetComments(c *gin.Context) {
	var req model.TweetCommentReq

	if err := request.GetQueryParamsFromUrl(c, &req); err != nil {
		return
	}

	token, err := token.GetTokenString(c)
	if err != nil {
		logger.Error("TweetHandler-GetTweetComments: get token from request error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(err.Error(), constants.GetLoveTweetsFailure))
		return
	}
	req.Token = token

	res, total, err := h.tweerService.GetTweetComments(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), constants.GetLoveTweetsFailure))
		return
	}
	c.JSON(http.StatusOK, common.NewPagingSuccessResponse(res, total))
}

// TweetCommentWSHandler godoc
//
//	@Summary		TweetCommentWSHandler
//	@Description	Establish a WebSocket connection to comment on tweets in real-time.
//	@Tags			Tweet
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string						true	"Bearer <your_token>"
//	@Param			userID			query		int							true	"User ID"
//	@Param			tweetID			query		int							true	"Tweet ID"
//	@Success		101				{string}	string						"WebSocket connection established"
//	@Failure		default			{object}	common.Response{data=nil}	"failure"
//	@Router			/ws/comment-tweet-ws [get]
func (h *TweetHandler) TweetCommentWSHandler(c *gin.Context) {
	var req model.CommentWSReq

	if err := request.GetQueryParamsFromUrl(c, &req); err != nil {
		return
	}
	if req.TweetID == 0 || req.Token == "" || req.UserID == 0 {
		return
	}
	// Upgrade HTTP connection to WebSocket
	conn, err := websocket.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error("Error when upgrade HTTP connection to WebSocket", zap.Error(err))
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(err.Error(), constants.ConnectWSFailure))
		return
	}

	h.tweerService.CommentWebSocket(c.Request.Context(), conn, req)
}
