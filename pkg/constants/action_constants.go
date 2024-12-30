package constants

type ActionTweet string

const (
	Love     ActionTweet = "Love"
	Bookmark ActionTweet = "Bookmark"
	Repost   ActionTweet = "Repost"
	Comment  ActionTweet = "Comment"
	Post     ActionTweet = "Post"
)

func (a ActionTweet) IsValid() bool {
	switch a {
	case Love, Bookmark, Repost, Post:
		return true
	default:
		return false
	}
}
