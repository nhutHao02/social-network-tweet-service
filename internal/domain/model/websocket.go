package model

import "time"

type CommentWSReq struct {
	TweetID int64 `form:"tweetID"`
	UserID  int64 `form:"userID"`
	Token   string
}

type IncomingMessageWSReq struct {
	Content string `json:"content"`
}

type OutgoingMessageWSRes struct {
	ID        int64     `json:"id" db:"ID"`
	From      *UserInfo `json:"from"`
	Content   string    `json:"content" db:"Content"`
	Timestamp time.Time `json:"timestamp" db:"Timestamp"`
}
