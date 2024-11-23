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

// GetLoveTweetsByUserID godoc
//
//	@Summary		GetLoveTweetsByUserID
//	@Description	Get Love Tweets By UserID
//	@Tags			Tweet
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string													true	"Bearer <your_token>"
//	@Param			userID			query		int														true	"User ID"
//	@Param			limit			query		int														true	"Limit"
//	@Param			page			query		int														true	"Page"
//	@Success		200				{object}	common.PagingSuccessResponse{data=[]model.GetTweetsRes}	"successful"
//	@Failure		default			{object}	common.Response{data=nil}								"failure"
//	@Router			/tweet/love [get]
func (h *TweetHandler) GetLoveTweetsByUserID(c *gin.Context) {
	var req model.GetLoveTweetsByUserIDReq

	err := request.GetQueryParamsFromUrl(c, &req)
	if err != nil {
		return
	}
	token, err := token.GetTokenString(c)
	if err != nil {
		logger.Error("TweetHandler-GetLoveTweetsByUserID: get token from request error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(err.Error(), constants.GetLoveTweetsFailure))
		return
	}

	req.Token = token

	res, total, err := h.tweerService.GetLoveTweetsByUserID(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), constants.GetLoveTweetsFailure))
		return
	}
	c.JSON(http.StatusOK, common.NewPagingSuccessResponse(res, total))
}
