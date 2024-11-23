package model

import (
	"time"

	"github.com/nhutHao02/social-network-tweet-service/pkg/constants"
)

type GetTweetByUserReq struct {
	UserID int `form:"userID" db:"UserID"`
	Page   int `form:"page" db:"Page"`
	Limit  int `form:"limit" db:"Limit"`
}

type GetTweetByUserRes struct {
	ID         int         `json:"id" db:"ID"`
	Content    *string     `json:"content" db:"Content"`
	UserID     int         `json:"userID" db:"UserID"`
	UserInfor  *UserInfo   `json:"userInfo"`
	CreatedAt  time.Time   `json:"createdAt" db:"CreatedAt"`
	UserAction *UserAction `json:"action"`
	Statistics *Statistics `json:"statistics"`
}

type Statistics struct {
	TotalLove     int `json:"totalLove" db:"TotalLove"`
	TotalComment  int `json:"totalComment" db:"TotalComment"`
	TotalBookmark int `json:"totalBookmark" db:"TotalBookmark"`
	TotalRepost   int `json:"totalRepost" db:"TotalRepost"`
}

type UserAction struct {
	Love     bool `json:"love" db:"Love"`
	Bookmark bool `json:"bookmark" db:"Bookmark"`
	Repost   bool `json:"repost" db:"Repost"`
}

type PostTweetReq struct {
	UserID   uint64  `json:"userId" db:"UserID"`
	Content  *string `json:"content" db:"Content"`
	UrlImg   *string `json:"urlImg" db:"UrlImg"`
	UrlVideo *string `json:"urlVideo" db:"UrlVideo"`
}

type GetTweetsReq struct {
	UserID int64 `db:"UserID"`
	Token  string
	Page   int `db:"Page"`
	Limit  int `db:"Limit"`
}

type GetTweetsRes struct {
	ID         int         `json:"id" db:"ID"`
	Content    *string     `json:"content" db:"Content"`
	UserID     int         `json:"userID" db:"UserID"`
	UserInfor  *UserInfo   `json:"userInfo"`
	CreatedAt  time.Time   `json:"createdAt" db:"CreatedAt"`
	UserAction *UserAction `json:"action"`
	Statistics *Statistics `json:"statistics"`
}

type GetActionTweetsByUserIDReq struct {
	UserID int                   `form:"userID" db:"UserID"`
	Action constants.ActionTweet `form:"action"`
	Page   int                   `form:"page" db:"Page"`
	Limit  int                   `form:"limit" db:"Limit"`
	Token  string
}
