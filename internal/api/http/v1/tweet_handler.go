package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nhutHao02/social-network-common-service/utils/token"
	"github.com/nhutHao02/social-network-tweet-service/internal/application"
	grpcUser "github.com/nhutHao02/social-network-user-service/pkg/grpc"
	"google.golang.org/grpc/metadata"
)

type TweetHandler struct {
	tweerService application.TweetService
	userClient   grpcUser.UserServiceClient
}

func NewTweetHandler(tweerService application.TweetService, userClient grpcUser.UserServiceClient) *TweetHandler {
	return &TweetHandler{tweerService: tweerService, userClient: userClient}
}

func (h *TweetHandler) GetTweetByUserID(c *gin.Context) {
	token, err := token.GetTokenString(c)

	// Create context with metadata
	md := metadata.Pairs("authorization", "Bearer "+token)
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	res, err := h.userClient.GetUserInfo(ctx, &grpcUser.GetUserRequest{UserID: 27})
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"id":   res.Id,
		"vlue": res.Email,
	})
}
