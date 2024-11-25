package constants

var (
	InvalidUserID        = "INVALID USER ID"
	PostTweetFailure     = "POST TWEET FAILURE"
	GetTweetsFailure     = "GET TWEETS FAILURE"
	GetLoveTweetsFailure = "GET LOVE TWEETS FAILURE"
	InvalidAction        = "INVALID ACTION"
	ActionTweetsFailure  = "ACTIONS TO TWEETS FAILURE"
	ConnectWSFailure     = "CONNECT WEBSOCKET IN COMMENT TWEET FAILURE"
)

var (
	BearerString = "Bearer "
)

type ActionTweet string

const (
	Love     ActionTweet = "Love"
	Bookmark ActionTweet = "Bookmark"
	Repost   ActionTweet = "Repost"
)

func (a ActionTweet) IsValid() bool {
	switch a {
	case Love, Bookmark, Repost:
		return true
	default:
		return false
	}
}
