package model

import "time"

type CommentWSReq struct {
	TweetID int64  `form:"tweetID"`
	UserID  int64  `form:"userID"`
	Token   string `form:"token"`
}

type IncomingMessageWSReq struct {
	Content string `json:"content"`
}

type OutgoingMessageWSRes struct {
	ID        int64     `json:"id" db:"ID"`
	UserID    int64     `json:"userId" db:"UserID"`
	From      *UserInfo `json:"userInfo"`
	Content   string    `json:"description" db:"Content"`
	Timestamp time.Time `json:"createdAt" db:"Timestamp"`
}
