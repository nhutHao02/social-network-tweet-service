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
