package model

import (
	"time"

	"github.com/nhutHao02/social-network-tweet-service/pkg/constants"
)

type Notification struct {
	UserID    int64                 `json:"userId"`
	AuthorID  int64                 `json:"authorId"`
	Type      constants.ActionTweet `json:"type"`
	CreatedAt time.Time             `json:"created_at"`
}
