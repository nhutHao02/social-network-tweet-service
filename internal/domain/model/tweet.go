package model

import "time"

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
