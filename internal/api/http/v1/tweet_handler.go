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
// @Summary     GetTweetByUserID
// @Description Get tweet by user id
// @Tags        Tweet
// @Accept      json
// @Produce     json
// @Param       Authorization header   string true "Bearer <your_token>"
// @Param       userID        query    int    true "User ID"
// @Success     200           {string} string "ok"
// @Router      /tweet [get]
func (h *TweetHandler) GetTweetByUserID(c *gin.Context) {
	var req model.GetTweetByUserReq

	err := request.GetQueryParamsFromUrl(c, &req)
	if err != nil {
		return
	}
	token, err := token.GetTokenString(c)
	if err != nil {
		logger.Error("TweetHandler-GetTweetByUserID: get token from request error", zap.Error(err))
		return
	}
	res, total, err := h.tweerService.GetTweetByUserID(c.Request.Context(), req, token)
	if err != nil {
		c.JSON(http.StatusOK, common.NewErrorResponse(err.Error(), "GetTweetByUserID failure"))
	}
	c.JSON(http.StatusOK, common.NewPagingSuccessResponse(res, total))
}

// PostTweet godoc
// @Summary     PostTweet
// @Description Post new Tweet
// @Tags        Tweet
// @Accept      json
// @Produce     json
// @Param       Authorization header string             true "Bearer <your_token>"
// @Param       body         body   model.PostTweetReq true "Post Tweet Request"
// @Success     200           {string} string "ok"
// @Router      /tweet [post]
func (h *TweetHandler) PostTweet(c *gin.Context) {
	var req model.PostTweetReq

	err := request.GetBodyJSON(c, &req)
	if err != nil {
		return
	}

	userId, err := token.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err.Error(), constants.PostTweetFailure))
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
