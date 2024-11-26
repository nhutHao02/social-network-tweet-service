package model

import "time"

type TweetCommentReq struct {
	TweetID int64 `form:"tweetID"`
	Page    int64 `form:"page"`
	Limit   int64 `form:"limit"`
	Token   string
}
type TweetCommentRes struct {
	ID          int64     `json:"id" db:"ID"`
	UserID      int64     `json:"userId" db:"UserID"`
	UserInfo    *UserInfo `json:"userInfo"`
	Description string    `json:"description" db:"Description"`
	CreatedAt   time.Time `json:"createdAt" db:"CreatedAt"`
}
